# Adminアプリのカスタムページの実装の仕組みを変更するの実装タスク一覧

## 概要
Adminアプリのカスタムページ（`dm_user_register.go`、`api_key.go`）の実装を、APIサーバーと同じレイヤー構造（pages -> usecase -> service -> repository -> db）に変更するためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: Repository層の拡張

#### タスク 1.1: DmUserRepositoryInterfaceへのCheckEmailExistsメソッドの追加
**目的**: `DmUserRepositoryInterface`にメールアドレスの重複チェックメソッドを追加する。

**作業内容**:
- `server/internal/repository/interfaces.go`を修正
- `CheckEmailExists(ctx context.Context, email string) (bool, error)`メソッドをインターフェースに追加

**実装内容**:
- 修正対象: `server/internal/repository/interfaces.go`
- インターフェース定義:
  ```go
  type DmUserRepositoryInterface interface {
      // ... 既存のメソッド ...
      CheckEmailExists(ctx context.Context, email string) (bool, error)
  }
  ```

**受け入れ基準**:
- [ ] `server/internal/repository/interfaces.go`の`DmUserRepositoryInterface`に`CheckEmailExists`メソッドが追加されている
- [ ] メソッドシグネチャが正しい（`ctx context.Context, email string) (bool, error)`）

- _Requirements: 3.2.1_
- _Design: 3.2.1_

---

#### タスク 1.2: DmUserRepositoryへのCheckEmailExistsメソッドの実装
**目的**: `DmUserRepository`にメールアドレスの重複チェックメソッドを実装する。

**作業内容**:
- `server/internal/repository/dm_user_repository.go`を修正
- `CheckEmailExists(ctx context.Context, email string) (bool, error)`メソッドを実装
- 既存の`checkEmailExistsSharded`関数のロジックを移行

**実装内容**:
- 修正対象: `server/internal/repository/dm_user_repository.go`
- メソッド定義:
  ```go
  func (r *DmUserRepository) CheckEmailExists(ctx context.Context, email string) (bool, error)
  ```
- 実装内容:
  - 全テーブルを検索（`tableCount`分ループ）
  - 各テーブルでメールアドレスの存在をチェック（`WHERE email = ?`）
  - 1つでも見つかったら`true`を返す
  - 全てのテーブルで見つからなかったら`false`を返す
  - 既存の`checkEmailExistsSharded`関数のロジックを移行

**受け入れ基準**:
- [ ] `server/internal/repository/dm_user_repository.go`に`CheckEmailExists`メソッドが実装されている
- [ ] 全シャードを検索する処理が実装されている
- [ ] メールアドレスの存在チェックが正しく動作する
- [ ] 既存の`checkEmailExistsSharded`関数のロジックが移行されている
- [ ] `context.Context`を受け取っている

- _Requirements: 3.2.1_
- _Design: 3.2.1_

---

### Phase 2: Service層の拡張

#### タスク 2.1: DmUserServiceInterfaceへのCheckEmailExistsメソッドの追加
**目的**: `DmUserServiceInterface`にメールアドレスの重複チェックメソッドを追加する。

**作業内容**:
- `server/internal/usecase/dm_user_usecase.go`を修正
- `CheckEmailExists(ctx context.Context, email string) (bool, error)`メソッドをインターフェースに追加

**実装内容**:
- 修正対象: `server/internal/usecase/dm_user_usecase.go`
- インターフェース定義:
  ```go
  type DmUserServiceInterface interface {
      // ... 既存のメソッド ...
      CheckEmailExists(ctx context.Context, email string) (bool, error)
  }
  ```

**受け入れ基準**:
- [ ] `server/internal/usecase/dm_user_usecase.go`の`DmUserServiceInterface`に`CheckEmailExists`メソッドが追加されている
- [ ] メソッドシグネチャが正しい（`ctx context.Context, email string) (bool, error)`）

- _Requirements: 3.2.1_
- _Design: 3.2.1_

