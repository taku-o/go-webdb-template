package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Auth0Validator はAuth0 JWT検証機能を提供
type Auth0Validator struct {
	jwks *keyfunc.JWKS
}

// NewAuth0Validator は新しいAuth0Validatorを作成
func NewAuth0Validator(issuerBaseURL string) (*Auth0Validator, error) {
	// JWKS URLの構築
	jwksURL := fmt.Sprintf("%s/.well-known/jwks.json", issuerBaseURL)

	// keyfuncのオプション設定
	options := keyfunc.Options{
		RefreshInterval:   time.Hour * 12,   // 12時間ごとに定期更新
		RefreshRateLimit:  time.Minute * 5,  // 再取得は最低5分あける（DoS対策）
		RefreshTimeout:    time.Second * 10, // 取得時のタイムアウト
		RefreshUnknownKID: true,             // 未知のKIDが来たら再取得する（重要！）
	}

	// JWKSの取得とキャッシュ
	jwks, err := keyfunc.Get(jwksURL, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get JWKS: %w", err)
	}

	return &Auth0Validator{
		jwks: jwks,
	}, nil
}

// ValidateAuth0JWT はAuth0 JWTを検証
func (v *Auth0Validator) ValidateAuth0JWT(tokenString string) (*jwt.Token, error) {
	// JWTの検証
	token, err := jwt.Parse(tokenString, v.jwks.Keyfunc)
	if err != nil {
		return nil, fmt.Errorf("failed to validate Auth0 JWT: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid Auth0 JWT")
	}

	return token, nil
}

// Close はリソースを解放
func (v *Auth0Validator) Close() {
	if v.jwks != nil {
		v.jwks.EndBackground()
	}
}
