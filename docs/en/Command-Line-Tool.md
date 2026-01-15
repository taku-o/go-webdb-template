**[日本語](../ja/Command-Line-Tool.md) | [English]**

# CLI Tool Documentation

## Overview

CLI tools for batch processing and cron job execution. Reuses existing service and repository layers, supporting non-interactive execution.

### Main Features

- **User List Output**: Retrieves user list from all shards and outputs in TSV format
- **Count Limit**: Controls output count with `--limit` flag
- **Cron Compatible**: Non-interactive execution with appropriate exit codes

## Directory Structure

```
server/
├── cmd/
│   ├── list-users/
│   │   ├── main.go          # CLI tool main
│   │   └── main_test.go     # Unit tests
│   └── generate-sample-data/
│       └── main.go          # Sample data generation tool
└── bin/                      # Built executables (.gitignore target)
    ├── list-users
    └── generate-sample-data
```

## Build Methods

### Development Environment

```bash
cd server
go build -o bin/list-users ./cmd/list-users
```

### Production Environment (Cross-compile)

```bash
cd server
GOOS=linux GOARCH=amd64 go build -o bin/list-users ./cmd/list-users
```

### Release Build (Optimized)

```bash
cd server
go build -ldflags="-s -w" -o bin/list-users ./cmd/list-users
```

## list-users Command

### Overview

Retrieves user list from all shards and outputs to standard output in TSV format.

### Usage

```bash
APP_ENV=<environment> ./bin/list-users [options]
```

### Options

| Option | Description | Default | Valid Range |
|--------|-------------|---------|-------------|
| `--limit` | Output count | 20 | 1-100 |

### Examples

```bash
# Default (20 records)
APP_ENV=develop ./bin/list-users

# Specify count
APP_ENV=develop ./bin/list-users --limit 50

# Maximum count (100 records)
APP_ENV=develop ./bin/list-users --limit 100

# Output to file
APP_ENV=develop ./bin/list-users --limit 100 > users.tsv
```

### Output Format

Output is in TSV (tab-separated values) format.

```
ID	Name	Email	CreatedAt	UpdatedAt
1234567890123456789	John Doe	john@example.com	2025-01-27T10:30:00Z	2025-01-27T10:30:00Z
1234567890123456790	Jane Smith	jane@example.com	2025-01-27T11:00:00Z	2025-01-27T11:00:00Z
```

| Field | Type | Description |
|-------|------|-------------|
| ID | int64 | User ID (timestamp-based) |
| Name | string | User name |
| Email | string | Email address |
| CreatedAt | RFC3339 | Creation datetime |
| UpdatedAt | RFC3339 | Update datetime |

### Exit Codes

| Code | Description |
|------|-------------|
| 0 | Normal exit |
| 1 | Error exit (config error, DB connection error, argument error, etc.) |

### Error Messages

Error messages are output to standard error.

```bash
# When limit value is invalid
$ APP_ENV=develop ./bin/list-users --limit 0
2025/01/27 10:30:00 Error: limit must be at least 1

# When limit value exceeds maximum (warning)
$ APP_ENV=develop ./bin/list-users --limit 200
2025/01/27 10:30:00 Warning: limit exceeds maximum (100), using 100
ID	Name	Email	CreatedAt	UpdatedAt
...
```

## Cron Configuration Examples

### Daily User List Backup at 3 AM

```cron
0 3 * * * cd /path/to/server && APP_ENV=production ./bin/list-users --limit 100 > /var/log/users_$(date +\%Y\%m\%d).tsv 2>> /var/log/list-users.log
```

### Environment Variable Setup

When running via cron, environment variables must be explicitly set.

```cron
APP_ENV=production
PATH=/usr/local/go/bin:/usr/bin:/bin

0 3 * * * cd /path/to/server && ./bin/list-users --limit 100 > /var/log/users.tsv 2>&1
```

## Testing

### Running Unit Tests

```bash
cd server
go test -v ./cmd/list-users/...
```

### Checking Test Coverage

```bash
cd server
go test -cover ./cmd/list-users/...
```

## Architecture

CLI tools reuse the existing layered architecture. Like the API server, they call the service layer through the usecase layer.

### list-dm-users Command

```
┌─────────────────────────────────────────────────────────────┐
│                    list-dm-users command                     │
│                    (cmd/list-dm-users/main.go)              │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│         Usecase Layer (internal/usecase/cli)                 │
│         - ListDmUsersUsecase.ListDmUsers()                  │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Service Layer (internal/service)                │
│              - DmUserService.ListDmUsers()                  │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Repository Layer (internal/repository)          │
│              - DmUserRepository.List()                      │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              DB Layer (internal/db)                          │
│              - GroupManager                                  │
└────────────────────────┬──────────────────────────────────┘
                         │
          ┌──────────────┴──────────────┐
          ▼                              ▼
    ┌─────────┐                    ┌─────────┐
    │ Shard 1 │                    │ Shard 2 │
    └─────────┘                    └─────────┘
```

### generate-sample-data Command

```
┌─────────────────────────────────────────────────────────────┐
│               generate-sample-data command                   │
│               (cmd/generate-sample-data/main.go)            │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│         Usecase Layer (internal/usecase/cli)                 │
│         - GenerateSampleUsecase.GenerateSampleData()        │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Service Layer (internal/service)                │
│              - GenerateSampleService.GenerateDmUsers()      │
│              - GenerateSampleService.GenerateDmPosts()      │
│              - GenerateSampleService.GenerateDmNews()       │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Repository Layer (internal/repository)          │
│              - DmUserRepository.InsertDmUsersBatch()        │
│              - DmPostRepository.InsertDmPostsBatch()        │
│              - DmNewsRepository.InsertDmNewsBatch()         │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              DB Layer (internal/db)                          │
│              - GroupManager                                  │
│              - TableSelector                                 │
└────────────────────────┬──────────────────────────────────┘
                         │
          ┌──────────────┼──────────────┐
          ▼              ▼              ▼
    ┌─────────┐    ┌─────────┐    ┌─────────┐
    │ Master  │    │ Shard 1 │    │ Shard 2 │  ...
    │(dm_news)│    │(dm_users│    │(dm_users│
    └─────────┘    │ dm_posts)│   │ dm_posts)│
                   └─────────┘    └─────────┘
```

### Layer Structure

| Layer | Directory | Role |
|-------|-----------|------|
| CLI Layer | cmd/list-dm-users/main.go | Entry point, validation, I/O control |
| Usecase Layer | internal/usecase/cli/ | CLI business logic coordination |
| Service Layer | internal/service/ | Domain logic, cross-shard operations |
| Repository Layer | internal/repository/ | Data access abstraction |
| DB Layer | internal/db/ | Sharding strategy, connection management |

## Related Documentation

- [Architecture.md](Architecture.md) - Architecture details
- [Sharding.md](Sharding.md) - Sharding strategy
- [Testing.md](Testing.md) - Testing strategy
