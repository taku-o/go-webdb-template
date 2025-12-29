# Gofakeitによる開発用サンプルデータ生成機能設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、Gofakeitライブラリを使用した開発用サンプルデータ生成機能の詳細設計を定義する。コマンドライン実行方式で、users、posts、newsテーブルに各100件程度のデータを生成する機能を実装する。

### 1.2 設計の範囲
- CLIツールのディレクトリ構造設計
- データ生成ロジックの詳細設計
- シャーディングテーブルへの分散ロジック設計
- バッチ処理の実装設計
- エラーハンドリング設計
- 既存アーキテクチャとの統合設計
- テスト戦略

### 1.3 設計方針
- **既存パターンの再利用**: 既存のCLI実装パターン（`server/cmd/list-users/`）に従う
- **既存アーキテクチャの活用**: 既存の設定読み込み、DB接続管理機能を再利用
- **パフォーマンス重視**: バッチ挿入による効率的なデータ生成
- **エラーハンドリング**: 適切なエラーメッセージと終了コード
- **進捗表示**: データ生成の進捗状況を適切に表示

## 2. アーキテクチャ設計

### 2.1 ディレクトリ構造

#### 2.1.1 変更前の構造
```
server/
├── cmd/
│   ├── server/
│   │   └── main.go          # サーバー起動コマンド
│   ├── admin/
│   │   └── main.go          # 管理画面起動コマンド
│   └── list-users/
│       └── main.go          # ユーザー一覧出力CLIツール
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
│   ├── list-users/
│   │   └── main.go          # ユーザー一覧出力CLIツール
│   └── generate-sample-data/
│       └── main.go          # サンプルデータ生成CLIツール（新規）
├── bin/                     # 実行ファイル生成先（既存）
│   ├── list-users           # ビルド後の実行ファイル
│   └── generate-sample-data # ビルド後の実行ファイル（新規）
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
│            generate-sample-data コマンド実行                  │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              1. 設定ファイルの読み込み                        │
│              - config.Load()                                │
│              - APP_ENV 環境変数の取得                        │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              2. GroupManagerの初期化                         │
│              - db.NewGroupManager(cfg)                        │
│              - groupManager.PingAll()                        │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              3. usersテーブルへのデータ生成                  │
│              - 32分割テーブル（users_000～users_031）       │
│              - 各テーブルに約3～4件ずつ生成                  │
│              - バッチ挿入（500件ずつ）                       │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              4. postsテーブルへのデータ生成                  │
│              - 32分割テーブル（posts_000～posts_031）       │
│              - 各テーブルに約3～4件ずつ生成                  │
│              - user_idは既存のusersテーブルから参照         │
│              - バッチ挿入（500件ずつ）                       │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              5. newsテーブルへのデータ生成                    │
│              - master DBのnewsテーブル                       │
│              - 100件を生成                                   │
│              - バッチ挿入（500件ずつ）                       │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              6. 生成完了メッセージの表示                     │
│              - 各テーブルの生成件数を表示                    │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              7. リソースのクリーンアップ                      │
│              - groupManager.CloseAll()                       │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              8. 正常終了（終了コード: 0）                    │
└─────────────────────────────────────────────────────────────┘
```

### 2.3 既存アーキテクチャとの統合

CLIツールは既存のアーキテクチャを直接利用する：

```
┌─────────────────────────────────────────────────────────────┐
│          CLIツール (cmd/generate-sample-data)                │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Config層 (internal/config)                      │
│              - config.Load()                                │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              DB層 (internal/db)                              │
│              - GroupManager                                 │
│              - GetMasterConnection()                        │
│              - GetShardingConnectionByTableNumber()         │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              データベース                                    │
│              - master DB (news)                             │
│              - sharding DB 1-4 (users, posts)               │
└─────────────────────────────────────────────────────────────┘
```

## 3. コンポーネント設計

### 3.1 main.goの設計

