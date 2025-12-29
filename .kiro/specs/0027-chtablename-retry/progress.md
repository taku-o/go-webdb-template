# テーブル名変更機能修正実装進捗（リトライ）

## 現在のステータス
**最終更新**: 2025-12-30
**全体進捗**: 全タスク実装済み（26/26タスク） - ユーザー確認待ち

## テスト結果
- **最新テスト実行**: 全テスト成功
- **コンパイルエラー**: なし
- **リンターエラー**: golangci-lint未インストールのため未確認

---

## Phase 1: CLIツールの修正

### タスク 1.1: server/cmd/generate-sample-data/main.go の修正
- **ステータス**: 実装済み
- **対象ファイル**: `server/cmd/generate-sample-data/main.go`
- **作業内容**:
  - 関数名の変更（7箇所: generateUsers → generateDmUsers, insertUsersBatch → insertDmUsersBatch, fetchUserIDs → fetchDmUserIDs, generatePosts → generateDmPosts, insertPostsBatch → insertDmPostsBatch）
  - 変数名の変更（複数箇所: userIDs → dmUserIDs, users → dmUsers, user → dmUser, posts → dmPosts, post → dmPost）
  - コメントの更新
- **備考**: 前回のセッションで実装済み

---

### タスク 1.2: server/cmd/server/main.go の修正
- **ステータス**: 実装済み
- **対象ファイル**: `server/cmd/server/main.go`
- **作業内容**:
  - 変数名の変更（6箇所: userRepo → dmUserRepo, userService → dmUserService, userHandler → dmUserHandler, postRepo → dmPostRepo, postService → dmPostService, postHandler → dmPostHandler）
  - 変数参照の更新（6箇所）
- **備考**: 前回のセッションで実装済み

---

### タスク 1.3: server/cmd/admin/main.go の修正
- **ステータス**: 実装済み
- **対象ファイル**: `server/cmd/admin/main.go`
- **作業内容**:
  - 関数呼び出しの変更（2箇所: UserRegisterPage → DmUserRegisterPage, UserRegisterCompletePage → DmUserRegisterCompletePage）
- **備考**: 前回のセッションで実装済み

---

### タスク 1.4: server/cmd/list-dm-users/main.go の修正
- **ステータス**: 実装済み
- **対象ファイル**: `server/cmd/list-dm-users/main.go`
- **作業内容**:
  - 関数名の変更（1箇所: printUsersTSV → printDmUsersTSV）
  - コメントの更新
- **備考**: 前回のセッションで実装済み

---

### タスク 1.5: server/cmd/list-dm-users/main_test.go の修正
- **ステータス**: 実装済み
- **対象ファイル**: `server/cmd/list-dm-users/main_test.go`
- **作業内容**:
  - 関数呼び出しの変更（2箇所）
  - テスト関数名の変更（2箇所: TestPrintUsersTSV → TestPrintDmUsersTSV, TestPrintUsersTSV_RFC3339Format → TestPrintDmUsersTSV_RFC3339Format）
  - 変数名の変更（複数箇所: users → dmUsers, user → dmUser）
- **備考**: 前回のセッションで実装済み

---

## Phase 2: Repository層の修正

### タスク 2.1: server/internal/repository/dm_post_repository_gorm.go の修正
- **ステータス**: 実装済み
- **対象ファイル**: `server/internal/repository/dm_post_repository_gorm.go`
- **作業内容**:
  - 変数名の変更（1箇所: tableUserPosts → tableDmUserPosts）
- **備考**: 前回のセッションで実装済み

---

### タスク 2.2: server/internal/repository/interfaces.go の修正
- **ステータス**: 実装済み
- **対象ファイル**: `server/internal/repository/interfaces.go`
- **作業内容**:
  - コメントの更新
- **備考**: 前回のセッションで実装済み

---

### タスク 2.3: server/internal/repository/dm_post_repository.go の修正
- **ステータス**: 実装済み
- **対象ファイル**: `server/internal/repository/dm_post_repository.go`
- **作業内容**:
  - コメントの更新
- **備考**: 前回のセッションで実装済み

---

