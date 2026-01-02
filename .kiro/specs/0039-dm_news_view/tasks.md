# dm_news_viewビュー作成実装タスク一覧

## 概要

Issue #78の対応として、Atlasのビュー作成機能を確認し、`dm_news`テーブルから`dm_news_view`ビューを作成する機能の実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、Atlas HCL形式でビューを定義し、マイグレーション生成・適用、CloudBeaverでの確認を行う。

## 実装フェーズ

### Phase 1: Atlasビュー定義機能の確認と調査

#### タスク 1.1: Atlasビュー定義の構文確認
**目的**: Atlas HCL形式でビューを定義する方法を確認する

**作業内容**:
- Atlas公式ドキュメントでビュー定義の構文を確認
- HCL形式でのビュー定義の記述方法を確認
- ビュー定義の例を確認
- ビュー定義の注意事項を確認

**受け入れ基準**:
- Atlas公式ドキュメントでビュー定義の構文を確認できている
- HCL形式でのビュー定義の記述方法を理解している
- ビュー定義の例を確認できている
- ビュー定義の注意事項を把握している
- _Requirements: REQ-2_

---

### Phase 2: ビュー定義の追加

#### タスク 2.1: db/schema/master.hclにビュー定義を追加
**目的**: `db/schema/master.hcl`に`dm_news_view`ビューの定義を追加する

**作業内容**:
- `db/schema/master.hcl`を開く
- 既存のテーブル定義（`dm_news`テーブル）の後にビュー定義を追加
- ビュー定義の記述:
  ```hcl
  view "dm_news_view" {
    schema = schema.main
    as     = "SELECT id, title, content, published_at FROM dm_news"
  }
  ```
- ビュー定義の構文が正しいことを確認

**受け入れ基準**:
- `db/schema/master.hcl`にビュー定義が追加されている
- ビュー定義の構文が正しい
- ビューに含まれるカラムが正しい（`id`, `title`, `content`, `published_at`）
- 既存のテーブル定義を変更していない
- _Requirements: REQ-1_

---

#### タスク 2.2: Atlasスキーマ検証の実行
**目的**: 追加したビュー定義が正しいことをAtlasスキーマ検証で確認する

**作業内容**:
- `atlas schema validate`コマンドを実行
- スキーマ検証の結果を確認
- エラーが発生した場合は、エラーメッセージを確認し、ビュー定義を修正
- スキーマ検証が正常に完了するまで繰り返す

**受け入れ基準**:
- `atlas schema validate`コマンドが正常に実行できる
- スキーマ検証が正常に完了する（エラーが発生しない）
- ビュー定義の構文が正しいことが確認できている
- _Requirements: REQ-1_

---

### Phase 3: マイグレーション生成と適用

#### タスク 3.1: マイグレーションファイルの生成
**目的**: `dm_news_view`ビューを作成するマイグレーションファイルを生成する

**作業内容**:
- `atlas migrate diff`コマンドを実行
  ```bash
  atlas migrate diff \
    --dir "file://db/migrations/master" \
    --to "file://db/schema/master.hcl" \
    --dev-url "sqlite://file?mode=memory"
  ```
- マイグレーションファイルが生成されることを確認
- マイグレーションファイルの内容を確認
- マイグレーションファイルに`CREATE VIEW`文が含まれていることを確認
- マイグレーションファイルが`db/migrations/master/`ディレクトリに配置されていることを確認

**受け入れ基準**:
- マイグレーションファイルが正常に生成される
- マイグレーションファイルの内容が正しい（`CREATE VIEW`文が含まれている）
- マイグレーションファイルが`db/migrations/master/`ディレクトリに配置されている
- マイグレーションファイルの命名規則に従っている
- _Requirements: REQ-3_

---

#### タスク 3.2: マイグレーションの適用
**目的**: 生成したマイグレーションファイルを適用してビューを作成する

**作業内容**:
- `atlas migrate apply`コマンドを実行
  ```bash
  atlas migrate apply \
    --dir "file://db/migrations/master" \
    --url "sqlite://server/data/master.db"
  ```
