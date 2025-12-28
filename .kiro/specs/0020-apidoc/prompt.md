/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/37 に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0020-apidocとしてください。"
think.

要件定義書を作成してください。

要件定義書を承認します。

/kiro:spec-design 0020-apidoc

設計書を承認します。

/kiro:spec-tasks 0020-apidoc

タスクリストを承認します。

この要件の作業用のgitブランチを切ってください。
ここまでの作業をcommitしてください。

_serena_indexing

/serena-initialize

/kiro:spec-impl 0020-apidoc

クライアントサーバーとAPIサーバーを再起動して。

Authorization: Bearer 123 となっているけど、
この123を、dummy_jwt_api_key と変えられる？

その123は、我々が指定した値ではない？
勝手に入る値？

現状の実装で問題ない

タグでPublic API、Private APIを追加したけど、
そうすると、users、postsなどを既に使っているから、
左のメニューで混在しちゃうんだね。これは良くない。

SummaryとDescriptionに追加する方式に変更しよう。

```
huma.Register(api, huma.Operation{
    OperationID: "get-user",
    Method:      http.MethodGet,
    Path:        "/user",
    Summary:     "[private] ユーザー情報の取得", // タイトルに含める
    Description: "**Access Level:** `private` (Requires Auth0 JWT)", // 詳細に太字で書く
}, ...)
```

APIサーバーが再起動していないかも。


today APIのDescriptionを修正。のみ、を削除。
```
(Auth0 JWT のみでアクセス可能)
->
(Auth0 JWT でアクセス可能)
```

publicなAPIのSummaryを修正。[public] を削除。
```
Summary: "[public] ユーザー情報の取得",
->
Summary: "ユーザー情報の取得",
```

なんか変になったんだけど？
まずDescriptionが表示されなくなった。

次にtoday APIの表示がおかしい。
```
[private] 今日の日付を取得
->
今日の日付を取得（Auth0認証必須）
```
と変わった。

OK。直りました。良い感じです。

ここまでの修正をcommitしてください。
その後、https://github.com/taku-o/go-webdb-template/issues/37 に対して
pull requestを作成してください。


/review 42





