# UUIDv7導入要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #57
- **Issueタイトル**: sonyflakeで生成したIDをシャーディングのキーに使うと同じテーブルにデータが振り分けられる
- **Feature名**: 0030-uuidv7
- **作成日**: 2025-01-27

### 1.2 目的
Sonyflakeで生成したIDをシャーディングのキーに使用すると、同一ミリ秒で大量のデータが生成されない限り、IDのあまりでデータを割り振るときに同じテーブルにデータが入ってしまう問題を解決する。
UUIDv7を使用したID生成方式に変更し、シャーディングキーの計算方法を改善することで、データの分散性を向上させる。

### 1.3 スコープ
- SonyflakeによるID生成をUUIDv7によるID生成に置き換える
- 対象カラム: `dm_users.id`, `dm_posts.id`, `dm_posts.user_id`
- IDの仕様: ハイフン抜き小文字32文字UUID
- シャーディングキーの計算方法を変更（UUIDの後ろ2文字から16進数として解釈）
- `github.com/google/uuid`ライブラリを使用する

**本実装の範囲外**:
- 他のテーブル（dm_newsなど）のID生成方式の変更
- データベース構造の変更（カラム型の変更はマイグレーションで対応）

**注意**: 既存データは破棄して良いため、マイグレーション時に既存データを削除する。

## 2. 背景・現状分析

### 2.1 現在の実装
- **ID生成方式**: Sonyflake（`github.com/sony/sonyflake`）を使用
- **IDの型**: `int64`
- **対象カラム**:
  - `dm_users.id`: `int64`型、Sonyflakeで生成
  - `dm_posts.id`: `int64`型、Sonyflakeで生成
  - `dm_posts.user_id`: `int64`型、`dm_users.id`を参照
- **シャーディングキーの計算**:
  - `GetTableNumber(id int64) int`: `id % int64(ts.tableCount)`で計算
  - `GetTableName(baseName string, id int64) string`: テーブル番号からテーブル名を生成
  - `GetShardingConnectionByID(id int64, tableName string)`: IDから接続を取得
- **テーブル分割**: 32テーブルに分割（`dm_users_000` ～ `dm_users_031`, `dm_posts_000` ～ `dm_posts_031`）

### 2.2 課題点
1. **SonyflakeのID構造の問題**: Sonyflake（やSnowflake）のIDは、時間（Timestamp）、マシンID、シーケンス（連番）のビット構成になっており、一般的な規模のシステムでは同一ミリ秒で大量のデータが生成されない限り、IDのあまりでデータを割り振るときに同じテーブルにデータが入ってしまう
2. **データ分散性の低下**: 同一ミリ秒で生成されたIDは、あまりが同じ値になりやすく、データが特定のテーブルに集中する
3. **シャーディングの効果が薄い**: データが均等に分散されず、特定のテーブルに負荷が集中する可能性がある

### 2.3 本実装による改善点
1. **データ分散性の向上**: UUIDv7を使用することで、時間ベースの一意性を保ちながら、より均等なデータ分散を実現
2. **シャーディングキーの改善**: UUIDの後ろ2文字（16進数）を使用することで、より均等な分散を実現
3. **将来の拡張性**: UUIDv7は時間順序性を保ちながら、より分散性の高いID生成が可能

## 3. 機能要件

### 3.1 UUIDv7ライブラリの導入

#### 3.1.1 ライブラリの選択
- **ライブラリ**: `github.com/google/uuid`
- **UUIDバージョン**: UUIDv7（`uuid.NewV7()`）
- **公式ドキュメント**: https://pkg.go.dev/github.com/google/uuid

#### 3.1.2 依存関係の追加
- `go.mod`に`github.com/google/uuid`の依存関係を追加
- `go get github.com/google/uuid`でライブラリを取得

### 3.2 ID生成ユーティリティの実装

