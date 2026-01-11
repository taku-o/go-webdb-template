/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/120
に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0058-layer-sampledataとしてください。"
think.


gofakeitをサービス層で使うなら、
サンプルデータ生成プログラム用の専用Repository層は新規に作らないで、
server/internal/repository/dm_user_repository.go
server/internal/repository/dm_post_repository.go
server/internal/repository/dm_news_repository.go
に
  - `InsertDmUsersBatch(ctx context.Context, tableName string, dmUsers []*model.DmUser) error`
  - `InsertDmPostsBatch(ctx context.Context, tableName string, dmPosts []*model.DmPost) error`
  - `InsertDmNewsBatch(ctx context.Context, dmNews []*model.DmNews) error`
メソッドを作るか利用することにしよう。


要件定義書を承認します。

/kiro:spec-design 0058-layer-sampledata

設計書を承認します。

/kiro:spec-tasks 0058-layer-sampledata

タスクリストを承認します。

/sdd-fix-plan

_serena_indexing

/serena-initialize

/kiro:spec-impl 0058-layer-sampledata



commitした後、
https://github.com/taku-o/go-webdb-template/issues/120
に対してpull requestを発行してください。

/review 121


