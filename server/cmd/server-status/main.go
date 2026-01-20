package main

import (
	"context"
	"fmt"
	"os"

	"github.com/taku-o/go-webdb-template/internal/service"
	"github.com/taku-o/go-webdb-template/internal/usecase/cli"
)

// printResults は結果を表形式で表示する
func printResults(results []service.ServerStatus) {
	// ヘッダー行
	fmt.Println("サーバー          | ポート | 状態")
	fmt.Println("------------------|-------|--------")

	// 各サーバーの状態を表示
	for _, result := range results {
		fmt.Printf("%-17s | %-5d | %s\n",
			result.Server.Name,
			result.Server.Port,
			result.Status,
		)
	}
}

func main() {
	// Service層の初期化
	serverStatusService := service.NewServerStatusService()

	// Usecase層の初期化
	serverStatusUsecase := cli.NewServerStatusUsecase(serverStatusService)

	// サーバー状態の確認
	ctx := context.Background()
	results, err := serverStatusUsecase.ListServerStatus(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// 結果を表形式で表示
	printResults(results)

	os.Exit(0)
}
