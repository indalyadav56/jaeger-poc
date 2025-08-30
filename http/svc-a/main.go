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
	tracer = otel.Tracer("svc-a")

	// Setup HTTP server
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/work", handleWork)

	port := ":8080"
	log.Printf("Service A starting on port %s", port)
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
	time.Sleep(10 * time.Millisecond)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Hello from Service A!", "service": "svc-a"}`))
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
	time.Sleep(50 * time.Millisecond)

	// Call downstream service
	downstreamURL := os.Getenv("DOWNSTREAM_URL")
	if downstreamURL != "" {
		callDownstreamService(ctx, downstreamURL)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Work completed by Service A!", "service": "svc-a"}`))
}

func callDownstreamService(ctx context.Context, url string) {
	// Create a span for the downstream call
	ctx, span := tracer.Start(ctx, "callDownstreamService")
	defer span.End()

	// Add attributes
	span.SetAttributes(
		semconv.HTTPURL(url),
		semconv.HTTPMethod("GET"),
	)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		span.RecordError(err)
		return
	}

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		span.RecordError(err)
		return
	}
	defer resp.Body.Close()

	// Add response attributes
	span.SetAttributes(
		semconv.HTTPStatusCode(resp.StatusCode),
	)

	log.Printf("Called downstream service: %s, status: %d", url, resp.StatusCode)
}
