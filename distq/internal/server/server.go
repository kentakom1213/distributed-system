package server

import "net/http"

func NewHandler() http.Handler {
	mux := http.NewServeMux()
	
	mux.HandleFunc("GET /help", handleHealth)
	
	return mux
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok\n"))
}
