package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))

func GenerateJWT(subject string) (string, error) {
	tokenDuration := 72 * time.Hour
	claims := jwt.MapClaims{
		"sub": subject,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(tokenDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
