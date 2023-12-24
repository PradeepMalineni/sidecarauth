package listener

import (
	"fmt"
	"net/http"
)

func main() {
	// Start the HTTP server on port 8080
	http.HandleFunc("/", handleRequest)
	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}
