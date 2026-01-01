package jobqueue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstants_JobTypeDelayPrint(t *testing.T) {
	// ジョブタイプが正しく定義されていること
	assert.Equal(t, "demo:delay_print", JobTypeDelayPrint)
}

func TestConstants_DefaultDelaySeconds(t *testing.T) {
	// デフォルト遅延時間が3分（180秒）であること
	assert.Equal(t, 180, DefaultDelaySeconds)
}

func TestConstants_DefaultMaxRetry(t *testing.T) {
	// デフォルト最大リトライ回数が10回であること
	assert.Equal(t, 10, DefaultMaxRetry)
}