### タスク 2.4: server/internal/repository/dm_user_repository_gorm_test.go の修正
- **ステータス**: 実装済み
- **対象ファイル**: `server/internal/repository/dm_user_repository_gorm_test.go`
- **作業内容**:
  - テスト関数名の変更（6箇所: TestUserRepositoryGORM_* → TestDmUserRepositoryGORM_*）
- **備考**: 前回のセッションで実装済み

---

### タスク 2.5: server/internal/repository/dm_user_repository_test.go の修正
- **ステータス**: 実装済み
- **対象ファイル**: `server/internal/repository/dm_user_repository_test.go`
- **作業内容**:
  - テスト関数名の変更（6箇所: TestUserRepository_* → TestDmUserRepository_*）
  - 変数名の変更（複数箇所: user → dmUser, users → dmUsers）
- **備考**: 前回のセッションで実装済み

---

## Phase 3: テストコードの修正

### タスク 3.1: server/test/integration/dm_user_flow_gorm_test.go の修正
- **ステータス**: 実装済み
- **対象ファイル**: `server/test/integration/dm_user_flow_gorm_test.go`
- **作業内容**:
  - テスト関数名の変更（2箇所: TestUserCRUDFlowGORM → TestDmUserCRUDFlowGORM, TestUserCrossShardOperationsGORM → TestDmUserCrossShardOperationsGORM）
  - コメントの更新
- **備考**: 前回のセッションで実装済み

---

### タスク 3.2: server/test/integration/dm_user_flow_test.go の修正
- **ステータス**: 実装済み
- **対象ファイル**: `server/test/integration/dm_user_flow_test.go`
- **作業内容**:
  - テスト関数名の変更（2箇所: TestUserCRUDFlow → TestDmUserCRUDFlow, TestUserCrossShardOperations → TestDmUserCrossShardOperations）
  - コメントの更新
- **備考**: 前回のセッションで実装済み

---

### タスク 3.3: server/test/integration/dm_post_flow_test.go の修正
- **ステータス**: 実装済み
- **対象ファイル**: `server/test/integration/dm_post_flow_test.go`
- **作業内容**:
  - テスト関数名の変更（2箇所: TestPostCRUDFlow → TestDmPostCRUDFlow, TestCrossShardJoin → TestDmPostCrossShardJoin）
  - フィクスチャ関数呼び出しの変更（fixtures.CreateTestUser → fixtures.CreateTestDmUser, fixtures.CreateTestPost → fixtures.CreateTestDmPost）
  - コメントの更新
- **備考**: 前回のセッションで実装済み

---

### タスク 3.4: server/test/integration/dm_post_flow_gorm_test.go の修正
- **ステータス**: 実装済み
- **対象ファイル**: `server/test/integration/dm_post_flow_gorm_test.go`
- **作業内容**:
  - テスト関数名の変更（2箇所: TestPostCRUDFlowGORM → TestDmPostCRUDFlowGORM, TestCrossShardJoinGORM → TestDmPostCrossShardJoinGORM）
  - コメントの更新
- **備考**: 前回のセッションで実装済み

---

### タスク 3.5: server/test/integration/sharding_test.go の修正
- **ステータス**: 実装済み（変更不要）
- **対象ファイル**: `server/test/integration/sharding_test.go`
- **作業内容**:
  - 確認の結果、既に正しい命名規則に従っていたため、変更不要
- **備考**: 前回のセッションで確認済み

---

### タスク 3.6: server/test/fixtures/dm_users.go の修正
- **ステータス**: 実装済み
- **対象ファイル**: `server/test/fixtures/dm_users.go`
- **作業内容**:
  - 関数名の変更（3箇所: CreateTestUser → CreateTestDmUser, CreateTestUserWithEmail → CreateTestDmUserWithEmail, CreateMultipleTestUsers → CreateMultipleTestDmUsers）
  - コメントの更新
- **備考**: 前回のセッションで実装済み

---

