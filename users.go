package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	auth "github.com/bitztec/chirpy/internal/auth"
	db "github.com/bitztec/chirpy/internal/database"
)

type UserPwdParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	params, err := UnmarshallUserParams(r)
	if err != nil {
		respondWithError(
			500,
			fmt.Sprintf("Error decoding parameters: %s", err),
			w)
		return
	}

	// Hash password when creating user
	hashedPwd, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("\nError creating user: %s", err)
		w.WriteHeader(400)
		return
	}

	// Create user
	userParams := db.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPwd,
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), userParams)
	if err != nil {
		log.Printf("\nError creating user: %s", err)
		w.WriteHeader(500)
		return
	}

	responseBody, err := json.Marshal(user.ToDTO())
	if err != nil {
		log.Printf("\nError marshalling JSON: %s", err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(responseBody)
}

func (cfg *apiConfig) logInHandler(w http.ResponseWriter, r *http.Request) {
	// Unmarshall request parameters
	params, err := UnmarshallUserParams(r)
	if err != nil {
		respondWithError(
			500,
			fmt.Sprintf("Unable to get parameters, %s", err),
			w)
		return
	}

	dbUser, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(401, "Incorrect email or password", w)
		return
	}

	validPwd, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil || !validPwd {
		respondWithError(401, "Incorrect email or password", w)
		return
	}

	respondWithJson(200, w, dbUser.ToDTO())
}

func UnmarshallUserParams(r *http.Request) (UserPwdParams, error) {
	decoder := json.NewDecoder(r.Body)
	params := UserPwdParams{}
	err := decoder.Decode(&params)
	return params, err
}
