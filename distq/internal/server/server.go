package server

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/kentakom1213/distributed-system/distq/internal/storage"
)

type Handler struct {
	storage *storage.Storage
}

func NewHandler(st *storage.Storage) http.Handler {
	h := &Handler{
		storage: st,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /help", h.handleHealth)
	mux.HandleFunc("POST /jobs", h.handleCreateJob)
	mux.HandleFunc("GET /jobs", h.handleListJobs)
	mux.HandleFunc("POST /jobs/claim", h.handlePickClaimJob)
	mux.HandleFunc("POST /jobs/{id}/ack", h.handleAckJob)

	return loggingMiddleware(mux)
}

type createJobRequest struct {
	Type        string          `json:"type"`
	Payload     json.RawMessage `json:"payload"`
	MaxAttempts int             `json:"max_attempts"`
}

func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok\n"))
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

func (h *Handler) handleListJobs(w http.ResponseWriter, r *http.Request) {
	status := storage.JobStatus(r.URL.Query().Get("status"))

	jobs, err := h.storage.ListJobs(r.Context(), storage.ListJobsParams{
		Status: status,
		Limit:  100,
	})
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, jobs)
}

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
