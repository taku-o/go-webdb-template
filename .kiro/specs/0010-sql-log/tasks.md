# SQLログ出力機能実装タスク一覧

## 概要
SQLログ出力機能の実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: 設定の拡張

#### - [ ] タスク 1.1: LoggingConfig構造体の拡張
**目的**: SQLログ設定項目を追加

**作業内容**:
- `server/internal/config/config.go`の`LoggingConfig`構造体に以下のフィールドを追加:
  ```go
  type LoggingConfig struct {
      Level          string `mapstructure:"level"`
      Format         string `mapstructure:"format"`
      Output         string `mapstructure:"output"`
      OutputDir      string `mapstructure:"output_dir"`        // 既存
      SQLLogEnabled  bool   `mapstructure:"sql_log_enabled"`   // 新規追加（オプション）
      SQLLogOutputDir string `mapstructure:"sql_log_output_dir"` // 新規追加（オプション）
  }
  ```

**受け入れ基準**:
- `LoggingConfig`構造体に`SQLLogEnabled`と`SQLLogOutputDir`フィールドが追加されている
- コンパイルエラーがない

---

#### - [ ] タスク 1.2: デフォルト値の設定
**目的**: SQLログ設定のデフォルト値を設定

**作業内容**:
- `server/internal/config/config.go`の`Load()`関数で、以下のデフォルト値を設定:
  - `SQLLogOutputDir`が空の場合は`OutputDir`と同じ値を使用（既存のデフォルト値設定ロジックを活用）
  - `SQLLogEnabled`は環境に応じて自動判定（develop/staging: true, production: false）

**受け入れ基準**:
- `SQLLogOutputDir`が空の場合、`OutputDir`と同じ値が設定される
- 環境判定ロジックが正しく動作する
- コンパイルエラーがない

---

#### - [ ] タスク 1.3: 設定ファイルの更新（develop環境）
**目的**: develop環境の設定ファイルにSQLログ設定を追加（オプション）

**作業内容**:
- `config/develop/config.yaml`の`logging`セクションに以下の項目を追加（オプション）:
  ```yaml
  logging:
    level: debug
    format: json
    output: file
    output_dir: logs
    sql_log_enabled: true  # オプション（環境別に自動判定する場合は省略可能）
    sql_log_output_dir: logs  # オプション（output_dirと同じ場合は省略可能）
  ```

**受け入れ基準**:
- `config/develop/config.yaml`にSQLログ設定が追加されている（オプション項目）
- YAML形式が正しい

---

#### - [ ] タスク 1.4: 設定ファイルの更新（staging環境）
**目的**: staging環境の設定ファイルにSQLログ設定を追加（オプション）

**作業内容**:
- `config/staging/config.yaml`の`logging`セクションに以下の項目を追加（オプション）:
  ```yaml
  logging:
    level: debug
    format: json
    output: file
    output_dir: logs
    sql_log_enabled: true  # オプション
    sql_log_output_dir: logs  # オプション
  ```

**受け入れ基準**:
- `config/staging/config.yaml`にSQLログ設定が追加されている（オプション項目）
- YAML形式が正しい

---

#### - [ ] タスク 1.5: 設定ファイルの更新（production環境）
**目的**: production環境の設定ファイルにSQLログ設定を追加（オプション、コメント付き）

**作業内容**:
- `config/production/config.yaml.example`の`logging`セクションに以下の項目を追加（コメント付きで説明）:
  ```yaml
  logging:
    level: info
    format: json
    output: file
    output_dir: /var/log/go-webdb-template
    # SQLログは本番環境では無効化されます（環境判定により自動的に無効化）
    # sql_log_enabled: false  # オプション（本番環境では自動的にfalse）
    # sql_log_output_dir: /var/log/go-webdb-template  # オプション
  ```

**受け入れ基準**:
- `config/production/config.yaml.example`にSQLログ設定が追加されている（コメント付き）
- コメントで説明が記載されている

---

### Phase 2: SQL Logger実装

#### - [ ] タスク 2.1: logger.goファイルの作成
**目的**: GORM Logger実装ファイルを作成

**作業内容**:
- `server/internal/db/logger.go`を作成
- パッケージ宣言とインポート文を追加:
  ```go
  package db

  import (
      "context"
      "fmt"
      "io"
      "os"
      "path/filepath"
      "regexp"
      "strings"
      "time"

      "github.com/example/go-webdb-template/internal/config"
      "github.com/sirupsen/logrus"
      "gorm.io/gorm/logger"
      "gopkg.in/natefinch/lumberjack.v2"
  )
  ```

