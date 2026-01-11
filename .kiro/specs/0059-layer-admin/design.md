# Adminアプリのカスタムページの実装の仕組みを変更するの設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、Adminアプリのカスタムページ（`dm_user_register.go`、`api_key.go`）の実装を、APIサーバーと同じレイヤー構造（pages -> usecase -> service -> repository -> db）に変更するための詳細設計を定義する。これにより、AdminアプリとAPIサーバーで一貫したアーキテクチャを実現し、コードの保守性と再利用性を向上させる。

### 1.2 設計の範囲
- Admin用usecase層（`server/internal/usecase/admin`）の設計
- Service層の拡張（`server/internal/service/api_key_service.go`の新規作成）の設計
- pages層の修正（`server/internal/admin/pages/dm_user_register.go`、`server/internal/admin/pages/api_key.go`）の設計
- 依存関係の注入設計（main.goでの初期化方法）
- テスト設計
- ドキュメント更新の設計

### 1.3 設計方針
- **一貫性**: APIサーバーと同じレイヤー構造を採用
- **既存コードの活用**: 既存のservice層（`DmUserService`）をそのまま使用
- **責務の明確化**: 各レイヤーの責務を明確に分離（pages層はバリデーションと入出力制御のみ）
- **テスト容易性**: 各層を独立してテストできる設計
- **後方互換性**: 既存のAdminページの動作（出力形式、エラーメッセージ）を維持
- **処理の分離**: api_key.goで鍵の生成とペイロードのデコードを2つの処理に分離

## 2. アーキテクチャ設計

### 2.1 全体構成

#### 2.1.1 dm_user_register.goのアーキテクチャ

```
┌─────────────────────────────────────────────────────────────┐
│      Pages Layer (internal/admin/pages/dm_user_register.go)  │
│  • エントリーポイント                                         │
│  • バリデーション（validateDmUserInput）                      │
│  • 入出力制御（renderDmUserRegisterForm）                     │
│  • usecase層の呼び出し                                       │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│   Usecase Layer (internal/usecase/admin)                     │
│  • DmUserRegisterUsecase                                     │
│  • RegisterDmUser()                                          │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│         Service Layer (internal/service)                        │
│  • DmUserService（既存）                                     │
│  • CreateDmUser()                                            │
│  • メールアドレスの重複チェック（既存のロジックを使用）      │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│      Repository Layer (internal/repository)                   │
│  • DmUserRepository（既存）                                  │
│  • Create()                                                  │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│         DB Layer (internal/db)                                 │
│  • GroupManager                                              │
│  • Sharding接続管理                                          │
└─────────────────────────────────────────────────────────────┘
```

#### 2.1.2 api_key.goのアーキテクチャ

```
┌─────────────────────────────────────────────────────────────┐
│         Pages Layer (internal/admin/pages/api_key.go)         │
│  • エントリーポイント                                         │
│  • 入出力制御（renderAPIKeyPage、renderAPIKeyResult）          │
│  • usecase層の呼び出し（2つの処理に分離）                     │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│   Usecase Layer (internal/usecase/admin)                     │
│  • APIKeyUsecase                                             │
│  • GenerateAPIKey()（鍵の生成）                               │
│  • DecodeAPIKeyPayload()（ペイロードのデコード）             │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│         Service Layer (internal/service)                        │
│  • APIKeyService（新規）                                      │
│  • GenerateAPIKey()（鍵の生成）                               │
│  • DecodeAPIKeyPayload()（ペイロードのデコード）             │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│         Auth Layer (internal/auth)                             │
│  • GeneratePublicAPIKey()                                     │
│  • ParseJWTClaims()                                           │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 データフロー

#### 2.2.1 dm_user_register.goのデータフロー

```
main.go（server/cmd/admin/main.go）
  ↓
GroupManagerの初期化（db.NewGroupManager(cfg)）
  ↓
Repository層の初期化
  - repository.NewDmUserRepository(groupManager)
  ↓
Service層の初期化
  - service.NewDmUserService(dmUserRepository)
  ↓
Usecase層の初期化
  - admin.NewDmUserRegisterUsecase(dmUserService)
  ↓
pages.DmUserRegisterPage(ctx, groupManager)
  ↓
