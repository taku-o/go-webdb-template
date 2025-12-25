# GoAdmin管理画面導入設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、GoAdmin管理画面の詳細設計を定義する。既存のGORM接続管理を活用しながら、独立した管理画面サービス（ポート8081）を実装し、複数シャード構成に対応したデータ管理機能とカスタムページを提供する。

### 1.2 設計の範囲
- 管理画面サービスのアーキテクチャ設計
- GoAdminフレームワークの統合設計
- 既存GORM接続との統合設計
- カスタムページの実装設計
- 認証・認可機能の設計
- シャーディング対応の設計
- 設定管理の設計
- エラーハンドリング設計
- セキュリティ設計
- テスト戦略

## 2. アーキテクチャ設計

### 2.1 全体アーキテクチャ

既存のメインサービスとは独立した管理画面サービスを実装する。

```
┌─────────────────────────────────────────────────────────────┐
│                    Main Service (Port 8080)                   │
│                  (既存のAPIサーバー)                          │
│                                                               │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐                  │
│  │   API    │  │ Service  │  │Repository│                  │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘                  │
│       │             │             │                          │
│       └─────────────┴─────────────┘                          │
│                       │                                      │
│       ┌──────────────▼──────────────┐                       │
│       │    GORM Manager              │                       │
│       │  (既存の接続管理)             │                       │
│       └──────────────┬──────────────┘                       │
└───────────────────────┼─────────────────────────────────────┘
                        │
         ┌──────────────┴──────────────┐
         │                              │
         ▼                              ▼
    ┌─────────┐                  ┌─────────┐
    │ Shard 1 │                  │ Shard 2 │
    └─────────┘                  └─────────┘

┌─────────────────────────────────────────────────────────────┐
│              Admin Service (Port 8081)                       │
│              (新規の管理画面サービス)                        │
│                                                               │
│  ┌─────────────────────────────────────────────────────┐  │
│  │              GoAdmin Engine                          │  │
│  │  • Admin Plugin (CRUD自動生成)                      │  │
│  │  • Custom Pages (カスタムページ)                    │  │
│  │  • Authentication (認証・認可)                      │  │
│  └──────────────────────┬──────────────────────────────┘  │
│                         │                                    │
│       ┌─────────────────▼─────────────────┐                │
│       │    GORM Manager (再利用)           │                │
│       │  (既存の接続管理を共有)             │                │
│       └─────────────────┬─────────────────┘                │
└──────────────────────────┼───────────────────────────────────┘
                           │
            ┌──────────────┴──────────────┐
            │                              │
            ▼                              ▼
       ┌─────────┐                  ┌─────────┐
       │ Shard 1 │                  │ Shard 2 │
       └─────────┘                  └─────────┘
```

### 2.2 管理画面サービスの構成

管理画面サービスは独立したエントリーポイントとして実装する。

```
server/cmd/admin/main.go
  ├─ 設定ファイル読み込み
  ├─ GORM Manager初期化（既存のものを再利用）
  ├─ GoAdmin Engine初期化
  ├─ テーブル設定（Users, Posts）
  ├─ カスタムページ登録
  ├─ 認証・認可設定
  └─ HTTPサーバー起動（ポート8081）
```

### 2.3 GoAdmin統合アーキテクチャ

GoAdminは既存のGORM接続をアダプターとして使用する。

```
GoAdmin Engine
  ├─ Admin Plugin
  │   ├─ Users Table Generator
  │   │   └─ GORM Adapter → GORM Manager → Shard 1/2
  │   └─ Posts Table Generator
  │       └─ GORM Adapter → GORM Manager → Shard 1/2
  ├─ Custom Pages
  │   ├─ Home Page
  │   ├─ User Register Page
  │   └─ User Register Complete Page
  └─ Authentication
      ├─ Login Handler
      └─ Session Manager
```

### 2.4 シャーディング対応アーキテクチャ

管理画面では全シャードのデータを統合表示する必要がある。

