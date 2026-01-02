# Redis遅延接続・自動再接続機能実装タスク一覧

## 概要
Redis遅延接続・自動再接続機能の実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: 設定構造体の拡張

#### - [ ] タスク 1.1: RedisClusterConfig構造体に接続オプションフィールドを追加
**目的**: Cache Server用（Cluster接続）の接続オプション設定を追加する

**作業内容**:
- `server/internal/config/config.go`を開く
- `RedisClusterConfig`構造体に以下のフィールドを追加:
  ```go
  MaxRetries      int           `mapstructure:"max_retries"`       // コマンド失敗時の最大リトライ数（デフォルト: 2）
  MinRetryBackoff time.Duration `mapstructure:"min_retry_backoff"` // リトライ間隔（最小）（デフォルト: 8ms）
  MaxRetryBackoff time.Duration `mapstructure:"max_retry_backoff"` // リトライ間隔（最大）（デフォルト: 512ms）
  DialTimeout     time.Duration `mapstructure:"dial_timeout"`      // 接続確立のタイムアウト（デフォルト: 5s）
  ReadTimeout     time.Duration `mapstructure:"read_timeout"`      // 読み取りタイムアウト（デフォルト: 3s）
  PoolSize        int           `mapstructure:"pool_size"`          // 接続プールサイズ（デフォルト: CPU数×10）
  PoolTimeout     time.Duration `mapstructure:"pool_timeout"`       // プールから接続を取り出す際の待機時間（デフォルト: 4s）
  ```
- 適切なコメントを追加
- `time`パッケージのインポートを確認（必要に応じて追加）

**受け入れ基準**:
- `RedisClusterConfig`構造体に接続オプションフィールドが追加されている
- すべてのフィールドに適切な`mapstructure`タグが設定されている
- 既存のコードスタイルに従っている

---

#### - [ ] タスク 1.2: RedisSingleConfig構造体に接続オプションフィールドを追加
**目的**: Jobqueue用（単一接続）の接続オプション設定を追加する

**作業内容**:
- `server/internal/config/config.go`を開く
- `RedisSingleConfig`構造体に以下のフィールドを追加:
  ```go
  MaxRetries      int           `mapstructure:"max_retries"`       // コマンド失敗時の最大リトライ数（デフォルト: 2）
  MinRetryBackoff time.Duration `mapstructure:"min_retry_backoff"` // リトライ間隔（最小）（デフォルト: 8ms）
  MaxRetryBackoff time.Duration `mapstructure:"max_retry_backoff"` // リトライ間隔（最大）（デフォルト: 512ms）
  DialTimeout     time.Duration `mapstructure:"dial_timeout"`      // 接続確立のタイムアウト（デフォルト: 5s）
  ReadTimeout     time.Duration `mapstructure:"read_timeout"`      // 読み取りタイムアウト（デフォルト: 3s）
  WriteTimeout    time.Duration `mapstructure:"write_timeout"`     // 書き込みタイムアウト（デフォルト: 3s）
  PoolSize        int           `mapstructure:"pool_size"`          // 接続プールサイズ（デフォルト: CPU数×10）
  PoolTimeout     time.Duration `mapstructure:"pool_timeout"`       // プールから接続を取り出す際の待機時間（デフォルト: 4s）
  ```
- 適切なコメントを追加
- `time`パッケージのインポートを確認（必要に応じて追加）

**受け入れ基準**:
- `RedisSingleConfig`構造体に接続オプションフィールドが追加されている
- すべてのフィールドに適切な`mapstructure`タグが設定されている
- `WriteTimeout`フィールドが含まれている（`asynq.RedisClientOpt`で使用）
- 既存のコードスタイルに従っている

---

#### - [ ] タスク 1.3: 設定ファイルの例を更新
**目的**: 設定ファイルに接続オプション設定の例を追加する

**作業内容**:
- `config/develop/cacheserver.yaml`を開く
- `redis.jobqueue`セクションに接続オプション設定の例を追加（コメント付き）:
  ```yaml
  redis:
    jobqueue:
      addr: "localhost:6379"
      # 接続オプション（オプション、未設定の場合はデフォルト値を使用）
      max_retries: 2
      min_retry_backoff: 8ms
      max_retry_backoff: 512ms
      dial_timeout: 5s
      read_timeout: 3s
      write_timeout: 3s
      pool_size: 10
      pool_timeout: 4s
  ```
