package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/taku-o/go-webdb-template/internal/config"
)

// contextKey はコンテキストキーの型
type contextKey string

// AllowedAccessLevelKey はJWTの許容する公開レベルを格納するコンテキストキー
const AllowedAccessLevelKey contextKey = "allowed_access_level"

// AccessLevel はAPI公開レベル
type AccessLevel string

const (
	AccessLevelPublic  AccessLevel = "public"
	AccessLevelPrivate AccessLevel = "private"
)

// validateScope はスコープを検証
func validateScope(claims *JWTClaims, method string) error {
	hasRead := false
	hasWrite := false

	for _, scope := range claims.Scope {
		if scope == "read" {
			hasRead = true
		}
		if scope == "write" {
			hasWrite = true
		}
	}

	// GETリクエストにはreadスコープが必要
	if method == "GET" && !hasRead {
		return errors.New("read scope required")
	}

	// POST/PUT/DELETEリクエストにはwriteスコープが必要
	if (method == "POST" || method == "PUT" || method == "DELETE") && !hasWrite {
		return errors.New("write scope required")
	}

	return nil
}

// NewHumaAuthMiddleware は新しいHuma形式の認証ミドルウェアを作成
func NewHumaAuthMiddleware(cfg *config.APIConfig, env string, auth0IssuerBaseURL string) func(ctx huma.Context, next func(huma.Context)) {
	validator := NewJWTValidator(cfg, env)

	// Auth0Validatorの初期化
	var auth0Validator *Auth0Validator
	if auth0IssuerBaseURL != "" {
		var err error
		auth0Validator, err = NewAuth0Validator(auth0IssuerBaseURL)
		if err != nil {
			// エラーハンドリング（起動時エラーとして処理）
			panic("failed to create Auth0Validator: " + err.Error())
		}
	}

	return func(ctx huma.Context, next func(huma.Context)) {
		path := ctx.URL().Path

		// OpenAPIドキュメントのパスは認証をスキップ
		if isOpenAPIPath(path) {
			next(ctx)
			return
		}

		// /api/で始まるパスのみ認証を適用
		if !strings.HasPrefix(path, "/api/") {
			next(ctx)
			return
		}

		// AuthorizationヘッダーからJWTトークンを取得
		authHeader := ctx.Header("Authorization")
		if authHeader == "" {
			writeHumaError(ctx, http.StatusUnauthorized, "Authorization header is required")
			return
		}

		// Bearerトークンの抽出
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			writeHumaError(ctx, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		tokenString := parts[1]

		// JWT種類の判別
		jwtType, err := DetectJWTType(tokenString)
		if err != nil {
			writeHumaError(ctx, http.StatusUnauthorized, "Invalid token format")
			return
		}

		var claims *JWTClaims
		var allowedAccessLevel AccessLevel

		// JWT種類に応じた検証
		switch jwtType {
		case JWTTypeAuth0:
			if auth0Validator == nil {
				writeHumaError(ctx, http.StatusUnauthorized, "Auth0 JWT validation is not configured")
				return
			}
			_, err = auth0Validator.ValidateAuth0JWT(tokenString)
			if err != nil {
				writeHumaError(ctx, http.StatusUnauthorized, "Invalid Auth0 JWT")
				return
			}
			// Auth0 JWTはpublicとprivateの両方にアクセス可能
			allowedAccessLevel = AccessLevelPrivate

		case JWTTypePublicAPIKey:
			claims, err = validator.ValidateJWT(tokenString)
			if err != nil {
				writeHumaError(ctx, http.StatusUnauthorized, "Invalid API key")
				return
			}
			// Public API Key JWTはpublicなAPIのみアクセス可能
			allowedAccessLevel = AccessLevelPublic

		default:
			writeHumaError(ctx, http.StatusUnauthorized, "Unknown JWT type")
			return
		}

		// スコープ検証（Public API Key JWTの場合のみ）
		if jwtType == JWTTypePublicAPIKey && claims != nil {
			if err := validateScope(claims, ctx.Method()); err != nil {
				writeHumaError(ctx, http.StatusForbidden, "Insufficient scope")
				return
			}
		}

		// JWTの許容する公開レベルをコンテキストに設定
		newCtx := context.WithValue(ctx.Context(), AllowedAccessLevelKey, allowedAccessLevel)
		ctx = huma.WithContext(ctx, newCtx)

		// 次のハンドラーを実行
		next(ctx)
	}
}

// writeHumaError はHumaコンテキストにエラーレスポンスを書き込む
func writeHumaError(ctx huma.Context, statusCode int, message string) {
	ctx.SetStatus(statusCode)
	ctx.SetHeader("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    statusCode,
		"message": message,
	}
	json.NewEncoder(ctx.BodyWriter()).Encode(response)
}

// isOpenAPIPath はOpenAPIドキュメントのパスかどうかを判定
func isOpenAPIPath(path string) bool {
	openAPIPaths := []string{"/docs", "/openapi.json", "/openapi.yaml", "/openapi-3.0.json", "/schemas"}
	for _, p := range openAPIPaths {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

// CheckAccessLevel はコンテキストから許容する公開レベルを取得し、エンドポイントの公開レベルと比較
// エンドポイントがprivateの場合、Auth0 JWTでのみアクセス可能
func CheckAccessLevel(ctx context.Context, endpointLevel AccessLevel) error {
	allowedLevel, ok := ctx.Value(AllowedAccessLevelKey).(AccessLevel)
	if !ok {
		return errors.New("access level not found in context")
	}

	// エンドポイントがprivateで、JWTがpublicの場合はエラー
	if endpointLevel == AccessLevelPrivate && allowedLevel == AccessLevelPublic {
		return errors.New("private API requires Auth0 authentication")
	}

	return nil
}

// GetAllowedAccessLevel はコンテキストから許容する公開レベルを取得
func GetAllowedAccessLevel(ctx context.Context) (AccessLevel, bool) {
	level, ok := ctx.Value(AllowedAccessLevelKey).(AccessLevel)
	return level, ok
}