```
一覧表示リクエスト
  ↓
GoAdmin Table Generator
  ↓
Cross-Shard Query Handler
  ├─ GetAllGORMConnections() で全シャード取得
  ├─ goroutine で並列クエリ実行
  │   ├─ Shard 1: SELECT * FROM users
  │   └─ Shard 2: SELECT * FROM users
  ├─ 結果をマージ
  └─ GoAdminに返却
```

**詳細表示・編集・削除**:
- シャードキー（user_id）に基づいて適切なシャードにルーティング
- `GetGORMByKey(user_id)` を使用

**新規作成**:
- シャードキー（user_id）に基づいて適切なシャードにルーティング
- `GetGORMByKey(user_id)` を使用

## 3. データモデル設計

### 3.1 既存モデルの再利用

既存のUserモデルとPostモデルをそのまま使用する。

#### 3.1.1 Userモデル
```go
// server/internal/model/user.go (既存)
type User struct {
    ID        int64     `gorm:"primaryKey" json:"id,string"`
    Name      string    `gorm:"type:varchar(100);not null" json:"name"`
    Email     string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_users_email" json:"email"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
```

#### 3.1.2 Postモデル
```go
// server/internal/model/post.go (既存)
type Post struct {
    ID        int64     `gorm:"primaryKey" json:"id,string"`
    UserID    int64     `gorm:"type:bigint;not null;index:idx_posts_user_id" json:"user_id,string"`
    Title     string    `gorm:"type:varchar(200);not null" json:"title"`
    Content   string    `gorm:"type:text;not null" json:"content"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
```

### 3.2 GoAdminテーブル設定

GoAdminのテーブル生成機能を使用してCRUD操作を自動生成する。

## 4. コンポーネント設計

### 4.1 エントリーポイント（main.go）

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/GoAdminGroup/go-admin/engine"
    "github.com/GoAdminGroup/go-admin/modules/config"
    "github.com/GoAdminGroup/go-admin/plugins/admin"
    "github.com/example/go-webdb-template/internal/config"
    "github.com/example/go-webdb-template/internal/db"
    "github.com/example/go-webdb-template/internal/admin"
)

func main() {
    // 設定ファイルの読み込み
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // GORM Managerの初期化（既存のものを再利用）
    gormManager, err := db.NewGORMManager(cfg)
    if err != nil {
        log.Fatalf("Failed to create GORM manager: %v", err)
    }
    defer gormManager.CloseAll()

    // すべてのShardへの接続確認
    if err := gormManager.PingAll(); err != nil {
        log.Fatalf("Failed to ping databases: %v", err)
    }
    log.Println("Successfully connected to all database shards")

    // GoAdmin Engineの初期化
    eng := engine.Default()

    // GoAdmin設定
    adminCfg := admin.NewConfig(cfg, gormManager)

    // Admin Pluginの登録
    adminPlugin := admin.NewAdmin(admin.Generators{
        // Usersテーブル設定
        "users": admin.GetUserTable(gormManager),
        // Postsテーブル設定
        "posts": admin.GetPostTable(gormManager),
    })

    // カスタムページの登録
    admin.RegisterCustomPages(eng, gormManager)

    // 認証・認可の設定
    admin.RegisterAuth(eng, adminCfg)

    // GoAdmin Engineの初期化
    if err := eng.AddConfig(adminCfg).AddPlugins(adminPlugin).Use(gormManager); err != nil {
        log.Fatalf("Failed to initialize GoAdmin: %v", err)
    }

    // HTTPサーバーの設定
    srv := &http.Server{
        Addr:         fmt.Sprintf(":%d", cfg.Admin.Port),
        Handler:      eng,
        ReadTimeout:  cfg.Server.ReadTimeout,
        WriteTimeout: cfg.Server.WriteTimeout,
    }

    // Graceful shutdown
    go func() {
        log.Printf("Starting admin server on port %d", cfg.Admin.Port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server failed: %v", err)
        }
    }()

    // シグナル待機
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("Shutting down admin server...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Printf("Server forced to shutdown: %v", err)
    }

    log.Println("Admin server exited")
}
```

### 4.2 GoAdmin設定（config.go）

