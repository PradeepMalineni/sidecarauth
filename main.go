// main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"sidecarauth/config"
	"sidecarauth/reverseproxy"
)

/*
func main() {
	// Initialize and use the auth sidecar proxy
	//authProxy := auth.NewAuthProxy()
	//authResult := authProxy.Authenticate("username", "password")
	fmt.Println("Authentication Result:", "hello2")
	fmt.Println("Authentication Result:", auth.X)

	// Initialize and use the cache module
	cacheInstance := cache.NewCache()
	cacheInstance.Set("key", "value")
	cachedValue, exists := cacheInstance.Get("key")
	if exists {
		fmt.Println("Cached Value:", cachedValue)
	} else {
		fmt.Println("Key not found in cache.")
	}
	fmt.Println("Starting the listening services for OAuth.")

	go_test.TestServices()

	//Start the Listner Function

	fmt.Println("Starting the Listner")
	listener.Listner()
}

*/

func main() {
	// Load configuration from file
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	//fmt.Println("Staring the test service on port 8443")
	//test_svcs.TestServices()

	// Create a reverse proxy that forwards requests to another service with a custom trust store
	proxy, err := reverseproxy.NewReverseProxy(cfg.TargetURL)
	if err != nil {
		log.Fatal("Error creating reverse proxy:", err)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: proxy,
	}

	fmt.Printf("Proxy server listening on %s...\n", server.Addr)
	err = server.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
	//select {}
}
