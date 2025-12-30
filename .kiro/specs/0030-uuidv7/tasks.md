# UUIDv7導入実装タスク一覧

## 概要
UUIDv7導入の実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: UUIDv7ライブラリの導入

#### - [ ] タスク 1.1: github.com/google/uuidライブラリの導入
**目的**: UUIDv7生成に必要なライブラリをプロジェクトに導入する

**作業内容**:
- `go get github.com/google/uuid`を実行してライブラリを取得
- `go.mod`に`github.com/google/uuid`の依存関係が追加されていることを確認
- `uuid.NewV7()`関数が使用できることを確認

**受け入れ基準**:
- `go.mod`に`github.com/google/uuid`の依存関係が追加されている
- `go get github.com/google/uuid`でライブラリが取得できる
- `uuid.NewV7()`関数が使用できる

---

### Phase 2: ID生成ユーティリティの実装

#### - [ ] タスク 2.1: uuidv7.goの作成
**目的**: UUIDv7生成関数を実装する

**作業内容**:
- `server/internal/util/idgen/uuidv7.go`を作成
- `GenerateUUIDv7() (string, error)`関数を実装:
  - `uuid.NewV7()`を使用してUUIDv7を生成
  - 生成されたUUIDからハイフンを削除（`strings.ReplaceAll(uuid.String(), "-", "")`）
  - 小文字に変換（`strings.ToLower()`）
  - 32文字の文字列として返す
- エラーハンドリングを実装（UUID生成エラーを適切に処理）

**受け入れ基準**:
- `server/internal/util/idgen/uuidv7.go`が作成されている
- `GenerateUUIDv7() (string, error)`関数が実装されている
- 生成されたUUIDからハイフンが削除されている
- 生成されたUUIDが小文字に変換されている
- 生成されたUUIDが32文字であること
- エラーハンドリングが適切に実装されている

---

#### - [ ] タスク 2.2: uuidv7_test.goの作成
**目的**: UUIDv7生成関数の単体テストを作成する

**作業内容**:
- `server/internal/util/idgen/uuidv7_test.go`を作成
- 以下のテストケースを実装:
  - UUIDが32文字であること
  - ハイフンが含まれていないこと
  - 小文字であること
  - 一意性（複数回生成して異なる値が返されること）
  - エラーハンドリング（エラーが適切に処理されること）

**受け入れ基準**:
- `server/internal/util/idgen/uuidv7_test.go`が作成されている
- 上記のテストケースが全て実装されている
- テストが全て通過する

---

### Phase 3: シャーディングキー計算関数の実装

#### - [ ] タスク 3.1: GetTableNumberFromUUID関数の実装
**目的**: UUID文字列からテーブル番号を計算する関数を実装する

**作業内容**:
- `server/internal/db/sharding.go`に`GetTableNumberFromUUID(uuid string) (int, error)`関数を追加
- 実装内容:
  - UUID文字列の長さをチェック（2文字以上であること）
  - 後ろ2文字を取得（`uuid[len(uuid)-2:]`）
  - 16進数として解釈（`strconv.ParseInt(suffix, 16, 64)`）
  - テーブル数（32）で割った余りを計算（`value % 32`）
  - テーブル番号（0～31）を返す
- エラーハンドリングを実装（無効なUUID文字列の場合）

**受け入れ基準**:
- `GetTableNumberFromUUID(uuid string) (int, error)`関数が実装されている
- UUIDの後ろ2文字を正しく取得できること
- 16進数として正しく解釈できること
- テーブル数（32）で割った余りが正しく計算できること
- テーブル番号が0～31の範囲内であること
- エラーハンドリングが適切に実装されている

---

#### - [ ] タスク 3.2: GetTableNameFromUUID関数の実装
**目的**: UUID文字列からテーブル名を生成する関数を実装する

