# Stage 1: Build the Go application with CGO support
FROM --platform=$BUILDPLATFORM golang:1.21 AS builder

# Set build arguments to properly cross-compile
ARG TARGETOS=linux
ARG TARGETARCH=amd64

WORKDIR /build

# Install necessary build dependencies, including cross-compilers
RUN apt-get update && apt-get install -y \
    gcc g++ \
    musl-dev \
    sqlite3 libsqlite3-dev \
    build-essential \
    crossbuild-essential-amd64 crossbuild-essential-arm64 \
    && rm -rf /var/lib/apt/lists/*

# Set the appropriate cross-compiler based on the target architecture
RUN if [ "$TARGETARCH" = "amd64" ]; then export CC=x86_64-linux-gnu-gcc; fi && \
    if [ "$TARGETARCH" = "arm64" ]; then export CC=aarch64-linux-gnu-gcc; fi

# Copy dependency files and download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the application source code
COPY api/ ./api/

# Install required build tools
RUN apt-get update && apt-get install -y \
    build-essential \
    gcc \
    g++ \
    libc6-dev \
    && rm -rf /var/lib/apt/lists/*

# Build the Go application with CGO enabled or disabled as needed
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags="-w -s" -o monitor ./api

# Stage 2: Create the runtime image
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache sqlite

WORKDIR /app
COPY --from=builder /build/monitor .

# Create a data directory
RUN mkdir -p /data
ENV DATA_DIR=/data

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./monitor"]