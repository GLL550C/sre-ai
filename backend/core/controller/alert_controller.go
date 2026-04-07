package controller

import (
	"core/model"
	"core/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AlertController handles alert HTTP requests
type AlertController struct {
	alertService *service.AlertService
	logger       *zap.Logger
}

// NewAlertController creates a new alert controller
func NewAlertController(alertService *service.AlertService, logger *zap.Logger) *AlertController {
	return &AlertController{
		alertService: alertService,
		logger:       logger,
	}
}

// GetAlerts handles GET /api/v1/alerts
func (c *AlertController) GetAlerts(ctx *gin.Context) {
	status := ctx.Query("status")
	severity := ctx.Query("severity")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	alerts, count, err := c.alertService.GetAlerts(status, severity, page, pageSize)
	if err != nil {
		c.logger.Error("Failed to get alerts", zap.Error(err))
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

// GetAlert handles GET /api/v1/alerts/:id
func (c *AlertController) GetAlert(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	alert, err := c.alertService.GetAlert(id)
	if err != nil {
		c.logger.Error("Failed to get alert", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": alert})
}

// CreateAlert handles POST /api/v1/alerts
func (c *AlertController) CreateAlert(ctx *gin.Context) {
	var alert model.Alert
	if err := ctx.ShouldBindJSON(&alert); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.alertService.CreateAlert(&alert); err != nil {
		c.logger.Error("Failed to create alert", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": alert})
}

// ReceiveWebhook handles POST /api/v1/alerts/webhook
func (c *AlertController) ReceiveWebhook(ctx *gin.Context) {
	var payload model.WebhookPayload

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.alertService.ProcessWebhook(&payload); err != nil {
		c.logger.Error("Failed to process webhook", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Webhook processed"})
}

// AcknowledgeAlert handles PUT /api/v1/alerts/:id/ack
func (c *AlertController) AcknowledgeAlert(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req struct {
		User string `json:"user"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		req.User = "unknown"
	}

	if err := c.alertService.AcknowledgeAlert(id, req.User); err != nil {
		c.logger.Error("Failed to acknowledge alert", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Alert acknowledged"})
}
