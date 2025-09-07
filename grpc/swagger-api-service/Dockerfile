FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . .

# install node + npm + swagger-combine + gettext for envsubst
RUN apk add --no-cache nodejs npm gettext \
    && npm install -g swagger-combine

RUN CGO_ENABLED=0 go build -o /swagger

EXPOSE 8085

# Use absolute path to the binary
CMD ["/swagger"]
