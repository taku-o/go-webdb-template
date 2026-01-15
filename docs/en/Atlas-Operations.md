**[日本語](../ja/Atlas-Operations.md) | [English]**

# Atlas Migration Operations Guide

## Overview

This project uses [Atlas](https://atlasgo.io/) for database schema management.
Atlas is a declarative schema management tool that defines schemas in HCL files
and automatically generates and applies migrations.

## Directory Structure

```
db/
├── schema/                    # Schema definition files (HCL)
│   ├── master.hcl            # Master DB schema definition
│   ├── sharding_1/           # Shard 1 schema (tables 000-007)
│   │   ├── _schema.hcl       # Schema definition (schema "public" {})
│   │   ├── dm_users.hcl      # dm_users_000 to dm_users_007
│   │   └── dm_posts.hcl      # dm_posts_000 to dm_posts_007
│   ├── sharding_2/           # Shard 2 schema (tables 008-015)
│   │   ├── _schema.hcl
│   │   ├── dm_users.hcl      # dm_users_008 to dm_users_015
│   │   └── dm_posts.hcl      # dm_posts_008 to dm_posts_015
│   ├── sharding_3/           # Shard 3 schema (tables 016-023)
│   │   ├── _schema.hcl
│   │   ├── dm_users.hcl      # dm_users_016 to dm_users_023
│   │   └── dm_posts.hcl      # dm_posts_016 to dm_posts_023
│   └── sharding_4/           # Shard 4 schema (tables 024-031)
│       ├── _schema.hcl
│       ├── dm_users.hcl      # dm_users_024 to dm_users_031
│       └── dm_posts.hcl      # dm_posts_024 to dm_posts_031
└── migrations/               # Migration files (including initial data)
    ├── master/               # Master DB migrations
    │   ├── YYYYMMDD_*.sql
    │   └── atlas.sum
    ├── sharding_1/           # webdb_sharding_1 migrations (tables 000-007)
    │   ├── YYYYMMDD_*.sql
    │   └── atlas.sum
    ├── sharding_2/           # webdb_sharding_2 migrations (tables 008-015)
    │   ├── YYYYMMDD_*.sql
    │   └── atlas.sum
    ├── sharding_3/           # webdb_sharding_3 migrations (tables 016-023)
    │   ├── YYYYMMDD_*.sql
    │   └── atlas.sum
    ├── sharding_4/           # webdb_sharding_4 migrations (tables 024-031)
    │   ├── YYYYMMDD_*.sql
    │   └── atlas.sum
    └── view_master/          # Master DB view migrations
        ├── YYYYMMDD_*.sql
        └── atlas.sum

config/
├── develop/atlas.hcl         # Development environment Atlas config
├── staging/atlas.hcl         # Staging environment Atlas config
└── production/atlas.hcl      # Production environment Atlas config
```

### Table Partitioning Rules

The sharding DB is divided into 4 databases, each containing 8 table partitions.

| Database | Table Range | Schema Directory | Migration Directory |
|----------|-------------|------------------|---------------------|
| webdb_sharding_1 | dm_users_000-007, dm_posts_000-007 | db/schema/sharding_1/ | db/migrations/sharding_1/ |
| webdb_sharding_2 | dm_users_008-015, dm_posts_008-015 | db/schema/sharding_2/ | db/migrations/sharding_2/ |
| webdb_sharding_3 | dm_users_016-023, dm_posts_016-023 | db/schema/sharding_3/ | db/migrations/sharding_3/ |
| webdb_sharding_4 | dm_users_024-031, dm_posts_024-031 | db/schema/sharding_4/ | db/migrations/sharding_4/ |

## PostgreSQL Container Configuration

| Container Name | Database Name | Host Port |
|----------------|---------------|-----------|
| postgres-master | webdb_master | 5432 |
| postgres-sharding-1 | webdb_sharding_1 | 5433 |
| postgres-sharding-2 | webdb_sharding_2 | 5434 |
| postgres-sharding-3 | webdb_sharding_3 | 5435 |
| postgres-sharding-4 | webdb_sharding_4 | 5436 |

## Basic Commands

### Migration Script

```bash
# Apply all migrations (default)
APP_ENV=develop ./scripts/migrate.sh all

# Master database only
APP_ENV=develop ./scripts/migrate.sh master

# Sharding databases only
APP_ENV=develop ./scripts/migrate.sh sharding
```

### Generating Migrations

After modifying schema definitions, generate diff migrations.

```bash
# Generate master DB migration
atlas migrate diff <migration_name> \
    --dir file://db/migrations/master \
    --to file://db/schema/master.hcl \
    --dev-url "postgres://webdb:webdb@localhost:5432/webdb_master?sslmode=disable"

# Generate sharding DB migrations (execute for each shard)
# Shard 1
atlas migrate diff <migration_name> \
    --dir file://db/migrations/sharding_1 \
    --to file://db/schema/sharding_1 \
    --dev-url "postgres://webdb:webdb@localhost:5433/webdb_sharding_1?sslmode=disable"

# Shard 2
atlas migrate diff <migration_name> \
    --dir file://db/migrations/sharding_2 \
    --to file://db/schema/sharding_2 \
    --dev-url "postgres://webdb:webdb@localhost:5434/webdb_sharding_2?sslmode=disable"

# Shard 3
atlas migrate diff <migration_name> \
    --dir file://db/migrations/sharding_3 \
    --to file://db/schema/sharding_3 \
    --dev-url "postgres://webdb:webdb@localhost:5435/webdb_sharding_3?sslmode=disable"

# Shard 4
atlas migrate diff <migration_name> \
    --dir file://db/migrations/sharding_4 \
    --to file://db/schema/sharding_4 \
    --dev-url "postgres://webdb:webdb@localhost:5436/webdb_sharding_4?sslmode=disable"
```

### Applying Migrations

```bash
# Apply migration to master DB
atlas migrate apply \
    --dir file://db/migrations/master \
    --url "postgres://webdb:webdb@localhost:5432/webdb_master?sslmode=disable"

# Apply migrations to sharding DBs (use corresponding migration directory for each shard)
for i in 1 2 3 4; do
    port=$((5432 + i))
    atlas migrate apply \
        --dir file://db/migrations/sharding_${i} \
        --url "postgres://webdb:webdb@localhost:${port}/webdb_sharding_${i}?sslmode=disable"
done
```

### Checking Migration Status

```bash
# Master DB migration status
atlas migrate status \
    --dir file://db/migrations/master \
    --url "postgres://webdb:webdb@localhost:5432/webdb_master?sslmode=disable"

# Sharding DB migration status (check each shard)
for i in 1 2 3 4; do
    echo "=== webdb_sharding_${i} ==="
    port=$((5432 + i))
    atlas migrate status \
        --dir file://db/migrations/sharding_${i} \
        --url "postgres://webdb:webdb@localhost:${port}/webdb_sharding_${i}?sslmode=disable"
done
```

## Environment-Specific Application

### Development Environment

```bash
# Migration using config file
atlas migrate apply \
    --config file://config/develop/atlas.hcl \
    --env master

atlas migrate apply \
    --config file://config/develop/atlas.hcl \
    --env sharding_1

# Or use the simple script
APP_ENV=develop ./scripts/migrate.sh all
```

### Staging Environment

```bash
# Set DB URL via environment variables or edit config file
atlas migrate apply \
    --config file://config/staging/atlas.hcl \
    --env master

# Apply to each shard
for env in sharding_1 sharding_2 sharding_3 sharding_4; do
    atlas migrate apply \
        --config file://config/staging/atlas.hcl \
        --env $env
done
```

### Production Environment

```bash
# In production, verify with dry-run before applying
atlas migrate apply \
    --config file://config/production/atlas.hcl \
    --env master \
    --dry-run

# Apply if no issues
atlas migrate apply \
    --config file://config/production/atlas.hcl \
    --env master
```

## Case-Specific Operations

### Adding a Table

1. Add table definition to schema file (`db/schema/master.hcl` or `db/schema/sharding_*/`)

```hcl
table "new_table" {
  schema = schema.public
  column "id" {
    null = false
    type = serial
  }
  column "name" {
    null = false
    type = text
  }
  primary_key {
    columns = [column.id]
  }
}
```

2. Generate migration

```bash
atlas migrate diff add_new_table \
    --dir file://db/migrations/master \
    --to file://db/schema/master.hcl \
    --dev-url "postgres://webdb:webdb@localhost:5432/webdb_master?sslmode=disable"
```

3. Review generated SQL and apply

### Adding a Column

1. Add column to table definition in schema file

```hcl
table "existing_table" {
  # ... existing columns ...

  column "new_column" {
    null = true
    type = text
  }
}
```

2. Generate and apply migration

### Deleting a Table

1. Remove table definition from schema file
2. Generate and apply migration

**Warning**: Table deletion also deletes data. Execute carefully in production.

### Adding an Index

```hcl
table "users" {
  # ... column definitions ...

  index "idx_users_email" {
    columns = [column.email]
  }

  # Unique index
  index "idx_users_unique_email" {
    unique  = true
    columns = [column.email]
  }
}
```

## Handling Irregular Cases

### Direct SQL Schema Changes

If SQL is executed directly outside of Atlas management, schema and migration history become inconsistent.

**Solution 1: Set Baseline**

```bash
# Set current DB schema as baseline
atlas migrate hash --dir file://db/migrations/master

# Manually update atlas_schema_revisions table
docker exec -i postgres-master psql -U webdb -d webdb_master -c "DELETE FROM atlas_schema_revisions"
atlas migrate apply \
    --dir file://db/migrations/master \
    --url "postgres://webdb:webdb@localhost:5432/webdb_master?sslmode=disable" \
    --baseline <version>
```

**Solution 2: Schema Sync**

```bash
# Inspect schema from current DB
atlas schema inspect \
    --url "postgres://webdb:webdb@localhost:5432/webdb_master?sslmode=disable" \
    --format hcl > db/schema/master_current.hcl

# Check diff and update master.hcl
diff db/schema/master.hcl db/schema/master_current.hcl
```

### Migration Application Failure

```bash
# Check migration status
atlas migrate status \
    --dir file://db/migrations/master \
    --url "postgres://webdb:webdb@localhost:5432/webdb_master?sslmode=disable"

# If partially applied migration exists, fix manually
# Execute rollback SQL as needed
```

### Database Initialization (Rebuild from Scratch)

```bash
# Stop PostgreSQL containers and delete data directories
./scripts/start-postgres.sh stop
rm -rf postgres/data/master postgres/data/sharding_*

# Restart PostgreSQL containers
./scripts/start-postgres.sh start

# Apply migrations (including initial data)
APP_ENV=develop ./scripts/migrate.sh all
```

## Troubleshooting

### Connection Error

**Symptom**: `connection refused` or `no such host`

**Solution**:
```bash
# Check PostgreSQL container status
./scripts/start-postgres.sh status
./scripts/start-postgres.sh health
```

### Checksum Error

**Symptom**: `checksum mismatch`

**Solution**:
```bash
atlas migrate hash --dir "file://db/migrations/master"
atlas migrate hash --dir "file://db/migrations/sharding_1"
atlas migrate hash --dir "file://db/migrations/sharding_2"
atlas migrate hash --dir "file://db/migrations/sharding_3"
atlas migrate hash --dir "file://db/migrations/sharding_4"
```

### Atlas CLI Error

**Symptom**: `atlas: command not found`

**Solution**:
```bash
# macOS
brew install ariga/tap/atlas

# Other OS
curl -sSf https://atlasgo.sh | sh
```

## Important Notes

- Do not edit migration files once generated
- `atlas.sum` file is used for migration integrity check
- Always backup before running migrations in production
- Sharding DBs must use corresponding migration directories for each shard
- When modifying sharding schema, update all 4 schema directories (sharding_1 to 4)

## Creating Data Update Migrations

Atlas only detects schema changes, so data update migrations must be created manually.

```bash
# Create empty migration file
atlas migrate new insert_data --dir file://db/migrations/master

# Add SQL to generated file
# Example: INSERT INTO table_name (col1, col2) VALUES ('value1', 'value2');

# Update atlas.sum
atlas migrate hash --dir file://db/migrations/master

# Apply migration
atlas migrate apply \
    --dir file://db/migrations/master \
    --url "postgres://webdb:webdb@localhost:5432/webdb_master?sslmode=disable"
```

## VIEW Management

### Overview

Views are managed by creating SQL files and applying them with Atlas.

### Directory Structure

View migrations are managed in a separate directory from table migrations.

```
db/
├── schema/
│   └── master.hcl              # Table definitions only (no view definitions)
└── migrations/
    ├── master/                 # Table migrations
    │   ├── YYYYMMDD_*.sql
    │   └── atlas.sum
    └── view_master/            # Master DB view migrations
        ├── YYYYMMDD_*.sql
        └── atlas.sum
```

### View Creation Steps

#### 1. Manually Create SQL File

Create SQL file in `db/migrations/view_master/`.

```bash
# Generate empty file with specified directory and name
atlas migrate new create_view_name \
    --dir "file://db/migrations/view_master"
```

```sql
-- Create dm_news_view
CREATE VIEW dm_news_view AS SELECT id, title, content, published_at FROM dm_news;
```

#### 2. Update Checksum

```bash
atlas migrate hash --dir "file://db/migrations/view_master"
```

#### 3. Apply Migration

View migrations are automatically applied by `scripts/migrate.sh`.

### Deleting a View

```sql
-- Drop dm_news_view
DROP VIEW IF EXISTS dm_news_view;
```

### Notes

- Views depend on base tables, so table migrations must be applied first
- History table (`atlas_schema_revisions`) is shared. Managed by filename (version number)
- Do not add view definitions to HCL files as it prevents use of `atlas migrate diff`