**受け入れ基準**:
- `server/internal/db/logger.go`ファイルが作成されている
- パッケージ宣言とインポート文が正しい
- コンパイルエラーがない

---

#### - [ ] タスク 2.2: SQLLogger構造体の定義
**目的**: SQLLogger構造体を定義

**作業内容**:
- `SQLLogger`構造体を定義:
  ```go
  type SQLLogger struct {
      logger    *logrus.Logger
      writer    io.WriteCloser
      shardID   int
      driver    string
      logLevel  logger.LogLevel
      outputDir string
  }
  ```

**受け入れ基準**:
- `SQLLogger`構造体が定義されている
- コンパイルエラーがない

---

#### - [ ] タスク 2.3: NewSQLLogger関数の実装
**目的**: SQL Loggerの初期化関数を実装

**作業内容**:
- `NewSQLLogger`関数を実装:
  - 環境判定（production環境ではnilを返す）
  - ログ出力先ディレクトリの取得と作成
  - `lumberjack`ライブラリの設定（日付別ファイル分割）
  - `logrus`ライブラリの設定
  - `SQLLogger`インスタンスの作成と返却

**受け入れ基準**:
- `NewSQLLogger`関数が実装されている
- production環境ではnilを返す
- develop/staging環境ではLoggerインスタンスを返す
- ログディレクトリが自動作成される
- コンパイルエラーがない

---

#### - [ ] タスク 2.4: GORM Loggerインターフェースの実装（LogMode）
**目的**: LogModeメソッドを実装

**作業内容**:
- `LogMode`メソッドを実装:
  ```go
  func (l *SQLLogger) LogMode(level logger.LogLevel) logger.Interface {
      newLogger := *l
      newLogger.logLevel = level
      return &newLogger
  }
  ```

**受け入れ基準**:
- `LogMode`メソッドが実装されている
- ログレベルが正しく設定される
- コンパイルエラーがない

---

#### - [ ] タスク 2.5: GORM Loggerインターフェースの実装（Info, Warn, Error）
**目的**: Info、Warn、Errorメソッドを実装

**作業内容**:
- `Info`メソッドを実装（ログレベルがInfo以上の場合にログ出力）
- `Warn`メソッドを実装（ログレベルがWarn以上の場合にログ出力）
- `Error`メソッドを実装（ログレベルがError以上の場合にログ出力）
- nilチェックを追加（production環境でnilの場合の処理）

**受け入れ基準**:
- `Info`、`Warn`、`Error`メソッドが実装されている
- ログレベルに応じて適切にログが出力される
- nilチェックが正しく動作する
- コンパイルエラーがない

---

#### - [ ] タスク 2.6: GORM Loggerインターフェースの実装（Trace）
**目的**: Traceメソッドを実装（SQLクエリトレース）

**作業内容**:
- `Trace`メソッドを実装:
  - SQLクエリと結果件数を取得（`fc`関数から）
  - 実行時間を計算（`begin`から現在時刻まで）
  - テーブル名を抽出（`extractTableName`関数を使用）
  - ログエントリを作成（シャードID、ドライバー名、テーブル名、結果件数、SQLクエリ、実行時間）
  - ログ出力
- nilチェックを追加

**受け入れ基準**:
- `Trace`メソッドが実装されている
- SQLクエリ情報が正しく取得される
- 実行時間が正しく計算される
- ログが正しく出力される
- コンパイルエラーがない

---

#### - [ ] タスク 2.7: Closeメソッドの実装
**目的**: Loggerのクローズ処理を実装

**作業内容**:
- `Close`メソッドを実装:
  ```go
  func (l *SQLLogger) Close() error {
      if l == nil || l.writer == nil {
          return nil
      }
      return l.writer.Close()
  }
  ```

**受け入れ基準**:
- `Close`メソッドが実装されている
- nilチェックが正しく動作する
- コンパイルエラーがない

---

#### - [ ] タスク 2.8: DSNフィルタリング機能の実装
**目的**: DSN文字列からpasswordなどの機密情報をフィルタリング

**作業内容**:
- `filterDSN`関数を実装:
  - PostgreSQL DSN: `password=[^ ]+` → `password=***`に置換
  - MySQL DSN: `:[^@]+@` → `:***@`に置換
  - SQLite DSN: そのまま返す（通常password情報は含まれない）

