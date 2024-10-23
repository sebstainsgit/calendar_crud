package main

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sebstainsgit/calendar/internal/database"
)

func (apiCfg *apiConfig) createJWT(UserID uuid.UUID, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "calendar",
		Subject:   UserID.String(),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiresIn) * time.Hour)),
	}

	JWToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := JWToken.SignedString([]byte(apiCfg.JWT_SECRET))

	if err != nil {
		return "", err
	}

	return token, nil
}

func (apiCfg *apiConfig) verifyJWT(strToken string, ctx context.Context) (database.User, error) {
	token, err := jwt.ParseWithClaims(strToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(apiCfg.JWT_SECRET), nil
	})

	if err != nil {
		return database.User{}, err
	}

	strID, err := token.Claims.GetSubject()

	if err != nil {
		return database.User{}, err
	}

	userID, err := uuid.Parse(strID)

	if err != nil {
		return database.User{}, err
	}

	user, err := apiCfg.DB.GetUserFromID(ctx, userID)

	if err != nil {
		return database.User{}, nil
	}

	return user, nil
}
