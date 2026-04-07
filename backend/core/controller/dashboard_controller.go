package controller

import (
	"core/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// DashboardController handles dashboard HTTP requests
type DashboardController struct {
	dashboardService *service.DashboardService
	logger           *zap.Logger
}

// NewDashboardController creates a new dashboard controller
func NewDashboardController(dashboardService *service.DashboardService, logger *zap.Logger) *DashboardController {
	return &DashboardController{
		dashboardService: dashboardService,
		logger:           logger,
	}
}

// GetDashboard handles GET /api/v1/dashboard
func (c *DashboardController) GetDashboard(ctx *gin.Context) {
	var tenantID *int64
	if tid := ctx.Query("tenant_id"); tid != "" {
		if id, err := strconv.ParseInt(tid, 10, 64); err == nil {
			tenantID = &id
		}
	}

	dashboard, err := c.dashboardService.GetDashboard(tenantID)
	if err != nil {
		c.logger.Error("Failed to get dashboard", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": dashboard})
}

// GetMetrics handles GET /api/v1/dashboard/metrics
func (c *DashboardController) GetMetrics(ctx *gin.Context) {
	// Return dashboard metrics summary
	metrics := map[string]interface{}{
		"cpu_usage":    65.5,
		"memory_usage": 72.3,
		"disk_usage":   45.8,
		"network_io":   120.5,
	}

	ctx.JSON(http.StatusOK, gin.H{"data": metrics})
}
