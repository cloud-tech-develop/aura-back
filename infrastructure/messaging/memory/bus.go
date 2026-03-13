package memory

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/cloud-tech-develop/aura-back/domain/events"
)

type subscription struct {
	eventName string
	handler   events.EventHandler
}

type MemoryEventBus struct {
	subscriptions []subscription
	channel       chan events.Event
	stopChan      chan struct{}
	wg            sync.WaitGroup
	mu            sync.RWMutex
	bufferSize    int
	workers       int
}

func NewMemoryEventBus(bufferSize, workers int) *MemoryEventBus {
	if bufferSize <= 0 {
		bufferSize = 100
	}
	if workers <= 0 {
		workers = 5
	}
	return &MemoryEventBus{
		subscriptions: make([]subscription, 0),
		channel:       make(chan events.Event, bufferSize),
		stopChan:      make(chan struct{}),
		bufferSize:    bufferSize,
		workers:       workers,
	}
}

func (b *MemoryEventBus) Publish(event events.Event) error {
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	select {
	case b.channel <- event:
		return nil
	case <-time.After(time.Second):
		return fmt.Errorf("timeout publishing event: %s", event.GetName())
	}
}

func (b *MemoryEventBus) Subscribe(eventName string, handler events.EventHandler) error {
	if eventName == "" {
		return fmt.Errorf("event name cannot be empty")
	}
	if handler == nil {
		return fmt.Errorf("handler cannot be nil")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.subscriptions = append(b.subscriptions, subscription{
		eventName: eventName,
		handler:   handler,
	})

	log.Printf("[EventBus] Subscribed handler to event: %s", eventName)
	return nil
}

func (b *MemoryEventBus) Unsubscribe(eventName string, handler events.EventHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	newSubs := make([]subscription, 0)
	for _, sub := range b.subscriptions {
		if sub.eventName == eventName && sub.handler == handler {
			continue
		}
		newSubs = append(newSubs, sub)
	}
	b.subscriptions = newSubs

	log.Printf("[EventBus] Unsubscribed handler from event: %s", eventName)
	return nil
}

func (b *MemoryEventBus) Start() error {
	log.Printf("[EventBus] Starting with %d workers, buffer size: %d", b.workers, b.bufferSize)

	for i := 0; i < b.workers; i++ {
		b.wg.Add(1)
		go b.worker(i)
	}

	return nil
}

func (b *MemoryEventBus) Stop() error {
	log.Println("[EventBus] Stopping...")
	close(b.stopChan)
	b.wg.Wait()
	log.Println("[EventBus] Stopped")
	return nil
}

func (b *MemoryEventBus) worker(id int) {
	defer b.wg.Done()

	log.Printf("[EventBus] Worker %d started", id)

	for {
		select {
		case event := <-b.channel:
			b.handleEvent(event)
		case <-b.stopChan:
			log.Printf("[EventBus] Worker %d stopped", id)
			return
		}
	}
}

func (b *MemoryEventBus) handleEvent(event events.Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	subs := make([]events.EventHandler, 0)
	for _, sub := range b.subscriptions {
		if sub.eventName == event.GetName() {
			subs = append(subs, sub.handler)
		}
	}

	if len(subs) == 0 {
		log.Printf("[EventBus] No handlers for event: %s", event.GetName())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, handler := range subs {
		select {
		case <-ctx.Done():
			log.Printf("[EventBus] Timeout handling event %s", event.GetName())
			return
		default:
			if err := handler.Handle(event); err != nil {
				log.Printf("[EventBus] Error handling event %s: %v", event.GetName(), err)
			}
		}
	}
}
