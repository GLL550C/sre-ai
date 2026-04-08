package service

import (
	"core/repository"
	"go.uber.org/zap"
)

// DashboardService 仪表板服务
type DashboardService struct {
	alertRepo   *repository.AlertRepository
	clusterRepo *repository.ClusterRepository
	logger      *zap.Logger
}

// NewDashboardService 创建仪表板服务
func NewDashboardService(alertRepo *repository.AlertRepository, clusterRepo *repository.ClusterRepository, logger *zap.Logger) *DashboardService {
	return &DashboardService{
		alertRepo:   alertRepo,
		clusterRepo: clusterRepo,
		logger:      logger,
	}
}

// GetDashboard 获取仪表板数据
func (s *DashboardService) GetDashboard(tenantID *int64) (map[string]interface{}, error) {
	// 获取集群数量
	clusters, err := s.clusterRepo.GetAll()
	if err != nil {
		s.logger.Error("获取集群列表失败", zap.Error(err))
	}

	// 获取告警统计
	firingCount, _ := s.alertRepo.GetCount("firing", "")
	resolvedCount, _ := s.alertRepo.GetCount("resolved", "")
	acknowledgedCount, _ := s.alertRepo.GetCount("acknowledged", "")

	// 获取严重级别统计
	criticalCount, _ := s.alertRepo.GetCount("", "critical")
	warningCount, _ := s.alertRepo.GetCount("", "warning")
	infoCount, _ := s.alertRepo.GetCount("", "info")

	dashboard := map[string]interface{}{
		"cluster_count":      len(clusters),
		"firing_alerts":      firingCount,
		"resolved_alerts":    resolvedCount,
		"acknowledged_count": acknowledgedCount,
		"critical_count":     criticalCount,
		"warning_count":      warningCount,
		"info_count":         infoCount,
	}

	return dashboard, nil
}
