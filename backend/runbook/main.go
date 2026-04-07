package main

import (
	"database/sql"
	"runbook/config"
	"runbook/controller"
	"runbook/middleware"
	"runbook/repository"
	"runbook/service"

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
	runbookRepo := repository.NewRunbookRepository(db, logger)

	// Initialize services
	runbookService := service.NewRunbookService(runbookRepo, redisClient, logger)

	// Initialize controllers
	runbookController := controller.NewRunbookController(runbookService, logger)

	// Setup router
	router := gin.New()
	router.Use(middleware.Logger(logger))
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.CORS())
	router.Use(middleware.NoCache())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "runbook"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Runbooks
		v1.GET("/runbooks", runbookController.GetRunbooks)
		v1.GET("/runbooks/:id", runbookController.GetRunbook)
		v1.POST("/runbooks", runbookController.CreateRunbook)
		v1.PUT("/runbooks/:id", runbookController.UpdateRunbook)
		v1.DELETE("/runbooks/:id", runbookController.DeleteRunbook)
		v1.GET("/runbooks/search", runbookController.SearchRunbooks)
	}

	// Start server
	logger.Info("Runbook service starting", zap.String("port", cfg.Server.Port))
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
