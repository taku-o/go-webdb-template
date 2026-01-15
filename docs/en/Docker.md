**[日本語](../ja/Docker.md) | [English]**

# Docker

This document explains the Docker environment for this project.

## Overview

The API server, Admin server, and Client server can run on Docker containers.

### Supported Environments

| Environment | Purpose | Database |
|-------------|---------|----------|
| develop | Development | PostgreSQL/MySQL |

### Architecture

```
┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│  Client      │    │  API Server  │    │ Admin Server │
│  (Port 3000) │───▶│  (Port 8080) │    │ (Port 8081)  │
└──────────────┘    └───────┬──────┘    └───────┬──────┘
                            │                    │
                    ┌───────┴────────────────────┘
                    ▼
            ┌──────────────┐    ┌──────────────┐
            │  PostgreSQL  │    │    Redis     │
            │  (postgres)  │    │   (redis)    │
            └──────────────┘    └──────────────┘
```

---

## Prerequisites

- Docker Desktop (macOS/Windows) or Docker Engine (Linux)
- Docker Compose v2 or higher
- Following ports must be available:
  - 8080 (API server)
  - 8081 (Admin server)
  - 3000 (Client server)
  - 5432 (PostgreSQL)
  - 6379 (Redis)

---

## Dockerfile Configuration

### Server Side (Go)

| File | Purpose | CGO | Base Image |
|------|---------|-----|------------|
| `server/Dockerfile` | All environments (develop/staging/production) | 0 | golang:1.24-alpine → alpine:latest |
| `server/Dockerfile.admin` | All environments (develop/staging/production) | 0 | golang:1.24-alpine → alpine:latest |

### Client Side (Next.js)

| File | Target | Purpose |
|------|--------|---------|
| `client/Dockerfile` | dev | Development (hot reload enabled) |
| `client/Dockerfile` | production | staging/production |

---

## Docker Compose Configuration Files

### File List

| File | Services | Environment |
|------|----------|-------------|
| `docker-compose.api.yml` | API Server | develop |
| `docker-compose.admin.yml` | Admin Server | develop |
| `docker-compose.client.yml` | Client | develop |

### Volume Mounts

| Path | Purpose | Mode |
|------|---------|------|
| `./config/{env}:/app/config/{env}` | Configuration files | Read-only |
| `./server/data:/app/server/data` | Data directory | Read-write |
| `./logs:/app/logs` | Log files | Read-write |

### Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `APP_ENV` | Environment specification | develop/staging/production |
| `REDIS_JOBQUEUE_ADDR` | Redis connection | redis:6379 |
| `NEXT_PUBLIC_API_URL` | API server URL | http://api:8080 |
| `NEXT_PUBLIC_API_KEY` | API key | (loaded from .env.local) |

---

## Build, Start, and Stop Commands

### Development Environment (develop)

```bash
# Build
docker-compose -f docker-compose.api.yml build
docker-compose -f docker-compose.admin.yml build
docker-compose -f docker-compose.client.yml build

# Start
docker-compose -f docker-compose.api.yml up -d
docker-compose -f docker-compose.admin.yml up -d
docker-compose -f docker-compose.client.yml up -d

# Stop
docker-compose -f docker-compose.api.yml down
docker-compose -f docker-compose.admin.yml down
docker-compose -f docker-compose.client.yml down

# View logs
docker-compose -f docker-compose.api.yml logs -f
docker-compose -f docker-compose.admin.yml logs -f
docker-compose -f docker-compose.client.yml logs -f
```

---

## Integration with Existing Services

### Required External Networks

Dockerized servers communicate with existing PostgreSQL and Redis containers on the same network.

```bash
# Start PostgreSQL container (creates postgres-network)
docker-compose -f docker-compose.postgres.yml up -d

# Start Redis container (creates redis-network)
docker-compose -f docker-compose.redis.yml up -d
```

### Startup Order

1. Start PostgreSQL and Redis containers
2. Start API server
3. Start Admin server
4. Start Client server

### Inter-service Communication

| Source | Destination | Hostname |
|--------|-------------|----------|
| API Server | PostgreSQL | postgres:5432 |
| API Server | Redis | redis:6379 |
| Admin Server | PostgreSQL | postgres:5432 |
| Client | API Server | api:8080 |

---

## PostgreSQL Container Management

### Configuration Overview

This project uses 1 master + 4 sharding PostgreSQL containers.

| Container Name | Database Name | Host Port |
|----------------|---------------|-----------|
| postgres-master | webdb_master | 5432 |
| postgres-sharding-1 | webdb_sharding_1 | 5433 |
| postgres-sharding-2 | webdb_sharding_2 | 5434 |
| postgres-sharding-3 | webdb_sharding_3 | 5435 |
| postgres-sharding-4 | webdb_sharding_4 | 5436 |

### Start/Stop Commands

```bash
# Start
./scripts/start-postgres.sh start

# Stop
./scripts/start-postgres.sh stop

# Status check
./scripts/start-postgres.sh status

# Health check
./scripts/start-postgres.sh health
```

### Migration

```bash
# Apply migrations after starting PostgreSQL
APP_ENV=develop ./scripts/migrate.sh all
```

See [Atlas-Operations.md](./Atlas-Operations.md) for migration details.

---

## Troubleshooting

### Network Connection Error

**Symptom**: `network postgres-network not found`

**Cause**: External network has not been created.

