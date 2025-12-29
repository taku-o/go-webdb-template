# テーブル名変更機能設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、テーブル名を`users`→`dm_users`、`posts`→`dm_posts`、`news`→`dm_news`に変更する機能の詳細設計を定義する。コードベース全体でのテーブル名参照の変更、スキーマ定義ファイルの変更、ファイル名の変更、import文の更新を体系的に実装する。

### 1.2 設計の範囲
- モデル層のテーブル名変更設計
- Repository層のテーブル名参照変更設計
- データベース層のテーブル名生成ロジック変更設計
- スキーマ定義ファイル（HCL）の変更設計
- ファイル名変更とimport文の更新設計
- テストコードの更新設計
- CLIツール、管理画面、ドキュメントの更新設計
- 実装順序と検証方法

### 1.3 設計方針
- **一貫性の確保**: コードベース全体でテーブル名の参照を一貫して変更する
- **段階的実装**: 影響範囲が広いため、層ごとに段階的に実装する
- **検証の徹底**: 各段階でテストを実行し、問題がないことを確認する
- **既存アーキテクチャの維持**: 既存のレイヤードアーキテクチャを維持する
- **ファイル名の統一**: テーブル名変更に合わせてファイル名も変更する

## 2. アーキテクチャ設計

### 2.1 変更前後の構造

#### 2.1.1 テーブル名の変更

```
変更前:
- users (32分割: users_000 ～ users_031)
- posts (32分割: posts_000 ～ posts_031)
- news (単一テーブル)

変更後:
- dm_users (32分割: dm_users_000 ～ dm_users_031)
- dm_posts (32分割: dm_posts_000 ～ dm_posts_031)
- dm_news (単一テーブル)
```

#### 2.1.2 ファイル名の変更

```
変更前:
server/internal/model/
  - user.go
  - post.go
  - news.go

server/internal/repository/
  - user_repository.go
  - user_repository_gorm.go
  - user_repository_test.go
  - user_repository_gorm_test.go
  - post_repository.go
  - post_repository_gorm.go
  - post_repository_test.go
  - post_repository_gorm_test.go
  - news_repository.go
  - news_repository_gorm.go
  - news_repository_test.go
  - news_repository_gorm_test.go

db/schema/sharding_*/
  - users.hcl
  - posts.hcl

変更後:
server/internal/model/
  - dm_user.go
  - dm_post.go
  - dm_news.go

server/internal/repository/
  - dm_user_repository.go
  - dm_user_repository_gorm.go
  - dm_user_repository_test.go
  - dm_user_repository_gorm_test.go
  - dm_post_repository.go
  - dm_post_repository_gorm.go
  - dm_post_repository_test.go
  - dm_post_repository_gorm_test.go
  - dm_news_repository.go
  - dm_news_repository_gorm.go
  - dm_news_repository_test.go
  - dm_news_repository_gorm_test.go

db/schema/sharding_*/
  - dm_users.hcl
  - dm_posts.hcl
```

### 2.2 変更フロー

```
┌─────────────────────────────────────────────────────────────┐
│           テーブル名変更実装フロー                           │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│  Phase 1: スキーマ定義ファイルの変更                          │
│  - db/schema/master.hcl: news → dm_news                    │
│  - db/schema/sharding_*/users.hcl: users_* → dm_users_*    │
│  - db/schema/sharding_*/posts.hcl: posts_* → dm_posts_*    │
│  - インデックス名の変更                                      │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│  Phase 2: データベース層の変更                               │
│  - server/internal/db/sharding.go                          │
│    - GetTableName()のベース名引数変更                       │
│    - ValidateTableName()の許可リスト更新                    │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│  Phase 3: モデル層の変更                                     │
│  - server/internal/model/*.go                               │
│    - TableName()メソッドの戻り値変更                        │
│  - ファイル名変更: user.go → dm_user.go など               │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│  Phase 4: Repository層の変更                                 │
│  - server/internal/repository/*_repository.go               │
│    - GetTableName()呼び出しの変更                           │
│    - 文字列リテラルの変更                                   │
│  - ファイル名変更: *_repository.go → dm_*_repository.go   │
│  - import文の更新（モデルファイル名変更に対応）             │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│  Phase 5: テストコードの更新                                 │
│  - server/test/testutil/db.go: スキーマ作成SQL             │
│  - server/test/integration/sharding_test.go                 │
│  - server/internal/db/sharding_test.go                      │
│  - import文の更新                                           │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│  Phase 6: CLIツール、管理画面、ドキュメントの更新            │
│  - server/cmd/generate-sample-data/main.go                   │
│  - server/internal/admin/tables.go                          │
│  - docs/Sharding.md                                         │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│  Phase 7: データベース再作成と検証                           │
│  - 既存データベースファイルの削除                           │
│  - マイグレーション適用スクリプトの実行                     │
│  - 全テストの実行と検証                                     │
└─────────────────────────────────────────────────────────────┘
```

