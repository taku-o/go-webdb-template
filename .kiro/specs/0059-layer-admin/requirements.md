# Adminアプリのカスタムページの実装の仕組みを変更するの要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Feature名**: 0059-layer-admin
- **作成日**: 2026-01-27
- **関連Issue**: https://github.com/taku-o/go-webdb-template/issues/122

### 1.2 目的
Adminアプリのカスタムページ（`dm_user_register.go`、`api_key.go`）の実装を、APIサーバーと同じレイヤー構造（pages -> usecase -> service -> repository -> db）に変更する。これにより、AdminアプリとAPIサーバーで一貫したアーキテクチャを実現し、コードの保守性と再利用性を向上させる。

### 1.3 スコープ
- Adminページのバリデーションと入出力制御をpages層に集約
- `server/internal/usecase/admin`ディレクトリにAdminアプリ用のusecaseを新規作成
- `server/internal/admin/pages/dm_user_register.go`を、pages -> usecase -> service -> repository -> db の流れで処理するように修正
- `server/internal/admin/pages/api_key.go`を、pages -> usecase -> service -> 鍵の生成、usecase -> service -> ペイロードのデコードという2つの処理に分ける
- 関連ドキュメントの更新（アーキテクチャ、プロジェクト構造）

**本実装の範囲外**:
- 既存のAPIサーバーのusecase層への影響（APIサーバーは変更しない）
- 他のAdminページへの影響（本実装は`dm_user_register.go`と`api_key.go`のみを対象）
- 既存のビジネスロジックの変更（既存のロジックを維持）

## 2. 背景・現状分析

### 2.1 現在の状況

#### 2.1.1 dm_user_register.goの現状
- **実装場所**: `server/internal/admin/pages/dm_user_register.go`に直接実装されている
- **レイヤー構造**: pages層にバリデーション、DB操作が直接実装されている
- **バリデーション**: pages層内で`validateDmUserInput`関数として実装
- **入出力制御**: pages層内で`renderDmUserRegisterForm`関数として実装
- **ビジネスロジック**: pages層内に直接実装（`checkEmailExistsSharded`、`insertDmUserSharded`）
- **usecase層**: Admin用のusecase層が存在しない（`server/internal/usecase/admin`ディレクトリが存在しない）
- **service層**: Admin用のservice層が存在しない（既存の`dm_user_service.go`は存在するが、Admin専用ではない）
- **repository層**: 既存のrepository層（`dm_user_repository.go`）は存在するが、pages層から直接使用されていない

#### 2.1.2 api_key.goの現状
- **実装場所**: `server/internal/admin/pages/api_key.go`に直接実装されている
- **レイヤー構造**: pages層に鍵の生成とペイロードのデコードが直接実装されている
- **バリデーション**: なし（POSTリクエストのみ）
- **入出力制御**: pages層内で`renderAPIKeyPage`、`renderAPIKeyResult`関数として実装
- **ビジネスロジック**: pages層内に直接実装（`generatePublicAPIKey`、`auth.ParseJWTClaims`を直接呼び出し）
- **usecase層**: Admin用のusecase層が存在しない
- **service層**: Admin用のservice層が存在しない（鍵の生成とペイロードのデコードが1つの処理として実装されている）

### 2.2 現在の処理内容

#### 2.2.1 dm_user_register.goの処理内容
1. **バリデーション**: 名前とメールアドレスの入力値をバリデーション
2. **メールアドレスの重複チェック**: 全シャードを検索してメールアドレスの重複をチェック
3. **dm_user登録**: UUIDv7でIDを生成し、シャーディング対応でdm_userを登録
4. **登録完了ページへリダイレクト**: クエリパラメータで情報を渡してリダイレクト

#### 2.2.2 api_key.goの処理内容
1. **鍵の生成**: JWTトークンを生成（`auth.GeneratePublicAPIKey`を使用）
2. **ペイロードのデコード**: 生成したJWTトークンのペイロードをデコード（`auth.ParseJWTClaims`を使用）
3. **結果の表示**: 生成したJWTトークンとペイロードを表示

