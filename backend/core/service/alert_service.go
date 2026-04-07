package service

import (
	"context"
	"core/model"
	"core/repository"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// AlertService handles alert business logic
type AlertService struct {
	alertRepo   *repository.AlertRepository
	clusterRepo *repository.ClusterRepository
	redisClient *redis.Client
	logger      *zap.Logger
}

// NewAlertService creates a new alert service
func NewAlertService(alertRepo *repository.AlertRepository, clusterRepo *repository.ClusterRepository, redisClient *redis.Client, logger *zap.Logger) *AlertService {
	return &AlertService{
		alertRepo:   alertRepo,
		clusterRepo: clusterRepo,
		redisClient: redisClient,
		logger:      logger,
	}
}

// GetAlerts retrieves alerts with filters
func (s *AlertService) GetAlerts(status, severity string, page, pageSize int) ([]model.Alert, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// Try cache first
	cacheKey := fmt.Sprintf("alerts:%s:%s:%d:%d", status, severity, page, pageSize)
	ctx := context.Background()
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var result struct {
			Alerts []model.Alert `json:"alerts"`
			Count  int           `json:"count"`
		}
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			return result.Alerts, result.Count, nil
		}
	}

	alerts, err := s.alertRepo.GetAlerts(status, severity, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.alertRepo.GetAlertCount(status, severity)
	if err != nil {
		return nil, 0, err
	}

	// Cache result
	result := struct {
		Alerts []model.Alert `json:"alerts"`
		Count  int           `json:"count"`
	}{alerts, count}
	resultJSON, _ := json.Marshal(result)
	s.redisClient.Set(ctx, cacheKey, resultJSON, 30*time.Second)

	return alerts, count, nil
}

// GetAlert retrieves an alert by ID
func (s *AlertService) GetAlert(id int64) (*model.Alert, error) {
	return s.alertRepo.GetAlertByID(id)
}

// CreateAlert creates a new alert
func (s *AlertService) CreateAlert(alert *model.Alert) error {
	return s.alertRepo.CreateAlert(alert)
}

// AcknowledgeAlert acknowledges an alert
func (s *AlertService) AcknowledgeAlert(id int64, user string) error {
	return s.alertRepo.AcknowledgeAlert(id, user)
}

// ProcessWebhook processes alertmanager webhook
func (s *AlertService) ProcessWebhook(payload *model.WebhookPayload) error {
	for _, webhookAlert := range payload.Alerts {
		// Check if alert already exists
		existingAlert, err := s.alertRepo.GetAlertByFingerprint(webhookAlert.Fingerprint)
		if err != nil && err.Error() != "sql: no rows in result set" {
			s.logger.Error("Failed to check existing alert", zap.Error(err))
			continue
		}

		labelsJSON, _ := json.Marshal(webhookAlert.Labels)

		if webhookAlert.Status == "firing" {
			if existingAlert != nil {
				// Update existing alert
				s.alertRepo.UpdateAlertStatus(existingAlert.ID, "firing")
			} else {
				// Create new alert
				severity := webhookAlert.Labels["severity"]
				if severity == "" {
					severity = "warning"
				}

				alert := &model.Alert{
					Fingerprint: webhookAlert.Fingerprint,
					Status:      "firing",
					Severity:    severity,
					Summary:     webhookAlert.Annotations["summary"],
					Description: webhookAlert.Annotations["description"],
					Labels:      labelsJSON,
					StartsAt:    webhookAlert.StartsAt,
				}

				if err := s.alertRepo.CreateAlert(alert); err != nil {
					s.logger.Error("Failed to create alert", zap.Error(err))
					continue
				}
			}
		} else if webhookAlert.Status == "resolved" {
			if existingAlert != nil {
				s.alertRepo.ResolveAlert(webhookAlert.Fingerprint)
			}
		}
	}

	// Invalidate cache
	ctx := context.Background()
	iter := s.redisClient.Scan(ctx, 0, "alerts:*", 0).Iterator()
	for iter.Next(ctx) {
		s.redisClient.Del(ctx, iter.Val())
	}

	return nil
}

// GetAlertStats gets alert statistics
func (s *AlertService) GetAlertStats() (map[string]interface{}, error) {
	ctx := context.Background()

	// Try cache
	cached, err := s.redisClient.Get(ctx, "alert:stats").Result()
	if err == nil {
		var stats map[string]interface{}
		if err := json.Unmarshal([]byte(cached), &stats); err == nil {
			return stats, nil
		}
	}

	// Calculate stats
	firingCount, _ := s.alertRepo.GetAlertCount("firing", "")
	criticalCount, _ := s.alertRepo.GetAlertCount("firing", "critical")
	warningCount, _ := s.alertRepo.GetAlertCount("firing", "warning")

	stats := map[string]interface{}{
		"firing":   firingCount,
		"critical": criticalCount,
		"warning":  warningCount,
	}

	// Cache stats
	statsJSON, _ := json.Marshal(stats)
	s.redisClient.Set(ctx, "alert:stats", statsJSON, 60*time.Second)

	return stats, nil
}
