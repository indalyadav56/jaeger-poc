package usecase

import (
	"auth-service/pb"
	"common-service/pkg/trace"
	"context"
	"fmt"
)

type authUsecase struct {
	UserGrpcClient pb.UserServiceClient
}

func NewAuthUsecase(userGrpcClient pb.UserServiceClient) *authUsecase {
	return &authUsecase{UserGrpcClient: userGrpcClient}
}

func (a *authUsecase) Login(ctx context.Context) (string, error) {
	ctx, span := trace.StartSpan(ctx, "AuthUsecase.Login")
	defer span.End()

	user, err := a.UserGrpcClient.GetUserByEmail(ctx, &pb.GetUserByEmailRequest{Email: "test@test.com"})
	if err != nil {
		return "", err
	}
	fmt.Println("User fetched successfully", user)
	return "", nil
}

func (a *authUsecase) Register(ctx context.Context) error {
	ctx, span := trace.StartSpan(ctx, "AuthUsecase.Register")
	defer span.End()

	user, err := a.UserGrpcClient.CreateUser(ctx, &pb.RegisterRequest{
		Username: "test",
		Password: "test",
	})
	if err != nil {
		return err
	}
	fmt.Println("User registered successfully", user)

	return nil
}