handleDmUserRegisterPost(ctx, groupManager)
  ↓
バリデーション（validateDmUserInput）
  ↓
usecase.RegisterDmUser(ctx, name, email)
  ↓
service.CheckEmailExists(ctx, email)（メールアドレスの重複チェック）
  ↓
repository.CheckEmailExists(ctx, email)（全シャード検索）
  ↓
service.CreateDmUser(ctx, req)
  ↓
repository.Create(ctx, req)
  ↓
DB操作（GroupManager経由）
  ↓
renderDmUserRegisterForm（成功時はリダイレクト）
```

#### 2.2.2 api_key.goのデータフロー

```
main.go（server/cmd/admin/main.go）
  ↓
Service層の初期化
  - service.NewAPIKeyService()
  ↓
Usecase層の初期化
  - admin.NewAPIKeyUsecase(apiKeyService)
  ↓
pages.APIKeyPage(ctx, conn)
  ↓
handleGenerateKey(ctx, cfg)
  ↓
usecase.GenerateAPIKey(ctx, env)（鍵の生成）
  ↓
service.GenerateAPIKey(ctx, secretKey, version, env, issuedAt)
  ↓
auth.GeneratePublicAPIKey(secretKey, version, env, issuedAt)
  ↓
usecase.DecodeAPIKeyPayload(ctx, token)（ペイロードのデコード）
  ↓
service.DecodeAPIKeyPayload(ctx, token)
  ↓
auth.ParseJWTClaims(token)
  ↓
renderAPIKeyResult（結果の表示）
```

### 2.3 レイヤー構造の比較

#### 修正前（dm_user_register.go）

```
pages.DmUserRegisterPage()
  ↓ (直接実装)
validateDmUserInput()
  ↓
checkEmailExistsSharded()（全シャード検索）
  ↓
insertDmUserSharded()（DB操作）
  ↓
renderDmUserRegisterForm()
```

#### 修正後（dm_user_register.go）

```
pages.DmUserRegisterPage()
  ↓
validateDmUserInput()（pages層でバリデーション）
  ↓
usecase.RegisterDmUser()
  ↓
service.CheckEmailExists()（メールアドレスの重複チェック）
  ↓
repository.CheckEmailExists()（全シャード検索）
  ↓
service.CreateDmUser()
  ↓
repository.Create()
  ↓
DB操作（GroupManager経由）
  ↓
renderDmUserRegisterForm()（pages層で入出力制御）
```

#### 修正前（api_key.go）

```
pages.APIKeyPage()
  ↓ (直接実装)
generatePublicAPIKey()（鍵の生成とペイロードのデコードが1つの処理）
  ↓
auth.GeneratePublicAPIKey()
  ↓
auth.ParseJWTClaims()
  ↓
renderAPIKeyResult()
```

#### 修正後（api_key.go）

```
pages.APIKeyPage()
  ↓
usecase.GenerateAPIKey()（鍵の生成）
  ↓
service.GenerateAPIKey()
  ↓
auth.GeneratePublicAPIKey()
  ↓
usecase.DecodeAPIKeyPayload()（ペイロードのデコード）
  ↓
service.DecodeAPIKeyPayload()
  ↓
auth.ParseJWTClaims()
  ↓
renderAPIKeyResult()（pages層で入出力制御）
```

## 3. 詳細設計

### 3.1 Admin用usecase層の設計

#### 3.1.1 `dm_user_register_usecase.go`の設計

**ファイルパス**: `server/internal/usecase/admin/dm_user_register_usecase.go`

**実装内容**:

```go
package admin

import (
	"context"

	"github.com/taku-o/go-webdb-template/internal/model"
	"github.com/taku-o/go-webdb-template/internal/usecase"
)

// DmUserRegisterUsecase はdm_user登録のビジネスロジックを担当
type DmUserRegisterUsecase struct {
	dmUserService usecase.DmUserServiceInterface
}

// NewDmUserRegisterUsecase は新しいDmUserRegisterUsecaseを作成
func NewDmUserRegisterUsecase(dmUserService usecase.DmUserServiceInterface) *DmUserRegisterUsecase {
	return &DmUserRegisterUsecase{
		dmUserService: dmUserService,
	}
}

