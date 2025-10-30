package models

import "time"

// Todo struct represents the todo list data
type Todo struct {
	ID          int64     `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    string    `json:"priority"`
	Status      string    `json:"status"`
	DueAt       time.Time `json:"dueAt"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

type PaginatedTodos struct {
	TotalCount int64   `json:"total_count"`
	Todos      []*Todo `json:"todos"`
}

type Request struct {
	// UserID string `json:"user_id" binding:"required"`
	Status string `json:"status"` // "ALL", "PENDING", "COMPLETED"
	Limit  int    `json:"limit" default:"10"`
	Offset int    `json:"offset" default:"0"`
}
