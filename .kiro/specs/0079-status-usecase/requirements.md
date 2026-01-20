# サーバー状態確認機能のリファクタリング要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Feature名**: 0079-status-usecase
- **作成日**: 2026-01-17
- **関連Issue**: https://github.com/taku-o/go-webdb-template/issues/161

### 1.2 目的
`server/cmd/server-status/main.go`の実装を、プロジェクトの標準的なアーキテクチャパターンに従って、usecase層とservice層を利用した3層構造にリファクタリングする。これにより、コードの保守性とテスタビリティを向上させ、プロジェクト全体のアーキテクチャの一貫性を保つ。

### 1.3 スコープ
- `server/cmd/server-status/main.go`を入出力制御のみを担当するように変更する
- `server/internal/usecase/cli/server_status_usecase.go`を新規作成する（CLI用usecase層）
- `server/internal/service/server_status_service.go`を新規作成する（service層）
- 既存の機能（サーバー状態の確認と表示）は維持する

**本実装の範囲外**:
- サーバー状態確認のロジック自体の変更（ポート接続による確認方法は維持）
- 表示形式の変更
- 確認対象サーバーの追加・削除
- その他の機能追加

## 2. 背景・現状分析

### 2.1 現在の状況
- **現在の実装**: `server/cmd/server-status/main.go`に全てのロジックが含まれている
  - サーバー情報の定義（`servers`変数）
  - サーバー状態確認ロジック（`checkServerStatus`関数）
  - 並列実行ロジック（`checkAllServers`関数）
  - 結果表示ロジック（`printResults`関数）
  - エントリーポイント（`main`関数）
- **問題点**:
  1. **アーキテクチャの不整合**: プロジェクトの標準的な3層構造（Handler/CLI → Usecase → Service → Repository）に従っていない
  2. **テスタビリティの低さ**: ロジックがmain.goに集約されており、単体テストが困難
  3. **再利用性の欠如**: 他のCLIコマンドやAPIからサーバー状態確認機能を再利用できない
  4. **責務の混在**: 入出力制御、ビジネスロジック、データアクセスが一つのファイルに混在している

### 2.2 プロジェクトの標準アーキテクチャ
プロジェクトでは以下の3層構造が標準パターンとして使用されている：

```
CLI/Handler層 (cmd/server-status/main.go)
  ↓
Usecase層 (internal/usecase/cli/server_status_usecase.go)
  ↓
Service層 (internal/service/server_status_service.go)
```

**既存のCLI実装例**:
- `server/cmd/list-dm-users/main.go`: CLIエントリーポイント
- `server/internal/usecase/cli/list_dm_users_usecase.go`: CLI用usecase
- `server/internal/service/dm_user_service.go`: Service層（既存）

### 2.3 本実装による改善点
1. **アーキテクチャの一貫性**: プロジェクト全体で統一されたアーキテクチャパターンに従う
2. **テスタビリティの向上**: 各層を独立してテストできるようになる
3. **再利用性の向上**: Service層を他のCLIコマンドやAPIから再利用可能
4. **保守性の向上**: 責務が明確に分離され、変更の影響範囲が限定される
5. **コードの可読性向上**: 各層の責務が明確になり、コードが理解しやすくなる

## 3. 機能要件

### 3.1 レイヤー分離

#### 3.1.1 CLI層（`server/cmd/server-status/main.go`）
**責務**:
- 入出力制御を担当
- usecaseを呼び出す
- usecaseから受け取ったリストをコンソール出力する

**実装内容**:
- usecaseのインスタンスを作成
- usecaseのメソッドを呼び出し
- 受け取った結果を表形式で表示（既存の`printResults`関数のロジックを維持）

**変更点**:
- サーバー状態確認ロジックを削除（usecaseに移譲）
- サーバー情報の定義を削除（service層に移譲）
- 並列実行ロジックを削除（service層に移譲）

#### 3.1.2 Usecase層（`server/internal/usecase/cli/server_status_usecase.go`）
**責務**:
- serviceに渡すパラメータを作る
- serviceから受け取ったリストをmain.goに渡す

**実装内容**:
- serviceのインスタンスを保持
- serviceのメソッドを呼び出し
- serviceから受け取った結果をそのまま返す（必要に応じて変換）

**インターフェース**:
```go
type ServerStatusUsecase struct {
    serverStatusService ServerStatusServiceInterface
}

func (u *ServerStatusUsecase) ListServerStatus(ctx context.Context) ([]ServerStatus, error)
```

