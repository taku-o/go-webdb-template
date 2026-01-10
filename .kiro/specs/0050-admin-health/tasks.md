# GoAdmin死活監視エンドポイントの実装タスク一覧

## 概要
Adminサーバーに`/health`エンドポイントを実装するためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: `/health`エンドポイントの実装

#### タスク 1.1: `/health`エンドポイントの実装
**目的**: Adminサーバーに`/health`エンドポイントを実装する。

**作業内容**:
- `server/cmd/admin/main.go`を開く
- GoAdmin Engineの初期化後、カスタムページの登録前に`/health`エンドポイントを追加
- エンドポイントは認証不要で、`200 OK`と`"OK"`を返す

**実装コード**:
```go
// Health check endpoint (認証不要)
app.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "text/plain")
    w.Write([]byte("OK"))
}).Methods("GET")
```

**実装位置**:
- GoAdmin Engineの初期化後（`eng.Use(app)`の後）
- カスタムページの登録前（`app.HandleFunc("/admin", ...)`の前）

**受け入れ基準**:
- `/health`エンドポイントが実装されている
- エンドポイントが認証なしでアクセス可能である
- エンドポイントが`200 OK`と`"OK"`を返す
- Content-Typeが`text/plain`である

- _Requirements: 3.1.1, 3.1.2, 3.1.3, 6.1_
- _Design: 3.1, 3.2, 3.3_

---

### Phase 2: 単体テストの実装

#### タスク 2.1: 単体テストファイルの作成
**目的**: `/health`エンドポイントの単体テストを実装する。

**作業内容**:
- `server/cmd/admin/main_test.go`を新規作成
- `TestHealthEndpoint`関数を実装
- エンドポイントが正常に動作することを確認するテストを追加

**テスト内容**:
- エンドポイントが`200 OK`を返すことを確認
- レスポンスボディが`"OK"`であることを確認
- Content-Typeが`text/plain`であることを確認
- 認証なしでアクセス可能であることを確認

**テストコード**:
```go
package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/gorilla/mux"
    "github.com/stretchr/testify/assert"
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

**受け入れ基準**:
- 単体テストが実装されている
- テストが正常に実行できる
- すべてのアサーションが成功する

- _Requirements: 6.4_
- _Design: 4.1_

---

#### タスク 2.2: 単体テストの実行と確認
**目的**: 実装した単体テストが正常に動作することを確認する。

**作業内容**:
- `server/cmd/admin`ディレクトリでテストを実行
- テストが正常に完了することを確認
- テスト結果を確認

**実行コマンド**:
```bash
cd server/cmd/admin
go test -v -run TestHealthEndpoint
```

**受け入れ基準**:
- テストが正常に実行できる
- すべてのテストが成功する
- エラーが発生しない

- _Requirements: 6.4_
- _Design: 4.3.1_

---

### Phase 3: 統合テストの実装（オプション）

#### タスク 3.1: 統合テストファイルの作成
**目的**: `/health`エンドポイントの統合テストを実装する（可能な場合）。

**作業内容**:
- `server/test/integration/admin_health_test.go`を新規作成
- `TestAdminHealth_HealthCheckNoAuth`関数を実装
- 実際のAdminサーバーを起動してテスト（可能な場合）

**注意事項**:
- 統合テストは、実際のサーバー起動が必要なため、実装が困難な場合はスキップ可能
- ローカル環境での動作確認で代替可能

**受け入れ基準**:
- 統合テストが実装されている（可能な場合）
- テストが正常に実行できる（可能な場合）

- _Requirements: 6.4_
- _Design: 4.2_

---

### Phase 4: 動作確認

#### タスク 4.1: ローカル環境での動作確認
**目的**: ローカル環境で`/health`エンドポイントが正常に動作することを確認する。

**作業内容**:
- Adminサーバーをローカル環境で起動
- `curl`コマンドで`/health`エンドポイントにアクセス
- レスポンスを確認

**実行コマンド**:
```bash
# Adminサーバーを起動（別ターミナル）
cd server/cmd/admin
go run main.go

# 別のターミナルでヘルスチェックを実行
curl http://localhost:8081/health
# 期待される出力: OK
```

**確認項目**:
- ステータスコードが`200 OK`である
- レスポンスボディが`"OK"`である
- Content-Typeが`text/plain`である
- 認証なしでアクセス可能である

**受け入れ基準**:
- ローカル環境で`curl http://localhost:8081/health`が正常に動作する
- レスポンスが正しい
- エラーが発生しない

