# APIサーバーにusecase層を導入する実装タスク一覧

## 概要
APIサーバーにusecase層を導入し、すべてのAPIエンドポイントでレイヤー構造を `handler -> usecase -> service -> repository -> model` に変更するためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: 基盤の準備

#### タスク 1.1: usecase層ディレクトリの作成
**目的**: usecase層を配置するディレクトリを作成する。

**作業内容**:
- `server/internal/usecase/` ディレクトリを作成

**実装内容**:
- ディレクトリパス: `server/internal/usecase/`

**受け入れ基準**:
- `server/internal/usecase/` ディレクトリが作成されている

- _Requirements: 3.1.1, 6.1_
- _Design: 2.2_

---

#### タスク 1.2: DateServiceの作成
**目的**: today API用の日付関連ドメインロジックを担当するservice層を作成する。

**作業内容**:
- `server/internal/service/date_service.go` を作成
- `DateService`構造体を定義
- `GetToday`メソッドを実装

**実装内容**:
- パッケージ名: `service`
- 構造体名: `DateService`
- コンストラクタ: `NewDateService() *DateService`
- メソッド: `GetToday(ctx context.Context) (string, error)`
  - 今日の日付をYYYY-MM-DD形式で返す
  - `time.Now().Format("2006-01-02")`を使用

**ファイル構成**:
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

**受け入れ基準**:
- `server/internal/service/date_service.go` が作成されている
- `DateService`構造体が定義されている
- `NewDateService`関数が実装されている
- `GetToday`メソッドが実装されている
- `GetToday`メソッドがYYYY-MM-DD形式の日付を返す

- _Requirements: 3.3.3, 6.2_
- _Design: 4.2.1_

---

### Phase 2: today APIの実装（パターン確立）

#### タスク 2.1: TodayUsecaseの作成
**目的**: today API用のビジネスロジックを担当するusecase層を作成する。

**作業内容**:
- `server/internal/usecase/today_usecase.go` を作成
- `TodayUsecase`構造体を定義
- `GetToday`メソッドを実装

**実装内容**:
- パッケージ名: `usecase`
- 構造体名: `TodayUsecase`
- 依存関係: `*service.DateService`
- コンストラクタ: `NewTodayUsecase(dateService *service.DateService) *TodayUsecase`
- メソッド: `GetToday(ctx context.Context) (string, error)`
  - `DateService.GetToday()`を呼び出して日付を取得
  - 現時点では薄い実装（service層の呼び出しのみ）

**ファイル構成**:
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

**受け入れ基準**:
- `server/internal/usecase/today_usecase.go` が作成されている
- `TodayUsecase`構造体が定義されている
- `NewTodayUsecase`関数が実装されている
- `GetToday`メソッドが実装されている
- `GetToday`メソッドが`DateService.GetToday()`を呼び出している

- _Requirements: 3.3.2, 6.1_
- _Design: 3.2.1_

---

#### タスク 2.2: TodayHandlerの修正
**目的**: `TodayHandler`をusecase層を呼び出すように修正する。

**作業内容**:
- `server/internal/api/handler/today_handler.go` を修正
- `TodayHandler`構造体に`todayUsecase`フィールドを追加
- `NewTodayHandler`関数を修正してusecase層を受け取るように変更
- `GetToday`メソッドを修正してusecase層を呼び出すように変更

**実装内容**:
- `TodayHandler`構造体:
  - 変更前: フィールドなし
  - 変更後: `todayUsecase *usecase.TodayUsecase`
- `NewTodayHandler`関数:
  - 変更前: `func NewTodayHandler() *TodayHandler`
  - 変更後: `func NewTodayHandler(todayUsecase *usecase.TodayUsecase) *TodayHandler`
- `GetToday`メソッド:
  - 認証チェックはhandler層で継続して行う
  - 日付取得処理をusecase層に委譲

**変更前**:
```go
type TodayHandler struct{}

func NewTodayHandler() *TodayHandler {
    return &TodayHandler{}
}

func (h *TodayHandler) GetToday(ctx context.Context) (string, error) {
    if err := auth.CheckAccessLevel(ctx, auth.AccessLevelPrivate); err != nil {
        return "", err
    }
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
    if err := auth.CheckAccessLevel(ctx, auth.AccessLevelPrivate); err != nil {
        return "", err
    }
    return h.todayUsecase.GetToday(ctx)
}
```

