# テーブル分割修正実装タスク一覧

## 概要
テーブル分割修正の実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: スキーマファイルの分析と分割準備

#### - [ ] タスク 1.1: 既存スキーマファイルの分析
**目的**: 既存の`db/schema/sharding.hcl`を分析し、各データベース用のテーブル範囲を特定

**作業内容**:
- `db/schema/sharding.hcl`を読み込み、全32テーブル（users_000-031, posts_000-031）の定義を確認
- 各データベースに割り当てるテーブル範囲を確認:
  - sharding_db_1.db: users_000-007, posts_000-007
  - sharding_db_2.db: users_008-015, posts_008-015
  - sharding_db_3.db: users_016-023, posts_016-023
  - sharding_db_4.db: users_024-031, posts_024-031
- テーブル定義の構造（カラム、インデックス、主キー）を確認

**受け入れ基準**:
- 各データベースに割り当てるテーブル範囲が明確になっている
- テーブル定義の構造が理解できている

---

### Phase 2: スキーマディレクトリとファイルの作成

#### - [ ] タスク 2.1: スキーマディレクトリの作成
**目的**: 各データベース用のスキーマディレクトリを作成

**作業内容**:
- `db/schema/sharding_1/`ディレクトリを作成
- `db/schema/sharding_2/`ディレクトリを作成
- `db/schema/sharding_3/`ディレクトリを作成
- `db/schema/sharding_4/`ディレクトリを作成

**受け入れ基準**:
- `db/schema/sharding_1/`ディレクトリが存在する
- `db/schema/sharding_2/`ディレクトリが存在する
- `db/schema/sharding_3/`ディレクトリが存在する
- `db/schema/sharding_4/`ディレクトリが存在する

---

#### - [ ] タスク 2.2: _schema.hclファイルの作成
**目的**: 各データベース用のスキーマ定義ファイルを作成

**作業内容**:
- `db/schema/sharding_1/_schema.hcl`を作成（`schema "main" {}`のみ）
- `db/schema/sharding_2/_schema.hcl`を作成（`schema "main" {}`のみ）
- `db/schema/sharding_3/_schema.hcl`を作成（`schema "main" {}`のみ）
- `db/schema/sharding_4/_schema.hcl`を作成（`schema "main" {}`のみ）

**受け入れ基準**:
- 各ディレクトリに`_schema.hcl`が存在する
- 各ファイルに`schema "main" {}`のみが含まれている

---

#### - [ ] タスク 2.3: users.hclファイルの作成
**目的**: 各データベース用のusersテーブル定義ファイルを作成

**作業内容**:
- `db/schema/sharding_1/users.hcl`を作成（users_000-007の定義）
- `db/schema/sharding_2/users.hcl`を作成（users_008-015の定義）
- `db/schema/sharding_3/users.hcl`を作成（users_016-023の定義）
- `db/schema/sharding_4/users.hcl`を作成（users_024-031の定義）
- 既存の`db/schema/sharding.hcl`から該当するテーブル定義を抽出
- foreign_key定義は含めない（postsテーブルとの参照関係は削除）

**受け入れ基準**:
- `db/schema/sharding_1/users.hcl`にusers_000-007のみが含まれている
- `db/schema/sharding_2/users.hcl`にusers_008-015のみが含まれている
- `db/schema/sharding_3/users.hcl`にusers_016-023のみが含まれている
- `db/schema/sharding_4/users.hcl`にusers_024-031のみが含まれている
- 各ファイルに必要なカラム、インデックス、主キーが正しく定義されている

---

#### - [ ] タスク 2.4: posts.hclファイルの作成
**目的**: 各データベース用のpostsテーブル定義ファイルを作成

**作業内容**:
- `db/schema/sharding_1/posts.hcl`を作成（posts_000-007の定義）
- `db/schema/sharding_2/posts.hcl`を作成（posts_008-015の定義）
- `db/schema/sharding_3/posts.hcl`を作成（posts_016-023の定義）
- `db/schema/sharding_4/posts.hcl`を作成（posts_024-031の定義）
- 既存の`db/schema/sharding.hcl`から該当するテーブル定義を抽出
- foreign_key定義は含めない（usersテーブルとの参照関係は削除）

