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
│  │                 Usecase Layer                           │ │
│  │  • Business logic                                       │ │
│  │  • Transaction management                               │ │
│  │  • Cross-service coordination                           │ │
│  └──────────────────────┬─────────────────────────────────┘ │
│                         │                                     │
│  ┌──────────────────────▼─────────────────────────────────┐ │
│  │                 Service Layer                           │ │
│  │  • Domain logic                                         │ │
│  │  • Domain-specific operations                           │ │
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
- Authentication and authorization checks

**Key Components**:
- `UserHandler`: User-related endpoints
- `PostHandler`: Post-related endpoints
- `Router`: Route definitions and middleware

**Constraints**:
- Minimal processing (business logic is delegated to usecase layer)
- Does not directly call service layer (goes through usecase layer)

### 2. Usecase Layer (`internal/usecase`)
**Location**: `internal/usecase/`

**Responsibilities**:
- Business logic implementation
- Transaction coordination
- Combining multiple services
- Business rule application

**Key Components**:
- `DmUserUsecase`: User business logic
- `DmPostUsecase`: Post business logic
- `EmailUsecase`: Email business logic
- `DmJobqueueUsecase`: Job queue business logic
- `TodayUsecase`: Date business logic

**Constraints**:
- Does not call other usecases
- Calls one or multiple services
- Does not directly call repository layer (goes through service layer)

### 3. Service Layer (`internal/service`)
**Location**: `internal/service/`

**Responsibilities**:
- Domain logic implementation
- Domain-specific operations
- Cross-shard operations
- Data transformation
- Domain-specific validation

**Key Components**:
- `DmUserService`: User domain logic
- `DmPostService`: Post domain logic
- `EmailService`: Email domain logic
- `DateService`: Date domain logic

**Constraints**:
- Focuses on single domain-specific processing
- Business logic is delegated to usecase layer

### 4. Repository Layer (`internal/repository`)
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

### 5. DB Layer (`internal/db`)
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

2. API Layer → Usecase Layer
   UserHandler.CreateUser()
   ↓
   DmUserUsecase.CreateDmUser(CreateDmUserRequest)

3. Usecase Layer → Service Layer
   Applies business rules
   ↓
   DmUserService.CreateDmUser(CreateDmUserRequest)

4. Service Layer → Repository Layer
   Validates domain rules
   ↓
   DmUserRepository.Create(user)

5. Repository Layer → DB Layer
   Constructs SQL query
   ↓
   DBManager.GetConnectionByKey(userID)

6. DB Layer
   Calculates shard ID using hash(userID)
   ↓
   Returns connection to appropriate shard

7. Repository Layer
   Executes INSERT statement
   ↓
   Returns created user

8. Service Layer → Usecase Layer
   Returns User
   ↓
   Usecase layer returns User

9. Usecase Layer → API Layer
   Returns User
   ↓
   UserHandler formats response

10. API Layer → Client
   HTTP 201 Created
   Body: {"id": 1, "name": "John", "email": "john@example.com", ...}
```

## Configuration Management

**Location**: `internal/config/`

The application uses environment-based configuration with YAML files:

- `config/develop/`: Development environment (PostgreSQL/MySQL)
- `config/staging/`: Staging environment (PostgreSQL/MySQL)
- `config/production/`: Production environment (PostgreSQL/MySQL)

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

### Table Sharding Rules

分散テーブル（dm_で始まるテーブル）のシャーディング規則:

| テーブル群 | シャーディングキー | 計算式 |
|---|---|---|
| dm_users_NNN | id | id % 32 |
| dm_posts_NNN | user_id | user_id % 32 |
| dm_news | なし（masterテーブル） | - |

**テーブル番号の決定規則**: dm_postsのシャーディングキーとしてuser_id（dm_usersのid）を使用することで、同じユーザーに属するdm_usersレコードとdm_postsレコードは同じテーブル番号のテーブルに配置されます。

シャーディングキーの計算方法（UUIDv7ベース）:
- UUIDの後ろ2文字を16進数として解釈
- テーブル数（32）で割った余りがテーブル番号

例:
- dm_users.id = "019b6f83add07d6586044649c19fa5c4" → 後ろ2文字 "c4" = 196 → 196 % 32 = 4 → dm_users_004
- dm_posts.user_id = "019b6f83add07d6586044649c19fa5c4" → 後ろ2文字 "c4" = 196 → 196 % 32 = 4 → dm_posts_004

## Identifier Generation

分散環境で一意なIDを生成するための規則を定義しています。

### UUIDv7（Primary Identifier）

すべての識別子には **UUIDv7** (`github.com/google/uuid`) を使用します。

**用途**:
- dm_users, dm_postsなどのプライマリキー
- 分散環境で一意性が保証されるID
- 32文字の16進数文字列（ハイフンなし小文字）として表現

**理由**:
- 分散環境でもID重複なく一意なIDを生成可能
- 時間順序性があり、生成順に並べ替え可能
- 標準的なUUID形式で互換性が高い
- UUIDv7は時間ベースでデータベースインデックスに優しい

**使用方法**:
```go
import "github.com/taku-o/go-webdb-template/internal/util/idgen"

