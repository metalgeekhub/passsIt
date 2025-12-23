package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"passIt/internal/auth"
	"passIt/internal/constant"
	"passIt/internal/database"
	"passIt/internal/models"
	"passIt/internal/services"
	"passIt/internal/store"
	"passIt/internal/utils"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	authClient   *auth.Client
	sessionStore store.SessionStore
	userService  services.UserService
	frontendURL  string
}

func NewAuthHandler(authClient *auth.Client, sessionStore store.SessionStore, dbService database.Service, frontendURL string) *AuthHandler {
	return &AuthHandler{
		authClient:   authClient,
		sessionStore: sessionStore,
		userService:  services.NewUserService(dbService, authClient),
		frontendURL:  frontendURL,
	}
}

// generateRandomSecureString creates a random secure string
func generateRandomSecureString() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (h *AuthHandler) ShowLoginPage(c *gin.Context) {
	// You can pass data as a map or struct if needed, or nil if not
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
	})
}

// LoginHandler initiates the OAuth2 authorization code flow with Keycloak.
// It generates a secure state parameter to prevent CSRF attacks and stores it
// in a secure cookie for later verification during the callback phase.
//
// Returns:
// - 302: Redirects to Keycloak login page
// - 500: Internal Server Error if state generation fails
// LoginHandler godoc
// @Summary      Initiate OAuth2 login
// @Description  Redirects to Keycloak for authentication
// @Tags         auth
// @Success      302 {string} string "Redirect to Keycloak"
// @Failure      500 {object} map[string]string
// @Router       /auth/login [get]
func (a *AuthHandler) LoginHandler(c *gin.Context) {
	state, err := generateRandomSecureString()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state"})
		return
	}

	// Store state in secure, signed cookie (simpler than Redis for state)
	// Cookie expires in 5 minutes - enough time for OAuth flow
	// Use SameSiteLaxMode for OAuth redirects (Strict blocks external redirects)
	// Use secure=false for localhost development (no HTTPS)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"oauth_state",  // name
		state,          // value
		300,            // maxAge: 5 minutes
		"/",            // path
		"",             // domain
		false,          // secure: false for localhost (set to true in production with HTTPS)
		true,           // httpOnly
	)

	// Build authentication URL
	authURL := a.authClient.Oauth.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("response_type", "code"),
		oauth2.SetAuthURLParam("scope", "openid profile email"),
	)

	// Redirect to Keycloak login page
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// LogoutHandler godoc
// @Summary      Logout user
// @Description  Logs the user out by deleting the session and clearing the cookie
// @Tags         auth
// @Success      302 {string} string "Redirect to Keycloak logout"
// @Router       /auth/logout [get]
func (a *AuthHandler) LogoutHandler(c *gin.Context) {
	sessionID, err := c.Cookie("session_id")
	log.Printf("Session_id=%s", sessionID)

	var sessionData *store.SessionData
	if err == nil && sessionID != "" {
		// Get session data before deleting
		sd, getErr := a.sessionStore.Get(c, sessionID)
		if getErr == nil {
			sessionData = sd
		}
		// Now delete session from store
		_ = a.sessionStore.Delete(c, sessionID)
	}

	// Clear the session_id cookie
	c.SetCookie(
		"session_id",
		"",
		-1, // MaxAge negative deletes the cookie
		"/",
		"",
		false, // secure: false for localhost
		true, // httpOnly
	)

	// Optionally, redirect to Keycloak's logout endpoint to log out globally:
	logoutURL := a.authClient.GetLogOutURL(sessionData.IDToken)
	c.Redirect(http.StatusTemporaryRedirect, logoutURL)
}

