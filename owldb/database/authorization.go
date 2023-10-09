// Or should userList map username to UserStruct then store all the user + token info in UserStruct?

package database

import (
	"encoding/json"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"time"
)

// Initialize a random number generator with a time-based seed
var seed = rand.New(rand.NewSource(time.Now().UnixNano()))

type authHandler struct {
	strlen     int
	charset    string
	tokenStore map[string]TokenInfo
}

func NewAuth() authHandler {
	var a authHandler
	// Define constants for token length and character set
	a.strlen = 15
	a.charset = "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz0123456789"
	// Map to store token information
	a.tokenStore = make(map[string]TokenInfo) // map token to TokenInfo struct (username + time)
	return authHandler{}
}

// Function to generate a random token
func (auth authHandler) makeToken() string {
	token := make([]byte, auth.strlen) // Initialize a byte array to hold the token
	for i := range token {
		token[i] = auth.charset[seed.Intn(len(auth.charset))] // Populate token with random characters from charset
	}
	slog.Info("Token made" + string(token))
	return string(token) // Convert byte array to string and return
}

// Struct to hold token information
type TokenInfo struct {
	Username string
	Created  time.Time
}

// HTTP handler function for authentication
func (auth authHandler) AuthFunction(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodOptions:
		// For the /auth endpoint, indicate that POST and DELETE are allowed.
		w.Header().Set("Allow", "POST, DELETE")
		w.WriteHeader(http.StatusOK)

	case http.MethodPost: // Handle POST method for user authentication
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var username string

		err2 := json.Unmarshal(body, &username)
		if err2 != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		// Get username from the query parameter
		if username == "" {
			http.Error(w, "Username is required", http.StatusBadRequest) // Return error if username is missing
			return
		}

		// ALSO NEED TO CHECK if user exists in the database here? or are all names valid?
		token := auth.makeToken()                                                   // Generate a new token
		auth.tokenStore[token] = TokenInfo{Username: username, Created: time.Now()} // Store the token and other info

		// Respond with the generated token
		response := marshalToken(token)
		w.Write(response)

	case http.MethodDelete: // Handle DELETE method for user de-authentication
		token := r.Header.Get("Authorization")[7:] // to get the token after "Bearer "
		// Get token from the Authorization header
		if token == "" {
			http.Error(w, "Token is required", http.StatusBadRequest) // Return error if token is missing
			return
		}
		if info, exists := auth.tokenStore[token]; exists { // Check if token exists
			if time.Since(info.Created).Hours() >= 1 { // Check token expiration
				delete(auth.tokenStore, token)                          // Remove expired token
				http.Error(w, "Token expired", http.StatusUnauthorized) // Return an expiration error

				return
			}
		} else {
			http.Error(w, "Invalid token", http.StatusUnauthorized) // Return an error for invalid token
			return
		}
		delete(auth.tokenStore, token) // Delete token if all checks pass

		w.Write([]byte("Logged out")) // Send logout confirmation

	default: // Handle unsupported HTTP methods
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func marshalToken(token string) []byte {
	tokenVal := map[string]string{"token": token}

	response, err := json.MarshalIndent(tokenVal, "", "  ")
	if err != nil {
		slog.Info("Token marshaling failed")
		return nil
	}
	return response
}

// need this case in NewHandler() in main.go
// http.HandleFunc("/auth", authorization.authHandler)  // Route /auth URL path to authHandler function if /auth in URL
// need to do OPTIONS ad well
// Use LOGGING
// need to check for token expiration each time for all incoming requests with the token in the header
// UserStruct with token and username