**作業内容**:
- `server/internal/db/sharding.go`に`GetTableNameFromUUID(baseName string, uuid string) (string, error)`関数を追加
- 実装内容:
  - `GetTableNumberFromUUID`を使用してテーブル番号を取得
  - テーブル名を生成（`fmt.Sprintf("%s_%03d", baseName, tableNumber)`形式）
  - テーブル名を返す
- エラーハンドリングを実装（`GetTableNumberFromUUID`のエラーを適切に処理）

**受け入れ基準**:
- `GetTableNameFromUUID(baseName string, uuid string) (string, error)`関数が実装されている
- テーブル名が正しく生成されること（例: `dm_users_000`, `dm_posts_015`）
- エラーハンドリングが適切に実装されている

---

#### - [ ] タスク 3.3: GetShardingConnectionByUUID関数の実装
**目的**: UUID文字列からシャーディング接続を取得する関数を実装する

**作業内容**:
- `server/internal/db/group_manager.go`に`GetShardingConnectionByUUID(uuid string, tableName string) (*GORMConnection, error)`関数を追加
- 実装内容:
  - `GetTableNumberFromUUID`を使用してテーブル番号を取得
  - テーブル番号からデータベースIDを取得（`GetDBID`を使用）
  - 適切な接続を返す（`GetShardingConnection`を使用）
- エラーハンドリングを実装（`GetTableNumberFromUUID`のエラーを適切に処理）

**受け入れ基準**:
- `GetShardingConnectionByUUID(uuid string, tableName string) (*GORMConnection, error)`関数が実装されている
- UUIDから正しい接続が取得できること
- エラーハンドリングが適切に実装されている

---

#### - [ ] タスク 3.4: シャーディングキー計算関数のテスト作成
**目的**: シャーディングキー計算関数の単体テストを作成する

**作業内容**:
- `server/internal/db/sharding_test.go`に新規関数のテストを追加
- 以下のテストケースを実装:
  - `GetTableNumberFromUUID`のテスト:
    - 後ろ2文字が00の場合（テーブル番号: 0）
    - 後ろ2文字が0fの場合（テーブル番号: 15）
    - 後ろ2文字が1fの場合（テーブル番号: 31）
    - 後ろ2文字が20の場合（テーブル番号: 0、32 % 32 = 0）
    - 短すぎるUUIDの場合（エラー）
    - 無効な16進数の場合（エラー）
  - `GetTableNameFromUUID`のテスト:
    - 様々なUUID文字列でテーブル名が正しく生成されること
  - `GetShardingConnectionByUUID`のテスト:
    - UUIDから正しい接続が取得できること

**受け入れ基準**:
- `sharding_test.go`に新規関数のテストが追加されている
- 上記のテストケースが全て実装されている
- テストが全て通過する

---

### Phase 4: モデル定義の変更

#### - [ ] タスク 4.1: DmUserモデルの変更
**目的**: DmUserモデルのIDの型を`int64`から`string`に変更する

**作業内容**:
- `server/internal/model/dm_user.go`を更新
- `ID int64` → `ID string`に変更
- JSONタグ: `json:"id,string"` → `json:"id"`
- GORMタグ: `gorm:"primaryKey"` → `gorm:"primaryKey;type:varchar(32)"`

**受け入れ基準**:
- `DmUser.ID`が`string`型に変更されている
- GORMタグが適切に設定されている（`type:varchar(32)`）
- JSONタグが適切に設定されている

---

#### - [ ] タスク 4.2: DmPostモデルの変更
**目的**: DmPostモデルのIDとUserIDの型を`int64`から`string`に変更する

**作業内容**:
- `server/internal/model/dm_post.go`を更新
- `ID int64` → `ID string`に変更
- `UserID int64` → `UserID string`に変更
- JSONタグ: `json:"id,string"` → `json:"id"`, `json:"user_id,string"` → `json:"user_id"`
- GORMタグ:
  - `ID`: `gorm:"primaryKey"` → `gorm:"primaryKey;type:varchar(32)"`
  - `UserID`: `gorm:"type:bigint;not null;index:idx_dm_posts_user_id"` → `gorm:"type:varchar(32);not null;index:idx_dm_posts_user_id"`

