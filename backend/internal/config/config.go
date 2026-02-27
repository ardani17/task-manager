package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	// Server
	AppPort string
	AppEnv  string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Redis
	RedisHost string
	RedisPort string

	// JWT
	JWTSecret string
	JWTExpiry string
}

var AppConfig *Config

func Load() *Config {
	// Load .env file (ignore error if not exists)
	_ = godotenv.Load()

	config := &Config{
		// Server
		AppPort: getEnv("APP_PORT", "8080"),
		AppEnv:  getEnv("APP_ENV", "development"),

		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5433"),
		DBUser:     getEnv("DB_USER", "taskmanager"),
		DBPassword: getEnv("DB_PASSWORD", "taskmanager123"),
		DBName:     getEnv("DB_NAME", "taskmanager"),

		// Redis
		RedisHost: getEnv("REDIS_HOST", "localhost"),
		RedisPort: getEnv("REDIS_PORT", "6380"),

		// JWT
		JWTSecret: getEnv("JWT_SECRET", "super-secret-key-change-in-production"),
		JWTExpiry: getEnv("JWT_EXPIRY", "24h"),
	}

	AppConfig = config
	log.Info().Msg("Configuration loaded successfully")
	return config
}

func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}

func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}
