package logging

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAccessLogger(t *testing.T) {
	t.Run("正常系: AccessLoggerが正常に作成される", func(t *testing.T) {
		tmpDir := t.TempDir()

		logger, err := NewAccessLogger("api", tmpDir)
		require.NoError(t, err)
		require.NotNil(t, logger)
		defer logger.Close()

		// ディレクトリが作成されていることを確認
		_, err = os.Stat(tmpDir)
		assert.NoError(t, err)
	})

	t.Run("正常系: ログディレクトリが自動作成される", func(t *testing.T) {
		tmpDir := t.TempDir()
		logDir := filepath.Join(tmpDir, "nested", "logs")

		logger, err := NewAccessLogger("api", logDir)
		require.NoError(t, err)
		require.NotNil(t, logger)
		defer logger.Close()

		// ネストされたディレクトリが作成されていることを確認
		_, err = os.Stat(logDir)
		assert.NoError(t, err)
	})

	t.Run("正常系: 相対パスが正しく処理される", func(t *testing.T) {
		// 一時ディレクトリを作成し、その中に相対パスでログを作成
		tmpDir := t.TempDir()
		originalWd, err := os.Getwd()
		require.NoError(t, err)

		err = os.Chdir(tmpDir)
		require.NoError(t, err)
		defer os.Chdir(originalWd)

		logger, err := NewAccessLogger("api", "logs")
		require.NoError(t, err)
		require.NotNil(t, logger)
		defer logger.Close()

		// logsディレクトリが作成されていることを確認
		_, err = os.Stat(filepath.Join(tmpDir, "logs"))
		assert.NoError(t, err)
	})
}

func TestAccessLogger_LogAccess(t *testing.T) {
	t.Run("正常系: アクセスログが正常に出力される", func(t *testing.T) {
		tmpDir := t.TempDir()

		logger, err := NewAccessLogger("api", tmpDir)
		require.NoError(t, err)
		require.NotNil(t, logger)

		// アクセスログを出力
		logger.LogAccess("GET", "/api/dm-users", "HTTP/1.1", 200, 15.2, "192.168.1.100", "Mozilla/5.0", "", "")

		// loggerをクローズしてファイルをフラッシュ
		err = logger.Close()
		require.NoError(t, err)

		// ログファイルが作成されていることを確認
		today := time.Now().Format("2006-01-02")
		logFileName := "api-access-" + today + ".log"
		logFilePath := filepath.Join(tmpDir, logFileName)

		_, err = os.Stat(logFilePath)
		require.NoError(t, err)

		// ログファイルの内容を確認
		content, err := os.ReadFile(logFilePath)
		require.NoError(t, err)

		logContent := string(content)
		assert.Contains(t, logContent, "GET")
		assert.Contains(t, logContent, "/api/dm-users")
		assert.Contains(t, logContent, "HTTP/1.1")
		assert.Contains(t, logContent, "200")
		assert.Contains(t, logContent, "192.168.1.100")
		assert.Contains(t, logContent, "Mozilla/5.0")
	})

	t.Run("正常系: ログエントリのフォーマットが正しい", func(t *testing.T) {
		tmpDir := t.TempDir()

		logger, err := NewAccessLogger("api", tmpDir)
		require.NoError(t, err)
		require.NotNil(t, logger)

		logger.LogAccess("POST", "/api/dm-users", "HTTP/1.1", 201, 23.5, "192.168.1.100", "curl/7.64.1", "", "")

		err = logger.Close()
		require.NoError(t, err)

		today := time.Now().Format("2006-01-02")
		logFileName := "api-access-" + today + ".log"
		logFilePath := filepath.Join(tmpDir, logFileName)

		content, err := os.ReadFile(logFilePath)
		require.NoError(t, err)

		logContent := string(content)
		// ログエントリがJSON形式であることを確認
		assert.True(t, strings.HasPrefix(logContent, "{"))
		assert.Contains(t, logContent, `"method":"POST"`)
		assert.Contains(t, logContent, `"path":"/api/dm-users"`)
		assert.Contains(t, logContent, `"status_code":201`)
		assert.Contains(t, logContent, `"response_time_ms":23.5`)
		assert.Contains(t, logContent, `"remote_ip":"192.168.1.100"`)
		assert.Contains(t, logContent, `"user_agent":"curl/7.64.1"`)
	})
}

func TestAccessLogger_Close(t *testing.T) {
	t.Run("正常系: Closeが正常に実行される", func(t *testing.T) {
		tmpDir := t.TempDir()

		logger, err := NewAccessLogger("api", tmpDir)
		require.NoError(t, err)
		require.NotNil(t, logger)

		err = logger.Close()
		assert.NoError(t, err)
	})

	t.Run("正常系: nilのloggerでもCloseがエラーにならない", func(t *testing.T) {
		logger := &AccessLogger{}
		err := logger.Close()
		assert.NoError(t, err)
	})
}

func TestCustomTextFormatter_Format(t *testing.T) {
	t.Run("正常系: ログエントリが正しくフォーマットされる", func(t *testing.T) {
		formatter := &CustomTextFormatter{}

		// ログエントリを作成
		fields := map[string]interface{}{
			"method":           "GET",
			"path":             "/api/dm-users",
			"protocol":         "HTTP/1.1",
			"status_code":      200,
			"response_time_ms": 15.2,
			"remote_ip":        "192.168.1.100",
			"user_agent":       "Mozilla/5.0",
			"headers":          "",
			"request_body":     "",
		}

		now := time.Now()
		formatted, err := formatter.Format(now, fields)
		require.NoError(t, err)

		result := string(formatted)
		// JSON形式であることを確認
		assert.True(t, strings.HasPrefix(result, "{"))
		assert.Contains(t, result, now.Format("2006-01-02"))
		assert.Contains(t, result, `"method":"GET"`)
		assert.Contains(t, result, `"path":"/api/dm-users"`)
		assert.Contains(t, result, `"protocol":"HTTP/1.1"`)
		assert.Contains(t, result, `"status_code":200`)
		assert.Contains(t, result, `"response_time_ms":15.2`)
		assert.Contains(t, result, `"remote_ip":"192.168.1.100"`)
		assert.Contains(t, result, `"user_agent":"Mozilla/5.0"`)
	})
}
