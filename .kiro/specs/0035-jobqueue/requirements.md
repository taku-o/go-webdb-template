# ジョブキュー機能要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #70
- **Issueタイトル**: キューを入れてバッググラウンドで処理を行う機能を用意する
- **Feature名**: 0035-jobqueue
- **作成日**: 2025-01-27

### 1.2 目的
キューにジョブを登録し、バックグラウンドで処理を行う機能を実装する。
これにより、非同期処理の基盤を構築し、将来的な拡張に対応できるようにする。

### 1.3 スコープ
- RedisをDockerで導入（キューのDBとして）
- Redis Insightを導入（データビューワとして）
- Asynqライブラリを使用したジョブキューシステムの実装
- クライアント側にジョブ登録ボタンの実装
- APIサーバー側でRedisにキューを登録する機能
- 3分後に標準出力に文字列を出力するジョブ処理の実装
- RedisとRedis Insightを起動する処理の実装

**本実装の範囲外**:
- 本番環境での本格的なジョブ処理機能
- ジョブの優先度設定機能
- ジョブの再試行機能（Asynqのデフォルト機能は使用可能）
- ジョブの履歴管理機能
- ジョブの監視・アラート機能
- 複数のジョブタイプの実装

## 2. 背景・現状分析

### 2.1 現在の実装
- **バックグラウンド処理**: 現在、バックグラウンド処理機能は実装されていない
- **キューシステム**: キューシステムは実装されていない
- **Redis**: Redisは現在導入されていない
- **Docker**: Docker Composeファイルが存在（mailpit、metabase、cloudbeaver用）
  - `docker-compose.mailpit.yml`: Mailpit用のDocker Compose設定
  - `docker-compose.metabase.yml`: Metabase用のDocker Compose設定
  - `docker-compose.cloudbeaver.yml`: CloudBeaver用のDocker Compose設定
- **起動スクリプト**: `scripts/`ディレクトリに各種サービスの起動スクリプトが存在
  - `scripts/start-mailpit.sh`: Mailpit起動スクリプト
  - `scripts/metabase-start.sh`: Metabase起動スクリプト
  - `scripts/cloudbeaver-start.sh`: CloudBeaver起動スクリプト

### 2.2 課題点
1. **非同期処理の不足**: 現在、すべての処理が同期的に実行されるため、時間のかかる処理がAPIレスポンスをブロックする
2. **バックグラウンド処理の基盤不足**: バックグラウンド処理を行うための基盤が存在しない
3. **ジョブ管理機能の不足**: ジョブを管理・監視する機能が存在しない
4. **キューシステムの不足**: ジョブをキューに登録して処理する機能が存在しない

### 2.3 本実装による改善点
1. **非同期処理の実現**: ジョブをキューに登録し、バックグラウンドで処理できるようになる
2. **バックグラウンド処理の基盤構築**: 将来的な拡張に対応できる基盤を構築
3. **ジョブ管理機能の提供**: Redis Insightを使用してジョブの状態を可視化できる
4. **キューシステムの実装**: Asynqを使用した堅牢なキューシステムを実装

## 3. 機能要件

### 3.1 Redis環境の構築

#### 3.1.1 Redisの導入
- RedisをDockerコンテナとして導入（1台のRedisサーバー）
- デフォルトポート: 6379
- データ永続化の設定（必須）
  - RDB（Redis Database Backup）またはAOF（Append Only File）を使用
  - Dockerボリュームを使用してデータを永続化
  - コンテナ再起動時もデータが保持されること
- 本番環境・staging環境でも使用される（Docker Composeと起動スクリプトは開発用途）

#### 3.1.2 Redis Insightの導入
- Redis InsightをDockerコンテナとして導入
- デフォルトポート: 8001（Web UI）
- `docker-compose.redis.yml`で起動した1台のRedisサーバーと接続
- データビューワとしての機能を提供
- 本番環境・staging環境でも使用される（Docker Composeと起動スクリプトは開発用途）

