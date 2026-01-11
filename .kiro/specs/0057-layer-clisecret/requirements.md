# server/cmd/generate-secretの構造修正の要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Feature名**: 0057-layer-clisecret
- **作成日**: 2026-01-27
- **関連Issue**: https://github.com/taku-o/go-webdb-template/issues/118

### 1.2 目的
`server/cmd/generate-secret`の実装を、APIサーバーと同じレイヤー構造（usecase -> service -> 秘密鍵生成処理）に変更する。これにより、CLIコマンドとAPIサーバーで一貫したアーキテクチャを実現し、秘密鍵生成処理を共通化してコードの保守性と再利用性を向上させる。

### 1.3 スコープ
- CLIコマンドのバリデーションと入出力制御をCLI層に集約
- `server/internal/usecase/cli`ディレクトリに`generate-secret`用のusecaseを追加
- 秘密鍵生成処理を`server/internal/auth/`に共通ライブラリとして実装
- main.goの簡素化（エントリーポイントと入出力制御のみ）
- 関連ドキュメントの更新（アーキテクチャ、プロジェクト構造、CLIツールドキュメント）

**本実装の範囲外**:
- 既存のJWT関連処理（`server/internal/auth/jwt.go`の`GeneratePublicAPIKey`など）への影響（既存実装をそのまま使用）
- 既存のAPIサーバーのusecase層への影響（APIサーバーは変更しない）
- 他のCLIコマンドへの影響（本実装は`generate-secret`のみを対象）

## 2. 背景・現状分析

### 2.1 現在の状況
- **CLI実装**: `server/cmd/generate-secret/main.go`に直接実装されている
- **レイヤー構造**: main.goに秘密鍵生成処理が直接実装されている（32バイトのランダムな秘密鍵を生成してBase64エンコード）
- **バリデーション**: なし（引数なし）
- **入出力制御**: main.go内で標準出力に直接出力
- **ビジネスロジック**: main.go内に直接実装（`crypto/rand`と`encoding/base64`を使用）
- **usecase層**: CLI用のusecase層が存在しない
- **共通ライブラリ**: 秘密鍵生成処理の共通ライブラリが存在しない

### 2.2 課題点
1. **レイヤー構造の不一致**: APIサーバーはusecase層を使用しているが、CLIは処理を直接実装している
2. **責務の混在**: main.goに秘密鍵生成処理が直接実装されている
3. **再利用性の低さ**: 秘密鍵生成処理がmain.goに直接実装されており、他のCLIコマンドから再利用できない
4. **テストの困難さ**: main.goにロジックが集中しているため、単体テストが困難

### 2.3 本実装による改善点
1. **アーキテクチャの一貫性**: APIサーバーとCLIで同じレイヤー構造を採用
2. **責務の明確化**: CLI層（main.go）はエントリーポイントと入出力制御のみを担当
3. **再利用性の向上**: 秘密鍵生成処理を共通ライブラリとして実装することで、CLIコマンドからも再利用可能
4. **テスト容易性**: usecase層とservice層を独立してテストできるようになる

## 3. 機能要件

### 3.1 秘密鍵生成処理の共通化

#### 3.1.1 共通ライブラリの作成
- **目的**: 秘密鍵生成処理を共通ライブラリとして実装
- **実装内容**:
  - `server/internal/auth/secret.go`を作成
  - `GenerateSecretKey() (string, error)`関数を実装
  - 32バイト（256ビット）のランダムな秘密鍵を生成してBase64エンコードして返す
  - `crypto/rand`と`encoding/base64`を使用
- **テスト**: `server/internal/auth/secret_test.go`を作成して単体テストを実装

### 3.2 Service層の作成

#### 3.2.1 SecretServiceの実装
- **目的**: 秘密鍵生成用のservice層を実装
- **実装内容**:
  - `server/internal/service/secret_service.go`を作成
  - `SecretService`構造体を定義
  - `GenerateSecretKey(ctx context.Context) (string, error)`メソッドを実装
  - 共通ライブラリ（`auth.GenerateSecretKey()`）を呼び出して結果を返す
- **テスト**: `server/internal/service/secret_service_test.go`を作成して単体テストを実装

### 3.3 CLI用usecase層の作成

#### 3.3.1 ディレクトリ構造
- **目的**: CLI用のusecase層を既存の`server/internal/usecase/cli`ディレクトリに追加
- **実装内容**:
  - `server/internal/usecase/cli/generate_secret_usecase.go`を作成
  - `server/internal/usecase/cli/generate_secret_usecase_test.go`を作成（テストファイル）

