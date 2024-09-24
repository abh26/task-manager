package controllers

import (
	"encoding/json"
	"fmt"
	"hardwaremonitoringexporter/models"
	"hardwaremonitoringexporter/services"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {
	// Declared a variable named task whose structure is of the request type
	var task models.Task
	// Here we decode entire raw json from the body into our variable task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	// Fetching userID from context after the middleware decodes the jwt token
	userID := r.Context().Value("userID").(uint)
	fmt.Println("UserId : ", userID)
	task.UserID = userID
	// Based on the userID and the request stored in task variable we create a task into our database system
	if err := services.CreateTask(&task); err != nil {
		http.Error(w, "Could not create task", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(task)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	var TaskReq models.GettaskRequest
	if err := json.NewDecoder(r.Body).Decode(&TaskReq); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	userID := r.Context().Value("userID").(uint)

	// Set default values for pagination
	page := 1
	pageSize := 10

	// check for page not equals to 0,if equal to 0 it will use the default value set above
	if TaskReq.Page != 0 {
		page = TaskReq.Page
	}

	// check for pagesize not equals to 0,if equal to 0 it will use the default value set above
	if TaskReq.PageSize != 0 {
		pageSize = TaskReq.PageSize
	}

	// Validate sortBy and sortOrder
	if TaskReq.SortBy == "" {
		TaskReq.SortBy = "createdAt" // default sort field
	}
	if TaskReq.SortOrder != "asc" && TaskReq.SortOrder != "desc" {
		TaskReq.SortOrder = "asc" // default sort order
	}

	// Fetch tasks with pagination and sorting
	tasks, totalTasks, err := services.GetTasksByUserID(userID, page, pageSize, TaskReq.SortBy, TaskReq.SortOrder)
	if err != nil {
		http.Error(w, "Could not fetch tasks", http.StatusInternalServerError)
		return
	}

	// Create response structure
	response := map[string]interface{}{
		"tasks":       tasks,
		"totalTasks":  totalTasks,
		"currentPage": page,
		"pageSize":    pageSize,
	}

	// Set response header and encode response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	// Define a variable of type models.Task to hold the incoming task data from the request body
	var task models.Task

	// Decode the JSON body of the request into the task variable
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		// If decoding fails, return a 400 Bad Request status with an error message
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Extract the userID from the request's context (set by the authentication middleware)
	userID := r.Context().Value("userID").(uint)

	// Assign the userID to the task's UserID field to ensure only the correct user can update their own task
	task.UserID = userID

	// Call the service function to update the task in the database
	if err := services.UpdateTask(&task); err != nil {
		// If the service returns an error, send a 500 Internal Server Error response
		http.Error(w, "Could not update task", http.StatusInternalServerError)
		return
	}

	// If the update is successful, return the updated task as JSON in the response body
	json.NewEncoder(w).Encode(task)
}

func MarkTasksAsDone(w http.ResponseWriter, r *http.Request) {
	// Decode the request body to get a list of task IDs
	var taskIDs []int
	if err := json.NewDecoder(r.Body).Decode(&taskIDs); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	ch := make(chan error, len(taskIDs)) // Channel to collect errors for each goroutine

	// Loop over each task ID and process them concurrently
	for _, id := range taskIDs {
		wg.Add(1)
		go func(taskID int) {
			defer wg.Done()
			// Update the task status to "done"
			err := services.UpdateTaskStatus(taskID, "done")
			ch <- err
		}(id) // Pass the task ID to the goroutine
	}

	wg.Wait() // Wait for all goroutines to complete
	close(ch) // Close the channel once all tasks are done

	// Check if any errors occurred
	for err := range ch {
		if err != nil {
			http.Error(w, "Error updating tasks", http.StatusInternalServerError)
			return
		}
	}

	// Send success response if all tasks were updated
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Tasks updated successfully")
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Extract URL variables
	vars := mux.Vars(r)
	taskID := vars["id"] // Get the "id" from the URL variables
	TaskIDConverted,_:=strconv.Atoi(taskID)
	// Retrieve the authenticated user's ID from the request's context (set by authentication middleware)
	userID := r.Context().Value("userID").(uint)

	// Call the service to delete the task by passing the task ID and user ID to ensure only the task owner can delete it
	if err := services.DeleteTask(TaskIDConverted, userID); err != nil {
		// If the deletion fails, respond with a 500 Internal Server Error and return an error message
		http.Error(w, "Could not delete task", http.StatusInternalServerError)
		return
	}

	// If the task is successfully deleted, return a 204 No Content status (indicating success with no response body)
	w.WriteHeader(http.StatusNoContent)
}
