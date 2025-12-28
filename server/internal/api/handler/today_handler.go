package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/example/go-webdb-template/internal/auth"
	humaapi "github.com/example/go-webdb-template/internal/api/huma"
)

// TodayHandler は今日の日付APIのハンドラー
type TodayHandler struct{}

// NewTodayHandler は新しいTodayHandlerを作成
func NewTodayHandler() *TodayHandler {
	return &TodayHandler{}
}

// GetToday は今日の日付を取得
func (h *TodayHandler) GetToday(ctx context.Context) (string, error) {
	// 公開レベルのチェック（privateエンドポイント）
	if err := auth.CheckAccessLevel(ctx, auth.AccessLevelPrivate); err != nil {
		return "", err
	}

	// 今日の日付をYYYY-MM-DD形式で返す
	return time.Now().Format("2006-01-02"), nil
}

// RegisterTodayEndpoints はHuma APIにToday APIエンドポイントを登録
func RegisterTodayEndpoints(api huma.API, h *TodayHandler) {
	// GET /api/today - 今日の日付取得（private API）
	huma.Register(api, huma.Operation{
		OperationID: "get-today",
		Method:      http.MethodGet,
		Path:        "/api/today",
		Summary:     "[private] 今日の日付を取得",
		Description: "**Access Level:** `private` (Auth0 JWT でアクセス可能)",
		Tags:        []string{"today"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, input *humaapi.GetTodayInput) (*humaapi.TodayOutput, error) {
		date, err := h.GetToday(ctx)
		if err != nil {
			return nil, huma.Error403Forbidden(err.Error())
		}

		resp := &humaapi.TodayOutput{}
		resp.Body.Date = date
		return resp, nil
	})
}
