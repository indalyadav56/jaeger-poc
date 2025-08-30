package domain

import "context"

type AuthUsecase interface {
	AuthLogin(ctx context.Context) (string, error)
}