**受け入れ基準**:
- `TodayHandler`構造体に`todayUsecase`フィールドが追加されている
- `NewTodayHandler`関数がusecase層を受け取るように修正されている
- `GetToday`メソッドがusecase層を呼び出すように修正されている
- 認証チェックがhandler層で継続して行われている
- 既存のtoday APIの動作が変わらない（同じレスポンスを返す）

- _Requirements: 3.3.4, 6.3_
- _Design: 5.2.1_

---

#### タスク 2.3: main.goの初期化処理の修正（today API用）
**目的**: `main.go`でtoday API用のusecase層とservice層を初期化し、handler層に注入する。

**作業内容**:
- `server/cmd/server/main.go` を修正
- `DateService`の初期化を追加
- `TodayUsecase`の初期化を追加
- `TodayHandler`の初期化を修正（usecase層を受け取るように変更）

**実装内容**:
- `DateService`の初期化:
  ```go
  dateService := service.NewDateService()
  ```
- `TodayUsecase`の初期化:
  ```go
  todayUsecase := usecase.NewTodayUsecase(dateService)
  ```
- `TodayHandler`の初期化:
  ```go
  // 変更前: todayHandler := handler.NewTodayHandler()
  // 変更後:
  todayHandler := handler.NewTodayHandler(todayUsecase)
  ```

**受け入れ基準**:
- `DateService`の初期化が追加されている
- `TodayUsecase`の初期化が追加されている
- `TodayHandler`の初期化がusecase層を受け取るように修正されている
- 初期化の順序が正しい（service → usecase → handler）

- _Requirements: 3.4.2, 6.4_
- _Design: 6.2.1_

---

#### タスク 2.4: TodayUsecaseのテスト作成
**目的**: `TodayUsecase`の単体テストを作成する。

**作業内容**:
- `server/internal/usecase/today_usecase_test.go` を作成
- `DateService`をモック化
- `GetToday`メソッドのテストを実装

**実装内容**:
- テストファイル: `server/internal/usecase/today_usecase_test.go`
- モック: `DateService`をモック化
- テストケース:
  - 正常系: 日付が正しく取得できること
  - 異常系: service層でエラーが発生した場合の処理

**受け入れ基準**:
- `server/internal/usecase/today_usecase_test.go` が作成されている
- `DateService`がモック化されている
- 正常系のテストが実装されている
- 異常系のテストが実装されている（必要に応じて）
- テストが正常に実行できる

- _Requirements: 6.5_
- _Design: 8.2.1_

---

#### タスク 2.5: TodayHandlerのテスト修正
**目的**: `TodayHandler`のテストをusecase層を使用するように修正する。

**作業内容**:
- `server/internal/api/handler/today_handler_test.go` を修正
- `TodayUsecase`をモック化
- テストケースを修正

**実装内容**:
- モック: `TodayUsecase`をモック化
- テストケース:
  - 正常系: usecase層を呼び出して日付が取得できること
  - 認証エラー: 認証チェックが正しく動作すること

**受け入れ基準**:
- `TodayHandler`のテストがusecase層を使用するように修正されている
- `TodayUsecase`がモック化されている
- 既存のテストが正常に動作する
- テストが正常に実行できる

- _Requirements: 6.5_
- _Design: 8.2.2_

---

#### タスク 2.6: today APIの動作確認
**目的**: today APIが正常に動作することを確認する。

**作業内容**:
- サーバーを起動
- today APIを実行
- レスポンスを確認

**実装内容**:
- APIエンドポイント: `GET /api/today`
- 認証: Auth0 JWT（private API）
- 期待されるレスポンス: `{"date": "YYYY-MM-DD"}`

**受け入れ基準**:
- today APIが正常に動作する
- レスポンスが正しい形式である
- 認証チェックが正しく動作する
- 既存のAPIの動作が変わらない

- _Requirements: 6.6_
- _Design: 9.2_

---

### Phase 3: dm_user APIの実装

#### タスク 3.1: DmUserUsecaseの作成
**目的**: dm_user API用のビジネスロジックを担当するusecase層を作成する。

**作業内容**:
- `server/internal/usecase/dm_user_usecase.go` を作成
- `DmUserUsecase`構造体を定義
- すべてのメソッドを実装

