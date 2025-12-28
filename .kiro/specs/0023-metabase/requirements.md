# Metabase導入要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #47
- **Issueタイトル**: Metabase の導入
- **Feature名**: 0023-metabase
- **作成日**: 2025-01-27

### 1.2 目的
非エンジニア向けのデータビューワアプリとしてMetabaseを導入し、Webベースのデータ可視化・分析ツールを提供する。
本プロジェクトでは管理系アプリを3つ用意する方針であり、GoAdminはカスタム処理用、CloudBeaverはデータ操作用、Metabaseはデータ可視化・分析用として役割を分担する。

### 1.3 スコープ
- MetabaseをDockerで動作させる
- CloudBeaverとMetabaseは片方ずつしか起動しない運用を実現（メモリ使用量の制約）
- 起動用のdocker-compose.ymlファイルをCloudBeaver用とMetabase用に分ける
- 起動コマンドをpackage.jsonに定義する
- Metabase上で設定した内容を可能であればGitに保存する
- 環境別設定（develop/staging/production）をサポートする

**本実装の範囲外**:
- Metabaseの詳細な設定やカスタマイズ（基本的な接続設定のみ）
- データベース構造の変更
- 既存のCloudBeaverやGoAdminの機能変更

## 2. 背景・現状分析

### 2.1 現在の実装
- **データベース構成**:
  - マスターデータベース: `server/data/master.db`（SQLite、newsテーブル、GoAdmin関連テーブル）
  - シャーディングデータベース: `server/data/sharding_db_1.db` ～ `server/data/sharding_db_4.db`（SQLite、usersテーブル32分割、postsテーブル32分割）
- **管理ツール**:
  - GoAdmin: カスタム処理用の管理画面（http://localhost:8081/admin）
  - CloudBeaver: データ操作用のWebベースツール（http://localhost:8978、Issue #27で導入済み）
  - データ可視化・分析ツール: 現時点では専用ツールなし
- **CloudBeaverの実装**:
  - `docker-compose.yml`でCloudBeaverサービスを定義
  - `scripts/cloudbeaver-start.sh`で起動スクリプトを提供
  - `cloudbeaver/config/${APP_ENV}/`で環境別設定を管理
  - `cloudbeaver/scripts/`でResource ManagerのスクリプトをGit管理
- **環境別設定**: `config/{env}/database.yaml`（develop, staging, production）

### 2.2 課題点
1. **データ可視化・分析ツールの不足**: 非エンジニアがデータを視覚的に分析するためのツールが存在しない
2. **ダッシュボード機能の不足**: データをダッシュボード形式で表示・共有する機能がない
3. **メモリ使用量の制約**: Metabaseはメモリをかなり使うため、CloudBeaverと同時に起動できない
4. **docker-compose.ymlの統合**: 現在はCloudBeaverのみがdocker-compose.ymlに定義されており、Metabaseを追加する際に運用上の制約がある

### 2.3 本実装による改善点
1. **非エンジニア向けのデータ可視化**: Metabaseにより、非エンジニアでもデータを視覚的に分析可能
2. **ダッシュボード機能**: データをダッシュボード形式で表示・共有可能
3. **運用の柔軟性**: CloudBeaverとMetabaseを個別に起動・停止可能
4. **設定のGit管理**: Metabaseの設定をGitで管理し、バージョン管理を実現

## 3. 機能要件

### 3.1 MetabaseのDockerセットアップ

#### 3.1.1 Docker Composeファイルの分離
- **CloudBeaver用**: `docker-compose.cloudbeaver.yml`（既存の`docker-compose.yml`をリネームまたは新規作成）
- **Metabase用**: `docker-compose.metabase.yml`（新規作成）
- **目的**: CloudBeaverとMetabaseを個別に起動・停止可能にする
- **運用**: 開発環境ではCloudBeaverとMetabaseは片方ずつしか起動しない

#### 3.1.2 Metabase Docker Composeファイルの作成
- **ファイル**: `docker-compose.metabase.yml`（プロジェクトルート）
- **内容**: 
  - Metabaseの公式Dockerイメージを使用
  - ポート設定: デフォルトポート8970をホストにマッピング
  - データベースファイルへのアクセス: `server/data/`ディレクトリをマウント
  - Metabase設定ディレクトリのマウント: `metabase/config/${APP_ENV}/`ディレクトリをマウント
  - 環境変数の設定（必要に応じて）

