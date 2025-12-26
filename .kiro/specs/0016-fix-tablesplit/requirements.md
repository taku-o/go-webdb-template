# テーブル分割修正要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #32
- **Issueタイトル**: shardingの各データベースにテーブル番号000から031までの全てのテーブルが作られている
- **Feature名**: 0016-fix-tablesplit
- **作成日**: 2025-01-27

### 1.2 目的
シャーディングデータベースごとに適切なテーブル範囲のみが作成されるよう、スキーマ定義ファイルをテーブルごとに分割し、マイグレーション管理システムを修正する。
現在、全4つのシャーディングデータベースに全32テーブル（000-031）が作成されている問題を解決し、各データベースに分割されたテーブルのみが作成されるようにする。

### 1.3 スコープ
- スキーマ定義ディレクトリの分割（データベースごとに個別のディレクトリを作成）
- スキーマ定義ファイルの分割（テーブルごとにファイルを分割：`_schema.hcl`, `users.hcl`, `posts.hcl`）
- マイグレーションディレクトリの分割（データベースごとに個別のディレクトリを作成）
- マイグレーション適用スクリプトの修正（データベースごとのマイグレーションディレクトリを参照するように修正）
- Atlas設定ファイルの修正（データベースごとのスキーマディレクトリとマイグレーションディレクトリを参照するように修正）
- 既存のマイグレーション履歴のクリーンアップ（必要に応じて）

## 2. 背景・現状分析

### 2.1 現在の実装
- **スキーマ定義ファイル**: `db/schema/sharding.hcl`（全32テーブルの定義を含む単一ファイル）
- **マイグレーションディレクトリ**: `db/migrations/sharding/`（全データベース共通）
- **マイグレーションファイル**: `db/migrations/sharding/20251226074934_initial.sql`（全32テーブルのCREATE文を含む）
- **マイグレーション適用スクリプト**: `scripts/migrate.sh`（全4つのデータベースに同じマイグレーションファイルを適用）
- **Atlas設定ファイル**: `config/{env}/atlas.hcl`（全4つのデータベースが同じスキーマファイルとマイグレーションディレクトリを参照）

### 2.2 課題点
1. **全テーブルの作成**: 各シャーディングデータベース（sharding_db_1.db ～ sharding_db_4.db）に、テーブル番号000から031までの全てのテーブルが作成されている
2. **期待される動作との不一致**: 
   - 期待: sharding_db_1.db → posts_000-007, users_000-007
   - 期待: sharding_db_2.db → posts_008-015, users_008-015
   - 期待: sharding_db_3.db → posts_016-023, users_016-023
   - 期待: sharding_db_4.db → posts_024-031, users_024-031
   - 実際: 全データベースに posts_000-031, users_000-031 が作成されている
3. **マイグレーション管理の不適切さ**: データベースごとに異なるテーブル範囲が必要なのに、全データベースに同じスキーマファイルとマイグレーションファイルを適用している
4. **スキーマファイルの保守性**: テーブル数が増えた際に、単一ファイルでの管理が困難になる

### 2.3 本実装による改善点
1. **適切なテーブル分割**: 各データベースに必要なテーブルのみが作成される
2. **スキーマ管理の明確化**: データベースごとに個別のスキーマディレクトリとファイルを管理
3. **保守性の向上**: テーブルごとにファイルを分割することで、将来のテーブル追加時に管理が容易になる
4. **拡張性の向上**: 新しいテーブルタイプを追加する際も、既存ファイルに影響を与えずに追加可能

## 3. 機能要件

### 3.1 スキーマ定義ディレクトリの分割

