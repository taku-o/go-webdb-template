# dm_newsテーブルidカラム定義変更設計書

## 1. 概要

### 1.1 設計の目的

要件定義書に基づき、Issue #59の対応として、dm_newsテーブルのidカラムの定義を`integer auto_increment = true`に変更する機能の詳細設計を定義する。0028-dmtable-defineで実装されたsonyflakeによるID生成を削除し、GORMのauto_increment機能に戻す。

### 1.2 設計の範囲

- テーブル定義の修正（Atlas形式）
- モデル定義の修正（GORMタグの修正）
- サンプルデータ生成コマンドの修正（sonyflake削除）
- GoAdmin管理画面の修正（sonyflake削除）
- マイグレーションSQLファイルの生成
- 既存テストの確認

### 1.3 設計方針

- **シンプルな実装**: auto_increment機能を活用し、ID生成ロジックを簡素化する
- **既存ロジックの維持**: ID生成方式のみを変更し、既存のビジネスロジックは維持する
- **段階的実装**: 各コンポーネントを段階的に実装し、各段階でテストを実行する
- **互換性の維持**: 既存のAPIインターフェースは変更しない
- **不要コードの削除**: sonyflake関連のコードを削除し、コードベースをクリーンに保つ

## 2. アーキテクチャ設計

### 2.1 ID生成システムの変更

#### 2.1.1 変更前（0028-dmtable-define）

```
┌─────────────────────────────────────────────────────────┐
│              ID生成システム（変更前）                    │
│                                                         │
│  ┌──────────────────────────────────────────────────┐  │
│  │  idgen.GenerateSonyflakeID()                    │  │
│  │  - sonyflakeインスタンスの管理                  │  │
│  │  - スレッドセーフなID生成                       │  │
│  └──────────────────────────────────────────────────┘  │
│                    │                                     │
│                    ▼                                     │
│  ┌──────────────────────────────────────────────────┐  │
│  │  github.com/sony/sonyflake                       │  │
│  │  - 分散環境で一意性が保証されるID生成             │  │
│  └──────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
                    │
        ┌───────────┴───────────┐
        ▼                       ▼
┌──────────────┐      ┌──────────────┐
│ サンプルデータ│      │ GoAdmin管理   │
│ 生成コマンド   │      │ 画面          │
│ generateDmNews│      │ GetDmNewsTable│
└──────────────┘      └──────────────┘
```

#### 2.1.2 変更後（本実装）

```
┌─────────────────────────────────────────────────────────┐
│              ID生成システム（変更後）                    │
│                                                         │
│  ┌──────────────────────────────────────────────────┐  │
│  │  GORM auto_increment                             │  │
│  │  - データベースのAUTO_INCREMENT機能を使用         │  │
│  │  - IDを設定せず、データベースに任せる             │  │
│  └──────────────────────────────────────────────────┘  │
│                    │                                     │
│                    ▼                                     │
│  ┌──────────────────────────────────────────────────┐  │
│  │  MySQL/SQLite AUTO_INCREMENT                     │  │
│  │  - データベース側で自動的にIDを生成               │  │
│  └──────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
                    │
        ┌───────────┴───────────┐
        ▼                       ▼
┌──────────────┐      ┌──────────────┐
│ サンプルデータ│      │ GoAdmin管理   │
│ 生成コマンド   │      │ 画面          │
│ generateDmNews│      │ GetDmNewsTable│
│ (ID設定なし)  │      │ (ID設定なし)  │
└──────────────┘      └──────────────┘
```

### 2.2 ID生成フローの変更

#### 2.2.1 サンプルデータ生成コマンドでのID生成（変更前）

```
1. generateDmNews() が呼び出される
   ↓
2. idgen.GenerateSonyflakeID() を呼び出し
   ↓
3. sonyflakeがIDを生成（int64）
   ↓
4. モデルのIDフィールドに設定
   ↓
5. データベースに保存（IDを含めてINSERT）
```

#### 2.2.2 サンプルデータ生成コマンドでのID生成（変更後）

```
1. generateDmNews() が呼び出される
   ↓
2. モデルのIDフィールドに値を設定しない（0のまま）
   ↓
3. データベースに保存（IDを指定せずにINSERT）
   ↓
4. データベースのAUTO_INCREMENTがIDを自動生成
   ↓
5. GORMが生成されたIDをモデルに反映
```

