/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/151
のissueの条件でNext.jsのコードを修正する要件定義書を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0073-missing-testとしてください。"
think.

サーバー側は範囲外ではありません。
>**本実装の範囲外**:
>- サーバー側（Go）のコード


非huma版って何？humaを通ってない処理がある？
>#### 3.4.1 対象ファイル（テストが不足している可能性があるファイル）
>- `server/internal/api/handler/dm_user_handler.go`: ユーザーハンドラー（非Huma版）
>- `server/internal/api/handler/dm_post_handler.go`: 投稿ハンドラー（非Huma版

要件定義書の
必要に応じて、の記載をいくつか消しました。

要件定義書を承認します。

/kiro:spec-design 0073-missing-test

クライアント側のテストのルールとして、useEffectを極力使わない、を追加。

設計書を承認します。

/kiro:spec-tasks 0073-missing-test

タスクリストを承認します。

/sdd-fix-plan
_serena_indexing
/serena-initialize

/kiro:spec-impl 0073-missing-test

