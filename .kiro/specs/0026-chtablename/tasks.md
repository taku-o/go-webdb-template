# テーブル名変更機能実装タスク一覧

## 概要
テーブル名を`users`→`dm_users`、`posts`→`dm_posts`、`news`→`dm_news`に変更する機能の実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: スキーマ定義ファイルの変更

#### - [x] タスク 1.1: マスターデータベーススキーマの変更
**目的**: `db/schema/master.hcl`でnewsテーブルをdm_newsに変更

**作業内容**:
- `db/schema/master.hcl`を開く
- `table "news"`を`table "dm_news"`に変更
- インデックス名を変更:
  - `idx_news_published_at` → `idx_dm_news_published_at`
  - `idx_news_author_id` → `idx_dm_news_author_id`
- HCLファイルの構文チェック

**受け入れ基準**:
- `table "dm_news"`が定義されている
- インデックス名が正しく変更されている
- HCLファイルの構文エラーがない

---

#### - [x] タスク 1.2: シャーディングデータベース1のスキーマ変更（users）
**目的**: `db/schema/sharding_1/users.hcl`でusersテーブルをdm_usersに変更

**作業内容**:
- `db/schema/sharding_1/users.hcl`を開く
- `table "users_000"`～`table "users_007"`を`table "dm_users_000"`～`table "dm_users_007"`に変更
- インデックス名を変更:
  - `idx_users_000_email` → `idx_dm_users_000_email`
  - （users_001～users_007も同様）
- ファイル名を変更: `users.hcl` → `dm_users.hcl`
- HCLファイルの構文チェック

**受け入れ基準**:
- `table "dm_users_000"`～`table "dm_users_007"`が定義されている
- インデックス名が正しく変更されている
- ファイル名が`dm_users.hcl`に変更されている
- HCLファイルの構文エラーがない

---

#### - [x] タスク 1.3: シャーディングデータベース1のスキーマ変更（posts）
**目的**: `db/schema/sharding_1/posts.hcl`でpostsテーブルをdm_postsに変更

**作業内容**:
- `db/schema/sharding_1/posts.hcl`を開く
- `table "posts_000"`～`table "posts_007"`を`table "dm_posts_000"`～`table "dm_posts_007"`に変更
- インデックス名を変更:
  - `idx_posts_000_user_id` → `idx_dm_posts_000_user_id`
  - `idx_posts_000_created_at` → `idx_dm_posts_000_created_at`
  - （posts_001～posts_007も同様）
- ファイル名を変更: `posts.hcl` → `dm_posts.hcl`
- HCLファイルの構文チェック

**受け入れ基準**:
- `table "dm_posts_000"`～`table "dm_posts_007"`が定義されている
- インデックス名が正しく変更されている
- ファイル名が`dm_posts.hcl`に変更されている
- HCLファイルの構文エラーがない

---

#### - [x] タスク 1.4: シャーディングデータベース2のスキーマ変更（users）
**目的**: `db/schema/sharding_2/users.hcl`でusersテーブルをdm_usersに変更

**作業内容**:
- `db/schema/sharding_2/users.hcl`を開く
- `table "users_008"`～`table "users_015"`を`table "dm_users_008"`～`table "dm_users_015"`に変更
- インデックス名を変更（users_008～users_015）
- ファイル名を変更: `users.hcl` → `dm_users.hcl`
- HCLファイルの構文チェック

**受け入れ基準**:
- `table "dm_users_008"`～`table "dm_users_015"`が定義されている
- インデックス名が正しく変更されている
- ファイル名が`dm_users.hcl`に変更されている
- HCLファイルの構文エラーがない

---

#### - [x] タスク 1.5: シャーディングデータベース2のスキーマ変更（posts）
**目的**: `db/schema/sharding_2/posts.hcl`でpostsテーブルをdm_postsに変更

