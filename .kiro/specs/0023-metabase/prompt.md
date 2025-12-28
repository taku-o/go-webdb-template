/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/47 に対応するための要件を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0023-metabaseとしてください。"
think.

Metabaseのportは3000ではなく、別のportを使用する設定にしておきましょう。
8970あたり。

同時起動の警告機能は不要。
> - **同時起動の防止**: オプションとして、CloudBeaverとMetabaseの同時起動を防止する機能を実装（設計フェーズで決定

こちらもやらなくて良い。
> - 運用上の制約（同時起動の制限）を記載

いろいろ作業する前に、今動いているCloudBeaverは止めてしまいましょう。

要件定義書を承認します。

/kiro:spec-design 0023-metabase

設計書を承認します。

Metabase用のドキュメントも作成したいけど、
Database-Viewer.mdとは別のファイルとして用意したい。
良いファイル名の案はあるかな？
Metabase.md でも良いけど。

Metabase.mdにしよう。

METABASE_PORT環境変数でportを管理するより、設定ファイルか、それが難しければ起動スクリプトに直書きでport番号を定義したい。
使用portが被るのを回避するのは運用者の責任。


これはどういうこと？
つまり、どこかのサーバーに置くstaging、production環境では参照できない想定ってこと？
> #### 7.2.2 ネットワークアクセス
> - Metabaseはローカルホスト（localhost）でのみアクセス可能
> - 外部からのアクセスは想定しない


これは何に使われる設定？
> `MB_DB_FILE`環境変数

これは環境変数でなくて、設定ファイルか、起動スクリプトに直書きする想定にしましょう。
極力、環境変数は使わない。
> `MB_DB_FILE`環境変数

設計書を承認します。

/kiro:spec-tasks 0023-metabase

Metabaseは起動にかなり時間がかかるらしい。
ヘルスチェックのstart_period 60sでは、足りない可能性がかなりある。

タスクリストを承認します。

この要件の作業用のgitブランチを切ってください。
ここまでの作業をcommitしてください。

_serena_indexing

/serena-initialize

/kiro:spec-impl 0023-metabase



metaadmin /



