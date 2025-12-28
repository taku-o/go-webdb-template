# APIレートリミット機能実装タスク一覧

## 概要
APIレートリミット機能の実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: 依存関係の追加

#### - [ ] タスク 1.1: ulule/limiterライブラリの追加
**目的**: レートリミット機能に必要なライブラリを追加

**作業内容**:
- `server/go.mod`に`github.com/ulule/limiter/v3`を追加
- `server/go.mod`に`github.com/redis/go-redis/v9`を追加
- `go mod tidy`を実行して依存関係を解決

**受け入れ基準**:
- `go.mod`に`github.com/ulule/limiter/v3`が追加されている
- `go.mod`に`github.com/redis/go-redis/v9`が追加されている
- `go mod tidy`が正常に実行される
- 既存の依存関係に影響がない

---

### Phase 2: 設定管理の実装

#### - [ ] タスク 2.1: APIConfig構造体にRateLimitConfigを追加、CacheServerConfig構造体を追加
**目的**: レートリミット設定とキャッシュサーバー設定を読み込むための構造体を追加

**作業内容**:
- `server/internal/config/config.go`の`APIConfig`構造体に`RateLimit RateLimitConfig`フィールドを追加
- `RateLimitConfig`構造体を定義:
  - `Enabled bool`
  - `RequestsPerMinute int`
  - `RequestsPerHour int`（オプション）
- `Config`構造体に`CacheServer CacheServerConfig`フィールドを追加
- `CacheServerConfig`構造体を定義:
  - `Redis RedisConfig`
- `RedisConfig`構造体を定義:
  - `Cluster RedisClusterConfig`
- `RedisClusterConfig`構造体を定義:
  - `Addrs []string`
- `mapstructure`タグを適切に設定

**受け入れ基準**:
- `APIConfig`構造体に`RateLimit`フィールドが追加されている
- `Config`構造体に`CacheServer`フィールドが追加されている
- `RateLimitConfig`、`CacheServerConfig`、`RedisConfig`、`RedisClusterConfig`構造体が定義されている
- `mapstructure`タグが正しく設定されている
- 既存の設定読み込みに影響がない

---

#### - [ ] タスク 2.2: config.goにcacheserver.yamlの読み込み処理を追加
**目的**: キャッシュサーバー設定ファイルを読み込む処理を追加

**作業内容**:
- `server/internal/config/config.go`の`Load()`関数を修正
- `database.yaml`の読み込み後に`cacheserver.yaml`の読み込み処理を追加
- `viper.MergeInConfig()`で`cacheserver.yaml`をマージ
- `cacheserver.yaml`が存在しない場合はエラーにしない（オプショナル）

**受け入れ基準**:
- `Load()`関数に`cacheserver.yaml`の読み込み処理が追加されている
- `cacheserver.yaml`が存在しない場合でもエラーにならない
- 設定が正しく読み込まれる

---

#### - [ ] タスク 2.3: 開発環境設定ファイルにレートリミット設定を追加
**目的**: 開発環境のレートリミット設定を追加

**作業内容**:
- `config/develop/config.yaml`の`api`セクションに`rate_limit`設定を追加:
  ```yaml
  api:
    rate_limit:
      enabled: true
      requests_per_minute: 60
      requests_per_hour: 1000
  ```

**受け入れ基準**:
- `config/develop/config.yaml`に`rate_limit`設定が追加されている
- 設定値が正しい（`enabled: true`、`requests_per_minute: 60`、`requests_per_hour: 1000`）
- YAMLの構文が正しい

---

#### - [ ] タスク 2.4: 開発環境キャッシュサーバー設定ファイルを作成
**目的**: 開発環境のキャッシュサーバー設定ファイルを作成

**作業内容**:
- `config/develop/cacheserver.yaml`を作成
- Redis Clusterの設定を追加（開発環境では空配列）:
  ```yaml
  redis:
    cluster:
      addrs: []
  ```

**受け入れ基準**:
- `config/develop/cacheserver.yaml`が作成されている
- 設定値が正しい（`addrs: []`）
- YAMLの構文が正しい

---

#### - [ ] タスク 2.5: Staging環境設定ファイルにレートリミット設定を追加
**目的**: Staging環境のレートリミット設定を追加

