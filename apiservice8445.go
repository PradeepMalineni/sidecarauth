package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Route for API
	http.HandleFunc("/api", apiHandler)

	// Specify the path to your SSL certificate and key files
	//certFile := "/Users/bhumi/goapp/server.cer"
	//keyFile := "/Users/bhumi/goapp/server.key"

	// Start the HTTPS server
	port1 := 8445
	fmt.Printf("Listening on :%d...\n", port1)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port1), nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	//response := fmt.Sprintf(`{"authtoken": "%s", "ResponseBody": "API handler"}`, authToken)

	IncomingHTTPMethod := r.Method

	response := fmt.Sprintf(`{"authtoken": "%s", "ResponseBody": "API handler","Method":"%s","port":"8445"}`, authToken, IncomingHTTPMethod)
	fmt.Fprintln(w, response)

}