## 3. 各層の詳細設計

### 3.1 スキーマ定義層の変更設計

#### 3.1.1 マスターデータベーススキーマ

**ファイル**: `db/schema/master.hcl`

**変更内容**:
```hcl
// 変更前
table "news" {
  // ...
  index "idx_news_published_at" {
    columns = [column.published_at]
  }
  index "idx_news_author_id" {
    columns = [column.author_id]
  }
}

// 変更後
table "dm_news" {
  // ...
  index "idx_dm_news_published_at" {
    columns = [column.published_at]
  }
  index "idx_dm_news_author_id" {
    columns = [column.author_id]
  }
}
```

#### 3.1.2 シャーディングデータベーススキーマ

**ファイル**: `db/schema/sharding_*/users.hcl`, `posts.hcl`

**変更内容**:
```hcl
// 変更前: users.hcl
table "users_000" {
  // ...
  index "idx_users_000_email" {
    unique = true
    columns = [column.email]
  }
}

// 変更後: dm_users.hcl
table "dm_users_000" {
  // ...
  index "idx_dm_users_000_email" {
    unique = true
    columns = [column.email]
  }
}
```

```hcl
// 変更前: posts.hcl
table "posts_000" {
  // ...
  index "idx_posts_000_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_000_created_at" {
    columns = [column.created_at]
  }
}

// 変更後: dm_posts.hcl
table "dm_posts_000" {
  // ...
  index "idx_dm_posts_000_user_id" {
    columns = [column.user_id]
  }
  index "idx_dm_posts_000_created_at" {
    columns = [column.created_at]
  }
}
```

**注意事項**:
- 4つのシャーディングデータベース（sharding_1～sharding_4）全てで変更が必要
- ファイル名も変更: `users.hcl` → `dm_users.hcl`, `posts.hcl` → `dm_posts.hcl`
- Atlas設定ファイル（`config/*/atlas.hcl`）でスキーマファイルの参照パスを更新する必要がある場合がある

### 3.2 データベース層の変更設計

#### 3.2.1 TableSelectorの変更

**ファイル**: `server/internal/db/sharding.go`

**変更内容**:
```go
// 変更前
func (ts *TableSelector) GetTableName(baseName string, id int64) string {
    tableNumber := ts.GetTableNumber(id)
    return fmt.Sprintf("%s_%03d", baseName, tableNumber)
}

// 使用例（変更前）
tableName := selector.GetTableName("users", userID)  // "users_005"

// 変更後（使用例のみ変更、実装は同じ）
tableName := selector.GetTableName("dm_users", userID)  // "dm_users_005"
```

**ValidateTableName関数の変更**:
```go
// 変更前
func ValidateTableName(tableName string, allowedBaseNames []string) bool {
    for _, baseName := range allowedBaseNames {
        // "users" をチェック
        for i := 0; i < 32; i++ {
            expectedName := fmt.Sprintf("%s_%03d", baseName, i)
            if tableName == expectedName {
                return true
            }
        }
    }
    return false
}

// 変更後: 許可リストに "dm_users", "dm_posts" を追加
// 使用例
allowedBaseNames := []string{"dm_users", "dm_posts"}
valid := ValidateTableName("dm_users_005", allowedBaseNames)  // true
```

**GetShardingTableName関数の変更**:
```go
// 変更前
func GetShardingTableName(baseName string, id int64) string {
    tableNumber := int(id % 32)
    return fmt.Sprintf("%s_%03d", baseName, tableNumber)
}

// 使用例（変更前）
tableName := db.GetShardingTableName("users", userID)  // "users_005"

// 変更後（使用例のみ変更、実装は同じ）
tableName := db.GetShardingTableName("dm_users", userID)  // "dm_users_005"
```

