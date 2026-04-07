package service

import (
	"context"
	"encoding/json"
	"fmt"
	"tenant/model"
	"tenant/repository"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// TenantService handles tenant business logic
type TenantService struct {
	tenantRepo  *repository.TenantRepository
	redisClient *redis.Client
	logger      *zap.Logger
}

// NewTenantService creates a new tenant service
func NewTenantService(tenantRepo *repository.TenantRepository, redisClient *redis.Client, logger *zap.Logger) *TenantService {
	return &TenantService{
		tenantRepo:  tenantRepo,
		redisClient: redisClient,
		logger:      logger,
	}
}

// GetTenants retrieves all active tenants
func (s *TenantService) GetTenants() ([]model.Tenant, error) {
	// Try cache first
	ctx := context.Background()
	cached, err := s.redisClient.Get(ctx, "tenants:active").Result()
	if err == nil {
		var tenants []model.Tenant
		if err := json.Unmarshal([]byte(cached), &tenants); err == nil {
			return tenants, nil
		}
	}

	tenants, err := s.tenantRepo.GetTenants(1)
	if err != nil {
		return nil, err
	}

	// Cache result
	tenantsJSON, _ := json.Marshal(tenants)
	s.redisClient.Set(ctx, "tenants:active", tenantsJSON, 5*time.Minute)

	return tenants, nil
}

// GetTenant retrieves a tenant by ID
func (s *TenantService) GetTenant(id int64) (*model.Tenant, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("tenant:%d", id)
	ctx := context.Background()
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var tenant model.Tenant
		if err := json.Unmarshal([]byte(cached), &tenant); err == nil {
			return &tenant, nil
		}
	}

	tenant, err := s.tenantRepo.GetTenantByID(id)
	if err != nil {
		return nil, err
	}

	// Cache result
	tenantJSON, _ := json.Marshal(tenant)
	s.redisClient.Set(ctx, cacheKey, tenantJSON, 10*time.Minute)

	return tenant, nil
}

// CreateTenant creates a new tenant
func (s *TenantService) CreateTenant(tenant *model.Tenant) error {
	tenant.Status = 1
	if err := s.tenantRepo.CreateTenant(tenant); err != nil {
		return err
	}

	// Invalidate cache
	ctx := context.Background()
	s.redisClient.Del(ctx, "tenants:active")

	return nil
}

// UpdateTenant updates a tenant
func (s *TenantService) UpdateTenant(tenant *model.Tenant) error {
	if err := s.tenantRepo.UpdateTenant(tenant); err != nil {
		return err
	}

	// Invalidate cache
	ctx := context.Background()
	cacheKey := fmt.Sprintf("tenant:%d", tenant.ID)
	s.redisClient.Del(ctx, cacheKey, "tenants:active")

	return nil
}

// DeleteTenant deletes a tenant
func (s *TenantService) DeleteTenant(id int64) error {
	if err := s.tenantRepo.DeleteTenant(id); err != nil {
		return err
	}

	// Invalidate cache
	ctx := context.Background()
	cacheKey := fmt.Sprintf("tenant:%d", id)
	s.redisClient.Del(ctx, cacheKey, "tenants:active")

	return nil
}
