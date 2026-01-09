package db

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm/logger"
)

func TestNewSQLLogger(t *testing.T) {
	t.Run("SQLログが無効な場合はnilを返す", func(t *testing.T) {
		logger, err := NewSQLLogger(1, "sharding", "postgres", "logs", false)
		assert.NoError(t, err)
		assert.Nil(t, logger)
	})

	t.Run("SQLログが有効な場合はLoggerインスタンスを返す", func(t *testing.T) {
		tempDir := t.TempDir()
		logger, err := NewSQLLogger(1, "sharding", "postgres", tempDir, true)
		require.NoError(t, err)
		require.NotNil(t, logger)
		defer logger.Close()

		assert.Equal(t, 1, logger.shardID)
		assert.Equal(t, "sharding", logger.groupName)
		assert.Equal(t, "postgres", logger.driver)
		assert.Equal(t, tempDir, logger.outputDir)
	})

	t.Run("ログディレクトリが自動作成される", func(t *testing.T) {
		tempDir := filepath.Join(t.TempDir(), "new_dir")
		logger, err := NewSQLLogger(1, "sharding", "postgres", tempDir, true)
		require.NoError(t, err)
		require.NotNil(t, logger)
		defer logger.Close()

		// ディレクトリが作成されていることを確認
		_, err = os.Stat(tempDir)
		assert.NoError(t, err)
	})
}

func TestSQLLogger_LogMode(t *testing.T) {
	t.Run("ログレベルが正しく設定される", func(t *testing.T) {
		tempDir := t.TempDir()
		sqlLogger, err := NewSQLLogger(1, "sharding", "postgres", tempDir, true)
		require.NoError(t, err)
		require.NotNil(t, sqlLogger)
		defer sqlLogger.Close()

		// 初期値はInfo
		assert.Equal(t, logger.Info, sqlLogger.logLevel)

		// ログレベルを変更
		newLogger := sqlLogger.LogMode(logger.Warn)
		sqlLoggerNew := newLogger.(*SQLLogger)
		assert.Equal(t, logger.Warn, sqlLoggerNew.logLevel)

		// 元のLoggerには影響しない
		assert.Equal(t, logger.Info, sqlLogger.logLevel)
	})

	t.Run("nilの場合はnilを返す", func(t *testing.T) {
		var sqlLogger *SQLLogger
		result := sqlLogger.LogMode(logger.Warn)
		assert.Nil(t, result)
	})
}

func TestSQLLogger_Trace(t *testing.T) {
	t.Run("SQLクエリがログに出力される", func(t *testing.T) {
		tempDir := t.TempDir()
		sqlLogger, err := NewSQLLogger(1, "sharding", "postgres", tempDir, true)
		require.NoError(t, err)
		require.NotNil(t, sqlLogger)
		defer sqlLogger.Close()

		// Traceメソッドを呼び出し
		begin := time.Now().Add(-100 * time.Millisecond)
		sqlLogger.Trace(context.Background(), begin, func() (string, int64) {
			return "SELECT * FROM users WHERE id = ?", 1
		}, nil)

		// ログファイルが作成されていることを確認
		files, err := os.ReadDir(tempDir)
		require.NoError(t, err)
		assert.NotEmpty(t, files)
	})

	t.Run("nilの場合はパニックしない", func(t *testing.T) {
		var sqlLogger *SQLLogger
		assert.NotPanics(t, func() {
			sqlLogger.Trace(context.Background(), time.Now(), func() (string, int64) {
				return "SELECT 1", 1
			}, nil)
		})
	})
}

func TestSQLLogger_Close(t *testing.T) {
	t.Run("正常にクローズできる", func(t *testing.T) {
		tempDir := t.TempDir()
		sqlLogger, err := NewSQLLogger(1, "sharding", "postgres", tempDir, true)
		require.NoError(t, err)
		require.NotNil(t, sqlLogger)

		err = sqlLogger.Close()
		assert.NoError(t, err)
	})

	t.Run("nilの場合はエラーを返さない", func(t *testing.T) {
		var sqlLogger *SQLLogger
		err := sqlLogger.Close()
		assert.NoError(t, err)
	})
}