**受け入れ基準**:
- `filterDSN`関数が実装されている
- PostgreSQL DSNのpasswordが正しくマスクされる
- MySQL DSNのpasswordが正しくマスクされる
- SQLite DSNは変更されない
- コンパイルエラーがない

---

#### - [ ] タスク 2.9: テーブル名抽出機能の実装
**目的**: SQLクエリからテーブル名を抽出

**作業内容**:
- `extractTableName`関数を実装:
  - SELECT文: `FROM table_name`から抽出
  - INSERT文: `INTO table_name`から抽出
  - UPDATE文: `UPDATE table_name`から抽出
  - DELETE文: `FROM table_name`から抽出
  - 抽出失敗時は"unknown"を返す

**受け入れ基準**:
- `extractTableName`関数が実装されている
- SELECT文からテーブル名が正しく抽出される
- INSERT文からテーブル名が正しく抽出される
- UPDATE文からテーブル名が正しく抽出される
- DELETE文からテーブル名が正しく抽出される
- 抽出失敗時は"unknown"が返される
- コンパイルエラーがない

---

#### - [ ] タスク 2.10: SQLTextFormatterの実装
**目的**: SQLログのテキストフォーマッターを実装

**作業内容**:
- `SQLTextFormatter`構造体を定義
- `Format`メソッドを実装:
  - ログエントリを以下の形式でフォーマット:
    ```
    [YYYY-MM-DD HH:MM:SS] [SHARD_ID] [DRIVER] [TABLE] ROWS_AFFECTED | SQL_QUERY | DURATION_MS
    ```
  - `logrus.Formatter`インターフェースを実装

**受け入れ基準**:
- `SQLTextFormatter`構造体が定義されている
- `Format`メソッドが実装されている
- ログフォーマットが要件定義書の形式と一致している
- コンパイルエラーがない

---

### Phase 3: GORM接続作成処理への統合

#### - [ ] タスク 3.1: createGORMConnection関数の拡張
**目的**: GORM接続作成時にSQL Loggerを設定

**作業内容**:
- `server/internal/db/connection.go`の`createGORMConnection`関数を拡張:
  - 関数シグネチャに`loggingCfg *config.LoggingConfig`と`env string`パラメータを追加
  - `NewSQLLogger`関数を呼び出してSQL Loggerを作成
  - Logger作成エラーは警告のみ（SQLクエリ実行には影響しない）
  - `gorm.Config`にLoggerを設定
  - 既存の接続プール設定は変更しない

**受け入れ基準**:
- `createGORMConnection`関数が拡張されている
  - SQL Loggerが正しく設定される
  - Logger作成エラーが警告として処理される
  - 既存の接続プール設定が維持されている
  - コンパイルエラーがない

---

#### - [ ] タスク 3.2: NewGORMConnection関数の拡張
**目的**: GORM接続作成時にSQL Loggerを設定

**作業内容**:
- `server/internal/db/connection.go`の`NewGORMConnection`関数を拡張:
  - 関数シグネチャに`loggingCfg *config.LoggingConfig`と`env string`パラメータを追加
  - `createGORMConnection`関数呼び出し時に`loggingCfg`と`env`を渡す
  - 既存のReader接続作成処理とdbresolver設定は変更しない

**受け入れ基準**:
- `NewGORMConnection`関数が拡張されている
- SQL Loggerが正しく設定される
- 既存のReader接続処理が維持されている
- コンパイルエラーがない

---

#### - [ ] タスク 3.3: NewGORMManager関数の拡張
**目的**: GORMManager初期化時にSQL Loggerを設定

**作業内容**:
- `server/internal/db/manager.go`の`NewGORMManager`関数を拡張:
  - 環境判定（`APP_ENV`環境変数から取得、デフォルトは"develop"）
  - 各シャードの接続作成時に`NewGORMConnection`関数に`loggingCfg`と`env`を渡す
  - 既存のシャーディングロジックは変更しない

**受け入れ基準**:
- `NewGORMManager`関数が拡張されている
- 環境判定が正しく動作する
- すべてのシャードに対してSQL Loggerが設定される
- 既存のシャーディングロジックが維持されている
- コンパイルエラーがない

---

### Phase 4: テスト実装

#### - [ ] タスク 4.1: SQLLoggerの単体テスト作成
**目的**: SQLLoggerの基本機能をテスト

**作業内容**:
- `server/internal/db/logger_test.go`を作成
- `NewSQLLogger`関数のテスト:
  - 正常系: Loggerの作成成功
  - 異常系: ディレクトリ作成失敗
  - 環境判定: production環境ではnilを返す
