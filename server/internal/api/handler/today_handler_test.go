package handler

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/taku-o/go-webdb-template/internal/auth"
	usecaseapi "github.com/taku-o/go-webdb-template/internal/usecase/api"
)

// MockDateService はDateServiceのモック（usecase用）
type MockDateService struct {
	GetTodayFunc func(ctx context.Context) (string, error)
}

func (m *MockDateService) GetToday(ctx context.Context) (string, error) {
	if m.GetTodayFunc != nil {
		return m.GetTodayFunc(ctx)
	}
	return time.Now().Format("2006-01-02"), nil
}

// createTodayHandlerWithMock はモックを使用してTodayHandlerを作成するヘルパー関数
func createTodayHandlerWithMock() *TodayHandler {
	mockDateService := &MockDateService{}
	todayUsecase := usecaseapi.NewTodayUsecase(mockDateService)
	return NewTodayHandler(todayUsecase)
}

func TestTodayHandler_GetToday(t *testing.T) {
	handler := createTodayHandlerWithMock()
	assert.NotNil(t, handler)
}

func TestTodayHandler_GetTodayResponse(t *testing.T) {
	handler := createTodayHandlerWithMock()

	// privateアクセスレベルのコンテキストを作成
	ctx := context.WithValue(context.Background(), auth.AllowedAccessLevelKey, auth.AccessLevelPrivate)

	date, err := handler.GetToday(ctx)
	assert.NoError(t, err)

	// 今日の日付が返されることを確認
	expected := time.Now().Format("2006-01-02")
	assert.Equal(t, expected, date)
}

func TestTodayHandler_GetTodayWithPublicAccessLevel(t *testing.T) {
	handler := createTodayHandlerWithMock()

	// publicアクセスレベルのコンテキストを作成
	ctx := context.WithValue(context.Background(), auth.AllowedAccessLevelKey, auth.AccessLevelPublic)

	_, err := handler.GetToday(ctx)
	// publicアクセスレベルではprivate APIにアクセスできないためエラー
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "private API requires Auth0 authentication")
}

func TestTodayHandler_GetTodayWithoutAccessLevel(t *testing.T) {
	handler := createTodayHandlerWithMock()

	// アクセスレベルがないコンテキスト
	ctx := context.Background()

	_, err := handler.GetToday(ctx)
	// アクセスレベルがない場合はエラー
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "access level not found in context")
}