#### 3.1.3 Service層（`server/internal/service/server_status_service.go`）
**責務**:
- `list({name, host, port})`を受け取って、結果のリストを返す
- サーバー状態確認のビジネスロジックを実装
- 並列実行による状態確認

**実装内容**:
- サーバー情報の定義（既存の`servers`変数の内容）
- サーバー状態確認ロジック（既存の`checkServerStatus`関数の内容）
- 並列実行ロジック（既存の`checkAllServers`関数の内容）

**インターフェース**:
```go
type ServerStatusService struct {
    // 必要に応じて依存関係を追加
}

func (s *ServerStatusService) ListServerStatus(ctx context.Context, servers []ServerInfo) ([]ServerStatus, error)
```

### 3.2 データ構造

#### 3.2.1 ServerInfo
サーバー情報を表す構造体（既存の定義を維持）

```go
type ServerInfo struct {
    Name    string // サーバー名
    Port    int    // ポート番号
    Address string // 接続先アドレス（通常は"localhost"）
}
```

#### 3.2.2 ServerStatus
サーバーの状態を表す構造体（既存の定義を維持）

```go
type ServerStatus struct {
    Server ServerInfo
    Status string // "起動中" または "停止中"
    Error  error  // エラー情報（デバッグ用、表示には使用しない）
}
```

### 3.3 実装の詳細

#### 3.3.1 サーバー情報の定義場所
- **現在**: `server/cmd/server-status/main.go`の`servers`変数
- **変更後**: `server/internal/service/server_status_service.go`内で定義

#### 3.3.2 状態確認ロジックの配置
- **現在**: `server/cmd/server-status/main.go`の`checkServerStatus`関数
- **変更後**: `server/internal/service/server_status_service.go`内に実装

#### 3.3.3 並列実行ロジックの配置
- **現在**: `server/cmd/server-status/main.go`の`checkAllServers`関数
- **変更後**: `server/internal/service/server_status_service.go`内に実装

#### 3.3.4 表示ロジックの配置
- **現在**: `server/cmd/server-status/main.go`の`printResults`関数
- **変更後**: `server/cmd/server-status/main.go`に維持（CLI層の責務）

## 4. 非機能要件

### 4.1 パフォーマンス
- **実行時間**: 既存の実装と同等のパフォーマンスを維持（並列実行による状態確認）
- **リソース使用量**: 既存の実装と同等のリソース使用量を維持

### 4.2 保守性
- **コードの簡潔性**: 各層の責務が明確に分離されている
- **拡張性**: 将来的にサーバーが追加された場合、service層の設定を変更するだけで対応できる
- **テスタビリティ**: 各層を独立してテストできる

### 4.3 一貫性
- **アーキテクチャパターン**: プロジェクトの標準的な3層構造に従う
- **命名規則**: 既存のCLI実装（`list_dm_users_usecase.go`など）と同様の命名規則を使用
- **コードスタイル**: 既存のコードスタイルに従う

### 4.4 動作環境
- **ローカル環境**: ローカル環境で動作する（既存の実装と同等）
- **OS**: macOS、Linuxで動作する（既存の実装と同等）

## 5. 制約事項

### 5.1 技術的制約
- **既存機能の維持**: 既存のサーバー状態確認機能は完全に維持する
- **表示形式の維持**: 既存の表形式での表示を維持する
- **確認対象サーバー**: 既存の13個のサーバーを確認対象とする（変更なし）

### 5.2 実装上の制約
- **後方互換性**: コマンドの実行方法と出力形式は既存と同一とする
- **既存コードへの影響**: 既存の他のコードに影響を与えない
- **依存関係**: 既存の依存関係を追加しない（標準ライブラリのみ使用）

### 5.3 動作環境
- **ローカル環境**: ローカル環境でのみ動作する（既存の実装と同等）
- **開発環境**: 開発環境（`APP_ENV=develop`）を想定

## 6. 受け入れ基準

### 6.1 レイヤー分離
- [ ] `server/cmd/server-status/main.go`が入出力制御のみを担当している
- [ ] `server/internal/usecase/cli/server_status_usecase.go`が作成されている
- [ ] `server/internal/service/server_status_service.go`が作成されている
- [ ] 各層の責務が明確に分離されている

