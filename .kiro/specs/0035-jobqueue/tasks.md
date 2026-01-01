# ジョブキュー機能実装タスク一覧

## 概要
ジョブキュー機能の実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: Docker Compose設定と起動スクリプト

#### - [x] タスク 1.1: docker-compose.redis.ymlの作成
**目的**: RedisをDocker Composeで起動するための設定ファイルを作成する

**作業内容**:
- `docker-compose.redis.yml`ファイルを新規作成
- Redis 7 Alpineイメージを使用
- ポート6379を公開
- データ永続化設定:
  - ボリュームマウント: `./redis/data/jobqueue:/data`
  - AOF（Append Only File）を有効化: `--appendonly yes`
  - AOF同期設定: `--appendfsync everysec`
- ネットワーク`redis-network`を作成（外部ネットワークとして使用可能）
- `restart: unless-stopped`を設定

**受け入れ基準**:
- `docker-compose.redis.yml`ファイルが作成されている
- Redisコンテナが起動できること
- データが`redis/data/jobqueue`ディレクトリに保存されること
- コンテナ再起動後もデータが保持されること

---

#### - [x] タスク 1.2: docker-compose.redis-insight.ymlの作成
**目的**: Redis InsightをDocker Composeで起動するための設定ファイルを作成する

**作業内容**:
- `docker-compose.redis-insight.yml`ファイルを新規作成
- Redis Insight最新イメージを使用
- ポート8001を公開（Web UI）
- 環境変数`REDIS_HOSTS`を設定（起動スクリプトで設定される）
- ボリューム`redis-insight-data`で設定を永続化
- `docker-compose.redis.yml`で作成されたネットワーク`redis-network`を外部ネットワークとして使用
- `restart: unless-stopped`を設定

**受け入れ基準**:
- `docker-compose.redis-insight.yml`ファイルが作成されている
- Redis Insightコンテナが起動できること
- `docker-compose.redis.yml`のネットワークに接続できること
- Web UIが`http://localhost:8001`でアクセスできること

---

#### - [x] タスク 1.3: scripts/start-redis.shの作成
**目的**: Redisを起動・停止するスクリプトを作成する

**作業内容**:
- `scripts/start-redis.sh`ファイルを新規作成
- 実行権限を付与（`chmod +x`）
- 既存の`start-mailpit.sh`と同じパターンで実装
- `start`コマンド: Docker ComposeでRedisを起動
- `stop`コマンド: Docker ComposeでRedisを停止
- 適切なフィードバックメッセージを出力

**受け入れ基準**:
- `scripts/start-redis.sh`ファイルが作成されている
- 実行権限が付与されている
- `./scripts/start-redis.sh start`でRedisが起動できること
- `./scripts/start-redis.sh stop`でRedisが停止できること
- 既存の起動スクリプトと同じパターンに従っている

---

#### - [x] タスク 1.4: scripts/start-redis-insight.shの作成
**目的**: Redis Insightを起動・停止するスクリプトを作成する

**作業内容**:
- `scripts/start-redis-insight.sh`ファイルを新規作成
- 実行権限を付与（`chmod +x`）
- 既存の起動スクリプトと同じパターンで実装
- `config/{env}/cacheserver.yaml`からRedis接続アドレスを読み取る（オプション）
- デフォルトでは`docker-compose.redis.yml`のコンテナ名（`redis`）を使用
- 環境変数`REDIS_HOSTS`を設定（形式: `local:ホスト名:ポート`）
- `start`コマンド: Docker ComposeでRedis Insightを起動
- `stop`コマンド: Docker ComposeでRedis Insightを停止
- 適切なフィードバックメッセージを出力（Web UIのURLを含む）

**受け入れ基準**:
- `scripts/start-redis-insight.sh`ファイルが作成されている
- 実行権限が付与されている
- `./scripts/start-redis-insight.sh start`でRedis Insightが起動できること
- `./scripts/start-redis-insight.sh stop`でRedis Insightが停止できること
- `docker-compose.redis.yml`で起動したRedisに接続できること
- 既存の起動スクリプトと同じパターンに従っている

---

### Phase 2: 設定管理

#### - [x] タスク 2.1: config.goのRedis設定構造体の拡張
**目的**: ジョブキュー用とデフォルト用のRedis設定を分離する構造体を追加する

**作業内容**:
- `server/internal/config/config.go`を開く
- `CacheServerConfig`構造体を確認
- `RedisConfig`構造体を追加または拡張:
  - `JobQueue RedisSingleConfig` - ジョブキュー用（単一接続）
  - `Default RedisClusterConfig` - デフォルト用（複数台対応）
