// Package handler предоставляет обработчики HTTP-запросов для filmhub.
package handler

import (
	"filmhub/internal/service"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

// @title FilmHub API
// @version 1.0
// @description API для работы с фильмами и актерами в FilmHub.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.github.com/avealice

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /api/
func (h *Handler) InitRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("docs/swagger.json")
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to open swagger.json: %v", err), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		w.Header().Set("Content-Type", "application/json")

		_, err = io.Copy(w, file)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to copy swagger.json to response: %v", err), http.StatusInternalServerError)
			return
		}
	})

	mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.Dir("swagger-ui"))))

	authMux := http.NewServeMux()
	authMux.HandleFunc("/sign-in", h.signIn)
	authMux.HandleFunc("/sign-up", h.signUp)

	mux.Handle("/auth/", http.StripPrefix("/auth", authMux))

	apiMux := http.NewServeMux()
	apiMux.Handle("/movies", h.userIdentity(http.HandlerFunc(h.getAllMovies)))
	apiMux.Handle("/movie/", h.userIdentity(http.HandlerFunc(h.movieHandle)))
	apiMux.Handle("/movie", h.userIdentity(http.HandlerFunc(h.createMovie)))
	apiMux.Handle("/movie/search", h.userIdentity(http.HandlerFunc(h.searchMovie)))

	apiMux.Handle("/actors", h.userIdentity(http.HandlerFunc(h.getAllActors)))
	apiMux.Handle("/actor", h.userIdentity(http.HandlerFunc(h.CreateActor)))
	apiMux.Handle("/actor/", h.userIdentity(http.HandlerFunc(h.actorHandle)))

	mux.Handle("/api/", http.StripPrefix("/api", apiMux))

	return mux
}
