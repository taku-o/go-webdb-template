# シャーディング規則修正実装タスク一覧

## 概要
シャーディング規則修正の実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: 設定構造体の拡張

#### - [ ] タスク 1.1: DatabaseGroupsConfig構造体の追加
**目的**: データベースグループ設定の構造体を定義

**作業内容**:
- `server/internal/config/config.go`に以下の構造体を追加:
  ```go
  type DatabaseGroupsConfig struct {
      Master   []ShardConfig        `mapstructure:"master"`
      Sharding ShardingGroupConfig  `mapstructure:"sharding"`
  }
  
  type ShardingGroupConfig struct {
      Databases []ShardConfig       `mapstructure:"databases"`
      Tables    []ShardingTableConfig `mapstructure:"tables"`
  }
  
  type ShardingTableConfig struct {
      Name        string `mapstructure:"name"`
      SuffixCount int    `mapstructure:"suffix_count"`
  }
  ```
- `ShardConfig`構造体に`TableRange [2]int`フィールドを追加

**受け入れ基準**:
- すべての構造体が定義されている
- コンパイルエラーがない
- 構造体のフィールドタグが正しく設定されている

---

#### - [ ] タスク 1.2: DatabaseConfig構造体の拡張
**目的**: DatabaseConfigにgroupsフィールドを追加

**作業内容**:
- `server/internal/config/config.go`の`DatabaseConfig`構造体に`Groups DatabaseGroupsConfig`フィールドを追加
- 既存の`Shards []ShardConfig`フィールドは後方互換性のために残す（非推奨コメントを追加）

**受け入れ基準**:
- `DatabaseConfig`構造体に`Groups`フィールドが追加されている
- 既存の`Shards`フィールドが残っている
- コンパイルエラーがない

---

### Phase 2: 設定ファイルの更新

#### - [ ] タスク 2.1: 開発環境設定ファイルの更新
**目的**: develop環境の設定ファイルにgroups構造を追加

**作業内容**:
- `config/develop/database.yaml`を更新
- `groups`構造を追加:
  ```yaml
  database:
    groups:
      master:
        - id: 1
          driver: sqlite3
          dsn: ./data/master.db
          writer_dsn: ./data/master.db
          reader_dsns:
            - ./data/master.db
          reader_policy: random
          max_connections: 10
          max_idle_connections: 5
          connection_max_lifetime: 300s
      sharding:
        databases:
          - id: 1
            driver: sqlite3
            dsn: ./data/sharding_db_1.db
            writer_dsn: ./data/sharding_db_1.db
            reader_dsns:
              - ./data/sharding_db_1.db
            reader_policy: random
            max_connections: 10
            max_idle_connections: 5
            connection_max_lifetime: 300s
            table_range: [0, 7]
          - id: 2
            driver: sqlite3
            dsn: ./data/sharding_db_2.db
            writer_dsn: ./data/sharding_db_2.db
            reader_dsns:
              - ./data/sharding_db_2.db
            reader_policy: random
            max_connections: 10
            max_idle_connections: 5
            connection_max_lifetime: 300s
            table_range: [8, 15]
          - id: 3
            driver: sqlite3
            dsn: ./data/sharding_db_3.db
            writer_dsn: ./data/sharding_db_3.db
            reader_dsns:
              - ./data/sharding_db_3.db
            reader_policy: random
            max_connections: 10
            max_idle_connections: 5
            connection_max_lifetime: 300s
            table_range: [16, 23]
          - id: 4
            driver: sqlite3
            dsn: ./data/sharding_db_4.db
            writer_dsn: ./data/sharding_db_4.db
            reader_dsns:
              - ./data/sharding_db_4.db
            reader_policy: random
            max_connections: 10
            max_idle_connections: 5
            connection_max_lifetime: 300s
            table_range: [24, 31]
        tables:
          - name: users
            suffix_count: 32
          - name: posts
            suffix_count: 32
  ```
- 既存の`shards`配列は後方互換性のために残す（コメントアウトまたは削除）

**受け入れ基準**:
- `config/develop/database.yaml`に`groups`構造が追加されている
- YAML形式が正しい
- すべての設定項目が正しく記述されている

