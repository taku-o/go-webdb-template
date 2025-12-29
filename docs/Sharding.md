# Database Sharding Documentation

## Overview

This project implements database sharding to distribute data across multiple database instances. The sharding architecture consists of two database groups:

- **Master Group**: Contains shared tables (e.g., `dm_news`) that don't require sharding
- **Sharding Group**: Contains partitioned tables (e.g., `dm_users_000` to `dm_users_031`, `dm_posts_000` to `dm_posts_031`)

## Architecture

### Database Groups

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           GroupManager                                   │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────────────┐    ┌─────────────────────────────────────────────┐ │
│  │  MasterManager  │    │              ShardingManager                │ │
│  │                 │    │                                             │ │
│  │  ┌───────────┐  │    │  8 Sharding Entries → 4 Databases           │ │
│  │  │ master.db │  │    │  (Connection Sharing for same DSN)         │ │
│  │  │(dm_news) │  │    │                                             │ │
│  │  └───────────┘  │    │  ┌──────────────────┐ ┌──────────────────┐ │ │
│  └─────────────────┘    │  │ DB 1             │ │ DB 2             │ │ │
│                         │  │ Entry 1: _000-003│ │ Entry 3: _008-011│ │ │
│                         │  │ Entry 2: _004-007│ │ Entry 4: _012-015│ │ │
│                         │  └──────────────────┘ └──────────────────┘ │ │
│                         │                                             │ │
│                         │  ┌──────────────────┐ ┌──────────────────┐ │ │
│                         │  │ DB 3             │ │ DB 4             │ │ │
│                         │  │ Entry 5: _016-019│ │ Entry 7: _024-027│ │ │
│                         │  │ Entry 6: _020-023│ │ Entry 8: _028-031│ │ │
│                         │  └──────────────────┘ └──────────────────┘ │ │
│                         └─────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────┘
```

### Sharding Entries and Connection Sharing

The system uses **8 sharding entries** distributed across **4 physical databases**:

| Entry ID | Table Range | Database | Connection |
|----------|-------------|----------|------------|
| Entry 1 | _000 〜 _003 | sharding_db_1 | Shared (conn1) |
| Entry 2 | _004 〜 _007 | sharding_db_1 | Shared (conn1) |
| Entry 3 | _008 〜 _011 | sharding_db_2 | Shared (conn2) |
| Entry 4 | _012 〜 _015 | sharding_db_2 | Shared (conn2) |
| Entry 5 | _016 〜 _019 | sharding_db_3 | Shared (conn3) |
| Entry 6 | _020 〜 _023 | sharding_db_3 | Shared (conn3) |
| Entry 7 | _024 〜 _027 | sharding_db_4 | Shared (conn4) |
| Entry 8 | _028 〜 _031 | sharding_db_4 | Shared (conn4) |

**Key Points**:
- 8 logical entries for fine-grained control
- 4 physical databases for actual storage
- Entries with same DSN share connections (connection pooling)
- Each entry handles 4 tables (for future scalability)

### Table Distribution

| Database | Entries | Table Range | Tables |
|----------|---------|-------------|--------|
| sharding_db_1 | 1, 2 | _000 〜 _007 | dm_users_000, dm_users_001, ..., dm_users_007 |
| sharding_db_2 | 3, 4 | _008 〜 _015 | dm_users_008, dm_users_009, ..., dm_users_015 |
| sharding_db_3 | 5, 6 | _016 〜 _023 | dm_users_016, dm_users_017, ..., dm_users_023 |
| sharding_db_4 | 7, 8 | _024 〜 _031 | dm_users_024, dm_users_025, ..., dm_users_031 |

## Sharding Strategy

### Table-Based Sharding

The application uses **table-based sharding** with 32 table partitions and 8 sharding entries.

**Algorithm**:
```go
// Table number calculation
tableNumber := id % 32  // Range: 0-31

// Table name generation
tableName := fmt.Sprintf("dm_users_%03d", tableNumber)  // e.g., "dm_users_005"

