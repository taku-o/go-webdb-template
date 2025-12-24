# GORM導入実装タスク一覧

## 概要
GORM導入の実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: 依存関係とインフラ準備

#### タスク 1.1: GORM関連パッケージの追加
**目的**: GORM本体とドライバー、プラグインをgo.modに追加

**作業内容**:
- `server/go.mod`に以下を追加:
  - `gorm.io/gorm` (最新安定版)
  - `gorm.io/driver/sqlite` (開発環境用)
  - `gorm.io/driver/postgres` (本番環境用)
  - `gorm.io/plugin/dbresolver` (Writer/Reader分離用)
  - `gorm.io/sharding` (シャーディング用)
- `go mod tidy`を実行して依存関係を解決
- `go mod download`でパッケージをダウンロード

**受け入れ基準**:
- `go mod tidy`がエラーなく実行できる
- すべての依存関係が正しく解決されている
- `go.sum`が更新されている

---

#### タスク 1.2: 設定構造の拡張
**目的**: ShardConfigにWriter/Reader設定を追加

**作業内容**:
- `server/internal/config/config.go`の`ShardConfig`構造体に以下を追加:
  - `WriterDSN string`: Writer接続用DSN
  - `ReaderDSNs []string`: Reader接続用DSNリスト
  - `ReaderPolicy string`: Reader選択ポリシー（"random" or "round_robin"）
- `GetWriterDSN()`メソッドを追加（後方互換性のため）
- `GetReaderDSNs()`メソッドを追加（後方互換性のため）

**受け入れ基準**:
- コンパイルエラーがない
- 既存の設定ファイルが動作する（後方互換性）
- 新しい設定項目が正しく読み込まれる

---

### Phase 2: モデル定義の更新

#### タスク 2.1: UserモデルにGORMタグを追加
**目的**: UserモデルにGORMタグを追加してテーブル定義を明確化

**作業内容**:
- `server/internal/model/user.go`の`User`構造体にGORMタグを追加:
  - `gorm:"primaryKey"` (IDフィールド)
  - `gorm:"type:varchar(100);not null"` (Nameフィールド)
  - `gorm:"type:varchar(255);not null;uniqueIndex:idx_users_email"` (Emailフィールド)
  - `gorm:"autoCreateTime"` (CreatedAtフィールド)
  - `gorm:"autoUpdateTime"` (UpdatedAtフィールド)
- `TableName()`メソッドを追加してテーブル名を明示的に指定

**受け入れ基準**:
- コンパイルエラーがない
- GORMタグが正しく設定されている
- 既存のJSONタグとdbタグが維持されている

---

#### タスク 2.2: PostモデルにGORMタグを追加
**目的**: PostモデルにGORMタグを追加してテーブル定義を明確化

**作業内容**:
- `server/internal/model/post.go`の`Post`構造体にGORMタグを追加:
  - `gorm:"primaryKey"` (IDフィールド)
  - `gorm:"type:bigint;not null;index:idx_posts_user_id"` (UserIDフィールド)
  - `gorm:"type:varchar(200);not null"` (Titleフィールド)
  - `gorm:"type:text;not null"` (Contentフィールド)
  - `gorm:"autoCreateTime"` (CreatedAtフィールド)
  - `gorm:"autoUpdateTime"` (UpdatedAtフィールド)
- `TableName()`メソッドを追加してテーブル名を明示的に指定

**受け入れ基準**:
- コンパイルエラーがない
- GORMタグが正しく設定されている
- 既存のJSONタグとdbタグが維持されている

---

#### タスク 2.3: UserPostモデルの確認
**目的**: UserPostモデル（JOIN結果用）がGORMで正しく動作することを確認

**作業内容**:
- `server/internal/model/post.go`の`UserPost`構造体を確認
- GORMタグが適切に設定されているか確認（必要に応じて追加）
- `TableName()`メソッドで空文字列を返すことを確認

**受け入れ基準**:
- UserPostモデルがGORMで正しく使用できる
- JOIN結果が正しくマッピングされる

---

### Phase 3: DB層のGORM移行

