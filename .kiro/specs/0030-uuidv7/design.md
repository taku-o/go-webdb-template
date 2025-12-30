# UUIDv7導入設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、Sonyflakeで生成したIDをシャーディングのキーに使用すると同じテーブルにデータが振り分けられる問題を解決するため、UUIDv7を使用したID生成方式に変更し、シャーディングキーの計算方法を改善するシステムの詳細設計を定義する。既存システム（dm_newsなど）との共存を保ちながら、dm_usersとdm_postsテーブルのID生成方式をUUIDv7に変更する。

### 1.2 設計の範囲
- UUIDv7ライブラリの導入
- ID生成ユーティリティの実装（UUIDv7）
- シャーディングキー計算ロジックの変更（UUIDベース）
- モデル定義の変更（IDの型を`int64`から`string`に）
- リポジトリ層の変更（ID生成とシャーディングキー計算）
- サービス層の変更（IDの型変更）
- API層の変更（IDの型変更）
- サンプルデータ生成コマンドの変更
- GoAdmin管理画面の変更
- クライアント側の変更（APIレスポンス型定義）
- Sonyflake関数の削除
- マイグレーションの作成
- ドキュメントの更新

### 1.3 設計方針
- **既存システムとの共存**: dm_newsなど他のテーブルは変更しない
- **後方互換性の維持**: 既存の`GetTableName(baseName string, id int64)`関数と`GetShardingConnectionByID(id int64, tableName string)`関数は残す
- **型の一貫性**: IDの型を`int64`から`string`に変更する際、すべての関連箇所で一貫して変更する
- **データ分散性の向上**: UUIDv7の後ろ2文字を使用したシャーディングキー計算により、より均等なデータ分散を実現
- **既存データの破棄**: 既存データは破棄し、マイグレーション時に削除する
- **Sonyflake関数の削除**: dm_usersとdm_postsでUUIDv7に置き換えられるため、Sonyflake関数は削除する

## 2. アーキテクチャ設計

### 2.1 ディレクトリ構造

#### 2.1.1 変更前の構造
```
server/
├── internal/
│   ├── util/
│   │   └── idgen/
│   │       ├── sonyflake.go          # Sonyflake関数（削除）
│   │       └── sonyflake_test.go     # Sonyflakeテスト（削除）
│   ├── db/
│   │   └── sharding.go               # シャーディング関数（既存関数は残す）
│   ├── model/
│   │   ├── dm_user.go                # ID: int64
│   │   └── dm_post.go                # ID, UserID: int64
│   ├── repository/
│   │   ├── dm_user_repository.go     # Sonyflake使用
│   │   └── dm_post_repository.go     # Sonyflake使用
│   ├── service/
│   │   ├── dm_user_service.go        # ID: int64
│   │   └── dm_post_service.go       # ID: int64
│   ├── api/
│   │   ├── handler/
│   │   │   ├── dm_user_handler.go    # ID: int64
│   │   │   └── dm_post_handler.go    # ID: int64
│   │   └── huma/
│   │       └── outputs.go             # ID: int64
│   └── admin/
│       └── pages/
│           └── dm_user_register.go   # UnixNano使用
├── cmd/
│   └── generate-sample-data/
│       └── main.go                    # Sonyflake使用
└── go.mod                              # sonyflake依存関係あり
```

#### 2.1.2 変更後の構造
```
server/
├── internal/
│   ├── util/
│   │   └── idgen/
│   │       ├── uuidv7.go              # UUIDv7関数（新規）
│   │       └── uuidv7_test.go        # UUIDv7テスト（新規）
│   ├── db/
│   │   └── sharding.go                # 新規関数追加、既存関数は残す
│   ├── model/
│   │   ├── dm_user.go                 # ID: string
│   │   └── dm_post.go                 # ID, UserID: string
│   ├── repository/
│   │   ├── dm_user_repository.go      # UUIDv7使用
│   │   └── dm_post_repository.go      # UUIDv7使用
│   ├── service/
│   │   ├── dm_user_service.go        # ID: string
│   │   └── dm_post_service.go       # ID: string
│   ├── api/
│   │   ├── handler/
│   │   │   ├── dm_user_handler.go    # ID: string
│   │   │   └── dm_post_handler.go    # ID: string
│   │   └── huma/
│   │       └── outputs.go              # ID: string
│   └── admin/
│       └── pages/
│           └── dm_user_register.go    # UUIDv7使用
├── cmd/
│   └── generate-sample-data/
│       └── main.go                     # UUIDv7使用
└── go.mod                               # uuid依存関係追加、sonyflake削除
```