- `RedisSingleConfig`構造体を追加:
  - `Addr string` - 単一Redis接続アドレス
- `RedisClusterConfig`構造体を追加:
  - `Cluster RedisClusterOptions` - Cluster設定
- `RedisClusterOptions`構造体を追加:
  - `Addrs []string` - Redis Clusterのアドレスリスト
- 適切な`mapstructure`タグを設定

**受け入れ基準**:
- `RedisConfig`構造体が拡張されている
- `RedisSingleConfig`、`RedisClusterConfig`、`RedisClusterOptions`構造体が定義されている
- すべてのフィールドに適切な`mapstructure`タグが設定されている
- 既存のコードスタイルに従っている

---

#### - [x] タスク 2.2: cacheserver.yamlの設定追加
**目的**: 各環境の設定ファイルにRedis接続設定を追加する

**作業内容**:
- `config/develop/cacheserver.yaml`を開く
- `redis`セクションを追加:
  - `jobqueue.addr: "localhost:6379"` - ジョブキュー用（単一接続）
  - `default.cluster.addrs: ["localhost:6379"]` - デフォルト用（複数台対応）
- `config/staging/cacheserver.yaml`に同様の設定を追加
- `config/production/cacheserver.yaml.example`に同様の設定を追加（コメント付き）

**受け入れ基準**:
- 3つの設定ファイルに適切な設定が追加されている
- 設定値が適切に設定されている
- YAML形式が正しい
- ジョブキュー用とデフォルト用の設定が分離されている

---

#### - [x] タスク 2.3: rate limit設定の拡張（config.yaml）
**目的**: rate limitにストレージタイプ設定を追加する

**作業内容**:
- `server/internal/config/config.go`を開く
- `RateLimitConfig`構造体を確認
- `StorageType string`フィールドを追加（`mapstructure:"storage_type"`タグ付き）
- `config/develop/config.yaml`の`api.rate_limit`セクションに`storage_type: "auto"`を追加
- `config/staging/config.yaml`に同様の設定を追加
- `config/production/config.yaml.example`に同様の設定を追加（コメント付き）

**受け入れ基準**:
- `RateLimitConfig`構造体に`StorageType`フィールドが追加されている
- 3つの設定ファイルに適切な設定が追加されている
- デフォルト値は`"auto"`（既存の動作を維持）

---

#### - [x] タスク 2.4: rate limit実装の修正
**目的**: rate limitがデフォルト用Redis設定を使用するように修正する

**作業内容**:
- `server/internal/ratelimit/middleware.go`を開く
- `initStore`関数を確認
- `storage_type`設定を確認する処理を追加:
  - `"memory"`が指定された場合はIn-Memoryストレージを使用
  - `"redis"`が指定された場合はRedisを使用（`cfg.CacheServer.Redis.Default.Cluster.Addrs`を使用）
  - `"auto"`の場合は既存の動作（設定に基づいて自動判定）
- `cfg.CacheServer.Redis.Default.Cluster.Addrs`を使用するように変更
- エラーハンドリングを追加（`storage_type`が`"redis"`だが設定がない場合）

**受け入れ基準**:
- `initStore`関数が修正されている
- デフォルト用Redis設定（`redis.default.cluster.addrs`）を使用している
- `storage_type`設定に基づいて適切なストレージが選択される
- エラーハンドリングが適切に実装されている

---

### Phase 3: Asynqライブラリの導入とジョブキューシステム実装

#### - [x] タスク 3.1: Asynqライブラリの導入
**目的**: Asynqライブラリをgo.modに追加する

**作業内容**:
- `go.mod`ファイルを開く
- `github.com/hibiken/asynq`を依存関係に追加
- `go mod tidy`を実行して依存関係を解決

**受け入れ基準**:
- `go.mod`に`github.com/hibiken/asynq`が追加されている
- `go mod tidy`がエラーなく実行できること
- 依存関係が正しく解決されている

---

#### - [x] タスク 3.2: jobqueue/constants.goの作成
**目的**: ジョブタイプ定数とデフォルト値を定義する

**作業内容**:
- `server/internal/service/jobqueue/`ディレクトリを新規作成
- `constants.go`ファイルを新規作成
- `JobTypeDelayPrint = "demo:delay_print"`定数を定義
- `DefaultDelaySeconds = 180`定数を定義（3分 = 180秒）
- `DefaultMaxRetry = 10`定数を定義
- 適切なコメントを追加

