**[日本語](../ja/Database-Viewer.md) | [English]**

# Database Viewer (CloudBeaver) Documentation

## Overview

CloudBeaver is a web-based database management tool. In this project, it is introduced as an admin application for data operations.

### Main Features

- **Web-based Database Management**: Access databases from browser
- **Visual Data Operations**: Easy table structure inspection, data viewing and editing
- **SQL Execution**: Execute SQL queries to manipulate data
- **Resource Manager**: Save and manage frequently used SQL scripts
- **Environment-specific Configuration**: Separate configurations for develop, staging, and production environments

### Role Division

This project provides two admin applications:

- **CloudBeaver**: For data operations (this tool)
- **GoAdmin**: For custom processing (see `docs/Admin.md`)

## Prerequisites

- Docker and Docker Compose installed
- Docker running properly
- Port 8978 available (no conflicts with other services)
- Database files exist in `server/data/` directory

## Starting

### Basic Startup

```bash
# Development environment (default)
npm run cloudbeaver:start

# Or explicitly specify environment
APP_ENV=develop npm run cloudbeaver:start
```

### Environment-specific Startup

```bash
# Development environment
APP_ENV=develop npm run cloudbeaver:start

# Staging environment
APP_ENV=staging npm run cloudbeaver:start

# Production environment
APP_ENV=production npm run cloudbeaver:start
```

**Note**: If environment variable `APP_ENV` is not set, it starts as `develop` environment by default.

### Startup Verification

After startup, access the following URL to verify CloudBeaver's Web UI is displayed:

- **URL**: http://localhost:8978

## Stopping

```bash
npm run cloudbeaver:stop
```

## Other Commands

### View Logs

```bash
npm run cloudbeaver:logs
```

### Restart

```bash
npm run cloudbeaver:restart
```

## Database Connection Setup

### Initial Setup

On first CloudBeaver startup, admin account setup and driver activation are required.

#### 1. Create Admin Account

1. Access http://localhost:8978
2. Setup wizard is displayed
3. Create admin account:
   - **Username**: `cbadmin`
   - **Password**: `Admin123`
4. Click "Next" → "Finish" to complete setup

#### 2. PostgreSQL Connection Setup

After setup, manually configure database connections from the Web UI.

### Master Database Connection