### タスク 3.7: server/test/fixtures/dm_posts.go の修正
- **ステータス**: 実装済み
- **対象ファイル**: `server/test/fixtures/dm_posts.go`
- **作業内容**:
  - 関数名の変更（3箇所: CreateTestPost → CreateTestDmPost, CreateTestPostWithContent → CreateTestDmPostWithContent, CreateMultipleTestPosts → CreateMultipleTestDmPosts）
  - コメントの更新
- **備考**: 前回のセッションで実装済み

---

### タスク 3.8: server/test/e2e/api_test.go の修正
- **ステータス**: 実装済み
- **対象ファイル**: `server/test/e2e/api_test.go`
- **作業内容**:
  - テスト関数名の変更（3箇所: TestUserAPI_CreateAndRetrieve → TestDmUserAPI_CreateAndRetrieve, TestUserAPI_UpdateAndDelete → TestDmUserAPI_UpdateAndDelete, TestPostAPI_CompleteFlow → TestDmPostAPI_CompleteFlow）
  - 変数名の変更（複数箇所: user → dmUser, userID → dmUserID, post → dmPost, postID → dmPostID, userPosts → dmUserPosts）
- **備考**: 前回のセッションで実装済み

---

## Phase 4: 管理画面の修正

### タスク 4.1: server/internal/admin/pages/dm_user_register.go の修正
- **ステータス**: 実装済み
- **対象ファイル**: `server/internal/admin/pages/dm_user_register.go`
- **作業内容**:
  - 関数名の変更（5箇所: UserRegisterPage → DmUserRegisterPage, handleUserRegisterPost → handleDmUserRegisterPost, validateUserInput → validateDmUserInput, insertUserSharded → insertDmUserSharded, renderUserRegisterForm → renderDmUserRegisterForm）
  - 変数名の変更（複数箇所: userID → dmUserID）
- **備考**: 今回のセッションで実装済み

---

### タスク 4.2: server/internal/admin/pages/dm_user_register_complete.go の修正
- **ステータス**: 実装済み
- **対象ファイル**: `server/internal/admin/pages/dm_user_register_complete.go`
- **作業内容**:
  - 関数名の変更（1箇所: UserRegisterCompletePage → DmUserRegisterCompletePage）
  - 変数名の変更（3箇所: userID → dmUserID, userName → dmUserName, userEmail → dmUserEmail）
- **備考**: 今回のセッションで実装済み

---

### タスク 4.3: server/internal/admin/pages/pages.go の修正
- **ステータス**: 実装済み
- **対象ファイル**: `server/internal/admin/pages/pages.go`
- **作業内容**:
  - 関数参照の変更（1箇所: UserRegisterCompletePage → DmUserRegisterCompletePage）
- **備考**: 今回のセッションで実装済み

---

### タスク 4.4: server/internal/admin/sharding.go の修正
- **ステータス**: 実装済み
- **対象ファイル**: `server/internal/admin/sharding.go`
- **作業内容**:
  - 関数名の変更（6箇所: FindUserAcrossShards → FindDmUserAcrossShards, FindPostAcrossShards → FindDmPostAcrossShards, CountUsersAcrossShards → CountDmUsersAcrossShards, CountPostsAcrossShards → CountDmPostsAcrossShards, GetShardForUserID → GetShardForDmUserID）
  - 構造体フィールド名の変更（ShardStats: UserCount → DmUserCount, PostCount → DmPostCount）
  - 変数名の変更（userCount → dmUserCount, postCount → dmPostCount）
- **備考**: 今回のセッションで実装済み

---

## Phase 5: 最終検証

### タスク 5.1: 全テストの実行
- **ステータス**: 実装済み
- **作業内容**:
  - 全テストの実行（`go test ./...`）
  - テスト結果: 全テスト成功
- **備考**: 今回のセッションで確認済み

---

### タスク 5.2: 旧命名規則の残存確認
- **ステータス**: 実装済み
- **作業内容**:
  - 旧命名規則の検索（grepコマンド）
  - 追加発見: dm_post_repository_test.go と dm_post_repository_gorm_test.go のテスト関数名に修正漏れがあり、修正済み
    - TestPostRepository_* → TestDmPostRepository_*
    - TestPostRepositoryGORM_* → TestDmPostRepositoryGORM_*
- **備考**: 今回のセッションで確認・修正済み