#### 3.2.1 UUIDv7生成関数の実装
- **ファイル**: `server/internal/util/idgen/uuidv7.go`（パッケージ名は`idgen`）
- **関数**: `GenerateUUIDv7() (string, error)`
- **仕様**:
  - `uuid.NewV7()`を使用してUUIDv7を生成
  - 生成されたUUIDからハイフンを削除
  - 小文字に変換
  - 32文字の文字列として返す
- **例**: `550e8400-e29b-41d4-a716-446655440000` → `550e8400e29b41d4a716446655440000`

#### 3.2.2 既存のSonyflake関数の削除
- **削除対象**: Sonyflake関数（`GenerateSonyflakeID()`）は、dm_usersとdm_postsでUUIDv7に置き換えられるため、削除する
- **削除ファイル**:
  - `server/internal/util/idgen/sonyflake.go`: Sonyflake関数の実装
  - `server/internal/util/idgen/sonyflake_test.go`: Sonyflake関数のテスト
- **依存関係の削除**: `go.mod`から`github.com/sony/sonyflake`の依存関係を削除
- **注意**: dm_newsではSonyflakeIDを使用していないため、削除しても問題ない

### 3.3 シャーディングキー計算の変更

#### 3.3.1 シャーディングキー計算ロジックの実装
- **ファイル**: `server/internal/db/sharding.go`
- **関数**: `GetTableNumberFromUUID(uuid string) int`
- **仕様**:
  - UUID文字列の後ろ2文字を取得（例: `550e8400e29b41d4a716446655440000` → `00`）
  - 16進数として解釈（`0x00` = 0）
  - テーブル数（32）で割った余りを計算（`0 % 32 = 0`）
  - テーブル番号（0～31）を返す
- **例**:
  - UUID: `550e8400e29b41d4a716446655440000` → 後ろ2文字: `00` → 16進数: 0 → テーブル番号: 0
  - UUID: `550e8400e29b41d4a71644665544000f` → 後ろ2文字: `0f` → 16進数: 15 → テーブル番号: 15
  - UUID: `550e8400e29b41d4a71644665544001f` → 後ろ2文字: `1f` → 16進数: 31 → テーブル番号: 31
  - UUID: `550e8400e29b41d4a716446655440020` → 後ろ2文字: `20` → 16進数: 32 → テーブル番号: 0（32 % 32 = 0）

#### 3.3.2 GetTableName関数の変更
- **既存**: `GetTableName(baseName string, id int64) string`
- **新規**: `GetTableNameFromUUID(baseName string, uuid string) string`
- **仕様**:
  - UUID文字列からテーブル番号を計算
  - テーブル名を生成（`baseName_%03d`形式）
- **後方互換性**: 既存の`GetTableName(baseName string, id int64) string`関数は残す（他のテーブルで使用されている可能性があるため）

#### 3.3.3 GetShardingConnectionByID関数の変更
- **既存**: `GetShardingConnectionByID(id int64, tableName string) (*GORMConnection, error)`
- **新規**: `GetShardingConnectionByUUID(uuid string, tableName string) (*GORMConnection, error)`
- **仕様**:
  - UUID文字列からテーブル番号を計算
  - テーブル番号からデータベースIDを取得
  - 適切な接続を返す
- **既存関数の保持**: 既存の`GetShardingConnectionByID(id int64, tableName string)`関数は残す（現在は使用箇所がないが、将来の使用に備えて保持する）

### 3.4 モデル定義の変更

#### 3.4.1 DmUserモデルの変更
- **ファイル**: `server/internal/model/dm_user.go`
- **変更内容**:
  - `ID int64` → `ID string`
  - JSONタグ: `json:"id,string"` → `json:"id"`
  - GORMタグ: `gorm:"primaryKey"` → `gorm:"primaryKey;type:varchar(32)"`
- **型**: `string`型（32文字のUUID文字列）

