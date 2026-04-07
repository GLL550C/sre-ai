package main

import (
	"core/ai"
	"core/config"
	"core/controller"
	"core/middleware"
	"core/repository"
	"core/service"
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
		zap.String("database", cfg.Database.Host),
	)

	// Initialize database
	db, err := initDB(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Initialize config manager
	configManager := config.NewConfigManager(db, redisClient, logger)
	configManager.Start()

	// Initialize repositories
	alertRepo := repository.NewAlertRepository(db, logger)
	ruleRepo := repository.NewRuleRepository(db, logger)
	clusterRepo := repository.NewClusterRepository(db, logger)
	analysisRepo := repository.NewAnalysisRepository(db, logger)
	dashboardRepo := repository.NewDashboardRepository(db, logger)
	aiModelRepo := repository.NewAIModelRepository(db, logger)
	configRepo := repository.NewConfigRepository(db, logger)

	// Initialize new config service and init config items
	configService := service.NewConfigService(configRepo, logger)
	if err := configService.InitConfigItems(); err != nil {
		logger.Warn("Failed to init config items", zap.Error(err))
	}

	// Initialize AI model service
	aiModelService := service.NewAIModelService(aiModelRepo, logger)

	// Initialize AI service from database config (or fallback to env)
	var aiAnalysisService *ai.AnalysisService
	aiModelConfig, err := aiModelService.GetActiveConfig()
	if err == nil && aiModelConfig != nil && aiModelConfig.IsEnabled {
		aiConfig := aiModelService.BuildAIConfigFromModel(aiModelConfig)
		aiAnalysisService, err = ai.NewAnalysisService(aiConfig, logger)
		if err != nil {
			logger.Warn("Failed to initialize AI service from DB config", zap.Error(err))
		} else {
			logger.Info("AI service initialized from database config",
				zap.String("provider", aiModelConfig.Provider),
				zap.String("model", aiModelConfig.Model),
			)
		}
	}

	// Fallback to environment config if DB config not available
	if aiAnalysisService == nil {
		aiConfig := config.LoadAIConfigFromEnv()
		aiAnalysisService, err = ai.NewAnalysisService(aiConfig, logger)
		if err != nil {
			logger.Warn("Failed to initialize AI service from env", zap.Error(err))
		}
	}

	// Initialize services
	alertService := service.NewAlertService(alertRepo, clusterRepo, redisClient, logger)
	ruleService := service.NewRuleService(ruleRepo, logger)
	clusterService := service.NewClusterService(clusterRepo, redisClient, logger)
	analysisService := service.NewAnalysisService(analysisRepo, alertRepo, clusterRepo, redisClient, logger, aiAnalysisService)
	dashboardService := service.NewDashboardService(dashboardRepo, logger)
	prometheusService := service.NewPrometheusService(clusterRepo, redisClient, logger)

	// Initialize controllers
	alertController := controller.NewAlertController(alertService, logger)
	ruleController := controller.NewRuleController(ruleService, logger)
	clusterController := controller.NewClusterController(clusterService, logger)
	analysisController := controller.NewAnalysisController(analysisService, logger)
	dashboardController := controller.NewDashboardController(dashboardService, logger)
	configController := controller.NewConfigController(configService, logger)
	prometheusController := controller.NewPrometheusController(prometheusService, logger)
	aiModelController := controller.NewAIModelController(aiModelService, logger)

	// Setup router
	router := gin.New()
	router.Use(middleware.Logger(logger))
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.CORS())
	router.Use(middleware.NoCache())

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "core"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Alerts
		v1.GET("/alerts", alertController.GetAlerts)
		v1.GET("/alerts/:id", alertController.GetAlert)
		v1.POST("/alerts", alertController.CreateAlert)
		v1.POST("/alerts/webhook", alertController.ReceiveWebhook)
		v1.PUT("/alerts/:id/ack", alertController.AcknowledgeAlert)

		// Rules
		v1.GET("/rules", ruleController.GetRules)
		v1.GET("/rules/:id", ruleController.GetRule)
		v1.POST("/rules", ruleController.CreateRule)
		v1.PUT("/rules/:id", ruleController.UpdateRule)
		v1.DELETE("/rules/:id", ruleController.DeleteRule)

		// Clusters
		v1.GET("/clusters", clusterController.GetClusters)
		v1.GET("/clusters/default", clusterController.GetDefaultCluster)
		v1.GET("/clusters/:id", clusterController.GetCluster)
		v1.POST("/clusters", clusterController.CreateCluster)
		v1.PUT("/clusters/:id", clusterController.UpdateCluster)
		v1.DELETE("/clusters/:id", clusterController.DeleteCluster)
		v1.POST("/clusters/:id/test", clusterController.TestCluster)
		v1.PUT("/clusters/:id/default", clusterController.SetDefaultCluster)

		// Analysis
		v1.GET("/analysis", analysisController.GetAnalysis)
		v1.GET("/analysis/:id", analysisController.GetAnalysisByID)
		v1.POST("/analysis", analysisController.CreateAnalysis)
		v1.DELETE("/analysis/:id", analysisController.DeleteAnalysis)
		v1.PUT("/analysis/:id/archive", analysisController.ArchiveAnalysis)
		v1.GET("/analysis/stats", analysisController.GetAnalysisStats)
		v1.POST("/analysis/compare", analysisController.CompareClusters)

		// AI Model Config endpoints
		v1.GET("/ai/configs", aiModelController.GetConfigs)
		v1.GET("/ai/configs/:id", aiModelController.GetConfig)
		v1.POST("/ai/configs", aiModelController.CreateConfig)
		v1.PUT("/ai/configs/:id", aiModelController.UpdateConfig)
		v1.DELETE("/ai/configs/:id", aiModelController.DeleteConfig)
		v1.POST("/ai/configs/:id/test", aiModelController.TestConfig)
		v1.PUT("/ai/configs/:id/default", aiModelController.SetDefaultConfig)
		v1.GET("/ai/configs/active", aiModelController.GetActiveConfig)

		// AI Chat endpoints
		v1.POST("/ai/chat", analysisController.Chat)
		v1.POST("/ai/chat/stream", analysisController.ChatStream)
		v1.GET("/ai/health", analysisController.AIHealth)
		v1.GET("/ai/model", analysisController.AIModelInfo)

		// Dashboard
		v1.GET("/dashboard", dashboardController.GetDashboard)
		v1.GET("/dashboard/metrics", dashboardController.GetMetrics)

		// Config (new hierarchical config system)
		v1.GET("/config/tree", configController.GetConfigTree)
		v1.GET("/config/items", configController.GetConfigItems)
		v1.GET("/config/items/:key", configController.GetConfigItem)
		v1.PUT("/config/items/:key", configController.UpdateConfigValue)
		v1.POST("/config/batch", configController.UpdateMultipleConfigs)
		v1.POST("/config/items/:key/reset", configController.ResetConfigToDefault)
		v1.GET("/config/settings/system", configController.GetSystemSettings)
		v1.GET("/config/settings/ai", configController.GetAISettings)
		v1.GET("/config/settings/notification", configController.GetNotificationSettings)
		v1.GET("/config/export", configController.ExportConfig)
		v1.POST("/config/import", configController.ImportConfig)

		// Legacy Config (keep for compatibility)
		v1.GET("/configs", configController.GetConfigs)
		v1.POST("/configs/reload", configController.ReloadConfig)

		// Prometheus proxy
		v1.GET("/prometheus/query", prometheusController.Query)
		v1.GET("/prometheus/query_range", prometheusController.QueryRange)
	}

	// Start server
	logger.Info("Core service starting", zap.String("port", cfg.Server.Port))
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}

func initDB(cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn := cfg.User + ":" + cfg.Password + "@tcp(" + cfg.Host + ":" + cfg.Port + ")/" + cfg.Name + "?parseTime=true&charset=utf8mb4"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	return db, nil
}
