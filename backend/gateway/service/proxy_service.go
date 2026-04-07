package service

import (
	"gateway/config"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// ProxyService handles request proxying to backend services
type ProxyService struct {
	logger      *zap.Logger
	redisClient *redis.Client
	services    config.ServicesConfig
}

// NewProxyService creates a new proxy service
func NewProxyService(logger *zap.Logger, redisClient *redis.Client, services config.ServicesConfig) *ProxyService {
	return &ProxyService{
		logger:      logger,
		redisClient: redisClient,
		services:    services,
	}
}

// ProxyRequest proxies a request to the specified service
func (ps *ProxyService) ProxyRequest(c *gin.Context, serviceName string) {
	var targetURL string
	var apiPrefix string
	switch serviceName {
	case "core":
		targetURL = ps.services.Core.URL
		apiPrefix = "/api/v1"
	case "runbook":
		targetURL = ps.services.Runbook.URL
		apiPrefix = "/api/v1"
	case "tenant":
		targetURL = ps.services.Tenant.URL
		apiPrefix = "/api/v1"
	default:
		ps.logger.Error("Unknown service", zap.String("service", serviceName))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unknown service"})
		return
	}

	// Build target URL
	path := strings.TrimPrefix(c.Param("path"), "/")
	target, err := url.Parse(targetURL)
	if err != nil {
		ps.logger.Error("Failed to parse target URL", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Add API prefix for backend service routing
	target.Path = apiPrefix + "/" + path
	target.RawQuery = c.Request.URL.RawQuery

	// Create new request
	proxyReq, err := http.NewRequest(c.Request.Method, target.String(), c.Request.Body)
	if err != nil {
		ps.logger.Error("Failed to create proxy request", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Copy headers
	for key, values := range c.Request.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}

	// Add X-Forwarded headers
	proxyReq.Header.Set("X-Forwarded-For", c.ClientIP())
	proxyReq.Header.Set("X-Forwarded-Host", c.Request.Host)
	proxyReq.Header.Set("X-Forwarded-Proto", c.Request.URL.Scheme)

	// Disable caching for API requests
	proxyReq.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	proxyReq.Header.Set("Pragma", "no-cache")

	// Execute request
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(proxyReq)
	if err != nil {
		ps.logger.Error("Failed to proxy request", zap.Error(err), zap.String("service", serviceName))
		c.JSON(http.StatusBadGateway, gin.H{"error": "Service unavailable"})
		return
	}
	defer resp.Body.Close()

	// Copy response headers (skip CORS headers to preserve gateway's CORS settings)
	for key, values := range resp.Header {
		// Skip CORS headers - these are set by gateway's CORS middleware
		lowerKey := strings.ToLower(key)
		if lowerKey == "access-control-allow-origin" ||
			lowerKey == "access-control-allow-credentials" ||
			lowerKey == "access-control-allow-headers" ||
			lowerKey == "access-control-allow-methods" ||
			lowerKey == "access-control-expose-headers" {
			continue
		}
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Force disable caching on response
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	// Ensure CORS headers are set (in case backend overwrote them)
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

	// Copy response status and body
	c.Status(resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ps.logger.Error("Failed to read response body", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Writer.Write(body)
}
