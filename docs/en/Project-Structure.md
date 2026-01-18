**[日本語](../ja/Project-Structure.md) | [English]**

# Project Structure Plan

## Overview

This document records the template creation plan for the go-webdb-template project.

## Project Structure

```
go-webdb-template/
├── server/                      # Golang server
│   ├── cmd/
│   │   ├── server/
│   │   │   └── main.go         # API server entry point
│   │   └── jobqueue/
│   │       └── main.go         # JobQueue server entry point
│   ├── internal/
│   │   ├── api/                # API definition layer
│   │   │   ├── handler/        # HTTP handlers
│   │   │   │   ├── user_handler.go
│   │   │   │   └── user_handler_test.go
│   │   │   └── router/         # Routing
│   │   ├── usecase/            # Business logic layer
│   │   │   ├── api/            # API usecase layer
│   │   │   │   ├── dm_user_usecase.go
│   │   │   │   ├── dm_user_usecase_test.go
│   │   │   │   ├── dm_post_usecase.go
│   │   │   │   ├── dm_post_usecase_test.go
│   │   │   │   ├── dm_jobqueue_usecase.go
│   │   │   │   ├── dm_jobqueue_usecase_test.go
│   │   │   │   ├── email_usecase.go
│   │   │   │   ├── email_usecase_test.go
│   │   │   │   ├── today_usecase.go
│   │   │   │   └── today_usecase_test.go
│   │   │   ├── jobqueue/       # JobQueue usecase layer
│   │   │   │   ├── delay_print.go
│   │   │   │   └── delay_print_test.go
│   │   │   ├── admin/          # Admin usecase layer
│   │   │   │   ├── dm_user_register_usecase.go
│   │   │   │   ├── dm_user_register_usecase_test.go
│   │   │   │   ├── api_key_usecase.go
│   │   │   │   └── api_key_usecase_test.go
│   │   │   └── cli/            # CLI usecase layer
│   │   │       ├── list_dm_users_usecase.go
│   │   │       ├── list_dm_users_usecase_test.go
│   │   │       ├── generate_secret_usecase.go
│   │   │       ├── generate_secret_usecase_test.go
│   │   │       ├── generate_sample_usecase.go
│   │   │       └── generate_sample_usecase_test.go
│   │   ├── service/            # Domain logic layer
│   │   │   ├── user_service.go
│   │   │   ├── user_service_test.go
│   │   │   ├── secret_service.go
│   │   │   ├── secret_service_test.go
│   │   │   ├── api_key_service.go
│   │   │   ├── api_key_service_test.go
│   │   │   ├── generate_sample_service.go
│   │   │   ├── generate_sample_service_test.go
│   │   │   ├── delay_print_service.go
│   │   │   └── delay_print_service_test.go
│   │   ├── repository/         # Database processing layer
│   │   │   ├── user_repository.go
│   │   │   ├── user_repository_test.go
│   │   │   ├── dm_news_repository.go
│   │   │   └── dm_news_repository_test.go
│   │   ├── sql/                # SQL definition layer
│   │   ├── model/              # Data models
│   │   ├── db/                 # DB connection management
│   │   │   ├── connection.go  # DB connection pool management
│   │   │   ├── connection_test.go
│   │   │   ├── sharding.go    # Sharding strategy
│   │   │   ├── sharding_test.go
│   │   │   └── manager.go     # DB manager
│   │   ├── auth/               # Authentication & secret key management
│   │   │   ├── jwt.go          # JWT verification & generation
│   │   │   ├── secret.go       # Secret key generation
│   │   │   └── secret_test.go  # Secret key generation test
│   │   └── config/             # Configuration loading
│   │       ├── config.go       # Config struct and loading
│   │       └── config_test.go
│   ├── test/                   # Test utilities
│   │   ├── integration/        # Integration tests
│   │   │   ├── api_test.go
│   │   │   └── sharding_test.go
│   │   ├── e2e/                # E2E tests
│   │   │   └── user_flow_test.go
│   │   ├── fixtures/           # Test data
│   │   │   ├── users.json
│   │   │   └── posts.json
│   │   └── testutil/           # Test helpers
│   │       ├── db.go           # Test DB setup
│   │       └── mock.go         # Mock objects
│   ├── go.mod
│   └── go.sum
│
├── client/                      # Next.js + TypeScript
│   ├── app/                     # App Router
│   │   ├── api/
│   │   │   └── auth/
│   │   │       ├── [...nextauth]/route.ts    # NextAuth auth routes
│   │   │       ├── profile/route.ts           # Profile API
│   │   │       └── token/route.ts             # Token API
│   │   ├── dm_email/send/page.tsx             # Email sending page
│   │   ├── dm_movie/upload/page.tsx           # Video upload page
│   │   ├── dm-jobqueue/page.tsx               # Job queue page
│   │   ├── dm-posts/page.tsx                  # Post management page
│   │   ├── dm-user-posts/page.tsx            # User-Post JOIN page
│   │   ├── dm-users/page.tsx                 # User management page
│   │   ├── layout.tsx                         # Root layout
│   │   ├── page.tsx                           # Top page
│   │   └── globals.css                        # Global styles
│   ├── components/
│   │   ├── ui/                                # shadcn/ui components
│   │   ├── layout/                            # Layout components
│   │   └── TodayApiButton.tsx                # TodayApiButton component
│   ├── lib/
│   │   ├── api.ts                             # API client
│   │   └── auth.ts                            # Auth helpers
│   ├── types/                                 # Type definitions
│   │   ├── dm_post.ts
│   │   ├── dm_user.ts
│   │   ├── jobqueue.ts
│   │   └── next-auth.d.ts
│   ├── e2e/                                   # E2E tests
│   │   ├── auth-flow.spec.ts
│   │   ├── user-flow.spec.ts
│   │   ├── post-flow.spec.ts
│   │   ├── cross-shard.spec.ts
│   │   ├── email-send.spec.ts
│   │   └── csv-download.spec.ts
│   ├── src/__tests__/                        # Unit/integration tests
│   │   ├── integration/
│   │   ├── components/
│   │   └── lib/
│   ├── auth.ts                                # NextAuth config
│   ├── jest.config.js                         # Jest config
│   ├── jest.setup.js                          # Jest setup
│   ├── jest.polyfills.js                      # Jest polyfills
│   ├── playwright.config.ts                   # Playwright config
│   └── package.json
│
├── config/                      # Environment-specific config files
│   ├── develop.yaml            # Development config
│   ├── staging.yaml            # Staging config
│   └── production.yaml         # Production config
│
├── db/
│   └── migrations/             # Migration SQL
│       ├── shard1/             # Shard 1 migrations
│       │   └── 001_init.sql
│       └── shard2/             # Shard 2 migrations
│           └── 001_init.sql
│
├── docs/
│   ├── plans/
│   │   └── project-structure.md  # This document
│   ├── Architecture.md         # Architecture description
│   ├── API.md                  # API specification
│   └── Sharding.md             # Sharding strategy document
│
├── .gitignore
└── README.md
```

