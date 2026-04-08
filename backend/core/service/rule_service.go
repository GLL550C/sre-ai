package service

import (
	"core/model"
	"core/repository"

	"go.uber.org/zap"
)

// RuleService 告警规则业务逻辑
type RuleService struct {
	repo   *repository.RuleRepository
	logger *zap.Logger
}

// NewRuleService 创建服务
func NewRuleService(repo *repository.RuleRepository, logger *zap.Logger) *RuleService {
	return &RuleService{repo: repo, logger: logger}
}

// GetAll 获取所有规则
func (s *RuleService) GetAll() ([]model.AlertRule, error) {
	return s.repo.GetAll(1)
}

// GetByID 根据ID获取
func (s *RuleService) GetByID(id int64) (*model.AlertRule, error) {
	return s.repo.GetByID(id)
}

// Create 创建规则
func (s *RuleService) Create(rule *model.AlertRule) error {
	rule.Status = 1
	return s.repo.Create(rule)
}

// Update 更新规则
func (s *RuleService) Update(rule *model.AlertRule) error {
	return s.repo.Update(rule)
}

// Delete 删除规则
func (s *RuleService) Delete(id int64) error {
	return s.repo.Delete(id)
}
