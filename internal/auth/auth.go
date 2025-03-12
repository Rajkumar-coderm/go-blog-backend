package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Load JWT secret from environment variables
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// Claims struct for JWT payload
type Claims struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a new JWT token with expiration
func GenerateJWT(id string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		ID:   id,
		Type: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   id,
			Issuer:    "go-blog-backend",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// GenerateRefreshToken creates a refresh token (valid for 7 days)
func GenerateRefreshToken(id string) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	claims := &Claims{
		ID:   id,
		Type: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   id,
			Issuer:    "go-blog-backend",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateJWT verifies the token and extracts claims
func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("authentication token is invalid")
	}

	return claims, nil
}