#### 3.1.3 Dockerイメージの選択
- **公式イメージ**: `metabase/metabase:latest` または最新の安定版
- **公式サイト**: https://www.metabase.com/
- **ドキュメント**: 公式Dockerドキュメントに従う

### 3.2 起動コマンドの定義

#### 3.2.1 package.jsonの更新
- **ファイル**: `package.json`（プロジェクトルート、既存ファイルを更新）
- **内容**: 
  - Metabase起動用のnpmスクリプトを追加
  - CloudBeaver起動用のnpmスクリプトを更新（docker-compose.cloudbeaver.ymlを使用）
  - 停止用スクリプトも更新
  - その他の管理用スクリプト（必要に応じて）

#### 3.2.2 npmスクリプトの定義
- **Metabase起動コマンド**: `npm run metabase:start`
  - Docker Composeを使用してMetabaseを起動
  - 環境変数`APP_ENV`で環境を指定（develop/staging/production）
  - デフォルトは`develop`環境
  - `docker-compose.metabase.yml`を使用
- **Metabase停止コマンド**: `npm run metabase:stop`
  - Docker Composeを使用してMetabaseを停止
  - `docker-compose.metabase.yml`を使用
- **Metabaseログ確認コマンド**: `npm run metabase:logs`（オプション）
  - Metabaseのログを確認
- **CloudBeaver起動コマンド**: `npm run cloudbeaver:start`（既存を更新）
  - `docker-compose.cloudbeaver.yml`を使用するように更新
- **CloudBeaver停止コマンド**: `npm run cloudbeaver:stop`（既存を更新）
  - `docker-compose.cloudbeaver.yml`を使用するように更新

#### 3.2.3 環境別制御
- **環境変数**: `APP_ENV`環境変数で環境を切り替え
  - `APP_ENV=develop`: 開発環境
  - `APP_ENV=staging`: ステージング環境
  - `APP_ENV=production`: 本番環境
  - デフォルト値: `develop`（環境変数が未設定の場合）
- **環境別設定の適用**:
  - 環境変数`APP_ENV`の値に基づいて、適切な設定を適用
  - データベースファイルのパスが環境によって異なる場合は、環境変数やDocker Composeの設定で制御
  - ポート番号は固定値（8970）を使用。ポート競合の回避は運用者の責任
  - 既存の`config/{env}/database.yaml`の設定を参考に、環境別のデータベースパスを決定

### 3.3 Metabase設定ディレクトリの管理

#### 3.3.1 ディレクトリ構造の作成
- **ディレクトリ**: `metabase/config/`
- **環境別サブディレクトリ**:
  - `metabase/config/develop/`: 開発環境用設定
  - `metabase/config/staging/`: ステージング環境用設定
  - `metabase/config/production/`: 本番環境用設定
- **目的**: Metabaseの設定ファイルを環境別に管理
- **Git管理**: このディレクトリをGitで管理し、設定のバージョン管理を実現

#### 3.3.2 Dockerボリュームマウント
- **マウント設定**: `metabase/config/${APP_ENV}/`ディレクトリをMetabaseコンテナ内の適切なパスにマウント
- **Metabaseの設定**: Metabaseの設定ファイルの保存先をマウントしたディレクトリに設定
- **注意**: Metabaseの公式ドキュメントに従って、設定ファイルの保存先を設定

### 3.4 データベース接続設定

#### 3.4.1 マスターデータベースへの接続
- **データベースファイル**: `server/data/master.db`
- **接続タイプ**: SQLite
- **接続名**: `master` または `Master Database`
- **設定方法**: MetabaseのWeb UIから手動で接続設定を行う
- **マウント設定**: `server/data/`ディレクトリをMetabaseコンテナ内にマウント

#### 3.4.2 シャーディングデータベースへの接続
- **データベースファイル**: 
  - `server/data/sharding_db_1.db`
  - `server/data/sharding_db_2.db`
  - `server/data/sharding_db_3.db`
  - `server/data/sharding_db_4.db`
- **接続タイプ**: SQLite
- **接続名**: `sharding_db_1`, `sharding_db_2`, `sharding_db_3`, `sharding_db_4` または `Sharding DB 1` ～ `Sharding DB 4`
- **設定方法**: MetabaseのWeb UIから手動で各データベースの接続設定を行う
- **マウント設定**: `server/data/`ディレクトリをMetabaseコンテナ内にマウント

