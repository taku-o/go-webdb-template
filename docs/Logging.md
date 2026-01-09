# ログ機能利用手順

## 概要

このドキュメントでは、go-webdb-templateのログ機能の利用手順を説明します。

本プロジェクトでは、以下の3種類のログをサポートしています：

1. **アクセスログ**: HTTPリクエスト/レスポンスのログ
2. **メール送信ログ**: メール送信処理のログ
3. **SQLログ**: データベースクエリのログ

## 設定方法

### 設定ファイル

ログ設定は`config/{env}/config.yaml`の`logging`セクションで行います。

```yaml
logging:
  level: debug                    # ログレベル（debug, info, warn, error）
  format: json                    # ログフォーマット（json, text）
  output: stdout                  # 出力先（stdout, file）
  output_dir: ../logs             # ログファイル出力先ディレクトリ
  sql_log_enabled: true           # SQLログの有効/無効
  mail_log_enabled: true          # メール送信ログの有効/無効
  mail_log_output_dir: ../logs    # メール送信ログ出力先（オプション、未設定時はoutput_dirを使用）
```

### ログレベル

- `debug`: デバッグ情報を含むすべてのログ
- `info`: 情報レベルのログ
- `warn`: 警告レベルのログ
- `error`: エラーレベルのログ

### ログフォーマット

- `json`: JSON形式で出力（構造化ログ）
- `text`: テキスト形式で出力（人間が読みやすい形式）

### 出力先

- `stdout`: 標準出力に出力
- `file`: ファイルに出力（`output_dir`で指定されたディレクトリ）

## アクセスログ

### 概要

アクセスログは、すべてのHTTPリクエストとレスポンスの情報を記録します。

### 出力形式（JSON）

```json
{
  "timestamp": "2025-01-27 10:30:45",
  "method": "POST",
  "path": "/api/users",
  "protocol": "HTTP/1.1",
  "status_code": 201,
  "response_time_ms": 45.23,
  "remote_ip": "192.168.1.100",
  "user_agent": "Mozilla/5.0...",
  "headers": "Content-Type: application/json; Authorization: Bearer ...",
  "request_body": "{\"name\":\"John\",\"email\":\"john@example.com\"}"
}
```

**フィールド説明**:
- `timestamp`: リクエスト受信時刻
- `method`: HTTPメソッド（GET, POST, PUT, DELETE等）
- `path`: リクエストパス
- `protocol`: HTTPプロトコルバージョン
- `status_code`: HTTPステータスコード
- `response_time_ms`: レスポンス時間（ミリ秒）
- `remote_ip`: クライアントのIPアドレス
- `user_agent`: ユーザーエージェント
- `headers`: リクエストヘッダー（セミコロン区切り）
- `request_body`: リクエストボディ（POST/PUT/PATCHの場合のみ、最大1MB）

### リクエストボディのログ出力条件

以下の条件を満たす場合のみ、リクエストボディがログに出力されます：

1. HTTPメソッドがPOST、PUT、またはPATCH
2. Content-Lengthが1MB以下
3. Content-Typeが画像/動画/音声/マルチパートではない

### ログファイルの場所

- 出力先: `{output_dir}/access.log`
- ローテーション: 日付別ファイル分割（`access-2025-01-27.log`形式）

## メール送信ログ

### 概要

メール送信ログは、メール送信処理の成功/失敗を記録します。

### 設定

メール送信ログを有効にするには、`mail_log_enabled: true`を設定します。

```yaml
logging:
  mail_log_enabled: true
  mail_log_output_dir: ../logs  # オプション、未設定時はoutput_dirを使用
```

### 出力形式（JSON）

```json
{
  "timestamp": "2025-01-27 10:30:45",
  "to": ["recipient@example.com"],
  "subject": "Welcome",
  "body": "Hello, John!",
  "body_truncated": false,
  "sender_type": "mock",
  "success": true,
  "error": ""
}
```

**エラー時**:

```json
{
  "timestamp": "2025-01-27 10:30:45",
  "to": ["recipient@example.com"],
  "subject": "Welcome",
  "body": "Hello, John!",
  "body_truncated": false,
  "sender_type": "ses",
  "success": false,
  "error": "failed to send email: connection timeout"
}
```

