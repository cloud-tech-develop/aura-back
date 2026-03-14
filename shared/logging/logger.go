package logging

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/cloud-tech-develop/aura-back/shared/events"
)

type LoggerHandler struct {
	LogDir string
}

func NewLoggerHandler(logDir string) *LoggerHandler {
	if logDir == "" {
		logDir = "logs"
	}
	return &LoggerHandler{LogDir: logDir}
}

func (h *LoggerHandler) Handle(event events.Event) error {
	payload := event.GetPayload()

	filename := h.getLogFilename(event.GetName())
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

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
	filename := fmt.Sprintf("system_%s.log", date)
	return filepath.Join(h.LogDir, filename)
}
