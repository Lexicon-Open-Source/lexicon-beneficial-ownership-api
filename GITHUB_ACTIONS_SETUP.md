# GitHub Actions Setup Complete ✅

This document summarizes the GitHub Actions CI/CD setup for the Lexicon Beneficial Ownership API.

## What Was Created

### 1. GitHub Actions Workflows

All workflows are in [.github/workflows/](.github/workflows/):

| File | Purpose |
|------|---------|
| **ci.yml** | Continuous Integration (lint, test, security) |
| **docker.yml** | Build Docker images, push to GHCR, and deploy to Dokploy |

### 2. Configuration Files

| File | Purpose |
|------|---------|
| **.golangci.yml** | Linter configuration |
| **Makefile** | Enhanced with build targets matching CI |
| **build-production.sh** | Docker build script with versioning |

### 3. Documentation

| File | Purpose |
|------|---------|
| **.github/README.md** | Quick start guide |
| **.github/WORKFLOWS.md** | Detailed workflow documentation |
| **DOKPLOY_DEPLOYMENT.md** | Dokploy deployment guide |

### 4. Code Changes

| File | Change |
|------|--------|
| **server.go** | Added `/health` endpoint for health checks |

## Quick Start

### 1. Enable Workflows

Workflows are automatically enabled when you push to GitHub. No action needed!

### 2. Add Required Secrets

Go to `Settings > Secrets and variables > Actions` and add:

```bash
# For Dokploy Deployment
DOKPLOY_API_URL=https://your-dokploy-instance.com
DOKPLOY_API_TOKEN=your-dokploy-api-token
DOKPLOY_APPLICATION_ID=your-application-id
APP_URL=https://your-app.com

# Note: GITHUB_TOKEN is automatically provided - no action needed
```

#### How to Get Dokploy Values

1. **DOKPLOY_API_URL**: Your Dokploy instance URL (e.g., `https://dokploy.yourdomain.com`)

2. **DOKPLOY_API_TOKEN**:
   - Log in to your Dokploy dashboard
   - Go to Profile Settings
   - Generate a new API token
   - Copy the token

3. **DOKPLOY_APPLICATION_ID**:
   - Use the API to get your application ID:

     ```bash
     curl -H "x-api-key: YOUR_API_TOKEN" \
       https://your-dokploy.com/api/project.all
     ```

   - Find your application in the response and copy its `applicationId`

4. **APP_URL**: The public URL where your application will be accessible

### 3. Set Workflow Permissions

Go to `Settings > Actions > General > Workflow permissions`:
- ✅ Select "Read and write permissions"
- ✅ Check "Allow GitHub Actions to create and approve pull requests"

## What Happens Automatically

### On Every Push/PR

```text
✓ Code is linted
✓ Security scan is performed
✓ Results are reported in PR
```

### On Push to Main

```text
✓ All CI checks pass
✓ Multi-platform Docker image is built
  - linux/amd64
  - linux/arm64
✓ Image pushed to GitHub Container Registry
✓ Automatic deployment to Dokploy
✓ Health check is performed
```

### On Version Tag (e.g., v1.0.0)

```text
✓ Multi-platform Docker image is built
✓ Tagged with version number (v1.0.0, v1.0, v1, latest)
✓ Pushed to GitHub Container Registry
✓ Automatic deployment to Dokploy
```

## Local Development

### Run CI Checks Locally

```bash
# Individual checks
make lint          # Lint code
make security      # Security scan

# Optional: Run tests locally
make test          # Run tests (not in CI)
```

### Docker Operations

```bash
# Build production image
make docker-build

# Run container
make docker-run

# Build and run
make docker-test
```

## Creating a Release

### Step 1: Create and Push Tag

