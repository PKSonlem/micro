package handler

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/timurzdev/mentorship-test-task/internal/entity"
)

// пишет ошибку на основе типа err
func ErrorResponse(w http.ResponseWriter, err error) {
	var validationErr validator.ValidationErrors

	if errors.Is(err, ErrBadRequest) ||
		errors.As(err, &validationErr) ||
		errors.Is(err, entity.ErrorCreatingHouse) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if errors.Is(err, ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if errors.Is(err, ErrForbidden) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if errors.Is(err, ErrUnauthorized) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
}

func SuccessResponse(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
