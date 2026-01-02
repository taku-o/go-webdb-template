# Redis遅延接続・自動再接続機能要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #74
- **Issueタイトル**: Redisに遅延接続＋Redisダウン後の再開時に再接続するように修正したい
- **Feature名**: 0037-redis-reconnect
- **作成日**: 2025-01-27

### 1.2 目的
`github.com/redis/go-redis/v9`の標準機能を活用して、Redisへの遅延接続とダウン後の自動再接続を実現する。
これにより、APIサーバー起動時にRedisが利用できない場合でもサーバーを起動でき、Redisが復旧した際に自動的に接続が確立されるようにする。

### 1.3 スコープ
- Redis接続オプションの設定追加（`MaxRetries`, `MinRetryBackoff`, `MaxRetryBackoff`, `DialTimeout`, `ReadTimeout`, `PoolSize`, `PoolTimeout`など）
- 遅延接続と自動再接続の動作確認（Redis Clusterを使用）
- Redis接続エラー時のリトライ機能の確認（標準機能の活用）
- Redis環境の構築（Dockerを使用、既存の`docker-compose.redis.yml`を活用）

**本実装の範囲外**:
- 本番環境での本格的な接続管理機能
- 接続監視・アラート機能
- 接続プールの動的な調整機能
- 複数Redisインスタンスへの接続管理の最適化

## 2. 背景・現状分析

### 2.1 現在の実装
- **Redis接続**: `server/internal/ratelimit/middleware.go`で`redis.NewClusterClient`を使用
- **接続オプション**: `ClusterOptions`に`Addrs`のみが設定されている（148-150行目、163-165行目）
- **起動時の接続確認**: 起動時のPing処理は実装されていない（遅延接続は既に実現されている可能性）
- **接続オプション設定**: 自動再接続・リトライの設定が実装されていない
- **Redis環境**: `docker-compose.redis.yml`が存在（0035-jobqueueで実装済み）

### 2.2 課題点
1. **接続オプションの不足**: 自動再接続・リトライの設定が実装されていない
2. **動作確認の不足**: 遅延接続と自動再接続が正しく動作するか確認されていない
3. **接続エラー時のリトライ設定の不足**: 接続エラー時のリトライ設定が実装されていない
4. **接続タイムアウト設定の不足**: 接続タイムアウトや読み取りタイムアウトの設定が実装されていない

### 2.3 本実装による改善点
1. **遅延接続の確認**: `redis.NewClusterClient`はデフォルトで遅延接続をサポートしていることを確認
2. **自動再接続の実現**: 適切な接続オプション設定により、Redisが復旧した際に自動的に再接続される
3. **接続オプションの最適化**: 適切な接続オプション設定により、パフォーマンスと安定性が向上
4. **接続エラー時のリトライ**: 接続エラー時に自動的にリトライする機能を実装
5. **動作確認環境の整備**: 既存のRedis環境を活用し、動作確認が可能になる

## 3. 機能要件

### 3.1 Redis接続オプションの設定追加

#### 3.1.1 接続オプションの確認
- `server/internal/ratelimit/middleware.go`の`initStore`関数で`redis.NewClusterClient`の接続オプションを確認
- 現在は`Addrs`のみが設定されていることを確認
- 接続オプションが設定ファイルから読み込めるか確認

#### 3.1.2 接続オプションの追加
以下の設定を`ClusterOptions`に追加する：
- `MaxRetries`: コマンド失敗時の最大リトライ数（デフォルト: 2）
- `MinRetryBackoff`: リトライ間隔（最小）（デフォルト: 8ms）
- `MaxRetryBackoff`: リトライ間隔（最大）（デフォルト: 512ms）
- `DialTimeout`: 接続確立のタイムアウト（デフォルト: 5秒）
- `ReadTimeout`: 読み取りタイムアウト（デフォルト: 3秒）
- `PoolSize`: 接続プールサイズ（デフォルト: CPU数×10）
- `PoolTimeout`: プールから接続を取り出す際の待機時間（デフォルト: 4秒）

#### 3.1.3 設定値の確認
- 設定ファイル（`config/{env}/cacheserver.yaml`）から接続オプション設定を読み込む
- デフォルト値が適切に設定されていることを確認
- 設定値が0以下の場合は適切なデフォルト値を設定

### 3.2 Redis環境の確認

#### 3.1.1 既存のDocker Compose設定の確認
- `docker-compose.redis.yml`ファイルが存在することを確認（0035-jobqueueで実装済み）
- Redisコンテナの定義を確認
- ネットワーク設定を確認
- ボリューム設定（データ永続化用）を確認
- デフォルトポート: 6379

#### 3.1.2 起動スクリプトの確認
- `scripts/start-redis.sh`が存在することを確認（0035-jobqueueで実装済み）
- Docker Composeを使用してRedisを起動することを確認
- 既存の起動スクリプト（`start-mailpit.sh`など）と同じパターンで実装されていることを確認

