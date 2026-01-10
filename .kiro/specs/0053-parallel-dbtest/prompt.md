/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/109
に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0053-parallel-dbtestとしてください。"
think.

os.TempDir()はプロセスが違うと、
違うディレクトリにならない？

gofrs/flockライブラリを導入する。
イメージとしては、こんな感じの処理を差し込みたい。

```
// ロックファイルのパスを指定
fileLock := flock.New("test_db.lock")

// ロックを取得（取得できるまでブロックする）
err := fileLock.Lock()
if err != nil {
    // ロック取得自体の失敗（権限不足など）のハンドリング
    panic(err)
}

defer fileLock.Unlock()
```

30秒待ってもロックが取れなければ
タイムアウトエラーにしちゃおう。


エラーメッセージには、{ロックファイルPATH}のロックが取れなかったので
タイムアウトしました。という感じのメッセージを返して、
ロックファイル消し忘れ時に気づけるようにしよう。

ロックファイルはあえて、.gitignoreに入れないで、
git statusした時にロックファイルの消し忘れに気づけるようにする。

.test-lock/test-db.lock -> test-db.lock
としたい。

要件定義書を承認します。

/kiro:spec-design 0053-parallel-dbtest

設計書を承認します。

/kiro:spec-tasks 0053-parallel-dbtest

docs/Testing.md にロックファイルの記載を追加したい。
書き込むドキュメントはここが適切か？

タスクリストにドキュメント修正の作業を追加してください。

タスクリストを承認します。

/sdd-fix-plan

_serena_indexing

/serena-initialize

/kiro:spec-impl 0053-parallel-dbtest

監視してるから、並列テストを3回ほど実行してくれますか？

server/test-db.lockってファイルが残ってる。

削除しないけど、ロックは取れてたの？

ファイルのあるなしじゃなくて、ファイルを握るかどうかで制御してるのか。

とすると、ファイルを消すのは都合が悪い？
ロスが大きい？たいしたことない？


ロックファイルを残すように変更した場合は、
ドキュメントの記述が違ってしまうのと、
目立たない場所にロックファイルを置きたい、というのがある。

ロックファイルの位置は
.test-lock/test-db.lock
にする。

修正計画を建てて

think.

ロックファイルを消さない方針にしたから、
.test-lock/test-db.lock は
.gitignoreに追加して、gitから見えないようにしよう。

修正作業お願いします。
think.

commitした後、
https://github.com/taku-o/go-webdb-template/issues/109
に対してpull requestを発行してください。

/review 110