---

#### - [ ] タスク 2.2: ステージング環境設定ファイルの更新
**目的**: staging環境の設定ファイルにgroups構造を追加

**作業内容**:
- `config/staging/database.yaml`を更新
- develop環境と同様の構造で、PostgreSQL用の設定を追加
- 環境変数（`DB_PASSWORD_*`）を使用

**受け入れ基準**:
- `config/staging/database.yaml`に`groups`構造が追加されている
- PostgreSQL用の設定が正しく記述されている
- 環境変数が使用されている

---

#### - [ ] タスク 2.3: 本番環境設定ファイルの更新
**目的**: production環境の設定ファイル例にgroups構造を追加

**作業内容**:
- `config/production/database.yaml.example`を更新
- staging環境と同様の構造で、コメント付きで説明を追加

**受け入れ基準**:
- `config/production/database.yaml.example`に`groups`構造が追加されている
- コメントで説明が記載されている

---

### Phase 3: テーブル選択ロジックの実装

#### - [ ] タスク 3.1: TableSelector構造体の実装
**目的**: テーブル選択ロジックを提供する構造体を実装

**作業内容**:
- `server/internal/db/sharding.go`に`TableSelector`構造体を追加
- 以下のメソッドを実装:
  - `NewTableSelector(tableCount, tablesPerDB int) *TableSelector`
  - `GetTableNumber(id int64) int`
  - `GetTableName(baseName string, id int64) string`
  - `GetDBID(tableNumber int) int`
  - `GetTableCount() int`

**受け入れ基準**:
- `TableSelector`構造体が実装されている
- すべてのメソッドが実装されている
- コンパイルエラーがない

---

#### - [ ] タスク 3.2: テーブル名生成ユーティリティ関数の実装
**目的**: テーブル名生成のユーティリティ関数を実装

**作業内容**:
- `server/internal/db/table_selector.go`を新規作成
- 以下の関数を実装:
  - `GetShardingTableName(baseName string, id int64) string`
  - `GetShardingTableNumber(id int64) int`
  - `GetShardingDBID(tableNumber int) int`
  - `ValidateTableName(tableName string, allowedBaseNames []string) bool`

**受け入れ基準**:
- `table_selector.go`ファイルが作成されている
- すべての関数が実装されている
- SQLインジェクション対策の検証関数が実装されている
- コンパイルエラーがない

---

#### - [ ] タスク 3.3: TableSelectorの単体テスト
**目的**: TableSelectorの動作を検証

**作業内容**:
- `server/internal/db/sharding_test.go`または新規テストファイルにテストを追加
- 以下のテストケースを実装:
  - テーブル番号の計算テスト
  - テーブル名の生成テスト
  - データベースIDの計算テスト
  - エッジケースのテスト（負の値、0、大きな値など）

**受け入れ基準**:
- テストファイルが作成されている
- すべてのテストケースが実装されている
- テストがパスする

---

### Phase 4: グループ別接続管理の実装

#### - [ ] タスク 4.1: MasterManagerの実装
**目的**: masterグループの接続管理を実装

**作業内容**:
- `server/internal/db/group_manager.go`を新規作成
- `MasterManager`構造体を実装:
  - `NewMasterManager(cfg *config.Config) (*MasterManager, error)`
  - `GetConnection() (*GORMConnection, error)`
  - `CloseAll() error`
  - `PingAll() error`

**受け入れ基準**:
- `MasterManager`構造体が実装されている
- すべてのメソッドが実装されている
- コンパイルエラーがない

---

#### - [ ] タスク 4.2: ShardingManagerの実装
**目的**: shardingグループの接続管理を実装

**作業内容**:
- `server/internal/db/group_manager.go`に`ShardingManager`構造体を追加
- 以下のメソッドを実装:
  - `NewShardingManager(cfg *config.Config) (*ShardingManager, error)`
  - `GetConnectionByTableNumber(tableNumber int) (*GORMConnection, error)`
  - `GetAllConnections() []*GORMConnection`
  - `CloseAll() error`
  - `PingAll() error`