// RegisterDmUser はユーザーを登録
func (u *DmUserRegisterUsecase) RegisterDmUser(ctx context.Context, name, email string) (string, error) {
	// メールアドレスの重複チェック
	exists, err := u.dmUserService.CheckEmailExists(ctx, email)
	if err != nil {
		return "", fmt.Errorf("failed to check email: %w", err)
	}
	if exists {
		return "", fmt.Errorf("this email address is already registered")
	}

	req := &model.CreateDmUserRequest{
		Name:  name,
		Email: email,
	}

	dmUser, err := u.dmUserService.CreateDmUser(ctx, req)
	if err != nil {
		return "", err
	}

	return dmUser.ID, nil
}
```

**設計のポイント**:
- 既存の`DmUserServiceInterface`を使用（`server/internal/usecase/dm_user_usecase.go`で定義）
- `RegisterDmUser`メソッドでメールアドレスの重複チェックを実施（`CheckEmailExists`を呼び出し）
- 重複チェック後に`CreateDmUserRequest`を作成してservice層を呼び出し
- エラーハンドリングはservice層から返されたエラーをそのまま返す
- `context.Context`を受け取る

#### 3.1.2 `api_key_usecase.go`の設計

**ファイルパス**: `server/internal/usecase/admin/api_key_usecase.go`

**実装内容**:

```go
package admin

import (
	"context"
	"os"
	"time"

	"github.com/taku-o/go-webdb-template/internal/auth"
	"github.com/taku-o/go-webdb-template/internal/config"
	"github.com/taku-o/go-webdb-template/internal/service"
)

// APIKeyUsecase はAPIキー発行のビジネスロジックを担当
type APIKeyUsecase struct {
	apiKeyService service.APIKeyServiceInterface
	cfg           *config.Config
}

// NewAPIKeyUsecase は新しいAPIKeyUsecaseを作成
func NewAPIKeyUsecase(apiKeyService service.APIKeyServiceInterface, cfg *config.Config) *APIKeyUsecase {
	return &APIKeyUsecase{
		apiKeyService: apiKeyService,
		cfg:           cfg,
	}
}

// GenerateAPIKey はAPIキーを生成
func (u *APIKeyUsecase) GenerateAPIKey(ctx context.Context, env string) (string, error) {
	if env == "" {
		env = "develop"
	}

	now := time.Now()
	token, err := u.apiKeyService.GenerateAPIKey(ctx, u.cfg.API.SecretKey, u.cfg.API.CurrentVersion, env, now.Unix())
	if err != nil {
		return "", err
	}

	return token, nil
}

