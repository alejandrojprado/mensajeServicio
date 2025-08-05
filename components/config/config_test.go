package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("AWS_REGION", "us-east-1")
	os.Setenv("PORT", "8080")
	os.Setenv("DDB_TABLE_MENSAJES", "test-messages-table")
	os.Setenv("DDB_TABLE_SEGUIDORES", "test-followers-table")
	os.Setenv("DDB_TABLE_TIMELINE", "test-timeline-table")
	os.Setenv("MAX_MESSAGE_LENGTH", "280")
	os.Setenv("DEFAULT_LIMIT", "20")

	defer func() {
		os.Unsetenv("AWS_REGION")
		os.Unsetenv("PORT")
		os.Unsetenv("DDB_TABLE_MENSAJES")
		os.Unsetenv("DDB_TABLE_SEGUIDORES")
		os.Unsetenv("DDB_TABLE_TIMELINE")
		os.Unsetenv("MAX_MESSAGE_LENGTH")
		os.Unsetenv("DEFAULT_LIMIT")
	}()

	config := LoadConfig()

	assert.NotNil(t, config)
	assert.Equal(t, "us-east-1", config.Region)
	assert.Equal(t, "8080", config.Port)
	assert.Equal(t, "test-messages-table", config.TableMensajesName)
	assert.Equal(t, "test-followers-table", config.TableSeguidoresName)
	assert.Equal(t, "test-timeline-table", config.TableTimelineName)
	assert.Equal(t, 280, config.MaxMessageLength)
	assert.Equal(t, 20, config.DefaultLimit)
}

func TestLoadConfig_Defaults(t *testing.T) {
	// Limpiar todas las variables de entorno
	os.Unsetenv("AWS_REGION")
	os.Unsetenv("PORT")
	os.Unsetenv("DDB_TABLE_MENSAJES")
	os.Unsetenv("DDB_TABLE_SEGUIDORES")
	os.Unsetenv("DDB_TABLE_TIMELINE")
	os.Unsetenv("MAX_MESSAGE_LENGTH")
	os.Unsetenv("DEFAULT_LIMIT")

	config := LoadConfig()

	assert.NotNil(t, config)
	assert.Equal(t, "us-east-1", config.Region)   // Valor por defecto
	assert.Equal(t, "80", config.Port)            // Valor por defecto
	assert.Equal(t, 280, config.MaxMessageLength) // Valor por defecto
	assert.Equal(t, 20, config.DefaultLimit)      // Valor por defecto
}
