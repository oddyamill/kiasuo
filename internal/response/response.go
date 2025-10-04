package response

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func response(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("JSON encoding error", "error", err)
	}
}

func Ok[T any](w http.ResponseWriter, body T) {
	response(w, http.StatusOK, body)
}

func Created[T any](w http.ResponseWriter, body T) {
	response(w, http.StatusCreated, body)
}

type errorResponse struct {
	Message string `json:"message"`
}

func errResponse(w http.ResponseWriter, status int, message string) {
	response(w, status, errorResponse{message})
}

func ErrBadRequest(w http.ResponseWriter, message string) {
	errResponse(w, http.StatusBadRequest, message)
}

func ErrInternalServer(w http.ResponseWriter, err error) {
	errResponse(w, http.StatusInternalServerError, err.Error())
}

func ErrNotFound(w http.ResponseWriter) {
	errResponse(w, http.StatusNotFound, "not found")
}

func ErrUnauthorized(w http.ResponseWriter) {
	errResponse(w, http.StatusUnauthorized, "unauthorized")
}
