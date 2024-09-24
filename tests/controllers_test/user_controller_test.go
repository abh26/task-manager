package controllers_test

import (
	"bytes"
	"encoding/json"
	"hardwaremonitoringexporter/controllers"
	"hardwaremonitoringexporter/models"
	"hardwaremonitoringexporter/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func setupuserRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/register", controllers.RegisterUser).Methods("POST")
	r.HandleFunc("/login", controllers.LoginUser).Methods("POST")
	return r
}

func TestRegisterUser(t *testing.T) {
	r := setupuserRouter()

	user := models.User{Username: "testuser", Password: "password123"}
	body, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	res := httptest.NewRecorder()

	// Public routes: These do not require authentication
	// Register a new user by providing details like username, password, etc.
	r.HandleFunc("/register", controllers.RegisterUser).Methods("POST")
	
	


	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var responseMessage string
	json.NewDecoder(res.Body).Decode(&responseMessage)
	assert.Equal(t, "User registered successfully", responseMessage)
}

func TestLoginUser(t *testing.T) {
	r := setupuserRouter()

	user := models.User{Username: "testuser", Password: "password123"}
	body, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	res := httptest.NewRecorder()

	// On successful login, a JWT token is generated and returned to the client.
	r.HandleFunc("/login", controllers.LoginUser).Methods("POST")
	// Mock the GenerateJWT service to simulate token generation
	services.GenerateJWT(user)

	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var response map[string]string
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, "mocked-token", response["token"])
}
