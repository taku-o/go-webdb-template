# テーブル名変更機能要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #54
- **Issueタイトル**: テーブル名を変更する (users, posts, news)
- **Feature名**: 0026-chtablename
- **作成日**: 2025-12-29

### 1.2 目的
現在のテーブル名（`users`, `posts`, `news`）は参考コード用のダミーテーブルであるため、新たにテーブルなどを追加していく際に開発の邪魔になる。そのため、テーブル名に`dm_`プレフィックスを付与して変更し、参考コード用であることを明確にする。

### 1.3 スコープ
- 対象テーブル: `users`, `posts`, `news`
- 変更内容:
  - `users` → `dm_users`（32分割テーブル: `users_000` → `dm_users_000`）
  - `posts` → `dm_posts`（32分割テーブル: `posts_000` → `dm_posts_000`）
  - `news` → `dm_news`
- コードベース全体でのテーブル名参照の変更
- スキーマ定義ファイル（HCL）の変更
- ファイル名の変更（モデル、Repository、スキーマ定義ファイル）
- テストコードの更新
- ドキュメントの更新

**本実装の範囲外**:
- テーブル構造の変更（名前のみ変更）
- データ移行処理（既存データベースのデータは破棄して良い）
- 本番環境でのデータ移行（開発環境のみを想定）

## 2. 背景・現状分析

### 2.1 現在の実装
- **テーブル名**:
  - `users`: 32分割テーブル（users_000～users_031）
  - `posts`: 32分割テーブル（posts_000～posts_031）
  - `news`: マスターデータベースの単一テーブル
- **テーブル名参照箇所**:
  - モデル層: `server/internal/model/*.go`の`TableName()`メソッド
  - Repository層: `server/internal/repository/*.go`での`GetTableName()`呼び出し
  - データベース層: `server/internal/db/sharding.go`でのテーブル名生成ロジック
  - スキーマ定義: `db/schema/**/*.hcl`でのテーブル定義
  - テストコード: `server/test/**/*.go`でのテーブル名参照
  - CLIツール: `server/cmd/generate-sample-data/main.go`でのテーブル名文字列
  - 管理画面: `server/internal/admin/tables.go`でのGoAdmin設定
  - ドキュメント: `docs/Sharding.md`でのテーブル名記載

### 2.2 課題点
1. **開発の邪魔になる**: 現在のテーブル名（`users`, `posts`, `news`）は参考コード用のダミーテーブルであるが、新たにテーブルなどを追加していく際に、実際のテーブル名と混同しやすい
2. **命名の明確性不足**: 参考コード用であることがテーブル名から判断できない
3. **将来の拡張性**: 実際のアプリケーション用のテーブルを追加する際に、名前の衝突を避ける必要がある

### 2.3 本実装による改善点
1. **命名の明確化**: `dm_`プレフィックスにより、参考コード用のダミーテーブルであることが明確になる
2. **開発効率の向上**: 実際のテーブル名とダミーテーブル名を区別しやすくなり、開発時の混乱を防ぐ
3. **将来の拡張性**: 実際のアプリケーション用のテーブル（例: `users`, `posts`）を追加する際に、名前の衝突を避けられる

## 3. 機能要件

### 3.1 モデル層の変更

#### 3.1.1 Userモデル
- **ファイル**: `server/internal/model/user.go`
- **変更内容**: `TableName()`メソッドの戻り値を`"users"`から`"dm_users"`に変更

#### 3.1.2 Postモデル
- **ファイル**: `server/internal/model/post.go`
- **変更内容**: `TableName()`メソッドの戻り値を`"posts"`から`"dm_posts"`に変更

#### 3.1.3 Newsモデル
- **ファイル**: `server/internal/model/news.go`
- **変更内容**: `TableName()`メソッドの戻り値を`"news"`から`"dm_news"`に変更

### 3.2 Repository層の変更

