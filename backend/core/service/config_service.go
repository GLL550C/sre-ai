package service

import (
	"core/model"
	"core/repository"
	"strconv"

	"go.uber.org/zap"
)

// ConfigService 配置服务
type ConfigService struct {
	repo   *repository.ConfigRepository
	logger *zap.Logger
}

// NewConfigService 创建服务
func NewConfigService(repo *repository.ConfigRepository, logger *zap.Logger) *ConfigService {
	return &ConfigService{repo: repo, logger: logger}
}

// GetByCategory 按分类获取
func (s *ConfigService) GetByCategory(category string) ([]model.SystemConfig, error) {
	return s.repo.GetByCategory(category)
}

// GetByKey 根据key获取
func (s *ConfigService) GetByKey(key string) (*model.SystemConfig, error) {
	return s.repo.GetByKey(key)
}

// GetValue 获取配置值
func (s *ConfigService) GetValue(key string) (string, error) {
	return s.repo.GetValue(key)
}

// Update 更新配置
func (s *ConfigService) Update(key, value, user string) error {
	return s.repo.Update(key, value, user)
}

// GetAll 获取所有配置
func (s *ConfigService) GetAll() ([]model.SystemConfig, error) {
	return s.repo.GetAll()
}

// GetValueAsBool 获取布尔值
func (s *ConfigService) GetValueAsBool(key string) bool {
	v, err := s.repo.GetValue(key)
	if err != nil {
		return false
	}
	return v == "true" || v == "1"
}

// GetValueAsInt 获取整数值
func (s *ConfigService) GetValueAsInt(key string) int {
	v, err := s.repo.GetValue(key)
	if err != nil {
		return 0
	}
	i, _ := strconv.Atoi(v)
	return i
}

// GetByCategoryAndSubCategory 按分类和子分类获取配置
func (s *ConfigService) GetByCategoryAndSubCategory(category, subCategory string) ([]model.SystemConfig, error) {
	return s.repo.GetByCategoryAndSubCategory(category, subCategory)
}

// GetAppName 获取应用名称
func (s *ConfigService) GetAppName() string {
	name, _ := s.repo.GetValue("app.name")
	if name == "" {
		return "SRE AI Platform"
	}
	return name
}
