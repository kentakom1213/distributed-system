package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/kentakom1213/distributed-system/distq/internal/storage"
)

type claimJobRequest struct {
	WorkerID     string `json:"worker_id"`
	LeaseSeconds int    `json:"lease_seconds"`
}

func (h *Handler) handlePickClaimJob(w http.ResponseWriter, r *http.Request) {
	var req claimJobRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if req.WorkerID == "" {
		writeJSONError(w, http.StatusBadRequest, "worker_id is required")
		return
	}

	leaseDuration := 30 * time.Second
	if req.LeaseSeconds > 0 {
		leaseDuration = time.Duration(req.LeaseSeconds) * time.Second
	}

	job, err := h.storage.PickNextJob(r.Context(), storage.PickNextJobParams{
		WorkerID:      req.WorkerID,
		LeaseDuration: leaseDuration,
	})
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if job == nil {
		slog.Info("no job available", "worker_id", req.WorkerID)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	slog.Info("job claimed",
		"job_id", job.ID,
		"type", job.Type,
		"worker_id", req.WorkerID,
		"lease_id", valueOrEmpty(job.LeaseID),
		"lease_until", job.LeaseUntil,
		"attempts", job.Attempts,
		"max_attempts", job.MaxAttempts,
	)

	writeJSON(w, http.StatusOK, job)
}
