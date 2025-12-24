# Database Sharding Documentation

## Overview

This project implements database sharding to distribute data across multiple database instances. Sharding improves scalability by allowing the system to handle more data and traffic than a single database server could support.

## Sharding Strategy

### Hash-Based Sharding

The application uses **hash-based sharding** with `user_id` as the shard key.

**Algorithm**:
```go
func (h *HashBasedSharding) GetShardID(key int64) int {
    hash := fnv.New32a()
    hash.Write([]byte(fmt.Sprintf("%d", key)))
    hashValue := hash.Sum32()
    shardID := int(hashValue%uint32(h.shardCount)) + 1
    return shardID
}
```

**Key Points**:
- Uses FNV-1a hash function for consistent distribution
- Shard ID range: 1 to N (where N is the shard count)
- Same `user_id` always maps to the same shard
- Posts are co-located with users (same shard)

### Why Hash-Based Sharding?

**Advantages**:
- Even distribution of data across shards
- Simple and predictable shard selection
- No hotspots if user IDs are well-distributed
- Deterministic: same key always goes to same shard

**Disadvantages**:
- Difficult to rebalance when adding/removing shards
- Range queries across shards are expensive
- Related data must share the same shard key

## Shard Configuration

### Development Environment

**File**: `config/develop.yaml`

```yaml
database:
  shards:
    - id: 1
      driver: sqlite3
      host: ./data/shard1.db
    - id: 2
      driver: sqlite3
      host: ./data/shard2.db
```

### Production Environment

**File**: `config/production.yaml`

```yaml
database:
  shards:
    - id: 1
      driver: postgres
      host: db-shard1.example.com
      port: 5432
      name: app_shard1
      user: app_user
      password: ${DB_SHARD1_PASSWORD}
    - id: 2
      driver: postgres
      host: db-shard2.example.com
      port: 5432
      name: app_shard2
      user: app_user
      password: ${DB_SHARD2_PASSWORD}
```

## Data Distribution

### User Table

Users are distributed across shards based on their ID:

```
User ID 1 → hash(1) % 2 = Shard 1 or 2
User ID 2 → hash(2) % 2 = Shard 1 or 2
...
```

Example distribution (with 2 shards):
- User ID 1 → Shard 2
- User ID 2 → Shard 1
- User ID 3 → Shard 2
- User ID 4 → Shard 1

### Post Table

Posts are co-located with users using `user_id` as the shard key:

```
Post with user_id=1 → Same shard as User 1
Post with user_id=2 → Same shard as User 2
```

**Benefit**: User and their posts are always on the same shard, enabling efficient single-shard queries.

## Query Patterns

### Single-Shard Queries

Operations that access data for a single user can be routed to a specific shard:

**Example**: Get User by ID
```go
// Automatically routes to correct shard
conn, err := dbManager.GetConnectionByKey(userID)
if err != nil {
    return nil, err
}

user, err := repo.GetByID(conn, userID)
```

**Queries**:
- `GetUserByID(id)` → Single shard
- `CreateUser(user)` → Single shard
- `GetPostsByUserID(userID)` → Single shard
- `UpdateUser(id, data)` → Single shard
- `DeleteUser(id)` → Single shard

### Cross-Shard Queries

Operations that need to access all data must query multiple shards:

**Example**: Get All Users
```go
func (r *UserRepository) GetAll() ([]*model.User, error) {
    var allUsers []*model.User

    // Query each shard
    for i := 1; i <= r.dbManager.GetShardCount(); i++ {
        conn, err := r.dbManager.GetConnection(i)
        if err != nil {
            return nil, err
        }

        // Query this shard
        users, err := queryUsersFromShard(conn)
        if err != nil {
            return nil, err
        }

        allUsers = append(allUsers, users...)
    }

    return allUsers, nil
}
```

**Queries**:
- `GetAllUsers()` → All shards
- `GetAllPosts()` → All shards
- `GetUserPosts()` (JOIN) → All shards