// DecodeAPIKeyPayload はAPIキーのペイロードをデコード
func (u *APIKeyUsecase) DecodeAPIKeyPayload(ctx context.Context, token string) (*auth.JWTClaims, error) {
	claims, err := u.apiKeyService.DecodeAPIKeyPayload(ctx, token)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
```

**設計のポイント**:
- `APIKeyServiceInterface`を依存として注入（新規作成）
- `GenerateAPIKey`と`DecodeAPIKeyPayload`を2つの処理に分離
- `config.Config`を依存として注入（SecretKey、CurrentVersionを取得するため）
- エラーハンドリングはservice層から返されたエラーをそのまま返す
- `context.Context`を受け取る

#### 3.1.3 Admin用usecase層のテスト設計

**テストファイル**:
- `server/internal/usecase/admin/dm_user_register_usecase_test.go`
- `server/internal/usecase/admin/api_key_usecase_test.go`

**テスト内容**:
- 正常系テスト
- エラーハンドリングのテスト
- service層のモックを使用したテスト

### 3.2 Service層の拡張設計

#### 3.2.1 既存のservice層の拡張

**既存のservice層**:
- `server/internal/service/dm_user_service.go`: 既存の`DmUserService`を使用
- `CreateDmUser`メソッドが既に実装されている
- **メールアドレスの重複チェック機能を追加する必要がある**

**追加する機能**:

1. **`DmUserServiceInterface`へのメソッド追加**:
   - `server/internal/usecase/dm_user_usecase.go`の`DmUserServiceInterface`に`CheckEmailExists(ctx context.Context, email string) (bool, error)`メソッドを追加

2. **`DmUserService`へのメソッド追加**:
   - `server/internal/service/dm_user_service.go`に`CheckEmailExists`メソッドを追加
   - repository層の`CheckEmailExists`を呼び出し

3. **`DmUserRepositoryInterface`へのメソッド追加**:
   - `server/internal/repository/interfaces.go`の`DmUserRepositoryInterface`に`CheckEmailExists(ctx context.Context, email string) (bool, error)`メソッドを追加

4. **`DmUserRepository`へのメソッド追加**:
   - `server/internal/repository/dm_user_repository.go`に`CheckEmailExists(ctx context.Context, email string) (bool, error)`メソッドを追加
   - 全シャードを検索してメールアドレスの重複をチェック
   - 既存の`checkEmailExistsSharded`関数のロジックを移行

**実装例**:

```go
// server/internal/repository/dm_user_repository.go に追加
// CheckEmailExists はメールアドレスが既に存在するかチェックする（全シャード検索）
func (r *DmUserRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	// 全テーブルを検索
	tableCount := r.tableSelector.GetTableCount()
	for tableNum := 0; tableNum < tableCount; tableNum++ {
		conn, err := r.groupManager.GetShardingConnection(tableNum)
		if err != nil {
			return false, fmt.Errorf("failed to get connection for table %d: %w", tableNum, err)
		}

		tableName := fmt.Sprintf("dm_users_%03d", tableNum)
		var count int64
		err = conn.DB.WithContext(ctx).Table(tableName).Where("email = ?", email).Count(&count).Error
		if err != nil {
			return false, fmt.Errorf("failed to check email in %s: %w", tableName, err)
		}

		if count > 0 {
			return true, nil
		}
	}

	return false, nil
}
```

```go
// server/internal/service/dm_user_service.go に追加
// CheckEmailExists はメールアドレスが既に存在するかチェックする
func (s *DmUserService) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	return s.dmUserRepo.CheckEmailExists(ctx, email)
}
```

**設計のポイント**:
- 既存の`DmUserService`を拡張（メソッドを追加）
- `DmUserServiceInterface`に`CheckEmailExists`メソッドを追加（`server/internal/usecase/dm_user_usecase.go`）
- `DmUserRepositoryInterface`に`CheckEmailExists`メソッドを追加（`server/internal/repository/interfaces.go`）
- `DmUserRepository`に`CheckEmailExists`メソッドを追加（全シャード検索用）
- 既存の`checkEmailExistsSharded`関数のロジックを移行

#### 3.2.2 `api_key_service.go`の設計

**ファイルパス**: `server/internal/service/api_key_service.go`

**実装内容**:

```go
package service

import (
	"context"

	"github.com/taku-o/go-webdb-template/internal/auth"
)

// APIKeyServiceInterface はAPIキーサービスのインターフェース
type APIKeyServiceInterface interface {
	GenerateAPIKey(ctx context.Context, secretKey, version, env string, issuedAt int64) (string, error)
	DecodeAPIKeyPayload(ctx context.Context, token string) (*auth.JWTClaims, error)
}

// APIKeyService はAPIキー発行のドメインロジックを担当
type APIKeyService struct{}

// NewAPIKeyService は新しいAPIKeyServiceを作成
func NewAPIKeyService() *APIKeyService {
	return &APIKeyService{}
}

// GenerateAPIKey はAPIキーを生成
func (s *APIKeyService) GenerateAPIKey(ctx context.Context, secretKey, version, env string, issuedAt int64) (string, error) {
	token, err := auth.GeneratePublicAPIKey(secretKey, version, env, issuedAt)
	if err != nil {
		return "", err
	}

	return token, nil
}

