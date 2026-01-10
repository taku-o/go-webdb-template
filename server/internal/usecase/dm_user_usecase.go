package usecase

import (
	"context"

	"github.com/taku-o/go-webdb-template/internal/model"
)

// DmUserServiceInterface はDmUserServiceのインターフェース
type DmUserServiceInterface interface {
	CreateDmUser(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error)
	GetDmUser(ctx context.Context, id string) (*model.DmUser, error)
	ListDmUsers(ctx context.Context, limit, offset int) ([]*model.DmUser, error)
	UpdateDmUser(ctx context.Context, id string, req *model.UpdateDmUserRequest) (*model.DmUser, error)
	DeleteDmUser(ctx context.Context, id string) error
}

// DmUserUsecase はdm_user関連のビジネスロジックを担当
type DmUserUsecase struct {
	dmUserService DmUserServiceInterface
}

// NewDmUserUsecase は新しいDmUserUsecaseを作成
func NewDmUserUsecase(dmUserService DmUserServiceInterface) *DmUserUsecase {
	return &DmUserUsecase{
		dmUserService: dmUserService,
	}
}

// CreateDmUser はユーザーを作成
func (u *DmUserUsecase) CreateDmUser(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error) {
	return u.dmUserService.CreateDmUser(ctx, req)
}

// GetDmUser はIDでユーザーを取得
func (u *DmUserUsecase) GetDmUser(ctx context.Context, id string) (*model.DmUser, error) {
	return u.dmUserService.GetDmUser(ctx, id)
}

// ListDmUsers はユーザー一覧を取得
func (u *DmUserUsecase) ListDmUsers(ctx context.Context, limit, offset int) ([]*model.DmUser, error) {
	return u.dmUserService.ListDmUsers(ctx, limit, offset)
}

// UpdateDmUser はユーザーを更新
func (u *DmUserUsecase) UpdateDmUser(ctx context.Context, id string, req *model.UpdateDmUserRequest) (*model.DmUser, error) {
	return u.dmUserService.UpdateDmUser(ctx, id, req)
}

// DeleteDmUser はユーザーを削除
func (u *DmUserUsecase) DeleteDmUser(ctx context.Context, id string) error {
	return u.dmUserService.DeleteDmUser(ctx, id)
}
