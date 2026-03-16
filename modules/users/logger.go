package users

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

	// Custom filename for users
	date := time.Now().Format("20060102")
	filename := fmt.Sprintf("users_%s.log", date)
	fullPath := filepath.Join(h.LogDir, filename)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	log.Printf("[Users Logger] Logged event: %s to %s", event.GetName(), filename)

	return nil
}