#### 3.4.2 DmPostモデルの変更
- **ファイル**: `server/internal/model/dm_post.go`
- **変更内容**:
  - `ID int64` → `ID string`
  - `UserID int64` → `UserID string`
  - JSONタグ: `json:"id,string"` → `json:"id"`, `json:"user_id,string"` → `json:"user_id"`
  - GORMタグ:
    - `ID`: `gorm:"primaryKey"` → `gorm:"primaryKey;type:varchar(32)"`
    - `UserID`: `gorm:"type:bigint;not null;index:idx_dm_posts_user_id"` → `gorm:"type:varchar(32);not null;index:idx_dm_posts_user_id"`
- **型**: `string`型（32文字のUUID文字列）

#### 3.4.3 リクエストモデルの変更
- **ファイル**: `server/internal/model/dm_user.go`, `server/internal/model/dm_post.go`
- **変更内容**:
  - `CreateDmPostRequest.UserID`: `int64` → `string`
  - バリデーション: `validate:"required,gt=0"` → `validate:"required,len=32"`（UUID文字列の長さチェック）

### 3.5 リポジトリ層の変更

#### 3.5.1 DmUserRepositoryの変更
- **ファイル**: `server/internal/repository/dm_user_repository.go`, `server/internal/repository/dm_user_repository_gorm.go`
- **変更内容**:
  - `Create`メソッド: `idgen.GenerateSonyflakeID()` → `idgen.GenerateUUIDv7()`
  - `GetTableName`呼び出し: `GetTableName("dm_users", user.ID)` → `GetTableNameFromUUID("dm_users", user.ID)`
  - `GetShardingConnectionByID`呼び出し: `GetShardingConnectionByID(user.ID, "dm_users")` → `GetShardingConnectionByUUID(user.ID, "dm_users")`
  - すべてのメソッドでIDの型を`int64`から`string`に変更

#### 3.5.2 DmPostRepositoryの変更
- **ファイル**: `server/internal/repository/dm_post_repository.go`, `server/internal/repository/dm_post_repository_gorm.go`
- **変更内容**:
  - `Create`メソッド: `idgen.GenerateSonyflakeID()` → `idgen.GenerateUUIDv7()`
  - `GetTableName`呼び出し: `GetTableName("dm_posts", req.UserID)` → `GetTableNameFromUUID("dm_posts", req.UserID)`
  - `GetShardingConnectionByID`呼び出し: `GetShardingConnectionByID(req.UserID, "dm_posts")` → `GetShardingConnectionByUUID(req.UserID, "dm_posts")`
  - すべてのメソッドでIDの型を`int64`から`string`に変更

### 3.6 サービス層の変更

#### 3.6.1 DmUserServiceの変更
- **ファイル**: `server/internal/service/dm_user_service.go`
- **変更内容**:
  - IDの型を`int64`から`string`に変更
  - リポジトリ層の変更に合わせて型を更新

#### 3.6.2 DmPostServiceの変更
- **ファイル**: `server/internal/service/dm_post_service.go`
- **変更内容**:
  - IDの型を`int64`から`string`に変更
  - リポジトリ層の変更に合わせて型を更新

### 3.7 API層の変更

#### 3.7.1 ハンドラーの変更
- **ファイル**: `server/internal/api/handler/dm_user_handler.go`, `server/internal/api/handler/dm_post_handler.go`
- **変更内容**:
  - IDの型を`int64`から`string`に変更
  - パスパラメータ、クエリパラメータの型を`string`に変更
  - バリデーションルールを更新（UUID文字列の形式チェック）

#### 3.7.2 出力モデルの変更
- **ファイル**: `server/internal/api/huma/outputs.go`
- **変更内容**:
  - `DmUserOutput.ID`: `int64` → `string`
  - `DmPostOutput.ID`: `int64` → `string`
  - `DmPostOutput.UserID`: `int64` → `string`
  - `DmUserPost`構造体のIDフィールドも`string`に変更

### 3.8 サンプルデータ生成コマンドの変更

