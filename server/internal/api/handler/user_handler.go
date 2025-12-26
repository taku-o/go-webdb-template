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

// CreateUser はユーザーを作成
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userService.CreateUser(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUser はユーザーを取得
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUser(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// ListUsers はユーザー一覧を取得
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
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

	users, err := h.userService.ListUsers(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// UpdateUser はユーザーを更新
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req model.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userService.UpdateUser(r.Context(), id, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// DeleteUser はユーザーを削除
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.userService.DeleteUser(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CreateUserEcho はユーザーを作成（Echo形式）
func (h *UserHandler) CreateUserEcho(c echo.Context) error {
	var req model.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	user, err := h.userService.CreateUser(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, user)
}

// GetUserEcho はユーザーを取得（Echo形式）
func (h *UserHandler) GetUserEcho(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	user, err := h.userService.GetUser(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

// ListUsersEcho はユーザー一覧を取得（Echo形式）
func (h *UserHandler) ListUsersEcho(c echo.Context) error {
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

	users, err := h.userService.ListUsers(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, users)
}

// UpdateUserEcho はユーザーを更新（Echo形式）
func (h *UserHandler) UpdateUserEcho(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	var req model.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	user, err := h.userService.UpdateUser(c.Request().Context(), id, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

// DeleteUserEcho はユーザーを削除（Echo形式）
func (h *UserHandler) DeleteUserEcho(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	if err := h.userService.DeleteUser(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// RegisterUserEndpoints はHuma APIにユーザーエンドポイントを登録
func RegisterUserEndpoints(api huma.API, h *UserHandler) {
	// POST /api/users - ユーザー作成
	huma.Register(api, huma.Operation{
		OperationID:   "create-user",
		Method:        http.MethodPost,
		Path:          "/api/users",
		Summary:       "ユーザーを作成",
		Tags:          []string{"users"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, input *humaapi.CreateUserInput) (*humaapi.UserOutput, error) {
		req := &model.CreateUserRequest{
			Name:  input.Body.Name,
			Email: input.Body.Email,
		}

		user, err := h.userService.CreateUser(ctx, req)
		if err != nil {
			return nil, huma.Error400BadRequest(err.Error())
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
		Tags:        []string{"users"},
	}, func(ctx context.Context, input *humaapi.GetUserInput) (*humaapi.UserOutput, error) {
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
		Tags:        []string{"users"},
	}, func(ctx context.Context, input *humaapi.ListUsersInput) (*humaapi.UsersOutput, error) {
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
		Tags:        []string{"users"},
	}, func(ctx context.Context, input *humaapi.UpdateUserInput) (*humaapi.UserOutput, error) {
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
		Tags:          []string{"users"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, input *humaapi.DeleteUserInput) (*struct{}, error) {
		err := h.userService.DeleteUser(ctx, input.ID)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		return nil, nil
	})
}
