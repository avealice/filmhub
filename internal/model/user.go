package model

// User represents a user in the system.
type User struct {
	Id       int    `json:"-"`                           // Unique identifier for the user
	Username string `json:"username" db:"username"`      // Username of the user
	Password string `json:"password" db:"password_hash"` // Password hash of the user
	Role     string `json:"role" db:"role"`              // Role of the user
}
