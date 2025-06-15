# Osprey Web Compiler Deployment Guide

## Security Features

The Osprey Web Compiler runs in sandbox mode for security:
- **Sandbox Mode**: All code execution uses `--sandbox` flag
- **HTTP Disabled**: No HTTP requests allowed
- **WebSocket Disabled**: No WebSocket connections allowed  
- **File System Disabled**: No file read/write operations
- **FFI Disabled**: No foreign function interface access

This ensures user code cannot perform potentially dangerous operations.

## AWS ECR + ECS Deployment (Production)

### Prerequisites
- AWS CLI installed and configured (`aws configure`)
- Docker installed
- AWS account with ECR and ECS permissions

### Quick AWS Deployment

1. **Clone the repository**:
   ```bash
   git clone <your-repo-url>
   cd <repo-name>/webcompiler
   ```

2. **Deploy to AWS**:
   ```bash
   ./deploy.sh --aws
   ```
   
   This will:
   - Create ECR repository if needed
   - Build the Docker image
   - Push to ECR
   - Create/update ECS cluster and task definition
   - Deploy or update the ECS service

3. **ECR only deployment**:
   ```bash
   ./deploy.sh --ecr-only
   ```
   
   For first-time ECS deployment, you'll need to manually create the service with your VPC details.

### Configuration

Set environment variables to customize deployment:

```bash
export AWS_REGION=us-west-2
export ECR_REPOSITORY=my-osprey-compiler
export IMAGE_TAG=v1.0.0
export ECS_CLUSTER=my-cluster
export ECS_SERVICE=my-service

./deploy.sh --aws
```

### Quick Commands

```bash
# Local deployment (default)
./deploy.sh

# AWS deployment
./deploy.sh --aws

# ECR only (no ECS)
./deploy.sh --ecr-only

# Skip build (quick redeploy)
./deploy.sh --aws --skip-build

# Help
./deploy.sh --help
```

### Manual ECR Operations

```bash
# List images
aws ecr list-images --repository-name osprey-web-compiler

# Pull image
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-east-1.amazonaws.com
docker pull <account-id>.dkr.ecr.us-east-1.amazonaws.com/osprey-web-compiler:latest

# Delete old images
aws ecr batch-delete-image --repository-name osprey-web-compiler --image-ids imageTag=old-tag
```

## Local Docker Deployment (Development)

### Quick Start

1. **Clone the repository** on your machine:
   ```bash
   git clone <your-repo-url>
   cd <repo-name>/webcompiler
   ```

2. **Run the deployment script**:
   ```bash
   ./deploy.sh
   ```

3. **Access the service**:
   - Local: `http://localhost:3001`
   - External: `http://your-host-ip:3001`

### Manual Local Deployment

```bash
# Build and start
docker-compose up --build -d

# View logs
docker-compose logs -f

# Stop service
docker-compose down

# Restart service
docker-compose restart
```

## Resource Requirements

### AWS ECS (Fargate)
- **CPU**: 256 CPU units (0.25 vCPU)
- **Memory**: 512 MB
- **Storage**: Ephemeral, ~200MB for runtime
- **Network**: Port 3001 exposed via ALB/NLB

### Local Docker
- **Memory**: ~256MB reserved, 512MB limit
- **CPU**: Shared, scales with load
- **Storage**: ~100MB for image
- **Network**: Port 3001 exposed

## Troubleshooting

### ECR Issues

```bash
# Check AWS credentials
aws sts get-caller-identity

# Check ECR repository
aws ecr describe-repositories --repository-names osprey-web-compiler

# Manual ECR login
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-east-1.amazonaws.com
```

### ECS Issues

```bash
# Check service status
aws ecs describe-services --cluster osprey-cluster --services osprey-web-compiler

# View logs
aws logs tail /ecs/osprey-web-compiler --follow

# Check task definition
aws ecs describe-task-definition --task-definition osprey-web-compiler
```

### Local Docker Issues

```bash
# Container won't start
docker-compose logs

# Out of memory
docker stats

# Port conflicts - change in docker-compose.yml
ports:
  - "8080:3001"  # Use port 8080 instead
```

## Production Considerations

### AWS Deployment
1. **Load Balancer**: Use ALB for HTTPS and domain routing
2. **Auto Scaling**: Configure ECS service auto scaling
3. **Monitoring**: CloudWatch logs and metrics
4. **Security**: VPC, security groups, IAM roles
5. **CI/CD**: Integrate with CodePipeline/GitHub Actions

### Local Deployment  
1. **Reverse Proxy**: Use nginx for SSL and domain routing
2. **Firewall**: Only expose necessary ports
3. **Monitoring**: Set up log aggregation
4. **Backups**: Consider volume mounts for persistent data
5. **Updates**: Automate rebuilds for security patches

## Security

- **Osprey Sandbox Mode**: All code execution runs with `--sandbox` flag
- **Restricted Operations**: HTTP, WebSocket, file system, and FFI access disabled
- Container runs as non-root user (UID 1001)
- Multi-stage build minimizes attack surface
- Health checks for reliability
- Resource limits prevent DoS
- ECR image scanning enabled
- CloudWatch logging for audit trails 