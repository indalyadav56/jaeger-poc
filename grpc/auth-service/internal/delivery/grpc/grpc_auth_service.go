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
	s.authUsecase.AuthLogin(ctx)
	return &pb.LoginResponse{Token: "token"}, nil
}
