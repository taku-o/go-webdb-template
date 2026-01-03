# Apache Superset導入設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、非エンジニア向けのデータビューワとしてApache Supersetを導入するための詳細設計を定義する。Docker Composeを使用してApache Supersetを構築し、既存のPostgreSQLデータベースに接続してデータを可視化できるようにする。

### 1.2 設計の範囲
- Apache Superset環境の構築（Docker Compose設定）
- PostgreSQL接続設定（Dockerネットワーク経由）
- 起動スクリプトの実装
- データ永続化設計
- 初期設定設計（デフォルト管理者アカウント）
- エラーハンドリング設計
- テスト戦略

### 1.3 設計方針
- **既存パターンの遵守**: 既存のDocker Compose設定や起動スクリプトのパターンに従う
- **シンプルな実装**: 最小限の機能のみを実装し、複雑な機能は将来の拡張項目とする
- **データ永続化**: Dockerボリュームを使用してデータを永続化
- **PostgreSQL接続**: Dockerネットワーク経由でPostgreSQLに接続
- **開発環境重視**: 本実装は主に開発環境での利用を想定

## 2. アーキテクチャ設計

### 2.1 ディレクトリ構造

#### 2.1.1 変更前の構造
```
go-webdb-template/
├── docker-compose.postgres.yml
├── docker-compose.metabase.yml
├── docker-compose.cloudbeaver.yml
├── scripts/
│   ├── start-postgres.sh
│   ├── metabase-start.sh
│   └── cloudbeaver-start.sh
└── ...
```

#### 2.1.2 変更後の構造
```
go-webdb-template/
├── docker-compose.postgres.yml
├── docker-compose.apache-superset.yml    # 新規: Apache Superset用Docker Compose設定
├── docker-compose.metabase.yml
├── docker-compose.cloudbeaver.yml
├── scripts/
│   ├── start-postgres.sh
│   ├── start-apache-superset.sh          # 新規: Apache Superset起動スクリプト
│   ├── metabase-start.sh
│   └── cloudbeaver-start.sh
├── apache-superset/                      # 新規: Apache Supersetデータディレクトリ
│   └── data/                             # データ永続化用
│       ├── superset.db                   # Superset内部データベース
│       ├── config/                       # 設定ファイル
│       └── uploads/                      # アップロードファイル
└── ...
```

### 2.2 システム構成図

```
┌─────────────────────────────────────────────────────────────┐
│                    Docker Host Machine                        │
│                                                               │
│  ┌──────────────────┐         ┌──────────────────────┐      │
│  │  PostgreSQL      │         │  Apache Superset     │      │
│  │  Container       │◄────────┤  Container           │      │
│  │  (port: 5432)    │         │  (port: 8088)        │      │
│  └──────────────────┘         └──────────────────────┘      │
│         │                              │                     │
│         │                              │                     │
│         └──────────────┬───────────────┘                     │
│                        │                                     │
│                  Docker Network                               │
│                  (postgres-network)                            │
│                                                               │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Volume Mounts                                        │   │
│  │  - ./postgres/data → /var/lib/postgresql/data        │   │
│  │  - ./apache-superset/data → /app/superset_home       │   │
│  └──────────────────────────────────────────────────────┘   │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

### 2.3 データフロー

```
┌─────────────────────────────────────────────────────────────┐
│              1. PostgreSQL起動                                │
│              ./scripts/start-postgres.sh start                │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              2. Apache Superset起動                          │
│              ./scripts/start-apache-superset.sh              │
│              - Docker Composeでコンテナ起動                   │
│              - データディレクトリをマウント                   │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              3. 初回起動時の初期化                            │
│              - デフォルト管理者アカウント作成（admin/admin）   │
│              - Superset内部データベース初期化                 │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              4. Web UIアクセス                                │
│              http://localhost:8088                            │
│              - 管理者アカウントでログイン                      │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              5. PostgreSQLデータソース接続設定                │
│              - Web UIから手動で接続設定                        │
│              - 接続情報: postgres:5432/webdb                 │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              6. データ可視化                                  │
│              - SQL Labでクエリ実行                             │
│              - チャート作成                                   │
│              - ダッシュボード作成                             │
└─────────────────────────────────────────────────────────────┘
```

## 3. Docker Compose設定設計

### 3.1 docker-compose.apache-superset.yml

```yaml
version: '3.8'

