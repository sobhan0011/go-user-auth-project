package user

import (
	"context"
	"errors"

	userdomain "dekamond/internal/domain/user"
)

var (
	ErrNotFound  = errors.New("user not found")
	ErrInvalidID = errors.New("invalid user id")
)

func (uuc *UserUsecase) GetByID(ctx context.Context, id string) (*userdomain.User, error) {
	if id == "" {
		return nil, ErrInvalidID
	}
	usr, err := uuc.users.GetByID(ctx, id)
	if err != nil || usr == nil {
		return nil, ErrNotFound
	}
	return usr, nil
}
