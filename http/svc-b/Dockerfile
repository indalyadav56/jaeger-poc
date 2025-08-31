FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /svc-b

EXPOSE 8081

CMD ["/svc-b"]