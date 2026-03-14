package enterprise

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

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
	log.Printf("[Enterprise Logger] Logged event: %s to %s", event.GetName(), filename)

	return nil
}

func (h *LoggerHandler) getLogFilename(eventName string) string {
	date := time.Now().Format("20060102")
	filename := fmt.Sprintf("enterprise_%s.log", date)
	return filepath.Join(h.LogDir, filename)
}

// Since I made logDir private in shared/logging, I'll update it to be public or add a method.
// Let's go back and make it public for simplicity in this POS context.
