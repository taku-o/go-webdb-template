/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/57 に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0030-uuidv7としてください。

ここ最近、要件定義書、設計書、タスクリストのフォーマットが揺れているので、
0023-metabaseあたりのフォーマットを参考にして、各仕様書を作成して欲しい。"
think.

既存のデータは維持不要。破棄して良いです。


dm_newsでは今はSonyflakeIDを使用していないと思う。
よっておそらく削除できる。
> #### 3.2.2 既存のSonyflake関数の扱い
> - **削除しない**: Sonyflake関数（`GenerateSonyflakeID()`）は他のテーブル（dm_newsなど）で使用されているため、削除しない
> - **共存**: UUIDv7関数とSonyflake関数は共存する

こちらも使用箇所はないが、しかし、まず確実に後々使用することになると思われるので残して置いてください。
> - **既存**: `GetShardingConnectionByID(id int64, tableName string) (*GORMConnection, error)`

GoAdminとクライアントの修正は必要ない？
特に、GoAdminは修正いるかも？

要件定義書を承認します。


/kiro:spec-design 0030-uuidv7
ここ最近、要件定義書、設計書、タスクリストのフォーマットが揺れているので、
0023-metabaseあたりのフォーマットを参考にして、各仕様書を作成して欲しい。
think.


設計書を承認します。


/kiro:spec-tasks 0030-uuidv7
ここ最近、要件定義書、設計書、タスクリストのフォーマットが揺れているので、
0023-metabaseあたりのフォーマットを参考にして、各仕様書を作成して欲しい。
think.


タスクリストを承認します。


/sdd-fix-plan

_serena_indexing

/serena-initialize

/kiro:spec-impl 0030-uuidv7

/clear

/serena-initialize

/kiro:spec-impl 0030-uuidv7
作業の継続お願いします。
think.


現在の作業進捗を
.kiro/specs/0030-uuidv7/progress.md か
.kiro/specs/0030-uuidv7/tasks.md に記録してください。

/kiro:spec-impl 0030-uuidv7 9

完了したタスクについてはtasks.mdのチェックを更新してください。

/kiro:spec-impl 0030-uuidv7 10

/kiro:spec-impl 0030-uuidv7 11

/kiro:spec-impl 0030-uuidv7 12

/kiro:spec-impl 0030-uuidv7 13

/kiro:spec-impl 0030-uuidv7 14

APIサーバ、クライアントサーバ、GoAdminサーバーを再起動してください。

server/cmd/generate-sample-data/ でサンプルデータを作成してください。

ここまでの修正をcommitしてください。
その後、https://github.com/taku-o/go-webdb-template/issues/57 に対して
pull requestを作成してください。








