# APIレートリミット機能要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #15
- **Issueタイトル**: APIサーバーへのリクエストにrate limit機能をつける
- **Feature名**: 0022-apilimit
- **作成日**: 2025-01-27

### 1.2 目的
APIの過剰な呼び出しを防止するために、レートリミット機能を導入する。IPアドレス単位でリクエスト数を制限し、制限を超過した場合は適切なHTTPステータスコードとヘッダーを返却する。開発環境ではIn-Memoryストレージ、staging/本番環境ではRedisを使用してレートリミットを管理する。

### 1.3 スコープ
- IPアドレス単位でのレートリミット実装
- 環境別ストレージ（開発: In-Memory、staging/production: Redis）
- 設定ファイルによる閾値管理
- レートリミット超過時の適切なレスポンス（HTTP 429、X-RateLimit-*ヘッダー）
- Public APIのみを制限（または/apiエンドポイント全体を制限）

**本実装の範囲外**:
- 認証ロジックの変更（既存の認証ミドルウェアは変更しない）
- APIエンドポイントの追加・削除
- アクセス制御ロジックの変更
- ユーザー単位のレートリミット（IPアドレスのみ）

## 2. 背景・現状分析

### 2.1 現在の実装
- **APIフレームワーク**: Echo + Huma v2を使用
- **認証機能**: Issue #19（0019-auth0-apicall）で実装済み
  - Auth0 JWTとPublic API Key JWTの両方をサポート
  - 認証ミドルウェア（`server/internal/auth/middleware.go`）でJWT検証を実施
  - `/api/`で始まるパスに認証を適用
- **APIエンドポイント**: 
  - Userエンドポイント（5つ）：全てpublicレベル
  - Postエンドポイント（6つ）：全てpublicレベル
  - Todayエンドポイント（1つ）：privateレベル
- **設定管理**: `server/internal/config/config.go`で設定を管理
- **ミドルウェア**: `server/internal/api/router/router.go`でEchoミドルウェアを適用

### 2.2 課題点
1. **レートリミット未実装**: APIの過剰な呼び出しを防止する機能がない
2. **DoS攻撃への脆弱性**: 無制限にAPIを呼び出すことが可能
3. **リソース保護の不足**: サーバーリソースを過剰に消費するリクエストを制限できない
4. **クライアントへの情報提供不足**: レートリミットの状態をクライアントに通知する仕組みがない

### 2.3 本実装による改善点
1. **過剰な呼び出しの防止**: IPアドレス単位でリクエスト数を制限
2. **DoS攻撃への対策**: レートリミットによりDoS攻撃の影響を軽減
3. **リソース保護**: サーバーリソースを適切に保護
4. **クライアントへの情報提供**: X-RateLimit-*ヘッダーにより制限状況を通知

## 3. 機能要件

### 3.1 レートリミットミドルウェアの実装

#### 3.1.1 基本機能
- IPアドレス単位でリクエスト数をカウント
- 設定された閾値を超過した場合、リクエストを拒否
- レートリミットの状態をレスポンスヘッダーで通知

#### 3.1.2 実装方法
- ulule/limiterライブラリを利用（Issue #15の記載に従う）
- Echoミドルウェアとして実装
- 認証ミドルウェアの前または後に適用（設計フェーズで決定）

#### 3.1.3 実装場所
- 新規ファイル: `server/internal/ratelimit/middleware.go`
- 統合場所: `server/internal/api/router/router.go`

### 3.2 IPアドレス単位の制限

#### 3.2.1 IPアドレスの取得
- リクエストのRemoteAddrからIPアドレスを取得
- X-Forwarded-Forヘッダーを考慮（プロキシ経由の場合）
- IPv4/IPv6の両方に対応

#### 3.2.2 制限の適用範囲
- Public APIのみを制限（推奨）
  - 各エンドポイントのアクセスレベル（public/private）を確認
  - publicレベルのエンドポイントのみにレートリミットを適用
- または、/apiエンドポイント全体を制限（実装が複雑になる場合はこちらを採用）

### 3.3 環境別ストレージ

