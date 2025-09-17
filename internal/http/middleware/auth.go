package middleware

import (
	"net/http"
	"strings"

	"dekamond/internal/http/handlers"
	"github.com/golang-jwt/jwt/v5"
)

func JwtAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				handlers.WriteJSON(w, http.StatusUnauthorized, handlers.ApiResponse{Error: "missing_token"})
				return
			}
			tokenStr := parts[1]
			_, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, http.ErrAbortHandler
				}
				return []byte(secret), nil
			})
			if err != nil {
				handlers.WriteJSON(w, http.StatusUnauthorized, handlers.ApiResponse{Error: "invalid_token"})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}