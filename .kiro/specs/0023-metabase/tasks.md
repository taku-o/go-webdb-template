# Metabase導入実装タスク一覧

## 概要
Metabase導入の実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: ディレクトリ構造の作成

#### - [ ] タスク 1.1: Metabase設定ディレクトリの作成
**目的**: Metabase設定ファイル用の環境別ディレクトリを作成

**作業内容**:
- `metabase/config/develop/`ディレクトリを作成
- `metabase/config/staging/`ディレクトリを作成
- `metabase/config/production/`ディレクトリを作成
- 各ディレクトリに`.gitkeep`ファイルを追加（空ディレクトリをGitに含める）

**受け入れ基準**:
- `metabase/config/develop/`ディレクトリが存在する
- `metabase/config/staging/`ディレクトリが存在する
- `metabase/config/production/`ディレクトリが存在する
- 各ディレクトリに`.gitkeep`ファイルが存在する

---

### Phase 2: docker-compose.ymlファイルの分離

#### - [ ] タスク 2.1: docker-compose.cloudbeaver.ymlの作成
**目的**: 既存の`docker-compose.yml`をCloudBeaver用にリネームまたは新規作成

**作業内容**:
- 既存の`docker-compose.yml`を`docker-compose.cloudbeaver.yml`にリネーム
- または、`docker-compose.cloudbeaver.yml`を新規作成して既存の内容をコピー
- ファイル内容に変更は不要（既存のCloudBeaver設定を維持）

**受け入れ基準**:
- `docker-compose.cloudbeaver.yml`が存在する
- CloudBeaverサービスの定義が正しく記述されている
- 既存の`docker-compose.yml`が存在しない（リネームした場合）または内容が一致している（新規作成した場合）

---

### Phase 3: Metabase用Docker Compose設定ファイルの作成

#### - [ ] タスク 3.1: docker-compose.metabase.ymlの作成
**目的**: Metabase用のDocker Compose設定ファイルを作成

**作業内容**:
- プロジェクトルートに`docker-compose.metabase.yml`を作成
- Metabaseサービスの定義を追加:
  - イメージ: `metabase/metabase:latest`
  - コンテナ名: `metabase`
  - 再起動ポリシー: `unless-stopped`
  - ポートマッピング: `8970:3000`（固定値）
  - ボリュームマウント:
    - `./server/data:/data:ro` (データベースファイル、読み取り専用)
    - `./metabase/config/${APP_ENV:-develop}:/metabase-data` (設定ディレクトリ、環境別)
  - 環境変数:
    - `APP_ENV=${APP_ENV:-develop}`
    - `MB_DB_FILE=/metabase-data/metabase.db` (固定値として記載)
  - ヘルスチェック設定:
    - test: `["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:3000/api/health"]`
    - interval: `30s`
    - timeout: `10s`
    - retries: `3`
    - start_period: `180s` (Metabaseは起動に時間がかかるため、余裕を持たせて設定)

**受け入れ基準**:
- `docker-compose.metabase.yml`が作成されている
- Metabaseサービスの定義が正しく記述されている
- ポートマッピングが固定値（8970:3000）として設定されている
- 環境別設定ディレクトリのマウントが正しく設定されている
- 環境変数の設定が正しい（APP_ENV、MB_DB_FILE）

---

### Phase 4: 起動スクリプトの作成と更新

#### - [ ] タスク 4.1: metabase-start.shスクリプトの作成
**目的**: Metabase起動用のシェルスクリプトを作成

**作業内容**:
- `scripts/metabase-start.sh`を作成
- スクリプトの内容:
  ```bash
  #!/bin/bash
  # Metabase起動スクリプト
  
  # APP_ENV環境変数が未設定の場合はdevelopをデフォルトとする
  export APP_ENV=${APP_ENV:-develop}
  
  # Docker ComposeでMetabaseを起動
  docker-compose -f docker-compose.metabase.yml up -d
  ```
- 実行権限を付与: `chmod +x scripts/metabase-start.sh`

**受け入れ基準**:
- `scripts/metabase-start.sh`が作成されている
- スクリプトの内容が正しい
- 実行権限が付与されている
- `docker-compose.metabase.yml`を指定している