#### 3.3.1 開発環境（develop）
- **ストレージ**: In-Memory（memory.NewStore()）
- **用途**: ローカル開発、テスト
- **特徴**: サーバー再起動でリセットされる

#### 3.3.2 Staging環境（staging）
- **ストレージ**: Redis Cluster
- **用途**: ステージング環境での動作確認
- **設定**: `config/staging/cacheserver.yaml`でRedis Clusterのアドレスを指定

#### 3.3.3 本番環境（production）
- **ストレージ**: Redis Cluster（AWS ElastiCache等を想定）
- **用途**: 本番環境でのレートリミット管理
- **設定**: `config/production/cacheserver.yaml`でRedis Clusterのアドレスを指定
- **実装**: `redis.NewClusterClient`を使用

#### 3.3.4 ストレージの切り替え
- 環境変数`APP_ENV`に基づいて自動的に切り替え
- `cacheserver.yaml`でRedis Clusterのアドレスが設定されている場合: Redis Clusterとして接続（`redis.NewClusterClient`を使用）
- `cacheserver.yaml`でRedis Clusterのアドレスが設定されていない場合: In-Memoryストレージを使用

### 3.4 設定ファイルによる閾値管理

#### 3.4.1 設定項目
- **requests_per_minute**: 1分あたりのリクエスト数上限
- **requests_per_hour**: 1時間あたりのリクエスト数上限（オプション）
- **enabled**: レートリミット機能の有効/無効

#### 3.4.2 設定ファイルの構造
```yaml
api:
  rate_limit:
    enabled: true
    requests_per_minute: 60
    requests_per_hour: 1000  # オプション
```

#### 3.4.3 設定の読み込み
- `server/internal/config/config.go`の`APIConfig`構造体に追加
- 環境別設定ファイル（`config/develop/config.yaml`等）で管理
- 開発中のサービスなので、閾値は緩く設定（テスト時は一時的に低く設定可能）

### 3.5 レスポンス仕様

#### 3.5.1 レートリミット超過時
- **HTTPステータスコード**: `429 Too Many Requests`
- **レスポンスボディ**: エラーメッセージを含むJSON形式
```json
{
  "code": 429,
  "message": "Too Many Requests"
}
```

#### 3.5.2 X-RateLimit-*ヘッダー
以下のヘッダーをすべてのレスポンスに付与：
- **X-RateLimit-Limit**: 制限値（例: `60`）
- **X-RateLimit-Remaining**: 残りリクエスト数（例: `45`）
- **X-RateLimit-Reset**: リセット時刻（Unix timestamp、例: `1706342400`）

#### 3.5.3 正常時のレスポンス
- レートリミット内の場合、通常通りAPIレスポンスを返却
- X-RateLimit-*ヘッダーは常に付与（クライアントが制限状況を把握できるように）

## 4. 非機能要件

### 4.1 パフォーマンス
- レートリミットチェックによるレスポンス時間への影響を最小化
- Redis使用時は接続プールを適切に管理
- In-Memory使用時はメモリ使用量を監視

### 4.2 互換性
- **既存機能への影響なし**: 認証ロジック、API動作、アクセス制御は一切変更しない
- **後方互換性の維持**: 既存のAPIクライアントへの影響なし
- **設定の後方互換性**: レートリミット設定が未指定の場合は機能を無効化

### 4.3 メンテナンス性
- **設定の一元管理**: レートリミット設定は`config.yaml`で一元管理
- **明確な実装**: ミドルウェアの実装が明確で理解しやすい
- **ログ出力**: レートリミット超過時のログ出力（オプション）

### 4.4 可用性
- Redis接続エラー時は適切にエラーハンドリング
- In-Memoryストレージへのフォールバック（オプション、設計フェーズで決定）

## 5. 制約事項

### 5.1 実装範囲の制約
- **Public APIのみ制限**: 各エンドポイントのアクセスレベルを確認してpublicのみに適用
- **または/apiエンドポイント全体を制限**: 実装が複雑になる場合はこちらを採用
- **IPアドレスのみ**: ユーザー単位のレートリミットは実装しない