**フィールド説明**:
- `timestamp`: メール送信時刻
- `to`: 送信先メールアドレスの配列
- `subject`: メール件名
- `body`: メール本文
- `body_truncated`: メール本文が切り詰められたかどうか（長すぎる場合）
- `sender_type`: 送信方式（`mock`, `mailpit`, `ses`）
- `success`: 送信成功/失敗
- `error`: エラーメッセージ（失敗時のみ）

### ログファイルの場所

- 出力先: `{mail_log_output_dir}/mail.log`（未設定時は`{output_dir}/mail.log`）
- ローテーション: 日付別ファイル分割（`mail-2025-01-27.log`形式）

## SQLログ

### 概要

SQLログは、データベースクエリの実行内容を記録します。GORMを使用したクエリがログに出力されます。

### 設定

SQLログを有効にするには、`sql_log_enabled: true`を設定します。

```yaml
logging:
  sql_log_enabled: true
  output_dir: ../logs
```

### 出力形式（テキスト）

```
[2025-01-27 10:30:45] [INFO] [master] [sqlite3] SELECT * FROM `dm_news` WHERE `dm_news`.`id` = 'abc123' LIMIT 1
[2025-01-27 10:30:45] [INFO] [sharding] [sqlite3] INSERT INTO `dm_users_000` (`id`,`name`,`email`,`created_at`,`updated_at`) VALUES ('def456','John','john@example.com','2025-01-27 10:30:45','2025-01-27 10:30:45')
```

**フォーマット**: `[時刻] [レベル] [グループ名] [ドライバー] SQLクエリ`

### ログファイルの場所

- 出力先: `{output_dir}/sql-{date}.log`（例: `sql-2025-01-27.log`）
- ローテーション: 日付別ファイル分割（1日1ファイル）

### ログレベル

SQLログは`Info`レベルで出力されます。GORMのログレベル設定に従います。

## ログローテーション

### アクセスログ・メール送信ログ

`lumberjack`ライブラリを使用してログローテーションを実装しています。

- **日付別分割**: 1日1ファイル（`access-2025-01-27.log`形式）
- **サイズ制限**: なし（日付別分割のみ使用）
- **保持期間**: 無制限（手動削除が必要）

### SQLログ

- **日付別分割**: 1日1ファイル（`sql-2025-01-27.log`形式）
- **保持期間**: 無制限（手動削除が必要）

## ログファイルの確認方法

### 開発環境

```bash
# アクセスログの確認
tail -f logs/access-$(date +%Y-%m-%d).log

# メール送信ログの確認
tail -f logs/mail-$(date +%Y-%m-%d).log

# SQLログの確認
tail -f logs/sql-$(date +%Y-%m-%d).log
```

### JSON形式のログを整形して確認

```bash
# jqを使用してJSONログを整形
tail -f logs/access-$(date +%Y-%m-%d).log | jq .

# 特定の条件でフィルタリング
tail -f logs/access-$(date +%Y-%m-%d).log | jq 'select(.status_code >= 400)'
```

## パフォーマンスへの影響

### アクセスログ

- リクエストボディの読み取りは、メモリ上で行われます
- 大きなリクエストボディ（1MB超）はログに出力されません
- ファイルI/Oは非同期で行われますが、高負荷時はパフォーマンスに影響する可能性があります

### メール送信ログ

- メール送信処理の一部として実行されるため、送信処理のオーバーヘッドは最小限です
- メール本文が長すぎる場合は切り詰められます

### SQLログ

- すべてのSQLクエリがログに出力されるため、高負荷時はログファイルが大きくなる可能性があります
- 本番環境では`sql_log_enabled: false`に設定することを推奨します

## 本番環境での推奨設定

```yaml
logging:
  level: info                    # デバッグ情報は出力しない
  format: json                   # 構造化ログでログ分析ツールと連携しやすい
  output: file                   # ファイルに出力
  output_dir: /var/log/app       # 適切なログディレクトリを指定
  sql_log_enabled: false         # SQLログは無効化（パフォーマンスとセキュリティのため）
  mail_log_enabled: true         # メール送信ログは有効化（監査のため）
```

## 注意事項

1. **ログファイルのディスク容量**: ログファイルは自動削除されないため、定期的なクリーンアップが必要です
2. **セキュリティ**: ログには認証トークンやリクエストボディが含まれる可能性があるため、適切なアクセス制御を設定してください
3. **パフォーマンス**: 高負荷環境では、ログ出力がパフォーマンスに影響する可能性があります
4. **SQLログ**: 本番環境ではSQLログを無効化することを推奨します（セキュリティとパフォーマンスのため）