### 2.2 ファイル構成

#### 2.2.1 新規作成ファイル
- **`server/internal/util/idgen/uuidv7.go`**: UUIDv7生成関数の実装
- **`server/internal/util/idgen/uuidv7_test.go`**: UUIDv7生成関数のテスト
- **Atlasマイグレーションファイル**: データベーススキーマの変更

#### 2.2.2 変更ファイル
- **モデル定義**: `dm_user.go`, `dm_post.go`（IDの型を`string`に変更）
- **リポジトリ層**: `dm_user_repository.go`, `dm_user_repository_gorm.go`, `dm_post_repository.go`, `dm_post_repository_gorm.go`（ID生成とシャーディングキー計算を変更）
- **サービス層**: `dm_user_service.go`, `dm_post_service.go`（IDの型を`string`に変更）
- **API層**: `dm_user_handler.go`, `dm_post_handler.go`, `outputs.go`（IDの型を`string`に変更）
- **シャーディング関連**: `sharding.go`（新規関数を追加）
- **サンプルデータ生成**: `generate-sample-data/main.go`（ID生成をUUIDv7に変更）
- **GoAdmin管理画面**: `dm_user_register.go`（ID生成をUUIDv7に変更）
- **依存関係**: `go.mod`（uuid追加、sonyflake削除）

#### 2.2.3 削除ファイル
- **`server/internal/util/idgen/sonyflake.go`**: Sonyflake関数の実装
- **`server/internal/util/idgen/sonyflake_test.go`**: Sonyflake関数のテスト

### 2.3 システム構成図

```
┌─────────────────────────────────────────────────────────┐
│                    APIリクエスト                          │
│              (ID: string型のUUID)                        │
└──────────────────┬────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────┐
│              API層 (Handler)                            │
│  - dm_user_handler.go                                   │
│  - dm_post_handler.go                                   │
│  - outputs.go (ID: string)                             │
└──────────────────┬────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────┐
│              サービス層 (Service)                        │
│  - dm_user_service.go (ID: string)                      │
│  - dm_post_service.go (ID: string)                      │
└──────────────────┬────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────┐
│             リポジトリ層 (Repository)                     │
│  - dm_user_repository.go                                 │
│  - dm_post_repository.go                                 │
│  - ID生成: GenerateUUIDv7()                              │
│  - シャーディング: GetTableNameFromUUID()                │
└──────────────────┬────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────┐
│            ID生成ユーティリティ (idgen)                  │
│  - uuidv7.go: GenerateUUIDv7() → string                │
│    1. uuid.NewV7()でUUIDv7を生成                        │
│    2. ハイフンを削除                                     │
│    3. 小文字に変換                                        │
│    4. 32文字の文字列として返す                           │
└──────────────────┬────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────┐
│        シャーディングキー計算 (sharding.go)              │
│  - GetTableNumberFromUUID(uuid string) int              │
│    1. UUIDの後ろ2文字を取得                              │
│    2. 16進数として解釈                                    │
│    3. 32で割った余りを計算                                │
│    4. テーブル番号（0～31）を返す                         │
│  - GetTableNameFromUUID(baseName, uuid) string          │
│  - GetShardingConnectionByUUID(uuid, tableName)          │
└──────────────────┬────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────┐
│              データベース (SQLite)                        │
│  - dm_users_000 ～ dm_users_031                         │
│  - dm_posts_000 ～ dm_posts_031                          │
│  - id: varchar(32) (UUID文字列)                         │
└─────────────────────────────────────────────────────────┘
```

### 2.4 データフロー