#### 3.8.1 generate-sample-dataコマンドの変更
- **ファイル**: `server/cmd/generate-sample-data/main.go`
- **変更内容**:
  - `generateDmUsers`関数: `idgen.GenerateSonyflakeID()` → `idgen.GenerateUUIDv7()`
  - `generateDmPosts`関数: `idgen.GenerateSonyflakeID()` → `idgen.GenerateUUIDv7()`
  - IDの型を`int64`から`string`に変更
  - シャーディングキーの計算をUUIDベースに変更

### 3.9 GoAdmin管理画面の変更

#### 3.9.1 dm_user_register.goの変更
- **ファイル**: `server/internal/admin/pages/dm_user_register.go`
- **変更内容**:
  - `insertDmUserSharded`関数: ID生成を`now.UnixNano()`から`idgen.GenerateUUIDv7()`に変更
  - `insertDmUserSharded`関数: 戻り値の型を`int64`から`string`に変更
  - `insertDmUserSharded`関数: テーブル番号の計算をUUIDベースに変更（`GetTableNumberFromUUID`を使用）
  - `handleDmUserRegisterPost`関数: リダイレクトURLのクエリパラメータでIDを`string`型として渡すように変更
  - `checkEmailExistsSharded`関数: 変更不要（メールアドレスの検索のみ）

### 3.10 クライアント側の変更

#### 3.10.1 APIレスポンス型の変更
- **影響範囲**: クライアント側のAPIレスポンス型定義
- **変更内容**:
  - `DmUser`型の`id`フィールド: `number` → `string`
  - `DmPost`型の`id`フィールド: `number` → `string`
  - `DmPost`型の`user_id`フィールド: `number` → `string`
  - APIクライアントコードでIDを`string`型として扱うように変更
  - 必要に応じて、IDのバリデーション（32文字のUUID文字列）を追加

#### 3.10.2 クライアント側の確認
- **確認項目**:
  - APIレスポンスの型定義が更新されていること
  - APIクライアントコードが`string`型のIDを正しく処理できること
  - IDの表示・入力フォームが`string`型に対応していること

### 3.11 Sonyflake関数の削除

#### 3.11.1 Sonyflake関数の削除
- **ファイル**: 
  - `server/internal/util/idgen/sonyflake.go`: 削除
  - `server/internal/util/idgen/sonyflake_test.go`: 削除
- **依存関係の削除**: 
  - `go.mod`から`github.com/sony/sonyflake`の依存関係を削除
  - `go mod tidy`を実行して依存関係を整理
- **理由**: dm_usersとdm_postsでUUIDv7に置き換えられるため、Sonyflake関数は不要になる
- **注意**: dm_newsではSonyflakeIDを使用していないため、削除しても問題ない

### 3.12 マイグレーションの作成

#### 3.12.1 データベーススキーマの変更
- **ファイル**: Atlasマイグレーションファイル（新規作成）
- **変更内容**:
  - `dm_users.id`: `bigint unsigned` → `varchar(32)`
  - `dm_posts.id`: `bigint unsigned` → `varchar(32)`
  - `dm_posts.user_id`: `bigint` → `varchar(32)`
- **既存データの扱い**: 既存データは破棄する（テーブルを再作成するか、既存データを削除する）

## 4. 非機能要件

### 4.1 パフォーマンス
- UUIDv7の生成速度はSonyflakeと同等またはそれ以上であること
- シャーディングキーの計算は高速であること（文字列操作は最小限に）

### 4.2 互換性
- 既存の`GetTableName(baseName string, id int64)`関数は残す（後方互換性のため）
- 既存の`GetShardingConnectionByID(id int64, tableName string)`関数は残す（現在は使用箇所がないが、将来の使用に備えて保持する）
- Sonyflake関数は削除する（dm_usersとdm_postsでUUIDv7に置き換えられるため）