```go
package admin

import (
    "github.com/GoAdminGroup/go-admin/modules/config"
    "github.com/example/go-webdb-template/internal/config"
    "github.com/example/go-webdb-template/internal/db"
)

// Config はGoAdmin設定を管理
type Config struct {
    appConfig *config.Config
    dbManager *db.GORMManager
}

// NewConfig は新しいGoAdmin設定を作成
func NewConfig(cfg *config.Config, dbManager *db.GORMManager) *Config {
    return &Config{
        appConfig: cfg,
        dbManager: dbManager,
    }
}

// GetGoAdminConfig はGoAdmin設定を返す
func (c *Config) GetGoAdminConfig() *config.Config {
    return &config.Config{
        Databases: config.DatabaseList{
            "default": c.getDatabaseConfig(),
        },
        UrlPrefix: "admin",
        Store: config.Store{
            Path:   "./uploads",
            Prefix: "uploads",
        },
        Language: "ja",
        // 認証設定
        Auth: config.Auth{
            SuccessCallback: "/admin/info/user",
            FailureCallback: "/admin/login",
        },
        // セッション設定
        Session: config.Session{
            LifeTime: 7200, // 2時間
        },
    }
}

// getDatabaseConfig はデータベース設定を返す
func (c *Config) getDatabaseConfig() config.Database {
    // 既存のGORM Managerを使用
    // GoAdminのアダプターとして設定
    return config.Database{
        Driver: "gorm",
        // GORM Managerをアダプターとして設定
        Connection: c.dbManager,
    }
}
```

### 4.3 テーブル設定（tables.go）

```go
package admin

import (
    "github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
    "github.com/GoAdminGroup/go-admin/template/types"
    "github.com/example/go-webdb-template/internal/db"
    "github.com/example/go-webdb-template/internal/model"
)

// GetUserTable はUsersテーブルの設定を返す
func GetUserTable(dbManager *db.GORMManager) *table.Generator {
    return table.NewDefaultTable(table.DefaultConfigWithDriver("gorm")).
        SetInfo().
        SetInfo().
        FieldInfo().
        Field("ID", "ID", db.Integer).
        Field("Name", "名前", db.Varchar).
        Field("Email", "メールアドレス", db.Varchar).
        Field("CreatedAt", "作成日時", db.Datetime).
        Field("UpdatedAt", "更新日時", db.Datetime).
        SetTable("users").
        SetTitle("ユーザー管理").
        SetDescription("ユーザー情報の管理").
        // シャーディング対応のクエリ関数を設定
        SetQueryFunc(func(ctx *context.Context) ([]map[string]interface{}, error) {
            return queryAllShards(dbManager, func(conn *db.GORMConnection) ([]model.User, error) {
                var users []model.User
                if err := conn.DB.Find(&users).Error; err != nil {
                    return nil, err
                }
                return users, nil
            })
        }).
        SetInsertFunc(func(ctx *context.Context, values map[string]interface{}) error {
            user := &model.User{
                Name:  values["name"].(string),
                Email: values["email"].(string),
            }
            // シャードキーに基づいて適切なシャードに保存
            gormDB, err := dbManager.GetGORMByKey(user.ID)
            if err != nil {
                return err
            }
            return gormDB.Create(user).Error
        }).
        SetUpdateFunc(func(ctx *context.Context, id int64, values map[string]interface{}) error {
            // 既存のユーザーを取得してシャードを特定
            user, err := findUserAcrossShards(dbManager, id)
            if err != nil {
                return err
            }
            gormDB, err := dbManager.GetGORMByKey(user.ID)
            if err != nil {
                return err
            }
            return gormDB.Model(&user).Updates(values).Error
        }).
        SetDeleteFunc(func(ctx *context.Context, id int64) error {
            // 既存のユーザーを取得してシャードを特定
            user, err := findUserAcrossShards(dbManager, id)
            if err != nil {
                return err
            }
            gormDB, err := dbManager.GetGORMByKey(user.ID)
            if err != nil {
                return err
            }
            return gormDB.Delete(&user).Error
        })
}

// GetPostTable はPostsテーブルの設定を返す
func GetPostTable(dbManager *db.GORMManager) *table.Generator {
    // 同様の実装
    // ...
}
```

