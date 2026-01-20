# サーバー状態確認機能のリファクタリング設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、`server/cmd/server-status/main.go`を3層構造（CLI → Usecase → Service）にリファクタリングするための詳細設計を定義する。プロジェクトの標準的なアーキテクチャパターンに従い、既存の機能を維持しながら、コードの保守性とテスタビリティを向上させる。

### 1.2 設計の範囲
- CLI層（`server/cmd/server-status/main.go`）のリファクタリング
- Usecase層（`server/internal/usecase/cli/server_status_usecase.go`）の新規作成
- Service層（`server/internal/service/server_status_service.go`）の新規作成
- 既存のサーバー状態確認機能の維持
- 各層の単体テストの実装

### 1.3 設計方針
- **3層構造の遵守**: CLI層 → Usecase層 → Service層の順で呼び出し
- **既存機能の維持**: 既存のサーバー状態確認機能と表示形式を完全に維持
- **インターフェースの使用**: Service層はインターフェースで定義し、テスト容易性を確保
- **既存パターンの遵守**: 既存のCLI実装（`list_dm_users_usecase.go`など）と同様の構造
- **標準ライブラリの使用**: 外部依存を追加せず、標準ライブラリのみを使用

## 2. アーキテクチャ設計

### 2.1 レイヤー構成

```
┌─────────────────────────────────────────────────────────────┐
│              CLI層 (server/cmd/server-status/main.go)        │
│              - 入出力制御                                      │
│              - usecaseを呼び出す                               │
│              - 結果を表形式で表示                              │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│    Usecase層 (server/internal/usecase/cli/                  │
│                server_status_usecase.go)                      │
│    - サーバー情報の定義                                        │
│    - serviceに渡すパラメータを作る                            │
│    - serviceから受け取ったリストをmain.goに渡す                │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│    Service層 (server/internal/service/                        │
│                server_status_service.go)                        │
│    - サーバー状態確認のビジネスロジック                        │
│    - 並列実行による状態確認                                    │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 データフロー

```
┌─────────────────────────────────────────────────────────────┐
│              サーバー状態確認のデータフロー                      │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  1. CLI層: main()                │
        │     - usecaseのインスタンス作成   │
        │     - usecase.ListServerStatus() │
        │       を呼び出し                  │
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  2. Usecase層: ListServerStatus()│
        │     - サーバー情報リストを取得    │
        │     - service.ListServerStatus() │
        │       (servers)を呼び出し         │
        │     - 結果をそのまま返す          │
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  3. Service層: ListServerStatus()│
        │     - サーバー一覧をパラメータとして│
        │       受け取る                    │
        │     - 並列実行で状態確認          │
        │     - []ServerStatusを返す       │
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  4. Usecase層: 結果を返す        │
        │     - []ServerStatusを返す       │
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  5. CLI層: 結果を表示            │
        │     - printResults()で表形式表示 │
        └─────────────────────────────────┘
```

### 2.3 ファイル構造

```
server/
├── cmd/
│   └── server-status/
│       └── main.go  (変更)
├── internal/
│   ├── usecase/
│   │   └── cli/
│   │       ├── server_status_usecase.go  (新規作成)
│   │       └── server_status_usecase_test.go  (新規作成)
│   └── service/
│       ├── server_status_service.go  (新規作成)
│       └── server_status_service_test.go  (新規作成)
```

## 3. 実装設計

### 3.1 データ構造

#### 3.1.1 ServerInfo
サーバー情報を表す構造体（既存の定義を維持）

**配置場所**: `server/internal/service/server_status_service.go`

```go
// ServerInfo はサーバー情報を表す
type ServerInfo struct {
	Name    string // サーバー名
	Port    int    // ポート番号
	Address string // 接続先アドレス（通常は"localhost"）
}
```

#### 3.1.2 ServerStatus
サーバーの状態を表す構造体（既存の定義を維持）

**配置場所**: `server/internal/service/server_status_service.go`

```go
// ServerStatus はサーバーの状態を表す
type ServerStatus struct {
	Server ServerInfo
	Status string // "起動中" または "停止中"
	Error  error  // エラー情報（デバッグ用、表示には使用しない）
}
```

### 3.2 Service層の実装

#### 3.2.1 ファイル構成

**ファイル**: `server/internal/service/server_status_service.go`

#### 3.2.2 インターフェース定義

```go
// ServerStatusServiceInterface はServerStatusServiceのインターフェース
type ServerStatusServiceInterface interface {
	ListServerStatus(ctx context.Context, servers []ServerInfo) ([]ServerStatus, error)
}
```

**注意**: インターフェースはusecase層で定義する（既存パターンに従う）

#### 3.2.3 Service実装

```go
package service

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	// connectionTimeout はTCP接続のタイムアウト時間
	connectionTimeout = 1 * time.Second
)

