# CloudBeaver導入設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、CloudBeaverをDockerで動作させ、Webベースのデータベース管理ツールを提供するシステムの詳細設計を定義する。既存システム（GoAdmin、Atlas）と共存し、データ操作用の管理アプリとして機能する。

### 1.2 設計の範囲
- CloudBeaverのDocker Compose設定
- 起動コマンドの定義（package.json）
- 環境別制御の実装（APP_ENV環境変数）
- Resource Manager用ディレクトリマウント
- データベース接続設定
- ドキュメント整備

### 1.3 設計方針
- **既存システムとの共存**: GoAdmin（カスタム処理用）とCloudBeaver（データ操作用）の役割分担を明確化
- **環境別制御**: 既存システムと同様に`APP_ENV`環境変数で環境を切り替え
- **Git管理**: Resource Managerに保存したスクリプトをGitで管理可能にする
- **シンプルな運用**: Docker Composeとnpmスクリプトによる簡単な起動・停止
- **開発環境優先**: 本実装は開発環境を優先し、ステージング・本番環境は必要に応じて対応

## 2. アーキテクチャ設計

### 2.1 ディレクトリ構造

#### 2.1.1 変更前の構造
```
.
├── server/
│   └── data/
│       ├── master.db
│       ├── sharding_db_1.db
│       ├── sharding_db_2.db
│       ├── sharding_db_3.db
│       └── sharding_db_4.db
├── config/
│   ├── develop/
│   │   └── database.yaml
│   ├── staging/
│   │   └── database.yaml
│   └── production/
│       └── database.yaml
└── client/
    └── package.json
```

#### 2.1.2 変更後の構造
```
.
├── docker-compose.yml          # 新規: CloudBeaver用Docker Compose設定
├── package.json                # 新規: npmスクリプト定義（プロジェクトルート）
├── cloudbeaver/                # 新規: CloudBeaver関連ディレクトリ
│   ├── config/                 # CloudBeaver設定ファイル（環境別）
│   │   ├── develop/            # 開発環境用設定
│   │   ├── staging/            # ステージング環境用設定
│   │   └── production/         # 本番環境用設定
│   └── scripts/                # Resource Manager用スクリプト保存ディレクトリ
├── server/
│   └── data/                   # 既存（維持）
│       ├── master.db
│       ├── sharding_db_1.db
│       ├── sharding_db_2.db
│       ├── sharding_db_3.db
│       └── sharding_db_4.db
├── config/
│   ├── develop/
│   │   └── database.yaml       # 既存（維持、参考として使用）
│   ├── staging/
│   │   └── database.yaml       # 既存（維持、参考として使用）
│   └── production/
│       └── database.yaml       # 既存（維持、参考として使用）
└── client/
    └── package.json            # 既存（維持）
```

### 2.2 ファイル構成

#### 2.2.1 Docker Compose設定ファイル
- **`docker-compose.yml`**: CloudBeaver用Docker Compose設定
  - CloudBeaverサービスの定義
  - ポートマッピング（8978:8978）
  - ボリュームマウント（データベースファイル、Resource Manager用ディレクトリ）
  - 環境変数の設定

#### 2.2.2 npmスクリプト定義ファイル
- **`package.json`**: プロジェクトルート用のnpmスクリプト定義
  - `cloudbeaver:start`: CloudBeaver起動スクリプト
  - `cloudbeaver:stop`: CloudBeaver停止スクリプト
  - `cloudbeaver:logs`: CloudBeaverログ確認スクリプト（オプション）

#### 2.2.3 CloudBeaver設定ディレクトリ
- **`cloudbeaver/config/{env}/`**: CloudBeaverの設定ファイルを保存するディレクトリ（環境別）
  - `cloudbeaver/config/develop/`: 開発環境用設定
  - `cloudbeaver/config/staging/`: ステージング環境用設定
  - `cloudbeaver/config/production/`: 本番環境用設定
  - Gitで管理
  - 接続情報、Resource Manager設定などが保存される

#### 2.2.4 Resource Manager用ディレクトリ
- **`cloudbeaver/scripts/`**: Resource Managerに保存したスクリプトを保存するディレクトリ
  - Gitで管理
  - スクリプトファイル（.sqlなど）を保存

### 2.3 システム構成図