### 4.4 カスタムページ設計

#### 4.4.1 最初のページ（home.go）

```go
package pages

import (
    "github.com/GoAdminGroup/go-admin/context"
    "github.com/GoAdminGroup/go-admin/modules/db"
    "github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
    "github.com/example/go-webdb-template/internal/db"
    "github.com/example/go-webdb-template/internal/model"
)

// HomePage は最初のページ（ランディングページ）を返す
func HomePage(ctx *context.Context) (types.Panel, error) {
    dbManager := ctx.UserValue["db_manager"].(*db.GORMManager)

    // 統計情報の取得
    stats := getStatistics(dbManager)

    // HTMLコンテンツの生成
    content := `
    <div class="container">
        <h1>管理画面トップ</h1>
        <div class="stats">
            <div class="stat-item">
                <h3>ユーザー数</h3>
                <p>{{.UserCount}}</p>
            </div>
            <div class="stat-item">
                <h3>投稿数</h3>
                <p>{{.PostCount}}</p>
            </div>
        </div>
        <div class="shard-info">
            <h3>シャード情報</h3>
            {{range .ShardStats}}
            <div class="shard-item">
                <h4>Shard {{.ShardID}}</h4>
                <p>ユーザー数: {{.UserCount}}</p>
                <p>投稿数: {{.PostCount}}</p>
            </div>
            {{end}}
        </div>
        <div class="navigation">
            <a href="/admin/info/users">ユーザー管理</a>
            <a href="/admin/info/posts">投稿管理</a>
            <a href="/admin/custom/user_register">ユーザー登録</a>
        </div>
    </div>
    `

    return types.Panel{
        Content:     template.HTML(content),
        Title:       "管理画面トップ",
        Description: "プロジェクト概要と統計情報",
    }, nil
}

// getStatistics は統計情報を取得
func getStatistics(dbManager *db.GORMManager) map[string]interface{} {
    // 全シャードからデータを取得して集計
    // ...
}
```

#### 4.4.2 ユーザー情報登録画面（user_register.go）

```go
package pages

import (
    "github.com/GoAdminGroup/go-admin/context"
    "github.com/GoAdminGroup/go-admin/template/types"
    "github.com/example/go-webdb-template/internal/db"
    "github.com/example/go-webdb-template/internal/model"
    "github.com/example/go-webdb-template/internal/service"
)

// UserRegisterPage はユーザー情報登録画面を返す
func UserRegisterPage(ctx *context.Context) (types.Panel, error) {
    if ctx.Method() == "POST" {
        // フォーム送信処理
        return handleUserRegister(ctx)
    }

    // GETリクエスト: フォーム表示
    content := `
    <div class="container">
        <h1>ユーザー情報登録</h1>
        <form method="POST" action="/admin/custom/user_register">
            <div class="form-group">
                <label for="name">ユーザー名</label>
                <input type="text" id="name" name="name" required>
            </div>
            <div class="form-group">
                <label for="email">メールアドレス</label>
                <input type="email" id="email" name="email" required>
            </div>
            <button type="submit">登録</button>
        </form>
    </div>
    `

    return types.Panel{
        Content:     template.HTML(content),
        Title:       "ユーザー情報登録",
        Description: "新規ユーザー情報を登録します",
    }, nil
}

// handleUserRegister はユーザー登録処理
func handleUserRegister(ctx *context.Context) (types.Panel, error) {
    name := ctx.PostFormValue("name")
    email := ctx.PostFormValue("email")

    // バリデーション
    if name == "" || email == "" {
        return types.Panel{
            Content: template.HTML("<p>エラー: 必須項目が入力されていません</p>"),
        }, nil
    }

    dbManager := ctx.UserValue["db_manager"].(*db.GORMManager)
    userRepo := repository.NewUserRepositoryGORM(dbManager)
    userService := service.NewUserService(userRepo)

    // ユーザー作成
    user, err := userService.CreateUser(&model.CreateUserRequest{
        Name:  name,
        Email: email,
    })
    if err != nil {
        return types.Panel{
            Content: template.HTML("<p>エラー: " + err.Error() + "</p>"),
        }, nil
    }

    // 登録完了画面にリダイレクト
    ctx.Redirect("/admin/custom/user_register_complete?id=" + fmt.Sprintf("%d", user.ID))
    return types.Panel{}, nil
}
```

