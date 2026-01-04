# Docker化実装タスク一覧

## 概要
Docker化の実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: APIサーバーDockerfile作成

#### タスク 1.1: server/Dockerfileの作成
**目的**: APIサーバー用のDockerfileを作成

**作業内容**:
- `server/Dockerfile`を新規作成
- マルチステージビルドを実装:
  - ビルドステージ: `golang:1.21-alpine`を使用
  - 実行ステージ: `alpine:latest`を使用
- CGO_ENABLED=0で静的リンクビルド
- 非rootユーザー（appuser）で実行
- ポート8080を公開
- 作業ディレクトリ: `/app`
- ログディレクトリ: `/app/logs`を作成（デフォルト設定でログファイルは`logs/`に出力されるため）

**受け入れ基準**:
- `server/Dockerfile`が作成されている
- マルチステージビルドが実装されている
- 非rootユーザーで実行される設定になっている
- ポート8080が公開されている

---

#### タスク 1.2: server/.dockerignoreの作成
**目的**: APIサーバー用の.dockerignoreファイルを作成

**作業内容**:
- `server/.dockerignore`を新規作成
- 以下のファイル・ディレクトリを除外:
  - `.git`, `.gitignore`
  - `*.md`
  - `*.test.go`
  - `*.sum`
  - `vendor/`
  - `.env`, `.env.local`
  - `*.log`
  - `coverage/`
  - `.idea/`, `.vscode/`

**受け入れ基準**:
- `server/.dockerignore`が作成されている
- 不要なファイルが除外されている

---

### Phase 2: AdminサーバーDockerfile作成

#### タスク 2.1: server/Dockerfile.adminの作成
**目的**: Adminサーバー用のDockerfileを作成

**作業内容**:
- `server/Dockerfile.admin`を新規作成
- APIサーバーと同様のマルチステージビルドを実装
- `cmd/admin/main.go`をビルド
- ポート8081を公開
- 非rootユーザー（appuser）で実行
- ログディレクトリ: `/app/logs`を作成（デフォルト設定でログファイルは`logs/`に出力されるため）

**受け入れ基準**:
- `server/Dockerfile.admin`が作成されている
- マルチステージビルドが実装されている
- ポート8081が公開されている

---

### Phase 3: クライアントサーバーDockerfile作成

#### タスク 3.1: client/Dockerfileの作成
**目的**: クライアントサーバー用のDockerfileを作成

**作業内容**:
- `client/Dockerfile`を新規作成
- マルチステージビルドを実装:
  - 依存関係インストールステージ: `node:22-alpine`
  - ビルドステージ: Next.jsアプリケーションをビルド
  - 開発モード実行ステージ: `dev`ターゲット
  - 本番モード実行ステージ: `production`ターゲット
- 開発モード: ホットリロード対応
- 本番モード: 最適化されたビルド成果物のみを含む
- 非rootユーザー（nodeuser）で実行（本番モード）
- ポート3000を公開

**受け入れ基準**:
- `client/Dockerfile`が作成されている
- マルチステージビルドが実装されている
- 開発モードと本番モードの両方が実装されている
- ポート3000が公開されている

---

#### タスク 3.2: client/.dockerignoreの作成
**目的**: クライアントサーバー用の.dockerignoreファイルを作成

**作業内容**:
- `client/.dockerignore`を新規作成
- 以下のファイル・ディレクトリを除外:
  - `.git`, `.gitignore`
  - `*.md`
  - `node_modules`
  - `.next`
  - `.env`, `.env.local`
  - `*.log`
  - `coverage/`
  - `.idea/`, `.vscode/`

**受け入れ基準**:
- `client/.dockerignore`が作成されている
- 不要なファイルが除外されている

---

### Phase 4: Docker Compose設定ファイル作成（開発環境）

#### タスク 4.1: docker-compose.api.develop.ymlの作成
**目的**: APIサーバー用のDocker Compose設定ファイル（開発環境）を作成

