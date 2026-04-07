package model

import (
	"time"
)

// AIModelConfig represents AI model configuration
type AIModelConfig struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Provider     string    `json:"provider"` // openai, claude, azure, custom
	Model        string    `json:"model"`    // gpt-4, claude-3-opus, etc.
	APIKey       string    `json:"api_key"`
	BaseURL      string    `json:"base_url"`
	MaxTokens    int       `json:"max_tokens"`
	Temperature  float64   `json:"temperature"`
	Timeout      int       `json:"timeout"`
	IsDefault    bool      `json:"is_default"`
	IsEnabled    bool      `json:"is_enabled"`
	Description  string    `json:"description"`
	CreatedBy    string    `json:"created_by"`
	UpdatedBy    string    `json:"updated_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ToAIConfig converts to ai config format
func (c *AIModelConfig) ToAIConfig() map[string]interface{} {
	return map[string]interface{}{
		"provider":    c.Provider,
		"model":       c.Model,
		"api_key":     c.APIKey,
		"base_url":    c.BaseURL,
		"max_tokens":  c.MaxTokens,
		"temperature": c.Temperature,
		"timeout":     c.Timeout,
		"enabled":     c.IsEnabled,
	}
}