### 2.3 課題点
1. **レイヤー構造の不一致**: APIサーバーはusecase層を使用しているが、Adminアプリは処理を直接実装している
2. **責務の混在**: pages層にビジネスロジックが直接実装されている
3. **再利用性の低さ**: Adminアプリの処理がpages層に直接実装されており、他のコンポーネントから再利用できない
4. **テストの困難さ**: pages層にロジックが集中しているため、単体テストが困難
5. **コードの重複**: メールアドレスの重複チェックやdm_user登録処理がpages層に直接実装されている
6. **処理の分離不足**: api_key.goで鍵の生成とペイロードのデコードが1つの処理として実装されている

### 2.4 本実装による改善点
1. **アーキテクチャの一貫性**: APIサーバーとAdminアプリで同じレイヤー構造を採用
2. **責務の明確化**: pages層はエントリーポイントと入出力制御のみを担当
3. **再利用性の向上**: Adminアプリの処理を各レイヤーに分離することで、他のコンポーネントからも再利用可能
4. **テスト容易性**: usecase層、service層、repository層を独立してテストできるようになる
5. **コードの整理**: ビジネスロジックを各レイヤーに分離することで、コードの重複を削減
6. **処理の分離**: api_key.goで鍵の生成とペイロードのデコードを2つの処理に分離

## 3. 機能要件

### 3.1 Admin用usecase層の作成

#### 3.1.1 ディレクトリ構造
- **目的**: Adminアプリ用のusecase層を新規作成
- **実装内容**:
  - `server/internal/usecase/admin`ディレクトリを新規作成
  - `server/internal/usecase/admin/dm_user_register_usecase.go`を作成
  - `server/internal/usecase/admin/dm_user_register_usecase_test.go`を作成（テストファイル）
  - `server/internal/usecase/admin/api_key_usecase.go`を作成
  - `server/internal/usecase/admin/api_key_usecase_test.go`を作成（テストファイル）

#### 3.1.2 DmUserRegisterUsecaseの実装
- **目的**: dm_user登録用のusecaseを実装
- **実装内容**:
  - `DmUserRegisterUsecase`構造体を定義
  - `DmUserServiceInterface`を依存として注入（既存の`server/internal/usecase/dm_user_usecase.go`で定義されているインターフェースを使用）
  - `RegisterDmUser(ctx context.Context, name, email string) (string, error)`メソッドを実装
  - service層の`CreateDmUser()`を呼び出し
  - エラーハンドリングを実装

#### 3.1.3 APIKeyUsecaseの実装
- **目的**: APIキー発行用のusecaseを実装
- **実装内容**:
  - `APIKeyUsecase`構造体を定義
  - `APIKeyServiceInterface`を依存として注入（新規作成）
  - `GenerateAPIKey(ctx context.Context, env string) (string, error)`メソッドを実装（鍵の生成）
  - `DecodeAPIKeyPayload(ctx context.Context, token string) (*auth.JWTClaims, error)`メソッドを実装（ペイロードのデコード）
  - service層の`GenerateAPIKey()`、`DecodeAPIKeyPayload()`を呼び出し
  - エラーハンドリングを実装

### 3.2 Service層の拡張

#### 3.2.1 既存のservice層の確認
- **目的**: 既存のservice層を確認し、必要に応じて拡張
- **実装内容**:
  - `server/internal/service/dm_user_service.go`が存在することを確認
  - 既存の`DmUserService`が`CreateDmUser`メソッドを実装していることを確認
  - 既存のservice層をそのまま使用（新規作成は不要）

#### 3.2.2 APIKeyServiceの作成
- **目的**: APIキー発行用のservice層を新規作成
- **実装内容**:
  - `server/internal/service/api_key_service.go`を作成
  - `server/internal/service/api_key_service_test.go`を作成（テストファイル）
  - `APIKeyService`構造体を定義
  - `APIKeyServiceInterface`を定義
  - `GenerateAPIKey(ctx context.Context, secretKey, version, env string, issuedAt int64) (string, error)`メソッドを実装（鍵の生成）
  - `DecodeAPIKeyPayload(ctx context.Context, token string) (*auth.JWTClaims, error)`メソッドを実装（ペイロードのデコード）
  - `auth.GeneratePublicAPIKey`、`auth.ParseJWTClaims`を呼び出し

### 3.3 Repository層の確認

#### 3.3.1 既存のrepository層の確認
- **目的**: 既存のrepository層を確認し、必要に応じて使用
- **実装内容**:
  - `server/internal/repository/dm_user_repository.go`が存在することを確認
  - 既存のrepository層をそのまま使用（新規作成は不要）
  - service層からrepository層を呼び出す（既存のパターンに従う）