### 3.3 モデル層の変更設計

#### 3.3.1 Userモデルの変更

**ファイル**: `server/internal/model/user.go` → `server/internal/model/dm_user.go`

**変更内容**:
```go
// 変更前
package model

func (User) TableName() string {
    return "users"
}

// 変更後
package model

func (User) TableName() string {
    return "dm_users"
}
```

**注意事項**:
- パッケージ名は`model`のまま（変更不要）
- ファイル名のみ変更: `user.go` → `dm_user.go`
- クラス名（`User`）は変更不要

#### 3.3.2 Postモデルの変更

**ファイル**: `server/internal/model/post.go` → `server/internal/model/dm_post.go`

**変更内容**:
```go
// 変更前
func (Post) TableName() string {
    return "posts"
}

// 変更後
func (Post) TableName() string {
    return "dm_posts"
}
```

#### 3.3.3 Newsモデルの変更

**ファイル**: `server/internal/model/news.go` → `server/internal/model/dm_news.go`

**変更内容**:
```go
// 変更前
func (News) TableName() string {
    return "news"
}

// 変更後
func (News) TableName() string {
    return "dm_news"
}
```

### 3.4 Repository層の変更設計

#### 3.4.1 UserRepositoryの変更

**ファイル**: 
- `server/internal/repository/user_repository.go` → `server/internal/repository/dm_user_repository.go`
- `server/internal/repository/user_repository_gorm.go` → `server/internal/repository/dm_user_repository_gorm.go`

**変更内容**:
```go
// 変更前: user_repository.go
tableName := r.tableSelector.GetTableName("users", user.ID)
conn, err := r.groupManager.GetShardingConnectionByID(user.ID, "users")

// 変更後: dm_user_repository.go
tableName := r.tableSelector.GetTableName("dm_users", user.ID)
conn, err := r.groupManager.GetShardingConnectionByID(user.ID, "dm_users")
```

**全テーブル検索時の変更**:
```go
// 変更前
tableName := fmt.Sprintf("users_%03d", tableNum)

// 変更後
tableName := fmt.Sprintf("dm_users_%03d", tableNum)
```

**import文の更新**:
```go
// 変更前
import (
    "github.com/taku-o/go-webdb-template/internal/model"
)

// 変更後: importパスは変更不要（パッケージ名は同じ）
// ただし、ファイル名が変更されたため、IDEやツールが自動的に認識する
import (
    "github.com/taku-o/go-webdb-template/internal/model"
)
```

#### 3.4.2 PostRepositoryの変更

**ファイル**: 
- `server/internal/repository/post_repository.go` → `server/internal/repository/dm_post_repository.go`
- `server/internal/repository/post_repository_gorm.go` → `server/internal/repository/dm_post_repository_gorm.go`

**変更内容**:
```go
// 変更前
tableName := r.tableSelector.GetTableName("posts", req.UserID)
conn, err := r.groupManager.GetShardingConnectionByID(req.UserID, "posts")

// 変更後
tableName := r.tableSelector.GetTableName("dm_posts", req.UserID)
conn, err := r.groupManager.GetShardingConnectionByID(req.UserID, "dm_posts")
```

#### 3.4.3 NewsRepositoryの変更

**ファイル**: 
- `server/internal/repository/news_repository.go` → `server/internal/repository/dm_news_repository.go`
- `server/internal/repository/news_repository_gorm.go` → `server/internal/repository/dm_news_repository_gorm.go`

**変更内容**:
```go
// 変更前: news_repository_gorm.go
if err := db.Table("news").CreateInBatches(batch, len(batch)).Error; err != nil {
    return fmt.Errorf("failed to create news: %w", err)
}

// 変更後: dm_news_repository_gorm.go
if err := db.Table("dm_news").CreateInBatches(batch, len(batch)).Error; err != nil {
    return fmt.Errorf("failed to create news: %w", err)
}
```

### 3.5 テストコードの変更設計

#### 3.5.1 テストユーティリティの変更

**ファイル**: `server/test/testutil/db.go`

