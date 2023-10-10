// Or should userList map username to UserStruct then store all the user + token info in UserStruct?

package database

import (
	"encoding/json"
	"fmt"
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

func newAuth() authHandler {
	var a authHandler
	// Define constants for token length and character set
	a.strlen = 15
	a.charset = "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz0123456789"
	// Map to store token information
	a.tokenStore = make(map[string]TokenInfo) // map token to TokenInfo struct (username + time)
	slog.Info("auth created")
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
func (auth authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("Hey, we made it this far..." + r.Method)
	switch r.Method {
	case http.MethodOptions:
		slog.Info("auth requests options")
		// For the /auth endpoint, indicate that POST and DELETE are allowed.
		auth.authPost(w, r)
		slog.Info("auth finished options")
	case http.MethodPost: // Handle POST method for user authentication
		slog.Info("auth requests post")
		auth.authPost(w, r)
		slog.Info("auth finished post")
	case http.MethodDelete: // Handle DELETE method for user de-authentication
		slog.Info("auth requests delete")
		auth.authDelete(w, r)
		slog.Info("auth finished delete")
	default: // Handle unsupported HTTP methods
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (auth authHandler) authOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "POST,DELETE")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, DELETE")
	w.WriteHeader(http.StatusOK)
}

func (auth authHandler) authPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slog.Info("Making it further...")

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, `"invalid user format"`, http.StatusBadRequest)
		return
	}

	slog.Info("body set" + fmt.Sprint(body))

	var thisToken TokenInfo
	err2 := json.Unmarshal(body, &thisToken)
	if err2 != nil {
		slog.Info("unmarshal failed")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get username from the query parameter
	if thisToken.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest) // Return error if username is missing
		return
	}

	slog.Info("username successful" + thisToken.Username)

	// ALSO NEED TO CHECK if user exists in the database here? or are all names valid?
	token := auth.makeToken() // Generate a new token
	thisToken.Created = time.Now()
	auth.tokenStore[token] = thisToken // Store the token and other info

	// Respond with the generated token
	response := marshalToken(token)
	w.Write(response)
}

func (auth authHandler) authDelete(w http.ResponseWriter, r *http.Request) {
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
	return
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
