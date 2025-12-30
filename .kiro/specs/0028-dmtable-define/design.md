# 分散テーブル環境対応テーブル設計修正設計書

## 1. 概要

### 1.1 設計の目的

要件定義書に基づき、Issue #52の対応として、分散テーブル環境に適したテーブル設計に修正する機能の詳細設計を定義する。auto_incrementによるID生成からsonyflakeによるID生成に変更し、分散環境で一意性が保証されるID生成方式を実装する。

### 1.2 設計の範囲

- sonyflakeライブラリの導入とID生成ユーティリティの実装
- テーブル定義の修正（Atlas形式）
- モデル定義の修正（GORMタグの修正）
- Repository層でのID生成統合
- サンプルデータ生成コマンドの修正
- sharding規則の定義
- テストの実装
- ドキュメントの更新

### 1.3 設計方針

- **分散環境対応**: sonyflakeを使用して分散環境で一意性が保証されるIDを生成する
- **既存ロジックの維持**: ID生成方式のみを変更し、既存のビジネスロジックは維持する
- **段階的実装**: 各コンポーネントを段階的に実装し、各段階でテストを実行する
- **一貫性の確保**: 全テーブルで統一されたID生成方式を採用する
- **互換性の維持**: 既存のAPIインターフェースは変更しない

## 2. アーキテクチャ設計

### 2.1 ID生成システムの構成

```
┌─────────────────────────────────────────────────────────┐
│              ID生成システム                              │
│                                                         │
│  ┌──────────────────────────────────────────────────┐  │
│  │  idgen.GenerateSonyflakeID()                              │  │
│  │  - sonyflakeインスタンスの管理                    │  │
│  │  - スレッドセーフなID生成                         │  │
│  └──────────────────────────────────────────────────┘  │
│                    │                                     │
│                    ▼                                     │
│  ┌──────────────────────────────────────────────────┐  │
│  │  github.com/sony/sonyflake                       │  │
│  │  - 分散環境で一意性が保証されるID生成             │  │
│  │  - 時系列順序が保たれる                           │  │
│  └──────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
                    │
        ┌───────────┴───────────┐
        ▼                       ▼
┌──────────────┐      ┌──────────────┐
│ Repository層 │      │ サンプルデータ│
│              │      │ 生成コマンド   │
│ Create()     │      │ generate*()   │
└──────────────┘      └──────────────┘
```

### 2.2 ID生成フロー

#### 2.2.1 Repository層でのID生成

```
1. Repository.Create() が呼び出される
   ↓
2. idgen.GenerateSonyflakeID() を呼び出し
   ↓
3. sonyflakeがIDを生成（int64）
   ↓
4. モデルのIDフィールドに設定
   ↓
5. データベースに保存
```

#### 2.2.2 サンプルデータ生成コマンドでのID生成

```
1. generateDmUsers() / generateDmPosts() / generateDmNews() が呼び出される
   ↓
2. 各エンティティ作成時に idgen.GenerateSonyflakeID() を呼び出し
   ↓
3. モデルのIDフィールドに設定
   ↓
4. バッチ挿入時にIDを含めてINSERT
```

### 2.3 Sharding規則の定義

#### 2.3.1 テーブルsharding規則

```
dm_users_NNN:
  - table_sharding_key: id
  - テーブル番号計算: id % DBShardingTableCount
  - 例: id = 123456789 → テーブル番号 = 123456789 % 32 = 13 → dm_users_013

dm_posts_NNN:
  - table_sharding_key: user_id
  - テーブル番号計算: user_id % DBShardingTableCount
  - 例: user_id = 123456789 → テーブル番号 = 123456789 % 32 = 13 → dm_posts_013
  - 重要な規則: dm_usersのIDとdm_postsのuser_idが同じ値であれば、同じテーブル番号になる
```

#### 2.3.2 Sharding規則の実装場所

- **ファイル**: `server/internal/db/sharding.go`
- **実装方法**: 定数または構造体として定義し、ドキュメントコメントを追加

## 3. データ構造設計

### 3.1 ID生成ユーティリティ

#### 3.1.1 パッケージ構造

```
server/internal/util/idgen/
  ├── sonyflake.go      # sonyflake ID生成関数の実装
  └── sonyflake_test.go # sonyflake ID生成の単体テスト
```

