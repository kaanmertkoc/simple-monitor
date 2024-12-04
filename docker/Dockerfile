FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o monitor ./api

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/
COPY --from=builder /app/monitor .

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget -qO- http://localhost:8080/health || exit 1

EXPOSE 8080

CMD ["./monitor"]