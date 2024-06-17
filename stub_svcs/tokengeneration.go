package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type TokenResponse struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	IssuedAt    int64  `json:"issued_at"`
	ExpiresIn   int64  `json:"expires_in"`
	Scope       string `json:"scope"`
}

func generateTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Generate a simple token with a 1-hour expiration
	token := fmt.Sprintf("%d", rand.Intn(10000))

	// Calculate the current time and the expiration time
	now := time.Now().Unix()
	issuedAt := now
	expiresIn := int64(120) // 1 hour in seconds
	//expiresIn := now.Add(time.Second * time.Duration(expiresIn)).Unix()

	// Create a TokenResponse struct
	response := TokenResponse{
		TokenType:   "TestBearer",
		AccessToken: token,
		IssuedAt:    issuedAt,
		ExpiresIn:   expiresIn,
		Scope:       "read write", // Adjust scope as needed
	}

	// Convert the struct to JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func main() {
	// Define the /generate-token route
	http.HandleFunc("/generate-token", generateTokenHandler)

	// Start the server on port 8080
	fmt.Println("Token Generation Service is running on :8080")
	http.ListenAndServe("0.0.0.0:8080", nil)
}
