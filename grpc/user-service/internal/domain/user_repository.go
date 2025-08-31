package domain

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context) error
}