**変更内容**:
```go
// 変更前
func InitMasterSchema(t *testing.T, database *gorm.DB) {
    schema := `
        CREATE TABLE IF NOT EXISTS news (
            // ...
        );
    `
    err := database.Exec(schema).Error
    require.NoError(t, err)
}

func InitShardingSchema(t *testing.T, database *gorm.DB, startTable, endTable int) {
    for i := startTable; i <= endTable; i++ {
        suffix := fmt.Sprintf("%03d", i)
        
        usersSchema := fmt.Sprintf(`
            CREATE TABLE IF NOT EXISTS users_%s (
                // ...
            );
        `, suffix)
        err := database.Exec(usersSchema).Error
        require.NoError(t, err)
        
        postsSchema := fmt.Sprintf(`
            CREATE TABLE IF NOT EXISTS posts_%s (
                // ...
                FOREIGN KEY (user_id) REFERENCES users_%s(id)
            );
        `, suffix, suffix)
        err = database.Exec(postsSchema).Error
        require.NoError(t, err)
    }
}

// 変更後
func InitMasterSchema(t *testing.T, database *gorm.DB) {
    schema := `
        CREATE TABLE IF NOT EXISTS dm_news (
            // ...
        );
    `
    err := database.Exec(schema).Error
    require.NoError(t, err)
}

func InitShardingSchema(t *testing.T, database *gorm.DB, startTable, endTable int) {
    for i := startTable; i <= endTable; i++ {
        suffix := fmt.Sprintf("%03d", i)
        
        usersSchema := fmt.Sprintf(`
            CREATE TABLE IF NOT EXISTS dm_users_%s (
                // ...
            );
        `, suffix)
        err := database.Exec(usersSchema).Error
        require.NoError(t, err)
        
        postsSchema := fmt.Sprintf(`
            CREATE TABLE IF NOT EXISTS dm_posts_%s (
                // ...
            );
        `, suffix)
        // 注意: 外部キー制約は削除（分散データ環境では使用しない）
        err = database.Exec(postsSchema).Error
        require.NoError(t, err)
    }
}
```

#### 3.5.2 統合テストの変更

**ファイル**: `server/test/integration/sharding_test.go`

**変更内容**:
```go
// 変更前
tableName := tableSelector.GetTableName("users", tc.id)
assert.Equal(t, "users_005", tableName)

// 変更後
tableName := tableSelector.GetTableName("dm_users", tc.id)
assert.Equal(t, "dm_users_005", tableName)
```

#### 3.5.3 ユニットテストの変更

**ファイル**: `server/internal/db/sharding_test.go`

**変更内容**:
```go
// 変更前
tests := []struct {
    baseName      string
    id            int64
    wantTableName string
}{
    {baseName: "users", id: 0, wantTableName: "users_000"},
    {baseName: "users", id: 1, wantTableName: "users_001"},
    {baseName: "posts", id: 15, wantTableName: "posts_015"},
}

// 変更後
tests := []struct {
    baseName      string
    id            int64
    wantTableName string
}{
    {baseName: "dm_users", id: 0, wantTableName: "dm_users_000"},
    {baseName: "dm_users", id: 1, wantTableName: "dm_users_001"},
    {baseName: "dm_posts", id: 15, wantTableName: "dm_posts_015"},
}
```

### 3.6 CLIツールの変更設計

#### 3.6.1 サンプルデータ生成ツールの変更

**ファイル**: `server/cmd/generate-sample-data/main.go`

**変更内容**:
```go
// 変更前
tableName := fmt.Sprintf("users_%03d", tableNumber)
if err := insertUsersBatch(conn.DB, tableName, users); err != nil {
    return nil, fmt.Errorf("failed to insert batch to %s: %w", tableName, err)
}

tableName := fmt.Sprintf("posts_%03d", tableNumber)
if err := insertPostsBatch(conn.DB, tableName, posts); err != nil {
    return fmt.Errorf("failed to insert batch to %s: %w", tableName, err)
}

if err := db.Table("news").CreateInBatches(batch, len(batch)).Error; err != nil {
    return fmt.Errorf("failed to create news: %w", err)
}

// 変更後
tableName := fmt.Sprintf("dm_users_%03d", tableNumber)
if err := insertUsersBatch(conn.DB, tableName, users); err != nil {
    return nil, fmt.Errorf("failed to insert batch to %s: %w", tableName, err)
}

tableName := fmt.Sprintf("dm_posts_%03d", tableNumber)
if err := insertPostsBatch(conn.DB, tableName, posts); err != nil {
    return fmt.Errorf("failed to insert batch to %s: %w", tableName, err)
}

if err := db.Table("dm_news").CreateInBatches(batch, len(batch)).Error; err != nil {
    return fmt.Errorf("failed to create news: %w", err)
}
```