- _Requirements: 6.3_
- _Design: 5.2.1_

---

#### タスク 4.2: Docker環境での動作確認
**目的**: Docker環境で`/health`エンドポイントが正常に動作することを確認する。

**作業内容**:
- docker-compose.admin.ymlを使用してコンテナを起動
- ヘルスチェックが正常に動作することを確認
- コンテナのステータスが`healthy`になることを確認

**実行コマンド**:
```bash
# docker-compose.admin.ymlを使用してコンテナを起動
docker-compose -f docker-compose.admin.yml up -d

# コンテナのステータスを確認（healthyになることを確認）
docker ps | grep admin-develop

# ヘルスチェックを手動で実行
docker exec admin-develop wget --quiet --tries=1 --spider http://localhost:8081/health
# 期待される出力: （エラーなし）
```

**確認項目**:
- コンテナが正常に起動する
- ヘルスチェックが正常に動作する
- コンテナのステータスが`healthy`になる
- ヘルスチェックのログにエラーが表示されない

**受け入れ基準**:
- docker-compose.admin.ymlを使用してコンテナを起動できる
- ヘルスチェックが正常に動作する
- コンテナのステータスが`healthy`になる
- ヘルスチェックのログにエラーが表示されない

- _Requirements: 3.2, 6.2, 6.3_
- _Design: 5.1, 5.2.2_

---

#### タスク 4.3: 既存エンドポイントの動作確認
**目的**: 既存のエンドポイントが正常に動作することを確認する。

**作業内容**:
- 既存のエンドポイント（`/admin`など）にアクセス
- 正常に動作することを確認

**確認項目**:
- `/admin`エンドポイントが正常に動作する
- その他の既存エンドポイントが正常に動作する
- エラーが発生しない

**受け入れ基準**:
- 既存のエンドポイント（`/admin`など）が正常に動作する
- エラーが発生しない

- _Requirements: 6.3_
- _Design: 6.1_

---

### Phase 5: 既存テストの確認

#### タスク 5.1: 既存テストの実行
**目的**: 既存のテストが全て失敗しないことを確認する。

**作業内容**:
- 既存のテストを実行
- すべてのテストが成功することを確認
- エラーが発生しないことを確認

**実行コマンド**:
```bash
# 既存のテストを実行
cd server
go test ./...

# または、特定のパッケージのテストを実行
go test ./cmd/admin/...
go test ./test/integration/...
```

**受け入れ基準**:
- 既存のテストが全て失敗しない
- すべてのテストが成功する
- エラーが発生しない

- _Requirements: 6.4_
- _Design: 7.2_

---

## 実装チェックリスト

### 実装項目
- [ ] `server/cmd/admin/main.go`に`/health`エンドポイントを追加
- [ ] エンドポイントが認証なしでアクセス可能であることを確認
- [ ] エンドポイントが`200 OK`と`"OK"`を返すことを確認
- [ ] Content-Typeが`text/plain`であることを確認

### テスト項目
- [ ] 単体テストを実装（`server/cmd/admin/main_test.go`）
- [ ] 単体テストが正常に実行できる
- [ ] 統合テストを実装（可能な場合）（`server/test/integration/admin_health_test.go`）
- [ ] 既存のテストが全て失敗しないことを確認

### 動作確認項目
- [ ] ローカル環境で`curl http://localhost:8081/health`が正常に動作する
- [ ] Docker環境で`wget --quiet --tries=1 --spider http://localhost:8081/health`が正常に動作する
- [ ] docker-compose.admin.ymlのヘルスチェックが正常に動作する
- [ ] コンテナのステータスが`healthy`になる
- [ ] 既存のエンドポイント（`/admin`など）が正常に動作する

## 参考情報

### 関連ドキュメント
- 要件定義書: `requirements.md`
- 設計書: `design.md`

### 既存実装の参考
- **APIサーバー**: `server/internal/api/router/router.go`の`/health`エンドポイント実装
- **APIサーバーのテスト**: `server/internal/api/router/router_test.go`の`TestHealthEndpoint`
- **統合テスト**: `server/test/integration/api_auth_test.go`の`TestAPIAuth_HealthCheckNoAuth`

### 技術スタック
- **言語**: Go
- **ルーター**: Gorilla Mux Router
- **フレームワーク**: GoAdmin Engine
- **コンテナ管理**: Docker Compose
- **ヘルスチェックツール**: wget
