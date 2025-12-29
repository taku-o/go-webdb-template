package handler

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	humaapi "github.com/taku-o/go-webdb-template/internal/api/huma"
	"github.com/taku-o/go-webdb-template/internal/auth"
	"github.com/taku-o/go-webdb-template/internal/model"
	"github.com/taku-o/go-webdb-template/internal/service"
)

// DmUserHandler はユーザーAPIのハンドラー
type DmUserHandler struct {
	dmUserService *service.DmUserService
}

// NewDmUserHandler は新しいDmUserHandlerを作成
func NewDmUserHandler(dmUserService *service.DmUserService) *DmUserHandler {
	return &DmUserHandler{
		dmUserService: dmUserService,
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

		dmUser, err := h.dmUserService.CreateDmUser(ctx, req)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		resp := &humaapi.DmUserOutput{}
		resp.Body = *dmUser
		return resp, nil
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

		dmUser, err := h.dmUserService.GetDmUser(ctx, input.ID)
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

		dmUsers, err := h.dmUserService.ListDmUsers(ctx, input.Limit, input.Offset)
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

		req := &model.UpdateDmUserRequest{
			Name:  input.Body.Name,
			Email: input.Body.Email,
		}

		dmUser, err := h.dmUserService.UpdateDmUser(ctx, input.ID, req)
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

		err := h.dmUserService.DeleteDmUser(ctx, input.ID)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		return nil, nil
	})
}
