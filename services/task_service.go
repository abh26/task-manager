package services

import (
	"errors"
	"hardwaremonitoringexporter/database"
	"hardwaremonitoringexporter/models"
)

// CreateTask inserts a new task into the database
func CreateTask(task *models.Task) error {
	return database.DB.Create(task).Error // Return any error that occurs during the creation
}

// GetTasksByUserID fetches tasks for a specific user with  pagination and sorting
func GetTasksByUserID(userID uint, page, pageSize int, sortBy, sortOrder string) ([]models.Task, int64, error) {
	var tasks []models.Task
	var totalTasks int64

	// Calculate the offset for pagination
	offset := (page - 1) * pageSize

	// Build the query with filtering, sorting, limit, and offset
	query := database.DB.Where("user_id = ?", userID).
		Order(sortBy + " " + sortOrder).
		Limit(pageSize).
		Offset(offset)

	// Execute the query to fetch tasks
	err := query.Find(&tasks).Error
	if err != nil {
		return nil, 0, err // Return error if fetching tasks fails
	}

	// Get the total count of tasks for pagination info
	err = database.DB.Model(&models.Task{}).Where("user_id = ?", userID).Count(&totalTasks).Error
	if err != nil {
		return nil, 0, err // Return error if counting tasks fails
	}

	return tasks, totalTasks, nil // Return tasks and total count
}

// UpdateTask updates an existing task in the database
func UpdateTask(task *models.Task) error {
	return database.DB.Save(task).Error // Return any error that occurs during the update
}

// DeleteTask removes a task from the database if it belongs to the user
func DeleteTask(taskID int, userID uint) error {
	return database.DB.Where("id = ? AND user_id = ?", taskID, userID).Delete(&models.Task{}).Error // Return any error during deletion
}

// UpdateTaskStatus changes the status of a task to a new valid status
func UpdateTaskStatus(taskID int, status string) error {
	// Validate status
	validStatuses := map[string]bool{"todo": true, "in progress": true, "done": true}
	if !validStatuses[status] {
		return errors.New("invalid status value") // Return error if status is invalid
	}

	// Fetch the task by ID
	var task models.Task
	if err := database.DB.First(&task, taskID).Error; err != nil {
		return err // Return error if task not found
	}

	// Update the task's status
	task.Status = status
	if err := database.DB.Save(&task).Error; err != nil {
		return err // Return error if saving the updated task fails
	}

	return nil // Return nil if successful
}

