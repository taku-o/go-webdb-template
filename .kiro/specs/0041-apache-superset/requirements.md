# Apache Superset導入要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #81
- **Issueタイトル**: Apache Supersetの導入
- **Feature名**: 0041-apache-superset
- **作成日**: 2025-01-27

### 1.2 目的
非エンジニア向けのデータビューワとしてApache Supersetを導入する。
Apache Supersetの使い勝手を確認し、データ可視化の基盤を構築する。

### 1.3 スコープ
- Apache SupersetをDockerで導入
- PostgreSQLデータベースに接続してデータを可視化
- Apache Supersetの起動スクリプトの実装
- Docker Compose設定の実装

**本実装の範囲外**:
- 本番環境での本格的なApache Superset運用
- 複数のデータソース接続（PostgreSQL以外）
- ユーザー認証・認可の詳細設定
- ダッシュボードの事前作成
- データ更新の自動化
- パフォーマンス最適化

## 2. 背景・現状分析

### 2.1 現在の実装
- **データビューワ**: Metabaseが既に導入されている（`docker-compose.metabase.yml`、`scripts/metabase-start.sh`）
- **PostgreSQL**: PostgreSQLが既に導入されている（`docker-compose.postgres.yml`、`scripts/start-postgres.sh`）
  - データベース名: `webdb`
  - ユーザー名: `webdb`
  - パスワード: `webdb`
  - ポート: `5432`
- **Docker**: Docker Composeファイルが存在（mailpit、metabase、cloudbeaver、postgres、redis等）
- **起動スクリプト**: `scripts/`ディレクトリに各種サービスの起動スクリプトが存在
  - `scripts/start-postgres.sh`: PostgreSQL起動スクリプト
  - `scripts/metabase-start.sh`: Metabase起動スクリプト
  - `scripts/cloudbeaver-start.sh`: CloudBeaver起動スクリプト

### 2.2 課題点
1. **非エンジニア向けデータビューワの不足**: 既存のMetabaseはあるが、Apache Supersetの使い勝手を確認したい
2. **データ可視化ツールの選択肢不足**: 複数のデータビューワを比較検討するための環境が不足
3. **Apache Supersetの評価環境不足**: Apache Supersetを試用するための環境が構築されていない

### 2.3 本実装による改善点
1. **非エンジニア向けデータビューワの提供**: Apache Supersetを導入し、非エンジニアでも使いやすいデータ可視化環境を提供
2. **データビューワの比較検討**: MetabaseとApache Supersetを比較検討できる環境を構築
3. **Apache Supersetの評価環境構築**: Apache Supersetの使い勝手を確認できる環境を構築

## 3. 機能要件

### 3.1 Apache Superset環境の構築

#### 3.1.1 Apache Supersetの導入
- Apache SupersetをDockerコンテナとして導入
- デフォルトポート: 8088（Web UI）
- データ永続化の設定（必須）
  - Dockerボリュームを使用してデータを永続化
  - コンテナ再起動時もデータが保持されること
  - 設定ファイル、データベース、アップロードファイル等を永続化
- 本番環境・staging環境でも使用される（Docker Composeと起動スクリプトは開発用途）

#### 3.1.2 PostgreSQL接続設定
- 既存のPostgreSQLデータベース（`docker-compose.postgres.yml`で起動）に接続
- 接続情報:
  - ホスト: `postgres`（Dockerネットワーク内）または`localhost`（ホストマシンから）
  - ポート: `5432`
  - データベース名: `webdb`
  - ユーザー名: `webdb`
  - パスワード: `webdb`
- Apache SupersetからPostgreSQLのデータを閲覧・可視化できること

#### 3.1.3 Docker Compose設定
- `docker-compose.apache-superset.yml`ファイルを作成（開発用途）
  - Apache Supersetサービスの定義
  - ネットワーク設定（既存のPostgreSQLと接続可能にする）
  - ボリューム設定（データ永続化用、必須）
    - データディレクトリをマウント（コンテナ内の`/app/superset_home`にマウント）
    - データはプロジェクトルートの`apache-superset/data`ディレクトリに保存される
    - ホストマシンから直接アクセス可能
  - 環境変数設定（データベース接続、セキュリティ設定等）
- Docker Composeと起動スクリプトは開発用途（本番・staging環境では別の方法でApache Supersetを起動）
- Apache Superset自体は本番環境・staging環境でも使用される

### 3.2 起動スクリプトの実装（開発用途）

#### 3.2.1 Apache Superset起動スクリプト
- `scripts/start-apache-superset.sh`を作成（開発用途）
- Docker Composeを使用してApache Supersetを起動
- 既存の起動スクリプト（`metabase-start.sh`など）と同じパターンで実装
- 本番・staging環境では別の方法でApache Supersetを起動

