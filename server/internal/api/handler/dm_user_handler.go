package handler

import (
	"context"
	"encoding/csv"
	"log"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	humaapi "github.com/taku-o/go-webdb-template/internal/api/huma"
	"github.com/taku-o/go-webdb-template/internal/auth"
	"github.com/taku-o/go-webdb-template/internal/model"
	usecaseapi "github.com/taku-o/go-webdb-template/internal/usecase/api"
)

// DmUserHandler はユーザーAPIのハンドラー
type DmUserHandler struct {
	dmUserUsecase *usecaseapi.DmUserUsecase
}

// NewDmUserHandler は新しいDmUserHandlerを作成
func NewDmUserHandler(dmUserUsecase *usecaseapi.DmUserUsecase) *DmUserHandler {
	return &DmUserHandler{
		dmUserUsecase: dmUserUsecase,
	}
}

// RegisterDmUserEndpoints はHuma APIにユーザーエンドポイントを登録
func RegisterDmUserEndpoints(api huma.API, h *DmUserHandler) {
	// POST /api/dm-users - ユーザー作成
	huma.Register(api, huma.Operation{
		OperationID:   "create-user",
		Method:        http.MethodPost,
		Path:          "/api/dm-users",
		Summary:       "ユーザーを作成",
		Description:   "**Access Level:** `public` (Public API Key JWT または Auth0 JWT でアクセス可能)",
		Tags:          []string{"users"},
		DefaultStatus: http.StatusCreated,
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, input *humaapi.CreateDmUserInput) (*humaapi.DmUserOutput, error) {
		// 公開レベルのチェック（publicエンドポイント）
		if err := auth.CheckAccessLevel(ctx, auth.AccessLevelPublic); err != nil {
			return nil, huma.Error403Forbidden(err.Error())
		}

		req := &model.CreateDmUserRequest{
			Name:  input.Body.Name,
			Email: input.Body.Email,
		}

		dmUser, err := h.dmUserUsecase.CreateDmUser(ctx, req)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		resp := &humaapi.DmUserOutput{}
		resp.Body = *dmUser
		return resp, nil
	})

	// GET /api/export/dm-users/csv - ユーザー情報をCSV形式でダウンロード
	huma.Register(api, huma.Operation{
		OperationID: "download-users-csv",
		Method:      http.MethodGet,
		Path:        "/api/export/dm-users/csv",
		Summary:     "ユーザー情報をCSV形式でダウンロード",
		Description: "**Access Level:** `public` (Public API Key JWT または Auth0 JWT でアクセス可能)",
		Tags:        []string{"users"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, input *struct{}) (*huma.StreamResponse, error) {
		// 公開レベルのチェック（publicエンドポイント）
		if err := auth.CheckAccessLevel(ctx, auth.AccessLevelPublic); err != nil {
			return nil, huma.Error403Forbidden(err.Error())
		}

		// ユーザー情報20件を取得
		users, err := h.dmUserUsecase.ListDmUsers(ctx, 20, 0)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		// ストリーミングレスポンスを返す
		return &huma.StreamResponse{
			Body: func(humaCtx huma.Context) {
				// ヘッダーを設定
				humaCtx.SetHeader("Content-Type", "text/csv; charset=utf-8")
				humaCtx.SetHeader("Content-Disposition", `attachment; filename="dm-users.csv"`)

				// BodyWriterを取得
				w := humaCtx.BodyWriter()

				// http.ResponseWriterを取り出してタイムアウトを設定
				if rw, ok := w.(http.ResponseWriter); ok {
					rc := http.NewResponseController(rw)
					if err := rc.SetWriteDeadline(time.Now().Add(3 * time.Minute)); err != nil {
						log.Printf("Warning: Failed to set write deadline: %v", err)
					}
				}

				// CSVエンコーダーを作成
				csvWriter := csv.NewWriter(w)
				defer csvWriter.Flush()

				// ヘッダー行を書き込み
				if err := csvWriter.Write([]string{
					"ID",
					"Name",
					"Email",
					"Created At",
					"Updated At",
				}); err != nil {
					log.Printf("Error writing CSV header: %v", err)
					return
				}

				// ユーザーデータを1件ずつCSV行として書き込み
				for _, user := range users {
					if err := csvWriter.Write([]string{
						user.ID,
						user.Name,
						user.Email,
						user.CreatedAt.Format(time.RFC3339),
						user.UpdatedAt.Format(time.RFC3339),
					}); err != nil {
						log.Printf("Error writing CSV row: %v", err)
						return
					}
				}
			},
		}, nil
	})

	// GET /api/dm-users/{id} - ユーザー取得
	huma.Register(api, huma.Operation{
		OperationID: "get-user",
		Method:      http.MethodGet,
		Path:        "/api/dm-users/{id}",
		Summary:     "ユーザーを取得",
		Description: "**Access Level:** `public` (Public API Key JWT または Auth0 JWT でアクセス可能)",
		Tags:        []string{"users"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, input *humaapi.GetDmUserInput) (*humaapi.DmUserOutput, error) {
		// 公開レベルのチェック（publicエンドポイント）
		if err := auth.CheckAccessLevel(ctx, auth.AccessLevelPublic); err != nil {
			return nil, huma.Error403Forbidden(err.Error())
		}

		// UUID文字列のバリデーション（32文字であること）
		if len(input.ID) != 32 {
			return nil, huma.Error400BadRequest("invalid id format: must be 32 characters")
		}

		dmUser, err := h.dmUserUsecase.GetDmUser(ctx, input.ID)
		if err != nil {
			return nil, huma.Error404NotFound(err.Error())
		}

		resp := &humaapi.DmUserOutput{}
		resp.Body = *dmUser
		return resp, nil
	})

	// GET /api/dm-users - ユーザー一覧取得
	huma.Register(api, huma.Operation{
		OperationID: "list-users",
		Method:      http.MethodGet,
		Path:        "/api/dm-users",
		Summary:     "ユーザー一覧を取得",
		Description: "**Access Level:** `public` (Public API Key JWT または Auth0 JWT でアクセス可能)",
		Tags:        []string{"users"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, input *humaapi.ListDmUsersInput) (*humaapi.DmUsersOutput, error) {
		// 公開レベルのチェック（publicエンドポイント）
		if err := auth.CheckAccessLevel(ctx, auth.AccessLevelPublic); err != nil {
			return nil, huma.Error403Forbidden(err.Error())
		}

		dmUsers, err := h.dmUserUsecase.ListDmUsers(ctx, input.Limit, input.Offset)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		resp := &humaapi.DmUsersOutput{}
		resp.Body = dmUsers
		return resp, nil
	})

	// PUT /api/dm-users/{id} - ユーザー更新
	huma.Register(api, huma.Operation{
		OperationID: "update-user",
		Method:      http.MethodPut,
		Path:        "/api/dm-users/{id}",
		Summary:     "ユーザーを更新",
		Description: "**Access Level:** `public` (Public API Key JWT または Auth0 JWT でアクセス可能)",
		Tags:        []string{"users"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, input *humaapi.UpdateDmUserInput) (*humaapi.DmUserOutput, error) {
		// 公開レベルのチェック（publicエンドポイント）
		if err := auth.CheckAccessLevel(ctx, auth.AccessLevelPublic); err != nil {
			return nil, huma.Error403Forbidden(err.Error())
		}

		// UUID文字列のバリデーション（32文字であること）
		if len(input.ID) != 32 {
			return nil, huma.Error400BadRequest("invalid id format: must be 32 characters")
		}

		req := &model.UpdateDmUserRequest{
			Name:  input.Body.Name,
			Email: input.Body.Email,
		}

		dmUser, err := h.dmUserUsecase.UpdateDmUser(ctx, input.ID, req)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		resp := &humaapi.DmUserOutput{}
		resp.Body = *dmUser
		return resp, nil
	})

	// DELETE /api/dm-users/{id} - ユーザー削除
	huma.Register(api, huma.Operation{
		OperationID:   "delete-user",
		Method:        http.MethodDelete,
		Path:          "/api/dm-users/{id}",
		Summary:       "ユーザーを削除",
		Description:   "**Access Level:** `public` (Public API Key JWT または Auth0 JWT でアクセス可能)",
		Tags:          []string{"users"},
		DefaultStatus: http.StatusNoContent,
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, input *humaapi.DeleteDmUserInput) (*struct{}, error) {
		// 公開レベルのチェック（publicエンドポイント）
		if err := auth.CheckAccessLevel(ctx, auth.AccessLevelPublic); err != nil {
			return nil, huma.Error403Forbidden(err.Error())
		}

		// UUID文字列のバリデーション（32文字であること）
		if len(input.ID) != 32 {
			return nil, huma.Error400BadRequest("invalid id format: must be 32 characters")
		}

		err := h.dmUserUsecase.DeleteDmUser(ctx, input.ID)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		return nil, nil
	})
}
