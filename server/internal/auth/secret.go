package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateSecretKey は32バイト（256ビット）のランダムな秘密鍵を生成してBase64エンコードして返す
func GenerateSecretKey() (string, error) {
	// 32バイト（256ビット）のランダムな秘密鍵を生成
	secretKey := make([]byte, 32)
	if _, err := rand.Read(secretKey); err != nil {
		return "", fmt.Errorf("failed to generate secret key: %w", err)
	}

	// Base64エンコード
	encoded := base64.StdEncoding.EncodeToString(secretKey)

	return encoded, nil
}
