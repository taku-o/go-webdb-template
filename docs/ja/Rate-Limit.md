**[日本語]** | [English](../en/Rate-Limit.md)

# APIレートリミット機能利用手順

## 概要

このドキュメントでは、go-webdb-templateのAPIレートリミット機能の利用手順を説明します。

レートリミット機能は、APIエンドポイントへのリクエストをIPアドレス単位で制限し、過剰なリクエストを防ぎます。

## 機能説明

### レートリミットの種類

本プロジェクトでは、以下の2種類のレートリミットをサポートしています：

1. **分制限（常に有効）**: 1分あたりのリクエスト数を制限
2. **時間制限（オプション）**: 1時間あたりのリクエスト数を制限（設定されている場合のみ）

### ストレージタイプ

レートリミットのカウンターは、環境に応じて異なるストレージを使用します：

- **In-Memory**: 開発環境（デフォルト）
- **Redis Cluster**: ステージング/本番環境（推奨）

## 設定方法

### 設定ファイル

レートリミット設定は`config/{env}/config.yaml`の`api.rate_limit`セクションで行います。

```yaml
api:
  rate_limit:
    enabled: true                    # レートリミットの有効/無効
    requests_per_minute: 60          # 1分あたりのリクエスト数
    requests_per_hour: 1000          # 1時間あたりのリクエスト数（オプション、0の場合は無効）
    storage_type: "auto"             # ストレージタイプ（"auto", "redis", "memory"）
```

### ストレージタイプの選択

- `auto`: 環境に応じて自動判定
  - Redis Cluster設定（`cacheserver.yaml`）が存在し、`addrs`が設定されている場合: Redis Cluster
  - それ以外: In-Memory
- `redis`: 強制的にRedis Clusterを使用（設定が存在しない場合はエラー）
- `memory`: 強制的にIn-Memoryを使用

### Redis Cluster設定

Redis Clusterを使用する場合、`config/{env}/cacheserver.yaml`で設定します。

```yaml
cache_server:
  redis:
    default:
      cluster:
        addrs:
          - "localhost:6379"
          - "localhost:6380"
          - "localhost:6381"
        max_retries: 2
        min_retry_backoff: 8ms
        max_retry_backoff: 512ms
        dial_timeout: 5s
        read_timeout: 3s
        pool_size: 100
        pool_timeout: 4s
```

**注意**: `addrs`が空または未設定の場合は、In-Memoryストレージが使用されます。

## レスポンスヘッダー

すべてのAPIレスポンスに、以下のレートリミット情報が含まれます。

### 分制限（常に付与）

| ヘッダー | 説明 | 例 |
|---------|------|-----|
| `X-RateLimit-Limit` | 1分あたりの制限値 | `60` |
| `X-RateLimit-Remaining` | 残りリクエスト数 | `45` |
| `X-RateLimit-Reset` | リセット時刻（Unix timestamp） | `1706342400` |

### 時間制限（`requests_per_hour`が設定されている場合のみ）

| ヘッダー | 説明 | 例 |
|---------|------|-----|
| `X-RateLimit-Hour-Limit` | 1時間あたりの制限値 | `1000` |
| `X-RateLimit-Hour-Remaining` | 残りリクエスト数 | `950` |
| `X-RateLimit-Hour-Reset` | リセット時刻（Unix timestamp） | `1706346000` |

## レートリミット超過時

制限を超過した場合、HTTP 429ステータスコードが返されます。

### レスポンス例

```json
{
  "code": 429,
  "message": "Too Many Requests"
}
```

### レスポンスヘッダー

レートリミット超過時も、レスポンスヘッダーにレートリミット情報が含まれます：

```
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1706342460
```

## 動作確認

### レートリミットヘッダーの確認

```bash
curl -i -H "Authorization: Bearer <YOUR_API_KEY>" http://localhost:8080/api/users
```

**レスポンス例**:

```
HTTP/1.1 200 OK
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 59
X-RateLimit-Reset: 1706342460
X-RateLimit-Hour-Limit: 1000
X-RateLimit-Hour-Remaining: 999
X-RateLimit-Hour-Reset: 1706346000
```

### レートリミット超過の確認

```bash
# 60回連続でリクエストを送信
for i in {1..61}; do
  curl -H "Authorization: Bearer <YOUR_API_KEY>" http://localhost:8080/api/users
  echo "Request $i"
done
```

61回目のリクエストでHTTP 429が返されます。

## Fail-Open方式

レートリミット機能は**fail-open方式**を採用しています。これは、レートリミットの初期化やチェック時にエラーが発生した場合、リクエストを許可することを意味します。

### エラー時の動作

- **ストレージ初期化エラー**: ログに記録し、すべてのリクエストを許可
- **レートリミットチェックエラー**: ログに記録し、リクエストを許可

この方式により、レートリミット機能の障害がアプリケーション全体の障害につながることを防ぎます。

## IPアドレスの取得方法

レートリミットはIPアドレス単位で適用されます。IPアドレスは以下の優先順位で取得されます：

1. `X-Forwarded-For`ヘッダー（プロキシ経由の場合）
2. `X-Real-IP`ヘッダー（プロキシ経由の場合）
3. リクエストの`RemoteAddr`

**注意**: IPアドレスが取得できない場合は、リクエストは許可されます。

## 環境別の推奨設定

### 開発環境

```yaml
api:
  rate_limit:
    enabled: true
    requests_per_minute: 60
    requests_per_hour: 0              # 時間制限は無効
    storage_type: "auto"              # In-Memoryが使用される
```

### ステージング環境

```yaml
api:
  rate_limit:
    enabled: true
    requests_per_minute: 60
    requests_per_hour: 1000
    storage_type: "auto"              # Redis Clusterが使用される
```

### 本番環境

```yaml
api:
  rate_limit:
    enabled: true
    requests_per_minute: 60
    requests_per_hour: 1000
    storage_type: "redis"             # 強制的にRedis Clusterを使用
```

## 注意事項

1. **In-Memoryストレージの制限**: In-Memoryストレージは、サーバー再起動時にカウンターがリセットされます。また、複数のサーバーインスタンス間でカウンターが共有されません。
2. **Redis Clusterの可用性**: Redis Clusterが利用できない場合、fail-open方式によりすべてのリクエストが許可されます。本番環境では、Redis Clusterの可用性を監視することを推奨します。
3. **IPアドレスの信頼性**: プロキシ経由でアクセスする場合、`X-Forwarded-For`ヘッダーが正しく設定されていることを確認してください。
4. **レートリミットの調整**: アプリケーションの使用パターンに応じて、`requests_per_minute`と`requests_per_hour`を適切に調整してください。

## トラブルシューティング

### レートリミットが機能しない

1. `enabled: true`が設定されているか確認
2. ログにエラーメッセージがないか確認
3. Redis Clusterを使用する場合、`cacheserver.yaml`の設定を確認

### レートリミットが厳しすぎる

1. `requests_per_minute`と`requests_per_hour`の値を増やす
2. クライアント側でリクエストをバッチ処理する

### Redis Cluster接続エラー

1. Redis Clusterが起動しているか確認
2. `cacheserver.yaml`の`addrs`設定を確認
3. ネットワーク接続を確認

## 関連ドキュメント

- [README.md](../README.md) - レートリミット機能の概要
- [Redis-Reconnection.md](Redis-Reconnection.md) - Redis接続設定（作成予定）
