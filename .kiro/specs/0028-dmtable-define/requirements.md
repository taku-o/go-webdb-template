# 分散テーブル環境対応テーブル設計修正要件定義書

## 1. 概要

### 1.1 プロジェクト情報

- **プロジェクト名**: go-webdb-template
- **Issue番号**: #52
- **Issueタイトル**: テーブル設計が分散テーブル環境にあった設計になっていない
- **Feature名**: 0028-dmtable-define
- **作成日**: 2025-01-27

### 1.2 目的

分散テーブル環境では、auto_incrementによるID生成は使用できない。本実装では、sonyflakeライブラリを導入し、分散環境で一意性が保証されるID生成方式に変更する。これにより、テーブル設計を分散テーブル環境に適したものに修正する。

### 1.3 スコープ

- sonyflakeライブラリの導入
- テーブル定義の修正（Atlas形式）
- dm_newsテーブルのidフィールド修正
- dm_users_NNNテーブルのidフィールド修正
- dm_posts_NNNテーブルのidフィールド修正
- モデル定義の修正（GORMタグの修正）
- ID生成ロジックの実装（sonyflakeを使用）
- Repository層でのID生成統合
- サンプルデータ生成コマンドの修正（server/cmd/generate-sample-data/）
- sharding規則の定義（server/internal/db/sharding.go）
- identifier生成ルールのドキュメント化（数値: sonyflake、文字列: UUIDv7）

**本実装の範囲外**:

- 既存データの移行（マイグレーション）
- 他のテーブル（GoAdmin関連テーブルなど）の修正
- パフォーマンス最適化

## 2. 背景・現状分析

### 2.1 現状の問題

分散テーブル環境では、複数のデータベースにまたがるテーブルが存在する。auto_incrementを使用すると、各データベースで独立したシーケンスが生成されるため、IDの一意性が保証されない。また、データを別のデータベースに移動する際に問題が発生する可能性がある。

### 2.2 現状のテーブル定義

#### 2.2.1 dm_newsテーブル（master.hcl）

- 現在: `id integer auto_increment = true`
- 修正後: `id BIGINT (unsigned) auto_increment = false`

#### 2.2.2 dm_users_NNNテーブル（sharding_*/dm_users.hcl）

- 現在: `id integer`（auto_incrementの記述なしだが、実際には使用されている可能性）
- 修正後: `id BIGINT (unsigned) auto_increment = false`
- table_sharding_key: `id`
- テーブル番号計算: `id % DBShardingTableCount`

#### 2.2.3 dm_posts_NNNテーブル（sharding_*/dm_posts.hcl）

- 現在: `id integer`（auto_incrementの記述なしだが、実際には使用されている可能性）
- 修正後: `id BIGINT (unsigned) auto_increment = false`
- table_sharding_key: `user_id`
- テーブル番号計算: `user_id % DBShardingTableCount`
- **重要な規則**: ある`dm_users`に紐付いた`dm_posts`は、同じテーブル番号のテーブルにデータが入る

### 2.3 現状のモデル定義

#### 2.3.1 DmNewsモデル（server/internal/model/dm_news.go）

- 現在: `gorm:"primaryKey;autoIncrement"`
- 修正後: `gorm:"primaryKey"`（autoIncrementを削除）

#### 2.3.2 DmUserモデル（server/internal/model/dm_user.go）

- 現在: `gorm:"primaryKey"`（autoIncrementの記述なし）
- 修正後: 変更不要（ただし、ID生成ロジックの追加が必要）

#### 2.3.3 DmPostモデル（server/internal/model/dm_post.go）

- 現在: `gorm:"primaryKey"`（autoIncrementの記述なし）
- 修正後: 変更不要（ただし、ID生成ロジックの追加が必要）

### 2.4 現状のID生成

現在、ID生成の明示的な実装は確認できていない。GORMのautoIncrement機能に依存している可能性が高い。

## 3. 要件定義

### 3.1 ライブラリ導入要件

#### REQ-1: sonyflakeライブラリの導入

- **要件**: `github.com/sony/sonyflake`ライブラリをプロジェクトに導入する
- **詳細**:
  - `go.mod`に依存関係を追加
  - 適切なバージョンを指定
- **受け入れ基準**:
  - `go mod tidy`が正常に実行できる
  - ビルドが正常に完了する

### 3.2 テーブル定義修正要件

#### REQ-2: dm_newsテーブル定義の修正

- **要件**: `db/schema/master.hcl`のdm_newsテーブルのidカラムを修正する
- **詳細**:
  - `type = integer` → `type = bigint`
  - `auto_increment = true` → `auto_increment = false`
  - unsigned属性を追加（Atlas形式で可能な場合）
