/kiro:steering
steeringを更新する必要があるなら更新してください。
think.

1. このプロジェクトのコードを分析し、
2. 次にdocsディレクトリのドキュメントを確認
3. docsディレクトリに足すべきファイルがあるなら、提案してください。
think.

この3つのファイルを生成したいです。
> File-Upload.md
> Logging.md
> Rate-Limit.md

このファイルの記述を修正したいです。
> docs/API.md
> 385-387行目: レートリミット機能について「Currently not implemented」と記載されているが、実際には実装済み。この部分を削除または更新する必要がある

think.

README.mdファイルも情報が古かったり、不足しているように感じる。
何か追加か、更新。あるいは思い切って削除すべき情報などあるなら提案して。
think.

そのREADME.mdの更新計画、すべて承認します。
README.mdを修正してください。


README.mdのAPIエンドポイントの記述が間違っているかも。
いくつかのエンドポイントは、dm_のプレフィックスをつけてある。
client/src/app/
think.

今、masterブランチで作業しているので、作業用ブランチを切りたいです。


docs/System-Configuration.drawio にシステム全体の図を
Draw.ioで書いています。
見た目を整えられますか？
think.

                 - [Relying Party (Auth0)] - [Identity Provider (Auth0 -> Partner)]

[Next.js] - [Go] - [GoAdmin]
                 - [CloudBeaver]
                 - [Metabase]
                 - [Redis Insight]

                 - [SQLite / PostgreSQL / MySQL]
                 - [Redis Standalone (データ永続化)]
                 - [Redis Cluster]

                 - [AWS SES / Simple Email Service]
                 - [AWS S3 / Tencent Cloud Cloud Object Storage]
                 - [AWS Elemental MediaConvert / Tencent Cloud Media Processing Service]


生成して貰った図を、次のように修正。
いくつか図形を削除して、スペースが出来たので、配置を調整してください。

Auth0 Identity Provider →  Identity Provider (Auth0 -> Partner Idp)
Partner Idp 削除
Mailpit 削除
Shard DB 1 →  Shard DB 1-4
Shard DB 2-4 削除
AWS SES →  AWS SES/Tencent Cloud SES
AWS S3 →　AWS S3/Tencent Cloud COS
AWS MediaConvert →  AWS MediaConvert/Tencent Cloud MPS
Tencent Cloud SES 削除
Tencent Cloud COS 削除
Tencent Cloud MPS 削除

docs/System-Configuration.drawio の図の情報から、
Port番号の情報は削除してください。
think.

docs/System-Configuration.drawio
GoAdminをサーバー層から、管理・可視化ツールに移動してください。
think.

ここまでの修正をcommitしてください。
その後、https://github.com/taku-o/go-webdb-template に対して
pull requestを作成してください。