## Layer Structure

### Server Side (Go)

1. **API Definition Layer** (`internal/api/`)
   - HTTP request/response processing
   - Routing definition
   - Validation (format checks)
   - Authentication & authorization checks

2. **Business Logic Layer** (`internal/usecase/`)
   - Core application logic
   - Transaction management
   - Processing combining multiple services

3. **Domain Logic Layer** (`internal/service/`)
   - Domain-specific logic
   - Domain-specific validation
   - Domain-specific business rules

4. **Database Processing Layer** (`internal/repository/`)
   - Database access
   - CRUD operation implementation
   - DB selection based on Shard Key

5. **SQL Definition Layer** (`internal/sql/`)
   - SQL query definitions
   - Query builder

6. **DB Connection Management Layer** (`internal/db/`)
   - Connection pool management for multiple DB shards
   - Sharding strategy implementation (Hash-based, Range-based, etc.)
   - DB connection lifecycle management

7. **Configuration Management Layer** (`internal/config/`)
   - Environment-specific config file loading
   - Config value validation
   - DB shard configuration management

### Client Side (Next.js + TypeScript)

- **App Router**: Page routing (`app/` directory)
- **Components**: Reusable UI components (including shadcn/ui)
- **Lib**: API calls and utility functions
- **Types**: TypeScript type definitions
- **Authentication**: Authentication via NextAuth (Auth.js) v5
- **UI**: shadcn/ui component library

## Data Models

### 1. User

| Column | Type | Description |
|--------|------|-------------|
| id | INTEGER | Primary key |
| name | TEXT | Username |
| email | TEXT | Email address |
| created_at | DATETIME | Creation datetime |
| updated_at | DATETIME | Update datetime |