**受け入れ基準**:
- `DmPost.ID`が`string`型に変更されている
- `DmPost.UserID`が`string`型に変更されている
- GORMタグが適切に設定されている（`type:varchar(32)`）
- JSONタグが適切に設定されている

---

#### - [ ] タスク 4.3: リクエストモデルの変更
**目的**: リクエストモデルのIDの型を`int64`から`string`に変更し、バリデーションルールを更新する

**作業内容**:
- `server/internal/model/dm_post.go`の`CreateDmPostRequest`を更新
- `UserID int64` → `UserID string`に変更
- バリデーション: `validate:"required,gt=0"` → `validate:"required,len=32"`（UUID文字列の長さチェック）

**受け入れ基準**:
- `CreateDmPostRequest.UserID`が`string`型に変更されている
- バリデーションルールが適切に設定されている（`len=32`）

---

### Phase 5: リポジトリ層の変更

#### - [ ] タスク 5.1: DmUserRepositoryの変更（標準版）
**目的**: DmUserRepository（標準版）のID生成とシャーディングキー計算をUUIDv7ベースに変更する

**作業内容**:
- `server/internal/repository/dm_user_repository.go`を更新
- `Create`メソッド:
  - `idgen.GenerateSonyflakeID()` → `idgen.GenerateUUIDv7()`に変更
  - `GetTableName("dm_users", user.ID)` → `GetTableNameFromUUID("dm_users", user.ID)`に変更
  - `GetShardingConnectionByID(user.ID, "dm_users")` → `GetShardingConnectionByUUID(user.ID, "dm_users")`に変更
- すべてのメソッドでIDの型を`int64`から`string`に変更:
  - `GetByID(ctx context.Context, id int64)` → `GetByID(ctx context.Context, id string)`
  - `Update(ctx context.Context, id int64, req *model.UpdateDmUserRequest)` → `Update(ctx context.Context, id string, req *model.UpdateDmUserRequest)`
  - `Delete(ctx context.Context, id int64)` → `Delete(ctx context.Context, id string)`
- エラーハンドリングを追加（`GetTableNameFromUUID`と`GetShardingConnectionByUUID`のエラーを適切に処理）

**受け入れ基準**:
- `DmUserRepository.Create`メソッドが`GenerateUUIDv7()`を使用している
- `GetTableNameFromUUID`が使用されている
- `GetShardingConnectionByUUID`が使用されている
- すべてのメソッドでIDの型が`string`に変更されている
- エラーハンドリングが適切に実装されている

---

#### - [ ] タスク 5.2: DmUserRepositoryの変更（GORM版）
**目的**: DmUserRepository（GORM版）のID生成とシャーディングキー計算をUUIDv7ベースに変更する

**作業内容**:
- `server/internal/repository/dm_user_repository_gorm.go`を更新
- `Create`メソッド:
  - `idgen.GenerateSonyflakeID()` → `idgen.GenerateUUIDv7()`に変更
  - `GetTableName("dm_users", user.ID)` → `GetTableNameFromUUID("dm_users", user.ID)`に変更
  - `GetShardingConnectionByID(user.ID, "dm_users")` → `GetShardingConnectionByUUID(user.ID, "dm_users")`に変更
- すべてのメソッドでIDの型を`int64`から`string`に変更
- エラーハンドリングを追加

**受け入れ基準**:
- `DmUserRepositoryGORM.Create`メソッドが`GenerateUUIDv7()`を使用している
- `GetTableNameFromUUID`が使用されている
- `GetShardingConnectionByUUID`が使用されている
- すべてのメソッドでIDの型が`string`に変更されている
- エラーハンドリングが適切に実装されている

---

#### - [ ] タスク 5.3: DmPostRepositoryの変更（標準版）
**目的**: DmPostRepository（標準版）のID生成とシャーディングキー計算をUUIDv7ベースに変更する

