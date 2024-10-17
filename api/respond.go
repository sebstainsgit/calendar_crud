package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)

	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error marshalling payload to JSON: %v", err)
		return
	}

	w.WriteHeader(code)

	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(data))
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}

	respondWithJSON(w, code, errResponse{Error: msg})
}