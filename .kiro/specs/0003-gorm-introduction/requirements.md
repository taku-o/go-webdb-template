# GORM導入要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #3
- **Issueタイトル**: gormの導入
- **作成日**: 2025-01-27

### 1.2 目的
データベース処理にGORMを導入し、以下の機能を実現する：
1. GORM本体によるORM機能の提供
2. Writer/Readerデータベースの使い分け（`gorm.io/plugin/dbresolver`）
3. シャーディング機能の統合（`gorm.io/sharding`）

### 1.3 スコープ
- 既存の`database/sql`ベースの実装をGORMに移行
- 既存のシャーディング戦略（Hash-based）をGORM Shardingプラグインに統合
- Writer/Reader分離機能の新規実装
- 既存のレイヤードアーキテクチャの維持
- 既存のAPIエンドポイントの動作保証

## 2. 背景・現状分析

### 2.1 現在の実装
- **データベース接続**: `database/sql`パッケージを使用
- **SQLクエリ**: 手動でのSQL構築とスキャン処理
- **シャーディング**: 独自実装の`HashBasedSharding`戦略
- **Writer/Reader分離**: 未実装（設定にも定義なし）

### 2.2 課題点
1. **保守性の低さ**: 手動SQL構築によるコードの冗長性とエラーの発生しやすさ
2. **型安全性の不足**: 実行時エラーのリスク（SQL構文エラー、型不一致など）
3. **コードの重複**: 類似したCRUD操作の繰り返し実装
4. **テストの複雑さ**: SQLクエリのモック化が困難
5. **Writer/Reader分離の未対応**: 読み取り専用レプリカの活用ができない

### 2.3 GORM導入による改善点
1. **型安全性の向上**: GORMの型安全なAPIによるコンパイル時エラー検出
2. **コードの簡潔化**: 宣言的なAPIによるコード量の削減
3. **保守性の向上**: モデル定義とクエリの一元管理
4. **テスト容易性**: GORMのモック機能によるテストの簡素化
5. **機能の拡張**: Writer/Reader分離、マイグレーション、リレーション管理などの機能追加

## 3. 機能要件

### 3.1 GORM本体の導入

#### 3.1.1 バージョン指定
- **GORM本体**: `gorm.io/gorm` (最新安定版)
- **SQLiteドライバ**: `gorm.io/driver/sqlite` (開発環境用)
- **PostgreSQLドライバ**: `gorm.io/driver/postgres` (本番環境用)
- **MySQLドライバ**: `gorm.io/driver/mysql` (本番環境用、必要に応じて)

#### 3.1.2 既存実装からの移行
- `server/internal/db/connection.go`: `*sql.DB`から`*gorm.DB`への移行
- `server/internal/db/manager.go`: GORM接続管理への移行
- `server/internal/repository/user_repository.go`: GORM APIへの置き換え
- `server/internal/repository/post_repository.go`: GORM APIへの置き換え

#### 3.1.3 モデル定義の更新
- `server/internal/model/user.go`: GORMタグの追加
  - テーブル名指定: `gorm:"table:users"`
  - 主キー指定: `gorm:"primaryKey"`
  - インデックス指定: `gorm:"index:idx_users_email"`
  - タイムスタンプ自動更新: `gorm:"autoUpdateTime"`
- `server/internal/model/post.go`: GORMタグの追加
  - 外部キー制約: `gorm:"foreignKey:UserID"`
  - インデックス指定: `gorm:"index:idx_posts_user_id"`

#### 3.1.4 リポジトリ層のGORM APIへの置き換え
- `Create`: `db.Create()`を使用
- `GetByID`: `db.First()`を使用
- `List`: `db.Find()`を使用
- `Update`: `db.Model().Updates()`を使用
- `Delete`: `db.Delete()`を使用

### 3.2 Writer/Reader分離（dbresolver）

#### 3.2.1 設定構造の拡張
- `server/internal/config/config.go`の`ShardConfig`に以下を追加：
  - `WriterDSN string`: Writer接続用DSN
  - `ReaderDSNs []string`: Reader接続用DSNリスト（複数可）
  - `ReaderPolicy string`: Reader選択ポリシー（"random", "round_robin"など）

