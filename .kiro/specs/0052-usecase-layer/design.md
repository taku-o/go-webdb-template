# APIサーバーにusecase層を導入する設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、APIサーバーにusecase層を導入し、ビジネスロジックとドメインロジックを明確に分離するための詳細設計を定義する。すべてのAPIエンドポイントにusecase層を導入し、レイヤー構造を `handler -> usecase -> service -> repository -> model` に変更する。

### 1.2 設計の範囲
- usecase層のディレクトリ構造と命名規則の定義
- すべてのAPIエンドポイントに対応するusecase層の設計
- service層の整理と新規作成が必要なserviceの設計
- handler層の修正設計
- 依存関係の注入設計
- 初期化処理の設計
- エラーハンドリング設計
- テスト設計

### 1.3 設計方針
- **薄いusecase層**: 現在の実装は複雑なことをしていないため、単純に薄いusecase層を差し込む
- **一貫性**: すべてのAPIエンドポイントで一貫したレイヤー構造を維持
- **後方互換性**: 既存のAPIの動作に影響を与えない
- **段階的な実装**: 各APIエンドポイントごとにusecase層を実装
- **既存コードの最小限の変更**: 既存のservice層は可能な限り変更せず、usecase層を追加する

## 2. アーキテクチャ設計

### 2.1 レイヤー構造

#### 2.1.1 変更前のレイヤー構造
```
api(handler) -> service(ビジネスロジック + ドメインロジック) -> repository -> model
```

#### 2.1.2 変更後のレイヤー構造
```
api(handler) -> usecase(ビジネスロジック) -> service(ドメインロジック) -> repository -> model
```

### 2.2 ディレクトリ構造

```
server/internal/
├── api/
│   └── handler/          # handler層（既存）
│       ├── today_handler.go
│       ├── dm_user_handler.go
│       ├── dm_post_handler.go
│       ├── email_handler.go
│       ├── dm_jobqueue_handler.go
│       └── upload_handler.go
├── usecase/              # usecase層（新規作成）
│   ├── today_usecase.go
│   ├── dm_user_usecase.go
│   ├── dm_post_usecase.go
│   ├── email_usecase.go
│   ├── dm_jobqueue_usecase.go
│   └── upload_usecase.go
├── service/              # service層（既存、必要に応じて整理）
│   ├── dm_user_service.go
│   ├── dm_post_service.go
│   ├── date_service.go   # 新規作成（today API用）
│   └── ...
└── repository/           # repository層（既存、変更なし）
    └── ...
```

### 2.3 各レイヤーの責務

#### 2.3.1 handler層（api層）
- **責務**:
  - HTTPリクエスト/レスポンスの処理
  - 認証・認可チェック（`auth.CheckAccessLevel`）
  - 入力値の形式バリデーション（UUID形式チェックなど）
  - エラーハンドリングとHTTPステータスコードの設定
  - usecase層の呼び出し
- **制約**:
  - ビジネスロジックは実装しない（usecase層に委譲）
  - service層を直接呼び出さない（usecase層を経由）

#### 2.3.2 usecase層
- **責務**:
  - ビジネスロジックの実装
  - 複数のserviceを組み合わせた処理
  - トランザクション管理（必要に応じて）
  - ビジネスルールの適用
- **制約**:
  - usecaseから別のusecaseは呼び出さない
  - usecaseは複数、もしくは一つのserviceを呼び出して処理を行う
  - repository層を直接呼び出さない（service層を経由）

#### 2.3.3 service層
- **責務**:
  - ドメインロジックの実装
  - 特定の分類（ドメイン）の処理を行う
  - ドメイン固有のバリデーション
  - ドメイン固有のビジネスルール
- **制約**:
  - 単一のドメインに特化した処理を行う
  - ビジネスロジックはusecase層に移譲

## 3. usecase層の設計

### 3.1 命名規則とファイル構造

#### 3.1.1 命名規則
- **ファイル名**: `{機能名}_usecase.go`
- **パッケージ名**: `usecase`
- **構造体名**: `{機能名}Usecase`（例: `TodayUsecase`, `DmUserUsecase`）
- **関数名**: `New{機能名}Usecase`（例: `NewTodayUsecase`, `NewDmUserUsecase`）

