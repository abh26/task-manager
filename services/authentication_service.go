package services

import (
	"hardwaremonitoringexporter/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var JwtKey = []byte("secret_key") // Secret key used for signing the JWT

// GenerateJWT creates a new JWT for the specified user
func GenerateJWT(user models.User) (string, error) {
	// Define claims, including userID and expiration time (24 hours)
	claims := &jwt.MapClaims{
		"userID": user.ID, // Store user ID in claims
		"exp":    time.Now().Add(24 * time.Hour).Unix(), // Set expiration
	}

	// Create a new token with the specified claims and signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey) // return the token
}

