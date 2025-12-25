# 設定ファイル分割・リファクタリング要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #8
- **Issueタイトル**: 設定ファイルの分割・リファクタリング
- **Feature名**: 0006-config-refactoring
- **作成日**: 2025-01-27

### 1.2 目的
現在、環境毎に1つのYAMLファイルにまとまっている設定ファイルを、以下の方針でリファクタリングする：
1. 環境別ディレクトリ構造への移行（`config/develop/`, `config/production/`, `config/staging/`）
2. データベース設定の分離（各環境ディレクトリ配下に`database.yaml`を配置）

これにより、設定ファイルの管理性と保守性を向上させる。

### 1.3 スコープ
- 既存の設定ファイル（`config/develop.yaml`, `config/production.yaml.example`, `config/staging.yaml`）の分割と移動
- データベース設定の分離（`database.yaml`）
- 設定読み込みロジック（`server/internal/config/config.go`）の修正
- 既存の設定構造体とAPIの互換性維持

## 2. 背景・現状分析

### 2.1 現在の実装
- **設定ファイル構造**: 環境毎に1つのYAMLファイル（`config/develop.yaml`, `config/production.yaml.example`, `config/staging.yaml`）
- **設定内容**: 各ファイルに以下の設定が含まれる：
  - `server`: サーバー設定（ポート、タイムアウト）
  - `admin`: 管理画面設定（ポート、認証、セッション）
  - `database`: データベース設定（シャード設定、接続情報）
  - `logging`: ロギング設定（レベル、フォーマット、出力先）
  - `cors`: CORS設定（許可オリジン、メソッド、ヘッダー）
- **設定読み込み**: `server/internal/config/config.go`の`Load()`関数でviperを使用して単一ファイルを読み込み

### 2.2 課題点
1. **設定ファイルの肥大化**: 1つのファイルに全設定が集約され、管理が困難
2. **データベース設定の分離不足**: データベース設定が他の設定と混在し、セキュリティ管理が困難
3. **環境別ディレクトリ構造の欠如**: 環境毎の設定ファイルが同一ディレクトリに配置され、整理が不十分
4. **設定ファイルの可読性**: ファイルサイズが大きくなり、必要な設定を見つけにくい

### 2.3 本実装による改善点
1. **明確なディレクトリ構造**: 環境別にディレクトリを分離し、設定ファイルの所在が明確になる
2. **データベース設定の分離**: データベース設定を別ファイルに分離し、セキュリティ管理が容易になる
3. **保守性の向上**: 設定カテゴリごとにファイルを分離し、変更影響範囲が明確になる
4. **設定ファイルの可読性向上**: ファイルサイズが小さくなり、必要な設定を見つけやすくなる

## 3. 機能要件

### 3.1 ディレクトリ構造の変更

#### 3.1.1 新規ディレクトリの作成
以下の3つのディレクトリを作成し、各環境の設定ファイルを移動する：
- `config/develop/`: 開発環境用設定ディレクトリ
- `config/production/`: 本番環境用設定ディレクトリ
- `config/staging/`: ステージング環境用設定ディレクトリ

#### 3.1.2 既存ファイルの移動
- `config/develop.yaml` → `config/develop/config.yaml`（データベース設定を除く）
- `config/production.yaml.example` → `config/production/config.yaml.example`（データベース設定を除く）
- `config/staging.yaml` → `config/staging/config.yaml`（データベース設定を除く）

### 3.2 データベース設定の分離

#### 3.2.1 データベース設定ファイルの作成
各環境ディレクトリ配下に`database.yaml`を作成し、データベース設定を分離する：
- `config/develop/database.yaml`: 開発環境用データベース設定
- `config/production/database.yaml.example`: 本番環境用データベース設定（テンプレート）
- `config/staging/database.yaml`: ステージング環境用データベース設定

