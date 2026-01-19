# 起動サーバー一覧表示機能の設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、13個のサーバーの起動状態を確認し、表形式で表示するコンソールプログラム`server-status`の詳細設計を定義する。Goプログラムとして実装し、ポートへのTCP接続によりサーバーの起動状態を判定する。

### 1.2 設計の範囲
- 13個のサーバー（API、Client、Admin、JobQueue、PostgreSQL、MySQL、Redis、Redis Cluster、Mailpit、CloudBeaver、Superset、Metabase、Redis Insight）の状態確認機能
- ポートへのTCP接続による起動状態の判定
- 表形式での結果表示
- 並列実行による高速化
- エラーハンドリング

### 1.3 設計方針
- **Goプログラムでの実装**: プロジェクトの技術スタックに一致し、クロスプラットフォーム対応
- **シンプルな実装**: 標準ライブラリを中心に使用し、外部依存を最小限に
- **並列実行**: goroutineを使用して複数サーバーの状態を並列に確認し、実行時間を短縮
- **明確な出力**: 表形式で見やすく、指定された順序で表示
- **エラーハンドリング**: 接続エラーを適切に処理し、停止中として表示

## 2. アーキテクチャ設計

### 2.1 プログラム構成

```
┌─────────────────────────────────────────────────────────────┐
│              server-status プログラム                          │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  1. サーバー定義の読み込み          │
        │     - サーバー名、ポート、確認方法   │
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  2. 並列状態確認                  │
        │     - goroutineで各サーバー確認   │
        │     - TCP接続試行                │
        │     - タイムアウト: 1秒           │
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  3. 結果の収集                    │
        │     - 各goroutineの結果を収集    │
        │     - 指定順序でソート            │
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  4. 表形式での表示                │
        │     - サーバー名、ポート、状態     │
        │     - 日本語で表示               │
        └─────────────────────────────────┘
```

### 2.2 状態確認フロー

```
┌─────────────────────────────────────────────────────────────┐
│              各サーバーの状態確認フロー                         │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  1. TCP接続の試行                │
        │     - net.DialTimeout()          │
        │     - タイムアウト: 1秒           │
        └─────────────────────────────────┘
                          │
        ┌─────────────────┴─────────────────┐
        │                                   │
        ▼                                   ▼
┌──────────────────┐            ┌──────────────────┐
│  接続成功         │            │  接続失敗         │
│  - 接続を即座に閉じる│            │  - 接続拒否       │
│  - 起動中と判定    │            │  - タイムアウト   │
│                   │            │  - その他のエラー │
│                   │            │  - 停止中と判定   │
└──────────────────┘            └──────────────────┘
```

### 2.3 並列実行フロー

```
┌─────────────────────────────────────────────────────────────┐
│              並列状態確認の実行フロー                          │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  全サーバー定義をループ            │
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  各サーバーに対してgoroutine起動  │
        │  - checkServerStatus()           │
        └─────────────────────────────────┘
                          │
        ┌─────────────────┴─────────────────┐
        │                                   │
        ▼                                   ▼
┌──────────────────┐            ┌──────────────────┐
│  goroutine 1     │            │  goroutine 2     │
│  (API)           │            │  (Client)         │
└──────────────────┘            └──────────────────┘
        │                                   │
        ▼                                   ▼
┌──────────────────┐            ┌──────────────────┐
│  goroutine 3      │            │  ...             │
│  (Admin)         │            │  goroutine 13    │
└──────────────────┘            └──────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  全goroutineの完了を待機          │
        │  - sync.WaitGroup                │
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  結果を指定順序でソート          │
        │  - サーバー定義の順序を維持       │
        └─────────────────────────────────┘
```

## 3. 実装設計

### 3.1 サーバー定義

#### 3.1.1 サーバー情報の構造体

**ファイル**: `server/cmd/server-status/main.go`

```go
// ServerInfo はサーバー情報を表す
type ServerInfo struct {
	Name    string // サーバー名
	Port    int    // ポート番号
	Address string // 接続先アドレス（通常は"localhost"）
}

// CheckMethod は確認方法を表す
type CheckMethod int

const (
	CheckMethodTCP CheckMethod = iota // TCP接続のみ
	CheckMethodHTTP                    // HTTPエンドポイント（将来の拡張用）
)
```

#### 3.1.2 サーバー定義リスト

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

### 3.2 状態確認の実装

#### 3.2.1 状態確認関数

```go
// ServerStatus はサーバーの状態を表す
type ServerStatus struct {
	Server ServerInfo
	Status string // "起動中" または "停止中"
	Error  error  // エラー情報（デバッグ用、表示には使用しない）
}

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

#### 3.2.2 並列実行の実装

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

### 3.3 結果表示の実装

#### 3.3.1 表形式表示関数

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

#### 3.3.2 列幅の調整

表形式の見やすさを確保するため、以下の列幅を使用：
- **サーバー名**: 17文字（左揃え）
- **ポート**: 5文字（右揃え）
- **状態**: 可変（「起動中」または「停止中」）

### 3.4 メイン関数の実装

#### 3.4.1 実装コード

```go
package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	// connectionTimeout はTCP接続のタイムアウト時間
	connectionTimeout = 1 * time.Second
)

