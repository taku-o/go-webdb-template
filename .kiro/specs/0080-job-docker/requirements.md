# JobQueueサーバーのDocker対応要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Feature名**: 0080-job-docker
- **作成日**: 2026-01-17
- **関連Issue**: https://github.com/taku-o/go-webdb-template/issues/163

### 1.2 目的
JobQueueサーバーをDocker上で動作するように対応する。既存のAPIサーバーやAdminサーバーと同様のDocker環境を構築し、開発環境での一貫性を保つ。

### 1.3 スコープ
- `server/Dockerfile.jobqueue`の新規作成（JobQueueサーバー用Dockerfile）
- `docker-compose.jobqueue.yml`の新規作成（JobQueueサーバー用docker-compose設定）
- `docker-compose.api.yml`のコンテナ名修正（`api-develop` → `api`）
- `docker-compose.admin.yml`のコンテナ名修正（`admin-develop` → `admin`）
- `docker-compose.client.yml`のコンテナ名修正（`client-develop` → `client`）
- ドキュメントの更新（`docs/ja/Docker.md`、`docs/en/Docker.md`、`README.md`、`README.ja.md`）

**本実装の範囲外**:
- JobQueueサーバーの機能変更
- その他のdocker-composeファイルの変更
- PostgreSQL/Redisコンテナの設定変更

## 2. 背景・現状分析

### 2.1 現在の状況
- **JobQueueサーバーの実装**: `server/cmd/jobqueue/main.go`に実装されている
- **起動方法**: `APP_ENV=develop go run ./cmd/jobqueue/main.go`でローカル環境で起動
- **ポート**: 8082
- **機能**:
  - HTTPサーバー（ポート8082、`/health`エンドポイント提供）
  - Asynqサーバー（Redisからジョブを取得して処理）
- **既存のDocker対応**:
  - APIサーバー: `docker-compose.api.yml`、`server/Dockerfile`
  - Adminサーバー: `docker-compose.admin.yml`、`server/Dockerfile.admin`
  - クライアント: `docker-compose.client.yml`、`client/Dockerfile`

### 2.2 課題点
1. **Docker環境の不整合**: JobQueueサーバーがDocker対応していないため、開発環境での一貫性が欠けている
2. **起動方法の違い**: APIサーバーとAdminサーバーはDockerで起動できるが、JobQueueサーバーはローカル環境でのみ起動可能
3. **環境変数の管理**: Docker環境での環境変数設定が統一されていない
4. **ドキュメントの不整合**: Docker環境に関するドキュメントにJobQueueサーバーの記載がない

### 2.3 本実装による改善点
1. **Docker環境の統一**: すべてのサーバー（API、Admin、JobQueue）をDocker環境で起動できるようになる
2. **起動方法の統一**: すべてのサーバーをdocker-composeで起動できるようになる
3. **環境変数の統一**: Docker環境での環境変数設定が統一される
4. **ドキュメントの整合性**: Docker環境に関するドキュメントが完全になる

## 3. 機能要件

### 3.1 Dockerfileの作成

#### 3.1.1 server/Dockerfile.jobqueueの作成
- **ベースイメージ**: 
  - ビルドステージ: `golang:1.24-alpine`
  - 実行ステージ: `alpine:latest`
- **CGO設定**: CGO_ENABLED=0
- **ビルドコマンド**: `CGO_ENABLED=0 go build -o jobqueue ./cmd/jobqueue/main.go`
- **パッケージインストール**: 
  - `ca-certificates`: SSL証明書用
  - `tzdata`: タイムゾーンデータ用
  - `wget`: ヘルスチェック用（既存のAPI/Adminサーバーと同様）
- **非rootユーザー**: `appuser`（UID 1000、GID 1000）
- **ディレクトリ構造**: 
  - `/app/server/`: バイナリ配置・作業ディレクトリ
  - `/app/config/`: 設定ファイル（`../config/`でアクセス）
  - `/app/logs/`: ログファイル（`../logs/`でアクセス）
  - `/app/server/data/`: データディレクトリ
