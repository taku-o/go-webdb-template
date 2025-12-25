# シャーディング規則修正要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #19
- **Issueタイトル**: シャーディング規則の修正
- **Feature名**: 0012-sharding
- **作成日**: 2025-01-27

### 1.2 目的
データベースとテーブルのシャーディング規則を大きく変更し、masterグループとshardingグループの2つのデータベースグループを導入する。
これにより、より柔軟なデータ分散戦略を実現し、スケーラビリティとパフォーマンスを向上させる。

### 1.3 スコープ
- データベースグループの概念導入（master/sharding）
- masterグループ: 1つのデータベース、通常のテーブル構造（newsテーブルを追加）
- shardingグループ: 4つのデータベース、32分割されたテーブル（_000-031のsuffix）
- 既存のusers/postsテーブルをshardingグループに移行（テーブル名にsuffixを付ける）
- テーブル選択ルールの実装: 格納先テーブル番号=ID(mod全テーブル数)
- テンプレートベースのマイグレーション管理システム
- 設定ファイルの拡張（データベースグループ対応）
- データベース接続管理の拡張（グループ別接続管理）
- Repository層の変更（グループ別テーブル選択）
- GoAdmin管理画面にnewsデータ参照ページの追加

## 2. 背景・現状分析

### 2.1 現在の実装
- **シャーディング戦略**: Hash-based sharding（`HashBasedSharding`）
- **シャード数**: 4つ（shard1, shard2, shard3, shard4）
- **シャードキー**: `user_id`
- **テーブル構造**: 各シャードに同じテーブル構造（users, posts）
- **設定ファイル**: `config/{env}/database.yaml`にシャード設定を定義
- **マイグレーションファイル**: `db/migrations/shard{N}/001_init.sql`で各シャードのスキーマを定義
- **データベース接続管理**: `server/internal/db/manager.go`で`Manager`/`GORMManager`が接続を管理
- **Repository層**: `server/internal/repository/`でシャード選択とクエリ実行

### 2.2 課題点
1. **柔軟性の不足**: すべてのテーブルが同じシャーディング戦略を強制される
2. **テーブル分割の制約**: テーブル単位での分割ができない（シャード単位のみ）
3. **データ分散の制約**: より細かい粒度でのデータ分散ができない
4. **マイグレーション管理の複雑さ**: 各シャードごとに個別のマイグレーションファイルが必要

### 2.3 本実装による改善点
1. **柔軟なデータ分散**: masterグループとshardingグループで異なる戦略を適用可能
2. **細かい粒度の分割**: テーブル単位で32分割することで、より均等なデータ分散を実現
3. **マイグレーション管理の簡素化**: テンプレートベースのマイグレーションで管理を簡素化
4. **スケーラビリティの向上**: テーブル単位での分割により、より細かい粒度でのスケーリングが可能

## 3. 機能要件

### 3.1 データベースグループの概念

#### 3.1.1 グループ定義
データベースを2つのグループに分類する：

1. **masterグループ**
   - データベース数: 1つ
   - テーブル構造: 通常のテーブル名（suffixなし）
   - 用途: シャーディング不要なデータ（newsテーブルなど）
   - テーブル例: `news`

2. **shardingグループ**
   - データベース数: 4つ
   - テーブル構造: テーブル名にsuffix（_000-031）を付ける
   - 用途: シャーディングが必要なデータ（users, postsなど）
   - テーブル例: `users_000`, `users_001`, ..., `users_031`, `posts_000`, `posts_001`, ..., `posts_031`

#### 3.1.2 データベースとテーブルの分散
shardingグループの4つのデータベースに、32個のテーブルを分散：

- **sharding_db_1**: テーブル _000 〜 _007（8テーブル）
- **sharding_db_2**: テーブル _008 〜 _015（8テーブル）
- **sharding_db_3**: テーブル _016 〜 _023（8テーブル）
- **sharding_db_4**: テーブル _024 〜 _031（8テーブル）

### 3.2 テーブル選択ルール

#### 3.2.1 ルール定義
shardingグループのテーブルは、以下のルールでデータを格納する：

```
格納先テーブル番号 = ID % 全テーブル数
```

- **ID**: レコードのID（usersテーブルの場合は`user_id`、postsテーブルの場合は`post_id`）
- **全テーブル数**: 32（_000から_031まで）
- **テーブル番号**: 0-31（_000=0, _001=1, ..., _031=31）

#### 3.2.2 実装例
```go
// テーブル番号の計算
tableNumber := id % 32  // 0-31の範囲

// テーブル名の生成
tableName := fmt.Sprintf("users_%03d", tableNumber)  // users_000, users_001, ...

// データベースの選択
dbNumber := tableNumber / 8  // 0-3の範囲（各DBに8テーブル）
dbID := dbNumber + 1  // 1-4の範囲
```