// Connection selection (O(1) lookup via tableNumberToDBID map)
conn := shardingManager.GetConnectionByTableNumber(tableNumber)
```

**Key Points**:
- Uses modulo 32 for even distribution
- Table number range: 0 to 31
- Same `id` always maps to the same table
- Posts use `user_id` as the sharding key
- 8 sharding entries, each handling 4 tables
- 4 physical databases, each containing 8 tables
- Connection sharing: entries with same DSN share connections

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
        # Entry 1,2 → sharding_db_1.db (tables _000-007)
        - id: 1
          driver: sqlite3
          dsn: ./data/sharding_db_1.db
          table_range: [0, 3]
        - id: 2
          driver: sqlite3
          dsn: ./data/sharding_db_1.db  # Same DSN = shared connection
          table_range: [4, 7]
        # Entry 3,4 → sharding_db_2.db (tables _008-015)
        - id: 3
          driver: sqlite3
          dsn: ./data/sharding_db_2.db
          table_range: [8, 11]
        - id: 4
          driver: sqlite3
          dsn: ./data/sharding_db_2.db  # Same DSN = shared connection
          table_range: [12, 15]
        # Entry 5,6 → sharding_db_3.db (tables _016-023)
        - id: 5
          driver: sqlite3
          dsn: ./data/sharding_db_3.db
          table_range: [16, 19]
        - id: 6
          driver: sqlite3
          dsn: ./data/sharding_db_3.db  # Same DSN = shared connection
          table_range: [20, 23]
        # Entry 7,8 → sharding_db_4.db (tables _024-031)
        - id: 7
          driver: sqlite3
          dsn: ./data/sharding_db_4.db
          table_range: [24, 27]
        - id: 8
          driver: sqlite3
          dsn: ./data/sharding_db_4.db  # Same DSN = shared connection
          table_range: [28, 31]

      tables:
        - name: dm_users
          suffix_count: 32
        - name: dm_posts
          suffix_count: 32
```

**Connection Sharing**: Entries 1 & 2 share the same DSN, so they share the same database connection. This reduces connection overhead while maintaining logical separation for future scalability.

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
        # Entry 1,2 → sharding1 (tables _000-007)
        - id: 1
          driver: postgres
          host: prod-db-sharding1.example.com
          name: app_sharding1
          password: ${DB_PASSWORD_SHARDING1}
          table_range: [0, 3]
        - id: 2
          driver: postgres
          host: prod-db-sharding1.example.com  # Same host = shared connection
          name: app_sharding1
          password: ${DB_PASSWORD_SHARDING1}
          table_range: [4, 7]
        # Entry 3,4 → sharding2 (tables _008-015)
        - id: 3
          driver: postgres
          host: prod-db-sharding2.example.com
          name: app_sharding2
          password: ${DB_PASSWORD_SHARDING2}
          table_range: [8, 11]
        - id: 4
          driver: postgres
          host: prod-db-sharding2.example.com  # Same host = shared connection
          name: app_sharding2
          password: ${DB_PASSWORD_SHARDING2}
          table_range: [12, 15]
        # Entry 5,6 → sharding3 (tables _016-023)
        - id: 5
          driver: postgres
          host: prod-db-sharding3.example.com
          name: app_sharding3
          password: ${DB_PASSWORD_SHARDING3}
          table_range: [16, 19]
        - id: 6
          driver: postgres
          host: prod-db-sharding3.example.com  # Same host = shared connection
          name: app_sharding3
          password: ${DB_PASSWORD_SHARDING3}
          table_range: [20, 23]
        # Entry 7,8 → sharding4 (tables _024-031)
        - id: 7
          driver: postgres
          host: prod-db-sharding4.example.com
          name: app_sharding4
          password: ${DB_PASSWORD_SHARDING4}
          table_range: [24, 27]
        - id: 8
          driver: postgres
          host: prod-db-sharding4.example.com  # Same host = shared connection
          name: app_sharding4
          password: ${DB_PASSWORD_SHARDING4}
          table_range: [28, 31]

      tables:
        - name: dm_users
          suffix_count: 32
        - name: dm_posts
          suffix_count: 32
