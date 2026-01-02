# Technology Stack and Architecture

## 開発環境サーバー構成

開発環境では3つのサーバーを起動する必要があります。

| サーバー | ポート | ディレクトリ | 起動コマンド |
|---------|-------|-------------|-------------|
| API サーバー | 8080 | `server/cmd/server` | `APP_ENV=develop go run ./cmd/server/main.go` |
| クライアント | 3000 | `client/` | `npm run dev` |
| Admin | 8081 | `server/cmd/admin` | `APP_ENV=develop go run ./cmd/admin/main.go` |

**注意**: 「サーバーを起動して」と言われた場合、上記3つ全てを起動すること。

## 技術スタック

### サーバー側

- **言語**: Go 1.21+
- **データベース**: 
  - 開発環境: SQLite3
  - 本番想定: PostgreSQL / MySQL
- **ルーティング**: `github.com/gorilla/mux`
- **DB接続**: `database/sql` + `github.com/mattn/go-sqlite3` / `github.com/lib/pq`
- **設定管理**: `github.com/spf13/viper` (YAML設定ファイル読み込み)
- **CORS**: `github.com/rs/cors`
- **Redis**: 
  - `github.com/redis/go-redis/v9` (Redisクライアント)
  - ジョブキュー用（単一接続）とレートリミット用（Cluster接続）の2種類の接続設定
- **ジョブキュー**: `github.com/hibiken/asynq` (Redisベースのジョブキュー)
- **メール送信**: 
  - `gopkg.in/mail.v2` (Mailpit用SMTP送信)
  - AWS SES SDK (本番環境用)
- **テスト**:
  - `testing` (標準ライブラリ)
  - `github.com/stretchr/testify` (アサーション、モック)
  - `net/http/httptest` (HTTPテスト)

### クライアント側

- **フレームワーク**: Next.js 14+ (App Router)
- **言語**: TypeScript 5+
- **UI開発**: Storybook（オプション）
- **テスト**:
  - Jest（ユニットテスト、統合テスト）
  - React Testing Library（コンポーネントテスト）
  - Playwright（E2Eテスト）
  - MSW（APIモック）

## アーキテクチャパターン

### レイヤードアーキテクチャ

```
┌─────────────────────────────────────┐
│         Client Layer                 │
│      (Next.js 14 + React)           │
└────────────────┬────────────────────┘
                 │ HTTP/REST
                 ▼
┌─────────────────────────────────────┐
│      Server Layer (Go)              │
│                                     │
│  ┌──────────────────────────────┐  │
│  │      API Layer                │  │
│  │  • HTTP Handlers              │  │
│  │  • Request validation         │  │
│  │  • Response formatting        │  │
│  └──────────┬───────────────────┘  │
│             │                        │
│  ┌──────────▼───────────────────┐  │
│  │    Service Layer              │  │
│  │  • Business logic            │  │
│  │  • Transaction management    │  │
│  │  • Cross-shard operations    │  │
│  └──────────┬───────────────────┘  │
│             │                        │
│  ┌──────────▼───────────────────┐  │
│  │  Repository Layer             │  │
│  │  • Data access abstraction    │  │
│  │  • CRUD operations           │  │
│  └──────────┬───────────────────┘  │
│             │                        │
│  ┌──────────▼───────────────────┐  │
│  │     DB Layer                  │  │
│  │  • Connection management     │  │
│  │  • Sharding strategy         │  │
│  └──────────┬───────────────────┘  │
└─────────────┼────────────────────────┘
              │
   ┌──────────┴──────────┐
   ▼                      ▼
┌─────────┐          ┌─────────┐
│ Shard 1 │          │ Shard 2 │
└─────────┘          └─────────┘
```

### レイヤー責務

#### 1. API Layer (`internal/api/`)
- HTTPリクエスト/レスポンスの処理
- ルーティング定義
- バリデーション
- エラーハンドリングとHTTPステータスコードマッピング

#### 2. Service Layer (`internal/service/`)
- ビジネスロジックの実装
- トランザクション管理
- クロスシャード操作
- データ変換
- メール送信処理（`email/`）
- ジョブキュー処理（`jobqueue/`）

#### 3. Repository Layer (`internal/repository/`)
- データアクセスの抽象化
- SQLクエリ構築
- CRUD操作
- ドメインモデルへの結果マッピング

#### 4. DB Layer (`internal/db/`)
- データベース接続管理
- シャーディング戦略の実装
- 接続プール管理
- シャードルーティング

#### 5. Config Layer (`internal/config/`)
- 環境別設定ファイルの読み込み
- 設定値のバリデーション
- DBシャード設定の管理

## シャーディング戦略

### Hash-Based Sharding

**アルゴリズム**:
```go
shard_id = hash(user_id) % shard_count + 1
```

**特徴**:
- FNV-1aハッシュ関数を使用
- シャードID範囲: 1 から N（Nはシャード数）
- 同じ`user_id`は常に同じシャードにマッピング
- 投稿はユーザーと同じシャードに配置（co-location）

