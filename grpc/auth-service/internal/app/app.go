package app

import (
	grpcservices "auth-service/internal/delivery/grpc"
	"auth-service/internal/usecase"
	"auth-service/pb"
	"bytes"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"

	"common-service/pkg/logger"
	"common-service/pkg/trace"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	ctx            context.Context
	grpcServer     *grpc.Server
	tp             *trace.Tracer
	userGrpcClient pb.UserServiceClient
	httpServer     *http.ServeMux
}

func NewApp(ctx context.Context) (*App, error) {
	// logging
	logger.InitLogger("debug")

	// tracing
	tp, err := trace.InitTracer(ctx, "jaeger:4318", "AUTH_SERVICE")
	if err != nil {
		log.Fatalf("Failed to initialize tracer: %v", err)
		return nil, err
	}

	// grpc client
	grpcClient, err := grpc.NewClient("user-service:50051",
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

	// grpc http gateway
	gwMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()} // disable TLS for local dev
	err = pb.RegisterAuthServiceHandlerFromEndpoint(ctx, gwMux, "localhost:50052", opts)
	if err != nil {
		log.Fatalf("failed to start HTTP gateway: %v", err)
	}

	// Standard mux http server
	httpServer := http.NewServeMux()

	// Mount grpc-gateway at root
	httpServer.Handle("/", gwMux)

	return &App{
		ctx:            ctx,
		grpcServer:     grpcServer,
		tp:             tp,
		userGrpcClient: userGrpcClient,
		httpServer:     httpServer,
	}, nil
}

func (a *App) Run() error {
	// Start server
	port := ":50052"

	fmt.Printf("gRPC server starting on port %s \n", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Start gRPC server in a separate goroutine
	go func() {
		log.Fatal(a.grpcServer.Serve(lis))
	}()

	// health check endpoint
	a.httpServer.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// // swagger endpoint
	// a.httpServer.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
	// 	file, err := os.OpenFile("api/swagger/auth.swagger.json", os.O_RDONLY, 0644)
	// 	if err != nil {
	// 		http.Error(w, "Swagger file not found", http.StatusNotFound)
	// 		return
	// 	}
	// 	defer file.Close()

	// 	stat, err := file.Stat()
	// 	if err != nil {
	// 		http.Error(w, "Could not obtain stat", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	http.ServeContent(w, r, "swagger.json", stat.ModTime(), file)
	// })
	a.httpServer.HandleFunc("/auth/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("api/swagger/auth.swagger.json")
		if err != nil {
			http.Error(w, "Swagger file not found", http.StatusNotFound)
			return
		}

		// Simple string replace (for Swagger 2.0)
		// Replace `"host": ""` with `"host": "auth-service:8080"`
		fixed := bytes.ReplaceAll(data, []byte(`"host": ""`), []byte(`"host": "auth-service:8080"`))

		w.Header().Set("Content-Type", "application/json")
		w.Write(fixed)
	})

	// Serve Swagger UI
	a.httpServer.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8081/swagger.json"), // must point to your swagger.json
	))

	handlerWithCORS := withCORS(a.httpServer)

	// Start HTTP server in a separate goroutine
	go func() {
		httpPort := ":8081"
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
