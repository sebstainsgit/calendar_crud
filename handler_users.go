package main

import (
	"fmt"
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
	//Automatically adds Elevation: "user" in sqlc function
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

	parameters, err := decodeParams[params](r.Body)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error decoding request body: %s", err))
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
	//After this point we assume the user is authed (has correct email and password)
	//Expires in 3 hours
	token, err := apiCfg.createJWT(user.UserID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating JWT: %s", err))
		return
	}

	refrToken, err := apiCfg.createRefrToken(r.Context(), user.UserID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating and saving refresh token: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, TokenResponse{JWT: token, RefreshToken: refrToken.RefrToken})
}

func (apiCfg *apiConfig) updateUserInfo(w http.ResponseWriter, r *http.Request, user database.User) {
	type params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	parameters, err := decodeParams[params](r.Body)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error decoding request body: %s", err))
		return
	}

	updatedUser, err := apiCfg.DB.UpdateUser(r.Context(), database.UpdateUserParams{
		UserID:    user.UserID,
		Name:      parameters.Name,
		Email:     parameters.Email,
		Password:  parameters.Password,
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error updating user info in DB: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, DBUserToLocalUser(updatedUser))
}

func (apiCfg *apiConfig) getAllUsers(w http.ResponseWriter, r *http.Request) {
	userArr, err := apiCfg.DB.GetAllUsers(r.Context())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting users from DB: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, DBUsersToLocalUsers(userArr))
}

func (apiCfg *apiConfig) deleteUserSelf(w http.ResponseWriter, r *http.Request, user database.User) {
	//Deletes logged in user
	err := apiCfg.DB.DeleteUser(r.Context(), user.UserID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error deleting user from database")
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")
}