**作業内容**:
- `db/schema/sharding_2/posts.hcl`を開く
- `table "posts_008"`～`table "posts_015"`を`table "dm_posts_008"`～`table "dm_posts_015"`に変更
- インデックス名を変更（posts_008～posts_015）
- ファイル名を変更: `posts.hcl` → `dm_posts.hcl`
- HCLファイルの構文チェック

**受け入れ基準**:
- `table "dm_posts_008"`～`table "dm_posts_015"`が定義されている
- インデックス名が正しく変更されている
- ファイル名が`dm_posts.hcl`に変更されている
- HCLファイルの構文エラーがない

---

#### - [x] タスク 1.6: シャーディングデータベース3のスキーマ変更（users）
**目的**: `db/schema/sharding_3/users.hcl`でusersテーブルをdm_usersに変更

**作業内容**:
- `db/schema/sharding_3/users.hcl`を開く
- `table "users_016"`～`table "users_023"`を`table "dm_users_016"`～`table "dm_users_023"`に変更
- インデックス名を変更（users_016～users_023）
- ファイル名を変更: `users.hcl` → `dm_users.hcl`
- HCLファイルの構文チェック

**受け入れ基準**:
- `table "dm_users_016"`～`table "dm_users_023"`が定義されている
- インデックス名が正しく変更されている
- ファイル名が`dm_users.hcl`に変更されている
- HCLファイルの構文エラーがない

---

#### - [x] タスク 1.7: シャーディングデータベース3のスキーマ変更（posts）
**目的**: `db/schema/sharding_3/posts.hcl`でpostsテーブルをdm_postsに変更

**作業内容**:
- `db/schema/sharding_3/posts.hcl`を開く
- `table "posts_016"`～`table "posts_023"`を`table "dm_posts_016"`～`table "dm_posts_023"`に変更
- インデックス名を変更（posts_016～posts_023）
- ファイル名を変更: `posts.hcl` → `dm_posts.hcl`
- HCLファイルの構文チェック

**受け入れ基準**:
- `table "dm_posts_016"`～`table "dm_posts_023"`が定義されている
- インデックス名が正しく変更されている
- ファイル名が`dm_posts.hcl`に変更されている
- HCLファイルの構文エラーがない

---

#### - [x] タスク 1.8: シャーディングデータベース4のスキーマ変更（users）
**目的**: `db/schema/sharding_4/users.hcl`でusersテーブルをdm_usersに変更

**作業内容**:
- `db/schema/sharding_4/users.hcl`を開く
- `table "users_024"`～`table "users_031"`を`table "dm_users_024"`～`table "dm_users_031"`に変更
- インデックス名を変更（users_024～users_031）
- ファイル名を変更: `users.hcl` → `dm_users.hcl`
- HCLファイルの構文チェック

**受け入れ基準**:
- `table "dm_users_024"`～`table "dm_users_031"`が定義されている
- インデックス名が正しく変更されている
- ファイル名が`dm_users.hcl`に変更されている
- HCLファイルの構文エラーがない

---

#### - [x] タスク 1.9: シャーディングデータベース4のスキーマ変更（posts）
**目的**: `db/schema/sharding_4/posts.hcl`でpostsテーブルをdm_postsに変更

**作業内容**:
- `db/schema/sharding_4/posts.hcl`を開く
- `table "posts_024"`～`table "posts_031"`を`table "dm_posts_024"`～`table "dm_posts_031"`に変更
- インデックス名を変更（posts_024～posts_031）
- ファイル名を変更: `posts.hcl` → `dm_posts.hcl`
- HCLファイルの構文チェック

**受け入れ基準**:
- `table "dm_posts_024"`～`table "dm_posts_031"`が定義されている
- インデックス名が正しく変更されている
- ファイル名が`dm_posts.hcl`に変更されている
- HCLファイルの構文エラーがない

---

#### - [x] タスク 1.10: Atlas設定ファイルの確認
**目的**: Atlas設定ファイルでスキーマファイルの参照パスを確認・更新

