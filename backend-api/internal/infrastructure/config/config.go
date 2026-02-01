package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Worker   WorkerConfig
	Storage  StorageConfig
	CORS     CORSConfig
	Logging  LoggingConfig
}

type ServerConfig struct {
	Port string
	Host string
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type WorkerConfig struct {
	PythonPath  string
	WorkerPath  string
	Concurrency int
	Timeout     time.Duration
}

type StorageConfig struct {
	Path          string
	MaxUploadSize int64
}

type CORSConfig struct {
	Origins []string
}

type LoggingConfig struct {
	Level  string
	Format string
}

// Load loads configuration from environment variables and .env file
func Load() (*Config, error) {
	// Set config file
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Read .env file if it exists (not required in production)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found; using environment variables only
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnvOrDefault("PORT", "8080"),
			Host: getEnvOrDefault("HOST", "0.0.0.0"),
			Env:  getEnvOrDefault("ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvOrDefault("DB_PORT", "5432"),
			User:     getEnvOrDefault("DB_USER", "manga_user"),
			Password: getEnvOrDefault("DB_PASSWORD", "secure_pass"),
			Name:     getEnvOrDefault("DB_NAME", "manga_translator"),
			SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Addr:     getEnvOrDefault("REDIS_ADDR", "localhost:6379"),
			Password: getEnvOrDefault("REDIS_PASSWORD", ""),
			DB:       viper.GetInt("REDIS_DB"),
		},
		Worker: WorkerConfig{
			PythonPath:  getEnvOrDefault("PYTHON_PATH", "python"),
			WorkerPath:  getEnvOrDefault("WORKER_PATH", "../ai-worker"),
			Concurrency: getIntOrDefault("WORKER_CONCURRENCY", 1),
			Timeout:     time.Duration(getIntOrDefault("WORKER_TIMEOUT", 600)) * time.Second,
		},
		Storage: StorageConfig{
			Path:          getEnvOrDefault("STORAGE_PATH", "./storage"),
			MaxUploadSize: int64(getIntOrDefault("MAX_UPLOAD_SIZE", 104857600)), // 100MB
		},
		CORS: CORSConfig{
			Origins: viper.GetStringSlice("CORS_ORIGINS"),
		},
		Logging: LoggingConfig{
			Level:  getEnvOrDefault("LOG_LEVEL", "info"),
			Format: getEnvOrDefault("LOG_FORMAT", "json"),
		},
	}

	// Set default CORS origins if not specified
	if len(cfg.CORS.Origins) == 0 {
		cfg.CORS.Origins = []string{"http://localhost:3000"}
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}
	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Redis.Addr == "" {
		return fmt.Errorf("redis address is required")
	}
	if c.Worker.PythonPath == "" {
		return fmt.Errorf("python path is required")
	}
	if c.Worker.WorkerPath == "" {
		return fmt.Errorf("worker path is required")
	}
	return nil
}

// DSN returns the database connection string
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

// getEnvOrDefault gets environment variable or returns default value
func getEnvOrDefault(key, defaultValue string) string {
	viper.SetDefault(key, defaultValue)
	return viper.GetString(key)
}

// getIntOrDefault gets integer environment variable or returns default value
func getIntOrDefault(key string, defaultValue int) int {
	viper.SetDefault(key, defaultValue)
	return viper.GetInt(key)
}
