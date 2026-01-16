/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/145
に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0070-dm_feedとしてください。"
think.

元の実装を残すこととか、安全性とか考えなくてよい。
シンボリックリンクは使わない。
>### 8.1 URLパスの変更方法
>- Next.jsのApp Routerでは、ディレクトリ名がURLパスになる
>- ただし、本実装ではファイルパスは変更せず、URLパスを変更する必要がある
>- そのため、`client/app/feed/`ディレクトリを`client/app/dm_feed/`に移動する必要がある
>- または、`client/app/dm_feed/`ディレクトリを作成し、既存のファイルをコピーまたはシンボリックリンクで対応する

要件定義書を承認します。

/kiro:spec-design 0070-dm_feed

設計書を承認します。

/kiro:spec-tasks 0070-dm_feed

タスクリストを承認します。

/sdd-fix-plan

_serena_indexing
/serena-initialize

/kiro:spec-impl 0070-dm_feed

commitして、
https://github.com/taku-o/go-webdb-template/issues/145 に向けた
pull requestを作成してください。

/review 146

