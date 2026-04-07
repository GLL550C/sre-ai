package config

import (
	"context"
	"core/model"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Log      LogConfig      `yaml:"log"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port         string `yaml:"port"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
}

// LogConfig represents logging configuration
type LogConfig struct {
	Level      string `yaml:"level"`
	Format     string `yaml:"format"`
	OutputPath string `yaml:"output_path"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Name         string `yaml:"name"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
}

// RedisConfig represents Redis configuration
type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// LoadConfig loads configuration from YAML file
func LoadConfig(path string) (*Config, error) {
	// If path is not provided, try default locations
	if path == "" {
		path = "config.yaml"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		// Try to load from config directory
		data, err = os.ReadFile("/app/config.yaml")
		if err != nil {
			// Return default config
			return defaultConfig(), nil
		}
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply environment variable overrides
	applyEnvOverrides(&cfg)

	return &cfg, nil
}

// defaultConfig returns default configuration
func defaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         "8081",
			ReadTimeout:  30,
			WriteTimeout: 30,
		},
		Log: LogConfig{
			Level:      "info",
			Format:     "json",
			OutputPath: "stdout",
		},
		Database: DatabaseConfig{
			Host:         "localhost",
			Port:         "3306",
			User:         "root",
			Password:     "root123",
			Name:         "sre_platform",
			MaxOpenConns: 25,
			MaxIdleConns: 10,
		},
		Redis: RedisConfig{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		},
	}
}

// applyEnvOverrides applies environment variable overrides to config
func applyEnvOverrides(cfg *Config) {
	if port := os.Getenv("PORT"); port != "" {
		cfg.Server.Port = port
	}
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		cfg.Log.Level = logLevel
	}
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		cfg.Database.Host = dbHost
	}
	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		cfg.Database.Port = dbPort
	}
	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		cfg.Database.User = dbUser
	}
	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		cfg.Database.Password = dbPassword
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		cfg.Database.Name = dbName
	}
	if redisURL := os.Getenv("REDIS_URL"); redisURL != "" {
		cfg.Redis.Addr = redisURL
	}
}

// InitLogger initializes the logger based on configuration
func InitLogger(cfg LogConfig) *zap.Logger {
	level := zap.InfoLevel
	switch cfg.Level {
	case "debug":
		level = zap.DebugLevel
		gin.SetMode(gin.DebugMode)
	case "info":
		level = zap.InfoLevel
		gin.SetMode(gin.ReleaseMode)
	case "warn":
		level = zap.WarnLevel
		gin.SetMode(gin.ReleaseMode)
	case "error":
		level = zap.ErrorLevel
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}

	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: cfg.Level == "debug",
		Encoding:    cfg.Format,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{cfg.OutputPath},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	return logger
}

// ConfigManager manages platform configuration with hot reload
type ConfigManager struct {
	db          *sql.DB
	redisClient *redis.Client
	logger      *zap.Logger
	configs     map[string]string
	mu          sync.RWMutex
	ticker      *time.Ticker
	stopCh      chan bool
}

// NewConfigManager creates a new config manager
func NewConfigManager(db *sql.DB, redisClient *redis.Client, logger *zap.Logger) *ConfigManager {
	return &ConfigManager{
		db:          db,
		redisClient: redisClient,
		logger:      logger,
		configs:     make(map[string]string),
		stopCh:      make(chan bool),
	}
}

// Start starts the config manager with periodic reload
func (cm *ConfigManager) Start() {
	// Load initial config
	if err := cm.LoadConfig(); err != nil {
		cm.logger.Error("Failed to load initial config", zap.Error(err))
	}

	// Start periodic reload (every 10 seconds)
	cm.ticker = time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-cm.ticker.C:
				if err := cm.LoadConfig(); err != nil {
					cm.logger.Error("Failed to reload config", zap.Error(err))
				}
			case <-cm.stopCh:
				cm.ticker.Stop()
				return
			}
		}
	}()

	cm.logger.Info("Config manager started")
}

// Stop stops the config manager
func (cm *ConfigManager) Stop() {
	close(cm.stopCh)
}

// LoadConfig loads configuration from database
func (cm *ConfigManager) LoadConfig() error {
	rows, err := cm.db.Query("SELECT config_key, config_value FROM platform_configs")
	if err != nil {
		return err
	}
	defer rows.Close()

	newConfigs := make(map[string]string)
	for rows.Next() {
		var config model.PlatformConfig
		if err := rows.Scan(&config.ConfigKey, &config.ConfigValue); err != nil {
			cm.logger.Error("Failed to scan config", zap.Error(err))
			continue
		}
		newConfigs[config.ConfigKey] = config.ConfigValue
	}

	// Update configs
	cm.mu.Lock()
	cm.configs = newConfigs
	cm.mu.Unlock()

	// Cache in Redis
	ctx := context.Background()
	configJSON, _ := json.Marshal(newConfigs)
	cm.redisClient.Set(ctx, "platform:configs", configJSON, 5*time.Minute)

	cm.logger.Info("Configuration loaded", zap.Int("count", len(newConfigs)))
	return nil
}

// Get gets a config value by key
func (cm *ConfigManager) Get(key string) string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.configs[key]
}

// GetAll gets all configs
func (cm *ConfigManager) GetAll() map[string]string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	result := make(map[string]string)
	for k, v := range cm.configs {
		result[k] = v
	}
	return result
}

// Reload forces a config reload
func (cm *ConfigManager) Reload() error {
	return cm.LoadConfig()
}
