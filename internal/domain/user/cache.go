package user

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("not found")

type CacheStore interface {
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Increment(ctx context.Context, key string) (int64, error)
	SetExpiry(ctx context.Context, key string, ttl time.Duration) error
}