**受け入れ基準**:
- `ShardingManager`構造体が実装されている
- すべてのメソッドが実装されている
- テーブル番号からデータベースIDを正しく決定できる
- コンパイルエラーがない

---

#### - [ ] タスク 4.3: GroupManagerの実装
**目的**: master/shardingグループを統合管理するGroupManagerを実装

**作業内容**:
- `server/internal/db/group_manager.go`に`GroupManager`構造体を追加
- 以下のメソッドを実装:
  - `NewGroupManager(cfg *config.Config) (*GroupManager, error)`
  - `GetMasterConnection() (*GORMConnection, error)`
  - `GetShardingConnection(tableNumber int) (*GORMConnection, error)`
  - `GetShardingConnectionByID(id int64, tableName string) (*GORMConnection, error)`
  - `GetShardingManager() *ShardingManager`（内部実装用）
  - `CloseAll() error`
  - `PingAll() error`

**受け入れ基準**:
- `GroupManager`構造体が実装されている
- すべてのメソッドが実装されている
- master/shardingグループの両方にアクセスできる
- コンパイルエラーがない

---

#### - [ ] タスク 4.4: GroupManagerの単体テスト
**目的**: GroupManagerの動作を検証

**作業内容**:
- `server/internal/db/group_manager_test.go`を新規作成
- 以下のテストケースを実装:
  - GroupManagerの初期化テスト
  - master接続の取得テスト
  - sharding接続の取得テスト（テーブル番号から）
  - 接続のクローズテスト
  - Pingテスト

**受け入れ基準**:
- テストファイルが作成されている
- すべてのテストケースが実装されている
- テストがパスする

---

### Phase 5: マイグレーションテンプレートの作成

#### - [ ] タスク 5.1: masterグループのマイグレーションファイル作成
**目的**: masterグループのnewsテーブル定義を作成

**作業内容**:
- `db/migrations/master/`ディレクトリを作成
- `db/migrations/master/001_init.sql`を作成
- newsテーブルの定義を記述:
  ```sql
  CREATE TABLE IF NOT EXISTS news (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      title TEXT NOT NULL,
      content TEXT NOT NULL,
      author_id INTEGER,
      published_at DATETIME,
      created_at DATETIME NOT NULL,
      updated_at DATETIME NOT NULL
  );
  
  CREATE INDEX IF NOT EXISTS idx_news_published_at ON news(published_at);
  CREATE INDEX IF NOT EXISTS idx_news_author_id ON news(author_id);
  ```

**受け入れ基準**:
- `db/migrations/master/001_init.sql`が作成されている
- newsテーブルの定義が正しく記述されている
- インデックスが定義されている

---

#### - [ ] タスク 5.2: usersテーブルのテンプレート作成
**目的**: usersテーブルのマイグレーションテンプレートを作成

**作業内容**:
- `db/migrations/sharding/templates/`ディレクトリを作成
- `db/migrations/sharding/templates/users.sql.template`を作成
- テンプレート内容:
  ```sql
  -- Users テーブル（テーブル名は{TABLE_NAME}に置換される）
  CREATE TABLE IF NOT EXISTS {TABLE_NAME} (
      id INTEGER PRIMARY KEY,
      name TEXT NOT NULL,
      email TEXT NOT NULL UNIQUE,
      created_at DATETIME NOT NULL,
      updated_at DATETIME NOT NULL
  );
  
  -- インデックス
  CREATE INDEX IF NOT EXISTS idx_{TABLE_NAME}_email ON {TABLE_NAME}(email);
  ```

**受け入れ基準**:
- `db/migrations/sharding/templates/users.sql.template`が作成されている
- `{TABLE_NAME}`プレースホルダーが使用されている
- テーブル定義が正しく記述されている

---

#### - [ ] タスク 5.3: postsテーブルのテンプレート作成
**目的**: postsテーブルのマイグレーションテンプレートを作成

