package controllers

import (
	"encoding/json"
	"hardwaremonitoringexporter/models"
	"hardwaremonitoringexporter/services"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
    // Decode the request body  to get the user data
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    // Hash the user's password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }
    user.Password = string(hashedPassword) //the password here is replaced with its hashed version

    // Create the user in the database, check if username is already taken
    if err := services.CreateUser(&user); err != nil {
        http.Error(w, "Username already taken", http.StatusBadRequest)
        return
    }

    // Send success response if the user is registered
    json.NewEncoder(w).Encode("User registered successfully")
}


func LoginUser(w http.ResponseWriter, r *http.Request) {
    var user, dbUser models.User

    // Decode the request body to get login credentials
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    // Fetch user from the database by username
    if err := services.GetUserByUsername(&dbUser, user.Username); err != nil {
        if err == gorm.ErrRecordNotFound {
            http.Error(w, "Invalid credentials", http.StatusUnauthorized) // Username not found
        } else {
            http.Error(w, "Server error", http.StatusInternalServerError)
        }
        return
    }

    // Compare the hashed password stored in the database with the provided password
    if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized) // Password mismatch
        return
    }

    // Generate JWT token on successful authentication
    token, err := services.GenerateJWT(dbUser)
    if err != nil {
        http.Error(w, "Error generating token", http.StatusInternalServerError)
        return
    }

    // Send the generated token as a response
    json.NewEncoder(w).Encode(map[string]string{"token": token})
}