#### タスク 3.1: GORM接続作成ヘルパー関数の実装
**目的**: GORM接続を作成するヘルパー関数を実装

**作業内容**:
- `server/internal/db/connection.go`に`createGORMConnection()`関数を追加:
  - ドライバーに応じて適切なDialectorを作成
  - GORMインスタンスを作成
  - 接続プール設定を適用
- SQLite、PostgreSQL、MySQLの各ドライバーに対応

**受け入れ基準**:
- 各ドライバーでGORM接続が正しく作成される
- 接続プール設定が正しく適用される
- エラーハンドリングが適切

---

#### タスク 3.2: GORMConnection構造体の実装
**目的**: 単一シャードのGORM接続を管理する構造体を実装

**作業内容**:
- `server/internal/db/connection.go`に`GORMConnection`構造体を追加:
  - `*gorm.DB`フィールド
  - `ShardID`フィールド
  - `Driver`フィールド
  - `config`フィールド
- `NewGORMConnection()`関数を実装:
  - Writer接続を作成
  - Reader接続を作成（複数可）
  - dbresolverプラグインを設定（Readerがある場合）
- `Close()`メソッドを実装
- `Ping()`メソッドを実装

**受け入れ基準**:
- GORM接続が正しく作成される
- Writer/Reader分離が正しく設定される
- 接続のクローズとPingが正しく動作する

---

#### タスク 3.3: GORMManager構造体の実装
**目的**: 複数シャードのGORM接続を管理する構造体を実装

**作業内容**:
- `server/internal/db/manager.go`に`GORMManager`構造体を追加:
  - `connections map[int]*GORMConnection`
  - `strategy ShardingStrategy`
  - `mu sync.RWMutex`
- `NewGORMManager()`関数を実装:
  - 各シャードのGORM接続を作成
  - エラー時のクリーンアップ処理
- `GetGORM(shardID int)`メソッドを実装
- `GetGORMByKey(key int64)`メソッドを実装
- `GetAllGORMConnections()`メソッドを実装
- `CloseAll()`メソッドを実装
- `PingAll()`メソッドを実装

**受け入れ基準**:
- 複数シャードのGORM接続が正しく管理される
- シャードキーに基づく接続取得が正しく動作する
- クロスシャードクエリ用の接続取得が正しく動作する

---

#### タスク 3.4: 既存Connection/Managerの非推奨化
**目的**: 既存のConnection/Managerを非推奨として残し、段階的移行を可能にする

**作業内容**:
- 既存の`Connection`構造体と`Manager`構造体に`Deprecated`コメントを追加
- 既存のメソッドに`Deprecated`コメントを追加
- 新規コードでは`GORMConnection`と`GORMManager`を使用することを明記

**受け入れ基準**:
- 既存コードが動作し続ける
- 非推奨の警告が明確に表示される

---

### Phase 4: Writer/Reader分離の実装

#### タスク 4.1: dbresolverプラグインの統合
**目的**: dbresolverプラグインをGORM接続に統合

**作業内容**:
- `NewGORMConnection()`内でdbresolverプラグインを設定:
  - Writer接続をSourceとして設定
  - Reader接続をReplicaとして設定
  - Reader選択ポリシーを設定（random or round_robin）
- プラグインの登録エラーハンドリング

**受け入れ基準**:
- dbresolverプラグインが正しく設定される
- 読み取り操作がReader接続を使用する
- 書き込み操作がWriter接続を使用する
- トランザクションがWriter接続を使用する

---

#### タスク 4.2: Writer/Reader分離のテスト
**目的**: Writer/Reader分離が正しく動作することを確認

**作業内容**:
- `server/internal/db/connection_test.go`にテストを追加:
  - Writer接続での書き込みテスト
  - Reader接続での読み取りテスト
  - トランザクションがWriter接続を使用するテスト
- 統合テストでWriter/Reader分離の動作を確認

**受け入れ基準**:
- すべてのテストがパスする
- Writer/Reader分離が正しく動作することを確認

---

### Phase 5: リポジトリ層のGORM API移行

#### タスク 5.1: UserRepositoryのGORM API移行
**目的**: UserRepositoryをGORM APIに移行