#### 3.1.3 Docker Compose設定
- `docker-compose.redis.yml`ファイルを作成（開発用途）
  - 1台のRedisサービスの定義
  - ネットワーク設定
  - ボリューム設定（データ永続化用、必須）
    - データディレクトリをマウント（コンテナ内の`/data`にマウント）
    - データはプロジェクトルートの`redis/data/jobqueue`ディレクトリに保存される
    - ホストマシンから直接アクセス可能
    - 永続化設定（RDB/AOF）を有効化
- `docker-compose.redis-insight.yml`ファイルを作成（開発用途）
  - Redis Insightサービスの定義
  - ネットワーク設定
  - `docker-compose.redis.yml`で起動した1台のRedisサーバーと接続する設定
- Docker Composeと起動スクリプトは開発用途（本番・staging環境では別の方法でRedisとRedis Insightを起動）
- RedisとRedis Insight自体は本番環境・staging環境でも使用される

### 3.2 起動スクリプトの実装（開発用途）

#### 3.2.1 Redis起動スクリプト
- `scripts/start-redis.sh`を作成（開発用途）
- Docker Composeを使用してRedisを起動
- 既存の起動スクリプト（`start-mailpit.sh`など）と同じパターンで実装
- 本番・staging環境では別の方法でRedisを起動

#### 3.2.2 Redis Insight起動スクリプト
- `scripts/start-redis-insight.sh`を作成（開発用途）
- Docker Composeを使用してRedis Insightを起動
- 既存の起動スクリプトと同じパターンで実装
- 本番・staging環境では別の方法でRedis Insightを起動

### 3.3 ジョブキューシステムの実装

#### 3.3.1 Asynqライブラリの導入
- `github.com/hibiken/asynq`ライブラリを使用
- クライアントとサーバーの実装
- ジョブの登録と処理の実装

#### 3.3.2 ジョブ処理の実装
- 指定された遅延時間後に標準出力に文字列を出力するジョブ処理を実装
- 遅延時間はジョブごとに設定可能
- デフォルトの遅延時間（3分 = 180秒）を定数で定義
- 大半のジョブはデフォルト値（3分）を使用
- 必要に応じてジョブごとに異なる遅延時間を設定可能
- ジョブ処理のエラーハンドリング

#### 3.3.3 ジョブ登録APIの実装
- APIエンドポイントを実装（参考コード用の名前を使用）
- クライアントからのリクエストを受け取り、ジョブをキューに登録
- ジョブ登録の成功・失敗をレスポンスで返す

### 3.4 クライアント側の実装

#### 3.4.1 ジョブ登録ボタンの実装
- クライアント側にボタンを追加
- ボタンを押したら、APIサーバーにリクエストを送信
- ジョブ登録の成功・失敗をユーザーに通知

#### 3.4.2 UIの実装
- 参考コードとして利用するため、邪魔にならない名前を使用
- 既存のUIパターンに従った実装

### 3.5 設定管理

#### 3.5.1 Redis接続設定
- `config/{env}/cacheserver.yaml`からRedis接続情報を読み取る
- ジョブキュー用とデフォルト（rate limit等）用のRedis設定を分離
  - ジョブキュー用: `redis.jobqueue.addr`（単一Redis接続、1台）
  - デフォルト用: `redis.default.cluster.addrs`（複数台対応可能、rate limit等で使用）
- デフォルト値: `localhost:6379`（`cacheserver.yaml`に設定がない場合）
- ジョブキュー用は単一Redis接続（1台）を想定
- デフォルト用は複数台対応（Cluster設定、rate limit等の他の処理でも使用）

#### 3.5.2 ジョブ処理設定
- デフォルトの遅延時間を定数で定義（3分 = 180秒）
- ジョブごとに遅延時間を設定可能
- ジョブ登録時に遅延時間を指定できる（指定がない場合はデフォルト値を使用）
- 設定ファイルからの読み込みは将来の拡張項目

## 4. 非機能要件