**作業内容**:
- `docker-compose.api.develop.yml`を新規作成
- `api`サービスを定義:
  - ビルドコンテキスト: `./server`
  - Dockerfile: `Dockerfile`
  - コンテナ名: `api-develop`
  - ポートマッピング: `8080:8080`
  - 環境変数: `APP_ENV=develop`
  - ボリュームマウント:
    - `./config/develop:/app/config/develop:ro`
    - `./server/data:/app/data`
    - `./logs:/app/logs`
    - `./logs:/app/logs`
  - ネットワーク: `postgres-network`, `redis-network`（外部ネットワーク）
  - ヘルスチェック設定
  - 再起動ポリシー: `unless-stopped`

**受け入れ基準**:
- `docker-compose.api.develop.yml`が作成されている
- APIサーバーサービスが正しく定義されている
- 既存のPostgreSQL、Redisネットワークと統合されている

---

#### タスク 4.2: docker-compose.client.develop.ymlの作成
**目的**: クライアントサーバー用のDocker Compose設定ファイル（開発環境）を作成

**作業内容**:
- `docker-compose.client.develop.yml`を新規作成
- `client`サービスを定義:
  - ビルドコンテキスト: `./client`
  - Dockerfile: `Dockerfile`
  - ターゲット: `dev`
  - ビルド引数: `NEXT_PUBLIC_API_URL=http://api:8080`
  - コンテナ名: `client-develop`
  - ポートマッピング: `3000:3000`
  - 環境変数:
    - `NODE_ENV=development`
    - `NEXT_PUBLIC_API_URL=http://api:8080`
  - ボリュームマウント:
    - `./client:/app`
    - `/app/node_modules`（ボリューム）
    - `/app/.next`（ボリューム）
  - ネットワーク: `postgres-network`（外部ネットワーク）
  - 依存関係: `api`サービス（depends_on）

**受け入れ基準**:
- `docker-compose.client.develop.yml`が作成されている
- クライアントサーバーサービスが正しく定義されている
- 開発モード用のボリュームマウントが設定されている
- APIサーバーへの依存関係が設定されている

---

#### タスク 4.3: docker-compose.admin.develop.ymlの作成
**目的**: Adminサーバー用のDocker Compose設定ファイル（開発環境）を作成

**作業内容**:
- `docker-compose.admin.develop.yml`を新規作成
- `admin`サービスを定義:
  - ビルドコンテキスト: `./server`
  - Dockerfile: `Dockerfile.admin`
  - コンテナ名: `admin-develop`
  - ポートマッピング: `8081:8081`
  - 環境変数: `APP_ENV=develop`
  - ボリュームマウント:
    - `./config/develop:/app/config/develop:ro`
    - `./server/data:/app/data`
    - `./logs:/app/logs`
    - `./logs:/app/logs`
  - ネットワーク: `postgres-network`（外部ネットワーク）
  - ヘルスチェック設定
  - 再起動ポリシー: `unless-stopped`

**受け入れ基準**:
- `docker-compose.admin.develop.yml`が作成されている
- Adminサーバーサービスが正しく定義されている
- 既存のPostgreSQLネットワークと統合されている

---

### Phase 5: Docker Compose設定ファイル作成（ステージング環境）

#### タスク 5.1: docker-compose.api.staging.ymlの作成
**目的**: APIサーバー用のDocker Compose設定ファイル（ステージング環境）を作成

**作業内容**:
- `docker-compose.api.staging.yml`を新規作成
- 開発環境用ファイルをベースに作成
- 環境変数: `APP_ENV=staging`
- 設定ファイルパス: `./config/staging:/app/config/staging:ro`
- コンテナ名: `api-staging`

**受け入れ基準**:
- `docker-compose.api.staging.yml`が作成されている
- ステージング環境用の設定が正しく定義されている

---

#### タスク 5.2: docker-compose.client.staging.ymlの作成
**目的**: クライアントサーバー用のDocker Compose設定ファイル（ステージング環境）を作成

**作業内容**:
- `docker-compose.client.staging.yml`を新規作成
- 開発環境用ファイルをベースに作成
- 環境変数: `NODE_ENV=production`（本番モード）
- ターゲット: `production`
- コンテナ名: `client-staging`
- 開発モード用のボリュームマウントを削除

**受け入れ基準**:
- `docker-compose.client.staging.yml`が作成されている
- ステージング環境用の設定が正しく定義されている
- 本番モード用の設定になっている

---

