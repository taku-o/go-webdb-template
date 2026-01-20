# サーバー状態確認機能のリファクタリング実装タスク一覧

## 概要
`server/cmd/server-status/main.go`を3層構造（CLI → Usecase → Service）にリファクタリングするための実装タスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: Service層の実装

#### - [ ] タスク 1.1: `server/internal/service/server_status_service.go`の作成
**目的**: サーバー状態確認のビジネスロジックを実装するService層を作成する。

**作業内容**:
- `server/internal/service/server_status_service.go`ファイルを作成
- パッケージ宣言とインポートを追加
- `ServerInfo`構造体を定義
- `ServerStatus`構造体を定義
- `ServerStatusService`構造体を定義
- `NewServerStatusService`関数を実装
- `checkServerStatus`メソッドを実装（既存の`checkServerStatus`関数の内容を移行）
- `checkAllServers`メソッドを実装（既存の`checkAllServers`関数の内容を移行）
- `ListServerStatus`メソッドを実装（サーバー一覧をパラメータとして受け取る）

**実装コード**:
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

**受け入れ基準**:
- `server/internal/service/server_status_service.go`ファイルが作成されている
- `ServerInfo`構造体が定義されている
- `ServerStatus`構造体が定義されている
- `ServerStatusService`構造体が定義されている
- `NewServerStatusService`関数が実装されている
- `checkServerStatus`メソッドが実装されている（既存のロジックを維持）
- `checkAllServers`メソッドが実装されている（既存のロジックを維持）
  - 各goroutineは`results[index]`に直接書き込むため、元の`servers`の順序が保たれる
- `ListServerStatus`メソッドが実装されている（サーバー一覧をパラメータとして受け取る）

- _Requirements: 3.1.3, 6.1_
- _Design: 3.2.1, 3.2.2, 3.2.3, 3.2.4_

---

#### - [ ] タスク 1.2: Service層のテスト実装
**目的**: Service層の単体テストを実装する。

**作業内容**:
- `server/internal/service/server_status_service_test.go`ファイルを作成
- `TestServerStatusService_ListServerStatus`テストを実装
- `TestServerStatusService_checkServerStatus`テストを実装
- テーブル駆動テストのパターンを使用
- `github.com/stretchr/testify/assert`を使用

**実装コード**:
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

**受け入れ基準**:
- `server/internal/service/server_status_service_test.go`ファイルが作成されている
- `TestServerStatusService_ListServerStatus`テストが実装されている
- `TestServerStatusService_checkServerStatus`テストが実装されている
- テーブル駆動テストのパターンが使用されている
- `APP_ENV=test go test ./server/internal/service/server_status_service_test.go`が正常に実行される

- _Requirements: 6.4_
- _Design: 4.1.1, 4.1.2, 4.1.3_

---

### Phase 2: Usecase層の実装

#### - [ ] タスク 2.1: `server/internal/usecase/cli/server_status_usecase.go`の作成
**目的**: CLI用のUsecase層を作成し、Service層を呼び出す実装を行う。

**作業内容**:
- `server/internal/usecase/cli/server_status_usecase.go`ファイルを作成
- パッケージ宣言とインポートを追加
- `ServerStatusServiceInterface`インターフェースを定義
- `ServerStatusUsecase`構造体を定義
- `NewServerStatusUsecase`関数を実装
- `getServers`メソッドを実装（サーバー一覧を定義）
- `ListServerStatus`メソッドを実装（サーバー一覧を取得してService層に渡す）

**実装コード**:
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

**受け入れ基準**:
- `server/internal/usecase/cli/server_status_usecase.go`ファイルが作成されている
- `ServerStatusServiceInterface`インターフェースが定義されている
- `ServerStatusUsecase`構造体が定義されている
- `NewServerStatusUsecase`関数が実装されている
- `getServers`メソッドが実装されている（13個のサーバーが定義されている）
- `ListServerStatus`メソッドが実装されている（サーバー一覧を取得してService層に渡す）

- _Requirements: 3.1.2, 6.1_
- _Design: 3.3.1, 3.3.2, 3.3.3, 3.3.4_

