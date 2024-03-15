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

func TestHandler_getAllActors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockActorService := mock_service.NewMockActor(ctrl)

	handler := &Handler{
		services: &service.Service{
			Actor: mockActorService,
		},
	}

	expectedActors := []model.Actor{
		{ID: 1, Name: "Actor 1"},
		{ID: 2, Name: "Actor 2"},
	}

	// Change the type of expectedActors to []model.ActorWithMovies
	expectedActorsWithMovies := make([]model.ActorWithMovies, len(expectedActors))
	for i, actor := range expectedActors {
		expectedActorsWithMovies[i] = model.ActorWithMovies{
			ID:   actor.ID,
			Name: actor.Name,
		}
	}

	mockActorService.EXPECT().GetAllActors().Return(expectedActorsWithMovies, nil)

	req := httptest.NewRequest("GET", "/api/actors", nil)
	w := httptest.NewRecorder()

	handler.getAllActors(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var responseActors []model.ActorWithMovies
	if err := json.NewDecoder(w.Body).Decode(&responseActors); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	if len(responseActors) != len(expectedActors) {
		t.Errorf("Expected %d actors, got %d", len(expectedActors), len(responseActors))
	}
}

func TestHandler_getAllActors_MethodNotAllowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockActorService := mock_service.NewMockActor(ctrl)

	handler := &Handler{
		services: &service.Service{
			Actor: mockActorService,
		},
	}

	req := httptest.NewRequest("POST", "/api/actors", nil) // Use POST method
	w := httptest.NewRecorder()

	handler.getAllActors(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestHandler_getAllActors_InternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockActorService := mock_service.NewMockActor(ctrl)

	handler := &Handler{
		services: &service.Service{
			Actor: mockActorService,
		},
	}

	// Simulate internal server error by making the mock service return an error
	mockActorService.EXPECT().GetAllActors().Return(nil, errors.New("database error"))

	req := httptest.NewRequest("GET", "/api/actors", nil) // Use GET method
	w := httptest.NewRecorder()

	handler.getAllActors(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestHandler_CreateActor_Successful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок сервиса Actor
	mockActorService := mock_service.NewMockActor(ctrl)
	mockActorService.EXPECT().CreateActor(gomock.Any()).Return(nil) // Ожидаем вызов CreateActor с любыми аргументами и возвращаем nil

	// Создаем Handler
	handler := &Handler{
		services: &service.Service{
			Actor: mockActorService,
		},
	}

	// Создаем тестируемый хендлер
	reqBody := `{"name":"Test Actor", "gender":"male", "birth_date":"2000-01-01"}`
	req := httptest.NewRequest("POST", "/actor", strings.NewReader(reqBody)) // Используем метод POST
	ctx := context.WithValue(req.Context(), userRoleCtx, "admin")            // Устанавливаем роль в контексте
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Вызываем тестируемый хендлер
	handler.CreateActor(w, req)

	// Проверяем код ответа
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Проверяем тело ответа
	expectedResponse := "actor created successfully"
	if w.Body.String() != expectedResponse {
		t.Errorf("Expected response body %q, got %q", expectedResponse, w.Body.String())
	}
}

func TestHandler_CreateActor_Unsuccessful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок сервиса Actor
	mockActorService := mock_service.NewMockActor(ctrl)
	mockActorService.EXPECT().CreateActor(gomock.Any()).Return(errors.New("Failed to create actor")) // Ожидаем вызов CreateActor с любыми аргументами и возвращаем ошибку

	// Создаем Handler
	handler := &Handler{
		services: &service.Service{
			Actor: mockActorService,
		},
	}

	// Создаем тестируемый хендлер
	reqBody := `{"name":"Test Actor", "gender":"male", "birth_date":"2000-01-01"}`
	req := httptest.NewRequest("POST", "/actor", strings.NewReader(reqBody)) // Используем метод POST
	ctx := context.WithValue(req.Context(), userRoleCtx, "admin")            // Устанавливаем роль в контексте
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Вызываем тестируемый хендлер
	handler.CreateActor(w, req)

	// Проверяем код ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	// Проверяем тело ответа
	expectedResponse := "{\"message\":\"Failed to create actor\"}"
	if w.Body.String() != expectedResponse {
		t.Errorf("Expected response body %q, got %q", expectedResponse, w.Body.String())
	}
}

func TestHandler_CreateActor_Unauthorized(t *testing.T) {
	// Создаем Handler
	handler := &Handler{}

	// Создаем тестируемый хендлер без установки роли в контексте
	req := httptest.NewRequest("POST", "/actor", nil)            // Используем метод DELETE
	ctx := context.WithValue(req.Context(), userRoleCtx, "user") // Устанавливаем роль в контексте
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Вызываем тестируемый хендлер
	handler.CreateActor(w, req)

	// Проверяем код ответа
	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status code %d, got %d", http.StatusForbidden, w.Code)
	}

	// Проверяем тело ответа
	expectedResponse := "{\"message\":\"only admin can create actors\"}"
	if w.Body.String() != expectedResponse {
		t.Errorf("Expected response body %q, got %q", expectedResponse, w.Body.String())
	}
}

