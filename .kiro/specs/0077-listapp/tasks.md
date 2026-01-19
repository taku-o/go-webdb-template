# 起動サーバー一覧表示機能の実装タスク一覧

## 概要
13個のサーバーの起動状態を確認し、表形式で表示するコンソールプログラム`server-status`の実装タスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: プロジェクト構造の準備

#### - [ ] タスク 1.1: `server/cmd/server-status`ディレクトリの作成
**目的**: server-statusプログラム用のディレクトリを作成する。

**作業内容**:
- `server/cmd/server-status`ディレクトリを作成
- ディレクトリが正しく作成されたことを確認

**受け入れ基準**:
- `server/cmd/server-status`ディレクトリが存在する

- _Requirements: 7.1_
- _Design: 4.1_

---

### Phase 2: サーバー定義の実装

#### - [ ] タスク 2.1: `main.go`ファイルの作成と基本構造の実装
**目的**: メインファイルを作成し、パッケージ宣言とインポートを実装する。

**作業内容**:
- `server/cmd/server-status/main.go`ファイルを作成
- パッケージ宣言を追加
- 必要な標準ライブラリのインポートを追加

**実装コード**:
```go
package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)
```

**受け入れ基準**:
- `main.go`ファイルが作成されている
- パッケージ宣言が正しい
- 必要なインポートが追加されている

- _Requirements: 7.1_
- _Design: 3.4.1_

---

#### - [ ] タスク 2.2: サーバー情報構造体の定義
**目的**: サーバー情報を表す構造体を定義する。

**作業内容**:
- `ServerInfo`構造体を定義
- フィールド: `Name`, `Port`, `Address`

**実装コード**:
```go
// ServerInfo はサーバー情報を表す
type ServerInfo struct {
	Name    string // サーバー名
	Port    int    // ポート番号
	Address string // 接続先アドレス（通常は"localhost"）
}
```

**受け入れ基準**:
- `ServerInfo`構造体が定義されている
- 必要なフィールドがすべて含まれている

- _Requirements: 3.1.1_
- _Design: 3.1.1_

---

#### - [ ] タスク 2.3: サーバー状態構造体の定義
**目的**: サーバーの状態を表す構造体を定義する。

**作業内容**:
- `ServerStatus`構造体を定義
- フィールド: `Server`, `Status`, `Error`

**実装コード**:
```go
// ServerStatus はサーバーの状態を表す
type ServerStatus struct {
	Server ServerInfo
	Status string // "起動中" または "停止中"
	Error  error  // エラー情報（デバッグ用、表示には使用しない）
}
```

**受け入れ基準**:
- `ServerStatus`構造体が定義されている
- 必要なフィールドがすべて含まれている

- _Requirements: 3.1.2_
- _Design: 3.2.1_

---

#### - [ ] タスク 2.4: サーバー定義リストの実装
**目的**: 確認対象の13個のサーバーを定義する。

**作業内容**:
- `servers`変数を定義
- 13個のサーバー（API、Client、Admin、JobQueue、PostgreSQL、MySQL、Redis、Redis Cluster、Mailpit、CloudBeaver、Superset、Metabase、Redis Insight）を指定された順序で定義

**実装コード**:
```go
// servers は確認対象のサーバーリスト（指定された順序）
var servers = []ServerInfo{
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
```

**受け入れ基準**:
- `servers`変数が定義されている
- 13個のサーバーがすべて含まれている
- サーバーが指定された順序で定義されている
- 各サーバーのポート番号が正しい

- _Requirements: 3.1.1, 6.3_
- _Design: 3.1.2_

---

### Phase 3: 状態確認機能の実装

#### - [ ] タスク 3.1: タイムアウト定数の定義
**目的**: TCP接続のタイムアウト時間を定数として定義する。

**作業内容**:
- `connectionTimeout`定数を定義
- 値は1秒に設定