---

#### タスク 2.2: DmUserServiceへのCheckEmailExistsメソッドの実装
**目的**: `DmUserService`にメールアドレスの重複チェックメソッドを実装する。

**作業内容**:
- `server/internal/service/dm_user_service.go`を修正
- `CheckEmailExists(ctx context.Context, email string) (bool, error)`メソッドを実装
- repository層の`CheckEmailExists`を呼び出し

**実装内容**:
- 修正対象: `server/internal/service/dm_user_service.go`
- メソッド定義:
  ```go
  func (s *DmUserService) CheckEmailExists(ctx context.Context, email string) (bool, error)
  ```
- 実装内容:
  - repository層の`CheckEmailExists`を呼び出し
  - エラーハンドリングを実装

**受け入れ基準**:
- [ ] `server/internal/service/dm_user_service.go`に`CheckEmailExists`メソッドが実装されている
- [ ] repository層の`CheckEmailExists`を呼び出している
- [ ] エラーハンドリングが実装されている
- [ ] `context.Context`を受け取っている

- _Requirements: 3.2.1_
- _Design: 3.2.1_

---

#### タスク 2.3: APIKeyServiceの新規作成
**目的**: APIキー発行用のservice層を新規作成する。

**作業内容**:
- `server/internal/service/api_key_service.go`を新規作成
- `APIKeyService`構造体を定義
- `APIKeyServiceInterface`を定義
- `GenerateAPIKey`、`DecodeAPIKeyPayload`メソッドを実装

**実装内容**:
- ファイルパス: `server/internal/service/api_key_service.go`
- パッケージ名: `service`
- インターフェース定義:
  ```go
  type APIKeyServiceInterface interface {
      GenerateAPIKey(ctx context.Context, secretKey, version, env string, issuedAt int64) (string, error)
      DecodeAPIKeyPayload(ctx context.Context, token string) (*auth.JWTClaims, error)
  }
  ```
- 構造体定義:
  ```go
  type APIKeyService struct{}
  ```
- コンストラクタ:
  ```go
  func NewAPIKeyService() *APIKeyService
  ```
- メソッド定義:
  ```go
  func (s *APIKeyService) GenerateAPIKey(ctx context.Context, secretKey, version, env string, issuedAt int64) (string, error)
  func (s *APIKeyService) DecodeAPIKeyPayload(ctx context.Context, token string) (*auth.JWTClaims, error)
  ```
- 実装内容:
  - `GenerateAPIKey`: `auth.GeneratePublicAPIKey`を呼び出し
  - `DecodeAPIKeyPayload`: `auth.ParseJWTClaims`を呼び出し
  - エラーハンドリングを実装

**受け入れ基準**:
- [ ] `server/internal/service/api_key_service.go`が作成されている
- [ ] `APIKeyServiceInterface`が定義されている
- [ ] `APIKeyService`構造体が定義されている
- [ ] `NewAPIKeyService`コンストラクタが実装されている
- [ ] `GenerateAPIKey`メソッドが実装されている
- [ ] `DecodeAPIKeyPayload`メソッドが実装されている
- [ ] `auth.GeneratePublicAPIKey`を呼び出している
- [ ] `auth.ParseJWTClaims`を呼び出している
- [ ] エラーハンドリングが実装されている
- [ ] `context.Context`を受け取っている

- _Requirements: 3.2.2, 6.2_
- _Design: 3.2.2_

---

### Phase 3: Admin用usecase層の作成

#### タスク 3.1: adminディレクトリの作成
**目的**: Adminアプリ用のusecase層のディレクトリを新規作成する。

**作業内容**:
- `server/internal/usecase/admin`ディレクトリを新規作成

**実装内容**:
- ディレクトリパス: `server/internal/usecase/admin`

**受け入れ基準**:
- [ ] `server/internal/usecase/admin`ディレクトリが作成されている

- _Requirements: 3.1.1, 6.1_
- _Design: 3.1.1_