id, err := idgen.GenerateUUIDv7()
if err != nil {
    return fmt.Errorf("failed to generate ID: %w", err)
}
// id = "019b6f83add07d6586044649c19fa5c4" (32文字)
```

**ID形式**:
- 32文字の16進数文字列
- ハイフンなし
- 小文字のみ
- 例: `019b6f83add07d6586044649c19fa5c4`

### シャーディングキーの計算

UUIDv7からシャーディングキー（テーブル番号）を計算する方法:

```go
// UUIDの後ろ2文字を16進数として解釈し、テーブル数で割った余りを計算
import "github.com/taku-o/go-webdb-template/internal/db"

selector := db.NewTableSelector(db.DBShardingTableCount, db.DBShardingTablesPerDB)
tableNumber, err := selector.GetTableNumberFromUUID(uuid)
// uuid = "019b6f83add07d6586044649c19fa5c4"
// 後ろ2文字 "c4" = 196 (10進数)
// 196 % 32 = 4 → テーブル番号は4
```

### JavaScript/フロントエンドでのID扱い

IDは文字列型として扱います。UUIDv7は32文字の文字列なので、JavaScriptの数値精度の問題はありません。

```go
type DmUser struct {
    ID    string `json:"id" gorm:"primaryKey;type:varchar(32)"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

APIレスポンス例:
```json
{
  "id": "019b6f83add07d6586044649c19fa5c4",
  "name": "John Doe",
  "email": "john@example.com"
}
```

フロントエンド側では、IDを文字列として扱い、比較やハッシュマップのキーとして使用してください。

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
3. **Usecase Layer**: Wraps errors with business context
4. **API Layer**: Converts errors to HTTP status codes
5. **Client**: Displays user-friendly error messages

## Security Considerations

1. **CORS**: Configured in router to allow specific origins
2. **Input Validation**: Performed at API layer (format validation) and Service layer (domain validation)
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
- Future enhancement: Add Redis/Memcached layer between Usecase and Service

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

## CLI Commands

### CLI Layer (`cmd/list-dm-users`)

**Location**: `cmd/list-dm-users/main.go`

**Responsibilities**:
- Command-line argument parsing
- Input validation
- Configuration loading
- Layer initialization (Repository → Service → Usecase)
- Usecase layer invocation
- Output formatting (TSV)

**Key Components**:
- `main()`: Entry point
- `validateLimit()`: Input validation
- `printDmUsersTSV()`: Output formatting

### CLI Usecase Layer (`internal/usecase/cli`)

**Location**: `internal/usecase/cli/`

**Responsibilities**:
- CLI-specific business logic coordination
- Service layer invocation for CLI commands

**Key Components**:
- `ListDmUsersUsecase`: User list retrieval for CLI

**Constraints**:
- Uses existing service layer interfaces
- Does not contain domain logic (delegates to service layer)

### CLI Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│              CLI Layer (cmd/list-dm-users/main.go)          │
│  • コマンドライン引数の解析                                  │
│  • 引数のバリデーション                                      │
│  • 設定ファイルの読み込み                                    │
│  • GroupManagerの初期化                                     │
│  • レイヤーの初期化（Repository → Service → Usecase）      │
│  • usecase層の呼び出し                                      │
│  • 結果の出力（TSV形式）                                    │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│         Usecase Layer (internal/usecase/cli)                  │
│  • ListDmUsersUsecase                                        │
│  • ビジネスロジックの調整（CLI用）                           │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│         Service Layer (internal/service)                      │
│  • DmUserService                                            │
│  • ドメインロジック                                          │
│  • クロスシャード操作                                        │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│      Repository Layer (internal/repository)                  │
│  • DmUserRepository                                         │
│  • データアクセスの抽象化                                    │
│  • CRUD操作                                                 │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│           DB Layer (internal/db)                              │
│  • GroupManager                                             │
│  • シャーディング戦略                                        │
│  • 接続プール管理                                            │
└────────────────────────┬────────────────────────────────────┘
                         │
          ┌──────────────┴──────────────┐
          ▼                              ▼
    ┌─────────┐                    ┌─────────┐
    │ Shard 1 │                    │ Shard 2 │
    └─────────┘                    └─────────┘
```

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
- `github.com/lib/pq`: PostgreSQL driver
- `github.com/rs/cors`: CORS middleware
- `gorm.io/gorm`: GORM ORM library (v1.25.12)
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
