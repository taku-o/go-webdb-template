# Metabase導入設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、MetabaseをDockerで動作させ、Webベースのデータ可視化・分析ツールを提供するシステムの詳細設計を定義する。既存システム（GoAdmin、CloudBeaver）と共存し、データ可視化・分析用の管理アプリとして機能する。

### 1.2 設計の範囲
- MetabaseのDocker Compose設定
- docker-compose.ymlファイルの分離（CloudBeaver用とMetabase用）
- 起動コマンドの定義（package.json）
- 環境別制御の実装（APP_ENV環境変数）
- Metabase設定ディレクトリの管理
- データベース接続設定
- ドキュメント整備

### 1.3 設計方針
- **既存システムとの共存**: GoAdmin（カスタム処理用）、CloudBeaver（データ操作用）、Metabase（データ可視化・分析用）の役割分担を明確化
- **docker-compose.ymlの分離**: CloudBeaver用とMetabase用でdocker-compose.ymlファイルを分けることで、個別に起動・停止可能にする
- **環境別制御**: 既存システムと同様に`APP_ENV`環境変数で環境を切り替え
- **Git管理**: Metabaseの設定ファイルをGitで管理可能にする
- **シンプルな運用**: Docker Composeとnpmスクリプトによる簡単な起動・停止
- **開発環境優先**: 本実装は開発環境を優先し、ステージング・本番環境は必要に応じて対応

## 2. アーキテクチャ設計

### 2.1 ディレクトリ構造

#### 2.1.1 変更前の構造
```
.
├── docker-compose.yml          # CloudBeaver用（既存）
├── package.json                # npmスクリプト定義（既存）
├── cloudbeaver/                # CloudBeaver関連ディレクトリ（既存）
│   ├── config/                 # CloudBeaver設定ファイル（環境別）
│   └── scripts/                # Resource Manager用スクリプト
├── server/
│   └── data/                   # データベースファイル
│       ├── master.db
│       ├── sharding_db_1.db
│       ├── sharding_db_2.db
│       ├── sharding_db_3.db
│       └── sharding_db_4.db
└── config/
    ├── develop/
    ├── staging/
    └── production/
```

#### 2.1.2 変更後の構造
```
.
├── docker-compose.cloudbeaver.yml  # CloudBeaver用（既存をリネーム）
├── docker-compose.metabase.yml     # Metabase用（新規）
├── package.json                    # npmスクリプト定義（更新）
├── scripts/
│   ├── cloudbeaver-start.sh        # CloudBeaver起動スクリプト（既存、更新）
│   └── metabase-start.sh           # Metabase起動スクリプト（新規）
├── cloudbeaver/                    # CloudBeaver関連ディレクトリ（既存）
│   ├── config/                     # CloudBeaver設定ファイル（環境別）
│   └── scripts/                    # Resource Manager用スクリプト
├── metabase/                       # Metabase関連ディレクトリ（新規）
│   └── config/                     # Metabase設定ファイル（環境別）
│       ├── develop/                # 開発環境用設定
│       ├── staging/                # ステージング環境用設定
│       └── production/             # 本番環境用設定
├── server/
│   └── data/                       # 既存（維持）
│       ├── master.db
│       ├── sharding_db_1.db
│       ├── sharding_db_2.db
│       ├── sharding_db_3.db
│       └── sharding_db_4.db
└── config/
    ├── develop/
    ├── staging/
    └── production/
```

### 2.2 ファイル構成

#### 2.2.1 Docker Compose設定ファイル
- **`docker-compose.cloudbeaver.yml`**: CloudBeaver用Docker Compose設定（既存の`docker-compose.yml`をリネーム）
  - CloudBeaverサービスの定義
  - ポートマッピング（8978:8978）
  - ボリュームマウント（データベースファイル、設定ディレクトリ）
  - 環境変数の設定
- **`docker-compose.metabase.yml`**: Metabase用Docker Compose設定（新規作成）
  - Metabaseサービスの定義
  - ポートマッピング（8970:3000）
  - ボリュームマウント（データベースファイル、設定ディレクトリ）
  - 環境変数の設定