**実装コード**:
```go
const (
	// connectionTimeout はTCP接続のタイムアウト時間
	connectionTimeout = 1 * time.Second
)
```

**受け入れ基準**:
- `connectionTimeout`定数が定義されている
- 値が1秒に設定されている

- _Requirements: 4.1_
- _Design: 3.4.1_

---

#### - [ ] タスク 3.2: 状態確認関数の実装
**目的**: 指定されたサーバーの状態をTCP接続で確認する関数を実装する。

**作業内容**:
- `checkServerStatus`関数を実装
- `net.DialTimeout`を使用してTCP接続を試行
- 接続成功時は即座に接続を閉じる
- 接続失敗時は停止中として返す

**実装コード**:
```go
// checkServerStatus は指定されたサーバーの状態を確認する
func checkServerStatus(server ServerInfo, timeout time.Duration) ServerStatus {
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
```

**受け入れ基準**:
- `checkServerStatus`関数が実装されている
- TCP接続を試行している
- 接続成功時に接続を閉じている
- 接続失敗時に停止中として返している
- エラーハンドリングが適切に実装されている

- _Requirements: 3.1.2, 3.1.3, 4.2, 6.1_
- _Design: 3.2.1_

---

#### - [ ] タスク 3.3: 並列実行関数の実装
**目的**: 全サーバーの状態を並列に確認する関数を実装する。

**作業内容**:
- `checkAllServers`関数を実装
- `sync.WaitGroup`を使用して並列実行を管理
- goroutineで各サーバーの状態を確認
- 結果を指定順序で返す

**実装コード**:
```go
// checkAllServers は全サーバーの状態を並列に確認する
func checkAllServers(servers []ServerInfo, timeout time.Duration) []ServerStatus {
	var wg sync.WaitGroup
	results := make([]ServerStatus, len(servers))
	
	for i, server := range servers {
		wg.Add(1)
		go func(index int, s ServerInfo) {
			defer wg.Done()
			results[index] = checkServerStatus(s, timeout)
		}(i, server)
	}
	
	wg.Wait()
	return results
}
```

**受け入れ基準**:
- `checkAllServers`関数が実装されている
- `sync.WaitGroup`を使用して並列実行を管理している
- goroutineで各サーバーの状態を確認している
- 結果が指定順序で返されている
- 全goroutineの完了を待機している

- _Requirements: 4.1, 8.1_
- _Design: 3.2.2_

---

### Phase 4: 結果表示機能の実装

#### - [ ] タスク 4.1: 結果表示関数の実装
**目的**: サーバー状態の結果を表形式で表示する関数を実装する。

**作業内容**:
- `printResults`関数を実装
- ヘッダー行を表示
- 各サーバーの状態を表形式で表示
- 列幅を適切に設定

**実装コード**:
```go
// printResults は結果を表形式で表示する
func printResults(results []ServerStatus) {
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
```

**受け入れ基準**:
- `printResults`関数が実装されている
- ヘッダー行が表示されている
- 各サーバーの状態が表形式で表示されている
- 列幅が適切に設定されている
- 日本語で表示されている

- _Requirements: 3.2.1, 3.2.2, 3.2.3, 6.2_
- _Design: 3.3.1_

---

### Phase 5: メイン関数の実装

#### - [ ] タスク 5.1: メイン関数の実装
**目的**: プログラムのエントリーポイントを実装する。

**作業内容**:
- `main`関数を実装
- `checkAllServers`を呼び出して全サーバーの状態を確認
- `printResults`を呼び出して結果を表示
- 正常終了

**実装コード**:
```go
func main() {
	// 全サーバーの状態を並列に確認
	results := checkAllServers(servers, connectionTimeout)
	
	// 結果を表形式で表示
	printResults(results)
	
	// 正常終了
	os.Exit(0)
}
```

**受け入れ基準**:
- `main`関数が実装されている
- `checkAllServers`を呼び出している
- `printResults`を呼び出している
- 正常終了している

