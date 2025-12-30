/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/52 に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0028-dmtable-defineとしてください。"
think.

JavaScriptの処理では、sonyflakeで生成したIDはJavaScriptの制限上、文字列として扱うこと。

server/cmd/generate-sample-data/ も修正する。

シャーディングの規則を次のように修正したので、
dm_users_NNN: table_sharding_key = id
dm_posts_NNN: table_sharding_key = user_id
ある dm_users に紐付いた dm_posts は、同じテーブル番号のテーブルにデータが入る。

identifier生成ルールとして、
* 数値のidentifierが必要な箇所はsonyflake (github.com/sony/sonyflake) を使用する。
* 文字列のidentifierが必要な箇所はUUIDv7 (github.com/google/uuid) を使用する。

これはドキュメントに記載する。

要件定義書を作成してください。

処理進んでる？

要件定義書のファイル作成処理が止まっていたから、Cursorを再起動したよ。
要件定義書を作成してください。

要件定義書のファイル作成処理が止まっていたから、要件定義書は私が代わりに作成しました。
作成しようとしていた内容と同じか要件定義書を確認してください。

要件定義書を承認します。

Agentモードに変更しました。
要件定義書を承認します。

spec.jsonがおかしかったので削除しました。
spec.jsonを作り直してください。

/kiro:spec-design 0028-dmtable-define

このサービスでIDは2種類使用する。
UUIDv7と、今回使用するsonyflakeです。

よって、GenerateID という関数名ではなく、
GenerateSonyflakeID という関数名でないと困る。


sf変数を持つ作りだから、
GenerateSonyflakeID -> GenerateID と戻して、
ファイル名の方を変更した方がいいか。
server/internal/util/idgen/idgen.go

このコード構成だと、
idgen.GenerateID() でsonyflakeのIDを生成する事になるんだっけ？

ごめんなさい、これやってください。
> 関数名をGenerateSonyflakeID()に戻す


設計書を承認します。

CursorでAcceptするまでファイルを変更しなくなったね？何か設定が変わった？

CursorでAcceptするまでファイルが変更されていないから、
Cursorのこの設定が変わってしまったかな？
> 自動適用から手動承認への切り替え

ひとまず置いとこう。

/kiro:spec-tasks 0028-dmtable-define

タスクリストのフォーマットがずいぶん変わったな。


UNSIGNED BIGINT はatlasだと、次のように定義すればいいらしい。

table "users" {
  schema = schema.public
  column "id" {
    type = bigint
    unsigned = true # ここで UNSIGNED を指定
  }
  column "name" {
    type = varchar(255)
  }
  primary_key {
    columns = [column.id]
  }
}


既存のデータの維持は考えなくて良い。
しかし、atlasのmigrations SQLを作成し直した時は、
db/migrations/master/20251229111855_initial_schema.sql の下の方に書いてある
初期データ用のSQLの移行を新しいファイルに移行するのを忘れてはいけない。

タスクリストを承認します。

/sdd-fix-plan









