**[日本語](../ja/Queue-Job.md) | [English]**

# Job Queue Feature Guide

## Overview

This document explains the usage of the job queue feature in go-webdb-template.

The job queue feature is a delayed job processing system using Redis and the Asynq library.

## Environment Setup

### 1. Start Redis

```bash
./scripts/start-redis.sh start
```

### 2. Start Redis Insight (Optional)

If you want to check Redis status:

```bash
./scripts/start-redis-insight.sh start
```

Web UI: http://localhost:8001

### 3. Start API Server

```bash
cd server
APP_ENV=develop go run ./cmd/server/main.go
```

### 4. Start JobQueue Server

```bash
cd server
APP_ENV=develop go run ./cmd/jobqueue/main.go
```

## Registering Jobs via API

### Endpoint

**POST** `/api/dm-jobqueue/register`

### Authentication

This endpoint requires authentication.

### Request Example

```bash
curl -X POST http://localhost:8080/api/dm-jobqueue/register \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -d '{
    "message": "Hello, World!",
    "delay_seconds": 10,
    "max_retry": 3
  }'
```

### Request Body

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| message | string | No | Message to output (default: "Job executed successfully") |
| delay_seconds | int | No | Delay time in seconds (default: 180 seconds) |
| max_retry | int | No | Maximum retry count (default: 10) |

### Response Example

```json
{
  "job_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "registered"
}
```

## Verifying Job Execution

Registered jobs are output to the JobQueue server's standard output after the specified delay time.

```
[2026-01-02 12:34:56] Hello, World!
```

**Note**: The JobQueue server must be running to process jobs.

## Shutdown Procedure

### Stop API Server

Stop with Ctrl+C or `kill <PID>`.

### Stop JobQueue Server

Stop with Ctrl+C or `kill <PID>`.

### Stop Redis

```bash
./scripts/start-redis.sh stop
```

### Stop Redis Insight

```bash
./scripts/start-redis-insight.sh stop
```

## Notes

### Orphaned Asynq Server Configuration Issue

The JobQueue server has a built-in Asynq server that registers itself in Redis (`asynq:servers:{hostname:pid:uuid}`).
When the JobQueue server doesn't shut down properly, server configuration may remain in Redis.

- Cases have been observed where jobs cannot be processed correctly in such situations.

**Verification Method**:

```bash
docker exec redis redis-cli KEYS "asynq:servers:*"
```

In normal state, there should be only one entry per running JobQueue server.

**Resolution**:

If old processes remain, terminate the relevant processes:

```bash
# Check PID
ps aux | grep "go run ./cmd/jobqueue/main.go"

# Terminate process
kill <PID>
```

### Behavior on Redis Connection Error

**API Server**:
The API server will start even if it cannot connect to Redis. However, the job registration API will return a 503 error.

**JobQueue Server**:
The JobQueue server will continue to start even if it cannot connect to Redis. It will output a warning log and automatically resume processing when the Redis connection is restored.

## Related Documentation

- [API Documentation](./API.md): API endpoint details
- [Architecture](./Architecture.md): System architecture description
