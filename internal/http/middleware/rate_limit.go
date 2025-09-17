package middleware

import (
	"fmt"
	"net/http"
	"time"

	"dekamond/internal/http/handlers"
)

func RateLimit(limiter RateLimiter, maxRequests int, window time.Duration, keyFunc func(r *http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := keyFunc(r)
			if key == "" {
				next.ServeHTTP(w, r)
				return
			}

			rlKey := fmt.Sprintf("rate_limit:%s", key)
			cnt, err := limiter.Increment(r.Context(), rlKey)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			if cnt == 1 {
				limiter.SetExpiry(r.Context(), rlKey, window)
			}

			if cnt > int64(maxRequests) {
				handlers.WriteJSON(w, http.StatusTooManyRequests, handlers.ApiResponse{
					Error: "rate limit exceeded",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}