**受け入れ基準**:
- `db/schema/sharding_1/posts.hcl`にposts_000-007のみが含まれている
- `db/schema/sharding_2/posts.hcl`にposts_008-015のみが含まれている
- `db/schema/sharding_3/posts.hcl`にposts_016-023のみが含まれている
- `db/schema/sharding_4/posts.hcl`にposts_024-031のみが含まれている
- 各ファイルに必要なカラム、インデックス、主キーが正しく定義されている
- foreign_key定義が含まれていない

---

#### - [ ] タスク 2.5: 既存スキーマファイルの削除
**目的**: 既存の`db/schema/sharding.hcl`を削除

**作業内容**:
- `db/schema/sharding.hcl`を削除
- 削除前に、新しいスキーマファイルが正しく作成されていることを確認

**受け入れ基準**:
- `db/schema/sharding.hcl`が削除されている
- 各データベース用のスキーマディレクトリとファイルが正しく作成されている

---

### Phase 3: マイグレーションディレクトリの分割とマイグレーションファイルの生成

#### - [ ] タスク 3.1: マイグレーションディレクトリの作成
**目的**: 各データベース用のマイグレーションディレクトリを作成

**作業内容**:
- `db/migrations/sharding_1/`ディレクトリを作成
- `db/migrations/sharding_2/`ディレクトリを作成
- `db/migrations/sharding_3/`ディレクトリを作成
- `db/migrations/sharding_4/`ディレクトリを作成

**受け入れ基準**:
- `db/migrations/sharding_1/`ディレクトリが存在する
- `db/migrations/sharding_2/`ディレクトリが存在する
- `db/migrations/sharding_3/`ディレクトリが存在する
- `db/migrations/sharding_4/`ディレクトリが存在する

---

#### - [ ] タスク 3.2: sharding_db_1.db用マイグレーションファイルの生成
**目的**: sharding_db_1.db用のマイグレーションファイルを生成

**作業内容**:
- Atlasでマイグレーションを生成:
  ```bash
  atlas migrate diff initial \
    --dir file://db/migrations/sharding_1 \
    --to file://db/schema/sharding_1 \
    --dev-url "sqlite://file?mode=memory"
  ```
- 生成されたマイグレーションファイルにテーブル000-007のみが含まれていることを確認
- `atlas.sum`が生成されていることを確認

**受け入れ基準**:
- `db/migrations/sharding_1/`にマイグレーションファイルが存在する
- マイグレーションファイルにテーブル000-007のみが含まれている
- `atlas.sum`が生成されている

---

#### - [ ] タスク 3.3: sharding_db_2.db用マイグレーションファイルの生成
**目的**: sharding_db_2.db用のマイグレーションファイルを生成

**作業内容**:
- Atlasでマイグレーションを生成:
  ```bash
  atlas migrate diff initial \
    --dir file://db/migrations/sharding_2 \
    --to file://db/schema/sharding_2 \
    --dev-url "sqlite://file?mode=memory"
  ```
- 生成されたマイグレーションファイルにテーブル008-015のみが含まれていることを確認
- `atlas.sum`が生成されていることを確認

**受け入れ基準**:
- `db/migrations/sharding_2/`にマイグレーションファイルが存在する
- マイグレーションファイルにテーブル008-015のみが含まれている
- `atlas.sum`が生成されている

---

#### - [ ] タスク 3.4: sharding_db_3.db用マイグレーションファイルの生成
**目的**: sharding_db_3.db用のマイグレーションファイルを生成

**作業内容**:
- Atlasでマイグレーションを生成:
  ```bash
  atlas migrate diff initial \
    --dir file://db/migrations/sharding_3 \
    --to file://db/schema/sharding_3 \
    --dev-url "sqlite://file?mode=memory"
  ```
- 生成されたマイグレーションファイルにテーブル016-023のみが含まれていることを確認
- `atlas.sum`が生成されていることを確認

**受け入れ基準**:
- `db/migrations/sharding_3/`にマイグレーションファイルが存在する
- マイグレーションファイルにテーブル016-023のみが含まれている
- `atlas.sum`が生成されている

---