#### 3.1.2 基本構造
```go
package usecase

import (
    "context"
    // 必要なimport
)

// {機能名}Usecase は{機能名}関連のビジネスロジックを担当
type {機能名}Usecase struct {
    // service層への依存
}

// New{機能名}Usecase は新しい{機能名}Usecaseを作成
func New{機能名}Usecase(/* service層の依存 */) *{機能名}Usecase {
    return &{機能名}Usecase{
        // 依存関係の注入
    }
}

// ビジネスロジックのメソッド
func (u *{機能名}Usecase) {メソッド名}(ctx context.Context, /* パラメータ */) (/* 戻り値 */, error) {
    // ビジネスロジックの実装
    // service層を呼び出して処理を実行
}
```

### 3.2 各APIエンドポイントのusecase層設計

#### 3.2.1 TodayUsecase

**ファイル**: `server/internal/usecase/today_usecase.go`

**依存関係**:
- `*service.DateService`（新規作成）

**メソッド**:
- `GetToday(ctx context.Context) (string, error)`: 今日の日付を取得

**実装内容**:
- `DateService.GetToday()`を呼び出して日付を取得
- 現時点では薄い実装（service層の呼び出しのみ）

**設計例**:
```go
package usecase

import (
    "context"
    "github.com/taku-o/go-webdb-template/internal/service"
)

type TodayUsecase struct {
    dateService *service.DateService
}

func NewTodayUsecase(dateService *service.DateService) *TodayUsecase {
    return &TodayUsecase{
        dateService: dateService,
    }
}

func (u *TodayUsecase) GetToday(ctx context.Context) (string, error) {
    return u.dateService.GetToday(ctx)
}
```

#### 3.2.2 DmUserUsecase

**ファイル**: `server/internal/usecase/dm_user_usecase.go`

**依存関係**:
- `*service.DmUserService`（既存）

**メソッド**:
- `CreateDmUser(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error)`
- `GetDmUser(ctx context.Context, id string) (*model.DmUser, error)`
- `ListDmUsers(ctx context.Context, limit, offset int) ([]*model.DmUser, error)`
- `UpdateDmUser(ctx context.Context, id string, req *model.UpdateDmUserRequest) (*model.DmUser, error)`
- `DeleteDmUser(ctx context.Context, id string) error`

**実装内容**:
- 既存の`DmUserService`のメソッドをそのまま呼び出す（薄い実装）
- 将来的にビジネスロジックが複雑になった場合に備えて、usecase層で実装

**設計例**:
```go
package usecase

import (
    "context"
    "github.com/taku-o/go-webdb-template/internal/model"
    "github.com/taku-o/go-webdb-template/internal/service"
)

type DmUserUsecase struct {
    dmUserService *service.DmUserService
}

func NewDmUserUsecase(dmUserService *service.DmUserService) *DmUserUsecase {
    return &DmUserUsecase{
        dmUserService: dmUserService,
    }
}

func (u *DmUserUsecase) CreateDmUser(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error) {
    return u.dmUserService.CreateDmUser(ctx, req)
}

// 他のメソッドも同様に実装
```

#### 3.2.3 DmPostUsecase

**ファイル**: `server/internal/usecase/dm_post_usecase.go`

**依存関係**:
- `*service.DmPostService`（既存）

**メソッド**:
- `CreateDmPost(ctx context.Context, req *model.CreateDmPostRequest) (*model.DmPost, error)`
- `GetDmPost(ctx context.Context, id string, userID string) (*model.DmPost, error)`
- `ListDmPosts(ctx context.Context, limit, offset int) ([]*model.DmPost, error)`
- `ListDmPostsByUser(ctx context.Context, userID string, limit, offset int) ([]*model.DmPost, error)`
- `GetDmUserPosts(ctx context.Context, limit, offset int) ([]*model.DmUserPost, error)`
- `UpdateDmPost(ctx context.Context, id string, userID string, req *model.UpdateDmPostRequest) (*model.DmPost, error)`
- `DeleteDmPost(ctx context.Context, id string, userID string) error`

**実装内容**:
- 既存の`DmPostService`のメソッドをそのまま呼び出す（薄い実装）

#### 3.2.4 EmailUsecase

