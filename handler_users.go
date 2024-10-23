package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sebstainsgit/calendar/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (apiCfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	//Turn into function
	type params struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var parameters params

	data, err := io.ReadAll(r.Body)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error reading request body: %s", err))
		return
	}

	err = json.Unmarshal(data, &parameters)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error unmarhsalling request body: %s", err))
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(parameters.Password), bcrypt.DefaultCost)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error hashing password")
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		UserID:    uuid.New(),
		Name:      parameters.Name,
		Email:     parameters.Email,
		Password:  string(hashed),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error writing user to DB: %s", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, DBUserToLocalUser(user))
}

func (apiCfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var parameters params

	data, err := io.ReadAll(r.Body)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error reading request body: %s", err))
		return
	}

	err = json.Unmarshal(data, &parameters)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error unmarhsalling request body: %s", err))
		return
	}

	user, err := apiCfg.DB.GetUserFromEmail(r.Context(), parameters.Email)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error getting user from DB: %s", err))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(parameters.Password))

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid password, please try again")
		return
	}
	//After this point we assume the user is authed (has correct email)
	//Expires in 3 hours
	token, err := apiCfg.createJWT(user.UserID, 3)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating JWT: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, JWTResponse{JWT: token})
}

func (apiCfg *apiConfig) getAllUsers(w http.ResponseWriter, r *http.Request) {
	userArr, err := apiCfg.DB.GetAllUsers(r.Context())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error getting users from DB")
	}

	respondWithJSON(w, http.StatusOK, DBUsersToLocalUsers(userArr))
}

func (apiCfg *apiConfig) deleteUserSelf(w http.ResponseWriter, r *http.Request, user database.User) {
	type params struct {
		User_ID string `json:"user_id"`
	}

	var parameters params

	data, err := io.ReadAll(r.Body)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error reading request body: %s", err))
		return
	}

	err = json.Unmarshal(data, &parameters)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error unmarhsalling request body: %s", err))
		return
	}
	//IF authenticated user is the user to be deleted, else ignore
	if user.UserID.String() != parameters.User_ID {
		respondWithError(w, http.StatusUnauthorized, "only user can delete their own account")
		return
	}

	err = apiCfg.DB.DeleteUser(r.Context(), user.UserID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error deleting user from database")
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")
}