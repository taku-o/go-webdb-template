/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/15 に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0022-apilimitとしてください。"
think.

要件定義書を作成してください。


本番では、複数台のRedisサーバーが使われると思います。
Redis Clusterを使う想定で。
```
Goでの対応: redis.NewClient ではなく redis.NewClusterClient を使います。

REDIS_ADDRS="host1:6379,host2:6379,host3:6379"
```

REDIS_ADDRSを使うなら、REDIS_URLの設定は使わない想定で良いよね？


レートリミット設定が未設定時は機能無効化にしましょう。
> ### 8.5 設定のデフォルト値
> - レートリミット設定が未指定の場合は機能を無効化
> - または、デフォルト値

レートリミットチェック時のエラーはRedisのサーバーと通信できないような場合だね。
この機能は、補助的な機能なので、APIへのリクエストは許可しちゃいましょう。
> レートリミットチェック時のエラーはログに記録し、リクエストは許可する（fail-open方式）または拒否する（fail-closed方式）を設計フェーズで決定

要件定義書を承認します。

/kiro:spec-design 0022-apilimit

設計書を承認します。

/kiro:spec-tasks 0022-apilimit

REDIS_ADDRSは環境変数でなく、設定で持ちたい。
新しく設定ファイルを作りたい。
config/{develop,staging,production}/cacheserver.yaml

タスクリストを承認します。

この要件の作業用のgitブランチを切ってください。
ここまでの作業をcommitしてください。

_serena_indexing

/serena-initialize

/kiro:spec-impl 0022-apilimit

まず、ドキュメントの更新はお願いします。
> 2. タスク 7.1: ドキュメントの更新（推奨）- ユーザー判断が必要

こちらは確認しなくていい。環境を用意するのが大変だから。
> 1. タスク 6.3: Redis Cluster環境での動作確認（tasks.mdに「オプション」と記載）

tasks.mdのチェックボックスを更新してください。

既存機能の動作確認したい。
APIサーバーと、クライアントサーバーを起動してください。

ソースコードの変更後、クライアントが動かないことがよくあるみたいなんだけど。
Failed to fetch RSC payload for http://localhost:3000/. Falling back to browser navigation.

あ、わかった。クライアントのポートを変えないでください。
3000を止めて。それからクライアントサーバーを起動。

ここまでの修正をcommitしてください。
その後、https://github.com/taku-o/go-webdb-template/issues/15 に対して
pull requestを作成してください。

/review 42