**ファイル**: `server/internal/usecase/email_usecase.go`

**依存関係**:
- `*email.EmailService`（既存）
- `*email.TemplateService`（既存）

**メソッド**:
- `SendEmail(ctx context.Context, to []string, template string, data map[string]interface{}) error`

**実装内容**:
- `TemplateService`でテンプレートをレンダリング
- `EmailService`でメールを送信
- 複数のserviceを組み合わせたビジネスロジック

**設計例**:
```go
package usecase

import (
    "context"
    "github.com/taku-o/go-webdb-template/internal/service/email"
)

type EmailUsecase struct {
    emailService    *email.EmailService
    templateService *email.TemplateService
}

func NewEmailUsecase(emailService *email.EmailService, templateService *email.TemplateService) *EmailUsecase {
    return &EmailUsecase{
        emailService:    emailService,
        templateService: templateService,
    }
}

func (u *EmailUsecase) SendEmail(ctx context.Context, to []string, template string, data map[string]interface{}) error {
    // テンプレートからメール本文を生成
    body, err := u.templateService.Render(template, data)
    if err != nil {
        return fmt.Errorf("failed to render template: %w", err)
    }

    // テンプレートから件名を取得
    subject, err := u.templateService.GetSubject(template)
    if err != nil {
        return fmt.Errorf("failed to get subject: %w", err)
    }

    // メール送信
    if err := u.emailService.SendEmail(ctx, to, subject, body); err != nil {
        return fmt.Errorf("failed to send email: %w", err)
    }

    return nil
}
```

#### 3.2.5 DmJobqueueUsecase

**ファイル**: `server/internal/usecase/dm_jobqueue_usecase.go`

**依存関係**:
- `*jobqueue.Client`（既存）

**メソッド**:
- `RegisterJob(ctx context.Context, message string, delaySeconds int, maxRetry int) (string, error)`

**実装内容**:
- `jobqueue.Client`を使用してジョブを登録
- ビジネスロジック（メッセージのデフォルト値設定など）を実装

**設計例**:
```go
package usecase

import (
    "context"
    "encoding/json"
    "github.com/taku-o/go-webdb-template/internal/service/jobqueue"
)

type DmJobqueueUsecase struct {
    jobQueueClient *jobqueue.Client
}

func NewDmJobqueueUsecase(jobQueueClient *jobqueue.Client) *DmJobqueueUsecase {
    return &DmJobqueueUsecase{
        jobQueueClient: jobQueueClient,
    }
}

func (u *DmJobqueueUsecase) RegisterJob(ctx context.Context, message string, delaySeconds int, maxRetry int) (string, error) {
    // メッセージの設定（デフォルト値）
    if message == "" {
        message = "Job executed successfully"
    }

    // ペイロードの作成
    payload := jobqueue.DelayPrintPayload{
        Message: message,
    }
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return "", fmt.Errorf("failed to marshal payload: %w", err)
    }

    // ジョブオプションの作成
    jobOpts := &jobqueue.JobOptions{
        DelaySeconds: delaySeconds,
        MaxRetry:     maxRetry,
    }

    // ジョブをキューに登録
    info, err := u.jobQueueClient.EnqueueJob(
        ctx,
        jobqueue.JobTypeDelayPrint,
        payloadBytes,
        jobOpts,
    )
    if err != nil {
        return "", fmt.Errorf("failed to enqueue job: %w", err)
    }

    return info.ID, nil
}
```

#### 3.2.6 UploadUsecase

**ファイル**: `server/internal/usecase/upload_usecase.go`

**依存関係**:
- `*storage.Storage`（既存、または新規作成）

**メソッド**:
- アップロード関連のビジネスロジック（必要に応じて）

**実装内容**:
- `UploadHandler`はTUSプロトコルを使用しているため、usecase層の実装は最小限
- 将来的にアップロード後の処理（通知、ログ記録など）が必要になった場合に備えて、usecase層を用意

**設計例**:
```go
package usecase

import (
    "context"
    // 必要なimport
)

type UploadUsecase struct {
    // 将来的に必要な依存関係を追加
}

func NewUploadUsecase(/* 依存関係 */) *UploadUsecase {
    return &UploadUsecase{
        // 依存関係の注入
    }
}

// 将来的にアップロード後の処理が必要になった場合に実装
```

