package models

// User struct represents the user data
type User struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login struct represents the user login data
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
