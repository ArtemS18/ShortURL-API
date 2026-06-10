package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type BaseResponse struct {
	OK   bool        `json:"ok"`
	Body interface{} `json:"body,omitempty"`
}

type ErrorResponse struct {
	Details string `json:"details,omitempty"`
}
type ValidationErrorResponse struct {
	Details string `json:"details,omitempty"`
	Field   string `json:"field"`
}

func JSONResponse(w http.ResponseWriter, status int, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := BaseResponse{OK: status >= 200 && status < 300, Body: obj}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("Failed to encode error response: %v", err)
	}

}
