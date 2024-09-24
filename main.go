package main

import (
	"log"
	"net/http"
	"hardwaremonitoringexporter/controllers"
	"hardwaremonitoringexporter/database"
	"hardwaremonitoringexporter/middlewares"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize the database connection
	// This will set up the necessary database configurations and open a connection to the database.
	database.InitDB()

	// A Gorilla Mux router to handle incoming HTTP requests.
	r := mux.NewRouter()

	// Public routes: These do not require authentication
	// Register a new user by providing details like username, password, etc.
	r.HandleFunc("/register", controllers.RegisterUser).Methods("POST")
	
	// Log in the user by verifying the provided credentials.
	// On successful login, a JWT token is generated and returned to the client.
	r.HandleFunc("/login", controllers.LoginUser).Methods("POST")

	// Protected routes: These routes require authentication
	// The API subrouter is prefixed with "/api" and uses the AuthMiddleware to ensure only authenticated users can access these routes.
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middlewares.AuthMiddleware) // Apply the AuthMiddleware to all routes under "/api"

	// Create a new task
	// This allows an authenticated user to create a new task by providing a title, description, and status.
	api.HandleFunc("/tasks", controllers.CreateTask).Methods("POST")

	// Get all tasks for the authenticated user with pagination and sorting.
	// The user can pass query parameters like page, pageSize, sortBy, and sortOrder.
	api.HandleFunc("/gettasks", controllers.GetTasks).Methods("GET")

	// Update an existing task identified by its ID
	// This allows the authenticated user to update the task's details such as title, description, and status.
	api.HandleFunc("/updatetasks/{id}", controllers.UpdateTask).Methods("PUT")

	// Concurrently mark multiple tasks as "done"
	// This endpoint demonstrates concurrency by updating multiple tasks' statuses concurrently using goroutines.
	api.HandleFunc("/concurrent", controllers.MarkTasksAsDone).Methods("PUT")

	// Delete a task by its ID
	// The authenticated user can delete a task by providing the task ID.
	api.HandleFunc("/deltasks/{id}", controllers.DeleteTask).Methods("DELETE")

	// Log a message indicating that the server has started and listen on port 3000 for incoming HTTP requests.
	log.Println("server started")
	log.Fatal(http.ListenAndServe(":3000", r)) // Start the server on port 3000 and bind the router.
}
