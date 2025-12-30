package idgen

import (
	"fmt"
	"sync"

	"github.com/sony/sonyflake"
)

var (
	sf   *sonyflake.Sonyflake
	once sync.Once
)

// initSonyflake はsonyflakeインスタンスを初期化する（初回のみ）
func initSonyflake() {
	once.Do(func() {
		st := sonyflake.Settings{}
		sf = sonyflake.NewSonyflake(st)
		if sf == nil {
			panic("failed to initialize sonyflake")
		}
	})
}

// GenerateSonyflakeID はsonyflakeを使用して一意のIDを生成する
func GenerateSonyflakeID() (int64, error) {
	initSonyflake()

	id, err := sf.NextID()
	if err != nil {
		return 0, fmt.Errorf("failed to generate ID: %w", err)
	}

	return int64(id), nil
}