**作業内容**:
- `db/migrations/sharding/templates/posts.sql.template`を作成
- テンプレート内容:
  ```sql
  -- Posts テーブル（テーブル名は{TABLE_NAME}に置換される）
  CREATE TABLE IF NOT EXISTS {TABLE_NAME} (
      id INTEGER PRIMARY KEY,
      user_id INTEGER NOT NULL,
      title TEXT NOT NULL,
      content TEXT NOT NULL,
      created_at DATETIME NOT NULL,
      updated_at DATETIME NOT NULL,
      FOREIGN KEY (user_id) REFERENCES users_{TABLE_SUFFIX}(id) ON DELETE CASCADE
  );
  
  -- インデックス
  CREATE INDEX IF NOT EXISTS idx_{TABLE_NAME}_user_id ON {TABLE_NAME}(user_id);
  CREATE INDEX IF NOT EXISTS idx_{TABLE_NAME}_created_at ON {TABLE_NAME}(created_at);
  ```
- 注意: `{TABLE_SUFFIX}`は`000`, `001`, ..., `031`に置換される

**受け入れ基準**:
- `db/migrations/sharding/templates/posts.sql.template`が作成されている
- `{TABLE_NAME}`と`{TABLE_SUFFIX}`プレースホルダーが使用されている
- テーブル定義が正しく記述されている

---

#### - [ ] タスク 5.4: マイグレーション生成ツールの実装
**目的**: テンプレートからマイグレーションファイルを生成するツールを実装

**作業内容**:
- `server/cmd/migrate-sharding/`ディレクトリを作成
- `server/cmd/migrate-sharding/main.go`を作成
- テンプレートファイルを読み込み、32個のテーブル定義を生成する機能を実装
- コマンドライン引数でテンプレートファイルと出力ディレクトリを指定

**受け入れ基準**:
- `server/cmd/migrate-sharding/main.go`が作成されている
- テンプレートから32個のテーブル定義が生成される
- コマンドライン引数が正しく処理される

---

#### - [ ] タスク 5.5: マイグレーション適用スクリプトの作成
**目的**: 生成されたマイグレーションファイルを各データベースに適用するスクリプトを作成

**作業内容**:
- `scripts/apply-sharding-migrations.sh`を作成
- テンプレートからマイグレーションファイルを生成
- 各データベースに適切なテーブルを適用（_000-007はsharding_db_1、など）

**受け入れ基準**:
- `scripts/apply-sharding-migrations.sh`が作成されている
- スクリプトが実行可能である
- 各データベースに適切なテーブルが適用される

---

### Phase 6: Repository層の変更

#### - [ ] タスク 6.1: UserRepositoryの変更（database/sql版）
**目的**: UserRepositoryで動的テーブル名を使用

**作業内容**:
- `server/internal/repository/user_repository.go`を更新
- `NewUserRepository`関数を`GroupManager`を受け取るように変更
- `TableSelector`を追加
- すべてのクエリで動的テーブル名を使用:
  - `Create`: テーブル名を動的生成
  - `GetByID`: テーブル名を動的生成
  - `List`: クロステーブルクエリを実装
  - `Update`: テーブル名を動的生成
  - `Delete`: テーブル名を動的生成

**受け入れ基準**:
- `UserRepository`が`GroupManager`を使用している
- すべてのメソッドで動的テーブル名が使用されている
- コンパイルエラーがない

---

#### - [ ] タスク 6.2: PostRepositoryの変更（database/sql版）
**目的**: PostRepositoryで動的テーブル名を使用

**作業内容**:
- `server/internal/repository/post_repository.go`を更新
- `NewPostRepository`関数を`GroupManager`を受け取るように変更
- `TableSelector`を追加
- すべてのクエリで動的テーブル名を使用

**受け入れ基準**:
- `PostRepository`が`GroupManager`を使用している
- すべてのメソッドで動的テーブル名が使用されている
- コンパイルエラーがない

---

#### - [ ] タスク 6.3: UserRepositoryGORMの変更
**目的**: UserRepositoryGORMで動的テーブル名を使用

**作業内容**:
- `server/internal/repository/user_repository_gorm.go`を更新
- `NewUserRepositoryGORM`関数を`GroupManager`を受け取るように変更
- GORMの`Table()`メソッドで動的テーブル名を指定
- すべてのクエリで動的テーブル名を使用

