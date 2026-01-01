# メール送信ログ機能要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #68
- **Issueタイトル**: メール送信ログ
- **Feature名**: 0034-maillog
- **作成日**: 2025-01-27

### 1.2 目的
メール送信時にデバッグ用のログを出力する機能を実装する。
これにより、開発環境とステージング環境でのメール送信内容を記録し、デバッグと動作確認を容易にする。

### 1.3 スコープ
- メール送信時のログ出力機能の実装
- 日付別ログファイル分割機能（`mail-YYYY-MM-DD.log`形式）
- ログ出力先の設定機能（デフォルト: `logs`ディレクトリ）
- 開発環境（develop）とステージング環境（staging）でのログ出力
- 本番環境（production）でのログ出力の無効化
- JSON形式でのログ出力
- メール内容の200文字制限（200文字以上は切り捨て）

**本実装の範囲外**:
- メール送信履歴のデータベース保存機能
- メール送信ログの検索・閲覧機能
- メール送信ログの削除・アーカイブ機能
- メール送信ログの分析・集計機能

## 2. 背景・現状分析

### 2.1 現在の実装
- **メール送信機能**: 0033-mailsenderで実装済み
  - `server/internal/service/email/email_service.go`: メール送信サービス
  - `server/internal/service/email/mock_sender.go`: 標準出力送信実装
  - `server/internal/service/email/mailpit_sender.go`: Mailpit送信実装
  - `server/internal/service/email/ses_sender.go`: AWS SES送信実装
- **ログ出力機能**: 既存のログ機能が実装済み
  - `server/internal/logging/access_logger.go`: アクセスログ出力機能
  - `server/internal/db/logger.go`: SQLログ出力機能
  - `lumberjack`と`logrus`ライブラリを使用
  - 日付別ファイル分割機能あり
- **設定**: `server/internal/config/config.go`に`LoggingConfig`構造体が存在
  - `OutputDir`フィールドが存在（アクセスログ用）
  - `SQLLogOutputDir`フィールドが存在（SQLログ用）

### 2.2 課題点
1. **メール送信内容の可視性不足**: メール送信時に送信内容が記録されないため、デバッグが困難
2. **送信内容の確認困難**: 開発環境やステージング環境でメール送信内容を確認する手段がない
3. **環境別制御の不足**: 開発環境と本番環境でメール送信ログの出力を切り替える機能がない
4. **ログ形式の不統一**: メール送信ログが存在しないため、他のログ（アクセスログ、SQLログ）と形式が統一されていない

### 2.3 本実装による改善点
1. **メール送信内容の可視化**: すべてのメール送信内容をログに記録し、デバッグを容易にする
2. **送信内容の確認**: 開発環境やステージング環境でメール送信内容をログファイルから確認できる
3. **環境別制御**: 開発環境とステージング環境でのみメール送信ログを出力し、本番環境では出力しない
4. **既存機能との統合**: 既存のアクセスログ機能やSQLログ機能と統合し、同じ`logs`ディレクトリを使用
5. **ログ形式の統一**: JSON形式でログを出力し、他のログと形式を統一

## 3. 機能要件

### 3.1 メール送信ログ出力機能

#### 3.1.1 基本機能
- メール送信時にログを出力する機能を実装
- すべてのメール送信実装（MockSender、MailpitSender、SESSender）でログを出力
- メール送信の成功・失敗に関わらずログを出力

#### 3.1.2 ログ出力内容
以下の情報をログに記録する：
- **送信時刻**: メール送信日時（タイムスタンプ）
- **メールタイトル**: メールの件名（subject）
- **メール内容**: メール本文（body）
  - 200文字以上の場合、200文字で切り捨て
  - 切り捨てた場合は末尾に`...`を付与
- **送信先**: メール送信先アドレス（to）
- **送信方式**: メール送信実装の種類（mock/mailpit/ses）
- **送信結果**: 送信成功・失敗の状態

#### 3.1.3 ログフォーマット
JSON形式で以下の情報を1行に記録する：
```json
{
  "timestamp": "2025-01-27 14:30:45",
  "to": ["recipient@example.com"],
  "subject": "Welcome to our service",
  "body": "This is a welcome email...",
  "body_truncated": false,
  "sender_type": "ses",
  "success": true,
  "error": null
}
```

例（200文字以上の場合）:
```json
{
  "timestamp": "2025-01-27 14:30:45",
  "to": ["recipient@example.com"],
  "subject": "Welcome to our service",
  "body": "This is a very long email body that contains more than 200 characters. When the body is longer than 200 characters, it will be truncated and the truncated flag will be set to true. This helps to keep log files manageable and prevents them from growing too large...",
  "body_truncated": true,
  "sender_type": "ses",
  "success": true,
  "error": null
}
```

#### 3.1.4 ログファイル名
- 形式: `mail-YYYY-MM-DD.log`
- 例: `mail-2025-01-27.log`
- 日付が変わったら自動的に新しいファイルに切り替える