---

#### タスク 3.2: DmUserRegisterUsecaseの実装
**目的**: dm_user登録用のusecaseを実装する。

**作業内容**:
- `server/internal/usecase/admin/dm_user_register_usecase.go`を新規作成
- `DmUserRegisterUsecase`構造体を定義
- `RegisterDmUser`メソッドを実装

**実装内容**:
- ファイルパス: `server/internal/usecase/admin/dm_user_register_usecase.go`
- パッケージ名: `admin`
- 構造体定義:
  ```go
  type DmUserRegisterUsecase struct {
      dmUserService usecase.DmUserServiceInterface
  }
  ```
- コンストラクタ:
  ```go
  func NewDmUserRegisterUsecase(dmUserService usecase.DmUserServiceInterface) *DmUserRegisterUsecase
  ```
- メソッド定義:
  ```go
  func (u *DmUserRegisterUsecase) RegisterDmUser(ctx context.Context, name, email string) (string, error)
  ```
- 実装内容:
  - メールアドレスの重複チェック（`dmUserService.CheckEmailExists`を呼び出し）
  - 重複している場合はエラーを返す
  - `CreateDmUserRequest`を作成
  - `dmUserService.CreateDmUser`を呼び出し
  - エラーハンドリングを実装

**受け入れ基準**:
- [ ] `server/internal/usecase/admin/dm_user_register_usecase.go`が作成されている
- [ ] `DmUserRegisterUsecase`構造体が定義されている
- [ ] `DmUserServiceInterface`を依存として注入している
- [ ] `NewDmUserRegisterUsecase`コンストラクタが実装されている
- [ ] `RegisterDmUser`メソッドが実装されている
- [ ] メールアドレスの重複チェックが実装されている
- [ ] `CreateDmUserRequest`を作成している
- [ ] `dmUserService.CreateDmUser`を呼び出している
- [ ] エラーハンドリングが実装されている
- [ ] `context.Context`を受け取っている

- _Requirements: 3.1.2, 6.1_
- _Design: 3.1.1_

---

#### タスク 3.3: APIKeyUsecaseの実装
**目的**: APIキー発行用のusecaseを実装する。

**作業内容**:
- `server/internal/usecase/admin/api_key_usecase.go`を新規作成
- `APIKeyUsecase`構造体を定義
- `GenerateAPIKey`、`DecodeAPIKeyPayload`メソッドを実装

**実装内容**:
- ファイルパス: `server/internal/usecase/admin/api_key_usecase.go`
- パッケージ名: `admin`
- 構造体定義:
  ```go
  type APIKeyUsecase struct {
      apiKeyService service.APIKeyServiceInterface
      cfg           *config.Config
  }
  ```
- コンストラクタ:
  ```go
  func NewAPIKeyUsecase(apiKeyService service.APIKeyServiceInterface, cfg *config.Config) *APIKeyUsecase
  ```
- メソッド定義:
  ```go
  func (u *APIKeyUsecase) GenerateAPIKey(ctx context.Context, env string) (string, error)
  func (u *APIKeyUsecase) DecodeAPIKeyPayload(ctx context.Context, token string) (*auth.JWTClaims, error)
  ```
- 実装内容:
  - `GenerateAPIKey`: 環境変数の取得、現在時刻の取得、`apiKeyService.GenerateAPIKey`を呼び出し
  - `DecodeAPIKeyPayload`: `apiKeyService.DecodeAPIKeyPayload`を呼び出し
  - エラーハンドリングを実装

**受け入れ基準**:
- [ ] `server/internal/usecase/admin/api_key_usecase.go`が作成されている
- [ ] `APIKeyUsecase`構造体が定義されている
- [ ] `APIKeyServiceInterface`を依存として注入している
- [ ] `config.Config`を依存として注入している
- [ ] `NewAPIKeyUsecase`コンストラクタが実装されている
- [ ] `GenerateAPIKey`メソッドが実装されている
- [ ] `DecodeAPIKeyPayload`メソッドが実装されている
- [ ] `apiKeyService.GenerateAPIKey`を呼び出している
- [ ] `apiKeyService.DecodeAPIKeyPayload`を呼び出している
- [ ] エラーハンドリングが実装されている
- [ ] `context.Context`を受け取っている

