package config

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"log"
	"os"
	"strings"
)

// Config contains the application configuration
type Config struct {
	ServerPort         string
	WhatsAppAPIKey     string
	WhatsAppBaseURL    string
	WebhookVerifyToken string
	LLMUrl             string
}

// Load loads configuration from environment variables or an .env file
func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
	setLogLevel()
	return &Config{
		ServerPort:         getEnv("SERVER_PORT", "8080"),
		WhatsAppAPIKey:     getEnv("WHATSAPP_API_KEY", ""),
		WhatsAppBaseURL:    getEnv("WHATSAPP_BASE_URL", "https://graph.facebook.com/v18.0"),
		WebhookVerifyToken: getEnv("WEBHOOK_VERIFY_TOKEN", ""),
		LLMUrl:             getEnv("LLM_URL", "http://localhost:8081/api/v1/chat/ask"),
	}
}

// getEnv retrieves environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// setLogLevel sets the log level defined in LOG_LEVEL environment variable
func setLogLevel() {
	levels := map[string]zerolog.Level{
		"debug": zerolog.DebugLevel,
		"info":  zerolog.InfoLevel,
		"warn":  zerolog.WarnLevel,
		"error": zerolog.ErrorLevel,
		"fatal": zerolog.FatalLevel,
		"panic": zerolog.PanicLevel,
	}
	levelEnv := strings.ToLower(getEnv("LOG_LEVEL", "info"))

	level, ok := levels[levelEnv]
	if !ok {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
}