#### 4.4.3 ユーザー情報登録完了画面（user_register_complete.go）

```go
package pages

import (
    "github.com/GoAdminGroup/go-admin/context"
    "github.com/GoAdminGroup/go-admin/template/types"
    "github.com/example/go-webdb-template/internal/db"
)

// UserRegisterCompletePage はユーザー情報登録完了画面を返す
func UserRegisterCompletePage(ctx *context.Context) (types.Panel, error) {
    userID := ctx.Query("id")
    
    dbManager := ctx.UserValue["db_manager"].(*db.GORMManager)
    // ユーザー情報を取得
    user := getUserByID(dbManager, userID)

    content := `
    <div class="container">
        <h1>ユーザー情報登録完了</h1>
        <div class="success-message">
            <p>ユーザー情報の登録が完了しました。</p>
        </div>
        <div class="user-info">
            <h3>登録されたユーザー情報</h3>
            <p>ID: {{.ID}}</p>
            <p>名前: {{.Name}}</p>
            <p>メールアドレス: {{.Email}}</p>
            <p>登録日時: {{.CreatedAt}}</p>
        </div>
        <div class="actions">
            <a href="/admin/info/users">一覧に戻る</a>
            <a href="/admin/custom/user_register">新規登録を続ける</a>
        </div>
    </div>
    `

    return types.Panel{
        Content:     template.HTML(content),
        Title:       "ユーザー情報登録完了",
        Description: "ユーザー情報の登録が成功しました",
    }, nil
}
```

### 4.5 認証・認可設計

#### 4.5.1 認証設定（auth.go）

```go
package admin

import (
    "github.com/GoAdminGroup/go-admin/modules/auth"
    "github.com/GoAdminGroup/go-admin/plugins/admin/modules/user"
    "golang.org/x/crypto/bcrypt"
)

// LoginHandler はログイン処理を行う
func LoginHandler(ctx *context.Context) {
    username := ctx.PostFormValue("username")
    password := ctx.PostFormValue("password")

    // 認証情報の検証（設定ファイルから読み込み）
    if validateCredentials(username, password) {
        // セッション作成
        session := auth.NewSession(ctx)
        session.Set("user", user.User{
            Id:       1,
            UserName: username,
            Name:     "Admin",
            Avatar:   "",
            Level:    "admin",
        })
        ctx.Redirect("/admin/info/user")
    } else {
        ctx.JSON(http.StatusUnauthorized, map[string]string{
            "error": "認証に失敗しました",
        })
    }
}

// validateCredentials は認証情報を検証
func validateCredentials(username, password string) bool {
    // 設定ファイルから認証情報を読み込み
    // 開発環境では簡易的な認証でも可
    // 本番環境では適切な認証を実装
    expectedUsername := os.Getenv("ADMIN_USERNAME")
    expectedPassword := os.Getenv("ADMIN_PASSWORD")

    if username == expectedUsername {
        // パスワードのハッシュ化（bcrypt）
        err := bcrypt.CompareHashAndPassword([]byte(expectedPassword), []byte(password))
        return err == nil
    }
    return false
}

// LogoutHandler はログアウト処理を行う
func LogoutHandler(ctx *context.Context) {
    session := auth.GetSession(ctx)
    session.Clear()
    ctx.Redirect("/admin/login")
}
```

#### 4.5.2 セッション管理

GoAdminのセッション機能を使用する。

```go
package admin

import (
    "github.com/GoAdminGroup/go-admin/modules/auth"
)

// RegisterAuth は認証・認可を登録
func RegisterAuth(eng *engine.Engine, cfg *Config) {
    eng.AddAuthFunc(func(ctx *context.Context) bool {
        // セッションチェック
        session := auth.GetSession(ctx)
        user := session.Get("user")
        return user != nil
    })
}
```

