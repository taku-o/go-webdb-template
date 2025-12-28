package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/example/go-webdb-template/internal/config"
)

// JWTType はJWTの種類
type JWTType string

const (
	JWTTypeAuth0        JWTType = "auth0"
	JWTTypePublicAPIKey JWTType = "public_api_key"
	JWTTypeUnknown      JWTType = "unknown"
)

// JWTClaims はJWTのクレーム構造
type JWTClaims struct {
	Issuer   string   `json:"iss"`
	Subject  string   `json:"sub"`
	Type     string   `json:"type"`
	Scope    []string `json:"scope"`
	IssuedAt int64    `json:"iat"`
	Version  string   `json:"version"`
	Env      string   `json:"env"`
	jwt.RegisteredClaims
}

// GetExpirationTime implements jwt.Claims interface
func (c JWTClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return nil, nil // Public APIキーは無期限
}

// GetIssuedAt implements jwt.Claims interface
func (c JWTClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	if c.IssuedAt == 0 {
		return nil, nil
	}
	return jwt.NewNumericDate(time.Unix(c.IssuedAt, 0)), nil
}

// GetNotBefore implements jwt.Claims interface
func (c JWTClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

// GetIssuer implements jwt.Claims interface
func (c JWTClaims) GetIssuer() (string, error) {
	return c.Issuer, nil
}

// GetSubject implements jwt.Claims interface
func (c JWTClaims) GetSubject() (string, error) {
	return c.Subject, nil
}

// GetAudience implements jwt.Claims interface
func (c JWTClaims) GetAudience() (jwt.ClaimStrings, error) {
	return nil, nil
}

// JWTValidator はJWT検証機能を提供
type JWTValidator struct {
	secretKey       string
	invalidVersions []string
	currentEnv      string
}

// NewJWTValidator は新しいJWTValidatorを作成
func NewJWTValidator(cfg *config.APIConfig, env string) *JWTValidator {
	return &JWTValidator{
		secretKey:       cfg.SecretKey,
		invalidVersions: cfg.InvalidVersions,
		currentEnv:      env,
	}
}

// ValidateJWT はJWTトークンを検証
func (v *JWTValidator) ValidateJWT(tokenString string) (*JWTClaims, error) {
	// JWTトークンをパース
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 署名アルゴリズムの検証
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(v.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT: %w", err)
	}

	// クレームの取得
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// クレームの検証
	if err := v.validateClaims(claims); err != nil {
		return nil, err
	}

	return claims, nil
}

// validateClaims はクレームを検証
func (v *JWTValidator) validateClaims(claims *JWTClaims) error {
	// issの検証
	if claims.Issuer != "go-webdb-template" {
		return errors.New("invalid issuer")
	}

	// typeの検証
	if claims.Type != "public" && claims.Type != "private" {
		return errors.New("invalid token type")
	}

	// versionの検証（無効バージョンリストとの照合）
	if v.IsVersionInvalid(claims.Version) {
		return errors.New("invalid token version")
	}

	// envの検証
	if claims.Env != v.currentEnv {
		return errors.New("token environment mismatch")
	}

	return nil
}

// IsVersionInvalid はバージョンが無効かどうかを判定
func (v *JWTValidator) IsVersionInvalid(version string) bool {
	for _, invalidVersion := range v.invalidVersions {
		if version == invalidVersion {
			return true
		}
	}
	return false
}

// ParseJWTClaims はJWTトークンからクレームをパース（表示用、署名検証なし）
func ParseJWTClaims(tokenString string) (*JWTClaims, error) {
	parser := jwt.NewParser()
	token, _, err := parser.ParseUnverified(tokenString, &JWTClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// GeneratePublicAPIKey はPublic JWTキーを生成
func GeneratePublicAPIKey(secretKey string, currentVersion string, env string, issuedAt int64) (string, error) {
	claims := &JWTClaims{
		Issuer:   "go-webdb-template",
		Subject:  "public_client",
		Type:     "public",
		Scope:    []string{"read", "write"},
		IssuedAt: issuedAt,
		Version:  currentVersion,
		Env:      env,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// DetectJWTType はJWTの種類を判別（署名検証前）
func DetectJWTType(tokenString string) (JWTType, error) {
	// 署名検証なしでパース
	parser := jwt.NewParser()
	token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return JWTTypeUnknown, fmt.Errorf("failed to parse JWT: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return JWTTypeUnknown, errors.New("invalid token claims")
	}

	// issによる判別
	issuer, ok := claims["iss"].(string)
	if !ok {
		return JWTTypeUnknown, errors.New("missing issuer claim")
	}

	if issuer == "go-webdb-template" {
		return JWTTypePublicAPIKey, nil
	}

	// Auth0のドメインパターンをチェック
	if strings.HasPrefix(issuer, "https://") &&
		(strings.Contains(issuer, ".auth0.com") ||
			strings.Contains(issuer, ".auth0.jp")) {
		return JWTTypeAuth0, nil
	}

	return JWTTypeUnknown, fmt.Errorf("unknown issuer: %s", issuer)
}
