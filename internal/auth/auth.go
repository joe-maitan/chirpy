package auth

import (
	"fmt"
	"log"
	"time"
	"net/http"
	"errors"

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return hashedPassword, nil
} // End HashPassword() func

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}

	return match, nil
} // End CheckPasswordHash() func

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject: userID.String(),
	})

	jwt, err := newToken.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return jwt, nil
} // End MakeJWT() func

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	}
	
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, keyFunc)
	if err != nil {
		return uuid.Nil, err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != "chirpy" {
		return uuid.Nil, errors.New("invalid issuer")
	}

	id, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}

	return id, nil
} // End ValidateJWT() func

func GetBearerToken(headers http.Header) (string, error) {
	// The outer loop iterates over the map
	for key, values := range headers {
		// The inner loop iterates over the slice of values for each header name
		for _, value := range values {
			if value.Contains("Bearer") {
				data := strings.Split(value, " ")
				tokenString := data[1]

				if tokenString == "" || tokenString == nil {
					return "", err
				} else {
					return tokenString, nil
				}
			}
		}
	}

	return "", err
} // End GetBearerToken() func