#### 2.4.1 ユーザー作成フロー
```
APIリクエスト（CreateDmUserRequest）
    ↓
Handler (dm_user_handler.go)
    ↓
Service (dm_user_service.go)
    ↓
Repository (dm_user_repository.go)
    ↓
1. GenerateUUIDv7() → UUID文字列を生成
    ↓
2. GetTableNumberFromUUID(uuid) → テーブル番号を計算
    ↓
3. GetTableNameFromUUID("dm_users", uuid) → テーブル名を取得
    ↓
4. GetShardingConnectionByUUID(uuid, "dm_users") → 接続を取得
    ↓
5. データベースにINSERT
    ↓
APIレスポンス（DmUserOutput: ID: string）
```

#### 2.4.2 投稿作成フロー
```
APIリクエスト（CreateDmPostRequest: UserID: string）
    ↓
Handler (dm_post_handler.go)
    ↓
Service (dm_post_service.go)
    ↓
Repository (dm_post_repository.go)
    ↓
1. GenerateUUIDv7() → 投稿IDを生成
    ↓
2. GetTableNumberFromUUID(userID) → テーブル番号を計算（user_idから）
    ↓
3. GetTableNameFromUUID("dm_posts", userID) → テーブル名を取得
    ↓
4. GetShardingConnectionByUUID(userID, "dm_posts") → 接続を取得
    ↓
5. データベースにINSERT
    ↓
APIレスポンス（DmPostOutput: ID, UserID: string）
```

#### 2.4.3 シャーディングキー計算フロー
```
UUID文字列: "550e8400e29b41d4a716446655440000"
    ↓
後ろ2文字を取得: "00"
    ↓
16進数として解釈: 0x00 = 0
    ↓
32で割った余りを計算: 0 % 32 = 0
    ↓
テーブル番号: 0
    ↓
テーブル名: "dm_users_000" または "dm_posts_000"
```

## 3. コンポーネント設計

### 3.1 UUIDv7生成関数

#### 3.1.1 uuidv7.goの構造
```go
package idgen

import (
    "strings"
    "github.com/google/uuid"
)

// GenerateUUIDv7 はUUIDv7を生成し、ハイフン抜き小文字32文字の文字列として返す
func GenerateUUIDv7() (string, error) {
    // 1. UUIDv7を生成
    u, err := uuid.NewV7()
    if err != nil {
        return "", fmt.Errorf("failed to generate UUIDv7: %w", err)
    }
    
    // 2. ハイフンを削除
    uuidStr := strings.ReplaceAll(u.String(), "-", "")
    
    // 3. 小文字に変換
    uuidStr = strings.ToLower(uuidStr)
    
    // 4. 32文字の文字列として返す
    return uuidStr, nil
}
```

#### 3.1.2 実装の詳細
- **UUID生成**: `uuid.NewV7()`を使用してUUIDv7を生成
- **ハイフン削除**: `strings.ReplaceAll(uuid.String(), "-", "")`でハイフンを削除
- **小文字変換**: `strings.ToLower()`で小文字に変換
- **エラーハンドリング**: UUID生成エラーは適切に処理し、エラーメッセージを返す
- **戻り値**: 32文字の文字列（例: `550e8400e29b41d4a716446655440000`）

### 3.2 シャーディングキー計算関数

#### 3.2.1 GetTableNumberFromUUID関数の構造
```go
// GetTableNumberFromUUID はUUID文字列からテーブル番号を取得
func (ts *TableSelector) GetTableNumberFromUUID(uuid string) (int, error) {
    // 1. UUID文字列の長さをチェック
    if len(uuid) < 2 {
        return 0, fmt.Errorf("invalid UUID string: length must be at least 2")
    }
    
    // 2. 後ろ2文字を取得
    suffix := uuid[len(uuid)-2:]
    
    // 3. 16進数として解釈
    value, err := strconv.ParseInt(suffix, 16, 64)
    if err != nil {
        return 0, fmt.Errorf("failed to parse UUID suffix as hex: %w", err)
    }
    
    // 4. テーブル数（32）で割った余りを計算
    tableNumber := int(value % int64(ts.tableCount))
    
    // 5. テーブル番号（0～31）を返す
    return tableNumber, nil
}
```

#### 3.2.2 GetTableNameFromUUID関数の構造
```go
// GetTableNameFromUUID はベース名とUUIDからテーブル名を生成
func (ts *TableSelector) GetTableNameFromUUID(baseName string, uuid string) (string, error) {
    tableNumber, err := ts.GetTableNumberFromUUID(uuid)
    if err != nil {
        return "", fmt.Errorf("failed to get table number from UUID: %w", err)
    }
    return fmt.Sprintf("%s_%03d", baseName, tableNumber), nil
}
```

