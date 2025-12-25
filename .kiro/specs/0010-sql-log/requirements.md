# SQLログ出力機能要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #16
- **Issueタイトル**: データベースのSQLクエリログを出力したい
- **Feature名**: 0010-sql-log
- **作成日**: 2025-01-27

### 1.2 目的
データベースに発行するSQLクエリをログに出力する機能を実装する。
これにより、開発環境とステージング環境でのSQLクエリの可視性を向上させ、デバッグとパフォーマンス分析を容易にする。

### 1.3 スコープ
- GORMのLoggerインターフェースを実装したカスタムSQL Loggerの作成
- 開発環境（develop）とステージング環境（staging）でのSQLログ出力
- 本番環境（production）でのSQLログ出力の無効化
- 日付別ログファイル分割機能
- ログ出力先の設定機能（デフォルト: `logs`ディレクトリ）
- DSN文字列からの機密情報（password等）のフィルタリング
- シャーディング構成（複数シャード）への対応

## 2. 背景・現状分析

### 2.1 現在の実装
- **ORM**: GORM（`gorm.io/gorm`）を使用
- **データベース接続**: `server/internal/db/connection.go`でGORM接続を管理
- **シャーディング**: 複数のシャード（shard1, shard2, shard3, shard4）に対応
- **ログ出力**: アクセスログ機能（0008-log-strategy）が実装済み
  - `server/internal/logging/access_logger.go`: アクセスログ出力機能
  - `lumberjack`と`logrus`ライブラリを使用
  - 日付別ファイル分割機能あり
- **GORM Logger**: 現在、GORMのLoggerは設定されていない（デフォルトのLoggerを使用）
- **設定**: `server/internal/config/config.go`に`LoggingConfig`構造体が存在
  - `OutputDir`フィールドが存在（アクセスログ用）

### 2.2 課題点
1. **SQLクエリの可視性不足**: データベースに発行されるSQLクエリが確認できないため、デバッグが困難
2. **パフォーマンス分析の困難**: どのSQLクエリが実行されているか、実行時間はどの程度かを把握できない
3. **環境別制御の不足**: 開発環境と本番環境でSQLログの出力を切り替える機能がない
4. **機密情報の漏洩リスク**: DSN文字列に含まれるpasswordなどの機密情報がログに出力される可能性がある

### 2.3 本実装による改善点
1. **SQLクエリの可視化**: すべてのSQLクエリをログに記録し、デバッグを容易にする
2. **パフォーマンス分析**: SQLクエリの実行時間や結果件数を記録し、パフォーマンス分析を可能にする
3. **環境別制御**: 開発環境とステージング環境でのみSQLログを出力し、本番環境では出力しない
4. **セキュリティ向上**: DSN文字列からpasswordなどの機密情報をフィルタリングしてログに出力
5. **既存機能との統合**: 既存のアクセスログ機能と統合し、同じ`logs`ディレクトリを使用

## 3. 機能要件

### 3.1 SQLログ出力機能

#### 3.1.1 基本機能
- GORMのLoggerインターフェース（`gorm.Interface`）を実装したカスタムLoggerを作成
- すべてのSQLクエリ（SELECT、INSERT、UPDATE、DELETE等）をログに記録
- 擬似的なSQL（GORMが生成するSQL）でもOK
- すべてのシャード（shard1, shard2, shard3, shard4）に対してLoggerを設定

#### 3.1.2 ログ出力内容
以下の情報をログに記録する：
- **接続先データベース**: シャードID、ドライバー名（sqlite3/postgres/mysql）
- **テーブル名**: SQLクエリから抽出可能な場合（GORMの情報から取得）
- **SQLクエリ**: パラメータがバインドされた実際のSQLクエリ
- **SQL結果の件数**: 取得可能な場合（`RowsAffected`、`Rows`等から取得）
- **実行時間**: SQLクエリの実行時間（ミリ秒またはマイクロ秒）
- **タイムスタンプ**: クエリ実行日時

#### 3.1.3 ログフォーマット
テキスト形式で以下の情報を1行に記録する：
```
[YYYY-MM-DD HH:MM:SS] [SHARD_ID] [DRIVER] [TABLE] ROWS_AFFECTED | SQL_QUERY | DURATION_MS
```

例:
```
[2025-01-27 14:30:45] [1] [sqlite3] [users] 1 | SELECT * FROM users WHERE id = ? | 2.5ms
[2025-01-27 14:30:46] [2] [sqlite3] [posts] 1 | INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?) | 3.2ms
```

### 3.2 環境別制御

#### 3.2.1 環境判定
- `APP_ENV`環境変数または設定ファイルから環境を判定
- 環境別の動作:
  - **開発環境（develop）**: SQLログを出力
  - **ステージング環境（staging）**: SQLログを出力
  - **本番環境（production）**: SQLログを出力しない

