package idgen

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// GenerateUUIDv7 はUUIDv7を生成し、ハイフン抜き小文字32文字の文字列として返す
func GenerateUUIDv7() (string, error) {
	// 1. UUIDv7を生成
	u, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("failed to generate UUIDv7: %w", err)
	}

	// 2. ハイフンを削除
	uuidStr := strings.ReplaceAll(u.String(), "-", "")

	// 3. 小文字に変換（uuid.Stringは既に小文字だが、念のため）
	uuidStr = strings.ToLower(uuidStr)

	// 4. 32文字の文字列として返す
	return uuidStr, nil
}
