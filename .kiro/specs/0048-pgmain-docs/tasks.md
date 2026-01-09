# ドキュメントの修正実装タスク一覧

## 概要
ドキュメントからSQLite関連の記述を削除し、必要に応じてPostgreSQL関連の記述を追加するためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: SQLite関連記述の検索と確認

#### タスク 1.1: 対象ファイルのリストアップ
**目的**: 修正対象となるドキュメントファイルをリストアップする。

**作業内容**:
- プロジェクトルート直下のMarkdownファイルを確認
- `docs/` ディレクトリ内のすべてのMarkdownファイルを確認
- 対象ファイルのリストを作成

**確認対象ファイル**:
- `README.md`
- `docs/Docker.md`
- `docs/API.md`
- `docs/Admin.md`
- `docs/Architecture.md`
- `docs/Atlas-Operations.md`
- `docs/Command-Line-Tool.md`
- `docs/Database-Viewer.md`
- `docs/File-Upload.md`
- `docs/Generate-Sample-Data.md`
- `docs/Logging.md`
- `docs/Metabase.md`
- `docs/Partner-Idp-Auth0-Login.md`
- `docs/Project-Structure.md`
- `docs/Queue-Job.md`
- `docs/Rate-Limit.md`
- `docs/Release-Check.md`
- `docs/Send-Mail.md`
- `docs/Sharding.md`
- `docs/Testing.md`
- その他のプロジェクトルート直下のMarkdownファイル

**除外ファイル**:
- `docs/Initial-Setup.md`（ユーザー作成ファイルのため対応不要）
- `docs/Spec-Driven-Development.md`（ユーザー作成ファイルのため対応不要）
- `.kiro/specs/` 内のファイル（履歴として残す）
- `cloudbeaver/config/` 内のファイル（ツールの設定）

**受け入れ基準**:
- すべての対象ファイルがリストアップされている
- 除外ファイルが明確に記載されている

- _Requirements: 3.1.1, 7.1_
- _Design: 2.1_

---

#### タスク 1.2: SQLite関連記述の網羅的検索
**目的**: すべての対象ファイルからSQLite関連の記述を検索する。

**作業内容**:
- `grep` コマンドを使用してSQLite関連の記述を検索
- 検索パターン: `SQLite|sqlite|SQLite3|sqlite3`
- 検索対象: `README.md` と `docs/` ディレクトリ内のすべてのMarkdownファイル
- 検索結果をファイルごとに記録

**検索コマンド例**:
```bash
grep -ri "SQLite\|sqlite" README.md docs/ --include="*.md"
```

**受け入れ基準**:
- すべての対象ファイルでSQLite関連の記述が検索されている
- 検索結果がファイルごとに記録されている
- 検索対象外のファイル（`docs/Initial-Setup.md`、`docs/Spec-Driven-Development.md`、`.kiro/specs/`、`cloudbeaver/config/`）は除外されている

- _Requirements: 8.1_
- _Design: 2.2.1_

---

#### タスク 1.3: 検索結果の確認と分類
**目的**: 検索結果を確認し、削除対象の記述を分類する。

**作業内容**:
- 検索結果を確認し、SQLite関連の記述を特定
- 記述の種類を分類:
  - 説明文
  - 設定例
  - 注意書き
  - 比較表の行
  - セットアップ手順
- 各記述の文脈を確認
- 削除対象の記述をリストアップ

**受け入れ基準**:
- すべてのSQLite関連の記述が特定されている
- 記述の種類が分類されている
- 各記述の文脈が確認されている

- _Requirements: 3.1.2_
- _Design: 2.3.1_

---

### Phase 2: SQLite関連記述の削除

#### タスク 2.1: README.mdの修正
**目的**: `README.md` からSQLite関連の記述を削除する。

**作業内容**:
- `README.md` を開く
- SQLite関連の記述を確認
- 以下の記述を削除:
  - SQLiteに関する説明文
  - SQLiteに関する設定例
  - SQLiteに関する注意書き
  - SQLiteに関する比較表の行
- 削除後の文脈が自然になるように調整
- 必要に応じてPostgreSQL関連の記述を追加

**確認項目**:
- プロジェクト概要でのデータベース記述
- セットアップ手順でのデータベース記述
- データベース接続情報の記述
- 機能説明でのデータベース記述

**受け入れ基準**:
- SQLite関連の記述が削除されている
- 削除後の文脈が自然である
- データベースとしてPostgreSQL/MySQLが明確に記載されている
- 必要に応じてPostgreSQL関連の記述が追加されている

