package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/cloud-tech-develop/aura-back/shared/events"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQEventBus struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	subscriptions map[string][]events.EventHandler
	exchange      string
	env           string
	tenant        string
	mu            sync.RWMutex
	stopChan      chan struct{}
	wg            sync.WaitGroup
	started       bool
}

// validEnvs are the only accepted APP_ENV values for exchange naming.
var validEnvs = map[string]bool{"dev": true, "prod": true}

// resolveEnv returns the APP_ENV to use for exchange naming.
// Falls back to "dev" and logs a warning when the variable is absent or unrecognised,
// preventing queues like "aura..something" from ever being created.
func resolveEnv() (string, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		log.Printf("[RabbitMQ] WARNING: APP_ENV not set, defaulting to 'dev'. " +
			"Set APP_ENV=dev or APP_ENV=prod to silence this warning.")
		return "dev", nil
	}
	if !validEnvs[env] {
		return "", fmt.Errorf("[RabbitMQ] invalid APP_ENV=%q: must be 'dev' or 'prod'", env)
	}
	return env, nil
}

// validateName checks that an exchange, queue or routing-key segment contains
// no empty parts (consecutive dots), which would indicate a missing env or tenant.
func validateName(label, name string) error {
	if strings.Contains(name, "..") {
		return fmt.Errorf("[RabbitMQ] malformed %s name %q: contains empty segment (consecutive dots)", label, name)
	}
	if strings.HasPrefix(name, ".") || strings.HasSuffix(name, ".") {
		return fmt.Errorf("[RabbitMQ] malformed %s name %q: starts or ends with a dot", label, name)
	}
	return nil
}

// NewRabbitMQEventBus creates a new RabbitMQ event bus (online/shared mode, no tenant).
// Exchange: aura.prod or aura.dev
// Routing keys include tenant slug extracted from each event's payload.
func NewRabbitMQEventBus() (*RabbitMQEventBus, error) {
	host := os.Getenv("RABBITMQ_HOST")
	user := os.Getenv("RABBITMQ_DEFAULT_USER")
	pass := os.Getenv("RABBITMQ_DEFAULT_PASS")

	env, err := resolveEnv()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("amqp://%s:%s@%s:5672/", user, pass, host)
	log.Printf("[RabbitMQ] Connecting to '%s' (exchange env: %s)", host, env)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("[RabbitMQ] failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("[RabbitMQ] failed to open channel: %w", err)
	}

	exchangeName := fmt.Sprintf("aura.%s", env)
	err = ch.ExchangeDeclare(exchangeName, "topic", true, false, false, false, nil)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("[RabbitMQ] failed to declare exchange: %w", err)
	}

	log.Printf("[RabbitMQ] Connected to %s, exchange: %s", host, exchangeName)

	return &RabbitMQEventBus{
		conn:          conn,
		channel:       ch,
		subscriptions: make(map[string][]events.EventHandler),
		exchange:      exchangeName,
		env:           env,
		stopChan:      make(chan struct{}),
	}, nil
}

// NewRabbitMQEventBusWithTenant creates a RabbitMQ event bus for a specific tenant (offline mode).
// All publishes and subscriptions are scoped to this tenant's routing key suffix.
func NewRabbitMQEventBusWithTenant(tenant string) (*RabbitMQEventBus, error) {
	bus, err := NewRabbitMQEventBus()
	if err != nil {
		return nil, err
	}
	bus.tenant = tenant
	return bus, nil
}