#### 3.2.3 GetShardingConnectionByUUID関数の構造
```go
// GetShardingConnectionByUUID はUUIDからshardingグループの接続を取得
func (gm *GroupManager) GetShardingConnectionByUUID(uuid string, tableName string) (*GORMConnection, error) {
    // 1. UUIDからテーブル番号を計算
    selector := NewTableSelector(DBShardingTableCount, DBShardingTablesPerDB)
    tableNumber, err := selector.GetTableNumberFromUUID(uuid)
    if err != nil {
        return nil, fmt.Errorf("failed to get table number from UUID: %w", err)
    }
    
    // 2. テーブル番号からデータベースIDを取得
    dbID := selector.GetDBID(tableNumber)
    
    // 3. 適切な接続を返す
    return gm.GetShardingConnection(dbID)
}
```

### 3.3 モデル定義の変更

#### 3.3.1 DmUserモデルの変更
```go
// 変更前
type DmUser struct {
    ID        int64     `json:"id,string" db:"id" gorm:"primaryKey"`
    Name      string    `json:"name" db:"name" gorm:"type:varchar(100);not null"`
    Email     string    `json:"email" db:"email" gorm:"type:varchar(255);not null;uniqueIndex:idx_dm_users_email"`
    CreatedAt time.Time `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
}

// 変更後
type DmUser struct {
    ID        string    `json:"id" db:"id" gorm:"primaryKey;type:varchar(32)"`
    Name      string    `json:"name" db:"name" gorm:"type:varchar(100);not null"`
    Email     string    `json:"email" db:"email" gorm:"type:varchar(255);not null;uniqueIndex:idx_dm_users_email"`
    CreatedAt time.Time `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
}
```

#### 3.3.2 DmPostモデルの変更
```go
// 変更前
type DmPost struct {
    ID        int64     `json:"id,string" db:"id" gorm:"primaryKey"`
    UserID    int64     `json:"user_id,string" db:"user_id" gorm:"type:bigint;not null;index:idx_dm_posts_user_id"`
    Title     string    `json:"title" db:"title" gorm:"type:varchar(200);not null"`
    Content   string    `json:"content" db:"content" gorm:"type:text;not null"`
    CreatedAt time.Time `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
}

