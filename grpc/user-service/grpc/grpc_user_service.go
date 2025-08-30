package grpcservices

import (
	"context"
	"fmt"
	"user-service/pb"

	"go.opentelemetry.io/otel"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
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

	tracer := otel.Tracer("User Service")
	_, span := tracer.Start(ctx, "GetUserByEmail Controller")
	defer span.End()

	return &pb.ApiResponse{Success: true, Message: "User fetched successfully"}, nil
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	fmt.Println("GetUser request received", req)
	return &pb.GetUserResponse{Success: true, Message: "User fetched successfully"}, nil
}
