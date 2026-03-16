package products

import (
	"fmt"
	"github.com/cloud-tech-develop/aura-back/shared/logging"
)

const (
	EventCategoryCreated = "category.created"
	EventBrandCreated    = "brand.created"
	EventProductCreated  = "product.created"
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
	// Type assertion to get event details
	switch e := event.(type) {
	case Category:
		fmt.Printf("[Products Logger] Category event: %+v\n", e)
	case Brand:
		fmt.Printf("[Products Logger] Brand event: %+v\n", e)
	case Product:
		fmt.Printf("[Products Logger] Product event: %+v\n", e)
	default:
		fmt.Printf("[Products Logger] Unknown event: %+v\n", e)
	}
}
