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
	GetAllMovies(sortBy, sortOrder string) ([]model.Movie, error)
	CreateMovie(movie model.InputMovie) error
	GetMovieByID(movieID int) (model.MovieWithActors, error)
	DeleteByID(movieID int) error
	UpdateMovie(movieID int, data model.InputMovie) error
	GetMoviesByTitle(title string) ([]model.MovieWithActors, error)
	GetMoviesByActor(actor string) ([]model.MovieWithActors, error)
}

type Actor interface {
	GetAllActors() ([]model.ActorWithMovies, error)
	CreateActor(actor model.InputActor) error
	Delete(actorID int) error
	Get(actorID int) (model.ActorWithMovies, error)
	Update(actorID int, data model.InputActor) error
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
