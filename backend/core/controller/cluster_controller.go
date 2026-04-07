package controller

import (
	"core/model"
	"core/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ClusterController handles Prometheus cluster HTTP requests
type ClusterController struct {
	clusterService *service.ClusterService
	logger         *zap.Logger
}

// NewClusterController creates a new cluster controller
func NewClusterController(clusterService *service.ClusterService, logger *zap.Logger) *ClusterController {
	return &ClusterController{
		clusterService: clusterService,
		logger:         logger,
	}
}

// GetClusters handles GET /api/v1/clusters
func (c *ClusterController) GetClusters(ctx *gin.Context) {
	clusters, err := c.clusterService.GetClusters()
	if err != nil {
		c.logger.Error("Failed to get clusters", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": clusters})
}

// GetCluster handles GET /api/v1/clusters/:id
func (c *ClusterController) GetCluster(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	cluster, err := c.clusterService.GetCluster(id)
	if err != nil {
		c.logger.Error("Failed to get cluster", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": cluster})
}

// CreateCluster handles POST /api/v1/clusters
func (c *ClusterController) CreateCluster(ctx *gin.Context) {
	var cluster model.PrometheusCluster
	if err := ctx.ShouldBindJSON(&cluster); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.clusterService.CreateCluster(&cluster); err != nil {
		c.logger.Error("Failed to create cluster", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": cluster})
}

// UpdateCluster handles PUT /api/v1/clusters/:id
func (c *ClusterController) UpdateCluster(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var cluster model.PrometheusCluster
	if err := ctx.ShouldBindJSON(&cluster); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cluster.ID = id
	if err := c.clusterService.UpdateCluster(&cluster); err != nil {
		c.logger.Error("Failed to update cluster", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": cluster})
}

// TestCluster handles GET /api/v1/clusters/:id/test
func (c *ClusterController) TestCluster(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	cluster, err := c.clusterService.GetCluster(id)
	if err != nil {
		c.logger.Error("Failed to get cluster", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	success, message := c.clusterService.TestCluster(cluster.URL)
	ctx.JSON(http.StatusOK, gin.H{"success": success, "message": message})
}

// SetDefaultCluster handles POST /api/v1/clusters/:id/default
func (c *ClusterController) SetDefaultCluster(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.clusterService.SetDefaultCluster(id); err != nil {
		c.logger.Error("Failed to set default cluster", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "已设为默认集群"})
}

// GetDefaultCluster handles GET /api/v1/clusters/default
func (c *ClusterController) GetDefaultCluster(ctx *gin.Context) {
	cluster, err := c.clusterService.GetDefaultCluster()
	if err != nil {
		c.logger.Error("Failed to get default cluster", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": cluster})
}

// DeleteCluster handles DELETE /api/v1/clusters/:id
func (c *ClusterController) DeleteCluster(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.clusterService.DeleteCluster(id); err != nil {
		c.logger.Error("Failed to delete cluster", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Cluster deleted"})
}
