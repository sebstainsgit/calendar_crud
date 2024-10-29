package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func makeVarChar() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func getBearerHeader(headers http.Header) (string, error) {
	val := headers.Get("Authorisation")

	if val == "" {
		return "", errors.New("no authentication info found")
	}

	//Expects [Bearer, {JWT}]
	vals := strings.Split(val, " ")

	if len(vals) != 2 {
		return "", errors.New("malformed auth header")
	}

	if vals[0] != "Bearer" {
		return "", errors.New("malformed first part of auth header")
	}

	return vals[1], nil
}

func decodeParams[T any](body io.ReadCloser) (T, error) {
	var parameters T

	data, err := io.ReadAll(body)

	if err != nil {
		return *new(T), err
	}

	err = json.Unmarshal(data, &parameters)

	if err != nil {
		return *new(T), err
	}

	return parameters, nil
}

func contains[T string | int | uuid.UUID](target T, arr []T) bool {
	for _, v := range arr {
		if target == v {
			return true
		}
	}
	return false
}

func (apiCfg *apiConfig) resGetUsersFromArr(ctx context.Context, concerns []uuid.UUID) ([]User, error) {
	users := []User{}

	for _, userID := range concerns {
		user, err := apiCfg.DB.GetUserByID(ctx, userID)

		if err != nil {
			return []User{}, err
		}

		users = append(users, DBUserToLocalUser(user))
	}

	return users, nil
}
