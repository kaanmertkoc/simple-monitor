FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .

# Install build dependencies (including SQLite)
RUN apk add --no-cache gcc musl-dev sqlite-dev

RUN go mod download
# Enable CGO for SQLite
RUN CGO_ENABLED=1 GOOS=linux go build -o monitor ./api

FROM alpine:latest

# Install SQLite and timezone data
RUN apk --no-cache add sqlite ca-certificates tzdata

WORKDIR /root/
COPY --from=builder /app/monitor .

# Create data directory
RUN mkdir -p /data

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget -qO- http://localhost:8080/health || exit 1

EXPOSE 8080

# Set data directory environment variable
ENV DATA_DIR=/data

CMD ["./monitor"]