**受け入れ基準**:
- `constants.go`ファイルが作成されている
- すべての定数が定義されている
- 適切なコメントが追加されている
- 既存のコードスタイルに従っている

---

#### - [x] タスク 3.3: jobqueue/client.goの作成
**目的**: AsynqクライアントをラップするClient構造体を実装する

**作業内容**:
- `client.go`ファイルを新規作成
- `JobOptions`構造体を定義:
  - `MaxRetry int` - 最大リトライ回数
  - `DelaySeconds int` - 遅延時間（秒）
- `Client`構造体を定義:
  - `client *asynq.Client` - Asynqクライアント
- `NewClient(cfg *config.Config) (*Client, error)`関数を実装:
  - `cfg.CacheServer.Redis.JobQueue.Addr`からRedis接続アドレスを取得
  - デフォルト値: `"localhost:6379"`
  - `asynq.NewClient`でクライアントを作成
- `EnqueueJob(ctx context.Context, jobType string, payload []byte, opts *JobOptions) (*asynq.TaskInfo, error)`メソッドを実装:
  - `opts`が`nil`の場合はデフォルト値を使用
  - `asynq.ProcessIn`で遅延時間を設定
  - `asynq.MaxRetry`で最大リトライ回数を設定
  - ジョブをキューに登録
- `Close() error`メソッドを実装

**受け入れ基準**:
- `client.go`ファイルが作成されている
- `JobOptions`構造体が定義されている
- `Client`構造体とメソッドが実装されている
- Redis接続設定が正しく読み込まれること
- ジョブ登録が正常に動作すること
- エラーハンドリングが適切に実装されている

---

#### - [x] タスク 3.4: jobqueue/processor.goの作成
**目的**: ジョブ処理関数を実装する

**作業内容**:
- `processor.go`ファイルを新規作成
- `DelayPrintPayload`構造体を定義:
  - `Message string` - 出力するメッセージ
- `ProcessDelayPrintJob(ctx context.Context, t *asynq.Task) error`関数を実装:
  - ペイロードの解析（JSON形式）
  - ペイロードがない場合はデフォルトメッセージを使用
  - 標準出力に文字列を出力（タイムスタンプ付き）
  - エラーハンドリング

**受け入れ基準**:
- `processor.go`ファイルが作成されている
- `DelayPrintPayload`構造体が定義されている
- `ProcessDelayPrintJob`関数が実装されている
- 標準出力に適切な形式で文字列が出力されること
- エラーハンドリングが適切に実装されている

---

#### - [x] タスク 3.5: jobqueue/server.goの作成
**目的**: AsynqサーバーをラップするServer構造体を実装する

**作業内容**:
- `server.go`ファイルを新規作成
- `Server`構造体を定義:
  - `server *asynq.Server` - Asynqサーバー
  - `mux *asynq.ServeMux` - ジョブハンドラーのマルチプレクサー
- `NewServer(cfg *config.Config) (*Server, error)`関数を実装:
  - `cfg.CacheServer.Redis.JobQueue.Addr`からRedis接続アドレスを取得
  - デフォルト値: `"localhost:6379"`
  - `asynq.NewServer`でサーバーを作成（Concurrency: 10）
  - `asynq.NewServeMux`でマルチプレクサーを作成
  - `JobTypeDelayPrint`のハンドラーを登録（`ProcessDelayPrintJob`）
- `Start() error`メソッドを実装:
  - バックグラウンドでジョブ処理を開始
- `Shutdown() error`メソッドを実装:
  - サーバーを停止

**受け入れ基準**:
- `server.go`ファイルが作成されている
- `Server`構造体とメソッドが実装されている
- Redis接続設定が正しく読み込まれること
- ジョブハンドラーが正しく登録されること
- サーバーが正常に起動・停止できること

---

### Phase 4: APIハンドラーの実装

#### - [x] タスク 4.1: dm_jobqueue_handler.goの作成
**目的**: ジョブ登録APIのハンドラーを実装する

**作業内容**:
- `server/internal/api/handler/dm_jobqueue_handler.go`ファイルを新規作成
- `DmJobqueueHandler`構造体を定義:
  - `jobQueueClient *jobqueue.Client` - ジョブキュークライアント
