package auth

import (
	"context"
	"passIt/internal/models"
)

// KeycloakClient interface for authentication operations
// This allows for easy mocking in tests and loose coupling
type KeycloakClient interface {
	// Authentication methods
	AuthCodeURL(state string) string
	GetLogOutURL(tokenHint string) string
	
	// User management in Keycloak (auth only)
	CreateKeycloakUser(ctx context.Context, user *models.User, password string) (string, error)
	UpdateKeycloakUser(ctx context.Context, user *models.User) error
	UpdatePassword(ctx context.Context, keycloakUserID string, newPassword string) error
	DeleteKeycloakUser(ctx context.Context, userID string) error
}

// Ensure Client implements KeycloakClient
var _ KeycloakClient = (*Client)(nil)
