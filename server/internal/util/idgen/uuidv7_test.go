package idgen

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateUUIDv7_ReturnsNonEmpty(t *testing.T) {
	uuid, err := GenerateUUIDv7()
	require.NoError(t, err)
	assert.NotEmpty(t, uuid)
}

func TestGenerateUUIDv7_Returns32Characters(t *testing.T) {
	uuid, err := GenerateUUIDv7()
	require.NoError(t, err)
	assert.Len(t, uuid, 32)
}

func TestGenerateUUIDv7_NoHyphens(t *testing.T) {
	uuid, err := GenerateUUIDv7()
	require.NoError(t, err)
	assert.NotContains(t, uuid, "-")
}

func TestGenerateUUIDv7_IsLowercase(t *testing.T) {
	uuid, err := GenerateUUIDv7()
	require.NoError(t, err)
	assert.Equal(t, strings.ToLower(uuid), uuid)
}

func TestGenerateUUIDv7_IsHexadecimal(t *testing.T) {
	uuid, err := GenerateUUIDv7()
	require.NoError(t, err)

	// すべての文字が16進数（0-9a-f）であることを確認
	for _, c := range uuid {
		assert.True(t, (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f'), "character %c is not hexadecimal", c)
	}
}

func TestGenerateUUIDv7_Uniqueness(t *testing.T) {
	const count = 100
	uuids := make(map[string]bool)

	for i := 0; i < count; i++ {
		uuid, err := GenerateUUIDv7()
		require.NoError(t, err)
		assert.False(t, uuids[uuid], "UUID %s was generated more than once", uuid)
		uuids[uuid] = true
	}
}
