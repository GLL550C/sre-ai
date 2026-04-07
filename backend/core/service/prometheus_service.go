package service

import (
	"context"
	"core/model"
	"core/repository"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// PrometheusService handles Prometheus interaction
type PrometheusService struct {
	clusterRepo *repository.ClusterRepository
	redisClient *redis.Client
	logger      *zap.Logger
}

// NewPrometheusService creates a new Prometheus service
func NewPrometheusService(clusterRepo *repository.ClusterRepository, redisClient *redis.Client, logger *zap.Logger) *PrometheusService {
	return &PrometheusService{
		clusterRepo: clusterRepo,
		redisClient: redisClient,
		logger:      logger,
	}
}

// Query executes a Prometheus query
func (s *PrometheusService) Query(clusterID int64, query string) (*model.PrometheusQueryResult, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("prometheus:query:%d:%s", clusterID, query)
	ctx := context.Background()
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var result model.PrometheusQueryResult
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			return &result, nil
		}
	}

	// Get cluster URL
	cluster, err := s.clusterRepo.GetClusterByID(clusterID)
	if err != nil {
		return nil, err
	}

	// Execute query
	url := fmt.Sprintf("%s/api/v1/query?query=%s", cluster.URL, query)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result model.PrometheusQueryResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	// Cache result
	resultJSON, _ := json.Marshal(result)
	s.redisClient.Set(ctx, cacheKey, resultJSON, 30*time.Second)

	return &result, nil
}

// QueryRange executes a Prometheus range query
func (s *PrometheusService) QueryRange(clusterID int64, query, start, end, step string) (*model.PrometheusQueryResult, error) {
	// Get cluster URL
	cluster, err := s.clusterRepo.GetClusterByID(clusterID)
	if err != nil {
		return nil, err
	}

	// Execute query
	url := fmt.Sprintf("%s/api/v1/query_range?query=%s&start=%s&end=%s&step=%s",
		cluster.URL, query, start, end, step)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result model.PrometheusQueryResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetMetrics retrieves available metrics from Prometheus
func (s *PrometheusService) GetMetrics(clusterID int64) ([]string, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("prometheus:metrics:%d", clusterID)
	ctx := context.Background()
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var metrics []string
		if err := json.Unmarshal([]byte(cached), &metrics); err == nil {
			return metrics, nil
		}
	}

	// Get cluster URL
	cluster, err := s.clusterRepo.GetClusterByID(clusterID)
	if err != nil {
		return nil, err
	}

	// Get metrics
	url := fmt.Sprintf("%s/api/v1/label/__name__/values", cluster.URL)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status string   `json:"status"`
		Data   []string `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	// Cache result
	metricsJSON, _ := json.Marshal(result.Data)
	s.redisClient.Set(ctx, cacheKey, metricsJSON, 5*time.Minute)

	return result.Data, nil
}