- _Requirements: 6.1, 6.2_
- _Design: 3.1_

---

#### タスク 2.3: docs/Docker.mdの修正
**目的**: `docs/Docker.md` からSQLite関連の記述を削除する。

**作業内容**:
- `docs/Docker.md` を開く
- SQLite関連の記述を確認
- 以下の記述を削除:
  - SQLiteに関する説明文
  - SQLiteに関する設定例
  - SQLiteに関する注意書き
- 削除後の文脈が自然になるように調整
- 必要に応じてPostgreSQL関連の記述を追加

**確認項目**:
- Dockerfileの説明でのデータベース記述
- docker-compose設定でのデータベース記述
- 環境変数でのデータベース記述

**受け入れ基準**:
- SQLite関連の記述が削除されている
- 削除後の文脈が自然である
- PostgreSQLを使用するDocker設定が明確に記載されている
- 必要に応じてPostgreSQL関連の記述が追加されている

- _Requirements: 6.1, 6.2_
- _Design: 3.3_

---

#### タスク 2.4: docs/API.mdの修正
**目的**: `docs/API.md` からSQLite関連の記述を削除する。

**作業内容**:
- `docs/API.md` を開く
- SQLite関連の記述を確認
- 以下の記述を削除:
  - SQLiteに関する説明文
  - SQLiteに関する設定例
  - SQLiteに関する注意書き
- 削除後の文脈が自然になるように調整
- 必要に応じてPostgreSQL関連の記述を追加

**確認項目**:
- APIエンドポイントの説明でのデータベース記述
- データベース接続の説明

**受け入れ基準**:
- SQLite関連の記述が削除されている
- 削除後の文脈が自然である
- PostgreSQLを使用するAPI設定が明確に記載されている
- 必要に応じてPostgreSQL関連の記述が追加されている

- _Requirements: 6.1, 6.2_
- _Design: 3.4_

---

#### タスク 2.5: docs/Admin.mdの修正
**目的**: `docs/Admin.md` からSQLite関連の記述を削除する。

**作業内容**:
- `docs/Admin.md` を開く
- SQLite関連の記述を確認
- 以下の記述を削除:
  - SQLiteに関する説明文
  - SQLiteに関する設定例
  - SQLiteに関する注意書き
- 削除後の文脈が自然になるように調整
- 必要に応じてPostgreSQL関連の記述を追加

**確認項目**:
- 管理画面の説明でのデータベース記述
- データベース接続の説明

**受け入れ基準**:
- SQLite関連の記述が削除されている
- 削除後の文脈が自然である
- PostgreSQLを使用する管理画面設定が明確に記載されている
- 必要に応じてPostgreSQL関連の記述が追加されている

- _Requirements: 6.1, 6.2_
- _Design: 3.5_

---

#### タスク 2.6: docs/Architecture.mdの修正
**目的**: `docs/Architecture.md` からSQLite関連の記述を削除する。

**作業内容**:
- `docs/Architecture.md` を開く
- SQLite関連の記述を確認
- 以下の記述を削除:
  - SQLiteに関する説明文
  - SQLiteに関する設定例
  - SQLiteに関する注意書き
- 削除後の文脈が自然になるように調整
- 必要に応じてPostgreSQL関連の記述を追加

**確認項目**:
- アーキテクチャ図でのデータベース記述
- データベース設計の説明

**受け入れ基準**:
- SQLite関連の記述が削除されている
- 削除後の文脈が自然である
- PostgreSQLを使用するアーキテクチャが明確に記載されている
- 必要に応じてPostgreSQL関連の記述が追加されている

- _Requirements: 6.1, 6.2_
- _Design: 3.6_

---

#### タスク 2.7: docs/Atlas-Operations.mdの修正
**目的**: `docs/Atlas-Operations.md` からSQLite関連の記述を削除する。

**作業内容**:
- `docs/Atlas-Operations.md` を開く
- SQLite関連の記述を確認
- 以下の記述を削除:
  - SQLiteに関する説明文
  - SQLiteに関する設定例
  - SQLiteに関する注意書き
- 削除後の文脈が自然になるように調整
- 必要に応じてPostgreSQL関連の記述を追加

**確認項目**:
- マイグレーション手順でのデータベース記述
- データベース接続の説明

**受け入れ基準**:
- SQLite関連の記述が削除されている
- 削除後の文脈が自然である
- PostgreSQLを使用するマイグレーション手順が明確に記載されている
- 必要に応じてPostgreSQL関連の記述が追加されている

- _Requirements: 6.1, 6.2_
- _Design: 3.7_

