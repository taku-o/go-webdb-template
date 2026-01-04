# Docker化要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #83
- **Issueタイトル**: APIサーバー、クライアントサーバー、GoAdminサーバーをDocker上で動作させる
- **Feature名**: 0042-to-docker
- **作成日**: 2025-01-27

### 1.2 目的
APIサーバー、クライアントサーバー、GoAdminサーバーをDockerコンテナ上で動作させることで、開発環境の統一と本番環境（AWSまたはTencent Cloud）へのデプロイを容易にする。
Dockerイメージを作成できるようにすることで、任意のコンテナ実行環境にデプロイ可能な状態を実現する。

### 1.3 スコープ
- APIサーバー（Go）のDocker化
- GoAdminサーバー（Go）のDocker化
- クライアントサーバー（Next.js）のDocker化
- Docker Composeによる3つのサーバーの一元管理
- Dockerイメージの最適化（マルチステージビルド）
- 既存のPostgreSQL、Redisコンテナとの統合
- 環境別設定（develop/staging/production）のサポート
- 本番環境（AWS/Tencent Cloud）へのデプロイ準備

**本実装の範囲外**:
- 本番環境での実際のデプロイ（準備のみ）
- CI/CDパイプラインの構築（将来の拡張項目）
- Kubernetesマニフェストの作成（将来の拡張項目）
- 監視・ログ集約システムの構築（将来の拡張項目）

## 2. 背景・現状分析

### 2.1 現在の実装
- **APIサーバー**: `server/cmd/server/main.go`（ポート8080）
  - 起動コマンド: `APP_ENV=develop go run cmd/server/main.go`
  - 設定ファイル: `config/{env}/config.yaml`, `config/{env}/database.yaml`, `config/{env}/cacheserver.yaml`
  - データベース: SQLite（開発環境）またはPostgreSQL/MySQL（本番想定）
- **GoAdminサーバー**: `server/cmd/admin/main.go`（ポート8081）
  - 起動コマンド: `APP_ENV=develop go run cmd/admin/main.go`
  - 設定ファイル: APIサーバーと同様
- **クライアントサーバー**: Next.js（ポート3000）
  - 起動コマンド: `npm run dev`
  - ビルドコマンド: `npm run build`
- **既存Docker環境**:
  - PostgreSQL: `docker-compose.postgres.yml`
  - Redis: `docker-compose.redis.yml`, `docker-compose.redis-cluster.yml`
  - CloudBeaver: `docker-compose.cloudbeaver.yml`
  - Metabase: `docker-compose.metabase.yml`
  - Apache Superset: `docker-compose.apache-superset.yml`
  - Mailpit: `docker-compose.mailpit.yml`

### 2.2 課題点
1. **開発環境の不統一**: 開発者ごとにGo/Node.jsのバージョンが異なる可能性がある
2. **デプロイの困難さ**: 本番環境（AWS/Tencent Cloud）へのデプロイ時に環境構築が複雑
3. **環境の再現性**: 開発環境と本番環境の差異による不具合のリスク
4. **スケーラビリティ**: 複数インスタンスの起動が困難
5. **既存インフラとの統合**: Docker化されていないアプリケーションサーバーと既存のDockerコンテナの連携が複雑

### 2.3 本実装による改善点
1. **開発環境の統一**: Dockerコンテナにより全開発者が同一環境で開発可能
2. **デプロイの容易さ**: Dockerイメージにより任意のコンテナ実行環境にデプロイ可能
3. **環境の再現性**: 開発環境と本番環境で同じコンテナイメージを使用可能
4. **スケーラビリティ**: コンテナを複数起動することで水平スケーリングが容易
5. **既存インフラとの統合**: Docker Composeにより既存のPostgreSQL、Redisと一元管理可能

## 3. 機能要件

### 3.1 APIサーバーのDocker化

