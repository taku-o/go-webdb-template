---
name: api-endpoint-creator
description: 新しいAPIエンドポイントを追加する際に使用。Huma API、Echo、Handler → Usecase → Service → Repository の4層構成、エンドポイント登録パターンを適用。REST API、新規エンドポイント、CRUD APIを追加する場合に使用。
allowed-tools: Read, Glob
---

# API エンドポイント作成パターン

このプロジェクトのAPI実装パターンを定義します。

## アーキテクチャ

4層構成。Handler は Usecase を保持し、Service を直接持たない。Service は Handler から直接呼ばず、Usecase 経由で呼ぶ。

```
Handler (API層) → Usecase (アプリケーションロジック) → Service (ビジネスロジック) → Repository (データアクセス)
```

## 使用フレームワーク

- **Echo**: HTTPルーター
- **Huma v2**: OpenAPI対応のAPIフレームワーク
- **humaecho**: EchoとHumaのアダプター

## ディレクトリ構成

```
server/internal/
├── api/
│   ├── handler/
│   │   ├── dm_user_handler.go
│   │   └── dm_post_handler.go
│   ├── huma/
│   │   ├── inputs.go        # Huma用の入力型定義
│   │   └── outputs.go       # Huma用の出力型定義
│   └── router/
│       └── router.go        # ルーター設定
├── usecase/
│   └── api/
│       ├── dm_user_usecase.go
│       └── dm_post_usecase.go
├── service/
│   ├── dm_user_service.go
│   └── dm_post_service.go
└── repository/
    ├── dm_user_repository.go
    └── dm_post_repository.go
```

## 参照ファイル

Handler:
- `server/internal/api/handler/dm_user_handler.go`
- `server/internal/api/handler/dm_post_handler.go`

Huma型定義:
- `server/internal/api/huma/inputs.go`
- `server/internal/api/huma/outputs.go`

Usecase:
- `server/internal/usecase/api/dm_user_usecase.go`
- `server/internal/usecase/api/dm_post_usecase.go`

ルーター:
- `server/internal/api/router/router.go`

Service:
- `server/internal/service/dm_user_service.go`（Handler から直接呼ばない。Usecase 経由で呼ぶ）

## Handler パターン

### 構造体定義

Handler は Usecase を保持する。Service を直接持たない。

```go
package handler

import (
    "context"
    "net/http"

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
```

### エンドポイント登録

Handler は Usecase を呼ぶ。Service は直接呼ばない。

```go
// RegisterDmUserEndpoints はHuma APIにユーザーエンドポイントを登録
func RegisterDmUserEndpoints(api huma.API, h *DmUserHandler) {
    // POST /api/dm-users - ユーザー作成
    huma.Register(api, huma.Operation{
        OperationID:   "create-user",
        Method:        http.MethodPost,
        Path:          "/api/dm-users",
        Summary:       "ユーザーを作成",
        Tags:          []string{"users"},
        DefaultStatus: http.StatusCreated,
        Security: []map[string][]string{
            {"bearerAuth": {}},
        },
    }, func(ctx context.Context, input *humaapi.CreateDmUserInput) (*humaapi.DmUserOutput, error) {
        // 認証・アクセスレベルチェックで拒否された場合は 403
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

    // GET /api/dm-users/{id} - ユーザー取得
    huma.Register(api, huma.Operation{
        OperationID: "get-user",
        Method:      http.MethodGet,
        Path:        "/api/dm-users/{id}",
        Summary:     "ユーザーを取得",
        Tags:        []string{"users"},
    }, func(ctx context.Context, input *humaapi.GetDmUserInput) (*humaapi.DmUserOutput, error) {
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
        Tags:        []string{"users"},
    }, func(ctx context.Context, input *humaapi.ListDmUsersInput) (*humaapi.DmUsersOutput, error) {
        users, err := h.dmUserUsecase.ListDmUsers(ctx, input.Limit, input.Offset)
        if err != nil {
            return nil, huma.Error500InternalServerError(err.Error())
        }

        resp := &humaapi.DmUsersOutput{}
        resp.Body = users
        return resp, nil
    })

    // PUT /api/dm-users/{id} - ユーザー更新
    huma.Register(api, huma.Operation{
        OperationID: "update-user",
        Method:      http.MethodPut,
        Path:        "/api/dm-users/{id}",
        Summary:     "ユーザーを更新",
        Tags:        []string{"users"},
    }, func(ctx context.Context, input *humaapi.UpdateDmUserInput) (*humaapi.DmUserOutput, error) {
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
        Tags:          []string{"users"},
        DefaultStatus: http.StatusNoContent,
    }, func(ctx context.Context, input *humaapi.DeleteDmUserInput) (*struct{}, error) {
        err := h.dmUserUsecase.DeleteDmUser(ctx, input.ID)
        if err != nil {
            return nil, huma.Error500InternalServerError(err.Error())
        }

        return nil, nil
    })
}
```

