package handler

import (
	"context"
	"encoding/json"
	"errors"
	"filmhub/internal/model"
	"filmhub/internal/service"
	mock_service "filmhub/internal/service/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestHandler_getAllMovies_NoSorting(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMovieService := mock_service.NewMockMovie(ctrl)

	handler := Handler{
		&service.Service{
			Movie: mockMovieService,
		},
	}

	req := httptest.NewRequest("GET", "/api/movies", nil)
	w := httptest.NewRecorder()

	expectedMovies := []model.Movie{
		{ID: 1, Title: "Movie 1", Description: "Description 1", ReleaseDate: "2022-01-01", Rating: 85},
		{ID: 2, Title: "Movie 2", Description: "Description 2", ReleaseDate: "2022-01-02", Rating: 79},
	}

	mockMovieService.EXPECT().GetAllMovies("rating", "desc").Return(expectedMovies, nil)

	handler.getAllMovies(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var responseMovies []model.Movie
	err := json.NewDecoder(w.Body).Decode(&responseMovies)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	if len(responseMovies) != len(expectedMovies) {
		t.Errorf("Expected %d movies, got %d", len(expectedMovies), len(responseMovies))
	}
	for i, movie := range expectedMovies {
		if movie.Title != responseMovies[i].Title || movie.Description != responseMovies[i].Description || movie.ReleaseDate != responseMovies[i].ReleaseDate || movie.Rating != responseMovies[i].Rating {
			t.Errorf("Expected movie %+v, got %+v", movie, responseMovies[i])
		}
	}
}

func TestHandler_getAllMovies_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMovieService := mock_service.NewMockMovie(ctrl)

	handler := Handler{
		&service.Service{
			Movie: mockMovieService,
		},
	}

	req := httptest.NewRequest("GET", "/api/movies", nil)
	w := httptest.NewRecorder()

	mockMovieService.EXPECT().GetAllMovies("rating", "desc").Return(nil, errors.New("service error"))

	handler.getAllMovies(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	expectedError := "Failed to get movies"
	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
	if errorMessage, ok := response["message"]; !ok || errorMessage != expectedError {
		t.Errorf("Expected error message %q, got %q", expectedError, errorMessage)
	}
}

func TestHandler_createMovie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок сервиса Movie
	mockMovieService := mock_service.NewMockMovie(ctrl)
	mockMovieService.EXPECT().CreateMovie(gomock.Any()).Return(nil)

	// Создаем Handler
	handler := &Handler{
		services: &service.Service{
			Movie: mockMovieService,
		},
	}

	// Создаем тестируемый хендлер
	reqBody := `{"title":"Test Movie", "actors":[{"name":"Actor 1", "gender":"female", "birth_date":"2003-9-2"}]}`
	req := httptest.NewRequest("POST", "/api/movie", strings.NewReader(reqBody))
	ctx := context.WithValue(req.Context(), userRoleCtx, "admin") // Устанавливаем роль в контексте
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Вызываем тестируемый хендлер
	handler.createMovie(w, req)

	// Проверяем код ответа
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Проверяем тело ответа
	expectedResponse := "movie created successfully"
	if w.Body.String() != expectedResponse {
		t.Errorf("Expected response body %q, got %q", expectedResponse, w.Body.String())
	}
}

func TestHandler_createMovie_UnsuccessfulCreation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок сервиса Movie
	mockMovieService := mock_service.NewMockMovie(ctrl)
	mockMovieService.EXPECT().CreateMovie(gomock.Any()).Return(errors.New("failed to create movie"))

	// Создаем Handler
	handler := &Handler{
		services: &service.Service{
			Movie: mockMovieService,
		},
	}

	// Создаем тестируемый хендлер
	reqBody := `{"title":"Test Movie", "actors":[{"name":"Actor 1", "gender":"female", "birth_date":"2003-9-2"}]}`
	req := httptest.NewRequest("POST", "/api/movie", strings.NewReader(reqBody))
	ctx := context.WithValue(req.Context(), userRoleCtx, "admin") // Устанавливаем роль в контексте
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Вызываем тестируемый хендлер
	handler.createMovie(w, req)

	// Проверяем код ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	// Проверяем тело ответа
	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response body: %v", err)
	}

	expectedResponse := "failed to create movie"
	if response["message"] != expectedResponse {
		t.Errorf("Expected response body %q, got %q", expectedResponse, response["message"])
	}
}

