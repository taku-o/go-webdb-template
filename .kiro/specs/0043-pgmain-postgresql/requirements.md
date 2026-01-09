# PostgreSQL起動スクリプト・マイグレーションスクリプト修正要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #86
- **親Issue番号**: #85
- **Issueタイトル**: PostgreSQLの起動スクリプトと、Atlasマイグレーションスクリプトの修正
- **Feature名**: 0043-pgmain-postgresql
- **作成日**: 2025-01-27

### 1.2 目的
開発環境、staging環境、production環境でPostgreSQLを利用する前提に移行するため、PostgreSQLの起動スクリプトとAtlasマイグレーションスクリプトを修正する。
既存のSQLiteベースのスクリプトをPostgreSQL対応に変更し、Docker環境でPostgreSQLを起動・管理できるようにする。

### 1.3 スコープ
- PostgreSQLの起動スクリプトの修正（Docker Composeベース）
- Atlasマイグレーションスクリプトの修正（PostgreSQL対応）
- `scripts/migrate.sh`の修正（PostgreSQL対応）
- HCLファイルのUUIDv7カラム定義の修正（`bigint`から`varchar(32)`に変更）
- データベースへのデータ投入確認まで実施

**本実装の範囲外**:
- APIサーバー、GoAdminサーバーのPostgreSQL接続設定変更（別Issue対応）
- 本番環境へのデプロイ（準備のみ）
- データ移行ツールの作成（既存データの移行は対象外）

## 2. 背景・現状分析

### 2.1 現在の実装
- **PostgreSQL起動**: `docker-compose.postgres.yml`で1台のPostgreSQLコンテナを定義
- **マイグレーションスクリプト**: `scripts/migrate.sh`がSQLite用のAtlasマイグレーションを実行
- **データベース構成**:
  - マスターデータベース: 1台（`master.db`）
  - シャーディングデータベース: 4台（`sharding_db_1.db` ～ `sharding_db_4.db`）
  - 論理シャーディング数: 8（物理DB 4台 × 2論理シャード）
- **マイグレーションディレクトリ**: 
  - `db/migrations/master/` - マスターデータベース用
  - `db/migrations/sharding_1/` ～ `db/migrations/sharding_4/` - シャーディングデータベース用
- **設定ファイル**: `config/{env}/database.yaml`にPostgreSQL設定がコメントアウトされている

### 2.2 課題点
1. **PostgreSQL起動スクリプトの不足**: 現在は1台のPostgreSQLコンテナのみ定義されており、master 1台 + sharding 4台の構成に対応していない
2. **マイグレーションスクリプトのSQLite依存**: `scripts/migrate.sh`がSQLite用のURL形式（`sqlite://`）を使用しており、PostgreSQLに対応していない
3. **環境別設定の未整備**: 開発環境、staging環境、production環境でのPostgreSQL起動設定が未整備
4. **データ投入確認の不足**: マイグレーション適用後のデータ投入確認手順が未整備

### 2.3 本実装による改善点
1. **PostgreSQL環境の整備**: Docker Composeによるmaster 1台 + sharding 4台のPostgreSQL環境構築
2. **マイグレーションスクリプトのPostgreSQL対応**: AtlasマイグレーションスクリプトをPostgreSQL対応に修正
3. **環境別設定の整備**: 開発環境、staging環境、production環境でのPostgreSQL起動設定を整備
4. **データ投入確認の実装**: マイグレーション適用後のデータ投入確認手順を実装

## 3. 機能要件

### 3.1 PostgreSQL起動スクリプトの修正

