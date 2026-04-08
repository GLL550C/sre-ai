package service

import (
	"context"
	"core/model"
	"core/repository"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// ClusterService Prometheus集群业务逻辑
type ClusterService struct {
	repo        *repository.ClusterRepository
	redisClient *redis.Client
	logger      *zap.Logger
}

// NewClusterService 创建服务
func NewClusterService(repo *repository.ClusterRepository, redisClient *redis.Client, logger *zap.Logger) *ClusterService {
	return &ClusterService{repo: repo, redisClient: redisClient, logger: logger}
}

// GetAll 获取所有集群
func (s *ClusterService) GetAll() ([]model.PrometheusCluster, error) {
	clusters, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	// 健康检查
	for i := range clusters {
		if s.checkHealth(&clusters[i]) {
			clusters[i].Status = 1
		} else {
			clusters[i].Status = 2
		}
	}
	return clusters, nil
}

// GetDefault 获取默认集群
func (s *ClusterService) GetDefault() (*model.PrometheusCluster, error) {
	// 先查缓存
	cacheKey := "cluster:default"
	ctx := context.Background()
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var c model.PrometheusCluster
		if err := json.Unmarshal([]byte(cached), &c); err == nil {
			if s.checkHealth(&c) {
				c.Status = 1
			} else {
				c.Status = 2
			}
			return &c, nil
		}
	}

	c, err := s.repo.GetDefault()
	if err != nil {
		return nil, err
	}

	if s.checkHealth(c) {
		c.Status = 1
	} else {
		c.Status = 2
	}

	// 缓存结果
	data, _ := json.Marshal(c)
	s.redisClient.Set(ctx, cacheKey, data, 5*time.Minute)
	return c, nil
}

// GetByID 根据ID获取
func (s *ClusterService) GetByID(id int64) (*model.PrometheusCluster, error) {
	// 先查缓存
	cacheKey := fmt.Sprintf("cluster:%d", id)
	ctx := context.Background()
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var c model.PrometheusCluster
		if err := json.Unmarshal([]byte(cached), &c); err == nil {
			return &c, nil
		}
	}

	c, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 缓存结果
	data, _ := json.Marshal(c)
	s.redisClient.Set(ctx, cacheKey, data, 5*time.Minute)
	return c, nil
}

// Create 创建集群
func (s *ClusterService) Create(c *model.PrometheusCluster) error {
	c.Status = 1
	if err := s.repo.Create(c); err != nil {
		return err
	}
	s.clearCache()
	return nil
}

// Update 更新集群
func (s *ClusterService) Update(c *model.PrometheusCluster) error {
	existing, err := s.repo.GetByID(c.ID)
	if err != nil {
		return err
	}

	// 保留未提供的字段
	if c.Name == "" {
		c.Name = existing.Name
	}
	if c.URL == "" {
		c.URL = existing.URL
	}
	if !c.IsDefault && existing.IsDefault {
		c.IsDefault = existing.IsDefault
	}
	if c.Status == 0 && existing.Status != 0 {
		c.Status = existing.Status
	}

	if err := s.repo.Update(c); err != nil {
		return err
	}

	s.clearCache()
	return nil
}

// Delete 删除集群
func (s *ClusterService) Delete(id int64) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	s.clearCache()
	return nil
}

// SetDefault 设置默认集群
func (s *ClusterService) SetDefault(id int64) error {
	if err := s.repo.ClearDefault(); err != nil {
		return err
	}
	if err := s.repo.SetDefault(id); err != nil {
		return err
	}
	ctx := context.Background()
	s.redisClient.Del(ctx, "cluster:default")
	return nil
}

// Test 测试连接
func (s *ClusterService) Test(url string) (bool, string) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url + "/api/v1/status/buildinfo")
	if err != nil {
		return false, err.Error()
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Sprintf("HTTP %d", resp.StatusCode)
	}
	return true, "连接成功"
}

// checkHealth 健康检查
func (s *ClusterService) checkHealth(c *model.PrometheusCluster) bool {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(c.URL + "/api/v1/status/buildinfo")
	if err != nil {
		s.logger.Debug("健康检查失败", zap.String("url", c.URL), zap.Error(err))
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// clearCache 清除缓存
func (s *ClusterService) clearCache() {
	ctx := context.Background()
	s.redisClient.Del(ctx, "cluster:default")
}