#### - [ ] タスク 3.5: sharding_db_4.db用マイグレーションファイルの生成
**目的**: sharding_db_4.db用のマイグレーションファイルを生成

**作業内容**:
- Atlasでマイグレーションを生成:
  ```bash
  atlas migrate diff initial \
    --dir file://db/migrations/sharding_4 \
    --to file://db/schema/sharding_4 \
    --dev-url "sqlite://file?mode=memory"
  ```
- 生成されたマイグレーションファイルにテーブル024-031のみが含まれていることを確認
- `atlas.sum`が生成されていることを確認

**受け入れ基準**:
- `db/migrations/sharding_4/`にマイグレーションファイルが存在する
- マイグレーションファイルにテーブル024-031のみが含まれている
- `atlas.sum`が生成されている

---

#### - [ ] タスク 3.6: 既存マイグレーションディレクトリの削除
**目的**: 既存の`db/migrations/sharding/`ディレクトリを削除

**作業内容**:
- `db/migrations/sharding/`ディレクトリを削除
- 削除前に、新しいマイグレーションディレクトリとファイルが正しく作成されていることを確認

**受け入れ基準**:
- `db/migrations/sharding/`ディレクトリが削除されている
- 各データベース用のマイグレーションディレクトリとファイルが正しく作成されている

---

### Phase 4: マイグレーション適用スクリプトの修正

#### - [ ] タスク 4.1: scripts/migrate.shの修正
**目的**: マイグレーション適用スクリプトを修正し、各データベースごとのマイグレーションディレクトリを参照するように変更

**作業内容**:
- `scripts/migrate.sh`の`migrate_sharding()`関数を修正
- 修正前: `--dir "file://$PROJECT_ROOT/db/migrations/sharding"`
- 修正後: `--dir "file://$PROJECT_ROOT/db/migrations/sharding_${db_id}"`
- 修正後のスクリプトが正常に動作することを確認

**受け入れ基準**:
- `scripts/migrate.sh`が各データベースごとのマイグレーションディレクトリを参照するように修正されている
- 修正後のスクリプトが正常に動作する

---

### Phase 5: Atlas設定ファイルの修正

#### - [ ] タスク 5.1: config/develop/atlas.hclの修正
**目的**: 開発環境用Atlas設定ファイルを修正し、各データベースごとのスキーマディレクトリとマイグレーションディレクトリを参照するように変更

**作業内容**:
- `config/develop/atlas.hcl`の各`env "sharding_*"`設定を修正
- 修正前: 
  - `src = "file://db/schema/sharding.hcl"`
  - `dir = "file://db/migrations/sharding"`
- 修正後:
  - `env "sharding_1"`: `src = "file://db/schema/sharding_1"`, `dir = "file://db/migrations/sharding_1"`
  - `env "sharding_2"`: `src = "file://db/schema/sharding_2"`, `dir = "file://db/migrations/sharding_2"`
  - `env "sharding_3"`: `src = "file://db/schema/sharding_3"`, `dir = "file://db/migrations/sharding_3"`
  - `env "sharding_4"`: `src = "file://db/schema/sharding_4"`, `dir = "file://db/migrations/sharding_4"`
- 修正後の設定ファイルが正しく動作することを確認

**受け入れ基準**:
- `config/develop/atlas.hcl`が各データベースごとのスキーマディレクトリとマイグレーションディレクトリを参照するように修正されている
- 修正後の設定ファイルが正しく動作する

---

#### - [ ] タスク 5.2: config/staging/atlas.hclの修正
**目的**: ステージング環境用Atlas設定ファイルを修正し、各データベースごとのスキーマディレクトリとマイグレーションディレクトリを参照するように変更

**作業内容**:
- `config/staging/atlas.hcl`の各`env "sharding_*"`設定を修正
- 修正内容は`config/develop/atlas.hcl`と同様
- データベースURLはステージング環境用のものを使用

**受け入れ基準**:
- `config/staging/atlas.hcl`が各データベースごとのスキーマディレクトリとマイグレーションディレクトリを参照するように修正されている
- 修正後の設定ファイルが正しく動作する

---