- `NewDmJobqueueHandler(jobQueueClient *jobqueue.Client) *DmJobqueueHandler`関数を実装
- `RegisterJobRequest`構造体を定義:
  - `Message string` - 出力するメッセージ（オプション）
  - `DelaySeconds int` - 遅延時間（秒、オプション）
  - `MaxRetry int` - 最大リトライ回数（オプション）
- `RegisterJobResponse`構造体を定義:
  - `JobID string` - ジョブID
  - `Status string` - ステータス（"registered"）
- `RegisterJob(ctx context.Context, req *RegisterJobRequest) (*RegisterJobResponse, error)`メソッドを実装:
  - Redis接続が利用できない場合のエラーハンドリング（503エラー）
  - メッセージの設定（デフォルト値）
  - ペイロードの作成（JSON形式）
  - `JobOptions`の作成
  - ジョブをキューに登録
  - レスポンスを返す
- `RegisterDmJobqueueEndpoints(api huma.API, h *DmJobqueueHandler)`関数を実装:
  - `POST /api/dm-jobqueue/register`エンドポイントを登録
  - Huma APIの`huma.Register`を使用
  - 参考コード用の名前を使用（`register-demo-job`、`jobqueue-demo`タグ）

**受け入れ基準**:
- `dm_jobqueue_handler.go`ファイルが作成されている
- すべての構造体とメソッドが実装されている
- Redis接続エラー時に503エラーを返すこと
- ジョブ登録が正常に動作すること
- エラーハンドリングが適切に実装されている
- 参考コード用の名前が使用されている

---

#### - [x] タスク 4.2: router.goへのエンドポイント登録
**目的**: ジョブキューエンドポイントをルーターに登録する

**作業内容**:
- `server/internal/api/router/router.go`を開く
- `RegisterDmJobqueueEndpoints`関数を呼び出す処理を追加
- 適切な位置に配置（既存のエンドポイント登録パターンに従う）

**受け入れ基準**:
- `router.go`にエンドポイント登録処理が追加されている
- エンドポイントが正しく登録されること
- 既存のパターンに従っている

---

### Phase 5: main.goへの統合

#### - [x] タスク 5.1: main.goへのAsynqクライアント初期化
**目的**: Asynqクライアントをmain.goで初期化する

**作業内容**:
- `server/cmd/server/main.go`を開く
- `jobqueue.NewClient(cfg)`を呼び出してクライアントを作成
- Redis接続エラーの場合:
  - 標準エラー出力に警告ログを出力
  - 起動処理は継続する（`log.Fatalf`を使用しない）
  - `jobQueueClient`を`nil`に設定
- エラーがない場合は`defer jobQueueClient.Close()`を追加

**受け入れ基準**:
- Asynqクライアントが初期化されている
- Redis接続エラー時に警告ログが出力されること
- エラーが発生してもAPIサーバーの起動が継続すること
- クライアントが適切にクローズされること

---

#### - [x] タスク 5.2: main.goへのAsynqサーバー初期化と起動
**目的**: Asynqサーバーをmain.goで初期化し、バックグラウンドで起動する

**作業内容**:
- `server/cmd/server/main.go`を開く
- `jobQueueClient`が`nil`でない場合のみサーバーを初期化
- `jobqueue.NewServer(cfg)`を呼び出してサーバーを作成
- Redis接続エラーの場合:
  - 標準エラー出力に警告ログを出力
  - 起動処理は継続する
- エラーがない場合:
  - バックグラウンドで`jobQueueServer.Start()`を実行（goroutine）
  - 起動エラーを標準エラー出力に記録

**受け入れ基準**:
- Asynqサーバーが初期化されている
- バックグラウンドでジョブ処理が開始されること
- Redis接続エラー時に警告ログが出力されること
- エラーが発生してもAPIサーバーの起動が継続すること

---

#### - [x] タスク 5.3: main.goへのハンドラー初期化とエンドポイント登録
**目的**: ジョブキューハンドラーを初期化し、エンドポイントを登録する

**作業内容**:
- `server/cmd/server/main.go`を開く
- `handler.NewDmJobqueueHandler(jobQueueClient)`を呼び出してハンドラーを作成（`jobQueueClient`が`nil`の場合も許可）
- `handler.RegisterDmJobqueueEndpoints(humaAPI, jobQueueHandler)`を呼び出してエンドポイントを登録

**受け入れ基準**:
- ジョブキューハンドラーが初期化されている
- エンドポイントが正しく登録されること
- `jobQueueClient`が`nil`の場合でもハンドラーが初期化できること

---

### Phase 6: クライアント側UIの実装