- `redis.default.cluster`セクションに接続オプション設定の例を追加（コメント付き）:
  ```yaml
    default:
      cluster:
        addrs:
          - "localhost:6379"
        # 接続オプション（オプション、未設定の場合はデフォルト値を使用）
        max_retries: 2
        min_retry_backoff: 8ms
        max_retry_backoff: 512ms
        dial_timeout: 5s
        read_timeout: 3s
        pool_size: 10
        pool_timeout: 4s
  ```
- 既存の設定との互換性を保つ（接続オプションはオプション）

**受け入れ基準**:
- 設定ファイルに接続オプション設定の例が追加されている
- コメントで「オプション」であることが明記されている
- 既存の設定（`addr`や`addrs`のみ）でも動作すること

---

#### - [ ] タスク 1.4: staging環境の設定ファイルに接続オプション設定を追加
**目的**: staging環境の設定ファイルに接続オプション設定の例を追加する

**作業内容**:
- `config/staging/cacheserver.yaml`を開く
- `redis.jobqueue`セクションに接続オプション設定の例を追加（コメント付き）:
  ```yaml
  redis:
    jobqueue:
      addr: "localhost:6379"
      # 接続オプション（オプション、未設定の場合はデフォルト値を使用）
      max_retries: 2
      min_retry_backoff: 8ms
      max_retry_backoff: 512ms
      dial_timeout: 5s
      read_timeout: 3s
      write_timeout: 3s
      pool_size: 10
      pool_timeout: 4s
  ```
- `redis.default.cluster`セクションに接続オプション設定の例を追加（コメント付き）:
  ```yaml
    default:
      cluster:
        addrs:
          - "localhost:6379"
        # 接続オプション（オプション、未設定の場合はデフォルト値を使用）
        max_retries: 2
        min_retry_backoff: 8ms
        max_retry_backoff: 512ms
        dial_timeout: 5s
        read_timeout: 3s
        pool_size: 10
        pool_timeout: 4s
  ```
- 既存の設定との互換性を保つ（接続オプションはオプション）

**受け入れ基準**:
- staging環境の設定ファイルに接続オプション設定の例が追加されている
- コメントで「オプション」であることが明記されている
- 既存の設定（`addr`や`addrs`のみ）でも動作すること

---

#### - [ ] タスク 1.5: production環境の設定ファイルに接続オプション設定を追加
**目的**: production環境の設定ファイルに接続オプション設定の例を追加する

**作業内容**:
- `config/production/cacheserver.yaml.example`を開く
- `redis.jobqueue`セクションに接続オプション設定の例を追加（コメント付き）:
  ```yaml
  redis:
    jobqueue:
      addr: "localhost:6379"
      # 接続オプション（オプション、未設定の場合はデフォルト値を使用）
      max_retries: 2
      min_retry_backoff: 8ms
      max_retry_backoff: 512ms
      dial_timeout: 5s
      read_timeout: 3s
      write_timeout: 3s
      pool_size: 10
      pool_timeout: 4s
  ```
- `redis.default.cluster`セクションに接続オプション設定の例を追加（コメント付き）:
  ```yaml
    default:
      cluster:
        addrs:
          - "localhost:6379"
        # 接続オプション（オプション、未設定の場合はデフォルト値を使用）
        max_retries: 2
        min_retry_backoff: 8ms
        max_retry_backoff: 512ms
        dial_timeout: 5s
        read_timeout: 3s
        pool_size: 10
        pool_timeout: 4s
  ```
- 既存の設定との互換性を保つ（接続オプションはオプション）

**受け入れ基準**:
- production環境の設定ファイルに接続オプション設定の例が追加されている
- コメントで「オプション」であることが明記されている
- 既存の設定（`addr`や`addrs`のみ）でも動作すること

---

### Phase 2: Cache Server用の接続オプション設定実装

