package db

import (
	"fmt"
	"hash/fnv"
	"strconv"
)

// =============================================================================
// シャーディング規則（Sharding Rules）
// =============================================================================
//
// このプロジェクトでは、分散テーブル環境で以下のシャーディング規則を使用します。
//
// ## テーブルシャーディングキー（Table Sharding Key）
//
// | テーブル群        | シャーディングキー | 計算式                                    |
// |-------------------|--------------------|-----------------------------------------|
// | dm_users_NNN      | id                 | UUID後ろ2文字(16進数) % DBShardingTableCount |
// | dm_posts_NNN      | user_id            | UUID後ろ2文字(16進数) % DBShardingTableCount |
//
// ## ID生成
//
// IDはUUIDv7（github.com/google/uuid）を使用して生成されます。
// - フォーマット: 32文字の16進数文字列（ハイフンなし小文字）
// - 例: "019b6f83add07d6586044649c19fa5c4"
// - 時系列順序が保証され、分散環境でも一意なIDが生成されます。
//
// ## シャーディングキーの計算方法
//
// UUIDからテーブル番号を計算するには、UUIDの後ろ2文字を16進数として解釈し、
// テーブル数（DBShardingTableCount = 32）で割った余りを使用します。
//
// 例: UUID = "019b6f83add07d6586044649c19fa5c4"
//   - 後ろ2文字: "c4"
//   - 16進数として解釈: 0xc4 = 196
//   - テーブル番号: 196 % 32 = 4
//   - テーブル名: dm_users_004
//
// =============================================================================

// シャーディング関連の定数
const (
	// DBShardingTableCount はshardingグループのテーブル総数
	DBShardingTableCount = 32
	// DBShardingTablesPerDB はデータベースあたりのテーブル数
	DBShardingTablesPerDB = 8
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
	tableCount  int // 全テーブル数（デフォルト: DBShardingTableCount）
	tablesPerDB int // データベースあたりのテーブル数（デフォルト: DBShardingTablesPerDB）
}

// NewTableSelector は新しいTableSelectorを作成
func NewTableSelector(tableCount, tablesPerDB int) *TableSelector {
	if tableCount <= 0 {
		tableCount = DBShardingTableCount
	}
	if tablesPerDB <= 0 {
		tablesPerDB = DBShardingTablesPerDB
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
	tableNumber := int(id % DBShardingTableCount)
	return fmt.Sprintf("%s_%03d", baseName, tableNumber)
}

// GetShardingTableNumber はIDからテーブル番号を取得
func GetShardingTableNumber(id int64) int {
	return int(id % DBShardingTableCount)
}

// GetShardingDBID はテーブル番号からデータベースIDを取得
func GetShardingDBID(tableNumber int) int {
	return (tableNumber / DBShardingTablesPerDB) + 1
}

// ValidateTableName はテーブル名が有効か検証（SQLインジェクション対策）
// allowedBaseNamesには "dm_users", "dm_posts" などのベース名を指定する
func ValidateTableName(tableName string, allowedBaseNames []string) bool {
	for _, baseName := range allowedBaseNames {
		// dm_users_000, dm_users_001, ..., dm_users_031 の形式をチェック
		for i := 0; i < DBShardingTableCount; i++ {
			expectedName := fmt.Sprintf("%s_%03d", baseName, i)
			if tableName == expectedName {
				return true
			}
		}
	}
	return false
}

// =============================================================================
// UUIDv7ベースのシャーディングキー計算関数
// =============================================================================

// GetTableNumberFromUUID はUUID文字列からテーブル番号を取得
// UUIDの後ろ2文字を16進数として解釈し、テーブル数で割った余りを返す
func (ts *TableSelector) GetTableNumberFromUUID(uuid string) (int, error) {
	// UUID文字列の長さをチェック
	if len(uuid) < 2 {
		return 0, fmt.Errorf("invalid UUID string: length must be at least 2, got %d", len(uuid))
	}

	// 後ろ2文字を取得
	suffix := uuid[len(uuid)-2:]

	// 16進数として解釈（大文字小文字両対応）
	value, err := strconv.ParseInt(suffix, 16, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse UUID suffix as hex: %w", err)
	}

	// テーブル数で割った余りを計算
	tableNumber := int(value % int64(ts.tableCount))

	return tableNumber, nil
}

// GetTableNameFromUUID はベース名とUUIDからテーブル名を生成
func (ts *TableSelector) GetTableNameFromUUID(baseName string, uuid string) (string, error) {
	tableNumber, err := ts.GetTableNumberFromUUID(uuid)
	if err != nil {
		return "", fmt.Errorf("failed to get table number from UUID: %w", err)
	}
	return fmt.Sprintf("%s_%03d", baseName, tableNumber), nil
}
