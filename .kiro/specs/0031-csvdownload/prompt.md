/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/62 に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0031-csvdownloadとしてください。"
think.

可能なら次のような実装を入れて。タイムアウト設定。
```
// Humaのハンドラー内でのイメージ
func(ctx context.Context, input *struct{}) (*huma.StreamResponse, error) {
    return &huma.StreamResponse{
        Body: func(w io.Writer) {
            // http.ResponseWriterを取り出す
            rw, ok := w.(http.ResponseWriter)
            if ok {
                rc := http.NewResponseController(rw)
                rc.SetWriteDeadline(time.Now().Add(3 * time.Minute))
            }

            // CSV書き込み処理...
        },
    }, nil
}
```

要件定義書を承認します。

/kiro:spec-design 0031-csvdownload

BOMは無し
> ### 10.4 CSV形式の互換性
> - **注意**: Excelなどの一部のアプリケーションはBOM（Byte Order Mark）を期待する場合がある
> - **対応**: 現時点ではBOMは追加しない（要件定義に従う）

APIサーバー全体のデフォルトのタイムアウトが設定されていないかも。
設定されていないなら設定したい。

IdleTimeoutってどんな時に使われる？

APIサーバーのデフォルトのIdleTimeout 120秒でお願いします。

設計書を承認します。

/kiro:spec-tasks 0031-csvdownload

タスクリストを承認します。

/sdd-fix-plan

_serena_indexing

/serena-initialize

/kiro:spec-impl 0031-csvdownload