#### 2.2.3 GoAdmin管理画面でのID生成（変更前）

```
1. 新規作成フォームが送信される
   ↓
2. FieldPostFilterFnが呼び出される
   ↓
3. idgen.GenerateSonyflakeID() を呼び出し
   ↓
4. sonyflakeがIDを生成（int64）
   ↓
5. フォーム値としてIDを設定
   ↓
6. データベースに保存
```

#### 2.2.4 GoAdmin管理画面でのID生成（変更後）

```
1. 新規作成フォームが送信される
   ↓
2. IDフィールドは非表示（FieldHide）
   ↓
3. IDを設定しない（データベースに任せる）
   ↓
4. データベースに保存（IDを指定せずにINSERT）
   ↓
5. データベースのAUTO_INCREMENTがIDを自動生成
```

## 3. データ構造設計

### 3.1 テーブル定義の変更

#### 3.1.1 dm_newsテーブル（master.hcl）

**変更前**:
```hcl
table "dm_news" {
  schema = schema.main
  column "id" {
    null           = false
    type           = bigint
    unsigned       = true
    auto_increment = false
  }
  // ... 他のカラム
}
```

**変更後**:
```hcl
table "dm_news" {
  schema = schema.main
  column "id" {
    null           = false
    type           = integer
    auto_increment = true
  }
  // ... 他のカラム
}
```

**変更点**:
- `type = bigint` → `type = integer`
- `unsigned = true` → 削除（integerはunsignedをサポートしない）
- `auto_increment = false` → `auto_increment = true`

### 3.2 モデル定義の変更

#### 3.2.1 DmNewsモデル（server/internal/model/dm_news.go）

**変更前**:
```go
type DmNews struct {
    ID          int64      `json:"id,string" db:"id" gorm:"primaryKey"`
    Title       string     `json:"title" db:"title" gorm:"type:varchar(255);not null"`
    Content     string     `json:"content" db:"content" gorm:"type:text;not null"`
    // ... 他のフィールド
}
```

**変更後**:
```go
type DmNews struct {
    ID          int64      `json:"id,string" db:"id" gorm:"primaryKey;autoIncrement"`
    Title       string     `json:"title" db:"title" gorm:"type:varchar(255);not null"`
    Content     string     `json:"content" db:"content" gorm:"type:text;not null"`
    // ... 他のフィールド
}
```

**変更点**:
- `gorm:"primaryKey"` → `gorm:"primaryKey;autoIncrement"`
- ID型は`int64`のまま（integerでもint64で問題ない）

## 4. 実装設計

### 4.1 テーブル定義の修正

#### 4.1.1 ファイル: `db/schema/master.hcl`

**修正内容**:
- `column "id"`の`type`を`bigint`から`integer`に変更
- `unsigned = true`を削除
- `auto_increment = false`を`auto_increment = true`に変更

**修正後のコード**:
```hcl
table "dm_news" {
  schema = schema.main
  column "id" {
    null           = false
    type           = integer
    auto_increment = true
  }
  // ... 他のカラムは変更なし
}
```

### 4.2 モデル定義の修正

#### 4.2.1 ファイル: `server/internal/model/dm_news.go`

**修正内容**:
- `ID`フィールドのGORMタグに`autoIncrement`を追加

**修正後のコード**:
```go
type DmNews struct {
    ID          int64      `json:"id,string" db:"id" gorm:"primaryKey;autoIncrement"`
    Title       string     `json:"title" db:"title" gorm:"type:varchar(255);not null"`
    Content     string     `json:"content" db:"content" gorm:"type:text;not null"`
    AuthorID    *int64     `json:"author_id,omitempty,string" db:"author_id" gorm:"index:idx_dm_news_author_id"`
    PublishedAt *time.Time `json:"published_at,omitempty" db:"published_at" gorm:"index:idx_dm_news_published_at"`
    CreatedAt   time.Time  `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   time.Time  `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
}
```

### 4.3 サンプルデータ生成コマンドの修正

#### 4.3.1 ファイル: `server/cmd/generate-sample-data/main.go`

**修正内容**:
- `generateDmNews`関数からsonyflakeによるID生成を削除
- モデルのIDフィールドに値を設定しない（0のまま）
- `idgen`パッケージのインポートが不要になった場合は削除（他の関数で使用されていない場合）

