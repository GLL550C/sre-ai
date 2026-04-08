package controller

import (
	"core/model"
	"core/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AlertController 告警控制器
type AlertController struct {
	service *service.AlertService
	logger  *zap.Logger
}

// NewAlertController 创建控制器
func NewAlertController(service *service.AlertService, logger *zap.Logger) *AlertController {
	return &AlertController{service: service, logger: logger}
}

// GetAlerts 获取告警列表
func (c *AlertController) GetAlerts(ctx *gin.Context) {
	status := ctx.Query("status")
	severity := ctx.Query("severity")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	alerts, count, err := c.service.GetAlerts(status, severity, page, pageSize)
	if err != nil {
		c.logger.Error("获取告警失败", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  alerts,
		"total": count,
		"page":  page,
		"size":  pageSize,
	})
}

// GetAlert 获取单个告警
func (c *AlertController) GetAlert(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	alert, err := c.service.GetByID(id)
	if err != nil {
		c.logger.Error("获取告警失败", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": alert})
}

// CreateAlert 创建告警
func (c *AlertController) CreateAlert(ctx *gin.Context) {
	var alert model.Alert
	if err := ctx.ShouldBindJSON(&alert); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.Create(&alert); err != nil {
		c.logger.Error("创建告警失败", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": alert})
}

// ReceiveWebhook 接收Webhook
func (c *AlertController) ReceiveWebhook(ctx *gin.Context) {
	var payload model.WebhookPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.ProcessWebhook(&payload); err != nil {
		c.logger.Error("处理Webhook失败", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Webhook已处理"})
}

// AcknowledgeAlert 确认告警
func (c *AlertController) AcknowledgeAlert(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	user := ctx.GetString("username")
	if user == "" {
		user = "unknown"
	}

	if err := c.service.Acknowledge(id, user); err != nil {
		c.logger.Error("确认告警失败", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "告警已确认"})
}
