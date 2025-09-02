package handlers

import (
	"crypto/rand"
	"encoding/base64"
)

// generateRandomSecureString creates a random secure string
func generateRandomSecureString() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// LoginHandler initiates the OAuth2 authorization code flow with Keycloak.
// It generates a secure state parameter to prevent CSRF attacks and stores it
// in Redis for later verification during the callback phase.
//
// Returns:
// - 302: Redirects to Keycloak login page
// - 500: Internal Server Error if state generation or storage fails
