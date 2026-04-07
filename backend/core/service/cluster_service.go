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

// ClusterService handles Prometheus cluster business logic
type ClusterService struct {
	clusterRepo *repository.ClusterRepository
	redisClient *redis.Client
	logger      *zap.Logger
}

// NewClusterService creates a new cluster service
func NewClusterService(clusterRepo *repository.ClusterRepository, redisClient *redis.Client, logger *zap.Logger) *ClusterService {
	return &ClusterService{
		clusterRepo: clusterRepo,
		redisClient: redisClient,
		logger:      logger,
	}
}

// GetClusters retrieves all active clusters with health check
func (s *ClusterService) GetClusters() ([]model.PrometheusCluster, error) {
	clusters, err := s.clusterRepo.GetActiveClusters()
	if err != nil {
		return nil, err
	}

	// Perform health check for each cluster
	for i := range clusters {
		healthy := s.checkClusterHealth(&clusters[i])
		if healthy {
			clusters[i].Status = 1
		} else {
			clusters[i].Status = 2 // 2 = unhealthy but still active in DB
		}
	}

	return clusters, nil
}

// GetDefaultCluster retrieves the default cluster
func (s *ClusterService) GetDefaultCluster() (*model.PrometheusCluster, error) {
	// Try cache first
	cacheKey := "cluster:default"
	ctx := context.Background()
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var cluster model.PrometheusCluster
		if err := json.Unmarshal([]byte(cached), &cluster); err == nil {
			return &cluster, nil
		}
	}

	cluster, err := s.clusterRepo.GetDefaultCluster()
	if err != nil {
		return nil, err
	}

	// Cache result
	clusterJSON, _ := json.Marshal(cluster)
	s.redisClient.Set(ctx, cacheKey, clusterJSON, 5*time.Minute)

	return cluster, nil
}

// TestCluster tests cluster connectivity
func (s *ClusterService) TestCluster(url string) (bool, string) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

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

// SetDefaultCluster sets a cluster as default
func (s *ClusterService) SetDefaultCluster(id int64) error {
	// Clear existing default
	if err := s.clusterRepo.ClearDefaultCluster(); err != nil {
		return err
	}

	// Set new default
	if err := s.clusterRepo.SetDefaultCluster(id); err != nil {
		return err
	}

	// Invalidate cache
	ctx := context.Background()
	s.redisClient.Del(ctx, "cluster:default")

	return nil
}

// checkClusterHealth checks if Prometheus cluster is reachable
func (s *ClusterService) checkClusterHealth(cluster *model.PrometheusCluster) bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Try to access Prometheus API
	resp, err := client.Get(cluster.URL + "/api/v1/status/buildinfo")
	if err != nil {
		s.logger.Debug("Cluster health check failed",
			zap.String("url", cluster.URL),
			zap.Error(err))
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// GetCluster retrieves a cluster by ID
func (s *ClusterService) GetCluster(id int64) (*model.PrometheusCluster, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("cluster:%d", id)
	ctx := context.Background()
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var cluster model.PrometheusCluster
		if err := json.Unmarshal([]byte(cached), &cluster); err == nil {
			return &cluster, nil
		}
	}

	cluster, err := s.clusterRepo.GetClusterByID(id)
	if err != nil {
		return nil, err
	}

	// Cache result
	clusterJSON, _ := json.Marshal(cluster)
	s.redisClient.Set(ctx, cacheKey, clusterJSON, 5*time.Minute)

	return cluster, nil
}

// CreateCluster creates a new Prometheus cluster
func (s *ClusterService) CreateCluster(cluster *model.PrometheusCluster) error {
	cluster.Status = 1
	if err := s.clusterRepo.CreateCluster(cluster); err != nil {
		return err
	}

	// Invalidate cache
	ctx := context.Background()
	s.redisClient.Del(ctx, "clusters:active")

	return nil
}

// UpdateCluster updates a Prometheus cluster
func (s *ClusterService) UpdateCluster(cluster *model.PrometheusCluster) error {
	// Get existing cluster to preserve status if not provided
	existing, err := s.clusterRepo.GetClusterByID(cluster.ID)
	if err != nil {
		return err
	}

	// Preserve status if not explicitly provided (status 0 means inactive, so we check if it was explicitly set)
	// If cluster.Status is 0 and existing.Status is not 0, keep existing status
	if cluster.Status == 0 && existing.Status != 0 {
		cluster.Status = existing.Status
	}

	if err := s.clusterRepo.UpdateCluster(cluster); err != nil {
		return err
	}

	// Invalidate cache
	ctx := context.Background()
	cacheKey := fmt.Sprintf("cluster:%d", cluster.ID)
	s.redisClient.Del(ctx, cacheKey, "clusters:active")

	return nil
}

// DeleteCluster deletes a Prometheus cluster
func (s *ClusterService) DeleteCluster(id int64) error {
	if err := s.clusterRepo.DeleteCluster(id); err != nil {
		return err
	}

	// Invalidate cache
	ctx := context.Background()
	cacheKey := fmt.Sprintf("cluster:%d", id)
	s.redisClient.Del(ctx, cacheKey, "clusters:active")

	return nil
}
