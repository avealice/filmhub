package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userRoleCtx         = "role"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

// userIdentity проверяет наличие и валидность токена аутентификации в заголовке запроса.
// Если токен корректен, устанавливает роль пользователя в контекст запроса.
// @Summary Проверка аутентификации пользователя
// @Description Middleware для проверки аутентификации пользователя и установки его роли в контекст запроса
// @Tags Authentication
// @Security BearerAuth
func (h *Handler) userIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(authorizationHeader)
		if header == "" {
			http.Error(w, "empty auth header", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "invalid auth header", http.StatusUnauthorized)
			return
		}

		if len(headerParts[1]) == 0 {
			http.Error(w, "token is empty", http.StatusUnauthorized)
			return
		}

		role, err := h.services.Authorization.ParseToken(headerParts[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userRoleCtx, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// getUserRole извлекает роль пользователя из контекста запроса.
// Если роль отсутствует или имеет неверный тип, возвращает ошибку.
// @Summary Извлечение роли пользователя
// @Description Извлекает роль пользователя из контекста запроса
// @Tags Authentication
func getUserRole(r *http.Request) (string, error) {
	role := r.Context().Value(userRoleCtx)
	if role == nil {
		return "", errors.New("user role not found")
	}

	roleStr, ok := role.(string)
	if !ok {
		return "", errors.New("user role is invalid type")
	}

	return roleStr, nil
}