**作業内容**:
- `server/internal/repository/dm_post_repository.go`を更新
- `Create`メソッド:
  - `idgen.GenerateSonyflakeID()` → `idgen.GenerateUUIDv7()`に変更
  - `GetTableName("dm_posts", req.UserID)` → `GetTableNameFromUUID("dm_posts", req.UserID)`に変更
  - `GetShardingConnectionByID(req.UserID, "dm_posts")` → `GetShardingConnectionByUUID(req.UserID, "dm_posts")`に変更
- すべてのメソッドでIDの型を`int64`から`string`に変更:
  - `GetByID(ctx context.Context, userID int64, postID int64)` → `GetByID(ctx context.Context, userID string, postID string)`
  - `GetByUserID(ctx context.Context, userID int64)` → `GetByUserID(ctx context.Context, userID string)`
  - `Update(ctx context.Context, userID int64, postID int64, req *model.UpdateDmPostRequest)` → `Update(ctx context.Context, userID string, postID string, req *model.UpdateDmPostRequest)`
  - `Delete(ctx context.Context, userID int64, postID int64)` → `Delete(ctx context.Context, userID string, postID string)`
- エラーハンドリングを追加

**受け入れ基準**:
- `DmPostRepository.Create`メソッドが`GenerateUUIDv7()`を使用している
- `GetTableNameFromUUID`が使用されている
- `GetShardingConnectionByUUID`が使用されている
- すべてのメソッドでIDの型が`string`に変更されている
- エラーハンドリングが適切に実装されている

---

#### - [ ] タスク 5.4: DmPostRepositoryの変更（GORM版）
**目的**: DmPostRepository（GORM版）のID生成とシャーディングキー計算をUUIDv7ベースに変更する

**作業内容**:
- `server/internal/repository/dm_post_repository_gorm.go`を更新
- `Create`メソッド:
  - `idgen.GenerateSonyflakeID()` → `idgen.GenerateUUIDv7()`に変更
  - `GetTableName("dm_posts", req.UserID)` → `GetTableNameFromUUID("dm_posts", req.UserID)`に変更
  - `GetShardingConnectionByID(req.UserID, "dm_posts")` → `GetShardingConnectionByUUID(req.UserID, "dm_posts")`に変更
- すべてのメソッドでIDの型を`int64`から`string`に変更
- エラーハンドリングを追加

**受け入れ基準**:
- `DmPostRepositoryGORM.Create`メソッドが`GenerateUUIDv7()`を使用している
- `GetTableNameFromUUID`が使用されている
- `GetShardingConnectionByUUID`が使用されている
- すべてのメソッドでIDの型が`string`に変更されている
- エラーハンドリングが適切に実装されている

---

### Phase 6: サービス層の変更

#### - [ ] タスク 6.1: DmUserServiceの変更
**目的**: DmUserServiceのIDの型を`int64`から`string`に変更する

**作業内容**:
- `server/internal/service/dm_user_service.go`を更新
- すべてのメソッドでIDの型を`int64`から`string`に変更:
  - `Create(ctx context.Context, req *model.CreateDmUserRequest)` → リポジトリ層の変更に合わせて型を更新
  - `GetByID(ctx context.Context, id int64)` → `GetByID(ctx context.Context, id string)`
  - `Update(ctx context.Context, id int64, req *model.UpdateDmUserRequest)` → `Update(ctx context.Context, id string, req *model.UpdateDmUserRequest)`
  - `Delete(ctx context.Context, id int64)` → `Delete(ctx context.Context, id string)`
- リポジトリ層の変更に合わせて型を更新

**受け入れ基準**:
- `DmUserService`のIDの型が`string`に変更されている
- すべてのメソッドで型が適切に更新されている

---

#### - [ ] タスク 6.2: DmPostServiceの変更
**目的**: DmPostServiceのIDの型を`int64`から`string`に変更する

