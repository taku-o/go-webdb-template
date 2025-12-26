package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/example/go-webdb-template/internal/config"
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
func NewHumaAuthMiddleware(cfg *config.APIConfig, env string) func(ctx huma.Context, next func(huma.Context)) {
	validator := NewJWTValidator(cfg, env)

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

		// JWT検証
		claims, err := validator.ValidateJWT(tokenString)
		if err != nil {
			writeHumaError(ctx, http.StatusUnauthorized, "Invalid API key")
			return
		}

		// スコープ検証
		if err := validateScope(claims, ctx.Method()); err != nil {
			writeHumaError(ctx, http.StatusForbidden, "Insufficient scope")
			return
		}

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
