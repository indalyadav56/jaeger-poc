FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /svc-c

EXPOSE 8082

CMD ["/svc-c"]