## 4. service層の設計

### 4.1 既存のservice層の整理

#### 4.1.1 DmUserService
- **現状**: ビジネスロジックとドメインロジックが混在
- **整理方針**: 現時点では変更せず、usecase層を追加することで責務を明確化
- **将来的な整理**: 必要に応じて、ドメインロジックのみを残し、ビジネスロジックをusecase層に移譲

#### 4.1.2 DmPostService
- **現状**: ビジネスロジックとドメインロジックが混在
- **整理方針**: 現時点では変更せず、usecase層を追加することで責務を明確化

#### 4.1.3 EmailService / TemplateService
- **現状**: ドメインロジックのみを担当（適切に分離されている）
- **整理方針**: 変更不要

### 4.2 新規作成が必要なservice層

#### 4.2.1 DateService

**ファイル**: `server/internal/service/date_service.go`

**責務**: 日付関連のドメインロジック

**メソッド**:
- `GetToday(ctx context.Context) (string, error)`: 今日の日付をYYYY-MM-DD形式で取得

**設計例**:
```go
package service

import (
    "context"
    "time"
)

type DateService struct{}

func NewDateService() *DateService {
    return &DateService{}
}

func (s *DateService) GetToday(ctx context.Context) (string, error) {
    return time.Now().Format("2006-01-02"), nil
}
```

## 5. handler層の修正設計

### 5.1 修正方針
- handler層からservice層への直接呼び出しを削除
- handler層からusecase層を呼び出すように変更
- 認証チェック、入力値の形式バリデーション、エラーハンドリングはhandler層で継続して行う

### 5.2 各handlerの修正内容

#### 5.2.1 TodayHandler

**変更前**:
```go
type TodayHandler struct{}

func (h *TodayHandler) GetToday(ctx context.Context) (string, error) {
    // 認証チェック
    if err := auth.CheckAccessLevel(ctx, auth.AccessLevelPrivate); err != nil {
        return "", err
    }
    // 日付取得（直接処理）
    return time.Now().Format("2006-01-02"), nil
}
```

**変更後**:
```go
type TodayHandler struct {
    todayUsecase *usecase.TodayUsecase
}

func NewTodayHandler(todayUsecase *usecase.TodayUsecase) *TodayHandler {
    return &TodayHandler{
        todayUsecase: todayUsecase,
    }
}

func (h *TodayHandler) GetToday(ctx context.Context) (string, error) {
    // 認証チェック（handler層で継続）
    if err := auth.CheckAccessLevel(ctx, auth.AccessLevelPrivate); err != nil {
        return "", err
    }
    // usecase層を呼び出し
    return h.todayUsecase.GetToday(ctx)
}
```

#### 5.2.2 DmUserHandler

**変更前**:
```go
type DmUserHandler struct {
    dmUserService *service.DmUserService
}
```

**変更後**:
```go
type DmUserHandler struct {
    dmUserUsecase *usecase.DmUserUsecase
}

func NewDmUserHandler(dmUserUsecase *usecase.DmUserUsecase) *DmUserHandler {
    return &DmUserHandler{
        dmUserUsecase: dmUserUsecase,
    }
}
```

**メソッドの変更**:
- `h.dmUserService.CreateDmUser(...)` → `h.dmUserUsecase.CreateDmUser(...)`
- 他のメソッドも同様に変更

#### 5.2.3 その他のhandler
- `DmPostHandler`: `dmPostService` → `dmPostUsecase`
- `EmailHandler`: `emailService`, `templateService` → `emailUsecase`
- `DmJobqueueHandler`: `jobQueueClient` → `dmJobqueueUsecase`
- `UploadHandler`: 現時点では変更不要（将来的にusecase層を追加する場合に備える）

## 6. 依存関係の注入設計

### 6.1 依存関係の構造

