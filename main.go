package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func main() {
	serverMux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: serverMux,
	}

	handler := http.FileServer(http.Dir("./"))
	cfg := apiConfig{fileServerHits: atomic.Int32{}}

	serverMux.Handle(
		"/app/",
		http.StripPrefix("/app", cfg.middlewareMetricsInc(handler)))
	serverMux.HandleFunc("GET /api/healthz", hearthBeatHandler)
	serverMux.HandleFunc("GET /api/metrics", cfg.fileServerHitsHandler)
	serverMux.HandleFunc("POST /api/reset", cfg.resetHitsHandler)
	log.Fatal(server.ListenAndServe())
}

func hearthBeatHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write([]byte("OK"))
}

func (cfg *apiConfig) fileServerHitsHandler(writer http.ResponseWriter,
	request *http.Request,
) {
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write(fmt.Append([]byte("Hits: "), cfg.fileServerHits.Load()))
}

func (cfg *apiConfig) resetHitsHandler(writer http.ResponseWriter, request *http.Request) {
	cfg.fileServerHits.Store(0)
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write([]byte("Hits reseted to 0"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
