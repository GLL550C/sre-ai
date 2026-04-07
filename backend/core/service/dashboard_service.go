package service

import (
	"core/model"
	"core/repository"

	"go.uber.org/zap"
)

// DashboardService handles dashboard business logic
type DashboardService struct {
	dashboardRepo *repository.DashboardRepository
	logger        *zap.Logger
}

// NewDashboardService creates a new dashboard service
func NewDashboardService(dashboardRepo *repository.DashboardRepository, logger *zap.Logger) *DashboardService {
	return &DashboardService{
		dashboardRepo: dashboardRepo,
		logger:        logger,
	}
}

// GetDashboard retrieves the default dashboard
func (s *DashboardService) GetDashboard(tenantID *int64) (*model.DashboardConfig, error) {
	return s.dashboardRepo.GetDefaultDashboard(tenantID)
}

// GetDashboardConfigs retrieves dashboard configurations
func (s *DashboardService) GetDashboardConfigs(tenantID *int64) ([]model.DashboardConfig, error) {
	return s.dashboardRepo.GetDashboardConfigs(tenantID)
}

// CreateDashboardConfig creates a new dashboard configuration
func (s *DashboardService) CreateDashboardConfig(config *model.DashboardConfig) error {
	return s.dashboardRepo.CreateDashboardConfig(config)
}

// UpdateDashboardConfig updates a dashboard configuration
func (s *DashboardService) UpdateDashboardConfig(config *model.DashboardConfig) error {
	return s.dashboardRepo.UpdateDashboardConfig(config)
}

// DeleteDashboardConfig deletes a dashboard configuration
func (s *DashboardService) DeleteDashboardConfig(id int64) error {
	return s.dashboardRepo.DeleteDashboardConfig(id)
}