### 3.4 pages層の修正

#### 3.4.1 dm_user_register.goの修正
- **目的**: pages層をエントリーポイントと入出力制御のみに限定
- **実装内容**:
  - `DmUserRegisterPage`関数を修正（バリデーションと入出力制御のみを担当）
  - `handleDmUserRegisterPost`関数を修正（usecase層を呼び出すように変更）
  - `validateDmUserInput`関数を維持（pages層でバリデーション）
  - `renderDmUserRegisterForm`関数を維持（pages層で入出力制御）
  - `checkEmailExistsSharded`関数を削除（service層に移動）
  - `insertDmUserSharded`関数を削除（service層に移動）
  - usecase層の`RegisterDmUser()`を呼び出し

#### 3.4.2 api_key.goの修正
- **目的**: pages層をエントリーポイントと入出力制御のみに限定
- **実装内容**:
  - `APIKeyPage`関数を修正（入出力制御のみを担当）
  - `handleGenerateKey`関数を修正（usecase層を呼び出すように変更）
  - `renderAPIKeyPage`関数を維持（pages層で入出力制御）
  - `renderAPIKeyResult`関数を維持（pages層で入出力制御）
  - `generatePublicAPIKey`関数を削除（service層に移動）
  - usecase層の`GenerateAPIKey()`、`DecodeAPIKeyPayload()`を呼び出し（2つの処理に分ける）

### 3.5 依存関係の注入

#### 3.5.1 usecase層の初期化
- **目的**: usecase層を適切に初期化してpages層で使用
- **実装内容**:
  - `DmUserRegisterUsecase`の初期化（`DmUserService`を依存として注入）
  - `APIKeyUsecase`の初期化（`APIKeyService`を依存として注入）
  - pages層でusecase層を使用する方法を確認（GoAdminのカスタムページ登録方法を確認）

#### 3.5.2 service層の初期化
- **目的**: service層を適切に初期化してusecase層で使用
- **実装内容**:
  - `DmUserService`の初期化（既存の初期化方法を確認）
  - `APIKeyService`の初期化（新規作成）
  - Repository層の初期化（既存の初期化方法を確認）

### 3.6 ドキュメントの更新

#### 3.6.1 アーキテクチャドキュメントの更新
- **目的**: Adminアプリのレイヤー構造をアーキテクチャドキュメントに反映
- **修正対象**: `docs/Architecture.md`
- **実装内容**:
  - Adminアプリのレイヤー構造を追加（usecase層を含む）
  - Adminアプリのアーキテクチャ図を更新
  - Admin用usecase層の説明を追加

#### 3.6.2 プロジェクト構造ドキュメントの更新
- **目的**: 新規作成するファイルをプロジェクト構造に反映
- **修正対象**: `docs/Project-Structure.md`
- **実装内容**:
  - `server/internal/usecase/admin`ディレクトリを追加
  - `server/internal/usecase/admin/dm_user_register_usecase.go`を追加
  - `server/internal/usecase/admin/api_key_usecase.go`を追加
  - `server/internal/service/api_key_service.go`を追加

#### 3.6.3 ファイル組織ドキュメントの更新
- **目的**: 新規作成するファイルをファイル組織に反映
- **修正対象**: `.kiro/steering/structure.md`
- **実装内容**:
  - `server/internal/usecase/admin`ディレクトリを追加
  - `server/internal/usecase/admin/dm_user_register_usecase.go`を追加
  - `server/internal/usecase/admin/api_key_usecase.go`を追加
  - `server/internal/service/api_key_service.go`を追加

## 4. 非機能要件

### 4.1 パフォーマンス
- **既存機能の維持**: 既存のAdminページのパフォーマンスを維持
- **オーバーヘッド**: usecase層、service層の追加によるパフォーマンスオーバーヘッドは無視できるレベル

### 4.2 信頼性
- **エラーハンドリング**: 既存のエラーハンドリングを維持
- **後方互換性**: 既存のAdminページの動作を維持（出力形式、エラーメッセージなど）
- **データ整合性**: 既存のデータ登録ロジックを維持（UUID生成、テーブル番号計算など）

