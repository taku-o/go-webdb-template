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

// UserHandler はユーザーAPIのハンドラー
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler は新しいUserHandlerを作成
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// RegisterUserEndpoints はHuma APIにユーザーエンドポイントを登録
func RegisterUserEndpoints(api huma.API, h *UserHandler) {
	// POST /api/users - ユーザー作成
	huma.Register(api, huma.Operation{
		OperationID:   "create-user",
		Method:        http.MethodPost,
		Path:          "/api/users",
		Summary:       "ユーザーを作成",
		Description:   "**Access Level:** `public` (Public API Key JWT または Auth0 JWT でアクセス可能)",
		Tags:          []string{"users"},
		DefaultStatus: http.StatusCreated,
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, input *humaapi.CreateUserInput) (*humaapi.UserOutput, error) {
		// 公開レベルのチェック（publicエンドポイント）
		if err := auth.CheckAccessLevel(ctx, auth.AccessLevelPublic); err != nil {
			return nil, huma.Error403Forbidden(err.Error())
		}

		req := &model.CreateUserRequest{
			Name:  input.Body.Name,
			Email: input.Body.Email,
		}

		user, err := h.userService.CreateUser(ctx, req)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		resp := &humaapi.UserOutput{}
		resp.Body = *user
		return resp, nil
	})

	// GET /api/users/{id} - ユーザー取得
	huma.Register(api, huma.Operation{
		OperationID: "get-user",
		Method:      http.MethodGet,
		Path:        "/api/users/{id}",
		Summary:     "ユーザーを取得",
		Description: "**Access Level:** `public` (Public API Key JWT または Auth0 JWT でアクセス可能)",
		Tags:        []string{"users"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, input *humaapi.GetUserInput) (*humaapi.UserOutput, error) {
		// 公開レベルのチェック（publicエンドポイント）
		if err := auth.CheckAccessLevel(ctx, auth.AccessLevelPublic); err != nil {
			return nil, huma.Error403Forbidden(err.Error())
		}

		user, err := h.userService.GetUser(ctx, input.ID)
		if err != nil {
			return nil, huma.Error404NotFound(err.Error())
		}

		resp := &humaapi.UserOutput{}
		resp.Body = *user
		return resp, nil
	})

	// GET /api/users - ユーザー一覧取得
	huma.Register(api, huma.Operation{
		OperationID: "list-users",
		Method:      http.MethodGet,
		Path:        "/api/users",
		Summary:     "ユーザー一覧を取得",
		Description: "**Access Level:** `public` (Public API Key JWT または Auth0 JWT でアクセス可能)",
		Tags:        []string{"users"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, input *humaapi.ListUsersInput) (*humaapi.UsersOutput, error) {
		// 公開レベルのチェック（publicエンドポイント）
		if err := auth.CheckAccessLevel(ctx, auth.AccessLevelPublic); err != nil {
			return nil, huma.Error403Forbidden(err.Error())
		}

		users, err := h.userService.ListUsers(ctx, input.Limit, input.Offset)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		resp := &humaapi.UsersOutput{}
		resp.Body = users
		return resp, nil
	})

	// PUT /api/users/{id} - ユーザー更新
	huma.Register(api, huma.Operation{
		OperationID: "update-user",
		Method:      http.MethodPut,
		Path:        "/api/users/{id}",
		Summary:     "ユーザーを更新",
		Description: "**Access Level:** `public` (Public API Key JWT または Auth0 JWT でアクセス可能)",
		Tags:        []string{"users"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, input *humaapi.UpdateUserInput) (*humaapi.UserOutput, error) {
		// 公開レベルのチェック（publicエンドポイント）
		if err := auth.CheckAccessLevel(ctx, auth.AccessLevelPublic); err != nil {
			return nil, huma.Error403Forbidden(err.Error())
		}

		req := &model.UpdateUserRequest{
			Name:  input.Body.Name,
			Email: input.Body.Email,
		}

		user, err := h.userService.UpdateUser(ctx, input.ID, req)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		resp := &humaapi.UserOutput{}
		resp.Body = *user
		return resp, nil
	})

	// DELETE /api/users/{id} - ユーザー削除
	huma.Register(api, huma.Operation{
		OperationID:   "delete-user",
		Method:        http.MethodDelete,
		Path:          "/api/users/{id}",
		Summary:       "ユーザーを削除",
		Description:   "**Access Level:** `public` (Public API Key JWT または Auth0 JWT でアクセス可能)",
		Tags:          []string{"users"},
		DefaultStatus: http.StatusNoContent,
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, input *humaapi.DeleteUserInput) (*struct{}, error) {
		// 公開レベルのチェック（publicエンドポイント）
		if err := auth.CheckAccessLevel(ctx, auth.AccessLevelPublic); err != nil {
			return nil, huma.Error403Forbidden(err.Error())
		}

		err := h.userService.DeleteUser(ctx, input.ID)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		return nil, nil
	})
}
