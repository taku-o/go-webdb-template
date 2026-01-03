# Apache Superset - データ可視化・分析ツール

## 概要

Apache Supersetは非エンジニア向けのデータビューワ・分析ツールです。Webブラウザからデータベースに接続し、クエリの作成・実行、ダッシュボードの作成・共有が可能です。

### 管理ツールの役割分担

本プロジェクトでは4つの管理ツールを用意しています：

| ツール | 役割 | ポート |
|--------|------|--------|
| GoAdmin | カスタム処理用の管理画面 | 8081 |
| CloudBeaver | データ操作用のWebベースツール | 8978 |
| Metabase | データ可視化・分析用ツール | 8970 |
| Apache Superset | データ可視化・分析用ツール | 8088 |

**重要**: CloudBeaver、Metabase、Apache Supersetはメモリ使用量が大きいため、開発環境では片方ずつしか起動しない運用を推奨します。

## 前提条件

- DockerおよびDocker Composeがインストールされていること
- Dockerが正常に動作していること
- ポート8088が使用可能であること
- PostgreSQLが起動していること（`docker-compose.postgres.yml`で起動）
- 十分なメモリがあること（Apache Supersetはメモリを多く使用します）

## 起動方法

### PostgreSQLの起動

Apache Supersetを使用するには、まずPostgreSQLを起動する必要があります：

```bash
# PostgreSQLを起動
./scripts/start-postgres.sh start
```

### Apache Supersetの起動

```bash
# Apache Supersetを起動
./scripts/start-apache-superset.sh
```

起動後、http://localhost:8088 にアクセスしてください。

## 停止方法

```bash
# Apache Supersetを停止
docker-compose -f docker-compose.apache-superset.yml down
```

## その他のコマンド

### ログ確認

```bash
docker-compose -f docker-compose.apache-superset.yml logs -f
```

### 再起動

```bash
# 停止
docker-compose -f docker-compose.apache-superset.yml down

# 起動
./scripts/start-apache-superset.sh
```

## データベース接続設定

### 初回設定

1. Apache Superset起動後、http://localhost:8088 にアクセス
2. 初回起動時は初期設定画面が表示されます
3. デフォルトの管理者アカウントでログイン:
   - **ユーザー名**: `admin`
   - **パスワード**: `admin`
4. 初回ログイン時にパスワード変更を求められる場合があります（開発環境では`admin`のままでも可）

### PostgreSQLデータベースの接続

1. 管理画面から「Settings」→「Database Connections」を選択
2. 「+ Database」ボタンをクリック
3. データベースタイプで「PostgreSQL」を選択
4. 接続情報を入力:
   - **Display Name**: `PostgreSQL - webdb` または任意の名前
   - **Host**: `postgres`（Dockerネットワーク内の場合）または`host.docker.internal`（ホストマシンのPostgreSQLに接続する場合）
   - **Port**: `5432`
   - **Database name**: `webdb`
   - **Username**: `webdb`
   - **Password**: `webdb`
5. 「Test Connection」をクリックして接続を確認
6. 「Connect」をクリックして接続設定を保存

**注意**: Dockerネットワーク内でPostgreSQLに接続する場合は、`postgres`をホスト名として使用します。ホストマシンから直接PostgreSQLに接続する場合は、`host.docker.internal`を使用します。

## クエリ作成方法

### SQL Labでのクエリ実行

1. 上部メニューから「SQL Lab」→「SQL Editor」を選択
2. データベースを選択
3. SQLクエリを入力
4. 「Run」ボタンでクエリを実行
5. 結果が表示されます
6. 「Save」でクエリを保存

### ビジュアライゼーションの作成

1. SQL Labでクエリを実行
2. 結果表示の下にある「Create Chart」ボタンをクリック
3. チャートタイプを選択（テーブル、棒グラフ、折れ線グラフ、円グラフなど）
4. データソース、メトリクス、ディメンションを設定
5. 「Create Chart」でチャートを作成
6. 「Save」でチャートを保存

## ダッシュボード作成方法

1. 上部メニューから「Dashboards」→「+ Dashboard」を選択
2. ダッシュボード名を入力
3. 保存したチャートを追加
4. レイアウトを調整
5. 「Save」でダッシュボードを保存

### 既存チャートの追加

1. ダッシュボード編集画面で「+ Add Chart」をクリック
2. 保存済みのチャートを選択
3. チャートがダッシュボードに追加されます

## データビューワの比較

本プロジェクトでは複数のデータビューワを用意しています。用途に応じて使い分けてください：

| 用途 | 推奨ツール |
|------|-----------|
| データの直接編集・操作 | CloudBeaver |
| テーブル構造の確認 | CloudBeaver |
| SQLスクリプトの管理 | CloudBeaver |
| データの可視化・グラフ作成 | Metabase / Apache Superset |
| ダッシュボードの作成・共有 | Metabase / Apache Superset |
| 非エンジニアへのデータ共有 | Metabase / Apache Superset |
| レポート作成 | Metabase / Apache Superset |
| 高度な可視化機能 | Apache Superset |

### MetabaseとApache Supersetの違い

