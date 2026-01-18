/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/153
に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0075-jobserverとしてください。"
think.

要件定義書のフォーマットは
.kiro/specs/0023-metabase/requirements.md あたりのフォーマットと合わせてください。


今、demo:delay_printという定義があるんだね。
では、その仕組みに乗ろう。
この要件での作業は不要だ。
>## キューの中身
>ジョブをRedisに登録して、それから後に処理を行うので、
>Redisに登録するタスクに、タスクの分類用のキーが欲しい。
>キーは定数で定義する。

issueに載っていたクラス名も、
server/internal/usecase/jobqueue/dm_console_hello.go
server/internal/service/dm_console_hello_service.go
->
server/internal/usecase/jobqueue/delay_print.go
server/internal/service/delay_print_service.go
としよう。


ジョブ消化の処理の流れとしては、
- `server/internal/service/jobqueue/server.go`: ジョブハンドラーの登録
- `server/internal/service/jobqueue/processor.go`: 入出力制御とusecase層の呼び出し
- server/internal/usecase/jobqueue/delay_print.go: サービス層を呼び出して処理を実現する。ビジネスロジック。
- server/internal/service/delay_print_service.go: ビジネスユーティリティロジック
としよう。

要件定義書を承認します。

/kiro:spec-design 0075-jobserver

設計書を承認します。

/kiro:spec-tasks 0075-jobserver

タスクリストのフォーマットは
.kiro/specs/0023-metabase/tasks.md あたりのフォーマットと合わせてください。

ドキュメントの更新のタスクを追加してください。
どのドキュメントを更新すべきかを調査して、タスクリストに追加する。

README.ja.mdドキュメントも修正が必要。

タスクリストを承認します。

/sdd-fix-plan

_serena_indexing
/serena-initialize

/kiro:spec-impl 0075-jobserver