#### 3.1.1 パッケージ構造
```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"
    
    "github.com/brianvoe/gofakeit/v6"
    "github.com/taku-o/go-webdb-template/internal/config"
    "github.com/taku-o/go-webdb-template/internal/db"
    "github.com/taku-o/go-webdb-template/internal/model"
)
```

#### 3.1.2 main関数の設計

```go
func main() {
    // 1. 設定ファイルの読み込み
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // 2. GroupManagerの初期化
    groupManager, err := db.NewGroupManager(cfg)
    if err != nil {
        log.Fatalf("Failed to create group manager: %v", err)
    }
    defer groupManager.CloseAll()
    
    // 3. データベース接続確認
    if err := groupManager.PingAll(); err != nil {
        log.Fatalf("Failed to ping databases: %v", err)
    }
    
    // 4. usersテーブルへのデータ生成
    userIDs, err := generateUsers(groupManager, 100)
    if err != nil {
        log.Fatalf("Failed to generate users: %v", err)
    }
    
    // 5. postsテーブルへのデータ生成
    if err := generatePosts(groupManager, userIDs, 100); err != nil {
        log.Fatalf("Failed to generate posts: %v", err)
    }
    
    // 6. newsテーブルへのデータ生成
    if err := generateNews(groupManager, 100); err != nil {
        log.Fatalf("Failed to generate news: %v", err)
    }
    
    // 7. 生成完了メッセージ
    log.Println("Sample data generation completed successfully")
    os.Exit(0)
}
```

### 3.2 データ生成関数の設計

#### 3.2.1 generateUsers関数

```go
// generateUsers はusersテーブルにデータを生成
// 戻り値: 生成されたuser_idのリスト（posts生成時に使用）
func generateUsers(groupManager *db.GroupManager, totalCount int) ([]int64, error) {
    const batchSize = 500
    const tableCount = 32
    countPerTable := totalCount / tableCount // 各テーブルに約3～4件
    
    var allUserIDs []int64
    ctx := context.Background()
    
    // 各テーブル（0～31）に対してデータ生成
    for tableNumber := 0; tableNumber < tableCount; tableNumber++ {
        // 接続を取得
        conn, err := groupManager.GetShardingConnection(tableNumber)
        if err != nil {
            return nil, fmt.Errorf("failed to get connection for table %d: %w", tableNumber, err)
        }
        
        // テーブル名を生成
        tableName := fmt.Sprintf("users_%03d", tableNumber)
        
        // バッチでデータ生成
        var batch []*model.User
        for i := 0; i < countPerTable; i++ {
            user := &model.User{
                Name:      gofakeit.Name(),
                Email:     gofakeit.Email(),
                CreatedAt: time.Now(),
                UpdatedAt: time.Now(),
            }
            batch = append(batch, user)
            
            // バッチサイズに達したら挿入
            if len(batch) >= batchSize {
                if err := insertBatch(conn.DB, tableName, batch); err != nil {
                    return nil, fmt.Errorf("failed to insert batch to %s: %w", tableName, err)
                }
                // 生成されたIDを取得（AUTOINCREMENT）
                // 注意: SQLiteではLastInsertId()を使用
                batch = nil
            }
        }
        
        // 残りのデータを挿入
        if len(batch) > 0 {
            if err := insertBatch(conn.DB, tableName, batch); err != nil {
                return nil, fmt.Errorf("failed to insert batch to %s: %w", tableName, err)
            }
        }
        
        log.Printf("Generated %d users in %s", countPerTable, tableName)
    }
    
    return allUserIDs, nil
}
```

#### 3.2.2 generatePosts関数

