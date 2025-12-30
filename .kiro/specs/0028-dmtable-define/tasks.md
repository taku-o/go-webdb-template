# 分散テーブル環境対応テーブル設計修正実装タスク一覧

## 概要

Issue #52の対応として、分散テーブル環境に適したテーブル設計に修正する機能の実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、auto_incrementによるID生成からsonyflakeによるID生成に変更し、分散環境で一意性が保証されるID生成方式を実装する。

## 実装フェーズ

### Phase 1: sonyflakeライブラリの導入とID生成ユーティリティの実装

#### タスク 1.1: sonyflakeライブラリの導入 ✓
**目的**: sonyflakeライブラリをプロジェクトに導入し、依存関係を解決する

**作業内容**:
- go.modにgithub.com/sony/sonyflakeの依存関係を追加
- go mod tidyを実行して依存関係を解決
- ビルドが正常に完了することを確認
- コンパイルエラーの確認

**受け入れ基準**:
- go.modにsonyflakeの依存関係が追加されている
- go mod tidyが正常に実行できる
- ビルドが正常に完了する
- コンパイルエラーがない
- _Requirements: REQ-1_

---

#### タスク 1.2: ID生成ユーティリティの実装 ✓
**目的**: sonyflakeを使用したID生成ユーティリティを実装する

**作業内容**:
- server/internal/util/idgen/sonyflake.goを作成
- sonyflakeインスタンスの初期化と管理機能を実装
- GenerateSonyflakeID()関数を実装（スレッドセーフなID生成）
- エラーハンドリングを実装
- 生成されるIDがint64型であることを確認
- コンパイルエラーの確認

**受け入れ基準**:
- ID生成関数が正常に動作する
- 生成されるIDが一意である
- エラーが適切に処理される
- コンパイルエラーがない
- _Requirements: REQ-7_

---

#### タスク 1.3: ID生成ユーティリティの単体テスト実装 ✓
**目的**: ID生成ユーティリティの動作を確認する単体テストを実装する

**作業内容**:
- server/internal/util/idgen/sonyflake_test.goを作成
- ID生成関数の正常動作を確認するテスト
- 生成されるIDの一意性を確認するテスト（複数回生成して重複がないことを確認）
- エラーハンドリングの動作を確認するテスト
- テストの実行

**受け入れ基準**:
- テストが正常に実行できる
- カバレッジが適切である
- 全てのテストが通過する
- _Requirements: REQ-10_

---

### Phase 2: テーブル定義の修正（Atlas形式）

#### タスク 2.1: dm_newsテーブル定義の修正 ✓
**目的**: master.hclのdm_newsテーブルのidカラムを修正する

**作業内容**:
- db/schema/master.hclのdm_newsテーブル定義を修正
- idカラムの型をintegerからbigintに変更
- auto_incrementをfalseに設定
- idカラムに`unsigned = true`を追加（Atlas形式: `column "id" { type = bigint, unsigned = true }`）
- Atlasスキーマ検証を実行

**受け入れ基準**:
- idカラムの型がbigintに変更されている
- auto_incrementがfalseに設定されている
- idカラムに`unsigned = true`が設定されている
- Atlasスキーマ検証が正常に完了する
- マイグレーションスクリプトが正常に生成される
- _Requirements: REQ-2_

---

#### タスク 2.2: dm_users_NNNテーブル定義の修正 ✓
**目的**: 全shardingグループのdm_users.hclのidカラムを修正する

**作業内容**:
- db/schema/sharding_1/dm_users.hclを修正
- db/schema/sharding_2/dm_users.hclを修正
- db/schema/sharding_3/dm_users.hclを修正
- db/schema/sharding_4/dm_users.hclを修正
- idカラムの型をintegerからbigintに変更
- auto_incrementをfalseに明示的に設定
- idカラムに`unsigned = true`を追加（Atlas形式: `column "id" { type = bigint, unsigned = true }`）
- Atlasスキーマ検証を実行

**受け入れ基準**:
- 全shardingグループのテーブル定義が修正されている
- idカラムの型がbigintに変更されている
- auto_incrementがfalseに設定されている
- idカラムに`unsigned = true`が設定されている
- Atlasスキーマ検証が正常に完了する
- _Requirements: REQ-3_