- _Requirements: 3.3.3, 6.3_
- _Design: 3.4.1_

---

### Phase 6: テストの実装

#### - [ ] タスク 6.1: `main_test.go`ファイルの作成
**目的**: テストファイルを作成する。

**作業内容**:
- `server/cmd/server-status/main_test.go`ファイルを作成
- パッケージ宣言を追加
- 必要なインポートを追加

**実装コード**:
```go
package main

import (
	"net"
	"testing"
	"time"
)
```

**受け入れ基準**:
- `main_test.go`ファイルが作成されている
- パッケージ宣言が正しい
- 必要なインポートが追加されている

- _Requirements: 6.4_
- _Design: 4.2.2, 6.1_

---

#### - [ ] タスク 6.2: モックサーバーの実装
**目的**: テスト用のモックサーバーを実装する。

**作業内容**:
- テスト用のTCPサーバーを起動する関数を実装
- テスト用のTCPサーバーを停止する関数を実装
- ポート番号を動的に割り当て

**実装コード**:
```go
// startMockServer はテスト用のTCPサーバーを起動する
func startMockServer(t *testing.T) (string, func()) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to start mock server: %v", err)
	}
	
	port := listener.Addr().(*net.TCPAddr).Port
	address := fmt.Sprintf("localhost:%d", port)
	
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	
	stop := func() {
		listener.Close()
	}
	
	return address, stop
}
```

**受け入れ基準**:
- モックサーバーを起動する関数が実装されている
- モックサーバーを停止する関数が実装されている
- ポート番号が動的に割り当てられている

- _Requirements: 6.4_
- _Design: 6.3_

---

#### - [ ] タスク 6.3: 状態確認関数のテスト
**目的**: `checkServerStatus`関数のテストを実装する。

**作業内容**:
- 起動中のサーバーへの接続テスト
- 停止中のサーバーへの接続テスト
- タイムアウトのテスト

**実装コード**:
```go
func TestCheckServerStatus(t *testing.T) {
	t.Run("起動中のサーバー", func(t *testing.T) {
		address, stop := startMockServer(t)
		defer stop()
		
		// サーバーが起動するまで少し待つ
		time.Sleep(100 * time.Millisecond)
		
		host, port, _ := net.SplitHostPort(address)
		portInt := 0
		fmt.Sscanf(port, "%d", &portInt)
		
		server := ServerInfo{
			Name:    "TestServer",
			Port:    portInt,
			Address: host,
		}
		
		result := checkServerStatus(server, 1*time.Second)
		
		if result.Status != "起動中" {
			t.Errorf("Expected status '起動中', got '%s'", result.Status)
		}
	})
	
	t.Run("停止中のサーバー", func(t *testing.T) {
		server := ServerInfo{
			Name:    "TestServer",
			Port:    99999, // 使用されていないポート
			Address: "localhost",
		}
		
		result := checkServerStatus(server, 1*time.Second)
		
		if result.Status != "停止中" {
			t.Errorf("Expected status '停止中', got '%s'", result.Status)
		}
	})
}
```

**受け入れ基準**:
- 起動中のサーバーへの接続テストが実装されている
- 停止中のサーバーへの接続テストが実装されている
- テストが正常に実行される

- _Requirements: 6.4_
- _Design: 6.1.1_

---

#### - [ ] タスク 6.4: 並列実行関数のテスト
**目的**: `checkAllServers`関数のテストを実装する。

**作業内容**:
- 複数サーバーの並列確認テスト
- 実行時間の確認テスト

**実装コード**:
```go
func TestCheckAllServers(t *testing.T) {
	t.Run("並列実行の確認", func(t *testing.T) {
		testServers := []ServerInfo{
			{Name: "Server1", Port: 99991, Address: "localhost"},
			{Name: "Server2", Port: 99992, Address: "localhost"},
			{Name: "Server3", Port: 99993, Address: "localhost"},
		}
		
		start := time.Now()
		results := checkAllServers(testServers, 1*time.Second)
		duration := time.Since(start)
		
		if len(results) != len(testServers) {
			t.Errorf("Expected %d results, got %d", len(testServers), len(results))
		}
		
		// 並列実行のため、順次実行より速いはず（ただし、すべて停止中なのでタイムアウトで1秒かかる）
		if duration > 2*time.Second {
			t.Errorf("Parallel execution took too long: %v", duration)
		}
	})
}
```

