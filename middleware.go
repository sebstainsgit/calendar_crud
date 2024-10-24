package main

import (
	"fmt"
	"net/http"
)

func (apiCfg *apiConfig) middlewareUserAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		JWT, err := getBearerHeader(r.Header)

		if err != nil {
			respondWithError(w, http.StatusForbidden, fmt.Sprintf("Missing or malformed header: %s", err))
			return
		}

		user, err := apiCfg.verifyJWT(JWT, r.Context())

		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error verifying JWT: %s", err))
			return
		}

		handler(w, r, user)
	}
}

func (apiCfg *apiConfig) middlewareAdminAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		JWT, err := getBearerHeader(r.Header)

		if err != nil {
			respondWithError(w, http.StatusForbidden, fmt.Sprintf("Missing or malformed header: %s", err))
			return
		}

		potentialAdmin, err := apiCfg.verifyJWT(JWT, r.Context())

		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error verifying JWT: %s", err))
			return
		}

		if potentialAdmin.Elevation != "admin" {
			respondWithError(w, http.StatusForbidden, "only admins are authorised to access this endpoint")
			return
		}

		handler(w, r)
	}
}

func (apiCfg *apiConfig) middlewareAdminAuthWithUser(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		JWT, err := getBearerHeader(r.Header)

		if err != nil {
			respondWithError(w, http.StatusForbidden, fmt.Sprintf("Missing or malformed header: %s", err))
			return
		}

		potentialAdmin, err := apiCfg.verifyJWT(JWT, r.Context())

		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error verifying JWT: %s", err))
			return
		}

		if potentialAdmin.Elevation != "admin" {
			respondWithError(w, http.StatusForbidden, "only admins are authorised to access this endpoint")
			return
		}

		handler(w, r, potentialAdmin)
	}
}