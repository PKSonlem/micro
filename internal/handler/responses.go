package handler

import (
	"errors"
	"net/http"
)

// пишет ошибку на основе типа err
func ErrorResponse(w http.ResponseWriter, err error) {
	if errors.Is(err, ErrBadRequest) {
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
	return
}

func SuccessResponse(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
	w.WriteHeader(http.StatusOK)
}
