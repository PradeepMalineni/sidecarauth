package listener

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

// Configuration is a struct to hold the configuration parameters.
type Configuration struct {
	TargetURL  string `json:"targetURL"`
	ServerPort int    `json:"serverPort"`
	CertFile   string `json:"certFile"`
	KeyFile    string `json:"keyFile"`
}

// ReverseProxy is a struct that implements the http.Handler interface
// and is used for forwarding requests to another HTTP service.
type ReverseProxy struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

func NewReverseProxy(targetURL string) *ReverseProxy {
	target, err := url.Parse(targetURL)
	if err != nil {
		panic(err)
	}

	return &ReverseProxy{
		target: target,
		proxy:  httputil.NewSingleHostReverseProxy(target),
	}
}

//Load the config file

func loadConfig(filename string) (*Configuration, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &Configuration{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (rp *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Modify the request if necessary before forwarding
	// For example, you might want to update headers or paths
	// r.Header.Add("X-Forwarded-Host", r.Host)
	log.Printf("Incoming Request: %s %s", r.Method, r.URL.Path)

	// Forward the request to the target service
	rp.proxy.ServeHTTP(w, r)
}

func isPortAvailable(port string) bool {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		// Port is not available
		return false
	}
	defer listener.Close()

	// Port is available
	return true
}

var port string = "8080"

func init() {

	if !isPortAvailable(port) {
		fmt.Printf("Error: Port %s is not available.\n", port)
		os.Exit(1) // Exit the program with error code 1
	}

	fmt.Printf("Port %s is available.\n", port)
}

func Listner() {
	// Start the HTTP server on port 8080
	config, err := loadConfig("listener/config.json")
	if err != nil {
		log.Fatal("Error loading configuration:", err)
		os.Exit(1) // Exit the program with error code 1
	}
	cert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
	if err != nil {
		log.Fatal(err)
	}
	// Create a TLS configuration with the loaded certificate and key
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	addr := ":" + port
	// Create a reverse proxy that forwards requests to another service
	targetURL := config.TargetURL // Replace with the target service URL
	proxy := NewReverseProxy(targetURL)

	// Create an HTTP server with the TLS configuration
	server := &http.Server{
		Addr: addr,
		//Handler:   http.HandlerFunc(handleRequest),
		TLSConfig: tlsConfig,
		Handler:   proxy,
	}

	// Start the HTTPS server
	fmt.Printf("Server listening on %s...\n", port)
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatal(err)
	}

}

/*
func handleRequest(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}
*/
