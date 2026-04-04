FROM golang:1.24.0-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
COPY pkg ./pkg
COPY services ./services

RUN go build -o /bin/auth-service ./services/auth/cmd/api

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /bin/auth-service /app/auth-service

EXPOSE 8081
EXPOSE 9090

ENTRYPOINT ["/app/auth-service"]