**受け入れ基準**:
- 複数サーバーの並列確認テストが実装されている
- 実行時間の確認テストが実装されている
- テストが正常に実行される

- _Requirements: 6.4_
- _Design: 6.2.1_

---

#### - [ ] タスク 6.5: 結果表示関数のテスト
**目的**: `printResults`関数のテストを実装する。

**作業内容**:
- 表形式の出力テスト
- 列幅の確認テスト

**実装コード**:
```go
func TestPrintResults(t *testing.T) {
	results := []ServerStatus{
		{
			Server: ServerInfo{Name: "API", Port: 8080, Address: "localhost"},
			Status: "起動中",
		},
		{
			Server: ServerInfo{Name: "Client", Port: 3000, Address: "localhost"},
			Status: "停止中",
		},
	}
	
	// 出力をキャプチャして確認（簡易的な実装）
	// 実際のテストでは、出力をキャプチャして検証する
	printResults(results)
	
	// このテストは主にコンパイルエラーがないことを確認する
}
```

**受け入れ基準**:
- 表形式の出力テストが実装されている
- テストが正常に実行される

- _Requirements: 6.4_
- _Design: 6.1.2_

---

### Phase 7: ドキュメントの更新

#### - [ ] タスク 7.1: `docs/ja/Command-Line-Tool.md`の更新
**目的**: 日本語版のCLIツールドキュメントにserver-statusコマンドの情報を追加する。

**作業内容**:
- `docs/ja/Command-Line-Tool.md`を開く
- server-statusコマンドのセクションを追加
- 概要、使用方法、実行例、出力形式を記載

**追加内容**:
- server-statusコマンドの概要
- ビルド方法
- 実行方法
- 出力形式の説明
- 実行例

**受け入れ基準**:
- server-statusコマンドの情報が追加されている
- 他のコマンド（list-dm-users、generate-sample-data）と同様の形式で記載されている
- 日本語で記載されている

- _Requirements: -_
- _Design: -_

---

#### - [ ] タスク 7.2: `docs/en/Command-Line-Tool.md`の更新
**目的**: 英語版のCLIツールドキュメントにserver-statusコマンドの情報を追加する。

**作業内容**:
- `docs/en/Command-Line-Tool.md`を開く
- server-statusコマンドのセクションを追加
- 概要、使用方法、実行例、出力形式を記載

**追加内容**:
- server-statusコマンドの概要
- ビルド方法
- 実行方法
- 出力形式の説明
- 実行例

**受け入れ基準**:
- server-statusコマンドの情報が追加されている
- 他のコマンド（list-dm-users、generate-sample-data）と同様の形式で記載されている
- 英語で記載されている

- _Requirements: -_
- _Design: -_

---

#### - [ ] タスク 7.3: `README.md`の更新
**目的**: 英語版のREADMEにserver-statusコマンドの情報を追加する。

**作業内容**:
- `README.md`を開く
- CLIツールのセクションにserver-statusコマンドを追加
- ビルド方法、実行方法、出力形式を記載

**追加内容**:
- server-statusコマンドのセクションを追加
- ビルド方法
- 実行方法
- 出力形式の説明

**受け入れ基準**:
- server-statusコマンドの情報が追加されている
- 他のCLIツールと同様の形式で記載されている
- 英語で記載されている

- _Requirements: -_
- _Design: -_

---

#### - [ ] タスク 7.4: `README.ja.md`の更新
**目的**: 日本語版のREADMEにserver-statusコマンドの情報を追加する。

