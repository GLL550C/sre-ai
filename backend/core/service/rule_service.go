package service

import (
	"core/model"
	"core/repository"

	"go.uber.org/zap"
)

// RuleService handles alert rule business logic
type RuleService struct {
	ruleRepo *repository.RuleRepository
	logger   *zap.Logger
}

// NewRuleService creates a new rule service
func NewRuleService(ruleRepo *repository.RuleRepository, logger *zap.Logger) *RuleService {
	return &RuleService{
		ruleRepo: ruleRepo,
		logger:   logger,
	}
}

// GetRules retrieves all active rules
func (s *RuleService) GetRules() ([]model.AlertRule, error) {
	return s.ruleRepo.GetRules(1)
}

// GetRule retrieves a rule by ID
func (s *RuleService) GetRule(id int64) (*model.AlertRule, error) {
	return s.ruleRepo.GetRuleByID(id)
}

// CreateRule creates a new alert rule
func (s *RuleService) CreateRule(rule *model.AlertRule) error {
	rule.Status = 1
	return s.ruleRepo.CreateRule(rule)
}

// UpdateRule updates an alert rule
func (s *RuleService) UpdateRule(rule *model.AlertRule) error {
	return s.ruleRepo.UpdateRule(rule)
}

// DeleteRule deletes an alert rule
func (s *RuleService) DeleteRule(id int64) error {
	return s.ruleRepo.DeleteRule(id)
}