- **ポート**: 8082をEXPOSE
- **CMD**: `["./jobqueue"]`

#### 3.1.2 既存Dockerfileとの整合性
- `server/Dockerfile`（APIサーバー用）と同様の構造を維持
- `server/Dockerfile.admin`（Adminサーバー用）と同様の構造を維持
- ベースイメージ、CGO設定、パッケージ、ディレクトリ構造を統一

### 3.2 docker-compose設定ファイルの作成

#### 3.2.1 docker-compose.jobqueue.ymlの作成
- **サービス名**: `jobqueue`
- **コンテナ名**: `jobqueue`
- **ビルド設定**:
  - `context: ./server`
  - `dockerfile: Dockerfile.jobqueue`
- **ポート設定**: `"8082:8082"`
- **環境変数**:
  - `APP_ENV=develop`
  - `REDIS_JOBQUEUE_ADDR=redis:6379`
  - データベースDSN（既存のAPI/Adminサーバーと同様の設定）
- **ボリュームマウント**:
  - `./config/develop:/app/config/develop:ro`: 設定ファイル（読み取り専用）
  - `./server/data:/app/server/data`: データディレクトリ（読み書き可）
  - `./logs:/app/logs`: ログファイル（読み書き可）
- **ネットワーク**:
  - `postgres-network`: PostgreSQL接続用（external: true）
  - `redis-network`: Redis接続用（external: true）
- **restart**: `unless-stopped`
- **healthcheck**:
  - `test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8082/health"]`
  - `interval: 30s`
  - `timeout: 10s`
  - `retries: 3`
  - `start_period: 40s`

#### 3.2.2 既存docker-composeファイルのコンテナ名修正
- **docker-compose.api.yml**: コンテナ名を`api-develop`から`api`に変更
- **docker-compose.admin.yml**: コンテナ名を`admin-develop`から`admin`に変更
- **docker-compose.client.yml**: コンテナ名を`client-develop`から`client`に変更
- **理由**: コンテナ名に環境名（`-develop`）を含めない統一的な命名規則に従う

#### 3.2.3 既存docker-composeファイルとの整合性
- `docker-compose.api.yml`と同様の構造を維持
- `docker-compose.admin.yml`と同様の構造を維持
- 環境変数、ボリュームマウント、ネットワーク設定を統一
- コンテナ名の命名規則を統一（環境名を含めない）

### 3.3 ドキュメントの更新

#### 3.3.1 docs/ja/Docker.mdの更新
- **Dockerfile構成表**: `server/Dockerfile.jobqueue`を追加
- **Docker Compose設定ファイル一覧**: `docker-compose.jobqueue.yml`を追加
- **ビルド・起動・停止コマンド**: JobQueueサーバー用のコマンドを追加
- **既存サービスとの統合**: JobQueueサーバーの起動順序を追加
- **サービス間通信**: JobQueueサーバーとRedisの接続を追加
- **概要**: JobQueueサーバーをDocker対応サーバーとして追加

#### 3.3.2 docs/en/Docker.mdの更新
- **Dockerfile構成表**: `server/Dockerfile.jobqueue`を追加
- **Docker Compose設定ファイル一覧**: `docker-compose.jobqueue.yml`を追加
- **ビルド・起動・停止コマンド**: JobQueueサーバー用のコマンドを追加
- **既存サービスとの統合**: JobQueueサーバーの起動順序を追加
- **サービス間通信**: JobQueueサーバーとRedisの接続を追加
- **概要**: JobQueueサーバーをDocker対応サーバーとして追加

#### 3.3.3 README.mdの更新
- **Docker環境の説明**: JobQueueサーバーのDocker起動方法を追加
- **ビルド・起動コマンド**: JobQueueサーバー用のコマンドを追加

#### 3.3.4 README.ja.mdの更新
- **Docker環境の説明**: JobQueueサーバーのDocker起動方法を追加
- **ビルド・起動コマンド**: JobQueueサーバー用のコマンドを追加

