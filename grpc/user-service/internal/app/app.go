package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	grpcservices "user-service/internal/delivery/grpc"
	"user-service/internal/repository"
	"user-service/internal/usecase"
	"user-service/pb"

	"common-service/pkg/logger"
	"common-service/pkg/trace"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	ctx        context.Context
	grpcServer *grpc.Server

	grpcClient        *grpc.ClientConn
	productGrpcClient pb.ProductServiceClient

	tp *trace.Tracer
}

func NewApp(ctx context.Context) *App {
	// logging
	logger.InitLogger("debug")

	// tracing
	tp, err := trace.InitTracer(ctx, "localhost:4318", "USER_SERVICE")
	if err != nil {
		log.Fatalf("Failed to initialize tracer: %v", err)
	}

	grpcClient, err := grpc.NewClient("localhost:50053",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		slog.Error("Failed to create grpc client", "error", err)
		return nil
	}
	productGrpcClient := pb.NewProductServiceClient(grpcClient)

	// grpc server
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	// repository
	userRepository := repository.NewUserRepository()

	// usecase
	userUsecase := usecase.NewUserUsecase(productGrpcClient, userRepository)

	// register services
	pb.RegisterUserServiceServer(grpcServer, grpcservices.NewUserService(userUsecase))

	return &App{
		ctx:               ctx,
		grpcClient:        grpcClient,
		productGrpcClient: productGrpcClient,
		tp:                tp,
		grpcServer:        grpcServer,
	}
}

func (a *App) Run() error {
	port := ":50051"

	fmt.Printf("User service starting on port %s \n", port)
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
