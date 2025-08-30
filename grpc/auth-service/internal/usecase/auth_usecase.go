package usecase

import (
	"context"
)

type authUsecase struct {
}

func NewAuthUsecase() *authUsecase {
	return &authUsecase{}
}

func (u *authUsecase) Login(ctx context.Context) (string, error) {
	return "", nil
}