#### 3.1.1 Docker Composeファイルの作成
- **ファイル**: `docker-compose.postgres.yml`を修正
- **構成**:
  - **masterグループ**: 1台のPostgreSQLコンテナ
    - サービス名: `postgres-master`
    - データベース名: `webdb_master`
    - ポート: `5432`（ホスト側ポートは`5432`）
  - **shardingグループ**: 4台のPostgreSQLコンテナ
    - サービス名: `postgres-sharding-1`, `postgres-sharding-2`, `postgres-sharding-3`, `postgres-sharding-4`
    - データベース名: `webdb_sharding_1`, `webdb_sharding_2`, `webdb_sharding_3`, `webdb_sharding_4`
    - ポート: `5433`, `5434`, `5435`, `5436`（ホスト側ポート）
- **ネットワーク**: 既存の`postgres-network`を使用
- **ボリューム**: 各PostgreSQLコンテナのデータを永続化
- **環境変数**: 
  - `POSTGRES_USER`: `webdb`
  - `POSTGRES_PASSWORD`: `webdb`
  - `POSTGRES_DB`: 各コンテナごとに設定

#### 3.1.2 起動スクリプトの作成
- **ファイル**: `scripts/start-postgres.sh`（新規作成または既存の修正）
- **機能**:
  - PostgreSQLコンテナの起動（`docker-compose -f docker-compose.postgres.yml up -d`）
  - PostgreSQLコンテナの停止（`docker-compose -f docker-compose.postgres.yml down`）
  - PostgreSQLコンテナの状態確認
  - ヘルスチェック確認
- **使用方法**: `./scripts/start-postgres.sh [start|stop|status|health]`

#### 3.1.3 環境別設定の対応
- **開発環境**: ローカルDocker環境での起動
- **staging環境**: staging環境用の設定ファイル（将来の拡張項目）
- **production環境**: production環境用の設定ファイル（将来の拡張項目）

### 3.2 Atlasマイグレーションスクリプトの修正

#### 3.2.1 scripts/migrate.shの新規作成
- **ファイル**: `scripts/migrate.sh`（既存スクリプトは破棄して新規作成）
- **変更内容**:
  - SQLite用のURL形式（`sqlite://`）からPostgreSQL用のURL形式（`postgres://`）に変更
  - PostgreSQL接続情報を設定ファイル（`config/{env}/database.yaml`）から読み込む
  - masterグループとshardingグループのAtlasマイグレーションをPostgreSQL対応に修正
  - 初期データ（`db/migrations/master/`内の初期データSQLファイル）を生SQLで適用（SQLite構文をPostgreSQL構文に変換）
    - ファイル名は適用順序に応じて変更可能（Atlasはファイル名順に適用）
  - Viewマイグレーション（`db/migrations/view_master/`）を生SQLで適用（Atlas PROライセンスが必要なため）
    - ファイル名は適用順序に応じて変更可能（Atlasはファイル名順に適用）

#### 3.2.2 マイグレーションURL形式
- **PostgreSQL形式**: `postgres://webdb:webdb@localhost:5432/webdb_master?sslmode=disable`
- **shardingグループ**: 各shardingデータベースごとに異なるポート・データベース名を使用

#### 3.2.3 設定ファイルからの読み込み
- **設定ファイル**: `config/{env}/database.yaml`からPostgreSQL接続情報を読み込む
- **環境変数**: `APP_ENV`環境変数で環境を指定（develop/staging/production）
- **読み込み項目**:
  - masterグループ: `database.groups.master[0]`から接続情報を取得
    - `host`, `port`, `user`, `password`, `name`（データベース名）
  - shardingグループ: `database.groups.sharding.databases[]`から各shardingデータベースの接続情報を取得
    - `host`, `port`, `user`, `password`, `name`（データベース名）
- **設定ファイルの構造**: 既存の`config/{env}/database.yaml`の構造を維持し、PostgreSQL設定のコメントアウトを解除して使用

### 3.4 HCLファイルのUUIDv7カラム定義修正

#### 3.4.1 db/schema/sharding_*/dm_posts.hclの修正
- **対象ファイル**: 
  - `db/schema/sharding_1/dm_posts.hcl`
  - `db/schema/sharding_2/dm_posts.hcl`
  - `db/schema/sharding_3/dm_posts.hcl`
  - `db/schema/sharding_4/dm_posts.hcl`
