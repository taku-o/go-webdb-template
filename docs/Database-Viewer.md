# データベースビューアー（CloudBeaver）ドキュメント

## 概要

CloudBeaverは、Webベースのデータベース管理ツールです。本プロジェクトでは、データ操作用の管理アプリとして導入されています。

### 主な機能

- **Webベースのデータベース管理**: ブラウザからデータベースにアクセス可能
- **視覚的なデータ操作**: テーブル構造の確認、データの閲覧・編集が容易
- **SQL実行**: SQLクエリを実行してデータを操作
- **Resource Manager**: よく使うSQLスクリプトを保存・管理
- **環境別設定**: develop、staging、production環境で設定を分離

### 役割分担

本プロジェクトでは管理系アプリを2つ用意しています：

- **CloudBeaver**: データ操作用（本ツール）
- **GoAdmin**: カスタム処理用（`docs/Admin.md`を参照）

## 前提条件

- DockerおよびDocker Composeがインストールされていること
- Dockerが正常に動作していること
- ポート8978が使用可能であること（他のサービスと競合しないこと）
- データベースファイルが`server/data/`ディレクトリに存在すること

## 起動方法

### 基本的な起動

```bash
# 開発環境（デフォルト）
npm run cloudbeaver:start

# または明示的に環境を指定
APP_ENV=develop npm run cloudbeaver:start
```

### 環境別の起動

```bash
# 開発環境
APP_ENV=develop npm run cloudbeaver:start

# ステージング環境
APP_ENV=staging npm run cloudbeaver:start

# 本番環境
APP_ENV=production npm run cloudbeaver:start
```

**注意**: 環境変数`APP_ENV`が未設定の場合は、デフォルトで`develop`環境として起動します。

### 起動確認

起動後、以下のURLにアクセスしてCloudBeaverのWeb UIが表示されることを確認してください：

- **URL**: http://localhost:8978

## 停止方法

```bash
npm run cloudbeaver:stop
```

## その他のコマンド

### ログ確認

```bash
npm run cloudbeaver:logs
```

### 再起動

```bash
npm run cloudbeaver:restart
```

## データベース接続設定

### 初回設定

CloudBeaver初回起動時には、管理者アカウントの設定とドライバーの有効化が必要です。

#### 1. 管理者アカウントの作成

1. http://localhost:8978 にアクセス
2. セットアップウィザードが表示される
3. 管理者アカウントを作成：
   - **ユーザー名**: `cbadmin`
   - **パスワード**: `Admin123`
4. 「Next」→「Finish」をクリックして設定を終える

#### 2. SQLiteドライバーの有効化

初期状態ではSQLiteドライバーが無効化されています。以下の手順で有効化してください。

1. 管理者アカウントでログイン
2. 左上のメニュー（≡）→「Administration」→「Server Configuration」を開く
3. 「DISABLED DRIVERS」セクションを確認
4. SQLiteが含まれている場合、SQLiteを選択して削除（有効化）
5. 「Save」をクリックして設定を保存

設定が終わったら、Web UIから手動でデータベース接続を設定します。

### マスターデータベースへの接続

1. CloudBeaver Web UI（http://localhost:8978）にアクセス
2. 「接続を追加」または「New Connection」をクリック
3. データベースタイプで「SQLite」を選択
4. 接続情報を入力：
   - **接続名**: `master` または `Master Database`
   - **データベースファイル**: `/data/master.db`
5. 「接続をテスト」をクリックして接続を確認
6. 「保存」をクリックして接続設定を保存

### シャーディングデータベースへの接続

同様の手順で、以下の4つのシャーディングデータベースに接続を追加します：

| 接続名 | データベースファイル |
|--------|---------------------|
| `sharding_db_1` | `/data/sharding_db_1.db` |
| `sharding_db_2` | `/data/sharding_db_2.db` |
| `sharding_db_3` | `/data/sharding_db_3.db` |
| `sharding_db_4` | `/data/sharding_db_4.db` |

### 接続設定の保存場所

接続設定は環境別の設定ディレクトリに保存されます：

