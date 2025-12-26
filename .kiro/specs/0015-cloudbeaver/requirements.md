# CloudBeaver導入要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #27
- **Issueタイトル**: データの管理機能としてのCloudBeaverの導入
- **Feature名**: 0015-cloudbeaver
- **作成日**: 2025-01-27

### 1.2 目的
データ操作用の管理アプリとしてCloudBeaverを導入し、Webベースのデータベース管理ツールを提供する。
本プロジェクトでは管理系アプリを2つ用意する方針であり、CloudBeaverはデータ操作用、GoAdminはカスタム処理用として役割を分担する。

### 1.3 スコープ
- CloudBeaverをDockerで動作させる
- 起動コマンドをpackage.jsonに定義する
- プロジェクト内のディレクトリをDockerでマウントさせて、Resource Managerに保存したスクリプトなどをGitに保存させる
- 起動させたCloudBeaverでデータベースを参照する

## 2. 背景・現状分析

### 2.1 現在の実装
- **データベース構成**:
  - マスターデータベース: `server/data/master.db`（SQLite、newsテーブル、GoAdmin関連テーブル）
  - シャーディングデータベース: `server/data/sharding_db_1.db` ～ `server/data/sharding_db_4.db`（SQLite、usersテーブル32分割、postsテーブル32分割）
- **管理ツール**:
  - GoAdmin: カスタム処理用の管理画面（http://localhost:8081/admin）
  - データベース管理ツール: 現時点では専用ツールなし
- **マイグレーション管理**: Atlas CLIを使用（Issue #26で導入済み）
- **環境別設定**: `config/{env}/database.yaml`（develop, staging, production）

### 2.2 課題点
1. **データベース操作の不便さ**: SQLiteデータベースファイルを直接操作するための専用ツールが存在しない
2. **SQL実行の手間**: データベースに対してSQLを実行する際、コマンドラインや専用クライアントが必要
3. **スクリプト管理の不足**: よく使うSQLスクリプトを保存・管理する仕組みがない
4. **視覚的なデータ確認の困難**: テーブル構造やデータを視覚的に確認するツールがない

### 2.3 本実装による改善点
1. **Webベースのデータベース管理**: ブラウザからデータベースにアクセス可能
2. **視覚的なデータ操作**: テーブル構造の確認、データの閲覧・編集が容易
3. **SQLスクリプトの管理**: Resource Manager機能により、よく使うSQLスクリプトを保存・管理可能
4. **Git管理によるスクリプトのバージョン管理**: Resource Managerに保存したスクリプトをプロジェクトディレクトリにマウントし、Gitで管理可能

## 3. 機能要件

### 3.1 CloudBeaverのDockerセットアップ

#### 3.1.1 Docker Composeファイルの作成
- **ファイル**: `docker-compose.yml`（プロジェクトルート）
- **内容**: 
  - CloudBeaverの公式Dockerイメージを使用
  - ポート設定: デフォルトポート8978をホストにマッピング
  - データベースファイルへのアクセス: `server/data/`ディレクトリをマウント
  - Resource Manager用ディレクトリのマウント: `cloudbeaver/scripts/`ディレクトリをマウント
  - 環境変数の設定（必要に応じて）

#### 3.1.2 Dockerイメージの選択
- **公式イメージ**: `dbeaver/cloudbeaver:latest` または最新の安定版
- **公式サイト**: https://cloudbeaver.io/
- **ドキュメント**: 公式Dockerドキュメントに従う

### 3.2 起動コマンドの定義

#### 3.2.1 package.jsonの作成
- **ファイル**: `package.json`（プロジェクトルート）
- **内容**: 
  - CloudBeaver起動用のnpmスクリプトを定義
  - 停止用スクリプトも定義
  - その他の管理用スクリプト（必要に応じて）

#### 3.2.2 npmスクリプトの定義
- **起動コマンド**: `npm run cloudbeaver:start`
  - Docker Composeを使用してCloudBeaverを起動
  - 環境変数`APP_ENV`で環境を指定（develop/staging/production）
  - デフォルトは`develop`環境
- **停止コマンド**: `npm run cloudbeaver:stop`
  - Docker Composeを使用してCloudBeaverを停止
- **ログ確認コマンド**: `npm run cloudbeaver:logs`（オプション）
  - CloudBeaverのログを確認

#### 3.2.3 環境別制御
- **環境変数**: `APP_ENV`環境変数で環境を切り替え
  - `APP_ENV=develop`: 開発環境
  - `APP_ENV=staging`: ステージング環境
  - `APP_ENV=production`: 本番環境
  - デフォルト値: `develop`（環境変数が未設定の場合）
