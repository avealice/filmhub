package model

// User represents a user in the system.
type User struct {
	ID       int    `json:"-" db:"id"`                   // Unique identifier for the user
	Username string `json:"username" db:"username"`      // Username of the user
	Password string `json:"password" db:"password_hash"` // Password hash of the user
	Role     string `json:"-" db:"role"`                 // Role of the user
}
