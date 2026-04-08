package controller

import (
	"core/model"
	"core/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AIModelController handles AI model config HTTP requests
type AIModelController struct {
	service *service.AIModelService
	logger  *zap.Logger
}

// NewAIModelController creates a new AI model controller
func NewAIModelController(service *service.AIModelService, logger *zap.Logger) *AIModelController {
	return &AIModelController{
		service: service,
		logger:  logger,
	}
}

// GetConfigs handles GET /api/v1/ai/configs
func (c *AIModelController) GetConfigs(ctx *gin.Context) {
	configs, err := c.service.GetAll()
	if err != nil {
		c.logger.Error("Failed to get AI configs", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": configs})
}

// GetConfig handles GET /api/v1/ai/configs/:id
func (c *AIModelController) GetConfig(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	config, err := c.service.GetByID(id)
	if err != nil {
		c.logger.Error("Failed to get AI config", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": config})
}

// CreateConfig handles POST /api/v1/ai/configs
func (c *AIModelController) CreateConfig(ctx *gin.Context) {
	var config model.AIModel
	if err := ctx.ShouldBindJSON(&config); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := ctx.GetString("user")
	if user == "" {
		user = "anonymous"
	}

	if err := c.service.Create(&config, user); err != nil {
		c.logger.Error("Failed to create AI config", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": config})
}

// UpdateConfig handles PUT /api/v1/ai/configs/:id
func (c *AIModelController) UpdateConfig(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var config model.AIModel
	if err := ctx.ShouldBindJSON(&config); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.ID = id

	user := ctx.GetString("user")
	if user == "" {
		user = "anonymous"
	}

	if err := c.service.Update(&config, user); err != nil {
		c.logger.Error("Failed to update AI config", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": config})
}

// DeleteConfig handles DELETE /api/v1/ai/configs/:id
func (c *AIModelController) DeleteConfig(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.service.Delete(id); err != nil {
		c.logger.Error("Failed to delete AI config", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Config deleted successfully"})
}

// TestConfig handles POST /api/v1/ai/configs/:id/test
func (c *AIModelController) TestConfig(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	success, message, err := c.service.Test(id)
	if err != nil {
		c.logger.Error("Failed to test AI config", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if success {
		ctx.JSON(http.StatusOK, gin.H{"success": true, "message": message})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "message": message})
	}
}

// SetDefaultConfig handles PUT /api/v1/ai/configs/:id/default
func (c *AIModelController) SetDefaultConfig(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.service.SetDefault(id); err != nil {
		c.logger.Error("Failed to set default AI config", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Config set as default"})
}

// GetActiveConfig handles GET /api/v1/ai/configs/active
func (c *AIModelController) GetActiveConfig(ctx *gin.Context) {
	config, err := c.service.GetActiveConfig()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No active config found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": config})
}
