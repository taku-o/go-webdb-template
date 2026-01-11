package admin

import (
	"context"
	"time"

	"github.com/taku-o/go-webdb-template/internal/auth"
	"github.com/taku-o/go-webdb-template/internal/config"
	"github.com/taku-o/go-webdb-template/internal/service"
)

// APIKeyUsecase はAPIキー発行のビジネスロジックを担当
type APIKeyUsecase struct {
	apiKeyService service.APIKeyServiceInterface
	cfg           *config.Config
}

// NewAPIKeyUsecase は新しいAPIKeyUsecaseを作成
func NewAPIKeyUsecase(apiKeyService service.APIKeyServiceInterface, cfg *config.Config) *APIKeyUsecase {
	return &APIKeyUsecase{
		apiKeyService: apiKeyService,
		cfg:           cfg,
	}
}

// GenerateAPIKey はAPIキーを生成
func (u *APIKeyUsecase) GenerateAPIKey(ctx context.Context, env string) (string, error) {
	if env == "" {
		env = "develop"
	}

	now := time.Now()
	token, err := u.apiKeyService.GenerateAPIKey(ctx, u.cfg.API.SecretKey, u.cfg.API.CurrentVersion, env, now.Unix())
	if err != nil {
		return "", err
	}

	return token, nil
}

// DecodeAPIKeyPayload はAPIキーのペイロードをデコード
func (u *APIKeyUsecase) DecodeAPIKeyPayload(ctx context.Context, token string) (*auth.JWTClaims, error) {
	claims, err := u.apiKeyService.DecodeAPIKeyPayload(ctx, token)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
