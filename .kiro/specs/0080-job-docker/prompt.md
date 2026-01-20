JobQueueサーバーのDocker対応

JobQueueサーバーをDocker上で動作するように対応する。
docker-compose.api.yml などを参考に、対応してください。

ドキュメントの更新を忘れないように気をつけてください。


/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/163
に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0080-job-dockerとしてください。"
think.

コンテナ名に-developはつけずにjobqueueとする。
> #### 3.2.1 docker-compose.jobqueue.ymlの作成
> - **コンテナ名**: `jobqueue-develop`

次のcomposeのコンテナ名に-developがついているのは誤り。
取り除く要件を追加してください。
> docker-compose.api.yml
> docker-compose.admin.yml

docker-compose.client.yml もコンテナ名に誤りがある。
こちらへの対応も要件に追加してください。

要件定義書を承認します。

/kiro:spec-design 0080-job-docker

jobqueueはPostgreSQLの設定は要らなそうだけど、どうかな？
内部で参照してる？
> #### 3.1.1 基本構造

ああ、そうでした。
今は使っているジョブはないが、
拡張したらほぼPostgreSQLを使うことになるのでした。
PostgreSQLの設定は残して置いてください。

設計書を承認します。

/kiro:spec-tasks 0080-job-docker

タスクリストを承認します。

/sdd-fix-plan

_serena_indexing
/serena-initialize

/kiro:spec-impl 0080-job-docker

今動いている
APIサーバー
Clientサーバー
Adminサーバー
JobQueueサーバーを
を止めてください。

その後、
APIサーバー
Clientサーバー
JobQueueサーバーのDocker版を起動してください。


healthチェックで http://localhost:3000/healthでなく、http://localhost:3000/ にリクエストを送信した。
この情報はどこから取得した？
古い情報が残っていそうだ。更新したい。
>⏺ Bash(curl -s -o /dev/null -w "%{http_code}" http://localhost:3000/ && echo "")
>  ⎿  eval export HOMEBREW_PREFIX=/opt/homebrew
>     200


お願いします。
>  追加で、docker-compose.client.ymlにはhealthcheck設定がありませんが、API/JobQueueには設定されています。一貫性のためClientにもhealthcheckを
>  追加しますか？

JobQueueサーバーをいったん止めて。

JobQueueサーバーが生きてるかも。JobQueueサーバーのportは使われていないのだが、
Redisにキーを登録しているJobQueueサーバーが確かに居る。

APIサーバーを止めてみて。

おそらく、APIサーバーのビルドしたイメージのバージョンが古い。
ビルドしなおして。


Docker版のAPIサーバーを起動して。
Docker版のJobQueueサーバーを起動して。

APIサーバー、Clientサーバー、JobQueueサーバーを止めて。

stagingに上がっている修正をcommitして、
https://github.com/taku-o/go-webdb-template/issues/163 に
対してpull requestを作成してください。

/review 164






