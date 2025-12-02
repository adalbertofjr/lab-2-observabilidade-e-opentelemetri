package api

import (
	"encoding/json"
	"net/http"
)

type HealthCheck struct {
	Status string `json:"status"`
}

func NewHealthCheck() *HealthCheck {
	return &HealthCheck{}
}

func (h *HealthCheck) HealthCheck(w http.ResponseWriter, r *http.Request) {
	healthCheck := HealthCheck{Status: "ok"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(healthCheck)
}
