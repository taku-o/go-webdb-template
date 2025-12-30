package idgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateSonyflakeID(t *testing.T) {
	// ID生成のテスト
	id1, err := GenerateSonyflakeID()
	require.NoError(t, err)
	assert.Greater(t, id1, int64(0), "ID should be greater than 0")
}

func TestGenerateSonyflakeID_Uniqueness(t *testing.T) {
	// 一意性のテスト
	ids := make(map[int64]bool)
	for i := 0; i < 1000; i++ {
		id, err := GenerateSonyflakeID()
		require.NoError(t, err)
		require.False(t, ids[id], "duplicate ID: %d", id)
		ids[id] = true
	}
}

func TestGenerateSonyflakeID_Incremental(t *testing.T) {
	// 時系列順序のテスト（後で生成したIDは前より大きい）
	id1, err := GenerateSonyflakeID()
	require.NoError(t, err)

	id2, err := GenerateSonyflakeID()
	require.NoError(t, err)

	assert.Greater(t, id2, id1, "ID2 should be greater than ID1")
}
