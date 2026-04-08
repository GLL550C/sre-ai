package controller

import (
	"core/model"
	"core/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ConfigController 配置控制器
type ConfigController struct {
	service *service.ConfigService
	logger  *zap.Logger
}

// NewConfigController 创建控制器
func NewConfigController(service *service.ConfigService, logger *zap.Logger) *ConfigController {
	return &ConfigController{service: service, logger: logger}
}

// GetByCategory 按分类获取配置
func (c *ConfigController) GetByCategory(ctx *gin.Context) {
	category := ctx.Query("category")
	if category == "" {
		category = "platform"
	}

	configs, err := c.service.GetByCategory(category)
	if err != nil {
		c.logger.Error("获取配置失败", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": configs})
}

// GetByKey 根据key获取配置
func (c *ConfigController) GetByKey(ctx *gin.Context) {
	key := ctx.Param("key")
	config, err := c.service.GetByKey(key)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "配置不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": config})
}

// Update 更新配置
func (c *ConfigController) Update(ctx *gin.Context) {
	key := ctx.Param("key")
	var req struct {
		Value string `json:"value" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := ctx.GetString("username")
	if user == "" {
		user = "system"
	}

	if err := c.service.Update(key, req.Value, user); err != nil {
		c.logger.Error("更新配置失败", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "配置已更新"})
}

// GetConfigItems 获取配置项列表(支持category和sub_category过滤)
func (c *ConfigController) GetConfigItems(ctx *gin.Context) {
	category := ctx.Query("category")
	subCategory := ctx.Query("sub_category")

	var configs []model.SystemConfig
	var err error

	if category != "" && subCategory != "" {
		configs, err = c.service.GetByCategoryAndSubCategory(category, subCategory)
	} else if category != "" {
		configs, err = c.service.GetByCategory(category)
	} else {
		configs, err = c.service.GetAll()
	}

	if err != nil {
		c.logger.Error("获取配置失败", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": configs})
}

// GetConfigItemByKey 根据key获取配置(用于 /config/items/:key 路由)
func (c *ConfigController) GetConfigItemByKey(ctx *gin.Context) {
	key := ctx.Param("key")
	config, err := c.service.GetByKey(key)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "配置不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": config})
}

// GetAppName 获取应用名称
func (c *ConfigController) GetAppName(ctx *gin.Context) {
	name := c.service.GetAppName()
	ctx.JSON(http.StatusOK, gin.H{"data": name})
}
