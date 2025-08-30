package app

import (
	"context"
	"fmt"
	"log"
	"net"
	grpcservices "product-service/internal/delivery/grpc"
	"product-service/pb"

	"common-service/pkg/logger"
	"common-service/pkg/trace"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type App struct {
	ctx        context.Context
	grpcServer *grpc.Server

	tp *trace.Tracer
}

func NewApp(ctx context.Context) *App {
	// logging
	logger.InitLogger("debug")

	// tracing
	tp, err := trace.InitTracer(ctx, "localhost:4318", "PRODUCT_SERVICE")
	if err != nil {
		log.Fatalf("Failed to initialize tracer: %v", err)
	}

	// grpc server
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	// register services
	pb.RegisterProductServiceServer(grpcServer, grpcservices.NewProductGrpcService())

	return &App{
		ctx:        ctx,
		grpcServer: grpcServer,
		tp:         tp,
	}
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
