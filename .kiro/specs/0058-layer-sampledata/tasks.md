# server/cmd/generate-sample-dataの構造修正の実装タスク一覧

## 概要
`server/cmd/generate-sample-data`の実装を、APIサーバーと同じレイヤー構造（usecase -> service -> repository -> db）に変更するためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: Repository層の拡張

#### タスク 1.1: dm_user_repository.goへのInsertDmUsersBatchメソッドの追加
**目的**: 既存の`dm_user_repository.go`にバッチ挿入メソッドを追加する。

**作業内容**:
- `server/internal/repository/dm_user_repository.go`を修正
- `InsertDmUsersBatch(ctx context.Context, tableName string, dmUsers []*model.DmUser) error`メソッドを追加
- 既存の`insertDmUsersBatch`関数のロジックを移行

**実装内容**:
- 修正対象: `server/internal/repository/dm_user_repository.go`
- メソッド定義:
  ```go
  func (r *DmUserRepository) InsertDmUsersBatch(ctx context.Context, tableName string, dmUsers []*model.DmUser) error
  ```
- 実装内容:
  - バッチサイズ500件ずつ挿入する処理を実装
  - テーブル名からテーブル番号を抽出（`extractTableNumber`関数を使用）
  - テーブル番号から接続を取得（`r.groupManager.GetShardingConnection()`）
  - 生SQLでバッチ挿入（動的テーブル名対応、IDを含む）
  - リトライ機能（`db.ExecuteWithRetry`）を使用
  - 既存の`insertDmUsersBatch`関数のロジックを移行

**受け入れ基準**:
- [ ] `server/internal/repository/dm_user_repository.go`に`InsertDmUsersBatch`メソッドが追加されている
- [ ] バッチサイズ500件ずつ挿入する処理が実装されている
- [ ] テーブル名からテーブル番号を抽出する処理が実装されている
- [ ] リトライ機能（`db.ExecuteWithRetry`）が使用されている
- [ ] 既存の`insertDmUsersBatch`関数のロジックが移行されている
- [ ] `context.Context`を受け取っている

- _Requirements: 3.1.1, 6.1_
- _Design: 3.1.1_

---

#### タスク 1.2: dm_post_repository.goへのInsertDmPostsBatchメソッドの追加
**目的**: 既存の`dm_post_repository.go`にバッチ挿入メソッドを追加する。

**作業内容**:
- `server/internal/repository/dm_post_repository.go`を修正
- `InsertDmPostsBatch(ctx context.Context, tableName string, dmPosts []*model.DmPost) error`メソッドを追加
- 既存の`insertDmPostsBatch`関数のロジックを移行

**実装内容**:
- 修正対象: `server/internal/repository/dm_post_repository.go`
- メソッド定義:
  ```go
  func (r *DmPostRepository) InsertDmPostsBatch(ctx context.Context, tableName string, dmPosts []*model.DmPost) error
  ```
- 実装内容:
  - バッチサイズ500件ずつ挿入する処理を実装
  - テーブル名からテーブル番号を抽出（`extractTableNumber`関数を使用）
  - テーブル番号から接続を取得（`r.groupManager.GetShardingConnection()`）
  - 生SQLでバッチ挿入（動的テーブル名対応、IDを含む）
  - リトライ機能（`db.ExecuteWithRetry`）を使用
  - 既存の`insertDmPostsBatch`関数のロジックを移行

**受け入れ基準**:
- [ ] `server/internal/repository/dm_post_repository.go`に`InsertDmPostsBatch`メソッドが追加されている
- [ ] バッチサイズ500件ずつ挿入する処理が実装されている
- [ ] テーブル名からテーブル番号を抽出する処理が実装されている
- [ ] リトライ機能（`db.ExecuteWithRetry`）が使用されている
- [ ] 既存の`insertDmPostsBatch`関数のロジックが移行されている
- [ ] `context.Context`を受け取っている

- _Requirements: 3.1.1, 6.1_
- _Design: 3.1.2_

---

