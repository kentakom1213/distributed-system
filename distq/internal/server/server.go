package server

import (
	"net/http"

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
	
	return mux
}

func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok\n"))
}
