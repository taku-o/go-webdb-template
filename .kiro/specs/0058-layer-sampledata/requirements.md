# server/cmd/generate-sample-dataの構造修正の要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Feature名**: 0058-layer-sampledata
- **作成日**: 2026-01-27
- **関連Issue**: https://github.com/taku-o/go-webdb-template/issues/120

### 1.2 目的
`server/cmd/generate-sample-data`の実装を、APIサーバーと同じレイヤー構造（usecase -> service -> repository -> db）に変更する。これにより、CLIコマンドとAPIサーバーで一貫したアーキテクチャを実現し、サンプルデータ生成処理を共通化してコードの保守性と再利用性を向上させる。

### 1.3 スコープ
- CLIコマンドのバリデーションと入出力制御をCLI層に集約
- `server/internal/usecase/cli`ディレクトリに`generate-sample-data`用のusecaseを追加
- `server/internal/service`ディレクトリにサンプルデータ生成用のservice層を追加
- 既存のrepositoryファイルにバッチ挿入メソッドを追加（`dm_user_repository.go`、`dm_post_repository.go`、`dm_news_repository.go`）
- main.goの簡素化（エントリーポイントと入出力制御のみ）
- 関連ドキュメントの更新（アーキテクチャ、プロジェクト構造、CLIツールドキュメント）

**本実装の範囲外**:
- 既存のAPIサーバーのusecase層への影響（APIサーバーは変更しない）
- 他のCLIコマンドへの影響（本実装は`generate-sample-data`のみを対象）
- サンプルデータ生成のロジック変更（既存の生成ロジックを維持）

## 2. 背景・現状分析

### 2.1 現在の状況
- **CLI実装**: `server/cmd/generate-sample-data/main.go`に直接実装されている
- **レイヤー構造**: main.goにサンプルデータ生成処理が直接実装されている
- **バリデーション**: なし（引数なし）
- **入出力制御**: main.go内で標準出力に直接出力
- **ビジネスロジック**: main.go内に直接実装（`gofakeit`を使用したデータ生成、バッチ挿入処理）
- **usecase層**: CLI用のusecase層が存在しない（`list_dm_users_usecase.go`は存在するが、`generate-sample-data`用は存在しない）
- **service層**: サンプルデータ生成用のservice層が存在しない
- **repository層**: 既存のrepository層（`dm_user_repository.go`、`dm_post_repository.go`）は存在するが、バッチ挿入メソッドが存在しない。`dm_news_repository.go`は存在しない

### 2.2 現在の処理内容
1. **dm_usersテーブルへのデータ生成**:
   - 32分割テーブル（dm_users_000〜dm_users_031）に分散
   - UUIDに基づいてテーブル番号を計算
   - バッチサイズ500件ずつ挿入
   - 合計100件を生成

2. **dm_postsテーブルへのデータ生成**:
   - 32分割テーブル（dm_posts_000〜dm_posts_031）に分散
   - user_idに基づいてテーブル番号を計算
   - バッチサイズ500件ずつ挿入
   - 合計100件を生成

3. **dm_newsテーブルへのデータ生成**:
   - masterデータベースの固定テーブル（dm_news）に挿入
   - バッチサイズ500件ずつ挿入
   - 合計100件を生成

### 2.3 課題点
1. **レイヤー構造の不一致**: APIサーバーはusecase層を使用しているが、CLIは処理を直接実装している
2. **責務の混在**: main.goにサンプルデータ生成処理が直接実装されている
3. **再利用性の低さ**: サンプルデータ生成処理がmain.goに直接実装されており、他のコンポーネントから再利用できない
4. **テストの困難さ**: main.goにロジックが集中しているため、単体テストが困難
5. **コードの重複**: バッチ挿入処理が各テーブルごとに重複実装されている

### 2.4 本実装による改善点
1. **アーキテクチャの一貫性**: APIサーバーとCLIで同じレイヤー構造を採用
2. **責務の明確化**: CLI層（main.go）はエントリーポイントと入出力制御のみを担当
3. **再利用性の向上**: サンプルデータ生成処理を各レイヤーに分離することで、他のコンポーネントからも再利用可能
4. **テスト容易性**: usecase層、service層、repository層を独立してテストできるようになる
5. **コードの整理**: バッチ挿入処理をrepository層に集約することで、コードの重複を削減