- **修正内容**:
  - `dm_posts_XXX.id`: `bigint` → `varchar(32)`
  - `dm_posts_XXX.user_id`: `bigint` → `varchar(32)`
  - `unsigned`属性を削除（PostgreSQLには存在しない）
  - `auto_increment`属性を削除（`id`カラム）

#### 3.4.2 db/schema/sharding_*/dm_users.hclの修正
- **対象ファイル**: 
  - `db/schema/sharding_1/dm_users.hcl`
  - `db/schema/sharding_2/dm_users.hcl`
  - `db/schema/sharding_3/dm_users.hcl`
  - `db/schema/sharding_4/dm_users.hcl`
- **修正内容**:
  - `dm_users_XXX.id`: `bigint` → `varchar(32)`
  - `unsigned`属性を削除（PostgreSQLには存在しない）
  - `auto_increment`属性を削除（`id`カラム）

**UUIDv7の仕様**:
- ハイフン抜き小文字で32文字の文字列
- PostgreSQLでは`varchar(32)`を使用

### 3.3 データ投入確認

#### 3.3.1 マイグレーション適用後の確認
- **マイグレーション適用**: `scripts/migrate.sh`を実行してマイグレーションを適用
- **テーブル存在確認**: 各データベースにテーブルが作成されていることを確認
- **データ投入**: テストデータを投入してデータベースにデータが入ることを確認
- **接続確認**: APIサーバーまたはGoAdminサーバーからPostgreSQLに接続できることを確認

## 4. 非機能要件

### 4.1 PostgreSQL環境の前提条件
- DockerおよびDocker Composeがインストールされていること
- Dockerが正常に動作していること
- ポート5432, 5433, 5434, 5435, 5436が使用可能であること（他のサービスと競合しないこと）
- 十分なメモリとディスク容量があること（5台のPostgreSQLコンテナを起動するため）

### 4.2 パフォーマンス
- **起動時間**: PostgreSQLコンテナの起動時間を最小化
- **マイグレーション時間**: Atlasマイグレーションの適用時間を最小化
- **接続時間**: データベース接続の確立時間を最小化

### 4.3 セキュリティ
- **パスワード管理**: 全環境（開発環境、staging環境、本番環境）で設定ファイルに固定パスワード（`webdb`）を記載する
- **ネットワーク**: Dockerネットワーク内での通信を前提
- **SSL/TLS**: 開発環境では`sslmode=disable`、本番環境では適切なSSL設定を推奨

### 4.4 データ永続化
- **ボリュームマウント**: 各PostgreSQLコンテナのデータをボリュームマウントで永続化
- **データディレクトリ**: `postgres/data/master/`, `postgres/data/sharding_1/` ～ `postgres/data/sharding_4/`にデータを保存

### 4.5 環境別対応
- **開発環境**: ローカルDocker環境での起動・マイグレーション
- **staging環境**: staging環境用の設定（将来の拡張項目）
- **production環境**: production環境用の設定（将来の拡張項目）

## 5. 制約事項

### 5.1 既存システムとの関係
- **既存のマイグレーションファイル**: 既存のマイグレーションファイル（`db/migrations/master/`, `db/migrations/sharding_*/`）は変更しない
- **既存の設定ファイル**: `config/{env}/database.yaml`の構造は維持（PostgreSQL設定のコメントアウト解除は別Issue対応）

### 5.2 技術スタック
- **PostgreSQL**: PostgreSQL 15-alpine（Dockerイメージ）
- **Atlas**: 既存のAtlas CLIを使用
- **Docker**: Docker Composeを使用
- **シェルスクリプト**: Bashスクリプトを使用

### 5.3 シャーディング構成
- **物理データベース数**: master 1台 + sharding 4台 = 合計5台
- **論理シャーディング数**: 8（各物理DBに2つの論理シャードを割り当て）
- **table_range**: 
  - sharding_1: [0, 7]（2つの論理シャード）
  - sharding_2: [8, 15]（2つの論理シャード）
  - sharding_3: [16, 23]（2つの論理シャード）
  - sharding_4: [24, 31]（2つの論理シャード）

