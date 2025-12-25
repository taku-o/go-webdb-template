# CLIツール対応設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、バッチ処理実行のためのCLIツール基盤の詳細設計を定義する。サンプルとしてユーザー一覧を出力するCLIツールを実装し、今後のバッチ処理実装の参考とする。

### 1.2 設計の範囲
- CLIツールのディレクトリ構造設計
- コマンドライン引数処理の設計
- 設定読み込みとDB接続の設計
- 出力形式の設計
- エラーハンドリング設計
- ビルドと実行ファイル生成の設計
- テスト戦略

## 2. アーキテクチャ設計

### 2.1 ディレクトリ構造

#### 2.1.1 変更前の構造
```
server/
├── cmd/
│   ├── server/
│   │   └── main.go          # サーバー起動コマンド
│   └── admin/
│       └── main.go          # 管理画面起動コマンド
├── internal/
│   ├── config/
│   ├── db/
│   ├── repository/
│   └── service/
└── ...
```

#### 2.1.2 変更後の構造
```
server/
├── cmd/
│   ├── server/
│   │   └── main.go          # サーバー起動コマンド
│   ├── admin/
│   │   └── main.go          # 管理画面起動コマンド
│   └── list-users/
│       └── main.go          # ユーザー一覧出力CLIツール（新規）
├── bin/                     # 実行ファイル生成先（新規、.gitignoreに追加）
│   └── list-users           # ビルド後の実行ファイル
├── internal/
│   ├── config/
│   ├── db/
│   ├── repository/
│   └── service/
└── ...
```

### 2.2 CLIツールの実行フロー

```
┌─────────────────────────────────────────────────────────────┐
│                    list-users コマンド実行                    │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              1. コマンドライン引数の解析                      │
│              - flag.Parse()                                 │
│              - --limit フラグの取得                          │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              2. 設定ファイルの読み込み                        │
│              - config.Load()                                │
│              - APP_ENV 環境変数の取得                        │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              3. DB接続の初期化                               │
│              - db.NewGORMManager(cfg)                        │
│              - gormManager.PingAll()                         │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              4. Repository層の初期化                         │
│              - repository.NewUserRepositoryGORM()           │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              5. Service層の初期化                           │
│              - service.NewUserService()                    │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              6. ユーザー一覧の取得                           │
│              - userService.ListUsers(ctx, limit, 0)         │
│              - クロスシャードクエリで全シャードから取得       │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              7. TSV形式での出力                              │
│              - ヘッダー行の出力                              │
│              - 各ユーザー情報の出力                          │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              8. リソースのクリーンアップ                      │
│              - gormManager.CloseAll()                       │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              9. 正常終了（終了コード: 0）                    │
└─────────────────────────────────────────────────────────────┘
```

### 2.3 既存アーキテクチャとの統合

CLIツールは既存のアーキテクチャを再利用する：

```
┌─────────────────────────────────────────────────────────────┐
│                    CLIツール (cmd/list-users)                │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Service層 (internal/service)                    │
│              - UserService.ListUsers()                      │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Repository層 (internal/repository)              │
│              - UserRepositoryGORM.List()                    │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              DB層 (internal/db)                              │
│              - GORMManager.GetAllGORMConnections()          │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              データベース（複数シャード）                    │
└─────────────────────────────────────────────────────────────┘
```

## 3. コンポーネント設計

### 3.1 main.goの設計

#### 3.1.1 パッケージ構造
```go
package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "os"
    
    "github.com/example/go-webdb-template/internal/config"
    "github.com/example/go-webdb-template/internal/db"
    "github.com/example/go-webdb-template/internal/repository"
    "github.com/example/go-webdb-template/internal/service"
)
```

#### 3.1.2 main関数の設計