**Solution**: Start PostgreSQL/Redis containers first.

```bash
docker-compose -f docker-compose.postgres.yml up -d
docker-compose -f docker-compose.redis.yml up -d
```

### Port Conflict

**Symptom**: `Bind for 0.0.0.0:8080 failed: port is already allocated`

**Cause**: The specified port is already in use.

**Solutions**:
1. Stop the process using the port
2. Or change the port mapping in docker-compose.yml

```bash
# Check which process is using the port
lsof -i :8080

# Stop the process
kill <PID>
```

### Image Rebuild

To rebuild without cache:

```bash
docker-compose -f docker-compose.api.yml build --no-cache
```

### Remove All Containers and Images

```bash
# Stop and remove all containers
docker stop $(docker ps -aq)
docker rm $(docker ps -aq)

# Remove unused images
docker image prune -a -f

# Remove build cache
docker builder prune -a -f
```

---

## Production Deployment

### Deployment Flow

1. Build production image
2. Tag the image
3. Push to container registry
4. Pull and start image in production environment

---

## Pushing to Container Registry

This section explains how to push Docker images to container registries for production deployment.

### Building Images

Build production images before pushing.

```bash
# API server
docker-compose -f docker-compose.api.yml build

# Admin server
docker-compose -f docker-compose.admin.yml build

# Client server
docker-compose -f docker-compose.client.yml build
```

### Tagging Images

Add tags for pushing to registry.

```bash
# API server
docker tag go-webdb-template-api:latest <registry>/api:<version>

# Admin server
docker tag go-webdb-template-admin:latest <registry>/admin:<version>

# Client server
docker tag go-webdb-template-client:latest <registry>/client:<version>
```

**Tag format examples**:
- `api:v1.0.0` - Version specified
- `api:latest` - Latest version
- `api:staging` - For staging environment

---

### AWS ECR (Elastic Container Registry)

#### Authentication

```bash
# Login to ECR with AWS CLI
aws ecr get-login-password --region <region> | docker login --username AWS --password-stdin <account-id>.dkr.ecr.<region>.amazonaws.com
```

**Parameters**:
- `<region>`: AWS region (e.g., `ap-northeast-1`)
- `<account-id>`: AWS account ID (12-digit number)

#### Create Repository (First time only)

```bash
# Create API server repository
aws ecr create-repository --repository-name api --region <region>

# Create Admin server repository
aws ecr create-repository --repository-name admin --region <region>

# Create Client server repository
aws ecr create-repository --repository-name client --region <region>
```

#### Tag and Push

```bash
# Tag
docker tag go-webdb-template-api:latest <account-id>.dkr.ecr.<region>.amazonaws.com/api:v1.0.0
docker tag go-webdb-template-admin:latest <account-id>.dkr.ecr.<region>.amazonaws.com/admin:v1.0.0
docker tag go-webdb-template-client:latest <account-id>.dkr.ecr.<region>.amazonaws.com/client:v1.0.0

# Push
docker push <account-id>.dkr.ecr.<region>.amazonaws.com/api:v1.0.0
docker push <account-id>.dkr.ecr.<region>.amazonaws.com/admin:v1.0.0
docker push <account-id>.dkr.ecr.<region>.amazonaws.com/client:v1.0.0
```

---

### Tencent Cloud TCR (Tencent Container Registry)

#### Authentication

```bash
# Login to TCR
docker login <registry-name>.tencentcloudcr.com --username <username>
```

**Parameters**:
- `<registry-name>`: TCR instance name
- `<username>`: Tencent Cloud account username or access key ID

Use the access key secret obtained from Tencent Cloud console as password.

#### Create Namespace (First time only)

Create namespace from Tencent Cloud console or use CLI.

#### Tag and Push

```bash
# Tag
docker tag go-webdb-template-api:latest <registry-name>.tencentcloudcr.com/<namespace>/api:v1.0.0
docker tag go-webdb-template-admin:latest <registry-name>.tencentcloudcr.com/<namespace>/admin:v1.0.0
docker tag go-webdb-template-client:latest <registry-name>.tencentcloudcr.com/<namespace>/client:v1.0.0

# Push
docker push <registry-name>.tencentcloudcr.com/<namespace>/api:v1.0.0
docker push <registry-name>.tencentcloudcr.com/<namespace>/admin:v1.0.0
docker push <registry-name>.tencentcloudcr.com/<namespace>/client:v1.0.0
```

---

### Docker Hub (Optional)

#### Authentication

```bash
# Login to Docker Hub
docker login --username <username>
```

Enter password or access token.

#### Tag and Push

```bash
# Tag
docker tag go-webdb-template-api:latest <username>/api:v1.0.0
docker tag go-webdb-template-admin:latest <username>/admin:v1.0.0
docker tag go-webdb-template-client:latest <username>/client:v1.0.0

# Push
docker push <username>/api:v1.0.0
docker push <username>/admin:v1.0.0
docker push <username>/client:v1.0.0
```

---

### Security Considerations

1. **Credential Management**
   - Manage credentials as secret variables in CI/CD environment
   - Do not commit credentials to source code
   - Automate authentication using IAM roles or service accounts

2. **Image Scanning**
   - Security scanning before push is advised
   - Use vulnerability scanning features of AWS ECR or TCR

3. **Tag Management**
   - Use specific version tags for production (avoid `latest`)
   - Semantic versioning (e.g., `v1.2.3`) is advised
