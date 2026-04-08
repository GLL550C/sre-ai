package controller

import (
	"core/model"
	"core/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RuleController 告警规则控制器
type RuleController struct {
	service *service.RuleService
	logger  *zap.Logger
}

// NewRuleController 创建控制器
func NewRuleController(service *service.RuleService, logger *zap.Logger) *RuleController {
	return &RuleController{service: service, logger: logger}
}

// GetRules 获取规则列表
func (c *RuleController) GetRules(ctx *gin.Context) {
	rules, err := c.service.GetAll()
	if err != nil {
		c.logger.Error("获取规则失败", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": rules})
}

// GetRule 获取单个规则
func (c *RuleController) GetRule(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	rule, err := c.service.GetByID(id)
	if err != nil {
		c.logger.Error("获取规则失败", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": rule})
}

// CreateRule 创建规则
func (c *RuleController) CreateRule(ctx *gin.Context) {
	var rule model.AlertRule
	if err := ctx.ShouldBindJSON(&rule); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.Create(&rule); err != nil {
		c.logger.Error("创建规则失败", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"data": rule})
}

// UpdateRule 更新规则
func (c *RuleController) UpdateRule(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var rule model.AlertRule
	if err := ctx.ShouldBindJSON(&rule); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rule.ID = id
	if err := c.service.Update(&rule); err != nil {
		c.logger.Error("更新规则失败", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": rule})
}

// DeleteRule 删除规则
func (c *RuleController) DeleteRule(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	if err := c.service.Delete(id); err != nil {
		c.logger.Error("删除规则失败", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "规则已删除"})
}