### Cross-Shard JOIN

**Challenge**: Users and Posts may be on different shards

**Solution**: Application-level JOIN

```go
func (r *PostRepository) GetUserPosts() ([]UserPost, error) {
    var result []UserPost

    // 1. Get all users from all shards
    users, err := r.userRepo.GetAll()

    // 2. Get all posts from all shards
    posts, err := r.GetAll()

    // 3. Join in memory
    for _, post := range posts {
        for _, user := range users {
            if post.UserID == user.ID {
                result = append(result, UserPost{
                    UserID:      user.ID,
                    UserName:    user.Name,
                    UserEmail:   user.Email,
                    PostID:      post.ID,
                    PostTitle:   post.Title,
                    PostContent: post.Content,
                    CreatedAt:   post.CreatedAt,
                })
                break
            }
        }
    }

    return result, nil
}
```

**Performance Note**: Application-level JOINs can be expensive for large datasets. Consider:
- Caching frequently accessed data
- Denormalizing data to avoid JOINs
- Using pagination to limit result size

## Schema Management

### Schema Consistency

Each shard must have an identical schema:

**Shard 1**: `db/migrations/shard1/001_init.sql`
**Shard 2**: `db/migrations/shard2/001_init.sql`

Both contain the same schema:
```sql
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### Schema Migrations

When evolving the schema:

1. Create migration files for each shard
2. Apply migrations to all shards
3. Ensure consistency across all shards
4. Test cross-shard queries after migration

**Future Enhancement**: Automated migration tool that applies changes to all shards.

## Adding/Removing Shards

### Adding a Shard

**Warning**: Adding shards to a hash-based sharding system requires data rebalancing.

**Steps**:
1. Add new shard configuration to config file
2. Create schema on new shard
3. Recalculate shard assignments for all users
4. Migrate data to new shards
5. Update application with new shard count

**Challenge**: Existing data must be redistributed because hash values will change.

### Alternative: Consistent Hashing

For easier shard management, consider implementing **consistent hashing**:
- Minimizes data movement when adding/removing shards
- Only ~1/N of data needs to move when adding a shard
- More complex implementation

## Performance Optimization

### Connection Pooling

Each shard maintains its own connection pool:

```go
type Connection struct {
    db       *sql.DB
    shardID  int
    config   ShardConfig
}

// Configure connection pool
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
```

### Parallel Queries

Cross-shard queries execute in parallel using goroutines:

```go
var wg sync.WaitGroup
resultChan := make(chan []*model.User, shardCount)

for i := 1; i <= shardCount; i++ {
    wg.Add(1)
    go func(shardID int) {
        defer wg.Done()
        users, _ := queryShardParallel(shardID)
        resultChan <- users
    }(i)
}

wg.Wait()
close(resultChan)

// Collect results
for users := range resultChan {
    allUsers = append(allUsers, users...)
}
```

### Indexes

Ensure appropriate indexes on each shard:

```sql
-- User lookups
CREATE INDEX idx_users_email ON users(email);

-- Post lookups by user
CREATE INDEX idx_posts_user_id ON posts(user_id);

-- Recent posts
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);
```

## Monitoring and Metrics

### Key Metrics to Monitor

1. **Shard Distribution**:
   - Number of users per shard
   - Number of posts per shard
   - Ensure even distribution

2. **Query Performance**:
   - Single-shard query latency
   - Cross-shard query latency
   - Connection pool utilization

3. **Shard Health**:
   - Connection errors per shard
   - Query errors per shard
   - Available connections per shard

### Monitoring Queries

```sql
-- Check user distribution
SELECT 'Shard 1' as shard, COUNT(*) as user_count FROM users;
-- Run on each shard

-- Check post distribution
SELECT 'Shard 1' as shard, COUNT(*) as post_count FROM posts;
-- Run on each shard
```

## Best Practices

### 1. Always Use Shard Keys

When designing queries, always include the shard key (`user_id`) when possible:

```go
// Good: Includes user_id for routing
DeletePost(postID, userID)