- マイグレーションが正常に適用されることを確認
- マイグレーション履歴が正常に記録されることを確認
- エラーが発生した場合は、エラーメッセージを確認し、対処する

**受け入れ基準**:
- マイグレーションが正常に適用される
- マイグレーション履歴が正常に記録される
- エラーが発生しない
- _Requirements: REQ-4_

---

### Phase 4: ビューの動作確認

#### タスク 4.1: ビュー存在確認
**目的**: データベースに`dm_news_view`ビューが作成されていることを確認する

**作業内容**:
- SQLiteコマンドまたはSQLクエリでビューの存在を確認
- 確認SQL:
  ```sql
  SELECT name FROM sqlite_master WHERE type='view' AND name='dm_news_view';
  ```
- ビューが存在することを確認

**受け入れ基準**:
- ビューがデータベースに存在する
- ビュー名が`dm_news_view`である
- _Requirements: REQ-6_

---

#### タスク 4.2: ビュー構造確認
**目的**: ビューのカラムが正しく定義されていることを確認する

**作業内容**:
- SQLiteコマンドまたはSQLクエリでビューの構造を確認
- 確認SQL:
  ```sql
  PRAGMA table_info(dm_news_view);
  ```
- ビューのカラム（`id`, `title`, `content`, `published_at`）が正しく定義されていることを確認
- ビューに含まれないカラム（`author_id`, `created_at`, `updated_at`）が除外されていることを確認

**受け入れ基準**:
- ビューのカラムが正しく定義されている（`id`, `title`, `content`, `published_at`）
- ビューに含まれないカラムが除外されている（`author_id`, `created_at`, `updated_at`）
- _Requirements: REQ-6_

---

#### タスク 4.3: ビューデータ確認
**目的**: ビューからデータを取得できることを確認する

**作業内容**:
- SQLクエリでビューからデータを取得
- 確認SQL:
  ```sql
  SELECT * FROM dm_news_view LIMIT 10;
  ```
- データが正しく取得できることを確認
- データが`dm_news`テーブルのデータと一致することを確認

**受け入れ基準**:
- ビューからデータを取得できる
- データが正しく表示される
- データが`dm_news`テーブルのデータと一致する
- _Requirements: REQ-6_

---

#### タスク 4.4: ビューデータ整合性確認
**目的**: ビューのデータがベーステーブルのデータと一致することを確認する

**作業内容**:
- SQLクエリでビューとベーステーブルのデータ件数を比較
- 確認SQL:
  ```sql
  SELECT COUNT(*) FROM dm_news_view;
  SELECT COUNT(*) FROM dm_news;
  ```
- データ件数が一致することを確認
- 必要に応じて、個別のレコードを比較して整合性を確認

**受け入れ基準**:
- ビューのデータ件数がベーステーブルのデータ件数と一致する
- ビューのデータがベーステーブルのデータと一致する
- _Requirements: REQ-6_

---

### Phase 5: CloudBeaverでの確認

#### タスク 5.1: CloudBeaver起動と接続確認
**目的**: CloudBeaverを起動し、マスターデータベースに接続できることを確認する

**作業内容**:
- CloudBeaverを起動
  ```bash
  npm run cloudbeaver:start
  ```
- CloudBeaverにアクセス（http://localhost:8978）
- マスターデータベース（`master.db`）に接続
- 接続が正常に確立されることを確認

**受け入れ基準**:
- CloudBeaverが正常に起動する
- CloudBeaverにアクセスできる
- マスターデータベースに接続できる
- _Requirements: REQ-5_

---

#### タスク 5.2: CloudBeaverでビュー表示確認
**目的**: CloudBeaverで`dm_news_view`ビューが表示されることを確認する

**作業内容**:
- CloudBeaverのデータベースツリーでビュー一覧を確認
- `dm_news_view`がビュー一覧に表示されることを確認
- ビューの構造（カラム一覧）が正しく表示されることを確認
- ビューのカラム（`id`, `title`, `content`, `published_at`）が正しく表示されることを確認

