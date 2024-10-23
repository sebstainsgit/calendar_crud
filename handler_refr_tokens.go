package main

import (
	"fmt"
	"net/http"
)

func (apiCfg *apiConfig) makeJWTfromRefrToken(w http.ResponseWriter, r *http.Request) {
	userRefrToken, err := getBearerHeader(r.Header)

	if err != nil {
		respondWithError(w, http.StatusForbidden, fmt.Sprintf("Missing or malformed header: %s", err))
		return
	}

	users_id, err := apiCfg.DB.UserIDFromRefrToken(r.Context(), userRefrToken)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error getting refresh token in DB: %s", err))
		return
	}

	newJWT, err := apiCfg.createJWT(users_id)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating new JWT from refresh token: %s", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, JWTResponse{JWT: newJWT})
}