#### 2.2.2 npmスクリプト定義ファイル
- **`package.json`**: プロジェクトルート用のnpmスクリプト定義（既存を更新）
  - `cloudbeaver:start`: CloudBeaver起動スクリプト（更新）
  - `cloudbeaver:stop`: CloudBeaver停止スクリプト（更新）
  - `metabase:start`: Metabase起動スクリプト（新規）
  - `metabase:stop`: Metabase停止スクリプト（新規）
  - `metabase:logs`: Metabaseログ確認スクリプト（新規）

#### 2.2.3 Metabase設定ディレクトリ
- **`metabase/config/{env}/`**: Metabaseの設定ファイルを保存するディレクトリ（環境別）
  - `metabase/config/develop/`: 開発環境用設定
  - `metabase/config/staging/`: ステージング環境用設定
  - `metabase/config/production/`: 本番環境用設定
  - Gitで管理
  - 接続情報、ダッシュボード設定などが保存される

#### 2.2.4 起動スクリプト
- **`scripts/cloudbeaver-start.sh`**: CloudBeaver起動スクリプト（既存、更新）
  - `docker-compose.cloudbeaver.yml`を使用するように更新
- **`scripts/metabase-start.sh`**: Metabase起動スクリプト（新規作成）
  - `docker-compose.metabase.yml`を使用

### 2.3 システム構成図

```
┌─────────────────────────────────────────────────────────┐
│                    開発者                                │
└──────────────────┬────────────────────────────────────┘
                    │
                    │ npm run metabase:start
                    │ (APP_ENV=develop)
                    │
                    ▼
┌─────────────────────────────────────────────────────────┐
│              package.json (プロジェクトルート)            │
│  - metabase:start                                        │
│  - metabase:stop                                         │
│  - metabase:logs                                         │
└──────────────────┬────────────────────────────────────┘
                    │
                    │ docker-compose -f docker-compose.metabase.yml up -d
                    │
                    ▼
┌─────────────────────────────────────────────────────────┐
│         docker-compose.metabase.yml                      │
│  - Metabaseサービス定義                                  │
│  - ポートマッピング: 8970:3000                          │
│  - ボリュームマウント:                                  │
│    - server/data/ → /data                                │
│    - metabase/config/{env}/ → /metabase-data            │
└──────────────────┬────────────────────────────────────┘
                    │
                    │ Docker Compose起動
                    │
                    ▼
┌─────────────────────────────────────────────────────────┐
│         Metabaseコンテナ (metabase/metabase)            │
│                                                          │
│  ┌──────────────────────────────────────────────────┐ │
│  │  Metabase Web UI (ポート3000→ホスト8970)         │ │
│  │  - データベース接続管理                            │ │
│  │  - クエリ作成・実行                                │ │
│  │  - ダッシュボード作成                              │ │
│  └──────────────────────────────────────────────────┘ │
│                                                          │
│  ┌──────────────────────────────────────────────────┐ │
│  │  マウントされたボリューム                          │ │
│  │  - /data: server/data/ (データベースファイル)      │ │
│  │  - /metabase-data: metabase/config/{env}/ (設定)  │ │
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

#### 2.4.1 Metabase起動フロー
```
開発者が npm run metabase:start を実行
    ↓
package.jsonのスクリプトが実行される
    ↓
scripts/metabase-start.shが実行される
    ↓
APP_ENV環境変数を取得（デフォルト: develop）
    ↓
docker-compose -f docker-compose.metabase.yml up -d を実行
    ↓
Docker ComposeがMetabaseコンテナを起動
    ↓
Metabaseがポート3000（ホスト側は8970）で起動
    ↓
開発者が http://localhost:8970 にアクセス
```

#### 2.4.2 データベース接続フロー
```
Metabase Web UIにアクセス
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
テーブル一覧、データ閲覧、クエリ作成、ダッシュボード作成が可能
```

#### 2.4.3 Metabase設定保存フロー
```
Metabase Web UIで接続設定やダッシュボードを作成
    ↓