### 4.1 パフォーマンス
- ジョブの登録は非同期で処理され、APIレスポンスをブロックしないこと
- Redisへの接続は効率的に行うこと

### 4.2 セキュリティ
- Redisへの接続は適切に保護すること（開発環境では簡易的な設定で可）
- 本番環境では適切な認証・認可を実装すること（将来の拡張項目）

### 4.3 可用性
- Redisが停止した場合のエラーハンドリングを実装
- ジョブ登録の失敗時は適切なエラーメッセージを返すこと

### 4.4 保守性
- 既存のDocker Compose設定パターンに従うこと
- 既存の起動スクリプトパターンに従うこと
- 参考コードとして利用するため、邪魔にならない名前を使用すること

### 4.5 拡張性
- 将来的に複数のジョブタイプを追加できるように設計すること
- ジョブ処理の設定を柔軟に変更できるようにすること

## 5. 技術仕様

### 5.1 サーバー側技術スタック
- **言語**: Go 1.21+
- **キューライブラリ**: `github.com/hibiken/asynq`
- **Redisクライアント**: Asynq内蔵のRedisクライアント
- **Webフレームワーク**: 既存のEcho + Huma API

### 5.2 インフラストラクチャ
- **Redis**: Docker Composeで起動
- **Redis Insight**: Docker Composeで起動
- **Docker Compose**: 既存のパターンに従う

### 5.3 クライアント側技術スタック
- **フレームワーク**: 既存のクライアントフレームワーク
- **HTTPクライアント**: 既存のHTTPクライアントライブラリ

### 5.4 ファイル構造
- **Docker Compose設定**: 
  - `docker-compose.redis.yml`（Redis用）
  - `docker-compose.redis-insight.yml`（Redis Insight用）
- **起動スクリプト**: `scripts/start-redis.sh`, `scripts/start-redis-insight.sh`
- **ジョブ処理実装**: `server/internal/service/jobqueue/`（新規作成）
- **APIハンドラー**: `server/internal/api/handler/dm_jobqueue_handler.go`（新規作成）
- **設定**: `server/internal/config/config.go`にRedis設定を追加

## 6. 受け入れ基準

### 6.1 機能要件
1. **Redis環境の構築**: RedisとRedis InsightがDocker Composeで起動できること
2. **データ永続化**: RedisのデータがDockerボリュームに永続化され、コンテナ再起動後もデータが保持されること
3. **起動スクリプト**: RedisとRedis Insightを起動するスクリプトが動作すること
4. **ジョブ登録**: クライアントからボタンを押してジョブを登録できること
5. **ジョブ処理**: 指定された遅延時間（デフォルト: 3分）後に標準出力に文字列が出力されること
6. **ジョブごとの遅延時間設定**: ジョブごとに異なる遅延時間を設定できること
7. **エラーハンドリング**: Redisが停止している場合、適切なエラーメッセージが返されること

### 6.2 非機能要件
1. **パフォーマンス**: ジョブ登録がAPIレスポンスをブロックしないこと
2. **可用性**: Redisが停止した場合のエラーハンドリングが実装されていること
3. **保守性**: 既存のパターンに従った実装であること
4. **拡張性**: 将来的に複数のジョブタイプを追加できる設計であること

## 7. 制約事項

1. **参考コードとしての利用**: 本実装は参考コードとして利用するため、URLやエンドポイントURLなどは後で本実装する際、邪魔にならない名前にする
2. **遅延時間の設定**: 遅延時間はジョブごとに設定可能。デフォルト値（3分）はコードに定数で定義する（設定ファイルからの読み込みは将来の拡張項目）
3. **開発環境での利用**: 本実装は主に開発環境での利用を想定（本番環境での利用は将来の拡張項目）
4. **シンプルな実装**: 最小限の機能のみを実装し、複雑な機能は将来の拡張項目とする

## 8. 将来の拡張項目（現時点では未実装）

以下の機能は将来の拡張として検討されていますが、現時点では実装対象外です：

