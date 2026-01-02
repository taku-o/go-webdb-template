/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/74 に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0037-redis-reconnectとしてください。"
think.

要件定義書を承認します。

Redisは2種類使われているよ。
jobqueue用に1台構成。
cache server用に複数台構成。

修正箇所は2箇所あるかもしれない。
think.

複数台のRedis構成の方は、現在rate limit機能で使用されている。
現在、develop環境ではRedisでなく、Memoryを使用する設定になっている。
config/develop/config.yaml の api/rate_limit/storage_typeの設定を変更しないと、
うまく動作確認できないかもしれない。
think.

設計書を承認します。
spec.jsonを更新したら、いったんユーザーに応答を返してください。

/kiro:spec-tasks 0037-redis-reconnect
think.

config/develop/cacheserver.yaml
config/develop/config.yaml
に設定項目を追加するよね？

その追加する設定項目を
config/staging/cacheserver.yaml
config/staging/config.yaml
config/production/cacheserver.yaml.example
config/production/config.yaml.example
にも追加するタスクを足して欲しい。

タスクリストを承認します。

/sdd-fix-plan

_serena_indexing

/serena-initialize

/kiro:spec-impl 0033-mailsender