## 4. 非機能要件

### 4.1 パフォーマンス
- **ビルド時間**: 既存のAPI/Adminサーバーと同等のビルド時間
- **イメージサイズ**: 既存のAPI/Adminサーバーと同等のイメージサイズ
- **起動時間**: 既存のAPI/Adminサーバーと同等の起動時間

### 4.2 互換性
- **既存機能の維持**: Docker環境でも既存のJobQueueサーバーの機能が正常に動作する
- **Redis接続**: Docker環境でのRedis接続が正常に動作する
- **設定ファイル**: 既存の設定ファイル（`config/develop/config.yaml`等）との互換性を維持
- **ボリュームマウント**: 既存のボリュームマウント設定との互換性を維持

### 4.3 セキュリティ
- **非rootユーザー**: 既存のAPI/Adminサーバーと同様に非rootユーザー（appuser）で実行
- **パッケージの更新**: Alpineのパッケージを最新の状態に保つ
- **CGO無効化**: CGO_ENABLED=0により、セキュリティリスクを低減

### 4.4 保守性
- **設定の統一**: 既存のAPI/Adminサーバーと同様の設定構造を維持
- **ドキュメントの整合性**: すべてのドキュメントでJobQueueサーバーのDocker対応を記載
- **コードの一貫性**: 既存のDockerfileとdocker-composeファイルと同様の構造を維持

### 4.5 動作環境
- **開発環境**: `APP_ENV=develop`を想定
- **ネットワーク**: 既存の`postgres-network`と`redis-network`を使用
- **依存関係**: PostgreSQLとRedisコンテナが起動している必要がある

## 5. 制約事項

### 5.1 既存システムとの関係
- **PostgreSQL接続**: 現在のジョブ処理ではPostgreSQLを使用していないが、将来的にジョブ処理でPostgreSQLを使用する可能性が高いため、拡張性を考慮してデータベースDSN環境変数を設定する（既存のAPI/Adminサーバーと同様）
- **Redis接続**: Docker環境では`redis:6379`に接続する必要がある（`REDIS_JOBQUEUE_ADDR=redis:6379`）
- **設定ファイル**: 既存の設定ファイル（`config/develop/config.yaml`等）との互換性を維持
- **ボリュームマウント**: 既存のボリュームマウント設定との互換性を維持

### 5.2 技術スタック
- **Goバージョン**: Go 1.24を維持
- **ベースイメージ**: `golang:1.24-alpine`（ビルドステージ）、`alpine:latest`（実行ステージ）
- **CGO設定**: CGO_ENABLED=0を維持
- **ビルドツール**: Alpine環境で必要なビルドツールをインストール

### 5.3 依存関係
- **Redis**: AsynqサーバーがRedisからジョブを取得するため、Redisコンテナが起動している必要がある
- **ネットワーク**: `postgres-network`と`redis-network`が存在する必要がある（external: true）

### 5.4 運用上の制約
- **起動順序**: PostgreSQLとRedisコンテナを起動してから、JobQueueサーバーを起動する必要がある
- **ヘルスチェック**: `/health`エンドポイントが正常に動作する必要がある
- **ログ出力**: ログファイルは`./logs`ディレクトリに出力される

## 6. 受け入れ基準

### 6.1 Dockerfileの作成
- [ ] `server/Dockerfile.jobqueue`が作成されている
- [ ] ベースイメージが`golang:1.24-alpine`（ビルドステージ）と`alpine:latest`（実行ステージ）である
- [ ] CGO_ENABLED=0でビルドされている
- [ ] 必要なパッケージ（`ca-certificates`, `tzdata`, `wget`）がインストールされている
- [ ] 非rootユーザー（appuser）が作成されている
- [ ] ディレクトリ構造が既存のAPI/Adminサーバーと同様である
- [ ] ポート8082がEXPOSEされている
- [ ] CMDが`["./jobqueue"]`である

