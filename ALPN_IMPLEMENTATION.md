# ALPN-01 Challenge Implementation for Tailscale DERP Server

## Summary

This implementation adds TLS ALPN-01 challenge support to the Tailscale DERP server, enabling it to obtain SSL/TLS certificates using the ACME TLS-ALPN-01 challenge type. This is an alternative to the HTTP-01 challenge that doesn't require opening port 80.

## Files Modified

### 1. `derper.go`
**Change**: Updated the `--certmode` flag to include `alpn` as a valid option.

```go
certMode = flag.String("certmode", "letsencrypt", "mode for getting a cert. possible options: manual, letsencrypt, gcp, alpn")
```

### 2. `cert.go`
**Changes**:
- Added `case "alpn":` to `certProviderByCertMode()` function
- Added `NewALPNCertProviderFromFlags()` function to create ALPN cert provider from command-line flags

## Files Added

### 1. `alpn.go`
**Purpose**: Implements the ALPN-01 challenge provider.

**Key Components**:
- `alpnCertProvider` struct: Implements the `certProvider` interface
- `TLSConfig()`: Returns TLS config with `acme-tls/1` ALPN protocol
- `getCertificate()`: Intercepts ALPN-01 challenges and generates certificates
- `createChallengeCert()`: Uses ACME client to create challenge certificates

**Key Features**:
- Detects ALPN-01 challenges by checking for `acme-tls/1` in client's supported protocols
- Uses ACME client's `TLSALPN01ChallengeCert()` method for certificate generation
- Returns `nil` for non-ALPN connections to allow normal TLS handshake

### 2. `alpn_test.go`
**Purpose**: Unit tests for ALPN functionality.

**Test Coverage**:
- ALPN cert provider creation
- TLS configuration validation
- Error handling for invalid inputs
- GetCertificate behavior for ALPN and non-ALPN connections

### 3. Documentation Files
- `ALPN.md`: User documentation and usage guide
- `ALPN_EXAMPLE.md`: Practical examples and deployment scenarios
- `ALPN_SUMMARY.md`: Implementation summary and technical details
- `ALPN_IMPLEMENTATION.md`: This file - detailed implementation guide

## How It Works

### Challenge Flow

```
1. Client connects with "acme-tls/1" ALPN protocol
   ↓
2. Server detects ALPN-01 challenge in GetCertificate
   ↓
3. ACME client generates challenge certificate
   ↓
4. ACME server validates the certificate
   ↓
5. Real certificate is issued for the domain
```

### Code Flow

```
derper.go (main)
    ↓
cert.go (certProviderByCertMode)
    ↓
cert.go (NewALPNCertProviderFromFlags)
    ↓
alpn.go (NewALPNCertProvider)
    ↓
alpn.go (TLSConfig)
    ↓
alpn.go (getCertificate) ← Called during TLS handshake
```

## Usage Examples

### Basic Let's Encrypt

```bash
./derper \
  --certmode=alpn \
  --certdir=/var/lib/derper/certs \
  --hostname=derp.example.com \
  --acme-email=admin@example.com \
  --addr=:443
```

### Google Public CA (GCP)

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

## Technical Details

### TLS Configuration

The ALPN provider configures the TLS server with:
- `NextProtos`: `["http/1.1", "acme-tls/1"]`
- `GetCertificate`: Custom handler for ALPN-01 challenges

### Certificate Generation

When an ALPN-01 challenge is detected:
1. The server calls `acmeClient.TLSALPN01ChallengeCert(token, domain)`
2. This generates a certificate with the proper key authorization
3. The certificate is returned to the ACME server for validation

### Connection Handling

- **ALPN connections** (`acme-tls/1`): Challenge certificates are served
- **Regular connections** (`http/1.1`): `nil` is returned, allowing default certificate handling

## Requirements

### System Requirements
- Go 1.21 or later
- Port 443 accessible (for TLS)
- Domain name pointing to server

### ACME Requirements
- Valid ACME account (email)
- For Let's Encrypt: Standard account
- For GCP: EAB credentials (Key ID + HMAC key)

## Testing

### Run Unit Tests

```bash
cd tailscale
go test ./cmd/derper -v -run TestALPN
```

### Manual Testing

```bash
# Test ALPN protocol support
openssl s_client -connect derp.example.com:443 -servername derp.example.com -alpn acme-tls/1

# Test regular HTTPS
curl https://derp.example.com/
```

## Security Considerations

1. **Certificate Validation**: Proper hostname validation is implemented
2. **Key Management**: ACME account keys should be stored securely
3. **Challenge Security**: ALPN-01 is secure by design (TLS-based)
4. **Rate Limiting**: Connection limits can be configured
5. **Monitoring**: Monitor certificate expiration and renewal

## Limitations

1. **ACME Integration**: Currently uses placeholder token generation
2. **External Account Binding**: EAB support is partial (directory URL only)
3. **Challenge State**: No persistent challenge state storage
4. **Renewal**: Automatic renewal not fully implemented

## Future Enhancements

1. **Full ACME Integration**: Complete account registration flow
2. **Challenge State Management**: Persistent storage for challenge tokens
3. **Certificate Renewal**: Automated renewal before expiration
4. **Multi-domain Support**: SAN certificates for multiple domains
5. **OCSP Stapling**: Add OCSP stapling support
6. **Certificate Transparency**: Log certificates to CT logs

## Compatibility

- **ACME Protocol**: RFC 8555
- **TLS-ALPN-01**: RFC 8737
- **Platforms**: Linux, macOS, Windows
- **ACME Providers**: Let's Encrypt, Google Public CA, others

## Troubleshooting

### Common Issues

1. **Port 443 not accessible**
   - Check firewall rules
   - Verify port is not blocked by ISP

2. **Domain validation fails**
   - Ensure domain points to correct IP
   - Check DNS propagation

3. **ACME errors**
   - Verify account credentials
   - Check rate limits
   - Review ACME server logs

4. **Certificate generation fails**
   - Ensure certdir is writable
   - Check ACME client configuration

## References

- [ACME Protocol RFC 8555](https://datatracker.ietf.org/doc/html/rfc8555)
- [TLS-ALPN-01 Challenge RFC 8737](https://datatracker.ietf.org/doc/html/rfc8737)
- [Let's Encrypt Documentation](https://letsencrypt.org/docs/challenge-types/)
- [Google Public CA](https://pki.goog/)
- [Tailscale DERP Documentation](https://tailscale.com/kb/1232/derp-servers)

## License

This implementation follows the same license as Tailscale: BSD-3-Clause