### 3.2 環境別制御

#### 3.2.1 環境判定
- `APP_ENV`環境変数または設定ファイルから環境を判定
- 環境別の動作:
  - **開発環境（develop）**: メール送信ログを出力
  - **ステージング環境（staging）**: メール送信ログを出力
  - **本番環境（production）**: メール送信ログを出力しない

#### 3.2.2 実装方法
- `config.Load()`で取得した設定の`APP_ENV`を確認
- production環境の場合はログを出力しない、またはLoggerを無効化する設定を適用

### 3.3 ログ出力先の設定機能

#### 3.3.1 デフォルト設定
- デフォルトのログ出力先: `logs`ディレクトリ（プロジェクトルート）
- ディレクトリが存在しない場合は自動作成

#### 3.3.2 設定による変更
- 設定ファイル（`config/{env}/config.yaml`）の`logging`セクションに`mail_log_output_dir`項目を追加
- 設定ファイルでログ出力先を変更可能
- **絶対パスと相対パスの両方をサポート**:
  - 相対パス: プロジェクトルートからの相対パス（例: `logs`）
  - 絶対パス: システムの絶対パス（例: `/var/log/go-webdb-template`）
- 設定が空の場合は、既存の`output_dir`設定を使用（アクセスログと同じディレクトリ）

#### 3.3.3 設定例
```yaml
logging:
  level: debug
  format: json
  output: file
  output_dir: logs  # アクセスログ用（デフォルト値）
  mail_log_output_dir: logs  # メール送信ログ用（オプション、空の場合はoutput_dirを使用）
  # mail_log_output_dir: /var/log/go-webdb-template  # カスタムパス
```

### 3.4 メール送信実装への統合

#### 3.4.1 EmailServiceへの統合
- `EmailService`にメール送信ログ機能を統合
- メール送信前にログを出力（送信前の情報を記録）
- メール送信後にログを出力（送信結果を記録）

#### 3.4.2 送信実装への統合
- `MockSender`、`MailpitSender`、`SESSender`の各実装でログ出力
- 各送信実装の`Send`メソッド内でログを出力
- 送信エラーが発生した場合もログに記録

## 4. 非機能要件

### 4.1 パフォーマンス
- メール送信ログの出力はメール送信処理のパフォーマンスに影響を与えないこと
- ログ出力がメール送信処理をブロックしないこと（同期的に出力）

### 4.2 セキュリティ
- メール本文に機密情報が含まれる可能性があるため、200文字制限により機密情報の漏洩リスクを軽減
- 本番環境ではログを出力しないことで、機密情報の漏洩リスクを回避

### 4.3 可用性
- ログ出力に失敗してもメール送信処理は継続すること
- ログ出力のエラーはメール送信のエラーとして扱わないこと

### 4.4 保守性
- 既存のログ機能（アクセスログ、SQLログ）と同じパターンで実装
- `lumberjack`と`logrus`ライブラリを使用
- 日付別ファイル分割機能を活用

## 5. 技術仕様

### 5.1 サーバー側技術スタック
- **言語**: Go 1.21+
- **ログライブラリ**: `github.com/sirupsen/logrus`
- **ログローテーション**: `gopkg.in/natefinch/lumberjack.v2`
- **JSONエンコーディング**: Go標準ライブラリの`encoding/json`

### 5.2 既存実装との統合
- `server/internal/logging/access_logger.go`の実装パターンを参考
- `server/internal/db/logger.go`の実装パターンを参考
- `server/internal/config/config.go`の設定構造体に項目を追加

### 5.3 ファイル構造
- **ログ出力先**: `server/logs/mail-YYYY-MM-DD.log`（デフォルト）
- **設定ファイル**: `server/config/{env}/config.yaml`

## 6. 受け入れ基準

### 6.1 機能要件
1. **メール送信ログ出力**: メール送信時にログが出力されること
2. **日付別ファイル分割**: 日付が変わったら自動的に新しいログファイルに切り替わること
3. **JSON形式**: ログがJSON形式で出力されること
4. **メール内容の切り捨て**: メール本文が200文字以上の場合、200文字で切り捨てられること
5. **環境別制御**: 開発環境とステージング環境でのみログが出力され、本番環境では出力されないこと
6. **設定による出力先変更**: 設定ファイルでログ出力先を変更できること

### 6.2 非機能要件
1. **パフォーマンス**: ログ出力がメール送信処理のパフォーマンスに影響を与えないこと
2. **エラーハンドリング**: ログ出力に失敗してもメール送信処理が継続すること
3. **既存機能との統合**: 既存のログ機能と統合され、同じ`logs`ディレクトリを使用すること

## 7. 制約事項

1. **メール内容の切り捨て**: メール本文は200文字で切り捨てられる（機密情報の漏洩リスク軽減のため）
2. **本番環境での無効化**: 本番環境ではメール送信ログを出力しない（セキュリティ上の理由）
3. **ログ出力の同期処理**: ログ出力は同期的に処理される（非同期処理は将来的な拡張項目）

