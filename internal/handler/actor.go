package handler

import (
	"encoding/json"
	"filmhub/internal/model"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// getAllActors получает всех актеров.
//
// @Summary Получить всех актеров.
// @Description Получить всех актеров из базы данных.
// @Tags /api/actors
// @Produce json
// @Success 200 {array} model.ActorWithMovies
// @Failure 401 {object} ErrorResponse "Пустой заголовок авторизации"
// @Failure 405 {object} ErrorResponse "Некорректный метод"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/actors [get]
// @Security ApiKeyAuth
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

// createActor создает актера.
//
// @Summary Создать актера.
// @Description Создает нового актера.
// @Tags /api/actor
// @Accept json
// @Produce json
// @Param actor body model.InputActor true "Данные нового актера"
// @Success 201 {string} string "Актер успешно создан"
// @Failure 401 {object} ErrorResponse "Пустой заголовок авторизации"
// @Failure 405 {object} ErrorResponse "Некорректный метод"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/actor [post]
// @Security ApiKeyAuth
func (h *Handler) CreateActor(w http.ResponseWriter, r *http.Request) {
	role, err := getUserRole(r)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, "failed to get user role")
		return
	}

	if role != "admin" {
		newErrorResponse(w, http.StatusForbidden, "only admin can create actors")
		return
	}

	var input model.InputActor
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Actor.CreateActor(input)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("actor created successfully"))
}

// deleteActor удаляет актера.
//
// @Summary Удалить актера.
// @Description Удаляет актера по его идентификатору.
// @Tags /api/actor/{id}
// @Param id path int true "Идентификатор актера"
// @Success 200 {string} string "Актер успешно удален"
// @Failure 400 {object} ErrorResponse "Некорректный запрос или данные"
// @Failure 401 {object} ErrorResponse "Пустой заголовок авторизации"
// @Failure 403 {object} ErrorResponse "Некорректная роль"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/actor/{id} [delete]
// @Security ApiKeyAuth
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

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("actor deleted successfully"))
}

// updateActor обновляет информацию об актере.
//
// @Summary Обновить информацию об актере.
// @Description Обновляет информацию об актере по его идентификатору.
// @Tags /api/actor/{id}
// @Accept json
// @Produce json
// @Param id path int true "Идентификатор актера"
// @Param actor body model.InputActor true "Новые данные актера"
// @Success 200 {string} string "Информация об актере успешно обновлена"
// @Failure 400 {object} ErrorResponse "Некорректный запрос или данные"
// @Failure 401 {object} ErrorResponse "Пустой заголовок авторизации"
// @Failure 403 {object} ErrorResponse "Некорректная роль"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/actor/{id} [put]
// @Security ApiKeyAuth
func (h *Handler) updateActor(w http.ResponseWriter, r *http.Request) {
	role, err := getUserRole(r)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, "failed to get user role")
		return
	}

	if role != "admin" {
		newErrorResponse(w, http.StatusForbidden, "only admin can update actors")
		return
	}

	path := strings.TrimSuffix(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) != 3 || parts[1] != "actor" {
		newErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	actodID, err := strconv.Atoi(parts[2])
	if err != nil || actodID < 0 {
		newErrorResponse(w, http.StatusBadRequest, "Invalid actor ID")
		return
	}

	var input model.InputActor
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Actor.Update(actodID, input)
	if err != nil {
		fmt.Println(err.Error())
		newErrorResponse(w, http.StatusInternalServerError, "Failed to update actor")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("actor updated successfully"))
}

// getActor получает информацию об актере.
//
// @Summary Получить информацию об актере.
// @Description Получает информацию об актере по его идентификатору.
// @Tags /api/actor/{id}
// @Produce json
// @Param id path int true "Идентификатор актера"
// @Success 200 {object} model.ActorWithMovies
// @Failure 400 {object} ErrorResponse "Некорректный запрос или данные"
// @Failure 401 {object} ErrorResponse "Пустой заголовок авторизации"
// @Failure 405 {object} ErrorResponse "Некорректный метод"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/actor/{id} [get]
// @Security ApiKeyAuth
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
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(actor)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, "Failed to encode movie to JSON")
		return
	}
}

func (h *Handler) actorHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
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