#### - [ ] タスク 5.3: config/production/atlas.hclの修正
**目的**: 本番環境用Atlas設定ファイルを修正し、各データベースごとのスキーマディレクトリとマイグレーションディレクトリを参照するように変更

**作業内容**:
- `config/production/atlas.hcl`の各`env "sharding_*"`設定を修正
- 修正内容は`config/develop/atlas.hcl`と同様
- データベースURLは本番環境用のものを使用

**受け入れ基準**:
- `config/production/atlas.hcl`が各データベースごとのスキーマディレクトリとマイグレーションディレクトリを参照するように修正されている
- 修正後の設定ファイルが正しく動作する

---

### Phase 6: 既存データベースのクリーンアップとマイグレーション適用（開発環境のみ）

#### - [ ] タスク 6.1: 開発環境での既存データベース削除
**目的**: 開発環境で既存のシャーディングデータベースファイルを削除

**作業内容**:
- 開発環境で既存のシャーディングデータベースファイルを削除:
  ```bash
  rm server/data/sharding_db_1.db
  rm server/data/sharding_db_2.db
  rm server/data/sharding_db_3.db
  rm server/data/sharding_db_4.db
  ```
- 削除前に、バックアップが必要な場合は取得

**受け入れ基準**:
- 開発環境で既存のシャーディングデータベースファイルが削除されている

---

#### - [ ] タスク 6.2: 修正後のマイグレーション適用
**目的**: 修正後のマイグレーションを適用して、正しいテーブル構造でデータベースを再作成

**作業内容**:
- `./scripts/migrate.sh sharding`を実行
- 各データベースに適切なテーブルのみが作成されることを確認
- 不要なテーブル（範囲外のテーブル）が作成されていないことを確認

**受け入れ基準**:
- マイグレーションが正常に適用される
- `sharding_db_1.db`にテーブル000-007のみが作成されている
- `sharding_db_2.db`にテーブル008-015のみが作成されている
- `sharding_db_3.db`にテーブル016-023のみが作成されている
- `sharding_db_4.db`にテーブル024-031のみが作成されている
- 各データベースに不要なテーブル（範囲外のテーブル）が作成されていない

---

### Phase 7: ドキュメントの更新

#### - [ ] タスク 7.1: docs/atlas-operations.mdの更新
**目的**: Atlas運用ガイドに新しいディレクトリ構造とマイグレーション適用方法を記載

**作業内容**:
- `docs/atlas-operations.md`を更新
- 新しいディレクトリ構造を記載:
  - `db/schema/sharding_1/` 〜 `sharding_4/`の構造
  - `db/migrations/sharding_1/` 〜 `sharding_4/`の構造
- マイグレーション生成方法を更新:
  - 各データベース用のスキーマディレクトリからマイグレーションを生成する方法
- マイグレーション適用方法を更新:
  - 各データベースごとのマイグレーションディレクトリを参照する方法
- Atlas設定ファイルの修正内容を記載

**受け入れ基準**:
- `docs/atlas-operations.md`に新しいディレクトリ構造が記載されている
- マイグレーション生成方法が更新されている
- マイグレーション適用方法が更新されている
- Atlas設定ファイルの修正内容が記載されている

---

## 受け入れ基準の確認

### 全体の受け入れ基準
- [ ] 各データベースに適切なテーブル範囲のみが作成されている
- [ ] 各データベースに不要なテーブル（範囲外のテーブル）が作成されていない
- [ ] マイグレーション適用スクリプトが正常に動作する
- [ ] Atlas設定ファイルが正しく各スキーマディレクトリを参照している
- [ ] ドキュメントが適切に更新されている

---

## 注意事項

### 実装時の注意
- スキーマファイルの分割時は、テーブル範囲を正確に確認すること
- マイグレーションファイルの生成時は、適切なスキーマディレクトリとマイグレーションディレクトリを指定すること
- 既存データベースの削除は開発環境のみで実施すること
- ステージング・本番環境では、適切な手順に従って対応すること

### テスト時の注意
- 各データベースに適切なテーブルのみが作成されていることを確認すること
- 不要なテーブルが作成されていないことを確認すること
- マイグレーション適用スクリプトが正常に動作することを確認すること
- Atlas設定ファイルが正しく各スキーマディレクトリを参照していることを確認すること

