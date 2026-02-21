package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) metricsHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(200)
	message := fmt.Sprintf(`<html>
  		<body>
    		<h1>Welcome, Chirpy Admin</h1>
    		<p>Chirpy has been visited %d times!</p>
  		</body>
	</html>`, cfg.fileServerHits.Load())

	writer.Write([]byte(message))
}

func (cfg *apiConfig) resetHandler(writer http.ResponseWriter, request *http.Request) {
	cfg.fileServerHits.Store(0)
	writer.Header().Add("Content-Type", "test/plain; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write([]byte("Hits reseted to 0"))
}

func (cfg *apiConfig) middlewareMetricsIncrement(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
