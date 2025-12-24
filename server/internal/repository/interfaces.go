package repository

import (
	"context"

	"github.com/example/go-webdb-template/internal/model"
)

// UserRepositoryInterface はUserRepositoryの共通インターフェース
type UserRepositoryInterface interface {
	Create(ctx context.Context, req *model.CreateUserRequest) (*model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	List(ctx context.Context, limit, offset int) ([]*model.User, error)
	Update(ctx context.Context, id int64, req *model.UpdateUserRequest) (*model.User, error)
	Delete(ctx context.Context, id int64) error
}

// PostRepositoryInterface はPostRepositoryの共通インターフェース
type PostRepositoryInterface interface {
	Create(ctx context.Context, req *model.CreatePostRequest) (*model.Post, error)
	GetByID(ctx context.Context, id int64, userID int64) (*model.Post, error)
	ListByUserID(ctx context.Context, userID int64, limit, offset int) ([]*model.Post, error)
	List(ctx context.Context, limit, offset int) ([]*model.Post, error)
	GetUserPosts(ctx context.Context, limit, offset int) ([]*model.UserPost, error)
	Update(ctx context.Context, id int64, userID int64, req *model.UpdatePostRequest) (*model.Post, error)
	Delete(ctx context.Context, id int64, userID int64) error
}
