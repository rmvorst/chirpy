package auth

import (
	"fmt"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const Issuer = "chirpy"

func HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    Issuer,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	}).SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid && claims.Issuer == Issuer {
		return uuid.Nil, fmt.Errorf("Token not valid")
	}
	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
