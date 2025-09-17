package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	redis *redis.Client
}

func NewRedisStore(redis *redis.Client) *RedisStore {
	return &RedisStore{redis: redis}
}

func (s *RedisStore) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return s.redis.Set(ctx, key, value, ttl).Err()
}

func (s *RedisStore) Get(ctx context.Context, key string) (string, error) {
	return s.redis.Get(ctx, key).Result()
}

func (s *RedisStore) Delete(ctx context.Context, key string) error {
	return s.redis.Del(ctx, key).Err()
}

func (s *RedisStore) Increment(ctx context.Context, key string) (int64, error) {
	return s.redis.Incr(ctx, key).Result()
}

func (s *RedisStore) SetExpiry(ctx context.Context, key string, ttl time.Duration) error {
	return s.redis.Expire(ctx, key, ttl).Err()
}