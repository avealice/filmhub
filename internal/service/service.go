package service

import (
	"filmhub/internal/model"
	"filmhub/internal/repository"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, string, error)
}

type Movie interface {
	GetAllMovies(sortBy, sortOrder string) ([]model.MovieWithActors, error)
	CreateMovie(movie model.Movie) error
	GetMovieByID(movieID int) (model.Movie, error)
	DeleteByID(movieID int) error
	CreateMovieForActor(actorID int, movie model.Movie) error
	UpdateMovie(movieID int, data model.Movie) error
	// GetMovie(data model.Movie) (model.Movie, error)
}

type Actor interface {
	CreateActor(actor model.Actor) error
	GetAllActors() ([]model.ActorWithMovies, error)
	Delete(actorID int) error
	Get(actorID int) (model.Actor, error)
}

type Service struct {
	Authorization
	Movie
	Actor
}

func NewService(r *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(r.Authorization),
		Movie:         NewMovieService(r.Movie),
		Actor:         NewActorService(r.Actor),
	}
}