- **開発環境**: `cloudbeaver/config/develop/`
- **ステージング環境**: `cloudbeaver/config/staging/`
- **本番環境**: `cloudbeaver/config/production/`

設定ファイルはGitで管理可能です。環境別に設定が分離されるため、各環境で個別に接続設定を行う必要があります。

## データベース操作

### テーブル一覧の表示

1. 接続したデータベースを選択
2. 左側のナビゲーションツリーから「Tables」を展開
3. テーブル一覧が表示されます

### データの閲覧

1. テーブルを選択
2. 「Data」タブをクリック
3. テーブルのデータが表示されます

### SQLクエリの実行

1. 接続したデータベースを選択
2. 「SQL Editor」タブをクリック
3. SQLクエリを入力
4. 「Execute」ボタンをクリックして実行
5. 結果が表示されます

**注意**: データベースファイルは読み取り専用でマウントされているため、データの変更はできません。データの変更が必要な場合は、既存のAPIや管理ツールを使用してください。

## Resource Manager

Resource Managerは、よく使うSQLスクリプトを保存・管理する機能です。

### スクリプトの作成

1. CloudBeaver Web UIで「Resource Manager」を開く
2. 「New Script」をクリック
3. スクリプト名とSQLクエリを入力
4. 「Save」をクリックして保存

### スクリプトの保存場所

Resource Managerに保存したスクリプトは、ユーザープロジェクトディレクトリに保存されます：

- `cloudbeaver/config/{env}/user-projects/{username}/`

例（開発環境、cbadminユーザーの場合）：
- `cloudbeaver/config/develop/user-projects/cbadmin/sql-1.sql`

スクリプトファイルはGitで管理可能です。

### スクリプトの使用

1. Resource Managerからスクリプトを選択
2. 「Execute」ボタンをクリックして実行
3. または、SQL Editorにスクリプトを読み込んで実行

### スクリプトの編集・削除

- **編集**: Resource Managerでスクリプトを選択し、「Edit」をクリック
- **削除**: Resource Managerでスクリプトを選択し、「Delete」をクリック

## 環境別設定

### 設定ディレクトリの構造

```
cloudbeaver/
├── config/
│   ├── develop/                    # 開発環境用設定
│   │   ├── GlobalConfiguration/    # 接続設定など
│   │   └── user-projects/          # ユーザー別スクリプト
│   │       └── cbadmin/            # cbadminユーザーのスクリプト
│   ├── staging/                    # ステージング環境用設定
│   └── production/                 # 本番環境用設定
└── scripts/                        # 共有スクリプト用（予備）
```

### 環境別設定の管理

- 各環境でCloudBeaverを起動すると、対応する設定ディレクトリがマウントされます
- 接続設定は環境別に保存されるため、環境ごとに異なる接続設定を管理できます
- 設定ファイルはGitで管理可能です

### 設定ファイルの確認

設定ファイルは以下のディレクトリに保存されます：

- 開発環境: `cloudbeaver/config/develop/`
- ステージング環境: `cloudbeaver/config/staging/`
- 本番環境: `cloudbeaver/config/production/`

## トラブルシューティング

### コンテナが起動しない

**問題**: `npm run cloudbeaver:start`を実行してもCloudBeaverが起動しない

**対処方法**:
1. Dockerが起動しているか確認
   ```bash
   docker ps
   ```
2. ポート8978が使用されていないか確認
   ```bash
   lsof -i :8978
   ```
3. ログを確認
   ```bash
   npm run cloudbeaver:logs
   ```
4. ポート番号を変更して起動（ポート競合の場合）
   ```bash
   CLOUDBEAVER_PORT=8979 npm run cloudbeaver:start
   ```

### データベースに接続できない

**問題**: CloudBeaverからデータベースに接続できない

**対処方法**:
1. データベースファイルが存在するか確認
   ```bash
   ls -la server/data/*.db
   ```
2. マウント設定を確認
   - `docker-compose.yml`の`volumes`セクションを確認
   - `./server/data:/data:ro`が正しく設定されているか確認