### 3.7 管理画面の変更設計

#### 3.7.1 GoAdmin設定の変更

**ファイル**: `server/internal/admin/tables.go`

**変更内容**:
```go
// 変更前
info.SetTable("news").SetTitle("ニュース").SetDescription("ニュース一覧")

// 変更後
info.SetTable("dm_news").SetTitle("ニュース").SetDescription("ニュース一覧")
```

### 3.8 ドキュメントの変更設計

#### 3.8.1 シャーディングドキュメントの更新

**ファイル**: `docs/Sharding.md`

**変更内容**:
- テーブル名の記載を全て更新: `users` → `dm_users`, `posts` → `dm_posts`, `news` → `dm_news`
- コード例のテーブル名を更新
- 図表や説明文のテーブル名を更新

## 4. ファイル名変更とimport文の更新設計

### 4.1 ファイル名変更の影響範囲

#### 4.1.1 モデルファイルの変更

**変更対象**:
- `server/internal/model/user.go` → `server/internal/model/dm_user.go`
- `server/internal/model/post.go` → `server/internal/model/dm_post.go`
- `server/internal/model/news.go` → `server/internal/model/dm_news.go`

**影響範囲**:
- Repositoryファイルでのimport（パッケージ名は同じなので、importパスは変更不要）
- テストファイルでのimport
- 他のパッケージでのimport

**注意事項**:
- Goのimportはパッケージ名（`model`）に基づくため、ファイル名変更だけではimport文の変更は不要
- ただし、IDEやツールがファイル名を参照する場合があるため、ビルドエラーが発生する可能性がある
- ファイル名変更後、`go mod tidy`を実行して依存関係を更新

#### 4.1.2 Repositoryファイルの変更

**変更対象**:
- `server/internal/repository/user_repository.go` → `server/internal/repository/dm_user_repository.go`
- `server/internal/repository/user_repository_gorm.go` → `server/internal/repository/dm_user_repository_gorm.go`
- `server/internal/repository/user_repository_test.go` → `server/internal/repository/dm_user_repository_test.go`
- `server/internal/repository/user_repository_gorm_test.go` → `server/internal/repository/dm_user_repository_gorm_test.go`
- （posts, newsも同様）

**影響範囲**:
- ハンドラーやサービス層でのimport
- テストファイルでのimport
- インターフェース定義ファイル（`interfaces.go`）での参照

**注意事項**:
- Repositoryファイルは他のパッケージからimportされる可能性がある
- ファイル名変更後、importしている箇所を確認し、必要に応じて更新

#### 4.1.3 スキーマ定義ファイルの変更

**変更対象**:
- `db/schema/sharding_1/users.hcl` → `db/schema/sharding_1/dm_users.hcl`
- `db/schema/sharding_1/posts.hcl` → `db/schema/sharding_1/dm_posts.hcl`
- （sharding_2, 3, 4も同様）

**影響範囲**:
- Atlas設定ファイル（`config/*/atlas.hcl`）でのスキーマファイル参照
- マイグレーション生成時の参照

**注意事項**:
- Atlas設定ファイルでスキーマファイルのパスを指定している場合、更新が必要
- ファイル名変更後、Atlas設定を確認し、必要に応じて更新

### 4.2 import文の更新方法

#### 4.2.1 Goのimport文

Goのimportはパッケージ名に基づくため、ファイル名変更だけではimport文の変更は基本的に不要:

```go
// 変更前後で同じ（パッケージ名は model のまま）
import (
    "github.com/taku-o/go-webdb-template/internal/model"
)
```

ただし、以下の場合に注意が必要:
- ファイル名を直接参照するツールやスクリプトがある場合
- テストファイルで特定のファイルを参照している場合

