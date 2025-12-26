/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/32 に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0016-fix-tablesplitとしてください。"
think.

この先、テーブルの数が増えてしまうと思いますので、
次のようにテーブル毎にファイルを作りましょう。
```
db/schema/
└── sharding_1/             # シャード1用のディレクトリ
    ├── _schema.hcl         # スキーマ定義（schema "main" {} のみ）
    ├── users.hcl           # users_000 〜 users_007
    └── posts.hcl           # posts_000 〜 posts_007
```

不要なファイルは基本的に削除の方向で。
> 既存のdb/schema/sharding.hclは削除またはアーカイブ

既存データベースの移行は考えなくて良い。
リセットして良い。

要件定義書を作成してください。

要件定義書を承認します。

spec.jsonを作成してください。

posts_xxxテーブルのforeign_keyの定義は消す。
何故なら、こういう分割したテーブル設計の場合、
参照先のデータが同じデータベースにないことはよくあることだからです。

foreign_keyを消すにあたって、
消す理由とか、コードに変なコメントとか残さなくていいよ。

cleanup-sharding-dbs.sh は不要。
staging、productionは安易にデータを消せず、develop環境はSQLiteのデータファイルを消せばいい。

設計書を承認します。

作業の進捗を管理するために、.kiro/specs/0016-fix-tablesplit/progress.mdを作成して、
タスクの処理の最中に進捗を記録してください。

タスクリストを承認します。
spec.jsonを更新するところまで作業を進めたら、
ユーザーに応答を返してください。

この要件の作業用のgitブランチを切ってください。
ここまでの作業をcommitしてください。







