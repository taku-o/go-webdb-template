# Testing Documentation

## Overview

This project implements a comprehensive testing strategy covering unit tests, integration tests, and end-to-end tests for both server and client components.

## Testing Strategy

```
┌─────────────────────────────────────────────────────────────┐
│                     Testing Pyramid                          │
│                                                               │
│                        ╱╲                                     │
│                       ╱  ╲                                    │
│                      ╱ E2E ╲     ← Few, slow, high value     │
│                     ╱────────╲                                │
│                    ╱          ╲                               │
│                   ╱ Integration╲  ← Some, medium speed       │
│                  ╱──────────────╲                             │
│                 ╱                ╲                            │
│                ╱   Unit Tests     ╲ ← Many, fast, focused    │
│               ╱────────────────────╲                          │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

## Server Testing (Go)

### Unit Tests

**Location**: `server/internal/*/`
**Naming**: `*_test.go`
**Framework**: Go testing + testify

#### Running Unit Tests

```bash
cd server
go test ./internal/... -v
```

#### Example: Repository Unit Test

**File**: `server/internal/repository/user_repository_test.go`

```go
package repository_test

import (
    "testing"
    "database/sql"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    _ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("sqlite3", ":memory:")
    require.NoError(t, err)

    // Create schema
    _, err = db.Exec(`
        CREATE TABLE users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            email TEXT NOT NULL UNIQUE,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
    require.NoError(t, err)

    return db
}

func TestUserRepository_Create(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    repo := repository.NewUserRepository(nil)

    user := &model.User{
        Name:  "Test User",
        Email: "test@example.com",
    }

    err := repo.Create(db, user)
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
    assert.NotZero(t, user.CreatedAt)
}

func TestUserRepository_GetByID(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    repo := repository.NewUserRepository(nil)

    // Insert test data
    result, err := db.Exec(
        "INSERT INTO users (name, email) VALUES (?, ?)",
        "Test User", "test@example.com",
    )
    require.NoError(t, err)
    id, _ := result.LastInsertId()

    // Test retrieval
    user, err := repo.GetByID(db, id)
    assert.NoError(t, err)
    assert.Equal(t, "Test User", user.Name)
    assert.Equal(t, "test@example.com", user.Email)
}
```

#### Example: Service Unit Test

**File**: `server/internal/service/user_service_test.go`

```go
package service_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Mock repository
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(db *sql.DB, user *model.User) error {
    args := m.Called(db, user)
    return args.Error(0)
}

func TestUserService_CreateUser(t *testing.T) {
    mockRepo := new(MockUserRepository)
    service := service.NewUserService(mockRepo, nil)

    req := &service.CreateUserRequest{
        Name:  "Test User",
        Email: "test@example.com",
    }

    mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

    user, err := service.CreateUser(req)
    assert.NoError(t, err)
    assert.Equal(t, "Test User", user.Name)

    mockRepo.AssertExpectations(t)
}
```

### Integration Tests

**Location**: `server/test/integration/`
**Framework**: Go testing

Integration tests verify multiple layers working together with a real database.

#### Running Integration Tests

```bash
cd server
go test ./test/integration/... -v
```

#### Example: Repository + DB Integration Test

**File**: `server/test/integration/user_flow_test.go`

```go
package integration_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestUserCRUDFlow(t *testing.T) {
    // Setup test database with sharding
    cfg := &config.Config{
        Database: config.DatabaseConfig{
            Shards: []config.ShardConfig{
                {ID: 1, Driver: "sqlite3", Host: ":memory:"},
                {ID: 2, Driver: "sqlite3", Host: ":memory:"},
            },
        },
    }

    dbManager, err := db.NewManager(cfg)
    require.NoError(t, err)
    defer dbManager.CloseAll()

    // Initialize schema on all shards
    // ... (schema setup)

    repo := repository.NewUserRepository(dbManager)
    service := service.NewUserService(repo, dbManager)

    // Test Create
    createReq := &service.CreateUserRequest{
        Name:  "Integration Test User",
        Email: "integration@example.com",
    }
    user, err := service.CreateUser(createReq)
    require.NoError(t, err)
    assert.NotZero(t, user.ID)

    // Test Read
    retrieved, err := service.GetUser(user.ID)
    require.NoError(t, err)
    assert.Equal(t, user.Name, retrieved.Name)

    // Test Update
    updateReq := &service.UpdateUserRequest{
        Name:  "Updated Name",
        Email: "updated@example.com",
    }
    updated, err := service.UpdateUser(user.ID, updateReq)
    require.NoError(t, err)
    assert.Equal(t, "Updated Name", updated.Name)

    // Test Delete
    err = service.DeleteUser(user.ID)
    assert.NoError(t, err)

    // Verify deletion
    _, err = service.GetUser(user.ID)
    assert.Error(t, err)
}
```

#### Example: Cross-Shard Query Test

**File**: `server/test/integration/sharding_test.go`

```go
func TestCrossShardQuery(t *testing.T) {
    // Setup multi-shard environment
    // ...

    // Create users on different shards
    user1, _ := service.CreateUser(&service.CreateUserRequest{
        Name: "User 1", Email: "user1@example.com",
    })
    user2, _ := service.CreateUser(&service.CreateUserRequest{
        Name: "User 2", Email: "user2@example.com",
    })

    // Verify they're on different shards
    shard1 := dbManager.GetShardingStrategy().GetShardID(user1.ID)
    shard2 := dbManager.GetShardingStrategy().GetShardID(user2.ID)
    // Note: May be same shard with 2 shards, that's okay

    // Test GetAll returns users from both shards
    allUsers, err := service.GetAllUsers()
    require.NoError(t, err)
    assert.GreaterOrEqual(t, len(allUsers), 2)
}
```

### E2E Tests

**Location**: `server/test/e2e/`
**Framework**: Go testing + HTTP client

E2E tests verify the complete API flow from HTTP request to database.

#### Running E2E Tests

```bash
cd server
go test ./test/e2e/... -v
```

#### Example: API E2E Test

**File**: `server/test/e2e/api_test.go`

```go
package e2e_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func setupTestServer(t *testing.T) *httptest.Server {
    // Initialize full application stack
    cfg, _ := config.Load()
    // ... setup

    router := router.SetupRouter(userHandler, postHandler)
    return httptest.NewServer(router)
}

func TestUserAPI_CreateAndRetrieve(t *testing.T) {
    server := setupTestServer(t)
    defer server.Close()

    // Create user
    createReq := map[string]string{
        "name":  "E2E Test User",
        "email": "e2e@example.com",
    }
    body, _ := json.Marshal(createReq)

    resp, err := http.Post(
        server.URL+"/api/users",
        "application/json",
        bytes.NewBuffer(body),
    )
    require.NoError(t, err)
    assert.Equal(t, http.StatusCreated, resp.StatusCode)

    var user map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&user)
    assert.Equal(t, "E2E Test User", user["name"])

    userID := int(user["id"].(float64))

    // Retrieve user
    resp, err = http.Get(server.URL + fmt.Sprintf("/api/users/%d", userID))
    require.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.StatusCode)

    var retrieved map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&retrieved)
    assert.Equal(t, "E2E Test User", retrieved["name"])
}
```

### Test Utilities

**Location**: `server/test/testutil/`

Shared test utilities and helpers.

**File**: `server/test/testutil/db.go`

```go
package testutil

func SetupTestShards(t *testing.T, shardCount int) *db.Manager {
    shards := make([]config.ShardConfig, shardCount)
    for i := 0; i < shardCount; i++ {
        shards[i] = config.ShardConfig{
            ID:     i + 1,
            Driver: "sqlite3",
            Host:   ":memory:",
        }
    }

    cfg := &config.Config{
        Database: config.DatabaseConfig{Shards: shards},
    }

    manager, err := db.NewManager(cfg)
    require.NoError(t, err)

    // Initialize schema on all shards
    for i := 1; i <= shardCount; i++ {
        conn, _ := manager.GetConnection(i)
        InitSchema(t, conn.DB())
    }

    return manager
}

func InitSchema(t *testing.T, db *sql.DB) {
    schema := `
        CREATE TABLE IF NOT EXISTS users (...);
        CREATE TABLE IF NOT EXISTS posts (...);
    `
    _, err := db.Exec(schema)
    require.NoError(t, err)
}
```

### Test Fixtures

**Location**: `server/test/fixtures/`

**File**: `server/test/fixtures/users.go`

```go
package fixtures

func CreateTestUser(t *testing.T, service *service.UserService, name string) *model.User {
    req := &service.CreateUserRequest{
        Name:  name,
        Email: name + "@example.com",
    }
    user, err := service.CreateUser(req)
    require.NoError(t, err)
    return user
}
```

---

## Client Testing (Next.js/React)

### Unit Tests

**Location**: `client/src/**/__tests__/`
**Naming**: `*.test.tsx`, `*.test.ts`
**Framework**: Jest + React Testing Library

#### Running Unit Tests

```bash
cd client
npm test
```

#### Example: Component Unit Test

**File**: `client/src/components/__tests__/UserCard.test.tsx`

```typescript
import { render, screen } from '@testing-library/react'
import { UserCard } from '../UserCard'

