package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taku-o/go-webdb-template/internal/service"
)

// MockServerStatusService はServerStatusServiceInterfaceのモック
type MockServerStatusService struct {
	ListServerStatusFunc func(servers []service.ServerInfo) ([]service.ServerStatus, error)
}

func (m *MockServerStatusService) ListServerStatus(servers []service.ServerInfo) ([]service.ServerStatus, error) {
	if m.ListServerStatusFunc != nil {
		return m.ListServerStatusFunc(servers)
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
				ListServerStatusFunc: func(servers []service.ServerInfo) ([]service.ServerStatus, error) {
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

			results, err := usecase.ListServerStatus()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, results)
			}
		})
	}
}