```
main.go
  ├── Repository層の初期化
  │   ├── dmUserRepo
  │   └── dmPostRepo
  ├── Service層の初期化
  │   ├── dmUserService (dmUserRepo)
  │   ├── dmPostService (dmPostRepo, dmUserRepo)
  │   └── dateService
  ├── Usecase層の初期化
  │   ├── todayUsecase (dateService)
  │   ├── dmUserUsecase (dmUserService)
  │   ├── dmPostUsecase (dmPostService)
  │   ├── emailUsecase (emailService, templateService)
  │   └── dmJobqueueUsecase (jobQueueClient)
  └── Handler層の初期化
      ├── todayHandler (todayUsecase)
      ├── dmUserHandler (dmUserUsecase)
      ├── dmPostHandler (dmPostUsecase)
      ├── emailHandler (emailUsecase)
      └── dmJobqueueHandler (dmJobqueueUsecase)
```

### 6.2 初期化処理の設計

#### 6.2.1 main.goの修正

**変更前**:
```go
// Service層の初期化
dmUserService := service.NewDmUserService(dmUserRepo)
dmPostService := service.NewDmPostService(dmPostRepo, dmUserRepo)

// Handler層の初期化
dmUserHandler := handler.NewDmUserHandler(dmUserService)
dmPostHandler := handler.NewDmPostHandler(dmPostService)
todayHandler := handler.NewTodayHandler()
```

**変更後**:
```go
// Service層の初期化
dmUserService := service.NewDmUserService(dmUserRepo)
dmPostService := service.NewDmPostService(dmPostRepo, dmUserRepo)
dateService := service.NewDateService()

// Usecase層の初期化
todayUsecase := usecase.NewTodayUsecase(dateService)
dmUserUsecase := usecase.NewDmUserUsecase(dmUserService)
dmPostUsecase := usecase.NewDmPostUsecase(dmPostService)
emailUsecase := usecase.NewEmailUsecase(emailService, templateService)
dmJobqueueUsecase := usecase.NewDmJobqueueUsecase(jobQueueClient)

// Handler層の初期化
todayHandler := handler.NewTodayHandler(todayUsecase)
dmUserHandler := handler.NewDmUserHandler(dmUserUsecase)
dmPostHandler := handler.NewDmPostHandler(dmPostUsecase)
emailHandler := handler.NewEmailHandler(emailUsecase)  // 修正が必要
dmJobqueueHandler := handler.NewDmJobqueueHandler(dmJobqueueUsecase)  // 修正が必要
```

#### 6.2.2 EmailHandlerの修正

**変更前**:
```go
type EmailHandler struct {
    emailService    *email.EmailService
    templateService *email.TemplateService
}

func NewEmailHandler(emailService *email.EmailService, templateService *email.TemplateService) *EmailHandler {
    return &EmailHandler{
        emailService:    emailService,
        templateService: templateService,
    }
}
```

**変更後**:
```go
type EmailHandler struct {
    emailUsecase *usecase.EmailUsecase
}

func NewEmailHandler(emailUsecase *usecase.EmailUsecase) *EmailHandler {
    return &EmailHandler{
        emailUsecase: emailUsecase,
    }
}
```

#### 6.2.3 DmJobqueueHandlerの修正

**変更前**:
```go
type DmJobqueueHandler struct {
    jobQueueClient *jobqueue.Client
}

func NewDmJobqueueHandler(jobQueueClient *jobqueue.Client) *DmJobqueueHandler {
    return &DmJobqueueHandler{
        jobQueueClient: jobQueueClient,
    }
}
```

**変更後**:
```go
type DmJobqueueHandler struct {
    dmJobqueueUsecase *usecase.DmJobqueueUsecase
}

func NewDmJobqueueHandler(dmJobqueueUsecase *usecase.DmJobqueueUsecase) *DmJobqueueHandler {
    return &DmJobqueueHandler{
        dmJobqueueUsecase: dmJobqueueUsecase,
    }
}
```

## 7. エラーハンドリング設計

### 7.1 エラーハンドリングの方針
- **handler層**: HTTPステータスコードの設定、エラーメッセージの整形
- **usecase層**: ビジネスロジックのエラーを返す（エラーメッセージはそのまま返す）
- **service層**: ドメインロジックのエラーを返す（エラーメッセージはそのまま返す）

### 7.2 エラーフロー

```
service層でエラー発生
  ↓
usecase層でエラーをそのまま返す
  ↓
handler層でエラーを受け取り、HTTPステータスコードを設定
  ↓
HTTPレスポンスとして返す
```