```go
func main() {
    // 1. コマンドライン引数の解析
    limit := flag.Int("limit", 20, "Number of users to output (default: 20, max: 100)")
    flag.Parse()
    
    // 2. 引数のバリデーション
    if *limit < 1 {
        log.Fatalf("Error: limit must be at least 1")
    }
    if *limit > 100 {
        *limit = 100
        log.Printf("Warning: limit exceeds maximum (100), using 100")
    }
    
    // 3. 設定ファイルの読み込み
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // 4. DB接続の初期化
    gormManager, err := db.NewGORMManager(cfg)
    if err != nil {
        log.Fatalf("Failed to create GORM manager: %v", err)
    }
    defer gormManager.CloseAll()
    
    // 5. 接続確認
    if err := gormManager.PingAll(); err != nil {
        log.Fatalf("Failed to ping databases: %v", err)
    }
    
    // 6. Repository層の初期化
    userRepo := repository.NewUserRepositoryGORM(gormManager)
    
    // 7. Service層の初期化
    userService := service.NewUserService(userRepo)
    
    // 8. ユーザー一覧の取得
    ctx := context.Background()
    users, err := userService.ListUsers(ctx, *limit, 0)
    if err != nil {
        log.Fatalf("Failed to list users: %v", err)
    }
    
    // 9. TSV形式での出力
    printUsersTSV(users)
    
    // 10. 正常終了
    os.Exit(0)
}
```

### 3.2 出力処理の設計

#### 3.2.1 printUsersTSV関数の設計

```go
func printUsersTSV(users []*model.User) {
    // ヘッダー行の出力
    fmt.Println("ID\tName\tEmail\tCreatedAt\tUpdatedAt")
    
    // 各ユーザー情報の出力
    for _, user := range users {
        fmt.Printf("%d\t%s\t%s\t%s\t%s\n",
            user.ID,
            user.Name,
            user.Email,
            user.CreatedAt.Format(time.RFC3339),
            user.UpdatedAt.Format(time.RFC3339),
        )
    }
}
```

#### 3.2.2 出力形式の詳細

- **区切り文字**: タブ文字（`\t`）
- **ヘッダー行**: 最初の1行目に出力
- **データ行**: 各ユーザー情報を1行ずつ出力
- **日時形式**: RFC3339形式（例: `2025-01-27T10:30:00Z`）
- **エスケープ**: タブ文字や改行文字が含まれる場合は、そのまま出力（TSV形式の標準的な扱い）

### 3.3 コマンドライン引数処理の設計

#### 3.3.1 flagパッケージの使用

```go
limit := flag.Int("limit", 20, "Number of users to output (default: 20, max: 100)")
flag.Parse()
```

#### 3.3.2 引数のバリデーション

| チェック項目 | 条件 | 処理 |
|------------|------|------|
| limit < 1 | 最小値チェック | エラーメッセージを出力し、`log.Fatalf()`で終了 |
| limit > 100 | 最大値チェック | 100に制限し、警告メッセージを出力 |

#### 3.3.3 使用方法の表示

```go
flag.Usage = func() {
    fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
    fmt.Fprintf(os.Stderr, "Options:\n")
    flag.PrintDefaults()
}
```

## 4. データモデル

### 4.1 出力データ構造

#### 4.1.1 TSV形式の構造

```
ID	Name	Email	CreatedAt	UpdatedAt
1234567890123456789	John Doe	john@example.com	2025-01-27T10:30:00Z	2025-01-27T10:30:00Z
1234567890123456790	Jane Smith	jane@example.com	2025-01-27T11:00:00Z	2025-01-27T11:00:00Z
```

#### 4.1.2 データ型

| フィールド | 型 | 説明 |
|----------|-----|------|
| ID | int64 | ユーザーID（タイムスタンプベース） |
| Name | string | ユーザー名 |
| Email | string | メールアドレス |
| CreatedAt | time.Time | 作成日時（RFC3339形式） |
| UpdatedAt | time.Time | 更新日時（RFC3339形式） |

### 4.2 クロスシャードクエリの動作

既存の`UserRepositoryGORM.List()`メソッドを使用：

```go
// 各シャードから並列にデータを取得
for _, conn := range connections {
    var shardUsers []*model.User
    if err := conn.DB.WithContext(ctx).
        Order("id").
        Limit(limit).
        Offset(offset).
        Find(&shardUsers).Error; err != nil {
        return nil, fmt.Errorf("failed to query shard %d: %w", conn.ShardID, err)
    }
    users = append(users, shardUsers...)
}
```