**注意**: パッケージ名は`idgen`のまま。将来的にUUIDv7用のファイル（`uuid.go`など）を追加することを想定。

#### 3.1.2 ID生成関数のシグネチャ

```go
package idgen

import (
    "github.com/sony/sonyflake"
)

var (
    sf *sonyflake.Sonyflake
    once sync.Once
)

// GenerateSonyflakeID はsonyflakeを使用して一意のIDを生成する
func GenerateSonyflakeID() (int64, error) {
    // sonyflakeインスタンスの初期化（初回のみ）
    // ID生成
    // エラーハンドリング
}
```

**注意**: パッケージ内でsonyflake用の関数として`GenerateSonyflakeID`を提供。将来的にUUIDv7用の関数（例: `GenerateUUIDv7()`）を追加することを想定。

### 3.2 テーブル定義の変更

#### 3.2.1 dm_newsテーブル（master.hcl）

```hcl
table "dm_news" {
  schema = schema.main
  column "id" {
    null           = false
    type           = bigint
    unsigned       = true
    auto_increment = false
  }
  // ... 他のカラム
}
```

#### 3.2.2 dm_users_NNNテーブル（sharding_*/dm_users.hcl）

```hcl
table "dm_users_000" {
  schema = schema.main
  column "id" {
    null           = false
    type           = bigint
    unsigned       = true
    auto_increment = false
  }
  // ... 他のカラム
}
// dm_users_001 〜 dm_users_031 も同様
```

#### 3.2.3 dm_posts_NNNテーブル（sharding_*/dm_posts.hcl）

```hcl
table "dm_posts_000" {
  schema = schema.main
  column "id" {
    null           = false
    type           = bigint
    unsigned       = true
    auto_increment = false
  }
  // ... 他のカラム
}
// dm_posts_001 〜 dm_posts_031 も同様
```

### 3.3 モデル定義の変更

#### 3.3.1 DmNewsモデル

```go
type DmNews struct {
    ID          int64      `json:"id,string" db:"id" gorm:"primaryKey"`  // autoIncrementを削除
    Title       string     `json:"title" db:"title" gorm:"type:varchar(255);not null"`
    // ... 他のフィールド
}
```

#### 3.3.2 DmUserモデル・DmPostモデル

- GORMタグの変更は不要（既に`autoIncrement`が設定されていない）
- ID生成ロジックの追加のみ必要

## 4. 実装設計

### 4.1 ID生成ユーティリティの実装

#### 4.1.1 ファイル: `server/internal/util/idgen/sonyflake.go`

```go
package idgen

import (
    "sync"
    "github.com/sony/sonyflake"
)

var (
    sf   *sonyflake.Sonyflake
    once sync.Once
)

// initSonyflake はsonyflakeインスタンスを初期化する（初回のみ）
func initSonyflake() {
    once.Do(func() {
        st := sonyflake.Settings{}
        sf = sonyflake.NewSonyflake(st)
        if sf == nil {
            panic("failed to initialize sonyflake")
        }
    })
}

// GenerateSonyflakeID はsonyflakeを使用して一意のIDを生成する
func GenerateSonyflakeID() (int64, error) {
    initSonyflake()
    
    id, err := sf.NextID()
    if err != nil {
        return 0, fmt.Errorf("failed to generate ID: %w", err)
    }
    
    return int64(id), nil
}
```

#### 4.1.2 エラーハンドリング

- sonyflakeの初期化失敗: panic（アプリケーション起動時の致命的エラー）
- ID生成失敗: errorを返す（呼び出し側で処理）

### 4.2 Repository層でのID生成統合

#### 4.2.1 DmUserRepositoryGORM.Create() の修正

**ファイル**: `server/internal/repository/dm_user_repository_gorm.go`

```go
import (
    "github.com/taku-o/go-webdb-template/internal/util/idgen"
)

func (r *DmUserRepositoryGORM) Create(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error) {
    // ID生成
    id, err := idgen.GenerateSonyflakeID()
    if err != nil {
        return nil, fmt.Errorf("failed to generate ID: %w", err)
    }
    
    user := &model.DmUser{
        ID:    id,  // 生成したIDを設定
        Name:  req.Name,
        Email: req.Email,
    }
    
    // 既存のロジック（テーブル名生成、接続取得、保存）
    // ...
}
```

#### 4.2.2 DmPostRepositoryGORM.Create() の修正