#### タスク 5.3: docker-compose.admin.staging.ymlの作成
**目的**: Adminサーバー用のDocker Compose設定ファイル（ステージング環境）を作成

**作業内容**:
- `docker-compose.admin.staging.yml`を新規作成
- 開発環境用ファイルをベースに作成
- 環境変数: `APP_ENV=staging`
- 設定ファイルパス: `./config/staging:/app/config/staging:ro`
- コンテナ名: `admin-staging`

**受け入れ基準**:
- `docker-compose.admin.staging.yml`が作成されている
- ステージング環境用の設定が正しく定義されている

---

### Phase 6: Docker Compose設定ファイル作成（本番環境）

#### タスク 6.1: docker-compose.api.production.ymlの作成
**目的**: APIサーバー用のDocker Compose設定ファイル（本番環境）を作成

**作業内容**:
- `docker-compose.api.production.yml`を新規作成
- ステージング環境用ファイルをベースに作成
- 環境変数: `APP_ENV=production`
- 設定ファイルパス: `./config/production:/app/config/production:ro`
- コンテナ名: `api-production`
- リソース制限を追加（必要に応じて）

**受け入れ基準**:
- `docker-compose.api.production.yml`が作成されている
- 本番環境用の設定が正しく定義されている

---

#### タスク 6.2: docker-compose.client.production.ymlの作成
**目的**: クライアントサーバー用のDocker Compose設定ファイル（本番環境）を作成

**作業内容**:
- `docker-compose.client.production.yml`を新規作成
- ステージング環境用ファイルをベースに作成
- 環境変数: `NODE_ENV=production`
- コンテナ名: `client-production`
- リソース制限を追加（必要に応じて）

**受け入れ基準**:
- `docker-compose.client.production.yml`が作成されている
- 本番環境用の設定が正しく定義されている

---

#### タスク 6.3: docker-compose.admin.production.ymlの作成
**目的**: Adminサーバー用のDocker Compose設定ファイル（本番環境）を作成

**作業内容**:
- `docker-compose.admin.production.yml`を新規作成
- ステージング環境用ファイルをベースに作成
- 環境変数: `APP_ENV=production`
- 設定ファイルパス: `./config/production:/app/config/production:ro`
- コンテナ名: `admin-production`
- リソース制限を追加（必要に応じて）

**受け入れ基準**:
- `docker-compose.admin.production.yml`が作成されている
- 本番環境用の設定が正しく定義されている

---

### Phase 7: 動作確認（開発環境）

#### タスク 7.1: APIサーバーDockerイメージのビルド確認
**目的**: APIサーバーのDockerイメージが正常にビルドされることを確認

**作業内容**:
- `docker-compose -f docker-compose.api.develop.yml build`を実行
- ビルドが正常に完了することを確認
- エラーがないことを確認

**受け入れ基準**:
- Dockerイメージが正常にビルドされる
- ビルドエラーがない

---

#### タスク 7.2: APIサーバーコンテナの起動確認
**目的**: APIサーバーコンテナが正常に起動することを確認

**作業内容**:
- 既存のPostgreSQL、Redisコンテナが起動していることを確認
- `docker-compose -f docker-compose.api.develop.yml up -d`を実行
- コンテナが正常に起動することを確認
- ポート8080でアクセスできることを確認
- ヘルスチェックが正常に動作することを確認
- ログを確認してエラーがないことを確認

**受け入れ基準**:
- APIサーバーコンテナが正常に起動する
- ポート8080でアクセスできる
- ヘルスチェックが正常に動作する
- PostgreSQL、Redisコンテナと通信できる

---

#### タスク 7.3: AdminサーバーDockerイメージのビルド確認
**目的**: AdminサーバーのDockerイメージが正常にビルドされることを確認

**作業内容**:
- `docker-compose -f docker-compose.admin.develop.yml build`を実行
- ビルドが正常に完了することを確認
- エラーがないことを確認

**受け入れ基準**:
- Dockerイメージが正常にビルドされる
- ビルドエラーがない

---

#### タスク 7.4: Adminサーバーコンテナの起動確認
**目的**: Adminサーバーコンテナが正常に起動することを確認

