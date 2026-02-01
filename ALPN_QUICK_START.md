# ALPN-01 Quick Start

## Installation

```bash
cd tailscale
go build ./cmd/derper
```

## Basic Usage

### Let's Encrypt

```bash
./derper \
  --certmode=alpn \
  --certdir=/var/lib/derper/certs \
  --hostname=your-domain.com \
  --acme-email=your-email@example.com \
  --addr=:443
```

### Google Public CA

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

## Development

```bash
./derper \
  --dev \
  --certmode=alpn \
  --certdir=/tmp/derper-certs \
  --hostname=localhost \
  --acme-email=test@example.com \
  --addr=:3340
```

## Testing

```bash
# Run unit tests
go test ./cmd/derper -v -run TestALPN

# Test ALPN protocol
openssl s_client -connect your-domain.com:443 -alpn acme-tls/1

# Test HTTPS
curl https://your-domain.com/
```

## Files

- `alpn.go` - Main implementation
- `alpn_test.go` - Unit tests
- `ALPN.md` - Documentation
- `ALPN_EXAMPLE.md` - Examples
- `ALPN_SUMMARY.md` - Summary
- `ALPN_IMPLEMENTATION.md` - Implementation details
- `ALPN_QUICK_START.md` - This file

## Key Points

✅ **No port 80 required** - Only port 443 needed  
✅ **TLS-based validation** - More secure than HTTP-01  
✅ **ACME standard** - Works with Let's Encrypt, GCP, etc.  
✅ **Firewall friendly** - Single port exposure  

## Requirements

- Go 1.21+
- Port 443 accessible
- Valid domain name
- ACME account (email)
- For GCP: EAB credentials

## Common Commands

```bash
# Check if server is running
curl https://your-domain.com/

# Check debug endpoints
curl https://your-domain.com/debug/

# View metrics
curl https://your-domain.com/debug/vars
```

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Port 443 blocked | Check firewall rules |
| Domain not resolving | Verify DNS settings |
| ACME errors | Check account credentials |
| Certificate not issued | Review ACME server logs |

## Next Steps

1. Read `ALPN.md` for detailed documentation
2. Check `ALPN_EXAMPLE.md` for deployment examples
3. Review `ALPN_IMPLEMENTATION.md` for technical details
4. Run tests with `go test ./cmd/derper -v -run TestALPN`