設定がMetabaseコンテナ内に保存される
    ↓
マウントされた /metabase-data ディレクトリに保存
    ↓
metabase/config/{env}/ ディレクトリに反映
    ↓
Gitで管理可能
```

## 3. コンポーネント設計

### 3.1 Docker Compose設定

#### 3.1.1 docker-compose.metabase.ymlの構造
```yaml
version: '3.8'

services:
  metabase:
    image: metabase/metabase:latest
    container_name: metabase
    restart: unless-stopped
    ports:
      - "8970:3000"
    volumes:
      - ./server/data:/data:ro
      - ./metabase/config/${APP_ENV:-develop}:/metabase-data
    environment:
      - APP_ENV=${APP_ENV:-develop}
      - MB_DB_FILE=/metabase-data/metabase.db
    # MB_DB_FILEは固定値として記載（環境変数ではなくdocker-compose.ymlに直接記載）
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:3000/api/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 180s
```

#### 3.1.2 設定項目の説明
- **image**: Metabaseの公式Dockerイメージ（`metabase/metabase:latest`）
- **container_name**: コンテナ名（`metabase`）
- **restart**: コンテナの再起動ポリシー（`unless-stopped`）
- **ports**: ポートマッピング（固定: 8970:3000）
  - ホスト側: 8970
  - コンテナ内: 3000（Metabaseのデフォルトポート）
- **volumes**:
  - `./server/data:/data:ro`: データベースファイルを読み取り専用でマウント
  - `./metabase/config/${APP_ENV:-develop}:/metabase-data`: Metabase設定ディレクトリを環境別にマウント
- **environment**: 環境変数の設定
  - `APP_ENV`: 環境名（develop/staging/production）
  - `MB_DB_FILE`: Metabaseの内部データベースファイルのパス（`/metabase-data/metabase.db`、固定値としてdocker-compose.ymlに直接記載）
- **healthcheck**: コンテナのヘルスチェック設定
  - MetabaseのヘルスチェックAPIエンドポイントを使用

#### 3.1.3 docker-compose.cloudbeaver.ymlの更新
既存の`docker-compose.yml`を`docker-compose.cloudbeaver.yml`にリネームし、起動スクリプトで明示的に指定するように更新する。

#### 3.1.4 環境別設定の考慮
- **設定ディレクトリ**: `APP_ENV`環境変数に基づいて、適切な設定ディレクトリをマウント
  - `APP_ENV=develop` → `./metabase/config/develop/`
  - `APP_ENV=staging` → `./metabase/config/staging/`
  - `APP_ENV=production` → `./metabase/config/production/`
- **データベースファイルのパス**: 環境によって異なる場合は、環境変数で制御
- **ポート番号**: 固定値（8970）を使用。ポート競合の回避は運用者の責任
- **既存設定ファイルの参照**: `config/{env}/database.yaml`の設定を参考に、環境別のデータベースパスを決定

### 3.2 package.jsonの設計

#### 3.2.1 package.jsonの構造（更新後）
```json
{
  "name": "go-webdb-template",
  "version": "1.0.0",
  "description": "Go Web DB Template - CloudBeaver and Metabase management scripts",
  "scripts": {
    "cloudbeaver:start": "./scripts/cloudbeaver-start.sh",
    "cloudbeaver:stop": "docker-compose -f docker-compose.cloudbeaver.yml down",
    "cloudbeaver:logs": "docker-compose -f docker-compose.cloudbeaver.yml logs -f cloudbeaver",
    "cloudbeaver:restart": "npm run cloudbeaver:stop && npm run cloudbeaver:start",
    "metabase:start": "./scripts/metabase-start.sh",
    "metabase:stop": "docker-compose -f docker-compose.metabase.yml down",
    "metabase:logs": "docker-compose -f docker-compose.metabase.yml logs -f metabase",
    "metabase:restart": "npm run metabase:stop && npm run metabase:start"
  }
}
```

#### 3.2.2 npmスクリプトの説明
- **cloudbeaver:start**: CloudBeaverを起動（更新）
  - `scripts/cloudbeaver-start.sh`スクリプトを実行
  - `docker-compose.cloudbeaver.yml`を使用するように更新
- **cloudbeaver:stop**: CloudBeaverを停止（更新）
  - `docker-compose.cloudbeaver.yml`を使用するように更新
- **cloudbeaver:logs**: CloudBeaverのログを確認（更新）
  - `docker-compose.cloudbeaver.yml`を使用するように更新
- **metabase:start**: Metabaseを起動（新規）
  - `scripts/metabase-start.sh`スクリプトを実行
  - スクリプト内で`APP_ENV`環境変数を設定（未設定の場合は`develop`をデフォルト）
  - Docker ComposeでMetabaseコンテナを起動
- **metabase:stop**: Metabaseを停止（新規）
  - Docker ComposeでMetabaseコンテナを停止
- **metabase:logs**: Metabaseのログを確認（新規）
  - Docker Composeのログを表示
- **metabase:restart**: Metabaseを再起動（新規）
  - 停止してから起動

#### 3.2.3 環境変数の扱い
- `scripts/metabase-start.sh`スクリプト内で環境変数を設定
- `APP_ENV`環境変数が未設定の場合は`develop`をデフォルトとする
- 環境変数はDocker Composeに渡され、コンテナ内で使用可能

#### 3.2.4 metabase-start.shスクリプトの構造
```bash
#!/bin/bash
# Metabase起動スクリプト

