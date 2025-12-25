# Architecture Documentation

## Overview

This project implements a database-sharded web application using Go for the backend and Next.js for the frontend. The architecture follows a layered design pattern with clear separation of concerns.

## System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         Client Layer                         │
│                      (Next.js 14 + React)                    │
└────────────────────────┬────────────────────────────────────┘
                         │ HTTP/REST
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                      Server Layer (Go)                       │
│                                                               │
│  ┌────────────────────────────────────────────────────────┐ │
│  │                    API Layer                            │ │
│  │  • HTTP Handlers (user_handler, post_handler)          │ │
│  │  • Request validation                                  │ │
│  │  • Response formatting                                 │ │
│  └──────────────────────┬─────────────────────────────────┘ │
│                         │                                     │
│  ┌──────────────────────▼─────────────────────────────────┐ │
│  │                 Service Layer                           │ │
│  │  • Business logic                                       │ │
│  │  • Transaction management                               │ │
│  │  • Cross-shard operations                               │ │
│  └──────────────────────┬─────────────────────────────────┘ │
│                         │                                     │
│  ┌──────────────────────▼─────────────────────────────────┐ │
│  │               Repository Layer                          │ │
│  │  • Data access abstraction                              │ │
│  │  • CRUD operations                                      │ │
│  │  • Query builders                                       │ │
│  └──────────────────────┬─────────────────────────────────┘ │
│                         │                                     │
│  ┌──────────────────────▼─────────────────────────────────┐ │
│  │                  DB Layer                               │ │
│  │  • Connection management                                │ │
│  │  • Sharding strategy                                    │ │
│  │  • Connection pooling                                   │ │
│  └──────────────────────┬─────────────────────────────────┘ │
└─────────────────────────┼───────────────────────────────────┘
                          │
         ┌────────────────┴────────────────┐
         ▼                                  ▼
    ┌─────────┐                        ┌─────────┐
    │ Shard 1 │                        │ Shard 2 │
    └─────────┘                        └─────────┘
```

## Layer Responsibilities

### 1. API Layer (`internal/api`)
**Location**: `internal/api/handler/`, `internal/api/router/`

**Responsibilities**:
- HTTP request handling
- Request validation and parsing
- Response serialization
- CORS configuration
- Error handling and HTTP status code mapping

**Key Components**:
- `UserHandler`: User-related endpoints
- `PostHandler`: Post-related endpoints
- `Router`: Route definitions and middleware

### 2. Service Layer (`internal/service`)
**Location**: `internal/service/`

**Responsibilities**:
- Business logic implementation
- Transaction coordination
- Cross-shard operations
- Data transformation
- Input validation

**Key Components**:
- `UserService`: User business logic
- `PostService`: Post business logic

### 3. Repository Layer (`internal/repository`)
**Location**: `internal/repository/`

**Responsibilities**:
- Data access abstraction
- SQL query construction
- CRUD operations
- Result mapping to domain models

**Key Components**:
- `UserRepository`: User data access (database/sql版)
- `PostRepository`: Post data access (database/sql版)
- `UserRepositoryGORM`: User data access (GORM版)
- `PostRepositoryGORM`: Post data access (GORM版)

**Note**: 現在は`database/sql`版のRepositoryを使用していますが、GORM版のRepositoryも実装済みです。将来的にService層をInterface化することで、GORM版への切り替えが可能です。

### 4. DB Layer (`internal/db`)
**Location**: `internal/db/`

**Responsibilities**:
- Database connection management
- Sharding strategy implementation
- Connection pooling
- Shard routing
- Writer/Reader分離 (GORM版)

**Key Components**:
- `Manager`: Multi-shard connection manager (database/sql版)
- `Connection`: Single database connection wrapper (database/sql版)
- `GORMManager`: Multi-shard GORM connection manager (GORM版)
- `GORMConnection`: Single GORM connection wrapper with Writer/Reader support
- `ShardingStrategy`: Shard selection logic
- `HashBasedSharding`: Hash-based sharding implementation (FNV-1a)

**GORM Features**:
- Writer/Reader分離: `gorm.io/plugin/dbresolver`を使用したRead/Write分離
- 接続管理: Writer用とReader用のDSNを個別に設定可能
- ポリシー設定: `random`または`round_robin`のReader選択ポリシー
- 後方互換性: 従来のDSN設定でも動作

## Data Flow

### Example: Creating a User

```
1. Client → API Layer
   POST /api/users
   Body: {"name": "John", "email": "john@example.com"}

2. API Layer → Service Layer
   UserHandler.CreateUser()
   ↓
   UserService.CreateUser(CreateUserRequest)

3. Service Layer → Repository Layer
   Validates business rules
   ↓
   UserRepository.Create(user)

4. Repository Layer → DB Layer
   Constructs SQL query
   ↓
   DBManager.GetConnectionByKey(userID)

5. DB Layer
   Calculates shard ID using hash(userID)
   ↓
   Returns connection to appropriate shard

6. Repository Layer
   Executes INSERT statement
   ↓
   Returns created user

7. Service Layer → API Layer
   Returns User
   ↓
   UserHandler formats response

8. API Layer → Client
   HTTP 201 Created
   Body: {"id": 1, "name": "John", "email": "john@example.com", ...}
