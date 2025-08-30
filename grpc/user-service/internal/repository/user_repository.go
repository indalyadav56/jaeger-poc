package repository

import (
	"common-service/pkg/trace"
	"context"
	"fmt"
)

type userRepository struct {
}

func NewUserRepository() *userRepository {
	return &userRepository{}
}

func (r *userRepository) CreateUser(ctx context.Context) (string, error) {
	ctx, span := trace.StartSpan(ctx, "UserRepository.CreateUser")
	defer span.End()

	fmt.Println("CreateUser repository")
	return "", nil
}
