package controller

import (
	"net/http"
	"strconv"
	"tenant/model"
	"tenant/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// TenantController handles tenant HTTP requests
type TenantController struct {
	tenantService *service.TenantService
	logger        *zap.Logger
}

// NewTenantController creates a new tenant controller
func NewTenantController(tenantService *service.TenantService, logger *zap.Logger) *TenantController {
	return &TenantController{
		tenantService: tenantService,
		logger:        logger,
	}
}

// GetTenants handles GET /api/v1/tenants
func (c *TenantController) GetTenants(ctx *gin.Context) {
	tenants, err := c.tenantService.GetTenants()
	if err != nil {
		c.logger.Error("Failed to get tenants", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": tenants})
}

// GetTenant handles GET /api/v1/tenants/:id
func (c *TenantController) GetTenant(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	tenant, err := c.tenantService.GetTenant(id)
	if err != nil {
		c.logger.Error("Failed to get tenant", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": tenant})
}

// CreateTenant handles POST /api/v1/tenants
func (c *TenantController) CreateTenant(ctx *gin.Context) {
	var tenant model.Tenant
	if err := ctx.ShouldBindJSON(&tenant); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.tenantService.CreateTenant(&tenant); err != nil {
		c.logger.Error("Failed to create tenant", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": tenant})
}

// UpdateTenant handles PUT /api/v1/tenants/:id
func (c *TenantController) UpdateTenant(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var tenant model.Tenant
	if err := ctx.ShouldBindJSON(&tenant); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenant.ID = id
	if err := c.tenantService.UpdateTenant(&tenant); err != nil {
		c.logger.Error("Failed to update tenant", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": tenant})
}

// DeleteTenant handles DELETE /api/v1/tenants/:id
func (c *TenantController) DeleteTenant(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.tenantService.DeleteTenant(id); err != nil {
		c.logger.Error("Failed to delete tenant", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tenant deleted"})
}