**実装内容**:
- パッケージ名: `usecase`
- 構造体名: `DmUserUsecase`
- 依存関係: `*service.DmUserService`
- コンストラクタ: `NewDmUserUsecase(dmUserService *service.DmUserService) *DmUserUsecase`
- メソッド:
  - `CreateDmUser(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error)`
  - `GetDmUser(ctx context.Context, id string) (*model.DmUser, error)`
  - `ListDmUsers(ctx context.Context, limit, offset int) ([]*model.DmUser, error)`
  - `UpdateDmUser(ctx context.Context, id string, req *model.UpdateDmUserRequest) (*model.DmUser, error)`
  - `DeleteDmUser(ctx context.Context, id string) error`
- 実装内容: 既存の`DmUserService`のメソッドをそのまま呼び出す（薄い実装）

**受け入れ基準**:
- `server/internal/usecase/dm_user_usecase.go` が作成されている
- `DmUserUsecase`構造体が定義されている
- すべてのメソッドが実装されている
- すべてのメソッドが`DmUserService`の対応するメソッドを呼び出している

- _Requirements: 3.3.2, 6.1_
- _Design: 3.2.2_

---

#### タスク 3.2: DmUserHandlerの修正
**目的**: `DmUserHandler`をusecase層を呼び出すように修正する。

**作業内容**:
- `server/internal/api/handler/dm_user_handler.go` を修正
- `DmUserHandler`構造体のフィールドを`dmUserService`から`dmUserUsecase`に変更
- `NewDmUserHandler`関数を修正してusecase層を受け取るように変更
- すべてのメソッドでservice層の呼び出しをusecase層の呼び出しに変更

**実装内容**:
- `DmUserHandler`構造体:
  - 変更前: `dmUserService *service.DmUserService`
  - 変更後: `dmUserUsecase *usecase.DmUserUsecase`
- `NewDmUserHandler`関数:
  - 変更前: `func NewDmUserHandler(dmUserService *service.DmUserService) *DmUserHandler`
  - 変更後: `func NewDmUserHandler(dmUserUsecase *usecase.DmUserUsecase) *DmUserHandler`
- すべてのメソッド:
  - `h.dmUserService.*` → `h.dmUserUsecase.*`

**受け入れ基準**:
- `DmUserHandler`構造体のフィールドが`dmUserUsecase`に変更されている
- `NewDmUserHandler`関数がusecase層を受け取るように修正されている
- すべてのメソッドでusecase層を呼び出すように修正されている
- 認証チェック、バリデーション、エラーハンドリングがhandler層で継続して行われている

- _Requirements: 3.3.4, 6.3_
- _Design: 5.2.2_

---

#### タスク 3.3: main.goの初期化処理の修正（dm_user API用）
**目的**: `main.go`でdm_user API用のusecase層を初期化し、handler層に注入する。

**作業内容**:
- `server/cmd/server/main.go` を修正
- `DmUserUsecase`の初期化を追加
- `DmUserHandler`の初期化を修正（usecase層を受け取るように変更）

**実装内容**:
- `DmUserUsecase`の初期化:
  ```go
  dmUserUsecase := usecase.NewDmUserUsecase(dmUserService)
  ```
- `DmUserHandler`の初期化:
  ```go
  // 変更前: dmUserHandler := handler.NewDmUserHandler(dmUserService)
  // 変更後:
  dmUserHandler := handler.NewDmUserHandler(dmUserUsecase)
  ```

**受け入れ基準**:
- `DmUserUsecase`の初期化が追加されている
- `DmUserHandler`の初期化がusecase層を受け取るように修正されている
- 初期化の順序が正しい（service → usecase → handler）

- _Requirements: 3.4.2, 6.4_
- _Design: 6.2.1_

---

#### タスク 3.4: DmUserUsecaseのテスト作成
**目的**: `DmUserUsecase`の単体テストを作成する。

**作業内容**:
- `server/internal/usecase/dm_user_usecase_test.go` を作成
- `DmUserService`をモック化
- すべてのメソッドのテストを実装

**受け入れ基準**:
- `server/internal/usecase/dm_user_usecase_test.go` が作成されている
- `DmUserService`がモック化されている
- すべてのメソッドのテストが実装されている
- テストが正常に実行できる

- _Requirements: 6.5_
- _Design: 8.2.1_

---

#### タスク 3.5: DmUserHandlerのテスト修正
**目的**: `DmUserHandler`のテストをusecase層を使用するように修正する。

