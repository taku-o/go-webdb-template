# dm_newsテーブルidカラム定義変更要件定義書

## 1. 概要

### 1.1 プロジェクト情報

- **プロジェクト名**: go-webdb-template
- **Issue番号**: #59
- **Issueタイトル**: newsテーブルのidカラムの定義を変更する
- **Feature名**: 0029-dm_news_id
- **作成日**: 2025-01-27

### 1.2 目的

dm_newsテーブルのidカラムの定義を`integer auto_increment = true`に変更する。0028-dmtable-defineで`bigint unsigned auto_increment = false`に変更され、sonyflakeによるID生成が実装されているが、auto_increment方式に戻す必要がある。

### 1.3 スコープ

- テーブル定義の修正（Atlas形式）
- dm_newsテーブルのidカラム定義の変更
- モデル定義の修正（GORMタグの修正）
- ID生成ロジックの修正（sonyflakeからauto_incrementへ）
- サンプルデータ生成コマンドの修正（server/cmd/generate-sample-data/）
- GoAdmin管理画面の修正（server/internal/admin/tables.go）
- マイグレーションSQLファイルの生成

**本実装の範囲外**:

- 既存データの移行（Issue記載により既存データは維持しなくて良い）
- 他のテーブル（dm_users、dm_postsなど）の修正
- sonyflakeライブラリの削除（他のテーブルで使用されているため）

## 2. 背景・現状分析

### 2.1 現状の問題

0028-dmtable-defineで、分散テーブル環境対応のためdm_newsテーブルのidカラムが`bigint unsigned auto_increment = false`に変更され、sonyflakeによるID生成が実装された。しかし、dm_newsテーブルはmasterグループに配置される単一テーブルであり、分散環境の制約を受けないため、auto_increment方式に戻す必要がある。

### 2.2 現状のテーブル定義

#### 2.2.1 dm_newsテーブル（master.hcl）

- 現在: `id bigint unsigned auto_increment = false`
- 修正後: `id integer auto_increment = true`

### 2.3 現状のモデル定義

#### 2.3.1 DmNewsモデル（server/internal/model/dm_news.go）

- 現在: `gorm:"primaryKey"`（autoIncrementタグなし）
- 修正後: `gorm:"primaryKey;autoIncrement"`
- ID型: `int64`（integerでもint64で問題ない）

### 2.4 現状のID生成

#### 2.4.1 サンプルデータ生成コマンド（server/cmd/generate-sample-data/main.go）

- 現在: `generateDmNews`関数でsonyflakeを使用してIDを生成し、モデルに設定
- 修正後: IDを設定せず、GORMのautoIncrementに任せる

#### 2.4.2 GoAdmin管理画面（server/internal/admin/tables.go）

- 現在: `GetDmNewsTable`関数のフォーム設定で、新規作成時にsonyflakeでIDを生成
- 修正後: ID生成ロジックを削除し、GORMのautoIncrementに任せる

## 3. 要件定義

### 3.1 テーブル定義修正要件

#### REQ-1: dm_newsテーブル定義の修正

- **要件**: `db/schema/master.hcl`のdm_newsテーブルのidカラムを修正する
- **詳細**:
  - `type = bigint` → `type = integer`
  - `unsigned = true` → 削除（integerはunsignedをサポートしない）
  - `auto_increment = false` → `auto_increment = true`
- **対象ファイル**:
  - `db/schema/master.hcl`
- **受け入れ基準**:
  - Atlasスキーマ検証が正常に完了する
  - マイグレーションスクリプトが正常に生成される
  - idカラムが`integer auto_increment = true`になっている

### 3.2 モデル定義修正要件

#### REQ-2: DmNewsモデルの修正

- **要件**: `server/internal/model/dm_news.go`のDmNews構造体のIDフィールドのGORMタグを修正する
- **詳細**:
  - `gorm:"primaryKey"` → `gorm:"primaryKey;autoIncrement"`
  - autoIncrementタグを追加
  - ID型は`int64`のまま（integerでもint64で問題ない）
- **対象ファイル**:
  - `server/internal/model/dm_news.go`
- **受け入れ基準**:
  - モデルファイルが正常にコンパイルできる
  - GORMがautoIncrementを使用する
  - IDフィールドに`autoIncrement`タグが設定されている

### 3.3 ID生成ロジック修正要件

#### REQ-3: サンプルデータ生成コマンドの修正