#### 3.2.1 UserRepository
- **ファイル**: `server/internal/repository/user_repository.go`, `user_repository_gorm.go`
- **変更内容**: 
  - `GetTableName("users", ...)`の呼び出しを`GetTableName("dm_users", ...)`に変更
  - 文字列リテラル`"users"`を`"dm_users"`に変更
  - テーブル名生成ロジック（`fmt.Sprintf("users_%03d", ...)`）を`fmt.Sprintf("dm_users_%03d", ...)`に変更

#### 3.2.2 PostRepository
- **ファイル**: `server/internal/repository/post_repository.go`, `post_repository_gorm.go`
- **変更内容**: 
  - `GetTableName("posts", ...)`の呼び出しを`GetTableName("dm_posts", ...)`に変更
  - 文字列リテラル`"posts"`を`"dm_posts"`に変更
  - テーブル名生成ロジック（`fmt.Sprintf("posts_%03d", ...)`）を`fmt.Sprintf("dm_posts_%03d", ...)`に変更

#### 3.2.3 NewsRepository
- **ファイル**: `server/internal/repository/news_repository.go`, `news_repository_gorm.go`
- **変更内容**: 
  - `Table("news")`の参照を`Table("dm_news")`に変更
  - 文字列リテラル`"news"`を`"dm_news"`に変更

### 3.3 データベース層の変更

#### 3.3.1 TableSelector
- **ファイル**: `server/internal/db/sharding.go`
- **変更内容**: 
  - `GetTableName()`メソッドのベース名引数（`"users"`, `"posts"`）を`"dm_users"`, `"dm_posts"`に変更
  - `GetShardingTableName()`関数のベース名引数を変更
  - `ValidateTableName()`関数の許可リスト（`allowedBaseNames`）に`"dm_users"`, `"dm_posts"`を追加

### 3.4 スキーマ定義の変更

#### 3.4.1 マスターデータベーススキーマ
- **ファイル**: `db/schema/master.hcl`
- **変更内容**: 
  - `table "news"`を`table "dm_news"`に変更
  - インデックス名の変更（`idx_news_*` → `idx_dm_news_*`）

#### 3.4.2 シャーディングデータベーススキーマ
- **ファイル**: 
  - `db/schema/sharding_1/users.hcl`, `posts.hcl`
  - `db/schema/sharding_2/users.hcl`, `posts.hcl`
  - `db/schema/sharding_3/users.hcl`, `posts.hcl`
  - `db/schema/sharding_4/users.hcl`, `posts.hcl`
- **変更内容**: 
  - `table "users_000"`～`table "users_031"`を`table "dm_users_000"`～`table "dm_users_031"`に変更
  - `table "posts_000"`～`table "posts_031"`を`table "dm_posts_000"`～`table "dm_posts_031"`に変更
  - インデックス名の変更（`idx_users_*` → `idx_dm_users_*`, `idx_posts_*` → `idx_dm_posts_*`）

### 3.5 テストコードの変更

#### 3.5.1 テストユーティリティ
- **ファイル**: `server/test/testutil/db.go`
- **変更内容**: 
  - スキーマ作成SQLのテーブル名を変更
  - `CREATE TABLE IF NOT EXISTS news` → `CREATE TABLE IF NOT EXISTS dm_news`
  - `CREATE TABLE IF NOT EXISTS users_%s` → `CREATE TABLE IF NOT EXISTS dm_users_%s`
  - `CREATE TABLE IF NOT EXISTS posts_%s` → `CREATE TABLE IF NOT EXISTS dm_posts_%s`
  - **注意**: 分散データ環境では外部キー制約を使用しないため、テストユーティリティに外部キー制約の定義がある場合は削除する

#### 3.5.2 統合テスト
- **ファイル**: `server/test/integration/sharding_test.go`
- **変更内容**: 
  - `GetTableName("users", ...)`の呼び出しを`GetTableName("dm_users", ...)`に変更
  - 期待値のテーブル名を変更（`"users_000"` → `"dm_users_000"`など）