**作業内容**:
- 既存のPostgreSQLコンテナが起動していることを確認
- `docker-compose -f docker-compose.admin.develop.yml up -d`を実行
- コンテナが正常に起動することを確認
- ポート8081でアクセスできることを確認
- ヘルスチェックが正常に動作することを確認
- ログを確認してエラーがないことを確認

**受け入れ基準**:
- Adminサーバーコンテナが正常に起動する
- ポート8081でアクセスできる
- ヘルスチェックが正常に動作する
- PostgreSQLコンテナと通信できる

---

#### タスク 7.5: クライアントサーバーDockerイメージのビルド確認
**目的**: クライアントサーバーのDockerイメージが正常にビルドされることを確認

**作業内容**:
- `docker-compose -f docker-compose.client.develop.yml build`を実行
- ビルドが正常に完了することを確認
- エラーがないことを確認

**受け入れ基準**:
- Dockerイメージが正常にビルドされる
- ビルドエラーがない

---

#### タスク 7.6: クライアントサーバーコンテナの起動確認
**目的**: クライアントサーバーコンテナが正常に起動することを確認

**作業内容**:
- APIサーバーコンテナが起動していることを確認
- `docker-compose -f docker-compose.client.develop.yml up -d`を実行
- コンテナが正常に起動することを確認
- ポート3000でアクセスできることを確認
- 開発モードでホットリロードが機能することを確認
- APIサーバーコンテナと通信できることを確認
- ログを確認してエラーがないことを確認

**受け入れ基準**:
- クライアントサーバーコンテナが正常に起動する
- ポート3000でアクセスできる
- 開発モードでホットリロードが機能する
- APIサーバーコンテナと通信できる

---

#### タスク 7.7: 3つのサーバーの同時起動確認
**目的**: 3つのサーバーを同時に起動できることを確認

**作業内容**:
- 既存のPostgreSQL、Redisコンテナが起動していることを確認
- APIサーバーを起動
- Adminサーバーを起動
- クライアントサーバーを起動
- すべてのコンテナが正常に起動することを確認
- 各サーバーが正常に動作することを確認

**受け入れ基準**:
- 3つのサーバーを同時に起動できる
- すべてのコンテナが正常に動作する

---

### Phase 8: 動作確認（環境別）

#### タスク 8.1: ステージング環境での動作確認
**目的**: ステージング環境で各サーバーが正常に動作することを確認

**作業内容**:
- `docker-compose -f docker-compose.api.staging.yml up -d`を実行
- `docker-compose -f docker-compose.admin.staging.yml up -d`を実行
- `docker-compose -f docker-compose.client.staging.yml up -d`を実行
- 各サーバーが正常に起動することを確認
- 環境別設定ファイルが正しく読み込まれることを確認

**受け入れ基準**:
- ステージング環境で各サーバーが正常に起動する
- 環境別設定ファイルが正しく読み込まれる

---

#### タスク 8.2: 本番環境での動作確認
**目的**: 本番環境で各サーバーが正常に動作することを確認

**作業内容**:
- `docker-compose -f docker-compose.api.production.yml up -d`を実行
- `docker-compose -f docker-compose.admin.production.yml up -d`を実行
- `docker-compose -f docker-compose.client.production.yml up -d`を実行
- 各サーバーが正常に起動することを確認
- 環境別設定ファイルが正しく読み込まれることを確認
- 本番モードで最適化されたビルドが実行されることを確認（クライアント）

**受け入れ基準**:
- 本番環境で各サーバーが正常に起動する
- 環境別設定ファイルが正しく読み込まれる
- 本番モードで最適化されたビルドが実行される

---

### Phase 9: Dockerイメージの最適化確認

#### タスク 9.1: イメージサイズの確認
**目的**: 最終イメージサイズが最小化されていることを確認

**作業内容**:
- `docker images`でイメージサイズを確認
- マルチステージビルドによりイメージサイズが最小化されていることを確認
- Alpine Linuxベースの軽量イメージになっていることを確認

**受け入れ基準**:
- イメージサイズが適切に最小化されている
- Alpine Linuxベースの軽量イメージになっている

---

#### タスク 9.2: ビルドキャッシュの確認
**目的**: ビルドキャッシュが正しく活用されていることを確認

**作業内容**:
- 2回目のビルドでキャッシュが使用されることを確認
- 依存関係のインストールがキャッシュされることを確認