#### 3.4.3 接続設定の管理
- **初回設定**: Metabase起動後、Web UIから手動で接続設定を行う
- **接続設定の保存**: Metabaseの設定ファイルに接続情報が保存される
- **設定ファイルの場所**: `metabase/config/{env}/`ディレクトリ（環境別）
  - `APP_ENV=develop` → `metabase/config/develop/`
  - `APP_ENV=staging` → `metabase/config/staging/`
  - `APP_ENV=production` → `metabase/config/production/`
- **Git管理**: 設定ファイルはGitで管理可能
- **環境別分離**: 環境ごとに異なる接続設定を管理可能

## 4. 非機能要件

### 4.1 Docker環境の前提条件
- DockerおよびDocker Composeがインストールされていること
- Dockerが正常に動作していること
- ポート8970が使用可能であること（他のサービスと競合しないこと）
- メモリが十分にあること（Metabaseはメモリをかなり使う）

### 4.2 ポート競合の回避
- デフォルトポート8970を使用
- ポートが競合する場合は、docker-compose.metabase.ymlでポート番号を変更可能にする
- ポート番号の変更方法をドキュメントに記載

### 4.3 データベースファイルへのアクセス権限
- DockerコンテナからSQLiteデータベースファイルにアクセス可能であること
- ファイルの読み書き権限が適切に設定されていること
- マウントしたディレクトリの権限設定を確認

### 4.4 Metabase設定ファイルのGit管理
- `metabase/config/`ディレクトリをGitで管理
- `.gitignore`で不要なファイルを除外（必要に応じて）
- 設定ファイルの命名規則を定義（オプション）

### 4.5 セキュリティ
- Metabaseは開発環境での使用を想定
- 本番環境での使用は想定しない（本番環境では適切なアクセス制御が必要）
- 認証設定は必要に応じて検討（Metabaseのデフォルト設定を確認）

### 4.6 メモリ使用量の制約
- **開発環境**: CloudBeaverとMetabaseは片方ずつしか起動しない
- **Staging/Production環境**: 動作させるサーバーを別にする運用
- **運用方針**: 開発者は必要に応じてCloudBeaverまたはMetabaseを選択して起動

## 5. 制約事項

### 5.1 既存システムとの関係
- **GoAdminとの共存**: GoAdminはカスタム処理用として維持し、Metabaseはデータ可視化・分析用として使用
- **CloudBeaverとの共存**: CloudBeaverはデータ操作用として維持し、Metabaseはデータ可視化・分析用として使用
- **データベース構造の変更なし**: 既存のデータベース構造は変更しない
- **Atlasとの関係**: Atlasはマイグレーション管理用として維持し、Metabaseはデータ可視化・分析用として使用

### 5.2 環境別の対応
- **環境切り替え**: `APP_ENV`環境変数で環境を切り替え（既存システムと同様）
- **開発環境**: 本実装は開発環境を優先
- **ステージング・本番環境**: 
  - データベースファイルのパスが環境によって異なる場合は、環境変数やDocker Composeの設定で制御
  - 本番環境では適切なアクセス制御が必要（本実装のスコープ外）
  - Staging/Production環境では動作させるサーバーを別にする運用

### 5.3 技術スタック
- **Metabase**: 最新の安定版を使用
- **Docker**: Docker Composeを使用
- **データベース**: SQLite（開発環境）
- 既存のGoバージョン（1.23.4）は維持

### 5.4 データベースファイルの場所
- データベースファイルは`server/data/`ディレクトリに配置されていることを前提
- データベースファイルの移動は行わない

### 5.5 運用上の制約
- **同時起動の制限**: 開発環境ではCloudBeaverとMetabaseは片方ずつしか起動しない
- **メモリ使用量**: Metabaseはメモリをかなり使うため、同時起動を避ける
- **docker-compose.ymlの分離**: CloudBeaver用とMetabase用でdocker-compose.ymlファイルを分ける

## 6. 受け入れ基準

### 6.1 Dockerセットアップ
- [ ] `docker-compose.metabase.yml`が作成されている
- [ ] `docker-compose.cloudbeaver.yml`が作成されている（既存の`docker-compose.yml`をリネームまたは新規作成）
- [ ] Metabaseの公式Dockerイメージが使用されている
- [ ] ポート8970がホストにマッピングされている
- [ ] `server/data/`ディレクトリがマウントされている
- [ ] `metabase/config/${APP_ENV}/`ディレクトリが環境別にマウントされている

