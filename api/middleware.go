package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) middlewareUserAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := GetApiKey(r.Header)

		if err != nil {
			respondWithError(w, http.StatusForbidden, fmt.Sprintf("Auth error: %s", err))
			return
		}

		user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)

		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn't get user: %s", err))
			return
		}

		handler(w, r, user)
	}
}
