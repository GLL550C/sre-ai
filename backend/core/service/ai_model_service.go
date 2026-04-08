package service

import (
	"core/ai"
	"core/config"
	"core/model"
	"core/repository"
	"fmt"

	"go.uber.org/zap"
)

// AIModelService AI模型业务逻辑
type AIModelService struct {
	repo   *repository.AIModelRepository
	logger *zap.Logger
}

// NewAIModelService 创建服务
func NewAIModelService(repo *repository.AIModelRepository, logger *zap.Logger) *AIModelService {
	return &AIModelService{repo: repo, logger: logger}
}

// GetAll 获取所有配置
func (s *AIModelService) GetAll() ([]model.AIModel, error) {
	return s.repo.GetAll()
}

// GetByID 根据ID获取
func (s *AIModelService) GetByID(id int64) (*model.AIModel, error) {
	return s.repo.GetByID(id)
}

// Create 创建配置
func (s *AIModelService) Create(m *model.AIModel, user string) error {
	m.CreatedBy = user
	m.UpdatedBy = user

	if m.Name == "" || m.Provider == "" || m.Model == "" || m.APIKey == "" {
		return fmt.Errorf("name, provider, model and api_key are required")
	}

	if m.MaxTokens == 0 {
		m.MaxTokens = 4000
	}
	if m.Temperature == 0 {
		m.Temperature = 0.7
	}
	if m.Timeout == 0 {
		m.Timeout = 60
	}

	return s.repo.Create(m)
}

// Update 更新配置
func (s *AIModelService) Update(m *model.AIModel, user string) error {
	m.UpdatedBy = user

	if m.Name == "" || m.Provider == "" || m.Model == "" || m.APIKey == "" {
		return fmt.Errorf("name, provider, model and api_key are required")
	}

	return s.repo.Update(m)
}

// Delete 删除配置
func (s *AIModelService) Delete(id int64) error {
	return s.repo.Delete(id)
}

// Test 测试配置
func (s *AIModelService) Test(id int64) (bool, string, error) {
	m, err := s.repo.GetByID(id)
	if err != nil {
		return false, "", err
	}

	if !m.IsEnabled {
		return false, "Config is disabled", nil
	}

	aiConfig := &config.AIConfig{
		Provider:    m.Provider,
		Model:       m.Model,
		APIKey:      m.APIKey,
		BaseURL:     m.BaseURL,
		MaxTokens:   m.MaxTokens,
		Temperature: m.Temperature,
		Timeout:     m.Timeout,
		Enabled:     m.IsEnabled,
	}

	tempService, err := ai.NewAnalysisService(aiConfig, s.logger)
	if err != nil {
		return false, fmt.Sprintf("Failed to create AI client: %v", err), nil
	}

	if err := tempService.Health(); err != nil {
		return false, fmt.Sprintf("AI health check failed: %v", err), nil
	}

	return true, "Connection successful", nil
}

// SetDefault 设置默认
func (s *AIModelService) SetDefault(id int64) error {
	return s.repo.SetDefault(id)
}

// GetDefault 获取默认配置
func (s *AIModelService) GetDefault() (*model.AIModel, error) {
	return s.repo.GetDefault()
}

// BuildAIConfig 构建AI配置
func (s *AIModelService) BuildAIConfig(m *model.AIModel) *config.AIConfig {
	return &config.AIConfig{
		Provider:    m.Provider,
		Model:       m.Model,
		APIKey:      m.APIKey,
		BaseURL:     m.BaseURL,
		MaxTokens:   m.MaxTokens,
		Temperature: m.Temperature,
		Timeout:     m.Timeout,
		Enabled:     m.IsEnabled,
	}
}

// GetActiveConfig 获取活跃配置(兼容旧代码)
func (s *AIModelService) GetActiveConfig() (*model.AIModel, error) {
	return s.repo.GetDefault()
}

// BuildAIConfigFromModel 从模型构建AI配置(兼容旧代码)
func (s *AIModelService) BuildAIConfigFromModel(m *model.AIModel) *config.AIConfig {
	return s.BuildAIConfig(m)
}