---

### タスク 5.3: コンパイルエラーの確認
- **ステータス**: 実装済み
- **作業内容**:
  - コンパイルの実行（`go build ./...`）
  - コンパイルエラー: なし
- **備考**: 今回のセッションで確認済み

---

### タスク 5.4: リンターエラーの確認
- **ステータス**: 実装済み
- **作業内容**:
  - golangci-lint 未インストールのため未確認
  - go vet は go test で自動的に実行され、エラーなし
- **備考**: 今回のセッションで確認済み

---

## エラー・修正履歴

### 2025-12-30
- Phase 5.2 の旧命名規則残存確認で追加の修正漏れを発見・修正:
  - `server/internal/repository/dm_post_repository_test.go`: テスト関数名 6件
  - `server/internal/repository/dm_post_repository_gorm_test.go`: テスト関数名 8件

---

## Phase 6: 追加命名規則修正（ユーザー要望）

### タスク 6.1: server/internal/admin/tables.go の修正
- **ステータス**: 未着手
- **対象ファイル**: `server/internal/admin/tables.go`
- **作業内容**:
  - 関数名の変更: `GetNewsTable` → `GetDmNewsTable`
  - マップキー参照の変更: `"dm-news": GetNewsTable` → `"dm-news": GetDmNewsTable`

---

### タスク 6.2: server/internal/model/dm_news.go の修正
- **ステータス**: 未着手
- **対象ファイル**: `server/internal/model/dm_news.go`
- **作業内容**:
  - 型名の変更: `CreateNewsRequest` → `CreateDmNewsRequest`
  - 型名の変更: `UpdateNewsRequest` → `UpdateDmNewsRequest`

---

### タスク 6.3: server/internal/repository/dm_post_repository.go の変数名修正
- **ステータス**: 未着手
- **対象ファイル**: `server/internal/repository/dm_post_repository.go`
- **作業内容**:
  - 変数名の変更: `userPosts` → `dmUserPosts`（GetUserPosts関数内）

---

### タスク 6.4: server/internal/api/huma/inputs.go の型名修正
- **ステータス**: 未着手
- **対象ファイル**: `server/internal/api/huma/inputs.go`
- **作業内容**:
  - `CreateUserInput` → `CreateDmUserInput`
  - `GetUserInput` → `GetDmUserInput`
  - `ListUsersInput` → `ListDmUsersInput`
  - `UpdateUserInput` → `UpdateDmUserInput`
  - `DeleteUserInput` → `DeleteDmUserInput`
  - `CreatePostInput` → `CreateDmPostInput`
  - `GetPostInput` → `GetDmPostInput`
  - `ListPostsInput` → `ListDmPostsInput`
  - `UpdatePostInput` → `UpdateDmPostInput`
  - `DeletePostInput` → `DeleteDmPostInput`
  - `GetUserPostsInput` → `GetDmUserPostsInput`

---

### タスク 6.5: server/internal/api/huma/outputs.go の型名修正
- **ステータス**: 未着手
- **対象ファイル**: `server/internal/api/huma/outputs.go`
- **作業内容**:
  - `UserOutput` → `DmUserOutput`
  - `UsersOutput` → `DmUsersOutput`
  - `DeleteUserOutput` → `DeleteDmUserOutput`
  - `PostOutput` → `DmPostOutput`
  - `PostsOutput` → `DmPostsOutput`
  - `DeletePostOutput` → `DeleteDmPostOutput`
  - `UserPostsOutput` → `DmUserPostsOutput`

---

### タスク 6.6: server/internal/api/handler/dm_user_handler.go の修正
- **ステータス**: 未着手
- **対象ファイル**: `server/internal/api/handler/dm_user_handler.go`
- **作業内容**:
  - 関数名の変更: `RegisterUserEndpoints` → `RegisterDmUserEndpoints`
  - 型参照の更新（humaapi.CreateUserInput → humaapi.CreateDmUserInput 等）

---

