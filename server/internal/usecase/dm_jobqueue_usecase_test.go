package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockJobQueueClient はテスト用のJobQueueClientモック
type MockJobQueueClient struct {
	EnqueueJobFunc func(ctx context.Context, jobType string, payload []byte, opts *JobOptions) (*JobInfo, error)
}

func (m *MockJobQueueClient) EnqueueJob(ctx context.Context, jobType string, payload []byte, opts *JobOptions) (*JobInfo, error) {
	if m.EnqueueJobFunc != nil {
		return m.EnqueueJobFunc(ctx, jobType, payload, opts)
	}
	return nil, nil
}

func TestDmJobqueueUsecase_RegisterJob(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		message        string
		delaySeconds   int
		maxRetry       int
		mockClient     JobQueueClientInterface
		wantJobID      string
		wantErr        bool
		wantErrMessage string
	}{
		{
			name:         "registers job successfully",
			message:      "Test message",
			delaySeconds: 10,
			maxRetry:     3,
			mockClient: &MockJobQueueClient{
				EnqueueJobFunc: func(ctx context.Context, jobType string, payload []byte, opts *JobOptions) (*JobInfo, error) {
					return &JobInfo{ID: "job-123"}, nil
				},
			},
			wantJobID: "job-123",
			wantErr:   false,
		},
		{
			name:         "uses default message when empty",
			message:      "",
			delaySeconds: 0,
			maxRetry:     0,
			mockClient: &MockJobQueueClient{
				EnqueueJobFunc: func(ctx context.Context, jobType string, payload []byte, opts *JobOptions) (*JobInfo, error) {
					// payloadにデフォルトメッセージが含まれていることを確認
					assert.Contains(t, string(payload), "Job executed successfully")
					return &JobInfo{ID: "job-456"}, nil
				},
			},
			wantJobID: "job-456",
			wantErr:   false,
		},
		{
			name:         "returns error when client is nil",
			message:      "Test message",
			delaySeconds: 10,
			maxRetry:     3,
			mockClient:   nil, // interface型のnil
			wantJobID:    "",
			wantErr:      true,
		},
		{
			name:         "returns error when enqueue fails",
			message:      "Test message",
			delaySeconds: 10,
			maxRetry:     3,
			mockClient: &MockJobQueueClient{
				EnqueueJobFunc: func(ctx context.Context, jobType string, payload []byte, opts *JobOptions) (*JobInfo, error) {
					return nil, errors.New("enqueue failed")
				},
			},
			wantJobID: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := NewDmJobqueueUsecase(tt.mockClient)

			jobID, err := usecase.RegisterJob(ctx, tt.message, tt.delaySeconds, tt.maxRetry)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantJobID, jobID)
			}
		})
	}
}