### 4.3 保守性
- **コードの可読性**: レイヤー構造が明確で、各レイヤーの責務が明確
- **一貫性**: APIサーバーと同じレイヤー構造を採用することで、コードベース全体の一貫性を向上
- **テスト容易性**: usecase層、service層を独立してテストできる
- **再利用性**: Adminアプリの処理を各レイヤーに分離することで、他のコンポーネントからも再利用可能

### 4.4 互換性
- **既存機能**: 既存のAPIサーバーに影響を与えない
- **Adminページの動作**: 既存のAdminページの動作（出力形式、エラーメッセージ）を維持
- **データ登録ロジック**: 既存のデータ登録ロジック（UUID生成、テーブル番号計算など）を維持

## 5. 制約事項

### 5.1 技術的制約
- **既存のデータ登録ロジック**: 既存のデータ登録ロジック（UUID生成、テーブル番号計算）を変更しない
- **既存の鍵生成ロジック**: 既存の鍵生成ロジック（`auth.GeneratePublicAPIKey`、`auth.ParseJWTClaims`）を変更しない
- **GoAdminのカスタムページ**: GoAdminのカスタムページ登録方法を維持

### 5.2 実装上の制約
- **ディレクトリ構造**: 既存のディレクトリ構造に従う（`server/internal/usecase/admin`、`server/internal/service`、`server/internal/repository`）
- **命名規則**: 既存の命名規則に従う（`DmUserRegisterUsecase`、`APIKeyUsecase`、`APIKeyService`）
- **インターフェース**: 既存の`DmUserServiceInterface`を使用。`APIKeyServiceInterface`を新規作成

### 5.3 動作環境
- **ローカル環境**: ローカル環境でAdminページが正常に動作することを確認
- **CI環境**: CI環境でもAdminページが正常に動作することを確認（該当する場合）
- **データベース**: PostgreSQL（master 1台 + sharding 4台）が正常に動作していることを前提

## 6. 受け入れ基準

### 6.1 Admin用usecase層の作成
- [ ] `server/internal/usecase/admin`ディレクトリが作成されている
- [ ] `server/internal/usecase/admin/dm_user_register_usecase.go`が作成されている
- [ ] `DmUserRegisterUsecase`構造体が定義されている
- [ ] `DmUserServiceInterface`を依存として注入している
- [ ] `RegisterDmUser(ctx context.Context, name, email string) (string, error)`メソッドが実装されている
- [ ] `server/internal/usecase/admin/dm_user_register_usecase_test.go`が作成されている
- [ ] `server/internal/usecase/admin/api_key_usecase.go`が作成されている
- [ ] `APIKeyUsecase`構造体が定義されている
- [ ] `APIKeyServiceInterface`を依存として注入している
- [ ] `GenerateAPIKey(ctx context.Context, env string) (string, error)`メソッドが実装されている
- [ ] `DecodeAPIKeyPayload(ctx context.Context, token string) (*auth.JWTClaims, error)`メソッドが実装されている
- [ ] `server/internal/usecase/admin/api_key_usecase_test.go`が作成されている
- [ ] 単体テストが実装されている

### 6.2 Service層の拡張
- [ ] `server/internal/service/api_key_service.go`が作成されている
- [ ] `APIKeyService`構造体が定義されている
- [ ] `APIKeyServiceInterface`が定義されている
- [ ] `GenerateAPIKey(ctx context.Context, secretKey, version, env string, issuedAt int64) (string, error)`メソッドが実装されている
- [ ] `DecodeAPIKeyPayload(ctx context.Context, token string) (*auth.JWTClaims, error)`メソッドが実装されている
- [ ] `auth.GeneratePublicAPIKey`、`auth.ParseJWTClaims`を呼び出している
- [ ] `server/internal/service/api_key_service_test.go`が作成されている
- [ ] 単体テストが実装されている

### 6.3 pages層の修正
- [ ] `server/internal/admin/pages/dm_user_register.go`が修正されている
- [ ] `DmUserRegisterPage`関数がバリデーションと入出力制御のみを担当している
- [ ] `handleDmUserRegisterPost`関数がusecase層を呼び出すように変更されている
- [ ] `validateDmUserInput`関数が維持されている（pages層でバリデーション）
- [ ] `renderDmUserRegisterForm`関数が維持されている（pages層で入出力制御）
- [ ] `checkEmailExistsSharded`関数が削除されている（service層に移動）
- [ ] `insertDmUserSharded`関数が削除されている（service層に移動）
- [ ] `server/internal/admin/pages/api_key.go`が修正されている
- [ ] `APIKeyPage`関数が入出力制御のみを担当している
- [ ] `handleGenerateKey`関数がusecase層を呼び出すように変更されている
- [ ] `renderAPIKeyPage`関数が維持されている（pages層で入出力制御）
- [ ] `renderAPIKeyResult`関数が維持されている（pages層で入出力制御）
- [ ] `generatePublicAPIKey`関数が削除されている（service層に移動）
- [ ] usecase層の`GenerateAPIKey()`、`DecodeAPIKeyPayload()`を呼び出している（2つの処理に分ける）