**作業内容**:
- `config/staging/config.yaml`の`api`セクションに`rate_limit`設定を追加:
  ```yaml
  api:
    rate_limit:
      enabled: true
      requests_per_minute: 60
      requests_per_hour: 1000
  ```

**受け入れ基準**:
- `config/staging/config.yaml`に`rate_limit`設定が追加されている
- 設定値が正しい（`enabled: true`、`requests_per_minute: 60`、`requests_per_hour: 1000`）
- YAMLの構文が正しい

---

#### - [ ] タスク 2.6: Staging環境キャッシュサーバー設定ファイルを作成
**目的**: Staging環境のキャッシュサーバー設定ファイルを作成

**作業内容**:
- `config/staging/cacheserver.yaml`を作成
- Redis Clusterの設定を追加:
  ```yaml
  redis:
    cluster:
      addrs:
        - host1:6379
        - host2:6379
        - host3:6379
  ```

**受け入れ基準**:
- `config/staging/cacheserver.yaml`が作成されている
- 設定値が正しい（Redis Clusterのアドレスリスト）
- YAMLの構文が正しい

---

#### - [ ] タスク 2.7: 本番環境設定ファイル例にレートリミット設定を追加
**目的**: 本番環境のレートリミット設定例を追加

**作業内容**:
- `config/production/config.yaml.example`の`api`セクションに`rate_limit`設定を追加:
  ```yaml
  api:
    rate_limit:
      enabled: true
      requests_per_minute: 60
      requests_per_hour: 1000
  ```

**受け入れ基準**:
- `config/production/config.yaml.example`に`rate_limit`設定が追加されている
- 設定値が正しい（`enabled: true`、`requests_per_minute: 60`、`requests_per_hour: 1000`）
- YAMLの構文が正しい

---

#### - [ ] タスク 2.8: 本番環境キャッシュサーバー設定ファイル例を作成
**目的**: 本番環境のキャッシュサーバー設定ファイル例を作成

**作業内容**:
- `config/production/cacheserver.yaml.example`を作成
- Redis Clusterの設定を追加:
  ```yaml
  redis:
    cluster:
      addrs:
        - host1:6379
        - host2:6379
        - host3:6379
  ```

**受け入れ基準**:
- `config/production/cacheserver.yaml.example`が作成されている
- 設定値が正しい（Redis Clusterのアドレスリスト）
- YAMLの構文が正しい

---

### Phase 3: レートリミットミドルウェアの実装

#### - [ ] タスク 3.1: ratelimitパッケージディレクトリの作成
**目的**: レートリミットミドルウェア用のパッケージディレクトリを作成

**作業内容**:
- `server/internal/ratelimit/`ディレクトリを作成
- ディレクトリが正しく作成されていることを確認

**受け入れ基準**:
- `server/internal/ratelimit/`ディレクトリが存在する
- ディレクトリのパーミッションが適切である

---

#### - [ ] タスク 3.2: ストレージ初期化関数の実装
**目的**: 環境に応じたストレージ（In-Memory/Redis Cluster）を初期化する関数を実装

**作業内容**:
- `server/internal/ratelimit/middleware.go`を作成
- `initStore(cfg *config.Config)`関数を実装:
  - `cfg.CacheServer.Redis.Cluster.Addrs`を確認
  - `addrs`が設定されている場合（空でない）: Redis Clusterストレージを初期化（`redis.NewClusterClient`を使用）
  - `addrs`が設定されていない場合（空）: In-Memoryストレージを初期化（`memory.NewStore()`を使用）
- エラーハンドリングを実装

**受け入れ基準**:
- `initStore()`関数が実装されている
- `addrs`が設定されている場合、Redis Clusterストレージが初期化される
- `addrs`が設定されていない場合、In-Memoryストレージが初期化される
- エラーハンドリングが適切に実装されている

---

#### - [ ] タスク 3.3: レートリミットミドルウェア関数の実装
**目的**: Echoミドルウェアとしてレートリミットチェックを実装