services:
  apache-superset:
    image: apache/superset:latest
    container_name: apache-superset
    restart: unless-stopped
    ports:
      - "8088:8088"
    volumes:
      - ./apache-superset/data:/app/superset_home
    environment:
      - SUPERSET_SECRET_KEY=${SUPERSET_SECRET_KEY:-dev-secret-key-change-in-production}
      - SUPERSET_CONFIG_PATH=/app/superset_home/config
    networks:
      - postgres-network
    depends_on:
      - postgres
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8088/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 180s

networks:
  postgres-network:
    external: true
    name: postgres-network
```

### 3.2 設定項目の説明

#### 3.2.1 イメージ
- **イメージ**: `apache/superset:latest`
  - 最新の安定版を使用
  - 必要に応じて特定バージョンを指定可能（例: `apache/superset:3.0.0`）

#### 3.2.2 ポート
- **ホストポート**: `8088`
- **コンテナポート**: `8088`
- Apache SupersetのWeb UIにアクセスするためのポート

#### 3.2.3 ボリューム
- **マウントパス**: `./apache-superset/data:/app/superset_home`
  - ホスト側: `./apache-superset/data`（プロジェクトルートからの相対パス）
  - コンテナ側: `/app/superset_home`（Apache Supersetのデフォルトデータディレクトリ）
  - データ永続化のため必須

#### 3.2.4 環境変数
- **SUPERSET_SECRET_KEY**: 
  - デフォルト値: `dev-secret-key-change-in-production`
  - 開発環境用の簡易的なシークレットキー
  - 本番環境では適切なランダムな値を設定すること
- **SUPERSET_CONFIG_PATH**: 
  - デフォルト値: `/app/superset_home/config`
  - 設定ファイルのパス（オプション）

#### 3.2.5 ネットワーク
- **ネットワーク名**: `postgres-network`
- **タイプ**: `external: true`
- PostgreSQLコンテナと同じネットワークに接続
- PostgreSQLに`postgres`というホスト名で接続可能

#### 3.2.6 依存関係
- **depends_on**: `postgres`
- PostgreSQLコンテナの起動を待つ（ただし、PostgreSQLの完全な起動は保証しない）
- 実際の接続はApache Superset起動後に手動で設定

#### 3.2.7 ヘルスチェック
- **テストコマンド**: `curl -f http://localhost:8088/health`
- **間隔**: 30秒
- **タイムアウト**: 10秒
- **リトライ**: 3回
- **開始待機時間**: 180秒（初回起動時は時間がかかるため）

### 3.3 PostgreSQLネットワーク設定

PostgreSQLのdocker-compose.postgres.ymlにネットワーク設定を追加する必要があります：

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: postgres
    ports:
      - "5432:5432"
    environment:
      POSTGRES_USER: webdb
      POSTGRES_PASSWORD: webdb
      POSTGRES_DB: webdb
    volumes:
      - ./postgres/data:/var/lib/postgresql/data
    restart: unless-stopped
    networks:
      - postgres-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U webdb"]
      interval: 10s
      timeout: 5s
      retries: 5

networks:
  postgres-network:
    name: postgres-network
    driver: bridge
```

**注意**: PostgreSQLのdocker-compose.postgres.ymlは既存ファイルのため、この変更は別途実装時に確認が必要です。

## 4. 起動スクリプト設計

### 4.1 scripts/start-apache-superset.sh

```bash
#!/bin/bash
# Apache Superset起動スクリプト

# Docker ComposeでApache Supersetを起動
docker-compose -f docker-compose.apache-superset.yml up -d

# 起動確認メッセージ
echo "Apache Superset started."
echo "Access URL: http://localhost:8088"
echo "Default credentials: admin/admin"
```

### 4.2 スクリプトの特徴

- **既存パターンに従う**: `metabase-start.sh`や`cloudbeaver-start.sh`と同じパターン
- **シンプルな実装**: Docker Composeコマンドを実行するだけ
- **ユーザーフレンドリー**: 起動後のアクセス情報を表示

### 4.3 停止スクリプト（オプション）

停止は`docker-compose down`コマンドで実行可能ですが、必要に応じてスクリプトを作成：

```bash
#!/bin/bash
# Apache Superset停止スクリプト

docker-compose -f docker-compose.apache-superset.yml down
echo "Apache Superset stopped."
```

## 5. データ永続化設計

### 5.1 データディレクトリ構造

```
apache-superset/
└── data/
    ├── superset.db              # Superset内部データベース（SQLite）
    ├── config/                  # 設定ファイル
    │   └── superset_config.py  # カスタム設定ファイル（オプション）
    └── uploads/                 # アップロードファイル
