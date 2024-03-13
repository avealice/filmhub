package handler

import (
	"encoding/json"
	"filmhub/internal/model"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) movieHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createMovie(w, r)
	case http.MethodDelete:
		h.deleteMovie(w, r)
	case http.MethodPut:
		h.updateMovie(w, r)
	case http.MethodGet:
		h.getMovie(w, r)
	default:
		newErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *Handler) getAllMovies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		newErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	if sortBy == "" {
		sortBy = "rating"
	}

	if sortOrder == "" {
		sortOrder = "desc"
	}

	movies, err := h.services.Movie.GetAllMovies(sortBy, sortOrder)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, "Failed to get movies")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func (h *Handler) createMovie(w http.ResponseWriter, r *http.Request) {
	role, err := getUserRole(r)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, "failed to get user role")
		return
	}

	if role != "admin" {
		newErrorResponse(w, http.StatusForbidden, "only admin can create movies")
		return
	}

	var input model.Movie
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Movie.CreateMovie(input)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, "failed to create movie")
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("movie created successfully"))
}

func (h *Handler) deleteMovie(w http.ResponseWriter, r *http.Request) {
	role, err := getUserRole(r)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, "failed to get user role")
		return
	}

	if role != "admin" {
		newErrorResponse(w, http.StatusForbidden, "only admin can delete movies")
		return
	}

	path := strings.TrimSuffix(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) != 3 || parts[1] != "movie" {
		newErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	movieID, err := strconv.Atoi(parts[2])
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, "Invalid movie ID")
		return
	}

	err = h.services.Movie.DeleteByID(movieID)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, "Failed to delete movie by ID")
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("movie deleted successfully"))
}

func (h *Handler) updateMovie(w http.ResponseWriter, r *http.Request) {
	role, err := getUserRole(r)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, "failed to get user role")
		return
	}

	if role != "admin" {
		newErrorResponse(w, http.StatusForbidden, "only admin can update movies")
		return
	}

	path := strings.TrimSuffix(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) != 3 || parts[1] != "movie" {
		newErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	movieID, err := strconv.Atoi(parts[2])
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, "Invalid movie ID")
		return
	}

	var input model.Movie
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.UpdateMovie(movieID, input)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, "Failed to update movie")
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("movie updated successfully"))

}

func (h *Handler) getMovie(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) != 3 || parts[1] != "movie" {
		newErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	movieID, err := strconv.Atoi(parts[2])
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, "Invalid movie ID")
		return
	}

	movie, err := h.services.Movie.GetMovieByID(movieID)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, "Failed to get movie by ID")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(movie)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, "Failed to encode movie to JSON")
		return
	}
}
