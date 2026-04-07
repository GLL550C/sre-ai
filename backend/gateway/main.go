package main

import (
	"gateway/config"
	"gateway/controller"
	"gateway/middleware"
	"gateway/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	// Initialize logger
	logger := config.InitLogger(cfg.Log)
	defer logger.Sync()

	logger.Info("Configuration loaded",
		zap.String("port", cfg.Server.Port),
		zap.String("log_level", cfg.Log.Level),
	)

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Initialize services
	proxyService := service.NewProxyService(logger, redisClient, cfg.Services)

	// Initialize controllers
	gatewayController := controller.NewGatewayController(proxyService, logger)

	// Setup router
	router := gin.New()
	router.Use(middleware.Logger(logger))
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.CORS())

	// Health check
	router.GET("/health", gatewayController.HealthCheck)
	router.GET("/metrics", gatewayController.Metrics)

	// API routes
	api := router.Group("/api")
	{
		// Core service routes
		api.Any("/core/*path", gatewayController.ProxyToCore)

		// Runbook service routes
		api.Any("/runbook/*path", gatewayController.ProxyToRunbook)

		// Tenant service routes
		api.Any("/tenant/*path", gatewayController.ProxyToTenant)
	}

	// Start server
	logger.Info("Gateway service starting", zap.String("port", cfg.Server.Port))
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