**ファイル**: `server/internal/repository/dm_post_repository_gorm.go`

```go
import (
    "github.com/taku-o/go-webdb-template/internal/util/idgen"
)

func (r *DmPostRepositoryGORM) Create(ctx context.Context, req *model.CreateDmPostRequest) (*model.DmPost, error) {
    // ID生成
    id, err := idgen.GenerateSonyflakeID()
    if err != nil {
        return nil, fmt.Errorf("failed to generate ID: %w", err)
    }
    
    post := &model.DmPost{
        ID:      id,  // 生成したIDを設定
        UserID:  req.UserID,
        Title:   req.Title,
        Content: req.Content,
    }
    
    // 既存のロジック（テーブル名生成、接続取得、保存）
    // ...
}
```

#### 4.2.3 database/sql版のRepositoryも同様に修正

- `server/internal/repository/dm_user_repository.go`
- `server/internal/repository/dm_post_repository.go`

### 4.3 サンプルデータ生成コマンドの修正

#### 4.3.1 generateDmUsers() の修正

**ファイル**: `server/cmd/generate-sample-data/main.go`

```go
import (
    "github.com/taku-o/go-webdb-template/internal/util/idgen"
)

func generateDmUsers(groupManager *db.GroupManager, totalCount int) ([]int64, error) {
    // ...
    var dmUsers []*model.DmUser
    for i := 0; i < countPerTable; i++ {
        // ID生成
        id, err := idgen.GenerateSonyflakeID()
        if err != nil {
            return nil, fmt.Errorf("failed to generate ID: %w", err)
        }
        
        dmUser := &model.DmUser{
            ID:        id,  // 生成したIDを設定
            Name:      gofakeit.Name(),
            Email:     gofakeit.Email(),
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        }
        dmUsers = append(dmUsers, dmUser)
        allDmUserIDs = append(allDmUserIDs, id)  // IDをリストに追加
    }
    
    // バッチ挿入（IDを含める）
    if err := insertDmUsersBatch(conn.DB, tableName, dmUsers); err != nil {
        return nil, fmt.Errorf("failed to insert batch to %s: %w", tableName, err)
    }
    
    // fetchDmUserIDs() は不要になるため削除
    // ...
}
```

#### 4.3.2 insertDmUsersBatch() の修正

```go
func insertDmUsersBatch(db *gorm.DB, tableName string, dmUsers []*model.DmUser) error {
    // ...
    // INSERT文にidカラムを追加
    query := fmt.Sprintf("INSERT INTO %s (id, name, email, created_at, updated_at) VALUES ", tableName)
    var values []interface{}
    var placeholders []string

    for _, dmUser := range batch {
        placeholders = append(placeholders, "(?, ?, ?, ?, ?)")
        values = append(values, dmUser.ID, dmUser.Name, dmUser.Email, dmUser.CreatedAt, dmUser.UpdatedAt)
    }
    // ...
}
```

#### 4.3.3 generateDmPosts() の修正

```go
func generateDmPosts(groupManager *db.GroupManager, dmUserIDs []int64, totalCount int) error {
    // ...
    for i := 0; i < countPerTable; i++ {
        dmUserID := dmUserIDs[gofakeit.IntRange(0, len(dmUserIDs)-1)]
        
        // ID生成
        id, err := idgen.GenerateSonyflakeID()
        if err != nil {
            return fmt.Errorf("failed to generate ID: %w", err)
        }
        
        dmPost := &model.DmPost{
            ID:        id,  // 生成したIDを設定
            UserID:    dmUserID,
            Title:     gofakeit.Sentence(5),
            Content:   gofakeit.Paragraph(3, 5, 10, "\n"),
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        }
        dmPosts = append(dmPosts, dmPost)
    }
    // ...
}
```

#### 4.3.4 insertDmPostsBatch() の修正

```go
func insertDmPostsBatch(db *gorm.DB, tableName string, dmPosts []*model.DmPost) error {
    // ...
    // INSERT文にidカラムを追加
    query := fmt.Sprintf("INSERT INTO %s (id, user_id, title, content, created_at, updated_at) VALUES ", tableName)
    var values []interface{}
    var placeholders []string

    for _, dmPost := range batch {
        placeholders = append(placeholders, "(?, ?, ?, ?, ?, ?)")
        values = append(values, dmPost.ID, dmPost.UserID, dmPost.Title, dmPost.Content, dmPost.CreatedAt, dmPost.UpdatedAt)
    }
    // ...
}
```

