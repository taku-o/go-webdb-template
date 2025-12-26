package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/danielgtaylor/huma/v2"
	humaapi "github.com/example/go-webdb-template/internal/api/huma"
	"github.com/example/go-webdb-template/internal/model"
	"github.com/example/go-webdb-template/internal/service"
	"github.com/gorilla/mux"
	"github.com/labstack/echo/v4"
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

// CreatePost は投稿を作成
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req model.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	post, err := h.postService.CreatePost(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

// GetPost は投稿を取得
func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(r.URL.Query().Get("user_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	post, err := h.postService.GetPost(r.Context(), id, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// ListPosts は投稿一覧を取得
func (h *PostHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
	limit := 20
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			limit = val
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil {
			offset = val
		}
	}

	// user_idが指定されている場合はユーザーの投稿のみ取得
	if userIDStr := r.URL.Query().Get("user_id"); userIDStr != "" {
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		posts, err := h.postService.ListPostsByUser(r.Context(), userID, limit, offset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
		return
	}

	posts, err := h.postService.ListPosts(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// GetUserPosts はユーザーと投稿をJOINして取得
func (h *PostHandler) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	limit := 20
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			limit = val
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil {
			offset = val
		}
	}

	userPosts, err := h.postService.GetUserPosts(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userPosts)
}

// UpdatePost は投稿を更新
func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(r.URL.Query().Get("user_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req model.UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	post, err := h.postService.UpdatePost(r.Context(), id, userID, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// DeletePost は投稿を削除
func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(r.URL.Query().Get("user_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.postService.DeletePost(r.Context(), id, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CreatePostEcho は投稿を作成（Echo形式）
func (h *PostHandler) CreatePostEcho(c echo.Context) error {
	var req model.CreatePostRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	post, err := h.postService.CreatePost(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, post)
}

// GetPostEcho は投稿を取得（Echo形式）
func (h *PostHandler) GetPostEcho(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid post ID"})
	}

	userID, err := strconv.ParseInt(c.QueryParam("user_id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	post, err := h.postService.GetPost(c.Request().Context(), id, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, post)
}

// ListPostsEcho は投稿一覧を取得（Echo形式）
func (h *PostHandler) ListPostsEcho(c echo.Context) error {
	limit := 20
	offset := 0

	if l := c.QueryParam("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			limit = val
		}
	}
	if o := c.QueryParam("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil {
			offset = val
		}
	}

	// user_idが指定されている場合はユーザーの投稿のみ取得
	if userIDStr := c.QueryParam("user_id"); userIDStr != "" {
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		}

		posts, err := h.postService.ListPostsByUser(c.Request().Context(), userID, limit, offset)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, posts)
	}

	posts, err := h.postService.ListPosts(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, posts)
}

// GetUserPostsEcho はユーザーと投稿をJOINして取得（Echo形式）
func (h *PostHandler) GetUserPostsEcho(c echo.Context) error {
	limit := 20
	offset := 0

	if l := c.QueryParam("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			limit = val
		}
	}
	if o := c.QueryParam("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil {
			offset = val
		}
	}

	userPosts, err := h.postService.GetUserPosts(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, userPosts)
}

// UpdatePostEcho は投稿を更新（Echo形式）
func (h *PostHandler) UpdatePostEcho(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid post ID"})
	}

	userID, err := strconv.ParseInt(c.QueryParam("user_id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	var req model.UpdatePostRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	post, err := h.postService.UpdatePost(c.Request().Context(), id, userID, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, post)
}

// DeletePostEcho は投稿を削除（Echo形式）
func (h *PostHandler) DeletePostEcho(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid post ID"})
	}

	userID, err := strconv.ParseInt(c.QueryParam("user_id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	if err := h.postService.DeletePost(c.Request().Context(), id, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
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
		userID, err := strconv.ParseInt(input.Body.UserID, 10, 64)
		if err != nil {
			return nil, huma.Error400BadRequest("Invalid user_id")
		}

		req := &model.CreatePostRequest{
			UserID:  userID,
			Title:   input.Body.Title,
			Content: input.Body.Content,
		}

		post, err := h.postService.CreatePost(ctx, req)
		if err != nil {
			return nil, huma.Error400BadRequest(err.Error())
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
