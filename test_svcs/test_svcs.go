package test_svcs

import (
	"fmt"
	"log"
	"net/http"
)

func TestServices() {
	// Route for authentication
	http.HandleFunc("/auth", authHandler)

	// Route for API
	http.HandleFunc("/api", apiHandler)

	// Specify the path to your SSL certificate and key files
	certFile := "/Users/siva/sidecarauth/certs/server.crt"
	keyFile := "/Users/siva/sidecarauth/certs/server.key"

	// Start the HTTPS server
	port1 := 8445
	fmt.Printf("Listening on :%d...\n", port1)
	err := http.ListenAndServeTLS(fmt.Sprintf(":%d", port1), certFile, keyFile, nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authentication handler")
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "API handler")
}
