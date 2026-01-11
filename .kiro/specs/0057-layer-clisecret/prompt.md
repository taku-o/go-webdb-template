/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/118
に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0057-layer-clisecretとしてください。"
think.

AdminサーバーにSecretを生成している箇所はないね。
Admin側の修正は要件から外しましょう。

要件定義書を承認します。

/kiro:spec-design 0057-layer-clisecret

作ったコマンドをビルドするときはバイナリは server/bin/ に出力するようにして。

設計書を承認します。

/kiro:spec-tasks 0057-layer-clisecret

タスクリストを承認します。

/sdd-fix-plan

_serena_indexing

/serena-initialize

/kiro:spec-impl 0057-layer-clisecret

server/list-dm-users でなくて、
server/bin/list-dm-users にビルドを出力して欲しい。


commitした後、
https://github.com/taku-o/go-webdb-template/issues/118
に対してpull requestを発行してください。

/review 119


