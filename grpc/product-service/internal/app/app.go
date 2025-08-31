package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"product-service/internal/config"
	grpcservices "product-service/internal/delivery/grpc"
	"product-service/internal/repository"
	"product-service/internal/usecase"
	"product-service/pb"

	"common-service/pkg/logger"
	"common-service/pkg/trace"

	"common-service/pkg/db/mongodb"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type App struct {
	ctx        context.Context
	grpcServer *grpc.Server

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
	tp, err := trace.InitTracer(ctx, cfg.App.Trace.Endpoint, "PRODUCT_SERVICE")
	if err != nil {
		log.Fatalf("Failed to initialize tracer: %v", err)
		return nil, err
	}

	// mongodb
	mongodbClient, err := mongodb.NewMongoDBClient(mongodb.Config{
		URI:      cfg.DB.MongoDB.URI,
		Database: cfg.DB.MongoDB.Database,
	})
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
		return nil, err
	}

	// repository
	productRepository := repository.NewMongodbProductRepository(mongodbClient)

	// usecase
	productUsecase := usecase.NewProductUsecase(productRepository)

	// grpc server
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	// register services
	pb.RegisterProductServiceServer(grpcServer, grpcservices.NewProductGrpcService(productUsecase))

	return &App{
		ctx:        ctx,
		grpcServer: grpcServer,
		tp:         tp,
	}, nil
}

func (a *App) Run() error {
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		return err
	}

	fmt.Println("Starting gRPC server on port 50053")
	if err := a.grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}

func (a *App) Shutdown() error {
	a.grpcServer.Stop()
	return nil
}
