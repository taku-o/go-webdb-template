# CloudBeaver導入実装タスク一覧

## 概要
CloudBeaver導入の実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: ディレクトリ構造の作成

#### - [x] タスク 1.1: CloudBeaver設定ディレクトリの作成
**目的**: CloudBeaver設定ファイル用の環境別ディレクトリを作成

**作業内容**:
- `cloudbeaver/config/develop/`ディレクトリを作成
- `cloudbeaver/config/staging/`ディレクトリを作成
- `cloudbeaver/config/production/`ディレクトリを作成
- 各ディレクトリに`.gitkeep`ファイルを追加（空ディレクトリをGitに含める）

**受け入れ基準**:
- `cloudbeaver/config/develop/`ディレクトリが存在する
- `cloudbeaver/config/staging/`ディレクトリが存在する
- `cloudbeaver/config/production/`ディレクトリが存在する
- 各ディレクトリに`.gitkeep`ファイルが存在する

---

#### - [x] タスク 1.2: Resource Manager用ディレクトリの作成
**目的**: Resource Managerに保存したスクリプト用のディレクトリを作成

**作業内容**:
- `cloudbeaver/scripts/`ディレクトリを作成
- `.gitkeep`ファイルを追加（空ディレクトリをGitに含める）

**受け入れ基準**:
- `cloudbeaver/scripts/`ディレクトリが存在する
- `.gitkeep`ファイルが存在する

---

### Phase 2: Docker Compose設定ファイルの作成

#### - [x] タスク 2.1: docker-compose.ymlの作成
**目的**: CloudBeaver用のDocker Compose設定ファイルを作成

**作業内容**:
- プロジェクトルートに`docker-compose.yml`を作成
- CloudBeaverサービスの定義を追加:
  - イメージ: `dbeaver/cloudbeaver:latest`
  - コンテナ名: `cloudbeaver`
  - 再起動ポリシー: `unless-stopped`
  - ポートマッピング: `${CLOUDBEAVER_PORT:-8978}:8978`
  - ボリュームマウント:
    - `./server/data:/data:ro` (データベースファイル、読み取り専用)
    - `./cloudbeaver/config/${APP_ENV:-develop}:/opt/cloudbeaver/workspace` (設定ディレクトリ、環境別)
    - `./cloudbeaver/scripts:/scripts` (Resource Manager用スクリプト)
  - 環境変数:
    - `APP_ENV=${APP_ENV:-develop}`
    - `CB_WORKSPACE=/opt/cloudbeaver/workspace`
  - ヘルスチェック設定

**受け入れ基準**:
- `docker-compose.yml`が作成されている
- CloudBeaverサービスの定義が正しく記述されている
- 環境別設定ディレクトリのマウントが正しく設定されている
- 環境変数の設定が正しい

---

### Phase 3: 起動スクリプトとpackage.jsonの作成

#### - [x] タスク 3.1: cloudbeaver-start.shスクリプトの作成
**目的**: CloudBeaver起動用のシェルスクリプトを作成

**作業内容**:
- `scripts/cloudbeaver-start.sh`を作成
- スクリプトの内容:
  ```bash
  #!/bin/bash
  # CloudBeaver起動スクリプト
  
  # APP_ENV環境変数が未設定の場合はdevelopをデフォルトとする
  export APP_ENV=${APP_ENV:-develop}
  
  # Docker ComposeでCloudBeaverを起動
  docker-compose up -d
  ```
- 実行権限を付与: `chmod +x scripts/cloudbeaver-start.sh`

**受け入れ基準**:
- `scripts/cloudbeaver-start.sh`が作成されている
- スクリプトの内容が正しい
- 実行権限が付与されている

---

#### - [x] タスク 3.2: package.jsonの作成
**目的**: プロジェクトルート用のpackage.jsonを作成し、npmスクリプトを定義

**作業内容**:
- プロジェクトルートに`package.json`を作成
- npmスクリプトを定義:
  - `cloudbeaver:start`: `./scripts/cloudbeaver-start.sh`
  - `cloudbeaver:stop`: `docker-compose down`
  - `cloudbeaver:logs`: `docker-compose logs -f cloudbeaver`
  - `cloudbeaver:restart`: `npm run cloudbeaver:stop && npm run cloudbeaver:start`

**受け入れ基準**:
- `package.json`が作成されている
- npmスクリプトが正しく定義されている
- 既存の`client/package.json`とは別ファイルである

---

### Phase 4: 動作確認

#### - [x] タスク 4.1: CloudBeaver起動確認（開発環境）
**目的**: 開発環境でCloudBeaverが正常に起動することを確認

**作業内容**:
- `npm run cloudbeaver:start`を実行
- または `APP_ENV=develop npm run cloudbeaver:start`を実行
- Docker ComposeでCloudBeaverコンテナが起動することを確認
- http://localhost:8978 にアクセスできることを確認
- CloudBeaverのWeb UIが表示されることを確認

**受け入れ基準**:
- CloudBeaverコンテナが正常に起動する
- http://localhost:8978 にアクセスできる
- CloudBeaverのWeb UIが表示される