- **環境別設定の適用**:
  - 環境変数`APP_ENV`の値に基づいて、適切な設定を適用
  - データベースファイルのパス、ポート番号などが環境によって異なる場合は、環境変数やDocker Composeの設定で制御
  - 既存の`config/{env}/database.yaml`の設定を参考に、環境別のデータベースパスを決定

### 3.3 Resource Manager用ディレクトリマウント

#### 3.3.1 ディレクトリ構造の作成
- **ディレクトリ**: `cloudbeaver/scripts/`
- **目的**: Resource Managerに保存したスクリプトをこのディレクトリに保存
- **Git管理**: このディレクトリをGitで管理し、スクリプトのバージョン管理を実現

#### 3.3.2 Dockerボリュームマウント
- **マウント設定**: `cloudbeaver/scripts/`ディレクトリをCloudBeaverコンテナ内の適切なパスにマウント
- **CloudBeaverの設定**: Resource Managerの保存先をマウントしたディレクトリに設定
- **注意**: CloudBeaverの公式ドキュメントに従って、Resource Managerの保存先を設定

### 3.4 データベース接続設定

#### 3.4.1 マスターデータベースへの接続
- **データベースファイル**: `server/data/master.db`
- **接続タイプ**: SQLite
- **接続名**: `master` または `Master Database`
- **設定方法**: CloudBeaverのWeb UIから手動で接続設定を行う
- **マウント設定**: `server/data/`ディレクトリをCloudBeaverコンテナ内にマウント

#### 3.4.2 シャーディングデータベースへの接続
- **データベースファイル**: 
  - `server/data/sharding_db_1.db`
  - `server/data/sharding_db_2.db`
  - `server/data/sharding_db_3.db`
  - `server/data/sharding_db_4.db`
- **接続タイプ**: SQLite
- **接続名**: `sharding_db_1`, `sharding_db_2`, `sharding_db_3`, `sharding_db_4` または `Sharding DB 1` ～ `Sharding DB 4`
- **設定方法**: CloudBeaverのWeb UIから手動で各データベースの接続設定を行う
- **マウント設定**: `server/data/`ディレクトリをCloudBeaverコンテナ内にマウント

#### 3.4.3 接続設定の管理
- **初回設定**: CloudBeaver起動後、Web UIから手動で接続設定を行う
- **接続設定の保存**: CloudBeaverの設定ファイルに接続情報が保存される
- **設定ファイルの場所**: `cloudbeaver/config/{env}/`ディレクトリ（環境別）
  - `APP_ENV=develop` → `cloudbeaver/config/develop/`
  - `APP_ENV=staging` → `cloudbeaver/config/staging/`
  - `APP_ENV=production` → `cloudbeaver/config/production/`
- **Git管理**: 設定ファイルはGitで管理可能
- **環境別分離**: 環境ごとに異なる接続設定を管理可能

## 4. 非機能要件

### 4.1 Docker環境の前提条件
- DockerおよびDocker Composeがインストールされていること
- Dockerが正常に動作していること
- ポート8978が使用可能であること（他のサービスと競合しないこと）

### 4.2 ポート競合の回避
- デフォルトポート8978を使用
- ポートが競合する場合は、docker-compose.ymlでポート番号を変更可能にする
- ポート番号の変更方法をドキュメントに記載

### 4.3 データベースファイルへのアクセス権限
- DockerコンテナからSQLiteデータベースファイルにアクセス可能であること
- ファイルの読み書き権限が適切に設定されていること
- マウントしたディレクトリの権限設定を確認

### 4.4 Resource ManagerスクリプトのGit管理
- `cloudbeaver/scripts/`ディレクトリをGitで管理
- `.gitignore`で不要なファイルを除外（必要に応じて）
- スクリプトファイルの命名規則を定義（オプション）

### 4.5 セキュリティ
- CloudBeaverは開発環境での使用を想定
- 本番環境での使用は想定しない（本番環境では適切なアクセス制御が必要）
- 認証設定は必要に応じて検討（CloudBeaverのデフォルト設定を確認）

## 5. 制約事項

### 5.1 既存システムとの関係
- **GoAdminとの共存**: GoAdminはカスタム処理用として維持し、CloudBeaverはデータ操作用として使用
- **データベース構造の変更なし**: 既存のデータベース構造は変更しない
- **Atlasとの関係**: Atlasはマイグレーション管理用として維持し、CloudBeaverはデータ操作・確認用として使用

### 5.2 環境別の対応
- **環境切り替え**: `APP_ENV`環境変数で環境を切り替え（既存システムと同様）
- **開発環境**: 本実装は開発環境を優先
- **ステージング・本番環境**: 
  - データベースファイルのパスが環境によって異なる場合は、環境変数やDocker Composeの設定で制御
  - 本番環境では適切なアクセス制御が必要（本実装のスコープ外）