#### - [x] タスク 6.1: クライアント側ジョブ登録ボタンの実装
**目的**: クライアント側にジョブ登録ボタンを追加する

**作業内容**:
- クライアント側の適切な場所にボタンを追加
- 参考コードとして利用するため、邪魔にならない名前を使用
- ボタンをクリックしたら`POST /api/dm-jobqueue/register`にリクエストを送信
- リクエストボディに`message`、`delay_seconds`、`max_retry`を含める（オプション）
- ジョブ登録の成功・失敗をユーザーに通知
- 既存のUIパターンに従った実装

**受け入れ基準**:
- ボタンが追加されている
- ボタンをクリックするとAPIリクエストが送信されること
- 成功・失敗の通知が表示されること
- 既存のUIパターンに従っている
- 参考コード用の名前が使用されている

---

### Phase 7: テスト実装

#### - [x] タスク 7.1: jobqueueパッケージの単体テスト
**目的**: jobqueueパッケージの各関数・メソッドの単体テストを実装する

**作業内容**:
- `server/internal/service/jobqueue/client_test.go`を作成
- `NewClient`関数のテスト（正常系・異常系）
- `EnqueueJob`メソッドのテスト（正常系・異常系、デフォルト値の確認）
- `Close`メソッドのテスト
- `server/internal/service/jobqueue/processor_test.go`を作成
- `ProcessDelayPrintJob`関数のテスト（正常系・異常系）
- `server/internal/service/jobqueue/constants_test.go`を作成（必要に応じて）
- テストは既存のテストパターンに従う

**受け入れ基準**:
- テストファイルが作成されている
- 主要な関数・メソッドのテストが実装されている
- テストが正常に実行できること
- 既存のテストパターンに従っている

---

#### - [x] タスク 7.2: APIハンドラーの単体テスト
**目的**: ジョブ登録APIハンドラーの単体テストを実装する

**作業内容**:
- `server/internal/api/handler/dm_jobqueue_handler_test.go`を作成
- `RegisterJob`メソッドのテスト:
  - 正常系（ジョブ登録成功）
  - 異常系（Redis接続エラー、ペイロードエラー）
  - デフォルト値の確認
- モックを使用して`jobQueueClient`を模擬
- テストは既存のテストパターンに従う

**受け入れ基準**:
- テストファイルが作成されている
- 主要なメソッドのテストが実装されている
- テストが正常に実行できること
- 既存のテストパターンに従っている

---

#### - [x] タスク 7.3: 統合テストの実装
**目的**: Redis接続からジョブ登録・処理までの統合テストを実装する

**作業内容**:
- 統合テストファイルを作成（適切な場所に配置）
- Redis接続のテスト
- ジョブ登録から処理までのフローのテスト
- 遅延時間のテスト
- リトライ設定のテスト
- エラーハンドリングのテスト
- テストは既存のテストパターンに従う

**受け入れ基準**:
- 統合テストファイルが作成されている
- 主要なフローのテストが実装されている
- テストが正常に実行できること
- 既存のテストパターンに従っている

---

#### - [x] タスク 7.4: E2Eテストの実装
**目的**: APIエンドポイントとクライアント側UIのE2Eテストを実装する

**作業内容**:
- E2Eテストファイルを作成（適切な場所に配置）
- APIエンドポイントのテスト（正常系・異常系）
- クライアント側UIのテスト（ボタンクリック、通知表示）
- テストは既存のテストパターンに従う

**受け入れ基準**:
- E2Eテストファイルが作成されている
- 主要なシナリオのテストが実装されている
- テストが正常に実行できること
- 既存のテストパターンに従っている

---

## 実装順序の推奨

1. **Phase 1**: Docker Compose設定と起動スクリプト（インフラ準備）
2. **Phase 2**: 設定管理（設定ファイルの準備）
3. **Phase 3**: Asynqライブラリの導入とジョブキューシステム実装（コア機能）
4. **Phase 4**: APIハンドラーの実装（API層）
5. **Phase 5**: main.goへの統合（統合）
6. **Phase 6**: クライアント側UIの実装（フロントエンド）
7. **Phase 7**: テスト実装（品質保証）

## 注意事項

- 各タスクは独立して実装可能な粒度で分解されている
- タスクの順序は推奨順序であり、必要に応じて調整可能
- 実装時は既存のコードスタイルとパターンに従うこと
- 参考コードとして利用するため、将来の本実装に影響しない名前を使用すること
- エラーハンドリングを適切に実装すること
- テストは実装と並行して進めることを推奨