**修正前**:
```go
func generateDmNews(groupManager *db.GroupManager, totalCount int) error {
    // ...
    var dmNews []*model.DmNews
    for i := 0; i < totalCount; i++ {
        id, err := idgen.GenerateSonyflakeID()
        if err != nil {
            return fmt.Errorf("failed to generate sonyflake ID: %w", err)
        }

        n := &model.DmNews{
            ID:          id,  // sonyflakeで生成したIDを設定
            Title:       gofakeit.Sentence(5),
            Content:     gofakeit.Paragraph(3, 5, 10, "\n"),
            // ...
        }
        dmNews = append(dmNews, n)
    }
    // ...
}
```

**修正後**:
```go
func generateDmNews(groupManager *db.GroupManager, totalCount int) error {
    // ...
    var dmNews []*model.DmNews
    for i := 0; i < totalCount; i++ {
        authorID := gofakeit.Int64()
        publishedAt := gofakeit.Date()

        n := &model.DmNews{
            // IDは設定しない（auto_incrementに任せる）
            Title:       gofakeit.Sentence(5),
            Content:     gofakeit.Paragraph(3, 5, 10, "\n"),
            AuthorID:    &authorID,
            PublishedAt: &publishedAt,
            CreatedAt:   time.Now(),
            UpdatedAt:   time.Now(),
        }
        dmNews = append(dmNews, n)
    }
    // ...
}
```

**注意事項**:
- `idgen`パッケージのインポートが他の関数（`generateDmUsers`、`generateDmPosts`）で使用されている場合は削除しない
- IDを設定しないことで、GORMがauto_incrementを使用してIDを自動生成する

### 4.4 GoAdmin管理画面の修正

#### 4.4.1 ファイル: `server/internal/admin/tables.go`

**修正内容**:
- `GetDmNewsTable`関数のフォーム設定から`FieldPostFilterFn`を削除
- IDフィールドの設定を簡素化（GORMのautoIncrementに任せる）
- `idgen`パッケージのインポートが不要になった場合は削除（他の関数で使用されていない場合）

**修正前**:
```go
formList.AddField("ID", "id", db.Int, form.Default).
    FieldNotAllowEdit().
    FieldHide().
    FieldPostFilterFn(func(value types.PostFieldModel) interface{} {
        // 新規作成時（値が空）の場合、sonyflakeでIDを生成
        if value.Value.Value() == "" {
            id, err := idgen.GenerateSonyflakeID()
            if err != nil {
                return 0
            }
            return fmt.Sprintf("%d", id)
        }
        return value.Value.Value()
    })
```

**修正後**:
```go
formList.AddField("ID", "id", db.Int, form.Default).
    FieldNotAllowEdit().
    FieldHide()
    // FieldPostFilterFnを削除（GORMのautoIncrementに任せる）
```

**注意事項**:
- `idgen`パッケージのインポートが他の関数で使用されている場合は削除しない
- IDフィールドは非表示のまま（`FieldHide()`）
- 編集不可のまま（`FieldNotAllowEdit()`）

### 4.5 マイグレーションSQLファイルの生成

#### 4.5.1 マイグレーションコマンド

```bash
atlas migrate diff \
  --dir "file://db/migrations/master" \
  --to "file://db/schema/master.hcl" \
  --dev-url "sqlite://file?mode=memory&_fk=1"
```

**生成されるマイグレーションSQL**:
```sql
-- Modify "dm_news" table
ALTER TABLE `dm_news` MODIFY COLUMN `id` INTEGER NOT NULL AUTO_INCREMENT;
```

**注意事項**:
- 既存データは維持しなくて良い（Issue記載）
- マイグレーションSQLファイルは`db/migrations/master/`に配置される
- マイグレーション実行前にデータベースのバックアップを推奨（ただし、既存データは維持しなくて良いため、必須ではない）

## 5. エラーハンドリング

### 5.1 テーブル定義の変更

- **Atlasスキーマ検証エラー**: スキーマ定義の構文エラーを確認し、修正する
- **マイグレーション生成エラー**: 既存のマイグレーションとの整合性を確認する

### 5.2 モデル定義の変更

- **コンパイルエラー**: GORMタグの構文エラーを確認し、修正する
- **型の不一致**: integer型とint64型の互換性を確認する（問題なし）

### 5.3 サンプルデータ生成コマンド

