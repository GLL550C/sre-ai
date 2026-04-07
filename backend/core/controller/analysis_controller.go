package controller

import (
	"context"
	"core/ai"
	"core/model"
	"core/service"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AnalysisController handles AI analysis HTTP requests
type AnalysisController struct {
	analysisService *service.AnalysisService
	logger          *zap.Logger
}

// NewAnalysisController creates a new analysis controller
func NewAnalysisController(analysisService *service.AnalysisService, logger *zap.Logger) *AnalysisController {
	return &AnalysisController{
		analysisService: analysisService,
		logger:          logger,
	}
}

// GetAnalysis handles GET /api/v1/analysis
func (c *AnalysisController) GetAnalysis(ctx *gin.Context) {
	clusterID, _ := strconv.ParseInt(ctx.Query("cluster_id"), 10, 64)
	analysisType := ctx.Query("type")
	status, _ := strconv.Atoi(ctx.DefaultQuery("status", "1"))
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	analyses, err := c.analysisService.GetAnalysis(clusterID, analysisType, status, page, pageSize)
	if err != nil {
		c.logger.Error("Failed to get analysis", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": analyses})
}

// GetAnalysisByID handles GET /api/v1/analysis/:id
func (c *AnalysisController) GetAnalysisByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	analysis, err := c.analysisService.GetAnalysisByID(id)
	if err != nil {
		c.logger.Error("Failed to get analysis", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": analysis})
}

// CreateAnalysis handles POST /api/v1/analysis
func (c *AnalysisController) CreateAnalysis(ctx *gin.Context) {
	var analysis model.AIAnalysis
	if err := ctx.ShouldBindJSON(&analysis); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set current user as creator
	analysis.CreatedBy = ctx.GetString("user") // Assuming auth middleware sets this
	if analysis.CreatedBy == "" {
		analysis.CreatedBy = "anonymous"
	}

	if err := c.analysisService.CreateAnalysis(&analysis); err != nil {
		c.logger.Error("Failed to create analysis", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": analysis})
}

// DeleteAnalysis handles DELETE /api/v1/analysis/:id
func (c *AnalysisController) DeleteAnalysis(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.analysisService.DeleteAnalysis(id); err != nil {
		c.logger.Error("Failed to delete analysis", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Analysis deleted successfully"})
}

// ArchiveAnalysis handles PUT /api/v1/analysis/:id/archive
func (c *AnalysisController) ArchiveAnalysis(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.analysisService.ArchiveAnalysis(id); err != nil {
		c.logger.Error("Failed to archive analysis", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Analysis archived successfully"})
}

// GetAnalysisStats handles GET /api/v1/analysis/stats
func (c *AnalysisController) GetAnalysisStats(ctx *gin.Context) {
	clusterID, _ := strconv.ParseInt(ctx.Query("cluster_id"), 10, 64)

	stats, err := c.analysisService.GetAnalysisStats(clusterID)
	if err != nil {
		c.logger.Error("Failed to get analysis stats", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": stats})
}

// CompareClusters handles POST /api/v1/analysis/compare
func (c *AnalysisController) CompareClusters(ctx *gin.Context) {
	var req struct {
		ClusterIDs []int64 `json:"cluster_ids" binding:"required,min=2"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	report, err := c.analysisService.CompareClusters(req.ClusterIDs)
	if err != nil {
		c.logger.Error("Failed to compare clusters", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": report})
}

// Chat handles POST /api/v1/ai/chat
func (c *AnalysisController) Chat(ctx *gin.Context) {
	var req struct {
		Messages []ai.Message `json:"messages" binding:"required,min=1"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.analysisService.Chat(context.Background(), req.Messages)
	if err != nil {
		c.logger.Error("AI chat failed", zap.Error(err))
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": resp})
}

// ChatStream handles POST /api/v1/ai/chat/stream
func (c *AnalysisController) ChatStream(ctx *gin.Context) {
	var req struct {
		Messages []ai.Message `json:"messages" binding:"required,min=1"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set SSE headers
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")

	stream, err := c.analysisService.ChatStream(context.Background(), req.Messages)
	if err != nil {
		c.logger.Error("AI chat stream failed", zap.Error(err))
		ctx.SSEvent("error", gin.H{"error": err.Error()})
		return
	}

	// Stream responses
	ctx.Stream(func(w io.Writer) bool {
		select {
		case resp, ok := <-stream:
			if !ok {
				return false
			}
			if resp.Error != "" {
				ctx.SSEvent("error", gin.H{"error": resp.Error})
				return false
			}
			if resp.Done {
				ctx.SSEvent("done", gin.H{"done": true})
				return false
			}
			ctx.SSEvent("message", gin.H{"content": resp.Content})
			return true
		case <-time.After(30 * time.Second):
			ctx.SSEvent("error", gin.H{"error": "timeout"})
			return false
		}
	})
}

// AIHealth handles GET /api/v1/ai/health
func (c *AnalysisController) AIHealth(ctx *gin.Context) {
	if err := c.analysisService.AIHealth(); err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

// AIModelInfo handles GET /api/v1/ai/model
func (c *AnalysisController) AIModelInfo(ctx *gin.Context) {
	info := c.analysisService.AIModelInfo()
	ctx.JSON(http.StatusOK, gin.H{"data": info})
}