#### 3.2.2 設定ファイルの更新
- `config/develop.yaml`: 各シャードにWriter/Reader設定を追加
- `config/staging.yaml`: 各シャードにWriter/Reader設定を追加
- `config/production.yaml.example`: 各シャードにWriter/Reader設定を追加

#### 3.2.3 接続管理の実装
- 各シャードごとにWriter/Reader接続を分離
- 読み取り操作（SELECT）は自動的にReader接続を使用
- 書き込み操作（INSERT, UPDATE, DELETE）は自動的にWriter接続を使用
- トランザクションは常にWriter接続を使用

#### 3.2.4 dbresolverプラグインの統合
- `gorm.io/plugin/dbresolver`プラグインを使用
- 各シャードのGORMインスタンスにプラグインを登録
- 既存の`Manager`構造体との統合

### 3.3 シャーディング（gorm.io/sharding）

#### 3.3.1 既存シャーディング戦略との統合
- 既存の`HashBasedSharding`戦略を維持
- シャードキー: `user_id`（既存と同じ）
- ハッシュアルゴリズム: FNV-1a（既存と同じ）

#### 3.3.2 GORM Shardingプラグインの設定
- `gorm.io/sharding`プラグインを使用
- 各シャードのGORMインスタンスを登録
- シャードキーに基づく自動ルーティングを実装

#### 3.3.3 クエリルーティング
- 単一シャードクエリ: シャードキー（`user_id`）を含むクエリは自動的に適切なシャードにルーティング
- クロスシャードクエリ: シャードキーを含まないクエリは全シャードに実行し、結果をマージ

#### 3.3.4 既存機能の維持
- `GetConnectionByKey()`: GORMインスタンスを返すように変更
- `GetAllConnections()`: 全シャードのGORMインスタンスを返すように変更
- クロスシャードJOIN: アプリケーションレベルでのJOIN実装を維持

## 4. 非機能要件

### 4.1 後方互換性
- 既存のAPIエンドポイントの動作を維持
- 既存のテストスイートが全てパスする
- 既存の設定ファイル構造との互換性（段階的移行を考慮）

### 4.2 パフォーマンス
- 既存実装と同等以上のパフォーマンスを維持
- 接続プールの適切な設定
- クエリの最適化（N+1問題の回避など）

### 4.3 テストカバレッジ
- 既存のテストカバレッジ（80%以上）を維持
- GORM実装に対する新規テストの追加
- Writer/Reader分離のテスト追加
- シャーディング機能のテスト追加

### 4.4 エラーハンドリング
- 既存のエラーハンドリングパターンを維持
- GORMのエラーを既存のエラー形式に変換
- 適切なエラーメッセージとログ出力

### 4.5 ドキュメント
- アーキテクチャドキュメントの更新
- APIドキュメントの更新
- 設定ファイルの説明追加
- 移行ガイドの作成

## 5. 制約事項

### 5.1 アーキテクチャ
- 既存のレイヤードアーキテクチャ（API → Service → Repository → DB）を維持
- 各レイヤーの責務を明確に分離

### 5.2 シャーディング戦略
- 既存のHash-based sharding戦略を維持
- シャードキー（`user_id`）の変更は行わない

### 5.3 設定ファイル
- 既存のYAML設定ファイル形式を維持
- 環境変数による設定上書き機能を維持

### 5.4 テスト
- 既存のテストフレームワーク（`testing`, `testify`）を維持
- 既存のテストディレクトリ構造を維持

### 5.5 データベース
- 開発環境: SQLite3
- 本番環境: PostgreSQL / MySQL
- 既存のマイグレーションファイル構造を維持

## 6. 受け入れ基準

### 6.1 機能要件
- [ ] GORM本体が正常に動作する
- [ ] Writer/Reader分離が正しく機能する（読み取りはReader、書き込みはWriter）
- [ ] シャーディングが正しく機能する（シャードキーに基づく自動ルーティング）
- [ ] 既存の全APIエンドポイントが正常に動作する
- [ ] クロスシャードクエリが正常に動作する

