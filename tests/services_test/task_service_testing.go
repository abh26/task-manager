package services_test

import (
	"hardwaremonitoringexporter/database"
	"hardwaremonitoringexporter/models"
	"hardwaremonitoringexporter/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateTask(t *testing.T) {
	task := &models.Task{Title: "Test Task", Description: "Test Description", UserID: 1}
	err := services.CreateTask(task)

	assert.NoError(t, err)     // Expect no error
	assert.NotZero(t, task.ID) // Expect the task ID to be set
}

func TestGetTasksByUserID(t *testing.T) {
	// Prepare test data
	userID := uint(1)
	services.CreateTask(&models.Task{Title: "Task 1", UserID: userID})
	services.CreateTask(&models.Task{Title: "Task 2", UserID: userID})

	tasks, totalTasks, err := services.GetTasksByUserID(userID, 1, 10, "created_at", "asc")

	assert.NoError(t, err)                // Expect no error
	assert.Equal(t, int64(2), totalTasks) // Expect total tasks count to match
	assert.Len(t, tasks, 2)               // Expect 2 tasks to be returned
}

func TestUpdateTask(t *testing.T) {
	task := &models.Task{Title: "Task to Update", UserID: 1}
	services.CreateTask(task)

	task.Title = "Updated Task Title"
	err := services.UpdateTask(task)

	assert.NoError(t, err) // Expect no error

	// Fetch the updated task
	var updatedTask models.Task
	err = database.DB.First(&updatedTask, task.ID).Error
	assert.NoError(t, err)                                   // Expect no error
	assert.Equal(t, "Updated Task Title", updatedTask.Title) // Expect the title to be updated
}

func TestDeleteTask(t *testing.T) {
	task := &models.Task{Title: "Task to Delete", UserID: 1}
	services.CreateTask(task)

	err := services.DeleteTask(int(task.ID), task.UserID)

	assert.NoError(t, err) // Expect no error

	// Verify the task is deleted
	var deletedTask models.Task
	err = database.DB.First(&deletedTask, task.ID).Error
	assert.Error(t, err)                         // Expect an error (task should not be found)
	assert.Equal(t, gorm.ErrRecordNotFound, err) // Expect the error to be record not found
}

func TestUpdateTaskStatus(t *testing.T) {
	task := &models.Task{Title: "Task Status Test", UserID: 1, Status: "todo"}
	services.CreateTask(task)

	// Test valid status update
	err := services.UpdateTaskStatus(int(task.ID), "done")
	assert.NoError(t, err) // Expect no error

	// Verify the status is updated
	var updatedTask models.Task
	err = database.DB.First(&updatedTask, task.ID).Error
	assert.NoError(t, err)                      // Expect no error
	assert.Equal(t, "done", updatedTask.Status) // Expect status to be "done"

	// Test invalid status
	err = services.UpdateTaskStatus(int(task.ID), "invalid_status")
	assert.Error(t, err)                                 // Expect error for invalid status
	assert.Equal(t, "invalid status value", err.Error()) // Check error message
}