### 4.3 データ整合性
- UUIDv7で生成されたIDは一意であること
- シャーディングキーの計算は正確であること（0～31の範囲内）
- マイグレーション後は新規データのみが存在すること（既存データは削除される）

### 4.4 エラーハンドリング
- UUID生成エラーは適切に処理されること
- 無効なUUID文字列が渡された場合はエラーを返すこと
- シャーディングキー計算時のエラーは適切に処理されること

### 4.5 テスト
- UUIDv7生成関数の単体テストを作成
- シャーディングキー計算関数の単体テストを作成
- リポジトリ層の統合テストを更新
- 既存のテストが全て通過すること

## 5. 制約事項

### 5.1 既存システムとの関係
- **Sonyflakeライブラリ**: dm_usersとdm_postsでUUIDv7に置き換えられるため、削除する
- **既存の関数**: 
  - `GetTableName(baseName string, id int64)`: 残す（後方互換性のため）
  - `GetShardingConnectionByID(id int64, tableName string)`: 残す（現在は使用箇所がないが、将来の使用に備えて保持する）
- **既存データ**: 既存データは破棄する（マイグレーション時に削除）

### 5.2 データベース制約
- **カラム型の変更**: `bigint`から`varchar(32)`への変更はマイグレーションで対応
- **既存データ**: 既存データは破棄する（マイグレーション時に削除）
- **インデックス**: 型変更に伴いインデックスを再作成する

### 5.3 技術スタック
- **UUIDライブラリ**: `github.com/google/uuid`を使用
- **Goバージョン**: 既存のGoバージョン（1.23.4）を維持
- **データベース**: SQLite（開発環境）

### 5.4 文字列長の制約
- **UUID文字列**: 32文字（ハイフン抜き）
- **データベースカラム**: `varchar(32)`で定義
- **バリデーション**: UUID文字列の長さは32文字であること

### 5.5 シャーディングキーの制約
- **テーブル数**: 32テーブル（0～31）
- **計算方法**: UUIDの後ろ2文字を16進数として解釈し、32で割った余りを使用
- **範囲**: 0～31の範囲内であること

## 6. 受け入れ基準

### 6.1 UUIDv7ライブラリの導入
- [ ] `github.com/google/uuid`ライブラリが`go.mod`に追加されている
- [ ] `go get github.com/google/uuid`でライブラリが取得できる
- [ ] `uuid.NewV7()`関数が使用できる

### 6.2 ID生成ユーティリティの実装
- [ ] `server/internal/util/idgen/uuidv7.go`が作成されている
- [ ] `GenerateUUIDv7() (string, error)`関数が実装されている
- [ ] 生成されたUUIDからハイフンが削除されている
- [ ] 生成されたUUIDが小文字に変換されている
- [ ] 生成されたUUIDが32文字であること
- [ ] 単体テストが作成されている
- [ ] `server/internal/util/idgen/sonyflake.go`が削除されている
- [ ] `server/internal/util/idgen/sonyflake_test.go`が削除されている
- [ ] `go.mod`から`github.com/sony/sonyflake`の依存関係が削除されている

### 6.3 シャーディングキー計算の変更
- [ ] `GetTableNumberFromUUID(uuid string) int`関数が実装されている
- [ ] UUIDの後ろ2文字を正しく取得できること
- [ ] 16進数として正しく解釈できること
- [ ] テーブル数（32）で割った余りが正しく計算できること
- [ ] テーブル番号が0～31の範囲内であること
- [ ] `GetTableNameFromUUID(baseName string, uuid string) string`関数が実装されている
- [ ] `GetShardingConnectionByUUID(uuid string, tableName string) (*GORMConnection, error)`関数が実装されている
- [ ] 単体テストが作成されている

### 6.4 モデル定義の変更
- [ ] `DmUser.ID`が`string`型に変更されている
- [ ] `DmPost.ID`が`string`型に変更されている
- [ ] `DmPost.UserID`が`string`型に変更されている
- [ ] GORMタグが適切に設定されている（`type:varchar(32)`）
- [ ] JSONタグが適切に設定されている
- [ ] リクエストモデルの型が適切に変更されている
- [ ] バリデーションルールが適切に設定されている