// DecodeAPIKeyPayload はAPIキーのペイロードをデコード
func (s *APIKeyService) DecodeAPIKeyPayload(ctx context.Context, token string) (*auth.JWTClaims, error) {
	claims, err := auth.ParseJWTClaims(token)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
```

**設計のポイント**:
- `APIKeyServiceInterface`を新規作成
- `GenerateAPIKey`と`DecodeAPIKeyPayload`を2つの処理に分離
- `auth.GeneratePublicAPIKey`、`auth.ParseJWTClaims`を呼び出し
- エラーハンドリングはauth層から返されたエラーをそのまま返す
- `context.Context`を受け取る（将来の拡張に備える）

#### 3.2.3 Service層のテスト設計

**テストファイル**:
- `server/internal/service/api_key_service_test.go`

**テスト内容**:
- 正常系テスト
- エラーハンドリングのテスト
- auth層のモックを使用したテスト（必要に応じて）

### 3.3 pages層の修正設計

#### 3.3.1 `dm_user_register.go`の修正設計

**ファイルパス**: `server/internal/admin/pages/dm_user_register.go`

**修正内容**:

```go
package pages

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/template/types"
	appdb "github.com/taku-o/go-webdb-template/internal/db"
	"github.com/taku-o/go-webdb-template/internal/usecase/admin"
)

// DmUserRegisterPage はユーザー登録ページを返す
func DmUserRegisterPage(ctx *context.Context, groupManager *appdb.GroupManager, dmUserRegisterUsecase *admin.DmUserRegisterUsecase) (types.Panel, error) {
	if ctx.Method() == http.MethodPost {
		return handleDmUserRegisterPost(ctx, dmUserRegisterUsecase)
	}
	return renderDmUserRegisterForm(ctx, "", "", nil)
}

// handleDmUserRegisterPost はPOSTリクエストを処理する
func handleDmUserRegisterPost(ctx *context.Context, dmUserRegisterUsecase *admin.DmUserRegisterUsecase) (types.Panel, error) {
	name := strings.TrimSpace(ctx.FormValue("name"))
	email := strings.TrimSpace(ctx.FormValue("email"))

	// バリデーション
	errors := validateDmUserInput(name, email)
	if len(errors) > 0 {
		return renderDmUserRegisterForm(ctx, name, email, errors)
	}

	// usecase層を呼び出し
	dmUserID, err := dmUserRegisterUsecase.RegisterDmUser(ctx.Request.Context(), name, email)
	if err != nil {
		return renderDmUserRegisterForm(ctx, name, email, []string{"ユーザー登録に失敗しました: " + err.Error()})
	}

	// 登録完了ページへリダイレクト（クエリパラメータで情報を渡す）
	redirectURL := fmt.Sprintf("/admin/dm-user/register/new?id=%s&name=%s&email=%s",
		url.QueryEscape(dmUserID),
		url.QueryEscape(name),
		url.QueryEscape(email),
	)

	// GoAdminのContent wrapperはctx.Redirectを上書きするため、
	// JavaScriptリダイレクトを使用
	return types.Panel{
		Title:       "リダイレクト中",
		Description: "",
		Content:     template.HTML(fmt.Sprintf(`<script>window.location.href='%s';</script>`, redirectURL)),
	}, nil
}

// validateDmUserInput は入力値をバリデーションする（維持）
func validateDmUserInput(name, email string) []string {
	var errors []string

	if name == "" {
		errors = append(errors, "名前は必須です")
	} else if len(name) > 100 {
		errors = append(errors, "名前は100文字以内で入力してください")
	}

	if email == "" {
		errors = append(errors, "メールアドレスは必須です")
	} else if !strings.Contains(email, "@") {
		errors = append(errors, "有効なメールアドレスを入力してください")
	} else if len(email) > 255 {
		errors = append(errors, "メールアドレスは255文字以内で入力してください")
	}

	return errors
}

// renderDmUserRegisterForm はユーザー登録フォームをレンダリングする（維持）
func renderDmUserRegisterForm(ctx *context.Context, name, email string, errors []string) (types.Panel, error) {
	// 既存の実装を維持
	// ...
}
```

**設計のポイント**:
- `DmUserRegisterPage`関数のシグネチャを変更（`dmUserRegisterUsecase`を追加）
- `handleDmUserRegisterPost`関数を修正（usecase層を呼び出すように変更）
- `validateDmUserInput`関数を維持（pages層でバリデーション）
- `renderDmUserRegisterForm`関数を維持（pages層で入出力制御）
- `checkEmailExistsSharded`関数を削除（service層に移動済み）
- `insertDmUserSharded`関数を削除（service層に移動済み）

#### 3.3.2 `api_key.go`の修正設計

**ファイルパス**: `server/internal/admin/pages/api_key.go`

**修正内容**:

```go
package pages

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/template/types"

	"github.com/taku-o/go-webdb-template/internal/auth"
	"github.com/taku-o/go-webdb-template/internal/config"
	"github.com/taku-o/go-webdb-template/internal/usecase/admin"
)