### 2. Post

| Column | Type | Description |
|--------|------|-------------|
| id | INTEGER | Primary key |
| user_id | INTEGER | User ID (foreign key) |
| title | TEXT | Title |
| content | TEXT | Body |
| created_at | DATETIME | Creation datetime |
| updated_at | DATETIME | Update datetime |

### 3. JOIN Function

Provides functionality to JOIN User and Post to retrieve and display users and their posts.

## Technology Stack

### Server Side

- **Language**: Go 1.21+
- **Database**: PostgreSQL or MySQL (all environments)
- **Routing**: gorilla/mux
- **DB Connection**: GORM + gorm.io/driver/postgres
- **Configuration**: spf13/viper (YAML config file loading)
- **Sharding**: Custom implementation (Hash-based sharding, 8 logical shards, 4 physical DBs)
- **Testing**:
  - testing (standard library)
  - testify (assertions, mocks)
  - httptest (HTTP testing)
  - go-sqlmock (DB mocking)

### Client Side

- **Framework**: Next.js 14+ (App Router)
- **Language**: TypeScript 5+
- **UI Components**: shadcn/ui
- **Authentication**: NextAuth (Auth.js) v5
- **Styling**: Tailwind CSS
- **Form Management**: react-hook-form
- **Validation**: zod
- **File Upload**: Uppy (TUS protocol)
- **Testing**:
  - Jest (unit tests, integration tests)
  - React Testing Library (component tests)
  - Playwright (E2E tests)
  - MSW (API mocking)

## Feature List

### Server API

- `GET /api/users` - Get user list
- `GET /api/users/:id` - Get user details
- `POST /api/users` - Create user
- `PUT /api/users/:id` - Update user
- `DELETE /api/users/:id` - Delete user
- `GET /api/posts` - Get post list
- `GET /api/posts/:id` - Get post details
- `POST /api/posts` - Create post
- `PUT /api/posts/:id` - Update post
- `DELETE /api/posts/:id` - Delete post
- `GET /api/user-posts` - Get User-Post JOIN result

### Client Screens

- `/` - Top page (feature list, auth status display)
- `/dm-users` - User list, create, edit, delete, CSV download
- `/dm-posts` - Post list, create, edit, delete
- `/dm-user-posts` - User-Post JOIN result display (cross-shard query)
- `/dm_email/send` - Email sending page
- `/dm_movie/upload` - Video upload page (TUS protocol)
- `/dm-jobqueue` - Job queue page

## Sharding Strategy

### Overview

Data is distributed across multiple DB servers to ensure scalability.

### Shard Key

- **User table**: Uses `user_id` as Shard Key
- **Post table**: Uses `user_id` as Shard Key (data from same user placed in same Shard)

### Sharding Method

**Hash-based Sharding** adopted:
```
shard_id = hash(user_id) % shard_count
```

### Shard Configuration (Example)

- **Shard 1**: Users with even user_id and their posts
- **Shard 2**: Users with odd user_id and their posts

### Cross-Shard Queries

When retrieving User-Post JOIN results:
- Retrieve data from each Shard in parallel
- Merge at application layer and return

### Implementation Notes

1. **Abstraction at Repository Layer**: Shard selection logic hidden in Repository layer
2. **Transactions**: Transaction support only within single Shard
3. **Extensibility**: Design that can handle Shard count changes (Consistent Hashing considered for future)

## Environment-Specific Configuration

### Config File Structure

Loads appropriate config file based on `APP_ENV` environment variable value.

```
APP_ENV=develop   → config/develop.yaml
APP_ENV=staging   → config/staging.yaml
APP_ENV=production → config/production.yaml
```

### Configuration Items

Each environment config file includes:

1. **Server Configuration**
   - Port number
   - Timeout values

2. **Database Configuration**
   - Connection info per Shard (host, port, DB name, auth info)
   - Connection pool settings (max connections, idle timeout, etc.)

3. **Logging Configuration**
   - Log level
   - Log output destination

4. **CORS Configuration**
   - Allowed origins

### Config File Example (develop.yaml)

