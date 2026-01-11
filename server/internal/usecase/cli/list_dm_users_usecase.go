package cli

import (
	"context"

	"github.com/taku-o/go-webdb-template/internal/model"
	"github.com/taku-o/go-webdb-template/internal/usecase"
)

// ListDmUsersUsecase はCLI用のdm_user一覧取得usecase
type ListDmUsersUsecase struct {
	dmUserService usecase.DmUserServiceInterface
}

// NewListDmUsersUsecase は新しいListDmUsersUsecaseを作成
func NewListDmUsersUsecase(dmUserService usecase.DmUserServiceInterface) *ListDmUsersUsecase {
	return &ListDmUsersUsecase{
		dmUserService: dmUserService,
	}
}

// ListDmUsers はユーザー一覧を取得
func (u *ListDmUsersUsecase) ListDmUsers(ctx context.Context, limit, offset int) ([]*model.DmUser, error) {
	return u.dmUserService.ListDmUsers(ctx, limit, offset)
}
