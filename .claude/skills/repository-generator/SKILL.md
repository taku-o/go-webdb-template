---
name: repository-generator
description: 新しいRepository層のコードを作成する際に使用。GORMを標準とし、GroupManager/TableSelector使用、UUIDベースのシャーディング対応、CRUDメソッドのパターンを適用。Repositoryの新規作成、データアクセス層の実装時に使用。
allowed-tools: Read, Glob
---

# Repository 生成パターン

このプロジェクトのRepository層実装パターンを定義します。**本プロジェクトでは GORM を標準とする。** 標準SQL版は他プロジェクト向けであり、本プロジェクトでは GORM 版のみを参照すること。

## ファイル構成

新しいエンティティを追加する場合、`dm_*` の命名規則に従う:

```
server/internal/repository/
├── dm_user_repository.go
├── dm_post_repository.go
└── dm_news_repository.go
```

## 参照ファイル

Repository実装の参照:
- `server/internal/repository/dm_user_repository.go`
- `server/internal/repository/dm_post_repository.go`
- `server/internal/repository/dm_news_repository.go`

モデル定義:
- `server/internal/model/` の DmUser, DmPost, DmNews 等（CreateDmUserRequest 等のリクエスト型を含む）

## 構造体パターン

```go
package repository

import (
    "context"
    "errors"
    "fmt"
    "time"

    "github.com/taku-o/go-webdb-template/internal/db"
    "github.com/taku-o/go-webdb-template/internal/model"
    "github.com/taku-o/go-webdb-template/internal/util/idgen"
    "gorm.io/gorm"
)

// DmUserRepository はユーザーのデータアクセスを担当
type DmUserRepository struct {
    groupManager  *db.GroupManager
    tableSelector *db.TableSelector
}

// NewDmUserRepository は新しいDmUserRepositoryを作成
func NewDmUserRepository(groupManager *db.GroupManager) *DmUserRepository {
    return &DmUserRepository{
        groupManager:  groupManager,
        tableSelector: db.NewTableSelector(db.DBShardingTableCount, db.DBShardingTablesPerDB),
    }
}
```

## CRUDメソッドパターン（GORM）

エンティティの ID は UUID 文字列（string）。生成には `idgen.GenerateUUIDv7()` を使用する。接続取得は `GetShardingConnectionByUUID(uuid, tableBaseName)`、テーブル名取得は `GetTableNameFromUUID(tableBaseName, uuid)`（戻り値 `(string, error)`）を使用する。

### Create

```go
func (r *DmUserRepository) Create(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error) {
    // ID生成（UUIDv7）
    id, err := idgen.GenerateUUIDv7()
    if err != nil {
        return nil, fmt.Errorf("failed to generate ID: %w", err)
    }

    user := &model.DmUser{
        ID:    id,
        Name:  req.Name,
        Email: req.Email,
    }

    // テーブル名の取得（戻り値 (string, error)）
    tableName, err := r.tableSelector.GetTableNameFromUUID("dm_users", user.ID)
    if err != nil {
        return nil, fmt.Errorf("failed to get table name: %w", err)
    }

    // 接続の取得（UUID・テーブルベース名で指定）
    conn, err := r.groupManager.GetShardingConnectionByUUID(user.ID, "dm_users")
    if err != nil {
        return nil, fmt.Errorf("failed to get sharding connection: %w", err)
    }

    err = db.ExecuteWithRetry(func() error {
        return conn.DB.WithContext(ctx).Table(tableName).Create(user).Error
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }

    return user, nil
}
```

### GetByID

```go
func (r *DmUserRepository) GetByID(ctx context.Context, id string) (*model.DmUser, error) {
    tableName, err := r.tableSelector.GetTableNameFromUUID("dm_users", id)
    if err != nil {
        return nil, fmt.Errorf("failed to get table name: %w", err)
    }

    conn, err := r.groupManager.GetShardingConnectionByUUID(id, "dm_users")
    if err != nil {
        return nil, fmt.Errorf("failed to get sharding connection: %w", err)
    }

    var user model.DmUser
    err = db.ExecuteWithRetry(func() error {
        return conn.DB.WithContext(ctx).Table(tableName).Where("id = ?", id).First(&user).Error
    })
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, fmt.Errorf("user not found: %s", id)
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }

    return &user, nil
}
```

