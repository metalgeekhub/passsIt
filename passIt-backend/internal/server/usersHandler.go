package server

import (
	"log"
	"net/http"
	"passIt/internal/models"
	codes "passIt/internal/passit-codes"
	"passIt/internal/store"

	"passIt/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LoginUserRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type FindUserByIdRequestBody struct {
	ID uuid.UUID `json:"id"`
}

type FindUserByEmailRequestBody struct {
	Email string `json:"email"`
}

type CreateUserRequestBody struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IsAdmin   bool   `json:"is_admin"`
}

type UpdateUserRequestBody struct {
	Password  string `json:"password,omitempty"` // Optional password
	Email     string `json:"email,omitempty"`
	Username  string `json:"username,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	IsAdmin   *bool  `json:"is_admin,omitempty"` // Pointer to distinguish between false and not provided
}

type CreateUserReturnBody struct {
	User            models.User `json:"user"`
	KeycloackUserID string      `json:"keycloak_user_id"`
}

// CreateUserHandler godoc
// @Summary      Create a new user (Admin only)
// @Description  Create a new user with username, email, password and admin status
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user body CreateUserRequestBody true "User creation data"
// @Success      200 {object} PassItResponseBody
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/users [post]
func (s *Server) CreateUserHandler(c *gin.Context) {
	var input CreateUserRequestBody

	if !utils.DecodeServerInput(c, &input) {
		return // Stop processing if decode fails
	}

	// Build user model from flat request
	user := models.User{
		Username:  input.Username,
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		IsAdmin:   input.IsAdmin,
		IsActive:  true,
	}

	log.Printf("Creating user: %+v\n", user)

	// Use service layer - handles all business logic and rollback
	err := s.userService.CreateUser(c, &user, input.Password)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	defer c.Request.Body.Close()

	c.JSON(http.StatusOK, PassItResponseBody{
		Code: codes.UserCreatedSuccessfully,
		Data: user,
	})
}

func (s *Server) FindUserByIdHandler(c *gin.Context) {
	// Get ID from query parameter
	idStr := c.Query("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id query parameter is required"})
		return
	}

	// Parse UUID
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format"})
		return
	}

	user, err := s.userService.GetUserByID(c, id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Server) FindUserByEmailHandler(c *gin.Context) {
	// Get email from query parameter
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email query parameter is required"})
		return
	}

	user, err := s.userService.GetUserByEmail(c, email)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUserByIdHandler godoc
// @Summary      Update user by ID (Admin only)
// @Description  Update user information including email, name, password, and admin status
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID"
// @Param        user body UpdateUserRequestBody true "User update data"
// @Success      200 {object} models.User
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/users/{id} [put]
func (s *Server) UpdateUserByIdHandler(c *gin.Context) {
	// Get ID from URL parameter
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format"})
		return
	}

	// Decode the update request
	var updateReq UpdateUserRequestBody
	if !utils.DecodeServerInput(c, &updateReq) {
		return // Stop processing if decode fails
	}

	// Get existing user
	existingUser, err := s.userService.GetUserByID(c, id)
	if err != nil {
		log.Printf("User not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Apply updates using a helper function
	applyUserUpdates(&existingUser, &updateReq)

	log.Printf("Received user update request for ID %s: %+v\n", id, existingUser)

	// Update in Keycloak first
	err = s.Keycloak.UpdateKeycloakUser(c, &existingUser)
	if err != nil {
		log.Printf("Failed to update user in Keycloak: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user in Keycloak"})
		return
	}

	// Update password in Keycloak if provided
	if updateReq.Password != "" {
		err = s.Keycloak.UpdatePassword(c, existingUser.KeycloackID, updateReq.Password)
		if err != nil {
			log.Printf("Failed to update password: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
			return
		}
	}

	// Update in database last
	err = s.userService.UpdateUser(c, &existingUser)
	if err != nil {
		log.Printf("Failed to update user in database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		// TODO: Implement proper rollback of Keycloak changes
		return
	}

	c.JSON(http.StatusOK, existingUser)
}

// applyUserUpdates applies non-empty fields from update request to existing user
func applyUserUpdates(user *models.User, update *UpdateUserRequestBody) {
	if update.Email != "" {
		user.Email = update.Email
	}
	if update.Username != "" {
		user.Username = update.Username
	}
	if update.FirstName != "" {
		user.FirstName = update.FirstName
	}
	if update.LastName != "" {
		user.LastName = update.LastName
	}
	if update.IsAdmin != nil {
		user.IsAdmin = *update.IsAdmin
	}
}

// DeleteUserByIdHandler godoc
// @Summary      Soft delete user by ID (Admin only)
// @Description  Deactivate a user by setting isActive to false and disabling in Keycloak
// @Tags         users
// @Produce      json
// @Param        id path string true "User ID"
// @Success      200 {object} PassItResponseBody
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/users/{id} [delete]
func (s *Server) DeleteUserByIdHandler(c *gin.Context) {
	// Get ID from URL parameter
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format"})
		return
	}

	// Get existing user
	existingUser, err := s.userService.GetUserByID(c, id)
	if err != nil {
		log.Printf("User not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if user is already inactive
	if !existingUser.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is already inactive"})
		return
	}

	// Deactivate user (soft delete)
	existingUser.IsActive = false

	// Update in database
	err = s.userService.UpdateUser(c, &existingUser)
	if err != nil {
		log.Printf("Failed to deactivate user in database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate user"})
		return
	}

	// Update Keycloak to disable the user
	err = s.Keycloak.UpdateKeycloakUser(c, &existingUser)
	if err != nil {
		// Rollback: reactivate user in database
		existingUser.IsActive = true
		rollbackErr := s.userService.UpdateUser(c, &existingUser)
		if rollbackErr != nil {
			log.Printf("CRITICAL: Failed to rollback user activation after Keycloak error: %v", rollbackErr)
		}
		log.Printf("Failed to disable user in Keycloak: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disable user in Keycloak"})
		return
	}

	c.JSON(http.StatusOK, PassItResponseBody{
		Code: codes.UserDeletedSuccessfully,
		Data: gin.H{
			"message": "User deactivated successfully",
			"user_id": existingUser.ID,
		},
	})
}

// GetAllUsersHandler godoc
// @Summary      Get all users (Admin only)
// @Description  Retrieve a list of all users in the system
// @Tags         users
// @Produce      json
// @Success      200 {array} models.User
// @Failure      500 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/users [get]
func (s *Server) GetAllUsersHandler(c *gin.Context) {
	users, err := s.userService.GetAllUsers(c)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	defer c.Request.Body.Close()

	c.JSON(http.StatusOK, users)
}

// GetInactiveUsersHandler godoc
// @Summary      Get all inactive users (Admin only)
// @Description  Retrieve a list of all deactivated/deleted users in the system
// @Tags         users
// @Produce      json
// @Success      200 {array} models.User
// @Failure      500 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/users/inactive [get]
func (s *Server) GetInactiveUsersHandler(c *gin.Context) {
	users, err := s.userService.GetInactiveUsers(c)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve inactive users"})
		return
	}

	defer c.Request.Body.Close()

	c.JSON(http.StatusOK, users)
}

// GetCurrentUserHandler godoc
// @Summary      Get current user profile
// @Description  Retrieve the profile of the currently authenticated user
// @Tags         users
// @Produce      json
// @Success      200 {object} models.User
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/users/me [get]
func (s *Server) GetCurrentUserHandler(c *gin.Context) {
	// Get session data from middleware
	sessionData, exists := c.Get("user_session")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No session found"})
		return
	}

	// Get user email from session
	session, ok := sessionData.(*store.SessionData)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid session data"})
		return
	}

	// Fetch full user data from database
	user, err := s.userService.GetUserByEmail(c, session.UserInfo.Email)
	if err != nil {
		log.Printf("Failed to get current user: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