## 3. 機能要件

### 3.1 Repository層の拡張

#### 3.1.1 既存のrepositoryファイルへのメソッド追加
- **目的**: 既存のrepositoryファイルにバッチ挿入メソッドを追加
- **実装内容**:
  - `server/internal/repository/dm_user_repository.go`に`InsertDmUsersBatch(ctx context.Context, tableName string, dmUsers []*model.DmUser) error`メソッドを追加
  - `server/internal/repository/dm_post_repository.go`に`InsertDmPostsBatch(ctx context.Context, tableName string, dmPosts []*model.DmPost) error`メソッドを追加
  - `server/internal/repository/dm_news_repository.go`を新規作成し、`InsertDmNewsBatch(ctx context.Context, dmNews []*model.DmNews) error`メソッドを実装
  - バッチサイズ500件ずつ挿入する処理を実装
  - 既存の`insertDmUsersBatch`、`insertDmPostsBatch`、`insertDmNewsBatch`関数のロジックを移行

#### 3.1.2 テストファイルの更新
- **目的**: 既存のrepositoryテストファイルにバッチ挿入メソッドのテストを追加
- **実装内容**:
  - `server/internal/repository/dm_user_repository_test.go`に`InsertDmUsersBatch`のテストを追加（存在する場合）
  - `server/internal/repository/dm_post_repository_test.go`に`InsertDmPostsBatch`のテストを追加（存在する場合）
  - `server/internal/repository/dm_news_repository_test.go`を新規作成し、`InsertDmNewsBatch`のテストを実装

### 3.2 Service層の作成

#### 3.2.1 ディレクトリ構造
- **目的**: サンプルデータ生成用のservice層を既存の`server/internal/service`ディレクトリに追加
- **実装内容**:
  - `server/internal/service/generate_sample_service.go`を作成
  - `server/internal/service/generate_sample_service_test.go`を作成（テストファイル）

#### 3.2.2 GenerateSampleServiceの実装
- **目的**: サンプルデータ生成用のservice層を実装
- **実装内容**:
  - `GenerateSampleService`構造体を定義
  - `DmUserRepository`、`DmPostRepository`、`DmNewsRepository`を依存として注入
  - `TableSelector`を依存として注入（テーブル番号計算用）
  - `GenerateDmUsers(ctx context.Context, totalCount int) ([]string, error)`メソッドを実装
  - `GenerateDmPosts(ctx context.Context, dmUserIDs []string, totalCount int) error`メソッドを実装
  - `GenerateDmNews(ctx context.Context, totalCount int) error`メソッドを実装
  - `gofakeit`を使用したデータ生成ロジックを実装
  - UUID生成、テーブル番号計算、データ生成、repository層への呼び出しを実装
  - 既存の`generateDmUsers`、`generateDmPosts`、`generateDmNews`関数のロジックを移行

### 3.3 CLI用usecase層の作成

#### 3.3.1 ディレクトリ構造
- **目的**: CLI用のusecase層を既存の`server/internal/usecase/cli`ディレクトリに追加
- **実装内容**:
  - `server/internal/usecase/cli/generate_sample_usecase.go`を作成
  - `server/internal/usecase/cli/generate_sample_usecase_test.go`を作成（テストファイル）

#### 3.3.2 GenerateSampleUsecaseの実装
- **目的**: CLI用のサンプルデータ生成usecaseを実装
- **実装内容**:
  - `GenerateSampleUsecase`構造体を定義
  - `GenerateSampleServiceInterface`を依存として注入
  - `GenerateSampleData(ctx context.Context, totalCount int) error`メソッドを実装
  - service層の`GenerateDmUsers()`、`GenerateDmPosts()`、`GenerateDmNews()`を順次呼び出し
  - エラーハンドリングを実装