func getEnv(key, defaultValue string) string {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Publish publishes an event with a tenant-aware routing key.
//
// Routing key resolution order:
//  1. b.tenant (set on offline instances after /offline/ping)
//  2. "tenant_slug" field in the event payload (used by online shared bus)
//  3. No tenant → routing key is just the event name
//
// Final format: {eventName}.{slug}  e.g. "product.offline.created.empresa_uno"
func (b *RabbitMQEventBus) Publish(event events.Event) error {
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	// Normalise payload to map[string]interface{}
	payload, ok := event.GetPayload().(map[string]interface{})
	if !ok {
		data, err := json.Marshal(event.GetPayload())
		if err != nil {
			payload = map[string]interface{}{"data": event.GetPayload()}
		} else {
			json.Unmarshal(data, &payload)
		}
	}

	// Determine the tenant for routing
	slug := b.tenant
	if slug == "" {
		// Online shared bus: extract from payload so each tenant's events are isolated
		if ts, ok := payload["tenant_slug"].(string); ok && ts != "" {
			slug = ts
		}
	}

	routingKey := event.GetName()
	if slug != "" {
		routingKey = fmt.Sprintf("%s.%s", event.GetName(), slug)
	}

	if err := validateName("routing key", routingKey); err != nil {
		return err
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = b.channel.PublishWithContext(ctx,
		b.exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish: %w", err)
	}

	log.Printf("[RabbitMQ] Published event: %s (routing key: %s)", event.GetName(), routingKey)
	return nil
}

// Subscribe registers a handler for an event. Supports wildcard suffix ".*":
//
//   - "product.online.created.*"   → queue aura.{env}.product.online.created
//                                     binding product.online.created.*  (online receives all tenants)
//   - "product.offline.created"    → queue aura.{env}.product.offline.created.{tenant}
//                                     binding product.offline.created.{tenant}  (offline: own tenant only)
//
// If Start() has already been called, a consumer goroutine is launched immediately.
func (b *RabbitMQEventBus) Subscribe(eventName string, handler events.EventHandler) error {
	if eventName == "" {
		return fmt.Errorf("event name cannot be empty")
	}
	if handler == nil {
		return fmt.Errorf("handler cannot be nil")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	// Strip wildcard suffix — used for queue name and internal dispatch key
	baseEventName := strings.TrimSuffix(eventName, ".*")
	isWildcard := baseEventName != eventName

	// If a consumer is already running for this event, just add the handler
	if _, exists := b.subscriptions[baseEventName]; exists {
		b.subscriptions[baseEventName] = append(b.subscriptions[baseEventName], handler)
		log.Printf("[RabbitMQ] Added handler to existing subscription: %s", baseEventName)
		return nil
	}

	// Build queue name
	queueName := fmt.Sprintf("aura.%s.%s", b.env, baseEventName)
	if b.tenant != "" {
		queueName = fmt.Sprintf("aura.%s.%s.%s", b.env, baseEventName, b.tenant)
	}

	// Build routing key for the queue binding
	routingKey := baseEventName
	switch {
	case isWildcard:
		routingKey = baseEventName + ".*"
	case b.tenant != "":
		routingKey = fmt.Sprintf("%s.%s", baseEventName, b.tenant)
	}

	if err := validateName("queue", queueName); err != nil {
		return err
	}
	if err := validateName("routing key", routingKey); err != nil {
		return err
	}

	if _, err := b.channel.QueueDeclare(queueName, true, false, false, false, nil); err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
	}
	if err := b.channel.QueueBind(queueName, routingKey, b.exchange, false, nil); err != nil {
		return fmt.Errorf("failed to bind queue %s (routing: %s): %w", queueName, routingKey, err)
	}

	b.subscriptions[baseEventName] = []events.EventHandler{handler}

	log.Printf("[RabbitMQ] Subscribed: %s  queue=%s  routing=%s", baseEventName, queueName, routingKey)

	// If already started, launch a consumer goroutine immediately
	if b.started {
		msgs, err := b.channel.Consume(queueName, "", false, false, false, false, nil)
		if err != nil {
			return fmt.Errorf("failed to start consumer for %s: %w", queueName, err)
		}
		b.wg.Add(1)
		go b.consume(msgs, baseEventName)
	}

	return nil
}

// Unsubscribe removes a handler from a subscription.
func (b *RabbitMQEventBus) Unsubscribe(eventName string, handler events.EventHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	baseEventName := strings.TrimSuffix(eventName, ".*")
	handlers := b.subscriptions[baseEventName]
	for i, h := range handlers {
		if h == handler {
			b.subscriptions[baseEventName] = append(handlers[:i], handlers[i+1:]...)
			log.Printf("[RabbitMQ] Unsubscribed from event: %s", baseEventName)
			return nil
		}
	}
	return fmt.Errorf("handler not found for event: %s", baseEventName)
}

// Start begins consuming messages for all registered subscriptions.
// Subscriptions added after Start() also get a consumer goroutine immediately.
func (b *RabbitMQEventBus) Start() error {
	log.Printf("[RabbitMQ] Starting consumer...")

	b.mu.Lock()
	b.started = true
	// Snapshot subscriptions while holding the lock
	type queueEntry struct {
		queueName string
		eventName string
	}
	entries := make([]queueEntry, 0, len(b.subscriptions))
	for eventName := range b.subscriptions {
		queueName := fmt.Sprintf("aura.%s.%s", b.env, eventName)
		if b.tenant != "" {
			queueName = fmt.Sprintf("aura.%s.%s.%s", b.env, eventName, b.tenant)
		}
		entries = append(entries, queueEntry{queueName: queueName, eventName: eventName})
	}
	b.mu.Unlock()

	for _, e := range entries {
		msgs, err := b.channel.Consume(e.queueName, "", false, false, false, false, nil)
		if err != nil {
			return fmt.Errorf("failed to consume %s: %w", e.queueName, err)
		}
		b.wg.Add(1)
		go b.consume(msgs, e.eventName)
	}

	return nil
}

// consume dispatches incoming RabbitMQ messages to registered handlers.
// Reads handlers dynamically from b.subscriptions so handlers added after
// startup are also invoked.
func (b *RabbitMQEventBus) consume(msgs <-chan amqp.Delivery, eventName string) {
	defer b.wg.Done()

	log.Printf("[RabbitMQ] Consumer started for: %s", eventName)

	for {
		select {
		case <-b.stopChan:
			log.Printf("[RabbitMQ] Consumer stopped for: %s", eventName)
			return
		case msg, ok := <-msgs:
			if !ok {
				log.Printf("[RabbitMQ] Channel closed for: %s", eventName)
				return
			}

			var payload map[string]interface{}
			if err := json.Unmarshal(msg.Body, &payload); err != nil {
				log.Printf("[RabbitMQ] Failed to unmarshal: %v", err)
				msg.Ack(false)
				continue
			}

			event := events.NewBaseEvent(eventName, payload)

			// Read handlers dynamically to pick up any added after startup
			b.mu.RLock()
			handlers := make([]events.EventHandler, len(b.subscriptions[eventName]))
			copy(handlers, b.subscriptions[eventName])
			b.mu.RUnlock()

			for _, handler := range handlers {
				if err := handler.Handle(&event); err != nil {
					log.Printf("[RabbitMQ] Error handling %s: %v", eventName, err)
				}
			}

			msg.Ack(false)
		}
	}
}

// Stop stops all consumer goroutines and closes the connection.
func (b *RabbitMQEventBus) Stop() error {
	log.Println("[RabbitMQ] Stopping...")

	close(b.stopChan)
	b.wg.Wait()

	if b.channel != nil {
		b.channel.Close()
	}
	if b.conn != nil {
		b.conn.Close()
	}

	log.Println("[RabbitMQ] Stopped")
	return nil
}

func (b *RabbitMQEventBus) GetExchange() string { return b.exchange }
func (b *RabbitMQEventBus) GetEnv() string      { return b.env }
func (b *RabbitMQEventBus) GetTenant() string   { return b.tenant }
