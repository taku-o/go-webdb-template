# JobQueueサーバー死活監視エンドポイントの実装タスク一覧

## 概要
JobQueueサーバーにHTTPサーバーを追加し、`/health`エンドポイントを実装するためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: 設定ファイルの拡張

#### - [x] タスク 1.1: `JobQueueConfig`構造体の追加
**目的**: 設定構造体に`JobQueueConfig`を追加し、設定ファイルから読み込めるようにする。

**作業内容**:
- `server/internal/config/config.go`を開く
- `JobQueueConfig`構造体を定義
- `Config`構造体に`JobQueue`フィールドを追加

**実装コード**:
```go
// JobQueueConfig はJobQueueサーバー設定
type JobQueueConfig struct {
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// Config構造体に追加
type Config struct {
	// ... 既存のフィールド ...
	JobQueue JobQueueConfig `mapstructure:"jobqueue"`
}
```

**受け入れ基準**:
- `JobQueueConfig`構造体が定義されている
- `Config`構造体に`JobQueue`フィールドが追加されている
- 設定ファイルから読み込める構造になっている

- _Requirements: 3.3.1, 6.3_
- _Design: 3.1.1_

---

#### - [x] タスク 1.2: `config/develop/config.yaml`に`jobqueue`セクションを追加
**目的**: 開発環境の設定ファイルに`jobqueue`セクションを追加する。

**作業内容**:
- `config/develop/config.yaml`を開く
- `jobqueue`セクションを追加
- ポート番号、タイムアウト設定を追加

**実装内容**:
```yaml
jobqueue:
  port: 8082
  read_timeout: 30s
  write_timeout: 30s
```

**受け入れ基準**:
- `config/develop/config.yaml`に`jobqueue`セクションが追加されている
- ポート番号が8082に設定されている
- タイムアウト設定が適切に設定されている

- _Requirements: 3.3.2, 6.3_
- _Design: 3.1.2_

---

#### - [x] タスク 1.3: `config/staging/config.yaml`に`jobqueue`セクションを追加
**目的**: ステージング環境の設定ファイルに`jobqueue`セクションを追加する。

**作業内容**:
- `config/staging/config.yaml`を開く
- `jobqueue`セクションを追加
- ポート番号、タイムアウト設定を追加

**実装内容**:
```yaml
jobqueue:
  port: 8082
  read_timeout: 30s
  write_timeout: 30s
```

**受け入れ基準**:
- `config/staging/config.yaml`に`jobqueue`セクションが追加されている
- ポート番号が8082に設定されている
- タイムアウト設定が適切に設定されている

- _Requirements: 3.3.2, 6.3_
- _Design: 3.1.2_

---

#### - [x] タスク 1.4: `config/production/config.yaml.example`に`jobqueue`セクションを追加
**目的**: 本番環境の設定ファイルテンプレートに`jobqueue`セクションを追加する。

**作業内容**:
- `config/production/config.yaml.example`を開く
- `jobqueue`セクションを追加
- ポート番号、タイムアウト設定を追加

**実装内容**:
```yaml
jobqueue:
  port: 8082
  read_timeout: 30s
  write_timeout: 30s
```

**受け入れ基準**:
- `config/production/config.yaml.example`に`jobqueue`セクションが追加されている
- ポート番号が8082に設定されている
- タイムアウト設定が適切に設定されている

- _Requirements: 3.3.2, 6.3_
- _Design: 3.1.2_

---

### Phase 2: HTTPサーバーの実装

#### - [x] タスク 2.1: `server/cmd/jobqueue/main.go`にHTTPサーバーの実装を追加
**目的**: JobQueueサーバーにHTTPサーバーを追加し、`/health`エンドポイントを実装する。

**作業内容**:
- `server/cmd/jobqueue/main.go`を開く
- 必要なパッケージをインポート（`context`、`fmt`、`net/http`、`time`）
- Asynqサーバーの初期化後、HTTPサーバーの初期化処理を追加
- `/health`エンドポイントを登録
- HTTPサーバーとAsynqサーバーを並行して起動（goroutineで起動）
- Graceful shutdownの実装を追加（HTTPサーバーとAsynqサーバーの両方）

