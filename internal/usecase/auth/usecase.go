package auth

import (
	"dekamond/internal/config"
	"time"

	userdomain "dekamond/internal/domain/user"
)


type AuthUsecase struct {
	users     userdomain.Repository
	cache     userdomain.CacheStore
	jwtSecret []byte
	tokenTTL  time.Duration
}

func New(users userdomain.Repository, cache userdomain.CacheStore, conf config.Config) *AuthUsecase {
	return &AuthUsecase{users: users, cache: cache, jwtSecret: []byte(conf.JWTSecret), tokenTTL: 24 * time.Hour}
}