#### 3.5.3 ユニットテスト
- **ファイル**: `server/internal/db/sharding_test.go`
- **変更内容**: 
  - テストケースのテーブル名を変更
  - `baseName: "users"` → `baseName: "dm_users"`
  - `baseName: "posts"` → `baseName: "dm_posts"`
  - 期待値のテーブル名を変更

### 3.6 CLIツールの変更

#### 3.6.1 サンプルデータ生成ツール
- **ファイル**: `server/cmd/generate-sample-data/main.go`
- **変更内容**: 
  - テーブル名の文字列リテラルを変更
  - `fmt.Sprintf("users_%03d", ...)` → `fmt.Sprintf("dm_users_%03d", ...)`
  - `fmt.Sprintf("posts_%03d", ...)` → `fmt.Sprintf("dm_posts_%03d", ...)`
  - `Table("news")` → `Table("dm_news")`

### 3.7 管理画面の変更

#### 3.7.1 GoAdmin設定
- **ファイル**: `server/internal/admin/tables.go`
- **変更内容**: 
  - `SetTable("news")`を`SetTable("dm_news")`に変更

### 3.8 ドキュメントの更新

#### 3.8.1 シャーディングドキュメント
- **ファイル**: `docs/Sharding.md`
- **変更内容**: 
  - テーブル名の記載を更新（`users` → `dm_users`, `posts` → `dm_posts`, `news` → `dm_news`）
  - コード例のテーブル名を更新

## 4. 非機能要件

### 4.1 既存機能への影響
- 既存の機能は全て正常に動作すること
- テーブル構造の変更は行わない（名前のみ変更）
- 既存のシャーディング構造は維持すること

### 4.2 データ移行
- **既存データベースのデータは破棄して良い**
- データ移行処理は実装しない
- 既存データベースファイルの再作成が可能
- スキーマ定義を変更してデータベースを再作成する方針で実装

### 4.3 テスト
- 既存のテストが全て通過すること
- テーブル名変更に関する新しいテストは不要（既存テストの更新のみ）

### 4.4 パフォーマンス
- テーブル名変更によるパフォーマンスへの影響はないこと
- 既存のクエリパフォーマンスは維持されること

## 5. 制約事項

### 5.1 技術的制約
- 既存のシャーディング構造（32分割、4データベース）は維持すること
- テーブル構造の変更は行わない（名前のみ変更）
- 既存のデータベース接続管理機能を再利用すること
- Atlasスキーマ定義（HCL）の変更が必要

### 5.2 プロジェクト制約
- 既存のレイヤードアーキテクチャを維持すること
- 既存のコーディング規約に従うこと
- 既存のテストパターンを維持すること

### 5.3 データベース制約
- 既存データベースのデータは破棄して良い
- マイグレーション戦略: テーブル名変更のマイグレーションではなく、スキーマ定義を変更してデータベースを再作成する方針

### 5.4 命名規則
- プレフィックス: `dm_`（dummyの略）
- 変更後のテーブル名:
  - `dm_users`（32分割: `dm_users_000`～`dm_users_031`）
  - `dm_posts`（32分割: `dm_posts_000`～`dm_posts_031`）
  - `dm_news`（単一テーブル）

## 6. 受け入れ基準

### 6.1 モデル層
- [ ] `User.TableName()`が`"dm_users"`を返す
- [ ] `Post.TableName()`が`"dm_posts"`を返す
- [ ] `News.TableName()`が`"dm_news"`を返す

### 6.2 Repository層
- [ ] UserRepositoryで`"dm_users"`を参照している
- [ ] PostRepositoryで`"dm_posts"`を参照している
- [ ] NewsRepositoryで`"dm_news"`を参照している
- [ ] 全32分割テーブル（`dm_users_000`～`dm_users_031`, `dm_posts_000`～`dm_posts_031`）が正しく参照されている