# APP_ENV環境変数が未設定の場合はdevelopをデフォルトとする
export APP_ENV=${APP_ENV:-develop}

# Docker ComposeでMetabaseを起動
docker-compose -f docker-compose.metabase.yml up -d
```

- スクリプトの実行権限を付与: `chmod +x scripts/metabase-start.sh`
- シェルスクリプトを使用することで、staging・本番環境でも動作する

#### 3.2.5 cloudbeaver-start.shスクリプトの更新
```bash
#!/bin/bash
# CloudBeaver起動スクリプト

# APP_ENV環境変数が未設定の場合はdevelopをデフォルトとする
export APP_ENV=${APP_ENV:-develop}

# Docker ComposeでCloudBeaverを起動
docker-compose -f docker-compose.cloudbeaver.yml up -d
```

- `docker-compose.yml`を`docker-compose.cloudbeaver.yml`に変更

### 3.3 Metabase設定ディレクトリ

#### 3.3.1 ディレクトリ構造
```
metabase/
└── config/                    # Metabase設定ファイル（環境別）
    ├── develop/               # 開発環境用設定
    │   └── .gitkeep          # Git管理用
    ├── staging/               # ステージング環境用設定
    │   └── .gitkeep          # Git管理用
    └── production/            # 本番環境用設定
        └── .gitkeep          # Git管理用
```

#### 3.3.2 Metabaseの設定
- Metabase起動後、Web UIから接続設定やダッシュボードを作成
- 設定はMetabaseコンテナ内の`/metabase-data`ディレクトリに保存される
- マウントされた`metabase/config/{env}/`ディレクトリに反映される
- 設定ファイルはGitで管理可能

#### 3.3.3 Git管理
- `metabase/config/`ディレクトリをGitで管理
- `.gitignore`で不要なファイルを除外（必要に応じて）
- 設定ファイルの命名規則は特に制限しない（Metabaseのデフォルトに従う）

### 3.4 データベース接続設定

#### 3.4.1 接続設定の手順
1. Metabase起動後、http://localhost:8970 にアクセス
2. 初回起動時は初期設定画面が表示される（管理者アカウント作成）
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
- Metabaseの設定ファイルに接続情報が保存される
- 設定ファイルは`metabase/config/{env}/`ディレクトリに保存される（環境別）
- 設定ファイルはGitで管理可能
- コンテナを再起動しても接続設定は保持される（マウントされたディレクトリに保存）
- 環境別に設定が分離されるため、環境ごとに異なる接続設定を管理可能

## 4. 環境別制御の実装

### 4.1 環境変数の取得

#### 4.1.1 シェルスクリプトでの環境変数取得
```bash
# scripts/metabase-start.sh
export APP_ENV=${APP_ENV:-develop}
docker-compose -f docker-compose.metabase.yml up -d
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
- 開発環境: `metabase/config/develop/`
- ステージング環境: `metabase/config/staging/`
- 本番環境: `metabase/config/production/`
- Docker Composeの`volumes`セクションで`${APP_ENV:-develop}`を使用して環境別ディレクトリを参照

