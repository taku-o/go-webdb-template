# GoAdmin死活監視エンドポイントの設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、Adminサーバーに`/health`エンドポイントを実装するための詳細設計を定義する。APIサーバーと同様のシンプルな実装を行い、docker-compose.admin.ymlのヘルスチェック機能を有効化する。

### 1.2 設計の範囲
- Adminサーバーに`/health`エンドポイントを実装する
- エンドポイントは認証不要でアクセス可能とする
- 単体テストと統合テストを実装する
- docker-compose.admin.ymlのヘルスチェック機能が正常に動作することを確認する

### 1.3 設計方針
- **シンプルな実装**: APIサーバーと同様のシンプルな実装を維持する
- **認証不要**: ヘルスチェック用のため、認証ミドルウェアを通過しない
- **一貫性**: APIサーバーの実装パターンに合わせる
- **テスト**: 単体テストと統合テストを実装する

## 2. アーキテクチャ設計

### 2.1 エンドポイントの配置

```
┌─────────────────────────────────────────────────────────────┐
│                    Adminサーバー (Port 8081)                  │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │    Gorilla Mux Router            │
        │    (app := mux.NewRouter())      │
        └─────────────────────────────────┘
                          │
        ┌─────────────────┴─────────────────┐
        │                                   │
        ▼                                   ▼
┌──────────────────┐            ┌──────────────────┐
│  /health         │            │  /admin/*         │
│  (認証不要)       │            │  (GoAdmin Engine)  │
│  - GET           │            │  - 認証必要        │
│  - 200 OK        │            │  - 各種ページ      │
│  - "OK"          │            │                   │
└──────────────────┘            └──────────────────┘
```

### 2.2 リクエストフロー

```
┌─────────────────────────────────────────────────────────────┐
│              /health エンドポイントのリクエストフロー           │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  HTTPリクエスト: GET /health     │
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  Gorilla Mux Router             │
        │  - ルーティング処理               │
        │  - 認証ミドルウェアを通過しない    │
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  ハンドラー関数                   │
        │  - 200 OKを返す                  │
        │  - "OK"を返す                   │
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  HTTPレスポンス                  │
        │  - Status: 200 OK               │
        │  - Body: "OK"                   │
        │  - Content-Type: text/plain      │
        └─────────────────────────────────┘
```

### 2.3 ミドルウェアチェーン

```
┌─────────────────────────────────────────────────────────────┐
│              Adminサーバーのミドルウェアチェーン                │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  HTTPリクエスト                   │
        └─────────────────────────────────┘
                          │
        ┌─────────────────┴─────────────────┐
        │                                   │
        ▼                                   ▼
┌──────────────────┐            ┌──────────────────┐
│  /health         │            │  /admin/*         │
│  (直接ルーティング)│            │  (GoAdmin Engine)  │
│                  │            │                   │
│  ミドルウェアなし  │            │  ┌─────────────┐  │
│                  │            │  │ 認証         │  │
│  ハンドラー直接    │            │  └─────────────┘  │
│                  │            │  ┌─────────────┐  │
│  200 OK + "OK"   │            │  │ アクセスログ │  │
│                  │            │  └─────────────┘  │
└──────────────────┘            │  ┌─────────────┐  │
                                │  │ GoAdmin     │  │
                                │  │ ページ処理   │  │
                                │  └─────────────┘  │
                                └──────────────────┘
```

## 3. 実装設計

### 3.1 エンドポイント実装

#### 3.1.1 実装場所
- **ファイル**: `server/cmd/admin/main.go`
- **実装位置**: GoAdmin Engineの初期化後、カスタムページの登録前

#### 3.1.2 実装コード

```go
// Health check endpoint (認証不要)
app.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "text/plain")
    w.Write([]byte("OK"))
}).Methods("GET")
```

#### 3.1.3 実装の詳細
- **ルーター**: Gorilla Mux Router（`app`）に直接登録
- **パス**: `/health`
- **メソッド**: `GET`のみ
- **認証**: 不要（認証ミドルウェアを通過しない）
- **レスポンス**: 
  - ステータスコード: `200 OK`
  - レスポンスボディ: `"OK"`（文字列）
  - Content-Type: `text/plain`

#### 3.1.4 実装位置の決定理由
- GoAdmin Engineの初期化後: GoAdmin Engineがルーターにマウントされる前に登録することで、GoAdmin Engineのミドルウェアチェーンを通過しない
- カスタムページの登録前: 既存のカスタムページ登録コードの前に配置することで、コードの可読性を維持

### 3.2 コード変更箇所

#### 3.2.1 変更前のコード構造

```go
// Gorilla Mux Router
app := mux.NewRouter()

// GoAdmin Engineの初期化
eng := engine.Default()
// ... GoAdmin Engineの設定 ...

// カスタムページの登録（Gorilla Mux用にContent関数を使用）
app.HandleFunc("/admin", ...).Methods("GET")
app.HandleFunc("/admin/", ...).Methods("GET")
// ... その他のカスタムページ ...
```

#### 3.2.2 変更後のコード構造