- _Requirements: 3.1.3, 6.1_
- _Design: 3.1.2_

---

### Phase 4: pages層の修正

#### タスク 4.1: dm_user_register.goの修正
**目的**: `dm_user_register.go`を修正し、usecase層を使用するように変更する。

**作業内容**:
- `server/internal/admin/pages/dm_user_register.go`を修正
- `DmUserRegisterPage`関数のシグネチャを変更（`dmUserRegisterUsecase`を追加）
- `handleDmUserRegisterPost`関数を修正（usecase層を呼び出すように変更）
- `checkEmailExistsSharded`関数を削除
- `insertDmUserSharded`関数を削除
- `validateDmUserInput`関数を維持
- `renderDmUserRegisterForm`関数を維持

**実装内容**:
- 修正対象: `server/internal/admin/pages/dm_user_register.go`
- 関数シグネチャの変更:
  ```go
  func DmUserRegisterPage(ctx *context.Context, groupManager *appdb.GroupManager, dmUserRegisterUsecase *admin.DmUserRegisterUsecase) (types.Panel, error)
  ```
- `handleDmUserRegisterPost`関数の修正:
  - usecase層の`RegisterDmUser`を呼び出すように変更
  - `checkEmailExistsSharded`の呼び出しを削除
  - `insertDmUserSharded`の呼び出しを削除
- 削除する関数:
  - `checkEmailExistsSharded`
  - `insertDmUserSharded`
- 維持する関数:
  - `validateDmUserInput`
  - `renderDmUserRegisterForm`

**受け入れ基準**:
- [ ] `server/internal/admin/pages/dm_user_register.go`が修正されている
- [ ] `DmUserRegisterPage`関数のシグネチャが変更されている（`dmUserRegisterUsecase`を追加）
- [ ] `handleDmUserRegisterPost`関数がusecase層を呼び出すように変更されている
- [ ] `checkEmailExistsSharded`関数が削除されている
- [ ] `insertDmUserSharded`関数が削除されている
- [ ] `validateDmUserInput`関数が維持されている
- [ ] `renderDmUserRegisterForm`関数が維持されている
- [ ] usecase層の`RegisterDmUser`を呼び出している

- _Requirements: 3.4.1, 6.3_
- _Design: 3.3.1_

---

#### タスク 4.2: api_key.goの修正
**目的**: `api_key.go`を修正し、usecase層を使用するように変更する（2つの処理に分ける）。

**作業内容**:
- `server/internal/admin/pages/api_key.go`を修正
- `APIKeyPage`関数のシグネチャを変更（`apiKeyUsecase`を追加）
- `handleGenerateKey`関数を修正（usecase層を呼び出すように変更、2つの処理に分ける）
- `generatePublicAPIKey`関数を削除
- `renderAPIKeyPage`関数を維持
- `renderAPIKeyResult`関数を維持

**実装内容**:
- 修正対象: `server/internal/admin/pages/api_key.go`
- 関数シグネチャの変更:
  ```go
  func APIKeyPage(ctx *context.Context, conn db.Connection, apiKeyUsecase *admin.APIKeyUsecase) (types.Panel, error)
  ```
- `handleGenerateKey`関数の修正:
  - usecase層の`GenerateAPIKey`を呼び出し（鍵の生成）
  - usecase層の`DecodeAPIKeyPayload`を呼び出し（ペイロードのデコード）
  - 2つの処理に分ける
- 削除する関数:
  - `generatePublicAPIKey`
- 維持する関数:
  - `renderAPIKeyPage`
  - `renderAPIKeyResult`