1. Access CloudBeaver Web UI (http://localhost:8978)
2. Click "Add Connection" or "New Connection"
3. Select "PostgreSQL" as database type
4. Enter connection info:
   - **Connection Name**: `master` or `Master Database`
   - **Host**: `postgres-master` (within Docker) or `localhost`
   - **Port**: `5432`
   - **Database Name**: `webdb_master`
   - **Username**: `webdb`
   - **Password**: `webdb`
5. Click "Test Connection" to verify connection
6. Click "Save" to save connection settings

### Sharding Database Connections

Similarly, add connections to the following 4 sharding databases:

| Connection Name | Host | Port | Database Name |
|-----------------|------|------|---------------|
| `sharding_db_1` | `postgres-sharding-1` | 5433 | `webdb_sharding_1` |
| `sharding_db_2` | `postgres-sharding-2` | 5434 | `webdb_sharding_2` |
| `sharding_db_3` | `postgres-sharding-3` | 5435 | `webdb_sharding_3` |
| `sharding_db_4` | `postgres-sharding-4` | 5436 | `webdb_sharding_4` |

### Connection Settings Storage

Connection settings are saved in environment-specific config directories:

- **Development**: `cloudbeaver/config/develop/`
- **Staging**: `cloudbeaver/config/staging/`
- **Production**: `cloudbeaver/config/production/`

Configuration files can be managed with Git. Since settings are separated by environment, connection settings must be configured individually for each environment.

## Database Operations

### Displaying Table List

1. Select connected database
2. Expand "Tables" from left navigation tree
3. Table list is displayed

### Viewing Data

1. Select table
2. Click "Data" tab
3. Table data is displayed

### Executing SQL Queries

1. Select connected database
2. Click "SQL Editor" tab
3. Enter SQL query
4. Click "Execute" button to run
5. Results are displayed

**Note**: Database files are mounted read-only, so data modification is not possible. For data changes, use existing APIs or management tools.

## Resource Manager

Resource Manager is a feature for saving and managing frequently used SQL scripts.

### Creating Scripts

1. Open "Resource Manager" in CloudBeaver Web UI
2. Click "New Script"
3. Enter script name and SQL query
4. Click "Save" to save

### Script Storage Location

Scripts saved in Resource Manager are stored in user project directories:

- `cloudbeaver/config/{env}/user-projects/{username}/`

Example (development environment, cbadmin user):
- `cloudbeaver/config/develop/user-projects/cbadmin/sql-1.sql`

Script files can be managed with Git.

### Using Scripts

1. Select script from Resource Manager
2. Click "Execute" button to run
3. Or load script into SQL Editor and execute

### Editing/Deleting Scripts

- **Edit**: Select script in Resource Manager, click "Edit"
- **Delete**: Select script in Resource Manager, click "Delete"

## Environment-specific Configuration

### Configuration Directory Structure

```
cloudbeaver/
└── config/
    ├── develop/                    # Development environment config
    │   ├── GlobalConfiguration/    # Connection settings etc.
    │   └── user-projects/          # User scripts
    │       └── cbadmin/            # cbadmin user scripts
    ├── staging/                    # Staging environment config
    └── production/                 # Production environment config
```

### Environment-specific Configuration Management

- When CloudBeaver starts for each environment, the corresponding config directory is mounted
- Connection settings are saved per environment, allowing different connection settings for each
- Configuration files can be managed with Git

### Checking Configuration Files

Configuration files are saved in the following directories:

- Development: `cloudbeaver/config/develop/`
- Staging: `cloudbeaver/config/staging/`
- Production: `cloudbeaver/config/production/`

## Troubleshooting

### Container Won't Start

**Problem**: CloudBeaver doesn't start when running `npm run cloudbeaver:start`

**Solutions**:
1. Check if Docker is running
   ```bash
   docker ps
   ```
2. Check if port 8978 is in use
   ```bash
   lsof -i :8978
   ```
3. Check logs
   ```bash
   npm run cloudbeaver:logs
   ```
4. Start with different port number (if port conflict)
   ```bash
   CLOUDBEAVER_PORT=8979 npm run cloudbeaver:start
   ```

### Cannot Connect to Database

**Problem**: Cannot connect to database from CloudBeaver

**Solutions**:
1. Check if database files exist
   ```bash
   ls -la server/data/*.db
   ```
2. Verify mount configuration
   - Check `volumes` section in `docker-compose.yml`
   - Verify `./server/data:/data:ro` is correctly configured
3. Verify file paths
   - Database file path in CloudBeaver is `/data/master.db` etc.
   - Use container mount path (`/data`)

### Cannot Save Scripts to Resource Manager

**Problem**: Cannot save scripts to Resource Manager

**Solutions**:
1. Check user project directory permissions
   ```bash
   ls -la cloudbeaver/config/develop/user-projects/
   ```
2. Verify mount configuration
   - Check `volumes` section in `docker-compose.yml`
   - Verify `./cloudbeaver/config/${APP_ENV:-develop}:/opt/cloudbeaver/workspace` is correctly configured

### Configuration Files Not Saving

**Problem**: Connection settings not saving, or not separated by environment

**Solutions**:
1. Check if environment variable `APP_ENV` is correctly set
   ```bash
   echo $APP_ENV
   ```
2. Check if config directory is mounted
   - Check `volumes` section in `docker-compose.yml`
   - Verify `./cloudbeaver/config/${APP_ENV:-develop}:/opt/cloudbeaver/workspace` is correctly configured
3. Check if config directory exists
   ```bash
   ls -la cloudbeaver/config/develop/
   ```

### Port Conflict

**Problem**: Port 8978 is already in use

**Solutions**:
1. Check process using the port
   ```bash
   lsof -i :8978
   ```
2. Start with different port number
   ```bash
   CLOUDBEAVER_PORT=8979 npm run cloudbeaver:start
   ```
3. Change port number in `docker-compose.yml` (for permanent change)

## Security Considerations

### Database File Access

- Database files are mounted read-only (`:ro` option)
- Prevents accidental database file modification from CloudBeaver
- For data changes, use existing APIs or management tools

### Authentication Settings

**Credentials** (Development):
- Username: `cbadmin`
- Password: `Admin123`

**Notes**:
- Not intended for production use (proper access control required for production)

### Network Access

- CloudBeaver is only accessible on localhost
- External access is not intended

## Configuration File Management

### Git Management

CloudBeaver configuration files can be managed with Git:

- **Configuration Files**: `cloudbeaver/config/{env}/`
- **Connection Settings**: `cloudbeaver/config/{env}/GlobalConfiguration/`
- **Scripts**: `cloudbeaver/config/{env}/user-projects/{username}/`

### Configuration File Structure

Configuration files are separated by environment:

- `cloudbeaver/config/develop/`: Development environment config
- `cloudbeaver/config/staging/`: Staging environment config
- `cloudbeaver/config/production/`: Production environment config

When CloudBeaver starts for each environment, the corresponding config directory is mounted and connection settings are saved.

### Sharing Configuration Files

By managing configuration files with Git, settings can be shared among team members. However, if sensitive information (passwords, etc.) is included, manage appropriately.

## Reference Information

### Related Documentation

- `README.md`: Project overview and setup instructions
- `docs/Admin.md`: GoAdmin admin panel documentation
- `docs/Sharding.md`: Sharding detailed specifications

### Technology Stack

- **CloudBeaver**: https://cloudbeaver.io/
- **Docker**: Using Docker Compose
- **Database**: PostgreSQL

### Reference Links

- CloudBeaver Official Site: https://cloudbeaver.io/
- CloudBeaver GitHub: https://github.com/dbeaver/cloudbeaver
- CloudBeaver Docker: https://hub.docker.com/r/dbeaver/cloudbeaver
- CloudBeaver Documentation: https://cloudbeaver.io/docs/
