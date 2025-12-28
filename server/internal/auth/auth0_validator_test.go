package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAuth0Validator(t *testing.T) {
	// 実際のAuth0のJWKSエンドポイントに接続するテスト
	// 注意: このテストはネットワーク接続が必要
	t.Run("valid issuer base URL", func(t *testing.T) {
		validator, err := NewAuth0Validator("https://dev-oaa5vtzmld4dsxtd.jp.auth0.com")
		if err != nil {
			// ネットワークエラーの場合はスキップ
			t.Skipf("Skipping test due to network error: %v", err)
		}
		defer validator.Close()
		assert.NotNil(t, validator)
	})
}

func TestAuth0Validator_ValidateAuth0JWT_InvalidToken(t *testing.T) {
	// Auth0Validatorの初期化
	validator, err := NewAuth0Validator("https://dev-oaa5vtzmld4dsxtd.jp.auth0.com")
	if err != nil {
		t.Skipf("Skipping test due to network error: %v", err)
	}
	defer validator.Close()

	// 不正なトークンの検証
	_, err = validator.ValidateAuth0JWT("invalid-token")
	assert.Error(t, err)
}
