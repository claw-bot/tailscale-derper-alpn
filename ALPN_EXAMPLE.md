# ALPN-01 Challenge Example

This document provides examples of using the ALPN-01 challenge with the DERP server.

## Basic Usage

### Let's Encrypt (Standard ACME)

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

## Development Mode

For testing ALPN-01 in development:

```bash
./derper \
  --dev \
  --certmode=alpn \
  --certdir=/tmp/derper-certs \
  --hostname=localhost \
  --acme-email=test@example.com \
  --addr=:3340
```

## Docker Example

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build ./cmd/derper

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/derper /usr/local/bin/derper
RUN mkdir -p /var/lib/derper/certs

EXPOSE 443/tcp
EXPOSE 3478/udp

CMD ["derper", \
     "--certmode=alpn", \
     "--certdir=/var/lib/derper/certs", \
     "--hostname=derp.example.com", \
     "--acme-email=admin@example.com", \
     "--addr=:443"]
```

## Systemd Service

Create `/etc/systemd/system/derper.service`:

```ini
[Unit]
Description=Tailscale DERP Server
After=network.target

[Service]
Type=simple
User=derper
Group=derper
ExecStart=/usr/local/bin/derper \
  --certmode=alpn \
  --certdir=/var/lib/derper/certs \
  --hostname=derp.example.com \
  --acme-email=admin@example.com \
  --addr=:443
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl enable derper
sudo systemctl start derper
```

## Firewall Configuration

Ensure port 443 is open:

```bash
# UFW (Ubuntu/Debian)
sudo ufw allow 443/tcp

# firewalld (CentOS/RHEL)
sudo firewall-cmd --permanent --add-port=443/tcp
sudo firewall-cmd --reload

# iptables
sudo iptables -A INPUT -p tcp --dport 443 -j ACCEPT
```

## Verification

### Check if ALPN is working

```bash
# Test with openssl
openssl s_client -connect derp.example.com:443 -servername derp.example.com -alpn acme-tls/1

# Test with curl (if available)
curl -v --alpn acme-tls/1 https://derp.example.com
```

### Check DERP server status

```bash
# Check if server is running
curl https://derp.example.com/

# Check debug endpoints
curl https://derp.example.com/debug/
```

## Troubleshooting

### Certificate not being issued

1. Check ACME account email
2. Verify domain points to your server
3. Check firewall rules
4. Review ACME server logs

### ALPN challenge fails

1. Ensure port 443 is accessible
2. Check that the server is listening on port 443
3. Verify the hostname is correct
4. Check for rate limiting

### GCP ACME issues

1. Verify EAB credentials are correct
2. Ensure the EAB key is properly base64-encoded
3. Check GCP ACME directory URL

## Performance Tuning

### Connection limits

```bash
./derper \
  --certmode=alpn \
  --certdir=/var/lib/derper/certs \
  --hostname=derp.example.com \
  --acme-email=admin@example.com \
  --addr=:443 \
  --accept-connection-limit=1000 \
  --accept-connection-burst=2000
```

### TCP settings

```bash
./derper \
  --certmode=alpn \
  --certdir=/var/lib/derper/certs \
  --hostname=derp.example.com \
  --acme-email=admin@example.com \
  --addr=:443 \
  --tcp-keepalive-time=5m \
  --tcp-user-timeout=30s
```

## Monitoring

### Prometheus metrics

The DERP server exposes metrics at `/debug/vars`:

```bash
curl https://derp.example.com/debug/vars
```

### Key metrics to monitor

- `tls_listener_counter_accepted_connections`
- `tls_listener_counter_rejected_connections`
- `gauge_derper_tls_active_version`
- `derper_tls_request_version`

## Security Considerations

1. **Use strong ACME credentials**: Keep your ACME account credentials secure
2. **Regular updates**: Keep the DERP server updated
3. **Firewall rules**: Only expose necessary ports
4. **Monitoring**: Set up alerts for certificate expiration
5. **Rate limiting**: Configure appropriate connection limits

## References

- [ACME Protocol RFC](https://datatracker.ietf.org/doc/html/rfc8555)
- [TLS-ALPN-01 Challenge RFC](https://datatracker.ietf.org/doc/html/rfc8737)
- [Let's Encrypt Documentation](https://letsencrypt.org/docs/)
- [Google Public CA Documentation](https://pki.goog/)