// CallbackHandler godoc
// @Summary      OAuth2 callback
// @Description  Handles the OAuth2 callback from Keycloak after authentication
// @Tags         auth
// @Param        code query string true "Authorization code"
// @Param        state query string true "State parameter"
// @Success      302 {string} string "Redirect to frontend"
// @Failure      500 {object} map[string]string
// @Router       /auth/callback [get]
func (a *AuthHandler) CallbackHandler(c *gin.Context) {
	// Check for error parameter from OAuth provider (e.g., user denied access)
		if errorParam := c.Query("error"); errorParam != "" {
		errorDesc := c.Query("error_description")
		log.Printf("OAuth error: %s - %s", errorParam, errorDesc)
		// Redirect to frontend with error
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=%s", a.frontendURL, errorParam))
		return
	}

	if err := a.validateStateSession(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate state session"})
		log.Printf("State validation error: %v", err)
		return
	}
	oauthToken, err := a.tokenExchange(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		log.Printf("Token exchange error: %v", err)
		return
	}
	userInfo, tokenID, err := a.validateAndGetClaimsIDToken(c, oauthToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate and get claims id token"})
		log.Printf("ID token validation error: %v", err)
		return
	}
	
	// Fetch user from database to get admin status
	dbUser, err := a.userService.GetUserByEmail(c, userInfo.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user from database"})
		log.Printf("Failed to fetch user: %v", err)
		return
	}
	
	sessionID, err := generateRandomSecureString()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate session ID"})
		log.Printf("Session ID generation error: %v", err)
		return
	}

	// rawIDToken, ok := oauthToken.Extra("id_token").(string)
	// if !ok {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get token ID"})
	// 	log.Printf("Failed to get token ID")
	// 	return
	// }
	// Create session data with admin status
	sessionData := store.SessionData{
		AccessToken: oauthToken.AccessToken, // From Keycloak
		IDToken:     tokenID,
		UserInfo: store.UserInfo{
			Username: userInfo.Username,
			Email:    userInfo.Email,
			IsAdmin:  dbUser.IsAdmin,
		},
		CreatedAt: time.Now(),
	}
	// Store session
	if err := a.sessionStore.Set(c, sessionID, sessionData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store session"})
		return
	}
	// Set SameSite attribute
	// Note: Gin handles SameSite through the Config struct
	c.SetSameSite(http.SameSiteLaxMode)
	// Set secure session cookie using Gin's methods
	c.SetCookie(
		"session_id",                            // name
		sessionID,                               // value
		int(constant.SessionDuration.Seconds()), // maxAge in seconds
		"/",                                     // path
		"",                                      // domain (empty means default to current domain)
		false,                                   // secure: false for localhost (set to true in production with HTTPS)
		true,                                    // httpOnly (prevents JavaScript access)
	)
	c.Redirect(http.StatusTemporaryRedirect, a.frontendURL)
}

func (a *AuthHandler) tokenExchange(c *gin.Context) (*oauth2.Token, error) {
	insecureCtx := utils.InsecureHttpContext(c)

	authorizationCode := c.Query("code")
	if authorizationCode == "" {
		return nil, errors.New("authorizationCode is required")
	}
	opts := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("grant_type", "authorization_code"),
	}
	oauth2Token, err := a.authClient.Oauth.Exchange(insecureCtx, authorizationCode, opts...)
	if err != nil {
		return nil, err
	}
	return oauth2Token, nil
}

type oidcClaims struct {
	Email    string `json:"email"`
	Username string `json:"preferred_username"`
}

// ValidateIDToken verifies the id token from the oauth2token
func (a *AuthHandler) validateAndGetClaimsIDToken(c *gin.Context, oauth2Token *oauth2.Token) (*oidcClaims, string, error) {
	insecureCtx := utils.InsecureHttpContext(c)

	// Get and validate the ID token - this proves the user's identity
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	log.Println("The rawIDToken:", rawIDToken)
	if !ok {
		return nil, "", errors.New("no ID token found")
	}
	// Verify the ID token using the OIDC provider's verifier
	verifier := a.authClient.Provider.Verifier(&oidc.Config{
		ClientID: a.authClient.Config.ClientID,
	})
	idToken, err := verifier.Verify(insecureCtx, rawIDToken)
	if err != nil {
		return nil, "", errors.New("failed to verify id token")
	}
	claims := oidcClaims{}
	if err := idToken.Claims(&claims); err != nil {
		return nil, "", errors.New("failed to get user info")
	}
	return &claims, rawIDToken, nil
}

func (a *AuthHandler) validateStateSession(c *gin.Context) error {
	// Get state from callback parameters
	stateParam := c.Query("state")
	if stateParam == "" {
		return errors.New("missing state parameter in callback")
	}

	// Retrieve stored state from cookie
	storedState, err := c.Cookie("oauth_state")
	if err != nil {
		return fmt.Errorf("failed to retrieve stored state: %w", err)
	}

	// Validate state match
	if storedState != stateParam {
		return errors.New("state parameter mismatch")
	}

	// Clean up used state cookie
	c.SetCookie(
		"oauth_state",
		"",
		-1, // Delete cookie
		"/",
		"",
		false, // secure: false for localhost
		true,
	)

	return nil
}

// SignupRequest represents the request body for public signup
type SignupRequest struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// SignupHandler godoc
// @Summary      Register new user
// @Description  Allows public users to self-register (always as non-admin)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user body SignupRequest true "User signup data"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /auth/signup [post]
func (a *AuthHandler) SignupHandler(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create user data - always set IsAdmin to false for public signups
	user := &models.User{
		Username:  req.Username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsAdmin:   false, // Public signups are never admin
		IsActive:  true,
	}

	err := a.userService.CreateUser(c, user, req.Password)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create user: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully. Please login.",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}
