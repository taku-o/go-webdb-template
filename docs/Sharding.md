# Database Sharding Documentation

## Overview

This project implements database sharding to distribute data across multiple database instances. The sharding architecture consists of two database groups:

- **Master Group**: Contains shared tables (e.g., `news`) that don't require sharding
- **Sharding Group**: Contains partitioned tables (e.g., `users_000` to `users_031`, `posts_000` to `posts_031`)

## Architecture

### Database Groups

```
┌─────────────────────────────────────────────────────────────┐
│                      GroupManager                           │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────────┐    ┌─────────────────────────────────┐ │
│  │  MasterManager  │    │       ShardingManager          │ │
│  │                 │    │                                 │ │
│  │  ┌───────────┐  │    │  ┌─────────┐  ┌─────────┐     │ │
│  │  │ master.db │  │    │  │  DB 1   │  │  DB 2   │     │ │
│  │  │ (news)    │  │    │  │ _000-007│  │ _008-015│     │ │
│  │  └───────────┘  │    │  └─────────┘  └─────────┘     │ │
│  └─────────────────┘    │                                 │ │
│                         │  ┌─────────┐  ┌─────────┐     │ │
│                         │  │  DB 3   │  │  DB 4   │     │ │
│                         │  │ _016-023│  │ _024-031│     │ │
│                         │  └─────────┘  └─────────┘     │ │
│                         └─────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### Table Distribution

| Database | Table Range | Tables |
|----------|------------|--------|
| sharding_db_1 | _000 〜 _007 | users_000, users_001, ..., users_007 |
| sharding_db_2 | _008 〜 _015 | users_008, users_009, ..., users_015 |
| sharding_db_3 | _016 〜 _023 | users_016, users_017, ..., users_023 |
| sharding_db_4 | _024 〜 _031 | users_024, users_025, ..., users_031 |

## Sharding Strategy

### Table-Based Sharding

The application uses **table-based sharding** with 32 table partitions.

**Algorithm**:
```go
// Table number calculation
tableNumber := id % 32  // Range: 0-31

// Table name generation
tableName := fmt.Sprintf("users_%03d", tableNumber)  // e.g., "users_005"

// Database selection
dbID := (tableNumber / 8) + 1  // Range: 1-4
```

**Key Points**:
- Uses modulo 32 for even distribution
- Table number range: 0 to 31
- Same `id` always maps to the same table
- Posts use `user_id` as the sharding key
- Each database holds 8 tables

### Why Table-Based Sharding?

**Advantages**:
- Even distribution of data across tables
- Simple and predictable table selection (O(1))
- No hotspots if IDs are well-distributed
- Easier to add more databases without data migration
- Flexible: can move tables between databases

**Disadvantages**:
- More tables to manage (32 tables per entity)
- Range queries require checking multiple tables
- Template-based migrations needed

## Configuration

### Development Environment

**File**: `config/develop/database.yaml`

```yaml
database:
  groups:
    master:
      - id: 1
        driver: sqlite3
        dsn: ./data/master.db
        writer_dsn: ./data/master.db
        reader_dsns:
          - ./data/master.db
        reader_policy: random
        max_connections: 10
        max_idle_connections: 5
        connection_max_lifetime: 300s

    sharding:
      databases:
        - id: 1
          driver: sqlite3
          dsn: ./data/sharding_db_1.db
          writer_dsn: ./data/sharding_db_1.db
          reader_dsns:
            - ./data/sharding_db_1.db
          reader_policy: random
          max_connections: 10
          max_idle_connections: 5
          connection_max_lifetime: 300s
          table_range: [0, 7]  # _000-007
        - id: 2
          driver: sqlite3
          dsn: ./data/sharding_db_2.db
          writer_dsn: ./data/sharding_db_2.db
          reader_dsns:
            - ./data/sharding_db_2.db
          reader_policy: random
          max_connections: 10
          max_idle_connections: 5
          connection_max_lifetime: 300s
          table_range: [8, 15]  # _008-015
        - id: 3
          driver: sqlite3
          dsn: ./data/sharding_db_3.db
          writer_dsn: ./data/sharding_db_3.db
          reader_dsns:
            - ./data/sharding_db_3.db
          reader_policy: random
          max_connections: 10
          max_idle_connections: 5
          connection_max_lifetime: 300s
          table_range: [16, 23]  # _016-023
        - id: 4
          driver: sqlite3
          dsn: ./data/sharding_db_4.db
          writer_dsn: ./data/sharding_db_4.db
          reader_dsns:
            - ./data/sharding_db_4.db
          reader_policy: random
          max_connections: 10
          max_idle_connections: 5
          connection_max_lifetime: 300s
          table_range: [24, 31]  # _024-031

      tables:
        - name: users
          suffix_count: 32
        - name: posts
          suffix_count: 32
