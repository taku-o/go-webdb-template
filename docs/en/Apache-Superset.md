**[日本語](../ja/Apache-Superset.md) | [English]**

# Apache Superset - Data Visualization & Analysis Tool

## Overview

Apache Superset is a data viewer and analysis tool for non-engineers. It connects to databases via web browser, enabling query creation/execution and dashboard creation/sharing.

### Management Tool Roles

This project provides four management tools:

| Tool | Role | Port |
|------|------|------|
| GoAdmin | Admin panel for custom processing | 8081 |
| CloudBeaver | Web-based tool for data operations | 8978 |
| Metabase | Data visualization & analysis tool | 8970 |
| Apache Superset | Data visualization & analysis tool | 8088 |

**Important**: CloudBeaver, Metabase, and Apache Superset have high memory usage, so running only one at a time is advised in development.

## Prerequisites

- Docker and Docker Compose installed
- Docker running properly
- Port 8088 available
- PostgreSQL running (started via `docker-compose.postgres.yml`)
- Sufficient memory (Apache Superset uses a lot of memory)

## Starting

### Start PostgreSQL

PostgreSQL must be running before using Apache Superset:

```bash
# Start PostgreSQL
./scripts/start-postgres.sh start
```

### Start Apache Superset

```bash
# Start Apache Superset
./scripts/start-apache-superset.sh
```

After startup, access http://localhost:8088.

## Stopping

```bash
# Stop Apache Superset
docker-compose -f docker-compose.apache-superset.yml down
```

## Other Commands

### View Logs

```bash
docker-compose -f docker-compose.apache-superset.yml logs -f
```

### Restart

```bash
# Stop
docker-compose -f docker-compose.apache-superset.yml down

# Start
./scripts/start-apache-superset.sh
```

## Database Connection Setup

### Initial Setup

1. After Apache Superset starts, access http://localhost:8088
2. Initial setup screen appears on first startup
3. Login with default admin account:
   - **Username**: `admin`
   - **Password**: `admin`
4. You may be prompted to change password on first login (keeping `admin` is fine for development)

### PostgreSQL Database Connection

