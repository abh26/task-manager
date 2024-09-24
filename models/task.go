package models
import (
    "gorm.io/gorm"
)

type Task struct {
    gorm.Model
    Title       string `json:"title" validate:"required"`
    Description string `json:"description" validate:"required"`
    Status      string `json:"status" validate:"required,oneof=todo in_progress done"`
    UserID      uint   `json:"user_id"`
}

type GettaskRequest struct {
    Page      int    `json:"page"`      // Page number for  pagination
    PageSize  int    `json:"pageSize"`  // Number of tasks per page
    SortBy    string `json:"sortBy"`    // Field to sort by (e.g., status, createdAt)
    SortOrder string `json:"sortOrder"` // Order to sort (asc or desc)
}
