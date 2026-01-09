/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/85 のsub issue
https://github.com/taku-o/go-webdb-template/issues/89 に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0046-pgmain-sampleとしてください。

issue 89の修正は、最終的に switch-to-postgresqlブランチに修正を取り込みます。"
think.

要件定義書を承認します。

/kiro:spec-design 0046-pgmain-sample

設計書を承認します。

/kiro:spec-tasks 0046-pgmain-sample

タスクリストを承認します。

この要件の作業用のgitブランチをswitch-to-postgresqlブランチから切ってください。
ここまでの作業をcommitしてください。
そこまで作業したら、いったんユーザーに応答を返してください。

_serena_indexing

/serena-initialize

/kiro:spec-impl 0046-pgmain-sample 1
/kiro:spec-impl 0046-pgmain-sample 2
/kiro:spec-impl 0046-pgmain-sample 3
/kiro:spec-impl 0046-pgmain-sample 4
/kiro:spec-impl 0046-pgmain-sample 5
/kiro:spec-impl 0046-pgmain-sample 6


author_idに入れるデータを修正してください。
> 修正するにはgofakeit.Int64()をgofakeit.Int32()に変更する必要があります。

CloudBeaverを起動してください。

/kiro:spec-impl 0046-pgmain-sample 7