#### 4.2.2 データベースファイルのパス
- 開発環境: `server/data/`ディレクトリ（既存のまま）
- ステージング・本番環境: 環境によって異なる可能性があるが、本実装では開発環境を優先
- 環境別のパスが必要な場合は、環境変数やDocker Composeの設定で制御

#### 4.2.3 ポート番号
- 固定値: 8970（ホスト側）、3000（コンテナ内）
- ポート番号は`docker-compose.metabase.yml`に固定値として定義
- ポート競合の回避は運用者の責任

### 4.3 環境別設定ファイルの参照

#### 4.3.1 既存設定ファイルの活用
- `config/{env}/database.yaml`の設定を参考に、環境別のデータベースパスを決定
- ただし、MetabaseはDockerコンテナ内で動作するため、マウントパスは固定
- 環境によってデータベースファイルの場所が異なる場合は、環境変数で制御

## 5. エラーハンドリング

### 5.1 Docker Compose起動時のエラー

#### 5.1.1 ポート競合
- **問題**: ポート8970が既に使用されている
- **対処**: ポート競合の回避は運用者の責任
  - 既存のサービスを停止する
  - または`docker-compose.metabase.yml`のポート番号を直接変更する

#### 5.1.2 ボリュームマウントエラー
- **問題**: マウントするディレクトリが存在しない
- **対処**: 必要なディレクトリを作成
  - `mkdir -p server/data`
  - `mkdir -p metabase/config/develop`
  - `mkdir -p metabase/config/staging`
  - `mkdir -p metabase/config/production`

#### 5.1.3 権限エラー
- **問題**: データベースファイルへのアクセス権限がない
- **対処**: ファイルの読み書き権限を確認
  - データベースファイルは読み取り専用でマウント（`:ro`オプション）

### 5.2 Metabase起動時のエラー

#### 5.2.1 コンテナ起動失敗
- **問題**: Metabaseコンテナが起動しない
- **対処**: ログを確認（`npm run metabase:logs`）
- **原因**: イメージのダウンロード失敗、リソース不足など

#### 5.2.2 データベース接続エラー
- **問題**: データベースに接続できない
- **対処**: 
  - データベースファイルのパスを確認
  - マウント設定を確認
  - ファイルの存在を確認

#### 5.2.3 設定ファイルの保存エラー
- **問題**: Metabaseの設定ファイルを保存できない
- **対処**: 
  - `metabase/config/{env}/`ディレクトリの権限を確認
  - マウント設定を確認

## 6. テスト戦略

### 6.1 単体テスト

#### 6.1.1 Docker Compose設定のテスト
- `docker-compose.metabase.yml`の構文チェック
- 環境変数の参照が正しいか確認

#### 6.1.2 npmスクリプトのテスト
- 各スクリプトが正しく実行されるか確認
- 環境変数のデフォルト値が正しいか確認

### 6.2 統合テスト

#### 6.2.1 Metabase起動テスト
- `npm run metabase:start`でMetabaseが起動するか確認
- http://localhost:8970 にアクセスできるか確認

#### 6.2.2 データベース接続テスト
- マスターデータベースに接続できるか確認
- シャーディングデータベースに接続できるか確認
- テーブル一覧が表示されるか確認
- データを閲覧できるか確認
- クエリを作成・実行できるか確認
- ダッシュボードを作成できるか確認

