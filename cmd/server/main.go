package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dekamond/internal/config"
	apphttp "dekamond/internal/http"
	"dekamond/internal/infra/cache"
	"dekamond/internal/infra/db/postgres"
)

func main() {
    conf := config.Load()

    pg, err := postgres.NewPostgres(conf)
    if err != nil {
        log.Fatalf("failed to connect to postgres: %v", err)
    }
    defer pg.Close()

    migrationPath := "/migrations"
    if _, err := os.Stat("/migrations"); err == nil {
        migrationPath = "/migrations"
    }
    if err := postgres.RunMigrations(conf.PostgresURL, migrationPath); err != nil {
        panic(err)
    }

    redis := cache.NewRedis(conf)
    defer func() {
        if err := redis.Close(); err != nil {
            log.Println("failed to close redis:", err)
        }
    }()

    router := apphttp.NewRouter(conf, pg, redis)

    server := &http.Server{
        Addr:              ":" + conf.HTTPPort,
        Handler:           router,
        ReadHeaderTimeout: 10 * time.Second,
        ReadTimeout:       15 * time.Second,
        WriteTimeout:      15 * time.Second,
        IdleTimeout:       60 * time.Second,
    }

    go func() {
        log.Printf("server listening on :%s", conf.HTTPPort)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen: %v", err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if err := server.Shutdown(ctx); err != nil {
        log.Printf("server shutdown error: %v", err)
    }
}


