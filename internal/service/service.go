package service

import (
	"filmhub/internal/model"
	"filmhub/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (string, error)
}

type Movie interface {
	GetAllMovies(sortBy, sortOrder string) ([]model.Movie, error)
	CreateMovie(movie model.InputMovie) error
	GetMovieByID(movieID int) (model.MovieWithActors, error)
	DeleteByID(movieID int) error
	UpdateMovie(movieID int, data model.InputMovie) error
	GetMoviesByActor(actor string) ([]model.MovieWithActors, error)
	GetMoviesByTitle(title string) ([]model.MovieWithActors, error)
}

type Actor interface {
	CreateActor(actor model.InputActor) error
	GetAllActors() ([]model.ActorWithMovies, error)
	Delete(actorID int) error
	Get(actorID int) (model.ActorWithMovies, error)
	Update(actorID int, data model.InputActor) error
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
