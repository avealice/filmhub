package handler

import (
	"bytes"
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

func normalizeJSON(input string) string {
	input = strings.ReplaceAll(input, "\n", "")
	return input
}

func TestHandler_signUp_Success(t *testing.T) {
	// Создаем мок сервиса Authorization
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_service.NewMockAuthorization(ctrl)

	// Определяем ожидаемые входные данные и выходные данные
	input := model.User{
		Username: "test",
		Password: "qwerty",
	}
	expectedID := 1

	// Устанавливаем ожидаемое поведение мок сервиса Authorization
	mockAuthService.EXPECT().CreateUser(input).Return(expectedID, nil)

	// Создаем Handler
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuthService,
		},
	}

	// Создаем фейковый запрос
	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/auth/sign-up", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Вызываем тестируемый хендлер
	handler.signUp(w, req)

	// Проверяем код ответа
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Проверяем тело ответа
	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	// Проверяем наличие ID пользователя в ответе
	id, ok := response["id"].(float64)
	if !ok {
		t.Error("ID not found in response")
	}

	// Проверяем, что ID пользователя соответствует ожидаемому значению
	if int(id) != expectedID {
		t.Errorf("Expected ID %d, got %d", expectedID, int(id))
	}
}

func TestHandler_signUp_InvalidInput(t *testing.T) {
	// Создаем мок сервиса Authorization
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_service.NewMockAuthorization(ctrl)

	// Создаем Handler
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuthService,
		},
	}

	// Ожидаем вызов CreateUser с пустым пользователем и возвратом ошибки
	mockAuthService.EXPECT().CreateUser(gomock.Any()).Return(0, errors.New("username and password cannot be empty")).Times(1)

	// Создаем фейковый запрос с невалидными данными
	body := []byte(`{"username": "", "password": ""}`)
	req := httptest.NewRequest("POST", "/auth/sign-up", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Вызываем тестируемый хендлер
	handler.signUp(w, req)

	// Проверяем код ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Проверяем тело ответа
	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	// Проверяем наличие сообщения об ошибке в ответе
	if _, ok := response["message"]; !ok {
		t.Error("Error message not found in response")
	}
}

func TestHandler_signUp_InternalServerError(t *testing.T) {
	// Создаем мок сервиса Authorization
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_service.NewMockAuthorization(ctrl)

	// Определяем ожидаемые входные данные
	input := model.User{
		Username: "test",
		Password: "qwerty",
	}

	// Устанавливаем ожидаемое поведение мок сервиса Authorization
	mockAuthService.EXPECT().CreateUser(input).Return(0, errors.New("internal error"))

	// Создаем Handler
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuthService,
		},
	}

	// Создаем фейковый запрос
	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/auth/sign-up", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Вызываем тестируемый хендлер
	handler.signUp(w, req)

	// Проверяем код ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	// Проверяем тело ответа
	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	// Проверяем наличие сообщения об ошибке в ответе
	if _, ok := response["message"]; !ok {
		t.Error("Error message not found in response")
	}
}

func TestHandler_signIn_Success(t *testing.T) {
	// Создаем мок сервиса Authorization
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_service.NewMockAuthorization(ctrl)

	// Определяем ожидаемые входные данные
	input := SignInInput{
		Username: "test",
		Password: "qwerty",
	}

	// Определяем ожидаемый токен
	expectedToken := "test_token"

	// Устанавливаем ожидаемое поведение мок сервиса Authorization
	mockAuthService.EXPECT().GenerateToken(input.Username, input.Password).Return(expectedToken, nil)

	// Создаем Handler
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuthService,
		},
	}

	// Создаем фейковый запрос
	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/auth/sign-in", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Вызываем тестируемый хендлер
	handler.signIn(w, req)

	// Проверяем код ответа
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Проверяем тело ответа
	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	// Проверяем наличие токена в ответе
	if _, ok := response["token"]; !ok {
		t.Error("Token not found in response")
	}

	// Проверяем соответствие полученного токена ожидаемому
	if response["token"].(string) != expectedToken {
		t.Errorf("Expected token %s, got %s", expectedToken, response["token"].(string))
	}
}

func TestHandler_signIn_BadRequest(t *testing.T) {
	// Создаем мок сервиса Authorization
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_service.NewMockAuthorization(ctrl)

	// Создаем Handler
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuthService,
		},
	}

	// Создаем фейковый запрос с некорректным телом
	req := httptest.NewRequest("POST", "/auth/sign-in", strings.NewReader("{"))
	w := httptest.NewRecorder()

	// Вызываем тестируемый хендлер
	handler.signIn(w, req)

	// Проверяем код ответа
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandler_signIn_InternalServerError(t *testing.T) {
	// Создаем мок сервиса Authorization
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_service.NewMockAuthorization(ctrl)

	// Определяем ожидаемые входные данные
	input := SignInInput{
		Username: "test",
		Password: "qwerty",
	}

	// Устанавливаем ожидаемое поведение мок сервиса Authorization
	mockAuthService.EXPECT().GenerateToken(input.Username, input.Password).Return("", errors.New("internal error"))

	// Создаем Handler
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuthService,
		},
	}

	// Создаем фейковый запрос
	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/auth/sign-in", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Вызываем тестируемый хендлер
	handler.signIn(w, req)

	// Проверяем код ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	// Проверяем тело ответа
	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	// Проверяем наличие сообщения об ошибке в ответе
	if _, ok := response["message"]; !ok {
		t.Error("Error message not found in response")
	}
}
