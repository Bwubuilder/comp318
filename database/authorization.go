// Or should userList map username to UserStruct then store all the user + token info in UserStruct?

package database

import (
	"encoding/json"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
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
	a := new(authHandler)
	a.tokenStore = make(map[string]string)
	return a
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
func (auth authHandler) handleAuthFunctions(w http.ResponseWriter, r *http.Request) {
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

func (auth *authHandler) authPost(w http.ResponseWriter, r *http.Request) {
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
	time.AfterFunc(1*time.Hour, func() { delete(auth.tokenStore, token) })

	auth.tokenStore[token] = d.Username // Store the token and other info
	slog.Info(auth.tokenStore[token])
	// Respond with the generated token
	response := marshalToken(token)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (auth *authHandler) authDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	if r.Header.Get("Authorization") == "" {
		w.Header().Add("WWW-Authenticate", "Bearer")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token := r.Header.Get("Authorization")[7:]

	if auth.tokenStore[token] == "" {
		http.Error(w, "No token", http.StatusUnauthorized)
	}
	// Get token from the Authorization header
	delete(auth.tokenStore, token)
	slog.Info(auth.tokenStore[token])
	w.WriteHeader(http.StatusNoContent)
}

// Packs the token up to be sent back to the user.
func marshalToken(token string) []byte {
	tokenVal := map[string]string{"token": token}

	response, err := json.MarshalIndent(tokenVal, "", "  ")
	if err != nil {
		slog.Info("Token marshaling failed")
		return nil
	}
	return response
}

func (auth *authHandler) checkToken(token string) bool {
	slog.Info(token)
	for k, v := range auth.tokenStore {
		slog.Info(k, v)
		if token[7:] == k {
			if v != "" {
				return true
			}
		}
	}
	return false
}

func (auth *authHandler) handleTokenFile(path string) {
	dat, err := os.ReadFile(path)
	if err != nil {
		slog.Info("Token file could not be read")
		return
	}

	var tokens map[string]string
	err = json.Unmarshal(dat, &tokens)
	if err != nil {
		slog.Info("Unmarshal Failed")
		return
	}

	for user, token := range tokens {
		auth.tokenStore[token] = user
	}
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
