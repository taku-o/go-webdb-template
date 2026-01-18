/kiro:spec-requirements "
JobQueueサーバーの死活監視
新しく作成したJobQueueサーバーにも死活監視を用意したい。
他のサーバーと合わせるなら /health だが、実現可能か？
ポートとか必要になるか？

実現可能なら、要件定義書を作成してください。

GitHub CLIは入っています。
cc-sddのfeature名は0076-jobqueue-healthとしてください。"
think.

https://github.com/taku-o/go-webdb-template/issues/155

Golang Asynqのジョブだけ消費するサーバーを作ったんだけど、死活監視ってどうするもの？
HTTPサーバーを建てて/healthでレスポンスを返すのは一般的？

OK。

要件定義書を承認します。

/kiro:spec-design 0076-jobqueue-health

config/develop/config.yaml だけでなく、
config/staging/config.yaml
config/production/config.yaml.example
にも設定を追加して。
staging、productionの動作確認まではしなくても良いけど。

設計書を承認します。

/kiro:spec-tasks 0076-jobqueue-health

更新すべきドキュメントがあったら更新するタスクを追加して欲しい。
まずは更新すべきドキュメントがあるかどうかの調査から。

タスクリストを承認します。

/sdd-fix-plan

_serena_indexing
/serena-initialize

/kiro:spec-impl 0076-jobqueue-health







