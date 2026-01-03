# Apache Superset導入実装タスク一覧

## 概要
Apache Superset導入の実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: PostgreSQLネットワーク設定

#### タスク 1.1: docker-compose.postgres.ymlへのネットワーク設定追加
**目的**: PostgreSQLコンテナをDockerネットワークに参加させ、Apache Supersetから接続可能にする

**作業内容**:
- `docker-compose.postgres.yml`ファイルを開く
- `services.postgres`セクションに`networks`セクションを追加:
  - `- postgres-network`
- `networks`セクションを追加（ファイル末尾）:
  - `postgres-network`ネットワークを定義
  - `name: postgres-network`
  - `driver: bridge`
- 既存の設定（ports、environment、volumes、healthcheck等）は変更しない

**受け入れ基準**:
- `docker-compose.postgres.yml`にネットワーク設定が追加されている
- 既存の設定が変更されていない
- PostgreSQLコンテナが`postgres-network`ネットワークに参加できること
- ネットワーク名が`postgres-network`であること

---

### Phase 2: Docker Compose設定とデータディレクトリ

#### タスク 2.1: docker-compose.apache-superset.ymlの作成
**目的**: Apache SupersetをDocker Composeで起動するための設定ファイルを作成する

**作業内容**:
- `docker-compose.apache-superset.yml`ファイルを新規作成
- `version: '3.8'`を指定
- `services.apache-superset`セクションを定義:
  - `image: apache/superset:latest`
  - `container_name: apache-superset`
  - `restart: unless-stopped`
  - `ports: "8088:8088"`
  - `volumes: ./apache-superset/data:/app/superset_home`
  - `environment`:
    - `SUPERSET_SECRET_KEY=${SUPERSET_SECRET_KEY:-dev-secret-key-change-in-production}`
    - `SUPERSET_CONFIG_PATH=/app/superset_home/config`
  - `networks: - postgres-network`
  - `depends_on: - postgres`
  - `healthcheck`:
    - `test: ["CMD", "curl", "-f", "http://localhost:8088/health"]`
    - `interval: 30s`
    - `timeout: 10s`
    - `retries: 3`
    - `start_period: 180s`
- `networks.postgres-network`セクションを定義:
  - `external: true`
  - `name: postgres-network`

**受け入れ基準**:
- `docker-compose.apache-superset.yml`ファイルが作成されている
- すべての設定項目が正しく定義されている
- ポート8088が公開されている
- データ永続化のボリューム設定が含まれている
- PostgreSQLネットワークに接続する設定が含まれている
- ヘルスチェック設定が含まれている

---

#### タスク 2.2: データディレクトリの.gitignore設定
**目的**: Apache Supersetのデータディレクトリのうち、ダッシュボード設定を含む`superset.db`をGit管理対象とし、その他のファイルを除外する

**作業内容**:
- `.gitignore`ファイルを開く
- `apache-superset/data/`を追加（ディレクトリ全体を除外）
- `!apache-superset/data/superset.db`を追加（`superset.db`を明示的にGit管理対象とする）
- 既存のMetabaseパターン（`!metabase/config/**/*.mv.db`）を参考にする

**受け入れ基準**:
- `.gitignore`に`apache-superset/data/`が追加されている
- `.gitignore`に`!apache-superset/data/superset.db`が追加されている
- 既存の設定が変更されていない
- `superset.db`がGit管理対象となること
- `apache-superset/data/uploads/`など、その他のファイルがGit管理から除外されること

---

### Phase 3: 起動スクリプトの実装

#### タスク 3.1: scripts/start-apache-superset.shの作成
**目的**: Apache Supersetを起動するスクリプトを作成する

**作業内容**:
- `scripts/start-apache-superset.sh`ファイルを新規作成
- 実行権限を付与（`chmod +x`）
- 既存の`metabase-start.sh`や`cloudbeaver-start.sh`と同じパターンで実装
- `docker-compose -f docker-compose.apache-superset.yml up -d`を実行
- 起動確認メッセージを出力:
  - "Apache Superset started."
  - "Access URL: http://localhost:8088"
  - "Default credentials: admin/admin"

