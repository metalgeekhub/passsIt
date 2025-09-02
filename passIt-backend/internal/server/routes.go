package server

import (
	"context"
	"net/http"

	"passIt/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func (s *Server) RegisterRoutes(ctx context.Context, redisClient *redis.Client) http.Handler {
	// gin.SetMode(gin.ReleaseMode) // Set Gin to release mode
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // TODO: Add your frontend URL from env variables
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	r.POST("/user", s.CreateUserHandler)

	auth := r.Group("/auth")
	{
		auth.GET("/login", s.LoginHandler)
	}

	// r.GET("/login", s.LoginUserHandler)

	r.Use(middleware.AuthenticationMiddleware()) // Apply the auth middleware to all routes

	r.GET("/health", s.healthHandler)
	r.GET("/user/find", s.FindUserByIdHandler)
	r.GET("/user", s.FindUserByEmailHandler)
	r.PUT("/user/update", s.UpdateUserByIdHandler)
	r.GET("/users", s.GetAllUsersHandler)
	r.DELETE("/user/delete", s.DeleteUserByIdHandler)

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