- **受け入れ基準**:
  - Atlasスキーマ検証が正常に完了する
  - マイグレーションスクリプトが正常に生成される

#### REQ-3: dm_users_NNNテーブル定義の修正

- **要件**: `db/schema/sharding_*/dm_users.hcl`の全テーブルのidカラムを修正する
- **詳細**:
  - `type = integer` → `type = bigint`
  - `auto_increment = false`を明示的に追加
  - unsigned属性を追加（Atlas形式で可能な場合）
  - table_sharding_keyとして`id`を定義（sharding.goに追加）
- **対象ファイル**:
  - `db/schema/sharding_1/dm_users.hcl`
  - `db/schema/sharding_2/dm_users.hcl`
  - `db/schema/sharding_3/dm_users.hcl`
  - `db/schema/sharding_4/dm_users.hcl`
- **受け入れ基準**:
  - 全shardingグループのテーブル定義が修正されている
  - Atlasスキーマ検証が正常に完了する

#### REQ-4: dm_posts_NNNテーブル定義の修正

- **要件**: `db/schema/sharding_*/dm_posts.hcl`の全テーブルのidカラムを修正する
- **詳細**:
  - `type = integer` → `type = bigint`
  - `auto_increment = false`を明示的に追加
  - unsigned属性を追加（Atlas形式で可能な場合）
  - table_sharding_keyとして`user_id`を定義（sharding.goに追加）
- **対象ファイル**:
  - `db/schema/sharding_1/dm_posts.hcl`
  - `db/schema/sharding_2/dm_posts.hcl`
  - `db/schema/sharding_3/dm_posts.hcl`
  - `db/schema/sharding_4/dm_posts.hcl`
- **受け入れ基準**:
  - 全shardingグループのテーブル定義が修正されている
  - Atlasスキーマ検証が正常に完了する

### 3.3 モデル定義修正要件

#### REQ-5: DmNewsモデルの修正

- **要件**: `server/internal/model/dm_news.go`のDmNews構造体のIDフィールドのGORMタグを修正する
- **詳細**:
  - `gorm:"primaryKey;autoIncrement"` → `gorm:"primaryKey"`
  - autoIncrementタグを削除
- **受け入れ基準**:
  - モデルファイルが正常にコンパイルできる
  - GORMがautoIncrementを使用しない

#### REQ-6: ID型の確認と統一

- **要件**: 全モデルのIDフィールドがint64型であることを確認する
- **詳細**:
  - DmUser.ID: int64（確認済み）
  - DmPost.ID: int64（確認済み）
  - DmNews.ID: int64（確認済み）
- **受け入れ基準**:
  - 全モデルのIDフィールドがint64型である

#### REQ-13: JavaScript/API側でのIDの文字列化要件

- **要件**: JavaScriptの制限により、sonyflakeで生成したIDは文字列として扱う
- **詳細**:
  - 全モデルのIDフィールドに`json:"id,string"`タグが設定されていることを確認
  - APIレスポンスでIDが文字列として返されることを確認
  - JavaScript側で数値精度の問題が発生しないことを保証
  - sonyflakeで生成されるIDは64ビット整数であり、JavaScriptのNumber型の安全な整数範囲（2^53-1）を超える可能性があるため、文字列として扱う必要がある
- **現状確認**:
  - DmUser.ID: `json:"id,string"`（確認済み）
  - DmPost.ID: `json:"id,string"`（確認済み）
  - DmNews.ID: `json:"id,string"`（確認済み）
- **受け入れ基準**:
  - 全モデルのIDフィールドに`json:"id,string"`タグが設定されている
  - APIレスポンスでIDが文字列として返される
  - JavaScript側でIDを文字列として扱える
  - 既存のAPIインターフェースとの互換性が維持されている

### 3.4 ID生成ロジック実装要件

#### REQ-7: ID生成ユーティリティの実装

- **要件**: sonyflakeを使用したID生成ユーティリティを実装する
- **詳細**:
  - パッケージ: `server/internal/util/idgen`（新規作成）
  - 関数: `GenerateSonyflakeID() (int64, error)`
  - ファイル: `server/internal/util/idgen/sonyflake.go`（パッケージ名は`idgen`）
  - sonyflakeインスタンスの初期化と管理
  - エラーハンドリングの実装
- **受け入れ基準**:
  - ID生成関数が正常に動作する
  - 生成されるIDが一意である
  - エラーが適切に処理される

#### REQ-8: Repository層でのID生成統合

