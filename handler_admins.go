package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sebstainsgit/calendar/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (apiCfg *apiConfig) createAdmin(w http.ResponseWriter, r *http.Request) {
	//Turn into function
	type params struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	parameters, err := decodeParams[params](r.Body)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error decoding request body: %s", err))
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(parameters.Password), bcrypt.DefaultCost)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error hashing password: %s", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		UserID:    uuid.New(),
		Name:      parameters.Name,
		Email:     parameters.Email,
		Password:  string(hashed),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Elevation: "admin",
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error writing user to DB: %s", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, DBUserToLocalUser(user))
}

func (apiCfg *apiConfig) deleteUser(w http.ResponseWriter, r *http.Request, admin database.User) {
	type params struct {
		UserID uuid.UUID `json:"user_id"`
	}

	parameters, err := decodeParams[params](r.Body)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error decoding request body: %s", err))
		return
	}

	err = apiCfg.DB.DeleteUser(r.Context(), parameters.UserID)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error deleting user from database")
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")
}