**受け入れ基準**:
- [ ] `server/internal/admin/pages/api_key.go`が修正されている
- [ ] `APIKeyPage`関数のシグネチャが変更されている（`apiKeyUsecase`を追加）
- [ ] `handleGenerateKey`関数がusecase層を呼び出すように変更されている
- [ ] 鍵の生成とペイロードのデコードが2つの処理に分かれている
- [ ] `generatePublicAPIKey`関数が削除されている
- [ ] `renderAPIKeyPage`関数が維持されている
- [ ] `renderAPIKeyResult`関数が維持されている
- [ ] usecase層の`GenerateAPIKey`を呼び出している
- [ ] usecase層の`DecodeAPIKeyPayload`を呼び出している

- _Requirements: 3.4.2, 6.3_
- _Design: 3.3.2_

---

### Phase 5: 依存関係の注入（main.goの修正）

#### タスク 5.1: main.goでの依存関係の初期化
**目的**: main.goでRepository層、Service層、Usecase層を初期化する。

**作業内容**:
- `server/cmd/admin/main.go`を修正
- Repository層の初期化（`DmUserRepository`）
- Service層の初期化（`DmUserService`、`APIKeyService`）
- Usecase層の初期化（`DmUserRegisterUsecase`、`APIKeyUsecase`）

**実装内容**:
- 修正対象: `server/cmd/admin/main.go`
- 初期化処理:
  ```go
  // Repository層の初期化
  dmUserRepository := repository.NewDmUserRepository(groupManager)
  
  // Service層の初期化
  dmUserService := service.NewDmUserService(dmUserRepository)
  apiKeyService := service.NewAPIKeyService()
  
  // Usecase層の初期化
  dmUserRegisterUsecase := admin.NewDmUserRegisterUsecase(dmUserService)
  apiKeyUsecase := admin.NewAPIKeyUsecase(apiKeyService, cfg)
  ```

**受け入れ基準**:
- [ ] `server/cmd/admin/main.go`が修正されている
- [ ] Repository層の初期化が実装されている（`DmUserRepository`）
- [ ] Service層の初期化が実装されている（`DmUserService`、`APIKeyService`）
- [ ] Usecase層の初期化が実装されている（`DmUserRegisterUsecase`、`APIKeyUsecase`）
- [ ] 依存関係が正しく注入されている

- _Requirements: 3.5.1, 3.5.2, 6.4_
- _Design: 3.4.1_

---

#### タスク 5.2: main.goでのpages層へのusecase層の渡し方の修正
**目的**: main.goでpages層の関数呼び出し時にusecase層を渡すように修正する。

**作業内容**:
- `server/cmd/admin/main.go`を修正
- `DmUserRegisterPage`関数呼び出し時に`dmUserRegisterUsecase`を渡す
- `APIKeyPage`関数呼び出し時に`apiKeyUsecase`を渡す

**実装内容**:
- 修正対象: `server/cmd/admin/main.go`
- 関数呼び出しの修正:
  ```go
  app.HandleFunc("/admin/dm-user/register", gorillaAdapter.Content(func(ctx gorillaAdapter.Context) (types.Panel, error) {
      return pages.DmUserRegisterPage(goadminContext.NewContext(ctx.Request), groupManager, dmUserRegisterUsecase)
  })).Methods("GET", "POST")
  
  app.HandleFunc("/admin/api-key", gorillaAdapter.Content(func(ctx gorillaAdapter.Context) (types.Panel, error) {
      return pages.APIKeyPage(goadminContext.NewContext(ctx.Request), conn, apiKeyUsecase)
  })).Methods("GET", "POST")
  ```

**受け入れ基準**:
- [ ] `server/cmd/admin/main.go`が修正されている
- [ ] `DmUserRegisterPage`関数呼び出し時に`dmUserRegisterUsecase`を渡している
- [ ] `APIKeyPage`関数呼び出し時に`apiKeyUsecase`を渡している
- [ ] 既存の動作が維持されている