- 複数のジョブタイプの実装
- ジョブの優先度設定機能
- ジョブの再試行機能（Asynqのデフォルト機能は使用可能だが、カスタマイズは未実装）
- ジョブの履歴管理機能
- ジョブの監視・アラート機能
- 設定ファイルからの遅延時間の読み込み
- 本番環境での本格的なジョブ処理機能
- Redis認証・認可の実装

## Project Description (Input)

キューを入れてバッググラウンドで処理を行う機能を用意する。

### 機能
* いったんキューにジョブを登録し、しばらく後の処理する機能を用意する。
    * 処理の遅延時間はコードに定数で定義して良い。

### 環境
* キューのDBとしてRedisをDockerで導入。
* Redisのデータビューワとして、Redis Insightを導入。
    * Redisを起動する処理を用意する
    * Redis Insightを起動する処理を用意する
* 処理を実装するためのライブラリとしてAsynqを利用。

### 処理フロー
* クライアントにボタンを用意。ボタンを押したら、ジョブを登録する。
* APIサーバー側でRedisにキューを登録して
* 3分後に標準出力に文字列を出力する

### 設計
* 用意した実装は参考コードとして利用する。
    * URLやエンドポイントURLなどはあとで本実装する際、邪魔にならない名前にしたい。
* まずは設計を提案してください。

## Requirements

### Requirement 1: Redis環境の構築
**Objective:** As a developer, I want to set up Redis as a job queue database, so that I can store and manage background jobs.

#### Acceptance Criteria
1. WHEN Redis is started THEN it SHALL be available as a Docker container
2. IF Redis is started THEN it SHALL listen on port 6379 by default
3. WHERE Redis is configured THEN it SHALL use Docker Compose for management
4. WHEN Redis is configured THEN it SHALL use Docker volumes for data persistence (required)
5. IF Redis data is persisted THEN it SHALL survive container restarts
6. WHERE Redis persistence is configured THEN it SHALL use RDB or AOF persistence method

### Requirement 2: Redis Insight環境の構築
**Objective:** As a developer, I want to set up Redis Insight as a data viewer, so that I can monitor and inspect job queue data.

#### Acceptance Criteria
1. WHEN Redis Insight is started THEN it SHALL be available as a Docker container
2. IF Redis Insight is started THEN it SHALL listen on port 8001 for Web UI
3. WHERE Redis Insight is configured THEN it SHALL connect to the Redis server started by docker-compose.redis.yml
4. WHEN Redis Insight is started THEN it SHALL provide data viewing capabilities
5. IF Redis Insight is configured THEN it SHALL use a separate docker-compose file from Redis
6. WHERE Redis Insight connects to Redis THEN it SHALL use the same Redis server (1 instance) as docker-compose.redis.yml

### Requirement 3: Docker Compose設定の実装
**Objective:** As a system, I want Redis and Redis Insight to be managed via separate Docker Compose files, so that they can be easily started and stopped independently.

#### Acceptance Criteria
1. WHEN docker-compose.redis.yml is created THEN it SHALL define a single Redis service
2. IF docker-compose.redis-insight.yml is created THEN it SHALL define Redis Insight service only
3. WHERE services are defined THEN they SHALL use appropriate Docker images
4. WHEN Redis service is defined THEN it SHALL use Docker volumes for data persistence (required)
5. WHERE Redis volumes are defined THEN they SHALL persist data across container restarts
6. WHEN Redis persistence is configured THEN it SHALL enable RDB or AOF persistence method
7. IF Redis Insight is configured THEN it SHALL connect to the Redis server started by docker-compose.redis.yml
8. WHERE Redis and Redis Insight are configured THEN they SHALL be in separate docker-compose files
9. WHEN Redis Insight connects to Redis THEN it SHALL use external network to connect to docker-compose.redis.yml's network

### Requirement 4: 起動スクリプトの実装
**Objective:** As a developer, I want to start Redis and Redis Insight easily, so that I can quickly set up the development environment.

