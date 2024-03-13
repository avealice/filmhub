package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	logrus.Error(message)

	errRes := errorResponse{Message: message}
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
