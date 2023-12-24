package listener

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

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
	cert, err := tls.LoadX509KeyPair("/Users/siva/sidecarauth/certs/server.crt", "/Users/siva/sidecarauth/certs/server.key")
	if err != nil {
		log.Fatal(err)
	}
	// Create a TLS configuration with the loaded certificate and key
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	addr := ":" + port

	// Create an HTTP server with the TLS configuration
	server := &http.Server{
		Addr:      addr,
		Handler:   http.HandlerFunc(handleRequest),
		TLSConfig: tlsConfig,
	}

	// Start the HTTPS server
	fmt.Printf("Server listening on %s...\n", port)
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatal(err)
	}

}

func handleRequest(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}
