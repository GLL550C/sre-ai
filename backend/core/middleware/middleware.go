package middleware

import (
	"core/service"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger returns a gin middleware for logging
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		fields := []zap.Field{
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
			zap.String("client_ip", clientIP),
			zap.String("method", method),
			zap.String("path", path),
		}

		if len(c.Errors) > 0 {
			fields = append(fields, zap.Strings("errors", c.Errors.Errors()))
			logger.Error("Request failed", fields...)
		} else if statusCode >= 500 {
			logger.Error("Server error", fields...)
		} else if statusCode >= 400 {
			logger.Warn("Client error", fields...)
		} else {
			logger.Info("Request processed", fields...)
		}
	}
}

// Recovery returns a gin middleware for panic recovery
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return gin.RecoveryWithWriter(&zapWriter{logger: logger})
}

type zapWriter struct {
	logger *zap.Logger
}

func (w *zapWriter) Write(p []byte) (n int, err error) {
	w.logger.Error("Panic recovered", zap.String("error", string(p)))
	return len(p), nil
}

// CORS returns a gin middleware for CORS
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// NoCache returns a gin middleware to disable caching
func NoCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Writer.Header().Set("Pragma", "no-cache")
		c.Writer.Header().Set("Expires", "0")
		c.Next()
	}
}

// JWTAuth returns a gin middleware for JWT authentication
func JWTAuth(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过登录和验证码接口 (支持直接访问和通过gateway代理)
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/v1/auth/") || strings.HasPrefix(path, "/api/core/auth/") {
			c.Next()
			return
		}

		// 跳过健康检查和metrics端点（Prometheus监控需要）
		if path == "/health" || path == "/metrics" {
			c.Next()
			return
		}

		// 跳过获取系统名称接口（登录页面需要显示）
		if strings.Contains(path, "/config/items/app.name") {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		user, err := authService.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("userID", user.ID)
		c.Set("username", user.Username)
		c.Set("role", user.Role)

		c.Next()
	}
}

// RequireRole returns a gin middleware that requires specific role
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}

		userRole := role.(string)
		for _, r := range roles {
			if userRole == r {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	}
}

// RequireAdmin returns a gin middleware that requires admin role
func RequireAdmin() gin.HandlerFunc {
	return RequireRole("admin")
}
