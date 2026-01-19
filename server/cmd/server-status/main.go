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

func main() {
	// 全サーバーの状態を並列に確認
	results := checkAllServers(servers, connectionTimeout)

	// 結果を表形式で表示
	printResults(results)
}
