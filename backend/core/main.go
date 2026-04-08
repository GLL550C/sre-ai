package main

import (
	"core/ai"
	"core/config"
	"core/controller"
	"core/middleware"
	"core/repository"
	"core/service"
	"database/sql"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		panic("加载配置失败: " + err.Error())
	}

	// 初始化日志
	logger := config.InitLogger(cfg.Log)
	defer logger.Sync()

	logger.Info("配置加载完成",
		zap.String("port", cfg.Server.Port),
		zap.String("database", cfg.Database.Host),
	)

	// 初始化数据库
	db, err := initDB(cfg.Database)
	if err != nil {
		logger.Fatal("连接数据库失败", zap.Error(err))
	}
	defer db.Close()

	// 初始化Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// 初始化仓库
	alertRepo := repository.NewAlertRepository(db, logger)
	ruleRepo := repository.NewRuleRepository(db, logger)
	clusterRepo := repository.NewClusterRepository(db, logger)
	analysisRepo := repository.NewAnalysisRepository(db, logger)
	configRepo := repository.NewConfigRepository(db, logger)
	userRepo := repository.NewUserRepository(db, logger)

	// 初始化AI模型仓库
	aiModelRepo := repository.NewAIModelRepository(db, logger)

	// 初始化服务
	alertService := service.NewAlertService(alertRepo, logger)
	ruleService := service.NewRuleService(ruleRepo, logger)
	clusterService := service.NewClusterService(clusterRepo, redisClient, logger)
	dashboardService := service.NewDashboardService(alertRepo, clusterRepo, logger)
	aiModelService := service.NewAIModelService(aiModelRepo, logger)
	analysisService := service.NewAnalysisService(analysisRepo, alertRepo, clusterRepo, redisClient, logger, nil, aiModelService)
	configService := service.NewConfigService(configRepo, logger)
	authService := service.NewAuthService(userRepo, redisClient, logger)
	userService := service.NewUserService(userRepo, authService, logger)

	// 初始化AI服务（如果配置了）
	aiAnalysisService, err := initAIService(db, logger)
	if err != nil {
		logger.Warn("AI服务初始化失败", zap.Error(err))
	}
	if aiAnalysisService != nil {
		// AI服务已通过依赖注入设置
		logger.Info("AI服务已初始化")
	}

	// 初始化控制器
	alertController := controller.NewAlertController(alertService, logger)
	ruleController := controller.NewRuleController(ruleService, logger)
	clusterController := controller.NewClusterController(clusterService, logger)
	dashboardController := controller.NewDashboardController(dashboardService, logger)
	analysisController := controller.NewAnalysisController(analysisService, logger)
	configController := controller.NewConfigController(configService, logger)
	authController := controller.NewAuthController(authService, userService, logger)
	userController := controller.NewUserController(userService, logger)
	aiModelController := controller.NewAIModelController(aiModelService, logger)

	// 设置路由
	router := gin.New()
	router.Use(middleware.Logger(logger))
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.CORS())
	router.Use(middleware.JWTAuth(authService))

	// 指标和健康检查
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "core"})
	})

	// API v1
	v1 := router.Group("/api/v1")
	{
		// 认证
		v1.GET("/auth/captcha", authController.GetCaptcha)
		v1.POST("/auth/login", authController.Login)
		v1.POST("/auth/logout", authController.Logout)
		v1.GET("/auth/me", authController.GetCurrentUser)
		v1.POST("/auth/change-password", authController.ChangePassword)

		// 用户管理
		v1.GET("/users", middleware.RequireAdmin(), userController.ListUsers)
		v1.GET("/users/:id", middleware.RequireAdmin(), userController.GetUser)
		v1.POST("/users", middleware.RequireAdmin(), userController.CreateUser)
		v1.PUT("/users/:id", middleware.RequireAdmin(), userController.UpdateUser)
		v1.DELETE("/users/:id", middleware.RequireAdmin(), userController.DeleteUser)
		v1.POST("/users/:id/reset-password", middleware.RequireAdmin(), userController.ResetPassword)

		// 告警
		v1.GET("/alerts", alertController.GetAlerts)
		v1.GET("/alerts/:id", alertController.GetAlert)
		v1.POST("/alerts", alertController.CreateAlert)
		v1.POST("/alerts/webhook", alertController.ReceiveWebhook)
		v1.PUT("/alerts/:id/ack", alertController.AcknowledgeAlert)

		// 规则
		v1.GET("/rules", ruleController.GetRules)
		v1.GET("/rules/:id", ruleController.GetRule)
		v1.POST("/rules", ruleController.CreateRule)
		v1.PUT("/rules/:id", ruleController.UpdateRule)
		v1.DELETE("/rules/:id", ruleController.DeleteRule)

		// 集群
		v1.GET("/clusters", clusterController.GetClusters)
		v1.GET("/clusters/default", clusterController.GetDefaultCluster)
		v1.GET("/clusters/:id", clusterController.GetCluster)
		v1.POST("/clusters", clusterController.CreateCluster)
		v1.PUT("/clusters/:id", clusterController.UpdateCluster)
		v1.DELETE("/clusters/:id", clusterController.DeleteCluster)
		v1.POST("/clusters/:id/test", clusterController.TestCluster)
		v1.PUT("/clusters/:id/default", clusterController.SetDefaultCluster)

		// 仪表板
		v1.GET("/dashboard", dashboardController.GetDashboard)
		v1.GET("/dashboard/metrics", dashboardController.GetMetrics)

		// AI分析
		v1.GET("/analysis", analysisController.GetAnalysis)
		v1.GET("/analysis/:id", analysisController.GetAnalysisByID)
		v1.POST("/analysis", analysisController.CreateAnalysis)
		v1.DELETE("/analysis/:id", analysisController.DeleteAnalysis)
		v1.GET("/analysis/stats", analysisController.GetAnalysisStats)
		v1.POST("/ai/chat", analysisController.Chat)
		v1.GET("/ai/health", analysisController.AIHealth)
		v1.GET("/ai/model", analysisController.AIModelInfo)

		// AI模型配置
		v1.GET("/ai/configs", aiModelController.GetConfigs)
		v1.GET("/ai/configs/:id", aiModelController.GetConfig)
		v1.POST("/ai/configs", aiModelController.CreateConfig)
		v1.PUT("/ai/configs/:id", aiModelController.UpdateConfig)
		v1.DELETE("/ai/configs/:id", aiModelController.DeleteConfig)
		v1.POST("/ai/configs/:id/test", aiModelController.TestConfig)
		v1.PUT("/ai/configs/:id/default", aiModelController.SetDefaultConfig)
		v1.GET("/ai/configs/active", aiModelController.GetActiveConfig)

		// 配置
		v1.GET("/config", configController.GetByCategory)
		v1.GET("/config/items", configController.GetConfigItems)
		v1.GET("/config/items/:key", configController.GetConfigItemByKey)
		v1.PUT("/config/items/:key", configController.Update)
		v1.GET("/config/app/name", configController.GetAppName)
		v1.GET("/config/:key", configController.GetByKey)
		v1.PUT("/config/:key", configController.Update)
	}

	// 启动服务
	logger.Info("Core服务启动", zap.String("port", cfg.Server.Port))
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		logger.Fatal("启动失败", zap.Error(err))
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
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	return db, nil
}

func initAIService(db *sql.DB, logger *zap.Logger) (*ai.AnalysisService, error) {
	// 从环境变量或配置加载AI配置
	aiConfig := &config.AIConfig{
		Provider:    getEnv("AI_PROVIDER", "openai"),
		Model:       getEnv("AI_MODEL", "gpt-4"),
		APIKey:      getEnv("AI_API_KEY", ""),
		BaseURL:     getEnv("AI_BASE_URL", ""),
		MaxTokens:   4000,
		Temperature: 0.7,
		Timeout:     60,
		Enabled:     getEnv("AI_API_KEY", "") != "",
	}

	if !aiConfig.Enabled {
		return nil, nil
	}

	return ai.NewAnalysisService(aiConfig, logger)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}