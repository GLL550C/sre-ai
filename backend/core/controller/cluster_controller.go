package controller

import (
	"core/model"
	"core/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ClusterController Prometheus集群控制器
type ClusterController struct {
	service *service.ClusterService
	logger  *zap.Logger
}

// NewClusterController 创建控制器
func NewClusterController(service *service.ClusterService, logger *zap.Logger) *ClusterController {
	return &ClusterController{service: service, logger: logger}
}

// GetClusters 获取所有集群
func (c *ClusterController) GetClusters(ctx *gin.Context) {
	clusters, err := c.service.GetAll()
	if err != nil {
		c.logger.Error("获取集群失败", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": clusters})
}

// GetCluster 获取单个集群
func (c *ClusterController) GetCluster(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	cluster, err := c.service.GetByID(id)
	if err != nil {
		c.logger.Error("获取集群失败", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": cluster})
}

// GetDefaultCluster 获取默认集群
func (c *ClusterController) GetDefaultCluster(ctx *gin.Context) {
	cluster, err := c.service.GetDefault()
	if err != nil {
		c.logger.Error("获取默认集群失败", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": cluster})
}

// CreateCluster 创建集群
func (c *ClusterController) CreateCluster(ctx *gin.Context) {
	var cluster model.PrometheusCluster
	if err := ctx.ShouldBindJSON(&cluster); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.Create(&cluster); err != nil {
		c.logger.Error("创建集群失败", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"data": cluster})
}

// UpdateCluster 更新集群
func (c *ClusterController) UpdateCluster(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var cluster model.PrometheusCluster
	if err := ctx.ShouldBindJSON(&cluster); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cluster.ID = id
	if err := c.service.Update(&cluster); err != nil {
		c.logger.Error("更新集群失败", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": cluster})
}

// DeleteCluster 删除集群
func (c *ClusterController) DeleteCluster(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	if err := c.service.Delete(id); err != nil {
		c.logger.Error("删除集群失败", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "集群已删除"})
}

// TestCluster 测试集群连接
func (c *ClusterController) TestCluster(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	cluster, err := c.service.GetByID(id)
	if err != nil {
		c.logger.Error("获取集群失败", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	success, message := c.service.Test(cluster.URL)
	ctx.JSON(http.StatusOK, gin.H{"success": success, "message": message})
}

// SetDefaultCluster 设置默认集群
func (c *ClusterController) SetDefaultCluster(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	if err := c.service.SetDefault(id); err != nil {
		c.logger.Error("设置默认集群失败", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "已设为默认集群"})
}