#### - [ ] タスク 2.1: initStore関数の1箇所目で接続オプションを設定
**目的**: `initStore`関数内の1箇所目（144行目付近）で接続オプションを設定する

**作業内容**:
- `server/internal/ratelimit/middleware.go`を開く
- `initStore`関数内の1箇所目（`storageType == "redis"`の分岐内、148-150行目付近）を修正
- `redis.NewClusterClient`の呼び出し前に接続オプションを設定:
  ```go
  clusterOpts := &redis.ClusterOptions{
      Addrs: cfg.CacheServer.Redis.Default.Cluster.Addrs,
  }

  // 接続オプションの設定（設定ファイルから読み込む、未設定の場合はデフォルト値を使用）
  if cfg.CacheServer.Redis.Default.Cluster.MaxRetries > 0 {
      clusterOpts.MaxRetries = cfg.CacheServer.Redis.Default.Cluster.MaxRetries
  } else {
      clusterOpts.MaxRetries = 2 // デフォルト値
  }

  if cfg.CacheServer.Redis.Default.Cluster.MinRetryBackoff > 0 {
      clusterOpts.MinRetryBackoff = cfg.CacheServer.Redis.Default.Cluster.MinRetryBackoff
  } else {
      clusterOpts.MinRetryBackoff = 8 * time.Millisecond // デフォルト値
  }

  if cfg.CacheServer.Redis.Default.Cluster.MaxRetryBackoff > 0 {
      clusterOpts.MaxRetryBackoff = cfg.CacheServer.Redis.Default.Cluster.MaxRetryBackoff
  } else {
      clusterOpts.MaxRetryBackoff = 512 * time.Millisecond // デフォルト値
  }

  if cfg.CacheServer.Redis.Default.Cluster.DialTimeout > 0 {
      clusterOpts.DialTimeout = cfg.CacheServer.Redis.Default.Cluster.DialTimeout
  } else {
      clusterOpts.DialTimeout = 5 * time.Second // デフォルト値
  }

  if cfg.CacheServer.Redis.Default.Cluster.ReadTimeout > 0 {
      clusterOpts.ReadTimeout = cfg.CacheServer.Redis.Default.Cluster.ReadTimeout
  } else {
      clusterOpts.ReadTimeout = 3 * time.Second // デフォルト値
  }

  if cfg.CacheServer.Redis.Default.Cluster.PoolSize > 0 {
      clusterOpts.PoolSize = cfg.CacheServer.Redis.Default.Cluster.PoolSize
  } else {
      clusterOpts.PoolSize = 10 * runtime.NumCPU() // デフォルト値: CPU数×10
  }

  if cfg.CacheServer.Redis.Default.Cluster.PoolTimeout > 0 {
      clusterOpts.PoolTimeout = cfg.CacheServer.Redis.Default.Cluster.PoolTimeout
  } else {
      clusterOpts.PoolTimeout = 4 * time.Second // デフォルト値
  }

  rdb := redis.NewClusterClient(clusterOpts)
  ```
- `time`と`runtime`パッケージのインポートを確認（必要に応じて追加）

**受け入れ基準**:
- 接続オプションが正しく設定されている
- 設定値が0以下の場合にデフォルト値が使用される
- `PoolSize`のデフォルト値が`runtime.NumCPU() * 10`であること
- 既存のコードスタイルに従っている

---

#### - [ ] タスク 2.2: initStore関数の2箇所目で接続オプションを設定
**目的**: `initStore`関数内の2箇所目（163行目付近）で接続オプションを設定する

**作業内容**:
- `server/internal/ratelimit/middleware.go`を開く
- `initStore`関数内の2箇所目（`storageType == "auto"`の分岐内、163-165行目付近）を修正
- タスク2.1と同じ接続オプション設定ロジックを適用
- コードの重複を避けるため、共通のヘルパー関数を作成することも検討可能

**受け入れ基準**:
- 接続オプションが正しく設定されている
- タスク2.1と同じロジックが適用されている
- 既存のコードスタイルに従っている

---

### Phase 3: Jobqueue用の接続オプション設定実装

