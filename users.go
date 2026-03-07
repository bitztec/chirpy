package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Printf("\nError decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), params.Email)
	if err != nil {
		fmt.Printf("\nError creating user: %s", err)
		w.WriteHeader(500)
		return
	}

	responseBody, err := json.Marshal(user.ToDTO())
	if err != nil {
		fmt.Printf("\nError marshalling JSON: %s", err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(responseBody)
}
