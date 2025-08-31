package domain

import "context"

type UserUsecase interface {
	CreateUser(ctx context.Context) error
}