**受け入れ基準**:
- `scripts/start-apache-superset.sh`ファイルが作成されている
- 実行権限が付与されている
- `./scripts/start-apache-superset.sh`でApache Supersetが起動できること
- 起動確認メッセージが表示されること
- 既存の起動スクリプトと同じパターンに従っている

---

### Phase 4: 動作確認とテスト

#### タスク 4.1: Apache Superset起動確認
**目的**: Apache Supersetが正常に起動することを確認する

**作業内容**:
- PostgreSQLを起動（`./scripts/start-postgres.sh start`）
- Apache Supersetを起動（`./scripts/start-apache-superset.sh`）
- コンテナの状態を確認（`docker-compose -f docker-compose.apache-superset.yml ps`）
- ログを確認（`docker-compose -f docker-compose.apache-superset.yml logs -f`）
- Web UIにアクセス（http://localhost:8088）
- ヘルスチェックが正常に動作することを確認

**受け入れ基準**:
- PostgreSQLが正常に起動していること
- Apache Supersetコンテナが正常に起動していること
- Web UIにアクセスできること
- ヘルスチェックが正常に動作すること
- エラーが発生していないこと

---

#### タスク 4.2: デフォルト管理者アカウント確認
**目的**: デフォルトの管理者アカウント（admin/admin）でログインできることを確認する

**作業内容**:
- Web UIにアクセス（http://localhost:8088）
- ログインフォームにアクセス
- ユーザー名: `admin`、パスワード: `admin`でログイン
- ログインが成功することを確認
- ダッシュボード画面が表示されることを確認

**受け入れ基準**:
- デフォルト管理者アカウントでログインできること
- ダッシュボード画面が表示されること
- エラーが発生していないこと

---

#### タスク 4.3: PostgreSQL接続設定確認
**目的**: Apache SupersetからPostgreSQLに接続できることを確認する

**作業内容**:
- Apache SupersetのWeb UIにログイン
- 「Settings」→「Database Connections」を選択
- 「+ Database」ボタンをクリック
- データベースタイプで「PostgreSQL」を選択
- 接続情報を入力:
  - **接続方法1（Dockerネットワーク経由）**:
    - Host: `postgres`
    - Port: `5432`
    - Database name: `webdb`
    - Username: `webdb`
    - Password: `webdb`
    - SQLAlchemy URI: `postgresql://webdb:webdb@postgres:5432/webdb`
  - **接続方法2（ホストマシン経由、接続方法1が失敗した場合）**:
    - Host: `host.docker.internal`（Mac/Windows）または`172.17.0.1`（Linux）
    - Port: `5432`
    - Database name: `webdb`
    - Username: `webdb`
    - Password: `webdb`
    - SQLAlchemy URI: `postgresql://webdb:webdb@host.docker.internal:5432/webdb`
- 「Test Connection」で接続を確認
- 「Connect」で接続設定を保存
- データソース一覧にPostgreSQLが表示されることを確認

**受け入れ基準**:
- PostgreSQLに接続できること
- 接続設定が保存されること
- データソース一覧にPostgreSQLが表示されること
- エラーが発生していないこと

---

#### タスク 4.4: データ閲覧確認
**目的**: Apache SupersetからPostgreSQLのデータを閲覧できることを確認する

**作業内容**:
- Apache SupersetのWeb UIにログイン
- 「SQL Lab」を選択
- データソースでPostgreSQLを選択
- 簡単なSQLクエリを実行（例: `SELECT * FROM information_schema.tables LIMIT 10;`）
- クエリ結果が表示されることを確認
- テーブル一覧が表示されることを確認

**受け入れ基準**:
- SQL LabからPostgreSQLに接続できること
- SQLクエリを実行できること
- クエリ結果が表示されること
- テーブル一覧が表示されること
- エラーが発生していないこと

---

#### タスク 4.5: データ永続化確認
**目的**: Apache Supersetのデータが永続化されることを確認する

