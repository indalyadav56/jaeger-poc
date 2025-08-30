package grpcservices

import (
	"context"
	"fmt"
	"user-service/internal/domain"
	"user-service/pb"

	"go.opentelemetry.io/otel"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	userUsecase domain.UserUsecase
}

func NewUserService(userUsecase domain.UserUsecase) *UserService {
	return &UserService{userUsecase: userUsecase}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	fmt.Println("CreateUser request received", req)

	tracer := otel.Tracer("User Service")
	_, span := tracer.Start(ctx, "CreateUser Controller")
	defer span.End()

	return &pb.RegisterResponse{Success: true, Message: "User created successfully"}, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.ApiResponse, error) {
	fmt.Println("GetUserByEmail request received", req)

	s.userUsecase.CreateUser(ctx)

	return &pb.ApiResponse{Success: true, Message: "User fetched successfully"}, nil
}
