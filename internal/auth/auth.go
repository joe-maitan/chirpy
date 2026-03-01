package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	// TokenTypeAccess -
	TokenTypeAccess TokenType = "chirpy-access"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Printf("auth.go - HashPassword() - Error hashing password: %v", err)
		return "", err
	}

	return hashedPassword, nil
} // End HashPassword() func

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		log.Printf("auth.go - CheckPasswordHash() - Error checking password hash: %v", err)
		return false, err
	}

	return match, nil
} // End CheckPasswordHash() func

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: string(TokenTypeAccess),
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject: userID.String(),
	})

	return newToken.SignedString([]byte(tokenSecret))
} // End MakeJWT() func

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)

	if err != nil {
		log.Printf("auth.go - ValidateJWT() - Error calling jwt.ParseWithClaims: %v", err)
		return uuid.Nil, err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		log.Printf("auth.go - ValidateJWT() - Error getting token subject: %v", err)
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		log.Printf("auth.go - ValidateJWT() - Error getting token issuer: %v", err)
		return uuid.Nil, err
	}

	if issuer != string(TokenTypeAccess) {
		log.Printf("auth.go - ValidateJWT() - Error checking token issuer: %v", err)
		return uuid.Nil, errors.New("invalid issuer")
	}

	id, err := uuid.Parse(userIDString)
	if err != nil {
		log.Printf("auth.go - ValidateJWT() - Error parsing userID: %v", err)
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}

	return id, nil
} // End ValidateJWT() func

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		log.Printf("auth.go - GetBearerToken() - No Authorization header included in request")
		return "", errors.New("no auth header included in request")
	}

	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		log.Printf("auth.go - GetBearerToken() - Malformed Authorization header: %s", authHeader)
		return "", errors.New("malformed authorization header")
	}

	log.Printf("auth.go - GetBearerToken() - Successfully retrieved bearer token from header: %s", splitAuth[1])
	return splitAuth[1], nil
}// End GetBearerToken() func

func MakeRefreshToken() (string, error) {
	byteArr := make([]byte, 32)
	rand.Read(byteArr)
	return hex.EncodeToString(byteArr), nil
} // End MakeRefreshToken() func

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		log.Printf("auth.go - GetAPIKey() - No Authorization header included in request")
		return "", errors.New("no auth header included in request")
	}

	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		log.Printf("auth.go - GetAPIKey() - Malformed Authorization header: %s", authHeader)
		return "", errors.New("malformed authorization header")
	}

	log.Printf("auth.go - GetAPIKey() - Successfully retrieved Apikey token from header: %s", splitAuth[1])
	return splitAuth[1], nil
} // End GetAPIKey() func
