package entity

import (
	"time"
)

type Task struct {
	ID          int64
	UserID      string
	Title       string
	Description string
	Priority    int
	Status      int
	DueAt       time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// User struct represents the user data
type User struct {
	ID       string
	UserName string
	Email    string
	Password string
}

// Login struct represents the user login data
type Login struct {
	Email    string
	Password string
}