## 5. シャーディング対応設計

### 5.1 クロスシャードクエリ実装

全シャードのデータを取得するためのヘルパー関数を実装する。

```go
package admin

import (
    "sync"
    "github.com/example/go-webdb-template/internal/db"
    "github.com/example/go-webdb-template/internal/model"
)

// queryAllShards は全シャードに対してクエリを実行し、結果をマージ
func queryAllShards(dbManager *db.GORMManager, queryFunc func(*db.GORMConnection) ([]model.User, error)) ([]model.User, error) {
    connections := dbManager.GetAllGORMConnections()
    
    var wg sync.WaitGroup
    results := make([][]model.User, len(connections))
    errors := make([]error, len(connections))

    // 並列実行
    for i, conn := range connections {
        wg.Add(1)
        go func(idx int, c *db.GORMConnection) {
            defer wg.Done()
            users, err := queryFunc(c)
            if err != nil {
                errors[idx] = err
                return
            }
            results[idx] = users
        }(i, conn)
    }

    wg.Wait()

    // エラーチェック
    for _, err := range errors {
        if err != nil {
            return nil, err
        }
    }

    // 結果をマージ
    var allUsers []model.User
    for _, users := range results {
        allUsers = append(allUsers, users...)
    }

    return allUsers, nil
}

// findUserAcrossShards は全シャードからユーザーを検索
func findUserAcrossShards(dbManager *db.GORMManager, userID int64) (*model.User, error) {
    connections := dbManager.GetAllGORMConnections()
    
    var wg sync.WaitGroup
    resultChan := make(chan *model.User, 1)
    errorChan := make(chan error, len(connections))

    // 並列検索
    for _, conn := range connections {
        wg.Add(1)
        go func(c *db.GORMConnection) {
            defer wg.Done()
            var user model.User
            if err := c.DB.Where("id = ?", userID).First(&user).Error; err != nil {
                if !errors.Is(err, gorm.ErrRecordNotFound) {
                    errorChan <- err
                }
                return
            }
            select {
            case resultChan <- &user:
            default:
            }
        }(conn)
    }

    // 最初に見つかったユーザーを返す
    go func() {
        wg.Wait()
        close(resultChan)
        close(errorChan)
    }()

    select {
    case user := <-resultChan:
        return user, nil
    case err := <-errorChan:
        return nil, err
    default:
        return nil, gorm.ErrRecordNotFound
    }
}
```

### 5.2 シャードキーに基づくルーティング

新規作成・更新・削除時はシャードキーに基づいて適切なシャードにルーティングする。

```go
// 新規作成時のシャード選択
func createUser(dbManager *db.GORMManager, user *model.User) error {
    // 一時的なIDを生成してシャードを決定
    // 実際のIDはデータベースで自動生成されるため、
    // 事前にシャードを決定する必要がある
    // または、デフォルトでShard 1に保存し、後でリバランス
    
    // 簡易実装: デフォルトでShard 1に保存
    gormDB, err := dbManager.GetGORM(1)
    if err != nil {
        return err
    }
    return gormDB.Create(user).Error
}
```

## 6. 設定管理設計

### 6.1 設定ファイル構造

既存の設定ファイルに管理画面用の設定を追加する。

```yaml
# config/develop.yaml
server:
  port: 8080
  read_timeout: 30s
  write_timeout: 30s

admin:
  port: 8081
  read_timeout: 30s
  write_timeout: 30s
  auth:
    username: admin
    password: ${ADMIN_PASSWORD}  # 環境変数から読み込み
  session:
    lifetime: 7200  # 2時間（秒）

database:
  shards:
    # 既存の設定
    ...
```

### 6.2 設定構造体の拡張

