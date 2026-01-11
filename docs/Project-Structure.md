# プロジェクト構造計画

## 概要

このドキュメントは、go-webdb-templateプロジェクトのひな形作成計画を記録したものです。

## プロジェクト構造

```
go-webdb-template/
├── server/                      # Golangサーバー
│   ├── cmd/
│   │   └── server/
│   │       └── main.go         # エントリーポイント
│   ├── internal/
│   │   ├── api/                # API定義層
│   │   │   ├── handler/        # HTTPハンドラー
│   │   │   │   ├── user_handler.go
│   │   │   │   └── user_handler_test.go
│   │   │   └── router/         # ルーティング
│   │   ├── usecase/            # ビジネスロジック層
│   │   │   ├── user_usecase.go
│   │   │   ├── user_usecase_test.go
│   │   │   ├── admin/          # Admin用usecase層
│   │   │   │   ├── dm_user_register_usecase.go
│   │   │   │   ├── dm_user_register_usecase_test.go
│   │   │   │   ├── api_key_usecase.go
│   │   │   │   └── api_key_usecase_test.go
│   │   │   └── cli/            # CLI用usecase層
│   │   │       ├── list_dm_users_usecase.go
│   │   │       ├── list_dm_users_usecase_test.go
│   │   │       ├── generate_secret_usecase.go
│   │   │       ├── generate_secret_usecase_test.go
│   │   │       ├── generate_sample_usecase.go
│   │   │       └── generate_sample_usecase_test.go
│   │   ├── service/            # ドメインロジック層
│   │   │   ├── user_service.go
│   │   │   ├── user_service_test.go
│   │   │   ├── secret_service.go
│   │   │   ├── secret_service_test.go
│   │   │   ├── api_key_service.go
│   │   │   ├── api_key_service_test.go
│   │   │   ├── generate_sample_service.go
│   │   │   └── generate_sample_service_test.go
│   │   ├── repository/         # データベース処理層
│   │   │   ├── user_repository.go
│   │   │   ├── user_repository_test.go
│   │   │   ├── dm_news_repository.go
│   │   │   └── dm_news_repository_test.go
│   │   ├── sql/                # SQL定義層
│   │   ├── model/              # データモデル
│   │   ├── db/                 # DB接続管理
│   │   │   ├── connection.go  # DB接続プール管理
│   │   │   ├── connection_test.go
│   │   │   ├── sharding.go    # Sharding戦略
│   │   │   ├── sharding_test.go
│   │   │   └── manager.go     # DBマネージャー
│   │   ├── auth/               # 認証・秘密鍵管理
│   │   │   ├── jwt.go          # JWT検証・生成
│   │   │   ├── secret.go       # 秘密鍵生成処理
│   │   │   └── secret_test.go  # 秘密鍵生成テスト
│   │   └── config/             # 設定読み込み
│   │       ├── config.go       # 設定構造体と読み込み処理
│   │       └── config_test.go
│   ├── test/                   # テストユーティリティ
│   │   ├── integration/        # 統合テスト
│   │   │   ├── api_test.go
│   │   │   └── sharding_test.go
│   │   ├── e2e/                # E2Eテスト
│   │   │   └── user_flow_test.go
│   │   ├── fixtures/           # テストデータ
│   │   │   ├── users.json
│   │   │   └── posts.json
│   │   └── testutil/           # テストヘルパー
│   │       ├── db.go           # テスト用DB準備
│   │       └── mock.go         # モックオブジェクト
│   ├── go.mod
│   └── go.sum
│
├── client/                      # Next.js + TypeScript
│   ├── src/
│   │   ├── app/                # App Router
│   │   │   ├── page.tsx        # トップページ
│   │   │   ├── users/          # ユーザー管理
│   │   │   ├── posts/          # 投稿管理
│   │   │   └── user-posts/     # ジョイン結果表示
│   │   ├── components/         # Reactコンポーネント
│   │   │   ├── UserList.tsx
│   │   │   └── UserList.test.tsx
│   │   ├── lib/                # API呼び出し等
│   │   │   ├── api.ts
│   │   │   └── api.test.ts
│   │   └── types/              # TypeScript型定義
│   ├── __tests__/              # Jestテスト
│   │   ├── integration/        # 統合テスト
│   │   │   └── api.test.ts
│   │   └── unit/               # ユニットテスト
│   │       └── utils.test.ts
│   ├── e2e/                    # E2Eテスト（Playwright）
│   │   ├── user-crud.spec.ts
│   │   └── post-crud.spec.ts
│   ├── .storybook/             # Storybook設定
│   ├── jest.config.js          # Jest設定
│   ├── playwright.config.ts    # Playwright設定
│   ├── package.json
│   ├── tsconfig.json
│   └── next.config.js
│
├── config/                      # 環境別設定ファイル
│   ├── develop.yaml            # 開発環境設定
│   ├── staging.yaml            # ステージング環境設定
│   └── production.yaml         # 本番環境設定
│
├── db/
│   └── migrations/             # マイグレーションSQL
│       ├── shard1/             # Shard 1用マイグレーション
│       │   └── 001_init.sql
│       └── shard2/             # Shard 2用マイグレーション
│           └── 001_init.sql
│
├── docs/
│   ├── plans/
│   │   └── project-structure.md  # このドキュメント
│   ├── Architecture.md         # アーキテクチャ説明
│   ├── API.md                  # API仕様
│   └── Sharding.md             # Sharding戦略ドキュメント
│
├── .gitignore
└── README.md
```