```
┌─────────────────────────────────────────────────────────┐
│                    開発者                                │
└──────────────────┬────────────────────────────────────┘
                    │
                    │ npm run cloudbeaver:start
                    │ (APP_ENV=develop)
                    │
                    ▼
┌─────────────────────────────────────────────────────────┐
│              package.json (プロジェクトルート)            │
│  - cloudbeaver:start                                     │
│  - cloudbeaver:stop                                      │
│  - cloudbeaver:logs                                      │
└──────────────────┬────────────────────────────────────┘
                    │
                    │ docker-compose up -d
                    │
                    ▼
┌─────────────────────────────────────────────────────────┐
│              docker-compose.yml                          │
│  - CloudBeaverサービス定義                                │
│  - ポートマッピング: 8978:8978                          │
│  - ボリュームマウント:                                  │
│    - server/data/ → /data                                │
│    - cloudbeaver/scripts/ → /scripts                    │
└──────────────────┬────────────────────────────────────┘
                    │
                    │ Docker Compose起動
                    │
                    ▼
┌─────────────────────────────────────────────────────────┐
│         CloudBeaverコンテナ (dbeaver/cloudbeaver)        │
│                                                          │
│  ┌──────────────────────────────────────────────────┐ │
│  │  CloudBeaver Web UI (ポート8978)                  │ │
│  │  - データベース接続管理                            │ │
│  │  - SQL実行                                         │ │
│  │  - Resource Manager                                │ │
│  └──────────────────────────────────────────────────┘ │
│                                                          │
│  ┌──────────────────────────────────────────────────┐ │
│  │  マウントされたボリューム                          │ │
│  │  - /data: server/data/ (データベースファイル)      │ │
│  │  - /scripts: cloudbeaver/scripts/ (スクリプト)    │ │
│  └──────────────────────────────────────────────────┘ │
└──────────────────┬────────────────────────────────────┘
                    │
                    │ SQLite接続
                    │
                    ▼
┌─────────────────────────────────────────────────────────┐
│              server/data/                                │
│  - master.db                                             │
│  - sharding_db_1.db                                      │
│  - sharding_db_2.db                                      │
│  - sharding_db_3.db                                      │
│  - sharding_db_4.db                                      │
└─────────────────────────────────────────────────────────┘
```

### 2.4 データフロー

#### 2.4.1 CloudBeaver起動フロー
```
開発者が npm run cloudbeaver:start を実行
    ↓
package.jsonのスクリプトが実行される
    ↓
APP_ENV環境変数を取得（デフォルト: develop）
    ↓
docker-compose.ymlに環境変数を渡す
    ↓
Docker ComposeがCloudBeaverコンテナを起動
    ↓
CloudBeaverがポート8978で起動
    ↓
開発者が http://localhost:8978 にアクセス
```

#### 2.4.2 データベース接続フロー
```
CloudBeaver Web UIにアクセス
    ↓
手動でデータベース接続を設定
    ↓
SQLiteドライバーを選択
    ↓
データベースファイルのパスを指定（/data/master.dbなど）
    ↓
接続設定を保存
    ↓
データベースに接続
    ↓
テーブル一覧、データ閲覧、SQL実行が可能
```

#### 2.4.3 Resource Managerスクリプト保存フロー
```
CloudBeaver Web UIでSQLスクリプトを作成
    ↓
Resource Managerに保存
    ↓
マウントされた /scripts ディレクトリに保存
    ↓
cloudbeaver/scripts/ ディレクトリに反映
    ↓
Gitで管理可能
```

## 3. コンポーネント設計

### 3.1 Docker Compose設定

#### 3.1.1 docker-compose.ymlの構造
```yaml
version: '3.8'

services:
  cloudbeaver:
    image: dbeaver/cloudbeaver:latest
    container_name: cloudbeaver
    restart: unless-stopped
    ports:
      - "${CLOUDBEAVER_PORT:-8978}:8978"
    volumes:
      - ./server/data:/data:ro
      - ./cloudbeaver/scripts:/scripts
      - ./cloudbeaver/config/${APP_ENV:-develop}:/opt/cloudbeaver/workspace
    environment:
      - APP_ENV=${APP_ENV:-develop}
      - CB_WORKSPACE=/opt/cloudbeaver/workspace
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8978"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
```

