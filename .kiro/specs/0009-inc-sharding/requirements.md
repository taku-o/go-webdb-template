# シャーディング数増加要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #17
- **Issueタイトル**: データベースのシャーディング数を増やす
- **Feature名**: 0009-inc-sharding
- **作成日**: 2025-01-27

### 1.2 目的
サンプルプロジェクトとはいえ、現在のシャーディング数が2つしかないのは少なすぎるため、シャーディング数を4つに増やす。
これにより、より実践的なシャーディング構成を提供し、スケーラビリティの理解を深める。

### 1.3 スコープ
- 開発環境（develop）の設定ファイルにshard3とshard4を追加
- ステージング環境（staging）の設定ファイルにshard3とshard4を追加
- 本番環境（production）の設定ファイル例にshard3とshard4を追加
- マイグレーションファイル（`db/migrations/shard3/`と`db/migrations/shard4/`）の作成
- ドキュメント（`docs/Sharding.md`）の更新
- 既存のshard1とshard2の設定・マイグレーションは変更しない

## 2. 背景・現状分析

### 2.1 現在の実装
- **シャーディング数**: 2つ（shard1, shard2）
- **シャーディング戦略**: Hash-based sharding（`HashBasedSharding`）
- **シャードキー**: `user_id`
- **シャード数決定方法**: 設定ファイルの`shards`配列の長さから自動的に決定される
- **設定ファイル**: `config/{env}/database.yaml`に定義
- **マイグレーションファイル**: 
  - `db/migrations/shard1/001_init.sql`: 初期化スクリプト（users, postsテーブル）
  - `db/migrations/shard1/002_goadmin.sql`: GoAdminフレームワーク用テーブル
  - `db/migrations/shard1/003_menu.sql`: アプリケーション用メニュー
  - `db/migrations/shard2/001_init.sql`: 初期化スクリプト（users, postsテーブル）

### 2.2 シャーディングロジック
現在の実装では、`HashBasedSharding`が以下のように動作する：

```go
func (h *HashBasedSharding) GetShardID(key int64) int {
    hash := fnv.New32a()
    hash.Write([]byte(fmt.Sprintf("%d", key)))
    hashValue := hash.Sum32()
    shardID := int(hashValue%uint32(h.shardCount)) + 1
    return shardID
}
```

- `shardCount`は設定ファイルの`shards`配列の長さから自動的に決定される
- シャードIDは1からN（Nはシャード数）の範囲で返される
- 同じ`user_id`は常に同じシャードにマッピングされる

### 2.3 設定ファイル構造
各環境の設定ファイル（`config/{env}/database.yaml`）には、以下のような構造でシャードが定義されている：

```yaml
database:
  shards:
    - id: 1
      driver: sqlite3  # または postgres
      dsn: ./data/shard1.db
      writer_dsn: ./data/shard1.db
      reader_dsns:
        - ./data/shard1.db
      reader_policy: random
      max_connections: 10
      max_idle_connections: 5
      connection_max_lifetime: 300s
    - id: 2
      # 同様の構造
```

### 2.4 課題点
1. **シャーディング数が少ない**: 2つのシャードでは実践的な構成として不十分
2. **マイグレーションファイルの不足**: shard3とshard4のマイグレーションファイルが存在しない
3. **設定ファイルの不足**: すべての環境でshard3とshard4の設定が不足している

### 2.5 本実装による改善点
1. **実践的な構成**: 4つのシャードにより、より実践的なシャーディング構成を提供
2. **スケーラビリティの理解**: 複数のシャードでの動作を確認できる
3. **データ分散の改善**: より多くのシャードにより、データ分散がより均等になる可能性

### 2.6 データ移行に関する方針
- **データ移行は不要**: 既存データの移行は行わない
- **データ損失を許容**: シャーディング数を増やすことで、既存データのシャード割り当てが変わる可能性があるが、データ損失を許容する
- **データリセットを許容**: 必要に応じて、既存のデータベースをリセットしても良い
- **サンプルプロジェクト**: サンプルプロジェクトであるため、既存データの保持は必須ではない

## 3. 機能要件

### 3.1 設定ファイル更新

#### 3.1.1 開発環境（develop）
- **ファイル**: `config/develop/database.yaml`
- **要件**: shard3とshard4を追加
- **設定内容**:
  - `id: 3`, `id: 4`
  - `driver: sqlite3`
  - `dsn: ./data/shard3.db`, `./data/shard4.db`
  - `writer_dsn`と`reader_dsns`も同様に設定
  - 既存のshard1とshard2の設定は変更しない

#### 3.1.2 ステージング環境（staging）
- **ファイル**: `config/staging/database.yaml`
- **要件**: shard3とshard4を追加
- **設定内容**:
  - `id: 3`, `id: 4`
  - `driver: postgres`
  - ホスト名、ポート、データベース名、ユーザー名を設定
  - `password: ${DB_PASSWORD_SHARD3}`, `${DB_PASSWORD_SHARD4}`（環境変数から読み込み）
  - `writer_dsn`と`reader_dsns`も同様に設定
  - 既存のshard1とshard2の設定は変更しない