```go
// server/internal/config/config.go
type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Admin    AdminConfig    `mapstructure:"admin"`  // 新規追加
    Database DatabaseConfig `mapstructure:"database"`
    Logging  LoggingConfig  `mapstructure:"logging"`
    CORS     CORSConfig     `mapstructure:"cors"`
}

// AdminConfig は管理画面設定
type AdminConfig struct {
    Port        int           `mapstructure:"port"`
    ReadTimeout time.Duration `mapstructure:"read_timeout"`
    WriteTimeout time.Duration `mapstructure:"write_timeout"`
    Auth        AuthConfig    `mapstructure:"auth"`
    Session     SessionConfig `mapstructure:"session"`
}

// AuthConfig は認証設定
type AuthConfig struct {
    Username string `mapstructure:"username"`
    Password string `mapstructure:"password"`
}

// SessionConfig はセッション設定
type SessionConfig struct {
    Lifetime int `mapstructure:"lifetime"`
}
```

## 7. エラーハンドリング設計

### 7.1 エラー分類

- **データベースエラー**: 接続エラー、クエリエラー、制約違反エラー
- **認証エラー**: 認証失敗、セッションタイムアウト
- **バリデーションエラー**: 入力値の検証エラー
- **シャーディングエラー**: シャード接続エラー、クロスシャードクエリエラー

### 7.2 エラーハンドリングパターン

```go
// エラーハンドリングの例
func handleError(ctx *context.Context, err error) {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        ctx.JSON(http.StatusNotFound, map[string]string{
            "error": "レコードが見つかりません",
        })
        return
    }
    
    if errors.Is(err, gorm.ErrDuplicatedKey) {
        ctx.JSON(http.StatusConflict, map[string]string{
            "error": "重複するレコードが存在します",
        })
        return
    }
    
    // その他のエラー
    log.Printf("Error: %v", err)
    ctx.JSON(http.StatusInternalServerError, map[string]string{
        "error": "内部エラーが発生しました",
    })
}
```

## 8. セキュリティ設計

### 8.1 認証・認可

- **認証**: ユーザー名とパスワードによる認証
- **認可**: 認証済みユーザーのみ管理画面にアクセス可能
- **セッション管理**: GoAdminのセッション機能を使用
- **パスワードハッシュ化**: bcryptを使用

### 8.2 セキュリティ対策

- **CSRF対策**: GoAdminのCSRF対策機能を使用
- **XSS対策**: 入力値のサニタイゼーション
- **SQLインジェクション対策**: GORMのパラメータ化クエリを使用
- **HTTPS対応**: 本番環境ではHTTPSを必須とする

## 9. テスト戦略

### 9.1 テストレベル

1. **ユニットテスト**: 各コンポーネントの単体テスト
2. **統合テスト**: GoAdminとGORM Managerの統合テスト
3. **E2Eテスト**: 管理画面の操作フローのテスト

### 9.2 テストカバレッジ

- 目標: 70%以上
- 重点領域:
  - 認証・認可機能
  - カスタムページ
  - シャーディング対応機能

## 10. 実装上の注意事項

### 10.1 GoAdminの統合

- GoAdminの最新バージョンを使用
- GORMアダプターの設定を正しく行う
- テーブル生成機能を活用してCRUD操作を自動化

### 10.2 シャーディング対応

- 全シャードのデータを取得する際は並列実行でパフォーマンスを維持
- シャードキーに基づくルーティングを正しく実装
- クロスシャードクエリの結果マージを適切に実装

### 10.3 認証・認可

- 開発環境では簡易的な認証でも可
- 本番環境では適切な認証を実装
- セッション管理を適切に実装

### 10.4 カスタムページ

- GoAdminのカスタムページ機能を活用
- HTMLテンプレートはGoAdminのテンプレート機能を使用
- フォーム送信は既存のService層を経由

## 11. 参考情報

### 11.1 GoAdmin公式ドキュメント
- GoAdmin公式サイト: https://github.com/GoAdminGroup/go-admin
- GoAdminドキュメント: https://book.go-admin.cn/
- GoAdmin GORM統合: https://book.go-admin.cn/guide/admin/plugins/gorm

### 11.2 既存実装
- `server/internal/db/manager.go`: GORM接続管理
- `server/internal/model/user.go`: Userモデル定義
- `server/internal/model/post.go`: Postモデル定義
- `server/cmd/server/main.go`: メインサービスエントリーポイント

