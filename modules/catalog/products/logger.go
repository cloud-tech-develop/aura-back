package products

import (
	"fmt"

	"github.com/cloud-tech-develop/aura-back/shared/events"
	"github.com/cloud-tech-develop/aura-back/shared/logging"
)

type LoggerHandler struct {
	*logging.LoggerHandler
}

func NewLoggerHandler(logDir string) *LoggerHandler {
	return &LoggerHandler{
		LoggerHandler: logging.NewLoggerHandler(logDir),
	}
}

func (l *LoggerHandler) Handle(event events.Event) error {
	payload := event.GetPayload()

	switch e := payload.(type) {
	case map[string]interface{}:
		fmt.Printf("[Catalog/Products Logger] Product event: %s - %+v\n", event.GetName(), e)
	case Product:
		fmt.Printf("[Catalog/Products Logger] Product event: %s - %+v\n", event.GetName(), e)
	default:
		fmt.Printf("[Catalog/Products Logger] Unknown event format: %s - %+v\n", event.GetName(), e)
	}

	return nil
}
