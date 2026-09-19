package auth

import (
	"time"
	"net/http"
	"errors"
	"strings"
	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeAccess TokenType = "chirpy-access"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}

func MakeJWT(
	userID uuid.UUID,
	tokenSecret string,
	expiresIn time.Duration,
) (string, error) {
	signingKey := []byte(tokenSecret)
	claims := &jwt.RegisteredClaims{
		Issuer: "chirpy-access",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject: userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claimsStruct, func(token *jwt.Token) (any, error) {return []byte(tokenSecret), nil })
	if err != nil {
		return uuid.Nil, err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, nil
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.UUID{}, nil
	}
	if issuer != string(TokenTypeAccess) {
		return uuid.UUID{}, nil
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.UUID{}, nil
	}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	auth:= headers.Get("Authorization")
	if auth == "" {
		return "", errors.New("no authorization header found")
	}
	trimmedAuth := strings.Split(auth, " ")
	if len(trimmedAuth) < 2 {
		return "", errors.New("invalid auth header")
	}

	return trimmedAuth[1], nil
}
