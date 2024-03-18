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
	userIDCtx           = "user_id"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

// userIdentity проверяет наличие и валидность токена аутентификации в заголовке запроса.
// Если токен корректен, устанавливает роль и id пользователя в контекст запроса.
// @Summary Проверка аутентификации пользователя
// @Description Middleware для проверки аутентификации пользователя и установки его роли и id в контекст запроса
// @Tags Authentication
// @Security ApiKeyAuth
func (h *Handler) userIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(authorizationHeader)
		if header == "" {
			http.Error(w, "empty auth header", http.StatusUnauthorized)
			return
		}

		var token string
		headerParts := strings.Split(header, " ")
		if len(headerParts) > 2 || len(headerParts) == 0 {
			http.Error(w, "invalid auth header", http.StatusUnauthorized)
			return
		} else if len(headerParts) == 2 {
			token = headerParts[1]
		} else {
			token = headerParts[0]
		}

		user_id, role, err := h.services.Authorization.ParseToken(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userRoleCtx, role)
		ctx = context.WithValue(ctx, userIDCtx, user_id)
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

// getUserID извлекает идентификатор пользователя из контекста запроса.
// Если идентификатор отсутствует или имеет неверный тип, возвращает ошибку.
// @Summary Извлечение идентификатора пользователя
// @Description Извлекает идентификатор пользователя из контекста запроса
// @Tags Authentication
func getUserID(r *http.Request) (int, error) {
	userID := r.Context().Value(userIDCtx)
	if userID == nil {
		return 0, errors.New("user ID not found")
	}

	userIDInt, ok := userID.(int)
	if !ok {
		return 0, errors.New("user ID is invalid type")
	}

	return userIDInt, nil
}