#### 3.1.3 本番環境（production）
- **ファイル**: `config/production/database.yaml.example`
- **要件**: shard3とshard4を追加
- **設定内容**:
  - `id: 3`, `id: 4`
  - `driver: postgres`
  - ホスト名、ポート、データベース名、ユーザー名を設定
  - `password: ${DB_PASSWORD_SHARD3}`, `${DB_PASSWORD_SHARD4}`（環境変数から読み込み）
  - `writer_dsn`と`reader_dsns`も同様に設定（複数のReaderも設定可能）
  - 既存のshard1とshard2の設定は変更しない

### 3.2 マイグレーションファイル作成

#### 3.2.1 shard3のマイグレーションファイル
- **ディレクトリ**: `db/migrations/shard3/`
- **ファイル**: `001_init.sql`
- **内容**: shard1の`001_init.sql`と同じスキーマ
  - `users`テーブル
  - `posts`テーブル
  - インデックス（`idx_users_email`, `idx_posts_user_id`, `idx_posts_created_at`）
- **注意**: コメント内の「Shard 1」を「Shard 3」に変更

#### 3.2.2 shard4のマイグレーションファイル
- **ディレクトリ**: `db/migrations/shard4/`
- **ファイル**: `001_init.sql`
- **内容**: shard1の`001_init.sql`と同じスキーマ
  - `users`テーブル
  - `posts`テーブル
  - インデックス（`idx_users_email`, `idx_posts_user_id`, `idx_posts_created_at`）
- **注意**: コメント内の「Shard 1」を「Shard 4」に変更

#### 3.2.3 GoAdmin関連マイグレーション
- **現状**: shard1には`002_goadmin.sql`と`003_menu.sql`が存在するが、shard2には存在しない
- **方針**: 本実装では、shard3とshard4にもGoAdmin関連のマイグレーションファイルは作成しない
- **理由**: GoAdminは管理画面用であり、すべてのシャードに必要ではない可能性があるため
- **将来の拡張**: 必要に応じて後で追加可能

### 3.3 既存機能の維持
- **既存のシャーディングロジック**: 変更不要（動的にシャード数を検出するため）
- **既存のデータベース接続処理**: 変更不要（設定ファイルから自動的に読み込まれる）
- **既存のテスト**: 4シャードでも動作することを確認するが、テストコードの変更は不要

## 4. 非機能要件

### 4.1 後方互換性
- 既存のshard1とshard2の設定・マイグレーションは変更しない
- 既存のデータベース接続処理は変更不要
- 既存のシャーディングロジックは変更不要（動的にシャード数を検出するため）

### 4.2 設定の一貫性
- すべての環境（develop/staging/production）で同じ構造でshard3とshard4を追加
- 既存のshard1とshard2の設定パターンに従う

### 4.3 マイグレーションファイルの一貫性
- shard3とshard4のマイグレーションファイルは、shard1とshard2と同じスキーマを使用
- スキーマの一貫性を保つ

### 4.4 ドキュメントの更新
- `docs/Sharding.md`を更新して、4シャード構成を反映
- 設定例を更新

## 5. 制約事項

### 5.1 技術的制約
- 既存のシャーディングロジック（`HashBasedSharding`）は変更しない
- 既存のデータベース接続処理（`GORMConnection`、`Connection`）は変更しない
- 設定ファイルの構造は既存の形式に従う

### 5.2 プロジェクト制約
- 既存のshard1とshard2の設定・マイグレーションは変更しない
- 既存のテストコードは変更しない（ただし、4シャードでも動作することを確認）

### 5.3 データ移行
- **データ移行は不要**: 本実装では、既存データの移行は行わない
- **データ損失を許容**: シャーディング数を増やすことで、既存データのシャード割り当てが変わる可能性があるが、データ損失を許容する
- **データリセットを許容**: 必要に応じて、既存のデータベースをリセットしても良い
- **新しいデータ**: 新しいデータは4つのシャードに分散される（ハッシュ値に基づいて）

### 5.4 ハッシュベースシャーディングの制約
- Hash-based shardingでは、シャード数を変更すると既存データのシャード割り当てが変わる
- **データ損失を許容**: 既存データの再配置は行わず、データ損失を許容する
- **データリセットを許容**: 必要に応じて、既存のデータベースをリセットしても良い

## 6. 受け入れ基準

### 6.1 設定ファイル
- [ ] `config/develop/database.yaml`にshard3とshard4が追加されている
- [ ] `config/staging/database.yaml`にshard3とshard4が追加されている
- [ ] `config/production/database.yaml.example`にshard3とshard4が追加されている
- [ ] すべての環境で既存のshard1とshard2の設定が変更されていない
- [ ] すべての環境でshard3とshard4の設定が既存のパターンに従っている