**作業内容**:
- `server/internal/api/handler/dm_user_handler_test.go` を修正
- `DmUserUsecase`をモック化
- テストケースを修正

**受け入れ基準**:
- `DmUserHandler`のテストがusecase層を使用するように修正されている
- `DmUserUsecase`がモック化されている
- 既存のテストが正常に動作する
- テストが正常に実行できる

- _Requirements: 6.5_
- _Design: 8.2.2_

---

#### タスク 3.6: dm_user APIの動作確認
**目的**: dm_user APIが正常に動作することを確認する。

**作業内容**:
- サーバーを起動
- dm_user APIのすべてのエンドポイントを実行
- レスポンスを確認

**実装内容**:
- APIエンドポイント:
  - `POST /api/dm-users`
  - `GET /api/dm-users/{id}`
  - `GET /api/dm-users`
  - `PUT /api/dm-users/{id}`
  - `DELETE /api/dm-users/{id}`
  - `GET /api/export/dm-users/csv`

**受け入れ基準**:
- すべてのdm_user APIエンドポイントが正常に動作する
- レスポンスが正しい形式である
- 既存のAPIの動作が変わらない

- _Requirements: 6.6_
- _Design: 9.2_

---

### Phase 4: dm_post APIの実装

#### タスク 4.1: DmPostUsecaseの作成
**目的**: dm_post API用のビジネスロジックを担当するusecase層を作成する。

**作業内容**:
- `server/internal/usecase/dm_post_usecase.go` を作成
- `DmPostUsecase`構造体を定義
- すべてのメソッドを実装

**実装内容**:
- パッケージ名: `usecase`
- 構造体名: `DmPostUsecase`
- 依存関係: `*service.DmPostService`
- コンストラクタ: `NewDmPostUsecase(dmPostService *service.DmPostService) *DmPostUsecase`
- メソッド:
  - `CreateDmPost(ctx context.Context, req *model.CreateDmPostRequest) (*model.DmPost, error)`
  - `GetDmPost(ctx context.Context, id string, userID string) (*model.DmPost, error)`
  - `ListDmPosts(ctx context.Context, limit, offset int) ([]*model.DmPost, error)`
  - `ListDmPostsByUser(ctx context.Context, userID string, limit, offset int) ([]*model.DmPost, error)`
  - `GetDmUserPosts(ctx context.Context, limit, offset int) ([]*model.DmUserPost, error)`
  - `UpdateDmPost(ctx context.Context, id string, userID string, req *model.UpdateDmPostRequest) (*model.DmPost, error)`
  - `DeleteDmPost(ctx context.Context, id string, userID string) error`
- 実装内容: 既存の`DmPostService`のメソッドをそのまま呼び出す（薄い実装）

**受け入れ基準**:
- `server/internal/usecase/dm_post_usecase.go` が作成されている
- `DmPostUsecase`構造体が定義されている
- すべてのメソッドが実装されている
- すべてのメソッドが`DmPostService`の対応するメソッドを呼び出している

- _Requirements: 3.3.2, 6.1_
- _Design: 3.2.3_

---

#### タスク 4.2: DmPostHandlerの修正
**目的**: `DmPostHandler`をusecase層を呼び出すように修正する。

**作業内容**:
- `server/internal/api/handler/dm_post_handler.go` を修正
- `DmPostHandler`構造体のフィールドを`dmPostService`から`dmPostUsecase`に変更
- `NewDmPostHandler`関数を修正してusecase層を受け取るように変更
- すべてのメソッドでservice層の呼び出しをusecase層の呼び出しに変更

**実装内容**:
- `DmPostHandler`構造体:
  - 変更前: `dmPostService *service.DmPostService`
  - 変更後: `dmPostUsecase *usecase.DmPostUsecase`
- `NewDmPostHandler`関数:
  - 変更前: `func NewDmPostHandler(dmPostService *service.DmPostService) *DmPostHandler`
  - 変更後: `func NewDmPostHandler(dmPostUsecase *usecase.DmPostUsecase) *DmPostHandler`
- すべてのメソッド:
  - `h.dmPostService.*` → `h.dmPostUsecase.*`

**受け入れ基準**:
- `DmPostHandler`構造体のフィールドが`dmPostUsecase`に変更されている
- `NewDmPostHandler`関数がusecase層を受け取るように修正されている
- すべてのメソッドでusecase層を呼び出すように修正されている
- 認証チェック、バリデーション、エラーハンドリングがhandler層で継続して行われている