#### 3.2.2 実装方法
- `config.Load()`で取得した設定の`APP_ENV`を確認
- production環境の場合はLoggerを設定しない、またはLoggerを無効化する設定を適用

### 3.3 セキュリティ要件

#### 3.3.1 機密情報のフィルタリング
- DSN文字列に含まれるpasswordなどの機密情報をフィルタリング
- ログに出力する前にDSNからpasswordを除去またはマスク
- フィルタリング対象:
  - `password=xxx` → `password=***` または除去
  - PostgreSQL DSN: `password=xxx` → `password=***`
  - MySQL DSN: `password=xxx` → `password=***`
  - SQLite DSN: password情報は通常含まれないが、念のためチェック

#### 3.3.2 実装方法
- DSN文字列をパースしてpassword部分を検出
- 正規表現またはURLパースを使用してpasswordをマスク
- ログ出力時にマスクされたDSNを使用

### 3.4 ログファイル管理

#### 3.4.1 ログ出力先
- ログ出力先: 設定ファイルで指定（デフォルト: `logs`ディレクトリ）
- 既存のアクセスログ機能（0008-log-strategy）と同じ`logs`ディレクトリを使用
- ディレクトリが存在しない場合は自動作成

#### 3.4.2 ログファイル名
- 形式: `sql-YYYY-MM-DD.log`
- 例: `sql-2025-01-27.log`
- 日付が変わったら自動的に新しいファイルに切り替える

#### 3.4.3 日付別ファイル分割
- **既存の`lumberjack`ライブラリを活用**: 0008-log-strategyで使用済み
- ファイル名パターンに日付フォーマット（`2006-01-02`）を含める
- ライブラリが自動的に日付変更を検知して新しいファイルに切り替える
- サーバーのローカルタイムゾーンを使用

### 3.5 設定機能

#### 3.5.1 設定ファイル
- `config/{env}/config.yaml`の`logging`セクションにSQLログ設定を追加
- 設定項目:
  - `sql_log_enabled`: 有効/無効の切り替え（オプション、環境別に自動判定も可能）
  - `sql_log_output_dir`: ログ出力先ディレクトリ（オプション、既存の`output_dir`を流用可能）

#### 3.5.2 設定例
```yaml
logging:
  level: debug
  format: json
  output: file
  output_dir: logs  # アクセスログとSQLログの共通出力先
  sql_log_enabled: true  # オプション（環境別に自動判定する場合は省略可能）
  sql_log_output_dir: logs  # オプション（output_dirと同じ場合は省略可能）
```

#### 3.5.3 デフォルト値
- `sql_log_enabled`: 環境に応じて自動判定（develop/staging: true, production: false）
- `sql_log_output_dir`: `output_dir`と同じ値を使用（デフォルト: `logs`）

## 4. 非機能要件

### 4.1 パフォーマンス
- SQLログ出力がSQLクエリ実行のボトルネックにならないよう、非同期書き込みまたはバッファリングを検討
- ログファイルへの書き込みエラーが発生しても、SQLクエリ実行には影響を与えない
- 本番環境ではSQLログを出力しないため、パフォーマンスへの影響はない

### 4.2 エラーハンドリング
- ログファイルの作成・書き込みエラーが発生した場合、標準エラー出力にエラーメッセージを出力
- ログファイルへの書き込みに失敗しても、SQLクエリ実行は継続する
- Loggerの初期化に失敗した場合、警告を出力してLoggerなしで動作を継続

### 4.3 既存機能への影響
- 既存のアクセスログ機能（0008-log-strategy）との共存
- 既存のログ出力（標準`log`パッケージによる出力）は維持
- 既存のデータベース接続処理への影響を最小化
- 既存のシャーディングロジックへの影響なし

### 4.4 シャーディング対応
- すべてのシャード（shard1, shard2, shard3, shard4）に対してLoggerを設定
- 各シャードのログを同じファイルに出力（シャードIDで識別可能）
- または、シャード別にログファイルを分割することも検討可能（将来の拡張）

## 5. 制約事項

### 5.1 技術的制約
- **GORM Loggerインターフェース**: GORMの`gorm.Interface`を実装する必要がある
- **既存ライブラリの活用**: 既存の`lumberjack`と`logrus`ライブラリを活用
- **既存の設定構造体**: `LoggingConfig`構造体を拡張して設定を追加
- **既存の接続処理**: `server/internal/db/connection.go`の`createGORMConnection`関数を拡張
- **ログフォーマット**: テキスト形式とする（将来的にJSON形式への拡張も検討可能）

### 5.2 プロジェクト制約
- 既存のアーキテクチャ（レイヤードアーキテクチャ）を維持
- 既存の設定読み込み機能（`internal/config`）を変更しない
- 既存のデータベース接続処理への影響を最小化
- 既存のシャーディングロジックへの影響なし