#### 4.2.2 検証方法

ファイル名変更後、以下のコマンドで検証:

```bash
# ビルドエラーの確認
go build ./...

# テストの実行
go test ./...

# 依存関係の更新
go mod tidy
```

## 5. 実装順序

### 5.1 Phase 1: スキーマ定義ファイルの変更

1. `db/schema/master.hcl`の変更
   - `table "news"` → `table "dm_news"`
   - インデックス名の変更
2. `db/schema/sharding_*/users.hcl`の変更
   - `table "users_*"` → `table "dm_users_*"`
   - インデックス名の変更
   - ファイル名変更: `users.hcl` → `dm_users.hcl`
3. `db/schema/sharding_*/posts.hcl`の変更
   - `table "posts_*"` → `table "dm_posts_*"`
   - インデックス名の変更
   - ファイル名変更: `posts.hcl` → `dm_posts.hcl`
4. Atlas設定ファイルの確認（必要に応じて更新）

**検証**:
- HCLファイルの構文チェック
- Atlas設定ファイルの確認

### 5.2 Phase 2: データベース層の変更

1. `server/internal/db/sharding.go`の変更
   - `ValidateTableName()`関数の許可リスト更新
   - コメントやドキュメントの更新

**検証**:
- ユニットテストの実行: `go test ./server/internal/db/...`
- テストケースの更新（Phase 5で実施）

### 5.3 Phase 3: モデル層の変更

1. `server/internal/model/user.go`の変更
   - `TableName()`メソッドの戻り値変更
   - ファイル名変更: `user.go` → `dm_user.go`
2. `server/internal/model/post.go`の変更
   - `TableName()`メソッドの戻り値変更
   - ファイル名変更: `post.go` → `dm_post.go`
3. `server/internal/model/news.go`の変更
   - `TableName()`メソッドの戻り値変更
   - ファイル名変更: `news.go` → `dm_news.go`

**検証**:
- ビルドエラーの確認: `go build ./server/internal/model/...`
- `go mod tidy`の実行

### 5.4 Phase 4: Repository層の変更

1. `server/internal/repository/user_repository.go`の変更
   - `GetTableName("users", ...)` → `GetTableName("dm_users", ...)`
   - 文字列リテラル`"users"` → `"dm_users"`
   - ファイル名変更: `user_repository.go` → `dm_user_repository.go`
2. `server/internal/repository/user_repository_gorm.go`の変更
   - 同様の変更
   - ファイル名変更: `user_repository_gorm.go` → `dm_user_repository_gorm.go`
3. `server/internal/repository/post_repository.go`の変更
   - 同様の変更
   - ファイル名変更: `post_repository.go` → `dm_post_repository.go`
4. `server/internal/repository/post_repository_gorm.go`の変更
   - 同様の変更
   - ファイル名変更: `post_repository_gorm.go` → `dm_post_repository_gorm.go`
5. `server/internal/repository/news_repository.go`の変更
   - `Table("news")` → `Table("dm_news")`
   - ファイル名変更: `news_repository.go` → `dm_news_repository.go`
6. `server/internal/repository/news_repository_gorm.go`の変更
   - 同様の変更
   - ファイル名変更: `news_repository_gorm.go` → `dm_news_repository_gorm.go`

**検証**:
- ビルドエラーの確認: `go build ./server/internal/repository/...`
- import文の確認

### 5.5 Phase 5: テストコードの更新

1. `server/test/testutil/db.go`の変更
   - スキーマ作成SQLのテーブル名変更
   - 外部キー制約の削除（分散データ環境では使用しない）
2. `server/test/integration/sharding_test.go`の変更
   - テーブル名参照の変更
   - 期待値の変更
3. `server/internal/db/sharding_test.go`の変更
   - テストケースのテーブル名変更
4. Repositoryテストファイルのファイル名変更
   - `*_repository_test.go` → `dm_*_repository_test.go`
   - `*_repository_gorm_test.go` → `dm_*_repository_gorm_test.go`

**検証**:
- 全テストの実行: `go test ./...`
- テストが全て通過することを確認

### 5.6 Phase 6: CLIツール、管理画面、ドキュメントの更新

1. `server/cmd/generate-sample-data/main.go`の変更
   - テーブル名文字列の変更
