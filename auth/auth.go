// auth/auth.go
package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sidecarauth/config"
	"strings"
	"sync"
	"time"
)

// TokenResponse holds the authentication token information
type TokenResponse struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	IssuedAt    int64  `json:"issued_at"`
	ExpiresIn   int64  `json:"expires_in"`
	Scope       string `json:"scope"`
}

// ErrorResponse holds information about an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// AuthHandler holds the state of the authentication process
type AuthHandler struct {
	tokenResponse TokenResponse
	config        config.AuthConfig
	mu            sync.Mutex // Mutex for thread-safe operations
}

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler(envConfig config.AuthConfig) *AuthHandler {
	return &AuthHandler{
		config: envConfig,
	}
}

// Initialize is called with the configuration values
func (a *AuthHandler) Initialize() {
	// Call the function to get the initial access token
	a.getAccessToken()
}

// GetAccessToken function to get the access token
func (a *AuthHandler) GetAccessToken() (TokenResponse, error) {
	// Lock to ensure thread-safe access
	a.mu.Lock()
	defer a.mu.Unlock()

	// Check if the token is expired or about to expire
	now := time.Now().Unix()
	if a.tokenResponse.AccessToken == "" || now >= a.tokenResponse.ExpiresIn+a.tokenResponse.IssuedAt {
		// Token is expired or about to expire, refresh it
		a.getAccessToken()
	}

	return a.tokenResponse, nil
}

func (a *AuthHandler) getAccessToken() {
	// Lock to ensure thread-safe access
	a.mu.Lock()
	defer a.mu.Unlock()

	// Check if the token is expired or about to expire
	now := time.Now().Unix()
	if now < a.tokenResponse.ExpiresIn+a.tokenResponse.IssuedAt-60 { // 60 seconds before expiration
		// Token is not close to expiration, no need to refresh
		return
	}

	// Use configURL and authHeader as needed
	payload := strings.NewReader("grant_type=client_credentials")

	req, _ := http.NewRequest("POST", a.config.TokenURL, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", a.config.AuthorizationHeader)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error performing HTTP request:", err)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	if err != nil {
		fmt.Println("Error performing HTTP request:", err)
		return
	}
	if res.StatusCode != http.StatusOK {
		fmt.Println("Received non-OK status:", res.Status)
		a.handleError(body)
		return
	}
	err = json.Unmarshal(body, &a.tokenResponse)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}
}

func (a *AuthHandler) handleError(body []byte) {
	var errorResponse ErrorResponse
	err := json.Unmarshal(body, &errorResponse)
	if err == nil && errorResponse.Error != "" {
		fmt.Println("Error response from the server:", errorResponse.Error)
	} else {
		fmt.Println("Unexpected response from the server")
	}
}