#### 3.2.2 データベース設定の内容
各`database.yaml`には、既存のYAMLファイルから`database`セクションの内容を移動する：
- `database.shards`: シャード設定のリスト
  - 各シャードの設定（id, driver, host, port, name, user, password, dsn, writer_dsn, reader_dsns, reader_policy, max_connections, max_idle_connections, connection_max_lifetime）

#### 3.2.3 メイン設定ファイルからの削除
各環境の`config.yaml`から`database`セクションを削除する。

### 3.3 設定読み込みロジックの修正

#### 3.3.1 複数ファイル読み込み対応
`server/internal/config/config.go`の`Load()`関数を修正し、以下の順序で設定ファイルを読み込む：
1. メイン設定ファイル（`config/{env}/config.yaml`）を読み込み
2. データベース設定ファイル（`config/{env}/database.yaml`）を読み込み
3. 両方の設定を統合して`Config`構造体にマッピング

#### 3.3.2 viperの設定パス変更
- `viper.AddConfigPath()`で環境別ディレクトリ（`config/{env}/`）を指定
- メイン設定ファイルとデータベース設定ファイルを順次読み込み
- `viper.MergeInConfig()`を使用してデータベース設定をメイン設定に統合

#### 3.3.3 後方互換性の維持
- 既存の`Config`構造体の定義は変更しない
- 既存の設定取得API（`cfg.Database.Shards`など）はそのまま動作する
- 環境変数による上書き機能（`DB_PASSWORD_SHARD*`）は維持する

## 4. 非機能要件

### 4.1 パフォーマンス
- 設定ファイルの読み込み時間は既存と同等以下を維持
- 複数ファイル読み込みによるオーバーヘッドを最小化

### 4.2 保守性
- 明確なディレクトリ構造を維持
- 設定ファイルの命名規則を統一（`config.yaml`, `database.yaml`）
- コメントやドキュメントで設定構造を説明

### 4.3 セキュリティ
- データベース設定ファイル（特に本番環境）は`.example`ファイルとして提供
- 環境変数による機密情報の上書き機能を維持
- `.gitignore`で実際のデータベース設定ファイルを除外（必要に応じて）

### 4.4 互換性
- 既存の設定構造体（`Config`, `DatabaseConfig`, `ShardConfig`など）との互換性を維持
- 既存の設定取得コードへの影響を最小化
- 環境変数`APP_ENV`による環境切り替え機能を維持

## 5. 制約事項

### 5.1 技術的制約
- viperライブラリの機能範囲内で実装
- Go言語の標準的な設定ファイル読み込みパターンに従う
- YAML形式の設定ファイルを維持

### 5.2 プロジェクト制約
- 既存の設定構造体の定義は変更しない
- 既存のAPIエンドポイントや機能への影響を最小化
- 後方互換性を維持

### 5.3 ディレクトリ構造
- 環境別ディレクトリ（`config/{env}/`）配下に設定ファイルを配置
- 既存の`config/`ディレクトリ構造を活用

## 6. 受け入れ基準

### 6.1 機能要件
- [ ] `config/develop/`, `config/production/`, `config/staging/`ディレクトリが作成されている
- [ ] 各環境の`config.yaml`ファイルが適切なディレクトリに配置されている（データベース設定を除く）
- [ ] 各環境の`database.yaml`ファイルが適切なディレクトリに配置されている
- [ ] `config.go`の`Load()`関数が複数ファイルを正しく読み込める
- [ ] 既存の設定取得コード（`cfg.Database.Shards`など）が正常に動作する
- [ ] 環境変数`APP_ENV`による環境切り替えが正常に動作する
- [ ] 環境変数によるパスワード上書き機能（`DB_PASSWORD_SHARD*`）が正常に動作する

### 6.2 非機能要件
- [ ] 既存のコードベースとの互換性が維持されている
- [ ] 既存のテストが正常に動作する
- [ ] 設定ファイルの読み込み時間が既存と同等以下である
- [ ] ディレクトリ構造が明確で保守しやすい

