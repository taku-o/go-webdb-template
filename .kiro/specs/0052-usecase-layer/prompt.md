/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/107
に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0052-usecase-layerとしてください。"
think.


このタスクの目的は全APIエンドポイントへの導入です。
> - 既存のすべてのAPIエンドポイントへのusecase層導入（today APIのみを対象）

要件定義書を承認します。

/kiro:spec-design 0052-usecase-layer

設計書を承認します。

/kiro:spec-tasks 0052-usecase-layer

アップロード処理は
今は UploadCompleteCallback の実装が定義されている箇所はない？
server/internal/api/handler/upload_handler.go

なら、今のところ、upload_handler.goは、
バリデーションと入出力の制御しかしてないってことになるね？

ならば、upload APIのusecase以降は用意しなくて良し！

タスクリストを承認します。

/sdd-fix-plan

_serena_indexing

/serena-initialize

/kiro:spec-impl 0052-usecase-layer