#### - [ ] タスク 3.1: NewClient関数で接続オプションを設定
**目的**: Jobqueue用の`NewClient`関数で接続オプションを設定する

**作業内容**:
- `server/internal/service/jobqueue/client.go`を開く
- `NewClient`関数内の`asynq.RedisClientOpt`設定部分を修正
- 接続オプションを設定:
  ```go
  redisOpt := asynq.RedisClientOpt{
      Addr: redisAddr,
  }

  // 接続オプションの設定（設定ファイルから読み込む、未設定の場合はデフォルト値を使用）
  if cfg.CacheServer.Redis.JobQueue.DialTimeout > 0 {
      redisOpt.DialTimeout = cfg.CacheServer.Redis.JobQueue.DialTimeout
  } else {
      redisOpt.DialTimeout = 5 * time.Second // デフォルト値
  }

  if cfg.CacheServer.Redis.JobQueue.ReadTimeout > 0 {
      redisOpt.ReadTimeout = cfg.CacheServer.Redis.JobQueue.ReadTimeout
  } else {
      redisOpt.ReadTimeout = 3 * time.Second // デフォルト値
  }

  if cfg.CacheServer.Redis.JobQueue.WriteTimeout > 0 {
      redisOpt.WriteTimeout = cfg.CacheServer.Redis.JobQueue.WriteTimeout
  } else {
      redisOpt.WriteTimeout = 3 * time.Second // デフォルト値
  }

  client := asynq.NewClient(redisOpt)

  // リトライ設定は、asynqが内部的に使用するgo-redisクライアントのオプションを直接設定
  // MakeRedisClient()でredis.Clientを取得して設定
  if redisClient, ok := redisOpt.MakeRedisClient().(*redis.Client); ok {
      if cfg.CacheServer.Redis.JobQueue.MaxRetries > 0 {
          redisClient.Options().MaxRetries = cfg.CacheServer.Redis.JobQueue.MaxRetries
      } else {
          redisClient.Options().MaxRetries = 2 // デフォルト値
      }

      if cfg.CacheServer.Redis.JobQueue.MinRetryBackoff > 0 {
          redisClient.Options().MinRetryBackoff = cfg.CacheServer.Redis.JobQueue.MinRetryBackoff
      } else {
          redisClient.Options().MinRetryBackoff = 8 * time.Millisecond // デフォルト値
      }

      if cfg.CacheServer.Redis.JobQueue.MaxRetryBackoff > 0 {
          redisClient.Options().MaxRetryBackoff = cfg.CacheServer.Redis.JobQueue.MaxRetryBackoff
      } else {
          redisClient.Options().MaxRetryBackoff = 512 * time.Millisecond // デフォルト値
      }

      if cfg.CacheServer.Redis.JobQueue.PoolSize > 0 {
          redisClient.Options().PoolSize = cfg.CacheServer.Redis.JobQueue.PoolSize
      } else {
          redisClient.Options().PoolSize = 10 * runtime.NumCPU() // デフォルト値: CPU数×10
      }

      if cfg.CacheServer.Redis.JobQueue.PoolTimeout > 0 {
          redisClient.Options().PoolTimeout = cfg.CacheServer.Redis.JobQueue.PoolTimeout
      } else {
          redisClient.Options().PoolTimeout = 4 * time.Second // デフォルト値
      }
  }
  ```
- `time`、`runtime`、`github.com/redis/go-redis/v9`パッケージのインポートを確認（必要に応じて追加）

**受け入れ基準**:
- 接続オプションが正しく設定されている
- `asynq.RedisClientOpt`で直接設定できるオプション（`DialTimeout`, `ReadTimeout`, `WriteTimeout`）が設定されている
- リトライ設定が`MakeRedisClient()`で取得した`redis.Client`のオプションに設定されている
- 設定値が0以下の場合にデフォルト値が使用される
- 既存のコードスタイルに従っている

---

#### - [ ] タスク 3.2: NewServer関数で接続オプションを設定
**目的**: Jobqueue用の`NewServer`関数で接続オプションを設定する