**受け入れ基準**:
- CloudBeaverで`dm_news_view`ビューが表示される
- ビューの構造（カラム一覧）が正しく表示される
- ビューのカラム（`id`, `title`, `content`, `published_at`）が正しく表示される
- _Requirements: REQ-5_

---

#### タスク 5.3: CloudBeaverでビューデータ表示確認
**目的**: CloudBeaverでビューのデータが表示されることを確認する

**作業内容**:
- CloudBeaverのUIでビューのデータを表示
- データが正しく表示されることを確認
- データが`dm_news`テーブルのデータと一致することを確認
- ビューのカラム（`id`, `title`, `content`, `published_at`）のみが表示されることを確認

**受け入れ基準**:
- CloudBeaverでビューのデータが表示される
- データが正しく表示される
- データが`dm_news`テーブルのデータと一致する
- ビューのカラム（`id`, `title`, `content`, `published_at`）のみが表示される
- _Requirements: REQ-5_

---

#### タスク 5.4: CloudBeaverでSQLクエリ実行確認
**目的**: CloudBeaverでビューに対してSQLクエリを実行できることを確認する

**作業内容**:
- CloudBeaverのSQLエディタでビューに対してSQLクエリを実行
- 実行SQL:
  ```sql
  SELECT * FROM dm_news_view LIMIT 10;
  ```
- SQLクエリが正常に実行されることを確認
- 結果が正しく表示されることを確認

**受け入れ基準**:
- CloudBeaverでビューに対してSQLクエリを実行できる
- SQLクエリが正常に実行される
- 結果が正しく表示される
- _Requirements: REQ-5_

---

### Phase 6: ドキュメント整備

#### タスク 6.1: ビュー作成手順のドキュメント化
**目的**: Atlasでビューを作成する手順をドキュメント化する

**作業内容**:
- `docs/Atlas.md`を確認（存在しない場合は作成）
- ビュー作成のセクションを追加
- ビュー作成の基本的な手順を記載
- ビュー定義の構文説明を記載
- マイグレーション生成と適用の手順を記載
- CloudBeaverでの確認手順を記載
- トラブルシューティング情報を記載

**受け入れ基準**:
- `docs/Atlas.md`にビュー作成のセクションが追加されている
- ビュー作成の基本的な手順が記載されている
- ビュー定義の構文説明が記載されている
- マイグレーション生成と適用の手順が記載されている
- CloudBeaverでの確認手順が記載されている
- トラブルシューティング情報が記載されている
- _Requirements: REQ-2_

---

#### タスク 6.2: ビュー定義の説明
**目的**: `dm_news_view`ビューの目的と構造を説明する

**作業内容**:
- `docs/Atlas.md`または別のドキュメントに`dm_news_view`ビューの説明を追加
- ビューの目的と用途を説明
- ビューに含まれるカラムと除外されるカラムの理由を説明
- ビューの構造を説明

**受け入れ基準**:
- `dm_news_view`ビューの説明が記載されている
- ビューの目的と用途が説明されている
- ビューに含まれるカラムと除外されるカラムの理由が説明されている
- ビューの構造が説明されている
- _Requirements: REQ-2_

---

## 実装順序

1. **Phase 1**: Atlasビュー定義機能の確認と調査
2. **Phase 2**: ビュー定義の追加
3. **Phase 3**: マイグレーション生成と適用
4. **Phase 4**: ビューの動作確認
5. **Phase 5**: CloudBeaverでの確認
6. **Phase 6**: ドキュメント整備

## 注意事項

- 各フェーズのタスクを順番に実行する
- 各タスクの受け入れ基準を満たしてから次のタスクに進む
- エラーが発生した場合は、エラーメッセージを確認し、適切に対処する
- 既存のスキーマ定義やマイグレーション履歴を変更しない
- ビューは読み取り専用として使用する（プログラムから参照しない）
