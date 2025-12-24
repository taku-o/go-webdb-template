package db

import (
	"fmt"
	"hash/fnv"
)

// ShardingStrategy はSharding戦略のインターフェース
type ShardingStrategy interface {
	GetShardID(key int64) int
}

// HashBasedSharding はHash-basedのSharding戦略
type HashBasedSharding struct {
	shardCount int
}

// NewHashBasedSharding は新しいHash-basedのSharding戦略を作成
func NewHashBasedSharding(shardCount int) *HashBasedSharding {
	if shardCount <= 0 {
		shardCount = 1
	}
	return &HashBasedSharding{
		shardCount: shardCount,
	}
}

// GetShardID はキーに基づいてShard IDを返す
// user_id をキーとして使用し、ハッシュ値からShard IDを決定
func (h *HashBasedSharding) GetShardID(key int64) int {
	hash := fnv.New32a()
	hash.Write([]byte(fmt.Sprintf("%d", key)))
	hashValue := hash.Sum32()

	// Shard IDは1から始まるので、1を加算
	shardID := int(hashValue%uint32(h.shardCount)) + 1
	return shardID
}

// GetShardCount はShard数を返す
func (h *HashBasedSharding) GetShardCount() int {
	return h.shardCount
}