#### タスク 1.3: dm_news_repository.goの新規作成
**目的**: `dm_news_repository.go`を新規作成し、バッチ挿入メソッドを実装する。

**作業内容**:
- `server/internal/repository/dm_news_repository.go`を新規作成
- `DmNewsRepository`構造体を定義
- `InsertDmNewsBatch(ctx context.Context, dmNews []*model.DmNews) error`メソッドを実装
- 既存の`insertDmNewsBatch`関数のロジックを移行

**実装内容**:
- ファイルパス: `server/internal/repository/dm_news_repository.go`
- パッケージ名: `repository`
- 構造体定義:
  ```go
  type DmNewsRepository struct {
      groupManager *db.GroupManager
  }
  ```
- コンストラクタ:
  ```go
  func NewDmNewsRepository(groupManager *db.GroupManager) *DmNewsRepository
  ```
- メソッド定義:
  ```go
  func (r *DmNewsRepository) InsertDmNewsBatch(ctx context.Context, dmNews []*model.DmNews) error
  ```
- 実装内容:
  - master接続を取得（`r.groupManager.GetMasterConnection()`）
  - バッチサイズ500件ずつ挿入する処理を実装
  - GORMの`CreateInBatches`を使用（固定テーブル名`dm_news`）
  - リトライ機能（`db.ExecuteWithRetry`）を使用
  - 既存の`insertDmNewsBatch`関数のロジックを移行

**受け入れ基準**:
- [ ] `server/internal/repository/dm_news_repository.go`が作成されている
- [ ] `DmNewsRepository`構造体が定義されている
- [ ] `NewDmNewsRepository`コンストラクタが実装されている
- [ ] `InsertDmNewsBatch(ctx context.Context, dmNews []*model.DmNews) error`メソッドが実装されている
- [ ] master接続を使用している
- [ ] バッチサイズ500件ずつ挿入する処理が実装されている
- [ ] GORMの`CreateInBatches`を使用している（固定テーブル名`dm_news`）
- [ ] リトライ機能（`db.ExecuteWithRetry`）が使用されている
- [ ] 既存の`insertDmNewsBatch`関数のロジックが移行されている
- [ ] `context.Context`を受け取っている

- _Requirements: 3.1.1, 6.1_
- _Design: 3.1.3_

---

#### タスク 1.4: dm_user_repository_test.goへのテスト追加（存在する場合）
**目的**: 既存の`dm_user_repository_test.go`に`InsertDmUsersBatch`のテストを追加する。

**作業内容**:
- `server/internal/repository/dm_user_repository_test.go`が存在する場合は修正
- `InsertDmUsersBatch`のテストケースを追加

**実装内容**:
- 修正対象: `server/internal/repository/dm_user_repository_test.go`（存在する場合）
- テストケース:
  1. 正常系: バッチ挿入が正常に完了する場合
  2. 正常系: 空のスライスを渡した場合
  3. 正常系: バッチサイズを超えるデータのテスト
  4. 異常系: エラーハンドリングのテスト
- テーブル駆動テストを使用
- `github.com/stretchr/testify/assert`を使用

**受け入れ基準**:
- [ ] `server/internal/repository/dm_user_repository_test.go`に`InsertDmUsersBatch`のテストが追加されている（存在する場合）
- [ ] 正常系のテストケースが実装されている
- [ ] 異常系のテストケースが実装されている
- [ ] 全てのテストケースが通過する

- _Requirements: 3.1.2, 6.1, 6.7_
- _Design: 3.1.4_

---

#### タスク 1.5: dm_post_repository_test.goへのテスト追加（存在する場合）
**目的**: 既存の`dm_post_repository_test.go`に`InsertDmPostsBatch`のテストを追加する。

**作業内容**:
- `server/internal/repository/dm_post_repository_test.go`が存在する場合は修正
- `InsertDmPostsBatch`のテストケースを追加