- _Requirements: 3.3.4, 6.3_
- _Design: 5.2.2_

---

#### タスク 4.3: main.goの初期化処理の修正（dm_post API用）
**目的**: `main.go`でdm_post API用のusecase層を初期化し、handler層に注入する。

**作業内容**:
- `server/cmd/server/main.go` を修正
- `DmPostUsecase`の初期化を追加
- `DmPostHandler`の初期化を修正（usecase層を受け取るように変更）

**実装内容**:
- `DmPostUsecase`の初期化:
  ```go
  dmPostUsecase := usecase.NewDmPostUsecase(dmPostService)
  ```
- `DmPostHandler`の初期化:
  ```go
  // 変更前: dmPostHandler := handler.NewDmPostHandler(dmPostService)
  // 変更後:
  dmPostHandler := handler.NewDmPostHandler(dmPostUsecase)
  ```

**受け入れ基準**:
- `DmPostUsecase`の初期化が追加されている
- `DmPostHandler`の初期化がusecase層を受け取るように修正されている
- 初期化の順序が正しい（service → usecase → handler）

- _Requirements: 3.4.2, 6.4_
- _Design: 6.2.1_

---

#### タスク 4.4: DmPostUsecaseのテスト作成
**目的**: `DmPostUsecase`の単体テストを作成する。

**作業内容**:
- `server/internal/usecase/dm_post_usecase_test.go` を作成
- `DmPostService`をモック化
- すべてのメソッドのテストを実装

**受け入れ基準**:
- `server/internal/usecase/dm_post_usecase_test.go` が作成されている
- `DmPostService`がモック化されている
- すべてのメソッドのテストが実装されている
- テストが正常に実行できる

- _Requirements: 6.5_
- _Design: 8.2.1_

---

#### タスク 4.5: DmPostHandlerのテスト修正
**目的**: `DmPostHandler`のテストをusecase層を使用するように修正する。

**作業内容**:
- `server/internal/api/handler/dm_post_handler_test.go` を修正
- `DmPostUsecase`をモック化
- テストケースを修正

**受け入れ基準**:
- `DmPostHandler`のテストがusecase層を使用するように修正されている
- `DmPostUsecase`がモック化されている
- 既存のテストが正常に動作する
- テストが正常に実行できる

- _Requirements: 6.5_
- _Design: 8.2.2_

---

#### タスク 4.6: dm_post APIの動作確認
**目的**: dm_post APIが正常に動作することを確認する。

**作業内容**:
- サーバーを起動
- dm_post APIのすべてのエンドポイントを実行
- レスポンスを確認

**実装内容**:
- APIエンドポイント:
  - `POST /api/dm-posts`
  - `GET /api/dm-posts/{id}`
  - `GET /api/dm-posts`
  - `GET /api/dm-posts/user/{user_id}`
  - `GET /api/dm-posts/user-posts`
  - `PUT /api/dm-posts/{id}`
  - `DELETE /api/dm-posts/{id}`

**受け入れ基準**:
- すべてのdm_post APIエンドポイントが正常に動作する
- レスポンスが正しい形式である
- 既存のAPIの動作が変わらない

- _Requirements: 6.6_
- _Design: 9.2_

---

### Phase 5: email APIの実装

#### タスク 5.1: EmailUsecaseの作成
**目的**: email API用のビジネスロジックを担当するusecase層を作成する。

**作業内容**:
- `server/internal/usecase/email_usecase.go` を作成
- `EmailUsecase`構造体を定義
- `SendEmail`メソッドを実装

**実装内容**:
- パッケージ名: `usecase`
- 構造体名: `EmailUsecase`
- 依存関係: `*email.EmailService`, `*email.TemplateService`
- コンストラクタ: `NewEmailUsecase(emailService *email.EmailService, templateService *email.TemplateService) *EmailUsecase`
- メソッド: `SendEmail(ctx context.Context, to []string, template string, data map[string]interface{}) error`
  - `TemplateService`でテンプレートをレンダリング
  - `TemplateService`で件名を取得
  - `EmailService`でメールを送信
  - 複数のserviceを組み合わせたビジネスロジック

**受け入れ基準**:
- `server/internal/usecase/email_usecase.go` が作成されている
- `EmailUsecase`構造体が定義されている
- `SendEmail`メソッドが実装されている
- `SendEmail`メソッドが`TemplateService`と`EmailService`を組み合わせて使用している