- **ID生成エラー**: sonyflakeによるID生成が削除されているため、ID生成エラーは発生しない
- **データベースエラー**: 既存のエラーハンドリングを維持する

### 5.4 GoAdmin管理画面

- **ID生成エラー**: sonyflakeによるID生成が削除されているため、ID生成エラーは発生しない
- **フォーム送信エラー**: 既存のエラーハンドリングを維持する

## 6. テスト戦略

### 6.1 単体テスト

#### 6.1.1 モデル定義のテスト

- **テスト内容**: DmNewsモデルが正常にコンパイルできることを確認
- **テストファイル**: 既存のテストファイルを確認（必要に応じて追加）

#### 6.1.2 サンプルデータ生成コマンドのテスト

- **テスト内容**: `generateDmNews`関数が正常に動作し、IDが自動生成されることを確認
- **テスト方法**: 実際にデータベースにデータを挿入し、IDが連番になっていることを確認

### 6.2 統合テスト

#### 6.2.1 データベース操作のテスト

- **テスト内容**: dm_newsテーブルへのINSERT操作が正常に動作し、IDが自動生成されることを確認
- **テスト方法**: 既存の統合テストを実行し、正常に動作することを確認

#### 6.2.2 GoAdmin管理画面のテスト

- **テスト内容**: GoAdmin管理画面で新規作成が正常に動作し、IDが自動生成されることを確認
- **テスト方法**: 実際にGoAdmin管理画面で新規作成を実行し、IDが連番になっていることを確認

### 6.3 既存テストの確認

- **テスト内容**: 既存の単体テスト・統合テストが正常に動作することを確認
- **テスト方法**: 既存のテストスイートを実行し、全てのテストが正常に動作することを確認

## 7. 実装順序

### 7.1 フェーズ1: テーブル定義とモデル定義の修正

1. `db/schema/master.hcl`のdm_newsテーブル定義を修正
2. `server/internal/model/dm_news.go`のDmNewsモデル定義を修正
3. コンパイルエラーがないことを確認

### 7.2 フェーズ2: ID生成ロジックの修正

1. `server/cmd/generate-sample-data/main.go`の`generateDmNews`関数を修正
2. `server/internal/admin/tables.go`の`GetDmNewsTable`関数を修正
3. 不要になった`idgen`パッケージのインポートを削除（他の関数で使用されていない場合）
4. コンパイルエラーがないことを確認

### 7.3 フェーズ3: マイグレーションSQLファイルの生成

1. AtlasでマイグレーションSQLファイルを生成
2. マイグレーションSQLファイルの内容を確認
3. マイグレーションが正常に適用できることを確認

### 7.4 フェーズ4: テストの実行

1. 既存の単体テストを実行
2. 既存の統合テストを実行
3. サンプルデータ生成コマンドを実行し、IDが自動生成されることを確認
4. GoAdmin管理画面で新規作成を実行し、IDが自動生成されることを確認

## 8. 注意事項

### 8.1 integer型の範囲制限

- **制約**: integer型は通常32ビット（-2,147,483,648 〜 2,147,483,647）または64ビット（-9,223,372,036,854,775,808 〜 9,223,372,036,854,775,807）の範囲を持つ
- **SQLite**: integer型は64ビット整数をサポート
- **MySQL**: integer型は32ビット整数（-2,147,483,648 〜 2,147,483,647）をサポート
- **注意**: 大量のデータがある場合は、integer型の範囲制限に注意が必要

### 8.2 既存データの扱い

- **Issue記載**: 既存データは維持しなくて良い
- **マイグレーション**: 既存データを削除する必要がある場合は、マイグレーションSQLに`TRUNCATE TABLE`または`DELETE FROM`を追加する

### 8.3 sonyflakeライブラリの扱い

- **削除しない**: sonyflakeライブラリは他のテーブル（dm_users、dm_posts）で使用されているため、削除しない
- **インポートの削除**: `idgen`パッケージのインポートが不要になった場合のみ削除する（他の関数で使用されていない場合）

### 8.4 互換性の維持

- **APIインターフェース**: 既存のAPIインターフェースは変更しない
- **データ構造**: 既存のデータ構造との互換性を維持する（ただし、既存データは維持しなくて良い）
- **JavaScript側**: IDの扱い（文字列として扱う）は既存の実装と互換性がある
