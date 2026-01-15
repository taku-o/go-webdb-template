**[日本語](../ja/Metabase.md) | [English]**

# Metabase - Data Visualization & Analysis Tool

## Overview

Metabase is a data viewer and analysis tool for non-engineers. It connects to databases via web browser, enabling query creation/execution and dashboard creation/sharing.

### Management Tool Roles

This project provides three management tools:

| Tool | Role | Port |
|------|------|------|
| GoAdmin | Admin panel for custom processing | 8081 |
| CloudBeaver | Web-based tool for data operations | 8978 |
| Metabase | Data visualization & analysis tool | 8970 |

**Important**: CloudBeaver and Metabase have high memory usage, so running only one at a time is advised in development.

## Prerequisites

- Docker and Docker Compose installed
- Docker running properly
- Port 8970 available
- Sufficient memory (Metabase uses a lot of memory)

## Starting

### Basic Startup (Development)

```bash
# Development environment (default)
npm run metabase:start
```

### Environment-specific Startup

```bash
# Explicitly specify development
APP_ENV=develop npm run metabase:start

# Staging environment
APP_ENV=staging npm run metabase:start

# Production environment
APP_ENV=production npm run metabase:start
```

After startup, access http://localhost:8970.

## Stopping

```bash
npm run metabase:stop
```

## Other Commands

### View Logs

```bash
npm run metabase:logs
```

### Restart

```bash
npm run metabase:restart
```

## Database Connection Setup

### Initial Setup

1. After Metabase starts, access http://localhost:8970
2. Initial setup screen appears on first startup
3. Create admin account:
   - Name
   - Email address
   - Password
4. Complete initial setup

### Development Admin Account

The following admin account is configured for development:

| Item | Value |
|------|-------|
| Email | `admin@example.com` |
| Password | `metaadmin123` |

### Master Database Connection

1. Select add database from admin panel
2. Select PostgreSQL
3. Enter connection info:
   - **Connection Name**: `master` or `Master Database`
   - **Host**: `postgres-master` (within Docker) or `localhost`
   - **Port**: `5432`
   - **Database Name**: `webdb_master`
   - **Username**: `webdb`
   - **Password**: `webdb`
4. Test connection
5. Save connection

### Sharding Database Connections

Add connections for each sharding database:

| Connection Name | Host | Port | Database Name |
|-----------------|------|------|---------------|
| `sharding_db_1` | `postgres-sharding-1` | 5433 | `webdb_sharding_1` |
| `sharding_db_2` | `postgres-sharding-2` | 5434 | `webdb_sharding_2` |
| `sharding_db_3` | `postgres-sharding-3` | 5435 | `webdb_sharding_3` |
| `sharding_db_4` | `postgres-sharding-4` | 5436 | `webdb_sharding_4` |

## Creating Queries

1. Select "New" → "Question" from left menu
2. Select database and table
3. Configure filters, aggregations, grouping
4. Click "Visualize" to select chart type
5. "Save" to save query

### Native Query (SQL)

1. Select "New" → "SQL Query"
2. Select database
3. Enter SQL query
4. "Run" to execute query
5. "Save" to save query

## Creating Dashboards

1. Select "New" → "Dashboard" from left menu
2. Enter dashboard name
3. Add saved queries
4. Adjust layout
5. "Save" to save dashboard

## CloudBeaver vs Metabase Usage

| Use Case | Tool |
|----------|------|
| Direct data editing/manipulation | CloudBeaver |
| Table structure inspection | CloudBeaver |
| SQL script management | CloudBeaver |
| Data visualization/chart creation | Metabase |
| Dashboard creation/sharing | Metabase |
| Sharing data with non-engineers | Metabase |
| Report creation | Metabase |

## Configuration File Management

Metabase configuration files are managed by environment:

```
metabase/
└── config/
    ├── develop/      # Development config
    ├── staging/      # Staging config
    └── production/   # Production config
```

- Connection settings and dashboard configs are saved in each environment directory
- Configuration files can be managed with Git
- Different connection settings can be managed per environment

## Troubleshooting

### Container Won't Start

1. Check if Docker is running
   ```bash
   docker ps
   ```

2. Check if port 8970 is in use
   ```bash
   lsof -i :8970
   ```

3. Check logs
   ```bash
   npm run metabase:logs
   ```

### Cannot Connect to Database

1. Check if database files exist
   ```bash
   ls -la server/data/
   ```

2. Verify connection settings path
   - Correct path: `/data/master.db`
   - Incorrect path: `server/data/master.db`

### Configuration Files Not Saving

1. Check config directory permissions
   ```bash
   ls -la metabase/config/
   ```

2. Check if directories exist
   ```bash
   mkdir -p metabase/config/develop
   mkdir -p metabase/config/staging
   mkdir -p metabase/config/production
   ```

### Out of Memory

Metabase uses a lot of memory. Check the following:

- CloudBeaver is stopped
- Sufficient system memory available
- Stop other memory-intensive applications

## Technical Specifications

### Docker Compose Configuration

- **Image**: `metabase/metabase:latest`
- **Port**: 8970 (host) → 3000 (container)
- **Volumes**:
  - `./server/data:/data:ro` - Database files (read-only)
  - `./metabase/config/${APP_ENV}:/metabase-data` - Configuration files

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_ENV` | Environment name | `develop` |
| `MB_DB_FILE` | Metabase internal DB | `/metabase-data/metabase.db` |

## Related Documentation

- [README.md](../README.md) - Project overview
- [Database-Viewer.md](Database-Viewer.md) - CloudBeaver details
- [Sharding.md](Sharding.md) - Sharding details
- [Admin.md](Admin.md) - GoAdmin admin panel details

## Reference Links

- [Metabase Official Site](https://www.metabase.com/)
- [Metabase Documentation](https://www.metabase.com/docs/)
- [Metabase Docker Documentation](https://www.metabase.com/docs/latest/installation-and-operation/running-metabase-on-docker)