**作業内容**:
- `server/internal/service/jobqueue/server.go`を開く
- `NewServer`関数内の`asynq.RedisClientOpt`設定部分を修正
- タスク3.1と同じ接続オプション設定ロジックを適用
- コードの重複を避けるため、共通のヘルパー関数を作成することも検討可能

**受け入れ基準**:
- 接続オプションが正しく設定されている
- タスク3.1と同じロジックが適用されている
- 既存のコードスタイルに従っている

---

### Phase 4: 既存環境の確認

#### - [ ] タスク 4.1: docker-compose.redis.ymlの確認
**目的**: 既存のRedis環境設定を確認する

**作業内容**:
- `docker-compose.redis.yml`ファイルを確認
- 以下の項目を確認:
  - Redisコンテナの定義が存在するか
  - ポート6379が公開されているか
  - データ永続化の設定が存在するか
  - ネットワーク設定が適切か
- 確認結果を記録

**受け入れ基準**:
- `docker-compose.redis.yml`ファイルが存在する
- すべての確認項目が適切である
- 確認結果が記録されている

---

#### - [ ] タスク 4.2: scripts/start-redis.shの確認
**目的**: 既存のRedis起動スクリプトを確認する

**作業内容**:
- `scripts/start-redis.sh`ファイルを確認
- 以下の項目を確認:
  - Docker Composeを使用してRedisを起動するか
  - 既存の起動スクリプト（`start-mailpit.sh`など）と同じパターンか
  - `start`/`stop`コマンドをサポートしているか
- 確認結果を記録

**受け入れ基準**:
- `scripts/start-redis.sh`ファイルが存在する
- すべての確認項目が適切である
- 確認結果が記録されている

---

#### - [ ] タスク 4.3: config/{env}/cacheserver.yamlの確認
**目的**: 既存の設定ファイルを確認する

**作業内容**:
- `config/develop/cacheserver.yaml`を確認
- 以下の項目を確認:
  - `redis.jobqueue.addr`が設定されているか
  - `redis.default.cluster.addrs`が設定されているか
  - 既存の設定が正しく読み込まれているか
- 確認結果を記録

**受け入れ基準**:
- 設定ファイルが存在する
- `redis.jobqueue.addr`が設定されている
- `redis.default.cluster.addrs`が設定されている（空配列でも可）
- 確認結果が記録されている

---

#### - [ ] タスク 4.4: config/{env}/config.yamlの確認
**目的**: Rate Limitのストレージタイプ設定を確認する

**作業内容**:
- `config/develop/config.yaml`を確認
- 以下の項目を確認:
  - `api.rate_limit.storage_type`が設定されているか
  - `storage_type: "auto"`の場合、`cacheserver.yaml`の`redis.default.cluster.addrs`が設定されているか確認
- 確認結果を記録

**受け入れ基準**:
- 設定ファイルが存在する
- `api.rate_limit.storage_type`が設定されている
- 確認結果が記録されている

---

### Phase 5: 動作確認前の設定変更

#### - [ ] タスク 5.1: Cache Server用の動作確認設定変更
**目的**: Cache Server用の動作確認のため、設定ファイルを変更する

**作業内容**:
- `config/develop/config.yaml`を開く
- `api.rate_limit.storage_type`を`"redis"`に変更するか、`config/develop/cacheserver.yaml`の`redis.default.cluster.addrs`に`["localhost:6379"]`を設定
- 変更内容を記録（動作確認後に元に戻すため）

**受け入れ基準**:
- 設定ファイルが変更されている
- Cache Server用のRedis接続が使用される設定になっている
- 変更内容が記録されている

---

### Phase 6: 動作確認

#### - [ ] タスク 6.1: Cache Server用の遅延接続の動作確認
**目的**: Cache Server用の遅延接続が正しく動作することを確認する

**作業内容**:
1. Redisを停止した状態でAPIサーバーを起動
2. サーバーが正常に起動することを確認（Redis接続なしで起動可能）
3. 最初のAPIリクエストを送信（レートリミットチェックが実行される）
4. ログで接続確立のタイミングを確認
5. Redisが接続されたことを確認