- _Requirements: 3.5.1, 6.4_
- _Design: 3.4.1_

---

### Phase 6: テストの実装

#### タスク 6.1: DmUserRepositoryのCheckEmailExistsのテスト追加（存在する場合）
**目的**: `DmUserRepository`の`CheckEmailExists`メソッドのテストを追加する。

**作業内容**:
- `server/internal/repository/dm_user_repository_test.go`が存在する場合は修正
- `CheckEmailExists`のテストケースを追加

**実装内容**:
- 修正対象: `server/internal/repository/dm_user_repository_test.go`（存在する場合）
- テストケース:
  1. 正常系: メールアドレスが存在する場合
  2. 正常系: メールアドレスが存在しない場合
  3. 正常系: 全シャードを検索する処理のテスト
  4. 異常系: エラーハンドリングのテスト
- テーブル駆動テストを使用
- `github.com/stretchr/testify/assert`を使用

**受け入れ基準**:
- [ ] `server/internal/repository/dm_user_repository_test.go`に`CheckEmailExists`のテストが追加されている（存在する場合）
- [ ] 正常系のテストケースが実装されている
- [ ] 異常系のテストケースが実装されている
- [ ] テーブル駆動テストを使用している

- _Requirements: 6.6_
- _Design: 3.5.2_

---

#### タスク 6.2: DmUserServiceのCheckEmailExistsのテスト追加（存在する場合）
**目的**: `DmUserService`の`CheckEmailExists`メソッドのテストを追加する。

**作業内容**:
- `server/internal/service/dm_user_service_test.go`が存在する場合は修正
- `CheckEmailExists`のテストケースを追加

**実装内容**:
- 修正対象: `server/internal/service/dm_user_service_test.go`（存在する場合）
- テストケース:
  1. 正常系: メールアドレスが存在する場合
  2. 正常系: メールアドレスが存在しない場合
  3. 異常系: repository層からエラーが返された場合
- repository層のモックを使用
- テーブル駆動テストを使用
- `github.com/stretchr/testify/assert`を使用

**受け入れ基準**:
- [ ] `server/internal/service/dm_user_service_test.go`に`CheckEmailExists`のテストが追加されている（存在する場合）
- [ ] 正常系のテストケースが実装されている
- [ ] 異常系のテストケースが実装されている
- [ ] repository層のモックを使用している
- [ ] テーブル駆動テストを使用している

- _Requirements: 6.6_
- _Design: 3.5.2_

---

#### タスク 6.3: APIKeyServiceのテスト実装
**目的**: `APIKeyService`のテストを実装する。

**作業内容**:
- `server/internal/service/api_key_service_test.go`を新規作成
- `GenerateAPIKey`、`DecodeAPIKeyPayload`のテストケースを実装

**実装内容**:
- ファイルパス: `server/internal/service/api_key_service_test.go`
- パッケージ名: `service`
- テストケース:
  1. `GenerateAPIKey`の正常系テスト
  2. `GenerateAPIKey`の異常系テスト（auth層からエラーが返された場合）
  3. `DecodeAPIKeyPayload`の正常系テスト
  4. `DecodeAPIKeyPayload`の異常系テスト（auth層からエラーが返された場合）
- テーブル駆動テストを使用
- `github.com/stretchr/testify/assert`を使用

**受け入れ基準**:
- [ ] `server/internal/service/api_key_service_test.go`が作成されている
- [ ] `GenerateAPIKey`のテストケースが実装されている
- [ ] `DecodeAPIKeyPayload`のテストケースが実装されている
- [ ] 正常系のテストケースが実装されている
- [ ] 異常系のテストケースが実装されている
- [ ] テーブル駆動テストを使用している

- _Requirements: 6.2, 6.6_
- _Design: 3.2.3_

---

#### タスク 6.4: DmUserRegisterUsecaseのテスト実装
**目的**: `DmUserRegisterUsecase`のテストを実装する。