#### 3.1.1 Dockerfileの作成
- **ファイル**: `server/Dockerfile`
- **内容**: 
  - マルチステージビルドを使用
  - ビルドステージ: `golang:1.21-alpine`を使用してGoアプリケーションをビルド
  - 実行ステージ: `alpine:latest`を使用して軽量な実行環境を構築
  - CGO_ENABLED=0で静的リンクビルド
  - 設定ファイルとデータディレクトリをコピー
  - ポート8080を公開
  - 非rootユーザーで実行（セキュリティベストプラクティス）

#### 3.1.2 環境変数と設定ファイルの管理
- **環境変数**: `APP_ENV`環境変数で環境を指定（develop/staging/production）
- **設定ファイル**: `config/{env}/`ディレクトリをボリュームマウント
- **データディレクトリ**: `server/data/`ディレクトリをボリュームマウント（SQLiteデータベースファイルの永続化）

#### 3.1.3 既存サービスとの連携
- **PostgreSQL**: Dockerネットワーク経由で`postgres`サービスに接続
- **Redis**: Dockerネットワーク経由で`redis`サービスに接続
- **ネットワーク**: 既存のdocker-composeネットワークを使用

### 3.2 AdminサーバーのDocker化

#### 3.2.1 Dockerfileの作成
- **ファイル**: `server/Dockerfile.admin`
- **内容**: 
  - APIサーバーと同様のマルチステージビルド
  - `cmd/admin/main.go`をビルド
  - ポート8081を公開
  - 設定ファイルとデータディレクトリをコピー

#### 3.2.2 環境変数と設定ファイルの管理
- **環境変数**: `APP_ENV`環境変数で環境を指定
- **設定ファイル**: `config/{env}/`ディレクトリをボリュームマウント
- **データディレクトリ**: `server/data/`ディレクトリをボリュームマウント

#### 3.2.3 既存サービスとの連携
- **PostgreSQL**: Dockerネットワーク経由で`postgres`サービスに接続
- **APIサーバー**: Dockerネットワーク経由で`api`サービスに接続（必要に応じて）

### 3.3 クライアントサーバーのDocker化

#### 3.3.1 Dockerfileの作成
- **ファイル**: `client/Dockerfile`
- **内容**: 
  - マルチステージビルドを使用
  - 依存関係インストールステージ: `node:22-alpine`を使用
  - ビルドステージ: Next.jsアプリケーションをビルド
  - 実行ステージ: 最適化されたビルド成果物のみを含む
  - ポート3000を公開
  - 開発モードと本番モードの両方に対応

#### 3.3.2 環境変数の管理
- **環境変数**: `NEXT_PUBLIC_API_URL`でAPIサーバーのURLを指定
- **デフォルト値**: Dockerネットワーク経由で`http://api:8080`に接続
- **開発モード**: ホットリロードを有効化

#### 3.3.3 APIサーバーとの連携
- **APIサーバー**: Dockerネットワーク経由で`api`サービスに接続
- **CORS設定**: APIサーバーのCORS設定でクライアントからのアクセスを許可

### 3.4 Docker Compose統合

#### 3.4.1 docker-composeファイルの構成
- **命名規則**: `docker-compose.{service}.{env}.yml`
  - `{service}`: `api`, `client`, `admin`
  - `{env}`: `develop`, `staging`, `production`
- **ファイル一覧**:
  - `docker-compose.api.develop.yml`: APIサーバー（開発環境）
  - `docker-compose.api.staging.yml`: APIサーバー（ステージング環境）
  - `docker-compose.api.production.yml`: APIサーバー（本番環境）
  - `docker-compose.client.develop.yml`: クライアントサーバー（開発環境）
  - `docker-compose.client.staging.yml`: クライアントサーバー（ステージング環境）
  - `docker-compose.client.production.yml`: クライアントサーバー（本番環境）
  - `docker-compose.admin.develop.yml`: Adminサーバー（開発環境）
  - `docker-compose.admin.staging.yml`: Adminサーバー（ステージング環境）
  - `docker-compose.admin.production.yml`: Adminサーバー（本番環境）