### 3.3 初期設定

#### 3.3.1 デフォルト管理者アカウント
- デフォルトの管理者アカウントを作成
  - ユーザー名: `admin`
  - パスワード: `admin`（開発環境用、本番環境では変更必須）
- 初回起動時に自動的に管理者アカウントが作成されること

#### 3.3.2 PostgreSQLデータソース接続
- Apache Superset起動後、PostgreSQLデータソースを手動で接続設定できること
- 接続情報は環境変数または設定ファイルで管理可能（将来の拡張項目）

## 4. 非機能要件

### 4.1 パフォーマンス
- Apache Supersetの起動は適切な時間内に完了すること
- データクエリの実行は適切な時間内に完了すること

### 4.2 セキュリティ
- 開発環境では簡易的な認証設定で可
- 本番環境では適切な認証・認可を実装すること（将来の拡張項目）
- デフォルトパスワードは開発環境のみで使用し、本番環境では変更必須

### 4.3 可用性
- PostgreSQLが停止した場合のエラーハンドリングを実装
- Apache Supersetが起動できない場合のエラーメッセージを表示

### 4.4 保守性
- 既存のDocker Compose設定パターンに従うこと
- 既存の起動スクリプトパターンに従うこと
- 既存のMetabaseやCloudBeaverと同じディレクトリ構造に従うこと

### 4.5 拡張性
- 将来的に複数のデータソースを追加できるように設計すること
- 設定ファイルから接続情報を読み込めるようにすること（将来の拡張項目）

## 5. 技術仕様

### 5.1 インフラストラクチャ
- **Apache Superset**: Docker Composeで起動
- **PostgreSQL**: 既存の`docker-compose.postgres.yml`で起動
- **Docker Compose**: 既存のパターンに従う

### 5.2 ファイル構造
- **Docker Compose設定**: `docker-compose.apache-superset.yml`（Apache Superset用）
- **起動スクリプト**: `scripts/start-apache-superset.sh`
- **データディレクトリ**: `apache-superset/data/`（プロジェクトルート）

### 5.3 Apache Superset設定
- **イメージ**: `apache/superset:latest`（または安定版）
- **ポート**: `8088`（Web UI）
- **データ永続化**: Dockerボリュームを使用
- **環境変数**: 
  - `SUPERSET_SECRET_KEY`: ランダムなシークレットキー（開発環境用）
  - `SUPERSET_CONFIG_PATH`: 設定ファイルパス（オプション）

## 6. 受け入れ基準

### 6.1 機能要件
1. **Apache Superset環境の構築**: Apache SupersetがDocker Composeで起動できること
2. **データ永続化**: Apache SupersetのデータがDockerボリュームに永続化され、コンテナ再起動後もデータが保持されること
3. **起動スクリプト**: Apache Supersetを起動するスクリプトが動作すること
4. **PostgreSQL接続**: Apache Supersetから既存のPostgreSQLデータベースに接続できること
5. **データ可視化**: Apache SupersetのWeb UIからPostgreSQLのデータを閲覧・可視化できること
6. **管理者アカウント**: デフォルトの管理者アカウント（admin/admin）でログインできること

### 6.2 非機能要件
1. **パフォーマンス**: Apache Supersetが適切な時間内に起動すること
2. **可用性**: PostgreSQLが停止した場合のエラーハンドリングが実装されていること
3. **保守性**: 既存のパターンに従った実装であること
4. **拡張性**: 将来的に複数のデータソースを追加できる設計であること

## 7. 制約事項

1. **開発環境での利用**: 本実装は主に開発環境での利用を想定（本番環境での利用は将来の拡張項目）
2. **シンプルな実装**: 最小限の機能のみを実装し、複雑な機能は将来の拡張項目とする
3. **既存PostgreSQLの利用**: 既存のPostgreSQLデータベース（`docker-compose.postgres.yml`）を使用する
4. **参考実装**: 本実装はApache Supersetの使い勝手を確認するための参考実装として利用する

## 8. 将来の拡張項目（現時点では未実装）

以下の機能は将来の拡張として検討されていますが、現時点では実装対象外です：

- 複数のデータソース接続（PostgreSQL以外）
- ユーザー認証・認可の詳細設定
- ダッシュボードの事前作成
- データ更新の自動化
- パフォーマンス最適化
- 設定ファイルからの接続情報読み込み
- 本番環境での本格的なApache Superset運用
- セキュリティ強化（本番環境用）

## Project Description (Input)

非エンジニア向けのデータビューワとしてApache Supersetを導入する。
Apache Supersetの使い勝手を確認したい。

Apache SupersetはDockerで導入する。
データベースはPostgreSQLを起動して、そのデータを見るものとする。