**注意**: 現在の実装では、各シャードから`limit`件ずつ取得してマージするため、実際の出力件数は`limit * シャード数`になる可能性がある。要件定義書では「limitパラメータで出力件数を制限」としているため、マージ後に`limit`件に制限する処理を追加する必要がある。

## 5. エラーハンドリング

### 5.1 コマンドライン引数のエラー

**エラーケース**:
- `limit`が1未満

**処理**:
```go
if *limit < 1 {
    log.Fatalf("Error: limit must be at least 1")
    os.Exit(1)
}
```

### 5.2 設定読み込みのエラー

**エラーケース**:
- 設定ファイルが存在しない
- 設定ファイルの形式が不正
- 環境変数`APP_ENV`が不正な値

**処理**:
```go
cfg, err := config.Load()
if err != nil {
    log.Fatalf("Failed to load config: %v", err)
    os.Exit(1)
}
```

### 5.3 DB接続のエラー

**エラーケース**:
- DB接続の初期化に失敗
- 接続確認（Ping）に失敗

**処理**:
```go
gormManager, err := db.NewGORMManager(cfg)
if err != nil {
    log.Fatalf("Failed to create GORM manager: %v", err)
    os.Exit(1)
}

if err := gormManager.PingAll(); err != nil {
    log.Fatalf("Failed to ping databases: %v", err)
    os.Exit(1)
}
```

### 5.4 クエリ実行のエラー

**エラーケース**:
- ユーザー一覧の取得に失敗
- クロスシャードクエリの実行に失敗

**処理**:
```go
users, err := userService.ListUsers(ctx, *limit, 0)
if err != nil {
    log.Fatalf("Failed to list users: %v", err)
    os.Exit(1)
}
```

### 5.5 出力処理のエラー

**エラーケース**:
- 標準出力への書き込みに失敗（通常は発生しない）

**処理**:
- Go言語の標準ライブラリは自動的にエラーを処理するため、明示的なエラーハンドリングは不要

### 5.6 終了コード

| 状況 | 終了コード | 説明 |
|------|----------|------|
| 正常終了 | 0 | ユーザー一覧が正常に出力された |
| 引数エラー | 1 | コマンドライン引数が不正 |
| 設定エラー | 1 | 設定ファイルの読み込みに失敗 |
| DB接続エラー | 1 | DB接続の初期化または接続確認に失敗 |
| クエリエラー | 1 | ユーザー一覧の取得に失敗 |

## 6. ビルドと実行ファイル生成

### 6.1 ビルドコマンド

#### 6.1.1 開発環境でのビルド

```bash
cd server
go build -o bin/list-users ./cmd/list-users
```

#### 6.1.2 本番環境でのビルド（クロスコンパイル）

```bash
cd server
GOOS=linux GOARCH=amd64 go build -o bin/list-users ./cmd/list-users
```

#### 6.1.3 リリースビルド（最適化）

```bash
cd server
go build -ldflags="-s -w" -o bin/list-users ./cmd/list-users
```

### 6.2 実行ファイルの配置

- **生成先**: `server/bin/list-users`
- **実行方法**: `./bin/list-users [options]`
- **環境変数**: `APP_ENV=production ./bin/list-users --limit 100`

### 6.3 .gitignoreの更新

```gitignore
# 実行ファイル
server/bin/
```

## 7. テスト戦略

### 7.1 ユニットテスト

#### 7.1.1 printUsersTSV関数のテスト

**テストケース**:
1. **正常系**: ユーザー一覧が正常にTSV形式で出力される
2. **正常系**: 空のユーザー一覧の場合、ヘッダーのみが出力される
3. **正常系**: 日時がRFC3339形式で出力される

**テスト実装場所**:
- `server/cmd/list-users/main_test.go`

**テスト実装例**:
```go
func TestPrintUsersTSV(t *testing.T) {
    users := []*model.User{
        {
            ID:        1234567890123456789,
            Name:      "John Doe",
            Email:     "john@example.com",
            CreatedAt: time.Date(2025, 1, 27, 10, 30, 0, 0, time.UTC),
            UpdatedAt: time.Date(2025, 1, 27, 10, 30, 0, 0, time.UTC),
        },
    }
    
    // 標準出力をキャプチャ
    // ... (実装)
    
    printUsersTSV(users)
    
    // 出力内容の検証
    // ... (実装)
}
```