**受け入れ基準**:
- ビルドキャッシュが正しく活用されている
- 2回目のビルドが高速化されている

---

### Phase 10: デプロイメント準備

#### タスク 10.1: イメージのタグ付け確認
**目的**: Dockerイメージにタグを付与できることを確認

**作業内容**:
- `docker tag`コマンドでイメージにタグを付与
- タグ形式: `{service-name}:{version}`または`{service-name}:latest`
- タグが正しく付与されることを確認

**受け入れ基準**:
- Dockerイメージにタグを付与できる
- タグが正しく付与される

---

#### タスク 10.2: コンテナレジストリへのプッシュ手順のドキュメント化
**目的**: コンテナレジストリへのプッシュ手順をドキュメントに記載（実際のプッシュ動作確認は不要）

**作業内容**:
- `docs/Docker.md`にコンテナレジストリへのプッシュ手順を記載:
  - AWS ECRへのプッシュ手順
  - Tencent Cloud TCRへのプッシュ手順
  - Docker Hubへのプッシュ手順（オプション）
  - 認証方法の説明
  - プッシュコマンドの例
- 実際のプッシュ動作確認は不要（ドキュメント記載のみ）

**受け入れ基準**:
- `docs/Docker.md`にコンテナレジストリへのプッシュ手順が記載されている
- 認証方法が説明されている
- プッシュコマンドの例が記載されている

---

### Phase 11: ドキュメント整備

#### タスク 11.1: docs/Docker.mdの作成
**目的**: Docker化に関する詳細ドキュメントを作成

**作業内容**:
- `docs/Docker.md`を新規作成
- 以下の内容を含める:
  - Docker化の概要
  - 前提条件
  - Dockerfileの説明
  - Docker Compose設定の説明
  - ビルド・起動・停止のコマンド
  - 環境別の起動方法
  - 既存サービスとの統合方法
  - トラブルシューティング情報
  - 本番環境へのデプロイ手順
  - コンテナレジストリへのプッシュ手順（AWS ECR、Tencent Cloud TCR、Docker Hub）

**受け入れ基準**:
- `docs/Docker.md`が作成されている
- 上記の内容が全て記載されている
- コンテナレジストリへのプッシュ手順が記載されている

---

#### タスク 11.2: README.mdの更新
**目的**: README.mdにDocker化に関する情報を追記

**作業内容**:
- `README.md`を確認
- Docker化に関する簡単な説明を追記
- Docker環境での起動方法を追記
- `docs/Docker.md`へのリンクを追加

**受け入れ基準**:
- `README.md`が更新されている
- Docker化に関する情報が追記されている
- `docs/Docker.md`へのリンクが追加されている

---

## 実装順序の推奨

1. **Phase 1-3**: Dockerfileの作成（APIサーバー、Adminサーバー、クライアントサーバー）
2. **Phase 4-6**: Docker Compose設定ファイルの作成（環境別）
3. **Phase 7**: 動作確認（開発環境）
4. **Phase 8**: 動作確認（環境別）
5. **Phase 9**: Dockerイメージの最適化確認
6. **Phase 10**: デプロイメント準備
7. **Phase 11**: ドキュメント整備

## 注意事項

- **既存サービスの起動**: 各サーバーを起動する前に、既存のPostgreSQL、Redisコンテナが起動している必要がある
- **ネットワーク統合**: 既存のdocker-composeネットワークを外部ネットワークとして参照する
- **環境変数**: `APP_ENV`環境変数で環境を切り替え（develop/staging/production）
- **ボリュームマウント**: 設定ファイル、データディレクトリ、ログディレクトリを適切にマウントする
  - 設定ファイル: `./config/{env}:/app/config/{env}:ro`
  - データディレクトリ: `./server/data:/app/data`
  - ログディレクトリ: `./logs:/app/logs`（デフォルト設定でログファイルは`logs/`に出力されるため）
- **起動順序**: APIサーバー → Adminサーバー → クライアントサーバーの順で起動することを推奨
- **ポート競合**: ポート8080、8081、3000が使用可能であることを確認
- **イメージサイズ**: マルチステージビルドによりイメージサイズを最小化する
- **セキュリティ**: 非rootユーザーで実行する設定を維持する
