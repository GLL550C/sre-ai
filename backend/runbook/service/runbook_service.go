package service

import (
	"context"
	"encoding/json"
	"fmt"
	"runbook/model"
	"runbook/repository"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RunbookService handles runbook business logic
type RunbookService struct {
	runbookRepo *repository.RunbookRepository
	redisClient *redis.Client
	logger      *zap.Logger
}

// NewRunbookService creates a new runbook service
func NewRunbookService(runbookRepo *repository.RunbookRepository, redisClient *redis.Client, logger *zap.Logger) *RunbookService {
	return &RunbookService{
		runbookRepo: runbookRepo,
		redisClient: redisClient,
		logger:      logger,
	}
}

// GetRunbooks retrieves runbooks with filters
func (s *RunbookService) GetRunbooks(alertName, severity string, page, pageSize int) ([]model.Runbook, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// Try cache first
	cacheKey := fmt.Sprintf("runbooks:%s:%s:%d:%d", alertName, severity, page, pageSize)
	ctx := context.Background()
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var result struct {
			Runbooks []model.Runbook `json:"runbooks"`
			Count    int             `json:"count"`
		}
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			return result.Runbooks, result.Count, nil
		}
	}

	runbooks, err := s.runbookRepo.GetRunbooks(alertName, severity, 1, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	// Count total (simplified, in real scenario would use count query)
	total := len(runbooks) // Simplified for demo

	// Cache result
	result := struct {
		Runbooks []model.Runbook `json:"runbooks"`
		Count    int             `json:"count"`
	}{runbooks, total}
	resultJSON, _ := json.Marshal(result)
	s.redisClient.Set(ctx, cacheKey, resultJSON, 5*time.Minute)

	return runbooks, total, nil
}

// GetRunbook retrieves a runbook by ID
func (s *RunbookService) GetRunbook(id int64) (*model.Runbook, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("runbook:%d", id)
	ctx := context.Background()
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var runbook model.Runbook
		if err := json.Unmarshal([]byte(cached), &runbook); err == nil {
			// Increment view count asynchronously
			go s.runbookRepo.IncrementViewCount(id)
			return &runbook, nil
		}
	}

	runbook, err := s.runbookRepo.GetRunbookByID(id)
	if err != nil {
		return nil, err
	}

	// Increment view count
	s.runbookRepo.IncrementViewCount(id)

	// Cache result
	runbookJSON, _ := json.Marshal(runbook)
	s.redisClient.Set(ctx, cacheKey, runbookJSON, 10*time.Minute)

	return runbook, nil
}

// CreateRunbook creates a new runbook
func (s *RunbookService) CreateRunbook(runbook *model.Runbook) error {
	runbook.Status = 1
	if err := s.runbookRepo.CreateRunbook(runbook); err != nil {
		return err
	}

	// Invalidate cache
	ctx := context.Background()
	iter := s.redisClient.Scan(ctx, 0, "runbooks:*", 0).Iterator()
	for iter.Next(ctx) {
		s.redisClient.Del(ctx, iter.Val())
	}

	return nil
}

// UpdateRunbook updates a runbook
func (s *RunbookService) UpdateRunbook(runbook *model.Runbook) error {
	if err := s.runbookRepo.UpdateRunbook(runbook); err != nil {
		return err
	}

	// Invalidate cache
	ctx := context.Background()
	cacheKey := fmt.Sprintf("runbook:%d", runbook.ID)
	s.redisClient.Del(ctx, cacheKey)

	iter := s.redisClient.Scan(ctx, 0, "runbooks:*", 0).Iterator()
	for iter.Next(ctx) {
		s.redisClient.Del(ctx, iter.Val())
	}

	return nil
}

// DeleteRunbook deletes a runbook
func (s *RunbookService) DeleteRunbook(id int64) error {
	if err := s.runbookRepo.DeleteRunbook(id); err != nil {
		return err
	}

	// Invalidate cache
	ctx := context.Background()
	cacheKey := fmt.Sprintf("runbook:%d", id)
	s.redisClient.Del(ctx, cacheKey)

	iter := s.redisClient.Scan(ctx, 0, "runbooks:*", 0).Iterator()
	for iter.Next(ctx) {
		s.redisClient.Del(ctx, iter.Val())
	}

	return nil
}

// SearchRunbooks searches runbooks by keyword
func (s *RunbookService) SearchRunbooks(keyword string, page, pageSize int) ([]model.Runbook, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.runbookRepo.SearchRunbooks(keyword, pageSize, offset)
}
