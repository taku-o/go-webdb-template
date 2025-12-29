# シャーディング数8対応要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #49
- **Issueタイトル**: シャーディング数を8にしたい
- **Feature名**: 0024-sharding8
- **作成日**: 2025-01-27

### 1.2 目的
shardingグループのデータベースのシャーディング数を4から8に増やし、将来のデータベース分割・合併作業を容易にするための基盤を整備する。
設定ファイル上では8つのシャーディングエントリを定義するが、実際のデータベース接続は4つのまま維持する。

### 1.3 スコープ
- 設定ファイルに8つのシャーディングエントリを追加（各4テーブルを担当）
- 各シャーディングエントリが適切なデータベースファイルに接続するように設定
- データベース接続管理ロジックの変更（table_rangeベースの接続選択）
- 既存のデータベースファイルはそのまま使用（データ移行は不要）
- 32分割のテーブル構造は維持（suffix_count: 32）

**本実装の範囲外**:
- データベースファイルの分割・合併（将来の作業のための準備）
- 既存データの移行（既存データは破棄して良い）
- テーブル数の変更（32分割は維持）

## 2. 背景・現状分析

### 2.1 現在の実装
- **シャーディング数**: 4つ（sharding 1-4）
- **テーブル分散**: 各シャーディングが8テーブルを担当
  - sharding 1: テーブル _000-007（sharding_db_1.db）
  - sharding 2: テーブル _008-015（sharding_db_2.db）
  - sharding 3: テーブル _016-023（sharding_db_3.db）
  - sharding 4: テーブル _024-031（sharding_db_4.db）
- **データベース接続**: 4つのデータベースファイル（sharding_db_1.db ～ sharding_db_4.db）
- **接続選択ロジック**: `dbID := (tableNumber / 8) + 1`でデータベースIDを決定
- **設定ファイル**: `config/{env}/database.yaml`に4つのシャーディングエントリを定義

### 2.2 課題点
1. **将来の拡張性**: データベースを分割・合併する際に、設定ファイルの変更が大規模になる可能性がある
2. **設定の柔軟性**: より細かい粒度でのシャーディング管理ができない
3. **運用の複雑さ**: データベース分割・合併時の作業が複雑になる可能性がある

### 2.3 本実装による改善点
1. **将来の拡張性向上**: 8つのシャーディングエントリにより、より細かい粒度での管理が可能
2. **データベース分割・合併の容易化**: 設定ファイル上で8つのエントリを管理することで、将来の分割・合併作業が容易になる
3. **設定の柔軟性向上**: 各シャーディングエントリが独立して管理できるため、柔軟な運用が可能

## 3. 機能要件

### 3.1 シャーディング数の変更

#### 3.1.1 シャーディングエントリの定義
設定ファイル上で8つのシャーディングエントリを定義する：

- **sharding 1**: テーブル _000-003（4テーブル）→ sharding_db_1.dbに接続
- **sharding 2**: テーブル _004-007（4テーブル）→ sharding_db_1.dbに接続
- **sharding 3**: テーブル _008-011（4テーブル）→ sharding_db_2.dbに接続
- **sharding 4**: テーブル _012-015（4テーブル）→ sharding_db_2.dbに接続
- **sharding 5**: テーブル _016-019（4テーブル）→ sharding_db_3.dbに接続
- **sharding 6**: テーブル _020-023（4テーブル）→ sharding_db_3.dbに接続
- **sharding 7**: テーブル _024-027（4テーブル）→ sharding_db_4.dbに接続
- **sharding 8**: テーブル _028-031（4テーブル）→ sharding_db_4.dbに接続

#### 3.1.2 データベース接続の実態
実際のデータベース接続は4つのまま維持する：

- **sharding_db_1.db**: sharding 1, 2が接続（テーブル _000-007）
- **sharding_db_2.db**: sharding 3, 4が接続（テーブル _008-015）
- **sharding_db_3.db**: sharding 5, 6が接続（テーブル _016-023）
- **sharding_db_4.db**: sharding 7, 8が接続（テーブル _024-031）

