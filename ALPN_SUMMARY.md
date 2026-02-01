# ALPN-01 Challenge Implementation Summary

## Overview

This implementation adds TLS ALPN-01 challenge support to the Tailscale DERP server, allowing it to obtain SSL/TLS certificates using the ACME TLS-ALPN-01 challenge type.

## Changes Made

### 1. New Files

#### `alpn.go`
- Implements `alpnCertProvider` struct that satisfies the `certProvider` interface
- Handles ALPN-01 challenge detection and certificate generation
- Intercepts TLS connections with `acme-tls/1` ALPN protocol
- Uses ACME client's `TLSALPN01ChallengeCert` method to generate challenge certificates

#### `alpn_test.go`
- Unit tests for ALPN cert provider functionality
- Tests error handling and edge cases
- Validates TLS configuration

#### `ALPN.md`
- Documentation for ALPN-01 challenge usage
- Implementation details and requirements
- Troubleshooting guide

#### `ALPN_EXAMPLE.md`
- Practical examples for different scenarios
- Docker and systemd deployment examples
- Monitoring and security considerations

### 2. Modified Files

#### `derper.go`
- Updated `--certmode` flag to include `alpn` option
- Added "alpn" to the list of supported certificate modes

#### `cert.go`
- Added `NewALPNCertProviderFromFlags()` function
- Integrated ALPN provider into `certProviderByCertMode()`
- Added support for both Let's Encrypt and GCP ACME directories

## How It Works

### Challenge Flow

1. **Client Connection**: When a TLS client connects with `acme-tls/1` ALPN protocol
2. **Challenge Detection**: Server detects ALPN-01 challenge via `GetCertificate`
3. **Certificate Generation**: ACME client generates challenge certificate with key authorization
4. **Validation**: ACME server validates the certificate
5. **Certificate Issuance**: Real certificate is issued for the domain

### Key Components

```go
type alpnCertProvider struct {
    acmeClient *acme.Client
    domain     string
    keyAuth    string
    challenge  *acme.Challenge
}
```

### TLS Configuration

```go
func (p *alpnCertProvider) TLSConfig() *tls.Config {
    return &tls.Config{
        NextProtos: []string{
            "http/1.1",
            "acme-tls/1", // ALPN protocol for ACME TLS-ALPN-01 challenge
        },
        GetCertificate: p.getCertificate,
    }
}
```

## Usage

### Basic Command

```bash
./derper \
  --certmode=alpn \
  --certdir=/var/lib/derper/certs \
  --hostname=derp.example.com \
  --acme-email=admin@example.com \
  --addr=:443
```

### With GCP ACME

```bash
./derper \
  --certmode=alpn \
  --certdir=/var/lib/derper/certs \
  --hostname=derp.example.com \
  --acme-email=admin@example.com \
  --acme-eab-kid=your-eab-key-id \
  --acme-eab-key=your-base64-encoded-eab-key \
  --addr=:443
```

## Requirements

- Port 443 must be accessible (for TLS connections)
- Valid domain name pointing to the server
- ACME account credentials (email)
- For GCP: EAB credentials (Key ID and HMAC key)

## Benefits

1. **No Port 80 Required**: Unlike HTTP-01, ALPN-01 only requires port 443
2. **Better Security**: Certificate validation happens over TLS
3. **Firewall Friendly**: Only one port needs to be open
4. **ACME Standard**: Compatible with Let's Encrypt and other ACME providers

## Limitations

1. **ACME Client Dependency**: Requires proper ACME client configuration
2. **Challenge Certificate**: Uses self-signed certificates during validation
3. **Token Management**: Token generation and management needs ACME server interaction
4. **Rate Limits**: Subject to ACME provider rate limits

## Testing

Run the unit tests:

```bash
cd tailscale
go test ./cmd/derper -v -run TestALPN
```

## Future Improvements

1. **Full ACME Integration**: Complete integration with ACME account registration
2. **Challenge State Management**: Better handling of challenge state and tokens
3. **Certificate Caching**: Improved certificate caching for ALPN-01
4. **Multi-domain Support**: Support for multiple domains in a single certificate
5. **Renewal Automation**: Automatic certificate renewal handling

## Compatibility

- **Go Version**: Requires Go 1.21 or later
- **ACME Protocol**: RFC 8555 (ACME)
- **TLS-ALPN-01**: RFC 8737
- **Platforms**: Linux, macOS, Windows (with proper TLS support)

## Security Considerations

1. **Certificate Validation**: Ensure proper hostname validation
2. **Key Management**: Secure handling of ACME account keys
3. **Challenge Security**: ALPN-01 is secure but requires proper implementation
4. **Rate Limiting**: Implement appropriate connection limits
5. **Monitoring**: Monitor for certificate expiration and renewal failures

## References

- [ACME Protocol RFC 8555](https://datatracker.ietf.org/doc/html/rfc8555)
- [TLS-ALPN-01 Challenge RFC 8737](https://datatracker.ietf.org/doc/html/rfc8737)
- [Let's Encrypt Documentation](https://letsencrypt.org/docs/challenge-types/)
- [Google Public CA](https://pki.goog/)
- [Tailscale DERP Documentation](https://tailscale.com/kb/1232/derp-servers)
