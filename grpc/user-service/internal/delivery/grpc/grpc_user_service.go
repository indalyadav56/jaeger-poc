package grpcservices

import (
	"context"
	"user-service/internal/domain"
	"user-service/pb"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	userUsecase domain.UserUsecase
}

func NewUserService(userUsecase domain.UserUsecase) *UserService {
	return &UserService{userUsecase: userUsecase}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	err := s.userUsecase.CreateUser(ctx)
	if err != nil {
		return &pb.RegisterResponse{Success: false, Message: err.Error()}, err
	}
	return &pb.RegisterResponse{Success: true, Message: "User created successfully"}, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.ApiResponse, error) {
	return &pb.ApiResponse{Success: true, Message: "User fetched successfully"}, nil
}