**作業内容**:
- `server/internal/service/dm_post_service.go`を更新
- すべてのメソッドでIDの型を`int64`から`string`に変更:
  - `Create(ctx context.Context, req *model.CreateDmPostRequest)` → リポジトリ層の変更に合わせて型を更新
  - `GetByID(ctx context.Context, userID int64, postID int64)` → `GetByID(ctx context.Context, userID string, postID string)`
  - `GetByUserID(ctx context.Context, userID int64)` → `GetByUserID(ctx context.Context, userID string)`
  - `Update(ctx context.Context, userID int64, postID int64, req *model.UpdateDmPostRequest)` → `Update(ctx context.Context, userID string, postID string, req *model.UpdateDmPostRequest)`
  - `Delete(ctx context.Context, userID int64, postID int64)` → `Delete(ctx context.Context, userID string, postID string)`
- リポジトリ層の変更に合わせて型を更新

**受け入れ基準**:
- `DmPostService`のIDの型が`string`に変更されている
- すべてのメソッドで型が適切に更新されている

---

### Phase 7: API層の変更

#### - [ ] タスク 7.1: DmUserHandlerの変更
**目的**: DmUserHandlerのIDの型を`int64`から`string`に変更する

**作業内容**:
- `server/internal/api/handler/dm_user_handler.go`を更新
- すべてのメソッドでIDの型を`int64`から`string`に変更:
  - パスパラメータ、クエリパラメータの型を`string`に変更
  - バリデーションルールを更新（UUID文字列の形式チェック、必要に応じて`len=32`を追加）
- サービス層の変更に合わせて型を更新

**受け入れ基準**:
- ハンドラーのIDの型が`string`に変更されている
- パスパラメータ、クエリパラメータの型が`string`に変更されている
- バリデーションルールが適切に設定されている

---

#### - [ ] タスク 7.2: DmPostHandlerの変更
**目的**: DmPostHandlerのIDの型を`int64`から`string`に変更する

**作業内容**:
- `server/internal/api/handler/dm_post_handler.go`を更新
- すべてのメソッドでIDの型を`int64`から`string`に変更:
  - パスパラメータ、クエリパラメータの型を`string`に変更
  - バリデーションルールを更新（UUID文字列の形式チェック、必要に応じて`len=32`を追加）
- サービス層の変更に合わせて型を更新

**受け入れ基準**:
- ハンドラーのIDの型が`string`に変更されている
- パスパラメータ、クエリパラメータの型が`string`に変更されている
- バリデーションルールが適切に設定されている

---

#### - [ ] タスク 7.3: 出力モデルの変更
**目的**: API出力モデルのIDの型を`int64`から`string`に変更する

**作業内容**:
- `server/internal/api/huma/outputs.go`を更新
- `DmUserOutput.ID`: `int64` → `string`に変更
- `DmPostOutput.ID`: `int64` → `string`に変更
- `DmPostOutput.UserID`: `int64` → `string`に変更
- `DmUserPost`構造体のIDフィールドも`string`に変更:
  - `PostID int64` → `PostID string`
  - `UserID int64` → `UserID string`

**受け入れ基準**:
- `DmUserOutput.ID`が`string`型に変更されている
- `DmPostOutput.ID`が`string`型に変更されている
- `DmPostOutput.UserID`が`string`型に変更されている
- `DmUserPost`構造体のIDフィールドが`string`型に変更されている

---

### Phase 8: サンプルデータ生成コマンドの変更

#### - [ ] タスク 8.1: generate-sample-dataコマンドの変更
**目的**: サンプルデータ生成コマンドのID生成をUUIDv7に変更する

**作業内容**:
- `server/cmd/generate-sample-data/main.go`を更新
- `generateDmUsers`関数:
  - `idgen.GenerateSonyflakeID()` → `idgen.GenerateUUIDv7()`に変更
  - IDの型を`int64`から`string`に変更
  - シャーディングキーの計算をUUIDベースに変更（`GetTableNumberFromUUID`を使用）
