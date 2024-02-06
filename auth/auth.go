// auth/auth.go
package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"sidecarauth/config"
	logger "sidecarauth/utility"
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
func NewAuthHandler(env string, envConfig config.AuthConfig) *AuthHandler {
	logger.LogF("Auth Module :  Authentication Handler initailizing:", env)
	return &AuthHandler{
		config: envConfig,
	}
}

// Initialize is called with the configuration values
func (a *AuthHandler) Initialize() error {
	// Call the function to get the initial access token
	logger.Log("Auth Module : getAccessToken func call")
	err := a.getAccessToken()
	if err != nil {
		logger.LogF("Auth Module get Acesstoken Initialize error", err)
		return err
	}
	return nil

}

// GetAccessToken function to get the access token
func (a *AuthHandler) GetAccessToken() (TokenResponse, error) {
	// Lock to ensure thread-safe access
	a.mu.Lock()
	defer a.mu.Unlock()
	logger.Log("Auth Module : GetAccessToken func call")

	// Check if the token is expired or about to expire
	now := time.Now().Unix()
	if a.tokenResponse.AccessToken == "" || now >= a.tokenResponse.ExpiresIn+a.tokenResponse.IssuedAt {
		// Token is expired or about to expire, refresh it
		logger.Log("Auth Module : GetAccessToken func call2")

		err := a.getAccessToken()
		if err != nil {
			logger.LogF("Auth Module get Acesstoken error", err)
			return TokenResponse{}, err
		}
	}
	return a.tokenResponse, nil
}

func (a *AuthHandler) getAccessToken() error {
	// Lock to ensure thread-safe access
	a.mu.Lock()
	defer a.mu.Unlock()
	logger.LogF("Token data", a.tokenResponse)
	if !isEmptyStruct(a.tokenResponse) {
		logger.Log("Auth Module : No token found")
		now := time.Now().Unix()
		if now < a.tokenResponse.ExpiresIn+a.tokenResponse.IssuedAt-60 { // 60 seconds before expiration
			// Token is not close to expiration, no need to refresh
			return nil
		}

	}
	// Check if the token is expired or about to expire
	logger.Log("Auth Module : New token request initiated")

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}

	// Use configURL and authHeader as needed
	payload := strings.NewReader("grant_type=client_credentials")

	req, err := http.NewRequest("POST", a.config.TokenURL, payload)
	if err != nil {
		return err
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", a.config.AuthorizationHeader)

	res, err := client.Do(req)
	if err != nil {
		//http.Error( "Backend service is unavailable", http.StatusInternalServerError)
		logger.LogF("Auth Module : Error performing HTTP request to IDP service:", err)
		return err
	}
	// Close the response body to ensure proper connection closure

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.LogF("Error performing HTTP request:", err)
		return err
	}
	if res.StatusCode != http.StatusOK {
		a.handleError(body)
		return err
	}
	err = json.Unmarshal(body, &a.tokenResponse)
	if err != nil {
		logger.LogF("Error unmarshalling JSON:", err)
		return err
	}

	return nil

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

func isEmptyStruct(s interface{}) bool {
	zeroValue := reflect.Zero(reflect.TypeOf(s))
	return reflect.DeepEqual(s, zeroValue.Interface())
}