- `LogMode`メソッドのテスト
- `Info`、`Warn`、`Error`メソッドのテスト
- `Trace`メソッドのテスト

**受け入れ基準**:
- `logger_test.go`が作成されている
- すべてのテストが正常に実行される
- テストカバレッジが適切である

---

#### - [ ] タスク 4.2: DSNフィルタリングの単体テスト作成
**目的**: DSNフィルタリング機能をテスト

**作業内容**:
- `filterDSN`関数のテストを追加:
  - PostgreSQL DSNのフィルタリングテスト
  - MySQL DSNのフィルタリングテスト
  - SQLite DSNのフィルタリングテスト（変更なし）

**受け入れ基準**:
- DSNフィルタリングのテストが作成されている
- すべてのテストが正常に実行される
- passwordが正しくマスクされることを確認

---

#### - [ ] タスク 4.3: テーブル名抽出の単体テスト作成
**目的**: テーブル名抽出機能をテスト

**作業内容**:
- `extractTableName`関数のテストを追加:
  - SELECT文からのテーブル名抽出テスト
  - INSERT文からのテーブル名抽出テスト
  - UPDATE文からのテーブル名抽出テスト
  - DELETE文からのテーブル名抽出テスト
  - 抽出失敗時の"unknown"返却テスト

**受け入れ基準**:
- テーブル名抽出のテストが作成されている
- すべてのテストが正常に実行される
- テーブル名が正しく抽出されることを確認

---

#### - [ ] タスク 4.4: GORM接続作成の統合テスト作成
**目的**: SQL Logger設定付きGORM接続作成をテスト

**作業内容**:
- 統合テストを作成:
  - Logger設定付きGORM接続の作成テスト
  - 複数シャードへのLogger設定テスト
  - 環境別のLogger有効/無効化テスト

**受け入れ基準**:
- 統合テストが作成されている
- すべてのテストが正常に実行される
- SQL Loggerが正しく設定されることを確認

---

#### - [ ] タスク 4.5: SQLログ出力の統合テスト作成
**目的**: 実際のSQLクエリ実行時のログ出力をテスト

**作業内容**:
- 統合テストを作成:
  - 実際のSQLクエリ実行時のログ出力テスト
  - ログファイルの作成確認テスト
  - ログフォーマットの確認テスト
  - ログ内容の検証テスト

**受け入れ基準**:
- 統合テストが作成されている
- すべてのテストが正常に実行される
- SQLログが正しく出力されることを確認

---

### Phase 5: 動作確認とドキュメント更新

#### - [ ] タスク 5.1: 開発環境での動作確認
**目的**: 開発環境でSQLログが正しく出力されることを確認

**作業内容**:
- 開発環境（`APP_ENV=develop`）でアプリケーションを起動
- データベース操作（SELECT、INSERT、UPDATE、DELETE）を実行
- `logs/sql-YYYY-MM-DD.log`ファイルが作成されることを確認
- ログ内容を確認:
  - シャードIDが正しく記録されている
  - ドライバー名が正しく記録されている
  - テーブル名が正しく記録されている
  - SQLクエリが正しく記録されている
  - 結果件数が正しく記録されている
  - 実行時間が正しく記録されている
  - ログフォーマットが要件定義書の形式と一致している

**受け入れ基準**:
- 開発環境でSQLログが正しく出力される
- ログファイルが正しく作成される
- ログ内容が要件定義書の要件を満たしている

---

#### - [ ] タスク 5.2: ステージング環境での動作確認
**目的**: ステージング環境でSQLログが正しく出力されることを確認

**作業内容**:
- ステージング環境（`APP_ENV=staging`）でアプリケーションを起動
- データベース操作を実行
- SQLログが正しく出力されることを確認

**受け入れ基準**:
- ステージング環境でSQLログが正しく出力される
- ログファイルが正しく作成される

---

#### - [ ] タスク 5.3: 本番環境での動作確認
**目的**: 本番環境でSQLログが出力されないことを確認

**作業内容**:
- 本番環境（`APP_ENV=production`）でアプリケーションを起動
- データベース操作を実行
- SQLログファイルが作成されないことを確認
- ログディレクトリにSQLログファイルが存在しないことを確認

**受け入れ基準**:
- 本番環境でSQLログが出力されない
- SQLログファイルが作成されない

---

#### - [ ] タスク 5.4: シャーディング対応の確認
**目的**: すべてのシャードでSQLログが正しく出力されることを確認