- `generateDmPosts`関数:
  - `idgen.GenerateSonyflakeID()` → `idgen.GenerateUUIDv7()`に変更
  - IDの型を`int64`から`string`に変更
  - シャーディングキーの計算をUUIDベースに変更（`GetTableNumberFromUUID`を使用）
- エラーハンドリングを追加

**受け入れ基準**:
- `generateDmUsers`関数が`GenerateUUIDv7()`を使用している
- `generateDmPosts`関数が`GenerateUUIDv7()`を使用している
- IDの型が`string`に変更されている
- シャーディングキーの計算がUUIDベースに変更されている
- エラーハンドリングが適切に実装されている

---

### Phase 9: GoAdmin管理画面の変更

#### - [ ] タスク 9.1: dm_user_register.goの変更
**目的**: GoAdmin管理画面のユーザー登録機能をUUIDv7ベースに変更する

**作業内容**:
- `server/internal/admin/pages/dm_user_register.go`を更新
- `insertDmUserSharded`関数:
  - ID生成を`now.UnixNano()`から`idgen.GenerateUUIDv7()`に変更
  - 戻り値の型を`int64`から`string`に変更
  - テーブル番号の計算をUUIDベースに変更（`GetTableNumberFromUUID`を使用）
  - エラーハンドリングを追加
- `handleDmUserRegisterPost`関数:
  - リダイレクトURLのクエリパラメータでIDを`string`型として渡すように変更（`url.QueryEscape`を使用）
- `checkEmailExistsSharded`関数: 変更不要（メールアドレスの検索のみ）

**受け入れ基準**:
- `insertDmUserSharded`関数が`GenerateUUIDv7()`を使用している
- `insertDmUserSharded`関数の戻り値の型が`string`に変更されている
- テーブル番号の計算がUUIDベースに変更されている（`GetTableNumberFromUUID`を使用）
- リダイレクトURLのクエリパラメータでIDが`string`型として渡されている
- エラーハンドリングが適切に実装されている

---

### Phase 10: クライアント側の変更

#### - [ ] タスク 10.1: APIレスポンス型定義の更新
**目的**: クライアント側のAPIレスポンス型定義を更新する

**作業内容**:
- クライアント側のAPIレスポンス型定義ファイルを確認
- `DmUser`型の`id`フィールド: `number` → `string`に変更
- `DmPost`型の`id`フィールド: `number` → `string`に変更
- `DmPost`型の`user_id`フィールド: `number` → `string`に変更

**受け入れ基準**:
- APIレスポンス型定義が更新されている（IDが`string`型）
- 型定義ファイルが正しく更新されている

---

#### - [ ] タスク 10.2: APIクライアントコードの更新
**目的**: APIクライアントコードでIDを`string`型として処理できるように更新する

**作業内容**:
- クライアント側のAPIクライアントコードを確認
- IDを`string`型として扱うように変更
- 必要に応じて、IDのバリデーション（32文字のUUID文字列）を追加

**受け入れ基準**:
- APIクライアントコードが`string`型のIDを正しく処理できること
- IDのバリデーションが適切に実装されている（必要に応じて）

---

#### - [ ] タスク 10.3: ID表示・入力フォームの対応
**目的**: IDの表示・入力フォームが`string`型に対応していることを確認する

**作業内容**:
- クライアント側のID表示・入力フォームを確認
- `string`型のIDを正しく表示・入力できることを確認
- 必要に応じて、フォームのバリデーションを追加

**受け入れ基準**:
- IDの表示・入力フォームが`string`型に対応していること
- フォームのバリデーションが適切に実装されている（必要に応じて）

---

### Phase 11: Sonyflake関数の削除

#### - [ ] タスク 11.1: sonyflake.goの削除
**目的**: Sonyflake関数の実装ファイルを削除する

**作業内容**:
- `server/internal/util/idgen/sonyflake.go`を削除
- コード内でSonyflake関数への参照が全て削除されていることを確認