### 3.3 既存テーブルの移行

#### 3.3.1 usersテーブル
- **移行先**: shardingグループ
- **テーブル名**: `users_000`から`users_031`まで32個
- **選択ルール**: `user_id % 32`でテーブル番号を決定

#### 3.3.2 postsテーブル
- **移行先**: shardingグループ
- **テーブル名**: `posts_000`から`posts_031`まで32個
- **選択ルール**: `post_id % 32`でテーブル番号を決定
- **注意**: `user_id`ではなく`post_id`でテーブルを選択（Issue要件に基づく）

### 3.4 masterグループのテーブル

#### 3.4.1 newsテーブル
- **データベース**: masterグループのデータベース
- **テーブル名**: `news`（suffixなし）
- **用途**: ニュース記事などのシャーディング不要なデータ
- **スキーマ**: 新規作成（要件定義に基づく）

#### 3.4.2 管理画面でのnewsデータ参照
- **GoAdmin管理画面**: newsテーブルのデータを参照するページを追加
- **実装内容**:
  - `server/internal/admin/tables.go`に`GetNewsTable`関数を追加
  - GoAdminのテーブルジェネレータに`news`を登録
  - 一覧表示、詳細表示、新規作成、編集、削除機能を提供
  - ホームページ（`server/internal/admin/pages/home.go`）にnewsの統計情報を追加（オプション）
- **表示項目**: ID、タイトル、内容、作成者ID、公開日時、作成日時、更新日時
- **フィルタリング**: タイトル、公開日時でフィルタリング可能
- **ソート**: ID、タイトル、公開日時、作成日時でソート可能

### 3.5 設定ファイルの拡張

#### 3.5.1 データベースグループ設定
設定ファイル（`config/{env}/database.yaml`）にグループ情報を追加：

```yaml
database:
  groups:
    master:
      - id: 1
        driver: sqlite3
        dsn: ./data/master.db
        # ... 既存の設定項目
    sharding:
      - id: 1
        driver: sqlite3
        dsn: ./data/sharding_db_1.db
        table_range: [0, 7]  # _000-007
        # ... 既存の設定項目
      - id: 2
        driver: sqlite3
        dsn: ./data/sharding_db_2.db
        table_range: [8, 15]  # _008-015
        # ... 既存の設定項目
      - id: 3
        driver: sqlite3
        dsn: ./data/sharding_db_3.db
        table_range: [16, 23]  # _016-023
        # ... 既存の設定項目
      - id: 4
        driver: sqlite3
        dsn: ./data/sharding_db_4.db
        table_range: [24, 31]  # _024-031
        # ... 既存の設定項目
```

#### 3.5.2 テーブル定義設定
shardingグループのテーブル定義を追加：

```yaml
database:
  groups:
    sharding:
      tables:
        - name: users
          suffix_count: 32  # _000-031
        - name: posts
          suffix_count: 32  # _000-031
```

### 3.6 マイグレーション管理

#### 3.6.1 テンプレートベースのマイグレーション
1つのテンプレートファイルから、32個のテーブルを生成するマイグレーションシステムを実装：

- **テンプレートファイル**: `db/migrations/sharding/templates/users.sql.template`
- **展開処理**: テンプレート内の`{TABLE_NAME}`を`users_000`, `users_001`, ..., `users_031`に置換
- **適用先**: 各データベースに適切なテーブルを生成（_000-007はsharding_db_1、など）

#### 3.6.2 マイグレーションファイル構造
```
db/migrations/
├── master/
│   └── 001_init.sql          # newsテーブル
└── sharding/
    ├── templates/
    │   ├── users.sql.template    # usersテーブルのテンプレート
    │   └── posts.sql.template     # postsテーブルのテンプレート
    └── 001_init_users.sql        # 生成されたマイグレーション（32テーブル分）
```

#### 3.6.3 マイグレーション適用ツール
マイグレーションツールまたはCLIコマンドを実装：

- テンプレートファイルを読み込み
- 32個のテーブル定義を生成
- 各データベースに適切なテーブルを適用

### 3.7 データベース接続管理の拡張

#### 3.7.1 グループ別接続管理
`Manager`/`GORMManager`を拡張して、グループ別の接続管理を実装：

- **MasterManager**: masterグループの接続を管理
- **ShardingManager**: shardingグループの接続を管理
- **統合Manager**: 両方のグループを統合管理

#### 3.7.2 接続取得メソッド
グループとテーブル番号に基づいて接続を取得：

