package handler

import (
	"filmhub/internal/service"
	"net/http"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() http.Handler {
	mux := http.NewServeMux()

	authMux := http.NewServeMux()
	authMux.HandleFunc("/sign-in", h.signIn)
	authMux.HandleFunc("/sign-up", h.signUp)

	mux.Handle("/auth/", http.StripPrefix("/auth", authMux))

	apiMux := http.NewServeMux()
	apiMux.Handle("/movies", h.userIdentity(http.HandlerFunc(h.getAllMovies)))
	apiMux.Handle("/movie/", h.userIdentity(http.HandlerFunc(h.movieHandle)))
	apiMux.Handle("/movie", h.userIdentity(http.HandlerFunc(h.createMovie)))

	apiMux.Handle("/actors", h.userIdentity(http.HandlerFunc(h.getAllActors)))
	apiMux.Handle("/actor", h.userIdentity(http.HandlerFunc(h.createActor)))
	apiMux.Handle("/actor/", h.userIdentity(http.HandlerFunc(h.actorHandle)))

	mux.Handle("/api/", http.StripPrefix("/api", apiMux))

	return mux
}
