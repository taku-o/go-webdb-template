/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/105
に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0051-testdbとしてください。"
think.


このシステムでは、外部キー制約を使用しない想定なので、気にしなくてよい。
仮に、外部制約が使われるケースがあるとしたら、それは絶対必須のデータであり、ほぼテスト対象でない。
むしろ一時的でも制約を無効化されると困る。
> - 外部キー制約がある場合は、制約を一時的に無効化してからクリア

要件定義書を承認します。

/kiro:spec-design 0051-testdb

設計書を承認します。

/kiro:spec-tasks 0051-testdb

今、scripts/migrate.shでデータベースを作成しているんだけど、
これのテスト用データベース向けスクリプトが必要かもしれない。
think.

タスクリストを承認します。

/sdd-fix-plan

_serena_indexing

/serena-initialize

/kiro:spec-impl 0051-testdb