- **要件**: 各RepositoryのCreateメソッドでID生成を実装する
- **詳細**:
  - DmUserRepository.Create: ID生成を追加
  - DmPostRepository.Create: ID生成を追加
  - DmNewsRepository.Create: ID生成を追加（存在する場合）
- **対象ファイル**:
  - `server/internal/repository/dm_user_repository.go`
  - `server/internal/repository/dm_user_repository_gorm.go`
  - `server/internal/repository/dm_post_repository.go`
  - `server/internal/repository/dm_post_repository_gorm.go`
- **受け入れ基準**:
  - 新規作成時にIDが自動生成される
  - 既存のテストが正常に動作する

#### REQ-14: サンプルデータ生成コマンドの修正

- **要件**: `server/cmd/generate-sample-data/main.go`でID生成を統合する
- **詳細**:
  - `generateDmUsers`: DmUser作成時にsonyflakeでIDを生成し、モデルに設定
  - `generateDmPosts`: DmPost作成時にsonyflakeでIDを生成し、モデルに設定
  - `generateDmNews`: DmNews作成時にsonyflakeでIDを生成し、モデルに設定
  - `insertDmUsersBatch`: INSERT文に`id`カラムを追加し、生成したIDを含める
  - `insertDmPostsBatch`: INSERT文に`id`カラムを追加し、生成したIDを含める
  - `insertDmNewsBatch`: GORMのCreateInBatchesを使用しているが、IDが設定されていることを確認
  - `fetchDmUserIDs`: 挿入後にIDを取得する必要がなくなるため、削除または変更を検討
  - ID生成ユーティリティ（`server/internal/util/idgen`）をインポートして使用
- **対象ファイル**:
  - `server/cmd/generate-sample-data/main.go`
- **受け入れ基準**:
  - サンプルデータ生成時にIDが正常に生成される
  - バッチ挿入が正常に動作する
  - 生成されたIDが一意である
  - コマンドが正常に実行できる

### 3.5 Sharding規則定義要件

#### REQ-9: sharding規則の定義

- **要件**: `server/internal/db/sharding.go`にテーブルsharding規則を定義する
- **詳細**:
  - dm_users_NNN: table_sharding_key = `id`
    - テーブル番号 = `id % DBShardingTableCount`
  - dm_posts_NNN: table_sharding_key = `user_id`
    - テーブル番号 = `user_id % DBShardingTableCount`
  - **重要な規則**: ある`dm_users`に紐付いた`dm_posts`は、同じテーブル番号のテーブルにデータが入る
    - `dm_users`のIDと`dm_posts`の`user_id`が同じ値であれば、同じテーブル番号になる
    - これにより、JOIN操作が効率的に行える
  - 定数または構造体として定義
  - ドキュメントコメントを追加
- **受け入れ基準**:
  - sharding規則が明確に定義されている
  - コードから規則が理解できる
  - `dm_users`のIDと`dm_posts`の`user_id`が同じ値の場合、同じテーブル番号になることが保証されている

### 3.6 テスト要件

#### REQ-10: 単体テストの実装

- **要件**: ID生成ユーティリティの単体テストを実装する
- **詳細**:
  - ID生成関数のテスト
  - 一意性のテスト
  - エラーハンドリングのテスト
- **受け入れ基準**:
  - テストが正常に実行できる
  - カバレッジが適切である

#### REQ-11: 統合テストの修正

- **要件**: 既存の統合テストが正常に動作することを確認する
- **詳細**:
  - ID生成が正常に動作することを確認
  - 既存のテストケースが正常に動作する
- **受け入れ基準**:
  - 既存の統合テストが全て正常に動作する

#### REQ-15: シャーディング規則の動作確認テスト

- **要件**: シャーディング規則が正しく動作することを確認するテストを実装する
- **詳細**:
  - `dm_users`のIDと`dm_posts`の`user_id`が同じ値の場合、同じテーブル番号になることを確認
  - `dm_users`を作成し、そのIDを使用して`dm_posts`を作成した場合、同じテーブル番号のテーブルにデータが入ることを確認
  - テーブル番号の計算ロジック（`id % DBShardingTableCount`）が正しく動作することを確認
- **受け入れ基準**:
  - シャーディング規則の動作確認テストが正常に実行できる
  - テストが全て正常に動作する

### 3.7 API/JavaScript互換性要件

#### REQ-13: JavaScript/API側でのIDの文字列化要件

