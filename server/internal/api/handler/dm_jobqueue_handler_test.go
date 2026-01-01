package handler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taku-o/go-webdb-template/internal/config"
	"github.com/taku-o/go-webdb-template/internal/service/jobqueue"
)

func TestNewDmJobqueueHandler(t *testing.T) {
	// nilのクライアントでもハンドラーが作成できること
	handler := NewDmJobqueueHandler(nil)
	assert.NotNil(t, handler)
}

func TestNewDmJobqueueHandler_WithClient(t *testing.T) {
	// クライアントを渡してハンドラーを作成
	cfg := &config.Config{
		CacheServer: config.CacheServerConfig{
			Redis: config.RedisConfig{
				JobQueue: config.RedisSingleConfig{
					Addr: "",
				},
			},
		},
	}

	client, err := jobqueue.NewClient(cfg)
	assert.NoError(t, err)
	defer client.Close()

	handler := NewDmJobqueueHandler(client)
	assert.NotNil(t, handler)
}

func TestDmJobqueueHandler_RegisterJob_NilClient(t *testing.T) {
	// nilのクライアントでジョブ登録を試みる
	handler := NewDmJobqueueHandler(nil)

	req := &RegisterJobRequest{
		Message: "Test message",
	}

	resp, err := handler.RegisterJob(context.Background(), req)

	// 503エラーが返されること
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Redis is not connected")
}

func TestDmJobqueueHandler_RegisterJob_DefaultMessage(t *testing.T) {
	// 空のメッセージでリクエスト（デフォルトメッセージが使用される）
	// 注意: 実際のRedis接続なしでは、このテストはジョブ登録の前にnilチェックで止まる
	handler := NewDmJobqueueHandler(nil)

	req := &RegisterJobRequest{
		Message: "", // 空のメッセージ
	}

	resp, err := handler.RegisterJob(context.Background(), req)

	// nilクライアントなので503エラー
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestRegisterJobRequest_Fields(t *testing.T) {
	// RegisterJobRequestのフィールドが正しく設定できること
	req := RegisterJobRequest{
		Message:      "Test message",
		DelaySeconds: 60,
		MaxRetry:     5,
	}

	assert.Equal(t, "Test message", req.Message)
	assert.Equal(t, 60, req.DelaySeconds)
	assert.Equal(t, 5, req.MaxRetry)
}

func TestRegisterJobRequest_DefaultValues(t *testing.T) {
	// RegisterJobRequestのデフォルト値（ゼロ値）
	req := RegisterJobRequest{}

	assert.Equal(t, "", req.Message)
	assert.Equal(t, 0, req.DelaySeconds)
	assert.Equal(t, 0, req.MaxRetry)
}

func TestRegisterJobResponse_Fields(t *testing.T) {
	// RegisterJobResponseのフィールドが正しく設定できること
	resp := RegisterJobResponse{
		JobID:  "test-job-id",
		Status: "registered",
	}

	assert.Equal(t, "test-job-id", resp.JobID)
	assert.Equal(t, "registered", resp.Status)
}

func TestRegisterJobInput_BodyField(t *testing.T) {
	// RegisterJobInputのBodyフィールドが正しく設定できること
	input := RegisterJobInput{
		Body: RegisterJobRequest{
			Message:      "Test",
			DelaySeconds: 30,
		},
	}

	assert.Equal(t, "Test", input.Body.Message)
	assert.Equal(t, 30, input.Body.DelaySeconds)
}

func TestRegisterJobOutput_BodyField(t *testing.T) {
	// RegisterJobOutputのBodyフィールドが正しく設定できること
	output := RegisterJobOutput{
		Body: RegisterJobResponse{
			JobID:  "abc123",
			Status: "registered",
		},
	}

	assert.Equal(t, "abc123", output.Body.JobID)
	assert.Equal(t, "registered", output.Body.Status)
}