#### 3.1.2 設定項目の説明
- **image**: CloudBeaverの公式Dockerイメージ
- **container_name**: コンテナ名（`cloudbeaver`）
- **restart**: コンテナの再起動ポリシー（`unless-stopped`）
- **ports**: ポートマッピング（デフォルト: 8978:8978、環境変数で変更可能）
- **volumes**:
  - `./server/data:/data:ro`: データベースファイルを読み取り専用でマウント
  - `./cloudbeaver/scripts:/scripts`: Resource Manager用スクリプトディレクトリをマウント
  - `./cloudbeaver/config/${APP_ENV:-develop}:/opt/cloudbeaver/workspace`: CloudBeaver設定ディレクトリを環境別にマウント
- **environment**: 環境変数の設定
  - `APP_ENV`: 環境名（develop/staging/production）
  - `CB_WORKSPACE`: CloudBeaverのワークスペースディレクトリ（`/opt/cloudbeaver/workspace`）
- **healthcheck**: コンテナのヘルスチェック設定

#### 3.1.3 環境別設定の考慮
- **設定ディレクトリ**: `APP_ENV`環境変数に基づいて、適切な設定ディレクトリをマウント
  - `APP_ENV=develop` → `./cloudbeaver/config/develop/`
  - `APP_ENV=staging` → `./cloudbeaver/config/staging/`
  - `APP_ENV=production` → `./cloudbeaver/config/production/`
- **データベースファイルのパス**: 環境によって異なる場合は、環境変数で制御
- **ポート番号**: 環境によって異なる場合は、`CLOUDBEAVER_PORT`環境変数で制御
- **既存設定ファイルの参照**: `config/{env}/database.yaml`の設定を参考に、環境別のデータベースパスを決定

### 3.2 package.jsonの設計

#### 3.2.1 package.jsonの構造
```json
{
  "name": "go-webdb-template",
  "version": "1.0.0",
  "description": "Go Web DB Template - CloudBeaver management scripts",
  "scripts": {
    "cloudbeaver:start": "./scripts/cloudbeaver-start.sh",
    "cloudbeaver:stop": "docker-compose down",
    "cloudbeaver:logs": "docker-compose logs -f cloudbeaver",
    "cloudbeaver:restart": "npm run cloudbeaver:stop && npm run cloudbeaver:start"
  }
}
```

#### 3.2.2 npmスクリプトの説明
- **cloudbeaver:start**: CloudBeaverを起動
  - `scripts/cloudbeaver-start.sh`スクリプトを実行
  - スクリプト内で`APP_ENV`環境変数を設定（未設定の場合は`develop`をデフォルト）
  - Docker ComposeでCloudBeaverコンテナを起動
- **cloudbeaver:stop**: CloudBeaverを停止
  - Docker ComposeでCloudBeaverコンテナを停止
- **cloudbeaver:logs**: CloudBeaverのログを確認
  - Docker Composeのログを表示
- **cloudbeaver:restart**: CloudBeaverを再起動
  - 停止してから起動

#### 3.2.3 環境変数の扱い
- `scripts/cloudbeaver-start.sh`スクリプト内で環境変数を設定
- `APP_ENV`環境変数が未設定の場合は`develop`をデフォルトとする
- 環境変数はDocker Composeに渡され、コンテナ内で使用可能

#### 3.2.4 cloudbeaver-start.shスクリプトの構造
```bash
#!/bin/bash
# CloudBeaver起動スクリプト

# APP_ENV環境変数が未設定の場合はdevelopをデフォルトとする
export APP_ENV=${APP_ENV:-develop}

# Docker ComposeでCloudBeaverを起動
docker-compose up -d
```

- スクリプトの実行権限を付与: `chmod +x scripts/cloudbeaver-start.sh`
- シェルスクリプトを使用することで、staging・本番環境でも動作する

### 3.3 Resource Manager設定

#### 3.3.1 ディレクトリ構造
```
cloudbeaver/
├── config/                    # CloudBeaver設定ファイル（環境別）
│   ├── develop/               # 開発環境用設定
│   │   └── .gitkeep          # Git管理用
│   ├── staging/               # ステージング環境用設定
│   │   └── .gitkeep          # Git管理用
│   └── production/            # 本番環境用設定
│       └── .gitkeep          # Git管理用
└── scripts/                   # Resource Manager用スクリプト
    ├── .gitkeep              # Git管理用（空ディレクトリをGitに含める）
    └── (スクリプトファイル)   # Resource Managerに保存されたスクリプト
```

#### 3.3.2 CloudBeaverの設定
- CloudBeaver起動後、Web UIからResource Managerの設定を行う
- Resource Managerの保存先を`/scripts`ディレクトリに設定
- スクリプトを保存すると、`cloudbeaver/scripts/`ディレクトリに反映される
- スクリプトファイルはGitで管理可能