// 変更後
type DmPost struct {
    ID        string    `json:"id" db:"id" gorm:"primaryKey;type:varchar(32)"`
    UserID    string    `json:"user_id" db:"user_id" gorm:"type:varchar(32);not null;index:idx_dm_posts_user_id"`
    Title     string    `json:"title" db:"title" gorm:"type:varchar(200);not null"`
    Content   string    `json:"content" db:"content" gorm:"type:text;not null"`
    CreatedAt time.Time `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
}
```

#### 3.3.3 リクエストモデルの変更
```go
// 変更前
type CreateDmPostRequest struct {
    UserID  int64  `json:"user_id,string" validate:"required,gt=0"`
    Title   string `json:"title" validate:"required,min=1,max=200"`
    Content string `json:"content" validate:"required,min=1"`
}

// 変更後
type CreateDmPostRequest struct {
    UserID  string `json:"user_id" validate:"required,len=32"`
    Title   string `json:"title" validate:"required,min=1,max=200"`
    Content string `json:"content" validate:"required,min=1"`
}
```

### 3.4 リポジトリ層の変更

#### 3.4.1 DmUserRepository.Createメソッドの変更
```go
// 変更前
func (r *DmUserRepository) Create(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error) {
    id, err := idgen.GenerateSonyflakeID()
    if err != nil {
        return nil, fmt.Errorf("failed to generate ID: %w", err)
    }
    
    user := &model.DmUser{
        ID:        id,  // int64
        // ...
    }
    
    tableName := r.tableSelector.GetTableName("dm_users", user.ID)
    conn, err := r.groupManager.GetShardingConnectionByID(user.ID, "dm_users")
    // ...
}

// 変更後
func (r *DmUserRepository) Create(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error) {
    id, err := idgen.GenerateUUIDv7()
    if err != nil {
        return nil, fmt.Errorf("failed to generate ID: %w", err)
    }
    
    user := &model.DmUser{
        ID:        id,  // string
        // ...
    }
    
    tableName, err := r.tableSelector.GetTableNameFromUUID("dm_users", user.ID)
    if err != nil {
        return nil, fmt.Errorf("failed to get table name: %w", err)
    }
    conn, err := r.groupManager.GetShardingConnectionByUUID(user.ID, "dm_users")
    // ...
}
```

#### 3.4.2 DmPostRepository.Createメソッドの変更
```go
// 変更前
func (r *DmPostRepository) Create(ctx context.Context, req *model.CreateDmPostRequest) (*model.DmPost, error) {
    id, err := idgen.GenerateSonyflakeID()
    if err != nil {
        return nil, fmt.Errorf("failed to generate ID: %w", err)
    }
    
    post := &model.DmPost{
        ID:        id,  // int64
        UserID:    req.UserID,  // int64
        // ...
    }
    
    tableName := r.tableSelector.GetTableName("dm_posts", req.UserID)
    conn, err := r.groupManager.GetShardingConnectionByID(req.UserID, "dm_posts")
    // ...
}

// 変更後
func (r *DmPostRepository) Create(ctx context.Context, req *model.CreateDmPostRequest) (*model.DmPost, error) {
    id, err := idgen.GenerateUUIDv7()
    if err != nil {
        return nil, fmt.Errorf("failed to generate ID: %w", err)
    }
    
    post := &model.DmPost{
        ID:        id,  // string
        UserID:    req.UserID,  // string
        // ...
    }
    
    tableName, err := r.tableSelector.GetTableNameFromUUID("dm_posts", req.UserID)
    if err != nil {
        return nil, fmt.Errorf("failed to get table name: %w", err)
    }
    conn, err := r.groupManager.GetShardingConnectionByUUID(req.UserID, "dm_posts")
    // ...
}
```

### 3.5 GoAdmin管理画面の変更

#### 3.5.1 insertDmUserSharded関数の変更
```go
// 変更前
func insertDmUserSharded(groupManager *appdb.GroupManager, name, email string) (int64, error) {
    now := time.Now()
    dmUserID := now.UnixNano()
    tableNumber := int(dmUserID % appdb.DBShardingTableCount)
    // ...
}

// 変更後
func insertDmUserSharded(groupManager *appdb.GroupManager, name, email string) (string, error) {
    now := time.Now()
    
    // UUIDv7でIDを生成
    dmUserID, err := idgen.GenerateUUIDv7()
    if err != nil {
        return "", fmt.Errorf("failed to generate UUIDv7: %w", err)
    }
    
    // UUIDからテーブル番号を計算
    selector := appdb.NewTableSelector(appdb.DBShardingTableCount, appdb.DBShardingTablesPerDB)
    tableNumber, err := selector.GetTableNumberFromUUID(dmUserID)
    if err != nil {
        return "", fmt.Errorf("failed to get table number: %w", err)
    }
    // ...
}
```

#### 3.5.2 handleDmUserRegisterPost関数の変更
```go
// 変更前
redirectURL := fmt.Sprintf("/admin/dm-user/register/new?id=%d&name=%s&email=%s",
    dmUserID,  // int64
    url.QueryEscape(name),
    url.QueryEscape(email),
)

// 変更後
redirectURL := fmt.Sprintf("/admin/dm-user/register/new?id=%s&name=%s&email=%s",
    url.QueryEscape(dmUserID),  // string
    url.QueryEscape(name),
    url.QueryEscape(email),
)
```

### 3.6 マイグレーション設計

#### 3.6.1 マイグレーションの内容
- **カラム型の変更**:
  - `dm_users.id`: `bigint unsigned` → `varchar(32)`
  - `dm_posts.id`: `bigint unsigned` → `varchar(32)`
  - `dm_posts.user_id`: `bigint` → `varchar(32)`
- **既存データの削除**: 既存データは破棄するため、テーブルを再作成するか、既存データを削除する
- **インデックスの再作成**: 型変更に伴いインデックスを再作成する

#### 3.6.2 マイグレーションの実装方法
- Atlasマイグレーションファイルを作成
- テーブルを再作成する方法:
  1. 既存テーブルを削除（`DROP TABLE`）
  2. 新しいスキーマでテーブルを作成（`CREATE TABLE`）
- または、既存データを削除してからカラム型を変更する方法:
  1. 既存データを削除（`DELETE FROM`）
  2. カラム型を変更（`ALTER TABLE ... MODIFY COLUMN`）

## 4. エラーハンドリング

### 4.1 UUID生成エラー

#### 4.1.1 エラーケース
- **UUID生成失敗**: `uuid.NewV7()`がエラーを返す場合
- **対処**: エラーメッセージを返し、処理を中断

#### 4.1.2 エラーハンドリングの実装
```go
func GenerateUUIDv7() (string, error) {
    u, err := uuid.NewV7()
    if err != nil {
        return "", fmt.Errorf("failed to generate UUIDv7: %w", err)
    }
    // ...
}
```

### 4.2 シャーディングキー計算エラー

#### 4.2.1 エラーケース
- **無効なUUID文字列**: 長さが2文字未満の場合
- **16進数解析エラー**: 後ろ2文字が16進数として解釈できない場合
- **対処**: エラーメッセージを返し、処理を中断

#### 4.2.2 エラーハンドリングの実装
```go
func (ts *TableSelector) GetTableNumberFromUUID(uuid string) (int, error) {
    if len(uuid) < 2 {
        return 0, fmt.Errorf("invalid UUID string: length must be at least 2")
    }
    
    suffix := uuid[len(uuid)-2:]
    value, err := strconv.ParseInt(suffix, 16, 64)
    if err != nil {
        return 0, fmt.Errorf("failed to parse UUID suffix as hex: %w", err)
    }
    // ...
}
```

### 4.3 バリデーションエラー

#### 4.3.1 API層でのバリデーション
- **UUID文字列の長さチェック**: `validate:"required,len=32"`
- **UUID文字列の形式チェック**: 必要に応じて、16進数文字（0-9a-f）のみであることを確認

#### 4.3.2 バリデーションの実装
```go
type CreateDmPostRequest struct {
    UserID  string `json:"user_id" validate:"required,len=32"`
    Title   string `json:"title" validate:"required,min=1,max=200"`
    Content string `json:"content" validate:"required,min=1"`
}
```

## 5. テスト戦略

### 5.1 単体テスト

#### 5.1.1 UUIDv7生成関数のテスト
- **テスト項目**:
  - UUIDが32文字であること
  - ハイフンが含まれていないこと
  - 小文字であること
  - 一意性（複数回生成して異なる値が返されること）
  - エラーハンドリング（エラーが適切に処理されること）

#### 5.1.2 シャーディングキー計算関数のテスト
- **テスト項目**:
  - UUIDの後ろ2文字を正しく取得できること
  - 16進数として正しく解釈できること
  - テーブル数（32）で割った余りが正しく計算できること
  - テーブル番号が0～31の範囲内であること
  - エラーハンドリング（無効なUUID文字列の場合）

#### 5.1.3 テストケース例
```go
func TestGetTableNumberFromUUID(t *testing.T) {
    selector := NewTableSelector(32, 8)
    
    tests := []struct {
        name      string
        uuid      string
        want      int
        wantError bool
    }{
        {"後ろ2文字が00", "550e8400e29b41d4a716446655440000", 0, false},
        {"後ろ2文字が0f", "550e8400e29b41d4a71644665544000f", 15, false},
        {"後ろ2文字が1f", "550e8400e29b41d4a71644665544001f", 31, false},
        {"後ろ2文字が20", "550e8400e29b41d4a716446655440020", 0, false}, // 32 % 32 = 0
        {"短すぎるUUID", "00", 0, true},
        {"無効な16進数", "550e8400e29b41d4a7164466554400gg", 0, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := selector.GetTableNumberFromUUID(tt.uuid)
            if (err != nil) != tt.wantError {
                t.Errorf("GetTableNumberFromUUID() error = %v, wantError %v", err, tt.wantError)
                return
            }
            if got != tt.want {
                t.Errorf("GetTableNumberFromUUID() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 5.2 統合テスト

#### 5.2.1 リポジトリ層の統合テスト
- **テスト項目**:
  - ユーザー作成時にUUIDv7でIDが生成されること
  - 投稿作成時にUUIDv7でIDが生成されること
  - シャーディングキーが正しく計算されること
  - データが正しいテーブルに保存されること
  - IDの型が`string`であること

#### 5.2.2 API層の統合テスト
- **テスト項目**:
  - APIリクエストでIDが`string`型として受け取れること
  - APIレスポンスでIDが`string`型として返されること
  - バリデーションが正しく動作すること

### 5.3 既存テストの更新

#### 5.3.1 既存テストの修正
- **dm_user_repository_test.go**: IDの型を`string`に変更
- **dm_post_repository_test.go**: IDの型を`string`に変更
- **sharding_test.go**: 新規関数のテストを追加

## 6. セキュリティ考慮事項

### 6.1 UUIDの一意性

#### 6.1.1 UUIDv7の特性
- **時間順序性**: UUIDv7は時間ベースの一意性を保つ
- **分散性**: より均等なデータ分散を実現
- **衝突の可能性**: 理論的には衝突の可能性があるが、実用的には無視できるレベル

### 6.2 バリデーション

#### 6.2.1 入力値の検証
- **UUID文字列の長さ**: 32文字であることを確認
- **UUID文字列の形式**: 16進数文字（0-9a-f）のみであることを確認（必要に応じて）

### 6.3 SQLインジェクション対策

#### 6.3.1 パラメータ化クエリ
- 既存の実装と同様に、パラメータ化クエリを使用
- UUID文字列を直接SQLに埋め込まない

## 7. 実装上の注意事項

### 7.1 型変換の一貫性

#### 7.1.1 すべての関連箇所で一貫して変更
- モデル定義、リポジトリ層、サービス層、API層でIDの型を`string`に統一
- JSONシリアライゼーション時の型変換に注意（`json:"id,string"`タグは不要になる）
- データベースへの保存時の型変換に注意（`varchar(32)`）

### 7.2 後方互換性の維持

#### 7.2.1 既存関数の保持
- `GetTableName(baseName string, id int64) string`: 残す（後方互換性のため）
- `GetShardingConnectionByID(id int64, tableName string) (*GORMConnection, error)`: 残す（将来の使用に備えて保持）

### 7.3 Sonyflake関数の削除

#### 7.3.1 削除手順
1. `server/internal/util/idgen/sonyflake.go`を削除
2. `server/internal/util/idgen/sonyflake_test.go`を削除
3. `go.mod`から`github.com/sony/sonyflake`の依存関係を削除
4. `go mod tidy`を実行して依存関係を整理
5. コード内でSonyflake関数への参照が全て削除されていることを確認
6. `docs/Architecture.md`からSonyflakeの記述を削除

### 7.4 マイグレーションの実装

#### 7.4.1 マイグレーション実行時の注意
- 既存データは破棄するため、テーブルを再作成するか、既存データを削除する
- マイグレーション実行前にバックアップを取得することを推奨（必要に応じて）
- インデックスを再作成する

### 7.5 ドキュメントの更新

#### 7.5.1 更新が必要なドキュメント
- **`docs/Architecture.md`**: UUIDv7の使用を記載し、Sonyflakeの削除を記載
- **`server/internal/db/sharding.go`**: コメントを更新（UUIDv7の使用を記載、Sonyflakeの記述を削除）
- **コード内のコメント**: 適切に更新

## 8. 参考情報

### 8.1 関連ドキュメント
- 要件定義書: `.kiro/specs/0030-uuidv7/requirements.md`
- アーキテクチャドキュメント: `docs/Architecture.md`
- シャーディング仕様: `server/internal/db/sharding.go`
- Sonyflake導入の要件定義書: `.kiro/specs/0028-dmtable-define/requirements.md`

### 8.2 技術スタック
- **UUIDライブラリ**: `github.com/google/uuid`
- **UUIDバージョン**: UUIDv7
- **Goバージョン**: 1.23.4
- **データベース**: SQLite（開発環境）

### 8.3 参考リンク
- UUIDv7仕様: https://www.ietf.org/rfc/rfc4122.txt
- github.com/google/uuid: https://pkg.go.dev/github.com/google/uuid
- UUIDv7の説明: https://www.ietf.org/archive/id/draft-ietf-uuidrev-rfc4122bis-00.html