**実装内容**:
- 修正対象: `server/internal/repository/dm_post_repository_test.go`（存在する場合）
- テストケース:
  1. 正常系: バッチ挿入が正常に完了する場合
  2. 正常系: 空のスライスを渡した場合
  3. 正常系: バッチサイズを超えるデータのテスト
  4. 異常系: エラーハンドリングのテスト
- テーブル駆動テストを使用
- `github.com/stretchr/testify/assert`を使用

**受け入れ基準**:
- [ ] `server/internal/repository/dm_post_repository_test.go`に`InsertDmPostsBatch`のテストが追加されている（存在する場合）
- [ ] 正常系のテストケースが実装されている
- [ ] 異常系のテストケースが実装されている
- [ ] 全てのテストケースが通過する

- _Requirements: 3.1.2, 6.1, 6.7_
- _Design: 3.1.4_

---

#### タスク 1.6: dm_news_repository_test.goの新規作成
**目的**: `dm_news_repository_test.go`を新規作成し、`InsertDmNewsBatch`のテストを実装する。

**作業内容**:
- `server/internal/repository/dm_news_repository_test.go`を新規作成
- `InsertDmNewsBatch`のテストケースを実装

**実装内容**:
- ファイルパス: `server/internal/repository/dm_news_repository_test.go`
- パッケージ名: `repository`
- テストケース:
  1. 正常系: バッチ挿入が正常に完了する場合
  2. 正常系: 空のスライスを渡した場合
  3. 正常系: バッチサイズを超えるデータのテスト
  4. 異常系: エラーハンドリングのテスト
- テーブル駆動テストを使用
- `github.com/stretchr/testify/assert`を使用

**受け入れ基準**:
- [ ] `server/internal/repository/dm_news_repository_test.go`が作成されている
- [ ] 正常系のテストケースが実装されている
- [ ] 異常系のテストケースが実装されている
- [ ] 全てのテストケースが通過する

- _Requirements: 3.1.2, 6.1, 6.7_
- _Design: 3.1.4_

---

### Phase 2: Service層の作成

#### タスク 2.1: generate_sample_service.goの実装
**目的**: サンプルデータ生成用のservice層を実装する。

**作業内容**:
- `server/internal/service/generate_sample_service.go`を作成
- `GenerateSampleService`構造体を定義
- `GenerateSampleServiceInterface`を定義
- `GenerateDmUsers`、`GenerateDmPosts`、`GenerateDmNews`メソッドを実装

**実装内容**:
- ファイルパス: `server/internal/service/generate_sample_service.go`
- パッケージ名: `service`
- インターフェース定義:
  ```go
  type GenerateSampleServiceInterface interface {
      GenerateDmUsers(ctx context.Context, totalCount int) ([]string, error)
      GenerateDmPosts(ctx context.Context, dmUserIDs []string, totalCount int) error
      GenerateDmNews(ctx context.Context, totalCount int) error
  }
  ```
- 構造体定義:
  ```go
  type GenerateSampleService struct {
      dmUserRepository *repository.DmUserRepository
      dmPostRepository *repository.DmPostRepository
      dmNewsRepository *repository.DmNewsRepository
      tableSelector    *db.TableSelector
  }
  ```
- コンストラクタ:
  ```go
  func NewGenerateSampleService(
      dmUserRepository *repository.DmUserRepository,
      dmPostRepository *repository.DmPostRepository,
      dmNewsRepository *repository.DmNewsRepository,
      tableSelector *db.TableSelector,
  ) *GenerateSampleService
  ```
- メソッド:
  - `GenerateDmUsers(ctx context.Context, totalCount int) ([]string, error)`
  - `GenerateDmPosts(ctx context.Context, dmUserIDs []string, totalCount int) error`
  - `GenerateDmNews(ctx context.Context, totalCount int) error`
- 実装内容:
  - `gofakeit`を使用したデータ生成ロジックを実装
  - UUID生成（`idgen.GenerateUUIDv7()`）
  - テーブル番号計算（`tableSelector.GetTableNumberFromUUID()`）
  - データ生成、repository層への呼び出しを実装
  - 既存の`generateDmUsers`、`generateDmPosts`、`generateDmNews`関数のロジックを移行

