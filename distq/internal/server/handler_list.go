package server

import (
	"net/http"

	"github.com/kentakom1213/distributed-system/distq/internal/storage"
)

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
