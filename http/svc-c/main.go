package main

import (
	"context"
	"log"

	"svc-c/pkg/tracing"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	ctx := context.Background()

	// tracing
	tp, err := tracing.NewOTLPTracerProvider(ctx, "SERVICE_C", "jaeger:4318")
	if err != nil {
		log.Fatalf("failed to create tracer provider: %v", err)
	}

	defer func() {
		_ = tp.Shutdown(ctx)
	}()

	// gin server
	router := gin.Default()

	router.Use(otelgin.Middleware("SERVICE_C"))

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	router.Run(":8082")
}
