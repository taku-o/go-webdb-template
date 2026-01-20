# 開発環境サーバー構成

開発環境では4つのサーバーを起動する必要があります。

## サーバー一覧

| サーバー | ポート | ディレクトリ | 起動コマンド | ヘルスチェック |
|---------|-------|-------------|-------------|---------------|
| API サーバー | 8080 | `server/cmd/server` | `APP_ENV=develop go run ./cmd/server/main.go` | `/health` |
| クライアント | 3000 | `client/` | `npm run dev` | `/health` |
| Admin | 8081 | `server/cmd/admin` | `APP_ENV=develop go run ./cmd/admin/main.go` | `/health` |
| JobQueue | 8082 | `server/cmd/jobqueue` | `APP_ENV=develop go run ./cmd/jobqueue/main.go` | `/health` |

## Docker版サーバー

| サーバー | docker-compose | コンテナ名 | ポート |
|---------|---------------|-----------|-------|
| API サーバー | `docker-compose.api.yml` | api | 8080 |
| クライアント | `docker-compose.client.yml` | client | 3000 |
| Admin | `docker-compose.admin.yml` | admin | 8081 |
| JobQueue | `docker-compose.jobqueue.yml` | jobqueue | 8082 |

## 注意事項

- 「サーバーを起動して」と言われた場合、上記4つ全てを起動すること
- クライアントはNext.jsアプリケーション（port 3000）
- API、Admin、JobQueueサーバーはGoアプリケーション
- 全サーバーに `/health` エンドポイントが存在する
- Docker版を起動する前にPostgreSQLとRedisコンテナを起動する必要がある
