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

func (a *authUsecase) AuthLogin(ctx context.Context) (string, error) {
	fmt.Println("AuthLogin usecase")

	// tracer := otel.Tracer("Auth Service")
	// _, span := tracer.Start(ctx, "AuthLogin")
	// defer span.End()

	ctx, span := trace.StartSpan(ctx, "AuthService.Login")
	defer span.End()

	user, err := a.UserGrpcClient.GetUserByEmail(ctx, &pb.GetUserByEmailRequest{Email: "test@test.com"})
	if err != nil {
		return "", err
	}
	fmt.Println("User fetched successfully", user)
	return "", nil
}
