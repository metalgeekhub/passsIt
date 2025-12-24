package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	//  "strconv"

	"passIt/internal/auth"
	"passIt/internal/database"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	App         *AppConfig
	Auth        *auth.Config
	DB          *database.DBConfig
	RedisClient *redis.Options
}
type AppConfig struct {
	Port                   int
	ENV                    string
	FrontendURL            string
	BootstrapAdminUsername string
	BootstrapAdminEmail    string
	BootstrapAdminPassword string
}

func LoadFromEnv() (*Config, error) {
	// Get the absolute path of the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Construct path to .env file
	envPath := filepath.Join(currentDir, ".env")
	err = godotenv.Load(envPath)

	// In Docker/production, .env file may not exist - that's ok, use environment variables directly
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables directly")
	}

	redisDB, err := strconv.Atoi(requireEnv("REDIS_DATABASE"))
	if err != nil {
		log.Fatal("failed to convert redis db")
	}

	port, err := strconv.Atoi(requireEnv("PORT"))
	if err != nil {
		log.Fatal("failed to convert PORT to int")
	}
	return &Config{
		App: &AppConfig{
			Port:                   port,
			ENV:                    requireEnv("ENV"),
			FrontendURL:            requireEnv("FRONTEND_URL"),
			BootstrapAdminUsername: os.Getenv("BOOTSTRAP_ADMIN_USERNAME"), // Optional
			BootstrapAdminEmail:    os.Getenv("BOOTSTRAP_ADMIN_EMAIL"),    // Optional
			BootstrapAdminPassword: os.Getenv("BOOTSTRAP_ADMIN_PASSWORD"), // Optional
		},
		DB: &database.DBConfig{
			Host:     requireEnv("DB_HOST"),
			Port:     requireEnv("DB_PORT"),
			Database: requireEnv("DB_DATABASE"),
			Username: requireEnv("DB_USERNAME"),
			Password: requireEnv("DB_PASSWORD"),
			Schema:   requireEnv("DB_SCHEMA"),
		},
		Auth: &auth.Config{
			BaseURL:       requireEnv("KEYCLOAK_URL"),
			ClientID:      requireEnv("KEYCLOAK_CLIENT_ID"),
			Realm:         requireEnv("KEYCLOAK_REALM"),
			ClientSecret:  requireEnv("KEYCLOAK_CLIENT_SECRET"),
			RedirectURL:   requireEnv("REDIRECT_URL"),
			AdminUsername: requireEnv("KEYCLOAK_ADMIN_USERNAME"),
			AdminPassword: requireEnv("KEYCLOAK_ADMIN_PASSWORD"),
			FrontendURL:   requireEnv("FRONTEND_URL"),
		},
		RedisClient: &redis.Options{
			Addr:     fmt.Sprintf("%s:%s", requireEnv("REDIS_HOST"), requireEnv("REDIS_PORT")),
			Username: requireEnv("REDIS_USERNAME"),
			Password: requireEnv("REDIS_PASSWORD"),
			DB:       redisDB,
		},
	}, nil
}

func requireEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("%s is required", key))
	}
	return value
}