---

#### - [ ] タスク 2.2: Usecase層のテスト実装
**目的**: Usecase層の単体テストを実装する。

**作業内容**:
- `server/internal/usecase/cli/server_status_usecase_test.go`ファイルを作成
- `MockServerStatusService`を実装
- `TestServerStatusUsecase_ListServerStatus`テストを実装
- テーブル駆動テストのパターンを使用
- `github.com/stretchr/testify/assert`を使用

**実装コード**:
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

**受け入れ基準**:
- `server/internal/usecase/cli/server_status_usecase_test.go`ファイルが作成されている
- `MockServerStatusService`が実装されている
- `TestServerStatusUsecase_ListServerStatus`テストが実装されている
- テーブル駆動テストのパターンが使用されている
- `APP_ENV=test go test ./server/internal/usecase/cli/server_status_usecase_test.go`が正常に実行される

- _Requirements: 6.4_
- _Design: 4.2.1, 4.2.2, 4.2.3_

---

### Phase 3: CLI層のリファクタリング

#### - [ ] タスク 3.1: `server/cmd/server-status/main.go`のリファクタリング
**目的**: CLI層を入出力制御のみを担当するようにリファクタリングする。

**作業内容**:
- `server/cmd/server-status/main.go`を開く
- 既存の`servers`変数を削除（Usecase層に移行）
- 既存の`checkServerStatus`関数を削除（Service層に移行）
- 既存の`checkAllServers`関数を削除（Service層に移行）
- `connectionTimeout`定数を削除（Service層に移行）
- `ServerInfo`構造体を削除（Service層に移行）
- `ServerStatus`構造体を削除（Service層に移行）
- `printResults`関数は維持（CLI層の責務）
- `main`関数をリファクタリング（Usecase層を呼び出すように変更）
- 必要なインポートを追加（`context`, `os`, `service`, `usecase/cli`）

**実装コード**:
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

**受け入れ基準**:
- `server/cmd/server-status/main.go`が入出力制御のみを担当している
- 既存の`servers`変数が削除されている
- 既存の`checkServerStatus`関数が削除されている
- 既存の`checkAllServers`関数が削除されている
- `connectionTimeout`定数が削除されている
- `ServerInfo`構造体が削除されている
- `ServerStatus`構造体が削除されている
- `printResults`関数が維持されている
- `main`関数がUsecase層を呼び出すように変更されている
- 必要なインポートが追加されている

- _Requirements: 3.1.1, 6.1_
- _Design: 3.4.1, 3.4.2, 3.4.3_

---

### Phase 4: 動作確認

#### - [ ] タスク 4.1: コマンド実行の確認
**目的**: リファクタリング後のコマンドが正常に実行されることを確認する。

**作業内容**:
- `go run ./server/cmd/server-status/main.go`を実行
- エラーが発生しないことを確認
- 表形式で結果が表示されることを確認
- 既存の表示形式が維持されていることを確認

**確認コマンド**:
```bash
go run ./server/cmd/server-status/main.go
```

**期待される結果**:
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

**受け入れ基準**:
- コマンド実行が正常に完了する
- エラーが発生しない
- 表形式で結果が表示される
- 既存の表示形式が維持されている
- 13個のサーバーが表示される
- サーバーが指定された順序で表示される

- _Requirements: 6.2, 6.3_
- _Design: 6.1.1_

---

#### - [ ] タスク 4.2: 既存機能の動作確認
**目的**: 既存のサーバー状態確認機能が正常に動作することを確認する。

**作業内容**:
- 全てのサーバーが起動している場合、全て「起動中」と表示されることを確認
- 一部のサーバーが停止している場合、該当サーバーが「停止中」と表示されることを確認
- 全てのサーバーが停止している場合、全て「停止中」と表示されることを確認
- 並列実行による状態確認が正常に動作することを確認

