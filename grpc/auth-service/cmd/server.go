package main

import (
	"context"
	"log"
	"net"

	grpcservices "auth-service/grpc"
	"auth-service/pb"
	"common-service/pkg/trace"

	"google.golang.org/grpc"
)

func main() {
	tp, err := trace.InitTracer("localhost:4318", "AUTH_SERVICE")
	if err != nil {
		log.Fatalf("Failed to initialize tracer: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// Start server
	port := ":50052"
	log.Printf("Service B gRPC server starting on port %s", port)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcClient, err := grpc.NewClient("user-service:50051")
	if err != nil {
		log.Fatalf("Failed to create user grpc client: %v", err)
	}
	userGrpcClient := pb.NewUserServiceClient(grpcClient)

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, &grpcservices.AuthService{UserGrpcClient: userGrpcClient})

	log.Fatal(grpcServer.Serve(lis))
}