#### 3.3.3 Git管理
- `cloudbeaver/scripts/`ディレクトリをGitで管理
- `.gitignore`で不要なファイルを除外（必要に応じて）
- スクリプトファイルの命名規則は特に制限しない（CloudBeaverのデフォルトに従う）

### 3.4 データベース接続設定

#### 3.4.1 接続設定の手順
1. CloudBeaver起動後、http://localhost:8978 にアクセス
2. 初回起動時は初期設定画面が表示される
3. データベース接続を追加
4. SQLiteドライバーを選択
5. 接続情報を入力:
   - **接続名**: `master` または `Master Database`
   - **データベースファイル**: `/data/master.db`
6. 接続をテスト
7. 接続を保存

#### 3.4.2 接続設定の一覧
- **マスターデータベース**:
  - 接続名: `master`
  - データベースファイル: `/data/master.db`
- **シャーディングデータベース**:
  - 接続名: `sharding_db_1`
  - データベースファイル: `/data/sharding_db_1.db`
  - 接続名: `sharding_db_2`
  - データベースファイル: `/data/sharding_db_2.db`
  - 接続名: `sharding_db_3`
  - データベースファイル: `/data/sharding_db_3.db`
  - 接続名: `sharding_db_4`
  - データベースファイル: `/data/sharding_db_4.db`

#### 3.4.3 接続設定の保存
- CloudBeaverの設定ファイルに接続情報が保存される
- 設定ファイルは`cloudbeaver/config/{env}/`ディレクトリに保存される（環境別）
- 設定ファイルはGitで管理可能
- コンテナを再起動しても接続設定は保持される（マウントされたディレクトリに保存）
- 環境別に設定が分離されるため、環境ごとに異なる接続設定を管理可能

## 4. 環境別制御の実装

### 4.1 環境変数の取得

#### 4.1.1 シェルスクリプトでの環境変数取得
```bash
# scripts/cloudbeaver-start.sh
export APP_ENV=${APP_ENV:-develop}
docker-compose up -d
```

- シェルスクリプト内で環境変数を設定
- `${APP_ENV:-develop}`で、`APP_ENV`が未設定の場合は`develop`をデフォルトとする
- 環境変数はDocker Composeに渡される

#### 4.1.2 Docker Composeでの環境変数参照
```yaml
environment:
  - APP_ENV=${APP_ENV:-develop}
```

- Docker Composeの`environment`セクションで環境変数を設定
- コンテナ内で`APP_ENV`環境変数が使用可能

### 4.2 環境別設定の適用

#### 4.2.1 設定ディレクトリの選択
- `APP_ENV`環境変数に基づいて、適切な設定ディレクトリをマウント
- 開発環境: `cloudbeaver/config/develop/`
- ステージング環境: `cloudbeaver/config/staging/`
- 本番環境: `cloudbeaver/config/production/`
- Docker Composeの`volumes`セクションで`${APP_ENV:-develop}`を使用して環境別ディレクトリを参照

#### 4.2.2 データベースファイルのパス
- 開発環境: `server/data/`ディレクトリ（既存のまま）
- ステージング・本番環境: 環境によって異なる可能性があるが、本実装では開発環境を優先
- 環境別のパスが必要な場合は、環境変数やDocker Composeの設定で制御

#### 4.2.3 ポート番号
- デフォルト: 8978
- 環境によって異なる場合は、`CLOUDBEAVER_PORT`環境変数で制御
- `docker-compose.yml`で`${CLOUDBEAVER_PORT:-8978}`として参照

### 4.3 環境別設定ファイルの参照

#### 4.3.1 既存設定ファイルの活用
- `config/{env}/database.yaml`の設定を参考に、環境別のデータベースパスを決定
- ただし、CloudBeaverはDockerコンテナ内で動作するため、マウントパスは固定
- 環境によってデータベースファイルの場所が異なる場合は、環境変数で制御

## 5. エラーハンドリング

### 5.1 Docker Compose起動時のエラー

#### 5.1.1 ポート競合
- **問題**: ポート8978が既に使用されている
- **対処**: `CLOUDBEAVER_PORT`環境変数でポート番号を変更
- **例**: `CLOUDBEAVER_PORT=8979 npm run cloudbeaver:start`

#### 5.1.2 ボリュームマウントエラー
- **問題**: マウントするディレクトリが存在しない
- **対処**: 必要なディレクトリを作成
  - `mkdir -p server/data`
  - `mkdir -p cloudbeaver/scripts`