**作業内容**:
- `config/*/atlas.hcl`ファイルを確認
- スキーマファイルの参照パスが正しいか確認
- 必要に応じて更新（ファイル名変更に対応）

**受け入れ基準**:
- Atlas設定ファイルでスキーマファイルが正しく参照されている
- ファイル名変更に対応している

---

### Phase 2: データベース層の変更

#### - [x] タスク 2.1: ValidateTableName関数の許可リスト更新
**目的**: `server/internal/db/sharding.go`の`ValidateTableName`関数で許可リストを更新

**作業内容**:
- `server/internal/db/sharding.go`を開く
- `ValidateTableName`関数を確認
- 許可リストに`"dm_users"`, `"dm_posts"`を追加
- コメントやドキュメントを更新

**受け入れ基準**:
- `ValidateTableName("dm_users_005", []string{"dm_users", "dm_posts"})`が`true`を返す
- `ValidateTableName("dm_posts_010", []string{"dm_users", "dm_posts"})`が`true`を返す
- ビルドエラーがない

---

#### - [x] タスク 2.2: データベース層のコメント・ドキュメント更新
**目的**: データベース層のコメントやドキュメントを更新

**作業内容**:
- `server/internal/db/sharding.go`のコメントを確認
- テーブル名に関するコメントを更新（`"users"` → `"dm_users"`など）
- 関数のドキュメントコメントを更新

**受け入れ基準**:
- コメントが正しく更新されている
- ドキュメントが最新の状態になっている

---

### Phase 3: モデル層の変更

#### - [x] タスク 3.1: UserモデルのTableNameメソッド変更
**目的**: `server/internal/model/user.go`の`TableName()`メソッドを変更

**作業内容**:
- `server/internal/model/user.go`を開く
- `TableName()`メソッドの戻り値を`"users"`から`"dm_users"`に変更
- ファイル名を変更: `user.go` → `dm_user.go`

**受け入れ基準**:
- `User.TableName()`が`"dm_users"`を返す
- ファイル名が`dm_user.go`に変更されている
- ビルドエラーがない

---

#### - [x] タスク 3.2: PostモデルのTableNameメソッド変更
**目的**: `server/internal/model/post.go`の`TableName()`メソッドを変更

**作業内容**:
- `server/internal/model/post.go`を開く
- `TableName()`メソッドの戻り値を`"posts"`から`"dm_posts"`に変更
- ファイル名を変更: `post.go` → `dm_post.go`

**受け入れ基準**:
- `Post.TableName()`が`"dm_posts"`を返す
- ファイル名が`dm_post.go`に変更されている
- ビルドエラーがない

---

#### - [x] タスク 3.3: NewsモデルのTableNameメソッド変更
**目的**: `server/internal/model/news.go`の`TableName()`メソッドを変更

**作業内容**:
- `server/internal/model/news.go`を開く
- `TableName()`メソッドの戻り値を`"news"`から`"dm_news"`に変更
- ファイル名を変更: `news.go` → `dm_news.go`

**受け入れ基準**:
- `News.TableName()`が`"dm_news"`を返す
- ファイル名が`dm_news.go`に変更されている
- ビルドエラーがない

---

#### - [x] タスク 3.4: モデル層のビルド確認
**目的**: モデル層の変更後、ビルドエラーがないことを確認

**作業内容**:
- `go build ./server/internal/model/...`を実行
- ビルドエラーがないことを確認
- `go mod tidy`を実行して依存関係を更新

**受け入れ基準**:
- ビルドエラーがない
- `go mod tidy`が正常に実行される

---

### Phase 4: Repository層の変更

#### - [x] タスク 4.1: UserRepositoryのテーブル名参照変更
**目的**: `server/internal/repository/user_repository.go`のテーブル名参照を変更