---

#### - [x] タスク 4.2: CloudBeaver起動確認（ステージング環境）
**目的**: ステージング環境でCloudBeaverが正常に起動することを確認

**作業内容**:
- `APP_ENV=staging npm run cloudbeaver:start`を実行
- Docker ComposeでCloudBeaverコンテナが起動することを確認
- 環境別設定ディレクトリ（`cloudbeaver/config/staging/`）がマウントされていることを確認

**受け入れ基準**:
- CloudBeaverコンテナが正常に起動する
- ステージング環境用の設定ディレクトリがマウントされている

---

#### - [x] タスク 4.3: CloudBeaver起動確認（本番環境）
**目的**: 本番環境でCloudBeaverが正常に起動することを確認

**作業内容**:
- `APP_ENV=production npm run cloudbeaver:start`を実行
- Docker ComposeでCloudBeaverコンテナが起動することを確認
- 環境別設定ディレクトリ（`cloudbeaver/config/production/`）がマウントされていることを確認

**受け入れ基準**:
- CloudBeaverコンテナが正常に起動する
- 本番環境用の設定ディレクトリがマウントされている

---

#### - [x] タスク 4.4: CloudBeaver停止確認
**目的**: CloudBeaverが正常に停止することを確認

**作業内容**:
- `npm run cloudbeaver:stop`を実行
- Docker ComposeでCloudBeaverコンテナが停止することを確認

**受け入れ基準**:
- CloudBeaverコンテナが正常に停止する

---

### Phase 5: データベース接続設定

#### - [ ] タスク 5.1: マスターデータベース接続設定
**目的**: CloudBeaverからマスターデータベースに接続できるように設定

**作業内容**:
- CloudBeaver Web UI（http://localhost:8978）にアクセス
- データベース接続を追加
- SQLiteドライバーを選択
- 接続情報を入力:
  - 接続名: `master` または `Master Database`
  - データベースファイル: `/data/master.db`
- 接続をテスト
- 接続を保存
- 接続設定が`cloudbeaver/config/develop/`ディレクトリに保存されることを確認

**受け入れ基準**:
- マスターデータベースに接続できる
- 接続設定が`cloudbeaver/config/develop/`ディレクトリに保存される
- 接続設定ファイルがGitで管理できる

---

#### - [ ] タスク 5.2: シャーディングデータベース接続設定
**目的**: CloudBeaverからシャーディングデータベース（4つ）に接続できるように設定

**作業内容**:
- CloudBeaver Web UIで各シャーディングデータベースの接続を追加:
  - `sharding_db_1`: `/data/sharding_db_1.db`
  - `sharding_db_2`: `/data/sharding_db_2.db`
  - `sharding_db_3`: `/data/sharding_db_3.db`
  - `sharding_db_4`: `/data/sharding_db_4.db`
- 各接続をテスト
- 各接続を保存
- 接続設定が`cloudbeaver/config/develop/`ディレクトリに保存されることを確認

**受け入れ基準**:
- 4つのシャーディングデータベースに接続できる
- 接続設定が`cloudbeaver/config/develop/`ディレクトリに保存される
- 接続設定ファイルがGitで管理できる

---

#### - [ ] タスク 5.3: データベース接続動作確認
**目的**: 接続したデータベースでテーブル一覧の表示、データの閲覧、SQL実行ができることを確認

**作業内容**:
- マスターデータベースに接続
- テーブル一覧が表示されることを確認
- テーブルのデータを閲覧できることを確認
- SQLクエリを実行できることを確認
- シャーディングデータベース（4つ）でも同様の確認を行う

**受け入れ基準**:
- 各データベースのテーブル一覧が表示される
- 各データベースのデータを閲覧できる
- SQLクエリを実行できる

---

### Phase 6: Resource Manager設定

#### - [ ] タスク 6.1: Resource Manager設定
**目的**: Resource Managerの保存先を`cloudbeaver/scripts/`ディレクトリに設定

**作業内容**:
- CloudBeaver Web UIでResource Managerの設定を確認
- Resource Managerの保存先を`/scripts`ディレクトリに設定（必要に応じて）
- テスト用のSQLスクリプトを作成してResource Managerに保存
- 保存したスクリプトが`cloudbeaver/scripts/`ディレクトリに反映されることを確認

**受け入れ基準**:
- Resource Managerにスクリプトを保存できる
- 保存したスクリプトが`cloudbeaver/scripts/`ディレクトリに反映される
- スクリプトファイルがGitで管理できる

---

#### - [ ] タスク 6.2: Resource Manager動作確認
**目的**: Resource Managerの作成・編集・削除が正常に動作することを確認

**作業内容**:
- Resource Managerでスクリプトを作成
- スクリプトを編集
- スクリプトを削除
- 各操作が正常に動作することを確認

**受け入れ基準**:
- スクリプトの作成・編集・削除が正常に動作する
- 変更が`cloudbeaver/scripts/`ディレクトリに反映される

---

### Phase 7: 環境別設定の確認

