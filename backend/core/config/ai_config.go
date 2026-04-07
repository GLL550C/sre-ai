package config

import (
	"encoding/json"
	"os"
)

// AIConfig represents AI model configuration
type AIConfig struct {
	Provider    string            `json:"provider" yaml:"provider"`       // openai, claude, azure, custom
	Model       string            `json:"model" yaml:"model"`             // gpt-4, gpt-3.5-turbo, claude-3-opus, etc.
	APIKey      string            `json:"api_key" yaml:"api_key"`         // API Key (can be from env)
	BaseURL     string            `json:"base_url" yaml:"base_url"`       // API Base URL
	MaxTokens   int               `json:"max_tokens" yaml:"max_tokens"`   // Max tokens per request
	Temperature float64           `json:"temperature" yaml:"temperature"` // Temperature (0-2)
	Timeout     int               `json:"timeout" yaml:"timeout"`         // Request timeout in seconds
	Enabled     bool              `json:"enabled" yaml:"enabled"`         // Enable AI analysis
	Extra       map[string]string `json:"extra" yaml:"extra"`             // Extra provider-specific config
}

// LoadAIConfigFromEnv loads AI config from environment variables
func LoadAIConfigFromEnv() *AIConfig {
	cfg := &AIConfig{
		Provider:    getEnv("AI_PROVIDER", "openai"),
		Model:       getEnv("AI_MODEL", "gpt-4"),
		APIKey:      getEnv("AI_API_KEY", ""),
		BaseURL:     getEnv("AI_BASE_URL", ""),
		MaxTokens:   getEnvInt("AI_MAX_TOKENS", 4000),
		Temperature: getEnvFloat("AI_TEMPERATURE", 0.7),
		Timeout:     getEnvInt("AI_TIMEOUT", 60),
		Enabled:     getEnvBool("AI_ENABLED", true),
	}

	// Load extra config from JSON env var
	if extraJSON := os.Getenv("AI_EXTRA_CONFIG"); extraJSON != "" {
		json.Unmarshal([]byte(extraJSON), &cfg.Extra)
	}

	return cfg
}

// getEnv gets environment variable with default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets integer environment variable
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		if err := json.Unmarshal([]byte(value), &result); err == nil {
			return result
		}
	}
	return defaultValue
}

// getEnvFloat gets float environment variable
func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		var result float64
		if err := json.Unmarshal([]byte(value), &result); err == nil {
			return result
		}
	}
	return defaultValue
}

// getEnvBool gets boolean environment variable
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1" || value == "yes"
	}
	return defaultValue
}
