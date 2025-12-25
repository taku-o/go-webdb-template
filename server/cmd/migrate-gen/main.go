package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// マイグレーション生成ツール
// テンプレートから32分割されたテーブル定義を生成する
func main() {
	// デフォルトパス
	templateDir := "../../../db/migrations/sharding/templates"
	outputDir := "../../../db/migrations/sharding/generated"

	// コマンドライン引数の処理
	if len(os.Args) > 1 {
		templateDir = os.Args[1]
	}
	if len(os.Args) > 2 {
		outputDir = os.Args[2]
	}

	// 出力ディレクトリの作成
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// テンプレートファイルの検索
	templates, err := filepath.Glob(filepath.Join(templateDir, "*.sql.template"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding templates: %v\n", err)
		os.Exit(1)
	}

	if len(templates) == 0 {
		fmt.Println("No templates found in", templateDir)
		os.Exit(0)
	}

	fmt.Printf("Found %d templates\n", len(templates))

	// 各テンプレートを処理
	for _, templatePath := range templates {
		baseName := filepath.Base(templatePath)
		tableName := strings.TrimSuffix(baseName, ".sql.template")

		fmt.Printf("Processing template: %s\n", baseName)

		// テンプレート読み込み
		templateContent, err := os.ReadFile(templatePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading template %s: %v\n", templatePath, err)
			continue
		}

		// 32個のテーブル定義を生成
		for i := 0; i < 32; i++ {
			suffix := fmt.Sprintf("%03d", i)
			fullTableName := fmt.Sprintf("%s_%s", tableName, suffix)

			// テンプレート変数を置換
			content := string(templateContent)
			content = strings.ReplaceAll(content, "{TABLE_NAME}", fullTableName)
			content = strings.ReplaceAll(content, "{TABLE_SUFFIX}", suffix)

			// 出力ファイル名
			outputPath := filepath.Join(outputDir, fmt.Sprintf("%s.sql", fullTableName))

			// ファイル書き込み
			if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", outputPath, err)
				continue
			}
		}

		fmt.Printf("  Generated 32 tables for %s\n", tableName)
	}

	fmt.Println("Done!")
}