### 6.2 機能の維持
- [ ] 既存のサーバー状態確認機能が正常に動作する
- [ ] 既存の表示形式（表形式）が維持されている
- [ ] 既存の13個のサーバーが確認対象として維持されている
- [ ] 並列実行による状態確認が正常に動作する

### 6.3 動作確認
- [ ] コマンド実行が正常に完了する（`go run ./cmd/server-status/main.go`）
- [ ] 全てのサーバーが起動している場合、全て「起動中」と表示される
- [ ] 一部のサーバーが停止している場合、該当サーバーが「停止中」と表示される
- [ ] 全てのサーバーが停止している場合、全て「停止中」と表示される
- [ ] サーバーが指定された順序で表示される

### 6.4 テスト
- [ ] usecase層の単体テストが実装されている
- [ ] service層の単体テストが実装されている
- [ ] 既存のテストが全て失敗しないことを確認
- [ ] テストカバレッジが適切に維持されている

### 6.5 コード品質
- [ ] プロジェクトの標準的なアーキテクチャパターンに従っている
- [ ] 既存のCLI実装（`list_dm_users_usecase.go`など）と同様の構造になっている
- [ ] コードスタイルが既存のコードと一致している
- [ ] 適切なコメントが追加されている

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### 変更するファイル
- `server/cmd/server-status/main.go`: 入出力制御のみを担当するように変更

#### 新規作成するファイル
- `server/internal/usecase/cli/server_status_usecase.go`: CLI用usecase層（新規作成）
- `server/internal/usecase/cli/server_status_usecase_test.go`: usecase層のテスト（新規作成）
- `server/internal/service/server_status_service.go`: service層（新規作成）
- `server/internal/service/server_status_service_test.go`: service層のテスト（新規作成）

#### 確認が必要なファイル
- なし（既存の他のコードへの影響はない）

### 7.2 既存機能への影響
- **既存のサーバー**: 影響なし（ポート接続のみで判定、ロジックは維持）
- **既存の機能**: 影響なし（リファクタリングのみ、機能は維持）
- **既存のCLIコマンド**: 影響なし（独立した実装）

### 7.3 テストへの影響
- **既存のテスト**: 影響なし（新規ファイルの追加のみ）
- **新規テスト**: usecase層とservice層のテストを追加

## 8. 実装上の注意事項

### 8.1 アーキテクチャパターンの遵守
- **3層構造**: CLI層 → Usecase層 → Service層の順で呼び出し
- **依存関係の方向**: CLI層がUsecase層に依存、Usecase層がService層に依存
- **インターフェースの使用**: Service層はインターフェースで定義し、テスト容易性を確保

### 8.2 既存ロジックの移行
- **サーバー情報の定義**: `servers`変数の内容をservice層に移行
- **状態確認ロジック**: `checkServerStatus`関数の内容をservice層に移行
- **並列実行ロジック**: `checkAllServers`関数の内容をservice層に移行
- **表示ロジック**: `printResults`関数はCLI層に維持

### 8.3 テストの実装
- **Usecase層のテスト**: Service層をモック化してテスト
- **Service層のテスト**: 実際のポート接続をテスト（統合テスト的なアプローチ）
- **テーブル駆動テスト**: 既存のテストパターンに従う

### 8.4 エラーハンドリング
- **Service層**: エラーをそのまま返す
- **Usecase層**: エラーにコンテキストを追加して返す
- **CLI層**: エラーを適切に表示する

## 9. 参考情報

### 9.1 既存実装の参考
- **CLI実装例**: `server/cmd/list-dm-users/main.go`
- **Usecase実装例**: `server/internal/usecase/cli/list_dm_users_usecase.go`
- **Service実装例**: `server/internal/service/dm_user_service.go`

### 9.2 アーキテクチャドキュメント
- `.kiro/steering/tech.md`: 技術スタックとアーキテクチャの詳細
- `.kiro/steering/structure.md`: ファイル組織とコードパターン

### 9.3 関連ドキュメント
- `.kiro/specs/0077-listapp/requirements.md`: 起動サーバー一覧表示機能の要件定義書（既存機能の参考）
- Issue #161: 本機能の元となる要望

### 9.4 技術スタック
- **言語**: Go 1.21+
- **標準ライブラリ**: `net`パッケージ（TCP接続）
- **テスト**: `testing`（標準ライブラリ）、`github.com/stretchr/testify`（アサーション、モック）
