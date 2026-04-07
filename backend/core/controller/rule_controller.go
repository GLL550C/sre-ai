package controller

import (
	"core/model"
	"core/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RuleController handles alert rule HTTP requests
type RuleController struct {
	ruleService *service.RuleService
	logger      *zap.Logger
}

// NewRuleController creates a new rule controller
func NewRuleController(ruleService *service.RuleService, logger *zap.Logger) *RuleController {
	return &RuleController{
		ruleService: ruleService,
		logger:      logger,
	}
}

// GetRules handles GET /api/v1/rules
func (c *RuleController) GetRules(ctx *gin.Context) {
	rules, err := c.ruleService.GetRules()
	if err != nil {
		c.logger.Error("Failed to get rules", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": rules})
}

// GetRule handles GET /api/v1/rules/:id
func (c *RuleController) GetRule(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	rule, err := c.ruleService.GetRule(id)
	if err != nil {
		c.logger.Error("Failed to get rule", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": rule})
}

// CreateRule handles POST /api/v1/rules
func (c *RuleController) CreateRule(ctx *gin.Context) {
	var rule model.AlertRule
	if err := ctx.ShouldBindJSON(&rule); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ruleService.CreateRule(&rule); err != nil {
		c.logger.Error("Failed to create rule", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": rule})
}

// UpdateRule handles PUT /api/v1/rules/:id
func (c *RuleController) UpdateRule(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var rule model.AlertRule
	if err := ctx.ShouldBindJSON(&rule); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rule.ID = id
	if err := c.ruleService.UpdateRule(&rule); err != nil {
		c.logger.Error("Failed to update rule", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": rule})
}

// DeleteRule handles DELETE /api/v1/rules/:id
func (c *RuleController) DeleteRule(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.ruleService.DeleteRule(id); err != nil {
		c.logger.Error("Failed to delete rule", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Rule deleted"})
}
