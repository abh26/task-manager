package services_test

import (
	"hardwaremonitoringexporter/models"
	"hardwaremonitoringexporter/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	user := &models.User{Username: "testuser", Password: "password123"}

	// Call the CreateUser function
	err := services.CreateUser(user)

	assert.NoError(t, err)     // Expect no error
	assert.NotZero(t, user.ID) // Expect the user ID to be set (indicating it was created)
}

func TestGetUserByUsername(t *testing.T) {
	// Prepare a user to be retrieved
	user := &models.User{Username: "testuser", Password: "password123"}
	services.CreateUser(user) // Ensure the user exists in the database

	// getting a user model ready for retrieval
	var retrievedUser models.User
	err := services.GetUserByUsername(&retrievedUser, user.Username)

	assert.NoError(t, err)                                 // Expect no error when retrieving
	assert.Equal(t, user.Username, retrievedUser.Username) // Expect the retrieved username to match
}

func TestGetUserByNonExistentUsername(t *testing.T) {
	var user models.User
	err := services.GetUserByUsername(&user, "nonexistentuser")

	assert.Error(t, err)                         // Expect an error when the user does not exist
	assert.Equal(t, gorm.ErrRecordNotFound, err) // Expect the error to be record not found
}