func TestHandler_createMovie_ForbiddenRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок сервиса Movie
	mockMovieService := mock_service.NewMockMovie(ctrl)

	// Создаем Handler
	handler := &Handler{
		services: &service.Service{
			Movie: mockMovieService,
		},
	}

	// Создаем тестируемый хендлер
	reqBody := `{"title":"Test Movie", "actors":[{"name":"Actor 1", "gender":"female", "birth_date":"2003-9-2"}]}`
	req := httptest.NewRequest("POST", "/api/movie", strings.NewReader(reqBody))
	ctx := context.WithValue(req.Context(), userRoleCtx, "user") // Устанавливаем роль пользователя в контексте
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Вызываем тестируемый хендлер
	handler.createMovie(w, req)

	// Проверяем код ответа
	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status code %d, got %d", http.StatusForbidden, w.Code)
	}

	// Проверяем тело ответа
	expectedResponse := "only admin can create movies"
	if !strings.Contains(w.Body.String(), expectedResponse) {
		t.Errorf("Expected response body to contain %q, got %q", expectedResponse, w.Body.String())
	}
}

func TestHandler_deleteMovie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок сервиса Movie
	mockMovieService := mock_service.NewMockMovie(ctrl)
	mockMovieService.EXPECT().DeleteByID(1).Return(nil) // Ожидаем вызов DeleteByID с аргументом 1

	// Создаем Handler
	handler := &Handler{
		services: &service.Service{
			Movie: mockMovieService,
		},
	}

	// Создаем тестируемый хендлер
	req := httptest.NewRequest("DELETE", "/movie/1", nil)         // Используем метод DELETE
	ctx := context.WithValue(req.Context(), userRoleCtx, "admin") // Устанавливаем роль в контексте
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Вызываем тестируемый хендлер
	handler.deleteMovie(w, req)

	// Проверяем код ответа
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Проверяем тело ответа
	expectedResponse := "movie deleted successfully"
	if w.Body.String() != expectedResponse {
		t.Errorf("Expected response body %q, got %q", expectedResponse, w.Body.String())
	}
}

func TestHandler_deleteMovie_Unsuccessful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок сервиса Movie
	mockMovieService := mock_service.NewMockMovie(ctrl)
	mockMovieService.EXPECT().DeleteByID(1).Return(errors.New("Failed to delete movie by ID")) // Ожидаем вызов DeleteByID с аргументом 1 и возвращаем ошибку

	// Создаем Handler
	handler := &Handler{
		services: &service.Service{
			Movie: mockMovieService,
		},
	}

	// Создаем тестируемый хендлер
	req := httptest.NewRequest("DELETE", "/movie/1", nil)         // Используем метод DELETE
	ctx := context.WithValue(req.Context(), userRoleCtx, "admin") // Устанавливаем роль в контексте
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Вызываем тестируемый хендлер
	handler.deleteMovie(w, req)

	// Проверяем код ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	// Проверяем тело ответа
	expectedResponse := "{\"message\":\"Failed to delete movie by ID\"}"
	actualResponse := strings.Trim(w.Body.String(), `"`) // Убираем кавычки из ответа
	if actualResponse != expectedResponse {
		t.Errorf("Expected response body %q, got %q", expectedResponse, actualResponse)
	}
}

func TestHandler_deleteMovie_NotAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMovieService := mock_service.NewMockMovie(ctrl)

	handler := &Handler{
		services: &service.Service{
			Movie: mockMovieService,
		},
	}

	req := httptest.NewRequest("DELETE", "/movie/1", nil)
	ctx := context.WithValue(req.Context(), userRoleCtx, "user") // Устанавливаем не администратора в контексте
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.deleteMovie(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status code %d, got %d", http.StatusForbidden, w.Code)
	}

	expectedResponse := "{\"message\":\"only admin can delete movies\"}"
	if w.Body.String() != expectedResponse {
		t.Errorf("Expected response body %q, got %q", expectedResponse, w.Body.String())
	}
}

func TestHandler_deleteMovie_InvalidMovieID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMovieService := mock_service.NewMockMovie(ctrl)

	handler := &Handler{
		services: &service.Service{
			Movie: mockMovieService,
		},
	}

	req := httptest.NewRequest("DELETE", "/movie/invalid_id", nil) // Используем некорректный формат идентификатора фильма
	ctx := context.WithValue(req.Context(), userRoleCtx, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.movieHandle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	expectedResponse := "{\"message\":\"Invalid movie ID\"}"
	if w.Body.String() != expectedResponse {
		t.Errorf("Expected response body %q, got %q", expectedResponse, w.Body.String())
	}
}