### 3.2 設定ファイルの変更

#### 3.2.1 設定ファイル構造
`config/{env}/database.yaml`の`sharding.databases`セクションを8つのエントリに拡張：

```yaml
database:
  groups:
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
          table_range: [0, 3]  # _000-003
        - id: 2
          driver: sqlite3
          dsn: ./data/sharding_db_1.db
          writer_dsn: ./data/sharding_db_1.db
          reader_dsns:
            - ./data/sharding_db_1.db
          reader_policy: random
          max_connections: 10
          max_idle_connections: 5
          connection_max_lifetime: 300s
          table_range: [4, 7]  # _004-007
        - id: 3
          driver: sqlite3
          dsn: ./data/sharding_db_2.db
          writer_dsn: ./data/sharding_db_2.db
          reader_dsns:
            - ./data/sharding_db_2.db
          reader_policy: random
          max_connections: 10
          max_idle_connections: 5
          connection_max_lifetime: 300s
          table_range: [8, 11]  # _008-011
        - id: 4
          driver: sqlite3
          dsn: ./data/sharding_db_2.db
          writer_dsn: ./data/sharding_db_2.db
          reader_dsns:
            - ./data/sharding_db_2.db
          reader_policy: random
          max_connections: 10
          max_idle_connections: 5
          connection_max_lifetime: 300s
          table_range: [12, 15]  # _012-015
        - id: 5
          driver: sqlite3
          dsn: ./data/sharding_db_3.db
          writer_dsn: ./data/sharding_db_3.db
          reader_dsns:
            - ./data/sharding_db_3.db
          reader_policy: random
          max_connections: 10
          max_idle_connections: 5
          connection_max_lifetime: 300s
          table_range: [16, 19]  # _016-019
        - id: 6
          driver: sqlite3
          dsn: ./data/sharding_db_3.db
          writer_dsn: ./data/sharding_db_3.db
          reader_dsns:
            - ./data/sharding_db_3.db
          reader_policy: random
          max_connections: 10
          max_idle_connections: 5
          connection_max_lifetime: 300s
          table_range: [20, 23]  # _020-023
        - id: 7
          driver: sqlite3
          dsn: ./data/sharding_db_4.db
          writer_dsn: ./data/sharding_db_4.db
          reader_dsns:
            - ./data/sharding_db_4.db
          reader_policy: random
          max_connections: 10
          max_idle_connections: 5
          connection_max_lifetime: 300s
          table_range: [24, 27]  # _024-027
        - id: 8
          driver: sqlite3
          dsn: ./data/sharding_db_4.db
          writer_dsn: ./data/sharding_db_4.db
          reader_dsns:
            - ./data/sharding_db_4.db
          reader_policy: random
          max_connections: 10
          max_idle_connections: 5
          connection_max_lifetime: 300s
          table_range: [28, 31]  # _028-031

      tables:
        - name: users
          suffix_count: 32
        - name: posts
          suffix_count: 32
```

#### 3.2.2 環境別設定ファイルの更新
以下の環境別設定ファイルを更新する：
- `config/develop/database.yaml`
- `config/staging/database.yaml`（存在する場合）
- `config/production/database.yaml.example`（存在する場合）
- `server/internal/config/testdata/develop/database.yaml`（テスト用設定ファイル）

### 3.3 データベース接続管理ロジックの変更

#### 3.3.1 接続選択ロジックの変更
`server/internal/db/group_manager.go`の`GetConnectionByTableNumber`メソッドを変更する：

**変更前**:
```go
// テーブル番号からデータベースIDを決定
dbID := (tableNumber / 8) + 1
```

**変更後**:
```go
// テーブル番号からデータベースIDを決定（table_rangeベース）
// 各シャーディングエントリのtable_rangeを確認して、該当するエントリのIDを返す
for dbID, tableRange := range sm.tableRange {
    if tableNumber >= tableRange[0] && tableNumber <= tableRange[1] {
        return sm.connections[dbID], nil
    }
}
```

