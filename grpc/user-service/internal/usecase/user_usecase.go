package usecase

import (
	"common-service/pkg/trace"
	"context"
	"fmt"
	"user-service/internal/domain"
	"user-service/pb"
)

type userUsecase struct {
	productGrpcClient pb.ProductServiceClient

	userRepository domain.UserRepository
}

func NewUserUsecase(productGrpcClient pb.ProductServiceClient, userRepository domain.UserRepository) *userUsecase {
	return &userUsecase{productGrpcClient: productGrpcClient, userRepository: userRepository}
}

func (u *userUsecase) CreateUser(ctx context.Context) (string, error) {
	fmt.Println("CreateUser usecase")

	ctx, span := trace.StartSpan(ctx, "UserService.CreateUser")
	defer span.End()

	// create user
	_, err := u.userRepository.CreateUser(ctx)
	if err != nil {
		return "", err
	}

	_, err = u.productGrpcClient.GetProduct(ctx, &pb.GetProductRequest{Id: "test-id"})
	if err != nil {
		return "", err
	}

	return "", nil
}
