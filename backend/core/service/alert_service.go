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

// AlertService 告警业务逻辑
type AlertService struct {
	repo        *repository.AlertRepository
	redisClient *redis.Client
	logger      *zap.Logger
}

// NewAlertService 创建服务
func NewAlertService(repo *repository.AlertRepository, logger *zap.Logger) *AlertService {
	return &AlertService{repo: repo, logger: logger}
}

// SetRedisClient 设置Redis客户端
func (s *AlertService) SetRedisClient(client *redis.Client) {
	s.redisClient = client
}

// GetAlerts 获取告警列表
func (s *AlertService) GetAlerts(status, severity string, page, pageSize int) ([]model.Alert, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// 尝试缓存
	if s.redisClient != nil {
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
	}

	alerts, err := s.repo.GetAlerts(status, severity, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.GetCount(status, severity)
	if err != nil {
		return nil, 0, err
	}

	// 缓存结果
	if s.redisClient != nil {
		result := struct {
			Alerts []model.Alert `json:"alerts"`
			Count  int           `json:"count"`
		}{alerts, count}
		resultJSON, _ := json.Marshal(result)
		ctx := context.Background()
		s.redisClient.Set(ctx, fmt.Sprintf("alerts:%s:%s:%d:%d", status, severity, page, pageSize), resultJSON, 30*time.Second)
	}

	return alerts, count, nil
}

// GetByID 根据ID获取
func (s *AlertService) GetByID(id int64) (*model.Alert, error) {
	return s.repo.GetByID(id)
}

// Create 创建告警
func (s *AlertService) Create(alert *model.Alert) error {
	return s.repo.Create(alert)
}

// Acknowledge 确认告警
func (s *AlertService) Acknowledge(id int64, user string) error {
	return s.repo.Acknowledge(id, user)
}

// ProcessWebhook 处理Webhook
func (s *AlertService) ProcessWebhook(payload *model.WebhookPayload) error {
	for _, wa := range payload.Alerts {
		existing, err := s.repo.GetByFingerprint(wa.Fingerprint)
		if err != nil && err.Error() != "sql: no rows in result set" {
			s.logger.Error("检查告警失败", zap.Error(err))
			continue
		}

		labelsJSON, _ := json.Marshal(wa.Labels)

		if wa.Status == "firing" {
			if existing != nil {
				s.repo.UpdateStatus(existing.ID, "firing")
			} else {
				severity := wa.Labels["severity"]
				if severity == "" {
					severity = "warning"
				}
				alert := &model.Alert{
					Fingerprint: wa.Fingerprint,
					Status:      "firing",
					Severity:    severity,
					Summary:     wa.Annotations["summary"],
					Description: wa.Annotations["description"],
					Labels:      labelsJSON,
					StartsAt:    wa.StartsAt,
				}
				if err := s.repo.Create(alert); err != nil {
					s.logger.Error("创建告警失败", zap.Error(err))
					continue
				}
			}
		} else if wa.Status == "resolved" {
			if existing != nil {
				s.repo.Resolve(wa.Fingerprint)
			}
		}
	}

	// 清除缓存
	if s.redisClient != nil {
		ctx := context.Background()
		iter := s.redisClient.Scan(ctx, 0, "alerts:*", 0).Iterator()
		for iter.Next(ctx) {
			s.redisClient.Del(ctx, iter.Val())
		}
	}

	return nil
}

// GetStats 获取统计
func (s *AlertService) GetStats() (map[string]interface{}, error) {
	ctx := context.Background()

	// 尝试缓存
	if s.redisClient != nil {
		cached, err := s.redisClient.Get(ctx, "alert:stats").Result()
		if err == nil {
			var stats map[string]interface{}
			if err := json.Unmarshal([]byte(cached), &stats); err == nil {
				return stats, nil
			}
		}
	}

	firing, _ := s.repo.GetCount("firing", "")
	critical, _ := s.repo.GetCount("firing", "critical")
	warning, _ := s.repo.GetCount("firing", "warning")

	stats := map[string]interface{}{
		"firing":   firing,
		"critical": critical,
		"warning":  warning,
	}

	// 缓存结果
	if s.redisClient != nil {
		statsJSON, _ := json.Marshal(stats)
		s.redisClient.Set(ctx, "alert:stats", statsJSON, 60*time.Second)
	}

	return stats, nil
}