```go
// Gorilla Mux Router
app := mux.NewRouter()

// GoAdmin Engineの初期化
eng := engine.Default()
// ... GoAdmin Engineの設定 ...

// Health check endpoint (認証不要)
app.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "text/plain")
    w.Write([]byte("OK"))
}).Methods("GET")

// カスタムページの登録（Gorilla Mux用にContent関数を使用）
app.HandleFunc("/admin", ...).Methods("GET")
app.HandleFunc("/admin/", ...).Methods("GET")
// ... その他のカスタムページ ...
```

### 3.3 実装の注意事項

#### 3.3.1 認証ミドルウェアの回避
- `/health`エンドポイントは認証ミドルウェアを通過しない
- GoAdmin Engineのミドルウェアチェーンを通過しない
- アクセスログミドルウェアは通過する可能性がある（実装による）

#### 3.3.2 エラーハンドリング
- エラーハンドリングは不要（常に成功を返す）
- サーバーが起動している限り、常に`200 OK`を返す

#### 3.3.3 パフォーマンス
- シンプルな実装のため、レスポンス時間は1ms以下を想定
- 追加のリソース消費は不要

## 4. テスト設計

### 4.1 単体テスト

#### 4.1.1 テストファイル
- **ファイル**: `server/cmd/admin/main_test.go`（新規作成）
- **テスト関数**: `TestHealthEndpoint`

#### 4.1.2 テスト内容
- エンドポイントが`200 OK`を返すことを確認
- レスポンスボディが`"OK"`であることを確認
- Content-Typeが`text/plain`であることを確認
- 認証なしでアクセス可能であることを確認

#### 4.1.3 テストコード例

```go
package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/gorilla/mux"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestHealthEndpoint(t *testing.T) {
    // Gorilla Mux Routerを作成
    app := mux.NewRouter()
    
    // Health check endpointを登録
    app.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Header().Set("Content-Type", "text/plain")
        w.Write([]byte("OK"))
    }).Methods("GET")
    
    // テストリクエストを作成
    req := httptest.NewRequest(http.MethodGet, "/health", nil)
    rec := httptest.NewRecorder()
    
    // リクエストを処理
    app.ServeHTTP(rec, req)
    
    // アサーション
    assert.Equal(t, http.StatusOK, rec.Code)
    assert.Equal(t, "OK", rec.Body.String())
    assert.Equal(t, "text/plain", rec.Header().Get("Content-Type"))
}
```

### 4.2 統合テスト

#### 4.2.1 テストファイル
- **ファイル**: `server/test/integration/admin_health_test.go`（新規作成）
- **テスト関数**: `TestAdminHealth_HealthCheckNoAuth`

#### 4.2.2 テスト内容
- 実際のAdminサーバーを起動してテスト
- 認証なしで`/health`エンドポイントにアクセス可能であることを確認
- レスポンスが正しいことを確認

#### 4.2.3 テストコード例

```go
package integration

import (
    "net/http"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestAdminHealth_HealthCheckNoAuth(t *testing.T) {
    // Adminサーバーを起動（テスト用の設定で）
    // 実際の実装では、テスト用のサーバー起動関数を使用
    
    // Health check should not require auth
    resp, err := http.Get("http://localhost:8081/health")
    require.NoError(t, err)
    defer resp.Body.Close()
    
    assert.Equal(t, http.StatusOK, resp.StatusCode)
    
    // レスポンスボディを確認
    body := make([]byte, 1024)
    n, err := resp.Body.Read(body)
    require.NoError(t, err)
    assert.Equal(t, "OK", string(body[:n]))
    assert.Equal(t, "text/plain", resp.Header.Get("Content-Type"))
}
```

### 4.3 テスト実行方法

#### 4.3.1 単体テストの実行
```bash
cd server/cmd/admin
go test -v -run TestHealthEndpoint
```

#### 4.3.2 統合テストの実行
```bash
cd server/test/integration
go test -v -run TestAdminHealth_HealthCheckNoAuth
```

## 5. docker-compose.admin.ymlとの連携

### 5.1 ヘルスチェック設定

#### 5.1.1 既存の設定
```yaml
healthcheck:
  test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8081/health"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 40s
```

#### 5.1.2 動作確認方法
```bash
# docker-compose.admin.ymlを使用してコンテナを起動
docker-compose -f docker-compose.admin.yml up -d

# ヘルスチェックの状態を確認
docker ps

# コンテナのログを確認
docker-compose -f docker-compose.admin.yml logs admin

# 手動でヘルスチェックを実行
docker exec admin-develop wget --quiet --tries=1 --spider http://localhost:8081/health
```

### 5.2 動作確認の手順

#### 5.2.1 ローカル環境での確認
```bash
# Adminサーバーを起動
cd server/cmd/admin
go run main.go

# 別のターミナルでヘルスチェックを実行
curl http://localhost:8081/health
# 期待される出力: OK
```

#### 5.2.2 Docker環境での確認
```bash
# docker-compose.admin.ymlを使用してコンテナを起動
docker-compose -f docker-compose.admin.yml up -d

# コンテナのステータスを確認（healthyになることを確認）
docker ps | grep admin-develop

# ヘルスチェックを手動で実行
docker exec admin-develop wget --quiet --tries=1 --spider http://localhost:8081/health
# 期待される出力: （エラーなし）
```

