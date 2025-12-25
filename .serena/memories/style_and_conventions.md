# コードスタイルとコーディング規約

## Go言語

### 命名規則
- パッケージ名: 小文字、短い名前（例: `db`, `config`, `model`）
- 構造体名: パスカルケース（例: `UserRepository`, `ShardConfig`）
- メソッド名: パスカルケース（例: `GetByID`, `CreateUser`）
- 変数名: キャメルケース（例: `userID`, `dbManager`）
- 定数名: パスカルケースまたは大文字スネークケース

### ファイル構成
- テストファイル: `*_test.go`
- GORM版Repository: `*_gorm.go`（例: `user_repository_gorm.go`）

### インポート順序
1. 標準ライブラリ
2. サードパーティライブラリ
3. プロジェクト内パッケージ

### エラーハンドリング
- エラーは明示的に返す
- ラップする場合は `fmt.Errorf("context: %w", err)` を使用

### テスト
- テーブル駆動テストを活用
- testifyを使用（assert, require）
- テスト関数名: `TestStructName_MethodName`

## TypeScript/React

### 命名規則
- コンポーネント名: パスカルケース（例: `UserCard`, `PostList`）
- 関数名: キャメルケース
- ファイル名: パスカルケース.tsx（コンポーネント）

### テスト
- Jest + React Testing Library
- テストファイル: `*.test.tsx`, `*.test.ts`

## ドキュメント
- 日本語で記述
- Markdownを使用

## Sharding規約
- Shard Key: `user_id` を使用
- Hash-based sharding: `shard_id = hash(user_id) % shard_count + 1`