- **インターフェース**: `GenerateSampleServiceInterface`を新規作成（`server/internal/service/generate_sample_service.go`で定義）

### 3.4 main.goの簡素化

#### 3.4.1 エントリーポイントの責務
- **目的**: main.goをエントリーポイントと入出力制御のみに限定
- **実装内容**:
  - 設定ファイルの読み込み（`config.Load()`）
  - `GroupManager`の初期化
  - データベース接続確認（`groupManager.PingAll()`）
  - usecase層の初期化と実行
  - 結果の出力（標準出力にログ出力）

#### 3.4.2 出力関数
- **目的**: サンプルデータ生成の進捗と結果をmain.goに保持
- **実装内容**:
  - usecase層から取得した結果を標準出力に出力
  - 既存のログ出力形式を維持（"Starting sample data generation..."、"Sample data generation completed successfully"など）

### 3.5 依存関係の注入

#### 3.5.1 usecase層の初期化
- **目的**: usecase層を適切に初期化してmain.goで使用
- **実装内容**:
  - Repository層の初期化（`repository.NewDmUserRepository(groupManager)`、`repository.NewDmPostRepository(groupManager)`、`repository.NewDmNewsRepository(groupManager)`）
  - Service層の初期化（`service.NewGenerateSampleService(dmUserRepository, dmPostRepository, dmNewsRepository, tableSelector)`）
  - Usecase層の初期化（`cli.NewGenerateSampleUsecase(generateSampleService)`）
  - usecase層の`GenerateSampleData()`メソッドを呼び出し

### 3.6 定数の定義

#### 3.6.1 定数の移動
- **目的**: 定数を適切な場所に定義
- **実装内容**:
  - `batchSize = 500`をrepository層に移動（またはservice層に移動）
  - `tableCount = 32`をservice層に移動（または定数として定義）
  - `totalCount = 100`をmain.goに保持（またはコマンドライン引数として受け取る）

### 3.7 ドキュメントの更新

#### 3.7.1 アーキテクチャドキュメントの更新
- **目的**: CLIコマンドのレイヤー構造をアーキテクチャドキュメントに反映
- **修正対象**: `docs/Architecture.md`
- **実装内容**:
  - CLIコマンドのレイヤー構造を追加（usecase層を含む）
  - CLIコマンドのアーキテクチャ図を更新
  - CLI用usecase層の説明を追加

#### 3.7.2 プロジェクト構造ドキュメントの更新
- **目的**: 新規作成するファイルをプロジェクト構造に反映
- **修正対象**: `docs/Project-Structure.md`
- **実装内容**:
  - `server/internal/usecase/cli/generate_sample_usecase.go`を追加
  - `server/internal/service/generate_sample_service.go`を追加
  - `server/internal/repository/dm_news_repository.go`を追加（新規作成）

#### 3.7.3 CLIツールドキュメントの更新
- **目的**: CLIコマンドのアーキテクチャ図を更新してusecase層を含める
- **修正対象**: `docs/Generate-Sample-Data.md`、`docs/Command-Line-Tool.md`（存在する場合）
- **実装内容**:
  - CLIコマンドのアーキテクチャ図を更新（usecase層を追加）
  - レイヤー構造の説明を更新（main.go → usecase → service → repository → db → 出力）

#### 3.7.4 ファイル組織ドキュメントの更新
- **目的**: 新規作成するファイルをファイル組織に反映
- **修正対象**: `.kiro/steering/structure.md`
- **実装内容**:
  - `server/internal/usecase/cli/generate_sample_usecase.go`を追加
  - `server/internal/service/generate_sample_service.go`を追加
  - `server/internal/repository/dm_news_repository.go`を追加（新規作成）

## 4. 非機能要件

### 4.1 パフォーマンス
- **既存機能の維持**: 既存のサンプルデータ生成処理のパフォーマンスを維持
- **オーバーヘッド**: usecase層、service層、repository層の追加によるパフォーマンスオーバーヘッドは無視できるレベル
- **バッチサイズ**: 既存のバッチサイズ（500件）を維持