## 6. 既存機能への影響

### 6.1 既存エンドポイントへの影響
- **影響なし**: 新規エンドポイントの追加のみ
- **既存のエンドポイント**: `/admin`、`/admin/dm-user/register`などは影響を受けない

### 6.2 GoAdmin Engineへの影響
- **影響なし**: GoAdmin Engineのミドルウェアチェーンを通過しない
- **既存の機能**: GoAdmin Engineの機能は影響を受けない

### 6.3 認証機能への影響
- **影響なし**: 認証ミドルウェアを通過しない
- **既存の認証機能**: 影響を受けない

### 6.4 アクセスログへの影響
- **影響あり（可能性）**: アクセスログミドルウェアを通過する可能性がある
- **対応**: アクセスログに記録されるが、問題なし（ヘルスチェックの記録は有用）

## 7. 実装上の注意事項

### 7.1 エンドポイント実装の注意事項
- **ルーターへの登録**: Gorilla Mux Routerに直接エンドポイントを登録する
- **認証ミドルウェア**: `/health`エンドポイントは認証ミドルウェアを通過しないようにする
- **GoAdmin Engine**: GoAdmin Engineのミドルウェアチェーンを通過しないようにする
- **実装の簡潔性**: APIサーバーと同様のシンプルな実装を維持する

### 7.2 テストの注意事項
- **単体テスト**: エンドポイントが正常に動作することを確認するテストを追加する
- **統合テスト**: 実際のサーバーを起動してテストする（可能な場合）
- **既存テスト**: 既存のテストが全て失敗しないことを確認する

### 7.3 動作確認の注意事項
- **ローカル環境**: ローカル環境で`curl http://localhost:8081/health`が正常に動作することを確認
- **Docker環境**: docker-compose.admin.ymlを使用してコンテナを起動し、ヘルスチェックが正常に動作することを確認
- **既存機能**: 既存のエンドポイント（`/admin`など）が正常に動作することを確認

## 8. 参考実装

### 8.1 APIサーバーの実装

#### 8.1.1 エンドポイント実装
```go
// Health check
e.GET("/health", func(c echo.Context) error {
    return c.String(http.StatusOK, "OK")
})
```

#### 8.1.2 テスト実装
```go
func TestHealthEndpoint(t *testing.T) {
    cfg := testutil.GetTestConfig()
    router := NewRouter(nil, nil, nil, nil, nil, cfg)
    
    req := httptest.NewRequest(http.MethodGet, "/health", nil)
    rec := httptest.NewRecorder()
    router.ServeHTTP(rec, req)
    
    assert.Equal(t, http.StatusOK, rec.Code)
    assert.Equal(t, "OK", rec.Body.String())
}
```

### 8.2 Adminサーバーでの実装方針
- APIサーバーと同様のシンプルな実装を維持
- Gorilla Mux Routerを使用するため、実装方法は異なるが、動作は同じ
- テストも同様のパターンで実装

## 9. 実装チェックリスト

### 9.1 実装項目
- [ ] `server/cmd/admin/main.go`に`/health`エンドポイントを追加
- [ ] エンドポイントが認証なしでアクセス可能であることを確認
- [ ] エンドポイントが`200 OK`と`"OK"`を返すことを確認
- [ ] Content-Typeが`text/plain`であることを確認

### 9.2 テスト項目
- [ ] 単体テストを実装（`server/cmd/admin/main_test.go`）
- [ ] 統合テストを実装（`server/test/integration/admin_health_test.go`）
- [ ] 既存のテストが全て失敗しないことを確認

### 9.3 動作確認項目
- [ ] ローカル環境で`curl http://localhost:8081/health`が正常に動作する
- [ ] Docker環境で`wget --quiet --tries=1 --spider http://localhost:8081/health`が正常に動作する
- [ ] docker-compose.admin.ymlのヘルスチェックが正常に動作する
- [ ] コンテナのステータスが`healthy`になる
- [ ] 既存のエンドポイント（`/admin`など）が正常に動作する

## 10. 参考情報

### 10.1 関連Issue
- GitHub Issue #103: GoAdmin死活監視エンドポイントの作成

### 10.2 既存実装の参考
- **APIサーバー**: `server/internal/api/router/router.go`の`/health`エンドポイント実装
- **APIサーバーのテスト**: `server/internal/api/router/router_test.go`の`TestHealthEndpoint`
- **統合テスト**: `server/test/integration/api_auth_test.go`の`TestAPIAuth_HealthCheckNoAuth`

### 10.3 技術スタック
- **言語**: Go
- **ルーター**: Gorilla Mux Router
- **フレームワーク**: GoAdmin Engine
- **コンテナ管理**: Docker Compose
- **ヘルスチェックツール**: wget

### 10.4 関連ドキュメント
- `docker-compose.admin.yml`: AdminサーバーのDocker Compose設定
- `server/cmd/admin/main.go`: Adminサーバーのメインエントリーポイント
- `server/internal/api/router/router.go`: APIサーバーのルーター実装（参考）
