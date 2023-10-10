// Or should userList map username to UserStruct then store all the user + token info in UserStruct?

package authorization

import (
	"encoding/json"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// Initialize a random number generator with a time-based seed
var seed = rand.New(rand.NewSource(time.Now().UnixNano()))

type authHandler struct {
	mu         sync.Mutex
	strlen     int
	charset    string
	tokenStore map[string]TokenInfo
}

func NewAuth() *authHandler {
	var a authHandler
	a.strlen = 15
	a.charset = "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz0123456789"
	a.tokenStore = make(map[string]TokenInfo) // map token to TokenInfo struct (username + time)
	return &a
}

// Function to generate a random token
func (auth *authHandler) makeToken() string {
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

func newTokenInfo() TokenInfo {
	var info TokenInfo
	info.Created = time.Now()
	return info
}

// HTTP handler function for authentication
func (auth *authHandler) HandleAuthFunctions(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		auth.authOptions(w, r)
	}
	slog.Info("Auth Method Called" + r.Method)
	slog.Info("Path" + r.URL.Path)
	logHeader(r)
	switch r.Method {
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

func (auth *authHandler) authOptions(w http.ResponseWriter, r *http.Request) {
	auth.mu.Lock()
	defer auth.mu.Unlock()
	w.Header().Set("Allow", "POST,DELETE")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, DELETE")
	w.WriteHeader(http.StatusOK)
	slog.Info("Auth options header written")
}

type userStore struct {
	username string
}

func (auth *authHandler) authPost(w http.ResponseWriter, r *http.Request) {
	//Detect if content-type is application/json
	if r.Header.Get("Content-Type") != "" {
		content := r.Header.Get("Content-Type")
		if content != "application/json" {
			http.Error(w, "Content header not JSON", http.StatusUnsupportedMediaType)
			return
		}
	}

	slog.Info("Making it further...")
	auth.mu.Lock()
	defer auth.mu.Unlock()
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, `"invalid user format"`, http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		http.Error(w, "Body of Request is empty", http.StatusBadRequest)
		return
	}

	var store userStore
	err2 := json.NewDecoder(r.Body).Decode(&store)
	if err2 != nil {
		slog.Info("decode failed")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get username from the query parameter
	if store.username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest) // Return error if username is missing
		return
	}

	slog.Info("username successful" + store.username)

	// ALSO NEED TO CHECK if user exists in the database here? or are all names valid?
	token := auth.makeToken() // Generate a new token

	thisToken := newTokenInfo()
	auth.tokenStore[token] = thisToken // Store the token and other info

	// Respond with the generated token
	response := marshalToken(token)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (auth *authHandler) authDelete(w http.ResponseWriter, r *http.Request) {
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
	auth.mu.Lock()
	defer auth.mu.Unlock()
	delete(auth.tokenStore, token) // Delete token if all checks pass

	w.Write([]byte("Logged out")) // Send logout confirmation
	return
}

func marshalToken(token string) []byte {
	slog.Info("We made it this far!")
	tokenVal := map[string]string{"token": token}

	response, err := json.MarshalIndent(tokenVal, "", "  ")
	if err != nil {
		slog.Info("Token marshaling failed")
		return nil
	}
	return response
}

func logHeader(r *http.Request) {
	for key, element := range r.Header {
		slog.Info("Header:", key, "Value", element)
	}
}

// need this case in NewHandler() in main.go
// http.HandleFunc("/auth", authorization.authHandler)  // Route /auth URL path to authHandler function if /auth in URL
// need to do OPTIONS ad well
// Use LOGGING
// need to check for token expiration each time for all incoming requests with the token in the header
// UserStruct with token and username