#### 5.1.3 権限エラー
- **問題**: データベースファイルへのアクセス権限がない
- **対処**: ファイルの読み書き権限を確認
  - データベースファイルは読み取り専用でマウント（`:ro`オプション）

### 5.2 CloudBeaver起動時のエラー

#### 5.2.1 コンテナ起動失敗
- **問題**: CloudBeaverコンテナが起動しない
- **対処**: ログを確認（`npm run cloudbeaver:logs`）
- **原因**: イメージのダウンロード失敗、リソース不足など

#### 5.2.2 データベース接続エラー
- **問題**: データベースに接続できない
- **対処**: 
  - データベースファイルのパスを確認
  - マウント設定を確認
  - ファイルの存在を確認

### 5.3 Resource Managerのエラー

#### 5.3.1 スクリプト保存エラー
- **問題**: Resource Managerにスクリプトを保存できない
- **対処**: 
  - `cloudbeaver/scripts/`ディレクトリの権限を確認
  - マウント設定を確認

## 6. テスト戦略

### 6.1 単体テスト

#### 6.1.1 Docker Compose設定のテスト
- `docker-compose.yml`の構文チェック
- 環境変数の参照が正しいか確認

#### 6.1.2 npmスクリプトのテスト
- 各スクリプトが正しく実行されるか確認
- 環境変数のデフォルト値が正しいか確認

### 6.2 統合テスト

#### 6.2.1 CloudBeaver起動テスト
- `npm run cloudbeaver:start`でCloudBeaverが起動するか確認
- http://localhost:8978 にアクセスできるか確認

#### 6.2.2 データベース接続テスト
- マスターデータベースに接続できるか確認
- シャーディングデータベースに接続できるか確認
- テーブル一覧が表示されるか確認
- データを閲覧できるか確認
- SQLクエリを実行できるか確認

#### 6.2.3 Resource Managerテスト
- Resource Managerにスクリプトを保存できるか確認
- 保存したスクリプトが`cloudbeaver/scripts/`ディレクトリに反映されるか確認
- スクリプトがGitで管理できるか確認

### 6.3 環境別テスト

#### 6.3.1 開発環境テスト
- `APP_ENV=develop npm run cloudbeaver:start`で起動できるか確認
- デフォルト（`APP_ENV`未設定）で起動できるか確認

#### 6.3.2 ステージング環境テスト
- `APP_ENV=staging npm run cloudbeaver:start`で起動できるか確認
- 環境別の設定が正しく適用されるか確認

#### 6.3.3 本番環境テスト
- `APP_ENV=production npm run cloudbeaver:start`で起動できるか確認
- 環境別の設定が正しく適用されるか確認

## 7. セキュリティ考慮事項

### 7.1 データベースファイルへのアクセス

#### 7.1.1 読み取り専用マウント
- データベースファイルは読み取り専用でマウント（`:ro`オプション）
- CloudBeaverからデータベースファイルを誤って変更することを防止

#### 7.1.2 アクセス制御
- CloudBeaverは開発環境での使用を想定
- 本番環境での使用は想定しない（本番環境では適切なアクセス制御が必要）

### 7.2 認証設定

#### 7.2.1 CloudBeaverの認証
- CloudBeaverのデフォルト認証設定を確認
- 必要に応じて認証を有効化

#### 7.2.2 ネットワークアクセス
- CloudBeaverはローカルホスト（localhost）でのみアクセス可能
- 外部からのアクセスは想定しない

## 8. 運用・保守

### 8.1 起動・停止

#### 8.1.1 起動
```bash
# 開発環境（デフォルト）
npm run cloudbeaver:start

# 環境を指定
APP_ENV=develop npm run cloudbeaver:start
APP_ENV=staging npm run cloudbeaver:start
APP_ENV=production npm run cloudbeaver:start
```

#### 8.1.2 停止
```bash
npm run cloudbeaver:stop
```

#### 8.1.3 再起動
```bash
npm run cloudbeaver:restart
```

#### 8.1.4 ログ確認
```bash
npm run cloudbeaver:logs
```

### 8.2 データベース接続管理

#### 8.2.1 接続設定の追加
- CloudBeaver Web UIから手動で接続設定を追加
- 接続設定はCloudBeaverコンテナ内に保存される

#### 8.2.2 接続設定の削除
- CloudBeaver Web UIから手動で接続設定を削除