#### 4.3.5 generateDmNews() の修正

```go
func generateDmNews(groupManager *db.GroupManager, totalCount int) error {
    // ...
    for i := 0; i < totalCount; i++ {
        // ID生成
        id, err := idgen.GenerateSonyflakeID()
        if err != nil {
            return fmt.Errorf("failed to generate ID: %w", err)
        }
        
        authorID := gofakeit.Int64()
        publishedAt := gofakeit.Date()

        n := &model.DmNews{
            ID:          id,  // 生成したIDを設定
            Title:       gofakeit.Sentence(5),
            Content:     gofakeit.Paragraph(3, 5, 10, "\n"),
            AuthorID:    &authorID,
            PublishedAt: &publishedAt,
            CreatedAt:   time.Now(),
            UpdatedAt:   time.Now(),
        }
        dmNews = append(dmNews, n)
    }
    // ...
}
```

#### 4.3.6 fetchDmUserIDs() の削除

- 挿入後にIDを取得する必要がなくなるため、この関数は削除する

### 4.4 Sharding規則の定義

#### 4.4.1 ファイル: `server/internal/db/sharding.go`

```go
// TableShardingRule はテーブルのsharding規則を定義する
type TableShardingRule struct {
    TableName      string // テーブル名（例: "dm_users", "dm_posts"）
    ShardingKey   string // shardingキー（例: "id", "user_id"）
    Description    string // 説明
}

// TableShardingRules は全テーブルのsharding規則を定義
var TableShardingRules = map[string]TableShardingRule{
    "dm_users": {
        TableName:    "dm_users",
        ShardingKey:  "id",
        Description:  "dm_users_NNNテーブルはidをshardingキーとして使用。テーブル番号 = id % DBShardingTableCount",
    },
    "dm_posts": {
        TableName:    "dm_posts",
        ShardingKey:  "user_id",
        Description:  "dm_posts_NNNテーブルはuser_idをshardingキーとして使用。テーブル番号 = user_id % DBShardingTableCount。dm_usersのIDとdm_postsのuser_idが同じ値であれば、同じテーブル番号になる。",
    },
}

// GetShardingKey はテーブル名からshardingキーを取得
func GetShardingKey(tableName string) (string, error) {
    rule, exists := TableShardingRules[tableName]
    if !exists {
        return "", fmt.Errorf("sharding rule not found for table: %s", tableName)
    }
    return rule.ShardingKey, nil
}
```

### 4.5 モデル定義の修正

#### 4.5.1 DmNewsモデルの修正

**ファイル**: `server/internal/model/dm_news.go`

```go
type DmNews struct {
    ID          int64      `json:"id,string" db:"id" gorm:"primaryKey"`  // autoIncrementを削除
    Title       string     `json:"title" db:"title" gorm:"type:varchar(255);not null"`
    Content     string     `json:"content" db:"content" gorm:"type:text;not null"`
    AuthorID    *int64     `json:"author_id,omitempty,string" db:"author_id" gorm:"index:idx_dm_news_author_id"`
    PublishedAt *time.Time `json:"published_at,omitempty" db:"published_at" gorm:"index:idx_dm_news_published_at"`
    CreatedAt   time.Time  `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   time.Time  `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
}
```

## 5. 実装順序

### Phase 1: ライブラリ導入とID生成ユーティリティの実装

1. **sonyflakeライブラリの導入**
   - `go.mod`に依存関係を追加
   - `go mod tidy`を実行
   - ビルド確認

2. **ID生成ユーティリティの実装**
   - `server/internal/util/idgen/sonyflake.go`を作成
   - `GenerateSonyflakeID()`関数を実装
   - 単体テストを実装（`sonyflake_test.go`）

### Phase 2: テーブル定義の修正

3. **dm_newsテーブル定義の修正**
   - `db/schema/master.hcl`を修正
   - Atlasスキーマ検証

4. **dm_users_NNNテーブル定義の修正**
   - `db/schema/sharding_*/dm_users.hcl`を修正（4ファイル）
   - Atlasスキーマ検証

5. **dm_posts_NNNテーブル定義の修正**
   - `db/schema/sharding_*/dm_posts.hcl`を修正（4ファイル）
   - Atlasスキーマ検証

