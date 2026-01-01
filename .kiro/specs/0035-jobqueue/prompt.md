/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/70 に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0035-jobqueueとしてください。"
think.

データ永続化設定をしないとデータが消えちゃうよね？
運用を考えると、データ永続化設定は必要。
> #### 3.1.1 Redisの導入
> - RedisをDockerコンテナとして導入
> - デフォルトポート: 6379
> - データ永続化の設定（オプション）
> - 開発環境での利用を想定


RedisとRedis Insightはおそらく別々のサーバで動作させることになる。
よって、RedisとRedis Insightは別のdocker-composeのファイルを用意して下さい。
> #### 3.1.3 Docker Compose設定
> - `docker-compose.redis.yml`ファイルを作成
> - RedisとRedis Insightのサービス定義


遅延時間については、ジョブごとに定義したい。
大半のジョブは同じ3分の定数を参照するとして。
> #### 3.5.2 ジョブ処理設定
> - 遅延時間の設定（デフォルト: 3分）
> - ジョブ処理の設定をコードに定数で定義


APIハンドラーは jobqueue_handler.go -> dm_jobqueue_handler.go とする
> ### 5.4 ファイル構造
> - **APIハンドラー**: `server/internal/api/handler/jobqueue_handler.go`（新規作成

要件定義書を承認します。

/kiro:spec-design 0035-jobqueue
think.

接続の設定は 環境変数`REDIS_HOSTS`でなく、config/{env}/cacheserver.yaml に定義する。
> 環境変数`REDIS_HOSTS`でRedis接続設定（別サーバーの場合は環境変数で設定可能）

既存のrate limitの機能がRedis or InMemoryで処理するようになっている。
rate limitがどちらのストレージを利用して動作するか。その設定を
config/{env}/config.yaml に追加したい。


docker-compose.redis.yml
は1台のRedisを起動する。開発環境での利用が主となる。

docker-compose.redis-insight.yml
はdocker-compose.redis.ymlと対になっていて、同じく1台のRedisサーバーと接続する想定で良い。


ちょっと方針を変更したい。
ジョブキューを管理するRedisは1台で処理する。
rate limitを処理するRedisは複数台になり得る。
つまり、Redisの環境設定を複数用意したい。

要件定義書、これは違う。
起動スクリプトは開発用途のため。Redis自体は本番でもstagingでも使用する。
> #### 3.1.1 Redisの導入
> - 開発環境での利用が主となる


要件定義書、ここも違う。
起動スクリプトは開発用途のため。Redis Insight自体は本番でもstagingでも使用する。
> #### 3.1.2 Redis Insightの導入
> - 開発環境での利用が主となる


要件定義書、複数台のRedisのキャッシュの方は、他の処理でも使いそうだから、
`redis.default.cluster.addrs`って名前にしておこう。
> - rate limit用: `redis.rate_limit.cluster.addrs`（複数台対応可能

Redisのデータはどこに保存される？
>    volumes:
>      - redis-data:/data

このdocker-composer.redis.ymlでRedisを起動した時、
Redisのデータはどこに保存される？

聞いているのは、設定ではなく、ディレクトリのPATHです。
この設定で起動したら、 **どこ** になるの？
PATHを書いて

> {プロジェクトルート}/redis-data

保存先ディレクトリを変更。
redis/data/jobqueue

設計書
Redisが起動していないケースも想定して。
APIサーバー起動時にRedisが起動していない場合、標準エラー出力にログを吐いて、
起動処理は継続する。
> #### 3.6.1 server/cmd/server/main.go への追加

ジョブタイプ
demo:delay_print の実装はどこに置かれる？

processor.goに処理が書かれる？


ジョブキューの処理にリトライ設定を入れたい。
リトライ回数と、遅延時間は定数でもっておいて。
ジョブ毎に設定を変えることになりそう。
```
client.Enqueue(task, 
    asynq.MaxRetry(10),                // 最大10回リトライ
    asynq.ProcessIn(5 * time.Minute),  // 初回実行を5分後に設定
)
```

もしかして、もう実装しようとしてる？
要件定義書、設計書の修正作業してたよね？

Contextが限界。次の会話のためにここまでの内容をまとめて。

設計書を承認します。
spec.jsonを更新してください。

/kiro:spec-tasks 0035-jobqueue
think.

タスクリストを承認します。

/sdd-fix-plan

_serena_indexing

/serena-initialize

/kiro:spec-impl 0035-jobqueue