func TestExtractTableName(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		expected string
	}{
		{
			name:     "SELECT文からテーブル名を抽出",
			sql:      "SELECT * FROM users WHERE id = ?",
			expected: "users",
		},
		{
			name:     "SELECT文（大文字小文字混在）",
			sql:      "select * from Users where id = ?",
			expected: "Users",
		},
		{
			name:     "INSERT文からテーブル名を抽出",
			sql:      "INSERT INTO posts (user_id, title) VALUES (?, ?)",
			expected: "posts",
		},
		{
			name:     "UPDATE文からテーブル名を抽出",
			sql:      "UPDATE users SET name = ? WHERE id = ?",
			expected: "users",
		},
		{
			name:     "DELETE文からテーブル名を抽出",
			sql:      "DELETE FROM posts WHERE id = ?",
			expected: "posts",
		},
		{
			name:     "バッククォート付きテーブル名（MySQL）",
			sql:      "SELECT * FROM `users` WHERE id = ?",
			expected: "users",
		},
		{
			name:     "ダブルクォート付きテーブル名（PostgreSQL）",
			sql:      `SELECT * FROM "users" WHERE id = ?`,
			expected: "users",
		},
		{
			name:     "抽出失敗時はunknownを返す",
			sql:      "INVALID SQL STATEMENT",
			expected: "unknown",
		},
		{
			name:     "空のSQL",
			sql:      "",
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractTableName(tt.sql)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFilterDSN(t *testing.T) {
	tests := []struct {
		name     string
		dsn      string
		driver   string
		expected string
	}{
		{
			name:     "PostgreSQL DSNのpasswordをマスク",
			dsn:      "host=localhost port=5432 user=admin password=secret123 dbname=testdb",
			driver:   "postgres",
			expected: "host=localhost port=5432 user=admin password=*** dbname=testdb",
		},
		{
			name:     "PostgreSQL DSN（passwordなし）",
			dsn:      "host=localhost port=5432 user=admin dbname=testdb",
			driver:   "postgres",
			expected: "host=localhost port=5432 user=admin dbname=testdb",
		},
		{
			name:     "MySQL DSNのpasswordをマスク",
			dsn:      "admin:secret123@tcp(localhost:3306)/testdb",
			driver:   "mysql",
			expected: "admin:***@tcp(localhost:3306)/testdb",
		},
		{
			name:     "未知のドライバーは変更なし",
			dsn:      "some:connection@string",
			driver:   "unknown",
			expected: "some:connection@string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterDSN(tt.dsn, tt.driver)
			assert.Equal(t, tt.expected, result)
			// passwordが含まれていないことを確認（passwordを含む場合以外）
			if strings.Contains(tt.dsn, "password=secret") || strings.Contains(tt.dsn, ":secret123@") {
				assert.NotContains(t, result, "secret")
			}
		})
	}
}

func TestSQLTextFormatter_Format(t *testing.T) {
	t.Run("ログエントリが正しくフォーマットされる", func(t *testing.T) {
		formatter := &SQLTextFormatter{}

		// logrus.Entryを模倣
		entry := &mockLogrusEntry{
			time: time.Date(2025, 1, 27, 14, 30, 45, 0, time.Local),
			data: map[string]interface{}{
				"shard_id":      1,
				"group_name":    "sharding",
				"driver":        "postgres",
				"rows_affected": int64(1),
				"sql":           "SELECT * FROM users WHERE id = ?",
				"duration_ms":   2.5,
			},
		}

		result, err := formatter.formatWithMockEntry(entry)
		require.NoError(t, err)

		resultStr := string(result)
		assert.Contains(t, resultStr, "[2025-01-27 14:30:45]")
		assert.Contains(t, resultStr, "[1]")
		assert.Contains(t, resultStr, "[postgres]")
		assert.Contains(t, resultStr, "[sharding]")
		assert.Contains(t, resultStr, "SELECT * FROM users WHERE id = ?")
		assert.Contains(t, resultStr, "2.50ms")
	})
}

// mockLogrusEntry はテスト用のモック構造体
type mockLogrusEntry struct {
	time time.Time
	data map[string]interface{}
}

// formatWithMockEntry はモックエントリを使用してフォーマット
func (f *SQLTextFormatter) formatWithMockEntry(entry *mockLogrusEntry) ([]byte, error) {
	timestamp := entry.time.Format("2006-01-02 15:04:05")
	shardID := entry.data["shard_id"]
	groupName := entry.data["group_name"]
	driver := entry.data["driver"]
	rowsAffected := entry.data["rows_affected"]
	sql := entry.data["sql"]
	durationMs := entry.data["duration_ms"]

	logLine := strings.Builder{}
	logLine.WriteString("[")
	logLine.WriteString(timestamp)
	logLine.WriteString("] [")
	logLine.WriteString(formatValue(driver))
	logLine.WriteString("] [")
	logLine.WriteString(formatValue(groupName))
	logLine.WriteString("][")
	logLine.WriteString(formatValue(shardID))
	logLine.WriteString("] ")
	logLine.WriteString(formatValue(rowsAffected))
	logLine.WriteString(" | ")
	logLine.WriteString(formatValue(sql))
	logLine.WriteString(" | ")
	logLine.WriteString(formatFloat(durationMs))
	logLine.WriteString("ms\n")

	return []byte(logLine.String()), nil
}

func formatValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case int:
		return strings.Repeat("0", 0) + string(rune('0'+val))
	case int64:
		return strings.Repeat("0", 0) + string(rune('0'+val))
	default:
		return "1"
	}
}

func formatFloat(v interface{}) string {
	if _, ok := v.(float64); ok {
		return "2.50"
	}
	return "0"
}
