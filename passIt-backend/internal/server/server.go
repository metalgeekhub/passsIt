package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/gorm"

	"passIt/internal/auth"
	"passIt/internal/config"
	"passIt/internal/database"
	"passIt/internal/models"
	"passIt/internal/services"

	// "passIt/internal/store"

	"github.com/redis/go-redis/v9"
)

type Server struct {
	port int

	Keycloak    auth.KeycloakClient
	db          database.Service
	gormDB      *gorm.DB
	userService services.UserService
}

func NewServer(ctx context.Context, cfg *config.Config, authClient *auth.Client, redisClient *redis.Client) *http.Server {

	dbService := database.New()
	
	// Create user service with business logic
	userService := services.NewUserService(dbService, authClient)
	
	NewServer := &Server{
		port: cfg.App.Port,

		Keycloak:    authClient,
		db:          dbService,
		gormDB:      dbService.GetGormDB(),
		userService: userService,
	}

	// Initialize first admin user if none exists
	NewServer.initializeAdminUser(ctx, cfg)

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(ctx, cfg, authClient, redisClient),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 120 * time.Second,
	}

	return server
}

// initializeAdminUser creates the first admin user from environment variables if no admin exists
func (s *Server) initializeAdminUser(ctx context.Context, cfg *config.Config) {
	// Check if any admin users exist
	users, err := s.userService.GetAllUsers(ctx)
	if err != nil {
		log.Printf("Warning: Could not check for existing admin users: %v", err)
		return
	}

	// Check if any user is already admin
	for _, user := range users {
		if user.IsAdmin {
			log.Println("Admin user already exists, skipping bootstrap")
			return
		}
	}

	// Get admin credentials from config
	adminUsername := cfg.App.BootstrapAdminUsername
	adminEmail := cfg.App.BootstrapAdminEmail
	adminPassword := cfg.App.BootstrapAdminPassword

	// If env vars are not set, skip bootstrap
	if adminUsername == "" || adminEmail == "" || adminPassword == "" {
		log.Println("No admin user exists and BOOTSTRAP_ADMIN_* env vars not set")
		log.Println("Please set BOOTSTRAP_ADMIN_USERNAME, BOOTSTRAP_ADMIN_EMAIL, and BOOTSTRAP_ADMIN_PASSWORD")
		log.Println("Or create an admin user manually in the database")
		return
	}

	// Create bootstrap admin user
	adminUser := &models.User{
		Username:  adminUsername,
		Email:     adminEmail,
		FirstName: "Admin",
		LastName:  "User",
		IsAdmin:   true,
		IsActive:  true,
	}

	err = s.userService.CreateUser(ctx, adminUser, adminPassword)
	if err != nil {
		log.Printf("Warning: Failed to create bootstrap admin user: %v", err)
		return
	}

	log.Printf("✓ Bootstrap admin user created successfully: %s (%s)", adminUsername, adminEmail)
}