---

#### タスク 2.8: docs/Command-Line-Tool.mdの修正
**目的**: `docs/Command-Line-Tool.md` からSQLite関連の記述を削除する。

**作業内容**:
- `docs/Command-Line-Tool.md` を開く
- SQLite関連の記述を確認
- 以下の記述を削除:
  - SQLiteに関する説明文
  - SQLiteに関する設定例
  - SQLiteに関する注意書き
- 削除後の文脈が自然になるように調整
- 必要に応じてPostgreSQL関連の記述を追加

**確認項目**:
- CLIツールの説明でのデータベース記述
- データベース接続の説明

**受け入れ基準**:
- SQLite関連の記述が削除されている
- 削除後の文脈が自然である
- PostgreSQLを使用するCLIツール設定が明確に記載されている
- 必要に応じてPostgreSQL関連の記述が追加されている

- _Requirements: 6.1, 6.2_
- _Design: 3.8_

---

#### タスク 2.9: docs/Database-Viewer.mdの修正
**目的**: `docs/Database-Viewer.md` からSQLite関連の記述を削除する。

**作業内容**:
- `docs/Database-Viewer.md` を開く
- SQLite関連の記述を確認
- 以下の記述を削除:
  - SQLiteに関する説明文
  - SQLiteに関する設定例
  - SQLiteに関する注意書き
- 削除後の文脈が自然になるように調整
- 必要に応じてPostgreSQL関連の記述を追加

**確認項目**:
- データベースビューアの説明でのデータベース記述
- 接続設定の説明

**受け入れ基準**:
- SQLite関連の記述が削除されている
- 削除後の文脈が自然である
- PostgreSQLを使用するデータベースビューア設定が明確に記載されている
- 必要に応じてPostgreSQL関連の記述が追加されている

- _Requirements: 6.1, 6.2_
- _Design: 3.9_

---

#### タスク 2.10: docs/Generate-Sample-Data.mdの修正
**目的**: `docs/Generate-Sample-Data.md` からSQLite関連の記述を削除する。

**作業内容**:
- `docs/Generate-Sample-Data.md` を開く
- SQLite関連の記述を確認
- 以下の記述を削除:
  - SQLiteに関する説明文
  - SQLiteに関する設定例
  - SQLiteに関する注意書き
- 削除後の文脈が自然になるように調整
- 必要に応じてPostgreSQL関連の記述を追加

**確認項目**:
- サンプルデータ生成の説明でのデータベース記述
- データベース接続の説明

**受け入れ基準**:
- SQLite関連の記述が削除されている
- 削除後の文脈が自然である
- PostgreSQLを使用するサンプルデータ生成手順が明確に記載されている
- 必要に応じてPostgreSQL関連の記述が追加されている

- _Requirements: 6.1, 6.2_
- _Design: 3.10_

---

#### タスク 2.11: その他のdocs/*.mdファイルの修正
**目的**: その他の`docs/`ディレクトリ内のMarkdownファイルからSQLite関連の記述を削除する。

**作業内容**:
- 以下のファイルを確認:
  - `docs/File-Upload.md`
  - `docs/Logging.md`
  - `docs/Metabase.md`
  - `docs/Partner-Idp-Auth0-Login.md`
  - `docs/Project-Structure.md`
  - `docs/Queue-Job.md`
  - `docs/Rate-Limit.md`
  - `docs/Release-Check.md`
  - `docs/Send-Mail.md`
  - `docs/Sharding.md`
  - `docs/Testing.md`
  - その他の`docs/`ディレクトリ内のMarkdownファイル（`docs/Initial-Setup.md`と`docs/Spec-Driven-Development.md`は除外）
- 各ファイルでSQLite関連の記述を確認
- 以下の記述を削除:
  - SQLiteに関する説明文
  - SQLiteに関する設定例
  - SQLiteに関する注意書き
- 削除後の文脈が自然になるように調整
- 必要に応じてPostgreSQL関連の記述を追加

**確認項目**:
- 各ドキュメントでのデータベース記述
- 設定例でのデータベース記述

**受け入れ基準**:
- すべてのファイルでSQLite関連の記述が削除されている
- 削除後の文脈が自然である
- PostgreSQL/MySQLに関する記述が適切である
- 必要に応じてPostgreSQL関連の記述が追加されている

- _Requirements: 6.1, 6.2_
- _Design: 3.11_

---

#### タスク 2.12: プロジェクトルート直下のその他のMarkdownファイルの修正
**目的**: プロジェクトルート直下のその他のMarkdownファイルからSQLite関連の記述を削除する。

