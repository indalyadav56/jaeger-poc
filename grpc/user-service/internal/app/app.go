package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"user-service/internal/config"
	grpcservices "user-service/internal/delivery/grpc"
	"user-service/internal/repository"
	"user-service/internal/usecase"
	"user-service/pb"

	"common-service/pkg/logger"
	"common-service/pkg/trace"

	"common-service/pkg/db"

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

func NewApp(ctx context.Context) (*App, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return nil, err
	}

	// logging
	logger.InitLogger("debug")

	// tracing
	tp, err := trace.InitTracer(ctx, cfg.App.Trace.Endpoint, "USER_SERVICE")
	if err != nil {
		log.Fatalf("Failed to initialize tracer: %v", err)
		return nil, err
	}

	// db
	dbConn, err := db.InitDB(ctx, "postgres", config.GetDatabaseDSN(cfg))
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// apply migrations
	err = db.ApplyMigrations(dbConn, "postgres", "migrations")
	if err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	// grpc client
	grpcClient, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.Clients.ProductService.Target, cfg.Clients.ProductService.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		slog.Error("Failed to create grpc client", "error", err)
		return nil, err
	}
	productGrpcClient := pb.NewProductServiceClient(grpcClient)

	// grpc server
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	// repository
	userRepository := repository.NewUserRepository(dbConn)

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
	}, nil
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
