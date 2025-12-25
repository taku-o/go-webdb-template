/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/21 に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0011-public-api-keyとしてください。"

* 使用できそうなライブラリは躊躇なく採用して良いです。
* クライアント側も、APIキーを使用するように修正する必要があります。

この要件が入った結果、テストの実行で問題になったりする？
> staging、productionのpublic APIキーは.gitignoreを設定して、commit出来ないようにする。

対策は、これにしたい。
テスト用のダミーAPIキーをtestdata/に配置

要件定義書を作成してください。

JWTのversionにセットする値はどこから持ってくる？
どこかに設定を用意する？
invalid_versionsの次の番号を使用する？


GoAdminキー発行ページには、
* JWTのペイロードも表示したい。
* JWTのダウンロードボタンも用意したい。

GoAdminキー発行ページは、
管理画面のメニューに追加したい。


> 8.4.1 メニュー項目の追加
> `parent_id`: 既存の「カスタムページ」カテゴリのID、または新しいカテゴリを作成
既存の「カスタムページ」カテゴリに入れてくれると助かる。


他に判断が必要な箇所はあるかな？

1. JWTの有効期限
は、少なくともpublic APIキーは無期限にしたい。
privateなAPIキーに有効期限を設けるなら、publicなAPIキーも有効期限が必要かな？
問題ないなら、public APIキーは、exp未定義で。

2. GoAdminキー発行ページのパス
/api-key にしましょう。

3. スコープの設定方法
固定（["read", "write"]のみ）

4. 秘密鍵の生成方法
複数のWebサーバーで動作することになるだろうから、
・コマンドラインツールで生成して、
・手動設定（設定ファイルに直接記述）
にするしかないか？

5. エラーレスポンスの形式
例: {"code": 401, "message": "..."}
こちらの方が嬉しいです。

6. クライアント側のAPIキー未設定時の動作
エラーを投げてリクエストを送信しない

7. 認証ミドルウェアの適用範囲
今は、
/api/*に適用
としましょう。

OKです。要件定義書を承認します。

/kiro:spec-design 0011-public-api-key

572行目、URLが
/admin/api-key?download=true&token=%s
となっているけど、もしかして、
/api-key?download=true&token=%s
ではない？

そうなんだ。じゃ、/api-key になっているURLは、/admin/api-key に揃えた方が良いね。
要件定義書のURLがいくつか /api-key になっているよ。修正して。

design.mdの785行目、
current_version の初期値はv2にしておこう。
そして、invalid_versions の初期値をv1にする。

OKです。設計書を承認します。

/kiro:spec-tasks 0011-public-api-key

OKです。タスクリストを承認します。

この要件の作業用のgitブランチを切ってください。
ここまでの作業をcommitしてください。



