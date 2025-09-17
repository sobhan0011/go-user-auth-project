package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"dekamond/internal/http/handlers"
)

func OTPRateLimit(limiter RateLimiter, maxRequests int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var body struct{ Phone string `json:"phone"` }
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				handlers.WriteJSON(w, http.StatusBadRequest, handlers.ApiResponse{Error: "invalid_phone"})
				return
			}

			phone := strings.TrimSpace(body.Phone)
			if phone == "" {
				handlers.WriteJSON(w, http.StatusBadRequest, handlers.ApiResponse{Error: "invalid_phone"})
				return
			}

			rlKey := fmt.Sprintf("otp:rl:%s", phone)
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

			bodyBytes, _ := json.Marshal(body)
			r.Body = &bodyReader{data: bodyBytes}
			
			next.ServeHTTP(w, r)
		})
	}
}

type bodyReader struct {
	data []byte
	pos  int
}

func (br *bodyReader) Read(p []byte) (n int, err error) {
	if br.pos >= len(br.data) {
		return 0, nil
	}
	n = copy(p, br.data[br.pos:])
	br.pos += n
	return n, nil
}

func (br *bodyReader) Close() error {
	return nil
}