1. Select "Settings" → "Database Connections" from admin panel
2. Click "+ Database" button
3. Select "PostgreSQL" as database type
4. Enter connection info:
   - **Display Name**: `PostgreSQL - webdb` or any name
   - **Host**: `postgres` (within Docker network) or `host.docker.internal` (connecting to host machine's PostgreSQL)
   - **Port**: `5432`
   - **Database name**: `webdb`
   - **Username**: `webdb`
   - **Password**: `webdb`
5. Click "Test Connection" to verify connection
6. Click "Connect" to save connection settings

**Note**: Use `postgres` as hostname when connecting to PostgreSQL within Docker network. Use `host.docker.internal` when connecting directly from host machine.

## Creating Queries

### Query Execution in SQL Lab

1. Select "SQL Lab" → "SQL Editor" from top menu
2. Select database
3. Enter SQL query
4. Click "Run" button to execute query
5. Results are displayed
6. "Save" to save query

### Creating Visualizations

1. Execute query in SQL Lab
2. Click "Create Chart" button below results
3. Select chart type (table, bar chart, line chart, pie chart, etc.)
4. Configure data source, metrics, and dimensions
5. "Create Chart" to create chart
6. "Save" to save chart

## Creating Dashboards

1. Select "Dashboards" → "+ Dashboard" from top menu
2. Enter dashboard name
3. Add saved charts
4. Adjust layout
5. "Save" to save dashboard

### Adding Existing Charts

1. Click "+ Add Chart" in dashboard edit screen
2. Select saved chart
3. Chart is added to dashboard

## Data Viewer Comparison

This project provides multiple data viewers. Use according to your needs:

| Use Case | Tool |
|----------|------|
| Direct data editing/manipulation | CloudBeaver |
| Table structure inspection | CloudBeaver |
| SQL script management | CloudBeaver |
| Data visualization/chart creation | Metabase / Apache Superset |
| Dashboard creation/sharing | Metabase / Apache Superset |
| Sharing data with non-engineers | Metabase / Apache Superset |
| Report creation | Metabase / Apache Superset |
| Advanced visualization features | Apache Superset |

### Differences Between Metabase and Apache Superset

| Item | Metabase | Apache Superset |
|------|----------|-----------------|
| License | AGPL v3 | Apache 2.0 |
| Commercial use restrictions | Source code disclosure required | No restrictions |
| SQL Lab feature | Available | Available (more advanced) |
| Visualization types | Standard charts | More chart types |
| Customizability | Moderate | High |
| Learning curve | Somewhat low | Somewhat high |

## Configuration File Management

Apache Superset configuration files are stored in the following directory:

```
apache-superset/
└── data/                    # Data directory
    ├── superset.db          # Superset internal database
    ├── config/              # Configuration files
    └── uploads/             # Uploaded files
```

- Configuration files and dashboard settings are stored in `apache-superset/data/` directory
- Data is persisted to Docker volumes and retained after container restart
- Configuration files can be managed with Git (be careful with sensitive information)

## Troubleshooting

### Container Won't Start

1. Check if Docker is running
   ```bash
   docker ps
   ```

2. Check if port 8088 is in use
   ```bash
   lsof -i :8088
   ```

3. Check logs
   ```bash
   docker-compose -f docker-compose.apache-superset.yml logs
   ```

### Cannot Connect to PostgreSQL

1. Check if PostgreSQL is running
   ```bash
   docker ps | grep postgres
   ```

2. Verify PostgreSQL connection info
   - Host: `postgres` (within Docker network) or `host.docker.internal` (from host machine)
   - Port: `5432`
   - Database name: `webdb`
   - Username: `webdb`
   - Password: `webdb`

3. Check Docker network
   ```bash
   docker network ls
   docker network inspect <network_name>
   ```

### Data Not Being Saved

1. Check data directory permissions
   ```bash
   ls -la apache-superset/data/
   ```

2. Check Docker volume mounts
   - Verify `volumes` section in `docker-compose.apache-superset.yml`
   - Verify `./apache-superset/data:/app/superset_home` is correctly configured

### Out of Memory

Apache Superset uses a lot of memory. Check the following:

- CloudBeaver and Metabase are stopped
- Sufficient system memory available
- Stop other memory-intensive applications

### First Startup is Slow

On first startup, Apache Superset takes time for database initialization and setup. Wait a few minutes before accessing.

## Technical Specifications

### Docker Compose Configuration

- **Image**: `apache/superset:latest` (or stable version)
- **Port**: 8088 (host) → 8088 (container)
- **Volumes**:
  - `./apache-superset/data:/app/superset_home` - Data directory (persisted)

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SUPERSET_SECRET_KEY` | Superset secret key | `dev-secret-key-change-in-production` |

### Database Connection Info

| Item | Value |
|------|-------|
| Host (within Docker network) | `postgres` |
| Host (from host machine) | `host.docker.internal` |
| Port | `5432` |
| Database name | `webdb` |
| Username | `webdb` |
| Password | `webdb` |

## Security Considerations

### Authentication Settings

**Credentials** (Development):
- Username: `admin`
- Password: `admin`

**Notes**:
- Intended for development environment use
- Implement proper password policy and access control for production
- Password change on first login is encouraged

### Network Access

- Apache Superset is only accessible on localhost
- External access is not intended
- Proper network and firewall configuration required for production

### Database Connection

- PostgreSQL connection info is sensitive
- Use environment variables or secret management systems for production
- Be careful not to commit connection info to Git

## Related Documentation

- [README.md](../../README.md) - Project overview
- [Database-Viewer.md](Database-Viewer.md) - CloudBeaver details
- [Metabase.md](Metabase.md) - Metabase details
- [Sharding.md](Sharding.md) - Sharding details
- [Admin.md](Admin.md) - GoAdmin admin panel details
- [License-Survey.md](License-Survey.md) - License survey results

## Reference Links

- [Apache Superset Official Site](https://superset.apache.org/)
- [Apache Superset GitHub](https://github.com/apache/superset)
- [Apache Superset Documentation](https://superset.apache.org/docs/)
- [Apache Superset Docker Documentation](https://superset.apache.org/docs/installation/installing-superset-using-docker-compose)