#### 3.4.2 各docker-composeファイルの内容
- **APIサーバー用**:
  - `api`サービス: APIサーバー
  - 既存のPostgreSQL、Redisサービスとの統合
  - 環境変数とボリュームマウントの設定
- **クライアントサーバー用**:
  - `client`サービス: クライアントサーバー
  - APIサーバーへの接続設定（環境変数で指定）
  - 環境変数とボリュームマウントの設定
- **Adminサーバー用**:
  - `admin`サービス: Adminサーバー
  - 既存のPostgreSQLサービスとの統合
  - 環境変数とボリュームマウントの設定

#### 3.4.3 サービス間の連携
- **ネットワーク**: 既存のdocker-composeネットワークを使用（外部ネットワークとして定義）
- **サービス名**: 既存のサービス名（`postgres`, `redis`）を参照
- **依存関係**: docker-composeの`depends_on`は使用せず、起動順序は運用で管理
- **環境変数**: 各サービスで必要な環境変数を個別に設定

#### 3.4.4 既存サービスとの統合
- **PostgreSQL**: 既存の`docker-compose.postgres.yml`と同一ネットワークで通信
- **Redis**: 既存の`docker-compose.redis.yml`と同一ネットワークで通信
- **ネットワーク**: 外部ネットワークとして既存のdocker-composeネットワークを参照

### 3.5 Dockerイメージの最適化

#### 3.5.1 マルチステージビルド
- **ビルドステージ**: ビルドツールを含む
- **実行ステージ**: 実行に必要なファイルのみを含む
- **イメージサイズ**: 最終イメージサイズを最小化

#### 3.5.2 .dockerignoreファイル
- **APIサーバー**: `server/.dockerignore`
- **クライアント**: `client/.dockerignore`
- **内容**: 不要なファイル（テストファイル、ドキュメント、.git等）を除外

#### 3.5.3 セキュリティベストプラクティス
- **非rootユーザー**: コンテナ内で非rootユーザーで実行
- **最小権限**: 必要最小限の権限のみ付与
- **イメージスキャン**: 脆弱性スキャンに対応（将来の拡張項目）

### 3.6 デプロイメント準備

#### 3.6.1 イメージのタグ付け
- **タグ形式**: `{service-name}:{version}`または`{service-name}:latest`
- **バージョン管理**: Gitタグやコミットハッシュを使用

#### 3.6.2 コンテナレジストリへのプッシュ
- **対応レジストリ**: Docker Hub、AWS ECR、Tencent Cloud TCR等
- **認証情報**: 環境変数やシークレット管理で管理

#### 3.6.3 本番環境用設定
- **環境変数**: 本番環境用の環境変数を設定
- **ヘルスチェック**: Docker healthcheckを実装
- **リソース制限**: CPU・メモリ制限を設定

## 4. 非機能要件

### 4.1 Docker環境の前提条件
- DockerおよびDocker Composeがインストールされていること
- Dockerが正常に動作していること
- ポート8080、8081、3000が使用可能であること（他のサービスと競合しないこと）
- 十分なメモリとディスク容量があること

### 4.2 パフォーマンス
- **ビルド時間**: マルチステージビルドとキャッシュによりビルド時間を短縮
- **起動時間**: コンテナの起動時間を最小化
- **イメージサイズ**: 最終イメージサイズを最小化（Alpine Linuxベース）

### 4.3 セキュリティ
- **非rootユーザー**: コンテナ内で非rootユーザーで実行
- **イメージスキャン**: 脆弱性スキャンに対応（将来の拡張項目）
- **シークレット管理**: 機密情報は環境変数やシークレット管理で管理

### 4.4 データ永続化
- **設定ファイル**: ボリュームマウントにより設定ファイルを永続化
- **データベースファイル**: SQLiteデータベースファイルをボリュームマウントで永続化
- **ログファイル**: ログファイルをボリュームマウントで永続化（必要に応じて）

### 4.5 環境別対応
- **環境変数**: `APP_ENV`環境変数で環境を切り替え（develop/staging/production）
- **設定ファイル**: `config/{env}/`ディレクトリを環境別に管理
- **デフォルト環境**: 環境変数が未設定の場合は`develop`環境を使用

