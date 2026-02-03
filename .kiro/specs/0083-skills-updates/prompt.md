
Claude Codeで使用しているSKILLS(.claude/skills/)の見直しを行いたい。

かなり古い時期に作っているため、プロジェクトの現状と合わない可能性がある。
もし不適切な状態であれば修正したい。

また不足していると考えられるSKILLがあるなら追加したい。

あとは次のようなSKILLを追加したい。
* クライアントのコードでuseEffectを使用しようとしていたら、本当にuseEffectを使うべき場所かどうか検討を促す。
* テストの実行で認証エラーが起きた時、APP_ENV=testを指定していない可能性を指摘する。

まずはSKILLSとプロジェクトの状況を調査、
どのように修正するべきかの計画を建ててください。
think.


/kiro:spec-requirements "作って貰ったPLANを元に、作業を行うための要件定義書を作成してください。
cc-sddのfeature名は0083-skills-updatesとします。
.claude/skills/SKILLS_REVIEW_PLAN.mdは、要件定義書を入れることになるディレクトリ、 .kiro/specs/0083-skills-updates/ に移動しましょう。"
think.

要件定義書・設計書・タスクリストのフォーマットは
.kiro/specs/0023-metabase/*.md に合わせて欲しい。

手順通り作業を進めます。
要件定義書の承認前なので、
設計書・タスクリストはいったん削除してください。

要件定義書に少し修正を入れました。
要件定義書を承認します。

/kiro:spec-design 0083-skills-updates

設計書を承認します。

/kiro:spec-tasks 0083-skills-updates

一部の「必要に応じて」の記載は、作業が必要なので「必要に応じて」の部分を消しました。

タスクリストを承認します。

/sdd-fix-plan

/kiro:spec-impl 0083-skills-updates






