package services

import (
	"context"
	"fmt"
	"log"
	"passIt/internal/auth"
	"passIt/internal/database"
	"passIt/internal/models"

	"github.com/google/uuid"
)

// UserService handles all user-related business logic
type UserService interface {
	CreateUser(ctx context.Context, user *models.User, password string) error
	GetUserByID(ctx context.Context, id uuid.UUID) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type userService struct {
	db       database.Service
	keycloak auth.KeycloakClient
}

// NewUserService creates a new user service
func NewUserService(db database.Service, keycloak auth.KeycloakClient) UserService {
	return &userService{
		db:       db,
		keycloak: keycloak,
	}
}

// CreateUser creates a user in Keycloak for authentication and stores user data in PostgreSQL
// PostgreSQL is the source of truth for user data
func (s *userService) CreateUser(ctx context.Context, user *models.User, password string) error {
	// Step 1: Create user in Keycloak (for authentication only)
	keycloakUserID, err := s.keycloak.CreateKeycloakUser(ctx, user, password)
	if err != nil {
		return fmt.Errorf("failed to create user in Keycloak: %w", err)
	}

	// Step 2: Store user data in PostgreSQL (source of truth)
	user.KeycloackID = keycloakUserID
	err = s.db.CreateUser(user)
	if err != nil {
		// Rollback: Delete from Keycloak since DB creation failed
		if deleteErr := s.keycloak.DeleteKeycloakUser(ctx, keycloakUserID); deleteErr != nil {
			log.Printf("Failed to rollback Keycloak user %s: %v", keycloakUserID, deleteErr)
		}
		return fmt.Errorf("failed to create user in database: %w", err)
	}

	return nil
}

// GetUserByID retrieves a user by ID from the database
func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (models.User, error) {
	user, err := s.db.FindUserById(id)
	if err != nil {
		return models.User{}, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

// GetUserByEmail retrieves a user by email from the database
func (s *userService) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	user, err := s.db.FindUserByEmail(email)
	if err != nil {
		return models.User{}, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

// GetAllUsers retrieves all users from the database
func (s *userService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	users, err := s.db.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}
	return users, nil
}

// UpdateUser updates user data in both PostgreSQL and Keycloak
// PostgreSQL is updated first as it's the source of truth
func (s *userService) UpdateUser(ctx context.Context, user *models.User) error {
	// Step 1: Update in PostgreSQL (source of truth)
	err := s.db.UpdateUserById(user)
	if err != nil {
		return fmt.Errorf("failed to update user in database: %w", err)
	}

	// Step 2: Get Keycloak ID if not present
	if user.KeycloackID == "" {
		err = s.db.GetKeycloakIDByUserID(user)
		if err != nil {
			return fmt.Errorf("failed to get Keycloak ID: %w", err)
		}
	}

	// Step 3: Sync to Keycloak (for authentication)
	err = s.keycloak.UpdateKeycloakUser(ctx, user)
	if err != nil {
		log.Printf("Warning: Failed to sync user to Keycloak: %v", err)
		// Don't fail the operation - PostgreSQL is source of truth
		// User update succeeded in DB, Keycloak sync can be retried
	}

	return nil
}

// DeleteUser deletes a user from both systems
func (s *userService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	// Get user to retrieve Keycloak ID
	user, err := s.db.FindUserById(id)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Delete from Keycloak first (can be recreated if DB delete fails)
	if user.KeycloackID != "" {
		err = s.keycloak.DeleteKeycloakUser(ctx, user.KeycloackID)
		if err != nil {
			log.Printf("Warning: Failed to delete user from Keycloak: %v", err)
			// Continue with DB deletion even if Keycloak fails
		}
	}

	// Delete from database (source of truth)
	// TODO: Implement soft delete in database layer
	return fmt.Errorf("delete not yet implemented")
}
