package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    HTTPPort      string
    JWTSecret     string
    PostgresURL   string
    RedisAddr     string
    RedisPassword string
    RedisDB       int
}

func Load() Config {
    _ = godotenv.Load()

    cfg := Config{
        HTTPPort:    getEnv("HTTP_PORT", "8080"),
        JWTSecret:   getEnv("JWT_SECRET", "sharing-the-secret-would-mean-a-lot-if-it-led-to-an-opportunity-to-join-you-in-dekamond"),
        PostgresURL: getEnv("POSTGRES_DSN", "postgres://postgres:postgres@localhost:5432/otpapp?sslmode=disable"),
        RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
        RedisPassword: getEnv("REDIS_PASSWORD", ""),
    }
    return cfg
}

func getEnv(key, defaultVal string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    if defaultVal == "" {
        log.Printf("warning: env %s not set", key)
    }
    return defaultVal
}