// ServerInfo はサーバー情報を表す
type ServerInfo struct {
	Name    string // サーバー名
	Port    int    // ポート番号
	Address string // 接続先アドレス（通常は"localhost"）
}

// ServerStatus はサーバーの状態を表す
type ServerStatus struct {
	Server ServerInfo
	Status string // "起動中" または "停止中"
	Error  error  // エラー情報（デバッグ用、表示には使用しない）
}

// ServerStatusService はサーバー状態確認のビジネスロジックを担当
type ServerStatusService struct {
	// 必要に応じて依存関係を追加（現時点では不要）
}

// NewServerStatusService は新しいServerStatusServiceを作成
func NewServerStatusService() *ServerStatusService {
	return &ServerStatusService{}
}

// checkServerStatus は指定されたサーバーの状態を確認する
func (s *ServerStatusService) checkServerStatus(server ServerInfo, timeout time.Duration) ServerStatus {
	address := fmt.Sprintf("%s:%d", server.Address, server.Port)

	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return ServerStatus{
			Server: server,
			Status: "停止中",
			Error:  err,
		}
	}

	// 接続が成功した場合は即座に閉じる
	conn.Close()

	return ServerStatus{
		Server: server,
		Status: "起動中",
		Error:  nil,
	}
}

// checkAllServers は全サーバーの状態を並列に確認する
// 注意: 各goroutineはresults[index]に直接書き込むため、元のserversの順序が保たれる
func (s *ServerStatusService) checkAllServers(servers []ServerInfo, timeout time.Duration) []ServerStatus {
	var wg sync.WaitGroup
	results := make([]ServerStatus, len(servers))

	for i, server := range servers {
		wg.Add(1)
		go func(index int, srv ServerInfo) {
			defer wg.Done()
			// indexは元のserversの順序に対応しているため、順序が保たれる
			results[index] = s.checkServerStatus(srv, timeout)
		}(i, server)
	}

	wg.Wait()
	return results
}

// ListServerStatus は指定されたサーバーリストの状態を確認して返す
func (s *ServerStatusService) ListServerStatus(ctx context.Context, servers []ServerInfo) ([]ServerStatus, error) {
	results := s.checkAllServers(servers, connectionTimeout)
	return results, nil
}
```

#### 3.2.4 実装の詳細
- **状態確認ロジック**: `checkServerStatus()`メソッドで実装（既存の`checkServerStatus`関数の内容を移行）
- **並列実行ロジック**: `checkAllServers()`メソッドで実装（既存の`checkAllServers`関数の内容を移行）
  - **順序の保持**: 各goroutineは`results[index]`に直接書き込むため、元の`servers`の順序が保たれる（`index`はループの`i`の値で、`servers`の順序に対応）
- **公開メソッド**: `ListServerStatus(ctx, servers)`メソッドでusecase層から呼び出し可能（サーバー一覧をパラメータとして受け取る）

### 3.3 Usecase層の実装

#### 3.3.1 ファイル構成

**ファイル**: `server/internal/usecase/cli/server_status_usecase.go`

#### 3.3.2 インターフェース定義

```go
// ServerStatusServiceInterface はServerStatusServiceのインターフェース
type ServerStatusServiceInterface interface {
	ListServerStatus(ctx context.Context, servers []service.ServerInfo) ([]service.ServerStatus, error)
}
```

**注意**: インターフェースはusecase層で定義する（既存パターンに従う）

#### 3.3.3 Usecase実装

```go
package cli

import (
	"context"

	"github.com/taku-o/go-webdb-template/internal/service"
)

// ServerStatusServiceInterface はServerStatusServiceのインターフェース
type ServerStatusServiceInterface interface {
	ListServerStatus(ctx context.Context, servers []service.ServerInfo) ([]service.ServerStatus, error)
}

// ServerStatusUsecase はCLI用のサーバー状態確認usecase
type ServerStatusUsecase struct {
	serverStatusService ServerStatusServiceInterface
}