#### 3.1.1 ディレクトリ構造の変更
- **現状**: `db/schema/sharding.hcl`（全データベース共通の単一ファイル）
- **変更後**: 
  ```
  db/schema/
  ├── sharding_1/             # シャード1用のディレクトリ
  │   ├── _schema.hcl         # スキーマ定義（schema "main" {} のみ）
  │   ├── users.hcl           # users_000 〜 users_007
  │   └── posts.hcl           # posts_000 〜 posts_007
  ├── sharding_2/             # シャード2用のディレクトリ
  │   ├── _schema.hcl         # スキーマ定義（schema "main" {} のみ）
  │   ├── users.hcl           # users_008 〜 users_015
  │   └── posts.hcl           # posts_008 〜 posts_015
  ├── sharding_3/             # シャード3用のディレクトリ
  │   ├── _schema.hcl         # スキーマ定義（schema "main" {} のみ）
  │   ├── users.hcl           # users_016 〜 users_023
  │   └── posts.hcl           # posts_016 〜 posts_023
  └── sharding_4/             # シャード4用のディレクトリ
      ├── _schema.hcl         # スキーマ定義（schema "main" {} のみ）
      ├── users.hcl           # users_024 〜 users_031
      └── posts.hcl           # posts_024 〜 posts_031
  ```

#### 3.1.2 ディレクトリの作成
- 各データベース用のスキーマディレクトリを作成
- 既存の`db/schema/sharding.hcl`は削除

### 3.2 スキーマ定義ファイルの分割

#### 3.2.1 テーブル範囲の定義
各データベースに作成するテーブル範囲を定義：

| データベース | テーブル番号範囲 | テーブル例 |
|------------|----------------|----------|
| sharding_db_1.db | _000 〜 _007 | users_000, users_001, ..., users_007, posts_000, posts_001, ..., posts_007 |
| sharding_db_2.db | _008 〜 _015 | users_008, users_009, ..., users_015, posts_008, posts_009, ..., posts_015 |
| sharding_db_3.db | _016 〜 _023 | users_016, users_017, ..., users_023, posts_016, posts_017, ..., posts_023 |
| sharding_db_4.db | _024 〜 _031 | users_024, users_025, ..., users_031, posts_024, posts_025, ..., posts_031 |

#### 3.2.2 スキーマファイルの構造

**`_schema.hcl`**: スキーマ定義のみ
```hcl
schema "main" {
}
```

**`users.hcl`**: usersテーブルの定義（該当するテーブル範囲のみ）
- sharding_1: users_000 〜 users_007
- sharding_2: users_008 〜 users_015
- sharding_3: users_016 〜 users_023
- sharding_4: users_024 〜 users_031

**`posts.hcl`**: postsテーブルの定義（該当するテーブル範囲のみ）
- sharding_1: posts_000 〜 posts_007
- sharding_2: posts_008 〜 posts_015
- sharding_3: posts_016 〜 posts_023
- sharding_4: posts_024 〜 posts_031

#### 3.2.3 Atlasの複数ファイル読み込み
- Atlasは、ディレクトリを指定すると、そのディレクトリ内のすべての`.hcl`ファイルを自動的に読み込む
- `src = "file://db/schema/sharding_1"`のようにディレクトリを指定することで、`_schema.hcl`, `users.hcl`, `posts.hcl`が自動的に結合される

### 3.3 マイグレーションディレクトリの分割

#### 3.3.1 ディレクトリ構造の変更
- **現状**: `db/migrations/sharding/`（全データベース共通）
- **変更後**: 
  - `db/migrations/sharding_1/`（sharding_db_1.db用）
  - `db/migrations/sharding_2/`（sharding_db_2.db用）
  - `db/migrations/sharding_3/`（sharding_db_3.db用）
  - `db/migrations/sharding_4/`（sharding_db_4.db用）

#### 3.3.2 マイグレーションファイルの生成
- 各データベース用のスキーマディレクトリから、Atlasでマイグレーションを生成
- 各マイグレーションディレクトリに、適切なテーブル範囲のみを含むマイグレーションファイルが作成される

### 3.4 マイグレーション適用スクリプトの修正

#### 3.4.1 scripts/migrate.shの修正
- 各データベースごとに個別のマイグレーションディレクトリを参照するように修正
- 修正前: `--dir "file://$PROJECT_ROOT/db/migrations/sharding"`
- 修正後: `--dir "file://$PROJECT_ROOT/db/migrations/sharding_${db_id}"`

#### 3.4.2 スクリプトの動作確認
- 修正後のスクリプトが正常に動作することを確認
- 各データベースに適切なテーブルのみが作成されることを確認

