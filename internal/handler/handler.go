package handler

import (
	"net/http"

	_ "github.com/avealice/filmhub/docs"
	"github.com/avealice/filmhub/internal/service"

	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
			return
		}

		httpSwagger.Handler(
			httpSwagger.URL("http://127.0.0.1:8000/swagger/doc.json"),
			httpSwagger.DeepLinking(true),
			httpSwagger.DocExpansion("none"),
			httpSwagger.DomID("swagger-ui"),
		).ServeHTTP(w, r)
	})

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