### 5.3 技術スタック
- **CloudBeaver**: 最新の安定版を使用
- **Docker**: Docker Composeを使用
- **データベース**: SQLite（開発環境）
- 既存のGoバージョン（1.23.4）は維持

### 5.4 データベースファイルの場所
- データベースファイルは`server/data/`ディレクトリに配置されていることを前提
- データベースファイルの移動は行わない

## 6. 受け入れ基準

### 6.1 Dockerセットアップ
- [ ] `docker-compose.yml`が作成されている
- [ ] CloudBeaverの公式Dockerイメージが使用されている
- [ ] ポート8978がホストにマッピングされている
- [ ] `server/data/`ディレクトリがマウントされている
- [ ] `cloudbeaver/config/{env}/`ディレクトリが環境別にマウントされている
- [ ] `cloudbeaver/scripts/`ディレクトリがマウントされている

### 6.2 起動コマンド
- [ ] プロジェクトルートに`package.json`が作成されている
- [ ] `scripts/cloudbeaver-start.sh`スクリプトが作成されている
- [ ] `scripts/cloudbeaver-start.sh`に実行権限が付与されている
- [ ] `npm run cloudbeaver:start`でCloudBeaverが起動する
- [ ] `APP_ENV=develop npm run cloudbeaver:start`で開発環境として起動する
- [ ] `APP_ENV=staging npm run cloudbeaver:start`でステージング環境として起動する
- [ ] `APP_ENV=production npm run cloudbeaver:start`で本番環境として起動する
- [ ] `npm run cloudbeaver:stop`でCloudBeaverが停止する
- [ ] 起動後、http://localhost:8978 にアクセスできる

### 6.3 データベース接続
- [ ] マスターデータベース（`server/data/master.db`）に接続できる
- [ ] シャーディングデータベース（`server/data/sharding_db_1.db` ～ `sharding_db_4.db`）に接続できる
- [ ] 各データベースのテーブル一覧が表示される
- [ ] 各データベースのデータを閲覧できる
- [ ] SQLクエリを実行できる

### 6.4 CloudBeaver設定ファイル管理
- [ ] `cloudbeaver/config/develop/`ディレクトリが作成されている
- [ ] `cloudbeaver/config/staging/`ディレクトリが作成されている
- [ ] `cloudbeaver/config/production/`ディレクトリが作成されている
- [ ] 接続設定が`cloudbeaver/config/{env}/`ディレクトリに保存される
- [ ] 設定ファイルがGitで管理できる
- [ ] 環境別に設定が分離されている

### 6.5 Resource Manager機能
- [ ] `cloudbeaver/scripts/`ディレクトリが作成されている
- [ ] Resource Managerに保存したスクリプトが`cloudbeaver/scripts/`ディレクトリに保存される
- [ ] 保存されたスクリプトがGitで管理できる
- [ ] スクリプトの作成・編集・削除が正常に動作する

### 6.6 ドキュメント
- [ ] `docs/Database-Viewer.md`が作成されている
- [ ] `docs/Database-Viewer.md`にCloudBeaverの詳細な使用方法が記載されている
- [ ] `README.md`にCloudBeaverの簡単な説明と起動方法が追記されている
- [ ] `README.md`に`docs/Database-Viewer.md`へのリンクが記載されている
- [ ] 環境別の起動方法（`APP_ENV`環境変数の使用方法）が記載されている
- [ ] データベース接続設定の手順が記載されている
- [ ] Resource Managerの使用方法が記載されている
- [ ] トラブルシューティング情報が記載されている

## 7. 影響範囲

### 7.1 新規追加が必要なディレクトリ・ファイル

#### ディレクトリ
- `cloudbeaver/config/`: CloudBeaver設定ファイル（環境別）
  - `cloudbeaver/config/develop/`: 開発環境用設定
  - `cloudbeaver/config/staging/`: ステージング環境用設定
  - `cloudbeaver/config/production/`: 本番環境用設定
- `cloudbeaver/scripts/`: Resource Manager用スクリプト保存ディレクトリ

#### ファイル
- `docker-compose.yml`: CloudBeaver用Docker Compose設定
- `package.json`: npmスクリプト定義（プロジェクトルート）
- `scripts/cloudbeaver-start.sh`: CloudBeaver起動スクリプト

### 7.2 変更が必要なファイル

#### ドキュメント
- `docs/Database-Viewer.md`: CloudBeaver専用の詳細ドキュメント（新規作成）
- `README.md`: CloudBeaverの簡単な説明と起動方法の追記

### 7.3 既存ファイルの扱い
- `server/data/*.db`: 既存のデータベースファイルはそのまま使用
- データベースファイルの移動や変更は行わない

