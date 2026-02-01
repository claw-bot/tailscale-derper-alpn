# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the derper binary
RUN CGO_ENABLED=0 GOOS=linux go build -o derper ./cmd/derper

# Runtime stage
FROM alpine:latest

# Install ca-certificates for TLS connections
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Create certificate directory
RUN mkdir -p /app/certs

# Copy binary from builder
COPY --from=builder /app/derper /app/derper

# Environment variables with defaults
ENV DERP_DOMAIN=your-hostname.com
ENV DERP_CERT_MODE=letsencrypt
ENV DERP_CERT_DIR=/app/certs
ENV DERP_ADDR=:443
ENV DERP_STUN=true
ENV DERP_STUN_PORT=3478
ENV DERP_HTTP_PORT=80
ENV DERP_VERIFY_CLIENTS=false
ENV DERP_VERIFY_CLIENT_URL=""
ENV DERP_ACME_EMAIL=""

# Expose ports
# 443 - HTTPS/TLS
# 80 - HTTP (for ACME HTTP-01 challenge if needed)
# 3478 - STUN (UDP)
EXPOSE 443/tcp 80/tcp 3478/udp

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:80/generate_204 || exit 1

# Run derper
CMD ["/app/derper", \
     "--hostname=${DERP_DOMAIN}", \
     "--certmode=${DERP_CERT_MODE}", \
     "--certdir=${DERP_CERT_DIR}", \
     "--a=${DERP_ADDR}", \
     "--stun=${DERP_STUN}", \
     "--stun-port=${DERP_STUN_PORT}", \
     "--http-port=${DERP_HTTP_PORT}", \
     "--verify-clients=${DERP_VERIFY_CLIENTS}", \
     "--verify-client-url=${DERP_VERIFY_CLIENT_URL}"]