---

#### タスク 2.3: dm_posts_NNNテーブル定義の修正 ✓
**目的**: 全shardingグループのdm_posts.hclのidカラムを修正する

**作業内容**:
- db/schema/sharding_1/dm_posts.hclを修正
- db/schema/sharding_2/dm_posts.hclを修正
- db/schema/sharding_3/dm_posts.hclを修正
- db/schema/sharding_4/dm_posts.hclを修正
- idカラムの型をintegerからbigintに変更
- auto_incrementをfalseに明示的に設定
- idカラムに`unsigned = true`を追加（Atlas形式: `column "id" { type = bigint, unsigned = true }`）
- Atlasスキーマ検証を実行

**受け入れ基準**:
- 全shardingグループのテーブル定義が修正されている
- idカラムの型がbigintに変更されている
- auto_incrementがfalseに設定されている
- idカラムに`unsigned = true`が設定されている
- Atlasスキーマ検証が正常に完了する
- _Requirements: REQ-4_

---

#### タスク 2.4: 初期データ用SQLの移行 ✓
**目的**: AtlasでマイグレーションSQLを作成し直した際に、既存の初期データ用SQLを新しいファイルに移行する

**作業内容**:
- db/migrations/master/20251229111855_initial_schema.sqlの下の方に書いてある初期データ用のSQLを確認
- 初期データ用SQLの内容を確認（GoAdminの初期データ: ロール、ユーザー、メニュー、権限、関連テーブルなど）
- AtlasでマイグレーションSQLを作成し直した後、既存のマイグレーションファイルが上書きされることを確認
- 初期データ用SQLを新しいマイグレーションファイル（例: `db/migrations/master/YYYYMMDDHHMMSS_seed_data.sql`）に移行
- 移行したSQLが正常に実行できることを確認
- 既存のマイグレーションファイルから初期データ用SQLが削除されていることを確認

**注意事項**:
- 既存のデータの維持は考えなくて良い（マイグレーション時のデータ移行は不要）
- ただし、初期データ用SQLの移行を忘れてはいけない
- 初期データ用SQLは、マイグレーション実行時に必要なデータ（GoAdminの管理者ユーザーなど）を含む

**受け入れ基準**:
- 初期データ用SQLが新しいマイグレーションファイルに移行されている
- 移行したSQLが正常に実行できる
- 既存のマイグレーションファイルから初期データ用SQLが削除されている
- マイグレーション実行時に初期データが正常に投入される

---

### Phase 3: モデル定義の修正

#### タスク 3.1: DmNewsモデルの修正 ✓
**目的**: DmNewsモデルのautoIncrementタグを削除する

**作業内容**:
- server/internal/model/dm_news.goを修正
- IDフィールドのgormタグからautoIncrementを削除
- primaryKeyタグのみを残す
- コンパイルエラーの確認

**受け入れ基準**:
- autoIncrementタグが削除されている
- primaryKeyタグのみが残っている
- モデルファイルが正常にコンパイルできる
- GORMがautoIncrementを使用しない
- _Requirements: REQ-5_

---

#### タスク 3.2: ID型の確認と統一 ✓
**目的**: 全モデルのIDフィールドの型とJSONタグを確認する

**作業内容**:
- server/internal/model/dm_user.goのIDフィールドを確認
- server/internal/model/dm_post.goのIDフィールドを確認
- server/internal/model/dm_news.goのIDフィールドを確認
- 全モデルのIDフィールドがint64型であることを確認
- 全モデルのIDフィールドにjson:"id,string"タグが設定されていることを確認
- JavaScript側でIDが文字列として扱われることを確認

**受け入れ基準**:
- 全モデルのIDフィールドがint64型である
- 全モデルのIDフィールドにjson:"id,string"タグが設定されている
- APIレスポンスでIDが文字列として返される
- JavaScript側でIDを文字列として扱える
- 既存のAPIインターフェースとの互換性が維持されている
- _Requirements: REQ-6, REQ-13_

---

### Phase 4: Repository層でのID生成統合

#### タスク 4.1: DmUserRepositoryGORMの修正 ✓
**目的**: DmUserRepositoryGORMのCreateメソッドでID生成を実装する