### 6.3 データベース層
- [ ] `TableSelector.GetTableName("dm_users", ...)`が正しく動作する
- [ ] `TableSelector.GetTableName("dm_posts", ...)`が正しく動作する
- [ ] `ValidateTableName()`が`"dm_users"`, `"dm_posts"`を許可する

### 6.4 スキーマ定義
- [ ] `db/schema/master.hcl`で`table "dm_news"`が定義されている
- [ ] `db/schema/sharding_*/users.hcl`で`table "dm_users_000"`～`table "dm_users_031"`が定義されている
- [ ] `db/schema/sharding_*/posts.hcl`で`table "dm_posts_000"`～`table "dm_posts_031"`が定義されている
- [ ] インデックス名が正しく変更されている（`idx_dm_users_*`, `idx_dm_posts_*`, `idx_dm_news_*`）

### 6.5 テストコード
- [ ] 全テストが通過する
- [ ] テストコード内のテーブル名参照が正しく更新されている

### 6.6 CLIツール
- [ ] `generate-sample-data`コマンドが正常に動作する
- [ ] 正しいテーブル名（`dm_users_*`, `dm_posts_*`, `dm_news`）にデータが生成される

### 6.7 管理画面
- [ ] GoAdminで`dm_news`テーブルが正しく表示される

### 6.8 ファイル名変更
- [ ] モデルファイル名が変更されている（`user.go` → `dm_user.go`, `post.go` → `dm_post.go`, `news.go` → `dm_news.go`）
- [ ] Repositoryファイル名が変更されている（`*_repository.go` → `dm_*_repository.go`）
- [ ] スキーマ定義ファイル名が変更されている（`users.hcl` → `dm_users.hcl`, `posts.hcl` → `dm_posts.hcl`）
- [ ] import文や参照箇所が正しく更新されている

### 6.9 ドキュメント
- [ ] `docs/Sharding.md`のテーブル名が更新されている

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### モデル層
- `server/internal/model/user.go`: `TableName()`メソッド
- `server/internal/model/post.go`: `TableName()`メソッド
- `server/internal/model/news.go`: `TableName()`メソッド

#### Repository層
- `server/internal/repository/user_repository.go`: テーブル名参照
- `server/internal/repository/user_repository_gorm.go`: テーブル名参照
- `server/internal/repository/post_repository.go`: テーブル名参照
- `server/internal/repository/post_repository_gorm.go`: テーブル名参照
- `server/internal/repository/news_repository.go`: テーブル名参照
- `server/internal/repository/news_repository_gorm.go`: テーブル名参照

#### データベース層
- `server/internal/db/sharding.go`: テーブル名生成ロジック、許可リスト

#### スキーマ定義
- `db/schema/master.hcl`: `table "news"` → `table "dm_news"`
- `db/schema/sharding_1/users.hcl`: `table "users_000"`～`table "users_007"` → `table "dm_users_000"`～`table "dm_users_007"`
- `db/schema/sharding_1/posts.hcl`: `table "posts_000"`～`table "posts_007"` → `table "dm_posts_000"`～`table "dm_posts_007"`
- `db/schema/sharding_2/users.hcl`: `table "users_008"`～`table "users_015"` → `table "dm_users_008"`～`table "dm_users_015"`
- `db/schema/sharding_2/posts.hcl`: `table "posts_008"`～`table "posts_015"` → `table "dm_posts_008"`～`table "dm_posts_015"`
- `db/schema/sharding_3/users.hcl`: `table "users_016"`～`table "users_023"` → `table "dm_users_016"`～`table "dm_users_023"`
- `db/schema/sharding_3/posts.hcl`: `table "posts_016"`～`table "posts_023"` → `table "dm_posts_016"`～`table "dm_posts_023"`
- `db/schema/sharding_4/users.hcl`: `table "users_024"`～`table "users_031"` → `table "dm_users_024"`～`table "dm_users_031"`
- `db/schema/sharding_4/posts.hcl`: `table "posts_024"`～`table "posts_031"` → `table "dm_posts_024"`～`table "dm_posts_031"`