- _Requirements: 3.3.2, 6.1_
- _Design: 3.2.4_

---

#### タスク 5.2: EmailHandlerの修正
**目的**: `EmailHandler`をusecase層を呼び出すように修正する。

**作業内容**:
- `server/internal/api/handler/email_handler.go` を修正
- `EmailHandler`構造体のフィールドを`emailService`, `templateService`から`emailUsecase`に変更
- `NewEmailHandler`関数を修正してusecase層を受け取るように変更
- `RegisterEmailEndpoints`内の処理をusecase層を呼び出すように変更

**実装内容**:
- `EmailHandler`構造体:
  - 変更前: `emailService *email.EmailService`, `templateService *email.TemplateService`
  - 変更後: `emailUsecase *usecase.EmailUsecase`
- `NewEmailHandler`関数:
  - 変更前: `func NewEmailHandler(emailService *email.EmailService, templateService *email.TemplateService) *EmailHandler`
  - 変更後: `func NewEmailHandler(emailUsecase *usecase.EmailUsecase) *EmailHandler`
- `RegisterEmailEndpoints`内の処理:
  - テンプレートレンダリング、件名取得、メール送信の処理をusecase層に委譲

**受け入れ基準**:
- `EmailHandler`構造体のフィールドが`emailUsecase`に変更されている
- `NewEmailHandler`関数がusecase層を受け取るように修正されている
- `RegisterEmailEndpoints`内の処理がusecase層を呼び出すように修正されている
- 認証チェック、バリデーション、エラーハンドリングがhandler層で継続して行われている

- _Requirements: 3.3.4, 6.3_
- _Design: 5.2.2_

---

#### タスク 5.3: main.goの初期化処理の修正（email API用）
**目的**: `main.go`でemail API用のusecase層を初期化し、handler層に注入する。

**作業内容**:
- `server/cmd/server/main.go` を修正
- `EmailUsecase`の初期化を追加
- `EmailHandler`の初期化を修正（usecase層を受け取るように変更）

**実装内容**:
- `EmailUsecase`の初期化:
  ```go
  emailUsecase := usecase.NewEmailUsecase(emailService, templateService)
  ```
- `EmailHandler`の初期化:
  ```go
  // 変更前: emailHandler := handler.NewEmailHandler(emailService, templateService)
  // 変更後:
  emailHandler := handler.NewEmailHandler(emailUsecase)
  ```

**受け入れ基準**:
- `EmailUsecase`の初期化が追加されている
- `EmailHandler`の初期化がusecase層を受け取るように修正されている
- 初期化の順序が正しい（service → usecase → handler）

- _Requirements: 3.4.2, 6.4_
- _Design: 6.2.1_

---

#### タスク 5.4: EmailUsecaseのテスト作成
**目的**: `EmailUsecase`の単体テストを作成する。

**作業内容**:
- `server/internal/usecase/email_usecase_test.go` を作成
- `EmailService`と`TemplateService`をモック化
- `SendEmail`メソッドのテストを実装

**受け入れ基準**:
- `server/internal/usecase/email_usecase_test.go` が作成されている
- `EmailService`と`TemplateService`がモック化されている
- `SendEmail`メソッドのテストが実装されている
- テストが正常に実行できる

- _Requirements: 6.5_
- _Design: 8.2.1_

---

#### タスク 5.5: EmailHandlerのテスト修正
**目的**: `EmailHandler`のテストをusecase層を使用するように修正する。

**作業内容**:
- `server/internal/api/handler/email_handler_test.go` を修正
- `EmailUsecase`をモック化
- テストケースを修正

**受け入れ基準**:
- `EmailHandler`のテストがusecase層を使用するように修正されている
- `EmailUsecase`がモック化されている
- 既存のテストが正常に動作する
- テストが正常に実行できる

- _Requirements: 6.5_
- _Design: 8.2.2_

---

#### タスク 5.6: email APIの動作確認
**目的**: email APIが正常に動作することを確認する。

**作業内容**:
- サーバーを起動
- email APIを実行
- レスポンスを確認

**実装内容**:
- APIエンドポイント: `POST /api/email/send`

**受け入れ基準**:
- email APIが正常に動作する
- レスポンスが正しい形式である
- 既存のAPIの動作が変わらない

- _Requirements: 6.6_
- _Design: 9.2_