**作業内容**:
- `NewRateLimitMiddleware()`関数を実装:
  - レートリミット設定が無効な場合（`enabled: false`または未設定）は、常に許可するミドルウェアを返す
  - レートリミット設定が有効な場合:
    - `limiter.Rate`を設定（`Period: time.Minute`、`Limit: requests_per_minute`）
    - `initStore(cfg)`でストレージを初期化（`cfg`は`config.Config`）
    - エラー時はfail-open方式でリクエストを許可
    - `limiter.New()`でlimiterインスタンスを作成
    - Echoミドルウェア関数を返す
- ミドルウェア関数内で:
  - `echo.Context.RealIP()`でIPアドレスを取得
  - IPアドレスが取得できない場合は許可
  - `instance.Get()`でレートリミットチェック
  - エラー時はfail-open方式でリクエストを許可
  - X-RateLimit-*ヘッダーを設定
  - レートリミット超過時はHTTP 429を返却

**受け入れ基準**:
- `NewRateLimitMiddleware()`関数が実装されている
- レートリミット設定が無効な場合、常に許可するミドルウェアが返される
- レートリミット設定が有効な場合、適切にレートリミットチェックが行われる
- IPアドレスが正しく取得できる
- レートリミット超過時にHTTP 429が返却される
- X-RateLimit-*ヘッダーが正しく設定される
- エラー時はfail-open方式でリクエストが許可される

---

#### - [ ] タスク 3.4: レスポンスヘッダーの実装
**目的**: X-RateLimit-*ヘッダーをレスポンスに付与

**作業内容**:
- ミドルウェア関数内でX-RateLimit-*ヘッダーを設定:
  - `X-RateLimit-Limit`: 制限値（`context.Limit`）
  - `X-RateLimit-Remaining`: 残りリクエスト数（`context.Remaining`）
  - `X-RateLimit-Reset`: リセット時刻（`context.Reset`、Unix timestamp）
- 正常時もレートリミット超過時もヘッダーを付与

**受け入れ基準**:
- X-RateLimit-Limitヘッダーが設定されている
- X-RateLimit-Remainingヘッダーが設定されている
- X-RateLimit-Resetヘッダーが設定されている
- 正常時もレートリミット超過時もヘッダーが付与される

---

#### - [ ] タスク 3.5: エラーハンドリングとログ出力の実装
**目的**: エラー時の適切な処理とログ出力を実装

**作業内容**:
- ストレージ初期化エラー時のログ出力を追加
- レートリミットチェックエラー時のログ出力を追加
- fail-open方式の実装（エラー時はリクエストを許可）
- ログレベルは適切に設定（エラー時はERROR、警告時はWARN）

**受け入れ基準**:
- ストレージ初期化エラー時にログが出力される
- レートリミットチェックエラー時にログが出力される
- エラー時はfail-open方式でリクエストが許可される
- ログレベルが適切に設定されている

---

### Phase 4: ルーターへの統合

#### - [ ] タスク 4.1: router.goにレートリミットミドルウェアを適用
**目的**: レートリミットミドルウェアをルーターに適用

**作業内容**:
- `server/internal/api/router/router.go`を修正
- `ratelimit`パッケージをインポート
- `NewRateLimitMiddleware()`を呼び出してミドルウェアを取得
- エラー時はログに記録し、サーバー起動を継続（fail-open方式）
- ミドルウェアを`e.Use()`で適用
- 適用順序: Recover → CORS → RateLimit → Auth

**受け入れ基準**:
- `ratelimit`パッケージがインポートされている
- `NewRateLimitMiddleware()`が呼び出されている
- ミドルウェアが`e.Use()`で適用されている
- 適用順序が正しい（Recover → CORS → RateLimit → Auth）
- エラー時はログに記録され、サーバー起動が継続される

---

### Phase 5: テスト

#### - [ ] タスク 5.1: レートリミットミドルウェアの単体テスト
**目的**: レートリミットミドルウェアの動作を確認

**作業内容**:
- `server/internal/ratelimit/middleware_test.go`を作成
- テストケース:
  - レートリミット設定が無効な場合、常に許可されることを確認
  - レートリミット内のリクエストが正常に処理されることを確認
  - レートリミット超過時にHTTP 429が返却されることを確認
  - X-RateLimit-*ヘッダーが正しく設定されることを確認
  - IPアドレスが正しく取得できることを確認
  - エラー時のfail-open動作を確認

**受け入れ基準**:
- テストファイルが作成されている
- すべてのテストケースが実装されている
- テストが正常に実行される
- すべてのテストがパスする

