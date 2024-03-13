package handler

import (
	"encoding/json"
	"filmhub/internal/model"
	"net/http"
)

type SignInInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		newErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var input SignInInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{"token": token}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		newErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var input model.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{"id": id}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}
