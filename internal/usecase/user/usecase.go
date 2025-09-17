package user

import (
	userdomain "dekamond/internal/domain/user"
)

type UserUsecase struct {
	users userdomain.Repository
}

func New(users userdomain.Repository) *UserUsecase {
	return &UserUsecase{users: users}
}
