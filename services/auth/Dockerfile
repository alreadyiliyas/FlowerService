FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
COPY services ./services

RUN go build -o /bin/auth-service ./services/auth/cmd/api

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /bin/auth-service /app/auth-service

EXPOSE 8080

ENTRYPOINT ["/app/auth-service"]