### 6.5 リポジトリ層の変更
- [ ] `DmUserRepository.Create`メソッドが`GenerateUUIDv7()`を使用している
- [ ] `DmPostRepository.Create`メソッドが`GenerateUUIDv7()`を使用している
- [ ] すべてのメソッドでIDの型が`string`に変更されている
- [ ] `GetTableNameFromUUID`が使用されている
- [ ] `GetShardingConnectionByUUID`が使用されている
- [ ] 統合テストが更新されている

### 6.6 サービス層の変更
- [ ] `DmUserService`のIDの型が`string`に変更されている
- [ ] `DmPostService`のIDの型が`string`に変更されている
- [ ] すべてのメソッドで型が適切に更新されている

### 6.7 API層の変更
- [ ] ハンドラーのIDの型が`string`に変更されている
- [ ] パスパラメータ、クエリパラメータの型が`string`に変更されている
- [ ] バリデーションルールが適切に設定されている
- [ ] 出力モデルのIDの型が`string`に変更されている

### 6.8 サンプルデータ生成コマンドの変更
- [ ] `generateDmUsers`関数が`GenerateUUIDv7()`を使用している
- [ ] `generateDmPosts`関数が`GenerateUUIDv7()`を使用している
- [ ] IDの型が`string`に変更されている
- [ ] シャーディングキーの計算がUUIDベースに変更されている

### 6.9 GoAdmin管理画面の変更
- [ ] `insertDmUserSharded`関数が`GenerateUUIDv7()`を使用している
- [ ] `insertDmUserSharded`関数の戻り値の型が`string`に変更されている
- [ ] テーブル番号の計算がUUIDベースに変更されている（`GetTableNumberFromUUID`を使用）
- [ ] リダイレクトURLのクエリパラメータでIDが`string`型として渡されている

### 6.10 クライアント側の変更
- [ ] APIレスポンス型定義が更新されている（IDが`string`型）
- [ ] APIクライアントコードが`string`型のIDを正しく処理できること
- [ ] IDの表示・入力フォームが`string`型に対応していること

### 6.11 Sonyflake関数の削除
- [ ] `server/internal/util/idgen/sonyflake.go`が削除されている
- [ ] `server/internal/util/idgen/sonyflake_test.go`が削除されている
- [ ] `go.mod`から`github.com/sony/sonyflake`の依存関係が削除されている
- [ ] `go mod tidy`が実行されている
- [ ] コード内でSonyflake関数への参照が全て削除されている

### 6.12 マイグレーションの作成
- [ ] Atlasマイグレーションファイルが作成されている
- [ ] `dm_users.id`が`varchar(32)`に変更されている
- [ ] `dm_posts.id`が`varchar(32)`に変更されている
- [ ] `dm_posts.user_id`が`varchar(32)`に変更されている
- [ ] 既存データが削除されていること（マイグレーション後に既存データが存在しないこと）
- [ ] マイグレーションが正常に実行できること

### 6.13 テスト
- [ ] UUIDv7生成関数の単体テストが作成されている
- [ ] シャーディングキー計算関数の単体テストが作成されている
- [ ] リポジトリ層の統合テストが更新されている
- [ ] 既存のテストが全て通過すること
- [ ] 新規作成したテストが全て通過すること

### 6.14 ドキュメント
- [ ] `docs/Architecture.md`が更新されている（UUIDv7の使用を記載、Sonyflakeの削除を記載）
- [ ] `server/internal/db/sharding.go`のコメントが更新されている（UUIDv7の使用を記載、Sonyflakeの記述を削除）
- [ ] コード内のコメントが適切に更新されている

## 7. 影響範囲

### 7.1 新規追加が必要なディレクトリ・ファイル