### 6.2 docker-compose設定ファイルの作成
- [ ] `docker-compose.jobqueue.yml`が作成されている
- [ ] サービス名が`jobqueue`である
- [ ] コンテナ名が`jobqueue`である
- [ ] `docker-compose.api.yml`のコンテナ名が`api`に修正されている
- [ ] `docker-compose.admin.yml`のコンテナ名が`admin`に修正されている
- [ ] `docker-compose.client.yml`のコンテナ名が`client`に修正されている
- [ ] ビルド設定が`context: ./server`と`dockerfile: Dockerfile.jobqueue`である
- [ ] ポート設定が`"8082:8082"`である
- [ ] 環境変数が適切に設定されている（`APP_ENV=develop`、`REDIS_JOBQUEUE_ADDR=redis:6379`等）
- [ ] ボリュームマウントが適切に設定されている
- [ ] ネットワーク設定が`postgres-network`と`redis-network`である
- [ ] ヘルスチェックが適切に設定されている

### 6.3 ビルドと動作確認
- [ ] `docker-compose.jobqueue.yml`でイメージが正常にビルドできる
- [ ] ビルドしたイメージからコンテナが正常に起動できる
- [ ] JobQueueサーバーが正常に動作する（HTTPサーバーとAsynqサーバー）
- [ ] `/health`エンドポイントが正常に動作する
- [ ] Redis接続が正常に動作する
- [ ] ヘルスチェックが正常に動作する

### 6.4 ドキュメント
- [ ] `docs/ja/Docker.md`が更新されている
- [ ] `docs/en/Docker.md`が更新されている
- [ ] `README.md`が更新されている
- [ ] `README.ja.md`が更新されている
- [ ] すべてのドキュメントでJobQueueサーバーのDocker対応が記載されている

### 6.5 既存機能との整合性
- [ ] 既存のAPI/AdminサーバーのDocker設定に影響がない
- [ ] 既存のdocker-composeファイルに影響がない
- [ ] 既存のDockerfileに影響がない
- [ ] 既存のドキュメントの他の部分に影響がない

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### 新規作成するファイル
- `server/Dockerfile.jobqueue`: JobQueueサーバー用Dockerfile（新規作成）
- `docker-compose.jobqueue.yml`: JobQueueサーバー用docker-compose設定（新規作成）

#### 更新するファイル
- `docker-compose.api.yml`: コンテナ名を`api-develop`から`api`に変更
- `docker-compose.admin.yml`: コンテナ名を`admin-develop`から`admin`に変更
- `docker-compose.client.yml`: コンテナ名を`client-develop`から`client`に変更
- `docs/ja/Docker.md`: Docker環境の説明を更新
- `docs/en/Docker.md`: Docker環境の説明を更新
- `README.md`: Docker環境の説明を更新
- `README.ja.md`: Docker環境の説明を更新

### 7.2 既存ファイルの扱い
- `server/Dockerfile`: 変更なし（APIサーバー用）
- `server/Dockerfile.admin`: 変更なし（Adminサーバー用）
- `docker-compose.api.yml`: コンテナ名を`api-develop`から`api`に変更
- `docker-compose.admin.yml`: コンテナ名を`admin-develop`から`admin`に変更
- `docker-compose.client.yml`: コンテナ名を`client-develop`から`client`に変更
- その他のdocker-composeファイル: 変更なし

### 7.3 既存機能への影響
- **既存のサーバー**: 影響なし（コンテナ名の変更のみ、機能に影響なし）
- **既存のDocker設定**: コンテナ名の変更により、既存のコンテナを停止して再起動する必要がある可能性がある
- **既存のドキュメント**: 更新が必要（JobQueueサーバーのDocker対応を追加、コンテナ名の変更を反映）

## 8. 実装上の注意事項

### 8.1 Dockerfileの作成
- **既存Dockerfileとの整合性**: `server/Dockerfile`と`server/Dockerfile.admin`を参考に、同様の構造を維持
- **ビルドコマンド**: `./cmd/jobqueue/main.go`をビルドする
- **バイナリ名**: `jobqueue`とする（既存の`server`や`admin`と同様）
- **パッケージインストール**: `wget`をインストールしてヘルスチェックに使用（既存のAPI/Adminサーバーと同様）

