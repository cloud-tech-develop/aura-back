package products

import (
	"fmt"

	"github.com/cloud-tech-develop/aura-back/shared/logging"
)

const (
	EventProductCreated = "catalog.product.created"
	EventProductUpdated = "catalog.product.updated"
	EventProductDeleted = "catalog.product.deleted"
)

type LoggerHandler struct {
	*logging.LoggerHandler
}

func NewLoggerHandler(logDir string) *LoggerHandler {
	return &LoggerHandler{
		LoggerHandler: logging.NewLoggerHandler(logDir),
	}
}

func (l *LoggerHandler) Handle(event interface{}) {
	switch e := event.(type) {
	case Product:
		fmt.Printf("[Catalog/Products Logger] Product event: %+v\n", e)
	default:
		fmt.Printf("[Catalog/Products Logger] Unknown event: %+v\n", e)
	}
}