**利点**:
- データの均等な分散
- シンプルで予測可能なシャード選択
- 決定論的（同じキーは常に同じシャード）

**制約**:
- シャード追加/削除時のリバランスが困難
- クロスシャードの範囲クエリは高コスト
- 関連データは同じシャードキーを共有する必要がある

## データフロー

### ユーザー作成の例

```
1. Client → API Layer
   POST /api/users
   Body: {"name": "John", "email": "john@example.com"}

2. API Layer → Service Layer
   UserHandler.CreateUser()
   ↓
   UserService.CreateUser(CreateUserRequest)

3. Service Layer → Repository Layer
   ビジネスルールの検証
   ↓
   UserRepository.Create(user)

4. Repository Layer → DB Layer
   SQLクエリ構築
   ↓
   DBManager.GetConnectionByKey(userID)

5. DB Layer
   hash(userID)でシャードIDを計算
   ↓
   適切なシャードへの接続を返却

6. Repository Layer
   INSERT文を実行
   ↓
   作成されたユーザーを返却

7. Service Layer → API Layer
   Userを返却
   ↓
   UserHandlerがレスポンスをフォーマット

8. API Layer → Client
   HTTP 201 Created
   Body: {"id": 1, "name": "John", ...}
```

## 設定管理

### 環境別設定

環境変数 `APP_ENV` の値に基づいて設定ファイルを読み込み:

- `APP_ENV=develop` → `config/develop.yaml`
- `APP_ENV=staging` → `config/staging.yaml`
- `APP_ENV=production` → `config/production.yaml`

### 設定項目

- サーバー設定（ポート、タイムアウト）
- データベース設定（各シャードの接続情報、接続プール設定）
- ログ設定（ログレベル、出力先、SQLログ、メールログ）
- CORS設定（許可するオリジン）
- Redis設定（ジョブキュー用、レートリミット用）
  - 接続オプション（MaxRetries、RetryBackoff、DialTimeout、ReadTimeout、PoolSize、PoolTimeout）
- メール送信設定（送信方式、SMTP設定、AWS SES設定）

## エラーハンドリング

### エラー伝播

1. **Repository Layer**: Go errorsを返却
2. **Service Layer**: コンテキスト付きでエラーをラップ
3. **API Layer**: エラーをHTTPステータスコードに変換
4. **Client**: ユーザーフレンドリーなエラーメッセージを表示

### HTTPエラーレスポンス

```go
// API Layer
w.WriteHeader(http.StatusBadRequest)
json.NewEncoder(w).Encode(map[string]string{
    "error": "Invalid request",
})
```

## セキュリティ考慮事項

1. **CORS**: ルーターで特定のオリジンを許可
2. **入力検証**: API層とService層の両方で実施
3. **SQLインジェクション防止**: パラメータ化クエリを使用
4. **環境変数**: 機密データは設定ファイルに保存（gitから除外）

## パフォーマンス最適化

### 接続プール

各シャードが独自の接続プールを維持:
- `SetMaxOpenConns(25)`
- `SetMaxIdleConns(5)`
- `SetConnMaxLifetime(5 * time.Minute)`

### 並列クエリ

クロスシャードクエリはgoroutineを使用して並列実行

### インデックス

各シャードに適切なインデックスを設定:
- `idx_users_email` on `users(email)`
- `idx_posts_user_id` on `posts(user_id)`
- `idx_posts_created_at` on `posts(created_at DESC)`

## テスト戦略

### テストピラミッド

```
        ╱╲
       ╱  ╲
      ╱ E2E ╲     ← 少数、低速、高価値
     ╱────────╲
    ╱          ╲
   ╱ Integration╲  ← 中程度、中速
  ╱──────────────╲
 ╱                ╲
╱   Unit Tests     ╲ ← 多数、高速、集中
╱────────────────────╲
```

### テストレベル

1. **ユニットテスト**: 各関数・メソッドの単体テスト（カバレッジ目標: 80%以上）
2. **統合テスト**: 複数レイヤーの組み合わせテスト
3. **E2Eテスト**: ユーザーシナリオベースのテスト

## 依存関係管理

### サーバー依存関係

- `github.com/spf13/viper`: 設定管理
- `github.com/gorilla/mux`: HTTPルーティング
- `github.com/mattn/go-sqlite3`: SQLiteドライバ（開発）
- `github.com/lib/pq`: PostgreSQLドライバ（本番）
- `github.com/rs/cors`: CORSミドルウェア
- `github.com/redis/go-redis/v9`: Redisクライアント
- `github.com/hibiken/asynq`: Redisベースのジョブキュー
- `gopkg.in/mail.v2`: SMTPメール送信（Mailpit用）
- AWS SES SDK: AWS SESメール送信（本番環境用）

### クライアント依存関係

- `next`: Reactフレームワーク
- `react`, `react-dom`: UIライブラリ
- `typescript`: 型安全性

