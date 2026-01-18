package service

import (
	"fmt"
	"os"
	"time"
)

// DelayPrintService はジョブ処理のビジネスユーティリティロジックを提供
type DelayPrintService struct{}

// NewDelayPrintService は新しいDelayPrintServiceを作成
func NewDelayPrintService() *DelayPrintService {
	return &DelayPrintService{}
}

// PrintMessage は標準出力にタイムスタンプ付きでメッセージを出力
func (s *DelayPrintService) PrintMessage(message string) error {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] %s\n", timestamp, message)
	os.Stdout.Sync() // バッファをフラッシュ
	return nil
}
