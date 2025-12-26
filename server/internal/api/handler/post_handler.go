package handler

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	humaapi "github.com/example/go-webdb-template/internal/api/huma"
	"github.com/example/go-webdb-template/internal/model"
	"github.com/example/go-webdb-template/internal/service"
)

// PostHandler は投稿APIのハンドラー
type PostHandler struct {
	postService *service.PostService
}

// NewPostHandler は新しいPostHandlerを作成
func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// RegisterPostEndpoints はHuma APIに投稿エンドポイントを登録
func RegisterPostEndpoints(api huma.API, h *PostHandler) {
	// POST /api/posts - 投稿作成
	huma.Register(api, huma.Operation{
		OperationID:   "create-post",
		Method:        http.MethodPost,
		Path:          "/api/posts",
		Summary:       "投稿を作成",
		Tags:          []string{"posts"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, input *humaapi.CreatePostInput) (*humaapi.PostOutput, error) {
		req := &model.CreatePostRequest{
			UserID:  input.Body.UserID,
			Title:   input.Body.Title,
			Content: input.Body.Content,
		}

		post, err := h.postService.CreatePost(ctx, req)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		resp := &humaapi.PostOutput{}
		resp.Body = *post
		return resp, nil
	})

	// GET /api/posts/{id} - 投稿取得
	huma.Register(api, huma.Operation{
		OperationID: "get-post",
		Method:      http.MethodGet,
		Path:        "/api/posts/{id}",
		Summary:     "投稿を取得",
		Tags:        []string{"posts"},
	}, func(ctx context.Context, input *humaapi.GetPostInput) (*humaapi.PostOutput, error) {
		post, err := h.postService.GetPost(ctx, input.ID, input.UserID)
		if err != nil {
			return nil, huma.Error404NotFound(err.Error())
		}

		resp := &humaapi.PostOutput{}
		resp.Body = *post
		return resp, nil
	})

	// GET /api/posts - 投稿一覧取得
	huma.Register(api, huma.Operation{
		OperationID: "list-posts",
		Method:      http.MethodGet,
		Path:        "/api/posts",
		Summary:     "投稿一覧を取得",
		Tags:        []string{"posts"},
	}, func(ctx context.Context, input *humaapi.ListPostsInput) (*humaapi.PostsOutput, error) {
		var posts []*model.Post
		var err error

		if input.UserID > 0 {
			posts, err = h.postService.ListPostsByUser(ctx, input.UserID, input.Limit, input.Offset)
		} else {
			posts, err = h.postService.ListPosts(ctx, input.Limit, input.Offset)
		}

		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		resp := &humaapi.PostsOutput{}
		resp.Body = posts
		return resp, nil
	})

	// PUT /api/posts/{id} - 投稿更新
	huma.Register(api, huma.Operation{
		OperationID: "update-post",
		Method:      http.MethodPut,
		Path:        "/api/posts/{id}",
		Summary:     "投稿を更新",
		Tags:        []string{"posts"},
	}, func(ctx context.Context, input *humaapi.UpdatePostInput) (*humaapi.PostOutput, error) {
		req := &model.UpdatePostRequest{
			Title:   input.Body.Title,
			Content: input.Body.Content,
		}

		post, err := h.postService.UpdatePost(ctx, input.ID, input.UserID, req)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		resp := &humaapi.PostOutput{}
		resp.Body = *post
		return resp, nil
	})

	// DELETE /api/posts/{id} - 投稿削除
	huma.Register(api, huma.Operation{
		OperationID:   "delete-post",
		Method:        http.MethodDelete,
		Path:          "/api/posts/{id}",
		Summary:       "投稿を削除",
		Tags:          []string{"posts"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, input *humaapi.DeletePostInput) (*struct{}, error) {
		err := h.postService.DeletePost(ctx, input.ID, input.UserID)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		return nil, nil
	})

	// GET /api/user-posts - ユーザーと投稿のJOIN結果取得
	huma.Register(api, huma.Operation{
		OperationID: "get-user-posts",
		Method:      http.MethodGet,
		Path:        "/api/user-posts",
		Summary:     "ユーザーと投稿のJOIN結果を取得",
		Tags:        []string{"posts"},
	}, func(ctx context.Context, input *humaapi.GetUserPostsInput) (*humaapi.UserPostsOutput, error) {
		userPosts, err := h.postService.GetUserPosts(ctx, input.Limit, input.Offset)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		resp := &humaapi.UserPostsOutput{}
		resp.Body = userPosts
		return resp, nil
	})
}
