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
