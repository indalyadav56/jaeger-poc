package grpcservices

import (
	"auth-service/internal/domain"
	"auth-service/pb"
	"context"
)

type authService struct {
	pb.UnimplementedAuthServiceServer
	authUsecase domain.AuthUsecase
}

func NewAuthService(userGrpcClient pb.UserServiceClient, authUsecase domain.AuthUsecase) *authService {
	return &authService{authUsecase: authUsecase}
}

func (s *authService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	s.authUsecase.Login(ctx)
	return &pb.LoginResponse{Token: "token"}, nil
}

func (s *authService) Register(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	s.authUsecase.Register(ctx)
	return &pb.RegisterUserResponse{Success: true, Message: "User registered successfully"}, nil
}
