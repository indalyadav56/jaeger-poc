package grpcservices

import (
	"auth-service/pb"
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	UserGrpcClient pb.UserServiceClient
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	fmt.Println("Login request received", req)

	tracer := otel.Tracer("Auth Service")
	_, span := tracer.Start(ctx, "Login Controller")
	defer span.End()

	user, err := s.UserGrpcClient.GetUserByEmail(ctx, &pb.GetUserByEmailRequest{Email: req.Username})
	if err != nil {
		return nil, err
	}
	fmt.Println("User fetched successfully", user)

	return &pb.LoginResponse{Token: "token"}, nil
}
