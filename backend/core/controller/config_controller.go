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

// NewConfigController 创建配置控制器
func NewConfigController(service *service.ConfigService, logger *zap.Logger) *ConfigController {
	return &ConfigController{
		service: service,
		logger:  logger,
	}
}

// GetConfigTree 获取配置树
func (c *ConfigController) GetConfigTree(ctx *gin.Context) {
	tree, err := c.service.GetConfigTree()
	if err != nil {
		c.logger.Error("Failed to get config tree", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": tree})
}

// GetConfigItems 获取配置项列表
func (c *ConfigController) GetConfigItems(ctx *gin.Context) {
	category := ctx.Query("category")
	subCategory := ctx.Query("sub_category")

	items, err := c.service.GetConfigItemsByCategory(category, subCategory)
	if err != nil {
		c.logger.Error("Failed to get config items", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": items})
}

// GetConfigItem 获取单个配置项
func (c *ConfigController) GetConfigItem(ctx *gin.Context) {
	key := ctx.Param("key")

	item, err := c.service.GetConfigItem(key)
	if err != nil {
		c.logger.Error("Failed to get config item", zap.Error(err))
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Config not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": item})
}

// UpdateConfigValue 更新配置值
func (c *ConfigController) UpdateConfigValue(ctx *gin.Context) {
	key := ctx.Param("key")

	var req struct {
		Value string `json:"value" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := ctx.GetString("user")
	if user == "" {
		user = "anonymous"
	}

	if err := c.service.UpdateConfigValue(key, req.Value, user); err != nil {
		c.logger.Error("Failed to update config", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Config updated successfully"})
}

// UpdateMultipleConfigs 批量更新配置
func (c *ConfigController) UpdateMultipleConfigs(ctx *gin.Context) {
	var req struct {
		Configs map[string]string `json:"configs" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := ctx.GetString("user")
	if user == "" {
		user = "anonymous"
	}

	if err := c.service.UpdateMultipleConfigs(req.Configs, user); err != nil {
		c.logger.Error("Failed to update configs", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Configs updated successfully"})
}

// ResetConfigToDefault 重置配置为默认值
func (c *ConfigController) ResetConfigToDefault(ctx *gin.Context) {
	key := ctx.Param("key")

	if err := c.service.ResetConfigToDefault(key); err != nil {
		c.logger.Error("Failed to reset config", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Config reset to default"})
}

// GetSystemSettings 获取系统设置
func (c *ConfigController) GetSystemSettings(ctx *gin.Context) {
	settings := c.service.GetSystemSettings()
	ctx.JSON(http.StatusOK, gin.H{"data": settings})
}

// GetAISettings 获取AI设置
func (c *ConfigController) GetAISettings(ctx *gin.Context) {
	settings := c.service.GetAISettings()
	ctx.JSON(http.StatusOK, gin.H{"data": settings})
}

// GetNotificationSettings 获取通知设置
func (c *ConfigController) GetNotificationSettings(ctx *gin.Context) {
	settings := c.service.GetNotificationSettings()
	ctx.JSON(http.StatusOK, gin.H{"data": settings})
}

// ExportConfig 导出配置
func (c *ConfigController) ExportConfig(ctx *gin.Context) {
	items, err := c.service.ExportConfig()
	if err != nil {
		c.logger.Error("Failed to export config", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("Content-Type", "application/json")
	ctx.Header("Content-Disposition", "attachment; filename=config.json")
	ctx.JSON(http.StatusOK, items)
}

// ImportConfig 导入配置
func (c *ConfigController) ImportConfig(ctx *gin.Context) {
	var items []model.ConfigItem
	if err := ctx.ShouldBindJSON(&items); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := ctx.GetString("user")
	if user == "" {
		user = "anonymous"
	}

	if err := c.service.ImportConfig(items, user); err != nil {
		c.logger.Error("Failed to import config", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Config imported successfully"})
}

// GetAIConfig 获取AI配置
func (c *ConfigController) GetAIConfig(ctx *gin.Context) {
	config, err := c.service.GetAIConfig()
	if err != nil {
		c.logger.Error("Failed to get AI config", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": config})
}

// GetConfigs 获取所有配置(兼容旧API)
func (c *ConfigController) GetConfigs(ctx *gin.Context) {
	configs, err := c.service.GetAllConfigValues()
	if err != nil {
		c.logger.Error("Failed to get configs", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": configs})
}

// ReloadConfig 重新加载配置(兼容旧API)
func (c *ConfigController) ReloadConfig(ctx *gin.Context) {
	// 新配置系统不需要手动reload，配置实时生效
	ctx.JSON(http.StatusOK, gin.H{"message": "Config reloaded successfully"})
}
