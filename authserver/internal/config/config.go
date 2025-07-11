package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Email    EmailConfig
	WebAuthn WebAuthnConfig
}

type ServerConfig struct {
	Host            string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	Environment     string
}

type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type RedisConfig struct {
	URL      string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

type EmailConfig struct {
	FromEmail    string
	FromName     string
	ResendAPIKey string
}

type WebAuthnConfig struct {
	RPID        string
	RPOrigins   []string
	RPDisplayName string
}

func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Host:            getEnv("HOST", "0.0.0.0"),
			Port:            getEnvAsInt("PORT", 8080),
			ReadTimeout:     getEnvAsDuration("READ_TIMEOUT", "30s"),
			WriteTimeout:    getEnvAsDuration("WRITE_TIMEOUT", "30s"),
			ShutdownTimeout: getEnvAsDuration("SHUTDOWN_TIMEOUT", "30s"),
			Environment:     getEnv("ENVIRONMENT", "development"),
		},
		Database: DatabaseConfig{
			URL:             getEnv("DATABASE_URL", "postgres://localhost/simple_auth_roles?sslmode=disable"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 25),
			ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", "5m"),
			ConnMaxIdleTime: getEnvAsDuration("DB_CONN_MAX_IDLE_TIME", "5m"),
		},
		Redis: RedisConfig{
			URL:      getEnv("REDIS_URL", ""),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
			Expiration: getEnvAsDuration("JWT_EXPIRATION", "7d"),
		},
		Email: EmailConfig{
			FromEmail:    getEnv("FROM_EMAIL", "auth@yourapp.com"),
			FromName:     getEnv("FROM_NAME", "Simple Auth"),
			ResendAPIKey: getEnv("RESEND_API_KEY", ""),
		},
		WebAuthn: WebAuthnConfig{
			RPID:          getEnv("WEBAUTHN_RPID", "localhost"),
			RPOrigins:     getEnvAsSlice("WEBAUTHN_RP_ORIGINS", []string{"http://localhost:3000"}),
			RPDisplayName: getEnv("WEBAUTHN_RP_DISPLAY_NAME", "Auth Template"),
		},
	}

	// Validate required config
	if config.JWT.Secret == "your-super-secret-jwt-key-change-in-production" && config.IsProduction() {
		return nil, fmt.Errorf("JWT_SECRET must be set in production")
	}

	return config, nil
}

func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	if duration, err := time.ParseDuration(defaultValue); err == nil {
		return duration
	}
	return 30 * time.Second
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}