---

### Phase 6: dm_jobqueue APIの実装

#### タスク 6.1: DmJobqueueUsecaseの作成
**目的**: dm_jobqueue API用のビジネスロジックを担当するusecase層を作成する。

**作業内容**:
- `server/internal/usecase/dm_jobqueue_usecase.go` を作成
- `DmJobqueueUsecase`構造体を定義
- `RegisterJob`メソッドを実装

**実装内容**:
- パッケージ名: `usecase`
- 構造体名: `DmJobqueueUsecase`
- 依存関係: `*jobqueue.Client`
- コンストラクタ: `NewDmJobqueueUsecase(jobQueueClient *jobqueue.Client) *DmJobqueueUsecase`
- メソッド: `RegisterJob(ctx context.Context, message string, delaySeconds int, maxRetry int) (string, error)`
  - メッセージのデフォルト値設定
  - ペイロードの作成
  - ジョブオプションの作成
  - ジョブの登録

**受け入れ基準**:
- `server/internal/usecase/dm_jobqueue_usecase.go` が作成されている
- `DmJobqueueUsecase`構造体が定義されている
- `RegisterJob`メソッドが実装されている
- `RegisterJob`メソッドがビジネスロジック（メッセージのデフォルト値設定など）を含んでいる

- _Requirements: 3.3.2, 6.1_
- _Design: 3.2.5_

---

#### タスク 6.2: DmJobqueueHandlerの修正
**目的**: `DmJobqueueHandler`をusecase層を呼び出すように修正する。

**作業内容**:
- `server/internal/api/handler/dm_jobqueue_handler.go` を修正
- `DmJobqueueHandler`構造体のフィールドを`jobQueueClient`から`dmJobqueueUsecase`に変更
- `NewDmJobqueueHandler`関数を修正してusecase層を受け取るように変更
- `RegisterJob`メソッドを修正してusecase層を呼び出すように変更

**実装内容**:
- `DmJobqueueHandler`構造体:
  - 変更前: `jobQueueClient *jobqueue.Client`
  - 変更後: `dmJobqueueUsecase *usecase.DmJobqueueUsecase`
- `NewDmJobqueueHandler`関数:
  - 変更前: `func NewDmJobqueueHandler(jobQueueClient *jobqueue.Client) *DmJobqueueHandler`
  - 変更後: `func NewDmJobqueueHandler(dmJobqueueUsecase *usecase.DmJobqueueUsecase) *DmJobqueueHandler`
- `RegisterJob`メソッド:
  - ビジネスロジックをusecase層に委譲

**受け入れ基準**:
- `DmJobqueueHandler`構造体のフィールドが`dmJobqueueUsecase`に変更されている
- `NewDmJobqueueHandler`関数がusecase層を受け取るように修正されている
- `RegisterJob`メソッドがusecase層を呼び出すように修正されている
- エラーハンドリングがhandler層で継続して行われている

- _Requirements: 3.3.4, 6.3_
- _Design: 5.2.2_

---

#### タスク 6.3: main.goの初期化処理の修正（dm_jobqueue API用）
**目的**: `main.go`でdm_jobqueue API用のusecase層を初期化し、handler層に注入する。

**作業内容**:
- `server/cmd/server/main.go` を修正
- `DmJobqueueUsecase`の初期化を追加
- `DmJobqueueHandler`の初期化を修正（usecase層を受け取るように変更）

**実装内容**:
- `DmJobqueueUsecase`の初期化:
  ```go
  dmJobqueueUsecase := usecase.NewDmJobqueueUsecase(jobQueueClient)
  ```
- `DmJobqueueHandler`の初期化:
  ```go
  // 変更前: dmJobqueueHandler := handler.NewDmJobqueueHandler(jobQueueClient)
  // 変更後:
  dmJobqueueHandler := handler.NewDmJobqueueHandler(dmJobqueueUsecase)
  ```

**受け入れ基準**:
- `DmJobqueueUsecase`の初期化が追加されている
- `DmJobqueueHandler`の初期化がusecase層を受け取るように修正されている
- 初期化の順序が正しい（service → usecase → handler）

- _Requirements: 3.4.2, 6.4_
- _Design: 6.2.1_

---

#### タスク 6.4: DmJobqueueUsecaseのテスト作成
**目的**: `DmJobqueueUsecase`の単体テストを作成する。

