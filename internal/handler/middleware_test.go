package handler

import (
	"filmhub/internal/service"
	mock_service "filmhub/internal/service/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestHandler_userIdentity_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_service.NewMockAuthorization(ctrl)
	mockAuthService.EXPECT().ParseToken(gomock.Any()).Return("user", nil)

	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuthService,
		},
	}

	identityHandler := handler.userIdentity(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if role, ok := r.Context().Value(userRoleCtx).(string); !ok || role != "user" {
			t.Error("Expected user role in context")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set(authorizationHeader, "Bearer valid_token")
	w := httptest.NewRecorder()

	identityHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHandler_userIdentity_AuthHeaderMissing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_service.NewMockAuthorization(ctrl)

	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuthService,
		},
	}

	identityHandler := handler.userIdentity(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called if auth header is missing")
		w.WriteHeader(http.StatusInternalServerError)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	identityHandler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestHandler_userIdentity_InvalidAuthHeader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_service.NewMockAuthorization(ctrl)

	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuthService,
		},
	}

	identityHandler := handler.userIdentity(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called if auth header is invalid")
		w.WriteHeader(http.StatusInternalServerError)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set(authorizationHeader, "invalid_token")
	w := httptest.NewRecorder()

	identityHandler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestHandler_userIdentity_AdminRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_service.NewMockAuthorization(ctrl)
	mockAuthService.EXPECT().ParseToken(gomock.Any()).Return("admin", nil)

	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuthService,
		},
	}

	identityHandler := handler.userIdentity(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if role, ok := r.Context().Value(userRoleCtx).(string); !ok || role != "admin" {
			t.Error("Expected admin role in context")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set(authorizationHeader, "Bearer admin_token")
	w := httptest.NewRecorder()

	identityHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}
