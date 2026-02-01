package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
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
		log.Println("Error in auth.MakeJWT: Making JWT failed")
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
		log.Println("Error in auth.ValidateJWT: Parsing provided token string to create JWT failed")
		return uuid.Nil, err
	}

	if !token.Valid && claims.Issuer == Issuer {
		log.Println("Error in auth.ValidateJWT: Token not valid")
		return uuid.Nil, fmt.Errorf("Token not valid")
	}
	id, err := uuid.Parse(claims.Subject)
	log.Println("Error in auth.ValidateJWT: Parsing user id from JWT string failed")
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	tokenString, found := strings.CutPrefix(authHeader, "Bearer ")
	if !found {
		log.Println("Error in auth.GetBearerToken: correct authorization header not found")
		return "", fmt.Errorf("Token String not found")
	}
	return tokenString, nil
}

func MakeRefreshToken() (string, error) {
	randData := make([]byte, 32)
	rand.Read(randData)
	tokenString := hex.EncodeToString(randData)

	return tokenString, nil
}