**受け入れ基準**:
- `UserRepositoryGORM`が`GroupManager`を使用している
- すべてのメソッドで動的テーブル名が使用されている
- コンパイルエラーがない

---

#### - [ ] タスク 6.4: PostRepositoryGORMの変更
**目的**: PostRepositoryGORMで動的テーブル名を使用

**作業内容**:
- `server/internal/repository/post_repository_gorm.go`を更新
- `NewPostRepositoryGORM`関数を`GroupManager`を受け取るように変更
- GORMの`Table()`メソッドで動的テーブル名を指定
- すべてのクエリで動的テーブル名を使用

**受け入れ基準**:
- `PostRepositoryGORM`が`GroupManager`を使用している
- すべてのメソッドで動的テーブル名が使用されている
- コンパイルエラーがない

---

### Phase 7: モデル層の追加

#### - [ ] タスク 7.1: Newsモデルの作成
**目的**: newsテーブル用のモデルを定義

**作業内容**:
- `server/internal/model/news.go`を新規作成
- `News`構造体を定義:
  ```go
  type News struct {
      ID          int64     `json:"id"`
      Title       string    `json:"title"`
      Content     string    `json:"content"`
      AuthorID    *int64    `json:"author_id,omitempty"`
      PublishedAt *time.Time `json:"published_at,omitempty"`
      CreatedAt   time.Time `json:"created_at"`
      UpdatedAt   time.Time `json:"updated_at"`
  }
  ```
- 必要に応じて`CreateNewsRequest`、`UpdateNewsRequest`構造体も定義

**受け入れ基準**:
- `server/internal/model/news.go`が作成されている
- `News`構造体が定義されている
- コンパイルエラーがない

---

### Phase 8: サービス層の更新

#### - [ ] タスク 8.1: main.goの更新
**目的**: main.goでGroupManagerを使用するように変更

**作業内容**:
- `server/cmd/server/main.go`を更新
- `NewGORMManager`の代わりに`NewGroupManager`を使用
- Repositoryの初期化時に`GroupManager`を渡す

**受け入れ基準**:
- `main.go`が`GroupManager`を使用している
- アプリケーションが正常に起動する
- コンパイルエラーがない

---

#### - [ ] タスク 8.2: admin/main.goの更新
**目的**: admin/main.goでGroupManagerを使用するように変更

**作業内容**:
- `server/cmd/admin/main.go`を更新
- `NewGORMManager`の代わりに`NewGroupManager`を使用
- 既存の機能が正常に動作することを確認

**受け入れ基準**:
- `admin/main.go`が`GroupManager`を使用している
- 管理画面が正常に起動する
- コンパイルエラーがない

---

### Phase 8.5: GoAdmin管理画面のnewsデータ参照ページ追加

#### - [ ] タスク 8.5.1: GetNewsTable関数の実装
**目的**: GoAdmin管理画面にnewsテーブルの参照機能を追加

**作業内容**:
- `server/internal/admin/tables.go`に`GetNewsTable`関数を追加
- users/postsテーブルと同様の実装パターンで実装
- 一覧表示、詳細表示、新規作成、編集、削除機能を提供
- フィルタリング（タイトル、公開日時）とソート（ID、タイトル、公開日時、作成日時）機能を実装

**受け入れ基準**:
- `GetNewsTable`関数が実装されている
- すべてのCRUD操作が可能である
- フィルタリングとソート機能が動作する
- コンパイルエラーがない

---

#### - [ ] タスク 8.5.2: テーブルジェネレータへの登録
**目的**: GoAdminのテーブルジェネレータにnewsを登録

**作業内容**:
- `server/internal/admin/tables.go`の`Generators`マップに`news`を追加
- `"news": GetNewsTable`を追加

**受け入れ基準**:
- `Generators`マップに`news`が追加されている
- 管理画面でnewsテーブルが表示される

---

#### - [ ] タスク 8.5.3: データベース接続設定の更新
**目的**: GoAdminのデータベース接続設定でmasterグループのデータベースを使用