func TestHandler_updateMovie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок сервиса Movie
	mockMovieService := mock_service.NewMockMovie(ctrl)
	mockMovieService.EXPECT().UpdateMovie(gomock.Any(), gomock.Any()).Return(nil) // Ожидаем вызов UpdateMovie с любыми аргументами и возвращаем nil

	// Создаем Handler
	handler := &Handler{
		services: &service.Service{
			Movie: mockMovieService,
		},
	}

	// Создаем тестируемый хендлер
	reqBody := `{"title":"Updated Movie", "actors":[{"name":"Actor 1", "gender":"female", "birth_date":"2003-9-2"}]}`
	req := httptest.NewRequest("PUT", "/movie/1", strings.NewReader(reqBody)) // Используем метод PUT
	ctx := context.WithValue(req.Context(), userRoleCtx, "admin")             // Устанавливаем роль в контексте
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Вызываем тестируемый хендлер
	handler.movieHandle(w, req)

	// Проверяем код ответа
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Проверяем тело ответа
	expectedResponse := "movie updated successfully"
	if w.Body.String() != expectedResponse {
		t.Errorf("Expected response body %q, got %q", expectedResponse, w.Body.String())
	}
}

func TestHandler_updateMovie_Unsuccessful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок сервиса Movie
	mockMovieService := mock_service.NewMockMovie(ctrl)
	mockMovieService.EXPECT().UpdateMovie(gomock.Any(), gomock.Any()).Return(errors.New("Failed to update movie")) // Ожидаем вызов UpdateMovie с любыми аргументами и возвращаем ошибку

	// Создаем Handler
	handler := &Handler{
		services: &service.Service{
			Movie: mockMovieService,
		},
	}

	// Создаем тестируемый хендлер
	reqBody := `{"title":"Updated Movie", "actors":[{"name":"Actor 1", "gender":"female", "birth_date":"2003-9-2"}]}`
	req := httptest.NewRequest("PUT", "/movie/1", strings.NewReader(reqBody)) // Используем метод PUT
	ctx := context.WithValue(req.Context(), userRoleCtx, "admin")             // Устанавливаем роль в контексте
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Вызываем тестируемый хендлер
	handler.updateMovie(w, req)

	// Проверяем код ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	// Проверяем тело ответа
	expectedResponse := "{\"message\":\"Failed to update movie\"}"
	if w.Body.String() != expectedResponse {
		t.Errorf("Expected response body %q, got %q", expectedResponse, w.Body.String())
	}
}

func TestHandler_getMovie_Unsuccessful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок сервиса Movie
	mockMovieService := mock_service.NewMockMovie(ctrl)
	mockMovieService.EXPECT().GetMovieByID(1).Return(model.MovieWithActors{}, errors.New("Failed to get movie")) // Ожидаем вызов GetMovieByID с аргументом 1 и возвращаем ошибку

	// Создаем Handler
	handler := &Handler{
		services: &service.Service{
			Movie: mockMovieService,
		},
	}

	// Создаем тестируемый хендлер
	req := httptest.NewRequest("GET", "/movie/1", nil) // Используем метод GET
	w := httptest.NewRecorder()

	// Вызываем тестируемый хендлер
	handler.getMovie(w, req)

	// Проверяем код ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	// Проверяем тело ответа
	expectedResponse := "{\"message\":\"Failed to get movie\"}"
	if w.Body.String() != expectedResponse {
		t.Errorf("Expected response body %q, got %q", expectedResponse, w.Body.String())
	}
}

func TestHandler_getMovie_Successful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок сервиса Movie
	mockMovieService := mock_service.NewMockMovie(ctrl)
	expectedMovie := &model.MovieWithActors{
		ID:     1,
		Title:  "Test Movie",
		Actors: []model.Actor{},
	}
	mockMovieService.EXPECT().GetMovieByID(1).Return(*expectedMovie, nil) // Ожидаем вызов GetMovieByID с аргументом 1 и возвращаем ожидаемый фильм

	// Создаем Handler
	handler := &Handler{
		services: &service.Service{
			Movie: mockMovieService,
		},
	}

	// Создаем тестируемый хендлер
	req := httptest.NewRequest("GET", "/movie/1", nil) // Используем метод GET
	w := httptest.NewRecorder()

	// Вызываем тестируемый хендлер
	handler.movieHandle(w, req)

	// Проверяем код ответа
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Проверяем тело ответа
	expectedResponse, err := json.Marshal(expectedMovie)
	if err != nil {
		t.Errorf("Error marshaling expected movie: %v", err)
	}
	expectedResponseString := string(expectedResponse) + "\n" // Добавляем новую строку в конце, чтобы совпадать с форматом тела ответа
	if w.Body.String() != expectedResponseString {
		t.Errorf("Expected response body %q, got %q", expectedResponseString, w.Body.String())
	}
}

