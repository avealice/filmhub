package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

// ErrorResponse представляет JSON-структуру ответа с сообщением об ошибке.
//
// @title ErrorResponse
// @description JSON-структура ответа с сообщением об ошибке.
type ErrorResponse struct {
	Message string `json:"message"`
}

// newErrorResponse создает новый JSON-ответ с сообщением об ошибке и отправляет его клиенту.
//
// @Summary Создать ответ с сообщением об ошибке.
// @Description Создает новый JSON-ответ с заданным статусом кода и сообщением об ошибке, затем отправляет его клиенту.
// @Tags Error Handling
func newErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	logrus.Error(message)

	errRes := ErrorResponse{Message: message}
	jsonResponse, err := json.Marshal(errRes)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(jsonResponse)
	if err != nil {
		logrus.Error("Failed to write response:", err)
	}
}