### 5.3 ディレクトリ構造
- ログファイルは`logs`ディレクトリ（プロジェクトルート）に出力
- 既存のアクセスログ機能と同じディレクトリを使用
- `.gitignore`で`logs/*`と`logs/**`が既に除外されている（変更不要）

### 5.4 設定ファイル形式
- YAML形式の設定ファイルを使用
- 環境別設定（develop/staging/production）を維持

## 6. 受け入れ基準

### 6.1 機能要件
- [ ] 開発環境（develop）でSQLログが出力される
- [ ] ステージング環境（staging）でSQLログが出力される
- [ ] 本番環境（production）でSQLログが出力されない
- [ ] すべてのSQLクエリ（SELECT、INSERT、UPDATE、DELETE等）がログに記録される
- [ ] ログに必要な情報（シャードID、ドライバー名、テーブル名、SQLクエリ、結果件数、実行時間）が記録される
- [ ] passwordなどの機密情報がログに出力されない（DSNからpasswordがフィルタリングされる）
- [ ] ログファイル名に日付（YYYY-MM-DD形式）が含まれる
- [ ] 日付が変わったら自動的に新しいログファイルに切り替わる
- [ ] 設定ファイルでログ出力先を変更できる
- [ ] デフォルトで`logs`ディレクトリにログが出力される
- [ ] `logs`ディレクトリが存在しない場合は自動作成される
- [ ] すべてのシャード（shard1, shard2, shard3, shard4）に対してLoggerが設定される

### 6.2 非機能要件
- [ ] SQLログ出力がSQLクエリ実行のパフォーマンスに大きな影響を与えない
- [ ] ログファイルへの書き込みエラーが発生してもSQLクエリ実行は継続する
- [ ] ログファイルの作成・書き込みエラーは標準エラー出力に記録される
- [ ] 既存のアクセスログ機能（0008-log-strategy）が正常に動作する
- [ ] 既存のログ出力機能（標準`log`パッケージ）が正常に動作する
- [ ] 既存のデータベース接続処理が正常に動作する

### 6.3 設定
- [ ] `config/{env}/config.yaml`の`logging`セクションにSQLログ設定が追加されている（オプション項目）
- [ ] 設定ファイルでログ出力先を変更できる
- [ ] デフォルト値（`logs`ディレクトリ）が正しく動作する
- [ ] 環境別の自動判定（develop/staging: 有効、production: 無効）が正しく動作する

### 6.4 セキュリティ
- [ ] DSN文字列に含まれるpasswordがログに出力されない
- [ ] PostgreSQL DSNのpasswordがマスクされる
- [ ] MySQL DSNのpasswordがマスクされる
- [ ] SQLite DSNにpasswordが含まれていないことを確認（念のため）

## 7. 影響範囲

### 7.1 新規追加が必要なディレクトリ・ファイル

#### ファイル
- `server/internal/db/logger.go`: GORM Logger実装（新規パッケージ内に追加）
  - `SQLLogger`構造体: GORMの`gorm.Interface`を実装
  - `NewSQLLogger`関数: Loggerの初期化
  - DSNフィルタリング機能
  - ログ出力処理

### 7.2 変更が必要なファイル

#### 設定ファイル
- `server/internal/config/config.go`: `LoggingConfig`構造体にSQLログ設定を追加
  - `SQLLogEnabled bool`: SQLログの有効/無効（オプション）
  - `SQLLogOutputDir string`: SQLログ出力先ディレクトリ（オプション）

#### コードファイル
- `server/internal/db/connection.go`: GORM接続作成時にLoggerを設定
  - `createGORMConnection`関数: Loggerの設定を追加
  - `NewGORMConnection`関数: Loggerの設定を追加
- `server/cmd/server/main.go`: SQL Loggerの初期化（オプション、設定から自動的に適用される場合は不要）
- `server/cmd/admin/main.go`: SQL Loggerの初期化（オプション、設定から自動的に適用される場合は不要）

#### 設定ファイル（YAML）
- `config/develop/config.yaml`: SQLログ設定を追加（オプション）
- `config/staging/config.yaml`: SQLログ設定を追加（オプション）
- `config/production/config.yaml.example`: SQLログ設定を追加（オプション、コメントで説明）

#### 依存関係
- `server/go.mod`: 既存の`lumberjack`と`logrus`ライブラリを使用（追加不要）

### 7.3 変更不要なファイル
- `server/internal/logging/access_logger.go`: アクセスログ機能は変更不要
- `server/internal/db/manager.go`: Managerの実装は変更不要
- `.gitignore`: `logs`ディレクトリの除外設定は既に存在（変更不要）

### 7.4 削除されるファイル
なし

