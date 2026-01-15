**[日本語](../ja/Release-Check.md) | [English]**

# Pre-Release Checklist

This document summarizes the verification items to perform before releasing to production.

---

## 1. Test Execution

### 1.1 Server-Side Tests

#### Unit Tests (Required)

```bash
cd server

# Run all unit tests
go test ./internal/... -v

# Run with coverage
go test ./internal/... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run specific packages only
go test ./internal/db/... -v          # Sharding logic
go test ./internal/repository/... -v  # Repository layer
go test ./internal/service/... -v     # Service layer
```

**Expected Result**: All tests PASS

#### Integration Tests (Required)

```bash
cd server

# Run integration tests
go test ./test/integration/... -v

# Output more detailed logs
go test ./test/integration/... -v -count=1
```

**Verification Items**:
- ✅ User CRUD Flow: User create/get/update/delete
- ✅ Post CRUD Flow: Post create/get/update/delete
- ✅ Cross-Shard Operations: Cross-shard query operation
- ✅ User-Post JOIN: User-Post JOIN operation

#### E2E Tests (Advised)

```bash
cd server

# Run E2E tests
go test ./test/e2e/... -v
```

**Verification Items**:
- ✅ All API endpoints working
- ✅ HTTP status code verification
- ✅ Error handling verification

---

### 1.2 Client-Side Tests

#### Unit Tests (Required)

```bash
cd client

# Run all unit tests
npm test

# Run with coverage
npm run test:coverage

# Watch mode (during development)
npm run test:watch

# Run excluding integration tests
npm test -- --testPathIgnorePatterns="integration"
```

**Expected Results**:
- API Client Tests: 7/7 passed
- Home Page Tests: 4/4 passed

#### E2E Tests (Advised)

```bash
cd client

# Start server and client, then run
npm run e2e

# Run in UI mode (for debugging)
npm run e2e:ui

# Run specific browser only
npx playwright test --project=chromium
```

**Verification Items**:
- ✅ User management flow
- ✅ Post management flow
- ✅ Cross-shard JOIN display

---

## 2. Build Verification

### 2.1 Server Build

```bash
cd server

# Run build
go build -o bin/server cmd/server/main.go

# Verify binary
ls -lh bin/server

# Verify operation
./bin/server
```

**Verification Items**:
- ✅ No build errors
- ✅ Binary starts normally
- ✅ Connection established to all Shards

### 2.2 Client Build

```bash
cd client

# Run production build
npm run build

# Verify build result
ls -lh .next

# Verify operation after build (advised)
npm start
```

**Verification Items**:
- ✅ No build errors
- ✅ No TypeScript type errors
- ✅ Optimized bundles generated

---

## 3. Environment Configuration Verification

### 3.1 Server Environment Variables

**Production config file**: `server/config/production.yaml`

```bash
# Verify environment variables
cat server/config/production.yaml

# Verify passwords and sensitive info are loaded from env vars
echo $DB_SHARD1_PASSWORD
echo $DB_SHARD2_PASSWORD
```

**Verification Items**:
- ✅ `production.yaml` created by copying `production.yaml.example`
- ✅ Database connection info correctly configured
- ✅ Production DB hostname, port, auth info correct
- ✅ Sensitive info loaded from environment variables
- ✅ `production.yaml` included in `.gitignore`

### 3.2 Client Environment Variables

```bash
# Verify environment variables
echo $NEXT_PUBLIC_API_BASE_URL
```

**Verification Items**:
- ✅ API base URL set to production environment URL
- ✅ Production environment variable file prepared

---

## 4. Database Verification

### 4.1 Migrations

```bash
# Use migration script
./scripts/migrate.sh all

# Or run directly with Atlas CLI (for PostgreSQL)
# Master
atlas migrate apply --dir file://db/migrations/master \
    --url "postgres://webdb:password@db-master.example.com:5432/webdb_master?sslmode=require"

# Sharding DBs
for i in 1 2 3 4; do
    atlas migrate apply --dir file://db/migrations/sharding \
        --url "postgres://webdb:password@db-sharding-${i}.example.com:5432/webdb_sharding_${i}?sslmode=require"
done
```

**Verification Items**:
- ✅ Same schema created in all Shards
- ✅ Tables, indexes, constraints created correctly

### 4.2 Connection Test

```bash
cd server

# Connection test in production environment
APP_ENV=production go run cmd/server/main.go
```

**Verification Items**:
- ✅ Log output: "Successfully connected to all database shards"
- ✅ Ping/connection to each Shard successful

---

## 5. Operation Verification (Manual Testing)

### 5.1 Server Startup

```bash
cd server

# Start in production environment
APP_ENV=production ./bin/server

# Or
APP_ENV=production go run cmd/server/main.go
```

**Verification Items**:
- ✅ Server starts normally
- ✅ Listening on port 8080
- ✅ No error logs output

### 5.2 Client Startup