```

### Production Environment

**File**: `config/production/database.yaml.example`

```yaml
database:
  groups:
    master:
      - id: 1
        driver: postgres
        host: prod-db-master.example.com
        port: 5432
        name: app_master
        user: prod_user
        password: ${DB_PASSWORD_MASTER}
        writer_dsn: host=prod-db-master-writer.example.com ...
        reader_dsns:
          - host=prod-db-master-reader1.example.com ...
        reader_policy: round_robin

    sharding:
      databases:
        - id: 1
          driver: postgres
          host: prod-db-shard1.example.com
          port: 5432
          name: app_shard1
          user: prod_user
          password: ${DB_PASSWORD_SHARD1}
          table_range: [0, 7]
        - id: 2
          driver: postgres
          host: prod-db-shard2.example.com
          port: 5432
          name: app_shard2
          user: prod_user
          password: ${DB_PASSWORD_SHARD2}
          table_range: [8, 15]
        - id: 3
          driver: postgres
          host: prod-db-shard3.example.com
          port: 5432
          name: app_shard3
          user: prod_user
          password: ${DB_PASSWORD_SHARD3}
          table_range: [16, 23]
        - id: 4
          driver: postgres
          host: prod-db-shard4.example.com
          port: 5432
          name: app_shard4
          user: prod_user
          password: ${DB_PASSWORD_SHARD4}
          table_range: [24, 31]

      tables:
        - name: users
          suffix_count: 32
        - name: posts
          suffix_count: 32
```

## Data Distribution

### Master Group Tables

The master group contains shared tables that don't require sharding:

**news table**:
```sql
CREATE TABLE IF NOT EXISTS news (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    author_id INTEGER,
    published_at DATETIME,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);
```

### Sharding Group Tables

Users and posts are distributed across 32 tables:

```
User ID 1  → 1 % 32 = 1  → users_001, DB 1
User ID 8  → 8 % 32 = 8  → users_008, DB 2
User ID 16 → 16 % 32 = 16 → users_016, DB 3
User ID 24 → 24 % 32 = 24 → users_024, DB 4
User ID 32 → 32 % 32 = 0  → users_000, DB 1
User ID 100 → 100 % 32 = 4 → users_004, DB 1
```

**Benefit**: Even distribution across tables and databases.

## Query Patterns

### Single-Table Queries

Operations that access data for a single entity use dynamic table names:

**Example**: Get User by ID
```go
// Calculate table name
tableName := db.GetShardingTableName("users", userID)  // e.g., "users_005"

// Get connection for this table
conn, err := groupManager.GetShardingConnectionByID(userID, "users")
if err != nil {
    return nil, err
}

// Query with dynamic table name
var user model.User
err = conn.DB.Table(tableName).Where("id = ?", userID).First(&user).Error
```

**Queries**:
- `GetUserByID(id)` → Single table
- `CreateUser(user)` → Single table
- `GetPostByID(postID, userID)` → Single table
- `UpdateUser(id, data)` → Single table
- `DeleteUser(id)` → Single table

### Cross-Table Queries

Operations that need to access all data must query all tables:

**Example**: Get All Users
```go
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
    connections := r.groupManager.GetAllShardingConnections()
    var allUsers []*model.User

    // Query each database
    for _, conn := range connections {
        // Query all tables in this database
        for i := conn.TableRange[0]; i <= conn.TableRange[1]; i++ {
            tableName := fmt.Sprintf("users_%03d", i)

            var users []*model.User
            err := conn.DB.Table(tableName).
                Order("id").
                Limit(limit).
                Offset(offset).
                Find(&users).Error
            if err != nil {
                continue
            }
            allUsers = append(allUsers, users...)
        }
    }

    return allUsers, nil
}
```

**Queries**:
- `GetAllUsers()` → All tables in all databases
- `GetAllPosts()` → All tables in all databases
- `GetUserPosts()` (JOIN) → All tables in all databases

### Cross-Table JOIN

**Challenge**: Users and Posts are in different tables

**Solution**: Application-level JOIN with table matching

```go
func (r *PostRepository) GetUserPosts(ctx context.Context, limit, offset int) ([]model.UserPost, error) {
    var result []model.UserPost

    connections := r.groupManager.GetAllShardingConnections()

    for _, conn := range connections {
        for i := conn.TableRange[0]; i <= conn.TableRange[1]; i++ {
            usersTable := fmt.Sprintf("users_%03d", i)
            postsTable := fmt.Sprintf("posts_%03d", i)

            // JOIN within same suffix (users_005 with posts_005)
            var userPosts []model.UserPost
            err := conn.DB.Table(postsTable).
                Select("users.id as user_id, users.name as user_name, ...").
                Joins(fmt.Sprintf("JOIN %s ON %s.user_id = %s.id",
                    usersTable, postsTable, usersTable)).
                Find(&userPosts).Error
            if err != nil {
                continue
            }
            result = append(result, userPosts...)
        }
    }

    return result, nil
}
```

**Note**: Users and their posts share the same table suffix (same sharding key), enabling efficient single-database JOINs.

## Migration Management

### Directory Structure

```
db/
└── migrations/
    ├── master/
    │   └── 001_init.sql          # news table
    └── sharding/
        ├── templates/
        │   ├── users.sql.template    # users table template
        │   └── posts.sql.template    # posts table template
        └── generated/              # generated migrations (optional)
