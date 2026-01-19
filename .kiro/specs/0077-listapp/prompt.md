/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/157
に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0077-listappとしてください。"
think.


確認したいサーバーは、API、Client、Admin、JobQueue、
PostgreSQL、MySQL、Redis、Redis Cluster、Mailpit、
CloudBeaver、Superset、Metabase、Redis Insight

並び順は、この通りで。
> #### 3.1.1 確認対象サーバー

Clientサーバーに/healthエンドポイントが追加されました。
要件定義書に反映してください。

プログラム名は、app-statusとしましょう。
>#### 3.3.3 実行方法
>- **シェルスクリプトの場合**: `./scripts/listapp.sh`または`bash scripts/listapp.sh`
>- **Goプログラムの場合**: `go run ./cmd/listapp/main.go`またはビルドして実行

失礼。やっぱりserver-statusという名前にしましょう。
>#### 3.3.3 実行方法
>- **シェルスクリプトの場合**: `./scripts/listapp.sh`または`bash scripts/listapp.sh`
>- **Goプログラムの場合**: `go run ./cmd/listapp/main.go`またはビルドして実行

要件定義書を承認します。

/kiro:spec-design 0077-listapp

設計書、一部を書き換えました。

設計書を承認します。

/kiro:spec-tasks 0077-listapp

ドキュメントの更新のタスクを追加してください。
更新すべきドキュメントを調べて、それを更新するタスクを追加する。
think.

タスクリストを承認します。

/sdd-fix-plan

_serena_indexing
/serena-initialize

/kiro:spec-impl 0077-listapp

実行バイナリは
server/bin
に出力してください。

タスク外の作業なのですが、
./scripts/start-cloudbeaver.sh を修正して、
./scriptsディレクトリの他のスクリプトと挙動を合わせて、

./scripts/start-cloudbeaver.sh start
./scripts/start-cloudbeaver.sh stop
が機能するように修正してください。


stagingに上がっている修正をcommitして、
https://github.com/taku-o/go-webdb-template/issues/157 に
対してpull requestを作成してください。

/review 160



