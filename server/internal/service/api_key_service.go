package service

import (
	"github.com/taku-o/go-webdb-template/internal/auth"
)

// APIKeyServiceInterface はAPIキーサービスのインターフェース
type APIKeyServiceInterface interface {
	GenerateAPIKey(secretKey, version, env string, issuedAt int64) (string, error)
	DecodeAPIKeyPayload(token string) (*auth.JWTClaims, error)
}

// APIKeyService はAPIキー発行のドメインロジックを担当
type APIKeyService struct{}

// NewAPIKeyService は新しいAPIKeyServiceを作成
func NewAPIKeyService() *APIKeyService {
	return &APIKeyService{}
}

// GenerateAPIKey はAPIキーを生成
func (s *APIKeyService) GenerateAPIKey(secretKey, version, env string, issuedAt int64) (string, error) {
	token, err := auth.GeneratePublicAPIKey(secretKey, version, env, issuedAt)
	if err != nil {
		return "", err
	}

	return token, nil
}

// DecodeAPIKeyPayload はAPIキーのペイロードをデコード
func (s *APIKeyService) DecodeAPIKeyPayload(token string) (*auth.JWTClaims, error) {
	claims, err := auth.ParseJWTClaims(token)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
