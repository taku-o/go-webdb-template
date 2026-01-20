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
		name    string
		server  ServerInfo
		want    string
		wantErr bool
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
