**[日本語](../ja/Generate-Sample-Data.md) | [English]**

# Sample Data Generation Feature

## Overview

A CLI tool for generating large amounts of sample data for development. Uses the Gofakeit library to generate realistic random data.

Uses PostgreSQL to insert data into master/sharding databases.

## Prerequisites

The following preparations are required before running the command.

### 1. Start PostgreSQL Containers

```bash
./scripts/start-postgres.sh start
```

**Connection Information** (Development environment):

| Database | Host | Port | User | Password | Database Name |
|----------|------|------|------|----------|---------------|
| Master | localhost | 5432 | webdb | webdb | webdb_master |
| Sharding 1 | localhost | 5433 | webdb | webdb | webdb_sharding_1 |
| Sharding 2 | localhost | 5434 | webdb | webdb | webdb_sharding_2 |
| Sharding 3 | localhost | 5435 | webdb | webdb | webdb_sharding_3 |
| Sharding 4 | localhost | 5436 | webdb | webdb | webdb_sharding_4 |

### 2. Apply Migrations

```bash
./scripts/migrate.sh
```

## Build Method

```bash
cd server
go build -o bin/generate-sample-data ./cmd/generate-sample-data
```

## Execution Method

```bash
# Run with go run
cd server
APP_ENV=develop go run cmd/generate-sample-data/main.go

# Run with built binary
APP_ENV=develop ./bin/generate-sample-data
```

## Generated Data

### dm_users Table

- **Target**: Sharding databases (distributed across 4 servers)
- **Tables**: 32 partitioned tables (dm_users_000 to dm_users_031)
- **Count**: 100 total (distributed to each table based on UUID)
- **Generated Fields**:
  - `id`: UUIDv7 (lowercase without hyphens, 32 characters)
  - `name`: Random name
  - `email`: Random email address
  - `created_at`, `updated_at`: Current time

### dm_posts Table

- **Target**: Sharding databases (distributed across 4 servers)
- **Tables**: 32 partitioned tables (dm_posts_000 to dm_posts_031)
- **Count**: 100 total (distributed to each table based on user_id)
- **Generated Fields**:
  - `id`: UUIDv7 (lowercase without hyphens, 32 characters)
  - `user_id`: Randomly selected from existing dm_users table
  - `title`: Random sentence of about 5 words
  - `content`: Random paragraph of 3-5 sentences, each about 10 words
  - `created_at`, `updated_at`: Current time

### dm_news Table

- **Target**: Master database (webdb_master)
- **Table**: dm_news (fixed table name)
- **Count**: 100 records
- **Generated Fields**:
  - `title`: Random sentence of about 5 words
  - `content`: Random paragraph of 3-5 sentences, each about 10 words
  - `author_id`: Random 32-bit integer
  - `published_at`: Random datetime
  - `created_at`, `updated_at`: Current time

## Execution Example

```
$ APP_ENV=develop go run cmd/generate-sample-data/main.go
2026/01/09 23:39:14 Starting sample data generation...
2026/01/09 23:39:14 Generated 2 dm_users in dm_users_001
2026/01/09 23:39:14 Generated 3 dm_users in dm_users_023
...
2026/01/09 23:39:14 Generated 5 dm_posts in dm_posts_010
2026/01/09 23:39:14 Generated 1 dm_posts in dm_posts_022
...
2026/01/09 23:39:14 Generated 100 dm_news articles
2026/01/09 23:39:14 Sample data generation completed successfully
```

## Data Verification

### Verify with CloudBeaver

```bash
./scripts/start-cloudbeaver.sh
```

Access http://localhost:8978 in your browser to verify the generated data.

### Verify with Command Line

```bash
# Check dm_news count
docker exec postgres-master psql -U webdb -d webdb_master -c "SELECT COUNT(*) FROM dm_news;"

# Check dm_users count (sharding-1 table)
docker exec postgres-sharding-1 psql -U webdb -d webdb_sharding_1 -c "SELECT COUNT(*) FROM dm_users_000;"
```

## Troubleshooting

### Connection Error

```
Failed to create group manager: ...
```

**Cause**: PostgreSQL containers are not running

**Solution**:
```bash
./scripts/start-postgres.sh start
```

### Table Does Not Exist Error

```
relation "dm_users_000" does not exist
```

**Cause**: Migrations have not been applied

**Solution**:
```bash
./scripts/migrate.sh
```

## Notes

- Designed for use in develop environment
- Does not delete existing data (adds only)
- Running multiple times adds more data

## Architecture

The generate-sample-data command uses the same layered architecture as the API server. It calls the service layer through the usecase layer.

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
| CLI Layer | cmd/generate-sample-data/main.go | Entry point, I/O control |
| Usecase Layer | internal/usecase/cli/ | CLI business logic coordination |
| Service Layer | internal/service/ | Data generation logic, gofakeit usage |
| Repository Layer | internal/repository/ | Batch insertion, data access abstraction |
| DB Layer | internal/db/ | Sharding strategy, connection management |

## Technical Specifications

- **Batch Size**: 500 records at a time
- **Sharding**: Distributed to 32 partitioned tables based on UUID
- **Library**: `github.com/brianvoe/gofakeit/v6`
- **Database**: PostgreSQL (1 master + 4 sharding)

## Related Documentation

- [Command-Line-Tool.md](./Command-Line-Tool.md) - Existing CLI tool documentation
- [Sharding.md](./Sharding.md) - Sharding strategy details
- [Database-Viewer.md](./Database-Viewer.md) - How to use CloudBeaver
