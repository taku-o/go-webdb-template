/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/54 に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0026-chtablenameとしてください。"
think.

既存のデータベースのデータは破棄して良い。

要件定義書を作成してください。


分散データ環境では外部キー制約を使わない。
同じデータベースに参照先のデータがないのが普通。
あるいはその制約のせいでデータを移動できなくなる。
> ### 8.3 外部キー制約の参照先テーブル名変更
> - postsテーブルの外部キー制約が参照するusersテーブル名も変更する必要がある
> - 例: `FOREIGN KEY (user_id) REFERENCES users_000(id)` → `FOREIGN KEY (user_id) REFERENCES dm_users_000(id)`

最終的に、これらのファイル名も変更してください。
server/internal/model/user.go
server/internal/repository/user_repository.go
db/schema/sharding_1/users.hcl

要件定義書を承認します。

/kiro:spec-design 0026-chtablename

設計書を承認します。

/kiro:spec-tasks 0026-chtablename

タスクリストを承認します。

/sdd-fix-plan

_serena_indexing

/serena-initialize

/kiro:spec-impl 0025-fakedata