#### 3.3.2 GenerateSecretUsecaseの実装
- **目的**: CLI用の秘密鍵生成usecaseを実装
- **実装内容**:
  - `GenerateSecretUsecase`構造体を定義
  - `SecretServiceInterface`を依存として注入
  - `GenerateSecret(ctx context.Context) (string, error)`メソッドを実装
  - service層の`GenerateSecretKey()`を呼び出して結果を返す
- **インターフェース**: `SecretServiceInterface`を新規作成（`server/internal/service/secret_service.go`で定義）

### 3.4 main.goの簡素化

#### 3.4.1 エントリーポイントの責務
- **目的**: main.goをエントリーポイントと入出力制御のみに限定
- **実装内容**:
  - 設定ファイルの読み込み（`config.Load()`）は不要（秘密鍵生成には設定ファイルが不要なため）
  - usecase層の初期化と実行
  - 結果の出力（標準出力に直接出力）

#### 3.4.2 出力関数
- **目的**: 秘密鍵の出力をmain.goに保持
- **実装内容**:
  - usecase層から取得した秘密鍵を標準出力に出力
  - 改行を付けて出力（既存の動作を維持）

### 3.5 依存関係の注入

#### 3.5.1 usecase層の初期化
- **目的**: usecase層を適切に初期化してmain.goで使用
- **実装内容**:
  - Service層の初期化（`service.NewSecretService()`）
  - Usecase層の初期化（`cli.NewGenerateSecretUsecase(secretService)`）
  - usecase層の`GenerateSecret()`メソッドを呼び出し

### 3.6 ドキュメントの更新

#### 3.6.1 アーキテクチャドキュメントの更新
- **目的**: CLIコマンドのレイヤー構造をアーキテクチャドキュメントに反映
- **修正対象**: `docs/Architecture.md`
- **実装内容**:
  - CLIコマンドのレイヤー構造を追加（usecase層を含む）
  - CLIコマンドのアーキテクチャ図を更新
  - CLI用usecase層の説明を追加
  - 秘密鍵生成処理の共通ライブラリの説明を追加

#### 3.6.2 プロジェクト構造ドキュメントの更新
- **目的**: 新規作成するファイルをプロジェクト構造に反映
- **修正対象**: `docs/Project-Structure.md`
- **実装内容**:
  - `server/internal/usecase/cli/generate_secret_usecase.go`を追加
  - `server/internal/service/secret_service.go`を追加
  - `server/internal/auth/secret.go`を追加

#### 3.6.3 CLIツールドキュメントの更新
- **目的**: CLIコマンドのアーキテクチャ図を更新してusecase層を含める
- **修正対象**: `docs/Command-Line-Tool.md`（存在する場合）
- **実装内容**:
  - CLIコマンドのアーキテクチャ図を更新（usecase層を追加）
  - レイヤー構造の説明を更新（main.go → usecase → service → auth → 出力）

#### 3.6.4 ファイル組織ドキュメントの更新
- **目的**: 新規作成するファイルをファイル組織に反映
- **修正対象**: `.kiro/steering/structure.md`
- **実装内容**:
  - `server/internal/usecase/cli/generate_secret_usecase.go`を追加
  - `server/internal/service/secret_service.go`を追加
  - `server/internal/auth/secret.go`を追加

## 4. 非機能要件

### 4.1 パフォーマンス
- **既存機能の維持**: 既存の秘密鍵生成処理のパフォーマンスを維持
- **オーバーヘッド**: usecase層とservice層の追加によるパフォーマンスオーバーヘッドは無視できるレベル

### 4.2 信頼性
- **エラーハンドリング**: 既存のエラーハンドリングを維持
- **後方互換性**: 既存のCLIコマンドの動作を維持（出力形式、エラーメッセージなど）
- **セキュリティ**: 秘密鍵生成処理のセキュリティレベルを維持（`crypto/rand`を使用）

### 4.3 保守性
- **コードの可読性**: レイヤー構造が明確で、各レイヤーの責務が明確
- **一貫性**: APIサーバーと同じレイヤー構造を採用することで、コードベース全体の一貫性を向上
- **テスト容易性**: usecase層とservice層を独立してテストできる
- **再利用性**: 秘密鍵生成処理を共通ライブラリとして実装することで、他のコンポーネントからも再利用可能

### 4.4 互換性
- **既存機能**: 既存のJWT関連処理（`server/internal/auth/jwt.go`）に影響を与えない
- **CLIコマンドの動作**: 既存のCLIコマンドの動作（出力形式、エラーメッセージ）を維持

## 5. 制約事項

### 5.1 技術的制約
- **既存のauth層**: 既存のJWT関連処理（`server/internal/auth/jwt.go`）を変更しない
- **秘密鍵生成アルゴリズム**: 32バイト（256ビット）のランダムな秘密鍵を生成してBase64エンコード（既存の動作を維持）
- **crypto/randの使用**: `crypto/rand`パッケージを使用（セキュリティ要件）