```

## Configuration Management

**Location**: `internal/config/`

The application uses environment-based configuration with YAML files:

- `config/develop.yaml`: Development environment (SQLite)
- `config/staging.yaml`: Staging environment (PostgreSQL)
- `config/production.yaml`: Production environment (PostgreSQL/MySQL)

Environment selection is controlled by the `APP_ENV` environment variable:

```go
env := os.Getenv("APP_ENV")
if env == "" {
    env = "develop"
}
```

## Database Sharding

See [Sharding.md](./Sharding.md) for detailed information about the sharding strategy.

### Key Points:
- Hash-based sharding using `user_id` as the shard key
- Automatic shard selection via `DBManager.GetConnectionByKey()`
- Cross-shard queries supported for read operations
- Each shard contains a complete schema

## Error Handling

### HTTP Error Responses

```go
// API Layer
w.WriteHeader(http.StatusBadRequest)
json.NewEncoder(w).Encode(map[string]string{
    "error": "Invalid request",
})
```

### Error Propagation

1. **Repository Layer**: Returns Go errors
2. **Service Layer**: Wraps errors with context
3. **API Layer**: Converts errors to HTTP status codes
4. **Client**: Displays user-friendly error messages

## Security Considerations

1. **CORS**: Configured in router to allow specific origins
2. **Input Validation**: Performed at both API and Service layers
3. **SQL Injection Prevention**: Using parameterized queries
4. **Environment Variables**: Sensitive data stored in config files (excluded from git)

## Scalability

### Horizontal Scaling
- Add more shards by updating configuration
- Shard count is configurable
- Each shard can be on a separate database server

### Vertical Scaling
- Connection pooling configured per shard
- Database-specific optimizations (indexes, query optimization)

### Caching Strategy
- Future enhancement: Add Redis/Memcached layer between Service and Repository

## Deployment

### Development
```bash
APP_ENV=develop go run cmd/server/main.go
```

### Staging
```bash
APP_ENV=staging go run cmd/server/main.go
```

### Production
```bash
APP_ENV=production ./server
```

## Monitoring and Logging

### Current Implementation
- Basic error logging to stderr
- Database connection status logging

### Future Enhancements
- Structured logging (e.g., logrus, zap)
- Request tracing
- Performance metrics
- Health check endpoints

## Testing Strategy

See [Testing.md](./Testing.md) for comprehensive testing documentation.

### Testing Layers
- Unit tests for each layer
- Integration tests for multi-layer interactions
- E2E tests for complete workflows

## Admin Panel (GoAdmin)

### Overview

GoAdminを使用した管理画面を提供しています。メインサービス（ポート8080）とは独立したサービスとしてポート8081で動作します。

### Architecture

```
┌─────────────────────────────────────────────────────────────┐
│              Admin Service (Port 8081)                       │
│                                                               │
│  ┌────────────────────────────────────────────────────────┐ │
│  │                  GoAdmin Engine                         │ │
│  │  • Admin Plugin (CRUD自動生成)                          │ │
│  │  • Custom Pages (カスタムページ)                         │ │
│  │  • Authentication (認証・認可)                          │ │
│  └──────────────────────┬─────────────────────────────────┘ │
│                         │                                     │
│  ┌──────────────────────▼─────────────────────────────────┐ │
│  │                 GORM Manager                            │ │
│  │  (既存の接続管理を再利用)                                 │ │
│  └──────────────────────┬─────────────────────────────────┘ │
└─────────────────────────┼───────────────────────────────────┘
                          │
         ┌────────────────┴────────────────┐
         ▼                                  ▼
    ┌─────────┐                        ┌─────────┐
    │ Shard 1 │                        │ Shard 2 │
    └─────────┘                        └─────────┘
```

### Components

**Location**: `internal/admin/`

- `config.go`: GoAdmin設定構造体
- `tables.go`: テーブルジェネレータ（Users, Posts）
- `sharding.go`: クロスシャードクエリヘルパー
- `auth/`: 認証・セッション管理
- `pages/`: カスタムページ（ダッシュボード、ユーザー登録）

### Features

1. **テーブル管理**: Users/Postsテーブルの一覧表示・CRUD操作
2. **シャーディング対応**: 全シャードのデータを統合表示
3. **認証・認可**: GoAdmin組み込み認証機能
4. **カスタムページ**: ダッシュボード、ユーザー登録フォーム

### Entry Point

**Location**: `cmd/admin/main.go`

```bash
APP_ENV=develop go run cmd/admin/main.go
```

## Dependencies

### Server Dependencies
- `github.com/spf13/viper`: Configuration management
- `github.com/gorilla/mux`: HTTP routing
- `github.com/mattn/go-sqlite3`: SQLite driver (development)
- `github.com/lib/pq`: PostgreSQL driver (production)
- `github.com/rs/cors`: CORS middleware
- `gorm.io/gorm`: GORM ORM library (v1.25.12)
- `gorm.io/driver/sqlite`: GORM SQLite driver
- `gorm.io/driver/postgres`: GORM PostgreSQL driver
- `gorm.io/plugin/dbresolver`: GORM Writer/Reader分離プラグイン
- `gorm.io/sharding`: GORM sharding plugin (将来使用予定)

### Client Dependencies
- `next`: React framework
- `react`, `react-dom`: UI library
- TypeScript: Type safety

## Future Improvements

1. **Authentication & Authorization**: JWT-based auth
2. **Rate Limiting**: API rate limiting middleware
3. **Caching**: Redis integration
4. **Message Queue**: Async processing with RabbitMQ/Kafka
5. **Service Discovery**: For multi-instance deployments
6. **API Versioning**: /v1, /v2 endpoints
7. **GraphQL**: Alternative to REST API
8. **WebSocket**: Real-time updates
