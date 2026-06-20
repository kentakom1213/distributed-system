package server

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/kentakom1213/distributed-system/distq/internal/storage"
)

type ackJobRequest struct {
	LeaseID string `json:"lease_id"`
}

func (h *Handler) handleAckJob(w http.ResponseWriter, r *http.Request) {
	jobID := r.PathValue("id")
	if jobID == "" {
		writeJSONError(w, http.StatusBadRequest, "job id is required")
		return
	}

	var req ackJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if req.LeaseID == "" {
		writeJSONError(w, http.StatusBadRequest, "lease_id is required")
		return
	}

	if err := h.storage.AckJob(r.Context(), jobID, req.LeaseID); err != nil {
		switch {
		case errors.Is(err, storage.ErrLeaseExpired):
			slog.Warn("ack lease expired",
				"job_id", jobID,
				"lease_id", req.LeaseID,
			)
			writeJSONError(w, http.StatusConflict, "lease expired")
			return

		case errors.Is(err, storage.ErrInvalidLease):
			slog.Warn("invalid ack lease",
				"job_id", jobID,
				"lease_id", req.LeaseID,
			)
			writeJSONError(w, http.StatusConflict, "invalid lease")
			return

		case errors.Is(err, storage.ErrJobNotFound):
			writeJSONError(w, http.StatusConflict, "job not found")
			return

		case errors.Is(err, storage.ErrJobNotRunning):
			writeJSONError(w, http.StatusConflict, "job not running")
			return

		default:
			writeJSONError(w, http.StatusConflict, err.Error())
			return
		}
	}

	slog.Info("job acknowledgement",
		"job_id", jobID,
		"lease_id", req.LeaseID,
	)

	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