### 8.3 Resource Manager管理

#### 8.3.1 スクリプトの作成
- CloudBeaver Web UIでSQLスクリプトを作成
- Resource Managerに保存
- `cloudbeaver/scripts/`ディレクトリに反映される

#### 8.3.2 スクリプトのGit管理
- `cloudbeaver/scripts/`ディレクトリをGitで管理
- スクリプトの変更をGitで追跡

### 8.4 トラブルシューティング

#### 8.4.1 コンテナが起動しない
- Dockerが起動しているか確認
- ポート8978が使用されていないか確認
- ログを確認（`npm run cloudbeaver:logs`）

#### 8.4.2 データベースに接続できない
- データベースファイルが存在するか確認
- マウント設定を確認
- ファイルのパスを確認

#### 8.4.3 Resource Managerにスクリプトを保存できない
- `cloudbeaver/scripts/`ディレクトリの権限を確認
- マウント設定を確認

## 9. 実装上の注意事項

### 9.1 Docker Compose設定

#### 9.1.1 イメージのバージョン
- `dbeaver/cloudbeaver:latest`を使用
- 必要に応じて特定のバージョンを指定

#### 9.1.2 ボリュームマウント
- データベースファイルは読み取り専用でマウント（`:ro`オプション）
- Resource Manager用ディレクトリは読み書き可能でマウント

#### 9.1.3 環境変数
- `APP_ENV`環境変数をコンテナに渡す
- ポート番号は`CLOUDBEAVER_PORT`環境変数で制御可能

### 9.2 package.jsonの作成

#### 9.2.1 依存関係
- 依存関係は不要（シェルスクリプトを使用）

#### 9.2.2 スクリプトの定義
- `scripts/cloudbeaver-start.sh`スクリプトを呼び出す
- 環境変数のデフォルト値はシェルスクリプト内で設定

#### 9.2.3 cloudbeaver-start.shスクリプトの作成
- `scripts/cloudbeaver-start.sh`を作成
- 実行権限を付与: `chmod +x scripts/cloudbeaver-start.sh`
- スクリプト内で`APP_ENV`環境変数を設定（未設定の場合は`develop`をデフォルト）

### 9.3 CloudBeaver設定ディレクトリの作成

#### 9.3.1 環境別設定ディレクトリの作成
- `cloudbeaver/config/develop/`ディレクトリを作成
- `cloudbeaver/config/staging/`ディレクトリを作成
- `cloudbeaver/config/production/`ディレクトリを作成
- 各ディレクトリに`.gitkeep`ファイルを追加（空ディレクトリをGitに含める）

#### 9.3.2 設定ファイルの管理
- CloudBeaver起動後、Web UIから接続設定を行う
- 接続設定は`cloudbeaver/config/{env}/`ディレクトリに保存される
- 設定ファイルはGitで管理可能
- 環境別に設定が分離されるため、環境ごとに異なる接続設定を管理可能

### 9.4 Resource Manager設定

#### 9.4.1 ディレクトリの作成
- `cloudbeaver/scripts/`ディレクトリを作成
- `.gitkeep`ファイルを追加（空ディレクトリをGitに含める）

#### 9.4.2 CloudBeaverの設定
- CloudBeaver起動後、Web UIからResource Managerの設定を行う
- 保存先を`/scripts`ディレクトリに設定

### 9.4 ドキュメント整備

#### 9.4.1 README.mdの更新
- CloudBeaverの起動方法を追記
- 環境別の起動方法を記載
- 基本的な使用方法を記載
- データベース接続設定の手順を記載
- Resource Managerの使用方法を記載

## 10. 参考情報

### 10.1 関連ドキュメント
- 要件定義書: `.kiro/specs/0015-cloudbeaver/requirements.md`
- プロジェクトREADME: `README.md`
- シャーディング仕様: `docs/Sharding.md`

### 10.2 技術スタック
- **CloudBeaver**: https://cloudbeaver.io/
- **Docker**: Docker Compose
- **npm**: npmスクリプト
- **cross-env**: クロスプラットフォーム環境変数設定

### 10.3 参考リンク
- CloudBeaver公式サイト: https://cloudbeaver.io/
- CloudBeaver GitHub: https://github.com/dbeaver/cloudbeaver
- CloudBeaver Docker: https://hub.docker.com/r/dbeaver/cloudbeaver
- CloudBeaver ドキュメント: https://cloudbeaver.io/docs/

