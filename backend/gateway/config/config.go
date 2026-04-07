package config

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Log      LogConfig      `yaml:"log"`
	Redis    RedisConfig    `yaml:"redis"`
	Services ServicesConfig `yaml:"services"`
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

// RedisConfig represents Redis configuration
type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// ServicesConfig represents backend services configuration
type ServicesConfig struct {
	Core    ServiceConfig `yaml:"core"`
	Runbook ServiceConfig `yaml:"runbook"`
	Tenant  ServiceConfig `yaml:"tenant"`
}

// ServiceConfig represents a single service configuration
type ServiceConfig struct {
	URL string `yaml:"url"`
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
			Port:         "8080",
			ReadTimeout:  30,
			WriteTimeout: 30,
		},
		Log: LogConfig{
			Level:      "info",
			Format:     "json",
			OutputPath: "stdout",
		},
		Redis: RedisConfig{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		},
		Services: ServicesConfig{
			Core:    ServiceConfig{URL: "http://localhost:8081"},
			Runbook: ServiceConfig{URL: "http://localhost:8082"},
			Tenant:  ServiceConfig{URL: "http://localhost:8083"},
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
	if redisURL := os.Getenv("REDIS_URL"); redisURL != "" {
		cfg.Redis.Addr = redisURL
	}
	if coreURL := os.Getenv("CORE_SERVICE_URL"); coreURL != "" {
		cfg.Services.Core.URL = coreURL
	}
	if runbookURL := os.Getenv("RUNBOOK_SERVICE_URL"); runbookURL != "" {
		cfg.Services.Runbook.URL = runbookURL
	}
	if tenantURL := os.Getenv("TENANT_SERVICE_URL"); tenantURL != "" {
		cfg.Services.Tenant.URL = tenantURL
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