func TestHandler_CreateActor_BadRequest(t *testing.T) {
	// Создаем Handler
	handler := &Handler{}

	// Создаем тестируемый хендлер с неправильным форматом JSON
	req := httptest.NewRequest("POST", "/actor", strings.NewReader("invalid json"))
	ctx := context.WithValue(req.Context(), userRoleCtx, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Вызываем тестируемый хендлер
	handler.CreateActor(w, req)

	// Проверяем код ответа
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Проверяем тело ответа
	expectedResponse := "{\"message\":\"invalid character 'i' looking for beginning of value\"}"
	if w.Body.String() != expectedResponse {
		t.Errorf("Expected response body %q, got %q", expectedResponse, w.Body.String())
	}
}

func TestHandler_deleteActor_Successful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок сервиса Actor
	mockActorService := mock_service.NewMockActor(ctrl)
	mockActorService.EXPECT().Delete(1).Return(nil) // Ожидаем вызов Delete с аргументом 1 и возвращаем nil

	// Создаем Handler
	handler := &Handler{
		services: &service.Service{
			Actor: mockActorService,
		},
	}

	// Создаем тестируемый хендлер
	req := httptest.NewRequest("DELETE", "/actor/1", nil)         // Используем метод DELETE
	ctx := context.WithValue(req.Context(), userRoleCtx, "admin") // Устанавливаем роль в контексте
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Вызываем тестируемый хендлер
	handler.deleteActor(w, req)

	// Проверяем код ответа
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Проверяем тело ответа
	expectedResponse := "actor deleted successfully"
	if w.Body.String() != expectedResponse {
		t.Errorf("Expected response body %q, got %q", expectedResponse, w.Body.String())
	}
}

func TestHandler_deleteActor_Unsuccessful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок сервиса Actor
	mockActorService := mock_service.NewMockActor(ctrl)
	mockActorService.EXPECT().Delete(1).Return(errors.New("Failed to delete actor")) // Ожидаем вызов Delete с аргументом 1 и возвращаем ошибку

	// Создаем Handler
	handler := &Handler{
		services: &service.Service{
			Actor: mockActorService,
		},
	}

	// Создаем тестируемый хендлер
	req := httptest.NewRequest("DELETE", "/actor/1", nil)         // Используем метод DELETE
	ctx := context.WithValue(req.Context(), userRoleCtx, "admin") // Устанавливаем роль в контексте
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Вызываем тестируемый хендлер
	handler.deleteActor(w, req)

	// Проверяем код ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	// Проверяем тело ответа
	expectedResponse := "{\"message\":\"Failed to delete actor\"}"
	if w.Body.String() != expectedResponse {
		t.Errorf("Expected response body %q, got %q", expectedResponse, w.Body.String())
	}
}

func TestHandler_deleteActor_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок сервиса Actor
	mockActorService := mock_service.NewMockActor(ctrl)

	// Создаем Handler
	handler := &Handler{
		services: &service.Service{
			Actor: mockActorService,
		},
	}

	// Создаем тестируемый хендлер
	req := httptest.NewRequest("DELETE", "/actor/notanumber", nil) // Используем метод DELETE с неверным идентификатором актера
	ctx := context.WithValue(req.Context(), userRoleCtx, "admin")  // Устанавливаем роль в контексте
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Вызываем тестируемый хендлер
	handler.deleteActor(w, req)

	// Проверяем код ответа
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Проверяем тело ответа
	expectedResponse := "{\"message\":\"Invalid actor ID\"}"
	if w.Body.String() != expectedResponse {
		t.Errorf("Expected response body %q, got %q", expectedResponse, w.Body.String())
	}
}

func TestHandler_updateActor_Successful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок сервиса Actor
	mockActorService := mock_service.NewMockActor(ctrl)
	mockActorService.EXPECT().Update(1, gomock.Any()).Return(nil) // Ожидаем вызов Update с аргументом 1 и любым значением ActorWithMovies и возвращаем nil

	// Создаем Handler
	handler := &Handler{
		services: &service.Service{
			Actor: mockActorService,
		},
	}

	// Создаем тестируемый хендлер
	actorData := `{"name":"Updated Actor","gender":"male","birth_date":"1990-01-01","movies":[]}`
	req := httptest.NewRequest("PUT", "/actor/1", strings.NewReader(actorData)) // Используем метод PUT
	ctx := context.WithValue(req.Context(), userRoleCtx, "admin")               // Устанавливаем роль в контексте
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Вызываем тестируемый хендлер
	handler.updateActor(w, req)

	// Проверяем код ответа
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Проверяем тело ответа
	expectedResponse := "actor updated successfully"
	if w.Body.String() != expectedResponse {
		t.Errorf("Expected response body %q, got %q", expectedResponse, w.Body.String())
	}
}