---

#### - [ ] タスク 4.2: cloudbeaver-start.shスクリプトの更新
**目的**: CloudBeaver起動スクリプトを`docker-compose.cloudbeaver.yml`を使用するように更新

**作業内容**:
- `scripts/cloudbeaver-start.sh`を更新
- `docker-compose up -d`を`docker-compose -f docker-compose.cloudbeaver.yml up -d`に変更

**受け入れ基準**:
- `scripts/cloudbeaver-start.sh`が更新されている
- `docker-compose.cloudbeaver.yml`を指定している
- 既存の機能（APP_ENV環境変数の設定）が維持されている

---

### Phase 5: package.jsonの更新

#### - [ ] タスク 5.1: package.jsonの更新
**目的**: Metabase用のnpmスクリプトを追加し、CloudBeaver用のスクリプトを更新

**作業内容**:
- 既存の`package.json`を更新
- Metabase用のnpmスクリプトを追加:
  - `metabase:start`: `./scripts/metabase-start.sh`
  - `metabase:stop`: `docker-compose -f docker-compose.metabase.yml down`
  - `metabase:logs`: `docker-compose -f docker-compose.metabase.yml logs -f metabase`
  - `metabase:restart`: `npm run metabase:stop && npm run metabase:start`
- CloudBeaver用のnpmスクリプトを更新:
  - `cloudbeaver:start`: 変更なし（`./scripts/cloudbeaver-start.sh`は既に更新済み）
  - `cloudbeaver:stop`: `docker-compose -f docker-compose.cloudbeaver.yml down`に変更
  - `cloudbeaver:logs`: `docker-compose -f docker-compose.cloudbeaver.yml logs -f cloudbeaver`に変更
  - `cloudbeaver:restart`: 変更なし（内部で`cloudbeaver:stop`と`cloudbeaver:start`を使用）

**受け入れ基準**:
- `package.json`が更新されている
- Metabase用のnpmスクリプトが正しく定義されている
- CloudBeaver用のnpmスクリプトが`docker-compose.cloudbeaver.yml`を使用するように更新されている

---

### Phase 6: 動作確認（開発環境）

#### - [ ] タスク 6.1: Metabase起動確認（開発環境）
**目的**: 開発環境でMetabaseが正常に起動することを確認

**作業内容**:
- `npm run metabase:start`を実行
- または `APP_ENV=develop npm run metabase:start`を実行
- Docker ComposeでMetabaseコンテナが起動することを確認
- http://localhost:8970 にアクセスできることを確認
- MetabaseのWeb UIが表示されることを確認
- 初回起動時は管理者アカウント作成画面が表示されることを確認

**受け入れ基準**:
- Metabaseコンテナが正常に起動する
- http://localhost:8970 にアクセスできる
- MetabaseのWeb UIが表示される

---

#### - [ ] タスク 6.2: Metabase停止確認
**目的**: Metabaseが正常に停止することを確認

**作業内容**:
- `npm run metabase:stop`を実行
- Docker ComposeでMetabaseコンテナが停止することを確認

**受け入れ基準**:
- Metabaseコンテナが正常に停止する

---

#### - [ ] タスク 6.3: CloudBeaver起動確認（docker-compose.cloudbeaver.yml使用）
**目的**: 更新後のCloudBeaver起動スクリプトが正常に動作することを確認

**作業内容**:
- `npm run cloudbeaver:start`を実行
- Docker ComposeでCloudBeaverコンテナが起動することを確認（`docker-compose.cloudbeaver.yml`を使用）
- http://localhost:8978 にアクセスできることを確認

**受け入れ基準**:
- CloudBeaverコンテナが正常に起動する
- `docker-compose.cloudbeaver.yml`が使用されている
- http://localhost:8978 にアクセスできる

---

#### - [ ] タスク 6.4: CloudBeaverとMetabaseの個別起動確認
**目的**: CloudBeaverとMetabaseを個別に起動・停止できることを確認

