/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/29 に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0025-fakedataとしてください。"
think.

要件定義書を作成してください。

/sdd-requirements-approve

/kiro:spec-design 0025-fakedata

/sdd-design-approve

/kiro:spec-tasks 0025-fakedata

/sdd-tasks-approve

/sdd-fix-plan

_serena_indexing

/serena-initialize

/kiro:spec-impl 0025-fakedata

今はデータが入っている状態？
APIサーバーとクライアントサーバーを起動して

クライアントのユーザー管理の画面、投稿管理の画面でエラーが出てる


分散テーブルデータの環境では、
IDはAUTOINCREMENTで生成するのが正しいです。

つまり、クロステーブルクエリという処理がおかしいです。
実装側に分散テーブル環境での基礎知識が足りないようです。

クロステーブルクエリで何をしているか、教えてください。


IDが重複していたら、何故エラーになるの？
リストに入れているだけじゃないの？


原因が見えた。テーブルのスキーマの設計が
分散テーブル環境の常識を知らないためだ。
各テーブルにはIDと、個別のデータを一意に決めるためのidentifierがいる。

usersのテーブルの場合は、
idの他に
user_id のようなカラムを設けて、
そこに全テーブルで一意のidentifierを入れるテーブル設計にしなければいけないのです。

shardingの規則も良くないね。
idで分割するような分け方はしない。

しかし、それは別のissueで対応しよう。

IDを参照するようなループじゃなくて、
単純ループの実装にできないかな？

ここまでの修正をcommitしてください。
その後、https://github.com/taku-o/go-webdb-template/issues/29 に対して
pull requestを作成してください。

/review 53



