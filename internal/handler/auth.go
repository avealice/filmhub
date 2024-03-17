package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/avealice/filmhub/internal/model"
	"github.com/sirupsen/logrus"
)

type SignInInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

// signIn авторизует пользователя и возвращает токен доступа.
//
// @Summary Авторизация пользователя
// @Description Авторизует пользователя с заданными учетными данными и возвращает токен доступа
// @Tags /auth/
// @Accept json
// @Produce json
// @Param request body SignInInput true "Данные для входа"
// @Success 200 {object} TokenResponse "Токен доступа"
// @Failure 400 {object} ErrorResponse "Некорректный запрос или данные"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /auth/sign-in [post]
func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		newErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var input SignInInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, errors.New("Invalid input").Error())
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, errors.New("User signin in unsuccessfully").Error())
		return
	}

	userID, _ := getUserID(r)

	logrus.WithField("user_id", userID).Info("User signed in successfully")

	response := TokenResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

type UserIDResponse struct {
	ID int `json:"id"`
}

// signUp регистрирует нового пользователя.
//
// @Summary Регистрация пользователя
// @Description Регистрирует нового пользователя с заданными данными
// @Tags /auth/
// @Accept json
// @Produce json
// @Param request body model.User true "Данные нового пользователя"
// @Success 201 {object} UserIDResponse "ID нового пользователя"
// @Failure 400 {object} ErrorResponse "Некорректный запрос или данные"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
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

	if len(input.Password) == 0 || len(input.Username) == 0 {
		newErrorResponse(w, http.StatusBadRequest, errors.New("Invalid input").Error())
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, errors.New("User signed up unsuccessfully").Error())
		return
	}

	logrus.WithField("userID", id).Info("User signed up successfully")

	response := UserIDResponse{ID: id}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}