### 5.4 運用上の制約
- **起動順序**: PostgreSQLコンテナを起動してからマイグレーションを適用する必要がある
- **データ投入**: マイグレーション適用後にデータ投入を確認する必要がある
- **ネットワーク**: Dockerネットワーク内での通信を前提

## 6. 受け入れ基準

### 6.1 PostgreSQL起動スクリプトの修正
- [ ] `docker-compose.postgres.yml`がmaster 1台 + sharding 4台の構成で定義されている
- [ ] `scripts/start-postgres.sh`が作成されている（または既存のスクリプトが修正されている）
- [ ] PostgreSQLコンテナが正常に起動する（`docker-compose -f docker-compose.postgres.yml up -d`）
- [ ] 各PostgreSQLコンテナが正常に動作していることを確認できる（ヘルスチェック）
- [ ] PostgreSQLコンテナを停止できる（`docker-compose -f docker-compose.postgres.yml down`）
- [ ] 各PostgreSQLコンテナのデータがボリュームマウントで永続化されている

### 6.2 HCLファイルのUUIDv7カラム定義修正
- [ ] `db/schema/sharding_1/dm_posts.hcl`のUUIDv7カラム定義が修正されている（`id`, `user_id`が`varchar(32)`）
- [ ] `db/schema/sharding_2/dm_posts.hcl`のUUIDv7カラム定義が修正されている（`id`, `user_id`が`varchar(32)`）
- [ ] `db/schema/sharding_3/dm_posts.hcl`のUUIDv7カラム定義が修正されている（`id`, `user_id`が`varchar(32)`）
- [ ] `db/schema/sharding_4/dm_posts.hcl`のUUIDv7カラム定義が修正されている（`id`, `user_id`が`varchar(32)`）
- [ ] `db/schema/sharding_1/dm_users.hcl`のUUIDv7カラム定義が修正されている（`id`が`varchar(32)`）
- [ ] `db/schema/sharding_2/dm_users.hcl`のUUIDv7カラム定義が修正されている（`id`が`varchar(32)`）
- [ ] `db/schema/sharding_3/dm_users.hcl`のUUIDv7カラム定義が修正されている（`id`が`varchar(32)`）
- [ ] `db/schema/sharding_4/dm_users.hcl`のUUIDv7カラム定義が修正されている（`id`が`varchar(32)`）
- [ ] すべてのHCLファイルから`unsigned`属性が削除されている
- [ ] すべてのHCLファイルから`auto_increment`属性が削除されている（`id`カラム）

### 6.3 Atlasマイグレーションスクリプトの修正
- [ ] `scripts/migrate.sh`がPostgreSQL対応の新規スクリプトとして作成されている（既存スクリプトは破棄）
- [ ] masterグループのAtlasマイグレーションがPostgreSQLで適用できる
- [ ] 初期データ（`db/migrations/master/`内の初期データSQLファイル）がPostgreSQLで適用できる（SQLite構文をPostgreSQL構文に変換、ファイル名は適用順序に応じて変更可能）
- [ ] Viewマイグレーション（`db/migrations/view_master/`）が生SQLで適用できる
- [ ] shardingグループのマイグレーションがPostgreSQLで適用できる（4台すべて）
- [ ] マイグレーション適用時にエラーが発生しない

### 6.4 データ投入確認
- [ ] マイグレーション適用後に各データベースにテーブルが作成されていることを確認できる
- [ ] テストデータを投入してデータベースにデータが入ることを確認できる
- [ ] APIサーバーまたはGoAdminサーバーからPostgreSQLに接続できることを確認できる（オプション）

