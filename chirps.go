package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/bitztec/chirpy/internal/database"
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

	dtoChirp := DTOChirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	respondWithJson(201, w, dtoChirp)
}

func (cfg *apiConfig) getAllChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(500, "", w)
		return
	}

	responseChirps := make([]DTOChirp, len(chirps))

	for i, chirp := range chirps {
		responseChirps[i] = DTOChirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}

	respondWithJson(200, w, responseChirps)
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