---

#### - [ ] タスク 5.2: ストレージの統合テスト
**目的**: In-MemoryストレージとRedis Clusterストレージの動作を確認

**作業内容**:
- ストレージの統合テストを作成
- テストケース:
  - `cacheserver.yaml`で`addrs`が空の場合、In-Memoryストレージが使用されることを確認
  - `cacheserver.yaml`で`addrs`が設定されている場合、Redis Clusterストレージが使用されることを確認（統合テスト環境が必要）
  - ストレージの切り替えが正しく動作することを確認

**受け入れ基準**:
- 統合テストが作成されている
- すべてのテストケースが実装されている
- テストが正常に実行される
- すべてのテストがパスする

---

### Phase 6: 動作確認

#### - [ ] タスク 6.1: 開発環境での動作確認
**目的**: 開発環境でレートリミット機能が正常に動作することを確認

**作業内容**:
- 開発環境でサーバーを起動
- レートリミット内のリクエストが正常に処理されることを確認
- レートリミット超過時にHTTP 429が返却されることを確認
- X-RateLimit-*ヘッダーが正しく付与されることを確認
- レートリミットが時間経過でリセットされることを確認
- In-Memoryストレージが使用されることを確認（`cacheserver.yaml`で`addrs`が空の場合）

**受け入れ基準**:
- サーバーが正常に起動する
- レートリミット内のリクエストが正常に処理される
- レートリミット超過時にHTTP 429が返却される
- X-RateLimit-*ヘッダーが正しく付与される
- レートリミットが時間経過でリセットされる
- In-Memoryストレージが使用される

---

#### - [ ] タスク 6.2: 既存機能への影響確認
**目的**: 既存のAPI動作に影響がないことを確認

**作業内容**:
- 既存のAPIエンドポイント（User、Post、Today）が正常に動作することを確認
- 認証機能が正常に動作することを確認
- アクセス制御が正常に動作することを確認
- 既存のAPIクライアントへの影響がないことを確認

**受け入れ基準**:
- 既存のAPIエンドポイントが正常に動作する
- 認証機能が正常に動作する
- アクセス制御が正常に動作する
- 既存のAPIクライアントへの影響がない

---

#### - [ ] タスク 6.3: Redis Cluster環境での動作確認（オプション）
**目的**: Redis Cluster環境でレートリミット機能が正常に動作することを確認

**作業内容**:
- Redis Cluster環境を構築（Docker等を使用）
- `config/{env}/cacheserver.yaml`にRedis Clusterのアドレスを設定
- サーバーを起動
- レートリミット機能が正常に動作することを確認
- Redis Clusterストレージが使用されることを確認

**受け入れ基準**:
- Redis Cluster環境が構築されている
- `cacheserver.yaml`にRedis Clusterのアドレスが設定されている
- サーバーが正常に起動する
- レートリミット機能が正常に動作する
- Redis Clusterストレージが使用される

---

### Phase 7: ドキュメント更新

#### - [ ] タスク 7.1: READMEまたはドキュメントの更新
**目的**: レートリミット機能の使用方法をドキュメントに記載

**作業内容**:
- READMEまたは関連ドキュメントにレートリミット機能の説明を追加
- 設定方法、環境変数の説明を追加
- 動作確認方法を追加

**受け入れ基準**:
- ドキュメントが更新されている
- 設定方法が記載されている
- 環境変数の説明が記載されている
- 動作確認方法が記載されている

---

## 実装順序の推奨

1. **Phase 1**: 依存関係の追加（必須）
2. **Phase 2**: 設定管理の実装（必須）
3. **Phase 3**: レートリミットミドルウェアの実装（必須）
4. **Phase 4**: ルーターへの統合（必須）
5. **Phase 5**: テスト（推奨）
6. **Phase 6**: 動作確認（必須）
7. **Phase 7**: ドキュメント更新（推奨）

## 注意事項

- 各タスクは独立して実装可能な粒度に分解
- タスクの実装順序は推奨順序に従うことを推奨
- テストは実装と並行して進めることを推奨
- 動作確認は各フェーズの完了後に実施することを推奨
- Redis Cluster環境での動作確認はオプション（統合テスト環境が必要）