### 4.2 信頼性
- **エラーハンドリング**: 既存のエラーハンドリングを維持
- **後方互換性**: 既存のCLIコマンドの動作を維持（出力形式、エラーメッセージなど）
- **データ整合性**: 既存のデータ生成ロジックを維持（UUID生成、テーブル番号計算など）

### 4.3 保守性
- **コードの可読性**: レイヤー構造が明確で、各レイヤーの責務が明確
- **一貫性**: APIサーバーと同じレイヤー構造を採用することで、コードベース全体の一貫性を向上
- **テスト容易性**: usecase層、service層、repository層を独立してテストできる
- **再利用性**: サンプルデータ生成処理を各レイヤーに分離することで、他のコンポーネントからも再利用可能

### 4.4 互換性
- **既存機能**: 既存のAPIサーバーに影響を与えない
- **CLIコマンドの動作**: 既存のCLIコマンドの動作（出力形式、エラーメッセージ）を維持
- **データ生成ロジック**: 既存のデータ生成ロジック（gofakeitの使用、バッチサイズ、テーブル分割など）を維持

## 5. 制約事項

### 5.1 技術的制約
- **既存のデータ生成ロジック**: 既存のデータ生成ロジック（gofakeitの使用、UUID生成、テーブル番号計算）を変更しない
- **バッチサイズ**: 既存のバッチサイズ（500件）を維持
- **テーブル分割**: 既存のテーブル分割（32テーブル）を維持
- **生成件数**: 既存の生成件数（合計100件）を維持（またはコマンドライン引数として受け取る）

### 5.2 実装上の制約
- **ディレクトリ構造**: 既存のディレクトリ構造に従う（`server/internal/usecase/cli`、`server/internal/service`、`server/internal/repository`）
- **命名規則**: 既存の命名規則に従う（`GenerateSampleUsecase`、`GenerateSampleService`）
- **インターフェース**: `GenerateSampleServiceInterface`を新規作成。既存のrepository（`DmUserRepository`、`DmPostRepository`、`DmNewsRepository`）を直接使用

### 5.3 動作環境
- **ローカル環境**: ローカル環境でCLIコマンドが正常に動作することを確認
- **CI環境**: CI環境でもCLIコマンドが正常に動作することを確認（該当する場合）
- **データベース**: PostgreSQL（master 1台 + sharding 4台）が正常に動作していることを前提

## 6. 受け入れ基準

### 6.1 Repository層の拡張
- [ ] `server/internal/repository/dm_user_repository.go`に`InsertDmUsersBatch(ctx context.Context, tableName string, dmUsers []*model.DmUser) error`メソッドが追加されている
- [ ] `server/internal/repository/dm_post_repository.go`に`InsertDmPostsBatch(ctx context.Context, tableName string, dmPosts []*model.DmPost) error`メソッドが追加されている
- [ ] `server/internal/repository/dm_news_repository.go`が作成されている
- [ ] `server/internal/repository/dm_news_repository.go`に`InsertDmNewsBatch(ctx context.Context, dmNews []*model.DmNews) error`メソッドが実装されている
- [ ] バッチサイズ500件ずつ挿入する処理が実装されている
- [ ] `server/internal/repository/dm_user_repository_test.go`に`InsertDmUsersBatch`のテストが追加されている（存在する場合）
- [ ] `server/internal/repository/dm_post_repository_test.go`に`InsertDmPostsBatch`のテストが追加されている（存在する場合）
- [ ] `server/internal/repository/dm_news_repository_test.go`が作成されている
- [ ] `server/internal/repository/dm_news_repository_test.go`に`InsertDmNewsBatch`のテストが実装されている