**受け入れ基準**:
- [ ] `server/internal/service/generate_sample_service.go`が作成されている
- [ ] `GenerateSampleService`構造体が定義されている
- [ ] `GenerateSampleServiceInterface`が定義されている
- [ ] `GenerateDmUsers(ctx context.Context, totalCount int) ([]string, error)`メソッドが実装されている
- [ ] `GenerateDmPosts(ctx context.Context, dmUserIDs []string, totalCount int) error`メソッドが実装されている
- [ ] `GenerateDmNews(ctx context.Context, totalCount int) error`メソッドが実装されている
- [ ] `gofakeit`を使用したデータ生成ロジックが実装されている
- [ ] UUID生成、テーブル番号計算、データ生成、repository層への呼び出しが実装されている
- [ ] 既存の`generateDmUsers`、`generateDmPosts`、`generateDmNews`関数のロジックが移行されている

- _Requirements: 3.2.1, 3.2.2, 6.2_
- _Design: 3.2.1_

---

#### タスク 2.2: generate_sample_service_test.goの実装
**目的**: Service層の単体テストを実装する。

**作業内容**:
- `server/internal/service/generate_sample_service_test.go`を作成
- `GenerateDmUsers`、`GenerateDmPosts`、`GenerateDmNews`のテストケースを実装

**実装内容**:
- ファイルパス: `server/internal/service/generate_sample_service_test.go`
- パッケージ名: `service`
- モック実装:
  - `MockDmUserRepository`
  - `MockDmPostRepository`
  - `MockDmNewsRepository`
- テストケース:
  1. `GenerateDmUsers`の正常系テスト
  2. `GenerateDmPosts`の正常系テスト
  3. `GenerateDmNews`の正常系テスト
  4. repository層のエラーハンドリングのテスト
  5. UUID生成エラーのテスト
  6. テーブル番号計算エラーのテスト
- テーブル駆動テストを使用
- `github.com/stretchr/testify/assert`を使用

**受け入れ基準**:
- [ ] `server/internal/service/generate_sample_service_test.go`が作成されている
- [ ] モック（`MockDmUserRepository`、`MockDmPostRepository`、`MockDmNewsRepository`）が実装されている
- [ ] 正常系のテストケースが実装されている
- [ ] 異常系のテストケースが実装されている
- [ ] 全てのテストケースが通過する

- _Requirements: 6.2, 6.7_
- _Design: 3.2.2_

---

### Phase 3: CLI用usecase層の作成

#### タスク 3.1: generate_sample_usecase.goの実装
**目的**: CLI用のサンプルデータ生成usecaseを実装する。

**作業内容**:
- `server/internal/usecase/cli/generate_sample_usecase.go`を作成
- `GenerateSampleUsecase`構造体を定義
- `GenerateSampleData`メソッドを実装

**実装内容**:
- ファイルパス: `server/internal/usecase/cli/generate_sample_usecase.go`
- パッケージ名: `cli`
- 構造体定義:
  ```go
  type GenerateSampleUsecase struct {
      generateSampleService service.GenerateSampleServiceInterface
  }
  ```
- コンストラクタ:
  ```go
  func NewGenerateSampleUsecase(generateSampleService service.GenerateSampleServiceInterface) *GenerateSampleUsecase
  ```
- メソッド定義:
  ```go
  func (u *GenerateSampleUsecase) GenerateSampleData(ctx context.Context, totalCount int) error
  ```
- 実装内容:
  - service層の`GenerateDmUsers()`、`GenerateDmPosts()`、`GenerateDmNews()`を順次呼び出し
  - エラーハンドリングを実装
  - service層から返されたエラーをそのまま返す（エラーのラップは不要）

**受け入れ基準**:
- [ ] `server/internal/usecase/cli/generate_sample_usecase.go`が作成されている
- [ ] `GenerateSampleUsecase`構造体が定義されている
- [ ] `GenerateSampleServiceInterface`を依存として注入している
- [ ] `GenerateSampleData(ctx context.Context, totalCount int) error`メソッドが実装されている
- [ ] service層の`GenerateDmUsers()`、`GenerateDmPosts()`、`GenerateDmNews()`を順次呼び出している
- [ ] エラーハンドリングが適切に実装されている

