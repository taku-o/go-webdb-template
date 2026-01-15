プロジェクトに英語のドキュメントを用意する。

日英で差分が生じたときのマスターのドキュメントは日本語版で、
内部コメントは今までと変わらず日本語を使用する。

README.md -> README.ja.md
README.md (新規作成。英語版を用意)

docs/*.md -> docs/ja/*.md
docs/en/*.md (新規作成。英語版を用意)

_config.yml descriptionを英語にする

/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/140
に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0068-endocとしてください。"
think.

要件定義書を承認します。

/kiro:spec-design 0068-endoc

各mdファイルの冒頭か末尾に言語切替のリンクが欲しい。
> **[日本語] | [English]**

設計書を承認します。

/kiro:spec-tasks 0068-endoc

タスクリストを承認します。

/sdd-fix-plan

_serena_indexing
/serena-initialize

/kiro:spec-impl 0068-endoc


このルールがあるのに、コメントが英語に変えられちゃってる。
>日英で差分が生じたときのマスターのドキュメントは日本語版で、
>内部コメントは今までと変わらず日本語を使用する。
docs/_config.yml


このルールをどこかに記録したい。
>日英で差分が生じたときのマスターのドキュメントは日本語版で、
>内部コメントは今までと変わらず日本語を使用する。


stagingに上がっている修正をcommitして、
https://github.com/taku-o/go-webdb-template/issues/140
にpull requestを発行してください。

/review 141