**作業内容**:
- server/internal/repository/dm_user_repository_gorm.goを修正
- Createメソッドでidgen.GenerateSonyflakeID()を呼び出し
- 生成したIDをモデルのIDフィールドに設定
- エラーハンドリングを実装
- 既存のテストを実行

**受け入れ基準**:
- 新規作成時にIDが自動生成される
- 生成されたIDが一意である
- 既存のテストが正常に動作する
- エラーが適切に処理される
- _Requirements: REQ-8_

---

#### タスク 4.2: DmPostRepositoryGORMの修正 ✓
**目的**: DmPostRepositoryGORMのCreateメソッドでID生成を実装する

**作業内容**:
- server/internal/repository/dm_post_repository_gorm.goを修正
- Createメソッドでidgen.GenerateSonyflakeID()を呼び出し
- 生成したIDをモデルのIDフィールドに設定
- エラーハンドリングを実装
- 既存のテストを実行

**受け入れ基準**:
- 新規作成時にIDが自動生成される
- 生成されたIDが一意である
- 既存のテストが正常に動作する
- エラーが適切に処理される
- _Requirements: REQ-8_

---

#### タスク 4.3: database/sql版のRepositoryの修正 ✓
**目的**: database/sql版のRepositoryのCreateメソッドでID生成を実装する

**作業内容**:
- server/internal/repository/dm_user_repository.goを修正
- server/internal/repository/dm_post_repository.goを修正
- Createメソッドでidgen.GenerateSonyflakeID()を呼び出し
- 生成したIDをモデルのIDフィールドに設定
- INSERT文にidカラムを追加
- エラーハンドリングを実装
- 既存のテストを実行

**受け入れ基準**:
- 新規作成時にIDが自動生成される
- 生成されたIDが一意である
- INSERT文にidカラムが含まれている
- 既存のテストが正常に動作する
- エラーが適切に処理される
- _Requirements: REQ-8_

---

### Phase 5: サンプルデータ生成コマンドの修正

#### タスク 5.1: generateDmUsers関数の修正 ✓
**目的**: generateDmUsers関数でID生成を統合する

**作業内容**:
- server/cmd/generate-sample-data/main.goのgenerateDmUsers関数を修正
- DmUser作成時にidgen.GenerateSonyflakeID()を呼び出し
- 生成したIDをモデルのIDフィールドに設定
- insertDmUsersBatch関数でINSERT文にidカラムを追加
- 生成したIDを含めてバッチ挿入を実行
- コンパイルエラーの確認

**受け入れ基準**:
- DmUser作成時にIDが生成される
- INSERT文にidカラムが含まれている
- バッチ挿入が正常に動作する
- 生成されたIDが一意である
- コンパイルエラーがない
- _Requirements: REQ-14_

---

#### タスク 5.2: generateDmPosts関数の修正 ✓
**目的**: generateDmPosts関数でID生成を統合する

**作業内容**:
- server/cmd/generate-sample-data/main.goのgenerateDmPosts関数を修正
- DmPost作成時にidgen.GenerateSonyflakeID()を呼び出し
- 生成したIDをモデルのIDフィールドに設定
- insertDmPostsBatch関数でINSERT文にidカラムを追加
- 生成したIDを含めてバッチ挿入を実行
- コンパイルエラーの確認

**受け入れ基準**:
- DmPost作成時にIDが生成される
- INSERT文にidカラムが含まれている
- バッチ挿入が正常に動作する
- 生成されたIDが一意である
- コンパイルエラーがない
- _Requirements: REQ-14_

---

#### タスク 5.3: generateDmNews関数の修正 ✓
**目的**: generateDmNews関数でID生成を統合する

**作業内容**:
- server/cmd/generate-sample-data/main.goのgenerateDmNews関数を修正
- DmNews作成時にidgen.GenerateSonyflakeID()を呼び出し
- 生成したIDをモデルのIDフィールドに設定
- GORMのCreateInBatchesを使用している場合、IDが設定されていることを確認
- コンパイルエラーの確認

**受け入れ基準**:
- DmNews作成時にIDが生成される
- IDがモデルに設定されている
- バッチ挿入が正常に動作する
- 生成されたIDが一意である
- コンパイルエラーがない
- _Requirements: REQ-14_

---

