package config

import (
	"os"
	"strconv"
)

type AppConfig struct {
	Env                 string
	Port                string
	TableMensajesName   string
	TableSeguidoresName string
	TableTimelineName   string
	Region              string
	BaseURL             string
	DefaultLimit        int
	MaxMessageLength    int
}

func LoadConfig() *AppConfig {
	defaultLimit, _ := strconv.Atoi(getEnv("DEFAULT_LIMIT", "20"))
	maxMessageLength, _ := strconv.Atoi(getEnv("MAX_MESSAGE_LENGTH", "280"))

	cfg := &AppConfig{
		Env:                 getEnv("ENV", "dev"),
		Port:                getEnv("PORT", "80"),
		TableMensajesName:   getEnv("DDB_TABLE_MENSAJES", "messages"),
		TableSeguidoresName: getEnv("DDB_TABLE_SEGUIDORES", "follows"),
		TableTimelineName:   getEnv("DDB_TABLE_TIMELINE", "timeline"),
		Region:              getEnv("AWS_REGION", "us-east-1"),
		BaseURL:             getEnv("BASE_URL", "http://localhost:8080/"),
		DefaultLimit:        defaultLimit,
		MaxMessageLength:    maxMessageLength,
	}
	return cfg
}

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
