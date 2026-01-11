package service

import (
	"context"

	"github.com/taku-o/go-webdb-template/internal/auth"
)

// SecretServiceInterface は秘密鍵生成サービスのインターフェース
type SecretServiceInterface interface {
	GenerateSecretKey(ctx context.Context) (string, error)
}

// SecretService は秘密鍵生成のビジネスロジックを担当
type SecretService struct{}

// NewSecretService は新しいSecretServiceを作成
func NewSecretService() *SecretService {
	return &SecretService{}
}

// GenerateSecretKey は秘密鍵を生成
func (s *SecretService) GenerateSecretKey(ctx context.Context) (string, error) {
	return auth.GenerateSecretKey()
}