### 6.4 依存関係の注入
- [ ] `DmUserRegisterUsecase`が適切に初期化されている（`DmUserService`を依存として注入）
- [ ] `APIKeyUsecase`が適切に初期化されている（`APIKeyService`を依存として注入）
- [ ] pages層でusecase層を使用している

### 6.5 動作確認
- [ ] ローカル環境でAdminページが正常に動作する
- [ ] 既存の出力形式（HTML出力）が維持されている
- [ ] 既存のエラーメッセージが維持されている
- [ ] dm_user登録が正常に動作する
- [ ] APIキー発行が正常に動作する
- [ ] 既存のテストが全て通過する
- [ ] CI環境でAdminページが正常に動作する（該当する場合）

### 6.6 テスト
- [ ] `server/internal/usecase/admin/dm_user_register_usecase_test.go`が作成されている
- [ ] `server/internal/usecase/admin/api_key_usecase_test.go`が作成されている
- [ ] `server/internal/service/api_key_service_test.go`が作成されている
- [ ] 各層の単体テストが実装されている
- [ ] 既存のテストが全て通過する

### 6.7 ドキュメントの更新
- [ ] `docs/Architecture.md`にAdminアプリのレイヤー構造が追加されている
- [ ] `docs/Architecture.md`のAdminアプリのアーキテクチャ図が更新されている
- [ ] `docs/Project-Structure.md`に新規作成するファイルが追加されている
- [ ] `.kiro/steering/structure.md`に新規作成するファイルが追加されている

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### 新規作成が必要なファイル
- `server/internal/usecase/admin/dm_user_register_usecase.go`: dm_user登録用のusecase層
- `server/internal/usecase/admin/dm_user_register_usecase_test.go`: dm_user登録用usecase層のテスト
- `server/internal/usecase/admin/api_key_usecase.go`: APIキー発行用のusecase層
- `server/internal/usecase/admin/api_key_usecase_test.go`: APIキー発行用usecase層のテスト
- `server/internal/service/api_key_service.go`: APIキー発行用のservice層
- `server/internal/service/api_key_service_test.go`: APIキー発行用service層のテスト

#### 修正が必要なファイル
- `server/internal/admin/pages/dm_user_register.go`: usecase層を使用するように修正
- `server/internal/admin/pages/api_key.go`: usecase層を使用するように修正（2つの処理に分ける）
- `docs/Architecture.md`: Adminアプリのレイヤー構造を追加
- `docs/Project-Structure.md`: 新規作成するファイルを追加
- `.kiro/steering/structure.md`: 新規作成するファイルを追加

### 7.2 既存機能への影響
- **既存のAPIサーバー**: 影響なし（APIサーバーは変更しない）
- **既存のAdminページの動作**: 影響なし（動作は維持される）
- **既存のデータ登録ロジック**: 影響なし（ロジックは維持される）

## 8. 実装上の注意事項

### 8.1 Admin用usecase層の実装
- **インターフェースの使用**: 既存の`DmUserServiceInterface`を使用。`APIKeyServiceInterface`を新規作成
- **依存関係の注入**: コンストラクタでservice層のインターフェースを注入
- **エラーハンドリング**: service層から返されたエラーをそのまま返す（エラーのラップは不要）
- **処理の分離**: api_key.goで鍵の生成とペイロードのデコードを2つの処理に分ける

### 8.2 Service層の実装
- **既存のservice層の使用**: `DmUserService`は既存のものを使用（新規作成は不要）
- **新規service層の作成**: `APIKeyService`を新規作成
- **依存関係の注入**: コンストラクタで依存関係を注入（必要に応じて）
- **エラーハンドリング**: 下位層から返されたエラーをそのまま返す（エラーのラップは不要）
- **鍵の生成とペイロードのデコード**: 2つの処理に分ける（`GenerateAPIKey`、`DecodeAPIKeyPayload`）

