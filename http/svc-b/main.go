package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func initTracer() (*sdktrace.TracerProvider, error) {
	ctx := context.Background()

	// Create OTLP exporter
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint("jaeger:4318"),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	// Create resource
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(os.Getenv("OTEL_SERVICE_NAME")),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// Set global TracerProvider
	otel.SetTracerProvider(tp)

	return tp, nil
}

func main() {
	// Initialize tracer
	tp, err := initTracer()
	if err != nil {
		log.Fatalf("Failed to initialize tracer: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// Get tracer
	tracer = otel.Tracer("svc-b")

	// Setup HTTP server
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/work", handleWork)

	port := ":8081"
	log.Printf("Service B starting on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Create a span for this request
	ctx, span := tracer.Start(ctx, "handleRoot")
	defer span.End()

	// Add some attributes to the span
	span.SetAttributes(
		semconv.HTTPMethod(r.Method),
		semconv.HTTPURL(r.URL.String()),
	)

	// Simulate some work
	time.Sleep(15 * time.Millisecond)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Hello from Service B!", "service": "svc-b"}`))
}

func handleWork(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Create a span for this request
	ctx, span := tracer.Start(ctx, "handleWork")
	defer span.End()

	// Add some attributes to the span
	span.SetAttributes(
		semconv.HTTPMethod(r.Method),
		semconv.HTTPURL(r.URL.String()),
	)

	// Simulate some work
	time.Sleep(100 * time.Millisecond)

	// Add some custom attributes
	span.SetAttributes(
		semconv.HTTPStatusCode(200),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Work completed by Service B!", "service": "svc-b", "processing_time_ms": 100}`))
}
