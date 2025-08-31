package usecase

import (
	"common-service/pkg/trace"
	"context"
	"user-service/internal/domain"
	"user-service/pb"
)

type userUsecase struct {
	productGrpcClient pb.ProductServiceClient
	userRepository    domain.UserRepository
}

func NewUserUsecase(productGrpcClient pb.ProductServiceClient, userRepository domain.UserRepository) *userUsecase {
	return &userUsecase{
		productGrpcClient: productGrpcClient,
		userRepository:    userRepository,
	}
}

func (u *userUsecase) CreateUser(ctx context.Context) error {
	ctx, span := trace.StartSpan(ctx, "UserUsecase.CreateUser")
	defer span.End()

	// create user
	err := u.userRepository.CreateUser(ctx)
	if err != nil {
		return err
	}

	// call product service
	_, err = u.productGrpcClient.GetProduct(ctx, &pb.GetProductRequest{Id: "test-id"})
	if err != nil {
		return err
	}

	return nil
}