### 5.2 技術的制約
- **ulule/limiterライブラリ**: Issue #15の記載に従い、利用可能なら利用
- **Echoフレームワーク**: 既存のEchoミドルウェアとして実装
- **Go言語**: 既存のGo言語実装を維持

### 5.3 設定の制約
- **閾値は緩く設定**: 開発中のサービスなので、テストしやすい閾値を設定
- **設定ファイルで管理**: 環境変数ではなく設定ファイルで管理（既存の設定管理方針に従う）

### 5.4 環境の制約
- **開発環境**: In-Memoryストレージのみ（Redis不要）
- **Staging/Production**: Redisが必要（Dockerで一時的に建てても良い）

## 6. 受け入れ基準

### 6.1 レートリミットミドルウェアの実装
- [ ] `server/internal/ratelimit/middleware.go`にレートリミットミドルウェアが実装されている
- [ ] ulule/limiterライブラリが使用されている
- [ ] Echoミドルウェアとして正しく実装されている

### 6.2 IPアドレス単位の制限
- [ ] IPアドレスが正しく取得できている（RemoteAddr、X-Forwarded-For対応）
- [ ] 同一IPアドレスからのリクエストが正しくカウントされている
- [ ] 異なるIPアドレスからのリクエストが独立してカウントされている

### 6.3 環境別ストレージ
- [ ] 開発環境でIn-Memoryストレージが使用されている
- [ ] Staging/Production環境でRedis Clusterが使用されている（`cacheserver.yaml`で設定時）
- [ ] 環境に応じて適切にストレージが切り替わる
- [ ] `cacheserver.yaml`でRedis Clusterのアドレスが未設定時はIn-Memoryストレージが使用される

### 6.4 設定ファイルによる閾値管理
- [ ] `config.yaml`にレートリミット設定が追加されている
- [ ] `server/internal/config/config.go`の`APIConfig`構造体に設定項目が追加されている
- [ ] 設定値が正しく読み込まれ、ミドルウェアで使用されている

### 6.5 レスポンス仕様
- [ ] レートリミット超過時にHTTP 429が返却される
- [ ] レスポンスボディに適切なエラーメッセージが含まれる
- [ ] X-RateLimit-Limitヘッダーが付与される
- [ ] X-RateLimit-Remainingヘッダーが付与される
- [ ] X-RateLimit-Resetヘッダーが付与される
- [ ] 正常時もX-RateLimit-*ヘッダーが付与される

### 6.6 動作確認
- [ ] レートリミット内のリクエストが正常に処理される
- [ ] レートリミット超過時に429エラーが返却される
- [ ] レートリミットが時間経過でリセットされる
- [ ] 開発環境でIn-Memoryストレージが動作する
- [ ] Staging/Production環境でRedisが動作する（設定されている場合）

### 6.7 既存機能の動作確認
- [ ] 既存のAPI動作に影響がないことを確認（認証、アクセス制御が正常に動作）
- [ ] 既存のAPIクライアントへの影響がないことを確認

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### サーバー側（Go）
- `server/internal/api/router/router.go` - レートリミットミドルウェアの適用
- `server/internal/config/config.go` - レートリミット設定とキャッシュサーバー設定の追加、`cacheserver.yaml`の読み込み

#### 設定ファイル
- `config/develop/config.yaml` - 開発環境のレートリミット設定
- `config/staging/config.yaml` - Staging環境のレートリミット設定
- `config/production/config.yaml.example` - 本番環境のレートリミット設定例
- `config/develop/cacheserver.yaml` - 開発環境のキャッシュサーバー設定（新規）
- `config/staging/cacheserver.yaml` - Staging環境のキャッシュサーバー設定（新規）
- `config/production/cacheserver.yaml.example` - 本番環境のキャッシュサーバー設定例（新規）

### 7.2 新規ファイル
- `server/internal/ratelimit/middleware.go` - レートリミットミドルウェアの実装
- `.kiro/specs/0022-apilimit/requirements.md` - 本要件定義書
- `.kiro/specs/0022-apilimit/spec.json` - 仕様書メタデータ

