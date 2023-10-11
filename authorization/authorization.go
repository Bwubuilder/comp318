// Or should userList map username to UserStruct then store all the user + token info in UserStruct?

package authorization

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

const charset = "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz0123456789"
const tokenLen = 15

type authHandler struct {
	tokenStore map[string]string
}

type UserFormat struct {
	Username string
}

func NewAuth() *authHandler {
	return new(authHandler)
}

// Function to generate a random token
func (auth authHandler) makeToken() string {
	token := make([]byte, tokenLen) // Initialize a byte array to hold the token
	for i := range token {
		token[i] = charset[seed.Intn(len(charset))] // Populate token with random characters from charset
	}
	slog.Info("Token made" + string(token))
	return string(token) // Convert byte array to string and return
}

// HTTP handler function for authentication
func (auth authHandler) HandleAuthFunctions(w http.ResponseWriter, r *http.Request) {
	slog.Info("Auth Method Called ", r.Method)
	slog.Info("Path ", r.URL.Path)
	logHeader(r)

	switch r.Method {
	case http.MethodPost: // Handle POST method for user authentication
		slog.Info("post request at /auth")
		auth.authPost(w, r)
		slog.Info("post finished")
	case http.MethodDelete: // Handle DELETE method for user de-authentication
		slog.Info("delete request at /auth")
		auth.authDelete(w, r)
		slog.Info("delete finished")
	case http.MethodOptions:
		slog.Info("auth requests options")
		auth.authOptions(w, r)
		slog.Info("options finished")
	default: // Handle unsupported HTTP methods
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (auth authHandler) authOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "POST,DELETE")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	slog.Info("Auth options header written")
	w.WriteHeader(http.StatusOK)
}

func (auth authHandler) authPost(w http.ResponseWriter, r *http.Request) {
	//Detect if content-type is application/json
	if r.Header.Get("Content-Type") != "" {
		content := r.Header.Get("Content-Type")
		if content != "application/json" {
			http.Error(w, "Content header not JSON", http.StatusUnsupportedMediaType)
			return
		}
	} else {
		slog.Info("Header contains no content type")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Info("Body could not be read")
		http.Error(w, `"invalid user format"`, http.StatusBadRequest)
		return
	}
	slog.Info("Read Body succeeded")
	r.Body.Close()

	var d UserFormat
	err = json.Unmarshal(body, &d)
	if err != nil {
		slog.Info("decode failed")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	slog.Info("Unmarshaled successfully")

	if d.Username == "" {
		slog.Info("No username")
		http.Error(w, "Username is required", http.StatusBadRequest) // Return error if username is missing
		return
	}

	slog.Info("Username exists")

	// ALSO NEED TO CHECK if user exists in the database here? or are all names valid?
	token := auth.makeToken() // Generate a new token

	auth.tokenStore[token] = d.Username // Store the token and other info
	// Respond with the generated token
	response := marshalToken(token)

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (auth authHandler) authDelete(w http.ResponseWriter, r *http.Request) {
	// Get token from the Authorization header
	token := r.Header.Get("Authorization")[7:] // to get the token after "Bearer "
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest) // Return error if token is missing
		return
	}
	if info, exists := auth.tokenStore[token]; exists { // Check if token exists
		delete(auth.tokenStore, info)
	} else {
		http.Error(w, "Invalid token", http.StatusUnauthorized) // Return an error for invalid token
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
