---
layout: default
title: Setup Guide
lang: en
---

# Setup Guide

This guide explains the detailed setup instructions to get the client server running.

---

## 1. Initial Setup

### Installing Package Applications

- **Docker**: [https://www.docker.com/](https://www.docker.com/)
- **Cursor**: [https://cursor.com/](https://cursor.com/)
- **Go**: [https://go.dev/dl/](https://go.dev/dl/)

### Installing Homebrew

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
eval "$(/opt/homebrew/bin/brew shellenv)"
```

#### Installing GitHub CLI

```bash
brew install gh
gh auth login
gh auth status
```

#### Installing Atlas

```bash
brew install ariga/tap/atlas
```

### Installing Node.js (nvm)

```bash
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.4/install.sh | bash
nvm ls-remote
nvm install v22.14.0
nvm use v22.14.0
nvm alias default v22.14.0
```

Add the following to `.bashrc`:

```bash
if [ -f ~/.nvm/nvm.sh ]
then
  source ~/.nvm/nvm.sh
fi
```

#### Installing Claude Code

```bash
npm install -g @anthropic-ai/claude-code
```

### Installing uv

```bash
brew install uv
```

#### Configuring Serena

Run the following in the project directory:

```bash
claude mcp add serena -- uvx --from git+https://github.com/oraios/serena serena-mcp-server --context ide-assistant --enable-web-dashboard false --project $(pwd)
```

Update Serena index:

```bash
uvx --from git+https://github.com/oraios/serena index-project
```

---

## 2. Installing Dependencies

### Server Side

```bash
cd server
go mod download
```

### Client Side

```bash
cd client
npm install --legacy-peer-deps
```

---

## 3. Database Setup

### Starting PostgreSQL

```bash
./scripts/start-postgres.sh start
```

**Connection Information**

| Database | Host | Port | User | Password | Database Name |
|----------|------|------|------|----------|---------------|
| Master | localhost | 5432 | webdb | webdb | webdb_master |
| Sharding 1 | localhost | 5433 | webdb | webdb | webdb_sharding_1 |
| Sharding 2 | localhost | 5434 | webdb | webdb | webdb_sharding_2 |
| Sharding 3 | localhost | 5435 | webdb | webdb | webdb_sharding_3 |
| Sharding 4 | localhost | 5436 | webdb | webdb | webdb_sharding_4 |

### Applying Migrations

```bash
./scripts/migrate.sh all
```

---

## 4. Starting Redis

```bash
# Start Redis
./scripts/start-redis.sh start
```

- Redis: http://localhost:6379

---

## 5. Auth0 Account Configuration

Refer to the following documentation to set up your Auth0 account:
- [Auth0 External ID Integration Setup Guide](https://github.com/taku-o/go-webdb-template/blob/master/docs/Partner-Idp-Auth0-Login.md)

---

## 6. Client Environment Variables Configuration

### Creating .env.local

Rename `client/.env.develop` to `client/.env.local` and set the following environment variables:

```
# Auth0 Configuration
AUTH0_ISSUER=https://your-tenant.auth0.com
AUTH0_CLIENT_ID=your-client-id
AUTH0_CLIENT_SECRET=your-client-secret
```

---

## 7. Starting Servers

### Starting API Server

```bash
cd server
APP_ENV=develop go run cmd/server/main.go
```

### Starting Admin Server

```bash
cd server
APP_ENV=develop go run cmd/admin/main.go
```

### Starting Client Server

```bash
cd client
npm run dev
```

---

## 8. URL Information

| Service | URL | Notes |
|---------|-----|-------|
| Client | http://localhost:3000 | Next.js Application |
| API Server doc | http://localhost:8080/docs | API Documentation UI |
| Admin Server | http://localhost:8081/admin | Admin Panel |

### Admin Server Credentials

| Item | Value |
|------|-------|
| ID | admin |
| Password | admin123 |

---

## Navigation

- [Home]({{ site.baseurl }}/pages/en/)
- [Project Overview]({{ site.baseurl }}/pages/en/about)
- [日本語]({{ site.baseurl }}/pages/ja/setup)