## 5. 制約事項

### 5.1 既存システムとの関係
- **既存のDocker環境**: PostgreSQL、Redis等の既存コンテナと統合
- **既存の設定ファイル**: `config/{env}/`ディレクトリの構造を維持
- **既存のデータベース**: SQLiteデータベースファイルの場所を維持（`server/data/`）
- **既存の起動方法**: Docker化後も既存の起動方法（`go run`、`npm run dev`）は維持可能

### 5.2 技術スタック
- **Go**: 既存のGoバージョン（1.21+）を維持
- **Node.js**: 既存のNode.jsバージョン（22+）を維持
- **Docker**: Docker Composeを使用
- **ベースイメージ**: Alpine Linuxベースの軽量イメージを使用

### 5.3 デプロイメント先
- **開発環境**: Docker Composeでローカル実行
- **本番環境**: AWS（ECS/EKS）またはTencent Cloud（TKE）へのデプロイを想定
- **コンテナレジストリ**: Docker Hub、AWS ECR、Tencent Cloud TCR等に対応

### 5.4 運用上の制約
- **同時起動**: 3つのサーバーを同時に起動可能
- **リソース制約**: 開発環境ではメモリ・CPU制限を考慮
- **ネットワーク**: Dockerネットワーク内での通信を前提

## 6. 受け入れ基準

### 6.1 APIサーバーのDocker化
- [ ] `server/Dockerfile`が作成されている
- [ ] `server/.dockerignore`が作成されている
- [ ] Dockerイメージが正常にビルドされる
- [ ] Dockerコンテナがポート8080で正常に起動する
- [ ] 環境変数`APP_ENV`で環境を切り替えられる
- [ ] 設定ファイルがボリュームマウントで読み込まれる
- [ ] データベースファイルがボリュームマウントで永続化される
- [ ] PostgreSQLコンテナと通信できる
- [ ] Redisコンテナと通信できる

### 6.2 AdminサーバーのDocker化
- [ ] `server/Dockerfile.admin`が作成されている
- [ ] Dockerイメージが正常にビルドされる
- [ ] Dockerコンテナがポート8081で正常に起動する
- [ ] 環境変数`APP_ENV`で環境を切り替えられる
- [ ] 設定ファイルがボリュームマウントで読み込まれる
- [ ] データベースファイルがボリュームマウントで永続化される
- [ ] PostgreSQLコンテナと通信できる
- [ ] APIサーバーコンテナと通信できる（必要に応じて）

### 6.3 クライアントサーバーのDocker化
- [ ] `client/Dockerfile`が作成されている
- [ ] `client/.dockerignore`が作成されている
- [ ] Dockerイメージが正常にビルドされる
- [ ] Dockerコンテナがポート3000で正常に起動する
- [ ] 環境変数`NEXT_PUBLIC_API_URL`でAPIサーバーのURLを設定できる
- [ ] 開発モードでホットリロードが機能する
- [ ] 本番モードで最適化されたビルドが実行される
- [ ] APIサーバーコンテナと通信できる

### 6.4 Docker Compose統合
- [ ] `docker-compose.api.develop.yml`が作成されている
- [ ] `docker-compose.api.staging.yml`が作成されている
- [ ] `docker-compose.api.production.yml`が作成されている
- [ ] `docker-compose.client.develop.yml`が作成されている
- [ ] `docker-compose.client.staging.yml`が作成されている
- [ ] `docker-compose.client.production.yml`が作成されている
- [ ] `docker-compose.admin.develop.yml`が作成されている
- [ ] `docker-compose.admin.staging.yml`が作成されている
- [ ] `docker-compose.admin.production.yml`が作成されている
- [ ] 各docker-composeファイルでサービスが正しく定義されている
- [ ] 既存のPostgreSQLコンテナと統合されている（外部ネットワーク参照）
- [ ] 既存のRedisコンテナと統合されている（外部ネットワーク参照）
- [ ] 環境変数が各ファイルで適切に設定されている
- [ ] ボリュームマウントが正しく設定されている
- [ ] 各サービスを個別に起動・停止できる

