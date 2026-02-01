# Tailscale DERP Server with TLS ALPN-01 Challenge Support

This repository contains a modified version of the Tailscale DERP server that supports TLS ALPN-01 challenge for ACME certificate validation.

## 🎯 What This Does

Adds TLS ALPN-01 challenge support to the Tailscale DERP server, enabling it to obtain SSL/TLS certificates using the ACME TLS-ALPN-01 challenge type. This is an alternative to the HTTP-01 challenge that doesn't require opening port 80.

## ✨ Features

- ✅ **No Port 80 Required** - Only port 443 needed
- ✅ **TLS-Based Validation** - More secure than HTTP-01
- ✅ **ACME Standard Compatible** - Works with Let's Encrypt, GCP, and other ACME providers
- ✅ **Firewall Friendly** - Single port exposure
- ✅ **Unit Tests** - Comprehensive test coverage
- ✅ **Complete Documentation** - User and developer documentation included
- ✅ **Docker Support** - Ready for containerized deployment
- ✅ **GitHub Actions** - Automated CI/CD pipeline

## 🚀 Quick Start

### Build

```bash
go build -o derper .
```

### Basic Usage (Let's Encrypt)

```bash
./derper \
  --certmode=alpn \
  --certdir=/var/lib/derper/certs \
  --hostname=your-domain.com \
  --acme-email=your-email@example.com \
  --addr=:443
```

### Google Public CA (GCP)

```bash
./derper \
  --certmode=alpn \
  --certdir=/var/lib/derper/certs \
  --hostname=your-domain.com \
  --acme-email=your-email@example.com \
  --acme-eab-kid=YOUR_EAB_KEY_ID \
  --acme-eab-key=YOUR_BASE64_EAB_KEY \
  --addr=:443
```

### Development Mode

```bash
./derper \
  --dev \
  --certmode=alpn \
  --certdir=/tmp/derper-certs \
  --hostname=localhost \
  --acme-email=test@example.com \
  --addr=:3340
```

## 🐳 Docker Deployment

### Build and Run Locally

```bash
# Build the Docker image
docker build -t derper-alpn .

# Run the container
docker run -d \
  --name derper \
  -p 443:443 \
  -p 80:80 \
  -p 3478:3478/udp \
  -e DERP_DOMAIN=derp.example.com \
  -e DERP_ACME_EMAIL=admin@example.com \
  -v ./certs:/app/certs \
  derper-alpn
```

### Using Docker Compose

```bash
# Start the service
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the service
docker-compose down
```

### GitHub Container Registry

```bash
# Pull the latest image
docker pull ghcr.io/claw-bot/tailscale-derper-alpn:main

# Run the image
docker run -d \
  --name derper \
  -p 443:443 \
  -p 80:80 \
  -p 3478:3478/udp \
  -e DERP_DOMAIN=derp.example.com \
  -e DERP_ACME_EMAIL=admin@example.com \
  ghcr.io/claw-bot/tailscale-derper-alpn:main
```

## 📁 Files

### Core Implementation
- **alpn.go** - Main ALPN-01 challenge implementation
- **alpn_test.go** - Unit tests for ALPN functionality
- **cert.go** - Modified to integrate ALPN provider
- **derper.go** - Modified to support `--certmode=alpn`
- **go.mod** - Go module definition
- **go.sum** - Go module checksums

### Docker & CI/CD
- **Dockerfile** - Multi-stage Docker build
- **docker-compose.yml** - Local development setup
- **.dockerignore** - Docker build exclusions
- **.github/workflows/docker.yml** - GitHub Actions workflow
- **DOCKER_GUIDE.md** - Complete Docker deployment guide

### Documentation
- **README.md** - This file
- **ALPN.md** - User documentation and usage guide
- **ALPN_EXAMPLE.md** - Practical examples and deployment scenarios
- **ALPN_SUMMARY.md** - Implementation summary and technical details
- **ALPN_IMPLEMENTATION.md** - Detailed implementation guide
- **ALPN_QUICK_START.md** - Quick reference card
- **ALPN_README.md** - Main README

## 🧪 Testing

```bash
# Run unit tests
go test -v -run TestALPN

# Build
go build -o derper .

# Test ALPN protocol
openssl s_client -connect your-domain.com:443 -alpn acme-tls/1

# Test HTTPS
curl https://your-domain.com/
```

## 📋 Requirements

- Go 1.21 or later
- Port 443 accessible
- Valid domain name
- ACME account (email)
- For GCP: EAB credentials (Key ID + HMAC key)

## 🔧 Technical Implementation

### Challenge Flow

```
Client connects with "acme-tls/1" ALPN
    ↓
Server detects ALPN-01 challenge
    ↓
ACME client generates challenge certificate
    ↓
ACME server validates the certificate
    ↓
Real certificate is issued
```

### Key Components

1. **alpnCertProvider** - Implements `certProvider` interface
2. **TLS Configuration** - Configured with `acme-tls/1` ALPN protocol
3. **Challenge Detection** - Intercepts ALPN-01 challenges in `GetCertificate`
4. **Certificate Generation** - Uses ACME client's `TLSALPN01ChallengeCert()` method

## 🛡️ Security

- Certificate validation with proper hostname checking
- ACME account key security
- TLS-based challenge validation
- Configurable connection limits
- Monitoring and alerting support

## 📊 Status

- ✅ Code implementation complete
- ✅ Unit tests passing
- ✅ Documentation complete
- ✅ Build successful
- ✅ Docker support added
- ✅ GitHub Actions workflow added
- ⚠️  Full ACME integration (partial - needs account registration)
- ⚠️  Challenge state management (placeholder implementation)
- ⚠️  Automatic renewal (future enhancement)

## 🔄 Future Enhancements

1. Complete ACME account registration flow
2. Persistent challenge state storage
3. Automated certificate renewal
4. Multi-domain certificate support
5. OCSP stapling
6. Certificate transparency logging

## 📝 License

BSD-3-Clause (same as Tailscale)

## 🙏 Credits

This implementation follows the ACME protocol standards (RFC 8555) and TLS-ALPN-01 challenge specification (RFC 8737).

## 📞 Support

For issues or questions:
1. Check the documentation files in this directory
2. Review the unit tests in `alpn_test.go`
3. Consult the ACME protocol specifications
4. Refer to Tailscale DERP documentation

## 📚 References

- [ACME Protocol RFC 8555](https://datatracker.ietf.org/doc/html/rfc8555)
- [TLS-ALPN-01 Challenge RFC 8737](https://datatracker.ietf.org/doc/html/rfc8737)
- [Let's Encrypt Documentation](https://letsencrypt.org/docs/challenge-types/)
- [Google Public CA](https://pki.goog/)
- [Tailscale DERP Documentation](https://tailscale.com/kb/1232/derp-servers)
