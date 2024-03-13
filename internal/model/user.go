package model

type User struct {
	Id       int    `json:"-"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password_hash"`
	Role     string `json:"role" db:"role"`
}
