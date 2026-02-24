package utils

import (
	"encoding/json"
	"net/http"

	"github.com/ilyas/flower/services/auth/internal/dto"
)

func Send(w http.ResponseWriter, statusCodeSucces int, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeSucces)
	_ = json.NewEncoder(w).Encode(dto.Response{
		Data:    data,
		Message: message,
	})
}
