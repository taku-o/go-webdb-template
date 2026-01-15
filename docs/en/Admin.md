**[日本語](../ja/Admin.md) | [English]**

# Admin Panel Documentation

## Overview

This is an admin panel using the GoAdmin framework. It operates on port 8081 as an independent service from the main service (port 8080).

### Main Features

- **Table Management**: List display and CRUD operations for Users/Posts tables
- **Sharding Support**: Integrated display of data from all shards
- **Authentication & Authorization**: GoAdmin built-in authentication
- **Custom Pages**: Dashboard, user registration form

## Starting the Admin Panel

### Prerequisites

- Database is set up
- GoAdmin migrations are applied

```bash
# Start PostgreSQL container
./scripts/start-postgres.sh start

# Apply migrations (first time only)
./scripts/migrate.sh all
```

### Start Admin Panel

```bash
cd server
APP_ENV=develop go run cmd/admin/main.go
```

Access the admin panel at http://localhost:8081/admin.

## Login

### Default Credentials (Development)

| Item | Value |
|------|-------|
| URL | http://localhost:8081/admin/login |
| Username | `admin` |
| Password | `admin123` |

### Changing Credentials

Credentials can be changed in the config file (`config/develop.yaml`).

```yaml
admin:
  port: 8081
  auth:
    username: admin
    password: your_password_here
  session:
    lifetime: 7200  # Session lifetime in seconds
```

### Production Settings

Use environment variables or a separate config file for production.

```yaml
# config/production.yaml
admin:
  port: 8081
  auth:
    username: ${ADMIN_USERNAME}
    password: ${ADMIN_PASSWORD}
  session:
    lifetime: 3600
```

## Table Management

### User Management (Users)

**URL**: http://localhost:8081/admin/info/users

#### List Display

| Column | Description | Operations |
|--------|-------------|------------|
| ID | User ID | Sortable |
| Name | Username | Sortable & Filterable |
| Email | Email address | Sortable & Filterable |
| Created At | Registration date | Sortable |
| Updated At | Last update date | Sortable |

#### Operations

- **Create New**: Click "New" button
- **Edit**: Click "Edit" icon in list
- **Delete**: Click "Delete" icon in list
- **Export**: Click "Export" button for CSV output

### Post Management (Posts)

**URL**: http://localhost:8081/admin/info/posts

#### List Display

| Column | Description | Operations |
|--------|-------------|------------|
| ID | Post ID | Sortable |
| User ID | Poster's user ID | Sortable & Filterable |
| Title | Post title | Sortable & Filterable |
| Content | Post body | - |
| Created At | Post date | Sortable |
| Updated At | Last update date | Sortable |

#### Operations

- **Create New**: Click "New" button
- **Edit**: Click "Edit" icon in list
- **Delete**: Click "Delete" icon in list
- **Export**: Click "Export" button for CSV output

## Custom Pages

### Dashboard

**URL**: http://localhost:8081/admin/

Top page displayed after login.

#### Display Content

- **Statistics**: User count, post count
- **Quick Actions**: Links to user registration, post creation
- **System Information**: Project name, GoAdmin version

### User Registration Page

**URL**: http://localhost:8081/admin/user/register

User registration page with custom form.

#### Input Fields

| Field | Required | Description |
|-------|----------|-------------|
| Name | Yes | Up to 100 characters |
| Email | Yes | Up to 255 characters, valid format |

#### Validation

- Name is required, up to 100 characters
- Email is required, valid format, up to 255 characters
- Email duplicate check

#### Completion Page

After successful registration, redirects to a completion page displaying the registered user information.

## Menu Structure

```
Admin Panel
├── Dashboard (/admin/)
├── Users
│   ├── List (/admin/info/users)
│   └── Register (/admin/user/register)
└── Posts
    └── List (/admin/info/posts)
```

## Troubleshooting

### Cannot Login

1. **Check Credentials**: Verify username/password in config file (`config/develop.yaml`)
2. **Check Database**: Verify GoAdmin migration (002_goadmin.sql) is applied
3. **Browser Cache**: Clear cookies and retry

```bash
# Check GoAdmin tables (PostgreSQL)
psql -h localhost -p 5432 -U webdb -d webdb_master -c "SELECT * FROM goadmin_users;"
```

### Admin Panel Won't Start

1. **Port Conflict**: Check if port 8081 is in use

```bash
lsof -i :8081
```

2. **Check Config File**: Verify no format errors in `config/develop.yaml`

3. **Check Dependencies**: Verify required packages are installed

```bash
cd server
go mod tidy
```

### Data Not Displaying

1. **Check Database Connection**: Verify no DB connection errors in server logs
2. **Check Tables**: Verify data exists in dm_users and dm_posts tables

```bash
# Check in PostgreSQL
psql -h localhost -p 5433 -U webdb -d webdb_sharding_1 -c "SELECT COUNT(*) FROM dm_users_000;"
psql -h localhost -p 5433 -U webdb -d webdb_sharding_1 -c "SELECT COUNT(*) FROM dm_posts_000;"
```

### Session Expires Quickly

Check `session.lifetime` in config file. Value is in seconds.

```yaml
admin:
  session:
    lifetime: 7200  # 2 hours
```

## Technical Specifications

### Libraries Used

- GoAdmin v1.2.26
- AdminLTE Theme
- Gorilla Mux (HTTP Router)

### Architecture

```
┌─────────────────────────────────────────────────────────────┐
│              Admin Service (Port 8081)                       │
│                                                               │
│  ┌────────────────────────────────────────────────────────┐ │
│  │                  GoAdmin Engine                         │ │
│  │  • Admin Plugin (Auto-generated CRUD)                   │ │
│  │  • Custom Pages                                         │ │
│  │  • Authentication                                       │ │
│  └──────────────────────┬─────────────────────────────────┘ │
│                         │                                     │
│  ┌──────────────────────▼─────────────────────────────────┐ │
│  │                 GORM Manager                            │ │
│  │  (Reuses existing connection management)                │ │
│  └──────────────────────┬─────────────────────────────────┘ │
└─────────────────────────┼───────────────────────────────────┘
                          │
         ┌────────────────┴────────────────┐
         ▼                                  ▼
    ┌─────────┐                        ┌─────────┐
    │ Shard 1 │                        │ Shard 2 │
    └─────────┘                        └─────────┘
```

### File Structure

```
server/
├── cmd/
│   └── admin/
│       └── main.go           # Entry point
└── internal/
    └── admin/
        ├── config.go         # GoAdmin config
        ├── tables.go         # Table generators
        ├── sharding.go       # Cross-shard queries
        ├── auth/
        │   ├── auth.go       # Authentication logic
        │   └── session.go    # Session management
        └── pages/
            ├── pages.go               # Custom page base
            ├── home.go                # Dashboard
            ├── user_register.go       # User registration
            └── user_register_complete.go  # Registration complete
```
