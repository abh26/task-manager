package services_test

import (
	"hardwaremonitoringexporter/models"
	"hardwaremonitoringexporter/services"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	// Create a dummy user for  testing
	user := models.User{Username: "abhishek", Password: "abh262000"}

	// Generate JWT for the mock user
	tokenString, err := services.GenerateJWT(user)

	// Check that no error occurred
	assert.NoError(t, err)

	// Parse the token to verify its claims
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return services.JwtKey, nil
	})

	// Check that no error occurred during parsing
	assert.NoError(t, err)

	// check for the token is valid
	assert.True(t, token.Valid)

	// Check claims
	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, float64(user.ID), claims["userID"]) // Check userID claim

	// Check expiration time
	exp, ok := claims["exp"].(float64)
	assert.True(t, ok)
	assert.WithinDuration(t, time.Now().Add(24*time.Hour), time.Unix(int64(exp), 0), time.Second) // Check expiration time
}
