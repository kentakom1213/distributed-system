package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/kentakom1213/distributed-system/distq/internal/storage"
)

type createJobRequest struct {
	Type        string          `json:"type"`
	Payload     json.RawMessage `json:"payload"`
	MaxAttempts int             `json:"max_attempts"`
}

func (h *Handler) handleCreateJob(w http.ResponseWriter, r *http.Request) {
	var req createJobRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	job, err := h.storage.CreateJob(r.Context(), storage.CreateJobParams{
		Type:        req.Type,
		Payload:     req.Payload,
		MaxAttempts: req.MaxAttempts,
	})
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	slog.Info("job created",
		"job_id", job.ID,
		"type", job.Type,
		"max_attempts", job.MaxAttempts,
	)

	writeJSON(w, http.StatusCreated, job)
}