### 8.2 docker-compose設定ファイルの作成
- **既存docker-composeファイルとの整合性**: `docker-compose.api.yml`と`docker-compose.admin.yml`を参考に、同様の構造を維持
- **環境変数**: 既存のAPI/Adminサーバーと同様の環境変数を設定
- **ボリュームマウント**: 既存のAPI/Adminサーバーと同様のボリュームマウントを設定
- **ネットワーク**: 既存の`postgres-network`と`redis-network`を使用（external: true）
- **ヘルスチェック**: `/health`エンドポイントを使用（既存のAPI/Adminサーバーと同様）
- **コンテナ名の統一**: 既存の`docker-compose.api.yml`、`docker-compose.admin.yml`、`docker-compose.client.yml`のコンテナ名から`-develop`を削除し、統一的な命名規則に従う

### 8.3 ビルドと動作確認
- **ビルド確認**: `docker-compose.jobqueue.yml`でイメージが正常にビルドできることを確認
- **起動確認**: ビルドしたイメージからコンテナが正常に起動できることを確認
- **接続確認**: Redisへの接続が正常に動作することを確認
- **ヘルスチェック**: `/health`エンドポイントが正常に動作することを確認
- **ジョブ処理**: Asynqサーバーが正常にジョブを処理できることを確認

### 8.4 ドキュメント整備
- **Dockerfile構成表**: `docs/ja/Docker.md`と`docs/en/Docker.md`のDockerfile構成表に`server/Dockerfile.jobqueue`を追加
- **Docker Compose設定ファイル一覧**: `docker-compose.jobqueue.yml`を追加
- **ビルド・起動・停止コマンド**: JobQueueサーバー用のコマンドを追加
- **既存サービスとの統合**: JobQueueサーバーの起動順序を追加
- **サービス間通信**: JobQueueサーバーとRedisの接続を追加
- **README**: Docker環境の説明にJobQueueサーバーを追加

### 8.5 既存機能との整合性確認
- **既存のDocker設定**: 既存のAPI/AdminサーバーのDocker設定に影響がないことを確認（コンテナ名の変更のみ）
- **既存のdocker-composeファイル**: コンテナ名の変更が正しく反映されていることを確認
- **既存のDockerfile**: 既存のDockerfileに影響がないことを確認
- **既存のドキュメント**: 既存のドキュメントの他の部分に影響がないことを確認
- **コンテナ名の変更**: 既存のコンテナを停止して再起動する必要がある可能性があることを確認

## 9. 参考情報

### 9.1 関連Issue
- GitHub Issue #163: JobQueueサーバーのDocker対応

### 9.2 既存ドキュメント
- `docs/ja/Docker.md`: Docker環境の説明（日本語）
- `docs/en/Docker.md`: Docker環境の説明（英語）
- `README.md`: プロジェクトの説明（英語）
- `README.ja.md`: プロジェクトの説明（日本語）
- `server/Dockerfile`: APIサーバー用Dockerfile
- `server/Dockerfile.admin`: Adminサーバー用Dockerfile
- `docker-compose.api.yml`: APIサーバー用docker-compose設定
- `docker-compose.admin.yml`: Adminサーバー用docker-compose設定

### 9.3 技術スタック
- **Go**: 1.24
- **ベースイメージ**: `golang:1.24-alpine`（ビルドステージ）、`alpine:latest`（実行ステージ）
- **CGO**: CGO_ENABLED=0
- **Redis**: AsynqサーバーがRedisからジョブを取得
- **HTTPサーバー**: ポート8082、`/health`エンドポイント提供

### 9.4 参考リンク
- Docker公式ドキュメント: https://docs.docker.com/
- Docker Compose公式ドキュメント: https://docs.docker.com/compose/
- Go公式ドキュメント: https://go.dev/doc/
- Alpine Linux公式ドキュメント: https://alpinelinux.org/documentation/
