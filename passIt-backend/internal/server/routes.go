package server

import (
	"context"
	"net/http"

	"passIt/internal/auth"
	"passIt/internal/handlers"
	"passIt/internal/middleware"
	"passIt/internal/store"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func (s *Server) RegisterRoutes(ctx context.Context, authClient *auth.Client, redisClient *redis.Client) http.Handler {
	// gin.SetMode(gin.ReleaseMode) // Set Gin to release mode
	r := gin.Default()
	r.LoadHTMLGlob("./internal/templates/*.*")

	authStore := store.NewAuthRedisManager(redisClient)
	sessionStore := store.NewSessionRedisManager(redisClient)

	authHandler := handlers.NewAuthHandler(authClient, authStore, sessionStore)
	// Initialize the auth middleware with your Keycloak configuration
	authMiddleware := middleware.NewAuthMiddleware(ctx, authClient, sessionStore)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // TODO: Add your frontend URL from env variables
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	// Serve login page
	r.GET("/", authHandler.ShowLoginPage)
	r.POST("/user", s.CreateUserHandler)

	auth := r.Group("/auth")
	{
		auth.GET("/login", authHandler.LoginHandler)
		auth.GET("/logout", authHandler.LogoutHandler)
		auth.GET("/callback", authHandler.CallbackHandler)
	}

	// r.GET("/login", s.LoginUserHandler)

	r.Use(authMiddleware.RequireAuth()) // Apply the auth middleware to all routes

	r.GET("/health", s.healthHandler)
	r.GET("/user/find", s.FindUserByIdHandler)
	r.GET("/user", s.FindUserByEmailHandler)
	r.PUT("/user/update", s.UpdateUserByIdHandler)
	r.GET("/users", s.GetAllUsersHandler)

	// api := r.Group("/api")

	// api.Use(auth())

	// api.GET("/user", s.FindUserByEmailHandler)
	// api.GET("/users", s.GetAllUsersHandler)
	// api.PUT("/user", s.UpdateUserByIdHandler)
	// api.GET("/user/find", s.FindUserByIdHandler)

	// tf := r.Group("/api/terraform")
	// tf.GET("/init", s.TerraformInitHandler)

	return r
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