### 3.5 Atlas設定ファイルの修正

#### 3.5.1 config/{env}/atlas.hclの修正
- 各データベース用の環境設定で、個別のスキーマディレクトリとマイグレーションディレクトリを参照するように修正
- 修正前: 
  - `src = "file://db/schema/sharding.hcl"`
  - `dir = "file://db/migrations/sharding"`
- 修正後: 
  - `env "sharding_1"`: 
    - `src = "file://db/schema/sharding_1"`
    - `dir = "file://db/migrations/sharding_1"`
  - `env "sharding_2"`: 
    - `src = "file://db/schema/sharding_2"`
    - `dir = "file://db/migrations/sharding_2"`
  - `env "sharding_3"`: 
    - `src = "file://db/schema/sharding_3"`
    - `dir = "file://db/migrations/sharding_3"`
  - `env "sharding_4"`: 
    - `src = "file://db/schema/sharding_4"`
    - `dir = "file://db/migrations/sharding_4"`

#### 3.5.2 環境別設定ファイルの修正
- `config/develop/atlas.hcl`
- `config/staging/atlas.hcl`
- `config/production/atlas.hcl`
- 全ての環境設定ファイルを修正

### 3.6 既存データベースのクリーンアップ

#### 3.6.1 既存データベースの削除
- 既存のシャーディングデータベース（sharding_db_1.db ～ sharding_db_4.db）を削除
- 修正後のマイグレーションを適用して、正しいテーブル構造で再作成
- **注意**: 既存データの移行は不要。データベースをリセットして再作成する

## 4. 非機能要件

### 4.1 マイグレーション履歴の整合性
- Atlasのマイグレーション履歴（`atlas.sum`）が各データベースごとに正しく管理されること
- 既存のマイグレーション履歴との整合性を保つこと

### 4.2 後方互換性
- 既存のマイグレーション適用スクリプトやAtlas設定ファイルの参照を更新すること
- 他のツールやドキュメントが既存のパスを参照している場合は、影響を確認すること

### 4.3 ドキュメントの更新
- `docs/atlas-operations.md`に新しいディレクトリ構造とマイグレーション適用方法を記載
- `README.md`に変更内容を追記（必要に応じて）

### 4.4 拡張性
- 新しいテーブルタイプ（例: `comments`）を追加する際は、各データベースのスキーマディレクトリに`comments.hcl`を追加するだけで対応可能
- 既存の`users.hcl`や`posts.hcl`に影響を与えずに拡張できる

## 5. 制約事項

### 5.1 既存システムとの関係
- **Atlas設定**: 既存のAtlas設定ファイルを修正する必要がある
- **マイグレーションスクリプト**: 既存のマイグレーション適用スクリプトを修正する必要がある
- **スキーマ定義**: 既存のスキーマ定義ファイル（`db/schema/sharding.hcl`）を分割する必要がある

### 5.2 データベースの再作成
- 既存のシャーディングデータベースを削除して再作成する必要がある
- 既存データの移行は不要。データベースをリセットして再作成する

### 5.3 マイグレーション履歴
- 既存のマイグレーション履歴（`atlas.sum`）を各データベースごとに再生成する必要がある
- マイグレーション履歴の整合性を保つこと

## 6. 受け入れ基準

### 6.1 スキーマ定義ディレクトリの分割
- [ ] `db/schema/sharding_1/`ディレクトリが作成されている
- [ ] `db/schema/sharding_2/`ディレクトリが作成されている
- [ ] `db/schema/sharding_3/`ディレクトリが作成されている
- [ ] `db/schema/sharding_4/`ディレクトリが作成されている
- [ ] 各ディレクトリに`_schema.hcl`, `users.hcl`, `posts.hcl`が存在する
- [ ] 既存の`db/schema/sharding.hcl`が削除されている

