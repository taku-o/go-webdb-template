package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/example/go-webdb-template/internal/config"
	"github.com/example/go-webdb-template/internal/db"
	"github.com/example/go-webdb-template/internal/model"
	"github.com/example/go-webdb-template/internal/repository"
	"github.com/example/go-webdb-template/internal/service"
)

func main() {
	// コマンドライン引数の解析
	limit := flag.Int("limit", 20, "Number of users to output (default: 20, max: 100)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	// 引数のバリデーション
	validatedLimit, err, warning := validateLimit(*limit)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	if warning {
		log.Printf("Warning: limit exceeds maximum (100), using 100")
	}

	// 設定ファイルの読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// GORM DB Managerの初期化
	gormManager, err := db.NewGORMManager(cfg)
	if err != nil {
		log.Fatalf("Failed to create GORM manager: %v", err)
	}
	defer gormManager.CloseAll()

	// すべてのShardへの接続確認
	if err := gormManager.PingAll(); err != nil {
		log.Fatalf("Failed to ping databases: %v", err)
	}

	// Repository層の初期化
	userRepo := repository.NewUserRepositoryGORM(gormManager)

	// Service層の初期化
	userService := service.NewUserService(userRepo)

	// ユーザー一覧の取得
	ctx := context.Background()
	users, err := userService.ListUsers(ctx, validatedLimit, 0)
	if err != nil {
		log.Fatalf("Failed to list users: %v", err)
	}

	// limit件に制限（クロスシャードクエリのため）
	if len(users) > validatedLimit {
		users = users[:validatedLimit]
	}

	// TSV形式での出力
	printUsersTSV(users)

	os.Exit(0)
}

// validateLimit validates the limit parameter and returns the validated limit,
// an error if invalid, and a boolean indicating if a warning was issued.
func validateLimit(limit int) (int, error, bool) {
	if limit < 1 {
		return 0, errors.New("limit must be at least 1"), false
	}
	if limit > 100 {
		return 100, nil, true
	}
	return limit, nil, false
}

// printUsersTSV prints users in TSV format to stdout.
func printUsersTSV(users []*model.User) {
	// ヘッダー行の出力
	fmt.Println("ID\tName\tEmail\tCreatedAt\tUpdatedAt")

	// 各ユーザー情報の出力
	for _, user := range users {
		fmt.Printf("%d\t%s\t%s\t%s\t%s\n",
			user.ID,
			user.Name,
			user.Email,
			user.CreatedAt.Format(time.RFC3339),
			user.UpdatedAt.Format(time.RFC3339),
		)
	}
}