**受け入れ基準**:
- サーバー起動時にはRedis接続が確立されていない
- 最初のRedisコマンド実行時に接続が確立される
- ログで接続確立のタイミングが確認できる

---

#### - [ ] タスク 6.2: Cache Server用の自動再接続の動作確認
**目的**: Cache Server用の自動再接続が正しく動作することを確認する

**作業内容**:
1. Redisを起動した状態でAPIサーバーを起動
2. 正常に動作することを確認
3. Redisを停止
4. APIリクエストを送信（エラーが発生することを確認）
5. Redisを再起動
6. 次のAPIリクエストを送信（自動的に再接続されることを確認）
7. ログで再接続のタイミングを確認

**受け入れ基準**:
- Redis停止時はエラーが発生する
- Redis再起動後、次のコマンド実行時に自動的に再接続される
- ログで再接続のタイミングが確認できる

---

#### - [ ] タスク 6.3: Cache Server用のリトライ機能の動作確認
**目的**: Cache Server用のリトライ機能が正しく動作することを確認する

**作業内容**:
1. Redisを停止した状態でAPIサーバーを起動
2. APIリクエストを送信（レートリミットチェックが実行される）
3. ログでリトライ処理を確認
4. リトライ間隔が`MinRetryBackoff`と`MaxRetryBackoff`の範囲内であることを確認

**受け入れ基準**:
- 接続エラー時に自動的にリトライが実行される
- 最大2回までリトライ（初回 + 2回のリトライ）
- リトライ間隔は8ms～512msの範囲で指数バックオフ
- ログでリトライ処理が確認できる

---

#### - [ ] タスク 6.4: Jobqueue用の遅延接続の動作確認
**目的**: Jobqueue用の遅延接続が正しく動作することを確認する

**作業内容**:
1. Redisを停止した状態でAPIサーバーを起動
2. サーバーが正常に起動することを確認（Redis接続なしで起動可能）
3. Jobqueueクライアントを作成（接続はまだ確立されていない）
4. ジョブを登録（最初のRedisコマンド実行時）
5. ログで接続確立のタイミングを確認
6. Redisが接続されたことを確認

**受け入れ基準**:
- サーバー起動時にはRedis接続が確立されていない
- 最初のRedisコマンド実行時に接続が確立される
- ログで接続確立のタイミングが確認できる

---

#### - [ ] タスク 6.5: Jobqueue用の自動再接続の動作確認
**目的**: Jobqueue用の自動再接続が正しく動作することを確認する

**作業内容**:
1. Redisを起動した状態でAPIサーバーを起動
2. 正常に動作することを確認
3. Redisを停止
4. ジョブ登録を試行（エラーが発生することを確認）
5. Redisを再起動
6. 次のジョブ登録を試行（自動的に再接続されることを確認）
7. ログで再接続のタイミングを確認

**受け入れ基準**:
- Redis停止時はエラーが発生する
- Redis再起動後、次のコマンド実行時に自動的に再接続される
- ログで再接続のタイミングが確認できる

---

#### - [ ] タスク 6.6: Jobqueue用のリトライ機能の動作確認
**目的**: Jobqueue用のリトライ機能が正しく動作することを確認する

**作業内容**:
1. Redisを停止した状態でAPIサーバーを起動
2. ジョブ登録を試行（最初のRedisコマンド実行時）
3. ログでリトライ処理を確認
4. リトライ間隔が`MinRetryBackoff`と`MaxRetryBackoff`の範囲内であることを確認

**受け入れ基準**:
- 接続エラー時に自動的にリトライが実行される
- 最大2回までリトライ（初回 + 2回のリトライ）
- リトライ間隔は8ms～512msの範囲で指数バックオフ
- ログでリトライ処理が確認できる

---

### Phase 7: 動作確認後の設定復元

#### - [ ] タスク 7.1: 動作確認後の設定復元
**目的**: 動作確認のために変更した設定を元に戻す

**作業内容**:
- タスク5.1で変更した設定ファイルを元に戻す
- 変更内容の記録を確認して、元の設定に復元

**受け入れ基準**:
- 設定ファイルが元の状態に復元されている
- 変更内容が記録されている

---

### Phase 8: テスト実装

