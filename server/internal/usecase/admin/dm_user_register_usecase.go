package admin

import (
	"context"
	"fmt"

	"github.com/taku-o/go-webdb-template/internal/model"
	"github.com/taku-o/go-webdb-template/internal/usecase"
)

// DmUserRegisterUsecase はdm_user登録のビジネスロジックを担当
type DmUserRegisterUsecase struct {
	dmUserService usecase.DmUserServiceInterface
}

// NewDmUserRegisterUsecase は新しいDmUserRegisterUsecaseを作成
func NewDmUserRegisterUsecase(dmUserService usecase.DmUserServiceInterface) *DmUserRegisterUsecase {
	return &DmUserRegisterUsecase{
		dmUserService: dmUserService,
	}
}

// RegisterDmUser はユーザーを登録
func (u *DmUserRegisterUsecase) RegisterDmUser(ctx context.Context, name, email string) (string, error) {
	// メールアドレスの重複チェック
	exists, err := u.dmUserService.CheckEmailExists(ctx, email)
	if err != nil {
		return "", fmt.Errorf("failed to check email: %w", err)
	}
	if exists {
		return "", fmt.Errorf("このメールアドレスは既に登録されています")
	}

	req := &model.CreateDmUserRequest{
		Name:  name,
		Email: email,
	}

	dmUser, err := u.dmUserService.CreateDmUser(ctx, req)
	if err != nil {
		return "", err
	}

	return dmUser.ID, nil
}
