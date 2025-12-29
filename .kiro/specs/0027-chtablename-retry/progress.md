# テーブル名変更機能修正実装進捗（リトライ）

## 現在のステータス
**最終更新**: 2025-01-27
**全体進捗**: 未着手（0/26タスク）

## テスト結果
- **最新テスト実行**: 未実行
- **コンパイルエラー**: 未確認
- **リンターエラー**: 未確認

---

## Phase 1: CLIツールの修正

### タスク 1.1: server/cmd/generate-sample-data/main.go の修正
- **ステータス**: 未着手
- **対象ファイル**: `server/cmd/generate-sample-data/main.go`
- **作業内容**:
  - 関数名の変更（7箇所: generateUsers, insertUsersBatch, fetchUserIDs, generatePosts, insertPostsBatch, generateNews, insertNewsBatch）
  - 変数名の変更（複数箇所: userIDs, users, user, posts, post, news）
  - コメントの更新
- **備考**: 

---

### タスク 1.2: server/cmd/server/main.go の修正
- **ステータス**: 未着手
- **対象ファイル**: `server/cmd/server/main.go`
- **作業内容**:
  - 変数名の変更（6箇所: userRepo, userService, userHandler, postRepo, postService, postHandler）
  - 変数参照の更新（6箇所）
- **備考**: 

---

### タスク 1.3: server/cmd/admin/main.go の修正
- **ステータス**: 未着手
- **対象ファイル**: `server/cmd/admin/main.go`
- **作業内容**:
  - 関数呼び出しの変更（2箇所: UserRegisterPage, UserRegisterCompletePage）
- **備考**: 関数定義のリネームが必要か確認（必要に応じてPhase 4で対応）

---

### タスク 1.4: server/cmd/list-dm-users/main.go の修正
- **ステータス**: 未着手
- **対象ファイル**: `server/cmd/list-dm-users/main.go`
- **作業内容**:
  - 関数名の変更（1箇所: printUsersTSV）
  - コメントの更新
- **備考**: 

---

### タスク 1.5: server/cmd/list-dm-users/main_test.go の修正
- **ステータス**: 未着手
- **対象ファイル**: `server/cmd/list-dm-users/main_test.go`
- **作業内容**:
  - 関数呼び出しの変更（2箇所）
  - テスト関数名の変更（2箇所: TestPrintUsersTSV, TestPrintUsersTSV_RFC3339Format）
  - 変数名の変更（複数箇所: users, user）
- **備考**: 

---

## Phase 2: Repository層の修正

### タスク 2.1: server/internal/repository/dm_post_repository_gorm.go の修正
- **ステータス**: 未着手
- **対象ファイル**: `server/internal/repository/dm_post_repository_gorm.go`
- **作業内容**:
  - 変数名の変更（1箇所: tableUserPosts）
- **備考**: 

---

### タスク 2.2: server/internal/repository/interfaces.go の修正
- **ステータス**: 未着手
- **対象ファイル**: `server/internal/repository/interfaces.go`
- **作業内容**:
  - コメントの更新
- **備考**: 

---

### タスク 2.3: server/internal/repository/dm_post_repository.go の修正
- **ステータス**: 未着手
- **対象ファイル**: `server/internal/repository/dm_post_repository.go`
- **作業内容**:
  - コメントの更新
- **備考**: 

---

### タスク 2.4: server/internal/repository/dm_user_repository_gorm_test.go の修正
- **ステータス**: 未着手
- **対象ファイル**: `server/internal/repository/dm_user_repository_gorm_test.go`
- **作業内容**:
  - テスト関数名の変更（6箇所）
- **備考**: 

---

### タスク 2.5: server/internal/repository/dm_user_repository_test.go の修正
- **ステータス**: 未着手
- **対象ファイル**: `server/internal/repository/dm_user_repository_test.go`
- **作業内容**:
  - テスト関数名の変更（6箇所）
  - 変数名の変更（1箇所: users）
- **備考**: 

---

## Phase 3: テストコードの修正

### タスク 3.1: server/test/integration/dm_user_flow_gorm_test.go の修正
- **ステータス**: 未着手
- **対象ファイル**: `server/test/integration/dm_user_flow_gorm_test.go`
- **作業内容**:
  - テスト関数名の変更（2箇所: TestUserCRUDFlowGORM, TestUserCrossShardOperationsGORM）
  - コメントの更新
- **備考**: 

---

### タスク 3.2: server/test/integration/dm_user_flow_test.go の修正
- **ステータス**: 未着手
- **対象ファイル**: `server/test/integration/dm_user_flow_test.go`
- **作業内容**:
  - テスト関数名の変更（2箇所: TestUserCRUDFlow, TestUserCrossShardOperations）
  - コメントの更新
- **備考**: 

---

### タスク 3.3: server/test/integration/dm_post_flow_test.go の修正
- **ステータス**: 未着手
- **対象ファイル**: `server/test/integration/dm_post_flow_test.go`
- **作業内容**:
  - テスト関数名の変更（2箇所: TestPostCRUDFlow, TestCrossShardJoin）
  - コメントの更新
- **備考**: 

---

