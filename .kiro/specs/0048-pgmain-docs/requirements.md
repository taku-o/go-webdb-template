# ドキュメントの修正要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #91
- **親Issue番号**: #85
- **Issueタイトル**: ドキュメントの修正
- **Feature名**: 0048-pgmain-docs
- **作成日**: 2025-01-27
- **対象ブランチ**: switch-to-postgresql

### 1.2 目的
プロジェクトでSQLiteを使用しないことになったため、ドキュメントからSQLite関連の記述を削除する。また、SQLiteの記述を削除した箇所に、今後メインで使用するPostgreSQL関連の記述がない場合は、PostgreSQL関連の記述を追加する。

### 1.3 スコープ
- プロジェクト内のドキュメントファイルからSQLite関連の記述を削除
- SQLiteの記述を削除した箇所に、PostgreSQL関連の記述がない場合は追加
- ドキュメントの整合性を保つ

**本実装の範囲外**:
- コード内のSQLite関連の記述（別issueで対応済み）
- 設定ファイル内のSQLite関連の記述（別issueで対応済み）
- `.kiro/specs/` 内の過去のspecファイル（履歴として残す）

## 2. 背景・現状分析

### 2.1 現在の状況
- **プロジェクト方針**: 開発環境、staging環境、production環境でPostgreSQLを利用する前提に変更された（Issue #85）
- **SQLiteの扱い**: 実際の開発現場でSQLiteが利用されることはないという想定から、SQLiteは利用しないことになった
- **既存の対応**: コードや設定ファイルからSQLite関連の記述は削除済み（Issue #85のsub issuesで対応）
- **残存する記述**: ドキュメント内にSQLite関連の記述が残っている可能性がある

### 2.2 課題点
1. **ドキュメントの不整合**: コードや設定ファイルからSQLiteが削除されているが、ドキュメントにSQLite関連の記述が残っている可能性がある
2. **情報の混乱**: ドキュメントを読む開発者がSQLiteを使用できると誤解する可能性がある
3. **PostgreSQL情報の不足**: SQLiteの記述を削除した箇所で、PostgreSQL関連の情報が不足している可能性がある

### 2.3 本実装による改善点
1. **ドキュメントの整合性**: ドキュメントと実装の整合性を保つ
2. **情報の明確化**: PostgreSQLをメインで使用することを明確にする
3. **開発者の混乱防止**: SQLiteが使用できないことを明確にする

## 3. 機能要件

### 3.1 ドキュメントの確認と修正

#### 3.1.1 対象ドキュメントの確認
以下のドキュメントディレクトリとファイルを確認し、SQLite関連の記述を検索する：
- `README.md`
- `docs/` ディレクトリ内のすべてのMarkdownファイル
- その他のプロジェクトルート直下のMarkdownファイル

#### 3.1.2 SQLite関連記述の削除
- **削除対象**: SQLite、sqlite、SQLite3、sqlite3などの記述
- **削除方法**: 
  - SQLiteに関する説明文を削除
  - SQLiteに関する設定例を削除
  - SQLiteに関する注意書きを削除
  - SQLiteに関する比較表の行を削除
- **注意事項**: 
  - 過去のspecファイル（`.kiro/specs/` 内）は履歴として残すため、削除しない
  - CloudBeaverの設定ファイル（`cloudbeaver/config/` 内）はツールの設定なので、削除しない

#### 3.1.3 PostgreSQL関連記述の追加
SQLiteの記述を削除した箇所で、以下の条件を満たす場合はPostgreSQL関連の記述を追加する：
- **追加条件**:
  - SQLiteの記述を削除した箇所に、PostgreSQL関連の記述がない
  - その箇所がデータベースに関する説明である
  - PostgreSQLの情報を追加することで、ドキュメントの理解が向上する
- **追加内容**:
  - PostgreSQLの接続情報
  - PostgreSQLの設定例
  - PostgreSQLに関する注意書き
  - PostgreSQLに関する説明文

### 3.2 ドキュメントの整合性確認

#### 3.2.1 データベースに関する記述の確認
- すべてのドキュメントで、データベースとしてPostgreSQL/MySQLが記載されていることを確認
- SQLiteがデータベースオプションとして記載されていないことを確認

#### 3.2.2 設定例の確認
- 設定例でSQLiteが使用されていないことを確認
- 設定例でPostgreSQLが使用されていることを確認

#### 3.2.3 手順の確認
- セットアップ手順でSQLiteが使用されていないことを確認
- セットアップ手順でPostgreSQLが使用されていることを確認

## 4. 非機能要件

### 4.1 ドキュメントの品質
- **正確性**: ドキュメントの内容が実装と一致していること
- **完全性**: 必要な情報が不足していないこと
- **一貫性**: ドキュメント間で記述が一貫していること

### 4.2 可読性
- **明確性**: PostgreSQLをメインで使用することが明確であること
- **理解しやすさ**: 開発者が誤解しないように記述されていること