#### テストコード
- `server/test/testutil/db.go`: スキーマ作成SQL
- `server/test/integration/sharding_test.go`: テーブル名参照
- `server/internal/db/sharding_test.go`: テストケース

#### CLIツール
- `server/cmd/generate-sample-data/main.go`: テーブル名文字列

#### 管理画面
- `server/internal/admin/tables.go`: GoAdmin設定

#### ドキュメント
- `docs/Sharding.md`: テーブル名記載

### 7.2 ファイル名の変更

以下のファイル名を変更する：

#### モデル層
- `server/internal/model/user.go` → `server/internal/model/dm_user.go`
- `server/internal/model/post.go` → `server/internal/model/dm_post.go`
- `server/internal/model/news.go` → `server/internal/model/dm_news.go`

#### Repository層
- `server/internal/repository/user_repository.go` → `server/internal/repository/dm_user_repository.go`
- `server/internal/repository/user_repository_gorm.go` → `server/internal/repository/dm_user_repository_gorm.go`
- `server/internal/repository/user_repository_test.go` → `server/internal/repository/dm_user_repository_test.go`
- `server/internal/repository/user_repository_gorm_test.go` → `server/internal/repository/dm_user_repository_gorm_test.go`
- `server/internal/repository/post_repository.go` → `server/internal/repository/dm_post_repository.go`
- `server/internal/repository/post_repository_gorm.go` → `server/internal/repository/dm_post_repository_gorm.go`
- `server/internal/repository/post_repository_test.go` → `server/internal/repository/dm_post_repository_test.go`
- `server/internal/repository/post_repository_gorm_test.go` → `server/internal/repository/dm_post_repository_gorm_test.go`
- `server/internal/repository/news_repository.go` → `server/internal/repository/dm_news_repository.go`
- `server/internal/repository/news_repository_gorm.go` → `server/internal/repository/dm_news_repository_gorm.go`
- `server/internal/repository/news_repository_test.go` → `server/internal/repository/dm_news_repository_test.go`
- `server/internal/repository/news_repository_gorm_test.go` → `server/internal/repository/dm_news_repository_gorm_test.go`

#### スキーマ定義
- `db/schema/sharding_1/users.hcl` → `db/schema/sharding_1/dm_users.hcl`
- `db/schema/sharding_1/posts.hcl` → `db/schema/sharding_1/dm_posts.hcl`
- `db/schema/sharding_2/users.hcl` → `db/schema/sharding_2/dm_users.hcl`
- `db/schema/sharding_2/posts.hcl` → `db/schema/sharding_2/dm_posts.hcl`
- `db/schema/sharding_3/users.hcl` → `db/schema/sharding_3/dm_users.hcl`
- `db/schema/sharding_3/posts.hcl` → `db/schema/sharding_3/dm_posts.hcl`
- `db/schema/sharding_4/users.hcl` → `db/schema/sharding_4/dm_users.hcl`
- `db/schema/sharding_4/posts.hcl` → `db/schema/sharding_4/dm_posts.hcl`

**注意**: ファイル名変更後、import文や参照箇所の更新が必要

### 7.3 新規追加が必要なファイル
なし（既存ファイルの変更とリネームのみ）

### 7.4 削除されるファイル
なし（既存ファイルは変更またはリネームのみ）

### 7.5 再利用する既存機能
- `server/internal/db/manager.go`: データベース接続管理機能（変更なし）
- `server/internal/config/config.go`: 設定読み込み機能（変更なし）
- 既存のレイヤードアーキテクチャ（変更なし）

## 8. 実装上の注意事項