## レイヤー構成

### サーバー側（Go）

1. **API定義層** (`internal/api/`)
   - HTTPリクエスト/レスポンスの処理
   - ルーティング定義
   - バリデーション（形式チェック）
   - 認証・認可チェック

2. **ビジネスロジック層** (`internal/usecase/`)
   - アプリケーションのコアロジック
   - トランザクション管理
   - 複数のserviceを組み合わせた処理

3. **ドメインロジック層** (`internal/service/`)
   - ドメイン固有のロジック
   - ドメイン固有のバリデーション
   - ドメイン固有のビジネスルール

4. **データベース処理層** (`internal/repository/`)
   - データベースへのアクセス
   - CRUD操作の実装
   - Shard Key に基づくDB選択

5. **SQL定義層** (`internal/sql/`)
   - SQL クエリの定義
   - クエリビルダー

6. **DB接続管理層** (`internal/db/`)
   - 複数DBシャードへの接続プール管理
   - Sharding戦略の実装（Hash-based, Range-based等）
   - DB接続のライフサイクル管理

7. **設定管理層** (`internal/config/`)
   - 環境別設定ファイルの読み込み
   - 設定値のバリデーション
   - DBシャード設定の管理

### クライアント側（Next.js + TypeScript）

- **App Router**: ページルーティング
- **Components**: 再利用可能なUIコンポーネント
- **Lib**: API呼び出しやユーティリティ関数
- **Types**: TypeScript型定義

## データモデル

### 1. User（ユーザー）

| カラム名 | 型 | 説明 |
|---------|-----|------|
| id | INTEGER | 主キー |
| name | TEXT | ユーザー名 |
| email | TEXT | メールアドレス |
| created_at | DATETIME | 作成日時 |
| updated_at | DATETIME | 更新日時 |

### 2. Post（投稿）

| カラム名 | 型 | 説明 |
|---------|-----|------|
| id | INTEGER | 主キー |
| user_id | INTEGER | ユーザーID（外部キー） |
| title | TEXT | タイトル |
| content | TEXT | 本文 |
| created_at | DATETIME | 作成日時 |
| updated_at | DATETIME | 更新日時 |

### 3. ジョイン機能

UserとPostをJOINして、ユーザーとその投稿一覧を取得・表示する機能を提供します。

## 使用技術スタック

### サーバー側

- **言語**: Go 1.21+
- **データベース**: PostgreSQL or MySQL（全環境）
- **ルーティング**: gorilla/mux
- **DB接続**: GORM + gorm.io/driver/postgres
- **設定管理**: spf13/viper（YAML設定ファイル読み込み）
- **Sharding**: 自前実装（Hash-based sharding、8論理シャード、4物理DB）
- **テスト**:
  - testing（標準ライブラリ）
  - testify（アサーション、モック）
  - httptest（HTTPテスト）
  - go-sqlmock（DBモック）

### クライアント側

