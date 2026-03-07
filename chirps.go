package main

import (
	"encoding/json"
	"net/http"
	"strings"

	database "github.com/bitztec/chirpy/internal/database"
	dto "github.com/bitztec/chirpy/internal/dataTransfer"
	"github.com/google/uuid"
)

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(500, "", w)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(400, "Chirp is too long", w)
		return
	}

	cleanChirp := cleanResponse(params.Body)
	dbParams := database.CreateChirpParams{
		Body:   cleanChirp,
		UserID: params.UserID,
	}

	dbChirp, err := cfg.dbQueries.CreateChirp(r.Context(), dbParams)
	if err != nil {
		respondWithError(500, "", w)
		return
	}

	respondWithJson(201, w, dbChirp.ToDTO())
}

func (cfg *apiConfig) getAllChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(500, "", w)
		return
	}

	responseChirps := make([]dto.DTOChirp, len(chirps))

	for i, chirp := range chirps {
		responseChirps[i] = chirp.ToDTO()
	}

	respondWithJson(200, w, responseChirps)
}

func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpID")
	if len(id) == 0 {
		respondWithError(404, "Chirp not found", w)
		return
	}

	dbID, err := uuid.Parse(id)
	if err != nil {
		respondWithError(404, "Chirp not found", w)
		return
	}

	dbChirp, err := cfg.dbQueries.GetChirpById(r.Context(), dbID)
	if err != nil {
		respondWithError(404, "Chirp not found", w)
		return
	}

	respondWithJson(200, w, dbChirp.ToDTO())
}

func cleanResponse(body string) string {
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