**受け入れ基準**:
- 全てのサーバーが起動している場合、全て「起動中」と表示される
- 一部のサーバーが停止している場合、該当サーバーが「停止中」と表示される
- 全てのサーバーが停止している場合、全て「停止中」と表示される
- 並列実行による状態確認が正常に動作する

- _Requirements: 6.2, 6.3_

---

### Phase 5: テスト実行

#### - [ ] タスク 5.1: 全テストの実行
**目的**: 全てのテストが正常に実行されることを確認する。

**作業内容**:
- `APP_ENV=test go test ./server/internal/service/server_status_service_test.go`を実行
- `APP_ENV=test go test ./server/internal/usecase/cli/server_status_usecase_test.go`を実行
- 既存のテストが全て失敗しないことを確認（`APP_ENV=test go test ./...`）

**確認コマンド**:
```bash
# Service層のテスト
APP_ENV=test go test ./server/internal/service/server_status_service_test.go

# Usecase層のテスト
APP_ENV=test go test ./server/internal/usecase/cli/server_status_usecase_test.go

# 全テストの実行
APP_ENV=test go test ./...
```

**受け入れ基準**:
- Service層のテストが正常に実行される
- Usecase層のテストが正常に実行される
- 既存のテストが全て失敗しない
- テストカバレッジが適切に維持されている

- _Requirements: 6.4_

---

### Phase 6: コード品質確認

#### - [ ] タスク 6.1: アーキテクチャパターンの確認
**目的**: プロジェクトの標準的なアーキテクチャパターンに従っていることを確認する。

**作業内容**:
- 3層構造（CLI → Usecase → Service）が正しく実装されていることを確認
- 依存関係の方向が正しいことを確認（CLI層がUsecase層に依存、Usecase層がService層に依存）
- インターフェースが適切に使用されていることを確認
- 既存のCLI実装（`list_dm_users_usecase.go`など）と同様の構造になっていることを確認

**受け入れ基準**:
- 3層構造が正しく実装されている
- 依存関係の方向が正しい
- インターフェースが適切に使用されている
- 既存のCLI実装と同様の構造になっている

- _Requirements: 6.5_

---

#### - [ ] タスク 6.2: コードスタイルの確認
**目的**: コードスタイルが既存のコードと一致していることを確認する。

**作業内容**:
- 命名規則が既存のコードと一致していることを確認
- コメントが適切に追加されていることを確認
- インポート順序が既存のコードと一致していることを確認

**受け入れ基準**:
- 命名規則が既存のコードと一致している
- コメントが適切に追加されている
- インポート順序が既存のコードと一致している

- _Requirements: 6.5_

---

## 受け入れ基準の確認

### レイヤー分離
- [ ] `server/cmd/server-status/main.go`が入出力制御のみを担当している
- [ ] `server/internal/usecase/cli/server_status_usecase.go`が作成されている
- [ ] `server/internal/service/server_status_service.go`が作成されている
- [ ] 各層の責務が明確に分離されている

### 機能の維持
- [ ] 既存のサーバー状態確認機能が正常に動作する
- [ ] 既存の表示形式（表形式）が維持されている
- [ ] 既存の13個のサーバーが確認対象として維持されている
- [ ] 並列実行による状態確認が正常に動作する

### 動作確認
- [ ] コマンド実行が正常に完了する（`go run ./cmd/server-status/main.go`）
- [ ] 全てのサーバーが起動している場合、全て「起動中」と表示される
- [ ] 一部のサーバーが停止している場合、該当サーバーが「停止中」と表示される
- [ ] 全てのサーバーが停止している場合、全て「停止中」と表示される
- [ ] サーバーが指定された順序で表示される

### テスト
- [ ] usecase層の単体テストが実装されている
- [ ] service層の単体テストが実装されている
- [ ] 既存のテストが全て失敗しないことを確認
- [ ] テストカバレッジが適切に維持されている

### コード品質
- [ ] プロジェクトの標準的なアーキテクチャパターンに従っている
- [ ] 既存のCLI実装（`list_dm_users_usecase.go`など）と同様の構造になっている
- [ ] コードスタイルが既存のコードと一致している
- [ ] 適切なコメントが追加されている