**実装コード**:
```go
// 3. HTTPサーバーの初期化
mux := http.NewServeMux()

// Health check endpoint (認証不要)
mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
})

httpServer := &http.Server{
	Addr:         fmt.Sprintf(":%d", cfg.JobQueue.Port),
	Handler:      mux,
	ReadTimeout:  cfg.JobQueue.ReadTimeout,
	WriteTimeout: cfg.JobQueue.WriteTimeout,
}

// 4. Asynqサーバーの起動（バックグラウンド）
go func() {
	log.Println("Starting job queue processing...")
	if err := jobQueueServer.Start(); err != nil {
		log.Printf("ERROR: Failed to start job queue server: %v", err)
	}
}()

// 5. HTTPサーバーの起動（バックグラウンド）
go func() {
	log.Printf("Starting HTTP server on port %d", cfg.JobQueue.Port)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server failed: %v", err)
	}
}()

// ... 既存のGraceful shutdown処理 ...

// 7. Asynqサーバーの停止
if err := jobQueueServer.Shutdown(); err != nil {
	log.Printf("JobQueue server shutdown error: %v", err)
}

// 8. HTTPサーバーの停止（30秒のタイムアウト）
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

if err := httpServer.Shutdown(ctx); err != nil {
	log.Printf("HTTP server forced to shutdown: %v", err)
}
```

**受け入れ基準**:
- HTTPサーバーの初期化処理が追加されている
- `/health`エンドポイントが実装されている
- HTTPサーバーとAsynqサーバーが並行して起動する
- Graceful shutdownが実装されている（両方のサーバーを適切に停止）

- _Requirements: 3.1.1, 3.2, 4.5, 6.1, 6.2_
- _Design: 3.2, 3.4_

---

### Phase 3: 単体テストの実装

#### - [x] タスク 3.1: `server/cmd/jobqueue/main_test.go`の作成
**目的**: `/health`エンドポイントの単体テストを実装する。

**作業内容**:
- `server/cmd/jobqueue/main_test.go`を新規作成
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

	"github.com/stretchr/testify/assert"
)