2. `server/internal/admin/tables.go`の変更
   - GoAdmin設定のテーブル名変更
3. `docs/Sharding.md`の更新
   - テーブル名の記載を更新

**検証**:
- CLIツールのビルド: `go build ./server/cmd/generate-sample-data/...`
- ドキュメントの確認

### 5.7 Phase 7: データベース再作成と検証

1. 既存データベースファイルの削除
   - `server/data/master.db`
   - `server/data/sharding_db_*.db`
2. マイグレーション適用スクリプトの実行
   - `./scripts/migrate.sh all`
3. 全テストの実行と検証
   - `go test ./...`
   - 統合テストの実行
   - CLIツールの動作確認

**検証**:
- 全テストが通過することを確認
- データベースに正しいテーブル名でテーブルが作成されていることを確認
- CLIツールが正常に動作することを確認

## 6. エラーハンドリング

### 6.1 ビルドエラー

**原因**:
- ファイル名変更によるimportエラー
- テーブル名参照の漏れ

**対応**:
- `go build ./...`でビルドエラーを確認
- エラーメッセージに基づいて修正
- `go mod tidy`で依存関係を更新

### 6.2 テストエラー

**原因**:
- テストコード内のテーブル名参照の漏れ
- 期待値の更新漏れ

**対応**:
- `go test ./...`でテストエラーを確認
- エラーメッセージに基づいて修正
- 各Phaseでテストを実行し、段階的に検証

### 6.3 データベースエラー

**原因**:
- スキーマ定義ファイルの構文エラー
- マイグレーション適用時のエラー

**対応**:
- Atlasの構文チェック: `atlas schema validate`
- マイグレーション適用前の確認
- エラーログの確認

## 7. テスト戦略

### 7.1 ユニットテスト

- **対象**: 各層の変更箇所
- **方法**: 既存のユニットテストを更新
- **検証項目**:
  - テーブル名生成ロジックの動作確認
  - モデルの`TableName()`メソッドの戻り値確認

### 7.2 統合テスト

- **対象**: Repository層、データベース層
- **方法**: 既存の統合テストを更新
- **検証項目**:
  - テーブル名が正しく生成されること
  - データベース操作が正常に動作すること

### 7.3 E2Eテスト

- **対象**: データベース再作成後の動作確認
- **方法**: マイグレーション適用後の全機能テスト
- **検証項目**:
  - 全テストが通過すること
  - CLIツールが正常に動作すること
  - 管理画面が正常に動作すること

## 8. 実装上の注意事項

### 8.1 テーブル名の一貫性

- コードベース全体でテーブル名の参照を一貫して変更する
- 正規表現検索（`"users"`, `"posts"`, `"news"`）を使用して漏れがないか確認
- 文字列リテラル、コメント、ドキュメントなど、全ての箇所で変更を確認

### 8.2 ファイル名変更の影響

- ファイル名変更後、import文が正しく動作するか確認
- `go mod tidy`を実行して依存関係を更新
- IDEやツールがファイル名を参照する場合があるため、ビルドエラーを確認

### 8.3 段階的実装

- 影響範囲が広いため、Phaseごとに段階的に実装する
- 各Phaseでテストを実行し、問題がないことを確認してから次のPhaseに進む
- 問題が発生した場合は、前のPhaseに戻って修正

### 8.4 データベース再作成

- 既存データベースのデータは破棄して良いため、既存データベースファイルを削除して再作成する
- マイグレーション適用スクリプトを実行してデータベースを再作成する
- 再作成後、全テストを実行して動作確認

## 9. 参考情報

### 9.1 関連ドキュメント
- 要件定義書: `.kiro/specs/0026-chtablename/requirements.md`
- Atlas運用ガイド: `docs/atlas-operations.md`
- シャーディング戦略: `docs/Sharding.md`

### 9.2 既存実装
- モデル層: `server/internal/model/*.go`
- Repository層: `server/internal/repository/*.go`
- データベース層: `server/internal/db/sharding.go`
- スキーマ定義: `db/schema/**/*.hcl`

### 9.3 技術スタック
- **Go**: 1.21+
- **GORM**: v1.25.12
- **データベース**: SQLite3（開発環境）
- **Atlas**: スキーマ管理ツール
