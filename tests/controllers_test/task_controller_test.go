package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"hardwaremonitoringexporter/controllers"
	"hardwaremonitoringexporter/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)
func setupRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(middlewareTest) // Add middleware if necessary
	return r
}

// Middleware for testing
func middlewareTest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), "userID", uint(1))) // Mock user ID
		next.ServeHTTP(w, r)
	})
}

// TestCreateTask tests the CreateTask controller.
// It sends a POST request with a sample task and checks if the response  status is 200 OK.
func TestCreateTask(t *testing.T) {
	// Setup router and define the route
	r := setupRouter()
	r.HandleFunc("/tasks", controllers.CreateTask).Methods("POST")

	// Create a sample task and marshal it to JSON
	task := models.Task{Title: "Test Task", Description: "Description of the task"}
	body, _ := json.Marshal(task)

	// Create a new HTTP request with the task data
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	res := httptest.NewRecorder()

	// Serve the HTTP request and record the response
	r.ServeHTTP(res, req)

	// response code is 200 OK
	assert.Equal(t, http.StatusOK, res.Code)
}

// TestGetTasks tests the GetTasks controller.
// It sends a GET request and checks if the response status is 200 OK.
func TestGetTasks(t *testing.T) {
	r := setupRouter()
	r.HandleFunc("/tasks", controllers.GetTasks).Methods("GET")

	req, _ := http.NewRequest("GET", "/tasks", nil)
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}

// TestUpdateTask tests the UpdateTask controller.
// It sends a PUT request with an updated task and checks if the response status is 200 OK.
func TestUpdateTask(t *testing.T) {
	r := setupRouter()
	r.HandleFunc("/tasks", controllers.UpdateTask).Methods("PUT")

	task := models.Task{UserID: 1, Title: "Updated Task"}
	body, _ := json.Marshal(task)

	req, _ := http.NewRequest("PUT", "/tasks", bytes.NewBuffer(body))
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}

// TestMarkTasksAsDone tests the MarkTasksAsDone controller.
// It sends a PUT request with task IDs and checks if the response status is 200 OK.
func TestMarkTasksAsDone(t *testing.T) {
	r := setupRouter()
	r.HandleFunc("/tasks/done", controllers.MarkTasksAsDone).Methods("PUT")

	taskIDs := []int{1, 2, 3}
	body, _ := json.Marshal(taskIDs)

	req, _ := http.NewRequest("PUT", "/tasks/done", bytes.NewBuffer(body))
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}

// TestDeleteTask tests the DeleteTask controller.
// It sends a DELETE request for a specific task ID and checks if the response status is 204 No Content.
func TestDeleteTask(t *testing.T) {
	r := setupRouter()
	r.HandleFunc("/tasks/{id}", controllers.DeleteTask).Methods("DELETE")

	req, _ := http.NewRequest("DELETE", "/tasks/1", nil)
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusNoContent, res.Code)
}

