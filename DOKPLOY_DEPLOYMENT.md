# Dokploy Deployment Guide

This guide explains how to deploy the Lexicon Beneficial Ownership API to Dokploy.

## Prerequisites

- Dokploy instance set up and running
- Git repository access
- PostgreSQL database (can be provisioned in Dokploy)
- Redis instance (can be provisioned in Dokploy)

## Deployment Steps

### 1. Create a New Application in Dokploy

1. Log in to your Dokploy dashboard
2. Click "New Application"
3. Select "Docker" as the deployment type
4. Connect your Git repository

### 2. Configure Build Settings

**Dockerfile Path:** `Dockerfile`

**Build Arguments (Optional):**
```
VERSION=1.0.0
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(git rev-parse --short HEAD)
```

**Docker Build Context:** `.` (root directory)

### 3. Configure Environment Variables

Add the following environment variables in Dokploy:

```bash
# App Configuration
APP_URL=your-domain.com
APP_LISTEN_HOST=0.0.0.0
APP_LISTEN_PORT=8080

# Database Configuration
POSTGRES_HOST=postgres-host
POSTGRES_PORT=5432
POSTGRES_DB_NAME=beneficial_ownership
POSTGRES_USERNAME=your-db-user
POSTGRES_PASSWORD=your-db-password

# Redis Configuration
REDIS_HOST=redis-host
REDIS_PORT=6379
REDIS_PASSWORD=your-redis-password

# API Security
API_KEY=your-secure-api-key
SALT=your-secure-salt

# External Services
CHATBOT_BASE_URL=https://chatbot-api.example.com
CHATBOT_API_KEY=your-chatbot-api-key

# URLs
BASE_URL=https://your-domain.com
CORS_ALLOWED_ORIGINS=https://frontend.example.com,https://admin.example.com
```

### 4. Port Configuration

- **Container Port:** `8080`
- **Public Port:** `80` or `443` (with SSL)

### 5. Health Check Configuration

Dokploy will use the built-in Docker health check defined in the Dockerfile:
- **Endpoint:** `/health`
- **Interval:** 30 seconds
- **Timeout:** 3 seconds
- **Retries:** 3
- **Start Period:** 5 seconds

### 6. Resource Allocation

Recommended resources for production:
- **CPU:** 0.5-1 vCPU
- **Memory:** 256MB-512MB
- **Storage:** 1GB (for logs and temporary files)

### 7. Domain and SSL

1. Add your custom domain in Dokploy
2. Enable SSL/TLS (Let's Encrypt is automatically configured)
3. Configure DNS to point to your Dokploy instance

## Database Setup

### Using Dokploy's PostgreSQL

1. Create a PostgreSQL service in Dokploy
2. Note the connection details
3. Update environment variables with the connection info

### Using External Database

Simply update the `POSTGRES_*` environment variables with your external database credentials.

## Redis Setup

### Using Dokploy's Redis

1. Create a Redis service in Dokploy
2. Note the connection details
3. Update environment variables with the connection info

## Deployment Strategy

### Zero-Downtime Deployment

Dokploy supports rolling deployments:
1. Build new image
2. Start new container
3. Health check passes
4. Route traffic to new container
5. Stop old container

### Rollback

If deployment fails:
1. Go to Dokploy dashboard
2. Navigate to your application
3. Click "Rollback" to previous version

## Monitoring

### Logs

View application logs in Dokploy:
```bash
# Via Dokploy dashboard
Application → Logs → Select time range
```

### Metrics

Monitor in Dokploy dashboard:
- CPU usage
- Memory usage
- Network traffic
- Request rate

## Troubleshooting

### Container Won't Start

1. Check environment variables are set correctly
2. Verify database connectivity
3. Check logs for errors

### Health Check Failing

1. Ensure `/health` endpoint exists in your application
2. Verify the application is listening on port 8080
3. Check if startup time exceeds start period

### Database Connection Issues

1. Verify database credentials
2. Check network connectivity between containers
3. Ensure database is running and accessible

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Deploy to Dokploy

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Trigger Dokploy Deployment
        run: |
          curl -X POST \
            ${{ secrets.DOKPLOY_WEBHOOK_URL }} \
            -H "Authorization: Bearer ${{ secrets.DOKPLOY_TOKEN }}"
```

## Security Considerations

1. **Environment Variables:** Never commit `.env` files
2. **Secrets:** Use Dokploy's secret management
3. **Non-root User:** Container runs as `appuser` (UID 1000)
4. **Read-only Filesystem:** Enabled in production compose file
5. **Resource Limits:** Prevent container from consuming all resources

## Scaling

To scale horizontally:
1. Go to Dokploy dashboard
2. Navigate to your application
3. Adjust replica count
4. Ensure your application is stateless

## Backup Strategy

### Database Backups

Configure regular backups in Dokploy:
- Frequency: Daily
- Retention: 7-30 days
- Storage: S3 or local storage

### Application State

Since the application is stateless, no backup needed for container itself.

## Updates and Maintenance

### Update Application

1. Push changes to Git repository
2. Dokploy auto-deploys (if webhook configured)
3. Or manually trigger deployment in dashboard

### Update Base Image

1. Update `FROM alpine:3.19` in Dockerfile
2. Test locally
3. Push and deploy

## Support

For issues or questions:
- Check Dokploy documentation: https://docs.dokploy.com
- Review application logs in Dokploy dashboard
- Check GitHub issues: https://github.com/Lexicon-Open-Source/lexicon-beneficial-ownership-api