**作業内容**:
- `server/internal/usecase/dm_jobqueue_usecase_test.go` を作成
- `jobqueue.Client`をモック化
- `RegisterJob`メソッドのテストを実装

**受け入れ基準**:
- `server/internal/usecase/dm_jobqueue_usecase_test.go` が作成されている
- `jobqueue.Client`がモック化されている
- `RegisterJob`メソッドのテストが実装されている
- テストが正常に実行できる

- _Requirements: 6.5_
- _Design: 8.2.1_

---

#### タスク 6.5: DmJobqueueHandlerのテスト修正
**目的**: `DmJobqueueHandler`のテストをusecase層を使用するように修正する。

**作業内容**:
- `server/internal/api/handler/dm_jobqueue_handler_test.go` を修正
- `DmJobqueueUsecase`をモック化
- テストケースを修正

**受け入れ基準**:
- `DmJobqueueHandler`のテストがusecase層を使用するように修正されている
- `DmJobqueueUsecase`がモック化されている
- 既存のテストが正常に動作する
- テストが正常に実行できる

- _Requirements: 6.5_
- _Design: 8.2.2_

---

#### タスク 6.6: dm_jobqueue APIの動作確認
**目的**: dm_jobqueue APIが正常に動作することを確認する。

**作業内容**:
- サーバーを起動
- dm_jobqueue APIを実行
- レスポンスを確認

**実装内容**:
- APIエンドポイント: `POST /api/dm-jobqueue/register`

**受け入れ基準**:
- dm_jobqueue APIが正常に動作する
- レスポンスが正しい形式である
- 既存のAPIの動作が変わらない

- _Requirements: 6.6_
- _Design: 9.2_

---

### Phase 7: 最終確認

#### タスク 7.1: すべてのAPIエンドポイントの動作確認
**目的**: すべてのAPIエンドポイントが正常に動作することを確認する。

**作業内容**:
- サーバーを起動
- すべてのAPIエンドポイントを実行
- レスポンスを確認

**実装内容**:
- すべてのAPIエンドポイント:
  - today API
  - dm_user API
  - dm_post API
  - email API
  - dm_jobqueue API
  - upload API（バリデーションと入出力の制御のみのため、usecase層は不要）

**受け入れ基準**:
- すべてのAPIエンドポイントが正常に動作する
- すべてのAPIエンドポイントのレスポンスが正しい形式である
- 既存のAPIの動作が変わらない

- _Requirements: 6.6_

---

#### タスク 7.2: 既存テストの動作確認
**目的**: 既存のテストがすべて正常に動作することを確認する。

**作業内容**:
- すべてのテストを実行
- テスト結果を確認

**実装内容**:
- すべてのテストファイルを実行
- テスト結果を確認

**受け入れ基準**:
- すべての既存テストが正常に動作する
- テストが全て失敗しないことを確認

- _Requirements: 6.5_

---

#### タスク 7.3: レイヤー構造の確認
**目的**: レイヤー構造が適切に実装されていることを確認する。

**作業内容**:
- 各レイヤーの実装を確認
- 依存関係を確認

**実装内容**:
- handler層がusecase層を呼び出していることを確認
- usecase層がservice層を呼び出していることを確認
- service層がrepository層を呼び出していることを確認
- usecase層がrepository層を直接呼び出していないことを確認

**受け入れ基準**:
- すべてのhandler層がusecase層を呼び出している
- すべてのusecase層がservice層を呼び出している
- usecase層がrepository層を直接呼び出していない
- レイヤー構造が適切に実装されている

- _Requirements: 6.6_

---

## 実装順序のまとめ

1. **Phase 1**: 基盤の準備（ディレクトリ作成、DateService作成）
2. **Phase 2**: today APIの実装（パターン確立）
3. **Phase 3**: dm_user APIの実装
4. **Phase 4**: dm_post APIの実装
5. **Phase 5**: email APIの実装
6. **Phase 6**: dm_jobqueue APIの実装
7. **Phase 7**: 最終確認

**注意**: upload APIはバリデーションと入出力の制御のみを行っており、ビジネスロジックがないため、usecase層の導入は不要です。

## 注意事項

- 各フェーズごとに動作確認を行い、問題がないことを確認してから次のフェーズに進む
- 既存のテストが正常に動作することを各フェーズで確認する
- 既存のAPIの動作が変わらないことを各フェーズで確認する
- 実装時は設計書の内容に従うこと
