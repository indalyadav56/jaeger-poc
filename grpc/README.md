# Distributed Tracing and Centerlized Swagger Service

A comprehensive microservices architecture demonstrating distributed tracing with Jaeger And centralized API documentation with Swagger and Scalar For API Doc, and gRPC communication between services.

## 🏗️ Architecture Overview

This project implements a microservices architecture with the following key components:

- **Distributed Tracing**: Jaeger for observability and request tracing
- **Centralized API Documentation**: Swagger UI and Scalar API Reference
- **gRPC Communication**: Inter-service communication using Protocol Buffers
- **Multi-Database Support**: PostgreSQL and MongoDB for different data requirements
- **Containerized Deployment**: Docker and Docker Compose for easy deployment

## 🚀 Services

### Core Services
- **User Service** (`:50051`) - User management with PostgreSQL
- **Product Service** (`:50053`) - Product catalog with MongoDB  
- **Auth Service** (`:50052, :8081`) - Authentication and authorization
- **Swagger API Service** (`:8085`) - Centralized API documentation

### Infrastructure Services
- **Jaeger** (`:16686`) - Distributed tracing and observability
- **PostgreSQL** (`:5432`) - Relational database for users
- **MongoDB** (`:27017`) - Document database for products
- **Elasticsearch** (`:9200`) - Store tracing data (optional)

## 🎯 Key Features

### Distributed Tracing
- **OpenTelemetry Integration**: Standardized tracing with OTLP exporter
- **Jaeger Visualization**: Complete request flow visualization
- **Span Correlation**: Automatic trace context propagation between services
- **Performance Monitoring**: Latency and dependency analysis

### Centralized API Documentation
- **Swagger UI**: Interactive API documentation at `/swagger/*`
- **Scalar API Reference**: Modern, responsive API reference at `/`
- **Unified Documentation**: Single source of truth for all service APIs
- **Real-time Updates**: Automatic synchronization with service changes

### gRPC Communication
- **Protocol Buffers**: Efficient serialization and type safety
- **Bidirectional Streaming**: Support for complex communication patterns
- **Service Discovery**: Built-in service resolution and load balancing
- **Health Checks**: gRPC health check protocol implementation

## 🛠️ Technology Stack

- **Language**: Go 1.24+
- **Communication**: gRPC with Protocol Buffers
- **Tracing**: Jaeger
- **Databases**: PostgreSQL, MongoDB
- **Containerization**: Docker + Docker Compose
- **API Documentation**: Swagger/OpenAPI + Scalar
- **Monitoring**: Jaeger UI, Health checks

## 📋 Prerequisites

- Docker and Docker Compose
- Go 1.24 or higher (for local development)
- Protocol Buffer compiler (protoc)
- Make (for build automation)

## 🚀 Quick Start

### 1. Clone the Repository
```bash
git clone <repository-url>
cd grpc
```

### 2. Start All Services
```bash
docker-compose up --build
```

### 3. Access Services
- **Jaeger UI**: http://localhost:16686
- **Scalar API Docs**: http://localhost:8085
- **Swagger API Docs**: http://localhost:8085/swagger/index.html
- **User Service**: localhost:50051
- **Product Service**: localhost:50053
- **Auth Service**: localhost:50052

### 4. View API Documentation
- **Swagger UI**: http://localhost:8085/swagger/
- **Scalar API Reference**: http://localhost:8085/

## 🔧 Development

### Local Development Setup
```bash
# Install dependencies
go mod download

# Generate Protocol Buffer code
make proto

# Run tests
make test

# Build services
make build
```

### Service Development
Each service follows a clean architecture pattern:
```
service/
├── api/           # Protocol Buffers and Swagger definitions
├── cmd/           # Application entry points
├── internal/      # Private application code
│   ├── app/       # Application configuration
│   ├── delivery/  # Transport layer (gRPC/HTTP)
│   ├── domain/    # Business logic and entities
│   ├── repository/# Data access layer
│   └── usecase/   # Application use cases
├── configs/       # Configuration files
└── pkg/           # Public packages
```

### Adding New Services
1. Create service directory with standard structure
2. Define Protocol Buffer interfaces
3. Implement gRPC service
4. Add tracing integration
5. Update docker-compose.yaml
6. Add to centralized Swagger documentation

## 📊 Observability

