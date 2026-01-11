/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/124
に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0060-layer-apiとしてください。"
think.

要件定義書を承認します。

/kiro:spec-design 0060-layer-api

後方互換性は不要。むしろ、admin、cliと実装が変わるので害の方が大きい。
パッケージ名は変更する。
>### 1.3 設計方針
>- **後方互換性**: パッケージ名は`package usecase`のまま維持（既存のコードとの互換性を保つ）

設計書を承認します。

/kiro:spec-tasks 0060-layer-api

タスクリストを承認します。

/sdd-fix-plan

_serena_indexing

/serena-initialize

/kiro:spec-impl 0060-layer-api

もしテストの実行で認証エラーが起きたのなら、APP_ENV=testを指定していない可能性があります。 .kiro/steering/tech.md を確認してください。

npm run test

動作確認は取れた。
残ったタスクはない？

commitした後、
https://github.com/taku-o/go-webdb-template/issues/124
に対してpull requestを発行してください。

/review 125