**作業内容**:
- `server/internal/repository/user_repository.go`を修正:
  - `dbManager *db.Manager`を`dbManager *db.GORMManager`に変更
  - `Create()`メソッド: `db.Create()`を使用
  - `GetByID()`メソッド: `db.First()`を使用
  - `List()`メソッド: `db.Find()`を使用（クロスシャードクエリ）
  - `Update()`メソッド: `db.Model().Updates()`を使用
  - `Delete()`メソッド: `db.Delete()`を使用
- GORMエラーを既存のエラー形式に変換
- コンテキストの適切な使用（`WithContext(ctx)`）

**受け入れ基準**:
- コンパイルエラーがない
- 既存のメソッドシグネチャが維持される
- GORM APIが正しく使用されている

---

#### タスク 5.2: PostRepositoryのGORM API移行
**目的**: PostRepositoryをGORM APIに移行

**作業内容**:
- `server/internal/repository/post_repository.go`を修正:
  - `dbManager *db.Manager`を`dbManager *db.GORMManager`に変更
  - `Create()`メソッド: `db.Create()`を使用
  - `GetByID()`メソッド: `db.First()`を使用
  - `ListByUserID()`メソッド: `db.Where("user_id = ?", userID).Find()`を使用
  - `List()`メソッド: `db.Find()`を使用（クロスシャードクエリ）
  - `GetUserPosts()`メソッド: クロスシャードJOIN（アプリケーションレベル）
  - `Update()`メソッド: `db.Model().Updates()`を使用
  - `Delete()`メソッド: `db.Delete()`を使用
- GORMエラーを既存のエラー形式に変換
- コンテキストの適切な使用

**受け入れ基準**:
- コンパイルエラーがない
- 既存のメソッドシグネチャが維持される
- GORM APIが正しく使用されている

---

#### タスク 5.3: エラーハンドリングの統一
**目的**: GORMエラーを既存のエラーハンドリングパターンに変換

**作業内容**:
- `server/internal/repository/`に`errors.go`を作成（必要に応じて）
- `convertGORMError()`関数を実装:
  - `gorm.ErrRecordNotFound`を既存のエラー形式に変換
  - その他のGORMエラーを適切にラップ
- 各Repositoryメソッドでエラー変換を使用

**受け入れ基準**:
- GORMエラーが既存のエラー形式に正しく変換される
- エラーメッセージが適切
- エラーハンドリングの一貫性が保たれる

---

### Phase 6: 設定ファイルの更新

#### タスク 6.1: develop.yamlの更新
**目的**: 開発環境の設定ファイルにWriter/Reader設定を追加

**作業内容**:
- `config/develop.yaml`を更新:
  - 各シャードに`writer_dsn`を追加
  - 各シャードに`reader_dsns`を追加（開発環境では同一DB）
  - 各シャードに`reader_policy`を追加（"random"）

**受け入れ基準**:
- 設定ファイルが正しく読み込まれる
- Writer/Reader設定が正しく適用される
- 既存の設定項目が維持される

---

#### タスク 6.2: staging.yamlの更新
**目的**: ステージング環境の設定ファイルにWriter/Reader設定を追加

**作業内容**:
- `config/staging.yaml`を更新:
  - 各シャードに`writer_dsn`を追加
  - 各シャードに`reader_dsns`を追加
  - 各シャードに`reader_policy`を追加

**受け入れ基準**:
- 設定ファイルが正しく読み込まれる
- Writer/Reader設定が正しく適用される

---

#### タスク 6.3: production.yaml.exampleの更新
**目的**: 本番環境の設定ファイル例にWriter/Reader設定を追加

**作業内容**:
- `config/production.yaml.example`を更新:
  - 各シャードに`writer_dsn`を追加（環境変数参照）
  - 各シャードに`reader_dsns`を追加（環境変数参照、複数可）
  - 各シャードに`reader_policy`を追加（"round_robin"）

**受け入れ基準**:
- 設定ファイル例が正しく記述されている
- 環境変数の参照方法が明確

---

### Phase 7: テストの更新