func TestHealthEndpoint(t *testing.T) {
	// HTTPサーバーを作成
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// テストリクエストを作成
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	// リクエストを処理
	mux.ServeHTTP(rec, req)

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

- _Requirements: 6.5_
- _Design: 4.1_

---

#### - [x] タスク 3.2: 単体テストの実行と確認
**目的**: 実装した単体テストが正常に動作することを確認する。

**作業内容**:
- `server/cmd/jobqueue`ディレクトリでテストを実行
- テストが正常に完了することを確認
- テスト結果を確認

**実行コマンド**:
```bash
cd server/cmd/jobqueue
go test -v -run TestHealthEndpoint
```

**受け入れ基準**:
- テストが正常に実行できる
- すべてのテストが成功する
- エラーが発生しない

- _Requirements: 6.5_
- _Design: 4.3.1_

---

### Phase 4: 統合テストの実装（オプション）

#### - [x] タスク 4.1: 統合テストファイルの作成
**目的**: `/health`エンドポイントの統合テストを実装する（可能な場合）。

**作業内容**:
- `server/test/integration/jobqueue_health_test.go`を新規作成
- `TestJobQueueHealth_HealthCheckNoAuth`関数を実装
- 実際のJobQueueサーバーを起動してテスト（可能な場合）

**注意事項**:
- 統合テストは、実際のサーバー起動が必要なため、実装が困難な場合はスキップ可能
- ローカル環境での動作確認で代替可能

**受け入れ基準**:
- 統合テストが実装されている（可能な場合）
- テストが正常に実行できる（可能な場合）

- _Requirements: 6.5_
- _Design: 4.2_

---

### Phase 5: 動作確認

#### - [ ] タスク 5.1: ローカル環境での動作確認
**目的**: ローカル環境で`/health`エンドポイントが正常に動作することを確認する。

**作業内容**:
- JobQueueサーバーをローカル環境で起動
- `curl`コマンドで`/health`エンドポイントにアクセス
- レスポンスを確認
- HTTPサーバーとAsynqサーバーの両方が動作することを確認

**実行コマンド**:
```bash
# JobQueueサーバーを起動（別ターミナル）
cd server
APP_ENV=develop go run ./cmd/jobqueue/main.go

# 別のターミナルでヘルスチェックを実行
curl http://localhost:8082/health
# 期待される出力: OK
```

**確認項目**:
- ステータスコードが`200 OK`である
- レスポンスボディが`"OK"`である
- Content-Typeが`text/plain`である
- 認証なしでアクセス可能である
- HTTPサーバーとAsynqサーバーの両方が動作している

**受け入れ基準**:
- ローカル環境で`curl http://localhost:8082/health`が正常に動作する
- レスポンスが正しい
- エラーが発生しない
- HTTPサーバーとAsynqサーバーの両方が動作している

- _Requirements: 6.4_
- _Design: 6.4_

---

#### - [ ] タスク 5.2: Graceful shutdownの動作確認
**目的**: Graceful shutdownが正常に動作することを確認する。

**作業内容**:
- JobQueueサーバーを起動
- SIGINTまたはSIGTERMシグナルを送信
- 両方のサーバー（HTTPサーバーとAsynqサーバー）が適切に停止することを確認

**実行コマンド**:
```bash
# JobQueueサーバーを起動
cd server
APP_ENV=develop go run ./cmd/jobqueue/main.go

# 別のターミナルでシグナルを送信
# Ctrl+C または kill -SIGTERM <PID>
```

**確認項目**:
- HTTPサーバーが適切に停止する
- Asynqサーバーが適切に停止する
- エラーログが出力されない
- 正常に終了する

**受け入れ基準**:
- Graceful shutdownが正常に動作する
- 両方のサーバーが適切に停止する
- エラーが発生しない

- _Requirements: 6.1, 6.4_
- _Design: 3.4.2, 6.4_

---

#### - [ ] タスク 5.3: 既存のAsynqサーバー機能の動作確認
**目的**: 既存のAsynqサーバーの機能が正常に動作することを確認する。

**作業内容**:
- JobQueueサーバーを起動
- 既存のジョブ登録APIからジョブを登録
- ジョブが正常に処理されることを確認

**確認項目**:
- ジョブが正常に登録される
- ジョブが正常に処理される
- 既存の機能に影響がない

**受け入れ基準**:
- 既存のAsynqサーバーの機能が正常に動作する
- ジョブが正常に処理される
- エラーが発生しない

- _Requirements: 6.4_
- _Design: 5.1_

---

### Phase 6: 既存テストの確認

#### - [x] タスク 6.1: 既存テストの実行
**目的**: 既存のテストが全て失敗しないことを確認する。

**作業内容**:
- 既存のテストを実行
- すべてのテストが成功することを確認
- エラーが発生しないことを確認

**実行コマンド**:
```bash
# 既存のテストを実行
cd server
APP_ENV=test go test ./...

# または、特定のパッケージのテストを実行
APP_ENV=test go test ./cmd/jobqueue/...
APP_ENV=test go test ./test/integration/...
```

**受け入れ基準**:
- 既存のテストが全て失敗しない
- すべてのテストが成功する
- エラーが発生しない

- _Requirements: 6.5_
- _Design: 6.3_

---

### Phase 7: ドキュメントの更新

#### - [x] タスク 7.1: `docs/ja/Queue-Job.md`の更新
**目的**: JobQueueサーバーの`/health`エンドポイントの情報を追加する。

**作業内容**:
- 「4. JobQueueサーバーの起動」セクションに、HTTPサーバーの起動について追記:
  - JobQueueサーバーはHTTPサーバー（ポート8082）とAsynqサーバーを並行して起動する
- 新しいセクション「ヘルスチェック」を追加:
  - `/health`エンドポイントの説明
  - ポート8082でHTTPサーバーが起動すること
  - 認証不要でアクセス可能であること
  - レスポンスが`200 OK`と`"OK"`であること

**追加内容例**:
```markdown
### ヘルスチェック

JobQueueサーバーには`/health`エンドポイントが用意されています。

**エンドポイント**: `GET http://localhost:8082/health`

**認証**: 不要

**レスポンス**:
- ステータスコード: `200 OK`
- レスポンスボディ: `"OK"`（文字列）
- Content-Type: `text/plain`

**使用例**:
```bash
curl http://localhost:8082/health
# 期待される出力: OK
```

**注意**: JobQueueサーバーはHTTPサーバー（ポート8082）とAsynqサーバーを並行して起動します。
```

**受け入れ基準**:
- `/health`エンドポイントの説明が追加されている
- ポート8082の情報が記載されている
- 使用例が記載されている

- _Requirements: 1.2, 3.2_
- _Design: 3.2_

---

#### - [x] タスク 7.2: `docs/en/Queue-Job.md`の更新
**目的**: 英語版のジョブキュー機能の利用手順に`/health`エンドポイントの情報を追加する。

**作業内容**:
- `docs/ja/Queue-Job.md`と同様の更新を英語版に適用
- 「4. Start JobQueue Server」セクションにHTTPサーバーの起動について追記
- 新しいセクション「Health Check」を追加

**受け入れ基準**:
- 英語版のドキュメントが日本語版と同様に更新されている
- `/health`エンドポイントの情報が正しく記載されている

- _Requirements: 1.2, 3.2_
- _Design: 3.2_

---

#### - [x] タスク 7.3: `README.ja.md`の更新
**目的**: プロジェクト概要にJobQueueサーバーの`/health`エンドポイントとポート情報を追加する。

**作業内容**:
- 「開発環境サーバー構成」セクション（または該当箇所）にJobQueueサーバーのポート情報を追加:
  - JobQueueサーバー: ポート8082
- 「APIエンドポイント一覧」セクションにJobQueueサーバーの`/health`エンドポイントを追加:
  - `GET http://localhost:8082/health` - ヘルスチェック（認証不要）
- 必要に応じて、他のサーバー（API、Admin）と同様に記載

**追加内容例**:
```markdown
### 開発環境サーバー構成

| サーバー | ポート | ディレクトリ | 起動コマンド |
|---------|-------|-------------|-------------|
| API サーバー | 8080 | `server/cmd/server` | `APP_ENV=develop go run ./cmd/server/main.go` |
| クライアント | 3000 | `client/` | `npm run dev` |
| Admin | 8081 | `server/cmd/admin` | `APP_ENV=develop go run ./cmd/admin/main.go` |
| JobQueue | 8082 | `server/cmd/jobqueue` | `APP_ENV=develop go run ./cmd/jobqueue/main.go` |
```

**受け入れ基準**:
- JobQueueサーバーのポート情報（8082）が記載されている
- `/health`エンドポイントがAPIエンドポイント一覧に追加されている
- 他のサーバーと一貫性のある記載になっている

- _Requirements: 1.2, 3.2_
- _Design: 3.2_

---

#### - [x] タスク 7.4: `README.md`の更新（該当する場合）
**目的**: 英語版のプロジェクト概要にJobQueueサーバーの`/health`エンドポイントとポート情報を追加する。

**作業内容**:
- `README.ja.md`と同様の更新を英語版に適用
- 開発環境サーバー構成にJobQueueサーバーを追加
- APIエンドポイント一覧にJobQueueサーバーの`/health`エンドポイントを追加

**受け入れ基準**:
- 英語版のドキュメントが日本語版と同様に更新されている
- JobQueueサーバーに関する情報が正しく記載されている

- _Requirements: 1.2, 3.2_
- _Design: 3.2_

---

## 実装チェックリスト

### 実装項目
- [x] `server/internal/config/config.go`に`JobQueueConfig`構造体を追加
- [x] `config/develop/config.yaml`に`jobqueue`セクションを追加
- [x] `config/staging/config.yaml`に`jobqueue`セクションを追加
- [x] `config/production/config.yaml.example`に`jobqueue`セクションを追加
- [x] `server/cmd/jobqueue/main.go`にHTTPサーバーの起動処理を追加
- [x] `server/cmd/jobqueue/main.go`に`/health`エンドポイントを追加
- [x] HTTPサーバーとAsynqサーバーを並行して起動する実装を追加
- [x] Graceful shutdownの実装を追加（HTTPサーバーとAsynqサーバーの両方）

### テスト項目
- [x] 単体テストを実装（`server/cmd/jobqueue/main_test.go`）
- [x] 単体テストが正常に実行できる
- [x] 統合テストを実装（可能な場合）（`server/test/integration/jobqueue_health_test.go`）
- [x] 既存のテストが全て失敗しないことを確認

### 動作確認項目
- [ ] ローカル環境で`curl http://localhost:8082/health`が正常に動作する
- [ ] JobQueueサーバーが起動した時、HTTPサーバーとAsynqサーバーの両方が動作する
- [ ] 既存のAsynqサーバーの機能が正常に動作することを確認
- [ ] Graceful shutdownが正常に動作する（両方のサーバーが適切に停止する）

### ドキュメント更新項目
- [x] `docs/ja/Queue-Job.md`に`/health`エンドポイントの情報を追加
- [x] `docs/en/Queue-Job.md`に`/health`エンドポイントの情報を追加
- [x] `README.ja.md`にJobQueueサーバーのポート情報と`/health`エンドポイントを追加
- [x] `README.md`にJobQueueサーバーのポート情報と`/health`エンドポイントを追加（該当する場合）

## 参考情報

### 関連ドキュメント
- 要件定義書: `requirements.md`
- 設計書: `design.md`

### 既存実装の参考
- **Adminサーバー**: `server/cmd/admin/main.go`の`/health`エンドポイント実装
- **Adminサーバーのテスト**: `server/cmd/admin/main_test.go`の`TestHealthEndpoint`
- **APIサーバー**: `server/internal/api/router/router.go`の`/health`エンドポイント実装

### 技術スタック
- **言語**: Go
- **HTTPサーバー**: 標準ライブラリの`net/http`
- **設定管理**: `github.com/spf13/viper`（既存システムと同様）
- **ジョブキュー**: `github.com/hibiken/asynq`（既存システムと同様）