**作業内容**:
- `server/internal/usecase/admin/dm_user_register_usecase_test.go`を新規作成
- `RegisterDmUser`のテストケースを実装

**実装内容**:
- ファイルパス: `server/internal/usecase/admin/dm_user_register_usecase_test.go`
- パッケージ名: `admin`
- テストケース:
  1. 正常系: ユーザー登録が正常に完了する場合
  2. 異常系: メールアドレスが重複している場合
  3. 異常系: service層からエラーが返された場合
- service層のモック（`MockDmUserServiceInterface`）を使用
- テーブル駆動テストを使用
- `github.com/stretchr/testify/assert`を使用

**受け入れ基準**:
- [ ] `server/internal/usecase/admin/dm_user_register_usecase_test.go`が作成されている
- [ ] `RegisterDmUser`のテストケースが実装されている
- [ ] 正常系のテストケースが実装されている
- [ ] 異常系のテストケースが実装されている
- [ ] service層のモックを使用している
- [ ] テーブル駆動テストを使用している

- _Requirements: 6.1, 6.6_
- _Design: 3.1.3_

---

#### タスク 6.5: APIKeyUsecaseのテスト実装
**目的**: `APIKeyUsecase`のテストを実装する。

**作業内容**:
- `server/internal/usecase/admin/api_key_usecase_test.go`を新規作成
- `GenerateAPIKey`、`DecodeAPIKeyPayload`のテストケースを実装

**実装内容**:
- ファイルパス: `server/internal/usecase/admin/api_key_usecase_test.go`
- パッケージ名: `admin`
- テストケース:
  1. `GenerateAPIKey`の正常系テスト
  2. `GenerateAPIKey`の異常系テスト（service層からエラーが返された場合）
  3. `DecodeAPIKeyPayload`の正常系テスト
  4. `DecodeAPIKeyPayload`の異常系テスト（service層からエラーが返された場合）
- service層のモック（`MockAPIKeyServiceInterface`）を使用
- テーブル駆動テストを使用
- `github.com/stretchr/testify/assert`を使用

**受け入れ基準**:
- [ ] `server/internal/usecase/admin/api_key_usecase_test.go`が作成されている
- [ ] `GenerateAPIKey`のテストケースが実装されている
- [ ] `DecodeAPIKeyPayload`のテストケースが実装されている
- [ ] 正常系のテストケースが実装されている
- [ ] 異常系のテストケースが実装されている
- [ ] service層のモックを使用している
- [ ] テーブル駆動テストを使用している

- _Requirements: 6.1, 6.6_
- _Design: 3.1.3_

---

### Phase 7: ドキュメントの更新

#### タスク 7.1: Architecture.mdの更新
**目的**: アーキテクチャドキュメントにAdminアプリのレイヤー構造を追加する。

**作業内容**:
- `docs/Architecture.md`を修正
- Adminアプリのレイヤー構造を追加（usecase層を含む）
- Adminアプリのアーキテクチャ図を更新
- Admin用usecase層の説明を追加

**実装内容**:
- 修正対象: `docs/Architecture.md`
- 追加内容:
  - Adminアプリのレイヤー構造の説明（pages -> usecase -> service -> repository -> db）
  - Adminアプリのアーキテクチャ図
  - Admin用usecase層の説明

**受け入れ基準**:
- [ ] `docs/Architecture.md`にAdminアプリのレイヤー構造が追加されている
- [ ] Adminアプリのアーキテクチャ図が更新されている
- [ ] Admin用usecase層の説明が追加されている

- _Requirements: 3.6.1, 6.7_
- _Design: 3.6.1_

---

#### タスク 7.2: Project-Structure.mdの更新
**目的**: プロジェクト構造ドキュメントに新規作成するファイルを追加する。

**作業内容**:
- `docs/Project-Structure.md`を修正
- `server/internal/usecase/admin`ディレクトリを追加
- 新規作成するファイルを追加

