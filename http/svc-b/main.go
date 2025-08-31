package main

import (
	"context"
	"log"
	"net/http"

	"svc-b/pkg/tracing"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	ctx := context.Background()

	// tracing
	tp, err := tracing.NewOTLPTracerProvider(ctx, "SERVICE_B", "jaeger:4318")
	if err != nil {
		log.Fatalf("failed to create tracer provider: %v", err)
	}

	defer func() {
		_ = tp.Shutdown(ctx)
	}()

	// gin server
	router := gin.Default()

	router.Use(otelgin.Middleware("SERVICE_B"))

	router.GET("/ping", func(c *gin.Context) {
		callServiceC(c.Request.Context())
		c.JSON(200, gin.H{"message": "pong"})
	})

	router.Run(":8081")
}

func callServiceC(c context.Context) {
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	req, _ := http.NewRequestWithContext(c, "GET", "http://svc-c:8082/ping", nil)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("failed to call service b: %v", err)
		return
	}
	defer resp.Body.Close()
}