### 5.2 実装上の制約
- **ディレクトリ構造**: `server/internal/usecase/cli`ディレクトリに追加（既存ディレクトリを使用）
- **命名規則**: 既存のusecase層の命名規則に従う（`GenerateSecretUsecase`）
- **インターフェース**: `SecretServiceInterface`を新規作成

### 5.3 動作環境
- **ローカル環境**: ローカル環境でCLIコマンドが正常に動作することを確認
- **CI環境**: CI環境でもCLIコマンドが正常に動作することを確認（該当する場合）

## 6. 受け入れ基準

### 6.1 秘密鍵生成処理の共通化
- [ ] `server/internal/auth/secret.go`が作成されている
- [ ] `GenerateSecretKey() (string, error)`関数が実装されている
- [ ] 32バイト（256ビット）のランダムな秘密鍵を生成してBase64エンコードして返している
- [ ] `crypto/rand`と`encoding/base64`を使用している
- [ ] `server/internal/auth/secret_test.go`が作成されている
- [ ] 単体テストが実装されている

### 6.2 Service層の作成
- [ ] `server/internal/service/secret_service.go`が作成されている
- [ ] `SecretService`構造体が定義されている
- [ ] `SecretServiceInterface`が定義されている
- [ ] `GenerateSecretKey(ctx context.Context) (string, error)`メソッドが実装されている
- [ ] 共通ライブラリ（`auth.GenerateSecretKey()`）を呼び出して結果を返している
- [ ] `server/internal/service/secret_service_test.go`が作成されている
- [ ] 単体テストが実装されている

### 6.3 CLI用usecase層の作成
- [ ] `server/internal/usecase/cli/generate_secret_usecase.go`が作成されている
- [ ] `GenerateSecretUsecase`構造体が定義されている
- [ ] `SecretServiceInterface`を依存として注入している
- [ ] `GenerateSecret(ctx context.Context) (string, error)`メソッドが実装されている
- [ ] service層の`GenerateSecretKey()`を呼び出して結果を返している
- [ ] `server/internal/usecase/cli/generate_secret_usecase_test.go`が作成されている
- [ ] 単体テストが実装されている

### 6.4 main.goの簡素化
- [ ] main.goがエントリーポイントと入出力制御のみを担当している
- [ ] usecase層を初期化して使用している
- [ ] 既存の出力形式（標準出力にBase64エンコードされた秘密鍵を出力）が維持されている
- [ ] 既存のエラーハンドリングが維持されている
- [ ] ビルド時のバイナリが`server/bin/generate-secret`に出力される

### 6.5 依存関係の注入
- [ ] Service層が適切に初期化されている
- [ ] Usecase層が適切に初期化されている
- [ ] usecase層の`GenerateSecret()`メソッドが呼び出されている

### 6.6 動作確認
- [ ] ローカル環境でCLIコマンドが正常に動作する
- [ ] 既存の出力形式（Base64エンコードされた秘密鍵）が維持されている
- [ ] 既存のエラーメッセージが維持されている
- [ ] 既存のテストが全て通過する
- [ ] CI環境でCLIコマンドが正常に動作する（該当する場合）

### 6.7 テスト
- [ ] `server/internal/auth/secret_test.go`が作成されている
- [ ] `server/internal/service/secret_service_test.go`が作成されている
- [ ] `server/internal/usecase/cli/generate_secret_usecase_test.go`が作成されている
- [ ] 各層の単体テストが実装されている
- [ ] 既存のテストが全て通過する

### 6.8 ドキュメントの更新
- [ ] `docs/Architecture.md`にCLIコマンドのレイヤー構造が追加されている
- [ ] `docs/Architecture.md`のCLIコマンドのアーキテクチャ図が更新されている
- [ ] `docs/Architecture.md`に秘密鍵生成処理の共通ライブラリの説明が追加されている
- [ ] `docs/Project-Structure.md`に新規作成するファイルが追加されている
- [ ] `docs/Command-Line-Tool.md`（存在する場合）のアーキテクチャ図が更新されている
- [ ] `docs/Command-Line-Tool.md`（存在する場合）のレイヤー構造の説明が更新されている
- [ ] `.kiro/steering/structure.md`に新規作成するファイルが追加されている

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### 新規作成が必要なファイル
- `server/internal/auth/secret.go`: 秘密鍵生成処理の共通ライブラリ
- `server/internal/auth/secret_test.go`: 秘密鍵生成処理のテスト
- `server/internal/service/secret_service.go`: 秘密鍵生成用のservice層
- `server/internal/service/secret_service_test.go`: service層のテスト
- `server/internal/usecase/cli/generate_secret_usecase.go`: CLI用usecase層の実装
- `server/internal/usecase/cli/generate_secret_usecase_test.go`: CLI用usecase層のテスト