```

### 5.2 永続化されるデータ

- **Superset内部データベース**: ユーザー情報、ダッシュボード、チャート、データソース接続情報など
- **設定ファイル**: カスタム設定（オプション）
- **アップロードファイル**: CSVファイルなどのアップロードデータ

### 5.3 データ永続化の仕組み

- Dockerボリュームを使用してホストマシンのディレクトリにマウント
- コンテナ再起動後もデータが保持される
- コンテナを削除してもデータは保持される（ボリュームが削除されない限り）

### 5.4 .gitignore設定

データディレクトリの大部分はGit管理から除外しますが、ダッシュボード設定を含む`superset.db`はGit管理対象とします：

```
apache-superset/data/
!apache-superset/data/superset.db
```

これにより、以下のファイルがGit管理されます：
- `apache-superset/data/superset.db` - Superset内部データベース（ダッシュボード設定、チャート設定、ユーザー情報、データソース接続情報を含む）

以下のファイルはGit管理から除外されます：
- `apache-superset/data/uploads/` - アップロードファイル
- `apache-superset/data/config/` - 設定ファイル（必要に応じて個別に追加可能）

**セキュリティ上の注意**:
- `superset.db`にはユーザー情報（パスワードハッシュ）とデータソース接続情報（パスワード含む可能性）が含まれます
- 本番環境では機密情報の取り扱いに注意してください

## 6. PostgreSQL接続設計

### 6.1 接続方法

Apache SupersetからPostgreSQLに接続する方法は2つあります：

#### 6.1.1 Dockerネットワーク経由（推奨）

- **ホスト名**: `postgres`
- **ポート**: `5432`
- **データベース名**: `webdb`
- **ユーザー名**: `webdb`
- **パスワード**: `webdb`
- **接続文字列**: `postgresql://webdb:webdb@postgres:5432/webdb`

#### 6.1.2 ホストマシン経由

- **ホスト名**: `host.docker.internal`（Mac/Windows）または`172.17.0.1`（Linux）
- **ポート**: `5432`
- **データベース名**: `webdb`
- **ユーザー名**: `webdb`
- **パスワード**: `webdb`
- **接続文字列**: `postgresql://webdb:webdb@host.docker.internal:5432/webdb`

### 6.2 接続設定手順

1. Apache SupersetのWeb UIにアクセス（http://localhost:8088）
2. 管理者アカウントでログイン（admin/admin）
3. 「Settings」→「Database Connections」を選択
4. 「+ Database」ボタンをクリック
5. データベースタイプで「PostgreSQL」を選択
6. 接続情報を入力（上記の接続方法を参照）
7. 「Test Connection」で接続を確認
8. 「Connect」で接続設定を保存

### 6.3 接続エラーの対処

- **接続できない場合**: 
  - PostgreSQLが起動しているか確認
  - Dockerネットワークが正しく設定されているか確認
  - 接続情報（ホスト名、ポート、認証情報）が正しいか確認
- **タイムアウトエラー**: 
  - PostgreSQLのヘルスチェックを確認
  - ネットワーク設定を確認

## 7. 初期設定設計

### 7.1 デフォルト管理者アカウント

Apache Supersetの初回起動時、デフォルトの管理者アカウントが自動的に作成されます：

- **ユーザー名**: `admin`
- **パスワード**: `admin`
- **メールアドレス**: 設定不要（開発環境）
- **ロール**: Administrator

### 7.2 初回起動時の処理

1. Apache Supersetコンテナ起動
2. データディレクトリが空の場合、初期化処理を実行
3. Superset内部データベース（SQLite）を作成
4. デフォルト管理者アカウントを作成
5. 初期設定を完了

### 7.3 パスワード変更

初回ログイン時にパスワード変更を求められる場合がありますが、開発環境では`admin`のままでも問題ありません。

## 8. エラーハンドリング設計

### 8.1 起動時のエラー

#### 8.1.1 ポート競合
- **問題**: ポート8088が既に使用されている
- **対処**: 使用しているプロセスを確認し、停止するかポート番号を変更

#### 8.1.2 データディレクトリの権限エラー
- **問題**: データディレクトリへの書き込み権限がない
- **対処**: ディレクトリの権限を確認し、適切な権限を設定

#### 8.1.3 イメージの取得エラー
- **問題**: Dockerイメージの取得に失敗
- **対処**: ネットワーク接続を確認し、Docker Hubへのアクセスを確認

### 8.2 実行時のエラー

#### 8.2.1 PostgreSQL接続エラー
- **問題**: PostgreSQLに接続できない
- **対処**: 
  - PostgreSQLが起動しているか確認
  - ネットワーク設定を確認
  - 接続情報を確認

