package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
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
}

// NewRabbitMQEventBus creates a new RabbitMQ event bus
// Exchange format: aura.prod or aura.dev
// Queues are tenant-specific: aura.prod.product.created.empresa_uno
func NewRabbitMQEventBus() (*RabbitMQEventBus, error) {
	host := os.Getenv("RABBITMQ_HOST")
	user := os.Getenv("RABBITMQ_DEFAULT_USER")
	pass := os.Getenv("RABBITMQ_DEFAULT_PASS")
	env := os.Getenv("APP_ENV") // "prod" or "dev"s

	url := fmt.Sprintf("amqp://%s:%s@%s:5672/", user, pass, host)
	log.Printf("[RabbitMQ] Connecting to '%s'", host)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("[RabbitMQ] failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("[RabbitMQ] failed to open channel: %w", err)
	}

	// Exchange pattern: aura.prod or aura.dev
	exchangeName := fmt.Sprintf("aura.%s", env)
	err = ch.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type (supports tenant-specific routing)
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
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

// NewRabbitMQEventBusWithTenant creates a RabbitMQ event bus for a specific tenant
// Exchange: aura.prod or aura.dev
// Queue pattern: aura.prod.{event}.{tenant}
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

// Publish publishes an event to RabbitMQ with tenant-specific routing key
// Routing key format: {eventname}.{tenant}  (e.g., product.created.empresa_uno)
// If tenant is set, it appends tenant to the routing key
func (b *RabbitMQEventBus) Publish(event events.Event) error {
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	// Build routing key with tenant if available
	routingKey := event.GetName()
	if b.tenant != "" {
		routingKey = fmt.Sprintf("%s.%s", event.GetName(), b.tenant)
	}

	// Get payload as map
	payload, ok := event.GetPayload().(map[string]interface{})
	if !ok {
		// Try to marshal if it's not a map
		data, err := json.Marshal(event.GetPayload())
		if err != nil {
			payload = map[string]interface{}{"data": event.GetPayload()}
		} else {
			json.Unmarshal(data, &payload)
		}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = b.channel.PublishWithContext(ctx,
		b.exchange, // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
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

	log.Printf("[RabbitMQ] Published event: %s to routing key: %s", event.GetName(), routingKey)
	return nil
}

// Subscribe subscribes to an event by name pattern with tenant-specific queue
// Queue: aura.{env}.{eventname}.{tenant}
// This allows offline versions to listen only to their tenant's events
func (b *RabbitMQEventBus) Subscribe(eventName string, handler events.EventHandler) error {
	if eventName == "" {
		return fmt.Errorf("event name cannot be empty")
	}
	if handler == nil {
		return fmt.Errorf("handler cannot be nil")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	// Build queue name with tenant
	queueName := fmt.Sprintf("aura.%s.%s", b.env, eventName)
	if b.tenant != "" {
		queueName = fmt.Sprintf("aura.%s.%s.%s", b.env, eventName, b.tenant)
	}

	_, err := b.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Routing key pattern - include tenant if available
	routingKey := eventName
	if b.tenant != "" {
		routingKey = fmt.Sprintf("%s.%s", eventName, b.tenant)
	}

	// Bind queue to exchange
	err = b.channel.QueueBind(
		queueName,  // queue name
		routingKey, // routing key
		b.exchange, // exchange
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	// Add handler to subscriptions
	b.subscriptions[eventName] = append(b.subscriptions[eventName], handler)

	log.Printf("[RabbitMQ] Subscribed to event: %s on queue: %s with routing key: %s", eventName, queueName, routingKey)
	return nil
}

// Unsubscribe removes a subscription
func (b *RabbitMQEventBus) Unsubscribe(eventName string, handler events.EventHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	handlers := b.subscriptions[eventName]
	for i, h := range handlers {
		if h == handler {
			b.subscriptions[eventName] = append(handlers[:i], handlers[i+1:]...)
			log.Printf("[RabbitMQ] Unsubscribed from event: %s", eventName)
			return nil
		}
	}
	return fmt.Errorf("handler not found for event: %s", eventName)
}

// Start starts consuming messages
func (b *RabbitMQEventBus) Start() error {
	log.Printf("[RabbitMQ] Starting consumer...")

	b.mu.RLock()
	defer b.mu.RUnlock()

	// Start a consumer for each subscription
	for eventName, handlers := range b.subscriptions {
		queueName := fmt.Sprintf("aura.%s.%s", b.env, eventName)
		if b.tenant != "" {
			queueName = fmt.Sprintf("aura.%s.%s.%s", b.env, eventName, b.tenant)
		}

		msgs, err := b.channel.Consume(
			queueName, // queue
			"",        // consumer
			false,     // auto-ack
			false,     // exclusive
			false,     // no-local
			false,     // no-wait
			nil,       // args
		)
		if err != nil {
			return fmt.Errorf("failed to consume: %w", err)
		}

		b.wg.Add(1)
		go b.consume(msgs, eventName, handlers)
	}

	return nil
}

func (b *RabbitMQEventBus) consume(msgs <-chan amqp.Delivery, eventName string, handlers []events.EventHandler) {
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

			// Unmarshal payload
			var payload map[string]interface{}
			if err := json.Unmarshal(msg.Body, &payload); err != nil {
				log.Printf("[RabbitMQ] Failed to unmarshal: %v", err)
				msg.Ack(false)
				continue
			}

			// Create event from message
			event := events.NewBaseEvent(eventName, payload)

			// Handle with all subscribed handlers
			for _, handler := range handlers {
				if err := handler.Handle(&event); err != nil {
					log.Printf("[RabbitMQ] Error handling %s: %v", eventName, err)
				}
			}

			msg.Ack(false)
		}
	}
}

// Stop stops the event bus
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

// GetExchange returns the exchange name
func (b *RabbitMQEventBus) GetExchange() string {
	return b.exchange
}

// GetEnv returns the environment
func (b *RabbitMQEventBus) GetEnv() string {
	return b.env
}

// GetTenant returns the tenant
func (b *RabbitMQEventBus) GetTenant() string {
	return b.tenant
}
