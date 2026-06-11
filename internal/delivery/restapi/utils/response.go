package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type BaseResponse struct {
	Success bool        `json:"success"`
	Body    interface{} `json:"body,omitempty"`
}

type BaseErrorResponse struct {
	Success bool        `json:"success"`
	Error   interface{} `json:"error,omitempty"`
}

type ErrorBody struct {
	Details string `json:"details,omitempty"`
}

type ValidationErrorResponse struct {
	Details string `json:"details,omitempty"`
	Field   string `json:"field,omitempty"`
}

func JSONResponse(w http.ResponseWriter, status int, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := BaseResponse{Success: status >= 200 && status < 400, Body: obj}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("Failed to encode error response: %v", err)
	}

}