#### 8.2.2 メモリ不足
- **問題**: Apache Supersetがメモリ不足で起動できない
- **対処**: 
  - 他のメモリを多く使用するコンテナ（Metabase、CloudBeaver）を停止
  - システムメモリを確認

### 8.3 ログ確認方法

```bash
# Apache Supersetのログを確認
docker-compose -f docker-compose.apache-superset.yml logs -f

# 特定のサービスのログのみ確認
docker-compose -f docker-compose.apache-superset.yml logs apache-superset
```

## 9. テスト戦略

### 9.1 機能テスト

#### 9.1.1 起動テスト
- Apache Supersetが正常に起動することを確認
- Web UIにアクセスできることを確認
- ヘルスチェックが正常に動作することを確認

#### 9.1.2 データ永続化テスト
- コンテナ再起動後もデータが保持されることを確認
- ダッシュボードやチャートが保持されることを確認

#### 9.1.3 PostgreSQL接続テスト
- PostgreSQLに接続できることを確認
- データを閲覧できることを確認
- クエリを実行できることを確認

### 9.2 統合テスト

#### 9.2.1 起動順序テスト
- PostgreSQLを先に起動
- Apache Supersetを起動
- 両方が正常に動作することを確認

#### 9.2.2 ネットワークテスト
- Dockerネットワーク経由でPostgreSQLに接続できることを確認
- 接続情報が正しく設定されていることを確認

### 9.3 パフォーマンステスト

#### 9.3.1 起動時間テスト
- Apache Supersetの起動時間を測定
- 初回起動と2回目以降の起動時間を比較

#### 9.3.2 メモリ使用量テスト
- Apache Supersetのメモリ使用量を確認
- 他のコンテナとの共存可能性を確認

## 10. セキュリティ考慮事項

### 10.1 認証設定

- **開発環境**: デフォルトのadmin/adminで問題なし
- **本番環境**: 適切なパスワードポリシーを実装
- **パスワード変更**: 初回ログイン時にパスワード変更を推奨

### 10.2 ネットワークセキュリティ

- **ローカルホストのみ**: 開発環境ではlocalhostでのみアクセス可能
- **外部アクセス**: 本番環境では適切なファイアウォール設定が必要

### 10.3 データベース接続

- **接続情報**: 機密情報のため、環境変数やシークレット管理システムを使用
- **Git管理**: 接続情報をGitにコミットしない

### 10.4 シークレットキー

- **開発環境**: 簡易的なシークレットキーで問題なし
- **本番環境**: 適切なランダムなシークレットキーを生成

## 11. 保守性・拡張性

### 11.1 既存パターンの遵守

- Docker Compose設定は既存のパターンに従う
- 起動スクリプトは既存のパターンに従う
- ディレクトリ構造は既存のパターンに従う

### 11.2 拡張性

- 将来的に複数のデータソースを追加可能
- 設定ファイルから接続情報を読み込む機能を追加可能
- 環境別設定（develop/staging/production）を追加可能

### 11.3 保守性

- シンプルな実装により、保守が容易
- 既存パターンに従うことで、チームメンバーが理解しやすい
- ドキュメントを充実させることで、運用が容易

## 12. 実装時の注意事項

### 12.1 PostgreSQLネットワーク設定

PostgreSQLのdocker-compose.postgres.ymlにネットワーク設定を追加する必要があります。既存ファイルの変更のため、実装時に確認が必要です。

### 12.2 データディレクトリの作成

初回起動前にデータディレクトリを作成する必要はありませんが、.gitignoreに追加する必要があります。

### 12.3 起動順序

PostgreSQLを先に起動してからApache Supersetを起動することを推奨します。

### 12.4 メモリ使用量

Apache Supersetはメモリを多く使用するため、MetabaseやCloudBeaverと同時に起動しないことを推奨します。

## 13. 参考情報

### 13.1 関連ドキュメント

- [要件定義書](requirements.md)
- [Apache Superset公式ドキュメント](docs/Apache-Superset.md)
- [License Survey](docs/License-Survey.md)

### 13.2 技術スタック

- **Apache Superset**: https://superset.apache.org/
- **Docker Compose**: 既存のパターンに従う
- **PostgreSQL**: 既存のdocker-compose.postgres.ymlを使用

### 13.3 参考リンク

- [Apache Superset公式サイト](https://superset.apache.org/)
- [Apache Superset GitHub](https://github.com/apache/superset)
- [Apache Superset Docker ドキュメント](https://superset.apache.org/docs/installation/installing-superset-using-docker-compose)
