package handler

import (
	"encoding/json"
	"filmhub/internal/model"
	"net/http"
	"strconv"
	"strings"
)

// getAllMovies возвращает список всех фильмов.
// @Summary Получить все фильмы
// @Description Получает список всех фильмов с возможностью сортировки.
// @Tags /api/movies
// @Produce json
// @Param sort_by query string false "Критерий сортировки (title, rating, release_date)"
// @Param sort_order query string false "Порядок сортировки (asc, desc)"
// @Success 200 {array} model.Movie "Список фильмов"
// @Failure 400 {object} ErrorResponse "Некорректный запрос или данные"
// @Failure 401 {object} ErrorResponse "Пустой заголовок авторизации"
// @Failure 405 {object} ErrorResponse "Некорректный метод"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/movies [get]
// @Security ApiKeyAuth
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

// createMovie создает новый фильм.
// @Summary Создать фильм
// @Description Создает новый фильм.
// @Tags /api/movie
// @Accept json
// @Produce json
// @Param movie body model.InputMovie true "Данные нового фильма"
// @Success 201 "Фильм создан успешно"
// @Failure 400 {object} ErrorResponse "Некорректный запрос или данные"
// @Failure 401 {object} ErrorResponse "Пустой заголовок авторизации"
// @Failure 403 {object} ErrorResponse "Некорректная роль"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/movie [post]
// @Security ApiKeyAuth
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

	var input model.InputMovie
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Movie.CreateMovie(input)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("movie created successfully"))
}

// deleteMovie удаляет фильм по его идентификатору.
// @Summary Удалить фильм
// @Description Удаляет фильм по его идентификатору.
// @Tags /api/movie/{id}
// @Param id path int true "Идентификатор фильма"
// @Success 200 "Фильм удален успешно"
// @Failure 400 {object} ErrorResponse "Некорректный запрос или данные"
// @Failure 401 {object} ErrorResponse "Пустой заголовок авторизации"
// @Failure 403 {object} ErrorResponse "Некорректная роль"
// @Failure 405 {object} ErrorResponse "Некорректный метод"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/movie/{id} [delete]
// @Security ApiKeyAuth
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

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("movie deleted successfully"))
}

// searchMovie выполняет поиск фильмов по указанным критериям.
// @Summary Поиск фильмов
// @Description Выполняет поиск фильмов по указанным критериям (название или актер).
// @Tags /api/movie/search
// @Produce json
// @Param title query string false "Название фильма для поиска"
// @Param actor query string false "Имя актера для поиска"
// @Success 200 {array} model.MovieWithActors "Список фильмов, удовлетворяющих критериям поиска"
// @Failure 400 {object} ErrorResponse "Некорректный запрос или данные"
// @Failure 401 {object} ErrorResponse "Пустой заголовок авторизации"
// @Failure 403 {object} ErrorResponse "Некорректная роль"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/movie/search [get]
// @Security ApiKeyAuth
func (h *Handler) searchMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		newErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	actor := r.URL.Query().Get("actor")
	title := r.URL.Query().Get("title")

	if (title == "" && actor == "") || (title != "" && actor != "") {
		newErrorResponse(w, http.StatusBadRequest, "Invalid search request")
		return
	}

	var movies []model.MovieWithActors
	var err error

	if title != "" {
		movies, err = h.services.Movie.GetMoviesByTitle(title)
	} else if actor != "" {
		movies, err = h.services.Movie.GetMoviesByActor(actor)
	}

	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

// updateMovie обновляет информацию о фильме.
// @Summary Обновить информацию о фильме
// @Description Обновляет информацию о фильме.
// @Tags /api/movie/{id}
// @Accept json
// @Produce json
// @Param id path int true "Идентификатор фильма"
// @Param movie body model.InputMovie true "Новые данные о фильме"
// @Success 200 "Фильм обновлен успешно"
// @Failure 400 {object} ErrorResponse "Некорректный запрос или данные"
// @Failure 401 {object} ErrorResponse "Пустой заголовок авторизации"
// @Failure 403 {object} ErrorResponse "Некорректная роль"
// @Failure 405 {object} ErrorResponse "Некорректный метод"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/movie/{id} [put]
// @Security ApiKeyAuth
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
	if err != nil || movieID < 0 {
		newErrorResponse(w, http.StatusBadRequest, "Invalid movie ID")
		return
	}

	var input model.InputMovie
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.UpdateMovie(movieID, input)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("movie updated successfully"))
}

// getMovie возвращает информацию о фильме по его идентификатору.
// @Summary Получить информацию о фильме
// @Description Получает информацию о фильме по его идентификатору.
// @Tags /api/movie/{id}
// @Produce json
// @Param id path int true "Идентификатор фильма"
// @Success 200 {object} model.MovieWithActors "Информация о фильме"
// @Failure 400 {object} ErrorResponse "Некорректный запрос или данные"
// @Failure 401 {object} ErrorResponse "Пустой заголовок авторизации"
// @Failure 403 {object} ErrorResponse "Некорректная роль"
// @Failure 405 {object} ErrorResponse "Некорректный метод"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/movie/{id} [get]
// @Security ApiKeyAuth
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
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(movie)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, "Failed to encode movie to JSON")
		return
	}
}

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