#### 修正が必要なファイル
- `server/cmd/generate-secret/main.go`: usecase層を使用するように修正
- `docs/Architecture.md`: CLIコマンドのレイヤー構造を追加
- `docs/Project-Structure.md`: 新規作成するファイルを追加
- `docs/Command-Line-Tool.md`（存在する場合）: CLIコマンドのアーキテクチャ図を更新
- `.kiro/steering/structure.md`: 新規作成するファイルを追加

### 7.2 既存機能への影響
- **既存のJWT関連処理**: 影響なし（既存実装をそのまま使用）
- **既存のAPIサーバー**: 影響なし（APIサーバーは変更しない）
- **既存のCLIコマンドの動作**: 影響なし（動作は維持される）

## 8. 実装上の注意事項

### 8.1 秘密鍵生成処理の共通化
- **セキュリティ**: `crypto/rand`パッケージを使用して安全な乱数生成を行う
- **出力形式**: Base64エンコードされた文字列を返す（既存の動作を維持）
- **エラーハンドリング**: 乱数生成に失敗した場合は適切なエラーを返す

### 8.2 Service層の実装
- **インターフェースの定義**: `SecretServiceInterface`を新規作成
- **依存関係の注入**: コンストラクタで依存関係を注入（現時点では依存関係なし）
- **エラーハンドリング**: 共通ライブラリから返されたエラーをそのまま返す（エラーのラップは不要）

### 8.3 usecase層の実装
- **インターフェースの使用**: 新規作成する`SecretServiceInterface`を使用
- **依存関係の注入**: コンストラクタで`SecretServiceInterface`を注入
- **エラーハンドリング**: service層から返されたエラーをそのまま返す（エラーのラップは不要）

### 8.4 main.goの修正
- **usecase層の初期化**: Service → Usecaseの順で初期化
- **既存の出力**: 標準出力にBase64エンコードされた秘密鍵を出力（既存の動作を維持）
- **エラーハンドリング**: 既存のエラーハンドリングを維持
- **ビルド出力先**: ビルド時のバイナリは`server/bin/generate-secret`に出力する

### 8.5 テストの実装
- **共通ライブラリのテスト**: 秘密鍵生成処理のテストを実装
- **service層のテスト**: `SecretServiceInterface`のモックを使用してテスト
- **usecase層のテスト**: `SecretServiceInterface`のモックを使用してテスト

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
- `.kiro/steering/structure.md`: ファイル組織とコードパターン
- `.kiro/steering/tech.md`: 技術スタックとアーキテクチャ

### 9.2 既存実装の参考
- `server/internal/usecase/cli/list_dm_users_usecase.go`: 既存のCLI用usecase層の実装パターン
- `server/cmd/generate-secret/main.go`: 既存のCLI実装
- `server/internal/auth/jwt.go`: 既存のauth層の実装パターン

### 9.3 技術スタック
- **言語**: Go
- **アーキテクチャ**: レイヤードアーキテクチャ（usecase -> service -> auth -> 出力）
- **テスト**: `testing`（標準ライブラリ）、`github.com/stretchr/testify`（アサーション、モック）
- **セキュリティ**: `crypto/rand`（安全な乱数生成）

### 9.4 レイヤー構造の比較

| 項目 | 現在（修正前） | 修正後 |
|------|---------------|--------|
| CLI層 | main.go（秘密鍵生成処理が直接実装） | main.go（エントリーポイント、入出力） |
| Usecase層 | なし | `usecase/cli/GenerateSecretUsecase` |
| Service層 | なし | `service.SecretService` |
| Auth層 | なし | `auth.GenerateSecretKey()` |
| 出力 | main.go内で直接出力 | main.go内で出力（usecase層から取得） |

### 9.5 APIサーバーとの比較

| 項目 | APIサーバー | CLI（修正後） |
|------|------------|--------------|
| エントリーポイント | `server/cmd/server/main.go` | `server/cmd/generate-secret/main.go` |
| バリデーション | API Layer（Handler） | CLI層（main.go、現時点では不要） |
| Usecase層 | `usecase.DmUserUsecase` | `usecase/cli.GenerateSecretUsecase` |
| Service層 | `service.DmUserService` | `service.SecretService` |
| Auth層 | `auth.GeneratePublicAPIKey` | `auth.GenerateSecretKey` |
| 出力 | HTTPレスポンス | 標準出力 |
