package services

import (
	"fmt"
	"hardwaremonitoringexporter/database"
	"hardwaremonitoringexporter/models"
)

// CreateUser inserts a new user into the database and logs user details
func CreateUser(user *models.User) error {
	fmt.Println("User to be created with details as follows:", user) // Log the  user details for debugging
	return database.DB.Create(user).Error // Return any error that occurs during the creation
}

// GetUserByUsername fetches a user from the database based on the provided username
func GetUserByUsername(user *models.User, username string) error {
	return database.DB.Where("username = ?", username).First(user).Error // Return any error if the user is not found or if there is a database error
}

