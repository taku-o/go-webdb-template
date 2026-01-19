package main

import (
	"fmt"
	"net"
	"testing"
	"time"
)

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

	t.Run("順序の維持", func(t *testing.T) {
		testServers := []ServerInfo{
			{Name: "Server1", Port: 99991, Address: "localhost"},
			{Name: "Server2", Port: 99992, Address: "localhost"},
			{Name: "Server3", Port: 99993, Address: "localhost"},
		}

		results := checkAllServers(testServers, 1*time.Second)

		for i, result := range results {
			if result.Server.Name != testServers[i].Name {
				t.Errorf("Expected server name '%s' at index %d, got '%s'",
					testServers[i].Name, i, result.Server.Name)
			}
		}
	})
}

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

func TestServersDefinition(t *testing.T) {
	t.Run("サーバー数の確認", func(t *testing.T) {
		expectedCount := 13
		if len(servers) != expectedCount {
			t.Errorf("Expected %d servers, got %d", expectedCount, len(servers))
		}
	})

	t.Run("サーバー順序の確認", func(t *testing.T) {
		expectedOrder := []string{
			"API",
			"Client",
			"Admin",
			"JobQueue",
			"PostgreSQL",
			"MySQL",
			"Redis",
			"Redis Cluster",
			"Mailpit",
			"CloudBeaver",
			"Superset",
			"Metabase",
			"Redis Insight",
		}

		for i, expected := range expectedOrder {
			if servers[i].Name != expected {
				t.Errorf("Expected server at index %d to be '%s', got '%s'",
					i, expected, servers[i].Name)
			}
		}
	})

	t.Run("サーバーポートの確認", func(t *testing.T) {
		expectedPorts := map[string]int{
			"API":           8080,
			"Client":        3000,
			"Admin":         8081,
			"JobQueue":      8082,
			"PostgreSQL":    5432,
			"MySQL":         3306,
			"Redis":         6379,
			"Redis Cluster": 7100,
			"Mailpit":       8025,
			"CloudBeaver":   8978,
			"Superset":      8088,
			"Metabase":      8970,
			"Redis Insight": 8001,
		}

		for _, server := range servers {
			expectedPort, ok := expectedPorts[server.Name]
			if !ok {
				t.Errorf("Unexpected server name: %s", server.Name)
				continue
			}
			if server.Port != expectedPort {
				t.Errorf("Expected port %d for server '%s', got %d",
					expectedPort, server.Name, server.Port)
			}
		}
	})
}