// NewServerStatusUsecase は新しいServerStatusUsecaseを作成
func NewServerStatusUsecase(serverStatusService ServerStatusServiceInterface) *ServerStatusUsecase {
	return &ServerStatusUsecase{
		serverStatusService: serverStatusService,
	}
}

// getServers は確認対象のサーバーリストを返す
func (u *ServerStatusUsecase) getServers() []service.ServerInfo {
	return []service.ServerInfo{
		{Name: "API", Port: 8080, Address: "localhost"},
		{Name: "Client", Port: 3000, Address: "localhost"},
		{Name: "Admin", Port: 8081, Address: "localhost"},
		{Name: "JobQueue", Port: 8082, Address: "localhost"},
		{Name: "PostgreSQL", Port: 5432, Address: "localhost"},
		{Name: "MySQL", Port: 3306, Address: "localhost"},
		{Name: "Redis", Port: 6379, Address: "localhost"},
		{Name: "Redis Cluster", Port: 7100, Address: "localhost"},
		{Name: "Mailpit", Port: 8025, Address: "localhost"},
		{Name: "CloudBeaver", Port: 8978, Address: "localhost"},
		{Name: "Superset", Port: 8088, Address: "localhost"},
		{Name: "Metabase", Port: 8970, Address: "localhost"},
		{Name: "Redis Insight", Port: 8001, Address: "localhost"},
	}
}

// ListServerStatus は全サーバーの状態を確認して返す
func (u *ServerStatusUsecase) ListServerStatus(ctx context.Context) ([]service.ServerStatus, error) {
	servers := u.getServers()
	return u.serverStatusService.ListServerStatus(ctx, servers)
}
```

#### 3.3.4 実装の詳細
- **サーバー情報の定義**: `getServers()`メソッドで定義（既存の`servers`変数の内容を移行）
- **Service層の呼び出し**: `serverStatusService.ListServerStatus(ctx, servers)`を呼び出す（サーバー一覧をパラメータとして渡す）
- **結果の返却**: service層から受け取った結果をそのまま返す（変換不要）
- **エラーハンドリング**: service層のエラーをそのまま返す（必要に応じてコンテキストを追加）

### 3.4 CLI層の実装

#### 3.4.1 ファイル構成

**ファイル**: `server/cmd/server-status/main.go`（変更）

#### 3.4.2 実装コード

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/taku-o/go-webdb-template/internal/service"
	"github.com/taku-o/go-webdb-template/internal/usecase/cli"
)

// printResults は結果を表形式で表示する
func printResults(results []service.ServerStatus) {
	// ヘッダー行
	fmt.Println("サーバー          | ポート | 状態")
	fmt.Println("------------------|-------|--------")

	// 各サーバーの状態を表示
	for _, result := range results {
		fmt.Printf("%-17s | %-5d | %s\n",
			result.Server.Name,
			result.Server.Port,
			result.Status,
		)
	}
}

func main() {
	// Service層の初期化
	serverStatusService := service.NewServerStatusService()

	// Usecase層の初期化
	serverStatusUsecase := cli.NewServerStatusUsecase(serverStatusService)

	// サーバー状態の確認
	ctx := context.Background()
	results, err := serverStatusUsecase.ListServerStatus(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// 結果を表形式で表示
	printResults(results)

	os.Exit(0)
}
```

#### 3.4.3 実装の詳細
- **Usecase層の呼び出し**: `serverStatusUsecase.ListServerStatus()`を呼び出す
- **結果の表示**: `printResults()`関数で表形式に表示（既存のロジックを維持）
- **エラーハンドリング**: エラーが発生した場合は標準エラー出力に表示し、終了コード1で終了

### 3.5 型のエクスポート

#### 3.5.1 Service層の型エクスポート
- `ServerInfo`: エクスポート（usecase層やCLI層で使用する可能性があるため）
- `ServerStatus`: エクスポート（usecase層やCLI層で使用するため）
- `ServerStatusService`: エクスポート（usecase層で使用するため）

#### 3.5.2 Usecase層の型エクスポート
- `ServerStatusUsecase`: エクスポート（CLI層で使用するため）
- `ServerStatusServiceInterface`: エクスポート（テストで使用するため）

## 4. テスト設計

### 4.1 Service層のテスト

#### 4.1.1 テストファイル

**ファイル**: `server/internal/service/server_status_service_test.go`

#### 4.1.2 テスト内容