### 6.5 ドキュメント
- [ ] `README.md`または`docs/`にPostgreSQL起動手順が記載されている
- [ ] `README.md`または`docs/`にPostgreSQLマイグレーション手順が記載されている
- [ ] 設定ファイル（`config/{env}/database.yaml`）の設定方法が記載されている

## 7. 影響範囲

### 7.1 新規追加が必要なディレクトリ・ファイル

#### ディレクトリ
- `postgres/data/master/`: master PostgreSQLのデータディレクトリ
- `postgres/data/sharding_1/` ～ `postgres/data/sharding_4/`: sharding PostgreSQLのデータディレクトリ

#### ファイル
- `scripts/start-postgres.sh`: PostgreSQL起動スクリプト（新規作成または既存の修正）

### 7.2 変更が必要なファイル

#### 設定ファイル
- `docker-compose.postgres.yml`: master 1台 + sharding 4台の構成に修正

#### スクリプト
- `scripts/migrate.sh`: PostgreSQL対応の新規スクリプトを作成（既存スクリプトは破棄）

#### HCLファイル
- `db/schema/sharding_1/dm_posts.hcl`: UUIDv7カラム定義を修正（`bigint` → `varchar(32)`）
- `db/schema/sharding_1/dm_users.hcl`: UUIDv7カラム定義を修正（`bigint` → `varchar(32)`）
- `db/schema/sharding_2/dm_posts.hcl`: UUIDv7カラム定義を修正（`bigint` → `varchar(32)`）
- `db/schema/sharding_2/dm_users.hcl`: UUIDv7カラム定義を修正（`bigint` → `varchar(32)`）
- `db/schema/sharding_3/dm_posts.hcl`: UUIDv7カラム定義を修正（`bigint` → `varchar(32)`）
- `db/schema/sharding_3/dm_users.hcl`: UUIDv7カラム定義を修正（`bigint` → `varchar(32)`）
- `db/schema/sharding_4/dm_posts.hcl`: UUIDv7カラム定義を修正（`bigint` → `varchar(32)`）
- `db/schema/sharding_4/dm_users.hcl`: UUIDv7カラム定義を修正（`bigint` → `varchar(32)`）

#### ドキュメント
- `README.md`: PostgreSQL起動・マイグレーション手順を追記

### 7.3 既存ファイルの扱い
- `db/migrations/master/`: 変更なし（既存のマイグレーションファイルを使用）
  - 初期データSQLファイル（現在: `20251230045548_seed_data.sql`）: 初期データが含まれている（SQLite構文をPostgreSQL構文に変換して適用、ファイル名は適用順序に応じて変更可能）
- `db/migrations/sharding_1/` ～ `db/migrations/sharding_4/`: 変更なし（既存のマイグレーションファイルを使用）
- `db/migrations/view_master/`: 変更なし（既存のViewマイグレーションファイルを使用、生SQLで適用、ファイル名は適用順序に応じて変更可能）
- `db/schema/sharding_*/dm_posts.hcl`: UUIDv7カラム定義を修正（`bigint` → `varchar(32)`）
- `db/schema/sharding_*/dm_users.hcl`: UUIDv7カラム定義を修正（`bigint` → `varchar(32)`）
- `scripts/migrate.sh`: 破棄して新規作成
- `config/{env}/database.yaml`: 変更なし（PostgreSQL設定のコメントアウト解除は別Issue対応）

## 8. 実装上の注意事項

### 8.1 Docker Compose設定
- **サービス名**: `postgres-master`, `postgres-sharding-1` ～ `postgres-sharding-4`
- **ネットワーク**: 既存の`postgres-network`を使用
- **ボリューム**: 各PostgreSQLコンテナのデータを個別のボリュームにマウント
- **環境変数**: 各コンテナで適切な環境変数を設定
- **ポート**: ホスト側ポートを競合しないように設定（5432, 5433, 5434, 5435, 5436）