### 6.2 Service層の作成
- [ ] `server/internal/service/generate_sample_service.go`が作成されている
- [ ] `GenerateSampleService`構造体が定義されている
- [ ] `GenerateSampleServiceInterface`が定義されている
- [ ] `GenerateDmUsers(ctx context.Context, totalCount int) ([]string, error)`メソッドが実装されている
- [ ] `GenerateDmPosts(ctx context.Context, dmUserIDs []string, totalCount int) error`メソッドが実装されている
- [ ] `GenerateDmNews(ctx context.Context, totalCount int) error`メソッドが実装されている
- [ ] `gofakeit`を使用したデータ生成ロジックが実装されている
- [ ] UUID生成、テーブル番号計算、データ生成、repository層への呼び出しが実装されている
- [ ] `server/internal/service/generate_sample_service_test.go`が作成されている
- [ ] 単体テストが実装されている

### 6.3 CLI用usecase層の作成
- [ ] `server/internal/usecase/cli/generate_sample_usecase.go`が作成されている
- [ ] `GenerateSampleUsecase`構造体が定義されている
- [ ] `GenerateSampleServiceInterface`を依存として注入している
- [ ] `GenerateSampleData(ctx context.Context, totalCount int) error`メソッドが実装されている
- [ ] service層の`GenerateDmUsers()`、`GenerateDmPosts()`、`GenerateDmNews()`を順次呼び出している
- [ ] `server/internal/usecase/cli/generate_sample_usecase_test.go`が作成されている
- [ ] 単体テストが実装されている

### 6.4 main.goの簡素化
- [ ] main.goがエントリーポイントと入出力制御のみを担当している
- [ ] usecase層を初期化して使用している
- [ ] 既存の出力形式（標準出力にログ出力）が維持されている
- [ ] 既存のエラーハンドリングが維持されている
- [ ] ビルド時のバイナリが`server/bin/generate-sample-data`に出力される

### 6.5 依存関係の注入
- [ ] `DmUserRepository`、`DmPostRepository`、`DmNewsRepository`が適切に初期化されている
- [ ] Service層が適切に初期化されている（`DmUserRepository`、`DmPostRepository`、`DmNewsRepository`を依存として注入）
- [ ] Usecase層が適切に初期化されている
- [ ] usecase層の`GenerateSampleData()`メソッドが呼び出されている

### 6.6 動作確認
- [ ] ローカル環境でCLIコマンドが正常に動作する
- [ ] 既存の出力形式（ログ出力）が維持されている
- [ ] 既存のエラーメッセージが維持されている
- [ ] dm_usersテーブルに100件のデータが生成される
- [ ] dm_postsテーブルに100件のデータが生成される
- [ ] dm_newsテーブルに100件のデータが生成される
- [ ] 既存のテストが全て通過する
- [ ] CI環境でCLIコマンドが正常に動作する（該当する場合）

### 6.7 テスト
- [ ] `server/internal/repository/dm_user_repository_test.go`に`InsertDmUsersBatch`のテストが追加されている（存在する場合）
- [ ] `server/internal/repository/dm_post_repository_test.go`に`InsertDmPostsBatch`のテストが追加されている（存在する場合）
- [ ] `server/internal/repository/dm_news_repository_test.go`が作成されている
- [ ] `server/internal/service/generate_sample_service_test.go`が作成されている
- [ ] `server/internal/usecase/cli/generate_sample_usecase_test.go`が作成されている
- [ ] 各層の単体テストが実装されている
- [ ] 既存のテストが全て通過する

### 6.8 ドキュメントの更新
- [ ] `docs/Architecture.md`にCLIコマンドのレイヤー構造が追加されている
- [ ] `docs/Architecture.md`のCLIコマンドのアーキテクチャ図が更新されている
- [ ] `docs/Project-Structure.md`に新規作成するファイルが追加されている
- [ ] `docs/Generate-Sample-Data.md`のアーキテクチャ図が更新されている（存在する場合）
- [ ] `docs/Generate-Sample-Data.md`のレイヤー構造の説明が更新されている（存在する場合）
- [ ] `docs/Command-Line-Tool.md`（存在する場合）のアーキテクチャ図が更新されている
- [ ] `docs/Command-Line-Tool.md`（存在する場合）のレイヤー構造の説明が更新されている
- [ ] `.kiro/steering/structure.md`に新規作成するファイルが追加されている

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### 新規作成が必要なファイル
- `server/internal/repository/dm_news_repository.go`: dm_news用のrepository層（新規作成）
- `server/internal/repository/dm_news_repository_test.go`: dm_news用repository層のテスト
- `server/internal/service/generate_sample_service.go`: サンプルデータ生成用のservice層
- `server/internal/service/generate_sample_service_test.go`: service層のテスト
- `server/internal/usecase/cli/generate_sample_usecase.go`: CLI用usecase層の実装
- `server/internal/usecase/cli/generate_sample_usecase_test.go`: CLI用usecase層のテスト

