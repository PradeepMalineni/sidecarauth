// main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sidecarauth/auth"
	"sidecarauth/config"
	"sidecarauth/service"
	"syscall"
	"time"
)

func main() {
	currentTime := time.Now()
	timestampFormat := "2006-01-02"
	logDir := "logs"

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, 0755)
		if err != nil {
			log.Fatalf("[%s]: Error creating directory: %v", currentTime.Format(timestampFormat), err)
		}
		fmt.Printf("[%s]: Created directory: %s\n", currentTime.Format(timestampFormat), logDir)
	}

	logFileName := fmt.Sprintf("%s/app_%s.log", logDir, currentTime.Format(timestampFormat))
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("[%s]: Error opening log file: %v", currentTime.Format(timestampFormat), err)
	}
	defer logFile.Close()

	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.Print("SideCarAuthSvcs Started")

	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if configFilePath == "" {
		log.Fatal("CONFIG_FILE_PATH environment variable not set.")
	}

	config, err := config.LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("[%s]: Error loading configuration %s:", currentTime.Format(timestampFormat), err)
	}

	authHandlers := make(map[string]*auth.AuthHandler)
	for env, envConfig := range config.AuthConfig {
		authHandler := auth.NewAuthHandler(envConfig)
		authHandlers[env] = authHandler
	}
	log.Printf("[%s]: Authentication Listeners enabled", currentTime.Format(timestampFormat))

	http.HandleFunc(config.ListenerConfig.ListenerURI, func(w http.ResponseWriter, r *http.Request) {
		handleHTTP(authHandlers, config, w, r)
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := &http.Server{Addr: ":8080"}

	go func() {
		fmt.Printf("Go HTTP Listener is listening on port 8080...\n")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting HTTP server: %v", err)
		}
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

	<-sigint
	log.Println("Shutting down gracefully...")

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}
}

func handleHTTP(authHandlers map[string]*auth.AuthHandler, config config.Config, w http.ResponseWriter, r *http.Request) {
	_, port, err := net.SplitHostPort(r.Host)
	if err != nil {
		log.Printf("[%s]: Error extracting port from host: %v\n", time.Now().Format("2006-01-02"), err)
		return
	}

	var env string
	for e, p := range config.ListenerConfig.PortNumber {
		if p == port {
			env = e
			break
		}
	}

	if env == "" {
		log.Printf("[%s]: No environment found for port: %s\n", time.Now().Format("2006-01-02"), port)
		return
	}

	authHandlers[env].Initialize()

	tokenResponse, err := newFunction(authHandlers, env)
	if err != nil {
		http.Error(w, "Error getting access token", http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(tokenResponse)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	httpMethod := r.Method
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "NA"
	}
	uri := r.URL.Path
	backendURL := config.ServiceConfig[env].ApiURL + uri
	accessToken := "Bearer " + tokenResponse.AccessToken
	formattedResponse, err := service.MakeRequest(context.Background(), backendURL, accessToken, httpMethod, contentType, string(payload), r.Header)
	if err != nil {
		log.Println("Error making request:", err)
		return
	}
	fmt.Fprintf(w, "\n\nFormatted Response: %s", formattedResponse)
}

func newFunction(authHandlers map[string]*auth.AuthHandler, env string) (auth.TokenResponse, error) {
	tokenResponse, err := authHandlers[env].GetAccessToken()
	return tokenResponse, err
}
