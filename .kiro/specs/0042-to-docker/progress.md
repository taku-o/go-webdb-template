# Docker化実装進捗

## 現在の状況

**作業中タスク**: Phase 7 - 動作確認（開発環境）
**最終更新**: 2026-01-04

### Phase 7 動作確認結果

| タスク | 状態 | 結果 |
|-------|------|------|
| 7.1 APIサーバーDockerイメージのビルド | 確認済み | ビルド成功 |
| 7.2 APIサーバーコンテナの起動 | 確認済み | 起動成功、`/health` → OK |
| 7.3 AdminサーバーDockerイメージのビルド | 確認済み | ビルド成功 |
| 7.4 Adminサーバーコンテナの起動 | 確認済み | 起動成功、`/admin` → 405（エンドポイント存在） |
| 7.5 クライアントサーバーDockerイメージのビルド | 確認済み | ビルド成功 |
| 7.6 クライアントサーバーコンテナの起動 | 確認済み | 起動成功、`/` → 200 OK |
| 7.7 3つのサーバーの同時起動 | 確認済み | 3サーバー同時稼働確認 |

**注意**: Admin/APIサーバーのDockerヘルスチェックは「unhealthy」表示。これはPhase 8で対応予定（Adminの`/health`エンドポイント未実装）。

## 作成済みファイル

| ファイル | 状態 | 備考 |
|---------|------|------|
| `server/Dockerfile` | 作成済み | staging/production用（CGO_ENABLED=0） |
| `server/Dockerfile.develop` | 作成済み | develop用（CGO_ENABLED=1、SQLite対応） |
| `server/.dockerignore` | 作成済み | |
| `server/Dockerfile.admin` | 作成済み | staging/production用（CGO_ENABLED=0） |
| `server/Dockerfile.admin.develop` | 作成済み | develop用（CGO_ENABLED=1、SQLite対応） |
| `client/Dockerfile` | 作成済み | dev/productionターゲット対応 |
| `client/.dockerignore` | 作成済み | |
| `docker-compose.api.develop.yml` | 作成済み | Dockerfile.develop使用 |
| `docker-compose.client.develop.yml` | 作成済み | devターゲット使用 |
| `docker-compose.admin.develop.yml` | 作成済み | Dockerfile.admin.develop使用 |
| `docker-compose.api.staging.yml` | 作成済み | Dockerfile使用 |
| `docker-compose.client.staging.yml` | 作成済み | productionターゲット使用 |
| `docker-compose.admin.staging.yml` | 作成済み | Dockerfile.admin使用 |
| `docker-compose.api.production.yml` | 作成済み | Dockerfile使用 |
| `docker-compose.client.production.yml` | 作成済み | productionターゲット使用 |
| `docker-compose.admin.production.yml` | 作成済み | Dockerfile.admin使用 |

## 設計書からの変更点

### 1. Goバージョンとベースイメージの変更

- **設計書**: `golang:1.21-alpine`, `alpine:latest`
- **実装**: `golang:1.24-bookworm`, `debian:bookworm-slim`
- **理由**:
  - go.modで`go 1.24.0`が指定されているため
  - Alpine Linux（musl libc）ではgo-sqlite3のビルドに失敗するため（`pread64`/`pwrite64`等のglibc固有関数が存在しない）
  - 環境統一のため、staging/productionも同様にDebian系に変更

### 1.1. .dockerignoreの修正

- **設計書**: `*.sum`を除外
- **実装**: `*.sum`を除外対象から削除、`bin/`、`admin`、`server/`を追加
- **理由**:
  - `go.sum`はGoの依存関係管理に必要なため（除外対象から削除）
  - ローカルでビルドされたバイナリがDockerコンテキストにコピーされるのを防ぐため（追加）

### 1.2. WORKDIRの変更

- **設計書**: `/app`
- **実装**: `/app/server`
- **理由**: アプリケーションが`server/`ディレクトリから起動される想定で、`../config/`や`../logs/`を参照するため

### 1.3. Redis接続設定の環境変数オーバーライド

- **設計書**: 記載なし
- **実装**: `viper.BindEnv`による環境変数オーバーライドを追加
- **変更ファイル**:
  - `server/internal/config/config.go` - `viper.BindEnv("cache_server.redis.jobqueue.addr", "REDIS_JOBQUEUE_ADDR")`を追加
  - 全docker-composeファイル - `REDIS_JOBQUEUE_ADDR=redis:6379`環境変数を追加
- **理由**: Docker環境では`localhost:6379`ではなく`redis:6379`でRedisに接続する必要があるため

### 1.4. クライアントpackage-lock.jsonの修正

- **問題**: `yaml@2.8.2`がpackage-lock.jsonに不足
- **解決**: `npm install yaml@2.8.2`で追加

### 1.5. クライアント環境変数の読み込み修正

- **問題**: `environment`セクションで`NEXT_PUBLIC_API_KEY=${NEXT_PUBLIC_API_KEY:-}`を指定していたため、ホスト環境変数がない場合に空文字列で上書きされ、`env_file`の値が無効化されていた
- **解決**: `environment`セクションとビルド引数から`NEXT_PUBLIC_API_KEY`を削除し、`env_file`のみで読み込むよう変更

### 1.6. クライアントDockerfileのユーザー設定修正

- **問題**: `production`ターゲットで`addgroup -g 1000 nodeuser`が失敗（node:22-alpineに既にuid/gid 1000の`node`ユーザーが存在）
- **解決**: 新規ユーザー作成をやめ、既存の`node`ユーザーを使用
- **影響範囲**: staging/production環境のクライアントビルドのみ（develop環境は`dev`ターゲットを使用するため影響なし）

### 2. 開発用Dockerfileの追加

設計書では1つのDockerfileを想定していたが、以下の理由で開発用を分離：

- **問題**: `github.com/mattn/go-sqlite3`はCGOが必要
- **設計書の矛盾**: Docker開発環境でSQLiteを使用するが、CGO_ENABLED=0を指定
- **解決策**: 環境別にDockerfileを分離

**追加ファイル**:
- `server/Dockerfile.develop` - APIサーバー開発用（CGO_ENABLED=1）
- `server/Dockerfile.admin.develop` - Adminサーバー開発用（CGO_ENABLED=1）

### 3. Dockerfile構成（変更後）

| ファイル | 用途 | CGO | DB |
|---------|------|-----|-----|
| `server/Dockerfile` | staging/production | 0 | PostgreSQL/MySQL |
| `server/Dockerfile.develop` | develop | 1 | SQLite |
| `server/Dockerfile.admin` | staging/production | 0 | PostgreSQL/MySQL |
| `server/Dockerfile.admin.develop` | develop | 1 | SQLite |

## 未着手タスク

- Phase 8: 動作確認（環境別：staging/production）
- Phase 9: Dockerイメージ確認
- Phase 10: デプロイメント準備
- Phase 11: ドキュメント整備

## 備考

- docker-compose設定ファイルは開発用は`Dockerfile.develop`を参照するよう変更済み
- ボリュームマウントは`./server/data:/app/server/data`に変更済み（WORKDIRの変更に伴う）