```

## Data Distribution

### Master Group Tables

The master group contains shared tables that don't require sharding:

**dm_news table**:
```sql
CREATE TABLE IF NOT EXISTS dm_news (
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
User ID 1  → 1 % 32 = 1  → dm_users_001, DB 1
User ID 8  → 8 % 32 = 8  → dm_users_008, DB 2
User ID 16 → 16 % 32 = 16 → dm_users_016, DB 3
User ID 24 → 24 % 32 = 24 → dm_users_024, DB 4
User ID 32 → 32 % 32 = 0  → dm_users_000, DB 1
User ID 100 → 100 % 32 = 4 → dm_users_004, DB 1
```

**Benefit**: Even distribution across tables and databases.

## Query Patterns

### Single-Table Queries

Operations that access data for a single entity use dynamic table names:

**Example**: Get User by ID
```go
// Calculate table name
tableName := db.GetShardingTableName("dm_users", userID)  // e.g., "dm_users_005"

// Get connection for this table
conn, err := groupManager.GetShardingConnectionByID(userID, "dm_users")
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
            tableName := fmt.Sprintf("dm_users_%03d", i)

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
            usersTable := fmt.Sprintf("dm_users_%03d", i)
            postsTable := fmt.Sprintf("dm_posts_%03d", i)

            // JOIN within same suffix (dm_users_005 with dm_posts_005)
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
        │   ├── dm_users.sql.template    # dm_users table template
        │   └── dm_posts.sql.template    # dm_posts table template
        └── generated/              # generated migrations (optional)
```

### Migration Templates

**dm_users.sql.template**:
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

**dm_posts.sql.template**:
```sql
CREATE TABLE IF NOT EXISTS {TABLE_NAME} (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_{TABLE_NAME}_user_id ON {TABLE_NAME}(user_id);
CREATE INDEX IF NOT EXISTS idx_{TABLE_NAME}_created_at ON {TABLE_NAME}(created_at);
```

### Migration Generation Tool

Generate migration files from templates:

```bash
cd server
go run cmd/migrate-gen/main.go \
    -template db/migrations/sharding/templates/dm_users.sql.template \
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
    sqlite3 server/data/sharding_db_1.db < db/migrations/sharding/generated/dm_users_$(printf "%03d" $i).sql
    sqlite3 server/data/sharding_db_1.db < db/migrations/sharding/generated/dm_posts_$(printf "%03d" $i).sql
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

// Get master connection (for dm_news table)
masterConn, err := groupManager.GetMasterConnection()

// Get sharding connection by ID
shardingConn, err := groupManager.GetShardingConnectionByID(userID, "dm_users")

// Get sharding connection by table number
shardingConn, err := groupManager.GetShardingConnection(tableNumber)

// Get all sharding connections (for cross-table queries)
connections := groupManager.GetAllShardingConnections()
```

### TableSelector

The `TableSelector` handles table name generation:

```go
tableSelector := db.NewTableSelector(32, 8)  // 32 tables, 8 sharding entries

// Get table number from ID
tableNumber := tableSelector.GetTableNumber(userID)  // e.g., 5

// Get table name
tableName := tableSelector.GetTableName("dm_users", userID)  // e.g., "dm_users_005"

// Get entry ID from table number (uses 4 tables per entry with 8 entries)
entryID := tableSelector.GetDBID(tableNumber)  // e.g., 2 for table 5
```

### ShardingManager Internal Architecture

The `ShardingManager` uses several internal maps for efficient connection management:

```go
type ShardingManager struct {
    connections       map[int]*GORMConnection  // Entry ID → Connection
    tableRange        map[int][2]int           // Entry ID → [min, max] table numbers
    connectionPool    map[string]*GORMConnection  // DSN → Connection (shared)
    tableNumberToDBID map[int]int              // Table number → Entry ID (O(1) lookup)
    mu                sync.RWMutex
}
```

**Connection Sharing**:
- Entries with the same DSN share the same connection object
- `connectionPool` stores unique connections by DSN
- `GetAllConnections()` returns only unique connections (4 connections for 8 entries)

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
tableName := db.GetShardingTableName("dm_users", userID)

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
SELECT 'dm_users_000' as table_name, COUNT(*) as count FROM dm_users_000
UNION ALL
SELECT 'dm_users_001', COUNT(*) FROM dm_users_001
-- ... for all tables
```

## Limitations

### Current Limitations

1. **No Distributed Transactions**: Transactions cannot span multiple databases
2. **Table Suffix Required**: Post operations require `user_id` to determine table
3. **Limited Range Queries**: Queries like "get users 1-100" require checking multiple tables
4. **32-Table Fixed**: Currently fixed at 32 table partitions

### Future Scalability

The 8-entry configuration allows for future database expansion:
- Currently: 8 entries → 4 databases (2 entries per DB share connections)
- Future: 8 entries → 8 databases (1 entry per DB, separate connections)
- Requires only configuration change, no code modification

### Backward Compatibility

The system maintains backward compatibility:
- Old `shards` configuration is still supported (deprecated)
- `GORMManager` is still available (deprecated)
- Migration path: Use `GroupManager` for new code
