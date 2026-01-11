package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/taku-o/go-webdb-template/internal/config"
	"github.com/taku-o/go-webdb-template/internal/db"
	"github.com/taku-o/go-webdb-template/internal/model"
	"github.com/taku-o/go-webdb-template/internal/repository"
	"github.com/taku-o/go-webdb-template/internal/service"
	"github.com/taku-o/go-webdb-template/internal/usecase/cli"
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

	// GroupManagerの初期化
	groupManager, err := db.NewGroupManager(cfg)
	if err != nil {
		log.Fatalf("Failed to create group manager: %v", err)
	}
	defer groupManager.CloseAll()

	// すべてのデータベースへの接続確認
	if err := groupManager.PingAll(); err != nil {
		log.Fatalf("Failed to ping databases: %v", err)
	}

	// Repository層の初期化
	dmUserRepo := repository.NewDmUserRepository(groupManager)

	// Service層の初期化
	dmUserService := service.NewDmUserService(dmUserRepo)

	// Usecase層の初期化
	listDmUsersUsecase := cli.NewListDmUsersUsecase(dmUserService)

	// ユーザー一覧の取得
	ctx := context.Background()
	dmUsers, err := listDmUsersUsecase.ListDmUsers(ctx, validatedLimit, 0)
	if err != nil {
		log.Fatalf("Failed to list users: %v", err)
	}

	// limit件に制限（クロスシャードクエリのため）
	if len(dmUsers) > validatedLimit {
		dmUsers = dmUsers[:validatedLimit]
	}

	// TSV形式での出力
	printDmUsersTSV(dmUsers)

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

// printDmUsersTSV prints dm_users in TSV format to stdout.
func printDmUsersTSV(dmUsers []*model.DmUser) {
	// ヘッダー行の出力
	fmt.Println("ID\tName\tEmail\tCreatedAt\tUpdatedAt")

	// 各ユーザー情報の出力
	for _, dmUser := range dmUsers {
		fmt.Printf("%s\t%s\t%s\t%s\t%s\n",
			dmUser.ID,
			dmUser.Name,
			dmUser.Email,
			dmUser.CreatedAt.Format(time.RFC3339),
			dmUser.UpdatedAt.Format(time.RFC3339),
		)
	}
}