- _Requirements: 3.3.1, 3.3.2, 6.3_
- _Design: 3.3.2_

---

#### タスク 3.2: generate_sample_usecase_test.goの実装
**目的**: CLI用usecase層の単体テストを実装する。

**作業内容**:
- `server/internal/usecase/cli/generate_sample_usecase_test.go`を作成
- `GenerateSampleData`のテストケースを実装

**実装内容**:
- ファイルパス: `server/internal/usecase/cli/generate_sample_usecase_test.go`
- パッケージ名: `cli`
- モック実装:
  ```go
  type MockGenerateSampleServiceInterface struct {
      GenerateDmUsersFunc func(ctx context.Context, totalCount int) ([]string, error)
      GenerateDmPostsFunc func(ctx context.Context, dmUserIDs []string, totalCount int) error
      GenerateDmNewsFunc  func(ctx context.Context, totalCount int) error
  }
  ```
- テストケース:
  1. 正常系: サンプルデータが正常に生成される場合
  2. 異常系: `GenerateDmUsers`でエラーが発生した場合
  3. 異常系: `GenerateDmPosts`でエラーが発生した場合
  4. 異常系: `GenerateDmNews`でエラーが発生した場合
- テーブル駆動テストを使用
- `github.com/stretchr/testify/assert`を使用

**受け入れ基準**:
- [ ] `server/internal/usecase/cli/generate_sample_usecase_test.go`が作成されている
- [ ] モック（`MockGenerateSampleServiceInterface`）が実装されている
- [ ] 正常系のテストケースが実装されている
- [ ] 異常系のテストケースが実装されている
- [ ] 全てのテストケースが通過する

- _Requirements: 6.3, 6.7_
- _Design: 3.3.3_

---

### Phase 4: main.goの簡素化

#### タスク 4.1: main.goの修正
**目的**: main.goをエントリーポイントと入出力制御のみに限定する。

**作業内容**:
- `server/cmd/generate-sample-data/main.go`を修正
- `internal/repository`パッケージをインポート
- `internal/service`パッケージをインポート
- `internal/usecase/cli`パッケージをインポート
- Repository層の初期化を追加
- Service層の初期化を追加
- Usecase層の初期化を追加
- サンプルデータ生成処理をusecase層の呼び出しに変更
- 既存の出力形式を維持

**実装内容**:
- 修正対象: `server/cmd/generate-sample-data/main.go`
- インポート追加:
  ```go
  "github.com/taku-o/go-webdb-template/internal/repository"
  "github.com/taku-o/go-webdb-template/internal/service"
  "github.com/taku-o/go-webdb-template/internal/usecase/cli"
  ```
- 初期化の追加:
  ```go
  // Repository層の初期化
  dmUserRepository := repository.NewDmUserRepository(groupManager)
  dmPostRepository := repository.NewDmPostRepository(groupManager)
  dmNewsRepository := repository.NewDmNewsRepository(groupManager)

  // Service層の初期化
  tableSelector := db.NewTableSelector(db.DBShardingTableCount, db.DBShardingTablesPerDB)
  generateSampleService := service.NewGenerateSampleService(
      dmUserRepository,
      dmPostRepository,
      dmNewsRepository,
      tableSelector,
  )

  // Usecase層の初期化
  generateSampleUsecase := cli.NewGenerateSampleUsecase(generateSampleService)
  ```
- サンプルデータ生成処理の変更:
  - 修正前: `generateDmUsers`、`generateDmPosts`、`generateDmNews`関数を直接呼び出し
  - 修正後: `usecase.GenerateSampleData(ctx, totalCount)`を呼び出し