**作業内容**:
- プロジェクトルート直下のMarkdownファイルを確認（`README.md`以外）
- 各ファイルでSQLite関連の記述を確認
- 以下の記述を削除:
  - SQLiteに関する説明文
  - SQLiteに関する設定例
  - SQLiteに関する注意書き
- 削除後の文脈が自然になるように調整
- 必要に応じてPostgreSQL関連の記述を追加

**確認項目**:
- 各ファイルでのデータベース記述
- 設定例でのデータベース記述

**受け入れ基準**:
- すべてのファイルでSQLite関連の記述が削除されている
- 削除後の文脈が自然である
- PostgreSQL/MySQLに関する記述が適切である
- 必要に応じてPostgreSQL関連の記述が追加されている

- _Requirements: 6.1, 6.2_
- _Design: 2.3.1_

---

### Phase 3: PostgreSQL関連記述の追加

#### タスク 3.1: 削除箇所の確認とPostgreSQL情報の不足確認
**目的**: SQLiteの記述を削除した箇所で、PostgreSQL関連の記述がない箇所を特定する。

**作業内容**:
- 各ファイルでSQLiteの記述を削除した箇所を確認
- 削除箇所にPostgreSQL関連の記述があるか確認
- PostgreSQL情報が不足している箇所をリストアップ
- 追加が必要な箇所を特定

**確認条件**:
- SQLiteの記述を削除した箇所に、PostgreSQL関連の記述がない
- その箇所がデータベースに関する説明である
- PostgreSQLの情報を追加することで、ドキュメントの理解が向上する

**受け入れ基準**:
- 削除箇所がすべて確認されている
- PostgreSQL情報が不足している箇所が特定されている
- 追加が必要な箇所がリストアップされている

- _Requirements: 3.1.3_
- _Design: 2.3.2_

---

#### タスク 3.2: PostgreSQL関連記述の追加
**目的**: 削除箇所でPostgreSQL情報が不足している場合、PostgreSQL関連の記述を追加する。

**作業内容**:
- タスク 3.1で特定した箇所にPostgreSQL関連の記述を追加
- 追加内容:
  - PostgreSQLの接続情報（ホスト、ポート、ユーザー、パスワード、データベース名）
  - PostgreSQLの設定例
  - PostgreSQLに関する注意書き
  - PostgreSQLに関する説明文
- 既存のPostgreSQL/MySQLに関する記述を参考にする
- ドキュメントの構造やフォーマットを維持する
- 他のドキュメントの記述と一貫性を保つ

**追加例**:
```markdown
PostgreSQLの接続情報:
- ホスト: localhost
- ポート: 5432
- ユーザー: webdb
- パスワード: webdb
- データベース名: webdb_master
```

**受け入れ基準**:
- 削除箇所でPostgreSQL情報が不足している箇所に、適切な記述が追加されている
- 追加した記述が既存の記述と一貫している
- ドキュメントの構造やフォーマットが維持されている

- _Requirements: 6.2_
- _Design: 4.2_

---

### Phase 4: 整合性確認

#### タスク 4.1: ドキュメント間の一貫性確認
**目的**: すべてのドキュメントで記述が一貫していることを確認する。

**作業内容**:
- すべてのドキュメントを読み直す
- データベースに関する記述の一貫性を確認:
  - すべてのドキュメントで、データベースとしてPostgreSQL/MySQLが記載されている
  - SQLiteがデータベースオプションとして記載されていない
- 設定例の一貫性を確認:
  - 設定例でPostgreSQLが使用されている
  - 設定例でSQLiteが使用されていない
- 手順の一貫性を確認:
  - セットアップ手順でPostgreSQLが使用されている
  - セットアップ手順でSQLiteが使用されていない

**受け入れ基準**:
- すべてのドキュメントで、データベースとしてPostgreSQL/MySQLが記載されている
- SQLiteがデータベースオプションとして記載されていない
- 設定例でPostgreSQLが使用されている
- セットアップ手順でPostgreSQLが使用されている
- ドキュメント間で記述が一貫している

- _Requirements: 6.3_
- _Design: 5.1_

---

#### タスク 4.2: 設定例とセットアップ手順の確認
**目的**: 設定例とセットアップ手順が正しいことを確認する。

**作業内容**:
- すべてのドキュメントの設定例を確認
- すべてのドキュメントのセットアップ手順を確認
- 設定例でPostgreSQLが使用されていることを確認
- セットアップ手順でPostgreSQLが使用されていることを確認
- 設定例とセットアップ手順に矛盾がないことを確認