6. **初期データ用SQLの移行**
   - AtlasでマイグレーションSQLを作成し直した際に、既存の初期データ用SQLを新しいファイルに移行
   - `db/migrations/master/20251229111855_initial_schema.sql`の下の方に書いてある初期データ用のSQL（GoAdminの初期データ）を新しいマイグレーションファイルに移行
   - **注意**: 既存のデータの維持は考えなくて良いが、初期データ用SQLの移行を忘れてはいけない

### Phase 3: モデル定義の修正

6. **DmNewsモデルの修正**
   - `server/internal/model/dm_news.go`を修正
   - autoIncrementタグを削除
   - コンパイル確認

### Phase 4: Repository層でのID生成統合

7. **DmUserRepositoryGORMの修正**
   - `server/internal/repository/dm_user_repository_gorm.go`を修正
   - ID生成を追加
   - テスト実行

8. **DmPostRepositoryGORMの修正**
   - `server/internal/repository/dm_post_repository_gorm.go`を修正
   - ID生成を追加
   - テスト実行

9. **database/sql版のRepositoryの修正**
   - `server/internal/repository/dm_user_repository.go`を修正
   - `server/internal/repository/dm_post_repository.go`を修正
   - テスト実行

### Phase 5: サンプルデータ生成コマンドの修正

10. **generate-sample-dataコマンドの修正**
    - `server/cmd/generate-sample-data/main.go`を修正
    - ID生成を統合
    - `fetchDmUserIDs()`を削除
    - コマンド実行確認

### Phase 6: Sharding規則の定義

11. **sharding規則の定義**
    - `server/internal/db/sharding.go`に規則を追加
    - ドキュメントコメントを追加

### Phase 7: テストの実装

12. **単体テストの実装**
    - ID生成ユーティリティのテスト
    - 一意性のテスト
    - エラーハンドリングのテスト

13. **統合テストの修正**
    - 既存の統合テストが正常に動作することを確認
    - ID生成が正常に動作することを確認

14. **シャーディング規則の動作確認テスト**
    - `dm_users`と`dm_posts`が同じテーブル番号に配置されることを確認
    - テーブル番号の計算ロジックが正しく動作することを確認

### Phase 8: ドキュメントの更新

15. **ドキュメントの更新**
    - `docs/Architecture.md`に「Identifier Generation」セクションを追加
    - identifier生成ルールを記載
    - ID生成方式の変更を記載
    - sharding規則の記載

## 6. エラーハンドリング設計

### 6.1 ID生成エラー

#### 6.1.1 エラーケース

- sonyflakeの初期化失敗: アプリケーション起動時の致命的エラー（panic）
- ID生成失敗: 呼び出し側でエラーを返す

#### 6.1.2 エラーハンドリング

```go
// ID生成ユーティリティ
func GenerateSonyflakeID() (int64, error) {
    initSonyflake()
    
    id, err := sf.NextID()
    if err != nil {
        return 0, fmt.Errorf("failed to generate ID: %w", err)
    }
    
    return int64(id), nil
}

// Repository層
func (r *DmUserRepositoryGORM) Create(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error) {
    id, err := idgen.GenerateSonyflakeID()
    if err != nil {
        return nil, fmt.Errorf("failed to generate ID: %w", err)
    }
    // ...
}
```

### 6.2 テーブル定義エラー

- Atlasスキーマ検証エラー: 修正が必要
- マイグレーションエラー: マイグレーションスクリプトの確認が必要

## 7. テスト設計

### 7.1 単体テスト

#### 7.1.1 ID生成ユーティリティのテスト

**ファイル**: `server/internal/util/idgen/sonyflake_test.go`

```go
func TestGenerateSonyflakeID(t *testing.T) {
    // ID生成のテスト
    id1, err := GenerateSonyflakeID()
    require.NoError(t, err)
    require.Greater(t, id1, int64(0))
    
    // 一意性のテスト
    ids := make(map[int64]bool)
    for i := 0; i < 1000; i++ {
        id, err := GenerateSonyflakeID()
        require.NoError(t, err)
        require.False(t, ids[id], "duplicate ID: %d", id)
        ids[id] = true
    }
}
```

### 7.2 統合テスト

#### 7.2.1 Repository層のテスト

- 既存の統合テストが正常に動作することを確認
- ID生成が正常に動作することを確認
- 生成されたIDが一意であることを確認