## 8. 実装上の注意事項

### 8.1 GORM Loggerインターフェースの実装
- GORMの`gorm.Interface`を実装する必要がある
- 主要なメソッド:
  - `LogMode(level logger.LogLevel) logger.Interface`: ログレベル設定
  - `Info(ctx context.Context, msg string, data ...interface{})`: 情報ログ
  - `Warn(ctx context.Context, msg string, data ...interface{})`: 警告ログ
  - `Error(ctx context.Context, msg string, data ...interface{})`: エラーログ
  - `Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error)`: SQLクエリトレース
- `Trace`メソッドでSQLクエリの詳細情報を取得できる

### 8.2 DSN文字列からのpassword除去方法
- **PostgreSQL DSN**: `host=xxx port=xxx user=xxx password=xxx dbname=xxx`
  - 正規表現: `password=[^ ]+` → `password=***`に置換
- **MySQL DSN**: `user:password@tcp(host:port)/dbname`
  - 正規表現: `:[^@]+@` → `:***@`に置換
- **SQLite DSN**: 通常password情報は含まれないが、念のためチェック
- URLパースを使用する方法も検討可能

### 8.3 テーブル名の抽出方法
- GORMの`Trace`メソッドの`fc`関数から取得できる情報を活用
- SQLクエリをパースしてテーブル名を抽出（`FROM table_name`、`INTO table_name`等）
- GORMのモデル情報からテーブル名を取得（可能な場合）

### 8.4 結果件数の取得方法
- GORMの`Trace`メソッドの`fc`関数の戻り値（`int64`）から取得
- `RowsAffected`や`Rows`の情報を活用
- SELECTクエリの場合は結果件数を取得できない場合がある

### 8.5 環境判定方法
- `config.Load()`で取得した設定から`APP_ENV`環境変数を確認
- または、設定ファイルから環境情報を取得
- production環境の場合はLoggerを設定しない、または無効化する設定を適用

### 8.6 ログライブラリの活用
- **既存の`lumberjack`ライブラリを活用**: 日付別ファイル分割機能
- **既存の`logrus`ライブラリを活用**: ログ出力機能
- アクセスログ機能（0008-log-strategy）と同じパターンで実装

### 8.7 シャーディング対応
- 各シャードのGORM接続に対して個別にLoggerを設定
- ログにシャードIDを含めて識別可能にする
- すべてのシャードのログを同じファイルに出力（シャードIDで識別）

### 8.8 エラーハンドリング
- Loggerの初期化に失敗した場合、警告を出力してLoggerなしで動作を継続
- ログファイルへの書き込みエラーは標準エラー出力に記録
- エラーが発生してもSQLクエリ実行には影響を与えない

## 9. 参考情報

### 9.1 関連Issue
- GitHub Issue #16: データベースのSQLクエリログを出力したい

### 9.2 既存ドキュメント
- `server/internal/db/connection.go`: GORM接続作成実装
- `server/internal/logging/access_logger.go`: アクセスログ実装（参考）
- `server/internal/config/config.go`: 設定読み込み実装
- `config/develop/config.yaml`: 開発環境設定ファイル
- `config/staging/config.yaml`: ステージング環境設定ファイル
- `config/production/config.yaml.example`: 本番環境設定ファイル例

### 9.3 既存実装
- `server/internal/db/connection.go`: `createGORMConnection`関数、`NewGORMConnection`関数
- `server/internal/db/manager.go`: `NewGORMManager`関数
- `server/internal/logging/access_logger.go`: アクセスログ実装（ログ出力パターンの参考）

### 9.4 GORM Logger
- **GORM Loggerインターフェース**: `gorm.io/gorm/logger.Interface`
- **公式ドキュメント**: https://gorm.io/docs/logger.html
- **カスタムLogger実装例**: GORMのドキュメントに実装例が記載されている

### 9.5 ログライブラリ
- **`logrus`** (github.com/sirupsen/logrus): 既存のアクセスログで使用
  - 構造化ログをサポート
  - テキストフォーマッターとJSONフォーマッターを提供
  - ファイル出力機能が標準でサポート
- **`lumberjack`** (gopkg.in/natefinch/lumberjack.v2): 既存のアクセスログで使用
  - ログローテーション機能を提供する軽量ライブラリ
  - 日付別ファイル分割機能が標準で提供
  - ファイル名に日付フォーマット（`2006-01-02`）を含める設定が可能
  - 自動的なファイル切り替えをサポート

### 9.6 Go言語標準ライブラリ
- `os`パッケージ: https://pkg.go.dev/os
- `time`パッケージ: https://pkg.go.dev/time
- `context`パッケージ: https://pkg.go.dev/context
- `regexp`パッケージ: https://pkg.go.dev/regexp（DSNフィルタリング用）