#### 6.2.3 設定ファイル管理テスト
- Metabaseの設定が`metabase/config/{env}/`ディレクトリに保存されるか確認
- 設定ファイルがGitで管理できるか確認
- 環境別に設定が分離されているか確認

### 6.3 環境別テスト

#### 6.3.1 開発環境テスト
- `APP_ENV=develop npm run metabase:start`で起動できるか確認
- デフォルト（`APP_ENV`未設定）で起動できるか確認

#### 6.3.2 ステージング環境テスト
- `APP_ENV=staging npm run metabase:start`で起動できるか確認
- 環境別の設定が正しく適用されるか確認

#### 6.3.3 本番環境テスト
- `APP_ENV=production npm run metabase:start`で起動できるか確認
- 環境別の設定が正しく適用されるか確認

### 6.4 docker-compose.ymlファイル分離テスト

#### 6.4.1 CloudBeaverとMetabaseの個別起動テスト
- CloudBeaverを起動できるか確認
- Metabaseを起動できるか確認
- 両方を個別に起動・停止できるか確認

## 7. セキュリティ考慮事項

### 7.1 データベースファイルへのアクセス

#### 7.1.1 読み取り専用マウント
- データベースファイルは読み取り専用でマウント（`:ro`オプション）
- Metabaseからデータベースファイルを誤って変更することを防止

#### 7.1.2 アクセス制御
- Metabaseは開発環境での使用を想定
- 本番環境での使用は想定しない（本番環境では適切なアクセス制御が必要）

### 7.2 認証設定

#### 7.2.1 Metabaseの認証
- Metabaseのデフォルト認証設定を確認
- 初回起動時に管理者アカウントを作成
- 必要に応じて認証を有効化

#### 7.2.2 ネットワークアクセス
- **開発環境**: 現在はローカルホスト（localhost）での開発が中心だが、サーバーに置く場合は外部から参照される可能性がある
- **Staging/Production環境**: 別サーバーで動作させるため、そのサーバーのIPアドレスやホスト名からアクセスする想定
- **アクセス制御**: 各環境で適切なアクセス制御（ファイアウォール、認証など）が必要。特にサーバーに置く場合は外部からのアクセスを想定したセキュリティ対策が必要

## 8. 運用・保守

### 8.1 起動・停止

#### 8.1.1 起動
```bash
# 開発環境（デフォルト）
npm run metabase:start

# 環境を指定
APP_ENV=develop npm run metabase:start
APP_ENV=staging npm run metabase:start
APP_ENV=production npm run metabase:start
```

#### 8.1.2 停止
```bash
npm run metabase:stop
```

#### 8.1.3 再起動
```bash
npm run metabase:restart
```

#### 8.1.4 ログ確認
```bash
npm run metabase:logs
```

### 8.2 データベース接続管理

#### 8.2.1 接続設定の追加
- Metabase Web UIから手動で接続設定を追加
- 接続設定はMetabaseコンテナ内に保存される

#### 8.2.2 接続設定の削除
- Metabase Web UIから手動で接続設定を削除

### 8.3 設定ファイル管理

#### 8.3.1 設定の作成
- Metabase Web UIで接続設定やダッシュボードを作成
- 設定は`metabase/config/{env}/`ディレクトリに保存される

#### 8.3.2 設定のGit管理
- `metabase/config/`ディレクトリをGitで管理
- 設定の変更をGitで追跡

### 8.4 トラブルシューティング

#### 8.4.1 コンテナが起動しない
- Dockerが起動しているか確認
- ポート8970が使用されていないか確認
- ログを確認（`npm run metabase:logs`）

#### 8.4.2 データベースに接続できない
- データベースファイルが存在するか確認
- マウント設定を確認
- ファイルのパスを確認

#### 8.4.3 設定ファイルを保存できない
- `metabase/config/{env}/`ディレクトリの権限を確認
- マウント設定を確認

## 9. 実装上の注意事項

### 9.1 Docker Compose設定

#### 9.1.1 イメージのバージョン
- `metabase/metabase:latest`を使用
- 必要に応じて特定のバージョンを指定