describe('UserCard', () => {
  it('renders user information', () => {
    const user = {
      id: 1,
      name: 'Test User',
      email: 'test@example.com',
      created_at: '2024-01-15T10:00:00Z',
      updated_at: '2024-01-15T10:00:00Z',
    }

    render(<UserCard user={user} />)

    expect(screen.getByText('Test User')).toBeInTheDocument()
    expect(screen.getByText('test@example.com')).toBeInTheDocument()
  })
})
```

#### Example: API Client Unit Test

**File**: `client/src/lib/__tests__/api.test.ts`

```typescript
import { apiClient } from '../api'

// Mock fetch
global.fetch = jest.fn()

describe('apiClient', () => {
  beforeEach(() => {
    (fetch as jest.Mock).mockClear()
  })

  it('creates a user', async () => {
    const mockUser = {
      id: 1,
      name: 'Test User',
      email: 'test@example.com',
      created_at: '2024-01-15T10:00:00Z',
      updated_at: '2024-01-15T10:00:00Z',
    }

    ;(fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => mockUser,
    })

    const result = await apiClient.createUser({
      name: 'Test User',
      email: 'test@example.com',
    })

    expect(result).toEqual(mockUser)
    expect(fetch).toHaveBeenCalledWith(
      'http://localhost:8080/api/users',
      expect.objectContaining({
        method: 'POST',
      })
    )
  })
})
```

### Integration Tests

**Location**: `client/src/__tests__/integration/`
**Framework**: Jest + React Testing Library + MSW (Mock Service Worker)

#### Example: Page Integration Test

**File**: `client/src/__tests__/integration/UsersPage.test.tsx`

```typescript
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { rest } from 'msw'
import { setupServer } from 'msw/node'
import UsersPage from '@/app/users/page'

