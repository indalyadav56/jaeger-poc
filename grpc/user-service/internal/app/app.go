package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
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

	httpServer *http.ServeMux
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

	// http server
	httpServer := http.NewServeMux()

	return &App{
		ctx:               ctx,
		grpcClient:        grpcClient,
		productGrpcClient: productGrpcClient,
		tp:                tp,
		grpcServer:        grpcServer,
		httpServer:        httpServer,
	}, nil
}

func (a *App) Run() error {
	port := ":50051"

	fmt.Printf("User service starting on port %s \n", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	go func() { log.Fatal(a.grpcServer.Serve(lis)) }()

	// swagger endpoint
	a.httpServer.HandleFunc("/user/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		file, err := os.OpenFile("api/swagger/swagger.json", os.O_RDONLY, 0644)
		if err != nil {
			http.Error(w, "Swagger file not found", http.StatusNotFound)
			return
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			http.Error(w, "Could not obtain stat", http.StatusInternalServerError)
			return
		}

		http.ServeContent(w, r, "swagger.json", stat.ModTime(), file)

	})

	handlerWithCORS := withCORS(a.httpServer)

	// Start HTTP server in a separate goroutine
	go func() {
		httpPort := ":8080"
		fmt.Printf("HTTP server starting on port %s \n", httpPort)
		if err := http.ListenAndServe(httpPort, handlerWithCORS); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	<-a.ctx.Done()

	return nil
}

func (a *App) Shutdown() error {
	if err := a.tp.TracerProvider.Shutdown(context.Background()); err != nil {
		slog.Error("Error shutting down tracer provider", "error", err)
	}
	return nil
}

// CORS middleware
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Pass to the next handler
		next.ServeHTTP(w, r)
	})
}
