**[日本語](../ja/Architecture.md) | [English]**

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
- `DmJobqueueUsecase`: Job queue business logic (for job registration)
- `DelayPrintUsecase`: Job processing business logic (for JobQueue server)
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
- `DelayPrintService`: Job processing business utility logic (e.g., outputting strings to stdout)
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
- `UserRepository`: User data access (database/sql version)
- `PostRepository`: Post data access (database/sql version)
- `UserRepositoryGORM`: User data access (GORM version)
- `PostRepositoryGORM`: Post data access (GORM version)

**Note**: Currently using `database/sql` version repositories, but GORM version is also implemented. Future interface abstraction at the Service layer will enable switching to GORM version.

### 5. DB Layer (`internal/db`)
**Location**: `internal/db/`

**Responsibilities**:
- Database connection management
- Sharding strategy implementation
- Connection pooling
- Shard routing
- Writer/Reader separation (GORM version)

**Key Components**:
- `Manager`: Multi-shard connection manager (database/sql version)
- `Connection`: Single database connection wrapper (database/sql version)
- `GORMManager`: Multi-shard GORM connection manager (GORM version)
- `GORMConnection`: Single GORM connection wrapper with Writer/Reader support
- `ShardingStrategy`: Shard selection logic
- `HashBasedSharding`: Hash-based sharding implementation (FNV-1a)

**GORM Features**:
- Writer/Reader separation: Read/Write separation using `gorm.io/plugin/dbresolver`
- Connection management: Separate DSN configuration for Writer and Reader
- Policy settings: `random` or `round_robin` Reader selection policy
- Backward compatibility: Works with traditional DSN settings

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

### Example: Job Processing Flow (JobQueue Server)

```
1. Redis → JobQueue Server
   Asynq Server retrieves job from Redis
   ↓
   Identify job type (JobTypeDelayPrint)

2. Processor Layer → Usecase Layer
   ProcessDelayPrintJob()
   ↓
   Parse and validate payload
   ↓
   DelayPrintUsecase.Execute(payload)

3. Usecase Layer → Service Layer
   Set default message
   ↓
   DelayPrintService.PrintMessage(message)

4. Service Layer
   Output string to stdout (with timestamp)
   ↓
   Flush buffer
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

Sharding rules for distributed tables (tables starting with dm_):

| Table Group | Sharding Key | Formula |
|---|---|---|
| dm_users_NNN | id | id % 32 |
| dm_posts_NNN | user_id | user_id % 32 |
| dm_news | None (master table) | - |

**Table Number Determination Rule**: Using user_id (dm_users id) as the sharding key for dm_posts ensures that dm_users records and dm_posts records belonging to the same user are placed in tables with the same table number.

Sharding key calculation (UUIDv7 based):
- Interpret the last 2 characters of UUID as hexadecimal
- Remainder of division by table count (32) is the table number

Example:
- dm_users.id = "019b6f83add07d6586044649c19fa5c4" → last 2 chars "c4" = 196 → 196 % 32 = 4 → dm_users_004
- dm_posts.user_id = "019b6f83add07d6586044649c19fa5c4" → last 2 chars "c4" = 196 → 196 % 32 = 4 → dm_posts_004

## Identifier Generation

Defines rules for generating unique IDs in distributed environments.

### UUIDv7 (Primary Identifier)

All identifiers use **UUIDv7** (`github.com/google/uuid`).

**Usage**:
- Primary keys for dm_users, dm_posts, etc.
- IDs with guaranteed uniqueness in distributed environments
- Represented as 32-character hexadecimal string (lowercase, no hyphens)

**Reasons**:
- Can generate unique IDs without duplication even in distributed environments
- Has temporal ordering, can be sorted by generation order
- High compatibility with standard UUID format
- UUIDv7 is time-based and database index friendly

**How to use**:
```go
import "github.com/taku-o/go-webdb-template/internal/util/idgen"

id, err := idgen.GenerateUUIDv7()
if err != nil {
    return fmt.Errorf("failed to generate ID: %w", err)
}
// id = "019b6f83add07d6586044649c19fa5c4" (32 characters)
```

**ID Format**:
- 32-character hexadecimal string
- No hyphens
- Lowercase only
- Example: `019b6f83add07d6586044649c19fa5c4`

### Sharding Key Calculation

How to calculate sharding key (table number) from UUIDv7:

```go
// Interpret UUID last 2 characters as hex and calculate remainder by table count
import "github.com/taku-o/go-webdb-template/internal/db"