```go
package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServerStatusService_ListServerStatus(t *testing.T) {
	tests := []struct {
		name    string
		servers []ServerInfo
		wantErr bool
	}{
		{
			name: "正常系: サーバーリストの状態を確認",
			servers: []ServerInfo{
				{Name: "Test", Port: 99999, Address: "localhost"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewServerStatusService()
			ctx := context.Background()

			results, err := service.ListServerStatus(ctx, tt.servers)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, results)
				assert.Equal(t, len(tt.servers), len(results))
			}
		})
	}
}

func TestServerStatusService_checkServerStatus(t *testing.T) {
	tests := []struct {
		name     string
		server   ServerInfo
		want     string
		wantErr  bool
	}{
		{
			name:    "停止中のサーバー",
			server:  ServerInfo{Name: "Test", Port: 99999, Address: "localhost"},
			want:    "停止中",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewServerStatusService()
			timeout := 1 * time.Second

			result := service.checkServerStatus(tt.server, timeout)

			assert.Equal(t, tt.server, result.Server)
			assert.Equal(t, tt.want, result.Status)
		})
	}
}
```

#### 4.1.3 テストの詳細
- **統合テスト的なアプローチ**: 実際のポート接続をテスト（モックサーバーは使用しない）
- **テーブル駆動テスト**: 既存のテストパターンに従う
- **アサーション**: `github.com/stretchr/testify/assert`を使用

### 4.2 Usecase層のテスト

#### 4.2.1 テストファイル

**ファイル**: `server/internal/usecase/cli/server_status_usecase_test.go`

#### 4.2.2 テスト内容

```go
package cli

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taku-o/go-webdb-template/internal/service"
)

// MockServerStatusService はServerStatusServiceInterfaceのモック
type MockServerStatusService struct {
	ListServerStatusFunc func(ctx context.Context, servers []service.ServerInfo) ([]service.ServerStatus, error)
}

func (m *MockServerStatusService) ListServerStatus(ctx context.Context, servers []service.ServerInfo) ([]service.ServerStatus, error) {
	if m.ListServerStatusFunc != nil {
		return m.ListServerStatusFunc(ctx, servers)
	}
	return nil, nil
}

func TestServerStatusUsecase_ListServerStatus(t *testing.T) {
	tests := []struct {
		name    string
		mock    *MockServerStatusService
		wantErr bool
	}{
		{
			name: "正常系: serviceから結果を受け取る",
			mock: &MockServerStatusService{
				ListServerStatusFunc: func(ctx context.Context, servers []service.ServerInfo) ([]service.ServerStatus, error) {
					return []service.ServerStatus{
						{
							Server: service.ServerInfo{Name: "Test", Port: 8080, Address: "localhost"},
							Status: "起動中",
							Error:  nil,
						},
					}, nil
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := NewServerStatusUsecase(tt.mock)
			ctx := context.Background()

			results, err := usecase.ListServerStatus(ctx)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, results)
			}
		})
	}
}
```

#### 4.2.3 テストの詳細
- **モックの使用**: Service層をモック化してテスト
- **テーブル駆動テスト**: 既存のテストパターンに従う
- **アサーション**: `github.com/stretchr/testify/assert`を使用

## 5. 実装上の注意事項

### 5.1 アーキテクチャパターンの遵守

#### 5.1.1 3層構造
- **CLI層**: 入出力制御のみを担当
- **Usecase層**: serviceにパラメータを渡し、結果を受け取る
- **Service層**: ビジネスロジックを実装

#### 5.1.2 依存関係の方向
- CLI層がUsecase層に依存
- Usecase層がService層に依存
- Service層は他の層に依存しない

#### 5.1.3 インターフェースの使用
- Service層のインターフェースはUsecase層で定義（既存パターンに従う）
- テスト容易性を確保するため、インターフェースを使用

### 5.2 既存ロジックの移行

#### 5.2.1 サーバー情報の定義
- **現在**: `server/cmd/server-status/main.go`の`servers`変数
- **変更後**: `server/internal/usecase/cli/server_status_usecase.go`の`getServers()`メソッド

#### 5.2.2 状態確認ロジック
- **現在**: `server/cmd/server-status/main.go`の`checkServerStatus`関数
- **変更後**: `server/internal/service/server_status_service.go`の`checkServerStatus()`メソッド

#### 5.2.3 並列実行ロジック
- **現在**: `server/cmd/server-status/main.go`の`checkAllServers`関数
- **変更後**: `server/internal/service/server_status_service.go`の`checkAllServers()`メソッド