### 7.3 変更なしのファイル
- `server/internal/auth/middleware.go` - 認証ミドルウェア（変更なし）
- その他のAPIハンドラー（変更なし）

### 7.4 依存関係の追加
- `go.mod`に`github.com/ulule/limiter/v3`を追加

## 8. 実装上の注意事項

### 8.1 ライブラリの選択
- ulule/limiterライブラリを利用（Issue #15の記載に従う）
- Redis Cluster対応: `redis.NewClusterClient`を使用してRedis Clusterに接続
- `config/{env}/cacheserver.yaml`でRedis Clusterのアドレスを指定（配列形式、例: `["host1:6379", "host2:6379", "host3:6379"]`）
- `cacheserver.yaml`でRedis Clusterのアドレスが設定されていない場合はIn-Memoryストレージを使用
- 利用できない場合は、echo-essentialのレートリミット機能を検討
- どちらも利用できない場合は、独自実装を検討（設計フェーズで決定）

### 8.2 ミドルウェアの適用順序
- 認証ミドルウェアの前または後に適用（設計フェーズで決定）
- レートリミットチェックは認証チェックの前に行うことを推奨（認証前に制限）

### 8.3 IPアドレスの取得
- `echo.Context.RealIP()`を使用してIPアドレスを取得（Echoの標準機能）
- X-Forwarded-Forヘッダーを適切に処理

### 8.4 エラーハンドリング
- Redis接続エラー時の適切な処理
- Redis Cluster接続エラー時の適切な処理（一部ノードの障害時も動作継続）
- レートリミットチェック時のエラー（Redisサーバーと通信できない場合など）はログに記録し、リクエストは許可する（fail-open方式）
  - レートリミット機能は補助的な機能のため、エラー時はAPIへのリクエストを許可する

### 8.5 設定のデフォルト値
- レートリミット設定が未指定の場合は機能を無効化
- 設定ファイルで`enabled: false`が指定されている場合も機能を無効化

### 8.6 動作確認
- 実装後は以下を確認：
  - レートリミット内のリクエストが正常に処理される
  - レートリミット超過時に429エラーが返却される
  - X-RateLimit-*ヘッダーが正しく付与される
  - 開発環境でIn-Memoryストレージが動作する
  - Staging/Production環境でRedis Clusterが動作する（`cacheserver.yaml`で設定時）
  - `cacheserver.yaml`でRedis Clusterのアドレスが未設定時はIn-Memoryストレージが動作する
  - 既存のAPI動作に影響がない

## 9. 将来のRedis移行ステップ

### 9.1 移行手順
1. Redis Cluster（AWS ElastiCache等）を構築
2. `config/{env}/cacheserver.yaml`にRedis Clusterのアドレスを設定（配列形式、例: `["host1:6379", "host2:6379", "host3:6379"]`）
3. `limiter/v3/drivers/store/redis`を有効化
4. `redis.NewClusterClient`を使用してRedis Clusterに接続

### 9.2 移行時の注意事項
- 既存のIn-Memory実装との互換性を維持
- 環境変数による自動切り替えを実装
- 移行時の動作確認を実施

## 10. 参考情報

### 10.1 関連Issue
- GitHub Issue #15: APIサーバーへのリクエストにrate limit機能をつける（0022-apilimit）
- GitHub Issue #19: Auth0から受け取ったJWTをAPIサーバーとの通信で利用する（0019-auth0-apicall）

### 10.2 既存ドキュメント
- `.kiro/specs/0019-auth0-apicall/requirements.md`: Auth0 API呼び出し機能の要件定義書
- `.kiro/specs/0013-echo-huma/requirements.md`: Echo + Huma導入の要件定義書

### 10.3 技術スタック
- **Go言語**: 1.24+
- **Echoフレームワーク**: v4
- **Humaフレームワーク**: v2
- **ulule/limiter**: v3（予定）
- **Redis**: Redis Cluster対応（`redis.NewClusterClient`を使用）

### 10.4 参考資料
- [ulule/limiter Documentation](https://github.com/ulule/limiter)
- [Echo Middleware](https://echo.labstack.com/docs/middleware)
- Issue #15のコメントに記載されている実装要件