selector := db.NewTableSelector(db.DBShardingTableCount, db.DBShardingTablesPerDB)
tableNumber, err := selector.GetTableNumberFromUUID(uuid)
// uuid = "019b6f83add07d6586044649c19fa5c4"
// last 2 chars "c4" = 196 (decimal)
// 196 % 32 = 4 → table number is 4
```

### Handling IDs in JavaScript/Frontend

IDs are handled as string type. Since UUIDv7 is a 32-character string, there are no JavaScript number precision issues.

```go
type DmUser struct {
    ID    string `json:"id" gorm:"primaryKey;type:varchar(32)"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

API response example:
```json
{
  "id": "019b6f83add07d6586044649c19fa5c4",
  "name": "John Doe",
  "email": "john@example.com"
}
```

On the frontend, handle IDs as strings and use them for comparison or as hashmap keys.

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

### CLI Layer (`cmd/generate-secret`)

**Location**: `cmd/generate-secret/main.go`

**Responsibilities**:
- Layer initialization (Service → Usecase)
- Usecase layer invocation
- Output to stdout

**Key Components**:
- `main()`: Entry point and output control

### CLI Usecase Layer (`internal/usecase/cli`)

**Location**: `internal/usecase/cli/`

**Responsibilities**:
- CLI-specific business logic coordination
- Service layer invocation for CLI commands

**Key Components**:
- `ListDmUsersUsecase`: User list retrieval for CLI
- `GenerateSecretUsecase`: Secret key generation for CLI
- `GenerateSampleUsecase`: Sample data generation for CLI

**Constraints**:
- Uses existing service layer interfaces
- Does not contain domain logic (delegates to service layer)

### Secret Key Generation (`internal/auth/secret.go`)

**Location**: `internal/auth/secret.go`

**Responsibilities**:
- Generate cryptographically secure random secret keys
- Base64 encode the generated keys

**Key Components**:
- `GenerateSecretKey()`: Generates a 32-byte (256-bit) random secret key and returns it as a Base64-encoded string

**Security**:
- Uses `crypto/rand` for secure random number generation

### CLI Architecture Diagrams

#### list-dm-users

```
┌─────────────────────────────────────────────────────────────┐
│              CLI Layer (cmd/list-dm-users/main.go)          │
│  • Command-line argument parsing                            │
│  • Argument validation                                      │
│  • Configuration file loading                               │
│  • GroupManager initialization                              │
│  • Layer initialization (Repository → Service → Usecase)   │
│  • Usecase layer invocation                                 │
│  • Result output (TSV format)                               │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│         Usecase Layer (internal/usecase/cli)                  │
│  • ListDmUsersUsecase                                        │
│  • Business logic coordination (for CLI)                     │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│         Service Layer (internal/service)                      │
│  • DmUserService                                            │
│  • Domain logic                                              │
│  • Cross-shard operations                                    │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│      Repository Layer (internal/repository)                  │
│  • DmUserRepository                                         │
│  • Data access abstraction                                   │
│  • CRUD operations                                          │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│           DB Layer (internal/db)                              │
│  • GroupManager                                             │
│  • Sharding strategy                                         │
│  • Connection pool management                                │
└────────────────────────┬────────────────────────────────────┘
                         │
          ┌──────────────┴──────────────┐
          ▼                              ▼
    ┌─────────┐                    ┌─────────┐
    │ Shard 1 │                    │ Shard 2 │
    └─────────┘                    └─────────┘
```

#### generate-secret

```
┌─────────────────────────────────────────────────────────────┐
│          CLI Layer (cmd/generate-secret/main.go)             │
│  • Entry point                                               │
│  • Layer initialization (Service → Usecase)                 │
│  • Usecase layer invocation                                  │
│  • Result output (stdout)                                    │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│      Usecase Layer (internal/usecase/cli)                    │
│  • GenerateSecretUsecase                                    │
│  • Business logic coordination (for CLI)                     │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│         Service Layer (internal/service)                      │
│  • SecretService                                            │
│  • Domain logic                                              │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│         Auth Layer (internal/auth)                            │
│  • GenerateSecretKey()                                       │
│  • Secret key generation (shared library)                    │
│  • crypto/rand + encoding/base64                            │
└─────────────────────────────────────────────────────────────┘
```

#### generate-sample-data

```
┌─────────────────────────────────────────────────────────────┐
│       CLI Layer (cmd/generate-sample-data/main.go)           │
│  • Entry point                                               │
│  • Configuration file loading                                │
│  • GroupManager initialization                               │
│  • Layer initialization (Repository → Service → Usecase)    │
│  • Usecase layer invocation                                  │
│  • Result output (log to stdout)                             │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│      Usecase Layer (internal/usecase/cli)                    │
│  • GenerateSampleUsecase                                    │
│  • Business logic coordination (for CLI)                     │
│  • GenerateDmUsers() → GenerateDmPosts() → GenerateDmNews() │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│         Service Layer (internal/service)                      │
│  • GenerateSampleService                                    │
│  • Domain logic                                              │
│  • Data generation using gofakeit                            │
│  • UUID generation, table number calculation                 │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│      Repository Layer (internal/repository)                   │
│  • DmUserRepository.InsertDmUsersBatch()                    │
│  • DmPostRepository.InsertDmPostsBatch()                    │
│  • DmNewsRepository.InsertDmNewsBatch()                     │
│  • Batch insert processing (500 records at a time)           │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│           DB Layer (internal/db)                              │
│  • GroupManager                                             │
│  • Sharding connection management                            │
│  • Master connection management                              │
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

Provides an admin panel using GoAdmin. Operates on port 8081 as an independent service from the main service (port 8080).

### Architecture

Custom pages for the Admin app (dm_user_register.go, api_key.go) adopt the same layer structure as the API server (pages → usecase → service → repository → db).

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
│  │              Pages Layer (internal/admin/pages)         │ │
│  │  • Entry point                                          │ │
│  │  • Validation                                           │ │
│  │  • I/O control (HTML rendering)                         │ │
│  └──────────────────────┬─────────────────────────────────┘ │
│                         │                                     │
│  ┌──────────────────────▼─────────────────────────────────┐ │
│  │           Usecase Layer (internal/usecase/admin)        │ │
│  │  • DmUserRegisterUsecase: User registration logic       │ │
│  │  • APIKeyUsecase: API key issuance logic               │ │
│  └──────────────────────┬─────────────────────────────────┘ │
│                         │                                     │
│  ┌──────────────────────▼─────────────────────────────────┐ │
│  │              Service Layer (internal/service)           │ │
│  │  • DmUserService: User domain logic                    │ │
│  │  • APIKeyService: API key domain logic                 │ │
│  └──────────────────────┬─────────────────────────────────┘ │
│                         │                                     │
│  ┌──────────────────────▼─────────────────────────────────┐ │
│  │           Repository Layer (internal/repository)        │ │
│  │  • DmUserRepository: User data access                  │ │
│  └──────────────────────┬─────────────────────────────────┘ │
│                         │                                     │
│  ┌──────────────────────▼─────────────────────────────────┐ │
│  │                 DB Layer (internal/db)                  │ │
│  │  • GroupManager: Sharding connection management        │ │
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

- `config.go`: GoAdmin configuration struct
- `tables.go`: Table generators (Users, Posts)
- `sharding.go`: Cross-shard query helpers
- `auth/`: Authentication and session management
- `pages/`: Custom pages (dashboard, user registration, API key issuance)

**Location**: `internal/usecase/admin/`

- `dm_user_register_usecase.go`: User registration business logic
- `api_key_usecase.go`: API key issuance business logic

**Location**: `internal/service/`

- `api_key_service.go`: API key issuance domain logic

### Features

1. **Table Management**: List display and CRUD operations for Users/Posts tables
2. **Sharding Support**: Integrated display of data from all shards
3. **Authentication & Authorization**: GoAdmin built-in authentication
4. **Custom Pages**: Dashboard, user registration form, API key issuance

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
- `gorm.io/plugin/dbresolver`: GORM Writer/Reader separation plugin
- `gorm.io/sharding`: GORM sharding plugin (for future use)

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
