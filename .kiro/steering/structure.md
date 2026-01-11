# File Organization and Code Patterns

## プロジェクト構造

```
go-webdb-template/
├── server/                      # Golangサーバー
│   ├── cmd/
│   │   └── server/
│   │       └── main.go         # エントリーポイント
│   ├── internal/                # 内部パッケージ
│   │   ├── api/                # API定義層
│   │   │   ├── handler/        # HTTPハンドラー
│   │   │   │   ├── user_handler.go
│   │   │   │   └── post_handler.go
│   │   │   └── router/         # ルーティング
│   │   │       └── router.go
│   │   ├── usecase/            # ビジネスロジック層
│   │   │   ├── dm_user_usecase.go
│   │   │   ├── dm_post_usecase.go
│   │   │   ├── email_usecase.go
│   │   │   ├── cli/            # CLI用usecase層
│   │   │   │   ├── list_dm_users_usecase.go
│   │   │   │   ├── list_dm_users_usecase_test.go
│   │   │   │   ├── generate_secret_usecase.go
│   │   │   │   └── generate_secret_usecase_test.go
│   │   │   └── ...
│   │   ├── service/            # ドメインロジック層
│   │   │   ├── dm_user_service.go
│   │   │   ├── dm_post_service.go
│   │   │   ├── secret_service.go
│   │   │   ├── secret_service_test.go
│   │   │   ├── email/          # メール送信サービス
│   │   │   │   ├── email_sender.go
│   │   │   │   ├── email_service.go
│   │   │   │   ├── mock_sender.go
│   │   │   │   ├── mailpit_sender.go
│   │   │   │   ├── ses_sender.go
│   │   │   │   └── template.go
│   │   │   └── jobqueue/      # ジョブキューサービス
│   │   │       ├── client.go
│   │   │       ├── server.go
│   │   │       ├── processor.go
│   │   │       ├── constants.go
│   │   │       └── redis_options.go
│   │   ├── repository/         # データベース処理層
│   │   │   ├── user_repository.go
│   │   │   ├── user_repository_test.go
│   │   │   ├── post_repository.go
│   │   │   └── post_repository_test.go
│   │   ├── model/               # データモデル
│   │   │   ├── user.go
│   │   │   └── post.go
│   │   ├── db/                 # DB接続管理
│   │   │   ├── connection.go  # DB接続プール管理
│   │   │   ├── manager.go      # DBマネージャー
│   │   │   ├── sharding.go     # Sharding戦略
│   │   │   └── sharding_test.go
│   │   ├── auth/               # 認証・秘密鍵管理
│   │   │   ├── jwt.go          # JWT検証・生成
│   │   │   ├── secret.go       # 秘密鍵生成処理
│   │   │   └── secret_test.go  # 秘密鍵生成テスト
│   │   └── config/             # 設定読み込み
│   │       └── config.go       # 設定構造体と読み込み処理
│   ├── test/                   # テストユーティリティ
│   │   ├── integration/        # 統合テスト
│   │   │   ├── user_flow_test.go
│   │   │   └── post_flow_test.go
│   │   ├── e2e/                # E2Eテスト
│   │   │   └── api_test.go
│   │   ├── fixtures/           # テストデータ
│   │   │   ├── users.go
│   │   │   └── posts.go
│   │   └── testutil/           # テストヘルパー
│   │       └── db.go           # テスト用DB準備
│   ├── go.mod
│   └── go.sum
│
├── client/                      # Next.js + TypeScript
│   ├── src/
│   │   ├── app/                # App Router
│   │   │   ├── page.tsx        # トップページ
│   │   │   ├── users/          # ユーザー管理
│   │   │   │   └── page.tsx
│   │   │   ├── posts/          # 投稿管理
│   │   │   │   └── page.tsx
│   │   │   └── user-posts/     # ジョイン結果表示
│   │   │       └── page.tsx
│   │   ├── lib/                # API呼び出し等
│   │   │   ├── api.ts
│   │   │   └── __tests__/
│   │   │       └── api.test.ts
│   │   └── types/              # TypeScript型定義
│   │       ├── user.ts
│   │       └── post.ts
│   ├── __tests__/              # Jestテスト
│   │   └── integration/
│   │       └── users-page.test.tsx
│   ├── e2e/                    # E2Eテスト（Playwright）
│   │   ├── user-flow.spec.ts
│   │   ├── post-flow.spec.ts
│   │   └── cross-shard.spec.ts
│   ├── jest.config.js          # Jest設定
│   ├── playwright.config.ts     # Playwright設定
│   ├── package.json
│   ├── tsconfig.json
│   └── next.config.js
│
├── config/                      # 環境別設定ファイル
│   ├── develop/                # 開発環境設定ディレクトリ
│   │   ├── config.yaml         # メイン設定（server, admin, logging, cors）
│   │   ├── database.yaml       # データベース設定
│   │   └── cacheserver.yaml    # Redis設定（ジョブキュー用、レートリミット用）
│   ├── production/             # 本番環境設定ディレクトリ
│   │   ├── config.yaml.example # メイン設定テンプレート
│   │   ├── database.yaml.example # データベース設定テンプレート
│   │   └── cacheserver.yaml.example # Redis設定テンプレート
│   └── staging/                # ステージング環境設定ディレクトリ
│       ├── config.yaml         # メイン設定
│       ├── database.yaml       # データベース設定
│       └── cacheserver.yaml    # Redis設定
│
├── db/
│   └── migrations/             # マイグレーションSQL
│       ├── shard1/             # Shard 1用マイグレーション
│       │   └── 001_init.sql
│       └── shard2/             # Shard 2用マイグレーション
│           └── 001_init.sql
│
├── docs/
│   ├── Architecture.md         # アーキテクチャ説明
│   ├── API.md                  # API仕様
│   ├── Sharding.md             # Sharding戦略ドキュメント
│   ├── Testing.md              # テストドキュメント
│   └── Project-Structure.md    # プロジェクト構造計画
│
├── docker-compose.redis.yml        # Redis（ジョブキュー用）Docker Compose設定
├── docker-compose.redis-cluster.yml # Redis Cluster（レートリミット用）Docker Compose設定
├── docker-compose.redis-insight.yml # Redis Insight（データビューワ）Docker Compose設定
├── scripts/
│   ├── start-redis.sh              # Redis起動スクリプト
│   ├── start-redis-cluster.sh       # Redis Cluster起動スクリプト
│   └── start-redis-insight.sh       # Redis Insight起動スクリプト
├── redis/                           # Redisデータ永続化ディレクトリ
│   └── data/
│       └── jobqueue/               # ジョブキュー用Redisデータ
├── .gitignore
└── README.md
```