```bash
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

### Step 2: Automatic Process

The following happens automatically:

1. **Build and Deploy Workflow** triggers
   - Builds multi-platform Docker image
   - Tags as `v1.0.0`, `v1.0`, `v1`, `latest`
   - Pushes to GitHub Container Registry
   - Deploys to Dokploy
   - Performs health check

### Step 3: Verify

1. Check GitHub Actions tab for status
2. Check container registry for new images: `ghcr.io/your-org/your-repo:v1.0.0`
3. Verify deployment with health check

## Docker Images

### Registry

Images are published to GitHub Container Registry:
```
ghcr.io/your-org/your-repo:latest
ghcr.io/your-org/your-repo:v1.0.0
```

### Pulling Images

```bash
# Latest
docker pull ghcr.io/your-org/your-repo:latest

# Specific version
docker pull ghcr.io/your-org/your-repo:v1.0.0
```

### Multi-Platform Support

Images are built for:
- linux/amd64 (Intel/AMD)
- linux/arm64 (ARM64/Apple Silicon)

## Deployment

### Automatic Deployment

Pushes to `main` branch automatically:

1. Build Docker image
2. Push to GitHub Container Registry
3. Deploy to Dokploy
4. Run health check

### Rollback

If deployment fails, Dokploy supports rollback to previous version in the Dokploy dashboard.

## Monitoring

### View Workflow Status

```
GitHub > Actions > Select workflow > View run
```

### View Logs

```
Workflow run > Select job > Expand steps
```

## Best Practices

1. **Before Pushing**
   ```bash
   make lint          # Check code quality
   make test          # Run tests locally (optional)
   ```

2. **Use Semantic Versioning**
   - `v1.0.0` - Major release
   - `v1.1.0` - Minor release (new features)
   - `v1.0.1` - Patch release (bug fixes)

3. **Write Good Commit Messages**
   ```
   feat: add user authentication
   fix: resolve database connection issue
   docs: update API documentation
   ```

4. **Review CI Results**
   - Don't merge if CI fails
   - Fix linting issues
   - Address security scan findings

## Troubleshooting

### Workflows Not Running

**Problem**: Workflows don't trigger
**Solution**: Check if Actions are enabled in repository settings

### Docker Push Fails

**Problem**: "denied: permission_denied"
**Solution**:
1. Check workflow permissions
2. Verify GITHUB_TOKEN has write access

### Tests Fail in CI

**Problem**: Tests pass locally but fail in CI
**Solution**:
1. Run `make ci` locally
2. Check for missing environment variables
3. Look for race conditions

### Deployment Fails

**Problem**: Deployment to Dokploy fails
**Solution**:
1. Verify secrets are set correctly
2. Check Dokploy webhook is active
3. Ensure health endpoint responds

## Next Steps

### 1. Customize Workflows

Edit workflow files to match your needs:
- Add more test jobs
- Change deployment strategy
- Add notification services

### 2. Add Status Badges

In your README.md:

```markdown
![CI](https://github.com/your-org/your-repo/workflows/CI/badge.svg)
![Docker](https://github.com/your-org/your-repo/workflows/Docker%20Build%20and%20Push/badge.svg)
```

### 3. Set Up Branch Protection

Require CI to pass before merging:
1. `Settings > Branches`
2. Add rule for `main`
3. Require status checks to pass

### 4. Configure Environments

Create staging and production environments:
1. `Settings > Environments`
2. Add protection rules
3. Add environment-specific secrets

## Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Workflow Guide](.github/WORKFLOWS.md)
- [Dokploy Deployment](DOKPLOY_DEPLOYMENT.md)
- [Makefile Reference](Makefile)

## Support

For issues or questions:
- Check workflow logs in Actions tab
- Review documentation in `.github/WORKFLOWS.md`
- Open an issue in the repository

## Summary

✅ **CI/CD Pipeline**: Fully automated
✅ **Docker-Only Deployment**: No binary builds
✅ **Multi-Platform Images**: linux/amd64, linux/arm64
✅ **Security**: Automated scanning
✅ **Deployment**: Automatic to Dokploy
✅ **Container Registry**: GitHub Container Registry (GHCR)

Your streamlined CI/CD pipeline is ready to use! 🚀
