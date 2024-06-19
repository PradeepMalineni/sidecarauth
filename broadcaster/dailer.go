package broadcast

import (
	"encoding/json"
	"log"
	"os"
	logger "sidecarauth/utility"

	"google.golang.org/grpc"
)

type Config struct {
	Application struct {
		Name    string `json:"name"`
		Servers []struct {
			IP string `json:"ip"`
		} `json:"servers"`
	} `json:"application"`
}

func Dailer() {
	// Load configuration
	logger.Log("GRPC Connection Dail")
	configFilePath := "/Users/siva/sidecar_auth_05_2024/sidecarauth/broadcaster/broadcaster-config.json"

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	// Start the server and clients
	servers := config.Application.Servers
	if len(servers) == 0 {
		log.Fatal("No servers configured")
	}

	primaryIP := servers[0].IP
	var server *grpc.Server

	if CheckServer(primaryIP) {
		server = startServer(primaryIP)
	} else {
		log.Printf("Primary server %s not available", primaryIP)
	}

	for _, srv := range servers[1:] {
		go startClient(srv.IP, primaryIP)
	}

	if server != nil {
		select {}
	} else {
		log.Fatal("No server started")
	}
}