### 8.3 pages層の修正
- **バリデーション**: pages層でバリデーションを実施（既存の`validateDmUserInput`を維持）
- **入出力制御**: pages層で入出力制御を実施（既存の`renderDmUserRegisterForm`、`renderAPIKeyPage`、`renderAPIKeyResult`を維持）
- **usecase層の呼び出し**: usecase層のメソッドを呼び出すように変更
- **既存の関数の削除**: `checkEmailExistsSharded`、`insertDmUserSharded`、`generatePublicAPIKey`を削除（service層に移動）

### 8.4 依存関係の注入
- **usecase層の初期化**: GoAdminのカスタムページ登録方法を確認し、usecase層を適切に初期化
- **service層の初期化**: 既存の初期化方法を確認し、service層を適切に初期化
- **Repository層の初期化**: 既存の初期化方法を確認し、repository層を適切に初期化

### 8.5 テストの実装
- **usecase層のテスト**: service層のインターフェースのモックを使用してテスト
- **service層のテスト**: 下位層のモックを使用してテスト（必要に応じて）

### 8.6 ドキュメントの更新
- **アーキテクチャドキュメント**: Adminアプリのレイヤー構造を明確に記載
- **プロジェクト構造ドキュメント**: 新規作成するファイルを反映
- **ファイル組織ドキュメント**: 新規作成するファイルを反映
- **一貫性**: 全てのドキュメントで同じレイヤー構造を記載

## 9. 参考情報

### 9.1 関連ドキュメント
- `docs/Architecture.md`: アーキテクチャドキュメント
- `docs/Project-Structure.md`: プロジェクト構造ドキュメント
- `.kiro/steering/structure.md`: ファイル組織とコードパターン
- `.kiro/steering/tech.md`: 技術スタックとアーキテクチャ

### 9.2 既存実装の参考
- `server/internal/usecase/dm_user_usecase.go`: 既存のusecase層の実装パターン
- `server/internal/usecase/cli/list_dm_users_usecase.go`: 既存のCLI用usecase層の実装パターン
- `server/internal/service/dm_user_service.go`: 既存のservice層の実装パターン
- `server/internal/admin/pages/dm_user_register.go`: 既存のAdminページの実装
- `server/internal/admin/pages/api_key.go`: 既存のAdminページの実装

### 9.3 技術スタック
- **言語**: Go
- **アーキテクチャ**: レイヤードアーキテクチャ（pages -> usecase -> service -> repository -> db）
- **テスト**: `testing`（標準ライブラリ）、`github.com/stretchr/testify`（アサーション、モック）
- **Adminフレームワーク**: GoAdmin

### 9.4 レイヤー構造の比較

| 項目 | 現在（修正前） | 修正後 |
|------|---------------|--------|
| Pages層 | pages.go（ビジネスロジックが直接実装） | pages.go（エントリーポイント、バリデーション、入出力） |
| Usecase層 | なし | `usecase/admin/DmUserRegisterUsecase`、`usecase/admin/APIKeyUsecase` |
| Service層 | なし（dm_user_service.goは存在するが使用されていない） | `service.DmUserService`（既存）、`service.APIKeyService`（新規） |
| Repository層 | なし（pages層から直接DB操作） | `repository.DmUserRepository`（既存） |
| DB層 | pages層内で直接使用 | `db.GroupManager`（既存） |
| 出力 | pages層内で直接出力 | pages層内で出力（usecase層から取得） |

### 9.5 APIサーバーとの比較

| 項目 | APIサーバー | Adminアプリ（修正後） |
|------|------------|---------------------|
| エントリーポイント | `server/cmd/server/main.go` | `server/internal/admin/pages/*.go` |
| バリデーション | API Layer（Handler） | Pages層 |
| Usecase層 | `usecase.DmUserUsecase` | `usecase/admin.DmUserRegisterUsecase`、`usecase/admin.APIKeyUsecase` |
| Service層 | `service.DmUserService` | `service.DmUserService`（既存）、`service.APIKeyService`（新規） |
| Repository層 | `repository.DmUserRepository` | `repository.DmUserRepository`（既存） |
| DB層 | `db.GroupManager` | `db.GroupManager` |
| 出力 | HTTPレスポンス | HTML（GoAdminのPanel） |