#### ファイル
- `server/internal/util/idgen/uuidv7.go`: UUIDv7生成関数
- `server/internal/util/idgen/uuidv7_test.go`: UUIDv7生成関数のテスト
- Atlasマイグレーションファイル（新規作成）

### 7.2 変更が必要なファイル

#### モデル定義
- `server/internal/model/dm_user.go`: IDの型を`string`に変更
- `server/internal/model/dm_post.go`: ID、UserIDの型を`string`に変更

#### リポジトリ層
- `server/internal/repository/dm_user_repository.go`: ID生成とシャーディングキー計算を変更
- `server/internal/repository/dm_user_repository_gorm.go`: ID生成とシャーディングキー計算を変更
- `server/internal/repository/dm_post_repository.go`: ID生成とシャーディングキー計算を変更
- `server/internal/repository/dm_post_repository_gorm.go`: ID生成とシャーディングキー計算を変更

#### サービス層
- `server/internal/service/dm_user_service.go`: IDの型を`string`に変更
- `server/internal/service/dm_post_service.go`: IDの型を`string`に変更

#### API層
- `server/internal/api/handler/dm_user_handler.go`: IDの型を`string`に変更
- `server/internal/api/handler/dm_post_handler.go`: IDの型を`string`に変更
- `server/internal/api/huma/outputs.go`: 出力モデルのIDの型を`string`に変更

#### シャーディング関連
- `server/internal/db/sharding.go`: シャーディングキー計算関数を追加

#### サンプルデータ生成
- `server/cmd/generate-sample-data/main.go`: ID生成をUUIDv7に変更

#### GoAdmin管理画面
- `server/internal/admin/pages/dm_user_register.go`: ID生成をUUIDv7に変更、IDの型を`string`に変更

#### クライアント側
- クライアント側のAPIレスポンス型定義: IDの型を`string`に変更
- クライアント側のAPIクライアントコード: IDの型を`string`として処理

#### テストファイル
- `server/internal/util/idgen/uuidv7_test.go`: 新規作成
- `server/internal/db/sharding_test.go`: シャーディングキー計算のテストを追加
- `server/internal/repository/dm_user_repository_test.go`: テストを更新
- `server/internal/repository/dm_post_repository_test.go`: テストを更新

#### ドキュメント
- `docs/Architecture.md`: UUIDv7の使用を記載

### 7.3 削除が必要なファイル
- `server/internal/util/idgen/sonyflake.go`: 削除する（UUIDv7に置き換えられるため）
- `server/internal/util/idgen/sonyflake_test.go`: 削除する（UUIDv7に置き換えられるため）

### 7.4 既存ファイルの扱い
- `server/internal/db/sharding.go`の既存関数: 
  - `GetTableName(baseName string, id int64)`: 削除しない（後方互換性のため）
  - `GetShardingConnectionByID(id int64, tableName string)`: 削除しない（現在は使用箇所がないが、将来の使用に備えて保持する）

## 8. 実装上の注意事項

### 8.1 UUIDv7生成の実装
- `uuid.NewV7()`を使用してUUIDv7を生成
- 生成されたUUIDからハイフンを削除（`strings.ReplaceAll(uuid.String(), "-", "")`）
- 小文字に変換（`strings.ToLower()`）
- 32文字の文字列として返す

### 8.2 シャーディングキー計算の実装
- UUID文字列の後ろ2文字を取得（`uuid[len(uuid)-2:]`）
- `strconv.ParseInt(uuidSuffix, 16, 64)`で16進数として解釈
- テーブル数（32）で割った余りを計算（`value % 32`）
- テーブル番号（0～31）を返す

### 8.3 エラーハンドリング
- UUID生成エラーは適切に処理し、エラーメッセージを返す
- 無効なUUID文字列が渡された場合はエラーを返す（長さチェック、16進数チェック）
- シャーディングキー計算時のエラーは適切に処理する

