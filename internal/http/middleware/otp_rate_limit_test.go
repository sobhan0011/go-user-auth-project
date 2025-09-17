package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	userdomain "dekamond/internal/domain/user"
)

type mockCacheStore struct {
	store map[string]string
}

func newMockCacheStore() *mockCacheStore {
	return &mockCacheStore{
		store: make(map[string]string),
	}
}

func (m *mockCacheStore) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	m.store[key] = value
	return nil
}

func (m *mockCacheStore) Get(ctx context.Context, key string) (string, error) {
	value, exists := m.store[key]
	if !exists {
		return "", userdomain.ErrNotFound
	}
	return value, nil
}

func (m *mockCacheStore) Delete(ctx context.Context, key string) error {
	delete(m.store, key)
	return nil
}

func (m *mockCacheStore) Increment(ctx context.Context, key string) (int64, error) {
	current, exists := m.store[key]
	if !exists {
		m.store[key] = "1"
		return 1, nil
	}
	count := len(current) + 1
	m.store[key] = string(make([]byte, count))
	return int64(count), nil
}

func (m *mockCacheStore) SetExpiry(ctx context.Context, key string, ttl time.Duration) error {
	return nil
}

func TestOTPRateLimit(t *testing.T) {
	tests := []struct {
		name           string
		phone          string
		requestCount   int
		maxRequests    int
		windowDuration time.Duration
		expectBlocked  bool
	}{
		{
			name:           "within rate limit",
			phone:          "+15551234567",
			requestCount:   2,
			maxRequests:    3,
			windowDuration: 10 * time.Minute,
			expectBlocked:  false,
		},
		{
			name:           "at rate limit",
			phone:          "+15551234567",
			requestCount:   3,
			maxRequests:    3,
			windowDuration: 10 * time.Minute,
			expectBlocked:  false,
		},
		{
			name:           "exceeds rate limit",
			phone:          "+15551234567",
			requestCount:   4,
			maxRequests:    3,
			windowDuration: 10 * time.Minute,
			expectBlocked:  true,
		},
		{
			name:           "different phone numbers",
			phone:          "+15551234568",
			requestCount:   4,
			maxRequests:    3,
			windowDuration: 10 * time.Minute,
			expectBlocked:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := newMockCacheStore()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			})

			middleware := OTPRateLimit(cache, tt.maxRequests, tt.windowDuration)(handler)

			for i := 0; i < tt.requestCount; i++ {
				body := fmt.Sprintf(`{"phone":"%s"}`, tt.phone)
				req := httptest.NewRequest("POST", "/api/auth/request-otp", strings.NewReader(body))
				req.Header.Set("Content-Type", "application/json")

				rr := httptest.NewRecorder()
				middleware.ServeHTTP(rr, req)

				if i < tt.maxRequests {
					if rr.Code != http.StatusOK {
						t.Errorf("Request %d: expected status %d, got %d", i+1, http.StatusOK, rr.Code)
					}
				} else {
					if rr.Code != http.StatusTooManyRequests {
						t.Errorf("Request %d: expected status %d, got %d", i+1, http.StatusTooManyRequests, rr.Code)
					}
				}
			}
		})
	}
}

func TestRateLimitKey(t *testing.T) {
	phone := "+15551234567"
	expected := "otp:rl:+15551234567"
	result := fmt.Sprintf("otp:rl:%s", phone)

	if result != expected {
		t.Errorf("rateLimitKey() = %s, want %s", result, expected)
	}
}