## 8. 将来の拡張項目（現時点では未実装）

以下の機能は将来の拡張として検討されていますが、現時点では実装対象外です：

- メール送信履歴のデータベース保存機能
- メール送信ログの検索・閲覧機能
- メール送信ログの削除・アーカイブ機能
- メール送信ログの分析・集計機能
- ログ出力の非同期処理
- メール内容の完全な記録（切り捨てなし）

## Project Description (Input)

メール送信ログを追加する。

### 要件
* デバッグ用のメール送信ログを追加する。
* ログはSQLログ、アクセスログと同じように、logsディレクトリに出力する。
    * mail-2025-12-25.log のようなファイル名を作成する。
* ログの出力内容は、送信時刻、メールタイトル、メール内容(200文字以上切り捨て)が入っていると良い。
* ログのフォーマットはJSON。
* メール送信ログは、develop、staging環境で出力する。

## Requirements

### Requirement 1: メール送信ログ出力機能の実装
**Objective:** As a developer, I want to log email sending activities, so that I can debug and verify email sending in development and staging environments.

#### Acceptance Criteria
1. WHEN an email is sent THEN the system SHALL log the email sending activity
2. IF email is sent successfully THEN the system SHALL log success status
3. IF email sending fails THEN the system SHALL log error status
4. WHERE email is sent THEN the system SHALL log timestamp, recipient, subject, and body
5. WHEN email body is longer than 200 characters THEN the system SHALL truncate it to 200 characters
6. IF email body is truncated THEN the system SHALL set `body_truncated` flag to true
7. WHERE email body is truncated THEN the system SHALL append `...` to the truncated body

### Requirement 2: 日付別ログファイル分割機能の実装
**Objective:** As a system, I want to split email logs by date, so that log files are manageable and easy to review.

#### Acceptance Criteria
1. WHEN email log is written THEN the system SHALL write to a file named `mail-YYYY-MM-DD.log`
2. IF date changes THEN the system SHALL automatically switch to a new log file
3. WHERE log file is created THEN the system SHALL use the current date in the filename
4. WHEN log file is created THEN the system SHALL use lumberjack library for date-based rotation

### Requirement 3: JSON形式でのログ出力
**Objective:** As a developer, I want email logs in JSON format, so that logs are easy to parse and analyze.

#### Acceptance Criteria
1. WHEN email log is written THEN the system SHALL output in JSON format
2. IF JSON log is written THEN it SHALL contain timestamp, to, subject, body, body_truncated, sender_type, success, and error fields
3. WHERE JSON log is written THEN it SHALL be written as a single line per log entry
4. WHEN JSON log is written THEN the system SHALL use Go standard library's encoding/json package

### Requirement 4: 環境別制御機能の実装
**Objective:** As a system, I want to control email log output based on environment, so that logs are only output in development and staging environments.

#### Acceptance Criteria
1. WHEN APP_ENV is "develop" THEN the system SHALL output email logs
2. WHEN APP_ENV is "staging" THEN the system SHALL output email logs
3. WHEN APP_ENV is "production" THEN the system SHALL NOT output email logs
4. IF environment is not specified THEN the system SHALL output email logs (default to develop behavior)
5. WHERE email log is disabled THEN the system SHALL not create log files

### Requirement 5: ログ出力先の設定機能
**Objective:** As a system administrator, I want to configure email log output directory, so that logs can be stored in a custom location.

#### Acceptance Criteria
1. WHEN mail_log_output_dir is configured THEN the system SHALL use the configured directory
2. IF mail_log_output_dir is empty THEN the system SHALL use output_dir from logging config
3. WHERE output directory is specified THEN the system SHALL support both absolute and relative paths
4. WHEN output directory does not exist THEN the system SHALL create it automatically
5. IF directory creation fails THEN the system SHALL return an appropriate error

### Requirement 6: メール送信実装への統合
**Objective:** As a system, I want email logging integrated into email sending implementations, so that all email sending activities are logged.

#### Acceptance Criteria
1. WHEN MockSender sends email THEN the system SHALL log the email sending activity
2. WHEN MailpitSender sends email THEN the system SHALL log the email sending activity
3. WHEN SESSender sends email THEN the system SHALL log the email sending activity
4. IF email sending fails THEN the system SHALL log the error
5. WHERE email is sent THEN the system SHALL log sender_type (mock/mailpit/ses)

### Requirement 7: エラーハンドリング
**Objective:** As a system, I want email logging to handle errors gracefully, so that log failures do not affect email sending.

#### Acceptance Criteria
1. WHEN log output fails THEN the system SHALL NOT fail email sending
2. IF log output fails THEN the system SHALL log the error to standard error output
3. WHERE log output fails THEN the system SHALL continue email sending process
4. WHEN log file cannot be created THEN the system SHALL return an appropriate error but continue email sending
