# 不足しているテストコードの作成の設計書

## Overview

### 目的
テストのない箇所にテストコードを用意する。特に、clientアプリのテストが不足している箇所に対して、画面が表示されているか確認する程度の簡易なテストでも良いので、テストコードを作成する。簡単にテストが作れるようならしっかりテストを作成する。

### ユーザー
- **開発者**: テストコードの追加により、コードの品質保証とリグレッション防止を実現
- **エンドユーザー**: 既存の機能が正常に動作することを維持

### 影響
現在のシステム状態を以下のように変更する：
- クライアント側: 新しいテストファイルの追加（`client/src/__tests__/`配下）
- サーバー側: 新しいテストファイルの追加（`server/internal/`配下の各パッケージ）

### Goals
- テストが不足している箇所を特定する
- クライアント側のページコンポーネント、コンポーネント、カスタムフック、ユーティリティ関数にテストを作成する
- サーバー側のハンドラー、サービス、リポジトリ、ユースケースにテストを作成する
- 既存のテストパターンに従ったテストコードを作成する

### Non-Goals
- 外部ライブラリのコードのテスト作成
- 既にテストが存在する箇所の再実装
- テストカバレッジの100%達成（簡易なテストでも良い）

## Architecture

### 設計方針

#### 1. テストが不足している箇所の特定方法
1. **クライアント側**:
   - `client/app/`配下の各ページファイルに対して、対応するテストファイルが存在するか確認
   - `client/components/`配下の各コンポーネントファイルに対して、対応するテストファイルが存在するか確認
   - `client/lib/hooks/`配下の各カスタムフックファイルに対して、対応するテストファイルが存在するか確認
   - `client/lib/utils.ts`に対して、対応するテストファイルが存在するか確認

2. **サーバー側**:
   - `server/internal/api/handler/`配下の各ハンドラーファイルに対して、対応するテストファイルが存在するか確認
   - `server/internal/service/`配下の各サービスファイルに対して、対応するテストファイルが存在するか確認
   - `server/internal/repository/`配下の各リポジトリファイルに対して、対応するテストファイルが存在するか確認
   - `server/internal/usecase/`配下の各ユースケースファイルに対して、対応するテストファイルが存在するか確認

#### 2. テストの作成方針
- **簡易なテスト**: 画面が表示されているか確認する程度のテストでも良い
- **しっかりしたテスト**: 簡単にテストが作れるようなら、主要な機能（フォーム送信、データ表示、クリック、入力など）をテストする

#### 3. テストの種類と優先順位
1. **高優先度**: 
   - クライアント側: ページコンポーネント（ユーザーが直接アクセスするページ）
   - サーバー側: ハンドラー、サービス（主要なビジネスロジック）
2. **中優先度**: 
   - クライアント側: 主要なコンポーネント（フォーム、カード、モーダルなど）
   - サーバー側: リポジトリ、ユースケース
3. **低優先度**: 
   - クライアント側: ユーティリティコンポーネント（アイコン、ローディングスピナーなど）
   - サーバー側: ユーティリティ関数、ヘルパー関数

## 詳細設計

### 1. クライアント側のテスト設計

#### 1.1 ページコンポーネントのテスト

##### テストファイルの配置
- `client/src/__tests__/integration/`配下に配置
- ファイル名: `{page-name}-page.test.tsx`（例: `dm-posts-page.test.tsx`）

##### テストの構造
```typescript
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import PageComponent from '@/app/{path}/page'

// MSWサーバーの設定
const server = setupServer(
  // APIモックの設定
)

beforeAll(() => server.listen())
afterEach(() => server.resetHandlers())
afterAll(() => server.close())

describe('{PageName} Integration', () => {
  it('displays page content', async () => {
    render(<PageComponent />)
    // 画面が表示されているか確認
    await waitFor(() => {
      expect(screen.getByText(/expected text/i)).toBeInTheDocument()
    })
  })

  // 必要に応じて、主要な機能のテストを追加
})
```

##### 簡易なテストの例
```typescript
it('displays page content', async () => {
  render(<PageComponent />)
  
  // ページの主要な要素が表示されているか確認
  await waitFor(() => {
    expect(screen.getByRole('heading')).toBeInTheDocument()
  })
})
```

