package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ArtemS18/ShortURL-API/internal/entity"
)

func WriteError(w http.ResponseWriter, msg string, status int) {

	w.WriteHeader(status)

	resp := BaseErrorResponse{
		Success: false,
		Error:   ErrorBody{Details: msg},
	}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		fmt.Printf("Failed to encode error response: %v", err)
	}
}

func HandelError(w http.ResponseWriter, err error) {
	var notFoundErr *entity.NotFoundError
	var alredyExitErr *entity.AlredyExitError
	switch {
	case errors.Is(err, entity.InvalidInput):

		var validateErr *entity.ValidationError
		if errors.As(err, &validateErr) {
			JSONResponse(w,
				http.StatusBadRequest,
				ValidationErrorResponse{
					Details: validateErr.Details,
					Field:   validateErr.Field,
				})
			return
		}
		WriteError(w, "bad request", http.StatusBadRequest)
	case errors.As(err, &notFoundErr):
		WriteError(w, fmt.Sprintf("%s not found", notFoundErr.Field), http.StatusNotFound)
	case errors.As(err, &alredyExitErr):
		WriteError(w, fmt.Sprintf("%s already exists", alredyExitErr.Field), http.StatusConflict)
	default:
		WriteError(w, "internal", http.StatusInternalServerError)
	}
}
