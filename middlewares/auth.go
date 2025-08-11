package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/viveksingh-01/lumina-api/utils"
)

var jwtSecret = os.Getenv("JWT_SECRET")

type contextKey string

const userContextKey = contextKey("user")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := GetClaimsFromAuthCookie(w, r)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusUnauthorized, utils.ErrorResponse{
				Error: "Unauthorized: " + err.Error(),
			})
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetClaimsFromAuthCookie(w http.ResponseWriter, r *http.Request) (*jwt.RegisteredClaims, error) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		return nil, fmt.Errorf("missing auth token")
	}

	tokenStr := strings.TrimSpace(cookie.Value)
	if tokenStr == "" {
		return nil, fmt.Errorf("empty token")
	}

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Additional claim checks (just to be safe)
	now := time.Now()
	if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(now) {
		return nil, fmt.Errorf("token expired")
	}
	if claims.IssuedAt != nil && claims.IssuedAt.Time.After(now) {
		return nil, fmt.Errorf("token issued in the future")
	}
	if claims.NotBefore != nil && claims.NotBefore.Time.After(now) {
		return nil, fmt.Errorf("token not active yet")
	}

	return claims, nil
}

func GetClaimsFromContext(ctx context.Context) *jwt.RegisteredClaims {
	if claims, ok := ctx.Value(userContextKey).(*jwt.RegisteredClaims); ok {
		return claims
	}
	return nil
}
