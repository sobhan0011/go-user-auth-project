package postgres

import (
	"context"
	"time"

	"dekamond/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres(cfg config.Config) (*pgxpool.Pool, error) {
	cfgPool, err := pgxpool.ParseConfig(cfg.PostgresURL)
	if err != nil {
		return nil, err
	}
	cfgPool.MaxConns = 5
	cfgPool.MinConns = 1
	cfgPool.MaxConnLifetime = time.Hour
	return pgxpool.NewWithConfig(context.Background(), cfgPool)
}