#### 3.3.2 接続の共有
同じデータベースファイルに接続する複数のシャーディングエントリは、同じ接続オブジェクトを共有する必要がある。
`NewShardingManager`メソッドで、同じDSNを持つエントリに対しては既存の接続を再利用する。

#### 3.3.3 接続管理の実装
- 接続オブジェクトはDSNをキーとして管理する
- 同じDSNを持つ複数のシャーディングエントリは、同じ接続オブジェクトを参照する
- 接続のクローズ時は、すべての参照がなくなった時点でクローズする

### 3.4 既存データの扱い

#### 3.4.1 データ移行の方針
- **既存データの移行は不要**: Issue要件に基づき、既存データの移行は行わない
- **データ損失を許容**: 既存データは消失しても構わない
- **新規データのみ対応**: 新規に作成されるデータのみ、新しいルールに従う

#### 3.4.2 データベースファイルの扱い
- 既存のデータベースファイル（sharding_db_1.db ～ sharding_db_4.db）はそのまま使用
- データベースファイルの分割・合併は行わない
- マイグレーションは既存のテーブル構造を維持

### 3.5 テーブル構造の維持

#### 3.5.1 テーブル分割数の維持
- 32分割のテーブル構造は維持（`suffix_count: 32`）
- テーブル名の形式は変更しない（`users_000` ～ `users_031`、`posts_000` ～ `posts_031`）

#### 3.5.2 テーブル選択ロジックの維持
- テーブル選択ロジック（`id % 32`）は変更しない
- テーブル番号の計算方法は変更しない

## 4. 非機能要件

### 4.1 パフォーマンス
- 接続選択の計算はO(n)で実行される（nはシャーディングエントリ数、最大8）
- 接続の再利用により、接続プールの効率が向上する
- 既存のクエリパフォーマンスは維持される

### 4.2 拡張性
- シャーディング数の変更（8以外）に対応できる設計とすること
- データベース数の変更（4以外）に対応できる設計とすること
- 設定ファイルで柔軟に構成変更できること

### 4.3 保守性
- 設定ファイルで柔軟に構成変更できること
- コードの可読性とテスト容易性を維持すること
- 接続管理ロジックが明確であること

### 4.4 後方互換性
- 既存のAPIエンドポイントの動作は維持されること
- 既存のテストコードは可能な限り動作すること（大幅な変更が必要な場合はテストコードの更新も許容）
- データベース構造の変更は行わない

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

### 5.4 データベースファイル
- 既存のデータベースファイル（sharding_db_1.db ～ sharding_db_4.db）はそのまま使用
- データベースファイルの分割・合併は行わない
- マイグレーションは既存のテーブル構造を維持

## 6. 受け入れ基準

### 6.1 設定ファイル
- [ ] `config/develop/database.yaml`に8つのシャーディングエントリが定義されている
- [ ] 各シャーディングエントリの`table_range`が正しく設定されている
- [ ] 各シャーディングエントリが適切なデータベースファイルに接続するように設定されている
- [ ] 環境別設定ファイル（staging, production）が更新されている（存在する場合）
- [ ] テスト用設定ファイルが更新されている

### 6.2 データベース接続管理
- [ ] `GetConnectionByTableNumber`メソッドが`table_range`ベースで接続を選択する
- [ ] 同じデータベースファイルに接続する複数のシャーディングエントリが、同じ接続オブジェクトを共有する
- [ ] 接続のクローズが適切に実装されている
- [ ] 既存の接続取得メソッド（`GetShardingConnectionByID`など）が正常に動作する

### 6.3 テーブル選択ロジック
- [ ] テーブル番号0-31が正しいシャーディングエントリにマッピングされる
- [ ] 各テーブル番号が正しいデータベースファイルに接続される
- [ ] テーブル選択ロジック（`id % 32`）が維持される