```go
// generatePosts はpostsテーブルにデータを生成
// userIDs: 既存のusersテーブルから取得したuser_idのリスト
func generatePosts(groupManager *db.GroupManager, userIDs []int64, totalCount int) error {
    const batchSize = 500
    const tableCount = 32
    countPerTable := totalCount / tableCount // 各テーブルに約3～4件
    
    ctx := context.Background()
    
    // 各テーブル（0～31）に対してデータ生成
    for tableNumber := 0; tableNumber < tableCount; tableNumber++ {
        // 接続を取得
        conn, err := groupManager.GetShardingConnection(tableNumber)
        if err != nil {
            return fmt.Errorf("failed to get connection for table %d: %w", tableNumber, err)
        }
        
        // テーブル名を生成
        tableName := fmt.Sprintf("posts_%03d", tableNumber)
        
        // バッチでデータ生成
        var batch []*model.Post
        for i := 0; i < countPerTable; i++ {
            // user_idをランダムに選択
            userID := userIDs[gofakeit.IntRange(0, len(userIDs)-1)]
            
            post := &model.Post{
                UserID:    userID,
                Title:     gofakeit.Sentence(5),
                Content:   gofakeit.Paragraph(3, 5, 10, "\n"),
                CreatedAt: time.Now(),
                UpdatedAt: time.Now(),
            }
            batch = append(batch, post)
            
            // バッチサイズに達したら挿入
            if len(batch) >= batchSize {
                if err := insertBatch(conn.DB, tableName, batch); err != nil {
                    return fmt.Errorf("failed to insert batch to %s: %w", tableName, err)
                }
                batch = nil
            }
        }
        
        // 残りのデータを挿入
        if len(batch) > 0 {
            if err := insertBatch(conn.DB, tableName, batch); err != nil {
                return fmt.Errorf("failed to insert batch to %s: %w", tableName, err)
            }
        }
        
        log.Printf("Generated %d posts in %s", countPerTable, tableName)
    }
    
    return nil
}
```

#### 3.2.3 generateNews関数

```go
// generateNews はnewsテーブルにデータを生成
func generateNews(groupManager *db.GroupManager, totalCount int) error {
    const batchSize = 500
    
    // master接続を取得
    conn, err := groupManager.GetMasterConnection()
    if err != nil {
        return fmt.Errorf("failed to get master connection: %w", err)
    }
    
    // バッチでデータ生成
    var batch []*model.News
    for i := 0; i < totalCount; i++ {
        news := &model.News{
            Title:       gofakeit.Sentence(5),
            Content:     gofakeit.Paragraph(3, 5, 10, "\n"),
            AuthorID:    func() *int64 { id := gofakeit.Int64(); return &id }(),
            PublishedAt: func() *time.Time { t := gofakeit.Date(); return &t }(),
            CreatedAt:   time.Now(),
            UpdatedAt:   time.Now(),
        }
        batch = append(batch, news)
        
        // バッチサイズに達したら挿入
        if len(batch) >= batchSize {
            if err := insertBatch(conn.DB, "news", batch); err != nil {
                return fmt.Errorf("failed to insert batch to news: %w", err)
            }
            batch = nil
        }
    }
    
    // 残りのデータを挿入
    if len(batch) > 0 {
        if err := insertBatch(conn.DB, "news", batch); err != nil {
            return fmt.Errorf("failed to insert batch to news: %w", err)
        }
    }
    
    log.Printf("Generated %d news articles", totalCount)
    return nil
}
```

### 3.3 バッチ挿入関数の設計

#### 3.3.1 insertBatch関数（汎用）

```go
// insertBatch はバッチでデータを挿入（汎用関数）
// 注意: GORMのCreateInBatchesを使用するか、生SQLでバッチ挿入を実装
func insertBatch(db *gorm.DB, tableName string, records interface{}) error {
    // トランザクション内で実行
    return db.Transaction(func(tx *gorm.DB) error {
        // GORMのCreateInBatchesを使用
        // 注意: 動的テーブル名の場合は生SQLを使用する必要がある
        return tx.Table(tableName).CreateInBatches(records, len(records)).Error
    })
}
```

#### 3.3.2 動的テーブル名対応のバッチ挿入

シャーディングテーブルの場合は、動的テーブル名を使用する必要があるため、生SQLでバッチ挿入を実装：

