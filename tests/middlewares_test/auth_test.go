package middlewares_test

import (
	"hardwaremonitoringexporter/middlewares"
	"hardwaremonitoringexporter/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)


func setupRouterWithAuth() *http.ServeMux {
	mux := http.NewServeMux()
	
	mux.Handle("/protected", middlewares.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract userID  from context and write it to the response
		userID := r.Context().Value("userID").(uint)
		w.Write([]byte("User ID: " + string(userID)))
	})))
	return mux
}

// TestAuthMiddleware_ValidToken tests the middleware with a valid JWT token
func TestAuthMiddleware_ValidToken(t *testing.T) {
	// Create a valid JWT token with userID 
	claims := jwt.MapClaims{"userID": 1}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(services.JwtKey)

	// Prepare a new HTTP request with the token in the Authorization header
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	res := httptest.NewRecorder() // Create a response recorder

	// Set up the router and serve the HTTP request
	router := setupRouterWithAuth()
	router.ServeHTTP(res, req)

	// Assert that the response status code is OK and the body contains the expected user ID
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "User ID: 1", res.Body.String())
}

// TestAuthMiddleware_NoToken tests the middleware when no token is provided
func TestAuthMiddleware_NoToken(t *testing.T) {
	req, _ := http.NewRequest("GET", "/protected", nil) // Create a request without the Authorization header
	res := httptest.NewRecorder()

	// Serve the request using the router
	router := setupRouterWithAuth()
	router.ServeHTTP(res, req)

	// Assert that the response status code is Unauthorized
	assert.Equal(t, http.StatusUnauthorized, res.Code)
}

// TestAuthMiddleware_InvalidToken tests the middleware with an invalid JWT token
func TestAuthMiddleware_InvalidToken(t *testing.T) {
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here") // Set an invalid token in the header
	res := httptest.NewRecorder()

	// Serve the request using the router
	router := setupRouterWithAuth()
	router.ServeHTTP(res, req)

	// Assert that the response status code is Unauthorized
	assert.Equal(t, http.StatusUnauthorized, res.Code)
}