### 4.3 保守性
- **更新容易性**: 将来の変更に対応しやすい構造であること
- **検索性**: SQLite関連の記述が残っていないことを確認しやすいこと

## 5. 制約事項

### 5.1 対象外のファイル
以下のファイルは修正対象外とする：
- `.kiro/specs/` 内の過去のspecファイル（履歴として残す）
- `cloudbeaver/config/` 内の設定ファイル（ツールの設定）
- コードファイル（別issueで対応済み）
- 設定ファイル（別issueで対応済み）

### 5.2 既存の記述の維持
- PostgreSQL/MySQLに関する既存の記述は維持する
- ドキュメントの構造やフォーマットは可能な限り維持する

### 5.3 ブランチ
- 修正は `switch-to-postgresql` ブランチに取り込む

## 6. 受け入れ基準

### 6.1 SQLite関連記述の削除
- [ ] `README.md` にSQLite関連の記述が残っていない
- [ ] `docs/` ディレクトリ内のすべてのMarkdownファイルにSQLite関連の記述が残っていない
- [ ] プロジェクトルート直下のMarkdownファイルにSQLite関連の記述が残っていない
- [ ] SQLiteに関する説明文が削除されている
- [ ] SQLiteに関する設定例が削除されている
- [ ] SQLiteに関する注意書きが削除されている

### 6.2 PostgreSQL関連記述の追加
- [ ] SQLiteの記述を削除した箇所で、PostgreSQL関連の記述がない場合は追加されている
- [ ] データベースに関する説明で、PostgreSQLが明確に記載されている
- [ ] 設定例でPostgreSQLが使用されている
- [ ] セットアップ手順でPostgreSQLが使用されている

### 6.3 ドキュメントの整合性
- [ ] すべてのドキュメントで、データベースとしてPostgreSQL/MySQLが記載されている
- [ ] SQLiteがデータベースオプションとして記載されていない
- [ ] ドキュメント間で記述が一貫している

### 6.4 動作確認
- [ ] ドキュメントを読んで、SQLiteが使用できないことが明確である
- [ ] ドキュメントを読んで、PostgreSQLがメインで使用されることが明確である
- [ ] ドキュメントに矛盾や不整合がない

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### 確認が必要なファイル
- `README.md`
- `docs/Initial-Setup.md`
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
- `docs/Spec-Driven-Development.md`
- `docs/Testing.md`
- その他のプロジェクトルート直下のMarkdownファイル

#### 変更しないファイル
- `.kiro/specs/` 内の過去のspecファイル（履歴として残す）
- `cloudbeaver/config/` 内の設定ファイル（ツールの設定）

### 7.2 既存ファイルの扱い
- 各ドキュメントファイルを確認し、SQLite関連の記述があれば削除
- SQLiteの記述を削除した箇所で、PostgreSQL関連の記述がない場合は追加
- 既存のPostgreSQL/MySQLに関する記述は維持

## 8. 実装上の注意事項

### 8.1 SQLite関連記述の検索
- **検索パターン**: SQLite、sqlite、SQLite3、sqlite3などの大文字小文字を区別しない検索
- **検索範囲**: ドキュメントファイル全体
- **検索ツール**: `grep` コマンドなどを使用して網羅的に検索

### 8.2 PostgreSQL関連記述の追加
- **追加判断**: SQLiteの記述を削除した箇所で、PostgreSQL関連の記述がない場合のみ追加
- **追加内容**: データベースに関する説明である場合のみ追加
- **追加方法**: 既存のPostgreSQL/MySQLに関する記述を参考にして追加

### 8.3 ドキュメントの整合性確認
- **確認方法**: すべてのドキュメントを読み直して、整合性を確認
- **確認項目**: 
  - データベースに関する記述が一貫しているか
  - 設定例が正しいか
  - セットアップ手順が正しいか

### 8.4 動作確認
- **確認方法**: ドキュメントを読んで、理解できるか確認
- **確認項目**: 
  - SQLiteが使用できないことが明確か
  - PostgreSQLがメインで使用されることが明確か
  - 矛盾や不整合がないか

## 9. 参考情報

### 9.1 関連Issue
- GitHub Issue #85: 開発環境はPostgreSQLを利用する前提とする
- GitHub Issue #91: ドキュメントの修正

### 9.2 既存ドキュメント
- `README.md`: プロジェクトの概要とセットアップ手順
- `docs/`: 機能別の詳細ドキュメント

### 9.3 技術スタック
- **データベース**: PostgreSQL/MySQL（全環境）
- **マイグレーションツール**: Atlas CLI
- **ドキュメント形式**: Markdown

### 9.4 参考リンク
- PostgreSQL公式ドキュメント: https://www.postgresql.org/docs/
- Atlas公式ドキュメント: https://atlasgo.io/