#### タスク 7.1: testutilの更新
**目的**: テストユーティリティをGORM対応に更新

**作業内容**:
- `server/test/testutil/db.go`を更新:
  - `SetupTestShards()`を`SetupTestGORMShards()`に変更または追加
  - GORM接続を作成するように変更
  - テスト用のGORMインスタンスを返す

**受け入れ基準**:
- テストユーティリティがGORMで動作する
- 既存のテストが動作し続ける

---

#### タスク 7.2: UserRepositoryテストの更新
**目的**: UserRepositoryのテストをGORM実装に合わせて更新

**作業内容**:
- `server/internal/repository/user_repository_test.go`を更新:
  - `testutil.SetupTestShards()`を`testutil.SetupTestGORMShards()`に変更
  - GORMエラーのアサーションを更新
  - テストケースが正しく動作することを確認

**受け入れ基準**:
- すべてのテストがパスする
- テストカバレッジが維持される

---

#### タスク 7.3: PostRepositoryテストの更新
**目的**: PostRepositoryのテストをGORM実装に合わせて更新

**作業内容**:
- `server/internal/repository/post_repository_test.go`を更新:
  - `testutil.SetupTestShards()`を`testutil.SetupTestGORMShards()`に変更
  - GORMエラーのアサーションを更新
  - テストケースが正しく動作することを確認

**受け入れ基準**:
- すべてのテストがパスする
- テストカバレッジが維持される

---

#### タスク 7.4: Writer/Reader分離のテスト追加
**目的**: Writer/Reader分離の動作を確認するテストを追加

**作業内容**:
- `server/internal/db/connection_test.go`にテストを追加:
  - Writer接続での書き込みテスト
  - Reader接続での読み取りテスト
  - トランザクションがWriter接続を使用するテスト
  - Reader選択ポリシーのテスト

**受け入れ基準**:
- Writer/Reader分離が正しく動作することを確認
- すべてのテストがパスする

---

#### タスク 7.5: シャーディングのテスト追加
**目的**: シャーディングの動作を確認するテストを追加

**作業内容**:
- `server/internal/db/manager_test.go`にテストを追加:
  - シャードキーに基づくルーティングテスト
  - クロスシャードクエリのテスト
  - データが正しいシャードに保存されることを確認するテスト

**受け入れ基準**:
- シャーディングが正しく動作することを確認
- すべてのテストがパスする

---

#### タスク 7.6: 統合テストの更新
**目的**: 統合テストをGORM実装に合わせて更新

**作業内容**:
- `server/test/integration/user_flow_test.go`を更新
- `server/test/integration/post_flow_test.go`を更新
- `server/test/e2e/api_test.go`を更新
- GORM実装で正常に動作することを確認

**受け入れ基準**:
- すべての統合テストがパスする
- E2Eテストがパスする

---

### Phase 8: メインアプリケーションの更新

#### タスク 8.1: main.goの更新
**目的**: メインアプリケーションでGORMManagerを使用するように更新

**作業内容**:
- `server/cmd/server/main.go`を更新:
  - `db.NewManager()`を`db.NewGORMManager()`に変更
  - `dbManager`の型を`*db.GORMManager`に変更
  - エラーハンドリングを確認

**受け入れ基準**:
- コンパイルエラーがない
- アプリケーションが正常に起動する
- データベース接続が正常に確立される

---

### Phase 9: ドキュメントの更新

#### タスク 9.1: Architecture.mdの更新
**目的**: アーキテクチャドキュメントにGORM導入を反映

**作業内容**:
- `docs/Architecture.md`を更新:
  - DB層の説明をGORMベースに更新
  - Writer/Reader分離の説明を追加
  - GORM Shardingプラグインの説明を追加
  - データフローの説明を更新

**受け入れ基準**:
- ドキュメントが最新の実装を反映している
- アーキテクチャ図が更新されている

---

#### タスク 9.2: Sharding.mdの更新
**目的**: シャーディングドキュメントにGORM Shardingプラグインの説明を追加