- 削除する内容:
  - `gofakeit`の直接使用（service層に移動）
  - `idgen`の直接使用（service層に移動）
  - データ生成処理の直接実装（service層に移動）
  - バッチ挿入処理の直接実装（repository層に移動）
  - `generateDmUsers`、`generateDmPosts`、`generateDmNews`関数（service層に移動）
  - `insertDmUsersBatch`、`insertDmPostsBatch`、`insertDmNewsBatch`関数（repository層に移動）
- 定数の保持:
  - `totalCount = 100`をmain.goに保持（またはコマンドライン引数として受け取る）

**受け入れ基準**:
- [ ] `internal/repository`パッケージがインポートされている
- [ ] `internal/service`パッケージがインポートされている
- [ ] `internal/usecase/cli`パッケージがインポートされている
- [ ] Repository層が適切に初期化されている
- [ ] Service層が適切に初期化されている
- [ ] Usecase層が適切に初期化されている
- [ ] サンプルデータ生成処理がusecase層の呼び出しに変更されている
- [ ] 既存の出力形式（標準出力にログ出力）が維持されている
- [ ] 既存のエラーハンドリング（`log.Fatalf`）が維持されている
- [ ] `gofakeit`、`idgen`の直接使用が削除されている
- [ ] データ生成処理の直接実装が削除されている
- [ ] バッチ挿入処理の直接実装が削除されている

- _Requirements: 3.4.1, 3.5.1, 6.4, 6.5_
- _Design: 3.4.1_

---

### Phase 5: ビルドと動作確認

#### タスク 5.1: ビルドの確認
**目的**: ビルドが正常に完了し、バイナリが正しい場所に出力されることを確認する。

**作業内容**:
- `server/bin/`ディレクトリが存在することを確認（存在しない場合は作成）
- ビルドコマンドを実行
- バイナリが`server/bin/generate-sample-data`に出力されることを確認

**実装内容**:
- ビルドコマンド:
  ```bash
  cd server
  go build -o bin/generate-sample-data ./cmd/generate-sample-data
  ```
- ビルド出力先: `server/bin/generate-sample-data`

**受け入れ基準**:
- [ ] ビルドが正常に完了する
- [ ] バイナリが`server/bin/generate-sample-data`に出力される
- [ ] ビルドエラーが発生しない

- _Requirements: 6.4_
- _Design: 3.4.2_

---

#### タスク 5.2: 動作確認
**目的**: CLIコマンドが正常に動作し、既存の動作を維持していることを確認する。

**作業内容**:
- ローカル環境でCLIコマンドを実行
- 既存の出力形式（ログ出力）が維持されていることを確認
- 既存のエラーメッセージが維持されていることを確認
- dm_usersテーブルに100件のデータが生成されることを確認
- dm_postsテーブルに100件のデータが生成されることを確認
- dm_newsテーブルに100件のデータが生成されることを確認

**実装内容**:
- 実行コマンド:
  ```bash
  cd server
  APP_ENV=develop ./bin/generate-sample-data
  ```
- 確認項目:
  - "Starting sample data generation..."が出力される
  - "Generated X dm_users in dm_users_XXX"が出力される
  - "Generated X dm_posts in dm_posts_XXX"が出力される
  - "Generated 100 dm_news articles"が出力される
  - "Sample data generation completed successfully"が出力される
  - データベースにデータが正しく生成される

**受け入れ基準**:
- [ ] ローカル環境でCLIコマンドが正常に動作する
- [ ] 既存の出力形式（ログ出力）が維持されている
- [ ] 既存のエラーメッセージが維持されている
- [ ] dm_usersテーブルに100件のデータが生成される
- [ ] dm_postsテーブルに100件のデータが生成される
- [ ] dm_newsテーブルに100件のデータが生成される

- _Requirements: 6.6_

---

#### タスク 5.3: 既存テストの確認
**目的**: 既存のテストが全て通過することを確認する。

**作業内容**:
- 既存のテストを実行
- 全てのテストが通過することを確認

**実装内容**:
- テストコマンド:
  ```bash
  cd server
  go test ./...
  ```