#### 3.1.3 設定ファイルの確認
- Redis用の設定ファイル（`config/{env}/cacheserver.yaml`）が存在することを確認
- Redis接続情報（Cluster Addrsなど）が設定されていることを確認

### 3.3 遅延接続と自動再接続の動作確認

#### 3.3.1 遅延接続の確認
- Redis環境でAPIサーバーを起動（Redis接続なしで起動可能であることを確認）
- 最初のRedisコマンド実行時に接続が確立されることを確認
- ログで接続確立のタイミングを確認

#### 3.3.2 自動再接続の確認
- Redisを停止して、Redisコマンド実行時にエラーが発生することを確認
- Redisを再起動して、次のRedisコマンド実行時に自動的に再接続されることを確認
- ログで再接続のタイミングを確認

### 3.4 Redis接続エラー時のリトライ機能の確認

#### 3.4.1 リトライ機能の確認
- `github.com/redis/go-redis/v9`の標準機能であるリトライ機能が動作することを確認
- `MaxRetries`設定により、接続エラー時にリトライが実行されることを確認
- リトライ間隔が`MinRetryBackoff`と`MaxRetryBackoff`の範囲内であることを確認

#### 3.4.2 リトライ設定の確認
- リトライ回数: 最大2回（初回 + 2回のリトライ）
- リトライ間隔: 8ms～512msの範囲で指数バックオフ
- リトライ対象: Redis接続エラー（ネットワークエラー、タイムアウトエラーなど）

## 4. 非機能要件

### 4.1 パフォーマンス
- 接続プール設定により、接続の再利用が適切に行われること
- リトライ処理が過度にパフォーマンスに影響を与えないこと

### 4.2 可用性
- Redis接続できない場合でもサーバーを起動できること
- Redisが復旧した際に自動的に再接続されること
- 接続エラー時に適切なリトライが実行されること

### 4.3 保守性
- 既存のDocker Compose設定パターンに従うこと
- 既存の起動スクリプトパターンに従うこと
- 接続オプション設定が設定ファイルから読み込めること

### 4.4 拡張性
- 将来的に接続監視・アラート機能を追加できる設計であること
- 接続プール設定を動的に調整できる設計であること（将来の拡張項目）

## 5. 技術仕様

### 5.1 サーバー側技術スタック
- **言語**: Go 1.21+
- **Redisクライアント**: `github.com/redis/go-redis/v9`
- **Redis**: Redis Cluster（既存実装）
- **Webフレームワーク**: 既存のEcho + Huma API

### 5.2 インフラストラクチャ
- **Redis**: Docker Composeで起動（既存の`docker-compose.redis.yml`を活用）
- **Docker Compose**: 既存のパターンに従う

### 5.3 ファイル構造
- **Docker Compose設定**: `docker-compose.redis.yml`（既存、確認のみ）
- **起動スクリプト**: `scripts/start-redis.sh`（既存、確認のみ）
- **Redis接続処理**: `server/internal/ratelimit/middleware.go`（修正）
- **設定**: `server/internal/config/config.go`（確認・修正）
- **設定ファイル**: `config/{env}/cacheserver.yaml`（確認・修正）

## 6. 受け入れ基準

### 6.1 機能要件
1. **Redis接続オプションの設定追加**: 接続オプションが適切に実装されていること
2. **Redis環境の確認**: 既存のRedis環境が適切に動作していること
3. **遅延接続の確認**: 最初のRedisコマンド実行時に接続が確立されること
4. **自動再接続の確認**: Redisが復旧した際に自動的に再接続されること
5. **接続リトライ機能**: Redis接続エラー時に8ms～512msの範囲で指数バックオフして最大2回までリトライすること

### 6.2 非機能要件
1. **パフォーマンス**: 接続プール設定により、接続の再利用が適切に行われること
2. **可用性**: Redis接続できない場合でもサーバーを起動できること
3. **保守性**: 既存のパターンに従った実装であること
4. **拡張性**: 将来的に接続監視・アラート機能を追加できる設計であること

## 7. 制約事項

1. **動作確認環境**: 既存のRedis環境（`docker-compose.redis.yml`）を使用する
2. **開発環境での利用**: 本実装は主に開発環境での利用を想定（本番環境での利用は将来の拡張項目）
3. **既存機能の保持**: 既存のRedis接続機能を壊さないこと
4. **後方互換性**: 既存の設定ファイルとの互換性を保つこと

## 8. 将来の拡張項目（現時点では未実装）

以下の機能は将来の拡張として検討されていますが、現時点では実装対象外です：

- 接続監視・アラート機能
- 接続プールの動的な調整機能
- 複数Redisインスタンスへの接続管理の最適化
- 本番環境での本格的な接続管理機能

## Project Description (Input)

https://github.com/taku-o/go-webdb-template/issues/74 に対応するための要件を作成してください。