**作業内容**:
- CloudBeaverを起動（`npm run cloudbeaver:start`）
- CloudBeaverが正常に起動することを確認
- CloudBeaverを停止（`npm run cloudbeaver:stop`）
- Metabaseを起動（`npm run metabase:start`）
- Metabaseが正常に起動することを確認
- Metabaseを停止（`npm run metabase:stop`）

**受け入れ基準**:
- CloudBeaverとMetabaseを個別に起動できる
- CloudBeaverとMetabaseを個別に停止できる
- 両方を同時に起動する必要はない（メモリ使用量の制約）

---

### Phase 7: 動作確認（環境別）

#### - [ ] タスク 7.1: Metabase起動確認（ステージング環境）
**目的**: ステージング環境でMetabaseが正常に起動することを確認

**作業内容**:
- `APP_ENV=staging npm run metabase:start`を実行
- Docker ComposeでMetabaseコンテナが起動することを確認
- 環境別設定ディレクトリ（`metabase/config/staging/`）がマウントされていることを確認

**受け入れ基準**:
- Metabaseコンテナが正常に起動する
- ステージング環境用の設定ディレクトリがマウントされている

---

#### - [ ] タスク 7.2: Metabase起動確認（本番環境）
**目的**: 本番環境でMetabaseが正常に起動することを確認

**作業内容**:
- `APP_ENV=production npm run metabase:start`を実行
- Docker ComposeでMetabaseコンテナが起動することを確認
- 環境別設定ディレクトリ（`metabase/config/production/`）がマウントされていることを確認

**受け入れ基準**:
- Metabaseコンテナが正常に起動する
- 本番環境用の設定ディレクトリがマウントされている

---

### Phase 8: データベース接続設定

#### - [ ] タスク 8.1: 管理者アカウントの作成
**目的**: Metabaseの初回起動時に管理者アカウントを作成

**作業内容**:
- Metabase Web UI（http://localhost:8970）にアクセス
- 初回起動時の管理者アカウント作成画面で、管理者情報を入力:
  - 名前
  - メールアドレス
  - パスワード
- 管理者アカウントを作成
- 管理者アカウントでログインできることを確認

**受け入れ基準**:
- 管理者アカウントが作成される
- 管理者アカウントでログインできる

---

#### - [ ] タスク 8.2: マスターデータベース接続設定
**目的**: Metabaseからマスターデータベースに接続できるように設定

**作業内容**:
- Metabase Web UIでデータベース接続を追加
- SQLiteドライバーを選択
- 接続情報を入力:
  - 接続名: `master` または `Master Database`
  - データベースファイル: `/data/master.db`
- 接続をテスト
- 接続を保存
- 接続設定が`metabase/config/develop/`ディレクトリに保存されることを確認

**受け入れ基準**:
- マスターデータベースに接続できる
- 接続設定が`metabase/config/develop/`ディレクトリに保存される
- 接続設定ファイルがGitで管理できる

---

#### - [ ] タスク 8.3: シャーディングデータベース接続設定
**目的**: Metabaseからシャーディングデータベース（4つ）に接続できるように設定

**作業内容**:
- Metabase Web UIで各シャーディングデータベースの接続を追加:
  - `sharding_db_1`: `/data/sharding_db_1.db`
  - `sharding_db_2`: `/data/sharding_db_2.db`
  - `sharding_db_3`: `/data/sharding_db_3.db`
  - `sharding_db_4`: `/data/sharding_db_4.db`
- 各接続をテスト
- 各接続を保存
- 接続設定が`metabase/config/develop/`ディレクトリに保存されることを確認

**受け入れ基準**:
- 4つのシャーディングデータベースに接続できる
- 接続設定が`metabase/config/develop/`ディレクトリに保存される
- 接続設定ファイルがGitで管理できる

---

#### - [ ] タスク 8.4: データベース接続動作確認
**目的**: 接続したデータベースでテーブル一覧の表示、データの閲覧、クエリ作成ができることを確認

**作業内容**:
- マスターデータベースに接続
- テーブル一覧が表示されることを確認
- テーブルのデータを閲覧できることを確認
- クエリを作成・実行できることを確認
- シャーディングデータベース（4つ）でも同様の確認を行う

