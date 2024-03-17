package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/avealice/filmhub/internal/model"
	"github.com/avealice/filmhub/internal/service"

	mock_service "github.com/avealice/filmhub/internal/service/mocks"

	"github.com/golang/mock/gomock"
)

func normalizeJSON(input string) string {
	input = strings.ReplaceAll(input, "\n", "")
	return input
}

func TestHandler_signUp_Success(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_service.NewMockAuthorization(ctrl)

	input := model.User{
		Username: "test",
		Password: "qwerty",
	}
	expectedID := 1

	mockAuthService.EXPECT().CreateUser(input).Return(expectedID, nil)

	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuthService,
		},
	}

	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/auth/sign-up", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.signUp(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	id, ok := response["id"].(float64)
	if !ok {
		t.Error("ID not found in response")
	}

	if int(id) != expectedID {
		t.Errorf("Expected ID %d, got %d", expectedID, int(id))
	}
}

func TestHandler_signUp_MethodNotAllowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_service.NewMockAuthorization(ctrl)

	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuthService,
		},
	}

	req := httptest.NewRequest("GET", "/auth/sign-up", nil)
	w := httptest.NewRecorder()

	handler.signUp(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}

	expectedResponse := "{\"message\":\"Method not allowed\"}"
	if w.Body.String() != expectedResponse {
		t.Errorf("Expected response body %q, got %q", expectedResponse, w.Body.String())
	}
}

func TestHandler_signUp_InternalServerError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_service.NewMockAuthorization(ctrl)

	input := model.User{
		Username: "test",
		Password: "qwerty",
	}

	mockAuthService.EXPECT().CreateUser(input).Return(0, errors.New("internal error"))

	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuthService,
		},
	}

	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/auth/sign-up", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.signUp(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	if _, ok := response["message"]; !ok {
		t.Error("Error message not found in response")
	}
}

func TestHandler_signIn_Success(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_service.NewMockAuthorization(ctrl)

	input := SignInInput{
		Username: "test",
		Password: "qwerty",
	}

	expectedToken := "test_token"

	mockAuthService.EXPECT().GenerateToken(input.Username, input.Password).Return(expectedToken, nil)

	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuthService,
		},
	}

	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/auth/sign-in", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.signIn(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	if _, ok := response["token"]; !ok {
		t.Error("Token not found in response")
	}

	if response["token"].(string) != expectedToken {
		t.Errorf("Expected token %s, got %s", expectedToken, response["token"].(string))
	}
}

func TestHandler_signIn_BadRequest(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_service.NewMockAuthorization(ctrl)

	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuthService,
		},
	}

	req := httptest.NewRequest("POST", "/auth/sign-in", strings.NewReader("{"))
	w := httptest.NewRecorder()

	handler.signIn(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandler_signIn_InternalServerError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_service.NewMockAuthorization(ctrl)

	input := SignInInput{
		Username: "test",
		Password: "qwerty",
	}

	mockAuthService.EXPECT().GenerateToken(input.Username, input.Password).Return("", errors.New("internal error"))

	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuthService,
		},
	}

	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/auth/sign-in", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.signIn(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	if _, ok := response["message"]; !ok {
		t.Error("Error message not found in response")
	}
}

func TestHandler_signIn_MethodNotAllowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_service.NewMockAuthorization(ctrl)

	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuthService,
		},
	}

	req := httptest.NewRequest("GET", "/auth/sign-in", nil)
	w := httptest.NewRecorder()

	handler.signIn(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}

	expectedResponse := "{\"message\":\"Method not allowed\"}"
	if w.Body.String() != expectedResponse {
		t.Errorf("Expected response body %q, got %q", expectedResponse, w.Body.String())
	}
}
