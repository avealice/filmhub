package repository

import (
	"filmhub/internal/model"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GetUser(username, password string) (model.User, error)
}

type Movie interface {
	GetAllMovies(sortBy, sortOrder string) ([]model.MovieWithActors, error)
	CreateMovie(movie model.Movie) error
	GetMovieByID(movieID int) (model.Movie, error)
	DeleteByID(movieID int) error
	CreateMovieForActor(actorID int, movie model.Movie) error
	UpdateMovie(movieID int, data model.Movie) error
	// DeleteMovie(id int) error
	// GetMovie(data model.Movie) (model.Movie, error)
}

type Actor interface {
	GetAllActors() ([]model.ActorWithMovies, error)
	CreateActor(actor model.Actor) error
	Delete(actorID int) error
	Get(actorID int) (model.Actor, error)
}

type Repository struct {
	Authorization
	Movie
	Actor
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Movie:         NewMoviePostgres(db),
		Actor:         NewActorPostgres(db),
	}
}