### 6.2 非機能要件
- [ ] 全既存テストがパスする
- [ ] 新規GORM実装のテストが追加され、パスする
- [ ] Writer/Reader分離のテストが追加され、パスする
- [ ] シャーディング機能のテストが追加され、パスする
- [ ] テストカバレッジが80%以上を維持

### 6.3 ドキュメント
- [ ] アーキテクチャドキュメント（`docs/Architecture.md`）が更新される
- [ ] 設定ファイルの説明が追加される
- [ ] 移行ガイドが作成される（必要に応じて）

### 6.4 コード品質
- [ ] 既存のコードスタイルに準拠
- [ ] 適切なエラーハンドリングが実装されている
- [ ] 適切なログ出力が実装されている

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### データベース層
- `server/internal/db/connection.go`: `*sql.DB`から`*gorm.DB`への移行
- `server/internal/db/manager.go`: GORM接続管理への移行、Writer/Reader分離の実装
- `server/internal/db/sharding.go`: GORM Shardingプラグインとの統合（必要に応じて）

#### 設定層
- `server/internal/config/config.go`: `ShardConfig`にWriter/Reader設定を追加

#### モデル層
- `server/internal/model/user.go`: GORMタグの追加
- `server/internal/model/post.go`: GORMタグの追加

#### リポジトリ層
- `server/internal/repository/user_repository.go`: GORM APIへの完全置き換え
- `server/internal/repository/post_repository.go`: GORM APIへの完全置き換え
- `server/internal/repository/user_repository_test.go`: GORM実装に合わせたテスト更新
- `server/internal/repository/post_repository_test.go`: GORM実装に合わせたテスト更新

#### 設定ファイル
- `config/develop.yaml`: Writer/Reader設定の追加
- `config/staging.yaml`: Writer/Reader設定の追加
- `config/production.yaml.example`: Writer/Reader設定の追加

#### 依存関係
- `server/go.mod`: GORM関連パッケージの追加
- `server/go.sum`: 依存関係の更新

### 7.2 新規追加が必要なファイル
- なし（既存ファイルの修正のみ）

### 7.3 削除されるファイル
- なし（既存ファイルの修正のみ）

### 7.4 ドキュメント更新
- `docs/Architecture.md`: GORM導入によるアーキテクチャ変更の反映
- `docs/Sharding.md`: GORM Shardingプラグインの説明追加
- `README.md`: 依存関係の更新（必要に応じて）

## 8. 実装上の注意事項

### 8.1 移行戦略
- 段階的な移行を検討（ただし、本要件では完全移行を想定）
- 既存の`database/sql`コードとGORMコードの共存期間を最小化

### 8.2 Writer/Reader分離の実装
- 開発環境ではWriter/Readerを同一データベースに設定可能（設定の簡素化）
- 本番環境では別々のデータベースインスタンスを想定

### 8.3 シャーディングの実装
- GORM Shardingプラグインの制約を理解し、既存のシャーディング戦略と整合性を保つ
- クロスシャードクエリのパフォーマンスに注意

### 8.4 エラーハンドリング
- GORMのエラーを既存のエラーハンドリングパターンに変換
- データベース接続エラー、クエリエラー、制約違反エラーなどの適切な処理

### 8.5 テスト
- GORMのモック機能を活用したテスト実装
- 統合テストでの実際のデータベース接続テスト
- Writer/Reader分離の動作確認テスト

## 9. 参考情報

### 9.1 GORM公式ドキュメント
- GORM本体: https://gorm.io/docs/
- dbresolverプラグイン: https://gorm.io/docs/dbresolver.html
- shardingプラグイン: https://gorm.io/docs/sharding.html

### 9.2 関連Issue
- GitHub Issue #3: gormの導入

### 9.3 既存ドキュメント
- `docs/Architecture.md`: システムアーキテクチャ
- `docs/Sharding.md`: シャーディング戦略の詳細
- `.kiro/steering/tech.md`: 技術スタックとアーキテクチャ

