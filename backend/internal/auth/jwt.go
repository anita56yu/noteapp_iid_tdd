package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT creates a new JWT for a given user ID.
func GenerateJWT(userID string, jwtSecret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token expires in 3 days
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