#### 5.2.4 表示ロジック
- **現在**: `server/cmd/server-status/main.go`の`printResults`関数
- **変更後**: `server/cmd/server-status/main.go`に維持（CLI層の責務）

### 5.3 エラーハンドリング

#### 5.3.1 Service層
- エラーをそのまま返す（コンテキストを追加しない）

#### 5.3.2 Usecase層
- Service層のエラーをそのまま返す（必要に応じてコンテキストを追加）

#### 5.3.3 CLI層
- エラーが発生した場合は標準エラー出力に表示し、終了コード1で終了

### 5.4 命名規則

#### 5.4.1 パッケージ名
- Service層: `service`
- Usecase層: `cli`
- CLI層: `main`

#### 5.4.2 型名
- Service層: `ServerStatusService`
- Usecase層: `ServerStatusUsecase`
- インターフェース: `ServerStatusServiceInterface`

#### 5.4.3 メソッド名
- Service層: `ListServerStatus()`, `checkServerStatus()`, `checkAllServers()`, `getServers()`
- Usecase層: `ListServerStatus()`

## 6. 動作確認設計

### 6.1 ローカル環境での動作確認

#### 6.1.1 コマンド実行

```bash
# 開発環境での実行
go run ./server/cmd/server-status/main.go
```

#### 6.1.2 期待される結果

```
サーバー          | ポート | 状態
------------------|-------|--------
API              | 8080  | 起動中
Client           | 3000  | 起動中
Admin            | 8081  | 停止中
JobQueue         | 8082  | 起動中
PostgreSQL       | 5432  | 起動中
MySQL            | 3306  | 停止中
Redis            | 6379  | 起動中
Redis Cluster    | 7100  | 停止中
Mailpit          | 8025  | 起動中
CloudBeaver      | 8978  | 停止中
Superset         | 8088  | 起動中
Metabase         | 8970  | 停止中
Redis Insight    | 8001  | 起動中
```

### 6.2 テスト実行

#### 6.2.1 Service層のテスト

```bash
APP_ENV=test go test ./server/internal/service/server_status_service_test.go
```

#### 6.2.2 Usecase層のテスト

```bash
APP_ENV=test go test ./server/internal/usecase/cli/server_status_usecase_test.go
```

### 6.3 既存機能への影響確認

- 既存のサーバー状態確認機能が正常に動作することを確認
- 既存の表示形式（表形式）が維持されていることを確認
- 既存の13個のサーバーが確認対象として維持されていることを確認
- 並列実行による状態確認が正常に動作することを確認

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### 変更するファイル
- `server/cmd/server-status/main.go`: 入出力制御のみを担当するように変更

#### 新規作成するファイル
- `server/internal/usecase/cli/server_status_usecase.go`: CLI用usecase層（新規作成）
- `server/internal/usecase/cli/server_status_usecase_test.go`: usecase層のテスト（新規作成）
- `server/internal/service/server_status_service.go`: service層（新規作成）
- `server/internal/service/server_status_service_test.go`: service層のテスト（新規作成）

### 7.2 既存機能への影響
- **既存のサーバー**: 影響なし（ポート接続のみで判定、ロジックは維持）
- **既存の機能**: 影響なし（リファクタリングのみ、機能は維持）
- **既存のCLIコマンド**: 影響なし（独立した実装）

### 7.3 テストへの影響
- **既存のテスト**: 影響なし（新規ファイルの追加のみ）
- **新規テスト**: usecase層とservice層のテストを追加

## 8. 参考情報

### 8.1 既存実装の参考
- **CLI実装例**: `server/cmd/list-dm-users/main.go`
- **Usecase実装例**: `server/internal/usecase/cli/list_dm_users_usecase.go`
- **Service実装例**: `server/internal/service/dm_user_service.go`

### 8.2 アーキテクチャドキュメント
- `.kiro/steering/tech.md`: 技術スタックとアーキテクチャの詳細
- `.kiro/steering/structure.md`: ファイル組織とコードパターン

### 8.3 関連ドキュメント
- `.kiro/specs/0079-status-usecase/requirements.md`: 本機能の要件定義書
- Issue #161: 本機能の元となる要望

### 8.4 技術スタック
- **言語**: Go 1.21+
- **標準ライブラリ**: `net`, `sync`, `time`, `fmt`, `context`
- **テスト**: `testing`（標準ライブラリ）、`github.com/stretchr/testify`（アサーション、モック）
- **外部依存**: なし
