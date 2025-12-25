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

// =============================================================================
// タスク3.1, 3.2: TableSelector - テーブル選択ロジック
// =============================================================================

// TableSelector はテーブル選択ロジックを提供
type TableSelector struct {
	tableCount  int // 全テーブル数（デフォルト: 32）
	tablesPerDB int // データベースあたりのテーブル数（デフォルト: 8）
}

// NewTableSelector は新しいTableSelectorを作成
func NewTableSelector(tableCount, tablesPerDB int) *TableSelector {
	if tableCount <= 0 {
		tableCount = 32
	}
	if tablesPerDB <= 0 {
		tablesPerDB = 8
	}

	return &TableSelector{
		tableCount:  tableCount,
		tablesPerDB: tablesPerDB,
	}
}

// GetTableNumber はIDからテーブル番号を取得
func (ts *TableSelector) GetTableNumber(id int64) int {
	return int(id % int64(ts.tableCount))
}

// GetTableName はベース名とIDからテーブル名を生成
func (ts *TableSelector) GetTableName(baseName string, id int64) string {
	tableNumber := ts.GetTableNumber(id)
	return fmt.Sprintf("%s_%03d", baseName, tableNumber)
}

// GetDBID はテーブル番号からデータベースIDを取得
func (ts *TableSelector) GetDBID(tableNumber int) int {
	return (tableNumber / ts.tablesPerDB) + 1
}

// GetTableCount は全テーブル数を返す
func (ts *TableSelector) GetTableCount() int {
	return ts.tableCount
}

// =============================================================================
// タスク3.2: テーブル名生成ユーティリティ関数
// =============================================================================

// GetShardingTableName はshardingグループのテーブル名を生成
func GetShardingTableName(baseName string, id int64) string {
	tableNumber := int(id % 32)
	return fmt.Sprintf("%s_%03d", baseName, tableNumber)
}

// GetShardingTableNumber はIDからテーブル番号を取得
func GetShardingTableNumber(id int64) int {
	return int(id % 32)
}

// GetShardingDBID はテーブル番号からデータベースIDを取得
func GetShardingDBID(tableNumber int) int {
	return (tableNumber / 8) + 1
}

// ValidateTableName はテーブル名が有効か検証（SQLインジェクション対策）
func ValidateTableName(tableName string, allowedBaseNames []string) bool {
	for _, baseName := range allowedBaseNames {
		// users_000, users_001, ..., users_031 の形式をチェック
		for i := 0; i < 32; i++ {
			expectedName := fmt.Sprintf("%s_%03d", baseName, i)
			if tableName == expectedName {
				return true
			}
		}
	}
	return false
}
