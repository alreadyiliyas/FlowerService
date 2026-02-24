package handlers

import (
	"net/http"

	"github.com/ilyas/flower/services/auth/internal/utils"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
	utils.Send(w, http.StatusOK, nil, "ok")
}

func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	utils.Send(w, http.StatusAlreadyReported, nil, "ready")
}
