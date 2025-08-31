package domain

import "context"

type AuthUsecase interface {
	Login(ctx context.Context) (string, error)
	Register(ctx context.Context) error
}
