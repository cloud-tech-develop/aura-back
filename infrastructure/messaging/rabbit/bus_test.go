package rabbit

import (
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestRabbitMQConnection(t *testing.T) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Printf("⚠️  No se pudo cargar .env: %v", err)
		return
	}

	bus, err := NewRabbitMQEventBus()
	if err != nil {
		t.Skipf("Skipping RabbitMQ test: %v", err)
		return
	}
	defer bus.Stop()

	assert.NotNil(t, bus)
	assert.Equal(t, "aura.dev", bus.GetExchange())
	assert.Equal(t, "dev", bus.GetEnv())
}

func TestRabbitMQConnectionWithEnv(t *testing.T) {
	t.Setenv("APP_ENV", "dev")

	bus, err := NewRabbitMQEventBus()
	if err != nil {
		t.Skipf("Skipping RabbitMQ test: %v", err)
		return
	}
	defer bus.Stop()

	assert.NotNil(t, bus)
	assert.Equal(t, "aura.dev", bus.GetExchange())
	assert.Equal(t, "dev", bus.GetEnv())
}

func TestRabbitMQEventBusWithTenant(t *testing.T) {
	t.Setenv("APP_ENV", "prod")

	bus, err := NewRabbitMQEventBusWithTenant("empresa_uno")
	if err != nil {
		t.Skipf("Skipping RabbitMQ test: %v", err)
		return
	}
	defer bus.Stop()

	assert.NotNil(t, bus)
	assert.Equal(t, "empresa_uno", bus.GetTenant())
}
