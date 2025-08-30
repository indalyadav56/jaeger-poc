package trace

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// StartSpan starts a span with given name.
// Use this in all layers (controller, usecase, repo).
func StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	tr := otel.Tracer("app")
	return tr.Start(ctx, name)
}
