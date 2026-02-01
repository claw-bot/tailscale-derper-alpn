# ALPN-01 Challenge Support

This document describes the ALPN-01 challenge support added to the DERP server.

## Overview

The ALPN-01 challenge is an ACME challenge type that validates domain ownership by serving a specific certificate over TLS with the `acme-tls/1` ALPN protocol. This is an alternative to the HTTP-01 challenge and doesn't require opening port 80.

## Usage

To use ALPN-01 challenge mode, start the DERP server with:

```bash
./derper \
  --certmode=alpn \
  --certdir=/path/to/certs \
  --hostname=your-domain.com \
  --acme-email=your-email@example.com
```

For GCP ACME directory (with EAB credentials):

```bash
./derper \
  --certmode=alpn \
  --certdir=/path/to/certs \
  --hostname=your-domain.com \
  --acme-email=your-email@example.com \
  --acme-eab-kid=your-eab-key-id \
  --acme-eab-key=your-base64-encoded-eab-key
```

## How It Works

1. **Challenge Detection**: When a TLS client connects with the `acme-tls/1` ALPN protocol, the server detects it as an ALPN-01 challenge.

2. **Certificate Generation**: The server generates a self-signed certificate that contains the key authorization in the certificate's Subject Alternative Name (SAN) extension.

3. **Validation**: The ACME server connects to your server with the `acme-tls/1` protocol and validates the certificate.

4. **Certificate Issuance**: Once validated, the ACME server issues a real certificate for your domain.

## Implementation Details

### Files Modified

- `derper.go`: Added `alpn` to the `--certmode` flag options
- `cert.go`: Added `NewALPNCertProviderFromFlags()` function
- `alpn.go`: New file implementing the ALPN-01 challenge provider

### Key Components

1. **alpnCertProvider**: Implements the `certProvider` interface for ALPN-01 challenges
2. **TLS Config**: Configured with `acme-tls/1` in the `NextProtos` list
3. **GetCertificate**: Intercepts ALPN-01 challenges and generates challenge certificates

### Certificate Format

The challenge certificate includes:
- Domain name as Common Name (CN)
- Domain name as Subject Alternative Name (DNS)
- Key authorization in an ACME-specific extension (OID: 1.3.6.1.5.5.7.1.31)

## Requirements

- Port 443 must be accessible (for TLS connections)
- A valid domain name
- ACME account credentials (email)
- For GCP: EAB credentials (Key ID and HMAC key)

## Limitations

- Currently supports only Let's Encrypt and Google Public CA (GCP)
- Requires the ACME client to be properly configured
- The challenge certificate is self-signed and short-lived (24 hours)

## Example

```bash
# Start DERP server with ALPN-01 challenge
./derper \
  --certmode=alpn \
  --certdir=/var/lib/derper/certs \
  --hostname=derp.example.com \
  --acme-email=admin@example.com \
  --addr=:443
```

## Troubleshooting

1. **Port 443 not accessible**: Ensure port 443 is open and not blocked by firewall
2. **Domain validation fails**: Verify the domain name is correct and points to your server
3. **ACME errors**: Check your ACME account credentials and rate limits
4. **Certificate generation fails**: Ensure the certdir is writable

## References

- [ACME Protocol](https://datatracker.ietf.org/doc/html/rfc8555)
- [TLS-ALPN-01 Challenge](https://datatracker.ietf.org/doc/html/rfc8737)
- [Let's Encrypt Documentation](https://letsencrypt.org/docs/challenge-types/)