**作業内容**:
- `server/internal/admin/config.go`の`getDatabaseConfig`関数を更新
- `c.appConfig.Database.Groups.Master[0].DSN`を使用するように変更
- 既存の`c.appConfig.Database.Shards[0].DSN`の使用は後方互換性のために残す（非推奨）

**受け入れ基準**:
- `getDatabaseConfig`関数がmasterグループのデータベースを使用している
- 管理画面が正常に動作する

---

#### - [ ] タスク 8.5.4: ホームページへの統計情報追加（オプション）
**目的**: ホームページにnewsの統計情報を表示

**作業内容**:
- `server/internal/admin/pages/home.go`の`HomePage`関数を更新
- `getTableCount(conn, "news")`でnewsの件数を取得
- HTMLコンテンツにnews統計情報のボックスを追加

**受け入れ基準**:
- ホームページにnewsの統計情報が表示される（オプション）
- 既存のusers/posts統計情報と同様の表示形式

---

### Phase 9: テストの実装

#### - [ ] タスク 9.1: Repository層の統合テスト更新
**目的**: Repository層のテストを新しいアーキテクチャに対応

**作業内容**:
- `server/internal/repository/*_test.go`を更新
- `GroupManager`を使用するようにテストを変更
- 動的テーブル名を使用したテストを追加
- クロステーブルクエリのテストを追加

**受け入れ基準**:
- すべてのテストが更新されている
- テストがパスする
- 新しいテストケースが追加されている

---

#### - [ ] タスク 9.2: 統合テストの実装
**目的**: 実際のデータベースを使用した統合テストを実装

**作業内容**:
- `server/test/integration/`に新しい統合テストを追加
- masterグループとshardingグループの両方をテスト
- テーブル選択ロジックの統合テスト
- クロステーブルクエリの統合テスト

**受け入れ基準**:
- 統合テストが実装されている
- テストがパスする
- カバレッジが適切である

---

### Phase 10: ドキュメントの更新

#### - [ ] タスク 10.1: Sharding.mdの更新
**目的**: 新しいアーキテクチャをドキュメントに反映

**作業内容**:
- `docs/Sharding.md`を更新
- master/shardingグループの説明を追加
- テーブル選択ルールの説明を追加
- マイグレーション手順を追加
- 設定例を更新

**受け入れ基準**:
- `docs/Sharding.md`が更新されている
- 新しいアーキテクチャが説明されている
- 設定例が正しく記載されている

---

#### - [ ] タスク 10.2: README.mdの更新
**目的**: セットアップ手順を更新

**作業内容**:
- `README.md`を更新
- データベースセットアップ手順を更新
- マイグレーション適用手順を追加
- 新しいディレクトリ構造を反映

**受け入れ基準**:
- `README.md`が更新されている
- セットアップ手順が正しく記載されている
- マイグレーション手順が記載されている

---

## 実装順序の推奨

1. **Phase 1-2**: 設定構造体と設定ファイルの拡張（基盤整備）
2. **Phase 3**: テーブル選択ロジックの実装（独立した機能）
3. **Phase 4**: グループ別接続管理の実装（基盤整備）
4. **Phase 5**: マイグレーションテンプレートの作成（データベース準備）
5. **Phase 6**: Repository層の変更（ビジネスロジック）
6. **Phase 7-8**: モデル層とサービス層の更新（統合）
7. **Phase 8.5**: GoAdmin管理画面のnewsデータ参照ページ追加（管理画面機能）
8. **Phase 9**: テストの実装（品質保証）
9. **Phase 10**: ドキュメントの更新（文書化）

## 注意事項

### 後方互換性
- 既存の`shards`配列は後方互換性のために残す
- 既存の`GORMManager`は非推奨として残す
- 段階的な移行を可能にする

### テスト戦略
- 各フェーズで単体テストを実装
- 統合テストで実際のデータベース操作を検証
- 既存のテストが可能な限り動作することを確認

### エラーハンドリング
- 接続エラー時の適切なエラーメッセージ
- テーブル選択エラー時の検証
- マイグレーションエラー時のログ記録