```go
// masterグループの接続取得
GetMasterConnection() (*Connection, error)

// shardingグループの接続取得（テーブル番号から自動的にDBを選択）
GetShardingConnection(tableNumber int) (*Connection, error)

// テーブル名から接続取得
GetConnectionByTableName(tableName string) (*Connection, error)
```

### 3.8 Repository層の変更

#### 3.8.1 テーブル名の動的生成
Repository層で、IDに基づいてテーブル名を動的に生成：

```go
// テーブル名の生成
func getTableName(baseName string, id int64) string {
    tableNumber := int(id % 32)
    return fmt.Sprintf("%s_%03d", baseName, tableNumber)
}

// データベースの選択
func getShardingDBID(tableNumber int) int {
    return (tableNumber / 8) + 1  // 1-4の範囲
}
```

#### 3.8.2 クエリの変更
既存のクエリを、動的なテーブル名を使用するように変更：

```go
// 変更前
query := "SELECT * FROM users WHERE id = ?"

// 変更後
tableName := getTableName("users", userID)
query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", tableName)
```

#### 3.8.3 クロステーブルクエリ
全32テーブルからデータを取得する場合の処理：

- 各データベースから並列にクエリを実行
- 結果をマージして返す

### 3.9 既存データの移行

#### 3.9.1 データ移行の方針
- **既存データの移行は不要**: Issue要件に基づき、既存データの移行は考えない
- **データ損失を許容**: 既存データは消失しても構わない
- **新規データのみ対応**: 新規に作成されるデータのみ、新しいルールに従う

## 4. 非機能要件

### 4.1 パフォーマンス
- テーブル選択の計算はO(1)で実行されること
- クロステーブルクエリは並列実行されること
- 接続プールはグループ別に管理されること

### 4.2 拡張性
- テーブル数の変更（32以外）に対応できる設計とすること
- データベース数の変更（4以外）に対応できる設計とすること
- 新しいテーブルの追加が容易であること

### 4.3 保守性
- テンプレートベースのマイグレーションで、スキーマ変更が容易であること
- 設定ファイルで柔軟に構成変更できること
- コードの可読性とテスト容易性を維持すること

### 4.4 後方互換性
- 既存のAPIエンドポイントの動作は維持されること
- 既存のテストコードは可能な限り動作すること（大幅な変更が必要な場合はテストコードの更新も許容）

## 5. 制約事項

### 5.1 技術的制約
- 既存のGORM v1.25.12を使用すること
- 既存のデータベースドライバ（sqlite3, postgres）をサポートすること
- 既存の設定ファイル構造（YAML）を基本とすること

### 5.2 プロジェクト制約
- 既存のレイヤードアーキテクチャを維持すること
- 既存のテストフレームワークを使用すること
- 既存のドキュメント構造を維持すること

### 5.3 データ移行
- **既存データの移行は不要**: Issue要件に基づき、既存データの移行は行わない
- **データ損失を許容**: 既存データは消失しても構わない
- **新規データのみ対応**: 新規に作成されるデータのみ、新しいルールに従う

## 6. 受け入れ基準

### 6.1 データベースグループ
- [ ] masterグループのデータベースが1つ作成されている
- [ ] shardingグループのデータベースが4つ作成されている
- [ ] 設定ファイルにグループ情報が定義されている

### 6.2 テーブル構造
- [ ] masterグループにnewsテーブルが作成されている
- [ ] shardingグループにusers_000からusers_031まで32個のテーブルが作成されている
- [ ] shardingグループにposts_000からposts_031まで32個のテーブルが作成されている
- [ ] 各データベースに適切なテーブルが分散されている（_000-007はsharding_db_1、など）

### 6.10 管理画面
- [ ] GoAdmin管理画面にnewsテーブルの参照ページが追加されている
- [ ] newsテーブルの一覧表示、詳細表示、新規作成、編集、削除が可能である
- [ ] フィルタリングとソート機能が動作する
- [ ] ホームページにnewsの統計情報が表示される（オプション）

### 6.3 テーブル選択ルール
- [ ] IDに基づいて正しいテーブル番号が計算される
- [ ] テーブル番号から正しいデータベースが選択される
- [ ] テーブル名が正しく生成される（users_000, users_001, ...）

### 6.4 マイグレーション管理
- [ ] テンプレートベースのマイグレーションシステムが実装されている
- [ ] テンプレートから32個のテーブル定義が生成される
- [ ] 各データベースに適切なテーブルが適用される

### 6.5 データベース接続管理
- [ ] グループ別の接続管理が実装されている
- [ ] masterグループの接続が取得できる
- [ ] shardingグループの接続がテーブル番号から取得できる

### 6.6 Repository層
- [ ] Repository層で動的なテーブル名が使用されている
- [ ] 既存のCRUD操作が正常に動作する
- [ ] クロステーブルクエリが正常に動作する

