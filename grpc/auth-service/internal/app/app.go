package app

import (
	grpcservices "auth-service/internal/delivery/grpc"
	"auth-service/internal/usecase"
	"auth-service/pb"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"

	"common-service/pkg/logger"
	"common-service/pkg/trace"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	ctx            context.Context
	grpcServer     *grpc.Server
	tp             *trace.Tracer
	userGrpcClient pb.UserServiceClient
}

func NewApp(ctx context.Context) *App {
	// logging
	logger.InitLogger("debug")

	// tracing
	tp, err := trace.InitTracer(ctx, "localhost:4318", "AUTH_SERVICE")
	if err != nil {
		log.Fatalf("Failed to initialize tracer: %v", err)
	}

	// grpc client
	grpcClient, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		log.Fatalf("Failed to create user grpc client: %v", err)
	}

	userGrpcClient := pb.NewUserServiceClient(grpcClient)

	// grpc server
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	// usecase
	authUsecase := usecase.NewAuthUsecase(userGrpcClient)

	// register services
	pb.RegisterAuthServiceServer(grpcServer, grpcservices.NewAuthService(userGrpcClient, authUsecase))

	return &App{
		ctx:            ctx,
		grpcServer:     grpcServer,
		tp:             tp,
		userGrpcClient: userGrpcClient,
	}
}

func (a *App) Run() error {
	// Start server
	port := ":50052"

	fmt.Printf("Service B gRPC server starting on port %s \n", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Fatal(a.grpcServer.Serve(lis))
	return nil
}

func (a *App) Shutdown() error {
	if err := a.tp.TracerProvider.Shutdown(context.Background()); err != nil {
		slog.Error("Error shutting down tracer provider", "error", err)
	}
	return nil
}
