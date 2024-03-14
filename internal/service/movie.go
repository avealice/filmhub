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

func (s *MovieService) GetAllMovies(sortBy, sortOrder string) ([]model.Movie, error) {
	return s.r.GetAllMovies(sortBy, sortOrder)
}

func (s *MovieService) CreateMovie(movie model.MovieWithActors) error {
	return s.r.CreateMovie(movie)
}

func (s *MovieService) GetMovieByID(movieID int) (model.MovieWithActors, error) {
	return s.r.GetMovieByID(movieID)
}

func (s *MovieService) DeleteByID(movieID int) error {
	return s.r.DeleteByID(movieID)
}

func (s *MovieService) UpdateMovie(movieID int, data model.MovieWithActors) error {
	return s.r.UpdateMovie(movieID, data)
}

func (s *MovieService) GetMoviesByTitle(title string) ([]model.MovieWithActors, error) {
	return s.r.GetMoviesByTitle(title)
}

func (s *MovieService) GetMoviesByActor(actor string) ([]model.MovieWithActors, error) {
	return s.r.GetMoviesByActor(actor)
}
