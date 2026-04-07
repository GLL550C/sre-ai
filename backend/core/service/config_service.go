package service

import (
	"core/model"
	"core/repository"
	"encoding/json"
	"fmt"
	"strconv"

	"go.uber.org/zap"
)

// ConfigService 配置服务
type ConfigService struct {
	repo   *repository.ConfigRepository
	logger *zap.Logger
}

// NewConfigService 创建配置服务
func NewConfigService(repo *repository.ConfigRepository, logger *zap.Logger) *ConfigService {
	return &ConfigService{
		repo:   repo,
		logger: logger,
	}
}

// InitConfigItems 初始化配置项
func (s *ConfigService) InitConfigItems() error {
	return s.repo.InitConfigItems()
}

// GetConfigTree 获取配置树
func (s *ConfigService) GetConfigTree() ([]model.ConfigTree, error) {
	return s.repo.GetConfigTree()
}

// GetConfigItemsByCategory 获取分类下的配置项
func (s *ConfigService) GetConfigItemsByCategory(category, subCategory string) ([]model.ConfigItem, error) {
	return s.repo.GetConfigItemsByCategory(category, subCategory)
}

// GetConfigItem 获取单个配置项
func (s *ConfigService) GetConfigItem(key string) (*model.ConfigItem, error) {
	return s.repo.GetConfigItem(key)
}

// UpdateConfigValue 更新配置值
func (s *ConfigService) UpdateConfigValue(key, value, user string) error {
	// 先获取配置项验证类型
	item, err := s.repo.GetConfigItem(key)
	if err != nil {
		return err
	}

	// 验证值类型
	if err := s.validateValue(item, value); err != nil {
		return err
	}

	return s.repo.UpdateConfigValue(key, value, user)
}

// UpdateMultipleConfigs 批量更新配置
func (s *ConfigService) UpdateMultipleConfigs(configs map[string]string, user string) error {
	for key, value := range configs {
		if err := s.UpdateConfigValue(key, value, user); err != nil {
			s.logger.Error("Failed to update config", zap.String("key", key), zap.Error(err))
			return fmt.Errorf("failed to update %s: %w", key, err)
		}
	}
	return nil
}

// GetConfigValue 获取配置值
func (s *ConfigService) GetConfigValue(key string) (string, error) {
	return s.repo.GetConfigValue(key)
}

// GetConfigValueAsBool 获取布尔配置值
func (s *ConfigService) GetConfigValueAsBool(key string) bool {
	value, err := s.repo.GetConfigValue(key)
	if err != nil {
		return false
	}
	return value == "true" || value == "1"
}

// GetConfigValueAsInt 获取整数配置值
func (s *ConfigService) GetConfigValueAsInt(key string) int {
	value, err := s.repo.GetConfigValue(key)
	if err != nil {
		return 0
	}
	i, _ := strconv.Atoi(value)
	return i
}

// GetConfigValueAsFloat 获取浮点配置值
func (s *ConfigService) GetConfigValueAsFloat(key string) float64 {
	value, err := s.repo.GetConfigValue(key)
	if err != nil {
		return 0
	}
	f, _ := strconv.ParseFloat(value, 64)
	return f
}

// GetAllConfigValues 获取所有配置值
func (s *ConfigService) GetAllConfigValues() (map[string]string, error) {
	return s.repo.GetAllConfigValues()
}

// ResetConfigToDefault 重置配置为默认值
func (s *ConfigService) ResetConfigToDefault(key string) error {
	return s.repo.ResetConfigToDefault(key)
}

// GetAIConfig 获取AI配置
func (s *ConfigService) GetAIConfig() (map[string]interface{}, error) {
	return s.repo.GetAIConfig()
}

// GetAIModelConfigs 获取所有AI模型配置
func (s *ConfigService) GetAIModelConfigs() ([]model.AIModelConfig, error) {
	// 从ai_model_configs表获取
	// 这里需要调用ai_model_repository的方法
	// 暂时返回空，后续整合
	return []model.AIModelConfig{}, nil
}

