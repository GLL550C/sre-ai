package controller

import (
	"core/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// PrometheusController handles Prometheus proxy HTTP requests
type PrometheusController struct {
	prometheusService *service.PrometheusService
	logger            *zap.Logger
}

// NewPrometheusController creates a new Prometheus controller
func NewPrometheusController(prometheusService *service.PrometheusService, logger *zap.Logger) *PrometheusController {
	return &PrometheusController{
		prometheusService: prometheusService,
		logger:            logger,
	}
}

// Query handles GET /api/v1/prometheus/query
func (c *PrometheusController) Query(ctx *gin.Context) {
	clusterID, _ := strconv.ParseInt(ctx.DefaultQuery("cluster_id", "1"), 10, 64)
	query := ctx.Query("query")

	if query == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter is required"})
		return
	}

	result, err := c.prometheusService.Query(clusterID, query)
	if err != nil {
		c.logger.Error("Failed to query Prometheus", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// QueryRange handles GET /api/v1/prometheus/query_range
func (c *PrometheusController) QueryRange(ctx *gin.Context) {
	clusterID, _ := strconv.ParseInt(ctx.DefaultQuery("cluster_id", "1"), 10, 64)
	query := ctx.Query("query")
	start := ctx.Query("start")
	end := ctx.Query("end")
	step := ctx.DefaultQuery("step", "15s")

	if query == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter is required"})
		return
	}

	result, err := c.prometheusService.QueryRange(clusterID, query, start, end, step)
	if err != nil {
		c.logger.Error("Failed to query Prometheus range", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}
