FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /svc-a

# run
FROM gcr.io/distroless/static-debian12
COPY --from=build /svc-a /svc-a
EXPOSE 8080
ENTRYPOINT ["/svc-a"]