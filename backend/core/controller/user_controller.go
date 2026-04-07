package controller

import (
	"core/model"
	"core/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserController 用户控制器
type UserController struct {
	userService *service.UserService
	logger      *zap.Logger
}

// NewUserController 创建用户控制器
func NewUserController(userService *service.UserService, logger *zap.Logger) *UserController {
	return &UserController{
		userService: userService,
		logger:      logger,
	}
}

// ListUsers 获取用户列表
func (c *UserController) ListUsers(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	users, total, err := c.userService.ListUsers(page, pageSize)
	if err != nil {
		c.logger.Error("Failed to list users", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"items":     users,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetUser 获取用户详情
func (c *UserController) GetUser(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := c.userService.GetUser(id)
	if err != nil {
		c.logger.Error("Failed to get user", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

// CreateUser 创建用户
func (c *UserController) CreateUser(ctx *gin.Context) {
	var req model.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 获取当前用户
	createdBy := "system"
	if username, exists := ctx.Get("username"); exists {
		createdBy = username.(string)
	}

	user, err := c.userService.CreateUser(&req, createdBy)
	if err != nil {
		c.logger.Error("Failed to create user", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"data": user,
	})
}

// UpdateUser 更新用户
func (c *UserController) UpdateUser(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		Email  string `json:"email"`
		Phone  string `json:"phone"`
		Role   string `json:"role"`
		Status int    `json:"status"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := c.userService.UpdateUser(id, req.Email, req.Phone, req.Role, req.Status)
	if err != nil {
		c.logger.Error("Failed to update user", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

// DeleteUser 删除用户
func (c *UserController) DeleteUser(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := c.userService.DeleteUser(id); err != nil {
		c.logger.Error("Failed to delete user", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

// ResetPassword 重置密码
func (c *UserController) ResetPassword(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		NewPassword string `json:"new_password" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "New password is required"})
		return
	}

	if err := c.userService.ResetPassword(id, req.NewPassword); err != nil {
		c.logger.Error("Failed to reset password", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Password reset successfully",
	})
}