### 6.7 設定ファイル
- [ ] すべての環境（develop/staging/production）で設定ファイルが更新されている
- [ ] 設定ファイルの構造が一貫している

### 6.8 テスト
- [ ] 単体テストが実装されている
- [ ] 統合テストが実装されている
- [ ] 既存のテストが可能な限り動作する（大幅な変更が必要な場合はテストコードの更新も許容）

### 6.9 ドキュメント
- [ ] `docs/Sharding.md`が更新されている
- [ ] 新しいアーキテクチャが文書化されている
- [ ] マイグレーション手順が文書化されている

## 7. 影響範囲

### 7.1 新規追加が必要なディレクトリ・ファイル

#### ディレクトリ
- `db/migrations/master/`: masterグループのマイグレーションファイル
- `db/migrations/sharding/templates/`: shardingグループのテンプレートファイル

#### ファイル
- `db/migrations/master/001_init.sql`: newsテーブルの定義
- `db/migrations/sharding/templates/users.sql.template`: usersテーブルのテンプレート
- `db/migrations/sharding/templates/posts.sql.template`: postsテーブルのテンプレート
- `server/internal/db/group_manager.go`: グループ別接続管理（新規または既存ファイルの拡張）

### 7.2 変更が必要なファイル

#### 設定ファイル
- `config/develop/database.yaml`: グループ情報を追加
- `config/staging/database.yaml`: グループ情報を追加
- `config/production/database.yaml.example`: グループ情報を追加

#### データベース接続管理
- `server/internal/config/config.go`: グループ設定の構造体を追加
- `server/internal/db/manager.go`: グループ別接続管理を追加
- `server/internal/db/connection.go`: グループ対応の接続処理を追加

#### Repository層
- `server/internal/repository/user_repository.go`: 動的テーブル名を使用
- `server/internal/repository/post_repository.go`: 動的テーブル名を使用
- `server/internal/repository/user_repository_gorm.go`: 動的テーブル名を使用
- `server/internal/repository/post_repository_gorm.go`: 動的テーブル名を使用

#### モデル層
- `server/internal/model/news.go`: newsモデルを追加（新規）

#### サービス層
- `server/internal/service/news_service.go`: newsサービスを追加（新規、オプション）

#### 管理画面
- `server/internal/admin/tables.go`: GetNewsTable関数を追加
- `server/internal/admin/pages/home.go`: news統計情報を追加（オプション）

#### テストファイル
- `server/internal/db/group_manager_test.go`: グループ管理のテスト（新規）
- `server/internal/repository/*_test.go`: 既存テストの更新

### 7.3 削除されるファイル
なし（既存ファイルは拡張または変更のみ）

### 7.4 ドキュメント
- `docs/Sharding.md`: 新しいアーキテクチャを反映
- `README.md`: セットアップ手順を更新

## 8. 実装上の注意事項

### 8.1 テーブル名の動的生成
- SQLインジェクション対策として、テーブル名は必ずホワイトリストで検証すること
- テーブル名の生成ロジックは一箇所に集約し、テスト容易性を確保すること

### 8.2 マイグレーション管理
- テンプレートファイルの構文エラーを検出できること
- マイグレーション適用時のエラーハンドリングを適切に実装すること
- ロールバック機能を考慮すること（将来の拡張）

### 8.3 データベース接続管理
- 接続プールはグループ別に管理し、リソースリークを防ぐこと
- 接続エラー時の適切なエラーハンドリングを実装すること

### 8.4 パフォーマンス
- テーブル選択の計算は軽量であること（O(1)）
- クロステーブルクエリは並列実行し、パフォーマンスを維持すること

### 8.5 テスト
- 単体テストでテーブル選択ロジックを検証すること
- 統合テストで実際のデータベース操作を検証すること
- 既存のテストが可能な限り動作することを確認すること

## 9. 参考情報

### 9.1 関連Issue
- GitHub Issue #19: シャーディング規則の修正

### 9.2 既存ドキュメント
- `docs/Sharding.md`: 既存のシャーディング戦略の詳細
- `docs/Architecture.md`: アーキテクチャの詳細
- `config/develop/database.yaml`: 開発環境設定ファイル

### 9.3 既存実装
- `server/internal/db/sharding.go`: 既存のシャーディング戦略
- `server/internal/db/manager.go`: 既存の接続管理
- `server/internal/repository/user_repository.go`: 既存のRepository実装

### 9.4 技術スタック
- **Go**: 1.21+
- **GORM**: v1.25.12
- **データベース**: SQLite3（開発環境）、PostgreSQL（本番環境）
- **設定管理**: viper（spf13/viper）

