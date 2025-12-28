package handler

import (
	"context"
	"testing"
	"time"

	"github.com/taku-o/go-webdb-template/internal/auth"
	"github.com/stretchr/testify/assert"
)

func TestTodayHandler_GetToday(t *testing.T) {
	handler := NewTodayHandler()
	assert.NotNil(t, handler)
}

func TestTodayHandler_GetTodayResponse(t *testing.T) {
	handler := NewTodayHandler()

	// privateアクセスレベルのコンテキストを作成
	ctx := context.WithValue(context.Background(), auth.AllowedAccessLevelKey, auth.AccessLevelPrivate)

	date, err := handler.GetToday(ctx)
	assert.NoError(t, err)

	// 今日の日付が返されることを確認
	expected := time.Now().Format("2006-01-02")
	assert.Equal(t, expected, date)
}

func TestTodayHandler_GetTodayWithPublicAccessLevel(t *testing.T) {
	handler := NewTodayHandler()

	// publicアクセスレベルのコンテキストを作成
	ctx := context.WithValue(context.Background(), auth.AllowedAccessLevelKey, auth.AccessLevelPublic)

	_, err := handler.GetToday(ctx)
	// publicアクセスレベルではprivate APIにアクセスできないためエラー
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "private API requires Auth0 authentication")
}

func TestTodayHandler_GetTodayWithoutAccessLevel(t *testing.T) {
	handler := NewTodayHandler()

	// アクセスレベルがないコンテキスト
	ctx := context.Background()

	_, err := handler.GetToday(ctx)
	// アクセスレベルがない場合はエラー
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "access level not found in context")
}
