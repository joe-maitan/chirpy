package auth

import (
	"log"
	"time"

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

	jwt, err := newToken.SignedString(tokenSecret)
	if err != nil {
		return "", err
	}

	return jwt, nil
} // End MakeJWT() func

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	// Correct signature: func(*jwt.Token) (interface{}, error)
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if err, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, err
		}
		return []byte("your-secret-key"), nil
	}
	
	token, err := jwt.ParseWithClaims(tokenString, jwt.RegisteredClaims{}, keyFunc)
} // End ValidateJWT() func