### 8.2 HCLファイルのUUIDv7カラム定義修正
- **対象ファイル**: `db/schema/sharding_*/dm_posts.hcl`, `db/schema/sharding_*/dm_users.hcl`
- **修正内容**: 
  - `dm_posts_XXX.id`: `bigint` → `varchar(32)`
  - `dm_posts_XXX.user_id`: `bigint` → `varchar(32)`
  - `dm_users_XXX.id`: `bigint` → `varchar(32)`
  - `unsigned`属性を削除（PostgreSQLには存在しない）
  - `auto_increment`属性を削除（`id`カラム）
- **UUIDv7の仕様**: ハイフン抜き小文字で32文字の文字列

### 8.3 マイグレーションスクリプト
- **新規作成**: 既存の`scripts/migrate.sh`は破棄し、PostgreSQL対応の新規スクリプトを作成
- **URL形式**: PostgreSQL用のURL形式（`postgres://user:password@host:port/dbname?sslmode=disable`）を使用
- **設定ファイル**: `config/{env}/database.yaml`から接続情報を読み込む
- **環境変数**: `APP_ENV`環境変数で環境を指定（develop/staging/production）
- **Atlasマイグレーション**: `db/migrations/master/`, `db/migrations/sharding_*/`のAtlasマイグレーションを適用
- **初期データ**: `db/migrations/master/`内の初期データSQLファイルを生SQLで適用（SQLite構文をPostgreSQL構文に変換）
  - `INSERT OR IGNORE` → `INSERT ... ON CONFLICT DO NOTHING`
  - `datetime('now')` → `NOW()`
  - ファイル名は適用順序に応じて変更可能（Atlasはファイル名順に適用）
- **Viewマイグレーション**: `db/migrations/view_master/`の生SQLを適用（Atlas PROライセンスが必要なため、生SQLで実行）
  - ファイル名は適用順序に応じて変更可能（Atlasはファイル名順に適用）
- **エラーハンドリング**: マイグレーション適用時のエラーを適切に処理

### 8.4 データ投入確認
- **テーブル確認**: `psql`コマンドまたはAtlas CLIでテーブル存在を確認
- **データ投入**: テストデータを投入してデータベースにデータが入ることを確認
- **接続確認**: APIサーバーまたはGoAdminサーバーからPostgreSQLに接続できることを確認

### 8.5 設定ファイルの管理
- **設定ファイル**: `config/{env}/database.yaml`から接続情報を読み込む
- **環境変数**: `APP_ENV`環境変数で環境を指定（develop/staging/production）
- **設定構造**: 既存の`config/{env}/database.yaml`の構造を維持
- **PostgreSQL設定**: `config/{env}/database.yaml`のPostgreSQL設定のコメントアウトを解除して使用

### 8.6 ドキュメント整備
- **起動手順**: PostgreSQLコンテナの起動・停止手順を記載
- **マイグレーション手順**: Atlasマイグレーションの適用手順を記載
- **設定ファイル**: `config/{env}/database.yaml`の設定方法を記載
- **トラブルシューティング**: よくある問題と解決方法を記載

## 9. 参考情報

### 9.1 関連Issue
- GitHub Issue #85: 開発環境はPostgreSQLを利用する前提とする
- GitHub Issue #86: PostgreSQLの起動スクリプトと、Atlasマイグレーションスクリプトの修正

### 9.2 既存ドキュメント
- `README.md`: プロジェクト概要とセットアップ手順
- `docs/Architecture.md`: システムアーキテクチャ
- `docs/Initial-Setup.md`: 初期セットアップ手順
- `config/{env}/database.yaml`: 環境別データベース設定

### 9.3 技術スタック
- **PostgreSQL**: 15-alpine（Dockerイメージ）
- **Atlas**: Atlas CLI（既存）
- **Docker**: Docker Compose（既存）

### 9.4 参考リンク
- PostgreSQL公式ドキュメント: https://www.postgresql.org/docs/
- Atlas公式ドキュメント: https://atlasgo.io/
- Docker Compose公式ドキュメント: https://docs.docker.com/compose/
