# Metabase - データ可視化・分析ツール

## 概要

Metabaseは非エンジニア向けのデータビューワ・分析ツールです。Webブラウザからデータベースに接続し、クエリの作成・実行、ダッシュボードの作成・共有が可能です。

### 管理ツールの役割分担

本プロジェクトでは3つの管理ツールを用意しています：

| ツール | 役割 | ポート |
|--------|------|--------|
| GoAdmin | カスタム処理用の管理画面 | 8081 |
| CloudBeaver | データ操作用のWebベースツール | 8978 |
| Metabase | データ可視化・分析用ツール | 8970 |

**重要**: CloudBeaverとMetabaseはメモリ使用量が大きいため、開発環境では片方ずつしか起動しない運用を推奨します。

## 前提条件

- DockerおよびDocker Composeがインストールされていること
- Dockerが正常に動作していること
- ポート8970が使用可能であること
- 十分なメモリがあること（Metabaseはメモリを多く使用します）

## 起動方法

### 基本的な起動（開発環境）

```bash
# 開発環境（デフォルト）
npm run metabase:start
```

### 環境別の起動

```bash
# 開発環境を明示的に指定
APP_ENV=develop npm run metabase:start

# ステージング環境
APP_ENV=staging npm run metabase:start

# 本番環境
APP_ENV=production npm run metabase:start
```

起動後、http://localhost:8970 にアクセスしてください。

## 停止方法

```bash
npm run metabase:stop
```

## その他のコマンド

### ログ確認

```bash
npm run metabase:logs
```

### 再起動

```bash
npm run metabase:restart
```

## データベース接続設定

### 初回設定

1. Metabase起動後、http://localhost:8970 にアクセス
2. 初回起動時は初期設定画面が表示されます
3. 管理者アカウントを作成:
   - 名前
   - メールアドレス
   - パスワード
4. 初期設定を完了

### 開発環境の管理者アカウント

開発環境では以下の管理者アカウントが設定されています：

| 項目 | 値 |
|------|-----|
| メールアドレス | `admin@example.com` |
| パスワード | `metaadmin123` |

### マスターデータベースの接続

1. 管理画面からデータベース追加を選択
2. PostgreSQLを選択
3. 接続情報を入力:
   - **接続名**: `master` または `Master Database`
   - **ホスト**: `postgres-master`（Docker内）または `localhost`
   - **ポート**: `5432`
   - **データベース名**: `webdb_master`
   - **ユーザー名**: `webdb`
   - **パスワード**: `webdb`
4. 接続をテスト
5. 接続を保存

### シャーディングデータベースの接続

各シャーディングデータベースに対して接続を追加します：

| 接続名 | ホスト | ポート | データベース名 |
|--------|--------|--------|--------------|
| `sharding_db_1` | `postgres-sharding-1` | 5433 | `webdb_sharding_1` |
| `sharding_db_2` | `postgres-sharding-2` | 5434 | `webdb_sharding_2` |
| `sharding_db_3` | `postgres-sharding-3` | 5435 | `webdb_sharding_3` |
| `sharding_db_4` | `postgres-sharding-4` | 5436 | `webdb_sharding_4` |

## クエリ作成方法

1. 左メニューから「新規」→「質問」を選択
2. データベースとテーブルを選択
3. フィルター、集計、グループ化などを設定
4. 「可視化」ボタンでグラフ形式を選択
5. 「保存」でクエリを保存

### ネイティブクエリ（SQL）

1. 「新規」→「SQL クエリ」を選択
2. データベースを選択
3. SQLクエリを入力
4. 「実行」でクエリを実行
5. 「保存」でクエリを保存

## ダッシュボード作成方法

1. 左メニューから「新規」→「ダッシュボード」を選択
2. ダッシュボード名を入力
3. 保存したクエリを追加
4. レイアウトを調整
5. 「保存」でダッシュボードを保存

## CloudBeaverとMetabaseの使い分け

| 用途 | 推奨ツール |
|------|-----------|
| データの直接編集・操作 | CloudBeaver |
| テーブル構造の確認 | CloudBeaver |
| SQLスクリプトの管理 | CloudBeaver |
| データの可視化・グラフ作成 | Metabase |
| ダッシュボードの作成・共有 | Metabase |
| 非エンジニアへのデータ共有 | Metabase |
| レポート作成 | Metabase |

## 設定ファイルの管理

Metabaseの設定ファイルは環境別に管理されます：

```
metabase/
└── config/
    ├── develop/      # 開発環境用設定
    ├── staging/      # ステージング環境用設定
    └── production/   # 本番環境用設定
```

- 接続設定やダッシュボード設定は各環境のディレクトリに保存されます
- 設定ファイルはGitで管理可能です
- 環境ごとに異なる接続設定を管理できます

## トラブルシューティング

### コンテナが起動しない

1. Dockerが起動しているか確認
   ```bash
   docker ps
   ```

2. ポート8970が使用されていないか確認
   ```bash
   lsof -i :8970
   ```

3. ログを確認
   ```bash
   npm run metabase:logs
   ```

### データベースに接続できない

1. データベースファイルが存在するか確認
   ```bash
   ls -la server/data/
   ```

2. 接続設定のパスを確認
   - 正しいパス: `/data/master.db`
   - 誤ったパス: `server/data/master.db`

### 設定ファイルが保存されない

1. 設定ディレクトリの権限を確認
   ```bash
   ls -la metabase/config/
   ```

2. ディレクトリが存在するか確認
   ```bash
   mkdir -p metabase/config/develop
   mkdir -p metabase/config/staging
   mkdir -p metabase/config/production
   ```

### メモリ不足

Metabaseはメモリを多く使用します。以下を確認してください：

- CloudBeaverが停止していること
- 十分なシステムメモリがあること
- 他のメモリを多く使用するアプリケーションを停止

## 技術仕様

### Docker Compose設定

- **イメージ**: `metabase/metabase:latest`
- **ポート**: 8970 (ホスト) → 3000 (コンテナ)
- **ボリューム**:
  - `./server/data:/data:ro` - データベースファイル（読み取り専用）
  - `./metabase/config/${APP_ENV}:/metabase-data` - 設定ファイル

### 環境変数

| 変数名 | 説明 | デフォルト値 |
|--------|------|-------------|
| `APP_ENV` | 環境名 | `develop` |
| `MB_DB_FILE` | Metabase内部DB | `/metabase-data/metabase.db` |

## 関連ドキュメント

- [README.md](../README.md) - プロジェクト概要
- [Database-Viewer.md](Database-Viewer.md) - CloudBeaverの詳細仕様
- [Sharding.md](Sharding.md) - シャーディングの詳細仕様
- [Admin.md](Admin.md) - GoAdmin管理画面の詳細

## 参考リンク

- [Metabase公式サイト](https://www.metabase.com/)
- [Metabase ドキュメント](https://www.metabase.com/docs/)
- [Metabase Docker ドキュメント](https://www.metabase.com/docs/latest/installation-and-operation/running-metabase-on-docker)