// APIKeyPage はAPIキー発行ページを返す
// 注意: RegisterCustomPagesで"/api-key"と登録すると、実際のURLは"/admin/api-key"になる
// HTML内のリンクも"/admin/api-key"とする必要がある
func APIKeyPage(ctx *context.Context, conn db.Connection, apiKeyUsecase *admin.APIKeyUsecase) (types.Panel, error) {
	// 設定を取得
	cfg, err := config.Load()
	if err != nil {
		return types.Panel{}, err
	}

	// POSTリクエスト: キー生成
	if ctx.Method() == http.MethodPost {
		return handleGenerateKey(ctx, cfg, apiKeyUsecase)
	}

	// GETリクエスト: フォーム表示
	return renderAPIKeyPage(ctx, cfg)
}

// handleGenerateKey はAPIキーを生成
func handleGenerateKey(ctx *context.Context, cfg *config.Config, apiKeyUsecase *admin.APIKeyUsecase) (types.Panel, error) {
	// 現在の環境を取得
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "develop"
	}

	// 鍵の生成（usecase層を呼び出し）
	token, err := apiKeyUsecase.GenerateAPIKey(ctx.Request.Context(), env)
	if err != nil {
		return types.Panel{}, err
	}

	// ペイロードのデコード（usecase層を呼び出し）
	claims, err := apiKeyUsecase.DecodeAPIKeyPayload(ctx.Request.Context(), token)
	if err != nil {
		return types.Panel{}, err
	}

	// 生成結果を表示
	return renderAPIKeyResult(ctx, token, claims)
}

// renderAPIKeyPage はAPIキー発行ページをレンダリング（維持）
func renderAPIKeyPage(ctx *context.Context, cfg *config.Config) (types.Panel, error) {
	// 既存の実装を維持
	// ...
}

// renderAPIKeyResult は生成結果をレンダリング（維持）
func renderAPIKeyResult(ctx *context.Context, token string, claims *auth.JWTClaims) (types.Panel, error) {
	// 既存の実装を維持
	// ...
}
```

**設計のポイント**:
- `APIKeyPage`関数のシグネチャを変更（`apiKeyUsecase`を追加）
- `handleGenerateKey`関数を修正（usecase層を呼び出すように変更、2つの処理に分離）
- `renderAPIKeyPage`関数を維持（pages層で入出力制御）
- `renderAPIKeyResult`関数を維持（pages層で入出力制御）
- `generatePublicAPIKey`関数を削除（service層に移動済み）

### 3.4 依存関係の注入設計

#### 3.4.1 main.goでの初期化設計

**ファイルパス**: `server/cmd/admin/main.go`

**修正内容**:

```go
package main

import (
	// ... 既存のimport ...

	"github.com/taku-o/go-webdb-template/internal/repository"
	"github.com/taku-o/go-webdb-template/internal/service"
	"github.com/taku-o/go-webdb-template/internal/usecase/admin"
)