**作業内容**:
- `server/internal/repository/user_repository.go`を開く
- `GetTableName("users", ...)`を`GetTableName("dm_users", ...)`に変更
- 文字列リテラル`"users"`を`"dm_users"`に変更
- `fmt.Sprintf("users_%03d", ...)`を`fmt.Sprintf("dm_users_%03d", ...)`に変更
- `GetShardingConnectionByID(..., "users")`を`GetShardingConnectionByID(..., "dm_users")`に変更
- ファイル名を変更: `user_repository.go` → `dm_user_repository.go`

**受け入れ基準**:
- 全てのテーブル名参照が`"dm_users"`に変更されている
- ファイル名が`dm_user_repository.go`に変更されている
- ビルドエラーがない

---

#### - [x] タスク 4.2: UserRepositoryGORMのテーブル名参照変更
**目的**: `server/internal/repository/user_repository_gorm.go`のテーブル名参照を変更

**作業内容**:
- `server/internal/repository/user_repository_gorm.go`を開く
- `GetTableName("users", ...)`を`GetTableName("dm_users", ...)`に変更
- `GetShardingConnectionByID(..., "users")`を`GetShardingConnectionByID(..., "dm_users")`に変更
- ファイル名を変更: `user_repository_gorm.go` → `dm_user_repository_gorm.go`

**受け入れ基準**:
- 全てのテーブル名参照が`"dm_users"`に変更されている
- ファイル名が`dm_user_repository_gorm.go`に変更されている
- ビルドエラーがない

---

#### - [x] タスク 4.3: PostRepositoryのテーブル名参照変更
**目的**: `server/internal/repository/post_repository.go`のテーブル名参照を変更

**作業内容**:
- `server/internal/repository/post_repository.go`を開く
- `GetTableName("posts", ...)`を`GetTableName("dm_posts", ...)`に変更
- 文字列リテラル`"posts"`を`"dm_posts"`に変更
- `fmt.Sprintf("posts_%03d", ...)`を`fmt.Sprintf("dm_posts_%03d", ...)`に変更
- `GetShardingConnectionByID(..., "posts")`を`GetShardingConnectionByID(..., "dm_posts")`に変更
- ファイル名を変更: `post_repository.go` → `dm_post_repository.go`

**受け入れ基準**:
- 全てのテーブル名参照が`"dm_posts"`に変更されている
- ファイル名が`dm_post_repository.go`に変更されている
- ビルドエラーがない

---

#### - [x] タスク 4.4: PostRepositoryGORMのテーブル名参照変更
**目的**: `server/internal/repository/post_repository_gorm.go`のテーブル名参照を変更

**作業内容**:
- `server/internal/repository/post_repository_gorm.go`を開く
- `GetTableName("posts", ...)`を`GetTableName("dm_posts", ...)`に変更
- `GetShardingConnectionByID(..., "posts")`を`GetShardingConnectionByID(..., "dm_posts")`に変更
- ファイル名を変更: `post_repository_gorm.go` → `dm_post_repository_gorm.go`

**受け入れ基準**:
- 全てのテーブル名参照が`"dm_posts"`に変更されている
- ファイル名が`dm_post_repository_gorm.go`に変更されている
- ビルドエラーがない

---

#### - [x] タスク 4.5: NewsRepositoryのテーブル名参照変更 (N/A - ファイルが存在しない)
**目的**: `server/internal/repository/news_repository.go`のテーブル名参照を変更

**作業内容**:
- `server/internal/repository/news_repository.go`を開く
- 文字列リテラル`"news"`を`"dm_news"`に変更
- ファイル名を変更: `news_repository.go` → `dm_news_repository.go`

**受け入れ基準**:
- 全てのテーブル名参照が`"dm_news"`に変更されている
- ファイル名が`dm_news_repository.go`に変更されている
- ビルドエラーがない

---

#### - [x] タスク 4.6: NewsRepositoryGORMのテーブル名参照変更 (N/A - ファイルが存在しない)
**目的**: `server/internal/repository/news_repository_gorm.go`のテーブル名参照を変更