func main() {
	// 全サーバーの状態を並列に確認
	results := checkAllServers(servers, connectionTimeout)
	
	// 結果を表形式で表示
	printResults(results)
	
	// 正常終了
	os.Exit(0)
}
```

#### 3.4.2 エラーハンドリング

- **接続エラー**: 停止中として表示（エラーメッセージは表示しない）
- **タイムアウト**: 停止中として表示
- **その他のエラー**: 停止中として表示

エラーの詳細は表示せず、状態のみを表示する（要件定義書に基づく）。

## 4. ファイル構成

### 4.1 新規作成ファイル

```
server/cmd/server-status/
├── main.go          # メイン実装ファイル
└── main_test.go     # テストファイル
```

### 4.2 ファイル内容

#### 4.2.1 main.go

- サーバー定義リスト
- 状態確認関数
- 並列実行関数
- 結果表示関数
- メイン関数

#### 4.2.2 main_test.go

- 状態確認関数のテスト
- 並列実行のテスト
- 結果表示のテスト

## 5. 非機能要件の実装

### 5.1 パフォーマンス

#### 5.1.1 並列実行
- goroutineを使用して全サーバーを並列に確認
- 実行時間は最も遅いサーバーの確認時間 + オーバーヘッド（約1秒）

#### 5.1.2 タイムアウト設定
- 各サーバーの確認は1秒でタイムアウト
- 全サーバーの確認は約1秒以内に完了（並列実行のため）

### 5.2 可用性

#### 5.2.1 エラーハンドリング
- 接続エラー、タイムアウト、その他のエラーを適切に処理
- エラーが発生してもプログラムは正常終了し、停止中として表示

#### 5.2.2 ネットワークエラー
- ネットワークエラーが発生した場合も停止中として表示
- プログラムはクラッシュせず、正常に終了

### 5.3 保守性

#### 5.3.1 コードの簡潔性
- 標準ライブラリのみを使用
- シンプルで理解しやすい実装
- コメントを適切に追加

#### 5.3.2 拡張性
- サーバー定義リストに追加するだけで新しいサーバーに対応可能
- 確認方法の拡張（HTTPエンドポイントなど）も容易

## 6. テスト設計

### 6.1 単体テスト

#### 6.1.1 状態確認関数のテスト
- 起動中のサーバーへの接続テスト
- 停止中のサーバーへの接続テスト
- タイムアウトのテスト

#### 6.1.2 結果表示関数のテスト
- 表形式の出力テスト
- 列幅の確認テスト

### 6.2 統合テスト

#### 6.2.1 並列実行のテスト
- 複数サーバーの並列確認テスト
- 実行時間の確認テスト

### 6.3 テストの実装方針

- モックサーバーを使用してテスト
- 実際のサーバーが起動していない環境でもテスト可能にする

## 7. 実装上の注意事項

### 7.1 ポート確認の注意事項

#### 7.1.1 タイムアウト設定
- 1秒のタイムアウトを設定し、応答のないサーバーを適切に検出
- タイムアウトが短すぎると誤検出の可能性があるため、1秒を維持

#### 7.1.2 接続の即座クローズ
- 接続が成功した場合は即座に閉じる
- リソースリークを防ぐ

#### 7.1.3 エラーの種類
- `ECONNREFUSED`: 接続拒否（サーバーが停止中）
- `timeout`: タイムアウト（サーバーが応答しない）
- その他のエラー: ネットワークエラーなど

### 7.2 表示の注意事項

#### 7.2.1 表形式の整列
- 列幅を適切に設定し、見やすく表示
- 日本語文字の幅を考慮

#### 7.2.2 順序の維持
- サーバー定義リストの順序を維持
- 並列実行の結果も定義順序で表示

### 7.3 実装方式の注意事項

#### 7.3.1 標準ライブラリの使用
- `net`パッケージを使用してTCP接続
- 外部依存を避ける

#### 7.3.2 クロスプラットフォーム対応
- macOS、Linuxで動作することを確認
- Windowsはオプション（要件定義書に記載）

## 8. 実行例

### 8.1 実行コマンド

```bash
# 開発環境での実行
go run ./server/cmd/server-status/main.go

# ビルドして実行
go build -o server-status ./server/cmd/server-status/main.go
./server-status
```

### 8.2 出力例

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

## 9. 参考情報

### 9.1 技術スタック
- **言語**: Go 1.21+
- **標準ライブラリ**: `net`, `sync`, `time`, `fmt`
- **外部依存**: なし

### 9.2 関連ドキュメント
- 要件定義書: `.kiro/specs/0077-listapp/requirements.md`
- サーバー構成: `.kiro/steering/tech.md`
- Issue #157: https://github.com/taku-o/go-webdb-template/issues/157

### 9.3 既存実装の参考
- `server/cmd/list-dm-users/main.go`: コマンドラインツールの実装パターン
- `server/cmd/generate-sample-data/main.go`: コマンドラインツールの実装パターン
