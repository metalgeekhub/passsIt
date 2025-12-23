package server

import (
	"context"
	"net/http"

	"passIt/internal/auth"
	"passIt/internal/config"
	"passIt/internal/handlers"
	"passIt/internal/middleware"
	"passIt/internal/store"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	_ "passIt/docs" // Import generated docs

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) RegisterRoutes(ctx context.Context, cfg *config.Config, authClient *auth.Client, redisClient *redis.Client) http.Handler {
	// gin.SetMode(gin.ReleaseMode) // Set Gin to release mode
	r := gin.Default()
	r.LoadHTMLGlob("./internal/templates/*.*")

	// No need for authStore - state is in cookies now (simpler!)
	sessionStore := store.NewSessionRedisManager(redisClient)

	authHandler := handlers.NewAuthHandler(authClient, sessionStore, s.db, cfg.App.FrontendURL)
	// Initialize the auth middleware with your Keycloak configuration
	authMiddleware := middleware.NewAuthMiddleware(ctx, authClient, sessionStore, s.db)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.App.FrontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	// Public routes - no authentication required
	r.GET("/", authHandler.ShowLoginPage)
	r.GET("/health", s.healthHandler)
	
	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Browser OAuth flow routes (for web frontend)
	auth := r.Group("/auth")
	{
		auth.GET("/login", authHandler.LoginHandler)
		auth.GET("/logout", authHandler.LogoutHandler)
		auth.GET("/callback", authHandler.CallbackHandler)
		auth.POST("/signup", authHandler.SignupHandler) // Public signup
	}

	// API routes - support both session cookies (browser) and Bearer tokens (API clients)
	api := r.Group("/api")
	api.Use(authMiddleware.RequireAuth()) // Apply auth middleware
	{
		// Available to all authenticated users
		api.GET("/users/me", s.GetCurrentUserHandler) // Get current user profile
		api.GET("/users/find", s.FindUserByIdHandler)
		api.GET("/users/by-email", s.FindUserByEmailHandler)
		
		// Admin-only endpoints
		adminAPI := api.Group("")
		adminAPI.Use(authMiddleware.RequireAdmin()) // Admin-only middleware
		{
			adminAPI.GET("/users", s.GetAllUsersHandler)
			adminAPI.POST("/users", s.CreateUserHandler)
			adminAPI.PUT("/users/:id", s.UpdateUserByIdHandler)
		}
	}

	return r
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
