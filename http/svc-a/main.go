package main

import (
	"context"
	"log"
	"net/http"

	"svc-a/pkg/tracing"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	ctx := context.Background()

	// tracing
	tp, err := tracing.NewOTLPTracerProvider(ctx, "SERVICE_A", "jaeger:4318")
	if err != nil {
		log.Fatalf("failed to create tracer provider: %v", err)
	}

	defer func() {
		_ = tp.Shutdown(ctx)
	}()

	// gin server
	router := gin.Default()

	router.Use(otelgin.Middleware("SERVICE_A"))

	router.GET("/ping", func(c *gin.Context) {
		callServiceB(c.Request.Context())
		c.JSON(200, gin.H{"message": "pong"})
	})

	router.Run(":8080")
}

func callServiceB(c context.Context) {
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	req, _ := http.NewRequestWithContext(c, "GET", "http://svc-b:8081/ping", nil)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("failed to call service b: %v", err)
		return
	}
	defer resp.Body.Close()
}