```bash
cd client

# Start with production build
NEXT_PUBLIC_API_BASE_URL=<production API URL> npm start
```

**Verification Items**:
- ✅ Client starts normally
- ✅ Accessible on port 3000

### 5.3 Basic Function Verification

#### User Management
1. Access http://localhost:3000/users
2. Create new user
3. Verify appears in user list
4. Delete user
5. Verify disappears from list

#### Post Management
1. Access http://localhost:3000/posts
2. Select user and create new post
3. Verify appears in post list
4. Delete post
5. Verify disappears from list

#### Cross-Shard Query
1. Access http://localhost:3000/user-posts
2. Verify User and Post are JOINed and displayed
3. Verify data retrieved from multiple Shards

---

## 6. Performance Testing (Advised)

### 6.1 Load Testing

```bash
# Simple load test using Apache Bench
ab -n 1000 -c 10 http://localhost:8080/api/users

# More detailed load testing (requires separate tools)
# Use k6, JMeter, Gatling, etc.
```

**Verification Items**:
- ✅ Response time within acceptable range
- ✅ Low error rate (<1%)
- ✅ No memory leaks

---

## 7. Security Check

### 7.1 Dependency Vulnerability Scan

```bash
# Server (Go)
cd server
go list -json -m all | nancy sleuth

# Client (npm)
cd client
npm audit
npm audit fix  # If auto-fix is possible
```

**Verification Items**:
- ✅ No critical vulnerabilities
- ✅ Dependencies are latest stable versions

### 7.2 Configuration File Verification

```bash
# Verify sensitive info not committed to Git
git log --all --full-history -- "*production.yaml"
git log --all --full-history -- "*.env"

# Verify .gitignore
cat .gitignore | grep -E "(production.yaml|.env|*.db)"
```

**Verification Items**:
- ✅ `production.yaml`, `.env` and other sensitive files not committed
- ✅ Database files not committed
- ✅ `.gitignore` properly configured

---

## 8. Deployment Preparation

### 8.1 Artifact Verification

```bash
# Server
ls -lh server/bin/server

# Client
ls -lh client/.next/

# Config files
ls -lh server/config/production.yaml
```

**Verification Items**:
- ✅ Server binary generated
- ✅ Client build generated
- ✅ Production config file prepared

### 8.2 Documentation Verification

```bash
# Verify required documentation is ready
ls -lh docs/
```

**Verification Items**:
- ✅ API.md: API documentation
- ✅ Architecture.md: Architecture documentation
- ✅ Sharding.md: Sharding strategy documentation
- ✅ Testing.md: Testing strategy documentation
- ✅ Project-Structure.md: Project structure
- ✅ Release-Check.md: This document

---

## 9. Release Checklist

Please verify the following items before release:

### Tests
- [ ] Server unit tests: PASS
- [ ] Server integration tests: PASS
- [ ] Client unit tests: PASS
- [ ] E2E tests: PASS (advised)

### Build
- [ ] Server build: Success
- [ ] Client build: Success
- [ ] TypeScript type check: No errors

### Environment Configuration
- [ ] Production config file created
- [ ] Environment variables configured
- [ ] Verified sensitive info not committed to Git

### Database
- [ ] All Shard connections verified
- [ ] Migrations executed
- [ ] Schema consistency verified

### Operation Verification
- [ ] Server startup verified
- [ ] Client startup verified
- [ ] Basic function manual testing done

### Security
- [ ] Dependency vulnerability scan performed
- [ ] Security configuration verified

### Documentation
- [ ] API documentation updated
- [ ] README.md updated
- [ ] Release notes created

---

## 10. Troubleshooting

### Tests Failing

```bash
# Clear cache and re-run
go clean -testcache
go test ./... -v

# Client
rm -rf client/node_modules client/.next
npm install
npm test
```

### Build Failing

```bash
# Reinstall dependencies
cd server && go mod tidy
cd client && rm -rf node_modules && npm install
```

### Database Connection Error

```bash
# Verify connection info
cat server/config/production.yaml

# Manual connection test
psql -h <host> -U <user> -d <database>

# Check logs
tail -f /var/log/app/server.log
```

---

## 11. Rollback Procedure

Emergency response procedure if issues occur after release:

1. **Immediately stop service**
   ```bash
   pkill -f server
   ```

2. **Revert to previous version**
   ```bash
   # Revert binary to previous version
   cp bin/server.backup bin/server

   # Revert config to previous version
   git checkout HEAD~1 server/config/production.yaml
   ```

3. **Restart service**
   ```bash
   ./bin/server
   ```

4. **Investigate issue**
   - Check log files
   - Collect error messages
   - Check database status

---

## Summary

Following this checklist for pre-release verification enables safe deployment to production.

**Important**:
- Release only after verifying all tests are successful
- Be especially careful during initial deployment to production
- Review rollback procedure in advance in case issues occur