### 6.3 ドキュメント
- [ ] README.mdに新しい設定ファイル構造が記載されている（必要に応じて）
- [ ] 設定ファイルの配置場所と読み込み順序が明確である

## 7. 影響範囲

### 7.1 新規追加が必要なディレクトリ・ファイル

#### ディレクトリ構造
- `config/develop/`: 開発環境用設定ディレクトリ
- `config/production/`: 本番環境用設定ディレクトリ
- `config/staging/`: ステージング環境用設定ディレクトリ

#### 設定ファイル
- `config/develop/config.yaml`: 開発環境用メイン設定（データベース設定を除く）
- `config/develop/database.yaml`: 開発環境用データベース設定
- `config/production/config.yaml.example`: 本番環境用メイン設定テンプレート（データベース設定を除く）
- `config/production/database.yaml.example`: 本番環境用データベース設定テンプレート
- `config/staging/config.yaml`: ステージング環境用メイン設定（データベース設定を除く）
- `config/staging/database.yaml`: ステージング環境用データベース設定

### 7.2 変更が必要なファイル

#### 設定読み込みロジック
- `server/internal/config/config.go`: `Load()`関数を修正し、複数ファイル読み込みに対応

### 7.3 削除されるファイル
- `config/develop.yaml`: `config/develop/config.yaml`に移動後、削除
- `config/production.yaml.example`: `config/production/config.yaml.example`に移動後、削除
- `config/staging.yaml`: `config/staging/config.yaml`に移動後、削除

### 7.4 ドキュメント更新
- `README.md`: 設定ファイル構造の変更を記載（必要に応じて）
- `.gitignore`: データベース設定ファイルの除外ルールを追加（必要に応じて）

## 8. 実装上の注意事項

### 8.1 viperの複数ファイル読み込み
- `viper.SetConfigName()`でメイン設定ファイル名を指定
- `viper.ReadInConfig()`でメイン設定を読み込み
- `viper.SetConfigName("database")`でデータベース設定ファイル名に変更
- `viper.MergeInConfig()`でデータベース設定をメイン設定に統合
- 統合後、`viper.Unmarshal()`で`Config`構造体にマッピング

### 8.2 設定ファイルのパス解決
- `viper.AddConfigPath()`で環境別ディレクトリ（`config/{env}/`）を指定
- 複数のパスを追加して、実行ディレクトリの違いに対応（既存の実装を維持）

### 8.3 データベース設定の統合
- メイン設定ファイルに`database`セクションがないことを前提とする
- データベース設定ファイルの`database`セクションをメイン設定に統合
- 統合後の構造は既存の`Config`構造体と互換性を保つ

### 8.4 エラーハンドリング
- メイン設定ファイルの読み込みエラーは既存と同様に処理
- データベース設定ファイルの読み込みエラーも適切に処理
- データベース設定ファイルが存在しない場合のフォールバック処理を検討（オプション）

### 8.5 環境変数の処理
- 環境変数による上書き機能（`viper.AutomaticEnv()`）は維持
- データベースパスワードの環境変数上書き（`DB_PASSWORD_SHARD*`）は既存ロジックを維持

## 9. 参考情報

### 9.1 viper公式ドキュメント
- Viper Configuration: https://github.com/spf13/viper
- Viper Multiple Config Files: https://github.com/spf13/viper#reading-multiple-config-files

### 9.2 関連Issue
- GitHub Issue #8: 設定ファイルの分割・リファクタリング

### 9.3 既存ドキュメント
- `server/internal/config/config.go`: 現在の設定読み込み実装
- `config/develop.yaml`: 開発環境設定ファイル
- `config/production.yaml.example`: 本番環境設定テンプレート
- `config/staging.yaml`: ステージング環境設定ファイル

### 9.4 既存実装
- `server/internal/config/config.go`: 設定読み込みロジック
- `server/internal/db/manager.go`: データベース接続管理（設定を使用）

