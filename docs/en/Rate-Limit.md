**[日本語](../ja/Rate-Limit.md) | [English]**

# API Rate Limiting Feature Guide

## Overview

This document explains the usage of the API rate limiting feature in go-webdb-template.

The rate limiting feature restricts requests to API endpoints on a per-IP address basis, preventing excessive requests.

## Feature Description

### Types of Rate Limiting

This project supports the following two types of rate limiting:

1. **Per-minute limit (always enabled)**: Limits requests per minute
2. **Per-hour limit (optional)**: Limits requests per hour (only when configured)

### Storage Types

Rate limit counters use different storage depending on the environment:

- **In-Memory**: Development environment (default)
- **Redis Cluster**: Staging/production environments (advised)

## Configuration

### Configuration File

Rate limit configuration is done in the `api.rate_limit` section of `config/{env}/config.yaml`.

```yaml
api:
  rate_limit:
    enabled: true                    # Rate limiting enabled/disabled
    requests_per_minute: 60          # Requests per minute
    requests_per_hour: 1000          # Requests per hour (optional, disabled if 0)
    storage_type: "auto"             # Storage type ("auto", "redis", "memory")
```

### Storage Type Selection

- `auto`: Auto-detect based on environment
  - If Redis Cluster config (`cacheserver.yaml`) exists with `addrs` configured: Redis Cluster
  - Otherwise: In-Memory
- `redis`: Force Redis Cluster (error if config doesn't exist)
- `memory`: Force In-Memory

### Redis Cluster Configuration

When using Redis Cluster, configure in `config/{env}/cacheserver.yaml`.

```yaml
cache_server:
  redis:
    default:
      cluster:
        addrs:
          - "localhost:6379"
          - "localhost:6380"
          - "localhost:6381"
        max_retries: 2
        min_retry_backoff: 8ms
        max_retry_backoff: 512ms
        dial_timeout: 5s
        read_timeout: 3s
        pool_size: 100
        pool_timeout: 4s
```

**Note**: If `addrs` is empty or not set, In-Memory storage is used.

## Response Headers

All API responses include the following rate limit information.

### Per-minute Limit (Always Included)

| Header | Description | Example |
|--------|-------------|---------|
| `X-RateLimit-Limit` | Limit per minute | `60` |
| `X-RateLimit-Remaining` | Remaining requests | `45` |
| `X-RateLimit-Reset` | Reset time (Unix timestamp) | `1706342400` |

### Per-hour Limit (Only When `requests_per_hour` is Set)

| Header | Description | Example |
|--------|-------------|---------|
| `X-RateLimit-Hour-Limit` | Limit per hour | `1000` |
| `X-RateLimit-Hour-Remaining` | Remaining requests | `950` |
| `X-RateLimit-Hour-Reset` | Reset time (Unix timestamp) | `1706346000` |

## Rate Limit Exceeded

When the limit is exceeded, HTTP 429 status code is returned.

### Response Example

```json
{
  "code": 429,
  "message": "Too Many Requests"
}
```

### Response Headers

Rate limit information is also included in response headers when limit is exceeded:

```
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1706342460
```

## Verification

### Checking Rate Limit Headers

```bash
curl -i -H "Authorization: Bearer <YOUR_API_KEY>" http://localhost:8080/api/users
```

**Response Example**:

```
HTTP/1.1 200 OK
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 59
X-RateLimit-Reset: 1706342460
X-RateLimit-Hour-Limit: 1000
X-RateLimit-Hour-Remaining: 999
X-RateLimit-Hour-Reset: 1706346000
```

### Verifying Rate Limit Exceeded

```bash
# Send 60 consecutive requests
for i in {1..61}; do
  curl -H "Authorization: Bearer <YOUR_API_KEY>" http://localhost:8080/api/users
  echo "Request $i"
done
```

The 61st request will return HTTP 429.

## Fail-Open Approach

The rate limiting feature adopts a **fail-open approach**. This means that if an error occurs during rate limit initialization or checking, the request is allowed.

### Behavior on Error

- **Storage initialization error**: Logged, all requests allowed
- **Rate limit check error**: Logged, request allowed

This approach prevents rate limiting feature failures from causing overall application failures.

## IP Address Retrieval

Rate limiting is applied per IP address. IP address is retrieved in the following priority:

1. `X-Forwarded-For` header (via proxy)
2. `X-Real-IP` header (via proxy)
3. Request's `RemoteAddr`

**Note**: If IP address cannot be retrieved, the request is allowed.

## Environment-Specific Suggested Settings

### Development Environment

```yaml
api:
  rate_limit:
    enabled: true
    requests_per_minute: 60
    requests_per_hour: 0              # Per-hour limit disabled
    storage_type: "auto"              # In-Memory will be used
```

### Staging Environment

```yaml
api:
  rate_limit:
    enabled: true
    requests_per_minute: 60
    requests_per_hour: 1000
    storage_type: "auto"              # Redis Cluster will be used
```

### Production Environment

```yaml
api:
  rate_limit:
    enabled: true
    requests_per_minute: 60
    requests_per_hour: 1000
    storage_type: "redis"             # Force Redis Cluster
```

## Notes

1. **In-Memory Storage Limitations**: In-Memory storage resets counters on server restart. Also, counters are not shared between multiple server instances.
2. **Redis Cluster Availability**: If Redis Cluster is unavailable, fail-open approach allows all requests. Monitoring Redis Cluster availability in production is advised.
3. **IP Address Reliability**: When accessing via proxy, ensure `X-Forwarded-For` header is correctly set.
4. **Rate Limit Adjustment**: Adjust `requests_per_minute` and `requests_per_hour` according to application usage patterns.

## Troubleshooting

### Rate Limiting Not Working

1. Verify `enabled: true` is set
2. Check logs for error messages
3. When using Redis Cluster, verify `cacheserver.yaml` configuration

### Rate Limiting Too Strict

1. Increase `requests_per_minute` and `requests_per_hour` values
2. Batch requests on client side

### Redis Cluster Connection Error

1. Verify Redis Cluster is running
2. Check `addrs` configuration in `cacheserver.yaml`
3. Verify network connectivity

## Related Documentation

- [README.md](../README.md) - Rate limiting feature overview
- [Redis-Reconnection.md](Redis-Reconnection.md) - Redis connection configuration (planned)