#### 9.1.2 ボリュームマウント
- データベースファイルは読み取り専用でマウント（`:ro`オプション）
- Metabase設定ディレクトリは読み書き可能でマウント

#### 9.1.3 環境変数
- `APP_ENV`環境変数をコンテナに渡す（環境別制御のため）
- `MB_DB_FILE`は固定値（`/metabase-data/metabase.db`）として`docker-compose.metabase.yml`に直接記載
- ポート番号は`docker-compose.metabase.yml`に固定値として定義（8970:3000）

### 9.2 package.jsonの更新

#### 9.2.1 依存関係
- 依存関係は不要（シェルスクリプトを使用）

#### 9.2.2 スクリプトの定義
- `scripts/metabase-start.sh`スクリプトを呼び出す
- 環境変数のデフォルト値はシェルスクリプト内で設定
- CloudBeaver用のスクリプトも更新（`docker-compose.cloudbeaver.yml`を使用）

#### 9.2.3 metabase-start.shスクリプトの作成
- `scripts/metabase-start.sh`を作成
- 実行権限を付与: `chmod +x scripts/metabase-start.sh`
- スクリプト内で`APP_ENV`環境変数を設定（未設定の場合は`develop`をデフォルト）

#### 9.2.4 cloudbeaver-start.shスクリプトの更新
- `scripts/cloudbeaver-start.sh`を更新
- `docker-compose.yml`を`docker-compose.cloudbeaver.yml`に変更

### 9.3 docker-compose.ymlファイルの分離

#### 9.3.1 既存ファイルのリネーム
- 既存の`docker-compose.yml`を`docker-compose.cloudbeaver.yml`にリネーム
- または、`docker-compose.cloudbeaver.yml`を新規作成して内容をコピー

#### 9.3.2 新規ファイルの作成
- `docker-compose.metabase.yml`を新規作成
- Metabaseサービスの定義を追加

### 9.4 Metabase設定ディレクトリの作成

#### 9.4.1 環境別設定ディレクトリの作成
- `metabase/config/develop/`ディレクトリを作成
- `metabase/config/staging/`ディレクトリを作成
- `metabase/config/production/`ディレクトリを作成
- 各ディレクトリに`.gitkeep`ファイルを追加（空ディレクトリをGitに含める）

#### 9.4.2 設定ファイルの管理
- Metabase起動後、Web UIから接続設定を行う
- 接続設定は`metabase/config/{env}/`ディレクトリに保存される
- 設定ファイルはGitで管理可能
- 環境別に設定が分離されるため、環境ごとに異なる接続設定を管理可能

### 9.5 ドキュメント整備

#### 9.5.1 README.mdの更新
- Metabaseの起動方法を追記
- 環境別の起動方法を記載
- 基本的な使用方法を記載
- データベース接続設定の手順を記載
- CloudBeaverとMetabaseの使い分けを記載

#### 9.5.2 docs/Metabase.mdの作成
- Metabase専用の詳細ドキュメントを作成
- 起動方法、接続設定、クエリ作成、ダッシュボード作成などの詳細を記載
- トラブルシューティング情報を記載

## 10. 参考情報

### 10.1 関連ドキュメント
- 要件定義書: `.kiro/specs/0023-metabase/requirements.md`
- CloudBeaver設計書: `.kiro/specs/0015-cloudbeaver/design.md`
- プロジェクトREADME: `README.md`
- シャーディング仕様: `docs/Sharding.md`

### 10.2 技術スタック
- **Metabase**: https://www.metabase.com/
- **Docker**: Docker Compose
- **npm**: npmスクリプト

### 10.3 参考リンク
- Metabase公式サイト: https://www.metabase.com/
- Metabase GitHub: https://github.com/metabase/metabase
- Metabase Docker: https://hub.docker.com/r/metabase/metabase
- Metabase ドキュメント: https://www.metabase.com/docs/
- Metabase Docker ドキュメント: https://www.metabase.com/docs/latest/installation-and-operation/running-metabase-on-docker
