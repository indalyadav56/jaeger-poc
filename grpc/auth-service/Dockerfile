# Build stage
FROM golang:1.24.5-alpine AS builder

# Install necessary packages
RUN apk add --no-cache git protobuf-dev make

# Set working directory
WORKDIR /app

# Copy go mod files from auth-service
COPY auth-service/go.mod auth-service/go.sum ./

# Copy common-service to parent directory (matching the replace directive)
COPY common-service/ ../common-service/

# Download dependencies
RUN go mod download

# Copy auth-service source code
COPY auth-service/ .

# Install protobuf plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

# Generate protobuf code
RUN make proto

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/server ./cmd/server.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/bin/server .
COPY --from=builder /app/configs /app/configs
COPY --from=builder /app/third_party /app/third_party
COPY --from=builder /app/api /app/api

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose ports
EXPOSE 50052

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8081/health || exit 1

# Run the application
CMD ["./server"]
