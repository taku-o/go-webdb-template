/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/130
に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0063-client2としてください。"
think.

client2のポートは選べるなら3000にして。

要件定義書を承認します。

/kiro:spec-design 0063-client2

コンポーネントってどんなのがある？
>**インストールするコンポーネント（必要最小限）**:

おそらく現clientの移行には次のコンポーネントが必要かな？
これらを入れてください。
* Alert Dialog
* Alert
* Button
* Select
* Input
* Form
* Field
* Card

Prismaは要らないから、消す作業は入れておきましょう。
>2. **不要な機能の削除（必要に応じて）**
>   - Prisma関連のファイルを削除または無効化
>   - 不要なコンポーネントを削除

これでイケルかな？
設計書を承認します。

/kiro:spec-tasks 0063-client2


今回は移行時だからREADMEの更新は無しにしよう。
代わりに、docs/Temp-Client2.md を作成。
client -> client2に移行が完成したら、docs/Temp-Client2.mdの内容をREADMEに移植する想定で。
>#### タスク 5.5: 基本的なREADMEの作成
>**目的**: プロジェクトの基本的なREADMEを作成する

タスクリストを承認します。

/sdd-fix-plan

/kiro:spec-impl 0063-client2 1.1
/kiro:spec-impl 0063-client2 1.2
/kiro:spec-impl 0063-client2 1.3

ここでいったんcommitしましょう。




