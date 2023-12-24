package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type TokenResponse struct {
	token_type   string `json:"token_type"`
	access_token string `json:"access_token"`
	issued_at    int64  `json:"issued_at"`
	expires_in   int64  `json:"expires_in"`
	scope        string `json:"scope"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {

	url := "https://apiidp-enterprise1-sandbox.wellsfargo.com/oauth/token"

	payload := strings.NewReader("grant_type=client_credentials")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic 123aqbc")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	if res.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected status code: %d\n", res.StatusCode)
		fmt.Println(string(body))
		handleError(body)
		return
	}
	var tokenResponse TokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	accessToken := tokenResponse.access_token
	expiresIn := tokenResponse.expires_in

	// Use the variables as needed
	fmt.Println("Access Token:", accessToken)
	fmt.Println("Expires In:", expiresIn)
	fmt.Println(res)
	fmt.Println(string(body))

}
func handleError(body []byte) {
	// Attempt to unmarshal the response into an ErrorResponse struct
	var errorResponse ErrorResponse
	err := json.Unmarshal(body, &errorResponse)
	if err == nil && errorResponse.Error != "" {
		// If successful, print the error message
		fmt.Println("Error response from the server:", errorResponse.Error)
	} else {
		// If unmarshaling into ErrorResponse fails, or the error field is not present, print a generic error message
		fmt.Println("Unexpected response from the server")

	}
}