| 項目 | Metabase | Apache Superset |
|------|----------|-----------------|
| ライセンス | AGPL v3 | Apache 2.0 |
| 商用利用時の制約 | ソースコード公開が必要 | 制約なし |
| SQL Lab機能 | あり | あり（より高度） |
| 可視化の種類 | 標準的なチャート | より多くのチャートタイプ |
| カスタマイズ性 | 中程度 | 高い |
| 学習曲線 | やや低い | やや高い |

## 設定ファイルの管理

Apache Supersetの設定ファイルは以下のディレクトリに保存されます：

```
apache-superset/
└── data/                    # データディレクトリ
    ├── superset.db          # Superset内部データベース
    ├── config/              # 設定ファイル
    └── uploads/             # アップロードファイル
```

- 設定ファイルやダッシュボード設定は`apache-superset/data/`ディレクトリに保存されます
- データはDockerボリュームに永続化され、コンテナ再起動後も保持されます
- 設定ファイルはGitで管理可能です（機密情報に注意）

## トラブルシューティング

### コンテナが起動しない

1. Dockerが起動しているか確認
   ```bash
   docker ps
   ```

2. ポート8088が使用されていないか確認
   ```bash
   lsof -i :8088
   ```

3. ログを確認
   ```bash
   docker-compose -f docker-compose.apache-superset.yml logs
   ```

### PostgreSQLに接続できない

1. PostgreSQLが起動しているか確認
   ```bash
   docker ps | grep postgres
   ```

2. PostgreSQLの接続情報を確認
   - ホスト: `postgres`（Dockerネットワーク内）または`host.docker.internal`（ホストマシンから）
   - ポート: `5432`
   - データベース名: `webdb`
   - ユーザー名: `webdb`
   - パスワード: `webdb`

3. Dockerネットワークを確認
   ```bash
   docker network ls
   docker network inspect <network_name>
   ```

### データが保存されない

1. データディレクトリの権限を確認
   ```bash
   ls -la apache-superset/data/
   ```

2. Dockerボリュームのマウントを確認
   - `docker-compose.apache-superset.yml`の`volumes`セクションを確認
   - `./apache-superset/data:/app/superset_home`が正しく設定されているか確認

### メモリ不足

Apache Supersetはメモリを多く使用します。以下を確認してください：

- CloudBeaverやMetabaseが停止していること
- 十分なシステムメモリがあること
- 他のメモリを多く使用するアプリケーションを停止

### 初回起動が遅い

Apache Supersetの初回起動時は、データベースの初期化やセットアップに時間がかかります。数分待ってからアクセスしてください。

## 技術仕様

### Docker Compose設定

- **イメージ**: `apache/superset:latest`（または安定版）
- **ポート**: 8088 (ホスト) → 8088 (コンテナ)
- **ボリューム**:
  - `./apache-superset/data:/app/superset_home` - データディレクトリ（永続化）

### 環境変数

| 変数名 | 説明 | デフォルト値 |
|--------|------|-------------|
| `SUPERSET_SECRET_KEY` | Supersetのシークレットキー | 自動生成 |
| `SUPERSET_CONFIG_PATH` | 設定ファイルパス | `/app/superset_home/config/` |

### データベース接続情報

| 項目 | 値 |
|------|-----|
| ホスト（Dockerネットワーク内） | `postgres` |
| ホスト（ホストマシンから） | `host.docker.internal` |
| ポート | `5432` |
| データベース名 | `webdb` |
| ユーザー名 | `webdb` |
| パスワード | `webdb` |

## セキュリティ考慮事項

### 認証設定

**認証情報**（開発環境）:
- ユーザー名: `admin`
- パスワード: `admin`

**注意事項**:
- 開発環境での使用を想定しています
- 本番環境では適切なパスワードポリシーとアクセス制御を実装してください
- 初回ログイン時にパスワード変更を推奨します

### ネットワークアクセス

- Apache Supersetはローカルホスト（localhost）でのみアクセス可能です
- 外部からのアクセスは想定していません
- 本番環境では適切なネットワーク設定とファイアウォール設定が必要です

### データベース接続

- PostgreSQLへの接続情報は機密情報です
- 本番環境では環境変数やシークレット管理システムを使用してください
- 接続情報をGitにコミットしないよう注意してください

## 関連ドキュメント

- [README.md](../README.md) - プロジェクト概要
- [Database-Viewer.md](Database-Viewer.md) - CloudBeaverの詳細仕様
- [Metabase.md](Metabase.md) - Metabaseの詳細仕様
- [Sharding.md](Sharding.md) - シャーディングの詳細仕様
- [Admin.md](Admin.md) - GoAdmin管理画面の詳細
- [License-Survey.md](License-Survey.md) - ライセンス調査結果

## 参考リンク

- [Apache Superset公式サイト](https://superset.apache.org/)
- [Apache Superset GitHub](https://github.com/apache/superset)
- [Apache Superset ドキュメント](https://superset.apache.org/docs/)
- [Apache Superset Docker ドキュメント](https://superset.apache.org/docs/installation/installing-superset-using-docker-compose)
