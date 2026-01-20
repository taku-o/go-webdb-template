package cli

import (
	"github.com/taku-o/go-webdb-template/internal/service"
)

// ServerStatusServiceInterface はServerStatusServiceのインターフェース
type ServerStatusServiceInterface interface {
	ListServerStatus(servers []service.ServerInfo) ([]service.ServerStatus, error)
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
func (u *ServerStatusUsecase) ListServerStatus() ([]service.ServerStatus, error) {
	servers := u.getServers()
	return u.serverStatusService.ListServerStatus(servers)
}