// ExportConfig 导出配置
func (s *ConfigService) ExportConfig() ([]model.ConfigItem, error) {
	return s.repo.ExportConfig()
}

// ImportConfig 导入配置
func (s *ConfigService) ImportConfig(items []model.ConfigItem, user string) error {
	return s.repo.ImportConfig(items, user)
}

// validateValue 验证配置值类型
func (s *ConfigService) validateValue(item *model.ConfigItem, value string) error {
	switch item.Type {
	case "string":
		return nil
	case "number":
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return fmt.Errorf("value must be a number")
		}
	case "boolean":
		if value != "true" && value != "false" {
			return fmt.Errorf("value must be true or false")
		}
	case "json":
		var js interface{}
		if err := json.Unmarshal([]byte(value), &js); err != nil {
			return fmt.Errorf("value must be valid JSON")
		}
	case "password":
		// 密码类型不验证格式
		return nil
	}

	// 验证可选项
	if item.Options != "" {
		var options []string
		if err := json.Unmarshal([]byte(item.Options), &options); err == nil {
			valid := false
			for _, opt := range options {
				if opt == value {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("value must be one of %v", options)
			}
		}
	}

	return nil
}

// GetSystemSettings 获取系统设置
func (s *ConfigService) GetSystemSettings() map[string]interface{} {
	settings := make(map[string]interface{})

	// 基础配置
	settings["app.name"], _ = s.GetConfigValue("app.name")
	settings["app.version"], _ = s.GetConfigValue("app.version")
	settings["app.timezone"], _ = s.GetConfigValue("app.timezone")
	settings["app.language"], _ = s.GetConfigValue("app.language")

	// 系统配置
	settings["system.session_timeout"] = s.GetConfigValueAsInt("system.session_timeout")
	settings["system.cache_ttl"] = s.GetConfigValueAsInt("system.cache_ttl")
	settings["system.log_level"], _ = s.GetConfigValue("system.log_level")
	settings["system.enable_registration"] = s.GetConfigValueAsBool("system.enable_registration")

	return settings
}

// GetNotificationSettings 获取通知设置
func (s *ConfigService) GetNotificationSettings() map[string]interface{} {
	emailHost, _ := s.GetConfigValue("notification.email_host")
	emailUser, _ := s.GetConfigValue("notification.email_user")
	webhookURL, _ := s.GetConfigValue("notification.webhook_url")
	return map[string]interface{}{
		"email_enabled":   s.GetConfigValueAsBool("notification.email_enabled"),
		"email_host":      emailHost,
		"email_port":      s.GetConfigValueAsInt("notification.email_port"),
		"email_user":      emailUser,
		"webhook_enabled": s.GetConfigValueAsBool("notification.webhook_enabled"),
		"webhook_url":     webhookURL,
	}
}

// GetAISettings 获取AI设置
func (s *ConfigService) GetAISettings() map[string]interface{} {
	systemPrompt, _ := s.GetConfigValue("ai.system_prompt")
	welcomeMessage, _ := s.GetConfigValue("ai.welcome_message")
	return map[string]interface{}{
		"analysis_timeout":     s.GetConfigValueAsInt("ai.analysis_timeout"),
		"auto_analysis":        s.GetConfigValueAsBool("ai.auto_analysis"),
		"confidence_threshold": s.GetConfigValueAsFloat("ai.confidence_threshold"),
		"include_metrics":      s.GetConfigValueAsBool("ai.include_metrics"),
		"include_logs":         s.GetConfigValueAsBool("ai.include_logs"),
		"chat_context_length":  s.GetConfigValueAsInt("ai.chat_context_length"),
		"chat_enable_memory":   s.GetConfigValueAsBool("ai.chat_enable_memory"),
		"system_prompt":        systemPrompt,
		"welcome_message":      welcomeMessage,
	}
}