const server = setupServer(
  rest.get('http://localhost:8080/api/users', (req, res, ctx) => {
    return res(ctx.json([
      { id: 1, name: 'User 1', email: 'user1@example.com' },
      { id: 2, name: 'User 2', email: 'user2@example.com' },
    ]))
  }),

  rest.post('http://localhost:8080/api/users', (req, res, ctx) => {
    return res(ctx.status(201), ctx.json({
      id: 3,
      name: req.body.name,
      email: req.body.email,
    }))
  })
)

beforeAll(() => server.listen())
afterEach(() => server.resetHandlers())
afterAll(() => server.close())

describe('UsersPage', () => {
  it('displays users and allows creation', async () => {
    render(<UsersPage />)

    // Wait for users to load
    await waitFor(() => {
      expect(screen.getByText('User 1')).toBeInTheDocument()
    })

    // Create new user
    const nameInput = screen.getByLabelText('名前')
    const emailInput = screen.getByLabelText('メールアドレス')
    const submitButton = screen.getByRole('button', { name: '作成' })

    await userEvent.type(nameInput, 'New User')
    await userEvent.type(emailInput, 'new@example.com')
    await userEvent.click(submitButton)

    // Verify new user appears
    await waitFor(() => {
      expect(screen.getByText('New User')).toBeInTheDocument()
    })
  })
})
```

### E2E Tests

**Location**: `client/e2e/`
**Framework**: Playwright

#### Running E2E Tests

```bash
cd client
npx playwright test
```

#### Example: E2E Test

**File**: `client/e2e/user-flow.spec.ts`

```typescript
import { test, expect } from '@playwright/test'

test.describe('User Management Flow', () => {
  test('should create, view, and delete user', async ({ page }) => {
    await page.goto('http://localhost:3000')

    // Navigate to users page
    await page.click('text=ユーザー管理')

    // Create user
    await page.fill('input[type="text"]', 'E2E Test User')
    await page.fill('input[type="email"]', 'e2e@example.com')
    await page.click('button:has-text("作成")')

    // Verify user appears
    await expect(page.locator('text=E2E Test User')).toBeVisible()
    await expect(page.locator('text=e2e@example.com')).toBeVisible()

    // Delete user
    await page.click('button:has-text("削除")')
    await page.click('button:has-text("OK")') // Confirm dialog

    // Verify user is removed
    await expect(page.locator('text=E2E Test User')).not.toBeVisible()
  })
})
```

---

## Test Coverage

### Measuring Coverage

**Server**:
```bash
cd server
go test ./internal/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

**Client**:
```bash
cd client
npm test -- --coverage
```

### Coverage Goals

- **Unit Tests**: 80%+ coverage
- **Integration Tests**: Critical paths covered
- **E2E Tests**: Main user flows covered

---

## Continuous Integration

### GitHub Actions Example

**File**: `.github/workflows/test.yml`

```yaml
name: Tests

on: [push, pull_request]

jobs:
  server-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run tests
        run: |
          cd server
          go test ./... -v -coverprofile=coverage.out
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./server/coverage.out

  client-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
      - name: Install dependencies
        run: |
          cd client
          npm ci
      - name: Run tests
        run: |
          cd client
          npm test -- --coverage
```

---

## Best Practices

1. **Test Independence**: Each test should be independent
2. **Clear Naming**: Use descriptive test names
3. **Arrange-Act-Assert**: Follow AAA pattern
4. **Mock External Dependencies**: Use mocks for external services
5. **Test Edge Cases**: Not just happy paths
6. **Fast Tests**: Keep unit tests fast (<100ms each)
7. **Cleanup**: Always cleanup resources (defer, afterEach)

---

## Running All Tests

```bash
# Server tests
cd server
go test ./... -v

# Client tests
cd client
npm test

# E2E tests
cd client
npx playwright test
```
