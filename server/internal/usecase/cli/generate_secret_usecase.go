package cli

import (
	"context"

	"github.com/taku-o/go-webdb-template/internal/service"
)

// GenerateSecretUsecase はCLI用の秘密鍵生成usecase
type GenerateSecretUsecase struct {
	secretService service.SecretServiceInterface
}

// NewGenerateSecretUsecase は新しいGenerateSecretUsecaseを作成
func NewGenerateSecretUsecase(secretService service.SecretServiceInterface) *GenerateSecretUsecase {
	return &GenerateSecretUsecase{
		secretService: secretService,
	}
}

// GenerateSecret は秘密鍵を生成
func (u *GenerateSecretUsecase) GenerateSecret(ctx context.Context) (string, error) {
	return u.secretService.GenerateSecretKey(ctx)
}
