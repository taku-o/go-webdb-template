package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/example/go-webdb-template/internal/config"
	"github.com/labstack/echo/v4"
)

// AuthMiddleware は認証ミドルウェア
type AuthMiddleware struct {
	validator *JWTValidator
}

// NewAuthMiddleware は新しい認証ミドルウェアを作成
func NewAuthMiddleware(cfg *config.APIConfig, env string) *AuthMiddleware {
	return &AuthMiddleware{
		validator: NewJWTValidator(cfg, env),
	}
}

// Middleware はHTTPミドルウェア関数を返す
func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// AuthorizationヘッダーからJWTトークンを取得
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.writeErrorResponse(w, http.StatusUnauthorized, "Authorization header is required")
			return
		}

		// Bearerトークンの抽出
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			m.writeErrorResponse(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		tokenString := parts[1]

		// JWT検証
		claims, err := m.validator.ValidateJWT(tokenString)
		if err != nil {
			m.writeErrorResponse(w, http.StatusUnauthorized, "Invalid API key")
			return
		}

		// スコープ検証
		if err := validateScope(claims, r.Method); err != nil {
			m.writeErrorResponse(w, http.StatusForbidden, "Insufficient scope")
			return
		}

		// 次のハンドラーを実行
		next.ServeHTTP(w, r)
	})
}

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

// writeErrorResponse はエラーレスポンスを書き込む
func (m *AuthMiddleware) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]interface{}{
		"code":    statusCode,
		"message": message,
	}

	json.NewEncoder(w).Encode(response)
}

// NewEchoAuthMiddleware は新しいEcho形式の認証ミドルウェアを作成
func NewEchoAuthMiddleware(cfg *config.APIConfig, env string) echo.MiddlewareFunc {
	validator := NewJWTValidator(cfg, env)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// AuthorizationヘッダーからJWTトークンを取得
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"code":    http.StatusUnauthorized,
					"message": "Authorization header is required",
				})
			}

			// Bearerトークンの抽出
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"code":    http.StatusUnauthorized,
					"message": "Invalid authorization header format",
				})
			}

			tokenString := parts[1]

			// JWT検証
			claims, err := validator.ValidateJWT(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"code":    http.StatusUnauthorized,
					"message": "Invalid API key",
				})
			}

			// スコープ検証
			if err := validateScope(claims, c.Request().Method); err != nil {
				return c.JSON(http.StatusForbidden, map[string]interface{}{
					"code":    http.StatusForbidden,
					"message": "Insufficient scope",
				})
			}

			// 次のハンドラーを実行
			return next(c)
		}
	}
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