##### しっかりしたテストの例
```typescript
it('creates a new item', async () => {
  const user = userEvent.setup()
  render(<PageComponent />)

  // フォーム入力
  const nameInput = screen.getByLabelText('名前')
  const emailInput = screen.getByLabelText('メールアドレス')
  const submitButton = screen.getByRole('button', { name: /作成/ })

  await user.type(nameInput, 'New Item')
  await user.type(emailInput, 'new@example.com')
  await user.click(submitButton)

  // 作成されたアイテムが表示されることを確認
  await waitFor(() => {
    expect(screen.getByText('New Item')).toBeInTheDocument()
  })
})
```

#### 1.2 コンポーネントのテスト

##### テストファイルの配置
- `client/src/__tests__/components/`配下に配置
- ファイル名: `{component-name}.test.tsx`（例: `feed-form.test.tsx`）

##### テストの構造
```typescript
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import Component from '@/components/{path}/{component-name}'

describe('ComponentName', () => {
  it('renders component', () => {
    render(<Component {...props} />)
    expect(screen.getByText(/expected text/i)).toBeInTheDocument()
  })

  // 必要に応じて、主要な機能のテストを追加
})
```

##### 簡易なテストの例
```typescript
it('renders component', () => {
  render(<FeedForm />)
  expect(screen.getByRole('form')).toBeInTheDocument()
})
```

##### しっかりしたテストの例
```typescript
it('submits form with valid data', async () => {
  const user = userEvent.setup()
  const onSubmit = jest.fn()
  
  render(<FeedForm onSubmit={onSubmit} />)

  const titleInput = screen.getByLabelText('タイトル')
  const contentInput = screen.getByLabelText('内容')
  const submitButton = screen.getByRole('button', { name: /投稿/ })

  await user.type(titleInput, 'Test Title')
  await user.type(contentInput, 'Test Content')
  await user.click(submitButton)

  expect(onSubmit).toHaveBeenCalledWith({
    title: 'Test Title',
    content: 'Test Content',
  })
})
```

#### 1.3 カスタムフックのテスト

##### テストファイルの配置
- `client/src/__tests__/lib/hooks/`配下に配置
- ファイル名: `{hook-name}.test.ts`（例: `use-intersection-observer.test.ts`）

##### テストの構造
```typescript
import { renderHook, act } from '@testing-library/react'
import useCustomHook from '@/lib/hooks/{hook-name}'

describe('useCustomHook', () => {
  it('returns expected value', () => {
    const { result } = renderHook(() => useCustomHook())
    expect(result.current).toBeDefined()
  })

  // 必要に応じて、主要な機能のテストを追加
})
```

#### 1.4 ユーティリティ関数のテスト

##### テストファイルの配置
- `client/src/__tests__/lib/`配下に配置
- ファイル名: `utils.test.ts`

##### テストの構造
```typescript
import { utilityFunction } from '@/lib/utils'

describe('utilityFunction', () => {
  it('returns expected result', () => {
    const result = utilityFunction(input)
    expect(result).toBe(expected)
  })
})
```

### 2. サーバー側のテスト設計

#### 2.1 ハンドラーのテスト

##### テストファイルの配置
- ハンドラーファイルと同じディレクトリに配置
- ファイル名: `{handler-name}_test.go`（例: `dm_user_handler_test.go`）

##### テストの構造
```go
package handler

import (
    "context"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/danielgtaylor/huma/v2"
    "github.com/danielgtaylor/huma/v2/adapters/humaecho"
    "github.com/labstack/echo/v4"
)

func TestDmUserHandler_CreateUser(t *testing.T) {
    tests := []struct {
        name        string
        setup       func() (*DmUserHandler, huma.API)
        requestBody string
        wantStatus   int
        wantErr     bool
    }{
        {
            name: "creates user successfully",
            setup: func() (*DmUserHandler, huma.API) {
                // モックの設定
            },
            requestBody: `{"name":"Test User","email":"test@example.com"}`,
            wantStatus: http.StatusCreated,
            wantErr:   false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // テスト実装
        })
    }
}
```

#### 2.2 サービスのテスト

##### テストファイルの配置
- サービスファイルと同じディレクトリに配置
- ファイル名: `{service-name}_test.go`（例: `dm_user_service_test.go`）

##### テストの構造
```go
package service

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestDmUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name    string
        req     *model.CreateDmUserRequest
        want    *model.DmUser
        wantErr bool
    }{
        {
            name: "creates user successfully",
            req: &model.CreateDmUserRequest{
                Name:  "Test User",
                Email: "test@example.com",
            },
            want: &model.DmUser{
                Name:  "Test User",
                Email: "test@example.com",
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // モックリポジトリの設定
            // サービスの実行
            // アサーション
        })
    }
}
```

