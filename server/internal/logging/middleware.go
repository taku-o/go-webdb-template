package logging

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	// MaxBodySize はログに出力するリクエストボディの最大サイズ（1MB）
	MaxBodySize = 1 * 1024 * 1024
)

// responseWriter はステータスコードを記録するResponseWriterラッパー
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader はステータスコードを記録してから元のWriteHeaderを呼び出す
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// AccessLogMiddleware はアクセスログを記録するHTTPミドルウェア
type AccessLogMiddleware struct {
	accessLogger *AccessLogger
}

// NewAccessLogMiddleware は新しいAccessLogMiddlewareを作成
func NewAccessLogMiddleware(accessLogger *AccessLogger) *AccessLogMiddleware {
	return &AccessLogMiddleware{
		accessLogger: accessLogger,
	}
}

// Middleware はHTTPミドルウェア関数を返す
func (m *AccessLogMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエスト開始時刻を記録
		startTime := time.Now()

		// ヘッダー情報を取得
		headers := formatHeaders(r.Header)

		// リクエストボディを取得（POST/PUT/PATCHの場合）
		requestBody := ""
		if shouldLogBody(r) {
			requestBody = readRequestBody(r)
		}

		// レスポンスライターをラップしてステータスコードを取得
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// 次のハンドラーを実行
		next.ServeHTTP(rw, r)

		// レスポンス時間を計算
		responseTime := time.Since(startTime)
		responseTimeMs := float64(responseTime.Nanoseconds()) / 1000000.0

		// リモートIPアドレスを取得
		remoteIP := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			remoteIP = forwarded
		}

		// User-Agentを取得
		userAgent := r.Header.Get("User-Agent")
		if userAgent == "" {
			userAgent = "-"
		}

		// アクセスログを出力
		m.accessLogger.LogAccess(
			r.Method,
			r.URL.Path,
			r.Proto,
			rw.statusCode,
			responseTimeMs,
			remoteIP,
			userAgent,
			headers,
			requestBody,
		)
	})
}

// formatHeaders はヘッダーを文字列にフォーマット
func formatHeaders(headers http.Header) string {
	var parts []string
	for key, values := range headers {
		for _, value := range values {
			parts = append(parts, key+": "+value)
		}
	}
	return strings.Join(parts, "; ")
}

// shouldLogBody はリクエストボディをログに出力すべきか判定
func shouldLogBody(r *http.Request) bool {
	// POST/PUT/PATCHのみ対象
	method := strings.ToUpper(r.Method)
	if method != "POST" && method != "PUT" && method != "PATCH" {
		return false
	}

	// Content-Lengthが大きすぎる場合は出力しない
	if r.ContentLength > MaxBodySize {
		return false
	}

	// 画像/動画/音声の場合は出力しない
	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "image/") ||
		strings.HasPrefix(contentType, "video/") ||
		strings.HasPrefix(contentType, "audio/") ||
		strings.HasPrefix(contentType, "multipart/form-data") {
		return false
	}

	return true
}

// readRequestBody はリクエストボディを読み取り、再度読み取れるように復元
func readRequestBody(r *http.Request) string {
	if r.Body == nil {
		return ""
	}

	// ボディを読み取り
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return "[read error]"
	}

	// ボディを復元（後続のハンドラーが読み取れるように）
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// サイズチェック（念のため）
	if len(bodyBytes) > MaxBodySize {
		return "[body too large]"
	}

	return string(bodyBytes)
}

// NewEchoAccessLogMiddleware はEcho用のアクセスログミドルウェアを作成
func NewEchoAccessLogMiddleware(accessLogger *AccessLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// リクエスト開始時刻を記録
			startTime := time.Now()

			req := c.Request()

			// ヘッダー情報を取得
			headers := formatHeaders(req.Header)

			// リクエストボディを取得（POST/PUT/PATCHの場合）
			requestBody := ""
			if shouldLogBody(req) {
				requestBody = readRequestBody(req)
			}

			// 次のハンドラーを実行
			err := next(c)

			// レスポンス時間を計算
			responseTime := time.Since(startTime)
			responseTimeMs := float64(responseTime.Nanoseconds()) / 1000000.0

			// リモートIPアドレスを取得
			remoteIP := c.RealIP()
			if remoteIP == "" {
				remoteIP = req.RemoteAddr
			}

			// User-Agentを取得
			userAgent := req.Header.Get("User-Agent")
			if userAgent == "" {
				userAgent = "-"
			}

			// ステータスコードを取得
			statusCode := c.Response().Status

			// アクセスログを出力
			accessLogger.LogAccess(
				req.Method,
				req.URL.Path,
				req.Proto,
				statusCode,
				responseTimeMs,
				remoteIP,
				userAgent,
				headers,
				requestBody,
			)

			return err
		}
	}
}