### 6.5 Dockerイメージの最適化
- [ ] マルチステージビルドが実装されている
- [ ] 最終イメージサイズが最小化されている
- [ ] .dockerignoreファイルが適切に設定されている
- [ ] ビルドキャッシュが活用されている
- [ ] 非rootユーザーで実行される

### 6.6 デプロイメント準備
- [ ] Dockerイメージにタグを付与できる
- [ ] コンテナレジストリにプッシュできる
- [ ] 本番環境用の環境変数が設定できる
- [ ] ヘルスチェックが実装されている
- [ ] リソース制限が設定できる

### 6.7 ドキュメント
- [ ] `docs/Docker.md`が作成されている
- [ ] Docker環境の構築手順が記載されている
- [ ] ビルド・起動・停止のコマンドが記載されている
- [ ] 環境別の起動方法が記載されている
- [ ] トラブルシューティング情報が記載されている
- [ ] 本番環境へのデプロイ手順が記載されている
- [ ] `README.md`にDocker化に関する情報が追記されている

## 7. 影響範囲

### 7.1 新規追加が必要なディレクトリ・ファイル

#### ディレクトリ
- なし（既存のディレクトリ構造を維持）

#### ファイル
- `server/Dockerfile`: APIサーバー用Dockerfile
- `server/Dockerfile.admin`: GoAdminサーバー用Dockerfile
- `server/.dockerignore`: APIサーバー用.dockerignore
- `client/Dockerfile`: クライアントサーバー用Dockerfile
- `client/.dockerignore`: クライアントサーバー用.dockerignore
- `docker-compose.api.develop.yml`: APIサーバー用Docker Compose設定（開発環境）
- `docker-compose.api.staging.yml`: APIサーバー用Docker Compose設定（ステージング環境）
- `docker-compose.api.production.yml`: APIサーバー用Docker Compose設定（本番環境）
- `docker-compose.client.develop.yml`: クライアントサーバー用Docker Compose設定（開発環境）
- `docker-compose.client.staging.yml`: クライアントサーバー用Docker Compose設定（ステージング環境）
- `docker-compose.client.production.yml`: クライアントサーバー用Docker Compose設定（本番環境）
- `docker-compose.admin.develop.yml`: GoAdminサーバー用Docker Compose設定（開発環境）
- `docker-compose.admin.staging.yml`: GoAdminサーバー用Docker Compose設定（ステージング環境）
- `docker-compose.admin.production.yml`: GoAdminサーバー用Docker Compose設定（本番環境）
- `docs/Docker.md`: Docker化に関するドキュメント

### 7.2 変更が必要なファイル

#### 設定ファイル
- なし（既存の設定ファイル構造を維持）

#### ドキュメント
- `README.md`: Docker化に関する情報を追記

### 7.3 既存ファイルの扱い
- `server/cmd/server/main.go`: 変更なし（Dockerfileから参照）
- `server/cmd/admin/main.go`: 変更なし（Dockerfile.adminから参照）
- `client/`: 変更なし（Dockerfileから参照）
- `config/{env}/`: 変更なし（ボリュームマウントで使用）
- `server/data/`: 変更なし（ボリュームマウントで使用）
- 既存のdocker-composeファイル: 変更なし（統合して使用）

## 8. 実装上の注意事項

### 8.1 Dockerfileの作成
- **マルチステージビルド**: ビルドステージと実行ステージを分離
- **ベースイメージ**: Alpine Linuxベースの軽量イメージを使用
- **CGO_ENABLED**: Goアプリケーションは`CGO_ENABLED=0`で静的リンクビルド
- **非rootユーザー**: セキュリティのため非rootユーザーで実行
- **作業ディレクトリ**: `/app`を作業ディレクトリとして使用

