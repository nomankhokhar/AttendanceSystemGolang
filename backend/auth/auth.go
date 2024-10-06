package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	secretKey = []byte("f98g7d8f6g8df87g6d87f6g87df78gdjdfgjdfjkgnkdjfngjndjf") 
)

// Claims struct to hold the JWT claims
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a new JWT token for a given email
func GenerateJWT(email string) (string, error) {
	claims := Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Token valid for 24 hours
		},
	}

	// Create the token using the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ValidateToken validates the JWT token
func ValidateToken(tokenString string) (Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return Claims{}, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return *claims, nil
	}

	return Claims{}, errors.New("invalid token")
}