### 6.2 スキーマ定義ファイルの内容
- [ ] `db/schema/sharding_1/users.hcl`にusers_000-007のみが含まれている
- [ ] `db/schema/sharding_1/posts.hcl`にposts_000-007のみが含まれている
- [ ] `db/schema/sharding_2/users.hcl`にusers_008-015のみが含まれている
- [ ] `db/schema/sharding_2/posts.hcl`にposts_008-015のみが含まれている
- [ ] `db/schema/sharding_3/users.hcl`にusers_016-023のみが含まれている
- [ ] `db/schema/sharding_3/posts.hcl`にposts_016-023のみが含まれている
- [ ] `db/schema/sharding_4/users.hcl`にusers_024-031のみが含まれている
- [ ] `db/schema/sharding_4/posts.hcl`にposts_024-031のみが含まれている

### 6.3 マイグレーションディレクトリの分割
- [ ] `db/migrations/sharding_1/`ディレクトリが作成されている
- [ ] `db/migrations/sharding_2/`ディレクトリが作成されている
- [ ] `db/migrations/sharding_3/`ディレクトリが作成されている
- [ ] `db/migrations/sharding_4/`ディレクトリが作成されている
- [ ] 既存の`db/migrations/sharding/`ディレクトリが削除されている

### 6.4 マイグレーションファイルの生成
- [ ] 各データベース用のスキーマディレクトリからマイグレーションが生成されている
- [ ] `db/migrations/sharding_1/`にテーブル000-007のみを含むマイグレーションファイルが存在する
- [ ] `db/migrations/sharding_2/`にテーブル008-015のみを含むマイグレーションファイルが存在する
- [ ] `db/migrations/sharding_3/`にテーブル016-023のみを含むマイグレーションファイルが存在する
- [ ] `db/migrations/sharding_4/`にテーブル024-031のみを含むマイグレーションファイルが存在する
- [ ] 各マイグレーションディレクトリに`atlas.sum`が生成されている

### 6.5 マイグレーション適用スクリプトの修正
- [ ] `scripts/migrate.sh`が各データベースごとのマイグレーションディレクトリを参照するように修正されている
- [ ] 修正後のスクリプトが正常に動作する

### 6.6 Atlas設定ファイルの修正
- [ ] `config/develop/atlas.hcl`が各データベースごとのスキーマディレクトリとマイグレーションディレクトリを参照するように修正されている
- [ ] `config/staging/atlas.hcl`が各データベースごとのスキーマディレクトリとマイグレーションディレクトリを参照するように修正されている
- [ ] `config/production/atlas.hcl`が各データベースごとのスキーマディレクトリとマイグレーションディレクトリを参照するように修正されている

### 6.7 データベースのテーブル構造
- [ ] `sharding_db_1.db`にテーブル000-007のみが作成されている
- [ ] `sharding_db_2.db`にテーブル008-015のみが作成されている
- [ ] `sharding_db_3.db`にテーブル016-023のみが作成されている
- [ ] `sharding_db_4.db`にテーブル024-031のみが作成されている
- [ ] 各データベースに不要なテーブル（範囲外のテーブル）が作成されていない

### 6.8 ドキュメントの更新
- [ ] `docs/atlas-operations.md`に新しいディレクトリ構造とマイグレーション適用方法が記載されている
- [ ] 変更内容が適切にドキュメント化されている

## 7. 影響範囲

### 7.1 新規追加が必要なディレクトリ・ファイル

#### ディレクトリ
- `db/schema/sharding_1/`: sharding_db_1.db用スキーマ
- `db/schema/sharding_2/`: sharding_db_2.db用スキーマ
- `db/schema/sharding_3/`: sharding_db_3.db用スキーマ
- `db/schema/sharding_4/`: sharding_db_4.db用スキーマ
- `db/migrations/sharding_1/`: sharding_db_1.db用マイグレーション
- `db/migrations/sharding_2/`: sharding_db_2.db用マイグレーション
- `db/migrations/sharding_3/`: sharding_db_3.db用マイグレーション
- `db/migrations/sharding_4/`: sharding_db_4.db用マイグレーション