### タスク 6.7: server/internal/api/handler/dm_post_handler.go の修正
- **ステータス**: 未着手
- **対象ファイル**: `server/internal/api/handler/dm_post_handler.go`
- **作業内容**:
  - 関数名の変更: `RegisterPostEndpoints` → `RegisterDmPostEndpoints`
  - 型参照の更新（humaapi.CreatePostInput → humaapi.CreateDmPostInput 等）

---

### タスク 6.8: server/internal/api/router/router.go の修正
- **ステータス**: 未着手
- **対象ファイル**: `server/internal/api/router/router.go`
- **作業内容**:
  - パラメータ名の変更: `userHandler` → `dmUserHandler`
  - パラメータ名の変更: `postHandler` → `dmPostHandler`
  - 関数呼び出しの変更: `RegisterUserEndpoints` → `RegisterDmUserEndpoints`
  - 関数呼び出しの変更: `RegisterPostEndpoints` → `RegisterDmPostEndpoints`

---

### タスク 6.9: server/internal/api/huma/huma_test.go の修正
- **ステータス**: 未着手
- **対象ファイル**: `server/internal/api/huma/huma_test.go`
- **作業内容**:
  - テスト関数名と型参照の変更（18件）:
    - `TestCreateUserInput` → `TestCreateDmUserInput`
    - `TestGetUserInput` → `TestGetDmUserInput`
    - `TestListUsersInput` → `TestListDmUsersInput`
    - `TestUpdateUserInput` → `TestUpdateDmUserInput`
    - `TestDeleteUserInput` → `TestDeleteDmUserInput`
    - `TestUserOutput` → `TestDmUserOutput`
    - `TestUsersOutput` → `TestDmUsersOutput`
    - `TestDeleteUserOutput` → `TestDeleteDmUserOutput`
    - `TestCreatePostInput` → `TestCreateDmPostInput`
    - `TestGetPostInput` → `TestGetDmPostInput`
    - `TestListPostsInput` → `TestListDmPostsInput`
    - `TestPostOutput` → `TestDmPostOutput`
    - `TestPostsOutput` → `TestDmPostsOutput`
    - `TestDeletePostOutput` → `TestDeleteDmPostOutput`
    - `TestGetUserPostsInput` → `TestGetDmUserPostsInput`
    - `TestUserPostsOutput` → `TestDmUserPostsOutput`

---

### タスク 6.10: server/internal/api/handler/dm_user_handler_huma_test.go の修正
- **ステータス**: 未着手
- **対象ファイル**: `server/internal/api/handler/dm_user_handler_huma_test.go`
- **作業内容**:
  - テスト関数名の変更: `TestRegisterUserEndpointsExists` → `TestRegisterDmUserEndpointsExists`
  - 関数参照の変更: `RegisterUserEndpoints` → `RegisterDmUserEndpoints`

---

### タスク 6.11: server/internal/api/handler/dm_post_handler_huma_test.go の修正
- **ステータス**: 未着手
- **対象ファイル**: `server/internal/api/handler/dm_post_handler_huma_test.go`
- **作業内容**:
  - テスト関数名の変更: `TestRegisterPostEndpointsExists` → `TestRegisterDmPostEndpointsExists`
  - 関数参照の変更: `RegisterPostEndpoints` → `RegisterDmPostEndpoints`

---

### タスク 6.12: server/internal/api/router/router_test.go の修正
- **ステータス**: 未着手
- **対象ファイル**: `server/internal/api/router/router_test.go`
- **作業内容**:
  - テスト関数名の変更: `TestRegisterUserEndpointsIntegration` → `TestRegisterDmUserEndpointsIntegration`
  - テスト関数名の変更: `TestRegisterPostEndpointsIntegration` → `TestRegisterDmPostEndpointsIntegration`
  - コメントの更新

---

### タスク 6.13: 全テストの実行と検証
- **ステータス**: 未着手
- **作業内容**:
  - 全テストの実行（`go test ./...`）
  - コンパイルの確認（`go build ./...`）
  - 旧命名規則の残存確認

---

## 実装サマリー
- **総タスク数**: 39タスク（Phase 1-5: 26タスク + Phase 6: 13タスク）
- **実装済みタスク数**: 26タスク
- **進捗率**: 66%（Phase 6 未着手）