- **要件**: JavaScriptの制限により、sonyflakeで生成したIDは文字列として扱う
- **詳細**:
  - 全モデルのIDフィールドに`json:"id,string"`タグが設定されていることを確認
  - APIレスポンスでIDが文字列として返されることを確認
  - JavaScript側で数値精度の問題が発生しないことを保証
  - sonyflakeで生成されるIDは64ビット整数であり、JavaScriptのNumber型の安全な整数範囲（2^53-1）を超える可能性があるため、文字列として扱う必要がある
- **現状確認**:
  - DmUser.ID: `json:"id,string"`（確認済み）
  - DmPost.ID: `json:"id,string"`（確認済み）
  - DmNews.ID: `json:"id,string"`（確認済み）
- **受け入れ基準**:
  - 全モデルのIDフィールドに`json:"id,string"`タグが設定されている
  - APIレスポンスでIDが文字列として返される
  - JavaScript側でIDを文字列として扱える
  - 既存のAPIインターフェースとの互換性が維持されている

### 3.8 ドキュメント要件

#### REQ-12: ドキュメントの更新

- **要件**: 関連ドキュメントを更新する
- **詳細**:
  - ID生成方式の変更を記載
  - sharding規則の記載
  - JavaScript側でのIDの扱い（文字列として扱うこと）を記載
  - **identifier生成ルールを記載**（`docs/Architecture.md`に新しいセクション「Identifier Generation」を追加）:
    - 数値のidentifierが必要な箇所はsonyflake (github.com/sony/sonyflake) を使用する
      - 用途: データベースの主キー（dm_users.id, dm_posts.id, dm_news.idなど）
      - 理由: 分散環境で一意性が保証される、時系列順序が保たれる
    - 文字列のidentifierが必要な箇所はUUIDv7 (github.com/google/uuid) を使用する
      - 用途: 文字列型のIDが必要な場合（APIキー、セッションIDなど）
      - 理由: 時系列順序が保たれる、グローバルに一意
  - マイグレーション手順の記載（必要に応じて）
- **対象ドキュメント**:
  - `docs/Architecture.md`（既存のアーキテクチャドキュメント）
    - 「Database Sharding」セクションの後に「Identifier Generation」セクションを追加
  - その他関連ドキュメント
- **受け入れ基準**:
  - ドキュメントが最新の状態である
  - identifier生成ルールが明確に記載されている
  - 各identifierタイプの用途と理由が説明されている

## 4. 非機能要件

### 4.1 パフォーマンス

- ID生成のオーバーヘッドは最小限に抑える
- sonyflakeは分散環境で効率的に動作する

### 4.2 互換性

- 既存のAPIインターフェースは変更しない
- 既存のデータ構造との互換性を維持する
- JavaScript側でのIDの扱い（文字列として扱う）は既存の実装と互換性がある

### 4.3 保守性

- コードは明確で理解しやすい
- 適切なコメントとドキュメントを提供する
- identifier生成ルールが明確にドキュメント化されている

## 5. 制約事項

### 5.1 技術的制約

- Atlas形式のテーブル定義を使用
- GORMを使用したデータアクセス
- Go言語の標準的な実装パターンに従う

### 5.2 既存システムとの整合性

- 既存のsharding機構との整合性を維持
- 既存のRepositoryインターフェースとの整合性を維持

## 6. リスクと対策

### 6.1 リスク

- 既存データとの互換性問題
- ID生成のパフォーマンス問題
- マイグレーション時のデータ整合性問題

### 6.2 対策

- 既存データの移行は別途検討（本実装の範囲外）
- パフォーマンステストを実施
- マイグレーション手順を明確に定義

## 7. 受け入れ基準サマリー

1. sonyflakeライブラリが正常に導入されている
2. 全テーブル定義が修正されている（dm_news, dm_users_NNN, dm_posts_NNN）
3. 全モデル定義が修正されている（autoIncrementタグの削除）
4. ID生成ユーティリティが実装されている
5. Repository層でID生成が統合されている
6. サンプルデータ生成コマンドでID生成が統合されている
7. sharding規則が定義されている（dm_users_NNN: table_sharding_key = id, dm_posts_NNN: table_sharding_key = user_id）
8. シャーディング規則が正しく動作する（dm_usersとdm_postsが同じテーブル番号に配置される）
9. JavaScript/API側でIDが文字列として扱われる（`json:"id,string"`タグの確認）
10. 単体テストが実装されている
11. 統合テストが正常に動作する
12. シャーディング規則の動作確認テストが正常に動作する
13. ドキュメントが更新されている
14. ビルドが正常に完了する
15. 既存のテストが全て正常に動作する
16. サンプルデータ生成コマンドが正常に実行できる