**受け入れ基準**:
- [ ] 既存のテストが全て通過する
- [ ] テストエラーが発生しない

- _Requirements: 6.6, 6.7_

---

### Phase 6: ドキュメント更新

#### タスク 6.1: Architecture.mdの更新
**目的**: アーキテクチャドキュメントにCLIコマンドのレイヤー構造を追加する。

**作業内容**:
- `docs/Architecture.md`を修正
- CLIコマンドのレイヤー構造を追加（usecase層を含む）
- CLIコマンドのアーキテクチャ図を更新
- CLI用usecase層の説明を追加

**実装内容**:
- 修正対象: `docs/Architecture.md`
- 追加内容:
  - CLIコマンドのレイヤー構造の説明
  - CLIコマンドのアーキテクチャ図（usecase層を含む）
  - CLI用usecase層の説明

**受け入れ基準**:
- [ ] `docs/Architecture.md`にCLIコマンドのレイヤー構造が追加されている
- [ ] CLIコマンドのアーキテクチャ図が更新されている（usecase層を含む）
- [ ] CLI用usecase層の説明が追加されている

- _Requirements: 3.7.1, 6.8_

---

#### タスク 6.2: Project-Structure.mdの更新
**目的**: プロジェクト構造ドキュメントに新規作成するファイルを追加する。

**作業内容**:
- `docs/Project-Structure.md`を修正
- `server/internal/usecase/cli/generate_sample_usecase.go`を追加
- `server/internal/service/generate_sample_service.go`を追加
- `server/internal/repository/dm_news_repository.go`を追加（新規作成）

**実装内容**:
- 修正対象: `docs/Project-Structure.md`
- 追加内容:
  - `server/internal/usecase/cli/generate_sample_usecase.go`の説明
  - `server/internal/service/generate_sample_service.go`の説明
  - `server/internal/repository/dm_news_repository.go`の説明

**受け入れ基準**:
- [ ] `docs/Project-Structure.md`に新規作成するファイルが追加されている
- [ ] `server/internal/usecase/cli/generate_sample_usecase.go`が追加されている
- [ ] `server/internal/service/generate_sample_service.go`が追加されている
- [ ] `server/internal/repository/dm_news_repository.go`が追加されている

- _Requirements: 3.7.2, 6.8_

---

#### タスク 6.3: Generate-Sample-Data.mdの更新（存在する場合）
**目的**: CLIツールドキュメントのアーキテクチャ図を更新してusecase層を含める。

**作業内容**:
- `docs/Generate-Sample-Data.md`が存在する場合は修正
- CLIコマンドのアーキテクチャ図を更新（usecase層を追加）
- レイヤー構造の説明を更新（main.go → usecase → service → repository → db → 出力）

**実装内容**:
- 修正対象: `docs/Generate-Sample-Data.md`（存在する場合）
- 更新内容:
  - CLIコマンドのアーキテクチャ図を更新（usecase層を追加）
  - レイヤー構造の説明を更新

**受け入れ基準**:
- [ ] `docs/Generate-Sample-Data.md`が存在する場合は、アーキテクチャ図が更新されている
- [ ] レイヤー構造の説明が更新されている（usecase層を含む）

- _Requirements: 3.7.3, 6.8_

---

#### タスク 6.4: Command-Line-Tool.mdの更新（存在する場合）
**目的**: CLIツールドキュメントのアーキテクチャ図を更新してusecase層を含める。

**作業内容**:
- `docs/Command-Line-Tool.md`が存在する場合は修正
- CLIコマンドのアーキテクチャ図を更新（usecase層を追加）
- レイヤー構造の説明を更新（main.go → usecase → service → repository → db → 出力）

**実装内容**:
- 修正対象: `docs/Command-Line-Tool.md`（存在する場合）
- 更新内容:
  - CLIコマンドのアーキテクチャ図を更新（usecase層を追加）
  - レイヤー構造の説明を更新

**受け入れ基準**:
- [ ] `docs/Command-Line-Tool.md`が存在する場合は、アーキテクチャ図が更新されている
- [ ] レイヤー構造の説明が更新されている（usecase層を含む）

