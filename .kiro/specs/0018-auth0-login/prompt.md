/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/30 に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0018-auth0-loginとしてください。

関連：Auth0セットアップ資料 docs/Partner-Idp-Auth0-Login.md "

think.

config/{env}/config.yamlは、gitにコミットしたいけっど、
他の設定項目と違って、Auth0の設定情報は外部に漏れると困る情報になる？

client secretは内部で作った情報だから、develop環境なら漏れても良いんだけど、
Auth0は外部サービスの情報だから、漏れたら迷惑になるかな、と。

これしたとき、将来、APIサーバのAPIを実行する時、HeaderにAPIキーを載せられる？
> JWTはHTTP-only Cookieに保存（セキュリティ考慮）

これなんだけど、両方のAPIキーを使いたい。
未ログイン時はNEXT_PUBLIC_API_KEY、
ログイン時はユーザーのJWTでAPIを実行する。
> 現在のapiClientとの関係
> 現在はNEXT_PUBLIC_API_KEYを使用（client/src/lib/api.ts）
> 将来JWTを使用する場合、apiClientを修正してJWTをHeaderに載せる（別issueで対応）

今回のissueのターゲットはログインして、
JWTをとって保存するまでね。
そのJWTを使ってAPIを叩くところは、次のissueで対応します。

ログイン状態の表示と、
ログアウトボタンも必要だったね。
ログアウトは HTTP-only Cookie からJWTを消すだけで実現できる？
実現できるなら、要件に取り込み。

要件定義書を作成してください

要件定義書を承認します。
spec.jsonを更新したら、ユーザーに応答を返してください。

/kiro:spec-design 0018-auth0-login

client/.env.development というファイルがあるけど、
この.env.localは、APIサーバー用になる？全体用になる？

今まで client/.env.develop だけでいいやと思ってたけど、
.env.local、.env.stagingとかひな形は、今回のissueで作っちゃってください。

設計書を承認します。

/kiro:spec-tasks 0018-auth0-login

config.yamlにも新たにAuth0の設定が入る？

今回のタスクで作業が発生しないなら、タスクから削除しましょう。

タスクリストを承認します。

この要件の作業用のgitブランチを切ってください。
ここまでの作業をcommitしてください。

_serena_indexing

/serena-initialize

/kiro:spec-impl 0018-auth0-login