- **フレームワーク**: Next.js 14+ (App Router)
- **言語**: TypeScript 5+
- **UI開発**: Storybook
- **スタイリング**: TailwindCSS（オプション）
- **テスト**:
  - Jest（ユニットテスト、統合テスト）
  - React Testing Library（コンポーネントテスト）
  - Playwright（E2Eテスト）
  - MSW（APIモック）

## 機能一覧

### サーバーAPI

- `GET /api/users` - ユーザー一覧取得
- `GET /api/users/:id` - ユーザー詳細取得
- `POST /api/users` - ユーザー作成
- `PUT /api/users/:id` - ユーザー更新
- `DELETE /api/users/:id` - ユーザー削除
- `GET /api/posts` - 投稿一覧取得
- `GET /api/posts/:id` - 投稿詳細取得
- `POST /api/posts` - 投稿作成
- `PUT /api/posts/:id` - 投稿更新
- `DELETE /api/posts/:id` - 投稿削除
- `GET /api/user-posts` - ユーザーと投稿のジョイン結果取得

### クライアント画面

- `/` - トップページ（機能一覧）
- `/users` - ユーザー一覧・作成・編集・削除
- `/posts` - 投稿一覧・作成・編集・削除
- `/user-posts` - ユーザーと投稿のジョイン結果表示

## Sharding戦略

### 概要

複数のDBサーバーにデータを分散させることで、スケーラビリティを確保します。

### Shard Key

- **Userテーブル**: `user_id` をShard Keyとして使用
- **Postテーブル**: `user_id` をShard Keyとして使用（同一ユーザーのデータは同じShardに配置）

### Sharding方式

**Hash-based Sharding** を採用：
```
shard_id = hash(user_id) % shard_count
```

### Shard構成（例）

- **Shard 1**: user_id が偶数のユーザーとその投稿
- **Shard 2**: user_id が奇数のユーザーとその投稿

### クロスシャードクエリ

ユーザーと投稿のジョイン結果を取得する場合：
- 各Shardから並列にデータを取得
- アプリケーション層でマージして返却

### 実装の注意点

1. **Repository層での抽象化**: Shard選択ロジックはRepository層で隠蔽
2. **トランザクション**: 単一Shard内でのみトランザクションをサポート
3. **拡張性**: Shard数の変更に対応できる設計（将来的にConsistent Hashingも検討）

## 環境別設定管理

### 設定ファイル構成

環境変数 `APP_ENV` の値に基づいて、適切な設定ファイルを読み込みます。

```
APP_ENV=develop   → config/develop.yaml
APP_ENV=staging   → config/staging.yaml
APP_ENV=production → config/production.yaml
```

### 設定項目

各環境の設定ファイルには以下を含めます：

1. **サーバー設定**
   - ポート番号
   - タイムアウト値

2. **データベース設定**
   - Shard毎の接続情報（ホスト、ポート、DB名、認証情報）
   - 接続プール設定（最大接続数、アイドルタイムアウト等）

3. **ログ設定**
   - ログレベル
   - ログ出力先

4. **CORS設定**
   - 許可するオリジン

### 設定ファイル例（develop.yaml）

```yaml
server:
  port: 8080
  read_timeout: 30s
  write_timeout: 30s

database:
  shards:
    - id: 1
      host: localhost
      port: 5432
      name: app_db_shard1
      user: dev_user
      password: dev_password
      max_connections: 10
    - id: 2
      host: localhost
      port: 5433
      name: app_db_shard2
      user: dev_user
      password: dev_password
      max_connections: 10

logging:
  level: debug
  format: json

cors:
  allowed_origins:
    - http://localhost:3000
```

### セキュリティ考慮

- 本番環境の設定ファイルはGitにコミットしない（`.gitignore`に追加）
- パスワード等の機密情報は環境変数で上書き可能にする
- 設定ファイルのテンプレート（`*.yaml.example`）のみをリポジトリに含める

## テスト戦略

### テストレベル

1. **ユニットテスト**
   - 各関数・メソッドの単体テスト
   - モックを使用して依存関係を切り離す
   - カバレッジ目標: 80%以上

2. **統合テスト**
   - 複数のレイヤーを組み合わせたテスト
   - 実際のDBを使用（テスト用DB）
   - API → Usecase → Service → Repository の流れを確認