**作業内容**:
- すべてのシャード（shard1, shard2, shard3, shard4）に対してデータベース操作を実行
- 各シャードのSQLクエリがログに記録されることを確認
- シャードIDが正しく記録されることを確認

**受け入れ基準**:
- すべてのシャードでSQLログが正しく出力される
- シャードIDが正しく記録される

---

#### - [ ] タスク 5.5: セキュリティ確認（DSNフィルタリング）
**目的**: DSN文字列からpasswordが正しくフィルタリングされることを確認

**作業内容**:
- PostgreSQL DSNを使用した接続でSQLログを確認
- MySQL DSNを使用した接続でSQLログを確認
- ログにpasswordが含まれていないことを確認
- passwordが`***`にマスクされていることを確認

**受け入れ基準**:
- DSN文字列に含まれるpasswordがログに出力されない
- passwordが正しくマスクされる

---

#### - [ ] タスク 5.6: 既存機能への影響確認
**目的**: 既存の機能が正常に動作することを確認

**作業内容**:
- 既存のアクセスログ機能（0008-log-strategy）が正常に動作することを確認
- 既存のログ出力機能（標準`log`パッケージ）が正常に動作することを確認
- 既存のデータベース接続処理が正常に動作することを確認
- 既存のシャーディングロジックが正常に動作することを確認

**受け入れ基準**:
- 既存のすべての機能が正常に動作する
- 既存のテストがすべて正常に実行される

---

#### - [ ] タスク 5.7: パフォーマンス確認
**目的**: SQLログ出力がパフォーマンスに大きな影響を与えないことを確認

**作業内容**:
- SQLクエリ実行時間を測定（SQLログ有効時と無効時）
- パフォーマンスへの影響を確認
- 本番環境ではSQLログを出力しないため、パフォーマンスへの影響はないことを確認

**受け入れ基準**:
- SQLログ出力がパフォーマンスに大きな影響を与えない
- 本番環境ではパフォーマンスへの影響がない

---

#### - [ ] タスク 5.8: エラーハンドリング確認
**目的**: エラーが発生してもSQLクエリ実行に影響しないことを確認

**作業内容**:
- Logger初期化エラーの場合、警告が出力されて動作を継続することを確認
- ログファイル書き込みエラーの場合、SQLクエリ実行が継続することを確認
- エラーログが標準エラー出力に記録されることを確認

**受け入れ基準**:
- エラーが発生してもSQLクエリ実行に影響しない
- エラーログが適切に記録される

---

#### - [ ] タスク 5.9: 日付別ファイル分割の確認
**目的**: 日付が変わったら自動的に新しいログファイルに切り替わることを確認

**作業内容**:
- 日付を跨いでSQLクエリを実行
- 新しい日付のログファイルが作成されることを確認
- 前日のログファイルが正しく保持されることを確認

**受け入れ基準**:
- 日付が変わったら自動的に新しいログファイルに切り替わる
- 前日のログファイルが正しく保持される

---

#### - [ ] タスク 5.10: 設定ファイルの動作確認
**目的**: 設定ファイルでログ出力先を変更できることを確認

**作業内容**:
- 設定ファイルで`sql_log_output_dir`を変更
- 指定したディレクトリにログファイルが作成されることを確認
- デフォルト値（`logs`ディレクトリ）が正しく動作することを確認

**受け入れ基準**:
- 設定ファイルでログ出力先を変更できる
- デフォルト値が正しく動作する

---

#### - [ ] タスク 5.11: ドキュメント更新（必要に応じて）
**目的**: 必要に応じてドキュメントを更新

**作業内容**:
- README.mdまたは関連ドキュメントにSQLログ出力機能の説明を追加（必要に応じて）
- 設定方法の説明を追加（必要に応じて）

**受け入れ基準**:
- ドキュメントが適切に更新されている（必要に応じて）

---

## 実装順序の推奨

1. **Phase 1: 設定の拡張** - 設定構造体と設定ファイルの更新
2. **Phase 2: SQL Logger実装** - コア機能の実装
3. **Phase 3: GORM接続作成処理への統合** - 既存コードへの統合
4. **Phase 4: テスト実装** - テストコードの作成
5. **Phase 5: 動作確認とドキュメント更新** - 最終確認とドキュメント更新

## 注意事項

- 既存のアクセスログ機能（0008-log-strategy）との共存を確認
- 既存のデータベース接続処理への影響を最小化
- エラーが発生してもSQLクエリ実行に影響を与えない
- 本番環境ではSQLログを出力しない（環境判定）
- DSN文字列からpasswordなどの機密情報を確実にフィルタリング