**作業内容**:
- `server/internal/repository/news_repository_gorm.go`を開く
- `Table("news")`を`Table("dm_news")`に変更
- ファイル名を変更: `news_repository_gorm.go` → `dm_news_repository_gorm.go`

**受け入れ基準**:
- 全てのテーブル名参照が`"dm_news"`に変更されている
- ファイル名が`dm_news_repository_gorm.go`に変更されている
- ビルドエラーがない

---

#### - [x] タスク 4.7: Repository層のビルド確認
**目的**: Repository層の変更後、ビルドエラーがないことを確認

**作業内容**:
- `go build ./server/internal/repository/...`を実行
- ビルドエラーがないことを確認
- import文が正しく動作することを確認

**受け入れ基準**:
- ビルドエラーがない
- import文が正しく動作する

---

### Phase 5: テストコードの更新

#### - [x] タスク 5.1: テストユーティリティのスキーマ作成SQL変更
**目的**: `server/test/testutil/db.go`のスキーマ作成SQLを変更

**作業内容**:
- `server/test/testutil/db.go`を開く
- `InitMasterSchema`関数の`CREATE TABLE IF NOT EXISTS news`を`CREATE TABLE IF NOT EXISTS dm_news`に変更
- `InitShardingSchema`関数の`CREATE TABLE IF NOT EXISTS users_%s`を`CREATE TABLE IF NOT EXISTS dm_users_%s`に変更
- `InitShardingSchema`関数の`CREATE TABLE IF NOT EXISTS posts_%s`を`CREATE TABLE IF NOT EXISTS dm_posts_%s`に変更
- 外部キー制約の定義を削除（分散データ環境では使用しない）

**受け入れ基準**:
- スキーマ作成SQLのテーブル名が正しく変更されている
- 外部キー制約が削除されている
- ビルドエラーがない

---

#### - [x] タスク 5.2: 統合テストのテーブル名参照変更
**目的**: `server/test/integration/sharding_test.go`のテーブル名参照を変更

**作業内容**:
- `server/test/integration/sharding_test.go`を開く
- `GetTableName("users", ...)`を`GetTableName("dm_users", ...)`に変更
- 期待値のテーブル名を変更（`"users_000"` → `"dm_users_000"`など）

**受け入れ基準**:
- 全てのテーブル名参照が正しく変更されている
- 期待値が正しく更新されている
- ビルドエラーがない

---

#### - [x] タスク 5.3: ユニットテストのテストケース変更
**目的**: `server/internal/db/sharding_test.go`のテストケースを変更

**作業内容**:
- `server/internal/db/sharding_test.go`を開く
- テストケースの`baseName: "users"`を`baseName: "dm_users"`に変更
- テストケースの`baseName: "posts"`を`baseName: "dm_posts"`に変更
- 期待値のテーブル名を変更（`"users_000"` → `"dm_users_000"`など）

**受け入れ基準**:
- 全てのテストケースが正しく変更されている
- 期待値が正しく更新されている
- ビルドエラーがない

---

#### - [x] タスク 5.4: Repositoryテストファイルのファイル名変更
**目的**: Repositoryテストファイルのファイル名を変更

**作業内容**:
- `server/internal/repository/user_repository_test.go`を`dm_user_repository_test.go`に変更
- `server/internal/repository/user_repository_gorm_test.go`を`dm_user_repository_gorm_test.go`に変更
- `server/internal/repository/post_repository_test.go`を`dm_post_repository_test.go`に変更
- `server/internal/repository/post_repository_gorm_test.go`を`dm_post_repository_gorm_test.go`に変更
- `server/internal/repository/news_repository_test.go`を`dm_news_repository_test.go`に変更
- `server/internal/repository/news_repository_gorm_test.go`を`dm_news_repository_gorm_test.go`に変更

