package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/danielgtaylor/huma/v2"
	humaapi "github.com/taku-o/go-webdb-template/internal/api/huma"
	"github.com/taku-o/go-webdb-template/internal/auth"
	"github.com/taku-o/go-webdb-template/internal/model"
	"github.com/taku-o/go-webdb-template/internal/service"
)

// DmPostHandler は投稿APIのハンドラー
type DmPostHandler struct {
	dmPostService *service.DmPostService
}

// NewDmPostHandler は新しいDmPostHandlerを作成
func NewDmPostHandler(dmPostService *service.DmPostService) *DmPostHandler {
	return &DmPostHandler{
		dmPostService: dmPostService,
	}
}

// RegisterDmPostEndpoints はHuma APIに投稿エンドポイントを登録
func RegisterDmPostEndpoints(api huma.API, h *DmPostHandler) {
	// POST /api/dm-posts - 投稿作成
	huma.Register(api, huma.Operation{
		OperationID:   "create-post",
		Method:        http.MethodPost,
		Path:          "/api/dm-posts",
		Summary:       "投稿を作成",
		Description:   "**Access Level:** `public` (Public API Key JWT または Auth0 JWT でアクセス可能)",
		Tags:          []string{"posts"},
		DefaultStatus: http.StatusCreated,
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, input *humaapi.CreateDmPostInput) (*humaapi.DmPostOutput, error) {
		// 公開レベルのチェック（publicエンドポイント）
		if err := auth.CheckAccessLevel(ctx, auth.AccessLevelPublic); err != nil {
			return nil, huma.Error403Forbidden(err.Error())
		}

		// string → int64 変換
		userID, err := strconv.ParseInt(input.Body.UserID, 10, 64)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid user_id format")
		}

		req := &model.CreateDmPostRequest{
			UserID:  userID,
			Title:   input.Body.Title,
			Content: input.Body.Content,
		}

		dmPost, err := h.dmPostService.CreateDmPost(ctx, req)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		resp := &humaapi.DmPostOutput{}
		resp.Body = *dmPost
		return resp, nil
	})

	// GET /api/dm-posts/{id} - 投稿取得
	huma.Register(api, huma.Operation{
		OperationID: "get-post",
		Method:      http.MethodGet,
		Path:        "/api/dm-posts/{id}",
		Summary:     "投稿を取得",
		Description: "**Access Level:** `public` (Public API Key JWT または Auth0 JWT でアクセス可能)",
		Tags:        []string{"posts"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, input *humaapi.GetDmPostInput) (*humaapi.DmPostOutput, error) {
		// 公開レベルのチェック（publicエンドポイント）
		if err := auth.CheckAccessLevel(ctx, auth.AccessLevelPublic); err != nil {
			return nil, huma.Error403Forbidden(err.Error())
		}

		// string → int64 変換
		id, err := strconv.ParseInt(input.ID, 10, 64)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid id format")
		}
		userID, err := strconv.ParseInt(input.UserID, 10, 64)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid user_id format")
		}

		dmPost, err := h.dmPostService.GetDmPost(ctx, id, userID)
		if err != nil {
			return nil, huma.Error404NotFound(err.Error())
		}

		resp := &humaapi.DmPostOutput{}
		resp.Body = *dmPost
		return resp, nil
	})

	// GET /api/dm-posts - 投稿一覧取得
	huma.Register(api, huma.Operation{
		OperationID: "list-posts",
		Method:      http.MethodGet,
		Path:        "/api/dm-posts",
		Summary:     "投稿一覧を取得",
		Description: "**Access Level:** `public` (Public API Key JWT または Auth0 JWT でアクセス可能)",
		Tags:        []string{"posts"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, input *humaapi.ListDmPostsInput) (*humaapi.DmPostsOutput, error) {
		// 公開レベルのチェック（publicエンドポイント）
		if err := auth.CheckAccessLevel(ctx, auth.AccessLevelPublic); err != nil {
			return nil, huma.Error403Forbidden(err.Error())
		}

		var dmPosts []*model.DmPost
		var err error

		if input.UserID != "" {
			// string → int64 変換
			userID, parseErr := strconv.ParseInt(input.UserID, 10, 64)
			if parseErr != nil {
				return nil, huma.Error400BadRequest("invalid user_id format")
			}
			dmPosts, err = h.dmPostService.ListDmPostsByUser(ctx, userID, input.Limit, input.Offset)
		} else {
			dmPosts, err = h.dmPostService.ListDmPosts(ctx, input.Limit, input.Offset)
		}

		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		resp := &humaapi.DmPostsOutput{}
		resp.Body = dmPosts
		return resp, nil
	})

	// PUT /api/dm-posts/{id} - 投稿更新
	huma.Register(api, huma.Operation{
		OperationID: "update-post",
		Method:      http.MethodPut,
		Path:        "/api/dm-posts/{id}",
		Summary:     "投稿を更新",
		Description: "**Access Level:** `public` (Public API Key JWT または Auth0 JWT でアクセス可能)",
		Tags:        []string{"posts"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, input *humaapi.UpdateDmPostInput) (*humaapi.DmPostOutput, error) {
		// 公開レベルのチェック（publicエンドポイント）
		if err := auth.CheckAccessLevel(ctx, auth.AccessLevelPublic); err != nil {
			return nil, huma.Error403Forbidden(err.Error())
		}

		// string → int64 変換
		id, err := strconv.ParseInt(input.ID, 10, 64)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid id format")
		}
		userID, err := strconv.ParseInt(input.UserID, 10, 64)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid user_id format")
		}

		req := &model.UpdateDmPostRequest{
			Title:   input.Body.Title,
			Content: input.Body.Content,
		}

		dmPost, err := h.dmPostService.UpdateDmPost(ctx, id, userID, req)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		resp := &humaapi.DmPostOutput{}
		resp.Body = *dmPost
		return resp, nil
	})

	// DELETE /api/dm-posts/{id} - 投稿削除
	huma.Register(api, huma.Operation{
		OperationID:   "delete-post",
		Method:        http.MethodDelete,
		Path:          "/api/dm-posts/{id}",
		Summary:       "投稿を削除",
		Description:   "**Access Level:** `public` (Public API Key JWT または Auth0 JWT でアクセス可能)",
		Tags:          []string{"posts"},
		DefaultStatus: http.StatusNoContent,
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, input *humaapi.DeleteDmPostInput) (*struct{}, error) {
		// 公開レベルのチェック（publicエンドポイント）
		if err := auth.CheckAccessLevel(ctx, auth.AccessLevelPublic); err != nil {
			return nil, huma.Error403Forbidden(err.Error())
		}

		// string → int64 変換
		id, err := strconv.ParseInt(input.ID, 10, 64)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid id format")
		}
		userID, err := strconv.ParseInt(input.UserID, 10, 64)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid user_id format")
		}

		err = h.dmPostService.DeleteDmPost(ctx, id, userID)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		return nil, nil
	})

	// GET /api/dm-user-posts - ユーザーと投稿のJOIN結果取得
	huma.Register(api, huma.Operation{
		OperationID: "get-user-posts",
		Method:      http.MethodGet,
		Path:        "/api/dm-user-posts",
		Summary:     "ユーザーと投稿のJOIN結果を取得",
		Description: "**Access Level:** `public` (Public API Key JWT または Auth0 JWT でアクセス可能)",
		Tags:        []string{"posts"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, input *humaapi.GetDmUserPostsInput) (*humaapi.DmUserPostsOutput, error) {
		// 公開レベルのチェック（publicエンドポイント）
		if err := auth.CheckAccessLevel(ctx, auth.AccessLevelPublic); err != nil {
			return nil, huma.Error403Forbidden(err.Error())
		}

		dmUserPosts, err := h.dmPostService.GetDmUserPosts(ctx, input.Limit, input.Offset)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		resp := &humaapi.DmUserPostsOutput{}
		resp.Body = dmUserPosts
		return resp, nil
	})
}
