/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/26 に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0014-db-atlasとしてください。"

今あるデータベースや、データベースの中身は破棄しても構いません。

要件定義書を作成してください


* 最終的に既存ツールが不要になったら削除する。
* 一通り完成したら、運用をある程度想定した操作の実験をやりたい。
    * develop環境のみで良し。
    * シナリオは 
        1. 0からのデータベースの初期化
        2. master側
            2-1. テーブルの追加
            2-2. テーブルにカラムを追加
            2-3. テーブルのデータを更新
            2-4. テーブルを削除
        3. sharding側
            2-1. テーブルの追加
            2-2. テーブルにカラムを追加
            2-3. テーブルのデータを更新
            2-4. テーブルを削除

* イレギュラーケースのシナリオの実験もいる。
    * 直接SQLを適用した後の作業 atlas migration hash ?

6.5のあと、Atlasでデータベースを構築したあとに、
既存のAPIサーバ、クライアント、管理画面が正しく動作するかを確認したい。

要件定義書を承認します。

/kiro:spec-design 0014-db-atlas

atlasはどのようにinstallすることになる？
プロジェクト内に置くのでなく、PATHにatlasが置かれることになる？

であれば、atlasのインストールは私がやっておきます。
```
brew install ariga/tap/atlas
```

システムの動作の要件のドキュメントがあったかな？
atlasの情報を追加したい。

atlasをインストールした。
atlasは見える？
which atlas

設計書を承認します

/kiro:spec-tasks 0014-db-atlas

タスクリストを承認します

この要件の作業用のgitブランチを切ってください。
ここまでの作業をcommitしてください。

_serena_indexing

/serena-initialize

/kiro:spec-impl 0014-db-atlas

db/seeds/master_seed.sqlは初期データ？

いろいろ作業する前に、
バックアップを兼ねたcommitをしておこうと思う。
git commitしなくて良いのは、
server/data/backup_before_atlas/
db/migrations_legacy/

db/seeds/master_seed.sql の中身は、
db/migrations/master/20251226074846_initial.sql
の後ろの方に入れちゃおう。
テーブル作ったけど、データが無いことを防ぐための手法らしい。

db/seedsの削除、および、ドキュメントの更新お願いします。
>  db/seeds/ ディレクトリは不要になりました。削除しますか？
>
>  また、README.mdとdocs/atlas-operations.mdにシードデータの手順が書いてあるので、そちらも更新が必要です。

では、
server/data/backup_before_atlas/
db/migrations_legacy/
を除いて、残りの編集をいったんcommitしてください。











今どこまで作業したかは、どこかに記録されていますか？
なければ、
.kiro/specs/0014-db-atlas/
に作業進捗の管理ファイルを作って、記録してください。