**作業内容**:
- `docs/Sharding.md`を更新:
  - GORM Shardingプラグインの使用について説明
  - 既存のHash-based sharding戦略との統合について説明
  - 設定方法の説明を追加

**受け入れ基準**:
- ドキュメントが最新の実装を反映している
- 設定方法が明確に記載されている

---

#### タスク 9.3: README.mdの更新
**目的**: READMEにGORM依存関係の情報を追加

**作業内容**:
- `README.md`を更新:
  - 依存関係セクションにGORM関連パッケージを追加
  - 必要に応じてセットアップ手順を更新

**受け入れ基準**:
- READMEが最新の依存関係を反映している
- セットアップ手順が正確

---

### Phase 10: 最終検証とクリーンアップ

#### タスク 10.1: 全テストの実行
**目的**: すべてのテストがパスすることを確認

**作業内容**:
- ユニットテストを実行: `go test ./...`
- 統合テストを実行
- E2Eテストを実行
- テストカバレッジを確認（80%以上を維持）

**受け入れ基準**:
- すべてのテストがパスする
- テストカバレッジが80%以上を維持
- エラーや警告がない

---

#### タスク 10.2: 動作確認
**目的**: アプリケーションが正常に動作することを確認

**作業内容**:
- サーバーを起動
- 各APIエンドポイントをテスト:
  - ユーザー作成、取得、更新、削除
  - 投稿作成、取得、更新、削除
  - クロスシャードクエリ
- Writer/Reader分離の動作を確認
- シャーディングの動作を確認

**受け入れ基準**:
- すべてのAPIエンドポイントが正常に動作する
- Writer/Reader分離が正しく機能する
- シャーディングが正しく機能する

---

#### タスク 10.3: コードレビューとリファクタリング
**目的**: コードの品質を確認し、必要に応じてリファクタリング

**作業内容**:
- コードスタイルの確認
- エラーハンドリングの確認
- ログ出力の確認
- 不要なコードの削除
- コメントの追加

**受け入れ基準**:
- コードスタイルが一貫している
- エラーハンドリングが適切
- コメントが適切に追加されている

---

#### タスク 10.4: 既存コードの削除（オプション）
**目的**: 既存のdatabase/sqlベースのコードを削除（完全移行後）

**作業内容**:
- 既存の`Connection`構造体と`Manager`構造体を削除（非推奨化後、十分な期間が経過した後）
- 既存の`*sql.DB`を使用しているコードを削除
- 不要なimportを削除

**注意**: このタスクは、GORM実装が完全に動作し、十分なテストが行われた後に実行すること。

**受け入れ基準**:
- 既存のdatabase/sqlコードが削除されている
- コンパイルエラーがない
- すべてのテストがパスする

---

## タスクの依存関係

```
Phase 1 (依存関係とインフラ準備)
  ↓
Phase 2 (モデル定義の更新)
  ↓
Phase 3 (DB層のGORM移行)
  ↓
Phase 4 (Writer/Reader分離の実装)
  ↓
Phase 5 (リポジトリ層のGORM API移行)
  ↓
Phase 6 (設定ファイルの更新)
  ↓
Phase 7 (テストの更新)
  ↓
Phase 8 (メインアプリケーションの更新)
  ↓
Phase 9 (ドキュメントの更新)
  ↓
Phase 10 (最終検証とクリーンアップ)
```

## 実装時の注意事項

1. **段階的な実装**: 各フェーズを順番に実装し、各フェーズでテストを実行
2. **後方互換性**: 既存のコードが動作し続けるように注意
3. **エラーハンドリング**: GORMエラーを既存のエラー形式に変換
4. **テストカバレッジ**: テストカバレッジを80%以上維持
5. **ドキュメント**: 実装と同時にドキュメントを更新

## 受け入れ基準（全体）

- [ ] すべてのタスクが完了している
- [ ] すべてのテストがパスする
- [ ] テストカバレッジが80%以上を維持
- [ ] 既存のAPIエンドポイントが正常に動作する
- [ ] Writer/Reader分離が正しく機能する
- [ ] シャーディングが正しく機能する
- [ ] ドキュメントが更新されている
- [ ] コードレビューが完了している

