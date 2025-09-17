package user

import (
	"context"
	"strings"

	userdomain "dekamond/internal/domain/user"
)

type ListQuery struct {
	Phone string
	Page  int
	Limit int
}

type Page[T any] struct {
	Items []T `json:"items"`
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func (uuc *UserUsecase) List(ctx context.Context, q ListQuery) (Page[userdomain.User], error) {
	q.Phone = strings.TrimSpace(q.Phone)

	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 20
	}
	offset := (q.Page - 1) * q.Limit

	items, total, err := uuc.users.List(ctx, q.Phone, q.Limit, offset)
	if err != nil {
		return Page[userdomain.User]{}, err
	}
	return Page[userdomain.User]{
		Items: items,
		Total: total,
		Page:  q.Page,
		Limit: q.Limit,
	}, nil
}