#### Acceptance Criteria
1. WHEN start-redis.sh is executed THEN it SHALL start Redis using Docker Compose
2. IF start-redis-insight.sh is executed THEN it SHALL start Redis Insight using Docker Compose
3. WHERE scripts are created THEN they SHALL follow existing script patterns
4. WHEN scripts are executed THEN they SHALL provide appropriate feedback

### Requirement 5: Asynqライブラリの導入とジョブキューシステムの実装
**Objective:** As a system, I want to use Asynq library for job queue management, so that I can reliably process background jobs.

#### Acceptance Criteria
1. WHEN Asynq is integrated THEN it SHALL use github.com/hibiken/asynq library
2. IF Asynq client is created THEN it SHALL connect to Redis
3. WHERE Asynq server is created THEN it SHALL process jobs from the queue
4. WHEN job is registered THEN it SHALL be stored in Redis queue
5. IF job processing fails THEN it SHALL handle errors appropriately

### Requirement 6: ジョブ処理の実装
**Objective:** As a system, I want to process jobs after a configurable delay, so that I can demonstrate background job processing with flexible timing.

#### Acceptance Criteria
1. WHEN job is processed THEN it SHALL wait for the specified delay time before execution
2. IF delay time is not specified THEN it SHALL use default delay time (3 minutes = 180 seconds)
3. WHERE default delay time is defined THEN it SHALL be defined as a constant in code
4. WHEN job is executed THEN it SHALL output a string to standard output
5. IF job has custom delay time THEN it SHALL use the specified delay time
6. WHEN job processing completes THEN it SHALL log the completion status

### Requirement 7: ジョブ登録APIの実装
**Objective:** As a client application, I want to register jobs via API with optional delay time, so that I can trigger background job processing with custom timing.

#### Acceptance Criteria
1. WHEN API endpoint is called THEN it SHALL register a job to Redis queue
2. IF delay time is specified in request THEN it SHALL use the specified delay time
3. WHERE delay time is not specified THEN it SHALL use default delay time (3 minutes)
4. IF job registration succeeds THEN it SHALL return success response
5. WHERE job registration fails THEN it SHALL return appropriate error response
6. WHEN API endpoint is implemented THEN it SHALL use reference code naming (not interfere with future implementation)
7. IF Redis is unavailable THEN it SHALL return appropriate error message

### Requirement 8: クライアント側ジョブ登録ボタンの実装
**Objective:** As a user, I want to trigger job registration via a button, so that I can test background job processing.

#### Acceptance Criteria
1. WHEN button is clicked THEN it SHALL send request to API server
2. IF job registration succeeds THEN it SHALL notify user of success
3. WHERE job registration fails THEN it SHALL notify user of error
4. WHEN button is implemented THEN it SHALL use reference code naming (not interfere with future implementation)

### Requirement 9: Redis接続設定の実装
**Objective:** As a system administrator, I want to configure Redis connection settings separately for job queue and rate limit, so that the system can connect to different Redis environments appropriately.

#### Acceptance Criteria
1. WHEN Redis connection is configured THEN it SHALL use config/{env}/cacheserver.yaml
2. IF job queue Redis connection is configured THEN it SHALL use redis.jobqueue.addr (single connection, 1 instance)
3. IF default Redis connection is configured THEN it SHALL use redis.default.cluster.addrs (multiple instances support, used by rate limit and other processes)
4. WHERE Redis connection is not configured THEN default SHALL be localhost:6379
5. WHEN Redis connection fails THEN it SHALL handle errors appropriately
6. WHERE job queue and default Redis are configured THEN they SHALL use separate Redis settings

### Requirement 10: エラーハンドリングの実装
**Objective:** As a system, I want to handle errors gracefully, so that failures do not crash the application.

#### Acceptance Criteria
1. WHEN Redis connection fails THEN the system SHALL return appropriate error message
2. IF job registration fails THEN the system SHALL return appropriate error response
3. WHERE job processing fails THEN the system SHALL log the error
4. WHEN Redis is unavailable THEN the system SHALL handle the error gracefully
