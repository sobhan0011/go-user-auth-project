package middleware

import (
	"context"
	"time"
)

type RateLimiter interface {
	Increment(ctx context.Context, key string) (int64, error)
	SetExpiry(ctx context.Context, key string, ttl time.Duration) error
}