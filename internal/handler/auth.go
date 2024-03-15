package handler

import (
	"encoding/json"
	"errors"
	"filmhub/internal/model"
	"net/http"
)

type SignInInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// signIn авторизует пользователя и возвращает токен доступа.
//
// @Summary Авторизация пользователя
// @Description Авторизует пользователя с заданными учетными данными и возвращает токен доступа
// @Tags Авторизация
// @Accept json
// @Produce json
// @Param request body SignInInput true "Данные для входа"
// @Success 200 {object} map[string]interface{} "map[token]: Токен доступа"
// @Failure 400 {object} errorResponse "Некорректный запрос или данные"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /auth/sign-in [post]
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

// signUp регистрирует нового пользователя.
//
// @Summary Регистрация пользователя
// @Description Регистрирует нового пользователя с заданными данными
// @Tags Регистрация
// @Accept json
// @Produce json
// @Param request body model.User true "Данные нового пользователя"
// @Success 201 {object} map[string]interface{} "map[id]: ID нового пользователя"
// @Failure 400 {object} errorResponse "Некорректный запрос или данные"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /auth/sign-up [post]
func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		newErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var input model.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, errors.New("Invalid input").Error())
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
