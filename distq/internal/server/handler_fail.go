package server

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/kentakom1213/distributed-system/distq/internal/storage"
)

type failJobRequest struct {
	LeaseID  string `json:"lease_id"`
	ErrorMsg string `json:"error"`
}

func (h *Handler) handleFailJob(w http.ResponseWriter, r *http.Request) {
	jobID := r.PathValue("id")
	if jobID == "" {
		writeJSONError(w, http.StatusBadRequest, "job id is required")
		return
	}

	var req failJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if req.LeaseID == "" {
		writeJSONError(w, http.StatusBadRequest, "lease_id is required")
		return
	}

	if req.ErrorMsg == "" {
		req.ErrorMsg = "job failed"
	}

	job, err := h.storage.FailJob(r.Context(), jobID, req.LeaseID, req.ErrorMsg)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrLeaseExpired):
			slog.Warn("fail lease expired",
				"job_id", jobID,
				"lease_id", req.LeaseID,
			)
			writeJSONError(w, http.StatusConflict, "lease expired")
			return

		case errors.Is(err, storage.ErrInvalidLease):
			slog.Warn("invalid fail lease",
				"job_id", jobID,
				"lease_id", req.LeaseID,
			)
			writeJSONError(w, http.StatusConflict, "invalid lease")
			return

		case errors.Is(err, storage.ErrJobNotFound):
			writeJSONError(w, http.StatusNotFound, "job not found")
			return

		case errors.Is(err, storage.ErrJobNotRunning):
			writeJSONError(w, http.StatusConflict, "job is not running")
			return

		default:
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	slog.Info("job failed",
		"job_id", jobID,
		"lease_id", req.LeaseID,
		"status", job.Status,
		"attempts", job.Attempts,
		"max_attempts", job.MaxAttempts,
		"error", req.ErrorMsg,
	)

	writeJSON(w, http.StatusOK, job)
}
