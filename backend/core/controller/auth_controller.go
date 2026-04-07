package controller

import (
	"core/model"
	"core/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthController 认证控制器
type AuthController struct {
	authService *service.AuthService
	userService *service.UserService
	logger      *zap.Logger
}

// NewAuthController 创建认证控制器
func NewAuthController(authService *service.AuthService, userService *service.UserService, logger *zap.Logger) *AuthController {
	return &AuthController{
		authService: authService,
		userService: userService,
		logger:      logger,
	}
}

// GetCaptcha 获取验证码
func (c *AuthController) GetCaptcha(ctx *gin.Context) {
	captcha, err := c.authService.GenerateCaptcha()
	if err != nil {
		c.logger.Error("Failed to generate captcha", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate captcha"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": captcha,
	})
}

// Login 用户登录
func (c *AuthController) Login(ctx *gin.Context) {
	var req model.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 获取客户端IP
	clientIP := ctx.ClientIP()

	resp, err := c.authService.Login(&req, clientIP)
	if err != nil {
		c.logger.Warn("Login failed", zap.String("username", req.Username), zap.Error(err))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": resp,
	})
}

// Logout 用户登出
func (c *AuthController) Logout(ctx *gin.Context) {
	// JWT是无状态的，客户端删除token即可
	// 这里可以添加token黑名单逻辑
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// GetCurrentUser 获取当前用户信息
func (c *AuthController) GetCurrentUser(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := c.userService.GetCurrentUser(userID.(int64))
	if err != nil {
		c.logger.Error("Failed to get current user", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

// ChangePassword 修改密码
func (c *AuthController) ChangePassword(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req model.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := c.authService.ChangePassword(userID.(int64), &req); err != nil {
		c.logger.Error("Failed to change password", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})
}

// RefreshToken 刷新token
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := c.userService.GetUser(userID.(int64))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	token, expiresAt, err := c.authService.GenerateToken(user)
	if err != nil {
		c.logger.Error("Failed to refresh token", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"token":      token,
			"expires_at": expiresAt,
		},
	})
}