#### 修正が必要なファイル
- `server/internal/repository/dm_user_repository.go`: `InsertDmUsersBatch`メソッドを追加
- `server/internal/repository/dm_post_repository.go`: `InsertDmPostsBatch`メソッドを追加
- `server/internal/repository/dm_user_repository_test.go`: `InsertDmUsersBatch`のテストを追加（存在する場合）
- `server/internal/repository/dm_post_repository_test.go`: `InsertDmPostsBatch`のテストを追加（存在する場合）
- `server/cmd/generate-sample-data/main.go`: usecase層を使用するように修正
- `docs/Architecture.md`: CLIコマンドのレイヤー構造を追加
- `docs/Project-Structure.md`: 新規作成するファイルを追加
- `docs/Generate-Sample-Data.md`: CLIコマンドのアーキテクチャ図を更新（存在する場合）
- `docs/Command-Line-Tool.md`（存在する場合）: CLIコマンドのアーキテクチャ図を更新
- `.kiro/steering/structure.md`: 新規作成するファイルを追加

### 7.2 既存機能への影響
- **既存のAPIサーバー**: 影響なし（APIサーバーは変更しない）
- **既存のCLIコマンドの動作**: 影響なし（動作は維持される）
- **既存のデータ生成ロジック**: 影響なし（ロジックは維持される）

## 8. 実装上の注意事項

### 8.1 Repository層の実装
- **既存のrepositoryファイルへの追加**: `dm_user_repository.go`、`dm_post_repository.go`にバッチ挿入メソッドを追加
- **新規repositoryファイルの作成**: `dm_news_repository.go`を新規作成（master接続を使用）
- **バッチサイズ**: 既存のバッチサイズ（500件）を維持
- **エラーハンドリング**: 既存のエラーハンドリングを維持
- **GORMの使用**: 既存のGORMの使用方法を維持
- **動的テーブル名**: 既存の動的テーブル名の生成方法を維持（dm_users、dm_postsは動的テーブル名、dm_newsは固定テーブル名）

### 8.2 Service層の実装
- **依存関係の注入**: コンストラクタで依存関係を注入（`DmUserRepository`、`DmPostRepository`、`DmNewsRepository`、`TableSelector`）
- **エラーハンドリング**: repository層から返されたエラーをそのまま返す（エラーのラップは不要）
- **データ生成ロジック**: 既存の`gofakeit`を使用したデータ生成ロジックを維持
- **UUID生成**: 既存のUUID生成方法（`idgen.GenerateUUIDv7()`）を維持
- **テーブル番号計算**: 既存のテーブル番号計算方法（`tableSelector.GetTableNumberFromUUID()`）を維持

### 8.3 usecase層の実装
- **インターフェースの使用**: 新規作成する`GenerateSampleServiceInterface`を使用
- **依存関係の注入**: コンストラクタで`GenerateSampleServiceInterface`を注入
- **エラーハンドリング**: service層から返されたエラーをそのまま返す（エラーのラップは不要）
- **処理の順序**: `GenerateDmUsers()` → `GenerateDmPosts()` → `GenerateDmNews()`の順で実行

### 8.4 main.goの修正
- **usecase層の初期化**: Repository → Service → Usecaseの順で初期化
- **既存の出力**: 標準出力にログ出力（既存の動作を維持）
- **エラーハンドリング**: 既存のエラーハンドリングを維持
- **ビルド出力先**: ビルド時のバイナリは`server/bin/generate-sample-data`に出力する
- **設定ファイルの読み込み**: 既存の設定ファイル読み込み方法を維持
- **GroupManagerの初期化**: 既存のGroupManager初期化方法を維持