## 8. 実装上の注意事項

### 8.1 Docker Compose設定
- CloudBeaverの公式Dockerイメージを使用
- ポートマッピングを適切に設定
- ボリュームマウントを適切に設定:
  - データベースファイル: `./server/data:/data:ro`
  - CloudBeaver設定ディレクトリ: `./cloudbeaver/config/${APP_ENV:-develop}:/opt/cloudbeaver/workspace`
  - Resource Manager用ディレクトリ: `./cloudbeaver/scripts:/scripts`
- 環境変数の設定:
  - `APP_ENV`: 環境名（develop/staging/production）
  - `CB_WORKSPACE`: CloudBeaverのワークスペースディレクトリ（`/opt/cloudbeaver/workspace`）

### 8.2 データベースファイルのマウント
- `server/data/`ディレクトリをCloudBeaverコンテナ内にマウント
- SQLiteデータベースファイルへのアクセス権限を確認
- マウントパスはCloudBeaverからアクセス可能なパスに設定

### 8.3 Resource Managerの設定
- CloudBeaverのResource Manager機能の保存先を`cloudbeaver/scripts/`ディレクトリに設定
- CloudBeaverの公式ドキュメントに従って設定
- 設定方法が不明な場合は、CloudBeaverの設定ファイルを確認

### 8.4 接続設定の手順
- CloudBeaver起動後、Web UIから手動でデータベース接続を設定
- 接続設定の手順をドキュメントに記載
- 接続設定は初回のみ必要（設定は`cloudbeaver/config/{env}/`ディレクトリに保存される）
- 環境別に設定が分離されるため、各環境で個別に接続設定を行う
- 設定ファイルはGitで管理可能

### 8.5 package.jsonの作成
- プロジェクトルートに`package.json`を作成
- npmスクリプトを定義（`cloudbeaver:start`, `cloudbeaver:stop`など）
- 既存の`client/package.json`とは別に、プロジェクトルート用の`package.json`を作成
- 環境変数`APP_ENV`をnpmスクリプトに渡す方法を実装
  - `scripts/cloudbeaver-start.sh`スクリプトを作成
  - スクリプト内で`APP_ENV`環境変数を設定（未設定の場合は`develop`をデフォルト）
  - スクリプトの実行権限を付与: `chmod +x scripts/cloudbeaver-start.sh`

### 8.7 環境別制御の実装
- **環境変数の取得**: `scripts/cloudbeaver-start.sh`スクリプト内で`APP_ENV`環境変数を取得
- **デフォルト値**: 環境変数が未設定の場合は`develop`をデフォルトとする
- **シェルスクリプトの使用**: `cross-env`の代わりにシェルスクリプトを使用（staging・本番環境でも動作）
- **Docker Composeへの環境変数の渡し方**:
  - `docker-compose.yml`で環境変数を参照可能にする
  - 環境変数をDocker Composeの`environment`セクションで設定
- **環境別設定の適用**:
  - データベースファイルのパスが環境によって異なる場合は、環境変数で制御
  - ポート番号が環境によって異なる場合は、環境変数で制御
  - 既存の`config/{env}/database.yaml`の設定を参考に、環境別のデータベースパスを決定

### 8.6 ドキュメント整備
- `README.md`にCloudBeaverの起動方法を追記
- 環境別の起動方法を記載:
  - `APP_ENV=develop npm run cloudbeaver:start`（開発環境）
  - `APP_ENV=staging npm run cloudbeaver:start`（ステージング環境）
  - `APP_ENV=production npm run cloudbeaver:start`（本番環境）
  - デフォルトは`develop`環境
- 基本的な使用方法（データベース接続、SQL実行、Resource Manager使用）を記載
- トラブルシューティング情報を記載（必要に応じて）

## 9. 参考情報

### 9.1 関連Issue
- GitHub Issue #27: データの管理機能としてのCloudBeaverの導入

### 9.2 既存ドキュメント
- `README.md`: プロジェクト概要とセットアップ手順
- `docs/Sharding.md`: シャーディングの詳細仕様
- `config/{env}/database.yaml`: データベース設定ファイル

### 9.3 技術スタック
- **CloudBeaver**: https://cloudbeaver.io/
- **Docker**: Docker Composeを使用
- **データベース**: SQLite（開発環境）

### 9.4 参考リンク
- CloudBeaver公式サイト: https://cloudbeaver.io/
- CloudBeaver GitHub: https://github.com/dbeaver/cloudbeaver
- CloudBeaver Docker: https://hub.docker.com/r/dbeaver/cloudbeaver
- CloudBeaver ドキュメント: https://cloudbeaver.io/docs/