### List（クロステーブルクエリ）

```go
func (r *DmUserRepository) List(ctx context.Context, limit, offset int) ([]*model.DmUser, error) {
    users := make([]*model.DmUser, 0)
    tableCount := r.tableSelector.GetTableCount()

    for tableNum := 0; tableNum < tableCount; tableNum++ {
        conn, err := r.groupManager.GetShardingConnection(tableNum)
        if err != nil {
            return nil, fmt.Errorf("failed to get connection for table %d: %w", tableNum, err)
        }
        tableName := fmt.Sprintf("dm_users_%03d", tableNum)

        var tableUsers []*model.DmUser
        err = db.ExecuteWithRetry(func() error {
            return conn.DB.WithContext(ctx).
                Table(tableName).
                Order("id").
                Limit(limit).
                Offset(offset).
                Find(&tableUsers).Error
        })
        if err != nil {
            return nil, fmt.Errorf("failed to list users: %w", err)
        }
        users = append(users, tableUsers...)
    }

    return users, nil
}
```

### Update

```go
func (r *DmUserRepository) Update(ctx context.Context, id string, req *model.UpdateDmUserRequest) (*model.DmUser, error) {
    tableName, err := r.tableSelector.GetTableNameFromUUID("dm_users", id)
    if err != nil {
        return nil, fmt.Errorf("failed to get table name: %w", err)
    }

    conn, err := r.groupManager.GetShardingConnectionByUUID(id, "dm_users")
    if err != nil {
        return nil, fmt.Errorf("failed to get sharding connection: %w", err)
    }

    updates := make(map[string]interface{})
    if req.Name != "" {
        updates["name"] = req.Name
    }
    if req.Email != "" {
        updates["email"] = req.Email
    }
    updates["updated_at"] = time.Now()

    var result *gorm.DB
    err = db.ExecuteWithRetry(func() error {
        result = conn.DB.WithContext(ctx).Table(tableName).Where("id = ?", id).Updates(updates)
        return result.Error
    })
    if err != nil {
        return nil, fmt.Errorf("failed to update user: %w", err)
    }
    if result.RowsAffected == 0 {
        return nil, fmt.Errorf("user not found: %s", id)
    }

    return r.GetByID(ctx, id)
}
```

### Delete

```go
func (r *DmUserRepository) Delete(ctx context.Context, id string) error {
    tableName, err := r.tableSelector.GetTableNameFromUUID("dm_users", id)
    if err != nil {
        return fmt.Errorf("failed to get table name: %w", err)
    }

    conn, err := r.groupManager.GetShardingConnectionByUUID(id, "dm_users")
    if err != nil {
        return fmt.Errorf("failed to get sharding connection: %w", err)
    }

    var result *gorm.DB
    err = db.ExecuteWithRetry(func() error {
        result = conn.DB.WithContext(ctx).Table(tableName).Where("id = ?", id).Delete(&model.DmUser{})
        return result.Error
    })
    if err != nil {
        return fmt.Errorf("failed to delete user: %w", err)
    }
    if result.RowsAffected == 0 {
        return fmt.Errorf("user not found: %s", id)
    }

    return nil
}
```

## エラーハンドリング

```go
// エラーラップのパターン
return nil, fmt.Errorf("failed to create user: %w", err)

// Not Found（GORM）
if errors.Is(err, gorm.ErrRecordNotFound) {
    return nil, fmt.Errorf("user not found: %s", id)
}
```

## 必要なモデル定義

エンティティの ID は UUID 文字列（string）として扱う。

```go
// server/internal/model/dm_user.go
package model

type DmUser struct {
    ID        string    `json:"id" gorm:"column:id;primaryKey"`
    Name      string    `json:"name" gorm:"column:name"`
    Email     string    `json:"email" gorm:"column:email"`
    CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
    UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

type CreateDmUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

type UpdateDmUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}
```