**作業内容**:
- Apache Supersetを起動
- デフォルト管理者アカウントでログイン
- PostgreSQLデータソースを接続設定
- コンテナを停止（`docker-compose -f docker-compose.apache-superset.yml down`）
- データディレクトリ（`apache-superset/data/`）の内容を確認
- コンテナを再起動（`./scripts/start-apache-superset.sh`）
- 再度ログインして、データソース接続設定が保持されていることを確認

**受け入れ基準**:
- データディレクトリにデータが保存されていること
- コンテナ再起動後もデータが保持されること
- データソース接続設定が保持されること
- エラーが発生していないこと

---

### Phase 5: エラーハンドリング確認

#### タスク 5.1: PostgreSQL未起動時のエラーハンドリング確認
**目的**: PostgreSQLが起動していない場合のエラーハンドリングを確認する

**作業内容**:
- PostgreSQLを停止（`./scripts/start-postgres.sh stop`）
- Apache Supersetを起動（`./scripts/start-apache-superset.sh`）
- Apache Supersetが正常に起動することを確認（PostgreSQLがなくても起動できる）
- Web UIにアクセスしてログインできることを確認
- PostgreSQLデータソース接続を試みる
- 接続エラーが適切に表示されることを確認
- PostgreSQLを再起動（`./scripts/start-postgres.sh start`）
- 再度接続を試みて、接続が成功することを確認

**受け入れ基準**:
- PostgreSQLが停止していてもApache Supersetが起動できること
- 接続エラーが適切に表示されること
- PostgreSQL起動後に接続が成功すること
- エラーメッセージが分かりやすいこと

---

#### タスク 5.2: ポート競合エラー確認
**目的**: ポート8088が既に使用されている場合のエラーハンドリングを確認する

**作業内容**:
- ポート8088を使用しているプロセスを確認（`lsof -i :8088`または`netstat -an | grep 8088`）
- ポート8088が使用されていないことを確認
- Apache Supersetを起動
- 正常に起動することを確認
- （オプション）別のプロセスでポート8088を使用してからApache Supersetを起動し、エラーメッセージを確認

**受け入れ基準**:
- ポート8088が使用されていない場合、正常に起動できること
- ポート競合エラーが適切に表示されること（該当する場合）
- エラーメッセージが分かりやすいこと

---

### Phase 6: ドキュメント確認

#### タスク 6.1: ドキュメントの整合性確認
**目的**: 実装内容とドキュメントの整合性を確認する

**作業内容**:
- `docs/Apache-Superset.md`を確認
- 実装内容とドキュメントの内容が一致していることを確認
- 必要に応じてドキュメントを更新
- 起動手順が正しく記載されていることを確認
- 接続設定手順が正しく記載されていることを確認

**受け入れ基準**:
- ドキュメントと実装内容が一致していること
- 起動手順が正しく記載されていること
- 接続設定手順が正しく記載されていること
- エラーハンドリングの説明が含まれていること

---

## 実装順序の推奨

1. **Phase 1**: PostgreSQLネットワーク設定（インフラ準備）
2. **Phase 2**: Docker Compose設定とデータディレクトリ（設定ファイル作成）
3. **Phase 3**: 起動スクリプトの実装（運用準備）
4. **Phase 4**: 動作確認とテスト（機能確認）
5. **Phase 5**: エラーハンドリング確認（堅牢性確認）
6. **Phase 6**: ドキュメント確認（ドキュメント整合性）

## 注意事項

- 各タスクは独立して実装可能な粒度で分解されている
- タスクの順序は推奨順序であり、必要に応じて調整可能
- 実装時は既存のコードスタイルとパターンに従うこと
- PostgreSQLを先に起動してからApache Supersetを起動すること
- ダッシュボード設定を含む`apache-superset/data/superset.db`はGit管理対象とすること
- `apache-superset/data/uploads/`など、その他のファイルはGit管理から除外すること
- エラーハンドリングを適切に実装すること
- テストは実装と並行して進めることを推奨
