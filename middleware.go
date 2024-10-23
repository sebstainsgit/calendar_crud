package main

import (
	"fmt"
	"net/http"
)

func (apiCfg *apiConfig) middlewareUserAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		JWT, err := getJWTHeader(r.Header)

		if err != nil {
			respondWithError(w, http.StatusForbidden, fmt.Sprintf("Missing or malformed header: %s", err))
			return
		}

		user, err := apiCfg.verifyJWT(JWT, r.Context())

		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error verifying JWT: %s", err))
		}
		handler(w, r, user)
	}
}