func main() {
	// ... 既存の初期化処理 ...

	// Repository層の初期化
	dmUserRepository := repository.NewDmUserRepository(groupManager)

	// Service層の初期化
	dmUserService := service.NewDmUserService(dmUserRepository)
	apiKeyService := service.NewAPIKeyService()

	// Usecase層の初期化
	dmUserRegisterUsecase := admin.NewDmUserRegisterUsecase(dmUserService)
	apiKeyUsecase := admin.NewAPIKeyUsecase(apiKeyService, cfg)

	// カスタムページの登録（Gorilla Mux用にContent関数を使用）
	app.HandleFunc("/admin/dm-user/register", gorillaAdapter.Content(func(ctx gorillaAdapter.Context) (types.Panel, error) {
		return pages.DmUserRegisterPage(goadminContext.NewContext(ctx.Request), groupManager, dmUserRegisterUsecase)
	})).Methods("GET", "POST")
	app.HandleFunc("/admin/api-key", gorillaAdapter.Content(func(ctx gorillaAdapter.Context) (types.Panel, error) {
		return pages.APIKeyPage(goadminContext.NewContext(ctx.Request), conn, apiKeyUsecase)
	})).Methods("GET", "POST")

	// ... 既存の処理 ...
}
```

**設計のポイント**:
- Repository層の初期化（`DmUserRepository`）
- Service層の初期化（`DmUserService`、`APIKeyService`）
- Usecase層の初期化（`DmUserRegisterUsecase`、`APIKeyUsecase`）
- pages層の関数呼び出し時にusecase層を渡す
- `groupManager`は既存のまま使用（`DmUserRegisterPage`で必要）

### 3.5 テスト設計

#### 3.5.1 usecase層のテスト設計

**テストファイル**:
- `server/internal/usecase/admin/dm_user_register_usecase_test.go`
- `server/internal/usecase/admin/api_key_usecase_test.go`

**テスト内容**:
- 正常系テスト
- service層のモックを使用したテスト
- エラーハンドリングのテスト

#### 3.5.2 service層のテスト設計

**テストファイル**:
- `server/internal/service/api_key_service_test.go`

**テスト内容**:
- 正常系テスト
- auth層のモックを使用したテスト（必要に応じて）
- エラーハンドリングのテスト

#### 3.5.3 pages層のテスト設計

**テスト内容**:
- 既存のテストを維持（必要に応じて修正）
- usecase層のモックを使用したテスト（必要に応じて）

### 3.6 ドキュメント更新の設計

#### 3.6.1 アーキテクチャドキュメントの更新

**ファイルパス**: `docs/Architecture.md`

**更新内容**:
- Adminアプリのレイヤー構造を追加（pages -> usecase -> service -> repository -> db）
- Adminアプリのアーキテクチャ図を更新
- Admin用usecase層の説明を追加

#### 3.6.2 プロジェクト構造ドキュメントの更新

**ファイルパス**: `docs/Project-Structure.md`

**更新内容**:
- `server/internal/usecase/admin`ディレクトリを追加
- `server/internal/usecase/admin/dm_user_register_usecase.go`を追加
- `server/internal/usecase/admin/api_key_usecase.go`を追加
- `server/internal/service/api_key_service.go`を追加

#### 3.6.3 ファイル組織ドキュメントの更新

**ファイルパス**: `.kiro/steering/structure.md`

**更新内容**:
- `server/internal/usecase/admin`ディレクトリを追加
- `server/internal/usecase/admin/dm_user_register_usecase.go`を追加
- `server/internal/usecase/admin/api_key_usecase.go`を追加
- `server/internal/service/api_key_service.go`を追加

## 4. 実装上の注意事項

### 4.1 Admin用usecase層の実装

- **インターフェースの使用**: 既存の`DmUserServiceInterface`を使用。`APIKeyServiceInterface`を新規作成
- **依存関係の注入**: コンストラクタでservice層のインターフェースを注入
- **エラーハンドリング**: service層から返されたエラーをそのまま返す（エラーのラップは不要）
- **処理の分離**: api_key.goで鍵の生成とペイロードのデコードを2つの処理に分ける

### 4.2 Service層の実装

- **既存のservice層の使用**: `DmUserService`は既存のものを使用（新規作成は不要）
- **新規service層の作成**: `APIKeyService`を新規作成
- **依存関係の注入**: コンストラクタで依存関係を注入（必要に応じて）
- **エラーハンドリング**: 下位層から返されたエラーをそのまま返す（エラーのラップは不要）
- **鍵の生成とペイロードのデコード**: 2つの処理に分ける（`GenerateAPIKey`、`DecodeAPIKeyPayload`）

### 4.3 pages層の修正

- **バリデーション**: pages層でバリデーションを実施（既存の`validateDmUserInput`を維持）
- **入出力制御**: pages層で入出力制御を実施（既存の`renderDmUserRegisterForm`、`renderAPIKeyPage`、`renderAPIKeyResult`を維持）
- **usecase層の呼び出し**: usecase層のメソッドを呼び出すように変更
- **既存の関数の削除**: `checkEmailExistsSharded`、`insertDmUserSharded`、`generatePublicAPIKey`を削除（service層/repository層に移動済み）
- **メールアドレスの重複チェック**: usecase層で`DmUserService.CheckEmailExists`を呼び出し（repository層で全シャード検索）
- **関数シグネチャの変更**: `DmUserRegisterPage`、`APIKeyPage`のシグネチャを変更（usecase層を追加）

### 4.4 依存関係の注入

- **usecase層の初期化**: main.goでusecase層を初期化
- **service層の初期化**: main.goでservice層を初期化
- **Repository層の初期化**: main.goでrepository層を初期化
- **pages層への渡し方**: pages層の関数呼び出し時にusecase層を渡す

### 4.5 テストの実装

- **usecase層のテスト**: service層のインターフェースのモックを使用してテスト
- **service層のテスト**: 下位層のモックを使用してテスト（必要に応じて）

### 4.6 ドキュメントの更新

- **アーキテクチャドキュメント**: Adminアプリのレイヤー構造を明確に記載
- **プロジェクト構造ドキュメント**: 新規作成するファイルを反映
- **ファイル組織ドキュメント**: 新規作成するファイルを反映
- **一貫性**: 全てのドキュメントで同じレイヤー構造を記載

## 5. 参考情報

### 5.1 関連ドキュメント
- `docs/Architecture.md`: アーキテクチャドキュメント
- `docs/Project-Structure.md`: プロジェクト構造ドキュメント
- `.kiro/steering/structure.md`: ファイル組織とコードパターン
- `.kiro/steering/tech.md`: 技術スタックとアーキテクチャ

### 5.2 既存実装の参考
- `server/internal/usecase/dm_user_usecase.go`: 既存のusecase層の実装パターン
- `server/internal/usecase/cli/list_dm_users_usecase.go`: 既存のCLI用usecase層の実装パターン
- `server/internal/service/dm_user_service.go`: 既存のservice層の実装パターン
- `server/internal/admin/pages/dm_user_register.go`: 既存のAdminページの実装
- `server/internal/admin/pages/api_key.go`: 既存のAdminページの実装

### 5.3 技術スタック
- **言語**: Go
- **アーキテクチャ**: レイヤードアーキテクチャ（pages -> usecase -> service -> repository -> db）
- **テスト**: `testing`（標準ライブラリ）、`github.com/stretchr/testify`（アサーション、モック）
- **Adminフレームワーク**: GoAdmin

### 5.4 レイヤー構造の比較

| 項目 | 現在（修正前） | 修正後 |
|------|---------------|--------|
| Pages層 | pages.go（ビジネスロジックが直接実装） | pages.go（エントリーポイント、バリデーション、入出力） |
| Usecase層 | なし | `usecase/admin/DmUserRegisterUsecase`、`usecase/admin/APIKeyUsecase` |
| Service層 | なし（dm_user_service.goは存在するが使用されていない） | `service.DmUserService`（既存）、`service.APIKeyService`（新規） |
| Repository層 | なし（pages層から直接DB操作） | `repository.DmUserRepository`（既存） |
| DB層 | pages層内で直接使用 | `db.GroupManager`（既存） |
| 出力 | pages層内で直接出力 | pages層内で出力（usecase層から取得） |

### 5.5 APIサーバーとの比較

| 項目 | APIサーバー | Adminアプリ（修正後） |
|------|------------|---------------------|
| エントリーポイント | `server/cmd/server/main.go` | `server/internal/admin/pages/*.go` |
| バリデーション | API Layer（Handler） | Pages層 |
| Usecase層 | `usecase.DmUserUsecase` | `usecase/admin.DmUserRegisterUsecase`、`usecase/admin.APIKeyUsecase` |
| Service層 | `service.DmUserService` | `service.DmUserService`（既存）、`service.APIKeyService`（新規） |
| Repository層 | `repository.DmUserRepository` | `repository.DmUserRepository`（既存） |
| DB層 | `db.GroupManager` | `db.GroupManager` |
| 出力 | HTTPレスポンス | HTML（GoAdminのPanel） |