### タスク 3.4: server/test/integration/dm_post_flow_gorm_test.go の修正
- **ステータス**: 未着手
- **対象ファイル**: `server/test/integration/dm_post_flow_gorm_test.go`
- **作業内容**:
  - テスト関数名の変更（2箇所: TestPostCRUDFlowGORM, TestCrossShardJoinGORM）
  - コメントの更新
- **備考**: 

---

### タスク 3.5: server/test/integration/sharding_test.go の修正
- **ステータス**: 未着手
- **対象ファイル**: `server/test/integration/sharding_test.go`
- **作業内容**:
  - テスト関数名の変更（1箇所: TestCrossTableQueryUsers）
  - 変数名の変更（複数箇所: testUsers, u）
  - コメントの更新
- **備考**: 

---

### タスク 3.6: server/test/fixtures/dm_users.go の修正
- **ステータス**: 未着手
- **対象ファイル**: `server/test/fixtures/dm_users.go`
- **作業内容**:
  - 関数名の変更（3箇所: CreateTestUser, CreateTestUserWithEmail, CreateMultipleTestUsers）
  - コメントの更新
  - この関数を使用している全てのテストコードの更新確認
- **備考**: 他のテストファイルで使用されている可能性があるため、慎重に確認

---

### タスク 3.7: server/test/fixtures/dm_posts.go の修正
- **ステータス**: 未着手
- **対象ファイル**: `server/test/fixtures/dm_posts.go`
- **作業内容**:
  - 関数名の変更（3箇所: CreateTestPost, CreateTestPostWithContent, CreateMultipleTestPosts）
  - 変数名の変更（複数箇所: post, posts）
  - コメントの更新
  - この関数を使用している全てのテストコードの更新確認
- **備考**: 他のテストファイルで使用されている可能性があるため、慎重に確認

---

### タスク 3.8: server/test/e2e/api_test.go の修正
- **ステータス**: 未着手
- **対象ファイル**: `server/test/e2e/api_test.go`
- **作業内容**:
  - テスト関数名の変更（3箇所: TestUserAPI_CreateAndRetrieve, TestUserAPI_UpdateAndDelete, TestPostAPI_CompleteFlow）
  - 変数名の変更（複数箇所: user, post）
  - コメントの更新
- **備考**: 

---

## Phase 4: 管理画面の修正

### タスク 4.1: server/internal/admin/pages/dm_user_register.go の修正
- **ステータス**: 未着手
- **対象ファイル**: `server/internal/admin/pages/dm_user_register.go`
- **作業内容**:
  - 関数定義の変更（1箇所: UserRegisterPage → DmUserRegisterPage、必要に応じて）
  - この関数を呼び出している全ての箇所の更新確認
- **備考**: 必要に応じて実施

---

### タスク 4.2: server/internal/admin/pages/dm_user_register_complete.go の修正
- **ステータス**: 未着手
- **対象ファイル**: `server/internal/admin/pages/dm_user_register_complete.go`
- **作業内容**:
  - 関数定義の変更（1箇所: UserRegisterCompletePage → DmUserRegisterCompletePage、必要に応じて）
  - この関数を呼び出している全ての箇所の更新確認
- **備考**: 必要に応じて実施

---

### タスク 4.3: server/internal/admin/pages/pages.go の修正
- **ステータス**: 未着手
- **対象ファイル**: `server/internal/admin/pages/pages.go`
- **作業内容**:
  - 関数参照の変更（1箇所: UserRegisterCompletePage → DmUserRegisterCompletePage、必要に応じて）
- **備考**: 必要に応じて実施

---

### タスク 4.4: server/internal/admin/sharding.go の修正
- **ステータス**: 未着手
- **対象ファイル**: `server/internal/admin/sharding.go`
- **作業内容**:
  - 関数名の変更（4箇所: FindUserAcrossShards, CountUsersAcrossShards, FindPostAcrossShards, CountPostsAcrossShards）
  - コメントの更新
  - この関数を使用している全ての箇所の更新確認
- **備考**: 他のファイルで使用されている可能性があるため、慎重に確認

---

## Phase 5: 最終検証

### タスク 5.1: 全テストの実行
- **ステータス**: 未着手
- **作業内容**:
  - 全テストの実行（`go test ./...`）
  - テスト結果の確認
  - 失敗したテストの修正（該当する場合）
- **備考**: 

---

### タスク 5.2: 旧命名規則の残存確認
- **ステータス**: 未着手
- **作業内容**:
  - 旧命名規則の検索（grepコマンド）
  - 検索結果の確認
  - 残存している旧命名規則の修正（該当する場合）
- **備考**: 

---

### タスク 5.3: コンパイルエラーの確認
- **ステータス**: 未着手
- **作業内容**:
  - コンパイルの実行（`go build ./...`）
  - コンパイルエラーの確認
  - エラーの修正（該当する場合）
- **備考**: 

---

### タスク 5.4: リンターエラーの確認
- **ステータス**: 未着手
- **作業内容**:
  - リンターの実行（`go vet ./...`）
  - リンターエラーの確認
  - エラーの修正（該当する場合）
- **備考**: 

---

## エラー・修正履歴

（実装中に発生したエラーや修正内容を記録）

---

## 次のアクション

Phase 1のタスク 1.1から順次実装を開始する。

### 実装サマリー
- **総タスク数**: 26タスク
- **完了タスク数**: 0タスク
- **進捗率**: 0%