3. ファイルのパスを確認
   - CloudBeaverでのデータベースファイルパスは`/data/master.db`など
   - コンテナ内のマウントパス（`/data`）を使用すること

### Resource Managerにスクリプトを保存できない

**問題**: Resource Managerにスクリプトを保存できない

**対処方法**:
1. ユーザープロジェクトディレクトリの権限を確認
   ```bash
   ls -la cloudbeaver/config/develop/user-projects/
   ```
2. マウント設定を確認
   - `docker-compose.yml`の`volumes`セクションを確認
   - `./cloudbeaver/config/${APP_ENV:-develop}:/opt/cloudbeaver/workspace`が正しく設定されているか確認

### 設定ファイルが保存されない

**問題**: 接続設定が保存されない、または環境別に分離されない

**対処方法**:
1. 環境変数`APP_ENV`が正しく設定されているか確認
   ```bash
   echo $APP_ENV
   ```
2. 設定ディレクトリがマウントされているか確認
   - `docker-compose.yml`の`volumes`セクションを確認
   - `./cloudbeaver/config/${APP_ENV:-develop}:/opt/cloudbeaver/workspace`が正しく設定されているか確認
3. 設定ディレクトリが存在するか確認
   ```bash
   ls -la cloudbeaver/config/develop/
   ```

### ポートが競合する

**問題**: ポート8978が既に使用されている

**対処方法**:
1. 使用しているプロセスを確認
   ```bash
   lsof -i :8978
   ```
2. ポート番号を変更して起動
   ```bash
   CLOUDBEAVER_PORT=8979 npm run cloudbeaver:start
   ```
3. `docker-compose.yml`でポート番号を変更（恒久的な変更が必要な場合）

## セキュリティ考慮事項

### データベースファイルへのアクセス

- データベースファイルは読み取り専用でマウントされています（`:ro`オプション）
- CloudBeaverからデータベースファイルを誤って変更することを防止しています
- データの変更が必要な場合は、既存のAPIや管理ツールを使用してください

### 認証設定

**認証情報**（開発環境）:
- ユーザー名: `cbadmin`
- パスワード: `Admin123`

**注意事項**:
- 本番環境での使用は想定していません（本番環境では適切なアクセス制御が必要です）

### ネットワークアクセス

- CloudBeaverはローカルホスト（localhost）でのみアクセス可能です
- 外部からのアクセスは想定していません

## 設定ファイルの管理

### Git管理

CloudBeaverの設定ファイルはGitで管理可能です：

- **設定ファイル**: `cloudbeaver/config/{env}/`
- **接続設定**: `cloudbeaver/config/{env}/GlobalConfiguration/`
- **スクリプト**: `cloudbeaver/config/{env}/user-projects/{username}/`

### 設定ファイルの構造

設定ファイルは環境別に分離されています：

- `cloudbeaver/config/develop/`: 開発環境用設定
- `cloudbeaver/config/staging/`: ステージング環境用設定
- `cloudbeaver/config/production/`: 本番環境用設定

各環境でCloudBeaverを起動すると、対応する設定ディレクトリがマウントされ、接続設定などが保存されます。

### 設定ファイルの共有

設定ファイルをGitで管理することで、チームメンバー間で設定を共有できます。ただし、機密情報（パスワードなど）が含まれる場合は、適切に管理してください。

## 参考情報

### 関連ドキュメント

- `README.md`: プロジェクト概要とセットアップ手順
- `docs/Admin.md`: GoAdmin管理画面のドキュメント
- `docs/Sharding.md`: シャーディングの詳細仕様

### 技術スタック

- **CloudBeaver**: https://cloudbeaver.io/
- **Docker**: Docker Composeを使用
- **データベース**: SQLite（開発環境）

### 参考リンク

- CloudBeaver公式サイト: https://cloudbeaver.io/
- CloudBeaver GitHub: https://github.com/dbeaver/cloudbeaver
- CloudBeaver Docker: https://hub.docker.com/r/dbeaver/cloudbeaver
- CloudBeaver ドキュメント: https://cloudbeaver.io/docs/

