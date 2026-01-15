**[日本語]** | [English](../en/Queue-Job.md)

# ジョブキュー機能利用手順

## 概要

このドキュメントでは、go-webdb-templateのジョブキュー機能の利用手順を説明します。

ジョブキュー機能は、Redis と Asynq ライブラリを使用した遅延ジョブ処理システムです。

## 環境構築

### 1. Redisの起動

```bash
./scripts/start-redis.sh start
```

### 2. Redis Insightの起動（オプション）

Redis の状態を確認したい場合：

```bash
./scripts/start-redis-insight.sh start
```

Web UI: http://localhost:8001

### 3. APIサーバーの起動

```bash
APP_ENV=develop go run ./cmd/server/main.go
```

## API経由でのジョブ登録

### エンドポイント

**POST** `/api/dm-jobqueue/register`

### 認証

このエンドポイントは認証が必要です。

### リクエスト例

```bash
curl -X POST http://localhost:8080/api/dm-jobqueue/register \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -d '{
    "message": "Hello, World!",
    "delay_seconds": 10,
    "max_retry": 3
  }'
```

### リクエストボディ

| フィールド | 型 | 必須 | 説明 |
|-----------|------|------|------|
| message | string | いいえ | 出力するメッセージ（デフォルト: "Job executed successfully"） |
| delay_seconds | int | いいえ | 遅延時間（秒、デフォルト: 180秒） |
| max_retry | int | いいえ | 最大リトライ回数（デフォルト: 10回） |

### レスポンス例

```json
{
  "job_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "registered"
}
```

## ジョブの実行確認

登録されたジョブは、指定した遅延時間後にAPIサーバーの標準出力に出力されます。

```
[2026-01-02 12:34:56] Hello, World!
```

## 停止手順

### APIサーバーの停止

Ctrl+C または `kill <PID>` で停止します。

### Redisの停止

```bash
./scripts/start-redis.sh stop
```

### Redis Insightの停止

```bash
./scripts/start-redis-insight.sh stop
```

## 注意事項

### 迷子のAsynqサーバー設定問題

各APIサーバーはAsynqサーバーを内蔵しており、Redisに自身を登録します（`asynq:servers:{hostname:pid:uuid}`）。
正しくAPIサーバーが終了しない時、Redisにサーバー設定が残ってしまう事がある。

- その際、ただしくジョブを処理できないケースが観測された。

**確認方法**:

```bash
docker exec redis redis-cli KEYS "asynq:servers:*"
```

正常な状態では、起動中のサーバー1つにつき1エントリのみです。

**対処方法**:

古いプロセスが残っている場合は、該当プロセスを終了してください：

```bash
# PIDを確認
ps aux | grep "go run ./cmd/server/main.go"

# プロセスを終了
kill <PID>
```

### Redis接続エラー時の動作

Redisに接続できない場合でも、APIサーバーは起動します。ただし、ジョブ登録APIは503エラーを返します。

## 関連ドキュメント

- [API Documentation](./API.md): APIエンドポイントの詳細
- [Architecture](./Architecture.md): システムアーキテクチャの説明
