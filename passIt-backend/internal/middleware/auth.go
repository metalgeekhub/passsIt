package middleware

import (
	"context"
	"log"
	"net/http"
	"passIt/internal/auth"
	"passIt/internal/database"
	"passIt/internal/store"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authClient   *auth.Client
	sessionStore store.SessionStore
	clientID     string
	dbService    database.Service
}

// NewAuthMiddleware creates a new authentication middleware with OIDC verification
func NewAuthMiddleware(c context.Context, authClient *auth.Client, sessionStore store.SessionStore, dbService database.Service) *AuthMiddleware {
	return &AuthMiddleware{
		authClient:   authClient,
		sessionStore: sessionStore,
		dbService:    dbService,
	}
}
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var accessToken string
		var authType string

		// Check for Authorization header first (for API clients)
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			accessToken = authHeader[7:]
			authType = "bearer"
		} else {
			// Fall back to session cookie (for browser clients)
			sessionID, err := c.Cookie("session_id")
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - no valid session or token"})
				c.Abort()
				return
			}
			// Get session data from Redis
			sessionData, err := m.sessionStore.Get(c, sessionID)
			if err != nil {
				// Clear invalid session cookie
				c.SetCookie("session_id", "", -1, "/", "", true, true)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - invalid session"})
				c.Abort()
				return
			}
			accessToken = sessionData.AccessToken
			authType = "session"
			c.Set("user_session", sessionData)
		}

		// Verify the access token using the OIDC provider
		token, err := m.authClient.Provider.Verifier(&oidc.Config{
			SkipClientIDCheck: true, // Access tokens don't require client ID check
		}).Verify(c, accessToken)

		if err != nil {
			// The token is invalid
			if authType == "session" {
				// Clean up session for browser clients
				sessionID, _ := c.Cookie("session_id")
				m.sessionStore.Delete(c, sessionID)
				c.SetCookie("session_id", "", -1, "/", "", true, true)
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - invalid token"})
			c.Abort()
			return
		}

		// Extract claims from the token
		var claims map[string]interface{}
		if err := token.Claims(&claims); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// For Bearer tokens, fetch user from database to get admin status
		if authType == "bearer" {
			// Get email from claims
			email, ok := claims["email"].(string)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - email not found in token"})
				c.Abort()
				return
			}

			// Fetch user from database
			user, err := m.dbService.FindUserByEmail(email)
			if err != nil {
				log.Printf("Failed to fetch user for Bearer token: %v", err)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - user not found"})
				c.Abort()
				return
			}

			// Create session data for Bearer token authentication
			sessionData := &store.SessionData{
				AccessToken: accessToken,
				UserInfo: store.UserInfo{
					Username: user.Username,
					Email:    user.Email,
					IsAdmin:  user.IsAdmin,
				},
			}
			c.Set("user_session", sessionData)
		}

		// Store the validated claims and auth type in the context
		c.Set("user_claims", claims)
		c.Set("auth_type", authType)
		c.Next()
	}
}

// RequireAdmin middleware ensures the user is an admin
func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get session data (set by RequireAuth middleware)
		sessionData, exists := c.Get("user_session")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden - admin access required"})
			c.Abort()
			return
		}

		session, ok := sessionData.(*store.SessionData)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid session data"})
			c.Abort()
			return
		}

		if !session.UserInfo.IsAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden - admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
