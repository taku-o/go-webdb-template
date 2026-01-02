/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/72 に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0036-db-reconnectとしてください。"
think.

要件定義書を承認します。

/kiro:spec-design 0036-db-reconnect
think.


DefaultConnectionMaxLifetimeは1時間程度にしておいて。
> DefaultConnectionMaxLifetime = 5 * time.Minute

設計書を承認します。

/kiro:spec-tasks 0036-db-reconnect
think.

全ての作業が終わったら、開発環境のデータベースの設定はSQLite版に戻したい。
開発環境のPostgreSQLの設定は消すしかないなら、コメントアウトする。
この作業を行うのは、最後の最後、ユーザーの確認作業が終わったらになるよ。
think.

タスクリストを承認します。

/sdd-fix-plan

_serena_indexing

/serena-initialize

/kiro:spec-impl 0036-db-reconnect