- _Requirements: 3.7.3, 6.8_

---

#### タスク 6.5: structure.mdの更新
**目的**: ファイル組織ドキュメントに新規作成するファイルを追加する。

**作業内容**:
- `.kiro/steering/structure.md`を修正
- `server/internal/usecase/cli/generate_sample_usecase.go`を追加
- `server/internal/service/generate_sample_service.go`を追加
- `server/internal/repository/dm_news_repository.go`を追加（新規作成）

**実装内容**:
- 修正対象: `.kiro/steering/structure.md`
- 追加内容:
  - `server/internal/usecase/cli/generate_sample_usecase.go`の説明
  - `server/internal/service/generate_sample_service.go`の説明
  - `server/internal/repository/dm_news_repository.go`の説明

**受け入れ基準**:
- [ ] `.kiro/steering/structure.md`に新規作成するファイルが追加されている
- [ ] `server/internal/usecase/cli/generate_sample_usecase.go`が追加されている
- [ ] `server/internal/service/generate_sample_service.go`が追加されている
- [ ] `server/internal/repository/dm_news_repository.go`が追加されている

- _Requirements: 3.7.4, 6.8_

---

## タスクの依存関係

```
Phase 1: Repository層の拡張
  ├─ タスク 1.1: dm_user_repository.goへのInsertDmUsersBatchメソッドの追加
  ├─ タスク 1.2: dm_post_repository.goへのInsertDmPostsBatchメソッドの追加
  ├─ タスク 1.3: dm_news_repository.goの新規作成
  ├─ タスク 1.4: dm_user_repository_test.goへのテスト追加（存在する場合）
  ├─ タスク 1.5: dm_post_repository_test.goへのテスト追加（存在する場合）
  └─ タスク 1.6: dm_news_repository_test.goの新規作成

Phase 2: Service層の作成
  ├─ タスク 2.1: generate_sample_service.goの実装（Phase 1に依存）
  └─ タスク 2.2: generate_sample_service_test.goの実装（タスク 2.1に依存）

Phase 3: CLI用usecase層の作成
  ├─ タスク 3.1: generate_sample_usecase.goの実装（Phase 2に依存）
  └─ タスク 3.2: generate_sample_usecase_test.goの実装（タスク 3.1に依存）

Phase 4: main.goの簡素化
  └─ タスク 4.1: main.goの修正（Phase 3に依存）

Phase 5: ビルドと動作確認
  ├─ タスク 5.1: ビルドの確認（Phase 4に依存）
  ├─ タスク 5.2: 動作確認（タスク 5.1に依存）
  └─ タスク 5.3: 既存テストの確認（Phase 1-3に依存）

Phase 6: ドキュメント更新
  ├─ タスク 6.1: Architecture.mdの更新（Phase 4に依存）
  ├─ タスク 6.2: Project-Structure.mdの更新（Phase 1-3に依存）
  ├─ タスク 6.3: Generate-Sample-Data.mdの更新（存在する場合）（Phase 4に依存）
  ├─ タスク 6.4: Command-Line-Tool.mdの更新（存在する場合）（Phase 4に依存）
  └─ タスク 6.5: structure.mdの更新（Phase 1-3に依存）
```

## 実装順序の推奨

1. **Phase 1**: Repository層の拡張（既存ファイルへの追加、新規ファイルの作成）
2. **Phase 2**: Service層の作成（Repository層に依存）
3. **Phase 3**: CLI用usecase層の作成（Service層に依存）
4. **Phase 4**: main.goの簡素化（全レイヤーに依存）
5. **Phase 5**: ビルドと動作確認（実装完了後）
6. **Phase 6**: ドキュメント更新（実装完了後）

## 注意事項

- 各タスクは独立して実装可能な粒度に分解されている
- テストは実装と同時に作成することを推奨
- 既存のコードスタイルに合わせて実装する
- 既存のエラーハンドリングパターンに従う
- 既存のテストパターンに従う
