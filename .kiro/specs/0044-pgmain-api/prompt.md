/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/85 のsub issue
https://github.com/taku-o/go-webdb-template/issues/87 に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0044-pgmain-apiとしてください。

issue 85の修正は、最終的に switch-to-postgresqlブランチに修正を取り込みます。"
think.

SQLiteは利用しなくなるので、SQLiteの設定はコメントアウトでなく、
削除する。

SQLite用のライブラリを読み込んでいたら取り除く。
ソースコード中にSQLite用の処理の分岐があったら、それも取り除く。

論理的なshardingグループのシャーディング数は8とする。
よって、config/develop/database.yaml の設定は8つないといけない。
現在SQLite版で4つしか指定していないのは、いつの間にか書き換えられたバグである。

要件定義書を承認します。

/kiro:spec-design 0044-pgmain-api

config/production/database.yaml.example というファイルがある。
修正漏れに注意して。

設計書を承認します。

/kiro:spec-tasks 0044-pgmain-api

タスクリストを承認します。

この要件の作業用のgitブランチをswitch-to-postgresqlブランチから切ってください。
ここまでの作業をcommitしてください。
そこまで作業したら、いったんユーザーに応答を返してください。


_serena_indexing

/serena-initialize

/kiro:spec-impl 0044-pgmain-api



