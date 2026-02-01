# Docker Deployment Guide

This guide explains how to build and deploy the Tailscale DERP server with ALPN-01 support using Docker.

## Prerequisites

- Docker installed on your system
- GitHub account (for GitHub Container Registry)
- Domain name pointing to your server

## Quick Start

### 1. Build and Run Locally

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

### 2. Using Docker Compose

```bash
# Create a docker-compose.yml file (already included)
# Edit the environment variables as needed

# Start the service
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the service
docker-compose down
```

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DERP_DOMAIN` | Yes | `your-hostname.com` | Server hostname |
| `DERP_ACME_EMAIL` | Yes | - | ACME account email |
| `DERP_CERT_MODE` | No | `letsencrypt` | Cert mode: `letsencrypt`, `alpn`, `manual` |
| `DERP_CERT_DIR` | No | `/app/certs` | Certificate directory |
| `DERP_ADDR` | No | `:443` | Listening address |
| `DERP_STUN` | No | `true` | Enable STUN server |
| `DERP_STUN_PORT` | No | `3478` | STUN port |
| `DERP_HTTP_PORT` | No | `80` | HTTP port (set to -1 to disable) |
| `DERP_VERIFY_CLIENTS` | No | `false` | Verify clients |
| `DERP_VERIFY_CLIENT_URL` | No | `""` | Client verification URL |

## GitHub Container Registry

### 1. Automatic Build (Recommended)

The GitHub workflow automatically builds and pushes images on every push to `main`:

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

### 2. Manual Build and Push

```bash
# Login to GitHub Container Registry
echo $GITHUB_TOKEN | docker login ghcr.io -u $GITHUB_USERNAME --password-stdin

# Build for multiple platforms
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t ghcr.io/claw-bot/tailscale-derper-alpn:latest \
  --push .

# Build with specific tag
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t ghcr.io/claw-bot/tailscale-derper-alpn:v1.0.0 \
  --push .
```

## Production Deployment

### Systemd Service

Create `/etc/systemd/system/derper.service`:

```ini
[Unit]
Description=Tailscale DERP Server
After=docker.service
Requires=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
ExecStart=/usr/bin/docker run --rm \
  --name derper \
  -p 443:443 \
  -p 80:80 \
  -p 3478:3478/udp \
  -e DERP_DOMAIN=derp.example.com \
  -e DERP_ACME_EMAIL=admin@example.com \
  -v /opt/derper/certs:/app/certs \
  ghcr.io/claw-bot/tailscale-derper-alpn:main

ExecStop=/usr/bin/docker stop derper
ExecStopPost=/usr/bin/docker rm derper

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl enable derper
sudo systemctl start derper
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: derper
spec:
  replicas: 1
  selector:
    matchLabels:
      app: derper
  template:
    metadata:
      labels:
        app: derper
    spec:
      containers:
      - name: derper
        image: ghcr.io/claw-bot/tailscale-derper-alpn:main
        ports:
        - containerPort: 443
          protocol: TCP
        - containerPort: 80
          protocol: TCP
        - containerPort: 3478
          protocol: UDP
        env:
        - name: DERP_DOMAIN
          value: "derp.example.com"
        - name: DERP_ACME_EMAIL
          value: "admin@example.com"
        volumeMounts:
        - name: certs
          mountPath: /app/certs
      volumes:
      - name: certs
        persistentVolumeClaim:
          claimName: derper-certs
---
apiVersion: v1
kind: Service
metadata:
  name: derper-service
spec:
  type: LoadBalancer
  selector:
    app: derper
  ports:
  - name: https
    port: 443
    targetPort: 443
    protocol: TCP
  - name: http
    port: 80
    targetPort: 80
    protocol: TCP
  - name: stun
    port: 3478
    targetPort: 3478
    protocol: UDP
```

## Monitoring

### Health Check

```bash
# Check if service is running
curl https://derp.example.com/generate_204

# Check ALPN protocol support
openssl s_client -connect derp.example.com:443 -alpn acme-tls/1
```

### Logs

```bash
# View container logs
docker logs -f derper

# View logs with systemd
journalctl -u derper -f
```

### Metrics

The DERP server exposes metrics at `/debug/vars`:

```bash
curl https://derp.example.com/debug/vars
```

## Troubleshooting

### Certificate Issues

```bash
# Check certificate directory
docker exec derper ls -la /app/certs

# View certificate details
docker exec derper cat /app/certs/derp.example.com.crt
```

### Port Conflicts

Ensure ports 443, 80, and 3478 are not in use:

```bash
# Check port usage
sudo netstat -tulpn | grep -E ':(443|80|3478)'
```

### ALPN-01 Challenge Fails

1. Ensure port 443 is accessible from the internet
2. Verify DNS points to your server
3. Check ACME account email is correct
4. Review container logs for errors

## Security Best Practices

1. **Use HTTPS only**: The server should only be accessed via HTTPS
2. **Firewall rules**: Only expose necessary ports
3. **Regular updates**: Keep Docker images updated
4. **Certificate monitoring**: Set up alerts for certificate expiration
5. **Resource limits**: Set CPU and memory limits for the container

## Performance Tuning

### Resource Limits

```bash
docker run -d \
  --name derper \
  --memory=512m \
  --cpus=1.0 \
  --restart=unless-stopped \
  -p 443:443 \
  -p 80:80 \
  -p 3478:3478/udp \
  -e DERP_DOMAIN=derp.example.com \
  -e DERP_ACME_EMAIL=admin@example.com \
  ghcr.io/claw-bot/tailscale-derper-alpn:main
```

### Connection Limits

Add to environment variables:
```bash
-e DERP_ACCEPT_CONNECTION_LIMIT=1000 \
-e DERP_ACCEPT_CONNECTION_BURST=2000 \
```

## CI/CD Integration

The GitHub workflow automatically:
1. Builds Docker images for multiple platforms (amd64, arm64)
2. Pushes to GitHub Container Registry
3. Runs security scans with Trivy
4. Generates attestations for provenance

### Triggering Builds

- **Push to main**: Builds `main` tag
- **Create tag**: Builds version tags (v1.0.0, etc.)
- **Pull request**: Builds test image (not pushed)

## Advanced Configuration

### Custom ACME Directory

For GCP ACME directory:

```bash
-e DERP_CERT_MODE=alpn \
-e DERP_ACME_EMAIL=admin@example.com \
-e DERP_ACME_EAB_KID=your-eab-key-id \
-e DERP_ACME_EAB_KEY=your-base64-eab-key \
```

### Client Verification

```bash
-e DERP_VERIFY_CLIENTS=true \
-e DERP_VERIFY_CLIENT_URL=http://your-admission-controller:8080/verify \
-v /var/run/tailscale/tailscaled.sock:/var/run/tailscale/tailscaled.sock \
```

## Support

For issues or questions:
1. Check the [ALPN.md](ALPN.md) documentation
2. Review the [ALPN_EXAMPLE.md](ALPN_EXAMPLE.md) examples
3. Consult the [Tailscale DERP documentation](https://tailscale.com/kb/1232/derp-servers)
4. Open an issue on GitHub
