FROM golang:1.21-alpine AS builder

WORKDIR /build

# Copy only the files needed for go mod download first
COPY go.mod go.sum ./
RUN go mod download

# Install only the necessary build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Copy source code
COPY api/ ./api/

# Build with optimization flags
RUN CGO_ENABLED=1 GOOS=linux GOARCH=$TARGETARCH \
    go build -ldflags="-w -s" -o monitor ./api

FROM alpine:latest

# Install only the required runtime dependencies
RUN apk add --no-cache sqlite tzdata

WORKDIR /app
COPY --from=builder /build/monitor .

# Create data directory
RUN mkdir -p /data

# Set data directory environment variable
ENV DATA_DIR=/data

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget -qO- http://localhost:8080/health || exit 1

CMD ["./monitor"]