**作業内容**:
- `README.ja.md`を開く
- CLIツールのセクションにserver-statusコマンドを追加
- ビルド方法、実行方法、出力形式を記載

**追加内容**:
- server-statusコマンドのセクションを追加
- ビルド方法
- 実行方法
- 出力形式の説明

**受け入れ基準**:
- server-statusコマンドの情報が追加されている
- 他のCLIツールと同様の形式で記載されている
- 日本語で記載されている

- _Requirements: -_
- _Design: -_

---

### Phase 8: 動作確認

#### - [ ] タスク 8.1: プログラムのビルド確認
**目的**: プログラムが正常にビルドできることを確認する。

**作業内容**:
- `go build ./server/cmd/server-status`を実行
- ビルドエラーがないことを確認

**受け入れ基準**:
- プログラムが正常にビルドできる
- ビルドエラーがない

- _Requirements: 6.4_
- _Design: 8.1_

---

#### - [ ] タスク 8.2: プログラムの実行確認
**目的**: プログラムが正常に実行できることを確認する。

**作業内容**:
- `go run ./server/cmd/server-status/main.go`を実行
- エラーなく実行できることを確認
- 出力が表形式で表示されることを確認

**受け入れ基準**:
- プログラムが正常に実行できる
- エラーなく実行できる
- 出力が表形式で表示される

- _Requirements: 6.3_
- _Design: 8.1_

---

#### - [ ] タスク 8.3: 全サーバー状態確認の動作確認
**目的**: 全13個のサーバーの状態が正しく確認できることを確認する。

**作業内容**:
- 一部のサーバーを起動
- プログラムを実行
- 起動中のサーバーが「起動中」と表示されることを確認
- 停止中のサーバーが「停止中」と表示されることを確認
- サーバーが指定された順序で表示されることを確認

**受け入れ基準**:
- 起動中のサーバーが「起動中」と表示される
- 停止中のサーバーが「停止中」と表示される
- サーバーが指定された順序で表示される
- 全13個のサーバーが表示される

- _Requirements: 6.1, 6.2, 6.3_
- _Design: 8.2_

---

#### - [ ] タスク 8.4: テストの実行確認
**目的**: すべてのテストが正常に実行されることを確認する。

**作業内容**:
- `APP_ENV=test go test ./server/cmd/server-status`を実行
- すべてのテストがパスすることを確認

**受け入れ基準**:
- すべてのテストが正常に実行される
- すべてのテストがパスする
- テストエラーがない

- _Requirements: 6.4_
- _Design: 6.1, 6.2_

---

#### - [ ] タスク 8.5: 既存テストへの影響確認
**目的**: 既存のテストが失敗しないことを確認する。

**作業内容**:
- `APP_ENV=test go test ./...`を実行
- 既存のテストがすべてパスすることを確認

**受け入れ基準**:
- 既存のテストがすべてパスする
- 既存のテストに影響がない

- _Requirements: 7.3_
- _Design: -_

---

## タスクの依存関係

```
Phase 1 (プロジェクト構造の準備)
  └─> Phase 2 (サーバー定義の実装)
      └─> Phase 3 (状態確認機能の実装)
          └─> Phase 4 (結果表示機能の実装)
              └─> Phase 5 (メイン関数の実装)
                  └─> Phase 6 (テストの実装)
                      ├─> Phase 7 (ドキュメントの更新)
                      └─> Phase 8 (動作確認)
```

## 実装の優先順位

1. **高優先度**: Phase 1-5（基本機能の実装）
2. **中優先度**: Phase 6（テストの実装）、Phase 7（ドキュメントの更新）
3. **低優先度**: Phase 8（動作確認）

## 注意事項

- 各タスクは独立して実装可能だが、依存関係に注意する
- テストは実装と並行して進めることも可能
- ドキュメント更新は実装が完了してから実施する（実装内容を確認してから記載する）
- 動作確認は実装とドキュメント更新が完了してから実施する
