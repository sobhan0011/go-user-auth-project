package user

import "context"

type Repository interface {
	GetByPhone(ctx context.Context, phone string) (*User, error)
	Create(ctx context.Context, phone string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	List(ctx context.Context, phone string, limit, offset int) ([]User, int, error)
}