**受け入れ基準**:
- 全てのテストファイル名が正しく変更されている
- ビルドエラーがない

---

#### - [x] タスク 5.5: テストコードのテーブル名参照変更
**目的**: Repositoryテストファイル内のテーブル名参照を変更

**作業内容**:
- 各Repositoryテストファイルを開く
- テーブル名参照を変更（`"users"` → `"dm_users"`, `"posts"` → `"dm_posts"`, `"news"` → `"dm_news"`）
- 期待値のテーブル名を変更

**受け入れ基準**:
- 全てのテストコード内のテーブル名参照が正しく変更されている
- 期待値が正しく更新されている
- ビルドエラーがない

---

#### - [x] タスク 5.6: 全テストの実行と検証
**目的**: テストコードの変更後、全テストが通過することを確認

**作業内容**:
- `go test ./...`を実行
- 全テストが通過することを確認
- エラーが発生した場合は修正

**受け入れ基準**:
- 全テストが通過する
- エラーがない

---

### Phase 6: CLIツール、管理画面、ドキュメントの更新

#### - [x] タスク 6.1: サンプルデータ生成ツールのテーブル名変更
**目的**: `server/cmd/generate-sample-data/main.go`のテーブル名を変更

**作業内容**:
- `server/cmd/generate-sample-data/main.go`を開く
- `fmt.Sprintf("users_%03d", ...)`を`fmt.Sprintf("dm_users_%03d", ...)`に変更
- `fmt.Sprintf("posts_%03d", ...)`を`fmt.Sprintf("dm_posts_%03d", ...)`に変更
- `Table("news")`を`Table("dm_news")`に変更

**受け入れ基準**:
- 全てのテーブル名参照が正しく変更されている
- ビルドエラーがない

---

#### - [x] タスク 6.2: CLIツールのビルド確認
**目的**: CLIツールの変更後、ビルドエラーがないことを確認

**作業内容**:
- `go build ./server/cmd/generate-sample-data/...`を実行
- ビルドエラーがないことを確認

**受け入れ基準**:
- ビルドエラーがない
- CLIツールが正常にビルドされる

---

#### - [x] タスク 6.3: GoAdmin設定のテーブル名変更
**目的**: `server/internal/admin/tables.go`のGoAdmin設定を変更

**作業内容**:
- `server/internal/admin/tables.go`を開く
- `SetTable("news")`を`SetTable("dm_news")`に変更

**受け入れ基準**:
- GoAdmin設定のテーブル名が正しく変更されている
- ビルドエラーがない

---

#### - [x] タスク 6.4: シャーディングドキュメントの更新
**目的**: `docs/Sharding.md`のテーブル名を更新

**作業内容**:
- `docs/Sharding.md`を開く
- テーブル名の記載を全て更新:
  - `users` → `dm_users`
  - `posts` → `dm_posts`
  - `news` → `dm_news`
- コード例のテーブル名を更新
- 図表や説明文のテーブル名を更新

**受け入れ基準**:
- 全てのテーブル名が正しく更新されている
- コード例が正しく更新されている
- ドキュメントが一貫している

---

### Phase 7: データベース再作成と検証

#### - [x] タスク 7.1: 既存データベースファイルの削除
**目的**: 既存データベースファイルを削除して再作成の準備

**作業内容**:
- `server/data/master.db`を削除
- `server/data/sharding_db_1.db`を削除
- `server/data/sharding_db_2.db`を削除
- `server/data/sharding_db_3.db`を削除
- `server/data/sharding_db_4.db`を削除

**受け入れ基準**:
- 全ての既存データベースファイルが削除されている

---

#### - [x] タスク 7.2: マイグレーション適用スクリプトの実行
**目的**: マイグレーション適用スクリプトを実行してデータベースを再作成

**作業内容**:
- `./scripts/migrate.sh all`を実行
- マイグレーションが正常に適用されることを確認
- エラーが発生した場合は確認・修正