#### 2.3 リポジトリのテスト

##### テストファイルの配置
- リポジトリファイルと同じディレクトリに配置
- ファイル名: `{repository-name}_test.go`（例: `dm_user_repository_test.go`）

##### テストの構造
```go
package repository

import (
    "context"
    "database/sql"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/taku-o/go-webdb-template/internal/model"
    "github.com/taku-o/go-webdb-template/test/testutil"
)

func TestDmUserRepository_Create(t *testing.T) {
    // テスト用DBの準備
    db := testutil.SetupTestDB(t)
    defer testutil.TeardownTestDB(t, db)

    repo := NewDmUserRepository()

    tests := []struct {
        name    string
        user    *model.DmUser
        wantErr bool
    }{
        {
            name: "creates user successfully",
            user: &model.DmUser{
                Name:  "Test User",
                Email: "test@example.com",
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := repo.Create(context.Background(), db, tt.user)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.NotEmpty(t, tt.user.ID)
            }
        })
    }
}
```

#### 2.4 ユースケースのテスト

##### テストファイルの配置
- ユースケースファイルと同じディレクトリに配置
- ファイル名: `{usecase-name}_test.go`（例: `dm_user_usecase_test.go`）

##### テストの構造
```go
package api

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/taku-o/go-webdb-template/internal/model"
)

// MockDmUserService はDmUserServiceのモック
type MockDmUserService struct {
    CreateDmUserFunc func(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error)
}

func (m *MockDmUserService) CreateDmUser(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error) {
    if m.CreateDmUserFunc != nil {
        return m.CreateDmUserFunc(ctx, req)
    }
    return nil, nil
}

func TestDmUserUsecase_CreateDmUser(t *testing.T) {
    tests := []struct {
        name        string
        mockFunc    func(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error)
        req         *model.CreateDmUserRequest
        want        *model.DmUser
        wantErr     bool
        expectedErr string
    }{
        {
            name: "creates user successfully",
            mockFunc: func(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error) {
                return &model.DmUser{
                    ID:    "user-001",
                    Name:  req.Name,
                    Email: req.Email,
                }, nil
            },
            req: &model.CreateDmUserRequest{
                Name:  "Test User",
                Email: "test@example.com",
            },
            want: &model.DmUser{
                ID:    "user-001",
                Name:  "Test User",
                Email: "test@example.com",
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockService := &MockDmUserService{
                CreateDmUserFunc: tt.mockFunc,
            }
            usecase := NewDmUserUsecase(mockService)

            got, err := usecase.CreateDmUser(context.Background(), tt.req)

            if tt.wantErr {
                assert.Error(t, err)
                if tt.expectedErr != "" {
                    assert.Contains(t, err.Error(), tt.expectedErr)
                }
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.want, got)
            }
        })
    }
}
```

## 実装上の注意事項

### 1. クライアント側のテスト

#### 1.1 useEffectの使用制限
- **テストコード内ではuseEffectを極力使わない**
- テストコードでも、本番コードと同様にuseEffectの使用を避ける
- 必要に応じて、`waitFor`や`findBy*`クエリを使用して非同期処理を待つ
- 絶対にuseEffectが必要な場合のみ使用（例: クリーンアップ処理）

#### 1.2 MSWの使用
- APIモックにはMSW（Mock Service Worker）を使用
- 既存のテストパターンに従う

#### 1.3 テストの実行
- `npm test`でテストを実行
- 既存のテストが正常に動作することを確認

#### 1.4 テストの構造
- `describe`ブロックでテストをグループ化
- テストの名前を明確にする

### 2. サーバー側のテスト

#### 2.1 テーブル駆動テストパターン
- Goの標準的なテーブル駆動テストパターンを使用
- 既存のテストパターンに従う

#### 2.2 テストの実行
- `APP_ENV=test go test ./...`でテストを実行
- 認証エラーが発生しないように`APP_ENV=test`を必ず指定

#### 2.3 モックの使用
- 必要に応じて、モックを使用
- 既存のテストパターンに従う

#### 2.4 データベーステスト
- テスト用データベースを使用
- `test/testutil`のヘルパー関数を使用

## テスト戦略

### 1. テストが不足している箇所の特定

#### 1.1 クライアント側
1. `client/app/`配下の各ページファイルを確認
2. `client/components/`配下の各コンポーネントファイルを確認
3. `client/lib/hooks/`配下の各カスタムフックファイルを確認
4. `client/lib/utils.ts`を確認
5. 対応するテストファイルが存在するか確認