**受け入れ基準**:
- `server/internal/util/idgen/sonyflake.go`が削除されている
- コード内でSonyflake関数への参照が全て削除されている

---

#### - [ ] タスク 11.2: sonyflake_test.goの削除
**目的**: Sonyflake関数のテストファイルを削除する

**作業内容**:
- `server/internal/util/idgen/sonyflake_test.go`を削除

**受け入れ基準**:
- `server/internal/util/idgen/sonyflake_test.go`が削除されている

---

#### - [ ] タスク 11.3: 依存関係の削除
**目的**: `go.mod`からSonyflakeの依存関係を削除する

**作業内容**:
- `go.mod`から`github.com/sony/sonyflake`の依存関係を削除
- `go mod tidy`を実行して依存関係を整理
- `go.sum`からもSonyflakeの依存関係が削除されていることを確認

**受け入れ基準**:
- `go.mod`から`github.com/sony/sonyflake`の依存関係が削除されている
- `go mod tidy`が実行されている
- `go.sum`からもSonyflakeの依存関係が削除されている

---

### Phase 12: マイグレーションの作成

#### - [ ] タスク 12.1: Atlasマイグレーションファイルの作成
**目的**: データベーススキーマの変更をマイグレーションファイルとして作成する

**作業内容**:
- Atlasマイグレーションファイルを作成
- カラム型の変更:
  - `dm_users.id`: `bigint unsigned` → `varchar(32)`
  - `dm_posts.id`: `bigint unsigned` → `varchar(32)`
  - `dm_posts.user_id`: `bigint` → `varchar(32)`
- 既存データの削除:
  - テーブルを再作成するか、既存データを削除する
- インデックスの再作成:
  - 型変更に伴いインデックスを再作成する

**受け入れ基準**:
- Atlasマイグレーションファイルが作成されている
- `dm_users.id`が`varchar(32)`に変更されている
- `dm_posts.id`が`varchar(32)`に変更されている
- `dm_posts.user_id`が`varchar(32)`に変更されている
- 既存データが削除される設定になっている
- インデックスが再作成される設定になっている

---

#### - [ ] タスク 12.2: マイグレーションの実行確認
**目的**: マイグレーションが正常に実行できることを確認する

**作業内容**:
- マイグレーションを実行
- マイグレーションが正常に完了することを確認
- 既存データが削除されていることを確認
- カラム型が正しく変更されていることを確認
- インデックスが再作成されていることを確認

**受け入れ基準**:
- マイグレーションが正常に実行できること
- 既存データが削除されていること（マイグレーション後に既存データが存在しないこと）
- カラム型が正しく変更されていること
- インデックスが再作成されていること

---

### Phase 13: テストの実装

#### - [ ] タスク 13.1: リポジトリ層の統合テスト更新
**目的**: リポジトリ層の統合テストをUUIDv7ベースに更新する

**作業内容**:
- `server/internal/repository/dm_user_repository_test.go`を更新:
  - IDの型を`int64`から`string`に変更
  - UUIDv7で生成されたIDを使用するように変更
  - シャーディングキー計算が正しく動作することを確認
- `server/internal/repository/dm_post_repository_test.go`を更新:
  - IDの型を`int64`から`string`に変更
  - UUIDv7で生成されたIDを使用するように変更
  - シャーディングキー計算が正しく動作することを確認

**受け入れ基準**:
- リポジトリ層の統合テストが更新されている
- IDの型が`string`に変更されている
- UUIDv7で生成されたIDを使用するテストが実装されている
- シャーディングキー計算が正しく動作することを確認するテストが実装されている
- テストが全て通過する

---

#### - [ ] タスク 13.2: 既存テストの確認
**目的**: 既存のテストが全て通過することを確認する

**作業内容**:
- 既存のテストを実行
- テストが全て通過することを確認
- 失敗したテストがあれば修正する

**受け入れ基準**:
- 既存のテストが全て通過すること
- 失敗したテストがないこと

---

