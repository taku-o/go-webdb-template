package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
)

func main() {
	// 32バイト（256ビット）のランダムな秘密鍵を生成
	secretKey := make([]byte, 32)
	if _, err := rand.Read(secretKey); err != nil {
		log.Fatalf("Failed to generate secret key: %v", err)
	}

	// Base64エンコード
	encoded := base64.StdEncoding.EncodeToString(secretKey)

	// 標準出力に表示
	fmt.Println(encoded)

	os.Exit(0)
}
