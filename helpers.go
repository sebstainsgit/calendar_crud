package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
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