// Bad: Would require querying all shards
DeletePost(postID)
```

### 2. Co-locate Related Data

Keep related entities on the same shard:
- Users and their Posts share the same shard key
- Comments could use the same user_id shard key
- User settings could use user_id as well

### 3. Denormalize When Necessary

If you frequently JOIN across shards:
- Consider denormalizing data
- Store user name with posts to avoid JOINs
- Use caching for frequently accessed data

### 4. Plan for Growth

- Monitor shard size and distribution
- Plan for rebalancing before reaching capacity
- Consider partitioning strategies for very large tables

### 5. Test Cross-Shard Operations

Thoroughly test:
- Cross-shard queries return correct results
- Transactions don't span shards (not supported)
- Error handling when shards are unavailable

## Limitations

### Current Limitations

1. **No Distributed Transactions**: Transactions cannot span multiple shards
2. **No Range Queries**: Queries like "get users 1-100" require checking all shards
3. **Rebalancing is Hard**: Adding/removing shards requires data migration
4. **JOIN Performance**: Application-level JOINs can be slow for large datasets

### Future Enhancements

1. **Global ID Generation**: Use a distributed ID generator (Snowflake, UUID)
2. **Consistent Hashing**: Easier shard management
3. **Read Replicas**: Add read replicas for each shard (GORM版で実装済み)
4. **Shard Rebalancing Tool**: Automate data migration
5. **Query Router**: Dedicated routing layer for shard selection

## GORM Support

### Writer/Reader Separation

GORM版のRepositoryは`gorm.io/plugin/dbresolver`を使用してWriter/Reader分離をサポートしています。

**設定例** (`config/production.yaml`):
```yaml
database:
  shards:
    - id: 1
      driver: postgres
      writer_dsn: host=prod-db-shard1-writer.example.com port=5432 user=prod_user password=${DB_PASSWORD_SHARD1} dbname=app_db_shard1 sslmode=require
      reader_dsns:
        - host=prod-db-shard1-reader1.example.com port=5432 user=prod_user password=${DB_PASSWORD_SHARD1} dbname=app_db_shard1 sslmode=require
        - host=prod-db-shard1-reader2.example.com port=5432 user=prod_user password=${DB_PASSWORD_SHARD1} dbname=app_db_shard1 sslmode=require
      reader_policy: round_robin
```

**機能**:
- **Write操作**: Create, Update, Delete → Writer DBに送信
- **Read操作**: Select, Find → Reader DBに送信
- **複数Reader**: 複数のReader DSNを設定可能
- **ロードバランシング**: `random`または`round_robin`ポリシー
- **後方互換性**: 従来の`dsn`設定も引き続きサポート

**使用例**:
```go
// GORMManager を使用
gormManager, err := db.NewGORMManager(cfg)
defer gormManager.CloseAll()

// Repository層での使用
userRepo := repository.NewUserRepositoryGORM(gormManager)

// Write操作 → Writer DBに送信
user, err := userRepo.Create(ctx, req)

// Read操作 → Reader DBに送信
user, err := userRepo.GetByID(ctx, id)
```

**利点**:
- Writeの負荷とReadの負荷を分離
- Read性能のスケールアウトが容易
- Readerの追加だけで読み込み性能を向上可能
- 既存のsharding戦略と組み合わせて使用可能

### GORM Migration Path

現在は`database/sql`版のRepositoryを使用していますが、GORM版も完全に実装されています。

**移行手順**:
1. Service層をInterface化
2. main.goでGORMManagerを使用
3. GORMRepositoryをインジェクト
4. Writer/Reader分離の設定を追加

**現状**:
- GORM版Repository: 実装済み ✅
- Writer/Reader分離: 実装済み ✅
- テスト: すべてパス ✅
- Service層Interface化: 未実装 (将来タスク)