### 6.2 マイグレーションファイル
- [ ] `db/migrations/shard3/001_init.sql`が作成されている
- [ ] `db/migrations/shard4/001_init.sql`が作成されている
- [ ] shard3とshard4のマイグレーションファイルがshard1とshard2と同じスキーマを使用している
- [ ] マイグレーションファイルのコメントが適切に更新されている

### 6.3 アプリケーション動作
- [ ] アプリケーションが4つのシャードに正常に接続できる
- [ ] 既存のシャーディングロジックが4シャードでも正常に動作する
- [ ] 新しいデータが4つのシャードに適切に分散される
- [ ] 既存のテストが4シャードでも正常に動作する
- [ ] 既存データの移行は行わず、データ損失を許容する（必要に応じてデータリセットも許容）

### 6.4 ドキュメント
- [ ] `docs/Sharding.md`が更新されている
- [ ] 設定例が4シャード構成を反映している
- [ ] シャーディング数の説明が更新されている

## 7. 影響範囲

### 7.1 新規追加が必要なディレクトリ・ファイル

#### ディレクトリ
- `db/migrations/shard3/`: shard3のマイグレーションファイル用ディレクトリ
- `db/migrations/shard4/`: shard4のマイグレーションファイル用ディレクトリ

#### ファイル
- `db/migrations/shard3/001_init.sql`: shard3の初期化スクリプト
- `db/migrations/shard4/001_init.sql`: shard4の初期化スクリプト

### 7.2 変更が必要なファイル

#### 設定ファイル
- `config/develop/database.yaml`: shard3とshard4を追加
- `config/staging/database.yaml`: shard3とshard4を追加
- `config/production/database.yaml.example`: shard3とshard4を追加

#### ドキュメント
- `docs/Sharding.md`: 4シャード構成を反映した更新

### 7.3 変更不要なファイル
- `server/internal/db/sharding.go`: 既存のロジックは動的にシャード数を検出するため変更不要
- `server/internal/config/config.go`: 既存の構造体定義は変更不要
- `server/internal/db/connection.go`: 既存の接続処理は変更不要
- 既存のテストコード: 変更不要（ただし、4シャードでも動作することを確認）

### 7.4 削除されるファイル
なし

## 8. 実装上の注意事項

### 8.1 設定ファイルの更新
- 既存のshard1とshard2の設定を変更しない
- shard3とshard4の設定は既存のパターンに従う
- 環境変数（`DB_PASSWORD_SHARD3`、`DB_PASSWORD_SHARD4`）の使用を考慮

### 8.2 マイグレーションファイルの作成
- shard1の`001_init.sql`をベースに作成
- コメント内の「Shard 1」を適切に変更（「Shard 3」「Shard 4」）
- スキーマの一貫性を保つ

### 8.3 データベース接続
- 既存の接続処理は変更不要（設定ファイルから自動的に読み込まれる）
- 4つのシャードすべてに接続できることを確認
- **データ損失を許容**: 既存データの移行は行わず、データ損失を許容する
- **データリセットを許容**: 必要に応じて、既存のデータベースをリセットしても良い

### 8.4 シャーディングロジック
- 既存の`HashBasedSharding`は変更不要（動的にシャード数を検出するため）
- 4つのシャードでも正常に動作することを確認

### 8.5 テスト
- 既存のテストコードは変更不要
- 4シャードでも動作することを確認
- 新しいデータが4つのシャードに適切に分散されることを確認

### 8.6 ドキュメント更新
- `docs/Sharding.md`の設定例を4シャード構成に更新
- シャーディング数の説明を更新

## 9. 参考情報

### 9.1 関連Issue
- GitHub Issue #17: データベースのシャーディング数を増やす

### 9.2 既存ドキュメント
- `docs/Sharding.md`: シャーディング戦略の詳細
- `config/develop/database.yaml`: 開発環境設定ファイル
- `config/staging/database.yaml`: ステージング環境設定ファイル
- `config/production/database.yaml.example`: 本番環境設定ファイル例

### 9.3 既存実装
- `server/internal/db/sharding.go`: `HashBasedSharding`の実装
- `server/internal/config/config.go`: 設定構造体の定義
- `server/internal/db/connection.go`: データベース接続処理
- `db/migrations/shard1/001_init.sql`: shard1の初期化スクリプト
- `db/migrations/shard2/001_init.sql`: shard2の初期化スクリプト

### 9.4 シャーディング戦略
- **方式**: Hash-based sharding
- **シャードキー**: `user_id`
- **ハッシュ関数**: FNV-1a
- **シャードID範囲**: 1からN（Nはシャード数）
- **データ分散**: ハッシュ値に基づいて均等に分散

### 9.5 データベースドライバ
- **開発環境**: SQLite3（`sqlite3`）
- **ステージング/本番環境**: PostgreSQL（`postgres`）
- **接続管理**: GORM（`gorm.io/gorm`）
- **Writer/Reader分離**: `gorm.io/plugin/dbresolver`を使用