### 8.5 テストの実装
- **repository層のテスト**: 既存のrepositoryテストファイルにバッチ挿入メソッドのテストを追加。`dm_news_repository_test.go`を新規作成
- **service層のテスト**: `DmUserRepository`、`DmPostRepository`、`DmNewsRepository`のモックを使用してテスト
- **usecase層のテスト**: `GenerateSampleServiceInterface`のモックを使用してテスト

### 8.6 ドキュメントの更新
- **アーキテクチャドキュメント**: CLIコマンドのレイヤー構造を明確に記載
- **プロジェクト構造ドキュメント**: 新規作成するファイルを反映
- **CLIツールドキュメント**: アーキテクチャ図を更新してusecase層を含める
- **ファイル組織ドキュメント**: 新規作成するファイルを反映
- **一貫性**: 全てのドキュメントで同じレイヤー構造を記載

## 9. 参考情報

### 9.1 関連ドキュメント
- `docs/Architecture.md`: アーキテクチャドキュメント
- `docs/Project-Structure.md`: プロジェクト構造ドキュメント
- `docs/Generate-Sample-Data.md`: サンプルデータ生成機能のドキュメント
- `.kiro/steering/structure.md`: ファイル組織とコードパターン
- `.kiro/steering/tech.md`: 技術スタックとアーキテクチャ

### 9.2 既存実装の参考
- `server/internal/usecase/cli/list_dm_users_usecase.go`: 既存のCLI用usecase層の実装パターン
- `server/internal/usecase/cli/generate_secret_usecase.go`: 既存のCLI用usecase層の実装パターン
- `server/cmd/generate-sample-data/main.go`: 既存のCLI実装
- `server/internal/service/secret_service.go`: 既存のservice層の実装パターン
- `server/internal/repository/dm_user_repository.go`: 既存のrepository層の実装パターン
- `server/internal/repository/dm_post_repository.go`: 既存のrepository層の実装パターン

### 9.3 技術スタック
- **言語**: Go
- **アーキテクチャ**: レイヤードアーキテクチャ（usecase -> service -> repository -> db -> 出力）
- **テスト**: `testing`（標準ライブラリ）、`github.com/stretchr/testify`（アサーション、モック）
- **データ生成**: `github.com/brianvoe/gofakeit/v6`（ランダムデータ生成）
- **UUID生成**: `github.com/taku-o/go-webdb-template/internal/util/idgen`（UUIDv7生成）

### 9.4 レイヤー構造の比較

| 項目 | 現在（修正前） | 修正後 |
|------|---------------|--------|
| CLI層 | main.go（サンプルデータ生成処理が直接実装） | main.go（エントリーポイント、入出力） |
| Usecase層 | なし | `usecase/cli/GenerateSampleUsecase` |
| Service層 | なし | `service.GenerateSampleService` |
| Repository層 | なし | `repository.DmUserRepository`、`repository.DmPostRepository`、`repository.DmNewsRepository`（既存を拡張） |
| DB層 | main.go内で直接使用 | `db.GroupManager`（既存） |
| 出力 | main.go内で直接出力 | main.go内で出力（usecase層から取得） |

### 9.5 APIサーバーとの比較

| 項目 | APIサーバー | CLI（修正後） |
|------|------------|--------------|
| エントリーポイント | `server/cmd/server/main.go` | `server/cmd/generate-sample-data/main.go` |
| バリデーション | API Layer（Handler） | CLI層（main.go、現時点では不要） |
| Usecase層 | `usecase.DmUserUsecase` | `usecase/cli.GenerateSampleUsecase` |
| Service層 | `service.DmUserService` | `service.GenerateSampleService` |
| Repository層 | `repository.UserRepository` | `repository.DmUserRepository`、`repository.DmPostRepository`、`repository.DmNewsRepository` |
| DB層 | `db.GroupManager` | `db.GroupManager` |
| 出力 | HTTPレスポンス | 標準出力（ログ） |