#### タスク 5.4: fetchDmUserIDs関数の削除または変更 ✓
**目的**: 挿入後にIDを取得する必要がなくなるため、fetchDmUserIDs関数を削除または変更する

**作業内容**:
- server/cmd/generate-sample-data/main.goのfetchDmUserIDs関数を確認
- 関数が使用されている箇所を確認
- 削除または変更を検討（使用されていない場合は削除、使用されている場合は変更）
- コマンドが正常に実行できることを確認
- 生成されたIDが一意であることを確認

**受け入れ基準**:
- fetchDmUserIDs関数が削除または変更されている
- コマンドが正常に実行できる
- 生成されたIDが一意である
- コンパイルエラーがない
- _Requirements: REQ-14_

---

### Phase 6: Sharding規則の定義

#### タスク 6.1: sharding規則の実装 ✓
**目的**: server/internal/db/sharding.goにテーブルsharding規則を定義する

**作業内容**:
- server/internal/db/sharding.goを修正
- dm_users_NNNのtable_sharding_keyとしてidを定義
- dm_posts_NNNのtable_sharding_keyとしてuser_idを定義
- テーブル番号計算ロジック（id % DBShardingTableCount、user_id % DBShardingTableCount）を実装
- あるdm_usersに紐付いたdm_postsが同じテーブル番号になることを保証する規則を実装
- ドキュメントコメントを追加
- コンパイルエラーの確認

**受け入れ基準**:
- sharding規則が明確に定義されている
- コードから規則が理解できる
- dm_usersのIDとdm_postsのuser_idが同じ値の場合、同じテーブル番号になることが保証されている
- ドキュメントコメントが追加されている
- コンパイルエラーがない
- _Requirements: REQ-9_

---

### Phase 7: 統合テストの修正と動作確認

#### タスク 7.1: 既存統合テストの修正 ✓
**目的**: 既存の統合テストが正常に動作することを確認する

**作業内容**:
- 既存の統合テストを実行
- ID生成が正常に動作することを確認
- 既存のテストケースが正常に動作することを確認
- テストが全て正常に実行できることを確認
- 必要に応じてテストコードを修正

**受け入れ基準**:
- ID生成が正常に動作する
- 既存の統合テストが全て正常に動作する
- テストが全て正常に実行できる
- _Requirements: REQ-11_

---

#### タスク 7.2: シャーディング規則の動作確認テスト実装 ✓
**目的**: シャーディング規則が正しく動作することを確認するテストを実装する

**作業内容**:
- シャーディング規則の動作確認テストを実装
- dm_usersのIDとdm_postsのuser_idが同じ値の場合、同じテーブル番号になることを確認するテスト
- dm_usersを作成し、そのIDを使用してdm_postsを作成した場合、同じテーブル番号のテーブルにデータが入ることを確認するテスト
- テーブル番号の計算ロジック（id % DBShardingTableCount）が正しく動作することを確認するテスト
- テストの実行

**受け入れ基準**:
- シャーディング規則の動作確認テストが正常に実行できる
- テストが全て正常に動作する
- dm_usersとdm_postsが同じテーブル番号に配置されることが確認できる
- _Requirements: REQ-15_

---

### Phase 8: ドキュメントの更新

#### タスク 8.1: Architecture.mdの更新 ✓
**目的**: 関連ドキュメントを更新し、ID生成方式の変更とidentifier生成ルールを記載する

**作業内容**:
- docs/Architecture.mdを修正
- 「Database Sharding」セクションの後に「Identifier Generation」セクションを追加
- 数値のidentifierが必要な箇所はsonyflake (github.com/sony/sonyflake) を使用することを記載
- 文字列のidentifierが必要な箇所はUUIDv7 (github.com/google/uuid) を使用することを記載
- 各identifierタイプの用途と理由を説明
- ID生成方式の変更を記載
- sharding規則の記載を更新
- JavaScript側でのIDの扱い（文字列として扱うこと）を記載

**受け入れ基準**:
- ドキュメントが最新の状態である
- identifier生成ルールが明確に記載されている
- 各identifierタイプの用途と理由が説明されている
- ID生成方式の変更が記載されている
- sharding規則が記載されている
- JavaScript側でのIDの扱いが記載されている
- _Requirements: REQ-12_