### 8.2 Docker Compose設定
- **ファイル構成**: サービス別・環境別にdocker-composeファイルを分離
- **命名規則**: `docker-compose.{service}.{env}.yml`
- **サービス名**: 各ファイルで`api`, `admin`, `client`を使用
- **ネットワーク**: 既存のdocker-composeネットワークを外部ネットワークとして参照
- **ボリューム**: 設定ファイルとデータディレクトリをマウント
- **環境変数**: 各ファイルで環境に応じた環境変数を設定
- **依存関係**: `depends_on`は使用せず、起動順序は運用で管理
- **外部ネットワーク**: 既存のPostgreSQL、Redisコンテナと通信するため、外部ネットワークを参照

### 8.3 環境変数の管理
- **APP_ENV**: 環境を指定（develop/staging/production）
- **NEXT_PUBLIC_API_URL**: クライアントからAPIサーバーへのURL
- **デフォルト値**: 環境変数が未設定の場合は適切なデフォルト値を設定

### 8.4 データ永続化
- **設定ファイル**: `config/{env}/`ディレクトリをボリュームマウント
- **データベースファイル**: `server/data/`ディレクトリをボリュームマウント
- **ログファイル**: 必要に応じてログディレクトリをボリュームマウント

### 8.5 既存サービスとの統合
- **PostgreSQL**: 既存の`docker-compose.postgres.yml`と統合
- **Redis**: 既存の`docker-compose.redis.yml`と統合
- **ネットワーク**: 同一Dockerネットワークで通信
- **サービス名**: 既存のサービス名（`postgres`, `redis`）を使用

### 8.6 ビルドと起動
- **ビルドコマンド**: `docker build`または`docker-compose -f docker-compose.{service}.{env}.yml build`
- **起動コマンド**: `docker-compose -f docker-compose.{service}.{env}.yml up`または`docker-compose -f docker-compose.{service}.{env}.yml up -d`
- **停止コマンド**: `docker-compose -f docker-compose.{service}.{env}.yml down`
- **ログ確認**: `docker-compose -f docker-compose.{service}.{env}.yml logs`または`docker logs`
- **複数サービスの起動**: 各サービスを個別に起動するか、起動スクリプトで一括起動
- **起動順序**: APIサーバー → Adminサーバー → クライアントサーバーの順で起動（運用で管理）

### 8.7 本番環境へのデプロイ準備
- **イメージタグ**: バージョン管理のため適切なタグを付与
- **コンテナレジストリ**: AWS ECR、Tencent Cloud TCR等にプッシュ
- **環境変数**: 本番環境用の環境変数を設定
- **ヘルスチェック**: Docker healthcheckを実装
- **リソース制限**: CPU・メモリ制限を設定

### 8.8 ドキュメント整備
- **Docker.md**: Docker環境の構築手順、ビルド・起動・停止コマンド、トラブルシューティング
- **README.md**: Docker化に関する簡単な説明と起動方法
- **環境別の起動方法**: `APP_ENV`環境変数の使用方法
- **本番環境へのデプロイ手順**: AWS/Tencent Cloudへのデプロイ方法

## 9. 参考情報

### 9.1 関連Issue
- GitHub Issue #83: APIサーバー、クライアントサーバー、GoAdminサーバーをDocker上で動作させる

### 9.2 既存ドキュメント
- `README.md`: プロジェクト概要とセットアップ手順
- `docs/Architecture.md`: システムアーキテクチャ
- `docs/Initial-Setup.md`: 初期セットアップ手順
- `config/{env}/`: 環境別設定ファイル

### 9.3 技術スタック
- **Go**: 1.21+
- **Node.js**: 22+
- **Docker**: Docker Composeを使用
- **ベースイメージ**: Alpine Linux

### 9.4 参考リンク
- Docker公式ドキュメント: https://docs.docker.com/
- Docker Compose公式ドキュメント: https://docs.docker.com/compose/
- Go公式Dockerイメージ: https://hub.docker.com/_/golang
- Node.js公式Dockerイメージ: https://hub.docker.com/_/node
- Alpine Linux公式サイト: https://alpinelinux.org/
