package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/avealice/filmhub/internal/model"

	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user model.User) (int, error) {
	var id int

	user.Role = "user"

	query := fmt.Sprintf("INSERT INTO %s (username, password_hash, role) values ($1, $2, $3) RETURNING id", usersTable)

	row := r.db.QueryRow(query, user.Username, user.Password, user.Role)
	if err := row.Scan(&id); err != nil {
		return -1, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT id, role FROM %s WHERE username=$1 AND password_hash=$2", usersTable)

	err := r.db.Get(&user, query, username, password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, errors.New("user not found")
		}
		return model.User{}, err
	}

	return user, nil
}
