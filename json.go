package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(code int, errorMsg string, w http.ResponseWriter) {
	log.Println(errorMsg)
	w.WriteHeader(code)

	if code >= 500 && code < 600 {
		err := "Error processing the request"
		w.Write(([]byte)(err))
		return
	}

	if code >= 400 && code < 500 {
		w.Write(([]byte)(errorMsg))
		return
	}
}

func respondWithJson(code int, w http.ResponseWriter, body interface{}) {
	data, err := json.Marshal(body)
	if err != nil {
		respondWithError(500, "", w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
