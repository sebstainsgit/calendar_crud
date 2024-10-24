package main

import (
	"fmt"
	"net/http"
	"time"
)

func (apiCfg *apiConfig) removeOldRefrTokens(w http.ResponseWriter, r *http.Request) {
	refrTokenArr, err := apiCfg.DB.GetAllRefrTokens(r.Context())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("unable to get all refresh tokens: %s", err))
		return
	}

	tokensDeleted := 0

	for _, token := range refrTokenArr {
		if token.Expires.After(time.Now().UTC()) {
			apiCfg.DB.DeleteRefrToken(r.Context(), token.RefrToken)
			tokensDeleted++
		}
	}

	type response struct {
		Tokens_Deleted int `json:"tokens_deleted"`
	}

	resp := response{Tokens_Deleted: tokensDeleted}

	respondWithJSON(w, http.StatusOK, resp)
}