#### - [ ] タスク 13.3: 新規テストの確認
**目的**: 新規作成したテストが全て通過することを確認する

**作業内容**:
- 新規作成したテストを実行:
  - `uuidv7_test.go`
  - `sharding_test.go`の新規関数のテスト
- テストが全て通過することを確認

**受け入れ基準**:
- 新規作成したテストが全て通過すること
- 失敗したテストがないこと

---

### Phase 14: ドキュメントの更新

#### - [ ] タスク 14.1: Architecture.mdの更新
**目的**: アーキテクチャドキュメントにUUIDv7の使用を記載し、Sonyflakeの削除を記載する

**作業内容**:
- `docs/Architecture.md`を更新
- UUIDv7の使用を記載:
  - ID生成方式としてUUIDv7を使用することを記載
  - `idgen.GenerateUUIDv7()`関数の使用方法を記載
  - UUIDv7の仕様（ハイフン抜き小文字32文字）を記載
- Sonyflakeの削除を記載:
  - Sonyflakeライブラリが削除されたことを記載
  - 削除理由を記載

**受け入れ基準**:
- `docs/Architecture.md`が更新されている
- UUIDv7の使用が記載されている
- Sonyflakeの削除が記載されている

---

#### - [ ] タスク 14.2: sharding.goのコメント更新
**目的**: シャーディング規則のドキュメントを更新する

**作業内容**:
- `server/internal/db/sharding.go`のコメントを更新
- UUIDv7の使用を記載:
  - ID生成方式としてUUIDv7を使用することを記載
  - シャーディングキーの計算方法（UUIDの後ろ2文字から16進数として解釈）を記載
- Sonyflakeの記述を削除:
  - Sonyflakeに関する記述を削除

**受け入れ基準**:
- `server/internal/db/sharding.go`のコメントが更新されている
- UUIDv7の使用が記載されている
- Sonyflakeの記述が削除されている

---

#### - [ ] タスク 14.3: コード内のコメント更新
**目的**: コード内のコメントを適切に更新する

**作業内容**:
- 変更したファイルのコメントを確認
- 必要に応じてコメントを更新:
  - ID生成方式の変更を反映
  - シャーディングキー計算方法の変更を反映
  - 型の変更を反映

**受け入れ基準**:
- コード内のコメントが適切に更新されている
- 変更内容がコメントに反映されている

---

## 実装順序の推奨

1. **Phase 1-2**: インフラストラクチャの構築（UUIDv7ライブラリの導入、ID生成ユーティリティの実装）
2. **Phase 3**: シャーディングキー計算関数の実装（新規関数の追加）
3. **Phase 4**: モデル定義の変更（型の変更）
4. **Phase 5-7**: アプリケーション層の変更（リポジトリ層、サービス層、API層）
5. **Phase 8-9**: サンプルデータ生成コマンドとGoAdmin管理画面の変更
6. **Phase 10**: クライアント側の変更
7. **Phase 11**: Sonyflake関数の削除（他の変更が完了してから）
8. **Phase 12**: マイグレーションの作成と実行
9. **Phase 13**: テストの実装と確認
10. **Phase 14**: ドキュメントの更新

## 注意事項

- **型変換の一貫性**: IDの型を`int64`から`string`に変更する際、すべての関連箇所で一貫して変更する必要がある
- **エラーハンドリング**: 新規関数（`GetTableNumberFromUUID`、`GetTableNameFromUUID`、`GetShardingConnectionByUUID`）のエラーハンドリングを適切に実装する
- **後方互換性**: 既存の`GetTableName(baseName string, id int64)`関数と`GetShardingConnectionByID(id int64, tableName string)`関数は残す
- **既存データの破棄**: 既存データは破棄するため、マイグレーション時に削除する
- **Sonyflake関数の削除**: 他の変更が完了してからSonyflake関数を削除する（依存関係の確認のため）
- **テストの実行**: 各フェーズの実装後、該当するテストを実行して動作確認を行う