**実装内容**:
- 修正対象: `docs/Project-Structure.md`
- 追加内容:
  - `server/internal/usecase/admin`ディレクトリ
  - `server/internal/usecase/admin/dm_user_register_usecase.go`
  - `server/internal/usecase/admin/api_key_usecase.go`
  - `server/internal/service/api_key_service.go`

**受け入れ基準**:
- [ ] `docs/Project-Structure.md`に新規作成するファイルが追加されている
- [ ] `server/internal/usecase/admin`ディレクトリが追加されている

- _Requirements: 3.6.2, 6.7_
- _Design: 3.6.2_

---

#### タスク 7.3: structure.mdの更新
**目的**: ファイル組織ドキュメントに新規作成するファイルを追加する。

**作業内容**:
- `.kiro/steering/structure.md`を修正
- `server/internal/usecase/admin`ディレクトリを追加
- 新規作成するファイルを追加

**実装内容**:
- 修正対象: `.kiro/steering/structure.md`
- 追加内容:
  - `server/internal/usecase/admin`ディレクトリ
  - `server/internal/usecase/admin/dm_user_register_usecase.go`
  - `server/internal/usecase/admin/api_key_usecase.go`
  - `server/internal/service/api_key_service.go`

**受け入れ基準**:
- [ ] `.kiro/steering/structure.md`に新規作成するファイルが追加されている
- [ ] `server/internal/usecase/admin`ディレクトリが追加されている

- _Requirements: 3.6.3, 6.7_
- _Design: 3.6.3_

---

### Phase 8: 動作確認

#### タスク 8.1: ローカル環境での動作確認
**目的**: ローカル環境でAdminページが正常に動作することを確認する。

**作業内容**:
- ローカル環境でAdminサーバーを起動
- dm_user登録ページが正常に動作することを確認
- APIキー発行ページが正常に動作することを確認
- 既存の出力形式（HTML出力）が維持されていることを確認
- 既存のエラーメッセージが維持されていることを確認

**実装内容**:
- 動作確認項目:
  1. dm_user登録が正常に動作する
  2. APIキー発行が正常に動作する
  3. 既存の出力形式（HTML出力）が維持されている
  4. 既存のエラーメッセージが維持されている
  5. メールアドレスの重複チェックが正常に動作する

**受け入れ基準**:
- [ ] ローカル環境でAdminページが正常に動作する
- [ ] dm_user登録が正常に動作する
- [ ] APIキー発行が正常に動作する
- [ ] 既存の出力形式（HTML出力）が維持されている
- [ ] 既存のエラーメッセージが維持されている
- [ ] メールアドレスの重複チェックが正常に動作する

- _Requirements: 6.5_
- _Design: 4.2_

---

#### タスク 8.2: 既存テストの確認
**目的**: 既存のテストが全て通過することを確認する。

**作業内容**:
- 既存のテストを実行
- 全てのテストが通過することを確認

**実装内容**:
- テスト実行:
  ```bash
  go test ./...
  ```

**受け入れ基準**:
- [ ] 既存のテストが全て通過する
- [ ] 新規追加したテストも全て通過する

- _Requirements: 6.5, 6.6_
- _Design: 4.2_

---

## 実装順序の推奨

1. **Phase 1**: Repository層の拡張（CheckEmailExistsメソッドの追加）
2. **Phase 2**: Service層の拡張（CheckEmailExistsメソッドの追加、APIKeyServiceの新規作成）
3. **Phase 3**: Admin用usecase層の作成
4. **Phase 4**: pages層の修正
5. **Phase 5**: 依存関係の注入（main.goの修正）
6. **Phase 6**: テストの実装
7. **Phase 7**: ドキュメントの更新
8. **Phase 8**: 動作確認

## 注意事項

- 各タスクは独立して実装可能な粒度に分解されている
- テストは実装と並行して進めることを推奨
- ドキュメントの更新は実装完了後に実施することを推奨
- 動作確認は各フェーズ完了時点で実施することを推奨