## 命名規則

### Go (サーバー側)

#### パッケージ名
- 小文字、単語区切りなし: `handler`, `service`, `repository`
- 内部パッケージは `internal/` 配下に配置

#### ファイル名
- スネークケース: `user_handler.go`, `post_repository.go`
- テストファイル: `*_test.go`

#### 型・関数名
- パブリック: パスカルケース: `UserHandler`, `CreateUser`
- プライベート: キャメルケース: `validateUser`, `getConnection`

#### 構造体
```go
// モデル
type User struct {
    ID        int64     `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// リクエスト/レスポンス
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

### TypeScript (クライアント側)

#### ファイル名
- ケバブケース: `user-card.tsx`, `api-client.ts`
- コンポーネント: パスカルケース: `UserCard.tsx`

#### 型・関数名
- パスカルケース: `User`, `Post`, `createUser`
- 変数: キャメルケース: `userList`, `apiClient`

#### コンポーネント
```typescript
// コンポーネント
export function UserCard({ user }: { user: User }) {
  return <div>{user.name}</div>
}

// API クライアント
export const apiClient = {
  createUser: async (data: CreateUserRequest): Promise<User> => {
    // ...
  }
}
```

## コードパターン

### Go パターン

#### Handler パターン
```go
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    
    user, err := h.service.CreateUser(&req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}
```

#### Repository パターン
```go
func (r *UserRepository) Create(db *sql.DB, user *model.User) error {
    query := `INSERT INTO users (name, email) VALUES (?, ?)`
    result, err := db.Exec(query, user.Name, user.Email)
    if err != nil {
        return err
    }
    
    id, err := result.LastInsertId()
    if err != nil {
        return err
    }
    
    user.ID = id
    return nil
}
```

#### Service パターン
```go
func (s *UserService) CreateUser(req *CreateUserRequest) (*model.User, error) {
    // バリデーション
    if req.Name == "" {
        return nil, errors.New("name is required")
    }
    
    // モデル作成
    user := &model.User{
        Name:  req.Name,
        Email: req.Email,
    }
    
    // シャード選択と保存
    conn, err := s.dbManager.GetConnectionByKey(user.ID)
    if err != nil {
        return nil, err
    }
    
    if err := s.repo.Create(conn.DB(), user); err != nil {
        return nil, err
    }
    
    return user, nil
}
```