### 8.1 シャーディングテーブルの一括変更
- users, postsテーブルは32分割テーブル（`dm_users_000`～`dm_users_031`, `dm_posts_000`～`dm_posts_031`）に変更する
- 4つのシャーディングデータベース（sharding_1～sharding_4）全てで変更が必要
- 各データベースに8テーブルずつ存在するため、合計32テーブル×2種類（users, posts）= 64テーブルの変更が必要

### 8.2 インデックス名の変更
- インデックス名もテーブル名に合わせて変更する必要がある
- 例: `idx_users_000_email` → `idx_dm_users_000_email`
- 例: `idx_posts_000_user_id` → `idx_dm_posts_000_user_id`
- 例: `idx_news_published_at` → `idx_dm_news_published_at`

### 8.3 スキーマ定義ファイルの変更
- Atlasスキーマ定義（HCL）ファイルを変更する
- 変更後、マイグレーションファイルを再生成する必要がある
- 既存データベースのデータは破棄して良いため、スキーマ定義を変更してデータベースを再作成する方針で実装

### 8.4 テーブル名の一貫性
- コードベース全体でテーブル名の参照を一貫して変更する
- 文字列リテラル、コメント、ドキュメントなど、全ての箇所で変更を確認する
- 正規表現検索（`"users"`, `"posts"`, `"news"`）を使用して漏れがないか確認する

### 8.5 テストコードの更新
- テストコード内のテーブル名参照も全て更新する
- 期待値のテーブル名も変更する
- テストが全て通過することを確認する

### 8.6 ファイル名変更とimport文の更新
- モデル、Repository、スキーマ定義ファイルのファイル名を変更する
- ファイル名変更後、以下の箇所でimport文や参照箇所の更新が必要:
  - モデルファイルをimportしている箇所（Repository、テストコードなど）
  - Repositoryファイルをimportしている箇所（ハンドラー、サービス層など）
  - スキーマ定義ファイルを参照している箇所（Atlas設定ファイルなど）
- Goのimport文はパッケージ名ではなくファイル名に基づくため、importパスも更新が必要

### 8.7 データベース再作成
- 既存データベースのデータは破棄して良いため、既存データベースファイルを削除して再作成する
- マイグレーション適用スクリプト（`scripts/migrate.sh`）を実行してデータベースを再作成する

## 9. 参考情報

### 9.1 関連Issue
- GitHub Issue #54: テーブル名を変更する (users, posts, news)

### 9.2 既存ドキュメント
- `docs/Sharding.md`: シャーディング戦略の詳細
- `docs/atlas-operations.md`: Atlasマイグレーション運用ガイド

### 9.3 既存実装
- `server/internal/model/user.go`: Userモデル
- `server/internal/model/post.go`: Postモデル
- `server/internal/model/news.go`: Newsモデル
- `server/internal/repository/user_repository.go`: UserRepository実装
- `server/internal/repository/post_repository.go`: PostRepository実装
- `server/internal/repository/news_repository.go`: NewsRepository実装
- `server/internal/db/sharding.go`: テーブル選択ロジック
- `db/schema/master.hcl`: マスターデータベーススキーマ定義
- `db/schema/sharding_*/users.hcl`: シャーディングデータベースusersテーブル定義
- `db/schema/sharding_*/posts.hcl`: シャーディングデータベースpostsテーブル定義

### 9.4 技術スタック
- **Go**: 1.21+
- **GORM**: v1.25.12
- **データベース**: SQLite3（開発環境）
- **Atlas**: スキーマ管理ツール

### 9.5 変更パターン
- テーブル名: `users` → `dm_users`, `posts` → `dm_posts`, `news` → `dm_news`
- インデックス名: `idx_users_*` → `idx_dm_users_*`, `idx_posts_*` → `idx_dm_posts_*`, `idx_news_*` → `idx_dm_news_*`

**注意**: 分散データ環境では外部キー制約を使用しないため、外部キー制約の変更は不要。既存のスキーマ定義（HCL）にも外部キー制約の定義は存在しない。