- **要件**: `server/cmd/generate-sample-data/main.go`の`generateDmNews`関数を修正する
- **詳細**:
  - sonyflakeによるID生成を削除
  - `idgen.GenerateSonyflakeID()`の呼び出しを削除
  - モデルのIDフィールドに値を設定しない（GORMのautoIncrementに任せる）
  - `idgen`パッケージのインポートが不要になった場合は削除（他の関数で使用されていない場合）
- **対象ファイル**:
  - `server/cmd/generate-sample-data/main.go`
- **受け入れ基準**:
  - サンプルデータ生成時にIDが自動生成される
  - sonyflakeによるID生成コードが削除されている
  - コマンドが正常に実行できる
  - 生成されたIDが連番になっている

#### REQ-4: GoAdmin管理画面の修正

- **要件**: `server/internal/admin/tables.go`の`GetDmNewsTable`関数を修正する
- **詳細**:
  - フォーム設定の`FieldPostFilterFn`でsonyflakeによるID生成を削除
  - `idgen.GenerateSonyflakeID()`の呼び出しを削除
  - IDフィールドの設定を簡素化（GORMのautoIncrementに任せる）
  - `idgen`パッケージのインポートが不要になった場合は削除（他の関数で使用されていない場合）
- **対象ファイル**:
  - `server/internal/admin/tables.go`
- **受け入れ基準**:
  - GoAdmin管理画面で新規作成時にIDが自動生成される
  - sonyflakeによるID生成コードが削除されている
  - 管理画面が正常に動作する
  - 新規作成時にIDが連番になっている

### 3.4 マイグレーション要件

#### REQ-5: マイグレーションSQLファイルの生成

- **要件**: AtlasでマイグレーションSQLファイルを生成する
- **詳細**:
  - `atlas migrate diff`コマンドでマイグレーションSQLファイルを生成
  - 既存データは維持しなくて良い（Issue記載）
  - マイグレーションSQLファイルは`db/migrations/master/`に配置
- **受け入れ基準**:
  - マイグレーションSQLファイルが正常に生成される
  - マイグレーションが正常に適用できる
  - idカラムが`INTEGER AUTO_INCREMENT`になっている

### 3.5 テスト要件

#### REQ-6: 既存テストの確認

- **要件**: 既存のテストが正常に動作することを確認する
- **詳細**:
  - 既存の単体テストが正常に動作する
  - 既存の統合テストが正常に動作する
  - サンプルデータ生成コマンドのテストが正常に動作する
- **受け入れ基準**:
  - 既存のテストが全て正常に動作する
  - テストエラーが発生しない

## 4. 非機能要件

### 4.1 互換性

- 既存のAPIインターフェースは変更しない
- 既存のデータ構造との互換性を維持する（ただし、既存データは維持しなくて良い）
- JavaScript側でのIDの扱い（文字列として扱う）は既存の実装と互換性がある

### 4.2 保守性

- コードは明確で理解しやすい
- 適切なコメントとドキュメントを提供する
- 不要になったコード（sonyflake関連）は削除する

## 5. 制約事項

### 5.1 技術的制約

- Atlas形式のテーブル定義を使用
- GORMを使用したデータアクセス
- Go言語の標準的な実装パターンに従う
- integer型はunsignedをサポートしない

### 5.2 既存システムとの整合性

- 既存のRepositoryインターフェースとの整合性を維持
- 既存のAPIインターフェースとの整合性を維持
- 他のテーブル（dm_users、dm_posts）は変更しない

## 6. リスクと対策

### 6.1 リスク

- 既存データとの互換性問題（Issue記載により既存データは維持しなくて良いため、リスクは低い）
- マイグレーション時のデータ整合性問題（既存データは維持しなくて良いため、リスクは低い）
- integer型の範囲制限（int64型を使用しているが、integer型に変更するため、範囲制限に注意が必要）

### 6.2 対策

- 既存データの移行は不要（Issue記載）
- マイグレーション手順を明確に定義
- integer型の範囲制限を確認（通常は問題ないが、大量のデータがある場合は注意）

## 7. 受け入れ基準サマリー

1. テーブル定義が`integer auto_increment = true`になっている
2. モデル定義に`autoIncrement`タグが追加されている
3. サンプルデータ生成コマンドでIDが自動生成される（sonyflakeによるID生成が削除されている）
4. GoAdmin管理画面でIDが自動生成される（sonyflakeによるID生成が削除されている）
5. マイグレーションSQLファイルが生成されている
6. 既存のテストが全て正常に動作する
7. ビルドが正常に完了する
8. サンプルデータ生成コマンドが正常に実行できる
9. GoAdmin管理画面が正常に動作する