### 7.2 統合テスト

#### 7.2.1 CLIツールの統合テスト

**テストケース**:
1. **正常系**: デフォルト引数で実行した場合、20件のユーザーが出力される
2. **正常系**: `--limit 10`で実行した場合、10件のユーザーが出力される
3. **正常系**: 環境変数`APP_ENV`で環境を切り替えられる
4. **異常系**: `--limit 0`で実行した場合、エラーが出力される
5. **異常系**: `--limit 200`で実行した場合、100件に制限される

**テスト実装場所**:
- `server/cmd/list-users/integration_test.go`

**テスト実装例**:
```go
func TestListUsersCLI(t *testing.T) {
    // テスト用の設定とDBを使用
    // ... (実装)
    
    // CLIツールを実行
    cmd := exec.Command("./bin/list-users", "--limit", "10")
    output, err := cmd.Output()
    
    // 出力内容の検証
    // ... (実装)
}
```

### 7.3 E2Eテスト

#### 7.3.1 cron実行シミュレーションテスト

**テストケース**:
1. 非対話的な実行が可能であること
2. 終了コードが適切に返されること
3. 標準出力と標準エラー出力が適切に分離されていること

**テスト実装場所**:
- 既存のE2Eテストスイートに追加

## 8. 実装上の注意事項

### 8.1 クロスシャードクエリのlimit処理

**問題点**:
- 既存の`UserRepositoryGORM.List()`は各シャードから`limit`件ずつ取得してマージするため、実際の出力件数は`limit * シャード数`になる可能性がある

**解決策**:
- マージ後に`limit`件に制限する処理を追加する
- または、`UserService.ListUsers()`の実装を確認し、必要に応じて修正する

**実装例**:
```go
users, err := userService.ListUsers(ctx, *limit, 0)
if err != nil {
    log.Fatalf("Failed to list users: %v", err)
}

// limit件に制限
if len(users) > *limit {
    users = users[:*limit]
}
```

### 8.2 リソース管理

**重要なポイント**:
- `defer gormManager.CloseAll()`でDB接続を確実にクローズする
- エラー発生時もリソースが適切に解放されるようにする

### 8.3 既存コードの再利用

**重要なポイント**:
- 既存の`config.Load()`、`db.NewGORMManager()`、`repository.NewUserRepositoryGORM()`、`service.NewUserService()`をそのまま使用する
- 新しいコードを追加せず、既存のアーキテクチャを維持する

### 8.4 出力形式の一貫性

**重要なポイント**:
- TSV形式の標準に従う（タブ区切り、改行で行を区切る）
- ヘッダー行を必ず含める
- 日時はRFC3339形式で統一する

### 8.5 エラーメッセージの明確性

**重要なポイント**:
- エラーメッセージは標準エラー出力（`log.Fatalf()`）に出力する
- エラーの原因が明確に分かるメッセージを出力する
- 使用方法（`flag.Usage()`）を表示する

## 9. 参考情報

### 9.1 関連Issue
- GitHub Issue #9: バッチ処理対応

### 9.2 既存ドキュメント
- `server/cmd/server/main.go`: サーバー起動コマンドの実装例
- `server/cmd/admin/main.go`: 管理画面起動コマンドの実装例
- `server/internal/config/config.go`: 設定読み込み実装
- `server/internal/db/manager.go`: DB接続管理実装
- `server/internal/service/user_service.go`: ユーザーService実装

### 9.3 既存実装
- `server/internal/repository/user_repository_gorm.go`: ユーザーRepository実装（GORM版）
- `server/internal/model/user.go`: ユーザーモデル定義

### 9.4 Go言語標準ライブラリ
- `flag`パッケージ: https://pkg.go.dev/flag
- `os`パッケージ: https://pkg.go.dev/os
- `context`パッケージ: https://pkg.go.dev/context
- `log`パッケージ: https://pkg.go.dev/log

