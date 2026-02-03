---
name: sharding-pattern
description: シャーディング対応のコードを書く際に使用。UUID（string）ベースのシャードキー、GetTableNameFromUUID/GetShardingConnectionByUUID、テーブル分割、クロスシャードクエリ、Repository実装時に自動適用。シャーディング、分散DB、テーブル分割、クロスDBクエリを扱う場合に使用。
allowed-tools: Read, Grep, Glob
---

# Sharding パターン

このプロジェクトのシャーディング実装パターンを定義します。

## シャーディング戦略

- **アルゴリズム**: UUID の後ろ2文字を16進数として解釈し、テーブル数で割った余りでテーブル番号を決定
- **Shard Key**: **UUID（string）が主**。エンティティの ID や user_id は UUID 文字列として扱う
- **計算式**: テーブル番号 = (UUID 後ろ2文字を16進数として解釈) % DBShardingTableCount
- **定数**: `db.DBShardingTableCount`（32）、`db.DBShardingTablesPerDB`（8）

## テーブル分割

- **全テーブル数**: 32（`db.DBShardingTableCount`）
- **DBあたりのテーブル数**: 8（`db.DBShardingTablesPerDB`）
- **テーブル名形式**: `{base_name}_{000-031}` (例: `dm_users_000`, `dm_users_031`, `dm_posts_004`)
- **テーブル名取得**: `GetTableNameFromUUID(baseName, uuid)` → 戻り値 `(string, error)`
- **接続取得**: `GetShardingConnectionByUUID(uuid, tableBaseName)` （例: `"dm_users"`）

## 参照ファイル

シャーディング実装の参照:
- `server/internal/db/sharding.go` - GetTableNameFromUUID, ValidateTableName, TableSelector, 定数（DBShardingTableCount, DBShardingTablesPerDB）
- `server/internal/db/group_manager.go` - GroupManager, GetShardingConnectionByUUID, GetShardingConnection

Repository実装の参照:
- `server/internal/repository/dm_user_repository.go`
- `server/internal/repository/dm_post_repository.go`

## コードパターン

### 1. TableSelector の初期化

定数を使用する。

```go
tableSelector := db.NewTableSelector(db.DBShardingTableCount, db.DBShardingTablesPerDB)
```

### 2. テーブル名の取得（UUID ベース）

戻り値は `(string, error)`。必ずエラーを確認すること。

```go
tableName, err := tableSelector.GetTableNameFromUUID("dm_users", uuid)
if err != nil {
    return nil, fmt.Errorf("failed to get table name: %w", err)
}
// 例: uuid="019b6f83add07d6586044649c19fa5c4" → "dm_users_004"
```

### 3. DB接続の取得（UUID ベース）

```go
conn, err := groupManager.GetShardingConnectionByUUID(uuid, "dm_users")
if err != nil {
    return nil, fmt.Errorf("failed to get sharding connection: %w", err)
}
```

### 4. クロステーブルクエリ（List 操作）

全シャードからデータを取得する場合:

```go
tableCount := r.tableSelector.GetTableCount()
for tableNum := 0; tableNum < tableCount; tableNum++ {
    conn, err := r.groupManager.GetShardingConnection(tableNum)
    if err != nil {
        return nil, fmt.Errorf("failed to get connection for table %d: %w", tableNum, err)
    }
    tableName := fmt.Sprintf("dm_users_%03d", tableNum)
    // クエリ実行...
}
```

### 5. SQLインジェクション対策（動的テーブル名の検証）

動的テーブル名を使用する場合は必ず検証:

```go
if !db.ValidateTableName(tableName, []string{"dm_users", "dm_posts"}) {
    return nil, fmt.Errorf("invalid table name: %s", tableName)
}
```

## マイグレーション

シャーディングDBのマイグレーションは4つのディレクトリに分かれています:
- `db/migrations/sharding_1/`
- `db/migrations/sharding_2/`
- `db/migrations/sharding_3/`
- `db/migrations/sharding_4/`

各DBで同じテーブル構造を持ちますが、テーブル名のサフィックスが異なります。
