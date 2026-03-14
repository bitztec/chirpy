package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"github.com/bitztec/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	serverMux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: serverMux,
	}

	handler := http.FileServer(http.Dir("./"))
	cfg := apiConfig{
		fileServerHits: atomic.Int32{},
		dbQueries:      database.New(db),
		platform:       os.Getenv("PLATFORM"),
	}

	serverMux.Handle(
		"/app/",
		http.StripPrefix("/app", cfg.middlewareMetricsIncrement(handler)))
	serverMux.HandleFunc("GET /api/healthz", hearthBeatHandler)
	serverMux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	serverMux.HandleFunc("POST /admin/reset", cfg.resetHandler)

	serverMux.HandleFunc("POST /api/users", cfg.createUserHandler)
	serverMux.HandleFunc("POST /api/login", cfg.logInHandler)

	serverMux.HandleFunc("POST /api/chirps", cfg.createChirpHandler)
	serverMux.HandleFunc("GET /api/chirps", cfg.getAllChirpsHandler)
	serverMux.HandleFunc("GET /api/chirps/{chirpID}", cfg.getChirpHandler)
	log.Fatal(server.ListenAndServe())
}

func hearthBeatHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write([]byte("OK"))
}

func validateChirpHandler(writer http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type validationResult struct {
		Error   string `json:"error"`
		Body    string `json:"cleaned_body"`
		isValid bool
	}

	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Printf("\nError decoding parameters: %s", err)
		writer.WriteHeader(500)
		return
	}

	respBody := validationResult{}

	if len(params.Body) > 140 {
		respBody.Error = "Chirp is too long"
	} else {
		respBody.Body = CleanResponse(params.Body)
		respBody.isValid = true
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		fmt.Printf("\nError marshalling JSON: %s", err)
		writer.WriteHeader(500)
		return
	}

	if respBody.isValid {
		writer.WriteHeader(200)
	} else {
		writer.WriteHeader(400)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(dat)
}

func CleanResponse(body string) string {
	parts := strings.Split(body, " ")
	newParts := make([]string, len(parts))

	for i, word := range parts {
		lower := strings.ToLower(word)
		if lower == "kerfuffle" ||
			lower == "sharbert" ||
			lower == "fornax" {
			word = "****"
		}

		newParts[i] = word
	}

	return strings.Join(newParts, " ")
}
