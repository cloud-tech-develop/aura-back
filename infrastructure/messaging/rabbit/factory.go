package rabbit

import (
	"fmt"
	"log"
	"os"

	"github.com/cloud-tech-develop/aura-back/infrastructure/messaging/memory"
	"github.com/cloud-tech-develop/aura-back/shared/events"
)

// EventBusFactory creates an event bus based on configuration
type EventBusFactory struct{}

// NewEventBusFactory creates a new EventBusFactory
func NewEventBusFactory() *EventBusFactory {
	return &EventBusFactory{}
}

// CreateEventBus creates either RabbitMQ or Memory event bus based on config
func (f *EventBusFactory) CreateEventBus() (events.EventBus, error) {
	useRabbitMQ := os.Getenv("USE_RABBITMQ")
	driver := os.Getenv("DATABASE_DRIVER")

	// Use RabbitMQ in production (postgres) mode if enabled and allowed
	if useRabbitMQ == "true" && driver == "postgres" {
		log.Println("[EventBus] Using RabbitMQ event bus")
		return NewRabbitMQEventBus()
	}

	// Fall back to memory event bus or return nil for automatic memory
	log.Println("[EventBus] Using memory event bus")
	return f.createMemoryBus(), nil
}

// CreateEventBusForTenant creates a RabbitMQ event bus for a specific tenant
// This is used in offline mode to listen to tenant-specific events
func (f *EventBusFactory) CreateEventBusForTenant(tenant string) (events.EventBus, error) {
	useRabbitMQ := os.Getenv("USE_RABBITMQ")

	if useRabbitMQ == "true" {
		log.Printf("[EventBus] Creating RabbitMQ event bus for tenant: %s", tenant)
		return NewRabbitMQEventBusWithTenant(tenant)
	}

	// Fall back to memory
	log.Printf("[EventBus] Creating memory event bus for tenant: %s", tenant)
	return f.createMemoryBus(), nil
}

// createMemoryBus creates a memory event bus
func (f *EventBusFactory) createMemoryBus() events.EventBus {
	bus := memory.NewMemoryEventBus(100, 5)
	if err := bus.Start(); err != nil {
		log.Printf("[EventBus] Failed to start memory bus: %v", err)
		return nil
	}
	return bus
}

// AsyncEventBus wraps an EventBus and publishes to both local (memory) and remote (RabbitMQ)
// Used in online mode to notify offline consumers
type AsyncEventBus struct {
	localBus  events.EventBus
	remoteBus *RabbitMQEventBus
}

// NewAsyncEventBus creates a new async event bus
func NewAsyncEventBus(local events.EventBus, remote *RabbitMQEventBus) *AsyncEventBus {
	return &AsyncEventBus{
		localBus:  local,
		remoteBus: remote,
	}
}

// Publish publishes to both local and remote
func (b *AsyncEventBus) Publish(event events.Event) error {
	// Always publish locally for local consumers
	if b.localBus != nil {
		if err := b.localBus.Publish(event); err != nil {
			log.Printf("[AsyncEventBus] Local publish error: %v", err)
		}
	}

	// Publish to RabbitMQ for remote/offline consumers
	if b.remoteBus != nil {
		if err := b.remoteBus.Publish(event); err != nil {
			log.Printf("[AsyncEventBus] Remote publish error: %v", err)
		}
	}

	return nil
}

// Subscribe subscribes to local bus
func (b *AsyncEventBus) Subscribe(eventName string, handler events.EventHandler) error {
	if b.localBus != nil {
		return b.localBus.Subscribe(eventName, handler)
	}
	return fmt.Errorf("local bus not available")
}

// Unsubscribe unsubscribes from local bus
func (b *AsyncEventBus) Unsubscribe(eventName string, handler events.EventHandler) error {
	if b.localBus != nil {
		return b.localBus.Unsubscribe(eventName, handler)
	}
	return fmt.Errorf("local bus not available")
}

// Start starts the local bus
func (b *AsyncEventBus) Start() error {
	if b.localBus != nil {
		return b.localBus.Start()
	}
	return nil
}

// Stop stops both buses
func (b *AsyncEventBus) Stop() error {
	if b.localBus != nil {
		b.localBus.Stop()
	}
	if b.remoteBus != nil {
		b.remoteBus.Stop()
	}
	return nil
}