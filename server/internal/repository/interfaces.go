package repository

import (
	"context"

	"github.com/taku-o/go-webdb-template/internal/model"
)

// DmUserRepositoryInterface はDmUserRepositoryの共通インターフェース
type DmUserRepositoryInterface interface {
	Create(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error)
	GetByID(ctx context.Context, id int64) (*model.DmUser, error)
	List(ctx context.Context, limit, offset int) ([]*model.DmUser, error)
	Update(ctx context.Context, id int64, req *model.UpdateDmUserRequest) (*model.DmUser, error)
	Delete(ctx context.Context, id int64) error
}

// DmPostRepositoryInterface はDmPostRepositoryの共通インターフェース
type DmPostRepositoryInterface interface {
	Create(ctx context.Context, req *model.CreateDmPostRequest) (*model.DmPost, error)
	GetByID(ctx context.Context, id int64, userID int64) (*model.DmPost, error)
	ListByUserID(ctx context.Context, userID int64, limit, offset int) ([]*model.DmPost, error)
	List(ctx context.Context, limit, offset int) ([]*model.DmPost, error)
	GetUserPosts(ctx context.Context, limit, offset int) ([]*model.DmUserPost, error)
	Update(ctx context.Context, id int64, userID int64, req *model.UpdateDmPostRequest) (*model.DmPost, error)
	Delete(ctx context.Context, id int64, userID int64) error
}