```go
// insertUsersBatch はusersテーブルにバッチでデータを挿入
func insertUsersBatch(db *gorm.DB, tableName string, users []*model.User) error {
    if len(users) == 0 {
        return nil
    }
    
    // バッチサイズを考慮して分割
    const batchSize = 500
    for i := 0; i < len(users); i += batchSize {
        end := i + batchSize
        if end > len(users) {
            end = len(users)
        }
        batch := users[i:end]
        
        // 生SQLでバッチ挿入
        query := fmt.Sprintf("INSERT INTO %s (name, email, created_at, updated_at) VALUES ", tableName)
        var values []interface{}
        var placeholders []string
        
        for _, user := range batch {
            placeholders = append(placeholders, "(?, ?, ?, ?)")
            values = append(values, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
        }
        
        query += strings.Join(placeholders, ", ")
        
        if err := db.Exec(query, values...).Error; err != nil {
            return fmt.Errorf("failed to insert batch: %w", err)
        }
    }
    
    return nil
}
```

## 4. データ生成ロジックの詳細設計

### 4.1 usersテーブルのデータ生成

#### 4.1.1 テーブル分散ロジック
- 32分割テーブル（users_000～users_031）に均等に分散
- 各テーブルに約3～4件ずつ生成（100件 ÷ 32テーブル）
- テーブル番号0～31をループして、各テーブルにデータを生成

#### 4.1.2 データ生成方法
- `name`: `gofakeit.Name()` - ランダムな名前
- `email`: `gofakeit.Email()` - ランダムなメールアドレス
- `created_at`, `updated_at`: `time.Now()` - 現在時刻
- `id`: データベースのAUTOINCREMENT機能を使用（明示的なID指定は不要）

#### 4.1.3 接続取得方法
```go
// テーブル番号から接続を取得
conn, err := groupManager.GetShardingConnection(tableNumber)
```

### 4.2 postsテーブルのデータ生成

#### 4.2.1 テーブル分散ロジック
- 32分割テーブル（posts_000～posts_031）に均等に分散
- 各テーブルに約3～4件ずつ生成（100件 ÷ 32テーブル）
- テーブル番号0～31をループして、各テーブルにデータを生成

#### 4.2.2 user_idの参照方法
- 既存のusersテーブルからuser_idを取得
- 全テーブルからuser_idを収集してリスト化
- ランダムにuser_idを選択してpostsに設定

#### 4.2.3 データ生成方法
- `user_id`: 既存のusersテーブルからランダムに選択
- `title`: `gofakeit.Sentence(5)` - 5単語程度の文
- `content`: `gofakeit.Paragraph(3, 5, 10, "\n")` - 3～5文、各文10単語程度
- `created_at`, `updated_at`: `time.Now()` - 現在時刻
- `id`: データベースのAUTOINCREMENT機能を使用（明示的なID指定は不要）

### 4.3 newsテーブルのデータ生成

#### 4.3.1 データ生成方法
- `title`: `gofakeit.Sentence(5)` - 5単語程度の文
- `content`: `gofakeit.Paragraph(3, 5, 10, "\n")` - 3～5文、各文10単語程度
- `author_id`: `gofakeit.Int64()` - ランダムな整数（NULL許容）
- `published_at`: `gofakeit.Date()` - ランダムな日時（NULL許容）
- `created_at`, `updated_at`: `time.Now()` - 現在時刻
- `id`: データベースのAUTOINCREMENT機能を使用（明示的なID指定は不要）

#### 4.3.2 接続取得方法
```go
// master接続を取得
conn, err := groupManager.GetMasterConnection()
```

## 5. バッチ処理の実装設計

### 5.1 バッチサイズ
- **バッチサイズ**: 500件ずつ（参考コードに基づく）
- 各テーブルへの挿入時にバッチ処理を使用

### 5.2 トランザクション管理
- 各バッチ挿入は独立したトランザクションで実行
- エラー発生時は適切にロールバック
- GORMの`Transaction()`メソッドを使用

### 5.3 パフォーマンス考慮
- バッチ挿入により、大量データ生成時のパフォーマンスを向上
- 100件程度のデータ生成は数秒以内に完了することを目標

## 6. エラーハンドリング設計