**受け入れ基準**:
- 設定例でPostgreSQLが使用されている
- セットアップ手順でPostgreSQLが使用されている
- 設定例とセットアップ手順に矛盾がない

- _Requirements: 3.2.2, 3.2.3_
- _Design: 5.1.2, 5.1.3_

---

### Phase 5: 最終確認

#### タスク 5.1: SQLite関連記述の再検索
**目的**: 修正後にSQLite関連の記述が残っていないことを確認する。

**作業内容**:
- `grep` コマンドを使用してSQLite関連の記述を再検索
- 検索パターン: `SQLite|sqlite|SQLite3|sqlite3`
- 検索対象: `README.md` と `docs/` ディレクトリ内のすべてのMarkdownファイル
- 検索結果を確認
- 検索対象外のファイル（`.kiro/specs/`、`cloudbeaver/config/`）の結果は除外

**検索コマンド例**:
```bash
grep -ri "SQLite\|sqlite" README.md docs/ --include="*.md"
```

**受け入れ基準**:
- `README.md` にSQLite関連の記述が残っていない
- `docs/` ディレクトリ内の対象MarkdownファイルにSQLite関連の記述が残っていない
- プロジェクトルート直下のMarkdownファイルにSQLite関連の記述が残っていない
- 検索対象外のファイル（`docs/Initial-Setup.md`、`docs/Spec-Driven-Development.md`、`.kiro/specs/`、`cloudbeaver/config/`）は除外されている

- _Requirements: 6.1_
- _Design: 5.2.1_

---

#### タスク 5.2: PostgreSQL関連記述の確認
**目的**: PostgreSQL関連の記述が適切であることを確認する。

**作業内容**:
- すべてのドキュメントを読み直す
- データベースに関する説明で、PostgreSQLが明確に記載されていることを確認
- 設定例でPostgreSQLが使用されていることを確認
- セットアップ手順でPostgreSQLが使用されていることを確認

**受け入れ基準**:
- データベースに関する説明で、PostgreSQLが明確に記載されている
- 設定例でPostgreSQLが使用されている
- セットアップ手順でPostgreSQLが使用されている

- _Requirements: 6.2_
- _Design: 5.2.2_

---

#### タスク 5.3: ドキュメントの読み直しと理解確認
**目的**: ドキュメントを読み直して、理解できることを確認する。

**作業内容**:
- すべてのドキュメントを読み直す
- SQLiteが使用できないことが明確であることを確認
- PostgreSQLがメインで使用されることが明確であることを確認
- 矛盾や不整合がないことを確認

**受け入れ基準**:
- ドキュメントを読んで、SQLiteが使用できないことが明確である
- ドキュメントを読んで、PostgreSQLがメインで使用されることが明確である
- ドキュメントに矛盾や不整合がない

- _Requirements: 6.4_
- _Design: 5.2.3_

---

## 実装上の注意事項

### 検索と確認
- **網羅的な検索**: `grep` コマンドなどでSQLite関連の記述を網羅的に検索
- **文脈の確認**: 検索結果を確認し、文脈を理解してから削除
- **影響範囲の確認**: 削除による影響範囲を確認

### 削除と追加
- **段階的な修正**: 検索 → 確認 → 削除 → 追加の順序で作業
- **文脈の維持**: 削除後も文脈が自然になるように調整
- **一貫性の確保**: 追加する記述が他のドキュメントと一貫していることを確認

### 整合性確認
- **全体の確認**: すべてのドキュメントを読み直して整合性を確認
- **設定例の確認**: 設定例が正しいことを確認
- **手順の確認**: セットアップ手順が正しいことを確認

### 最終確認
- **再検索**: 修正後に再度SQLite関連の記述を検索して、残っていないことを確認
- **読み直し**: ドキュメントを読み直して、理解できることを確認
- **矛盾の確認**: 矛盾や不整合がないことを確認

## 参考情報

### 関連Issue
- GitHub Issue #85: 開発環境はPostgreSQLを利用する前提とする
- GitHub Issue #91: ドキュメントの修正

### 既存ドキュメント
- `README.md`: プロジェクトの概要とセットアップ手順
- `docs/`: 機能別の詳細ドキュメント

### 技術スタック
- **データベース**: PostgreSQL/MySQL（全環境）
- **マイグレーションツール**: Atlas CLI
- **ドキュメント形式**: Markdown

### 参考リンク
- PostgreSQL公式ドキュメント: https://www.postgresql.org/docs/
- Atlas公式ドキュメント: https://atlasgo.io/
