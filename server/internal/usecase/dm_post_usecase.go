package usecase

import (
	"context"

	"github.com/taku-o/go-webdb-template/internal/model"
)

// DmPostServiceInterface はDmPostServiceのインターフェース
type DmPostServiceInterface interface {
	CreateDmPost(ctx context.Context, req *model.CreateDmPostRequest) (*model.DmPost, error)
	GetDmPost(ctx context.Context, id string, userID string) (*model.DmPost, error)
	ListDmPosts(ctx context.Context, limit, offset int) ([]*model.DmPost, error)
	ListDmPostsByUser(ctx context.Context, userID string, limit, offset int) ([]*model.DmPost, error)
	GetDmUserPosts(ctx context.Context, limit, offset int) ([]*model.DmUserPost, error)
	UpdateDmPost(ctx context.Context, id string, userID string, req *model.UpdateDmPostRequest) (*model.DmPost, error)
	DeleteDmPost(ctx context.Context, id string, userID string) error
}

// DmPostUsecase は投稿のビジネスロジックを担当するユースケース層
type DmPostUsecase struct {
	dmPostService DmPostServiceInterface
}

// NewDmPostUsecase は新しいDmPostUsecaseを作成
func NewDmPostUsecase(dmPostService DmPostServiceInterface) *DmPostUsecase {
	return &DmPostUsecase{
		dmPostService: dmPostService,
	}
}

// CreateDmPost は投稿を作成
func (u *DmPostUsecase) CreateDmPost(ctx context.Context, req *model.CreateDmPostRequest) (*model.DmPost, error) {
	return u.dmPostService.CreateDmPost(ctx, req)
}

// GetDmPost はIDで投稿を取得
func (u *DmPostUsecase) GetDmPost(ctx context.Context, id string, userID string) (*model.DmPost, error) {
	return u.dmPostService.GetDmPost(ctx, id, userID)
}

// ListDmPosts は投稿一覧を取得
func (u *DmPostUsecase) ListDmPosts(ctx context.Context, limit, offset int) ([]*model.DmPost, error) {
	return u.dmPostService.ListDmPosts(ctx, limit, offset)
}

// ListDmPostsByUser はユーザーIDで投稿一覧を取得
func (u *DmPostUsecase) ListDmPostsByUser(ctx context.Context, userID string, limit, offset int) ([]*model.DmPost, error) {
	return u.dmPostService.ListDmPostsByUser(ctx, userID, limit, offset)
}

// GetDmUserPosts はユーザーと投稿をJOINして取得
func (u *DmPostUsecase) GetDmUserPosts(ctx context.Context, limit, offset int) ([]*model.DmUserPost, error) {
	return u.dmPostService.GetDmUserPosts(ctx, limit, offset)
}

// UpdateDmPost は投稿を更新
func (u *DmPostUsecase) UpdateDmPost(ctx context.Context, id string, userID string, req *model.UpdateDmPostRequest) (*model.DmPost, error) {
	return u.dmPostService.UpdateDmPost(ctx, id, userID, req)
}

// DeleteDmPost は投稿を削除
func (u *DmPostUsecase) DeleteDmPost(ctx context.Context, id string, userID string) error {
	return u.dmPostService.DeleteDmPost(ctx, id, userID)
}