Redisでも、アプリ起動時に接続しないで、必要になった時に接続を作る。
かつ、Redisがダウンした後、Redisが再起動した時に接続するように修正したい。

調べたところ、次の調査結果が出た。
```
1. 接続の挙動：デフォルトで「遅延接続」
redis.NewClient を実行した瞬間には、実はまだRedisサーバーに接続していません。
内部的に接続プールが作成されるだけで、実際のTCP接続は、
最初のコマンド（Ping, Get, Set など）が実行されたタイミングで初めて行われます。

2. 自動再接続：標準でサポート
import "github.com/redis/go-redis/v9"

rdb := redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "", 
    DB:       0,

    // --- 自動再接続・リトライの設定 ---
    MaxRetries:      2,                      // コマンド失敗時の最大リトライ数
    MinRetryBackoff: 8 * time.Millisecond,   // リトライ間隔（最小）
    MaxRetryBackoff: 512 * time.Millisecond, // リトライ間隔（最大）
    
    // --- 接続管理の設定 ---
    DialTimeout:  5 * time.Second,  // 接続確立のタイムアウト
    ReadTimeout:  3 * time.Second,  // 読み取りタイムアウト
    PoolSize:     10,               // CPU数に応じた適切なプールサイズ
    PoolTimeout:  4 * time.Second,  // プールから接続を取り出す際の待機時間
})
```

## Requirements

### Requirement 1: Redis接続オプションの設定追加
**Objective:** As a system, I want proper Redis connection options to be configured, so that lazy connection and auto-reconnection can work correctly with appropriate retry settings.

#### Acceptance Criteria
1. WHEN Redis connection options are checked THEN they SHALL be implemented in initStore function
2. IF connection options are checked THEN they SHALL include MaxRetries, MinRetryBackoff, MaxRetryBackoff, DialTimeout, ReadTimeout, PoolSize, PoolTimeout
3. WHERE connection options are read THEN they SHALL be read from config file (config/{env}/cacheserver.yaml)
4. WHEN connection options are not configured THEN default values SHALL be used
5. IF connection options are 0 or negative THEN appropriate default values SHALL be set
6. WHERE ClusterOptions is configured THEN it SHALL include all necessary connection options for auto-reconnection and retry

### Requirement 2: Redis環境の確認
**Objective:** As a developer, I want to verify that existing Redis environment is properly set up, so that I can test lazy connection and auto-reconnection features.

#### Acceptance Criteria
1. WHEN docker-compose.redis.yml is checked THEN it SHALL exist and define Redis service
2. IF Redis service is defined THEN it SHALL use appropriate Docker image
3. WHERE Redis is configured THEN it SHALL listen on port 6379 by default
4. WHEN Redis volumes are defined THEN they SHALL persist data across container restarts
5. IF start-redis.sh is checked THEN it SHALL start Redis using Docker Compose
6. WHERE start-redis.sh is executed THEN it SHALL follow existing script patterns
7. WHEN Redis configuration is checked THEN it SHALL be in config/{env}/cacheserver.yaml
8. IF Redis configuration is checked THEN it SHALL include Cluster Addrs

### Requirement 3: 遅延接続の動作確認
**Objective:** As a developer, I want to verify that lazy connection works correctly, so that I can confirm connections are established only when needed.

#### Acceptance Criteria
1. WHEN API server starts without Redis connection THEN it SHALL start successfully
2. IF first Redis command is executed THEN the connection SHALL be established automatically
3. WHERE connection is established THEN it SHALL be logged appropriately
4. WHEN lazy connection works THEN the connection timing SHALL be verified in logs

### Requirement 4: 自動再接続の動作確認
**Objective:** As a developer, I want to verify that auto-reconnection works correctly, so that I can confirm connections are automatically re-established when Redis becomes available.

#### Acceptance Criteria
1. WHEN Redis is stopped THEN Redis commands SHALL fail with appropriate errors
2. IF Redis is restarted THEN the next Redis command SHALL automatically re-establish connection
3. WHERE auto-reconnection works THEN the reconnection timing SHALL be verified in logs
4. WHEN Redis becomes available THEN connections SHALL be automatically re-established

### Requirement 5: Redis接続エラー時のリトライ機能の確認
**Objective:** As a system, I want Redis connection errors to be retried automatically, so that temporary connection failures can be recovered without manual intervention.

#### Acceptance Criteria
1. WHEN retry functionality is verified THEN it SHALL use github.com/redis/go-redis/v9 standard retry feature
2. IF Redis connection fails THEN it SHALL retry up to 2 times (initial attempt + 2 retries)
3. WHERE retry is executed THEN it SHALL wait between 8ms and 512ms with exponential backoff
4. WHEN retry is configured THEN it SHALL be configured via MaxRetries, MinRetryBackoff, MaxRetryBackoff
5. IF retry is executed THEN it SHALL handle connection errors (network errors, timeout errors, etc.)
6. WHERE all retries fail THEN it SHALL return appropriate error message
