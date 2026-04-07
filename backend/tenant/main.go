package main

import (
	"database/sql"
	"tenant/config"
	"tenant/controller"
	"tenant/middleware"
	"tenant/repository"
	"tenant/service"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
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

	// Initialize repositories
	tenantRepo := repository.NewTenantRepository(db, logger)

	// Initialize services
	tenantService := service.NewTenantService(tenantRepo, redisClient, logger)

	// Initialize controllers
	tenantController := controller.NewTenantController(tenantService, logger)

	// Setup router
	router := gin.New()
	router.Use(middleware.Logger(logger))
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.CORS())
	router.Use(middleware.NoCache())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "tenant"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Tenants
		v1.GET("/tenants", tenantController.GetTenants)
		v1.GET("/tenants/:id", tenantController.GetTenant)
		v1.POST("/tenants", tenantController.CreateTenant)
		v1.PUT("/tenants/:id", tenantController.UpdateTenant)
		v1.DELETE("/tenants/:id", tenantController.DeleteTenant)
	}

	// Start server
	logger.Info("Tenant service starting", zap.String("port", cfg.Server.Port))
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
