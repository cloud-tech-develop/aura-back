package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/cloud-tech-develop/aura-back/domain/enterprise"
	"github.com/cloud-tech-develop/aura-back/domain/events"
)

type LoggerHandler struct {
	logDir string
}

func NewLoggerHandler(logDir string) *LoggerHandler {
	if logDir == "" {
		logDir = "logs"
	}
	return &LoggerHandler{logDir: logDir}
}

func (h *LoggerHandler) Handle(event events.Event) error {
	payload, ok := event.GetPayload().(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid payload type")
	}

	filename := h.getLogFilename(event.GetName())
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	logEntry := map[string]interface{}{
		"event":     event.GetName(),
		"timestamp": event.GetTimestamp().Format(time.RFC3339),
		"payload":   payload,
	}

	jsonEntry, err := json.Marshal(logEntry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	fmt.Fprintln(file, string(jsonEntry))
	log.Printf("[LoggerHandler] Logged event: %s to %s", event.GetName(), filename)

	return nil
}

func (h *LoggerHandler) getLogFilename(eventName string) string {
	date := time.Now().Format("20060102")
	filename := fmt.Sprintf("enterprise_%s.log", date)
	return filepath.Join(h.logDir, filename)
}

type EnterpriseLoggerHandler struct {
	*LoggerHandler
}

func NewEnterpriseLoggerHandler(logDir string) *EnterpriseLoggerHandler {
	return &EnterpriseLoggerHandler{
		LoggerHandler: NewLoggerHandler(logDir),
	}
}

func (h *EnterpriseLoggerHandler) HandleEnterpriseCreated(e *enterprise.Enterprise) error {
	event := enterprise.NewEnterpriseCreatedEvent(e)
	return h.Handle(event)
}

func (h *EnterpriseLoggerHandler) HandleEnterpriseUpdated(e *enterprise.Enterprise) error {
	event := enterprise.NewEnterpriseUpdatedEvent(e)
	return h.Handle(event)
}

func (h *EnterpriseLoggerHandler) HandleEnterpriseDeleted(e *enterprise.Enterprise) error {
	event := enterprise.NewEnterpriseDeletedEvent(e)
	return h.Handle(event)
}
