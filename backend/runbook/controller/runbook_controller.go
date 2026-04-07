package controller

import (
	"net/http"
	"runbook/model"
	"runbook/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RunbookController handles runbook HTTP requests
type RunbookController struct {
	runbookService *service.RunbookService
	logger         *zap.Logger
}

// NewRunbookController creates a new runbook controller
func NewRunbookController(runbookService *service.RunbookService, logger *zap.Logger) *RunbookController {
	return &RunbookController{
		runbookService: runbookService,
		logger:         logger,
	}
}

// GetRunbooks handles GET /api/v1/runbooks
func (c *RunbookController) GetRunbooks(ctx *gin.Context) {
	alertName := ctx.Query("alert_name")
	severity := ctx.Query("severity")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	runbooks, total, err := c.runbookService.GetRunbooks(alertName, severity, page, pageSize)
	if err != nil {
		c.logger.Error("Failed to get runbooks", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  runbooks,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// GetRunbook handles GET /api/v1/runbooks/:id
func (c *RunbookController) GetRunbook(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	runbook, err := c.runbookService.GetRunbook(id)
	if err != nil {
		c.logger.Error("Failed to get runbook", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": runbook})
}

// CreateRunbook handles POST /api/v1/runbooks
func (c *RunbookController) CreateRunbook(ctx *gin.Context) {
	var runbook model.Runbook
	if err := ctx.ShouldBindJSON(&runbook); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.runbookService.CreateRunbook(&runbook); err != nil {
		c.logger.Error("Failed to create runbook", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": runbook})
}

// UpdateRunbook handles PUT /api/v1/runbooks/:id
func (c *RunbookController) UpdateRunbook(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var runbook model.Runbook
	if err := ctx.ShouldBindJSON(&runbook); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	runbook.ID = id
	if err := c.runbookService.UpdateRunbook(&runbook); err != nil {
		c.logger.Error("Failed to update runbook", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": runbook})
}

// DeleteRunbook handles DELETE /api/v1/runbooks/:id
func (c *RunbookController) DeleteRunbook(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.runbookService.DeleteRunbook(id); err != nil {
		c.logger.Error("Failed to delete runbook", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Runbook deleted"})
}

// SearchRunbooks handles GET /api/v1/runbooks/search
func (c *RunbookController) SearchRunbooks(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	if keyword == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Keyword is required"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	runbooks, total, err := c.runbookService.SearchRunbooks(keyword, page, pageSize)
	if err != nil {
		c.logger.Error("Failed to search runbooks", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  runbooks,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}
