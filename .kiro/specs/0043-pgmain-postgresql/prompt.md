/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/85 のsub issue
https://github.com/taku-o/go-webdb-template/issues/86 に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0043-pgmain-postgresqlとしてください。"
think.

このSQLite、PostgreSQL切替機能は不要です。
> #### 3.2.1 scripts/migrate.shの修正
>   - 環境変数`DB_TYPE`でSQLite/PostgreSQLを切り替え可能にする（後方互換性のため）


環境変数はなるべく使用しない。
設定ファイルから読み込む。

scripts 内には、起動スクリプト系を集めたいので、
scripts/verify-postgres.shは作成しない。
> #### 3.3.2 確認スクリプトの作成（オプション）
> - **ファイル**: `scripts/verify-postgres.sh`（オプション）

staging環境、本番環境も設定ファイルに固定パスワードを記載する。
> ### 4.3 セキュリティ
> - **パスワード管理**: 開発環境では固定パスワード（`webdb`）を使用、本番環境では設定ファイルで管理（gitから除外）
> - **ネットワーク**: Dockerネットワーク内での通信を前提
> - **SSL/TLS**: 開発環境では`sslmode=disable`、本番環境では適切なSSL設定を推奨

要件定義書を承認します。
spec.jsonを更新してから、ユーザーに応答を返してください。

/kiro:spec-design 0043-pgmain-postgresql

既存のAtlasマイグレーションスクリプトは破棄しても良い。
db/migrations/master/20251230045548_seed_data.sql に初期データが入っていることに注意。

Viewに関しては生SQLを入れる。
AtlasでViewを使うにはPROライセンスが必要なため。
db/migrations/view_master/20260103030225_create_dm_news_view.sql
think.


Atlasでマイグレーションを適用するとき、ファイル名順にマイグレーションが適用される。
適用タイミングを調整するために、ファイル名は変更してよい。
db/migrations/master/20251230045548_seed_data.sql
db/migrations/view_master/20260103030225_create_dm_news_view.sql


hclの次のカラムの定義が間違っているようだ。
UUIDv7のハイフン抜き小文字なので、32文字の文字列となる。
これもこのタスクで修正してして欲しい。

db/schema/sharding_1/dm_posts.hcl
dm_posts_000.id
dm_posts_000.user_id

db/schema/sharding_1/dm_posts.hcl
dm_users_000.id

設計書を承認します。
spec.jsonを更新してから、ユーザーに応答を返してください。

/kiro:spec-tasks 0043-pgmain-postgresql

過去に1台構成でPostgreSQLを使用していた。
過去のデータが残っていたなら、それは破棄して良い。

タスクリストを承認します。

/sdd-fix-plan

_serena_indexing

/serena-initialize

/kiro:spec-impl 0043-pgmain-postgresql





