package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	// Basic health check function
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(http.StatusText(200)))
} // End HandlerReadiness() func

func (cfg *apiConfig) HandlerMetrics(w http.ResponseWriter, r *http.Request) {
	// Display metrics to the admin user
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileServerHits.Load())))
} // End handlerMetrics() func

func (cfg *apiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
} // End MiddlewareMetricsInc() func