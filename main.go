package main

import (
	"log"
	"net/http"
)

func main() {
	serverMux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: serverMux,
	}

	serverMux.Handle(
		"/app/",
		http.StripPrefix("/app", http.FileServer(http.Dir("./"))))
	serverMux.HandleFunc("/healthz", hearthBeatHandler)
	log.Fatal(server.ListenAndServe())
}

func hearthBeatHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write([]byte("OK"))
}
