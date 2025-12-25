package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/example/go-webdb-template/internal/config"
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
		if err := m.validateScope(claims, r.Method); err != nil {
			m.writeErrorResponse(w, http.StatusForbidden, "Insufficient scope")
			return
		}

		// 次のハンドラーを実行
		next.ServeHTTP(w, r)
	})
}

// validateScope はスコープを検証
func (m *AuthMiddleware) validateScope(claims *JWTClaims, method string) error {
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
