package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "Environment variable exists",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "env_value",
			expected:     "env_value",
		},
		{
			name:         "Environment variable does not exist",
			key:          "NON_EXISTENT_KEY",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
		{
			name:         "Empty environment variable",
			key:          "EMPTY_KEY",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up environment
			os.Unsetenv(tt.key)
			
			// Set environment variable if needed
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnv(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLoadConfig_WithDefaults(t *testing.T) {
	// Clean up environment variables
	envVars := []string{
		"SERVER_PORT",
		"WHATSAPP_API_KEY",
		"WHATSAPP_BASE_URL",
		"WEBHOOK_VERIFY_TOKEN",
		"LLM_URL",
		"LLM_BEARER_TOKEN",
	}

	for _, env := range envVars {
		os.Unsetenv(env)
	}

	// Note: We can't easily test Load() function because it calls godotenv.Load()
	// and log.Fatal() which would exit the test. Instead we test the structure
	// and individual functions.

	// Test default values by calling getEnv directly
	assert.Equal(t, "8080", getEnv("SERVER_PORT", "8080"))
	assert.Equal(t, "", getEnv("WHATSAPP_API_KEY", ""))
	assert.Equal(t, "https://graph.facebook.com/v20.0", getEnv("WHATSAPP_BASE_URL", "https://graph.facebook.com/v20.0"))
	assert.Equal(t, "", getEnv("WEBHOOK_VERIFY_TOKEN", ""))
	assert.Equal(t, "http://localhost:8081/api/v1/chat/ask", getEnv("LLM_URL", "http://localhost:8081/api/v1/chat/ask"))
	assert.Equal(t, "", getEnv("LLM_BEARER_TOKEN", ""))
}

func TestLoadConfig_WithEnvironmentVariables(t *testing.T) {
	// Set environment variables
	testEnvVars := map[string]string{
		"SERVER_PORT":          "9090",
		"WHATSAPP_API_KEY":     "test-whatsapp-key",
		"WHATSAPP_BASE_URL":    "https://test.whatsapp.com",
		"WEBHOOK_VERIFY_TOKEN": "test-webhook-token",
		"LLM_URL":              "https://test.llm.com/api",
		"LLM_BEARER_TOKEN":     "test-llm-token",
	}

	// Set environment variables
	for key, value := range testEnvVars {
		os.Setenv(key, value)
		defer os.Unsetenv(key)
	}

	// Test that getEnv returns the environment values
	assert.Equal(t, "9090", getEnv("SERVER_PORT", "8080"))
	assert.Equal(t, "test-whatsapp-key", getEnv("WHATSAPP_API_KEY", ""))
	assert.Equal(t, "https://test.whatsapp.com", getEnv("WHATSAPP_BASE_URL", "https://graph.facebook.com/v20.0"))
	assert.Equal(t, "test-webhook-token", getEnv("WEBHOOK_VERIFY_TOKEN", ""))
	assert.Equal(t, "https://test.llm.com/api", getEnv("LLM_URL", "http://localhost:8081/api/v1/chat/ask"))
	assert.Equal(t, "test-llm-token", getEnv("LLM_BEARER_TOKEN", ""))
}

func TestConfigStruct(t *testing.T) {
	config := Config{
		ServerPort:         "8080",
		WhatsAppAPIKey:     "test-key",
		WhatsAppBaseURL:    "https://test.com",
		WebhookVerifyToken: "verify-token",
		LLMUrl:             "https://llm.test.com",
		LLMBearerToken:     "bearer-token",
	}

	assert.Equal(t, "8080", config.ServerPort)
	assert.Equal(t, "test-key", config.WhatsAppAPIKey)
	assert.Equal(t, "https://test.com", config.WhatsAppBaseURL)
	assert.Equal(t, "verify-token", config.WebhookVerifyToken)
	assert.Equal(t, "https://llm.test.com", config.LLMUrl)
	assert.Equal(t, "bearer-token", config.LLMBearerToken)
}
