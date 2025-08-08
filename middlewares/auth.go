package middlewares

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/viveksingh-01/lumina-api/utils"
)

var jwtSecret = os.Getenv("JWT_SECRET")

type contextKey string

const userContextKey = contextKey("user")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			utils.SendErrorResponse(w, http.StatusUnauthorized, utils.ErrorResponse{
				Error: "Unauthorized: missing auth token",
			})
			return
		}

		tokenStr := strings.TrimSpace(cookie.Value)
		if tokenStr == "" {
			utils.SendErrorResponse(w, http.StatusUnauthorized, utils.ErrorResponse{
				Error: "Unauthorized: empty token",
			})
			return
		}

		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			utils.SendErrorResponse(w, http.StatusUnauthorized, utils.ErrorResponse{
				Error: "Unauthorized: invalid token",
			})
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