#### 1.2 サーバー側
1. `server/internal/api/handler/`配下の各ハンドラーファイルを確認
2. `server/internal/service/`配下の各サービスファイルを確認
3. `server/internal/repository/`配下の各リポジトリファイルを確認
4. `server/internal/usecase/`配下の各ユースケースファイルを確認
5. 対応するテストファイルが存在するか確認

### 2. テストの作成

#### 2.1 簡易なテスト
- 画面が表示されているか確認する程度のテスト
- コンポーネントが正常にレンダリングされることを確認

#### 2.2 しっかりしたテスト
- 主要な機能（フォーム送信、データ表示、クリック、入力など）をテスト
- エラーハンドリングをテスト
- エッジケースをテスト

### 3. テストの実行と確認

#### 3.1 クライアント側
- `npm test`でテストを実行
- 既存のテストが正常に動作することを確認

#### 3.2 サーバー側
- `APP_ENV=test go test ./...`でテストを実行
- 既存のテストが正常に動作することを確認

## 移行計画

### Phase 1: テストが不足している箇所の特定
1. クライアント側のページコンポーネントを確認
2. クライアント側のコンポーネントを確認
3. クライアント側のカスタムフックとユーティリティ関数を確認
4. サーバー側のハンドラー、サービス、リポジトリ、ユースケースを確認

### Phase 2: 高優先度のテスト作成
1. クライアント側: ページコンポーネントのテスト作成
2. サーバー側: ハンドラー、サービスのテスト作成

### Phase 3: 中優先度のテスト作成
1. クライアント側: 主要なコンポーネントのテスト作成
2. サーバー側: リポジトリ、ユースケースのテスト作成

### Phase 4: 低優先度のテスト作成
1. クライアント側: ユーティリティコンポーネントのテスト作成
2. サーバー側: ユーティリティ関数、ヘルパー関数のテスト作成

### Phase 5: テストの実行と確認
1. 作成したテストが正常に実行されることを確認
2. 既存のテストが正常に動作することを確認
3. テストエラーが発生しないことを確認

## リスクと対策

### リスク1: テストが不足している箇所の見落とし
**対策**: 体系的にファイルを確認し、対応するテストファイルが存在するか確認

### リスク2: 既存のテストの失敗
**対策**: 既存のテストが正常に動作することを確認し、必要に応じて修正

### リスク3: テストの実行時間の増加
**対策**: テストの実行時間を考慮し、必要に応じて最適化

### リスク4: サーバー側のテストで認証エラー
**対策**: `APP_ENV=test`を必ず指定してテストを実行

## 参考情報

### 関連ドキュメント
- Next.js App Routerドキュメント
- React Testing Libraryドキュメント
- Jestドキュメント
- Playwrightドキュメント
- Go Testingドキュメント
- testifyドキュメント
- 既存のプロジェクトドキュメント

### 関連Issue
- https://github.com/taku-o/go-webdb-template/issues/151: 本設計書の元となったIssue

### 技術スタック
- **クライアント側**:
  - フレームワーク: Next.js 14+ (App Router)
  - 言語: TypeScript 5+
  - テストフレームワーク: Jest、React Testing Library、Playwright
  - APIモック: MSW (Mock Service Worker)
- **サーバー側**:
  - 言語: Go 1.21+
  - テストフレームワーク: Go標準テスト（`testing`パッケージ）、`testify`
  - HTTPテスト: `net/http/httptest`

### 既存のテストパターン
- クライアント側:
  - Jest統合テスト: `client/src/__tests__/integration/users-page.test.tsx`を参考
  - Jestコンポーネントテスト: `client/src/__tests__/components/TodayApiButton.test.tsx`を参考
  - Jestライブラリテスト: `client/src/__tests__/lib/api.test.ts`を参考
  - Playwright E2Eテスト: `client/e2e/`配下のファイルを参考
- サーバー側:
  - Handlerテスト: `server/internal/api/handler/today_handler_test.go`を参考
  - Usecaseテスト: `server/internal/usecase/api/dm_user_usecase_test.go`を参考
  - Serviceテスト: `server/internal/service/api_key_service_test.go`を参考
  - Repositoryテスト: `server/internal/repository/dm_user_repository_test.go`を参考
  - 統合テスト: `server/test/integration/`配下のファイルを参考