```yaml
server:
  port: 8080
  read_timeout: 30s
  write_timeout: 30s

database:
  shards:
    - id: 1
      host: localhost
      port: 5432
      name: app_db_shard1
      user: dev_user
      password: dev_password
      max_connections: 10
    - id: 2
      host: localhost
      port: 5433
      name: app_db_shard2
      user: dev_user
      password: dev_password
      max_connections: 10

logging:
  level: debug
  format: json

cors:
  allowed_origins:
    - http://localhost:3000
```

### Security Considerations

- Do not commit production config files to Git (add to `.gitignore`)
- Make sensitive info like passwords overridable via environment variables
- Include only config file templates (`*.yaml.example`) in repository

## Testing Strategy

### Test Levels

1. **Unit Tests**
   - Individual function/method tests
   - Isolate dependencies using mocks
   - Coverage target: 80%+

2. **Integration Tests**
   - Tests combining multiple layers
   - Use actual DB (test DB)
   - Verify API → Usecase → Service → Repository flow

3. **E2E Tests**
   - User scenario-based tests
   - Browser automation testing from frontend to backend

### Server-Side Testing Policy

#### Unit Tests
- Place `*_test.go` in each layer (handler, usecase, service, repository)
- Utilize table-driven tests
- Use testify/mock or go-sqlmock for mocking

```go
// Example: user_service_test.go
func TestUserService_GetUser(t *testing.T) {
    tests := []struct {
        name    string
        userID  int64
        want    *model.User
        wantErr bool
    }{
        // Test cases
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

#### Integration Tests
- Place in `test/integration/`
- Use test PostgreSQL (testutil.SetupTestGroupManager)
- Load test data from `test/fixtures/`
- Rollback with transactions (cleanup)

#### E2E Tests
- Place in `test/e2e/`
- Start server using httptest.Server
- Test API with actual HTTP requests

### Client-Side Testing Policy

#### Component Tests
- Place `.test.tsx` for each component
- Test rendering and interaction with React Testing Library
- Mock APIs with MSW

```typescript
// Example: UserList.test.tsx
describe('UserList', () => {
  it('displays user list', async () => {
    render(<UserList />)
    expect(await screen.findByText('John Doe')).toBeInTheDocument()
  })
})
```

#### Integration Tests
- Place in `__tests__/integration/`
- Test multiple component coordination
- Include API client tests

#### E2E Tests
- Place in `e2e/`
- Browser automation with Playwright
- Test entire user flows (create → list → edit → delete)

### Sharding Feature Tests

#### Sharding Strategy Tests
- Hash function consistency tests
- Shard selection logic tests
- Data distribution across multiple Shards verification

#### Cross-Shard Query Tests
- Data retrieval from each Shard
- Merge processing accuracy
- Parallel execution performance tests

### Test Execution in CI/CD

```yaml
# GitHub Actions example
- name: Run unit tests
  run: go test -v -cover ./...

- name: Run integration tests
  run: go test -v -tags=integration ./test/integration/...

- name: Run E2E tests
  run: go test -v -tags=e2e ./test/e2e/...
```

### Coverage Measurement

- Server side: `go test -coverprofile=coverage.out ./...`
- Client side: `npm run test:coverage`
- Visualize coverage reports in CI/CD

## Development Guidelines

1. **Large-scale Project Support**: Adopt scalable structure even for small samples
2. **Layer Separation**: Clearly separate responsibilities, improve maintainability
3. **Type Safety**: Ensure safety with TypeScript type definitions
4. **Test-Driven**: Ensure quality with unit/integration/E2E tests, target 80%+ coverage
5. **Documentation**: Improve visibility with Storybook and documentation
6. **Docker Environment**: Use PostgreSQL containers for easy environment setup
7. **Sharding Support**: Design assuming horizontal partitioning across multiple DBs
8. **Environment Separation**: Enable config switching for develop/staging/production

## Next Steps

1. Create directory structure
2. Create environment-specific config files (develop/staging/production.yaml)
3. Create server-side basic files (Go)
   - Implement config management layer with tests
   - Implement DB connection management layer (with Sharding support) with tests
   - Implement Repository/Service/API layers with unit tests
4. Create database migration files (for each Shard)
5. Create client-side basic files (Next.js)
   - Create components and tests
   - Create API client and tests
6. Implement integration tests and E2E tests
7. Implement each layer and verify Sharding operation
8. Check and improve test coverage
9. Setup Storybook
10. CI/CD configuration (GitHub Actions, etc.)
11. Create documentation (Architecture.md, API.md, Sharding.md, Testing.md)