3. **E2Eテスト**
   - ユーザーシナリオベースのテスト
   - ブラウザ自動化でフロントエンドからバックエンドまで通しでテスト

### サーバー側テスト方針

#### ユニットテスト
- 各レイヤー（handler, usecase, service, repository）に `*_test.go` を配置
- テーブル駆動テストを活用
- モックは testify/mock や go-sqlmock を使用

```go
// 例: user_service_test.go
func TestUserService_GetUser(t *testing.T) {
    tests := []struct {
        name    string
        userID  int64
        want    *model.User
        wantErr bool
    }{
        // テストケース
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // テスト実装
        })
    }
}
```

#### 統合テスト
- `test/integration/` に配置
- テスト用PostgreSQLを使用（testutil.SetupTestGroupManager）
- テストデータは `test/fixtures/` から読み込み
- トランザクションでロールバック（クリーンアップ）

#### E2Eテスト
- `test/e2e/` に配置
- httptest.Server を使用してサーバーを起動
- 実際のHTTPリクエストでAPIをテスト

### クライアント側テスト方針

#### コンポーネントテスト
- 各コンポーネントに `.test.tsx` を配置
- React Testing Library でレンダリングとインタラクションをテスト
- MSWでAPIをモック

```typescript
// 例: UserList.test.tsx
describe('UserList', () => {
  it('displays user list', async () => {
    render(<UserList />)
    expect(await screen.findByText('John Doe')).toBeInTheDocument()
  })
})
```

#### 統合テスト
- `__tests__/integration/` に配置
- 複数コンポーネントの連携をテスト
- APIクライアントのテストも含む

#### E2Eテスト
- `e2e/` に配置
- Playwright でブラウザ自動化
- ユーザーフロー全体（作成→一覧→編集→削除）をテスト

### Sharding機能のテスト

#### Sharding戦略のテスト
- Hash関数の一貫性テスト
- Shard選択ロジックのテスト
- 複数Shardへのデータ分散確認

#### クロスシャードクエリのテスト
- 各Shardからのデータ取得
- マージ処理の正確性
- 並列実行のパフォーマンステスト

### CI/CDでのテスト実行

```yaml
# GitHub Actionsの例
- name: Run unit tests
  run: go test -v -cover ./...

- name: Run integration tests
  run: go test -v -tags=integration ./test/integration/...

- name: Run E2E tests
  run: go test -v -tags=e2e ./test/e2e/...
```

### カバレッジ計測

- サーバー側: `go test -coverprofile=coverage.out ./...`
- クライアント側: `npm run test:coverage`
- カバレッジレポートをCI/CDで可視化

## 開発方針

1. **大規模プロジェクト対応**: 小規模サンプルだが、スケーラブルな構成を採用
2. **レイヤー分離**: 責務を明確に分離し、保守性を向上
3. **型安全性**: TypeScriptによる型定義で安全性を確保
4. **テスト駆動**: ユニット/統合/E2Eテストで品質を担保、カバレッジ80%以上を目標
5. **ドキュメント**: Storybookとドキュメントで可視性を向上
6. **Docker環境**: PostgreSQLコンテナを使用し、環境構築を容易に
7. **Sharding対応**: 複数DBへの水平分割を想定した設計
8. **環境分離**: develop/staging/production環境で設定を切り替え可能に

## 次のステップ

1. ディレクトリ構造の作成
2. 環境別設定ファイルの作成（develop/staging/production.yaml）
3. サーバー側の基本ファイル作成（Go）
   - 設定管理層の実装とテスト
   - DB接続管理層の実装（Sharding対応）とテスト
   - Repository/Service/API層の実装とユニットテスト
4. データベースマイグレーションファイル作成（各Shard用）
5. クライアント側の基本ファイル作成（Next.js）
   - コンポーネントとテストの作成
   - APIクライアントとテストの作成
6. 統合テストとE2Eテストの実装
7. 各レイヤーの実装とSharding動作確認
8. テストカバレッジの確認と改善
9. Storybookセットアップ
10. CI/CD設定（GitHub Actions等）
11. ドキュメント作成（Architecture.md、API.md、Sharding.md、Testing.md）
