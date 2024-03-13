package handler

import (
	"encoding/json"
	"filmhub/internal/model"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) actorHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createMovieForActor(w, r)
	case http.MethodDelete:
		h.deleteActor(w, r)
	case http.MethodPut:
		h.updateActor(w, r)
	case http.MethodGet:
		h.getActor(w, r)
	default:
		newErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *Handler) getAllActors(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		newErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	actors, err := h.services.Actor.GetAllActors()
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, "Failed to get actors")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actors)
}

func (h *Handler) createMovieForActor(w http.ResponseWriter, r *http.Request) {
	role, err := getUserRole(r)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, "failed to get user role")
		return
	}

	if role != "admin" {
		newErrorResponse(w, http.StatusForbidden, "only admin can create movie")
		return
	}

	path := strings.TrimSuffix(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) != 4 || parts[3] != "movie" {
		newErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	actorId, err := strconv.Atoi(parts[2])
	if err != nil {
		newErrorResponse(w, http.StatusMethodNotAllowed, "invalid actor id")
		return
	}

	var input model.Movie
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Movie.CreateMovieForActor(actorId, input)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, "failed to create movie")
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("mnovie created successfully"))
}

func (h *Handler) createActor(w http.ResponseWriter, r *http.Request) {
	role, err := getUserRole(r)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, "failed to get user role")
		return
	}

	if role != "admin" {
		newErrorResponse(w, http.StatusForbidden, "only admin can create actors")
		return
	}

	var input model.Actor
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Actor.CreateActor(input)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, "failed to create actor")
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("actor created successfully"))
}

func (h *Handler) deleteActor(w http.ResponseWriter, r *http.Request) {
	role, err := getUserRole(r)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, "failed to get user role")
		return
	}

	if role != "admin" {
		newErrorResponse(w, http.StatusForbidden, "only admin can delete actors")
		return
	}

	path := strings.TrimSuffix(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) != 3 || parts[1] != "actor" {
		newErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	actorID, err := strconv.Atoi(parts[2])
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, "Invalid actor ID")
		return
	}

	err = h.services.Actor.Delete(actorID)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, "Failed to delete actor")
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("actor deleted successfully"))
}

func (h *Handler) updateActor(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) getActor(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) != 3 || parts[1] != "actor" {
		newErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	actorID, err := strconv.Atoi(parts[2])
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, "Invalid actor ID")
		return
	}

	actor, err := h.services.Actor.Get(actorID)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, "Failed to get actor")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(actor)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, "Failed to encode movie to JSON")
		return
	}
}