#### - [ ] タスク 8.1: 設定読み込みテストの実装
**目的**: 設定構造体が正しく設定ファイルから読み込まれることを確認するテストを実装する

**作業内容**:
- `server/internal/config/config_test.go`を開く
- `RedisClusterConfig`構造体の設定読み込みテストを追加
- `RedisSingleConfig`構造体の設定読み込みテストを追加
- デフォルト値が正しく適用されることを確認するテストを追加

**受け入れ基準**:
- テストが実装されている
- 設定ファイルから正しく読み込まれることを確認できる
- デフォルト値が正しく適用されることを確認できる

---

#### - [ ] タスク 8.2: 接続オプション設定テストの実装
**目的**: 接続オプションが正しく設定されることを確認するテストを実装する

**作業内容**:
- `server/internal/ratelimit/middleware_test.go`を開く（存在しない場合は作成）
- `initStore`関数で接続オプションが正しく設定されることを確認するテストを追加
- `server/internal/service/jobqueue/client_test.go`を開く
- `NewClient`関数で接続オプションが正しく設定されることを確認するテストを追加
- `server/internal/service/jobqueue/server_test.go`を開く（存在しない場合は作成）
- `NewServer`関数で接続オプションが正しく設定されることを確認するテストを追加
- 設定値が0以下の場合にデフォルト値が使用されることを確認するテストを追加

**受け入れ基準**:
- テストが実装されている
- 接続オプションが正しく設定されることを確認できる
- デフォルト値が正しく適用されることを確認できる

---

#### - [ ] タスク 8.3: 統合テストの実装
**目的**: 遅延接続、自動再接続、リトライ機能の統合テストを実装する

**作業内容**:
- 統合テストファイルを作成（適切な場所に配置）
- 遅延接続のテストを実装:
  - Redis停止状態でAPIサーバーが起動できることを確認
  - 最初のRedisコマンド実行時に接続が確立されることを確認
- 自動再接続のテストを実装:
  - Redis停止→再起動時に自動的に再接続されることを確認
  - 再接続後の動作が正常であることを確認
- リトライ機能のテストを実装:
  - 接続エラー時にリトライが実行されることを確認
  - リトライ間隔が適切であることを確認

**受け入れ基準**:
- 統合テストが実装されている
- 遅延接続、自動再接続、リトライ機能が正しく動作することを確認できる

---

### Phase 9: 最終確認

#### - [ ] タスク 9.1: 既存テストの実行
**目的**: 既存のテストが全て失敗しないことを確認する

**作業内容**:
- プロジェクトルートで`npm test`または`go test ./...`を実行
- テスト結果を確認
- 失敗したテストがあれば原因を調査して修正

**受け入れ基準**:
- 既存のテストが全て成功する
- 新しく追加したテストも成功する

---

#### - [ ] タスク 9.2: リントチェック
**目的**: コードのリントチェックを実行する

**作業内容**:
- プロジェクトルートで`npm run lint`または`golangci-lint run`を実行
- リント結果を確認
- エラーがあれば修正

**受け入れ基準**:
- リントエラーが0件である
- 警告があれば適切に対応する

---

#### - [ ] タスク 9.3: ビルド確認
**目的**: プロジェクトが正常にビルドできることを確認する

**作業内容**:
- プロジェクトルートで`npm run build`または`go build ./...`を実行
- ビルド結果を確認
- エラーがあれば修正

**受け入れ基準**:
- ビルドが成功する
- 警告があれば適切に対応する

---

## 実装の注意事項

### コードスタイル
- 既存のコードスタイルに従うこと
- 適切なコメントを追加すること
- エラーハンドリングを適切に実装すること

### 後方互換性
- 既存の設定ファイルとの互換性を保つこと（`addr`や`addrs`のみでも動作する）
- 既存のRedis接続機能を壊さないこと

### テスト
- 単体テストと統合テストを実装すること
- 既存のテストが失敗しないことを確認すること

### 動作確認
- Cache Server用とJobqueue用の両方で動作確認を行うこと
- 動作確認前の設定変更と動作確認後の設定復元を忘れないこと