#### - [x] タスク 7.1: 環境別設定ディレクトリの動作確認
**目的**: 環境別に設定ディレクトリが正しくマウントされ、設定が分離されることを確認

**作業内容**:
- 開発環境でCloudBeaverを起動し、接続設定を行う
- 接続設定が`cloudbeaver/config/develop/`ディレクトリに保存されることを確認
- ステージング環境でCloudBeaverを起動
- ステージング環境用の設定ディレクトリ（`cloudbeaver/config/staging/`）がマウントされていることを確認
- 本番環境でCloudBeaverを起動
- 本番環境用の設定ディレクトリ（`cloudbeaver/config/production/`）がマウントされていることを確認

**受け入れ基準**:
- 環境別に設定ディレクトリが正しくマウントされる
- 環境別に設定が分離される
- 各環境の設定ファイルがGitで管理できる

---

### Phase 8: ドキュメント整備

#### - [x] タスク 8.1: Database-Viewer.mdの作成
**目的**: CloudBeaver専用の詳細ドキュメントを作成

**作業内容**:
- `docs/Database-Viewer.md`を作成
- 以下の内容を含める:
  - 概要（主な機能、役割分担）
  - 前提条件
  - 起動方法（基本的な起動、環境別の起動）
  - 停止方法
  - その他のコマンド（ログ確認、再起動）
  - データベース接続設定（初回設定、マスターデータベース、シャーディングデータベース）
  - データベース操作（テーブル一覧、データ閲覧、SQL実行）
  - Resource Manager（スクリプトの作成、保存場所、使用方法）
  - 環境別設定（設定ディレクトリの構造、管理方法）
  - トラブルシューティング（コンテナ起動、データベース接続、Resource Manager、設定ファイル）
  - セキュリティ考慮事項（データベースファイルへのアクセス、認証設定、ネットワークアクセス）
  - 設定ファイルの管理（Git管理、設定ファイルの構造、共有方法）
  - 参考情報（関連ドキュメント、技術スタック、参考リンク）

**受け入れ基準**:
- `docs/Database-Viewer.md`が作成されている
- すべての主要な機能が記載されている
- 起動方法が記載されている
- 環境別の起動方法が記載されている
- データベース接続設定の手順が記載されている
- Resource Managerの使用方法が記載されている
- トラブルシューティング情報が記載されている

---

#### - [x] タスク 8.2: README.mdへの追記
**目的**: README.mdにCloudBeaverの簡単な説明とDatabase-Viewer.mdへのリンクを追記

**作業内容**:
- `README.md`にCloudBeaverのセクションを追加
- 簡単な説明を記載
- 起動方法を簡潔に記載:
  - 開発環境: `npm run cloudbeaver:start` または `APP_ENV=develop npm run cloudbeaver:start`
  - ステージング環境: `APP_ENV=staging npm run cloudbeaver:start`
  - 本番環境: `APP_ENV=production npm run cloudbeaver:start`
  - デフォルトは`develop`環境
- 停止方法を記載: `npm run cloudbeaver:stop`
- 詳細な使用方法は`docs/Database-Viewer.md`を参照する旨を記載

**受け入れ基準**:
- `README.md`にCloudBeaverのセクションが追加されている
- 起動方法が記載されている
- 環境別の起動方法が記載されている
- `docs/Database-Viewer.md`へのリンクが記載されている

---

### Phase 9: 最終確認

#### - [ ] タスク 9.1: 受け入れ基準の確認
**目的**: 要件定義書の受け入れ基準をすべて満たしていることを確認

**作業内容**:
- 要件定義書の受け入れ基準（6.1～6.6）を確認
- 各項目が満たされていることを確認
- 不足している項目があれば対応

**受け入れ基準**:
- すべての受け入れ基準が満たされている

---

#### - [ ] タスク 9.2: Git管理の確認
**目的**: CloudBeaver関連のファイルがGitで管理できることを確認

**作業内容**:
- `cloudbeaver/config/`ディレクトリがGitで管理できることを確認
- `cloudbeaver/scripts/`ディレクトリがGitで管理できることを確認
- `.gitignore`で不要なファイルが除外されていることを確認（必要に応じて）

**受け入れ基準**:
- CloudBeaver関連のファイルがGitで管理できる
- 不要なファイルが除外されている（必要に応じて）

---

## 実装順序の推奨

1. **Phase 1**: ディレクトリ構造の作成
2. **Phase 2**: Docker Compose設定ファイルの作成
3. **Phase 3**: 起動スクリプトとpackage.jsonの作成
4. **Phase 4**: 動作確認（開発環境）
5. **Phase 5**: データベース接続設定
6. **Phase 6**: Resource Manager設定
7. **Phase 7**: 環境別設定の確認
8. **Phase 8**: ドキュメント整備
9. **Phase 9**: 最終確認

## 注意事項

- 各タスクは順番に実行することを推奨
- タスクの実行前に、前のタスクが完了していることを確認
- エラーが発生した場合は、原因を特定してから次のタスクに進む
- 環境別の動作確認は、開発環境で動作確認が完了してから実施