### 6.1 エラー処理の階層

```
┌─────────────────────────────────────────────────────────────┐
│              エラー発生箇所                                  │
│              - 設定読み込みエラー                            │
│              - DB接続エラー                                  │
│              - データ生成エラー                               │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              エラーメッセージの生成                          │
│              - エラー発生箇所を特定                           │
│              - 詳細なエラーメッセージを生成                   │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              ログ出力                                        │
│              - log.Fatalf()でエラーメッセージを表示          │
│              - 終了コード: 非ゼロ（1）                       │
└─────────────────────────────────────────────────────────────┘
```

### 6.2 エラーメッセージの形式

```go
// 設定読み込みエラー
log.Fatalf("Failed to load config: %v", err)

// DB接続エラー
log.Fatalf("Failed to create group manager: %v", err)
log.Fatalf("Failed to ping databases: %v", err)

// データ生成エラー
log.Fatalf("Failed to generate users: %v", err)
log.Fatalf("Failed to generate posts: %v", err)
log.Fatalf("Failed to generate news: %v", err)

// バッチ挿入エラー
return fmt.Errorf("failed to insert batch to %s: %w", tableName, err)
```

### 6.3 進捗表示

```go
// 各テーブルの生成件数を表示
log.Printf("Generated %d users in %s", countPerTable, tableName)
log.Printf("Generated %d posts in %s", countPerTable, tableName)
log.Printf("Generated %d news articles", totalCount)

// 生成完了メッセージ
log.Println("Sample data generation completed successfully")
```

## 7. テスト戦略

### 7.1 単体テスト
- 各データ生成関数の単体テスト
- バッチ挿入関数の単体テスト
- エラーハンドリングのテスト

### 7.2 統合テスト
- 実際のデータベースを使用した統合テスト
- データ生成後のデータ検証
- シャーディングテーブルへの分散確認

### 7.3 テストデータ
- テスト用のデータベースを使用
- テスト後にデータをクリーンアップ

## 8. 実装上の注意事項

### 8.1 シャーディングテーブルへの分散
- users, postsテーブルは32分割テーブルに適切に分散すること
- テーブル番号0～31をループして、各テーブルに均等にデータを生成すること
- 各テーブルに約3～4件ずつ生成すること（100件 ÷ 32テーブル）

### 8.2 postsテーブルのuser_id参照
- postsテーブルのuser_idは既存のusersテーブルから取得すること
- usersテーブルの生成後にpostsテーブルを生成すること
- 全テーブルからuser_idを収集してリスト化すること

### 8.3 バッチ挿入の実装
- バッチサイズは500件ずつ（参考コードに基づく）
- 動的テーブル名の場合は生SQLでバッチ挿入を実装すること
- トランザクション管理を適切に行うこと

### 8.4 Gofakeitの使用方法
- 各フィールドに適切なGofakeit関数を使用すること
- リアルなデータを生成するため、適切な関数を選択すること

### 8.5 エラーハンドリング
- データベース接続エラー、データ生成エラーなど、適切なエラーハンドリングを実装すること
- エラー発生時は詳細なエラーメッセージを表示すること
- 終了コードを適切に設定すること（エラー時は非ゼロ）

## 9. 参考情報

### 9.1 既存実装
- `server/cmd/list-users/main.go`: 既存のCLI実装パターン
- `server/internal/config/config.go`: 設定読み込み機能
- `server/internal/db/manager.go`: データベース接続管理機能
- `server/internal/model/user.go`: Userモデル
- `server/internal/model/post.go`: Postモデル
- `server/internal/model/news.go`: Newsモデル

### 9.2 技術スタック
- **Go**: 1.21+
- **GORM**: v1.25.12
- **データベース**: SQLite3（開発環境）
- **Gofakeit**: v6（新規追加）

### 9.3 参考コード
Issue #29に記載されている参考コードを基に実装する：
- バッチサイズ: 500件ずつ
- バッチ挿入の実装パターン
- Gofakeitを使用したデータ生成パターン