**受け入れ基準**:
- 各データベースのテーブル一覧が表示される
- 各データベースのデータを閲覧できる
- クエリを作成・実行できる

---

#### - [ ] タスク 8.5: ダッシュボード作成確認
**目的**: ダッシュボードを作成できることを確認

**作業内容**:
- Metabase Web UIでダッシュボードを作成
- クエリ結果をダッシュボードに追加
- ダッシュボードが保存されることを確認
- ダッシュボード設定が`metabase/config/develop/`ディレクトリに保存されることを確認

**受け入れ基準**:
- ダッシュボードを作成できる
- ダッシュボード設定が`metabase/config/develop/`ディレクトリに保存される
- ダッシュボード設定ファイルがGitで管理できる

---

### Phase 9: 環境別設定の確認

#### - [ ] タスク 9.1: 環境別設定ディレクトリの動作確認
**目的**: 環境別に設定ディレクトリが正しくマウントされ、設定が分離されることを確認

**作業内容**:
- 開発環境でMetabaseを起動し、接続設定を行う
- 接続設定が`metabase/config/develop/`ディレクトリに保存されることを確認
- ステージング環境でMetabaseを起動
- ステージング環境用の設定ディレクトリ（`metabase/config/staging/`）がマウントされていることを確認
- 本番環境でMetabaseを起動
- 本番環境用の設定ディレクトリ（`metabase/config/production/`）がマウントされていることを確認

**受け入れ基準**:
- 環境別に設定ディレクトリが正しくマウントされる
- 環境別に設定が分離される
- 各環境の設定ファイルがGitで管理できる

---

### Phase 10: ドキュメント整備

#### - [ ] タスク 10.1: Metabase.mdの作成
**目的**: Metabase専用の詳細ドキュメントを作成

**作業内容**:
- `docs/Metabase.md`を作成
- 以下の内容を含める:
  - 概要（主な機能、役割分担）
  - 前提条件
  - 起動方法（基本的な起動、環境別の起動）
  - 停止方法
  - その他のコマンド（ログ確認、再起動）
  - データベース接続設定（初回設定、マスターデータベース、シャーディングデータベース）
  - クエリ作成方法
  - ダッシュボード作成方法
  - CloudBeaverとMetabaseの使い分け
  - トラブルシューティング情報

**受け入れ基準**:
- `docs/Metabase.md`が作成されている
- 上記の内容が全て記載されている
- 既存の`docs/Database-Viewer.md`と整合性が取れている

---

#### - [ ] タスク 10.2: README.mdの更新
**目的**: README.mdにMetabaseの簡単な説明と起動方法を追記

**作業内容**:
- `README.md`を確認
- Metabaseの簡単な説明を追記（データ可視化・分析ツールとして）
- Metabaseの起動方法を追記:
  - 基本的な起動方法: `npm run metabase:start`
  - 環境別の起動方法: `APP_ENV=develop npm run metabase:start`など
- `docs/Metabase.md`へのリンクを追加
- CloudBeaverとMetabaseの使い分けを追記

**受け入れ基準**:
- `README.md`が更新されている
- Metabaseの説明が追記されている
- Metabaseの起動方法が記載されている
- `docs/Metabase.md`へのリンクが追加されている
- CloudBeaverとMetabaseの使い分けが記載されている

---

## 実装順序の推奨

1. **Phase 1-5**: インフラストラクチャの構築（ディレクトリ作成、Docker Compose設定、スクリプト作成）
2. **Phase 6-7**: 動作確認（開発環境、環境別）
3. **Phase 8**: データベース接続設定（手動作業）
4. **Phase 9**: 環境別設定の確認
5. **Phase 10**: ドキュメント整備

## 注意事項

- **Phase 8のデータベース接続設定**: 手動でMetabase Web UIから行う必要がある
- **環境別設定**: 各環境（develop/staging/production）で個別に接続設定を行う必要がある
- **メモリ使用量**: CloudBeaverとMetabaseは同時に起動しない（メモリ使用量の制約）
- **ポート番号**: ポート8970は固定値として使用。ポート競合の回避は運用者の責任