### 6.2 起動コマンド
- [ ] `scripts/metabase-start.sh`スクリプトが作成されている
- [ ] `scripts/metabase-start.sh`に実行権限が付与されている
- [ ] `package.json`に`metabase:start`スクリプトが追加されている
- [ ] `package.json`の`cloudbeaver:start`スクリプトが`docker-compose.cloudbeaver.yml`を使用するように更新されている
- [ ] `npm run metabase:start`でMetabaseが起動する
- [ ] `APP_ENV=develop npm run metabase:start`で開発環境として起動する
- [ ] `APP_ENV=staging npm run metabase:start`でステージング環境として起動する
- [ ] `APP_ENV=production npm run metabase:start`で本番環境として起動する
- [ ] `npm run metabase:stop`でMetabaseが停止する
- [ ] 起動後、http://localhost:8970 にアクセスできる
- [ ] CloudBeaverとMetabaseを個別に起動・停止できる

### 6.3 データベース接続
- [ ] マスターデータベース（`server/data/master.db`）に接続できる
- [ ] シャーディングデータベース（`server/data/sharding_db_1.db` ～ `sharding_db_4.db`）に接続できる
- [ ] 各データベースのテーブル一覧が表示される
- [ ] 各データベースのデータを閲覧できる
- [ ] クエリを作成・実行できる
- [ ] ダッシュボードを作成できる

### 6.4 Metabase設定ファイル管理
- [ ] `metabase/config/develop/`ディレクトリが作成されている
- [ ] `metabase/config/staging/`ディレクトリが作成されている
- [ ] `metabase/config/production/`ディレクトリが作成されている
- [ ] 接続設定が`metabase/config/{env}/`ディレクトリに保存される
- [ ] 設定ファイルがGitで管理できる
- [ ] 環境別に設定が分離されている

### 6.5 docker-compose.ymlファイルの分離確認
- [ ] docker-compose.ymlファイルがCloudBeaver用とMetabase用に分離されている

### 6.6 ドキュメント
- [ ] `docs/Metabase.md`が作成されている
- [ ] `docs/Metabase.md`にMetabaseの詳細な使用方法が記載されている
- [ ] `README.md`にMetabaseの簡単な説明と起動方法が追記されている
- [ ] `README.md`に`docs/Metabase.md`へのリンクが記載されている
- [ ] 環境別の起動方法（`APP_ENV`環境変数の使用方法）が記載されている
- [ ] データベース接続設定の手順が記載されている
- [ ] CloudBeaverとMetabaseの使い分けが記載されている
- [ ] トラブルシューティング情報が記載されている

## 7. 影響範囲

### 7.1 新規追加が必要なディレクトリ・ファイル

#### ディレクトリ
- `metabase/config/`: Metabase設定ファイル（環境別）
  - `metabase/config/develop/`: 開発環境用設定
  - `metabase/config/staging/`: ステージング環境用設定
  - `metabase/config/production/`: 本番環境用設定

#### ファイル
- `docker-compose.metabase.yml`: Metabase用Docker Compose設定
- `docker-compose.cloudbeaver.yml`: CloudBeaver用Docker Compose設定（既存の`docker-compose.yml`をリネームまたは新規作成）
- `scripts/metabase-start.sh`: Metabase起動スクリプト

### 7.2 変更が必要なファイル

#### 設定ファイル
- `package.json`: Metabase用のnpmスクリプトを追加、CloudBeaver用のスクリプトを更新

#### ドキュメント
- `docs/Metabase.md`: Metabase専用の詳細ドキュメント（新規作成）
- `README.md`: Metabaseの簡単な説明と起動方法の追記

### 7.3 既存ファイルの扱い
- `server/data/*.db`: 既存のデータベースファイルはそのまま使用
- データベースファイルの移動や変更は行わない
- `docker-compose.yml`: CloudBeaver用の`docker-compose.cloudbeaver.yml`にリネームまたは新規作成

## 8. 実装上の注意事項

### 8.1 Docker Compose設定
- Metabaseの公式Dockerイメージを使用
- ポートマッピングを適切に設定（デフォルト: 8970:3000）
- ボリュームマウントを適切に設定:
  - データベースファイル: `./server/data:/data:ro`
  - Metabase設定ディレクトリ: `./metabase/config/${APP_ENV:-develop}:/metabase-data`
- 環境変数の設定:
  - `APP_ENV`: 環境名（develop/staging/production）
  - Metabase固有の環境変数（必要に応じて）
- ヘルスチェック設定（オプション）

