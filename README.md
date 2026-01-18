**[日本語](README.ja.md) | [English]**

# Go DB Project Sample

A template project for Go API server and database configuration designed to handle large-scale users and high-traffic operations.

Database partitioning with numbered DB tables. Server architecture scalable by adding API servers and web servers.
Maintainable source code structure with rate limiting, various logging, database schema management and other operational features.

**GitHub Pages**: [https://taku-o.github.io/go-webdb-template/pages/en/](https://taku-o.github.io/go-webdb-template/pages/en/)

## Project Overview

- **Server**: Go language, Layered Architecture, Database Sharding support
- **Client**: Next.js 14 (App Router), TypeScript
- **Database**: PostgreSQL/MySQL (all environments)
- **Testing**: Go testing, Jest, Playwright

## Features

- ✅ **Sharding Support**: Table-based sharding (32 partitions) for data distribution across multiple DBs
- ✅ **GORM Support**: Writer/Reader separation supported (GORM v1.25.12)
- ✅ **GoAdmin Dashboard**: Web-based admin panel for data management
- ✅ **Layer Separation**: Clear responsibility separation with API, Usecase, Service, Repository, and DB layers
- ✅ **Environment-specific Config**: Configuration switching for develop/staging/production environments
- ✅ **Type Safety**: Type definitions with TypeScript
- ✅ **Testing**: Unit/Integration/E2E test support
- ✅ **Rate Limiting**: API call restriction by IP address (using ulule/limiter)
- ✅ **Job Queue**: Background job processing using Redis + Asynq
- ✅ **Email Sending**: Email functionality supporting stdout, Mailpit, and AWS SES
- ✅ **File Upload**: Large file upload via TUS protocol (local/S3 storage support)
- ✅ **Logging**: Access logs, email logs, SQL logs output
- ✅ **Docker Support**: Dockerized API server, Admin server, and client server

## Setup

### Prerequisites

- Go 1.21+
- Node.js 18+
- Docker (for PostgreSQL container)
- Atlas CLI (for database migration management)
  - Installation: `brew install ariga/tap/atlas` (macOS)
  - Verify installation: `atlas version`
  - Details: https://atlasgo.io/
- Redis (optional, for job queue functionality)
  - Can be started using Docker (`./scripts/start-redis.sh`)

### 1. Install Dependencies

#### Server

```bash
cd server
go mod download
```

#### Client

```bash
cd client
npm install
```

### 2. Database Setup

This project uses PostgreSQL and manages migrations with [Atlas](https://atlasgo.io/).

#### Start PostgreSQL

```bash
# Start PostgreSQL containers (master + 4 sharding instances)
./scripts/start-postgres.sh start
```

**Connection Info** (Development):

| Database | Host | Port | User | Password | Database Name |
|----------|------|------|------|----------|---------------|
| Master | localhost | 5432 | webdb | webdb | webdb_master |
| Sharding 1 | localhost | 5433 | webdb | webdb | webdb_sharding_1 |
| Sharding 2 | localhost | 5434 | webdb | webdb | webdb_sharding_2 |
| Sharding 3 | localhost | 5435 | webdb | webdb | webdb_sharding_3 |
| Sharding 4 | localhost | 5436 | webdb | webdb | webdb_sharding_4 |

#### Apply Migrations

```bash
# Apply migrations to all databases (including initial data)
./scripts/migrate.sh all
```

#### Stop PostgreSQL

```bash
./scripts/start-postgres.sh stop
```

#### Start MySQL (Optional)

You can use MySQL instead of PostgreSQL.

```bash
# Start MySQL containers (master + 4 sharding instances)
./scripts/start-mysql.sh start
```

**Connection Info** (Development):

| Database | Host | Port | User | Password | Database Name |
|----------|------|------|------|----------|---------------|
| Master | localhost | 3306 | webdb | webdb | webdb_master |
| Sharding 1 | localhost | 3307 | webdb | webdb | webdb_sharding_1 |
| Sharding 2 | localhost | 3308 | webdb | webdb | webdb_sharding_2 |
| Sharding 3 | localhost | 3309 | webdb | webdb | webdb_sharding_3 |
| Sharding 4 | localhost | 3310 | webdb | webdb | webdb_sharding_4 |

#### Apply MySQL Migrations

```bash
# Development environment migrations
atlas migrate apply --dir "file://db/migrations/master-mysql" --url "mysql://webdb:webdb@localhost:3306/webdb_master"
atlas migrate apply --dir "file://db/migrations/sharding_1-mysql" --url "mysql://webdb:webdb@localhost:3307/webdb_sharding_1"
atlas migrate apply --dir "file://db/migrations/sharding_2-mysql" --url "mysql://webdb:webdb@localhost:3308/webdb_sharding_2"
atlas migrate apply --dir "file://db/migrations/sharding_3-mysql" --url "mysql://webdb:webdb@localhost:3309/webdb_sharding_3"
atlas migrate apply --dir "file://db/migrations/sharding_4-mysql" --url "mysql://webdb:webdb@localhost:3310/webdb_sharding_4"

# Test environment migrations
./scripts/migrate-test-mysql.sh
```

#### Stop MySQL

```bash
./scripts/start-mysql.sh stop
```

#### Switching Database Type

Switch databases in `config/{env}/config.yaml` using `DB_TYPE`:

```yaml
# Use PostgreSQL (default)
DB_TYPE: postgresql

# Use MySQL
DB_TYPE: mysql
```

When using MySQL, settings are loaded from `database.mysql.yaml`.

#### Sharding Configuration

This project distributes **8 logical shards** across **4 physical databases**:

| Logical Shard ID | Table Range | Physical Database |
|------------------|-------------|-------------------|
| 1 | _000 to _003 | webdb_sharding_1 (port 5433) |
| 2 | _004 to _007 | webdb_sharding_1 (port 5433) |
| 3 | _008 to _011 | webdb_sharding_2 (port 5434) |
| 4 | _012 to _015 | webdb_sharding_2 (port 5434) |
| 5 | _016 to _019 | webdb_sharding_3 (port 5435) |
| 6 | _020 to _023 | webdb_sharding_3 (port 5435) |
| 7 | _024 to _027 | webdb_sharding_4 (port 5436) |
| 8 | _028 to _031 | webdb_sharding_4 (port 5436) |

#### Generating Migrations for Schema Changes

```bash
# After modifying master.hcl
atlas migrate diff <migration_name> \
    --dir file://db/migrations/master \
    --to file://db/schema/master.hcl \
    --dev-url "postgres://webdb:webdb@localhost:5432/webdb_master?sslmode=disable"

# After modifying sharding.hcl
atlas migrate diff <migration_name> \
    --dir file://db/migrations/sharding \
    --to file://db/schema/sharding.hcl \
    --dev-url "postgres://webdb:webdb@localhost:5433/webdb_sharding_1?sslmode=disable"
```

For details, see [docs/en/Atlas-Operations.md](docs/en/Atlas-Operations.md).

#### Lazy Connection & Auto-Reconnection

This project implements the following features:

- **Lazy Connection**: DB connection is established on first query, not at server startup
- **Auto-Reconnection**: Automatically reconnects when database recovers
- **Retry Functionality**: Retries up to 3 times with 1-second intervals on connection errors

### 3. Start Server

#### API Server

```bash
cd server
APP_ENV=develop go run cmd/server/main.go
```

Server starts at http://localhost:8080.

#### JobQueue Server

```bash
cd server
APP_ENV=develop go run cmd/jobqueue/main.go
```

The JobQueue server runs an HTTP server (port 8082) and an Asynq server in parallel. The HTTP server provides a `/health` endpoint for health monitoring. The Asynq server processes jobs from Redis. Make sure Redis is running before starting the JobQueue server.

### 4. Start Admin Panel

```bash
cd server
APP_ENV=develop go run cmd/admin/main.go
```

Admin panel starts at http://localhost:8081/admin.

**Credentials** (Development):
- Username: `admin`
- Password: `admin123`

**Database Connection**:
The GoAdmin server uses PostgreSQL. It reads connection info from the config file (`config/{env}/database.yaml`) and connects to the master database (`webdb_master`).

Before starting, verify:
1. PostgreSQL container is running (`./scripts/start-postgres.sh start`)
2. Migrations are applied (`./scripts/migrate.sh all`)

For details, see [docs/en/Admin.md](docs/en/Admin.md).

### 5. Start Database Viewer (CloudBeaver)

```bash
# Development (default)
npm run cloudbeaver:start

# Start with specific environment
APP_ENV=develop npm run cloudbeaver:start
APP_ENV=staging npm run cloudbeaver:start
APP_ENV=production npm run cloudbeaver:start

# Stop
npm run cloudbeaver:stop
```

Database viewer starts at http://localhost:8978.

**Credentials** (Development):
- Username: `cbadmin`
- Password: `Admin123`

**Main Features**:
- Operate database from web browser
- View table structure and data
- Execute SQL queries
- Save and manage SQL scripts (Resource Manager)

For details, see [docs/en/Database-Viewer.md](docs/en/Database-Viewer.md).

### 6. Start Data Visualization Tool (Metabase)

```bash
# Development (default)
npm run metabase:start

# Start with specific environment
APP_ENV=develop npm run metabase:start
APP_ENV=staging npm run metabase:start
APP_ENV=production npm run metabase:start

# Stop
npm run metabase:stop
```

Metabase starts at http://localhost:8970.

**Main Features**:
- Data visualization and chart creation
- Dashboard creation and sharing
- Data analysis for non-engineers

**CloudBeaver vs Metabase**:
- **CloudBeaver**: Direct data editing/manipulation, table structure inspection
- **Metabase**: Data visualization/analysis, dashboard creation

**Note**: CloudBeaver and Metabase have high memory usage, so running only one at a time is advised in development.

For details, see [docs/en/Metabase.md](docs/en/Metabase.md).

### 7. Start Redis (For Job Queue)

```bash
# Start Redis
./scripts/start-redis.sh start

# Start Redis Insight (optional, data viewer)
./scripts/start-redis-insight.sh start
```

Redis starts at http://localhost:6379.
Redis Insight starts at http://localhost:8001.

For details, see [docs/en/Queue-Job.md](docs/en/Queue-Job.md).

### 8. Start Mailpit (Optional, For Email)

```bash
./scripts/start-mailpit.sh start
```

Mailpit starts at http://localhost:8025.

For details, see [docs/en/Send-Mail.md](docs/en/Send-Mail.md).

### 9. Start Client

#### Install Dependencies

```bash
cd client
npm install --legacy-peer-deps
```

**Note**: Use `--legacy-peer-deps` flag if there are peer dependency conflicts.

#### Set Environment Variables

**Generate AUTH_SECRET**:
```bash
# Run from project root
npm run cli:generate-secret
```
Copy the generated secret key.

Create `.env.local` with the following environment variables:
```
# NextAuth (Auth.js)
AUTH_SECRET=<secret key generated by npm run cli:generate-secret>
AUTH_URL=http://localhost:3000

# Auth0 Settings
AUTH0_ISSUER=https://your-tenant.auth0.com
AUTH0_CLIENT_ID=your-client-id
AUTH0_CLIENT_SECRET=your-client-secret
AUTH0_AUDIENCE=https://your-api-audience

# API Settings
NEXT_PUBLIC_API_KEY=your-api-key
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080

# Test Environment (required for tests)
APP_ENV=test
```

**Notes**:
- `AUTH_SECRET` is generated using `npm run cli:generate-secret` (uses `server/cmd/generate-secret`).
- `APP_ENV=test` is required for test execution (`npm test`, `npm run e2e`).

#### Auth0 Application Settings

Configure the following URLs in Auth0 Dashboard (`Applications > [Target App] > Settings`):

**Allowed Callback URLs:**
```
http://localhost:3000/api/auth/callback/auth0
```

**Allowed Logout URLs:**
```
http://localhost:3000
```

**Allowed Web Origins:**
```
http://localhost:3000
```

#### Start Development Server

```bash
cd client
npm run dev
```

Client starts at http://localhost:3000.

#### Tech Stack

- **Framework**: Next.js 14+ (App Router)
- **Language**: TypeScript 5+
- **UI Components**: shadcn/ui
- **Authentication**: NextAuth (Auth.js) v5
- **Styling**: Tailwind CSS
- **Form Management**: react-hook-form
- **Validation**: zod
- **File Upload**: Uppy (TUS protocol)
- **Testing**: Playwright (E2E), Jest (unit/integration), MSW (API mocking)

#### Available Scripts

**Development**:
- `npm run dev` - Start development server (port 3000)
- `npm run build` - Run production build
- `npm run start` - Start production build (port 3000)
- `npm run lint` - Run ESLint
- `npm run type-check` - Run TypeScript type checking
- `npm run format` - Check formatting with Prettier
- `npm run format:write` - Apply formatting with Prettier

**Testing**:
- `npm test` - Run Jest tests (unit/integration)
- `npm run test:watch` - Run Jest tests in watch mode
- `npm run test:coverage` - Get Jest test coverage
- `npm run e2e` - Run Playwright E2E tests
- `npm run e2e:ui` - Run Playwright E2E tests in UI mode
- `npm run e2e:headed` - Run Playwright E2E tests in headed mode

**Note**: `APP_ENV=test` is automatically set during test execution (included in `package.json` scripts).

### 10. Docker Environment (Optional)

You can start servers in Docker environment.

```bash
# Start PostgreSQL and Redis containers first
docker-compose -f docker-compose.postgres.yml up -d
docker-compose -f docker-compose.redis.yml up -d

# Build and start API server
docker-compose -f docker-compose.api.yml build
docker-compose -f docker-compose.api.yml up -d

# Build and start Admin server
docker-compose -f docker-compose.admin.yml build
docker-compose -f docker-compose.admin.yml up -d

# Build and start client server
docker-compose -f docker-compose.client.yml build
docker-compose -f docker-compose.client.yml up -d
```

**Access URLs after startup**:
- API Server: http://localhost:8080
- Admin Server: http://localhost:8081/admin
- Client: http://localhost:3000

For details, see [docs/en/Docker.md](docs/en/Docker.md).

## API Endpoints

### Basic Endpoints

#### User Related

- `GET /api/dm-users` - Get user list
- `GET /api/dm-users/{id}` - Get user
- `POST /api/dm-users` - Create user
- `PUT /api/dm-users/{id}` - Update user
- `DELETE /api/dm-users/{id}` - Delete user
- `GET /api/export/dm-users/csv` - Download user info as CSV

#### Post Related

- `GET /api/dm-posts` - Get post list
- `GET /api/dm-posts/{id}` - Get post
- `POST /api/dm-posts` - Create post
- `PUT /api/dm-posts/{id}` - Update post
- `DELETE /api/dm-posts/{id}` - Delete post
- `GET /api/dm-user-posts` - Join users and posts (cross-shard query)

#### Others

- `GET /api/today` - Get today's date (private API, Auth0 JWT required)
- `GET /health` - Health check (no authentication required, API Server)

#### JobQueue Server

- `GET http://localhost:8082/health` - Health check (no authentication required, JobQueue Server)

### Feature Endpoints

- `POST /api/email/send` - Send email
- `POST /api/dm-jobqueue/register` - Register job
- `POST /api/upload/dm_movie` - Upload file (TUS protocol)

### OpenAPI Specification

- `GET /docs` - API Documentation UI (Stoplight Elements)
- `GET /openapi.json` - OpenAPI 3.1 (JSON)
- `GET /openapi.yaml` - OpenAPI 3.1 (YAML)
- `GET /openapi-3.0.json` - OpenAPI 3.0.3 (JSON)

*OpenAPI document endpoints are accessible without authentication.

For details, see [docs/en/API.md](docs/en/API.md).

## Feature Documentation

See the following documents for detailed usage instructions:

- [Job Queue](docs/en/Queue-Job.md) - Background job processing using Redis + Asynq
- [Email Sending](docs/en/Send-Mail.md) - Email sending with stdout, Mailpit, and AWS SES support
- [File Upload](docs/en/File-Upload.md) - Large file upload via TUS protocol
- [Logging](docs/en/Logging.md) - Access logs, email logs, SQL logs
- [Rate Limiting](docs/en/Rate-Limit.md) - Detailed API rate limit configuration
- [Docker](docs/en/Docker.md) - Starting and deploying in Docker environment

## API Rate Limiting

Requests to API endpoints are rate limited by IP address.

### Response Headers

All API responses include the following headers:

**Per-minute limit (always included):**

| Header | Description | Example |
|--------|-------------|---------|
| `X-RateLimit-Limit` | Limit per minute | `60` |
| `X-RateLimit-Remaining` | Remaining requests | `45` |
| `X-RateLimit-Reset` | Reset time (Unix timestamp) | `1706342400` |

**Per-hour limit (only when `requests_per_hour` is configured):**

| Header | Description | Example |
|--------|-------------|---------|
| `X-RateLimit-Hour-Limit` | Limit per hour | `1000` |
| `X-RateLimit-Hour-Remaining` | Remaining requests | `950` |
| `X-RateLimit-Hour-Reset` | Reset time (Unix timestamp) | `1706346000` |

### Rate Limit Exceeded

When limit is exceeded, HTTP 429 status code is returned:

```json
{
  "code": 429,
  "message": "Too Many Requests"
}
```

### Configuration

Rate limit settings are managed in `config/{env}/config.yaml`:

```yaml
api:
  rate_limit:
    enabled: true
    requests_per_minute: 60
    requests_per_hour: 1000
```

### Storage

Rate limit counters use different storage based on environment:

| Environment | Storage | Config File |
|-------------|---------|-------------|
| develop | In-Memory | `config/develop/cacheserver.yaml` |
| staging | Redis Cluster | `config/staging/cacheserver.yaml` |
| production | Redis Cluster | `config/production/cacheserver.yaml` |

When using Redis Cluster, configure addresses in `cacheserver.yaml`:

```yaml
redis:
  cluster:
    addrs:
      - host1:6379
      - host2:6379
      - host3:6379
```

In-Memory storage is used when `addrs` is empty or not configured.

### Verification

```bash
# Check rate limit headers
curl -i -H "Authorization: Bearer <YOUR_API_KEY>" http://localhost:8080/api/users

# Example response headers
# X-RateLimit-Limit: 60
# X-RateLimit-Remaining: 59
# X-RateLimit-Reset: 1706342460
```

For details, see [docs/en/Rate-Limit.md](docs/en/Rate-Limit.md).

## API Authentication

Access to API endpoints (`/api/*`) requires a JWT-based Public API key.

### Issuing API Keys

API keys can be issued from the GoAdmin dashboard.

1. Login to admin panel (http://localhost:8081/admin)
2. Select "Custom Pages" → "Issue API Key" from side menu
3. Click "Issue API Key" button
4. Download or copy the generated JWT token

### Generating Secret Key

A tool is provided to generate the secret key used for API key signing.

```bash
cd server
go run cmd/generate-secret/main.go
```

Set the generated secret key in `config/{env}/config.yaml` under `api.secret_key`.

### API Request Authentication

Send the JWT token in the `Authorization` header for API requests.

```bash
curl -H "Authorization: Bearer <YOUR_API_KEY>" http://localhost:8080/api/users
```

### Client Configuration

In Next.js client, set the API key in the `NEXT_PUBLIC_API_KEY` environment variable.

```bash
# client/.env.local
NEXT_PUBLIC_API_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

For detailed client setup instructions, see section 9 "Start Client".

### Error Responses

- `401 Unauthorized` - API key is invalid or not set
- `403 Forbidden` - Insufficient scope (GET without read scope, POST/PUT/DELETE without write scope)

Error response format:
```json
{
  "code": 401,
  "message": "Invalid API key"
}
```

## CLI Tools

Batch processing CLI tools are available.

### Sample Data Generation (generate-sample-data)

Generates sample data for development. Uses PostgreSQL to insert data into master/sharding databases.

#### Prerequisites

1. PostgreSQL container is running
2. Migrations are applied

```bash
# Start PostgreSQL
./scripts/start-postgres.sh start

# Apply migrations
./scripts/migrate.sh
```

#### Execution

```bash
cd server
APP_ENV=develop go run cmd/generate-sample-data/main.go
```

#### Generated Data

| Table | Database | Count |
|-------|----------|-------|
| dm_users_000-031 | sharding (distributed across 4 instances) | 100 records |
| dm_posts_000-031 | sharding (distributed across 4 instances) | 100 records |
| dm_news | master | 100 records |

For details, see [docs/en/Generate-Sample-Data.md](docs/en/Generate-Sample-Data.md).

### Secret Key Generation (generate-secret)

Generates a secret key for API key signing.

#### Execution

```bash
cd server
go run cmd/generate-secret/main.go
```

#### Output

A Base64-encoded 32-byte (256-bit) random secret key is displayed on standard output.

### User List Output (list-users)

Outputs user list in TSV format.

#### Build

```bash
cd server
go build -o bin/list-users ./cmd/list-users
```

#### Execution

```bash
# Default (20 records)
APP_ENV=develop ./bin/list-users

# Specify count (max 100)
APP_ENV=develop ./bin/list-users --limit 50
```

#### Options

| Option | Description | Default | Range |
|--------|-------------|---------|-------|
| `--limit` | Output count | 20 | 1-100 |

#### Output Format

TSV (tab-separated) format with the following fields:
- ID, Name, Email, CreatedAt, UpdatedAt

## Sharding Strategy

Uses table-based sharding (32 partitions, 8 logical shards).

```
table_number = id % 32      # 0-31
table_name = "dm_users_" + sprintf("%03d", table_number)  # dm_users_000 ~ dm_users_031
logical_shard_id = (table_number / 4) + 1  # 1-8
```

**Database Groups**:
- **Master Group**: Stores shared tables (dm_news) (PostgreSQL, port 5432)
- **Sharding Group**: Distributes 32-partitioned tables (dm_users_000-031, dm_posts_000-031) across 8 logical shards → 4 physical DBs

| Logical Shard | Table Range | Physical DB (PostgreSQL) |
|---------------|-------------|--------------------------|
| 1 | _000 to _003 | webdb_sharding_1 (port 5433) |
| 2 | _004 to _007 | webdb_sharding_1 (port 5433) |
| 3 | _008 to _011 | webdb_sharding_2 (port 5434) |
| 4 | _012 to _015 | webdb_sharding_2 (port 5434) |
| 5 | _016 to _019 | webdb_sharding_3 (port 5435) |
| 6 | _020 to _023 | webdb_sharding_3 (port 5435) |
| 7 | _024 to _027 | webdb_sharding_4 (port 5436) |
| 8 | _028 to _031 | webdb_sharding_4 (port 5436) |

For details, see [docs/en/Sharding.md](docs/en/Sharding.md).

## Configuration File Structure

Configuration files are organized by environment directories.

### Directory Structure

```
config/
├── develop/                  # Development environment config directory
│   ├── config.yaml           # Main config (server, admin, logging, cors, api)
│   ├── database.yaml         # Database config (groups structure)
│   └── cacheserver.yaml      # Cache server config (Redis Cluster)
├── production/               # Production environment config directory
│   ├── config.yaml.example   # Main config template
│   ├── database.yaml.example # Database config template
│   └── cacheserver.yaml.example # Cache server config template
└── staging/                  # Staging environment config directory
    ├── config.yaml           # Main config
    ├── database.yaml         # Database config
    └── cacheserver.yaml      # Cache server config

db/
└── migrations/
    ├── master/               # Master group migrations
    │   └── 001_init.sql      # news table
    └── sharding/             # Sharding group migrations
        ├── templates/        # Template files
        │   ├── users.sql.template
        │   └── posts.sql.template
        └── generated/        # Generated migrations
```

### Config File Loading Order

1. Load main config file (`config/{env}/config.yaml`)
2. Merge database config file (`config/{env}/database.yaml`)
3. Merge cache server config file (`config/{env}/cacheserver.yaml`) (optional)
4. Map unified config to `Config` struct
5. Override passwords with environment variables (`DB_PASSWORD_SHARD*`)

### Environment Switching

Switch environments using `APP_ENV` environment variable:

```bash
APP_ENV=develop go run cmd/server/main.go    # Development
APP_ENV=staging go run cmd/server/main.go    # Staging
APP_ENV=production go run cmd/server/main.go # Production
```

## GORM Support

GORM-based Repository with Writer/Reader separation is implemented.

### Writer/Reader Separation Example

`config/production/database.yaml`:
```yaml
database:
  shards:
    - id: 1
      driver: postgres
      writer_dsn: host=writer.example.com port=5432 user=app password=xxx dbname=db sslmode=require
      reader_dsns:
        - host=reader1.example.com port=5432 user=app password=xxx dbname=db sslmode=require
        - host=reader2.example.com port=5432 user=app password=xxx dbname=db sslmode=require
      reader_policy: round_robin
```

### Key Dependencies

- `gorm.io/gorm` v1.25.12
- `gorm.io/driver/postgres`
- `gorm.io/plugin/dbresolver` (Writer/Reader separation)
- `gorm.io/sharding` (for future use)
- `github.com/labstack/echo/v4` v4.13.3 (HTTP router)
- `github.com/danielgtaylor/huma/v2` v2.34.1 (OpenAPI spec auto-generation)
- `github.com/ulule/limiter/v3` v3.11.2 (rate limiting)
- `github.com/redis/go-redis/v9` v9.17.2 (Redis Cluster connection)

For details, see [docs/en/Architecture.md](docs/en/Architecture.md).

## License

MIT License
