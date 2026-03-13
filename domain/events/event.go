package events

import "time"

type Event interface {
	GetName() string
	GetPayload() interface{}
	GetTimestamp() time.Time
}

type EventHandler interface {
	Handle(event Event) error
}

type EventBus interface {
	Publish(event Event) error
	Subscribe(eventName string, handler EventHandler) error
	Unsubscribe(eventName string, handler EventHandler) error
	Start() error
	Stop() error
}

type BaseEvent struct {
	name      string
	payload   interface{}
	timestamp time.Time
}

func NewBaseEvent(name string, payload interface{}) BaseEvent {
	return BaseEvent{
		name:      name,
		payload:   payload,
		timestamp: time.Now(),
	}
}

func (e BaseEvent) GetName() string {
	return e.name
}

func (e BaseEvent) GetPayload() interface{} {
	return e.payload
}

func (e BaseEvent) GetTimestamp() time.Time {
	return e.timestamp
}
