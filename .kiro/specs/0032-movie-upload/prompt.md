/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/64 に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0032-movie-uploadとしてください。"
think.

3.1.2記載の開発環境のエンドポイントは、これはAPIサーバー側？
http://localhost:8080/api/upload/dm_movieとして欲しい。
> #### 3.1.2 TUSクライアントの設定
> - **実装内容**:
>   - uppy を使用する場合: `@uppy/core` と `@uppy/tus` を使用
>   - tus-js-client を使用する場合: `tus-js-client` ライブラリを使用
>   - エンドポイント: `http://localhost:8080/files/` (開発環境)
>   - チャンクサイズ: 5MB
>   - リトライ遅延: [0, 1000, 3000, 5000] ms


ファイルの検証は必要なんだけど、処理が重いから、
本番はAWSの機能を使って、ファイルアップロード完了後にS3上でやろうと考えている。
よってコード上では実装は不要だ。
> ### 4.2 セキュリティ
> - ファイル検証: アップロードされたファイルの検証（将来的な拡張項目）

この中で、ファイルサイズ制限と、追加でファイル拡張子の制限は入れておきたい。
> ## 8. 将来の拡張項目（現時点では未実装）
> 
> 以下の機能は将来の拡張として検討されていますが、現時点では実装対象外です：
> 
> - ファイルサイズ制限の設定

Goバックエンドで、処理完了後に通知を受けるURLが必要？
> handler, err := handler.NewHandler(handler.Config{
>     BasePath:              "/api/upload/dm_movie",
>     StoreComposer:         composer,
>     NotifyCompleteUrls:    []string{"http://localhost:8080/hook"}, // 完了後の処理
> })

疑問の段階だから、まだ要件定義書を修正しないで。
・この機能は必要なの？
・何に使われるの？

NotifyCompleteUrlsが無い場合は、アップロード完了時に何らかの処理が差し込めない？

NotifyCompleteUrlsの実装でなくて、
この実装を入れておいて。
今は特に何かするわけじゃないけど、最終的には現在のステータスはアップロード完了の状態です、みたいに表示することになるだろうから。
> // HTTPフックサーバーを設定
> hookHandler := hooks.NewHandler()
> hookHandler.PostFinish = func(event hooks.HookEvent) {
>     // アップロード完了時の処理
>     // ファイルID、パス、サイズなどの情報が取得可能
> }

要件定義書を承認します。

/kiro:spec-design 0032-movie-upload

config/staging/config.yamlを追加して。
> #### 2.2.1 変更ファイル
> - **サーバー側**:
>   - `server/internal/api/router/router.go`: TUSエンドポイントの登録を追加
>   - `server/internal/config/config.go`: `UploadConfig`構造体を追加
>   - `config/develop/config.yaml`: アップロード設定を追加
>   - `config/production/config.yaml.example`: アップロード設定を追加


#### 決定2: 環境別ストレージ対応
- **コンテキスト**: 開発環境と本番環境で異なるストレージが必要
- **代替案**:
  1. 常にS3を使用（開発環境でもS3が必要）
  2. 常にローカルファイルシステム（本番環境でスケーラビリティの問題）
- **選択アプローチ**: 環境に応じてストレージを切り替え（開発: ローカル、本番: S3）
- **根拠**: 開発環境での簡易性と本番環境でのスケーラビリティを両立
- **トレードオフ**: ストレージ抽象化の実装コスト vs 環境別最適化の利点


疑問：これはどういう話？
今の想定だとサーバーでシングルトンで、1件のみ想定しているってこと？
それとも複数アップロードしようとした時、区別する方法がないってこと？
> ### 11.2 複数ファイルの同時アップロード
> 現時点では1ファイルずつだが、将来的に複数ファイルの同時アップロードに対応可能。

つまり、ここで言う、複数アップロードってのは、
Aファイル＋Bファイル＋Cファイル、みたいに1回で3ファイル送り込むような処理のことを言っている？
Aファイル→Bファイル→Cファイルでなくて。

Aファイル＋Bファイル＋Cファイルのパターンは対応しないよ。
対象が動画ファイルで、巨大だからね。
> **将来の拡張**: 将来的には、クライアント側のUIを拡張して複数ファイルの並列アップロード（Aファイル＋Bファイル＋Cファイルを同時にアップロード）に対応可能です。この場合、複数のuppyインスタンスまたはtus-js-clientインスタンスを使用して、各ファイルを独立したTUSセッションでアップロードします。サーバー側の変更は不要です（既に対応可能

設計書を承認します。

/kiro:spec-tasks 0032-movie-upload


allowed_extensionsはmp4だけにしておいて。
> - `allowed_extensions: ["mp4", "mov", "avi", "mkv"]` (動画ファイル形式の例)

max_file_sizeは 2GBにしておいて。
>  - `max_file_size: 10737418240` (10GB、例)

タスクリストを承認します。

/sdd-fix-plan

_serena_indexing

/serena-initialize

/kiro:spec-impl 0032-movie-upload