#### ファイル
- `db/schema/sharding_1/_schema.hcl`: スキーマ定義
- `db/schema/sharding_1/users.hcl`: users_000-007の定義
- `db/schema/sharding_1/posts.hcl`: posts_000-007の定義
- `db/schema/sharding_2/_schema.hcl`: スキーマ定義
- `db/schema/sharding_2/users.hcl`: users_008-015の定義
- `db/schema/sharding_2/posts.hcl`: posts_008-015の定義
- `db/schema/sharding_3/_schema.hcl`: スキーマ定義
- `db/schema/sharding_3/users.hcl`: users_016-023の定義
- `db/schema/sharding_3/posts.hcl`: posts_016-023の定義
- `db/schema/sharding_4/_schema.hcl`: スキーマ定義
- `db/schema/sharding_4/users.hcl`: users_024-031の定義
- `db/schema/sharding_4/posts.hcl`: posts_024-031の定義
- 各マイグレーションディレクトリのマイグレーションファイルと`atlas.sum`

### 7.2 変更が必要なファイル

#### スクリプト
- `scripts/migrate.sh`: 各データベースごとのマイグレーションディレクトリを参照するように修正

#### 設定ファイル
- `config/develop/atlas.hcl`: 各データベースごとのスキーマディレクトリとマイグレーションディレクトリを参照するように修正
- `config/staging/atlas.hcl`: 各データベースごとのスキーマディレクトリとマイグレーションディレクトリを参照するように修正
- `config/production/atlas.hcl`: 各データベースごとのスキーマディレクトリとマイグレーションディレクトリを参照するように修正

#### ドキュメント
- `docs/atlas-operations.md`: 新しいディレクトリ構造とマイグレーション適用方法を記載

### 7.3 削除が必要なファイル

#### ファイル・ディレクトリ
- `db/schema/sharding.hcl`: 既存の共通スキーマファイル（削除）
- `db/migrations/sharding/`: 既存の共通マイグレーションディレクトリ（削除）

## 8. 実装上の注意事項

### 8.1 スキーマファイルの分割
- 既存の`db/schema/sharding.hcl`から、各データベース用のスキーマファイルを抽出する際は、テーブル範囲を正確に確認すること
- 各スキーマファイルに必要なテーブルのみが含まれていることを確認すること
- `_schema.hcl`には`schema "main" {}`のみを含めること

### 8.2 マイグレーションファイルの生成
- 各データベース用のスキーマディレクトリから、Atlasでマイグレーションを生成すること
- マイグレーション生成時は、適切なスキーマディレクトリとマイグレーションディレクトリを指定すること

### 8.3 マイグレーション履歴の管理
- 各データベースごとに`atlas.sum`を生成する必要がある
- 既存のマイグレーション履歴との整合性を保つこと

### 8.4 データベースの再作成
- 既存のシャーディングデータベースを削除して再作成する
- 既存データの移行は不要。データベースをリセットして再作成する

### 8.5 テストの実施
- 修正後のマイグレーション適用スクリプトが正常に動作することを確認すること
- 各データベースに適切なテーブルのみが作成されることを確認すること
- 不要なテーブルが作成されていないことを確認すること
- Atlas設定ファイルが正しく各スキーマディレクトリを参照していることを確認すること

### 8.6 拡張時の注意
- 新しいテーブルタイプを追加する際は、各データベースのスキーマディレクトリに新しい`.hcl`ファイルを追加する
- 既存のファイル（`users.hcl`, `posts.hcl`）は変更しない

## 9. 参考情報

### 9.1 関連Issue
- GitHub Issue #32: shardingの各データベースにテーブル番号000から031までの全てのテーブルが作られている

### 9.2 既存ドキュメント
- `docs/atlas-operations.md`: Atlas運用ガイド
- `docs/Sharding.md`: シャーディングの詳細仕様
- `.kiro/specs/0012-sharding/`: シャーディング実装の仕様書
- `.kiro/specs/0014-db-atlas/`: Atlas導入の仕様書

### 9.3 技術スタック
- **Atlas**: データベーススキーマ管理ツール（複数HCLファイルの自動読み込みをサポート）
- **SQLite**: 開発環境のデータベース
- **Bash**: マイグレーション適用スクリプト