## Huma 入出力型パターン

入力型は `server/internal/api/huma/inputs.go`、出力型は `server/internal/api/huma/outputs.go` に定義する。

```go
// server/internal/api/huma/inputs.go

// 入力型
type CreateDmUserInput struct {
    Body struct {
        Name  string `json:"name" required:"true" maxLength:"100"`
        Email string `json:"email" required:"true" format:"email" maxLength:"255"`
    }
}

type GetDmUserInput struct {
    ID string `path:"id" doc:"ユーザーID（文字列形式）"`
}

type ListDmUsersInput struct {
    Limit  int `query:"limit" default:"20" minimum:"1" maximum:"100"`
    Offset int `query:"offset" default:"0" minimum:"0"`
}

type UpdateDmUserInput struct {
    ID   string `path:"id"`
    Body struct {
        Name  string `json:"name,omitempty" maxLength:"100"`
        Email string `json:"email,omitempty" format:"email" maxLength:"255"`
    }
}

type DeleteDmUserInput struct {
    ID string `path:"id"`
}
```

```go
// server/internal/api/huma/outputs.go

import "github.com/taku-o/go-webdb-template/internal/model"

// 出力型
type DmUserOutput struct {
    Body model.DmUser
}

type DmUsersOutput struct {
    Body []*model.DmUser
}
```

## Service パターン

Service は Handler から直接呼ばない。Usecase 経由で呼ぶ。Usecase が Service を保持し、Handler は Usecase のみを保持する。

```go
package service

import (
    "context"

    "github.com/taku-o/go-webdb-template/internal/model"
    "github.com/taku-o/go-webdb-template/internal/repository"
)

type DmUserService struct {
    dmUserRepo *repository.DmUserRepository
}

func NewDmUserService(dmUserRepo *repository.DmUserRepository) *DmUserService {
    return &DmUserService{
        dmUserRepo: dmUserRepo,
    }
}

func (s *DmUserService) CreateDmUser(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error) {
    return s.dmUserRepo.Create(ctx, req)
}

func (s *DmUserService) GetDmUser(ctx context.Context, id string) (*model.DmUser, error) {
    return s.dmUserRepo.GetByID(ctx, id)
}

// ListDmUsers, UpdateDmUser, DeleteDmUser も同様...
```

## ルーターへの登録

`router.go` でエンドポイントを登録:

```go
// Humaエンドポイントの登録
handler.RegisterDmUserEndpoints(humaAPI, dmUserHandler)
handler.RegisterDmPostEndpoints(humaAPI, dmPostHandler)
```

## エラーレスポンス

```go
// 403 Forbidden（認証・アクセスレベル拒否）
return nil, huma.Error403Forbidden(err.Error())

// 404 Not Found
return nil, huma.Error404NotFound(err.Error())

// 500 Internal Server Error
return nil, huma.Error500InternalServerError(err.Error())

// 400 Bad Request
return nil, huma.Error400BadRequest("Invalid input")
```