### TypeScript パターン

#### API クライアント
```typescript
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

export const apiClient = {
  async createUser(data: CreateUserRequest): Promise<User> {
    const response = await fetch(`${API_BASE_URL}/api/users`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    })
    
    if (!response.ok) {
      throw new Error('Failed to create user')
    }
    
    return response.json()
  }
}
```

#### Next.js App Router パターン
```typescript
// app/users/page.tsx
export default async function UsersPage() {
  const users = await apiClient.getUsers()
  
  return (
    <div>
      <h1>ユーザー一覧</h1>
      {users.map(user => (
        <UserCard key={user.id} user={user} />
      ))}
    </div>
  )
}
```

## テストパターン

### Go テスト

#### テーブル駆動テスト
```go
func TestUserRepository_GetByID(t *testing.T) {
    tests := []struct {
        name    string
        userID  int64
        want    *model.User
        wantErr bool
    }{
        {
            name:   "existing user",
            userID: 1,
            want:   &model.User{ID: 1, Name: "Test"},
        },
        {
            name:    "non-existing user",
            userID:  999,
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // テスト実装
        })
    }
}
```

### TypeScript テスト

#### React Testing Library
```typescript
describe('UserCard', () => {
  it('renders user information', () => {
    const user = { id: 1, name: 'Test User', email: 'test@example.com' }
    render(<UserCard user={user} />)
    
    expect(screen.getByText('Test User')).toBeInTheDocument()
    expect(screen.getByText('test@example.com')).toBeInTheDocument()
  })
})
```

#### Playwright E2E
```typescript
test('should create user', async ({ page }) => {
  await page.goto('http://localhost:3000/users')
  
  await page.fill('input[name="name"]', 'E2E Test User')
  await page.fill('input[name="email"]', 'e2e@example.com')
  await page.click('button[type="submit"]')
  
  await expect(page.locator('text=E2E Test User')).toBeVisible()
})
```

## ディレクトリ配置ルール

### サーバー側

- **`cmd/`**: エントリーポイント（main.go）
- **`internal/`**: 内部パッケージ（外部からインポート不可）
  - **`api/`**: HTTPハンドラーとルーター
  - **`usecase/`**: ビジネスロジック
  - **`service/`**: ドメインロジック
  - **`repository/`**: データアクセス
  - **`model/`**: ドメインモデル
  - **`db/`**: データベース接続管理
  - **`config/`**: 設定管理
- **`test/`**: テストユーティリティと統合テスト

### クライアント側

- **`src/app/`**: Next.js App Router ページ
- **`src/lib/`**: ユーティリティとAPIクライアント
- **`src/types/`**: TypeScript型定義
- **`__tests__/`**: Jestテスト
- **`e2e/`**: Playwright E2Eテスト

## インポート規則

### Go

```go
import (
    // 標準ライブラリ
    "context"
    "database/sql"
    "encoding/json"
    "net/http"
    
    // サードパーティ
    "github.com/gorilla/mux"
    "github.com/stretchr/testify/assert"
    
    // 内部パッケージ
    "your-project/internal/model"
    "your-project/internal/service"
)
```

### TypeScript

```typescript
// 外部ライブラリ
import { useState } from 'react'
import { render, screen } from '@testing-library/react'

// 内部モジュール（パスエイリアス）
import { apiClient } from '@/lib/api'
import { User } from '@/types/user'
```

## コメント規則

### Go

```go
// User represents a user in the system.
type User struct {
    ID int64 `json:"id"`
}

// CreateUser creates a new user in the appropriate shard.
// It returns the created user with assigned ID.
func (s *UserService) CreateUser(req *CreateUserRequest) (*model.User, error) {
    // ...
}
```

### TypeScript

```typescript
/**
 * Creates a new user via the API.
 * @param data - User creation data
 * @returns The created user
 */
export async function createUser(data: CreateUserRequest): Promise<User> {
  // ...
}
```

## エラーハンドリングパターン

### Go

```go
// Repository層: エラーをそのまま返す
if err != nil {
    return nil, err
}

// Service層: エラーにコンテキストを追加
if err != nil {
    return nil, fmt.Errorf("failed to create user: %w", err)
}

// Handler層: HTTPステータスコードに変換
if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
}
```

### TypeScript

```typescript
try {
  const user = await apiClient.createUser(data)
  // 成功処理
} catch (error) {
  if (error instanceof Error) {
    // エラー処理
    console.error('Failed to create user:', error.message)
  }
}
```

