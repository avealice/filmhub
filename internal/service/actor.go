package service

import (
	"filmhub/internal/model"
	"filmhub/internal/repository"
)

type ActorService struct {
	r repository.Actor
}

func NewActorService(r repository.Actor) *ActorService {
	return &ActorService{
		r: r,
	}
}

func (s *ActorService) CreateActor(actor model.InputActor) error {
	return s.r.CreateActor(actor)
}

func (s *ActorService) GetAllActors() ([]model.ActorWithMovies, error) {
	return s.r.GetAllActors()
}

func (s *ActorService) Delete(actorID int) error {
	return s.r.Delete(actorID)
}

func (s *ActorService) Get(actorID int) (model.ActorWithMovies, error) {
	return s.r.Get(actorID)
}

func (s *ActorService) Update(actorID int, data model.InputActor) error {
	return s.r.Update(actorID, data)
}
