package main

import (
	"context"
	"log"
	"os"

	"github.com/taku-o/go-webdb-template/internal/config"
	"github.com/taku-o/go-webdb-template/internal/db"
	"github.com/taku-o/go-webdb-template/internal/repository"
	"github.com/taku-o/go-webdb-template/internal/service"
	"github.com/taku-o/go-webdb-template/internal/usecase/cli"
)

const (
	totalCount = 100
)

func main() {
	// 1. 設定ファイルの読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. GroupManagerの初期化
	groupManager, err := db.NewGroupManager(cfg)
	if err != nil {
		log.Fatalf("Failed to create group manager: %v", err)
	}
	defer groupManager.CloseAll()

	// 3. データベース接続確認
	if err := groupManager.PingAll(); err != nil {
		log.Fatalf("Failed to ping databases: %v", err)
	}

	log.Println("Starting sample data generation...")

	// 4. Repository層の初期化
	dmUserRepository := repository.NewDmUserRepository(groupManager)
	dmPostRepository := repository.NewDmPostRepository(groupManager)
	dmNewsRepository := repository.NewDmNewsRepository(groupManager)

	// 5. Service層の初期化
	tableSelector := db.NewTableSelector(db.DBShardingTableCount, db.DBShardingTablesPerDB)
	generateSampleService := service.NewGenerateSampleService(
		dmUserRepository,
		dmPostRepository,
		dmNewsRepository,
		tableSelector,
	)

	// 6. Usecase層の初期化
	generateSampleUsecase := cli.NewGenerateSampleUsecase(generateSampleService)

	// 7. サンプルデータの生成
	ctx := context.Background()
	if err := generateSampleUsecase.GenerateSampleData(ctx, totalCount); err != nil {
		log.Fatalf("Failed to generate sample data: %v", err)
	}

	// 8. 生成完了メッセージ
	log.Println("Sample data generation completed successfully")
	os.Exit(0)
}