**受け入れ基準**:
- マイグレーションが正常に適用される
- データベースに正しいテーブル名でテーブルが作成されている

---

#### - [x] タスク 7.3: データベーススキーマの確認
**目的**: データベースに正しいテーブル名でテーブルが作成されていることを確認

**作業内容**:
- マスターデータベースに`dm_news`テーブルが存在することを確認
- 各シャーディングデータベースに`dm_users_*`, `dm_posts_*`テーブルが存在することを確認
- インデックス名が正しく変更されていることを確認

**受け入れ基準**:
- 全てのテーブルが正しい名前で作成されている
- インデックス名が正しく変更されている

---

#### - [x] タスク 7.4: 全テストの実行と検証
**目的**: データベース再作成後、全テストが通過することを確認

**作業内容**:
- `go test ./...`を実行
- 全テストが通過することを確認
- エラーが発生した場合は確認・修正

**受け入れ基準**:
- 全テストが通過する
- エラーがない

---

#### - [x] タスク 7.5: CLIツールの動作確認
**目的**: CLIツールが正常に動作することを確認

**作業内容**:
- `APP_ENV=develop ./bin/generate-sample-data`を実行
- 正常にデータが生成されることを確認
- 正しいテーブル名（`dm_users_*`, `dm_posts_*`, `dm_news`）にデータが生成されることを確認

**受け入れ基準**:
- CLIツールが正常に動作する
- 正しいテーブル名にデータが生成される

---

#### - [x] タスク 7.6: 管理画面の動作確認
**目的**: GoAdmin管理画面が正常に動作することを確認

**作業内容**:
- 管理画面を起動
- `dm_news`テーブルが正しく表示されることを確認
- データの表示・操作が正常に動作することを確認

**受け入れ基準**:
- 管理画面が正常に動作する
- `dm_news`テーブルが正しく表示される

---

#### - [x] タスク 7.7: 最終検証とドキュメント確認
**目的**: 全ての変更が正しく実装されていることを最終確認

**作業内容**:
- コードベース全体でテーブル名の参照を確認（正規表現検索）
- ファイル名変更が正しく実施されていることを確認
- ドキュメントが最新の状態になっていることを確認
- 受け入れ基準を全て満たしていることを確認

**受け入れ基準**:
- 全てのテーブル名参照が正しく変更されている
- 全てのファイル名が正しく変更されている
- ドキュメントが最新の状態になっている
- 受け入れ基準を全て満たしている

---

## 実装上の注意事項

### テーブル名の一貫性
- コードベース全体でテーブル名の参照を一貫して変更する
- 正規表現検索（`"users"`, `"posts"`, `"news"`）を使用して漏れがないか確認
- 文字列リテラル、コメント、ドキュメントなど、全ての箇所で変更を確認

### ファイル名変更の影響
- ファイル名変更後、import文が正しく動作するか確認
- `go mod tidy`を実行して依存関係を更新
- IDEやツールがファイル名を参照する場合があるため、ビルドエラーを確認

### 段階的実装
- 影響範囲が広いため、Phaseごとに段階的に実装する
- 各Phaseでテストを実行し、問題がないことを確認してから次のPhaseに進む
- 問題が発生した場合は、前のPhaseに戻って修正

### データベース再作成
- 既存データベースのデータは破棄して良いため、既存データベースファイルを削除して再作成する
- マイグレーション適用スクリプトを実行してデータベースを再作成する
- 再作成後、全テストを実行して動作確認

## 参考情報

### 関連ドキュメント
- 要件定義書: `.kiro/specs/0026-chtablename/requirements.md`
- 設計書: `.kiro/specs/0026-chtablename/design.md`

### 既存実装
- モデル層: `server/internal/model/*.go`
- Repository層: `server/internal/repository/*.go`
- データベース層: `server/internal/db/sharding.go`
- スキーマ定義: `db/schema/**/*.hcl`
