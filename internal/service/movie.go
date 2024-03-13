package service

import (
	"filmhub/internal/model"
	"filmhub/internal/repository"
)

type MovieService struct {
	r repository.Movie
}

func NewMovieService(r repository.Movie) *MovieService {
	return &MovieService{
		r: r,
	}
}

func (s *MovieService) GetAllMovies(sortBy, sortOrder string) ([]model.MovieWithActors, error) {
	return s.r.GetAllMovies(sortBy, sortOrder)
}

func (s *MovieService) CreateMovie(movie model.Movie) error {
	return s.r.CreateMovie(movie)
}

func (s *MovieService) GetMovieByID(movieID int) (model.Movie, error) {
	return s.r.GetMovieByID(movieID)
}

func (s *MovieService) DeleteByID(movieID int) error {
	return s.r.DeleteByID(movieID)
}

func (s *MovieService) CreateMovieForActor(actorID int, movie model.Movie) error {
	return s.r.CreateMovieForActor(actorID, movie)
}

func (s *MovieService) UpdateMovie(movieID int, data model.Movie) error {
	return s.r.UpdateMovie(movieID, data)
}
