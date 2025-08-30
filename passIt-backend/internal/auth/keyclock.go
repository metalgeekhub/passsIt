package auth

import (
	"context"
	"crypto/tls"
	"os"
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
	BaseURL      string // Authorization base url
	ClientID     string // client id oauth
	RedirectURL  string // valid redirect url
	ClientSecret string // keycloak client secret
	Realm        string // keycloak realm
}

// Client struct holds all components needed for authentication
type Client struct {
	Provider *oidc.Provider        // Handles OIDC protocol operations with Keycloak
	OIDC     *oidc.IDTokenVerifier // Verifies JWT tokens from Keycloak
	Oauth    oauth2.Config         // Manages OAuth2 flow (authorization codes, tokens)
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

	// Create ID token verifier
	verifier := provider.Verifier(&oidc.Config{
		ClientID: config.ClientID,
	})

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

	// Return initialized client with all required components
	return &Client{
		// oauth2Config: Used for OAuth2 operations like:
		// - Generating login URL (AuthCodeURL)
		// - Exchanging auth code for tokens (Exchange)
		// - Managing token refresh
		Oauth: oauth2Config,

		// verifier: Used to validate tokens:
		// - Verifies JWT signature
		// - Validates token claims (exp, iss, aud)
		// - Extracts user information
		OIDC: verifier,

		// provider: Keycloak OIDC provider that:
		// - Provides endpoint URLs (auth, token)
		// - Handles OIDC protocol details
		// - Manages provider metadata
		Provider: provider,
	}, nil
}

// AuthCodeURL generates the login URL for OAuth2 authorization code flow.
// It returns a URL that the user should be redirected to for authentication.
// The state parameter is a random string that will be validated in the callback
// to prevent CSRF attacks.
func (c *Client) AuthCodeURL(state string) string {
	return c.Oauth.AuthCodeURL(state)
}

func (c *Client) CreateKeycloakUser(ctx context.Context, user *models.User, password string) (string, error) {
	realm := os.Getenv("KEYCLOAK_REALM")
	// Initialize gocloak client
	client := gocloak.NewClient(os.Getenv("KEYCLOAK_URL"))
	restyClient := resty.New()
	restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetRestyClient(restyClient)

	// Admin login to Keycloak
	token, err := client.LoginAdmin(
		ctx,
		os.Getenv("KEYCLOAK_ADMIN_USERNAME"),
		os.Getenv("KEYCLOAK_ADMIN_PASSWORD"),
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
	userID, err := client.CreateUser(ctx, token.AccessToken, realm, kcUser)
	if err != nil {
		return "", fmt.Errorf("failed to create user in keycloak: %w", err)
	}

	// Set password for the user
	cred := gocloak.CredentialRepresentation{
		Type:      gocloak.StringP("password"),
		Value:     gocloak.StringP(password),
		Temporary: gocloak.BoolP(false),
	}
	err = client.SetPassword(ctx, token.AccessToken, userID, realm, *cred.Value, false)
	if err != nil {
		return "", fmt.Errorf("failed to set password in keycloak: %w", err)
	}

	return userID, nil
}
