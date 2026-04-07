package service

import (
	"core/ai"
	"core/config"
	"core/model"
	"core/repository"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// AIModelService handles AI model config business logic
type AIModelService struct {
	repo      *repository.AIModelRepository
	logger    *zap.Logger
	aiService *ai.AnalysisService
}

// NewAIModelService creates a new AI model service
func NewAIModelService(repo *repository.AIModelRepository, logger *zap.Logger) *AIModelService {
	return &AIModelService{
		repo:   repo,
		logger: logger,
	}
}

// GetAllConfigs retrieves all AI model configs
func (s *AIModelService) GetAllConfigs() ([]model.AIModelConfig, error) {
	return s.repo.GetAllConfigs()
}

// GetConfigByID retrieves an AI model config by ID
func (s *AIModelService) GetConfigByID(id int64) (*model.AIModelConfig, error) {
	return s.repo.GetConfigByID(id)
}

// CreateConfig creates a new AI model config
func (s *AIModelService) CreateConfig(config *model.AIModelConfig, user string) error {
	config.CreatedBy = user
	config.UpdatedBy = user
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()

	// Validate required fields
	if config.Name == "" || config.Provider == "" || config.Model == "" || config.APIKey == "" {
		return fmt.Errorf("name, provider, model and api_key are required")
	}

	// Set defaults
	if config.MaxTokens == 0 {
		config.MaxTokens = 4000
	}
	if config.Temperature == 0 {
		config.Temperature = 0.7
	}
	if config.Timeout == 0 {
		config.Timeout = 60
	}

	return s.repo.CreateConfig(config)
}

// UpdateConfig updates an AI model config
func (s *AIModelService) UpdateConfig(config *model.AIModelConfig, user string) error {
	config.UpdatedBy = user
	config.UpdatedAt = time.Now()

	// Validate required fields
	if config.Name == "" || config.Provider == "" || config.Model == "" || config.APIKey == "" {
		return fmt.Errorf("name, provider, model and api_key are required")
	}

	return s.repo.UpdateConfig(config)
}

// DeleteConfig deletes an AI model config
func (s *AIModelService) DeleteConfig(id int64) error {
	return s.repo.DeleteConfig(id)
}

// TestConfig tests an AI model config
func (s *AIModelService) TestConfig(id int64) (bool, string, error) {
	cfg, err := s.repo.GetConfigByID(id)
	if err != nil {
		return false, "", err
	}

	if !cfg.IsEnabled {
		return false, "Config is disabled", nil
	}

	// Create temporary AI config
	aiConfig := &config.AIConfig{
		Provider:    cfg.Provider,
		Model:       cfg.Model,
		APIKey:      cfg.APIKey,
		BaseURL:     cfg.BaseURL,
		MaxTokens:   cfg.MaxTokens,
		Temperature: cfg.Temperature,
		Timeout:     cfg.Timeout,
		Enabled:     cfg.IsEnabled,
	}

	// Create temporary AI service for testing
	tempService, err := ai.NewAnalysisService(aiConfig, s.logger)
	if err != nil {
		return false, fmt.Sprintf("Failed to create AI client: %v", err), nil
	}

	// Test health
	if err := tempService.Health(); err != nil {
		return false, fmt.Sprintf("AI health check failed: %v", err), nil
	}

	return true, "Connection successful", nil
}

// SetDefaultConfig sets a config as default
func (s *AIModelService) SetDefaultConfig(id int64) error {
	return s.repo.SetDefaultConfig(id)
}

// GetActiveConfig gets the active AI config (default enabled config)
func (s *AIModelService) GetActiveConfig() (*model.AIModelConfig, error) {
	return s.repo.GetDefaultConfig()
}

// BuildAIConfigFromModel builds AI config from model config
func (s *AIModelService) BuildAIConfigFromModel(cfg *model.AIModelConfig) *config.AIConfig {
	return &config.AIConfig{
		Provider:    cfg.Provider,
		Model:       cfg.Model,
		APIKey:      cfg.APIKey,
		BaseURL:     cfg.BaseURL,
		MaxTokens:   cfg.MaxTokens,
		Temperature: cfg.Temperature,
		Timeout:     cfg.Timeout,
		Enabled:     cfg.IsEnabled,
	}
}