func TestHandler_updateActor_Unsuccessful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок сервиса Actor
	mockActorService := mock_service.NewMockActor(ctrl)
	mockActorService.EXPECT().Update(1, gomock.Any()).Return(errors.New("Failed to update actor")) // Ожидаем вызов Update с аргументом 1 и любым значением ActorWithMovies и возвращаем ошибку

	// Создаем Handler
	handler := &Handler{
		services: &service.Service{
			Actor: mockActorService,
		},
	}

	// Создаем тестируемый хендлер
	actorData := `{"name":"Updated Actor","gender":"male","birth_date":"1990-01-01","movies":[]}`
	req := httptest.NewRequest("PUT", "/actor/1", strings.NewReader(actorData)) // Используем метод PUT
	ctx := context.WithValue(req.Context(), userRoleCtx, "admin")               // Устанавливаем роль в контексте
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Вызываем тестируемый хендлер
	handler.updateActor(w, req)

	// Проверяем код ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	// Проверяем тело ответа
	expectedResponse := "{\"message\":\"Failed to update actor\"}"
	if w.Body.String() != expectedResponse {
		t.Errorf("Expected response body %q, got %q", expectedResponse, w.Body.String())
	}
}

func TestHandler_updateActor_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок сервиса Actor
	mockActorService := mock_service.NewMockActor(ctrl)

	// Создаем Handler
	handler := &Handler{
		services: &service.Service{
			Actor: mockActorService,
		},
	}

	// Создаем тестируемый хендлер с неверным идентификатором актера
	actorData := `{"name":"Updated Actor","gender":"male","birth_date":"1990-01-01","movies":[]}`
	req := httptest.NewRequest("PUT", "/actor/notanumber", strings.NewReader(actorData))
	ctx := context.WithValue(req.Context(), userRoleCtx, "admin") // Устанавливаем роль в контексте
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Вызываем тестируемый хендлер
	handler.updateActor(w, req)

	// Проверяем код ответа
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Проверяем тело ответа
	expectedResponse := "{\"message\":\"Invalid actor ID\"}"
	if w.Body.String() != expectedResponse {
		t.Errorf("Expected response body %q, got %q", expectedResponse, w.Body.String())
	}
}
func TestHandler_getActor_Successful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок сервиса Actor
	mockActorService := mock_service.NewMockActor(ctrl)
	expectedActor := &model.ActorWithMovies{
		ID:        1,
		Name:      "Test Actor",
		Gender:    "male",
		BirthDate: "1990-01-01",
	}
	mockActorService.EXPECT().Get(1).Return(*expectedActor, nil) // Ожидаем вызов Get с аргументом 1 и возвращаем ожидаемого актера

	// Создаем Handler
	handler := &Handler{
		services: &service.Service{
			Actor: mockActorService,
		},
	}

	// Создаем тестируемый хендлер
	req := httptest.NewRequest("GET", "/actor/1", nil) // Используем метод GET
	w := httptest.NewRecorder()

	// Вызываем тестируемый хендлер
	handler.getActor(w, req)

	// Проверяем код ответа
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Проверяем тело ответа
	expectedResponse, err := json.Marshal(expectedActor)
	if err != nil {
		t.Errorf("Error marshaling expected actor: %v", err)
	}
	if w.Body.String() != string(expectedResponse)+"\n" {
		t.Errorf("Expected response body %q, got %q", expectedResponse, w.Body.String())
	}
}

func TestHandler_getActor_Unsuccessful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок сервиса Actor
	mockActorService := mock_service.NewMockActor(ctrl)
	mockActorService.EXPECT().Get(1).Return(model.ActorWithMovies{}, errors.New("Failed to get actor")) // Ожидаем вызов Get с аргументом 1 и возвращаем ошибку

	// Создаем Handler
	handler := &Handler{
		services: &service.Service{
			Actor: mockActorService,
		},
	}

	// Создаем тестируемый хендлер
	req := httptest.NewRequest("GET", "/actor/1", nil) // Используем метод GET
	w := httptest.NewRecorder()

	// Вызываем тестируемый хендлер
	handler.getActor(w, req)

	// Проверяем код ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	// Проверяем тело ответа
	expectedResponse := "{\"message\":\"Failed to get actor\"}"
	if w.Body.String() != expectedResponse {
		t.Errorf("Expected response body %q, got %q", expectedResponse, w.Body.String())
	}
}

func TestHandler_getActor_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок сервиса Actor
	mockActorService := mock_service.NewMockActor(ctrl)

	// Создаем Handler
	handler := &Handler{
		services: &service.Service{
			Actor: mockActorService,
		},
	}

	// Создаем тестируемый хендлер с неверным идентификатором актера
	req := httptest.NewRequest("GET", "/actor/notanumber", nil)
	w := httptest.NewRecorder()

	// Вызываем тестируемый хендлер
	handler.getActor(w, req)

	// Проверяем код ответа
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Проверяем тело ответа
	expectedResponse := "{\"message\":\"Invalid actor ID\"}"
	if w.Body.String() != expectedResponse {
		t.Errorf("Expected response body %q, got %q", expectedResponse, w.Body.String())
	}
}
