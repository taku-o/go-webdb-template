/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/27 に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0015-cloudbeaverとしてください。"
think.

要件定義書を作成してください

起動コマンドnpm run cloudbeaver:startを実行した時、develop、staging、production。
この環境の違いはどうやって制御する？環境変数？

要件定義書を承認します。

/kiro:spec-design 0015-cloudbeaver

cross-env はstaging、本番環境で入れられなそうだから、
cross-envの使用の想定は除外しよう。

CloudBeaverの設定ファイルもgit管理したい。
develop、staging、productionで環境分かれてるけど、なんとかなるだろうか？
- CloudBeaverの設定ファイルに接続情報が保存される

設計書を承認します。

/kiro:spec-tasks 0015-cloudbeaver

README.md以外に
専用のドキュメントを作って欲しい。
Database-Viewer.md

タスクリストを承認します。

この要件の作業用のgitブランチを切ってください。
ここまでの作業をcommitしてください。







