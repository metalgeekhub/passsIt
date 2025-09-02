package server

import (
	"log"
	"net/http"
	"passIt/internal/models"
	codes "passIt/internal/passit-codes"

	"passIt/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FindUserByIdRequestBody struct {
	ID uuid.UUID `json:"id"`
}

type FindUserByEmailRequestBody struct {
	Email string `json:"email"`
}

type CreateUserRequestBody struct {
	User     models.User `json:"user"`
	Password string      `json:"password"`
}

type CreateUserReturnBody struct {
	User            models.User `json:"user"`
	KeycloackUserID string      `json:"keycloak_user_id"`
}

type DeleteUserByIdRequestBody struct {
	ID uuid.UUID `json:"id"`
}

type LoginUserRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) LoginHandler(c *gin.Context) {
	var input LoginUserRequestBody

	if !utils.DecodeServerInput(c, &input) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		c.Abort()
		return
	}

	// TODO: hash password and check with database if it is correct
	// Then get from database the userID
	user, err := s.db.FindUserByEmail(input.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User with this email is not exists."})
		c.Abort()
		return
	}

	token, err := utils.GenerateJWTToken(user.ID)
	if err != nil {
		log.Printf("Error when generating JWT token for user: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Error generating JWT token for this user."})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (s *Server) CreateUserHandler(c *gin.Context) {
	var input CreateUserRequestBody
	// var keycloackClient keyclock.KeycloakClient
	// k := keycloackClient.NewKeycloakClient()
	log.Printf("Request Body: %+v\n", c.Request.Body)

	if !utils.DecodeServerInput(c, &input) {
		return // Stop processing if decode fails
	}

	log.Printf("Input: %+v\n", input)

	user := input.User
	// password := input.Password

	// REMOVED THE KEYCLOAK LOGIC

	log.Printf("Creating user: %+v\n", user)

	err := s.db.CreateUser(&user)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	defer c.Request.Body.Close()

	c.JSON(http.StatusOK, PassItResponseBody{
		Code: codes.UserCreatedSuccessfully,
		Data: user,
	},
	)
}

func (s *Server) FindUserByIdHandler(c *gin.Context) {
	var input FindUserByIdRequestBody

	if !utils.DecodeServerInput(c, &input) {
		return // Stop processing if decode fails
	}
	user, err := s.db.FindUserById(input.ID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	defer c.Request.Body.Close()

	c.JSON(http.StatusOK, user)
}

func (s *Server) FindUserByEmailHandler(c *gin.Context) {
	var input FindUserByEmailRequestBody

	if !utils.DecodeServerInput(c, &input) {
		return // Stop processing if decode fails
	}
	user, err := s.db.FindUserByEmail(input.Email)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	defer c.Request.Body.Close()

	c.JSON(http.StatusOK, user)
}

func (s *Server) UpdateUserByIdHandler(c *gin.Context) {
	var user models.User

	if !utils.DecodeServerInput(c, &user) {
		return // Stop processing if decode fails
	}

	log.Printf("Received user update request: %+v\n", user)

	err := s.db.UpdateUserById(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	err = s.db.GetKeycloakIDByUserID(&user)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Keycloak ID"})
		return
	}

	// log.Printf("Fetched Keycloak ID: %s for user ID: %s\n", keycloakId, user.ID)
	// user.KeycloackID = keycloakId

	// REMOVED THE KEYCLOAK LOGIC

	defer c.Request.Body.Close()

	c.JSON(http.StatusOK, user)
}

func (s *Server) GetAllUsersHandler(c *gin.Context) {
	users, err := s.db.GetAllUsers()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	defer c.Request.Body.Close()

	c.JSON(http.StatusOK, users)
}

func (s *Server) DeleteUserByIdHandler(c *gin.Context) {
	var input DeleteUserByIdRequestBody
	if !utils.DecodeServerInput(c, &input) {
		return // Stop processing if decode fails
	}
	// Get Keycloak ID before deleting user
	// TODO: Make this function to return Keycloak ID directly or create another one.
	err := s.db.GetKeycloakIDByUserID(&models.User{ID: input.ID})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Keycloak ID"})
		return
	}
	// REMOVED THE KEYCLOAK LOGIC
	// Delete user in local DB
	err = s.db.DeleteUserById(input.ID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
}