#### 7.2.2 シャーディング規則の動作確認テスト

**ファイル**: `server/test/integration/sharding_rule_test.go`（新規作成）

```go
func TestShardingRule_DmUsersAndDmPostsSameTable(t *testing.T) {
    // dm_usersを作成
    user, err := userRepo.Create(ctx, &model.CreateDmUserRequest{...})
    require.NoError(t, err)
    
    // テーブル番号を計算
    userTableNumber := int(user.ID % db.DBShardingTableCount)
    
    // dm_postsを作成（user.IDをuser_idとして使用）
    post, err := postRepo.Create(ctx, &model.CreateDmPostRequest{
        UserID: user.ID,
        ...
    })
    require.NoError(t, err)
    
    // テーブル番号を計算
    postTableNumber := int(post.UserID % db.DBShardingTableCount)
    
    // 同じテーブル番号であることを確認
    require.Equal(t, userTableNumber, postTableNumber, 
        "dm_users and dm_posts should be in the same table number")
}
```

## 8. パフォーマンス考慮

### 8.1 ID生成のパフォーマンス

- sonyflakeは分散環境で効率的に動作する
- ID生成のオーバーヘッドは最小限に抑える
- シングルトンパターンでsonyflakeインスタンスを管理

### 8.2 メモリ使用量

- sonyflakeインスタンスは1つのみ（シングルトン）
- メモリ使用量は最小限

## 9. セキュリティ考慮

### 9.1 IDの一意性

- sonyflakeは分散環境で一意性が保証される
- 時系列順序が保たれるため、推測が困難

### 9.2 エラーハンドリング

- ID生成失敗時の適切なエラーハンドリング
- ログ出力による問題の追跡

## 10. 互換性考慮

### 10.1 API互換性

- 既存のAPIインターフェースは変更しない
- IDは文字列として返される（`json:"id,string"`タグ）
- JavaScript側での互換性を維持

### 10.2 データ互換性

- 既存のデータ構造との互換性を維持
- 既存のテストが正常に動作することを確認

## 11. ドキュメント設計

### 11.1 Architecture.mdへの追加

**ファイル**: `docs/Architecture.md`

「Database Sharding」セクションの後に「Identifier Generation」セクションを追加：

```markdown
## Identifier Generation

本プロジェクトでは、分散テーブル環境に対応するため、以下のidentifier生成ルールを採用しています。

### 数値のidentifier

数値のidentifierが必要な箇所は**sonyflake** (github.com/sony/sonyflake) を使用します。

- **用途**: データベースの主キー（dm_users.id, dm_posts.id, dm_news.idなど）
- **理由**: 
  - 分散環境で一意性が保証される
  - 時系列順序が保たれる
  - 64ビット整数として生成される

### 文字列のidentifier

文字列のidentifierが必要な箇所は**UUIDv7** (github.com/google/uuid) を使用します。

- **用途**: 文字列型のIDが必要な場合（APIキー、セッションIDなど）
- **理由**: 
  - 時系列順序が保たれる
  - グローバルに一意
  - 標準的なUUID形式

### JavaScript側での扱い

sonyflakeで生成されるIDは64ビット整数であり、JavaScriptのNumber型の安全な整数範囲（2^53-1）を超える可能性があるため、**文字列として扱います**。

- APIレスポンスでは`json:"id,string"`タグにより文字列として返される
- JavaScript側では文字列として扱う必要がある
```

## 12. リスクと対策

### 12.1 リスク

- sonyflakeの初期化失敗
- ID生成のパフォーマンス問題
- 既存データとの互換性問題

### 12.2 対策

- sonyflakeの初期化失敗時はpanic（アプリケーション起動時の致命的エラー）
- パフォーマンステストを実施
- 既存データの移行は別途検討（本実装の範囲外）

## 13. 受け入れ基準

1. sonyflakeライブラリが正常に導入されている
2. ID生成ユーティリティが実装されている
3. 全テーブル定義が修正されている
4. 全モデル定義が修正されている
5. Repository層でID生成が統合されている
6. サンプルデータ生成コマンドでID生成が統合されている
7. sharding規則が定義されている
8. 単体テストが実装されている
9. 統合テストが正常に動作する
10. シャーディング規則の動作確認テストが正常に動作する
11. ドキュメントが更新されている
12. ビルドが正常に完了する
13. 既存のテストが全て正常に動作する