```

### Migration Templates

**users.sql.template**:
```sql
CREATE TABLE IF NOT EXISTS {TABLE_NAME} (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_{TABLE_NAME}_email ON {TABLE_NAME}(email);
```

**posts.sql.template**:
```sql
CREATE TABLE IF NOT EXISTS {TABLE_NAME} (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users_{TABLE_SUFFIX}(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_{TABLE_NAME}_user_id ON {TABLE_NAME}(user_id);
CREATE INDEX IF NOT EXISTS idx_{TABLE_NAME}_created_at ON {TABLE_NAME}(created_at);
```

### Migration Generation Tool

Generate migration files from templates:

```bash
cd server
go run cmd/migrate-gen/main.go \
    -template db/migrations/sharding/templates/users.sql.template \
    -output db/migrations/sharding/generated/
```

### Migration Application

Apply migrations using the script:

```bash
# Apply all migrations
./scripts/migrate.sh

# Or apply manually:

# Master database
sqlite3 server/data/master.db < db/migrations/master/001_init.sql

# Sharding databases (example for DB 1, tables _000-007)
for i in {0..7}; do
    sqlite3 server/data/sharding_db_1.db < db/migrations/sharding/generated/users_$(printf "%03d" $i).sql
    sqlite3 server/data/sharding_db_1.db < db/migrations/sharding/generated/posts_$(printf "%03d" $i).sql
done
```

## Connection Management

### GroupManager

The `GroupManager` provides unified access to all database connections:

```go
// Initialize GroupManager
groupManager, err := db.NewGroupManager(cfg)
if err != nil {
    log.Fatal(err)
}
defer groupManager.CloseAll()

// Get master connection (for news table)
masterConn, err := groupManager.GetMasterConnection()

// Get sharding connection by ID
shardingConn, err := groupManager.GetShardingConnectionByID(userID, "users")

// Get sharding connection by table number
shardingConn, err := groupManager.GetShardingConnection(tableNumber)

// Get all sharding connections (for cross-table queries)
connections := groupManager.GetAllShardingConnections()
```

### TableSelector

The `TableSelector` handles table name generation:

```go
tableSelector := db.NewTableSelector(32, 8)  // 32 tables, 8 per database

// Get table number from ID
tableNumber := tableSelector.GetTableNumber(userID)  // e.g., 5

// Get table name
tableName := tableSelector.GetTableName("users", userID)  // e.g., "users_005"

// Get database ID from table number
dbID := tableSelector.GetDBID(tableNumber)  // e.g., 1
```

## GORM Support

### Writer/Reader Separation

GORM repositories use `gorm.io/plugin/dbresolver` for Writer/Reader separation.

**Configuration**:
```yaml
database:
  groups:
    sharding:
      databases:
        - id: 1
          driver: postgres
          writer_dsn: host=prod-db-shard1-writer.example.com ...
          reader_dsns:
            - host=prod-db-shard1-reader1.example.com ...
            - host=prod-db-shard1-reader2.example.com ...
          reader_policy: round_robin  # or "random"
```

**Behavior**:
- **Write operations**: Create, Update, Delete → Writer DB
- **Read operations**: Select, Find → Reader DB
- **Load balancing**: Supports `random` or `round_robin` policies

### Dynamic Table Names with GORM

```go
// Use Table() method for dynamic table names
tableName := db.GetShardingTableName("users", userID)

var user model.User
err := conn.DB.Table(tableName).Where("id = ?", userID).First(&user).Error

// For Create/Update/Delete
err := conn.DB.Table(tableName).Create(&user).Error
```

## Best Practices

### 1. Always Use Sharding Keys

When designing queries, always include the sharding key:

```go
// Good: Includes user_id for routing
GetPost(postID, userID)
DeletePost(postID, userID)

// Bad: Would require querying all tables
GetPost(postID)  // Which table contains this post?
```

### 2. Co-locate Related Data

Keep related entities using the same sharding key:
- Posts use `user_id` as the sharding key
- Same `user_id` → same table suffix → same database
- Enables efficient JOINs within a single database

### 3. Use Template-Based Migrations

Always use the template system for schema changes:
- Modify the template file
- Regenerate migration files
- Apply to all databases

### 4. Monitor Table Distribution

Check data distribution across tables:
```sql
-- Check user count per table (run in each database)
SELECT 'users_000' as table_name, COUNT(*) as count FROM users_000
UNION ALL
SELECT 'users_001', COUNT(*) FROM users_001
-- ... for all tables
```

## Limitations

### Current Limitations

1. **No Distributed Transactions**: Transactions cannot span multiple databases
2. **Table Suffix Required**: Post operations require `user_id` to determine table
3. **Limited Range Queries**: Queries like "get users 1-100" require checking multiple tables
4. **32-Table Fixed**: Currently fixed at 32 table partitions

### Backward Compatibility

The system maintains backward compatibility:
- Old `shards` configuration is still supported (deprecated)
- `GORMManager` is still available (deprecated)
- Migration path: Use `GroupManager` for new code