### 7.3 エラーハンドリングの実装例

**handler層**:
```go
result, err := h.usecase.SomeMethod(ctx, input)
if err != nil {
    // エラーの種類に応じてHTTPステータスコードを設定
    if errors.Is(err, ErrNotFound) {
        return nil, huma.Error404NotFound(err.Error())
    }
    return nil, huma.Error500InternalServerError(err.Error())
}
```

## 8. テスト設計

### 8.1 テストの方針
- **各レイヤーの独立テスト**: 各レイヤーを独立してテスト可能にする
- **モックの使用**: 下位レイヤーをモック化してテストする
- **既存テストの維持**: 既存のテストが正常に動作することを確認

### 8.2 テスト構造

#### 8.2.1 usecase層のテスト
- **ファイル**: `server/internal/usecase/{機能名}_usecase_test.go`
- **モック**: service層をモック化
- **テスト内容**: ビジネスロジックのテスト

**テスト例**:
```go
func TestTodayUsecase_GetToday(t *testing.T) {
    // service層のモックを作成
    mockDateService := &MockDateService{}
    mockDateService.On("GetToday", mock.Anything).Return("2026-01-10", nil)

    // usecase層を作成
    usecase := NewTodayUsecase(mockDateService)

    // テスト実行
    result, err := usecase.GetToday(context.Background())
    
    // アサーション
    assert.NoError(t, err)
    assert.Equal(t, "2026-01-10", result)
    mockDateService.AssertExpectations(t)
}
```

#### 8.2.2 handler層のテスト
- **ファイル**: `server/internal/api/handler/{機能名}_handler_test.go`
- **モック**: usecase層をモック化
- **テスト内容**: HTTPリクエスト/レスポンスのテスト、認証チェックのテスト

**テスト例**:
```go
func TestTodayHandler_GetToday(t *testing.T) {
    // usecase層のモックを作成
    mockUsecase := &MockTodayUsecase{}
    mockUsecase.On("GetToday", mock.Anything).Return("2026-01-10", nil)

    // handler層を作成
    handler := NewTodayHandler(mockUsecase)

    // テスト実行
    result, err := handler.GetToday(context.Background())
    
    // アサーション
    assert.NoError(t, err)
    assert.Equal(t, "2026-01-10", result)
    mockUsecase.AssertExpectations(t)
}
```

### 8.3 既存テストの修正
- **handler層のテスト**: usecase層をモック化するように修正
- **統合テスト**: 既存の統合テストは正常に動作することを確認

## 9. 実装順序

### 9.1 実装フェーズ

#### フェーズ1: 基盤の準備
1. `server/internal/usecase/` ディレクトリの作成
2. `DateService`の作成（today API用）

#### フェーズ2: today APIの実装（パターン確立）
1. `TodayUsecase`の作成
2. `TodayHandler`の修正
3. `main.go`の初期化処理の修正
4. テストの作成・修正

#### フェーズ3: その他のAPIエンドポイントの実装
1. `DmUserUsecase`の作成と`DmUserHandler`の修正
2. `DmPostUsecase`の作成と`DmPostHandler`の修正
3. `EmailUsecase`の作成と`EmailHandler`の修正
4. `DmJobqueueUsecase`の作成と`DmJobqueueHandler`の修正
5. `UploadUsecase`の作成（必要に応じて）
6. 各APIのテストの作成・修正

### 9.2 実装時の注意事項
- **段階的な実装**: 1つのAPIエンドポイントずつ実装し、動作確認を行う
- **既存テストの確認**: 各フェーズで既存のテストが正常に動作することを確認
- **API動作の確認**: 各APIエンドポイントが正常に動作することを確認

## 10. 参考情報

### 10.1 関連ドキュメント
- 要件定義書: `.kiro/specs/0052-usecase-layer/requirements.md`
- Issue: GitHub Issue #107

### 10.2 既存実装の参考
- handler層: `server/internal/api/handler/`
- service層: `server/internal/service/`
- repository層: `server/internal/repository/`

### 10.3 技術スタック
- **言語**: Go
- **フレームワーク**: Huma API, Echo
- **テストフレームワーク**: testify/mock, testify/assert
