package auth

import (
	"context"
	"crypto/tls"
	"log"
	"net/url"
	"passIt/internal/models"
	"passIt/internal/utils"

	// "crypto/rsa"

	// "crypto/tls"
	// "crypto/x509"

	// // "encoding/base64"
	// "encoding/json"
	// "encoding/pem"
	// "errors"
	"fmt"
	// "log"

	// "net/http"
	// "passIt/internal/models"
	// "sync"

	"github.com/Nerzal/gocloak/v13"
	"github.com/go-resty/resty/v2"

	// "github.com/golang-jwt/jwt/v5"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type Config struct {
	BaseURL       string // Authorization base url
	ClientID      string // client id oauth
	RedirectURL   string // valid redirect url
	ClientSecret  string // keycloak client secret
	Realm         string // keycloak realm
	AdminUsername string // keycloak admin username
	AdminPassword string // keycloak admin password
	FrontendURL   string // frontend URL for redirects
}

// Client struct holds all components needed for authentication
type Client struct {
	Client      *gocloak.GoCloak      // gocloak client for Keycloak admin operations
	Provider    *oidc.Provider        // Handles OIDC protocol operations with Keycloak
	Oauth       *oauth2.Config        // OAuth2 configuration for token exchange
	Keycloak    KeycloakClient        // Keycloak admin client
	FrontendURL string                // Frontend URL for redirects
	Config      *Config               // Store config for admin operations
}

func New(ctx context.Context, config *Config) (*Client, error) {
	// Create an insecure HTTP client (for self-signed certs only)
	// insecureHttpClient := &http.Client{
	// 	Timeout: 10 * time.Second,
	// 	Transport: &http.Transport{
	// 		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 👈 disables cert check
	// 	},
	// }

	// Use the custom HTTP client in the context
	insecureCtx := utils.InsecureHttpContext(ctx)
	// Construct the provider URL using Keycloak realm
	providerURL := fmt.Sprintf("%s/realms/%s", config.BaseURL, config.Realm)

	provider, err := oidc.NewProvider(insecureCtx, providerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %v", err)
	}

	// Configure an OpenID Connect aware OAuth2 client with specific scopes:
	// - oidc.ScopeOpenID: Required for OpenID Connect authentication, provides subject ID (sub)
	// - "roles": Keycloak-specific scope to get user roles in the token
	oauth2Config := oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes: []string{
			oidc.ScopeOpenID, // Required for OIDC authentication
			"roles",          // Request user roles from Keycloak
		},
	}

	client := gocloak.NewClient(config.BaseURL)
	restyClient := resty.New()
	restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetRestyClient(restyClient)

	// Return initialized client with all required components
	// Note: The returned Client implements KeycloakClient interface
	authClient := &Client{
		Client:   client,
		Config:   config,
		Oauth:    &oauth2Config,
		Provider: provider,
		FrontendURL: config.FrontendURL,
	}
	
	// Set Keycloak field to point to itself (implements KeycloakClient)
	authClient.Keycloak = authClient
	
	return authClient, nil
}

func (c *Client) GetLogOutURL(tokenHint string) string {
	authEndpoint := c.Oauth.Endpoint.AuthURL[:len(c.Oauth.Endpoint.AuthURL)-len("/auth")]
	log.Println(tokenHint)
	return fmt.Sprintf("%s/logout?id_token_hint=%s&post_logout_redirect_uri=%s", authEndpoint, tokenHint, url.QueryEscape(c.FrontendURL))
}

// AuthCodeURL generates the login URL for OAuth2 authorization code flow.
// It returns a URL that the user should be redirected to for authentication.
// The state parameter is a random string that will be validated in the callback
// to prevent CSRF attacks.
func (c *Client) AuthCodeURL(state string) string {
	return c.Oauth.AuthCodeURL(state)
}

func (c *Client) CreateKeycloakUser(ctx context.Context, user *models.User, password string) (string, error) {
	realm := c.Config.Realm

	// Admin login to Keycloak
	token, err := c.Client.LoginAdmin(
		ctx,
		c.Config.AdminUsername,
		c.Config.AdminPassword,
		realm,
	)
	if err != nil {
		return "", fmt.Errorf("keycloak admin login failed: %w", err)
	}

	// Prepare Keycloak user
	kcUser := gocloak.User{
		Username:  gocloak.StringP(user.Username),
		Email:     gocloak.StringP(user.Email),
		FirstName: gocloak.StringP(user.FirstName),
		LastName:  gocloak.StringP(user.LastName),
		Enabled:   gocloak.BoolP(true),
	}

	// Create user in Keycloak
	userID, err := c.Client.CreateUser(ctx, token.AccessToken, realm, kcUser)
	if err != nil {
		return "", fmt.Errorf("failed to create user in keycloak: %w", err)
	}

	// Set password for the user
	cred := gocloak.CredentialRepresentation{
		Type:      gocloak.StringP("password"),
		Value:     gocloak.StringP(password),
		Temporary: gocloak.BoolP(false),
	}
	err = c.Client.SetPassword(ctx, token.AccessToken, userID, realm, *cred.Value, false)
	if err != nil {
		return "", fmt.Errorf("failed to set password in keycloak: %w", err)
	}

	return userID, nil
}

func (c *Client) UpdateKeycloakUser(ctx context.Context, user *models.User) error {
	realm := c.Config.Realm

	// Admin login to Keycloak
	token, err := c.Client.LoginAdmin(
		ctx,
		c.Config.AdminUsername,
		c.Config.AdminPassword,
		realm,
	)
	if err != nil {
		return fmt.Errorf("keycloak admin login failed: %w", err)
	}
	// Prepare Keycloak user
	kcUser := gocloak.User{
		ID:        gocloak.StringP(user.KeycloackID),
		Username:  gocloak.StringP(user.Username),
		Email:     gocloak.StringP(user.Email),
		FirstName: gocloak.StringP(user.FirstName),
		LastName:  gocloak.StringP(user.LastName),
		Enabled:   gocloak.BoolP(user.IsActive),
	}
	// Update user in Keycloak
	err = c.Client.UpdateUser(ctx, token.AccessToken, realm, kcUser)
	if err != nil {
		return fmt.Errorf("failed to update user in keycloak: %w", err)
	}
	return nil
}

func (c *Client) DeleteKeycloakUser(ctx context.Context, userID string) error {
	realm := c.Config.Realm

	// Admin login to Keycloak
	token, err := c.Client.LoginAdmin(
		ctx,
		c.Config.AdminUsername,
		c.Config.AdminPassword,
		realm,
	)
	if err != nil {
		return fmt.Errorf("keycloak admin login failed: %w", err)
	}

	// Delete user from Keycloak
	err = c.Client.DeleteUser(ctx, token.AccessToken, realm, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user from keycloak: %w", err)
	}

	return nil
}

// UpdatePassword updates a user's password in Keycloak
func (c *Client) UpdatePassword(ctx context.Context, keycloakUserID string, newPassword string) error {
	realm := c.Config.Realm

	// Admin login to Keycloak
	token, err := c.Client.LoginAdmin(
		ctx,
		c.Config.AdminUsername,
		c.Config.AdminPassword,
		realm,
	)
	if err != nil {
		return fmt.Errorf("keycloak admin login failed: %w", err)
	}

	// Set new password for the user
	err = c.Client.SetPassword(ctx, token.AccessToken, keycloakUserID, realm, newPassword, false)
	if err != nil {
		return fmt.Errorf("failed to update password in keycloak: %w", err)
	}

	return nil
}