### 6.4 既存機能の動作確認
- [ ] 既存のCRUD操作が正常に動作する
- [ ] クロステーブルクエリが正常に動作する
- [ ] 既存のAPIエンドポイントが正常に動作する

### 6.5 テスト
- [ ] 単体テストが実装されている
- [ ] 統合テストが実装されている
- [ ] 既存のテストが可能な限り動作する（大幅な変更が必要な場合はテストコードの更新も許容）

### 6.6 ドキュメント
- [ ] `docs/Sharding.md`が更新されている
- [ ] 新しいシャーディング構成が文書化されている
- [ ] 設定ファイルの変更内容が文書化されている

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### 設定ファイル
- `config/develop/database.yaml`: 8つのシャーディングエントリに拡張
- `config/staging/database.yaml`: 8つのシャーディングエントリに拡張（存在する場合）
- `config/production/database.yaml.example`: 8つのシャーディングエントリに拡張（存在する場合）
- `server/internal/config/testdata/develop/database.yaml`: 8つのシャーディングエントリに拡張

#### データベース接続管理
- `server/internal/db/group_manager.go`: `GetConnectionByTableNumber`メソッドの変更、接続共有の実装
- `server/internal/db/sharding.go`: `GetShardingDBID`関数の変更（必要に応じて）

#### テストファイル
- `server/internal/db/group_manager_test.go`: 接続選択ロジックのテスト更新
- `server/test/integration/sharding_test.go`: 統合テストの更新（必要に応じて）

### 7.2 ドキュメント
- `docs/Sharding.md`: 新しいシャーディング構成を反映

### 7.3 新規追加が必要なファイル
なし（既存ファイルの変更のみ）

### 7.4 削除されるファイル
なし（既存ファイルは変更のみ）

## 8. 実装上の注意事項

### 8.1 接続選択ロジックの変更
- `table_range`ベースの接続選択を実装する
- テーブル番号が複数の`table_range`に該当しないように、範囲が重複しないことを確認する
- パフォーマンスを考慮し、必要に応じてインデックスやキャッシュを検討する

### 8.2 接続の共有
- 同じDSNを持つ複数のシャーディングエントリが、同じ接続オブジェクトを共有する
- 接続のクローズ時は、すべての参照がなくなった時点でクローズする
- 接続の参照カウントを管理する必要がある

### 8.3 設定ファイルの整合性
- 8つのシャーディングエントリの`table_range`が、0-31の範囲を完全にカバーすることを確認する
- 範囲の重複や欠落がないことを確認する
- 各シャーディングエントリが適切なデータベースファイルに接続することを確認する

### 8.4 テスト
- 単体テストで接続選択ロジックを検証する
- 統合テストで実際のデータベース操作を検証する
- 既存のテストが可能な限り動作することを確認する

### 8.5 パフォーマンス
- 接続選択の計算はO(n)で実行される（nはシャーディングエントリ数、最大8）
- 必要に応じて、テーブル番号からシャーディングエントリIDへのマッピングをキャッシュする

## 9. 参考情報

### 9.1 関連Issue
- GitHub Issue #49: シャーディング数を8にしたい

### 9.2 既存ドキュメント
- `docs/Sharding.md`: 既存のシャーディング戦略の詳細
- `.kiro/specs/0012-sharding/requirements.md`: シャーディング規則修正の要件定義書
- `.kiro/specs/0012-sharding/design.md`: シャーディング規則修正の設計書

### 9.3 既存実装
- `server/internal/db/group_manager.go`: 既存の接続管理
- `server/internal/db/sharding.go`: 既存のシャーディング戦略
- `config/develop/database.yaml`: 既存の設定ファイル

### 9.4 技術スタック
- **Go**: 1.21+
- **GORM**: v1.25.12
- **データベース**: SQLite3（開発環境）、PostgreSQL（本番環境）
- **設定管理**: viper（spf13/viper）
