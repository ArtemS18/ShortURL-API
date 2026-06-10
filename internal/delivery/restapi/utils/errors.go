package utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/ArtemS18/ShortURL-API/internal/entity"
)

func WriteError(w http.ResponseWriter, msg string, status int) {

	w.WriteHeader(status)

	resp := ErrorResponse{Details: msg}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode error response: %v", err)
	}
}

func HandelError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, entity.InvalidInput):

		var validateErr *entity.ValidationError
		if errors.As(err, &validateErr) {
			JSONResponse(w,
				http.StatusBadRequest,
				ValidationErrorResponse{
					Details: "validation error",
					Field:   validateErr.Field,
				})
			return
		}
		WriteError(w, "bad request", http.StatusBadRequest)
	case errors.Is(err, entity.NotFoundError):
		WriteError(w, "slug not found", http.StatusNotFound)
	case errors.Is(err, entity.AlredyExitError):
		WriteError(w, "slug already exists", http.StatusConflict)
	default:
		WriteError(w, "internal", http.StatusInternalServerError)
	}
}