Apache Supersetの起動スクリプトも欲しい。

## Requirements

### Requirement 1: Apache Superset環境の構築
**Objective:** As a developer, I want to set up Apache Superset as a data visualization tool, so that non-engineers can view and analyze database data easily.

#### Acceptance Criteria
1. WHEN Apache Superset is started THEN it SHALL be available as a Docker container
2. IF Apache Superset is started THEN it SHALL listen on port 8088 for Web UI
3. WHERE Apache Superset is configured THEN it SHALL use Docker Compose for management
4. WHEN Apache Superset is configured THEN it SHALL use Docker volumes for data persistence (required)
5. IF Apache Superset data is persisted THEN it SHALL survive container restarts
6. WHERE Apache Superset persistence is configured THEN it SHALL store data in apache-superset/data directory

### Requirement 2: PostgreSQL接続設定
**Objective:** As a user, I want Apache Superset to connect to the existing PostgreSQL database, so that I can visualize PostgreSQL data.

#### Acceptance Criteria
1. WHEN Apache Superset is started THEN it SHALL be able to connect to PostgreSQL database
2. IF Apache Superset connects to PostgreSQL THEN it SHALL use the existing PostgreSQL instance (docker-compose.postgres.yml)
3. WHERE PostgreSQL connection is configured THEN it SHALL use connection info: host=postgres/localhost, port=5432, database=webdb, user=webdb, password=webdb
4. WHEN PostgreSQL connection is established THEN Apache Superset SHALL be able to browse and visualize PostgreSQL data
5. IF PostgreSQL is not available THEN Apache Superset SHALL handle the error gracefully

### Requirement 3: Docker Compose設定の実装
**Objective:** As a system, I want Apache Superset to be managed via Docker Compose, so that it can be easily started and stopped.

#### Acceptance Criteria
1. WHEN docker-compose.apache-superset.yml is created THEN it SHALL define Apache Superset service
2. IF docker-compose.apache-superset.yml is created THEN it SHALL use appropriate Docker image (apache/superset:latest or stable version)
3. WHERE Apache Superset service is defined THEN it SHALL use Docker volumes for data persistence (required)
4. WHEN Apache Superset volumes are defined THEN they SHALL persist data across container restarts
5. IF Apache Superset is configured THEN it SHALL connect to existing PostgreSQL via Docker network
6. WHERE Apache Superset network is configured THEN it SHALL be able to access PostgreSQL service

### Requirement 4: 起動スクリプトの実装
**Objective:** As a developer, I want to start Apache Superset easily, so that I can quickly set up the development environment.

#### Acceptance Criteria
1. WHEN start-apache-superset.sh is executed THEN it SHALL start Apache Superset using Docker Compose
2. IF start-apache-superset.sh is executed THEN it SHALL provide appropriate feedback
3. WHERE script is created THEN it SHALL follow existing script patterns (metabase-start.sh, cloudbeaver-start.sh)
4. WHEN script is executed THEN it SHALL use docker-compose.apache-superset.yml file

### Requirement 5: 初期設定の実装
**Objective:** As a user, I want to access Apache Superset with default admin account, so that I can start using it immediately after setup.

#### Acceptance Criteria
1. WHEN Apache Superset is started for the first time THEN it SHALL create default admin account automatically
2. IF default admin account is created THEN it SHALL use username=admin, password=admin (development only)
3. WHERE admin account is created THEN it SHALL be usable for login
4. WHEN PostgreSQL data source is configured THEN it SHALL be configurable manually via Web UI
5. IF data source connection fails THEN Apache Superset SHALL display appropriate error message

### Requirement 6: データ可視化機能
**Objective:** As a non-engineer user, I want to view and visualize PostgreSQL data in Apache Superset, so that I can analyze data without SQL knowledge.

#### Acceptance Criteria
1. WHEN PostgreSQL data source is connected THEN Apache Superset SHALL be able to browse database tables
2. IF tables are browsed THEN Apache Superset SHALL display table structure and data
3. WHERE data is visualized THEN Apache Superset SHALL provide chart creation capabilities
4. WHEN charts are created THEN Apache Superset SHALL allow saving and sharing dashboards
5. IF data query fails THEN Apache Superset SHALL display appropriate error message

### Requirement 7: エラーハンドリングの実装
**Objective:** As a system, I want to handle errors gracefully, so that failures do not crash the application.

#### Acceptance Criteria
1. WHEN PostgreSQL connection fails THEN the system SHALL return appropriate error message
2. IF Apache Superset startup fails THEN the system SHALL display appropriate error message
3. WHERE data query fails THEN Apache Superset SHALL log the error and display user-friendly message
4. WHEN Docker container fails THEN the system SHALL handle the error gracefully