func TestHandler_searchMovie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок сервиса Movie
	mockMovieService := mock_service.NewMockMovie(ctrl)
	expectedMovies := []model.MovieWithActors{
		{
			Title:       "Movie 1",
			Description: "Description 1",
			ReleaseDate: "2022-01-01",
			Rating:      8,
			Actors: []model.Actor{
				{Name: "Actor 1"},
			},
		},
		{
			Title:       "Movie 2",
			Description: "Description 2",
			ReleaseDate: "2022-02-01",
			Rating:      7,
			Actors: []model.Actor{
				{Name: "Actor 2"},
			},
		},
	}
	mockMovieService.EXPECT().GetMoviesByTitle("Movie").Return(expectedMovies, nil) // Ожидаем вызов GetMoviesByTitle с аргументом "Movie" и возвращаем ожидаемые фильмы

	// Создаем Handler
	handler := &Handler{
		services: &service.Service{
			Movie: mockMovieService,
		},
	}

	// Создаем тестируемый хендлер
	req := httptest.NewRequest("GET", "/movie?title=Movie", nil) // Используем метод GET с параметром title
	w := httptest.NewRecorder()

	// Вызываем тестируемый хендлер
	handler.searchMovie(w, req)

	// Проверяем код ответа
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Проверяем тело ответа
	expectedResponse, err := json.Marshal(expectedMovies)
	if err != nil {
		t.Errorf("Error marshaling expected movie: %v", err)
	}
	if w.Body.String() != string(expectedResponse)+"\n" {
		t.Errorf("Expected response body %q, got %q", expectedResponse, w.Body.String())
	}
}

func TestHandler_searchMovie_EmptyParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMovieService := mock_service.NewMockMovie(ctrl)

	handler := &Handler{
		services: &service.Service{
			Movie: mockMovieService,
		},
	}

	// Проверка случая, когда оба параметра пустые строки
	req := httptest.NewRequest("GET", "/movie", nil)
	w := httptest.NewRecorder()

	handler.searchMovie(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d for empty params, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandler_searchMovie_InvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMovieService := mock_service.NewMockMovie(ctrl)

	handler := &Handler{
		services: &service.Service{
			Movie: mockMovieService,
		},
	}

	// Проверка случая, когда оба параметра указаны
	req := httptest.NewRequest("GET", "/movie?title=Movie&actor=Actor", nil)
	w := httptest.NewRecorder()

	handler.searchMovie(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d for invalid request, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandler_updateMovie_MethodNotAllowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок сервиса Movie
	mockMovieService := mock_service.NewMockMovie(ctrl)

	// Создаем Handler
	handler := &Handler{
		services: &service.Service{
			Movie: mockMovieService,
		},
	}

	// Создаем тестируемый хендлер
	req := httptest.NewRequest("POST", "/api/movie/1", nil)       // Используем метод POST, который не поддерживается для обновления фильма
	ctx := context.WithValue(req.Context(), userRoleCtx, "admin") // Устанавливаем роль в контексте
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Вызываем тестируемый хендлер
	handler.updateMovie(w, req)

	// Проверяем код ответа
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}

	// Проверяем тело ответа
	expectedResponse := "{\"message\":\"Method not allowed\"}"
	if w.Body.String() != expectedResponse {
		t.Errorf("Expected response body %q, got %q", expectedResponse, w.Body.String())
	}
}

func TestHandler_updateMovie_InvalidMovieID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок сервиса Movie
	mockMovieService := mock_service.NewMockMovie(ctrl)

	// Создаем Handler
	handler := &Handler{
		services: &service.Service{
			Movie: mockMovieService,
		},
	}

	// Создаем тестируемый хендлер
	reqBody := `{"title":"Updated Movie", "actors":[{"name":"Actor 1", "gender":"female", "birth_date":"2003-9-2"}]}`
	req := httptest.NewRequest("PUT", "/movie/invalid_id", strings.NewReader(reqBody)) // Передаем некорректный идентификатор фильма
	ctx := context.WithValue(req.Context(), userRoleCtx, "admin")                      // Устанавливаем роль в контексте
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Вызываем тестируемый хендлер
	handler.updateMovie(w, req)

	// Проверяем код ответа
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Проверяем тело ответа
	expectedResponse := "{\"message\":\"Invalid movie ID\"}"
	if w.Body.String() != expectedResponse {
		t.Errorf("Expected response body %q, got %q", expectedResponse, w.Body.String())
	}
}
