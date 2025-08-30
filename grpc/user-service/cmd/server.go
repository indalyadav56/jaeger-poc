package main

import (
	"context"
	"log"
	"net"
	grpcservices "user-service/grpc"
	"user-service/pb"

	"common-service/pkg/trace"

	"google.golang.org/grpc"
)

func main() {
	tp, err := trace.InitTracer("localhost:4318", "USER_SERVICE")
	if err != nil {
		log.Fatalf("Failed to initialize tracer: %v", err)
	}

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	port := ":50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcServer, &grpcservices.UserService{})

	log.Printf("User service starting on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