### Distributed Tracing
The system implements comprehensive distributed tracing:

```go
// Initialize tracer in each service
tracer, err := trace.InitTracer(ctx, "jaeger:4318", "service-name")
if err != nil {
    log.Fatal(err)
}
defer tracer.Shutdown(ctx)

// Create spans for operations
ctx, span := otel.Tracer("").Start(ctx, "operation-name")
defer span.End()
```

### Trace Propagation
- **Automatic Context Propagation**: Trace context flows through gRPC calls
- **Span Correlation**: Related operations are linked across services
- **Performance Metrics**: Latency, throughput, and error rates
- **Dependency Mapping**: Service interaction visualization

### Monitoring Endpoints
- **Health Checks**: `/health` endpoint on each service
- **Metrics**: OpenTelemetry metrics collection
- **Logs**: Structured logging with trace correlation

## 📚 API Documentation

### Centralized Swagger Service
The Swagger API service provides unified documentation:

- **Single Entry Point**: All service APIs in one location
- **Interactive Testing**: Try API endpoints directly from documentation
- **Schema Validation**: Automatic request/response validation
- **Code Generation**: Client SDK generation capabilities

### Documentation Features
- **Swagger UI**: Traditional OpenAPI documentation
- **Scalar API Reference**: Modern, responsive interface
- **Real-time Updates**: Automatic synchronization with service changes
- **Multi-format Support**: JSON, YAML, and interactive formats

## 🗄️ Database Design

### PostgreSQL (User Service)
- **Users Table**: User authentication and profile data
- **Migrations**: Version-controlled schema changes
- **ACID Compliance**: Transactional data integrity

### MongoDB (Product Service)
- **Products Collection**: Flexible product catalog structure
- **Indexing**: Optimized query performance
- **Scalability**: Horizontal scaling capabilities

## 🔒 Security

### Authentication
- **JWT Tokens**: Secure token-based authentication
- **Service-to-Service**: Mutual TLS for inter-service communication
- **Environment Variables**: Secure configuration management

### Data Protection
- **Input Validation**: Protocol Buffer schema validation
- **SQL Injection Prevention**: Parameterized queries
- **NoSQL Injection Prevention**: Document validation

## 🚀 Deployment

### Production Considerations
- **Environment Configuration**: Separate configs for dev/staging/prod
- **Resource Limits**: Docker resource constraints
- **Health Checks**: Comprehensive health monitoring
- **Logging**: Centralized log aggregation
- **Monitoring**: Prometheus metrics and Grafana dashboards

### Scaling
- **Horizontal Scaling**: Service replication
- **Load Balancing**: gRPC load balancer integration
- **Database Scaling**: Read replicas and sharding
- **Cache Layer**: Redis for performance optimization

## 🧪 Testing

### Test Strategy
- **Unit Tests**: Individual component testing
- **Integration Tests**: Service interaction testing
- **End-to-End Tests**: Complete workflow validation
- **Performance Tests**: Load and stress testing

### Running Tests
```bash
# Run all tests
make test

# Run specific service tests
cd user-service && go test ./...

# Run with coverage
go test -cover ./...
```

## 📈 Performance

### Optimization Techniques
- **Connection Pooling**: Database connection management
- **Caching**: Redis for frequently accessed data
- **Compression**: gRPC compression for large payloads
- **Async Processing**: Non-blocking operations

### Monitoring
- **Jaeger Traces**: Request flow analysis
- **Database Metrics**: Query performance monitoring
- **Service Metrics**: Response time and throughput
- **Resource Usage**: CPU, memory, and network monitoring

## 🤝 Contributing

### Development Workflow
1. Fork the repository
2. Create a feature branch
3. Implement changes with tests
4. Update documentation
5. Submit a pull request

### Code Standards
- **Go Formatting**: `gofmt` and `goimports`
- **Linting**: `golangci-lint` configuration
- **Documentation**: Comprehensive code comments
- **Testing**: Minimum 80% test coverage

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **Jaeger**: Distributed tracing platform
- **OpenTelemetry**: Observability framework
- **gRPC**: High-performance RPC framework
- **Swagger**: API documentation tools

## 📞 Support

For questions and support:
- **Issues**: GitHub issue tracker
- **Documentation**: This README and inline code comments
- **Community**: Open source community contributions

---

**Happy Tracing! 🚀**
