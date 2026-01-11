package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/taku-o/go-webdb-template/internal/service"
	"github.com/taku-o/go-webdb-template/internal/usecase/cli"
)

func main() {
	// Service層の初期化
	secretService := service.NewSecretService()

	// Usecase層の初期化
	generateSecretUsecase := cli.NewGenerateSecretUsecase(secretService)

	// 秘密鍵の生成
	ctx := context.Background()
	secretKey, err := generateSecretUsecase.GenerateSecret(ctx)
	if err != nil {
		log.Fatalf("Failed to generate secret key: %v", err)
	}

	// 標準出力に表示
	fmt.Println(secretKey)

	os.Exit(0)
}
