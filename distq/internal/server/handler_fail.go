package server

import "net/http"

type FailJobRequest struct {
	LeaseID string `json:"lease_id"`
	Error   string `json:"error"`
}

func (h *Handler) handleFailJob(w http.ResponseWriter, r *http.Request) {

}