### 8.2 データベースファイルのマウント
- `server/data/`ディレクトリをMetabaseコンテナ内にマウント
- SQLiteデータベースファイルへのアクセス権限を確認
- マウントパスはMetabaseからアクセス可能なパスに設定

### 8.3 Metabase設定ファイルの管理
- Metabaseの設定ファイルの保存先を`metabase/config/{env}/`ディレクトリに設定
- Metabaseの公式ドキュメントに従って設定
- 設定方法が不明な場合は、Metabaseの設定ファイルを確認

### 8.4 接続設定の手順
- Metabase起動後、Web UIから手動でデータベース接続を設定
- 接続設定の手順をドキュメントに記載
- 接続設定は初回のみ必要（設定は`metabase/config/{env}/`ディレクトリに保存される）
- 環境別に設定が分離されるため、各環境で個別に接続設定を行う
- 設定ファイルはGitで管理可能

### 8.5 package.jsonの更新
- 既存の`package.json`を更新
- npmスクリプトを追加（`metabase:start`, `metabase:stop`, `metabase:logs`など）
- CloudBeaver用のスクリプトを更新（`docker-compose.cloudbeaver.yml`を使用）
- 環境変数`APP_ENV`をnpmスクリプトに渡す方法を実装
  - `scripts/metabase-start.sh`スクリプトを作成
  - スクリプト内で`APP_ENV`環境変数を設定（未設定の場合は`develop`をデフォルト）
  - スクリプトの実行権限を付与: `chmod +x scripts/metabase-start.sh`

### 8.6 docker-compose.ymlファイルの分離
- 既存の`docker-compose.yml`を`docker-compose.cloudbeaver.yml`にリネームまたは新規作成
- `docker-compose.metabase.yml`を新規作成
- 各docker-compose.ymlファイルで適切なサービス名を設定
- 起動スクリプトで適切なdocker-compose.ymlファイルを指定

### 8.7 環境別制御の実装
- **環境変数の取得**: `scripts/metabase-start.sh`スクリプト内で`APP_ENV`環境変数を取得
- **デフォルト値**: 環境変数が未設定の場合は`develop`をデフォルトとする
- **シェルスクリプトの使用**: `cross-env`の代わりにシェルスクリプトを使用（staging・本番環境でも動作）
- **Docker Composeへの環境変数の渡し方**:
  - `docker-compose.metabase.yml`で環境変数を参照可能にする
  - 環境変数をDocker Composeの`environment`セクションで設定
- **環境別設定の適用**:
  - データベースファイルのパスが環境によって異なる場合は、環境変数で制御
  - ポート番号は固定値（8970）を使用。ポート競合の回避は運用者の責任
  - 既存の`config/{env}/database.yaml`の設定を参考に、環境別のデータベースパスを決定

### 8.8 ドキュメント整備
- `README.md`にMetabaseの起動方法を追記
- 環境別の起動方法を記載:
  - `APP_ENV=develop npm run metabase:start`（開発環境）
  - `APP_ENV=staging npm run metabase:start`（ステージング環境）
  - `APP_ENV=production npm run metabase:start`（本番環境）
  - デフォルトは`develop`環境
- 基本的な使用方法（データベース接続、クエリ作成、ダッシュボード作成）を記載
- CloudBeaverとMetabaseの使い分けを記載
- トラブルシューティング情報を記載（必要に応じて）

## 9. 参考情報

### 9.1 関連Issue
- GitHub Issue #47: Metabase の導入
- GitHub Issue #27: データの管理機能としてのCloudBeaverの導入（0015-cloudbeaver）

### 9.2 既存ドキュメント
- `README.md`: プロジェクト概要とセットアップ手順
- `docs/Database-Viewer.md`: CloudBeaverの詳細仕様
- `docs/Sharding.md`: シャーディングの詳細仕様
- `config/{env}/database.yaml`: データベース設定ファイル
- `.kiro/specs/0015-cloudbeaver/requirements.md`: CloudBeaver導入の要件定義書

### 9.3 技術スタック
- **Metabase**: https://www.metabase.com/
- **Docker**: Docker Composeを使用
- **データベース**: SQLite（開発環境）

### 9.4 参考リンク
- Metabase公式サイト: https://www.metabase.com/
- Metabase GitHub: https://github.com/metabase/metabase
- Metabase Docker: https://hub.docker.com/r/metabase/metabase
- Metabase ドキュメント: https://www.metabase.com/docs/
- Metabase Docker ドキュメント: https://www.metabase.com/docs/latest/installation-and-operation/running-metabase-on-docker