### 8.4 型変換の注意
- `int64`から`string`への変更は、すべての関連箇所で一貫して行う
- JSONシリアライゼーション時の型変換に注意（`json:"id,string"`タグは不要になる）
- データベースへの保存時の型変換に注意（`varchar(32)`）

### 8.5 テストの実装
- UUIDv7生成関数のテスト: 生成されたUUIDが32文字であること、ハイフンが含まれていないこと、小文字であることを確認
- シャーディングキー計算のテスト: 様々なUUID文字列でテーブル番号が正しく計算されることを確認
- リポジトリ層のテスト: ID生成とシャーディングキー計算が正しく動作することを確認

### 8.6 GoAdmin管理画面の実装
- `insertDmUserSharded`関数で`idgen.GenerateUUIDv7()`を使用してIDを生成
- 生成されたUUIDからテーブル番号を計算（`GetTableNumberFromUUID`を使用）
- リダイレクトURLのクエリパラメータでIDを`string`型として渡す
- 既存の`checkEmailExistsSharded`関数は変更不要（メールアドレスの検索のみ）

### 8.7 クライアント側の実装
- APIレスポンス型定義を更新（IDを`string`型に変更）
- APIクライアントコードでIDを`string`型として処理
- IDの表示・入力フォームが`string`型に対応していることを確認
- 必要に応じて、IDのバリデーション（32文字のUUID文字列）を追加

### 8.8 Sonyflake関数の削除
- `server/internal/util/idgen/sonyflake.go`を削除
- `server/internal/util/idgen/sonyflake_test.go`を削除
- `go.mod`から`github.com/sony/sonyflake`の依存関係を削除
- `go mod tidy`を実行して依存関係を整理
- コード内でSonyflake関数への参照が全て削除されていることを確認
- `docs/Architecture.md`からSonyflakeの記述を削除

### 8.9 マイグレーションの実装
- Atlasマイグレーションファイルを作成
- カラム型を`bigint`から`varchar(32)`に変更
- 既存データを削除する（テーブルを再作成するか、既存データを削除する）
- マイグレーション実行前にバックアップを取得することを推奨（必要に応じて）

### 8.10 後方互換性の維持
- 既存の`GetTableName(baseName string, id int64)`関数は残す（後方互換性のため）
- 既存の`GetShardingConnectionByID(id int64, tableName string)`関数は残す（現在は使用箇所がないが、将来の使用に備えて保持する）
- 新しい関数（`GetTableNameFromUUID`, `GetShardingConnectionByUUID`）を追加
- Sonyflake関数は削除する（dm_usersとdm_postsでUUIDv7に置き換えられるため）

### 8.11 ドキュメントの更新
- `docs/Architecture.md`にUUIDv7の使用を記載し、Sonyflakeの削除を記載
- `server/internal/db/sharding.go`のコメントを更新（UUIDv7の使用を記載、Sonyflakeの記述を削除）
- コード内のコメントを適切に更新

## 9. 参考情報

### 9.1 関連Issue
- GitHub Issue #57: sonyflakeで生成したIDをシャーディングのキーに使うと同じテーブルにデータが振り分けられる

### 9.2 既存ドキュメント
- `docs/Architecture.md`: アーキテクチャドキュメント
- `server/internal/db/sharding.go`: シャーディング規則のドキュメント
- `.kiro/specs/0028-dmtable-define/requirements.md`: Sonyflake導入の要件定義書

### 9.3 技術スタック
- **UUIDライブラリ**: `github.com/google/uuid`
- **UUIDバージョン**: UUIDv7
- **Goバージョン**: 1.23.4
- **データベース**: SQLite（開発環境）

### 9.4 参考リンク
- UUIDv7仕様: https://www.ietf.org/rfc/rfc4122.txt
- github.com/google/uuid: https://pkg.go.dev/github.com/google/uuid
- UUIDv7の説明: https://www.ietf.org/archive/id/draft-ietf-uuidrev-rfc4122bis-00.html
