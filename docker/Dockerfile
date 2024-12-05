# Build stage
FROM golang:1.21-alpine AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY api/ ./api/
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o monitor ./api/main.go

# Runtime stage
FROM alpine:latest
RUN apk add --no-cache sqlite
WORKDIR /app
COPY --from=builder /build/monitor .
RUN mkdir -p /data
ENV DATA_DIR=/data
EXPOSE 8080
CMD ["./monitor"]