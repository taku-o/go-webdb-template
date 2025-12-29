# テーブル名変更機能修正要件定義書（リトライ）

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #54
- **Issueタイトル**: テーブル名を変更する (users, posts, news)
- **Feature名**: 0027-chtablename-retry
- **作成日**: 2025-01-27
- **関連Feature**: 0026-chtablename（前回の実装）

### 1.2 目的
Issue #54の対応として、テーブル名を`users`, `posts`, `news`から`dm_users`, `dm_posts`, `dm_news`に変更する作業を実施したが、変数名や関数名の修正漏れが多数存在する。本実装では、これらの修正漏れを全て修正し、命名の一貫性を確保する。

### 1.3 スコープ
- 変数名の修正漏れの修正
- 関数名の修正漏れの修正
- テストコード内の関数名参照の修正
- 命名規則の一貫性確保

**本実装の範囲外**:
- テーブル名の変更（既に実施済み）
- モデル名の変更（既に実施済み）
- ファイル名の変更（既に実施済み）
- 新機能の追加

## 2. 背景・現状分析

### 2.1 前回実装の状況
- Feature 0026-chtablenameでテーブル名変更を実施
- モデル名、Repository名、ファイル名などの主要な変更は完了
- しかし、変数名や関数名の修正漏れが多数存在

### 2.2 修正漏れの詳細

#### 2.2.1 server/cmd/generate-sample-data/main.go
以下の関数名が修正漏れ：
- `generateUsers` → `generateDmUsers`（45行目、67行目）
- `insertUsersBatch` → `insertDmUsersBatch`（97行目、201行目）
- `fetchUserIDs` → `fetchDmUserIDs`（102行目、292行目）
- `generatePosts` → `generateDmPosts`（51行目、117行目）
- `insertPostsBatch` → `insertDmPostsBatch`（153行目、235行目）
- `generateNews` → `generateDmNews`（56行目、165行目）
- `insertNewsBatch` → `insertDmNewsBatch`（191行目、269行目）

以下の変数名が修正漏れ：
- `posts` → `dmPosts`（136行目、148行目、152行目、153行目、241行目、243行目、244行目、246行目、253行目）
- `post` → `dmPost`（141行目、253行目）
- `news` → `dmNews`（173行目、186行目、190行目、191行目、275行目、277行目、278行目、280行目）

コメント内の修正漏れ：
- コメント内の`posts` → `dm_posts`、`news` → `dm_news`

#### 2.2.2 server/cmd/admin/main.go
以下の関数呼び出しが修正漏れ：
- `pages.UserRegisterPage` → `pages.DmUserRegisterPage`（91行目）
- `pages.UserRegisterCompletePage` → `pages.DmUserRegisterCompletePage`（94行目）

**注意**: 関数定義自体もリネームが必要な可能性がある（`server/internal/admin/pages/dm_user_register.go`、`dm_user_register_complete.go`）

#### 2.2.3 server/cmd/server/main.go
以下の変数名が修正漏れ：
- `userRepo` → `dmUserRepo`（42行目）
- `userService` → `dmUserService`（46行目）
- `userHandler` → `dmUserHandler`（50行目）
- `postRepo` → `dmPostRepo`（43行目）
- `postService` → `dmPostService`（47行目）
- `postHandler` → `dmPostHandler`（51行目）

#### 2.2.4 server/cmd/list-dm-users/main.go
以下の関数名が修正漏れ：
- `printUsersTSV` → `printDmUsersTSV`（75行目、93行目）

#### 2.2.5 server/cmd/list-dm-users/main_test.go
以下の関数呼び出しが修正漏れ：
- `printUsersTSV` → `printDmUsersTSV`（72行目、130行目）

#### 2.2.6 server/internal/repository/dm_post_repository_gorm.go
以下の変数名が修正漏れ：
- `tableUserPosts` → `tableDmUserPosts`（149行目、165行目、170行目）

#### 2.2.7 server/test/integration/dm_user_flow_gorm_test.go
以下のテスト関数名が修正漏れ：
- `TestUserCRUDFlowGORM` → `TestDmUserCRUDFlowGORM`（15行目）
- `TestUserCrossShardOperationsGORM` → `TestDmUserCrossShardOperationsGORM`（77行目）

#### 2.2.8 server/test/integration/dm_user_flow_test.go
以下のテスト関数名が修正漏れ：
- `TestUserCRUDFlow` → `TestDmUserCRUDFlow`（16行目）
- `TestUserCrossShardOperations` → `TestDmUserCrossShardOperations`（78行目）

#### 2.2.9 server/test/integration/dm_post_flow_test.go
以下のテスト関数名が修正漏れ：
- `TestPostCRUDFlow` → `TestDmPostCRUDFlow`（17行目）
- `TestCrossShardJoin` → `TestDmCrossShardJoin`（78行目）

コメント内の修正漏れ：
- コメント内の`User` → `DmUser`、`Post` → `DmPost`

#### 2.2.10 server/test/integration/sharding_test.go
以下のテスト関数名・変数名が修正漏れ：
- `TestCrossTableQueryUsers` → `TestCrossTableQueryDmUsers`（114行目）
- `testUsers` → `testDmUsers`（121行目）
- `u` → `dmU`（変数名、135行目など）
- コメント内の`User` → `DmUser`

#### 2.2.11 server/test/fixtures/dm_users.go
以下の関数名が修正漏れ：
- `CreateTestUser` → `CreateTestDmUser`（14行目）
- `CreateTestUserWithEmail` → `CreateTestDmUserWithEmail`（25行目）
- `CreateMultipleTestUsers` → `CreateMultipleTestDmUsers`（36行目）
- コメント内の`User` → `DmUser`

#### 2.2.12 server/test/e2e/api_test.go
以下のテスト関数名・変数名が修正漏れ：
- `TestUserAPI_CreateAndRetrieve` → `TestDmUserAPI_CreateAndRetrieve`（75行目）
- `TestUserAPI_UpdateAndDelete` → `TestDmUserAPI_UpdateAndDelete`（123行目）
- `TestPostAPI_CompleteFlow` → `TestDmPostAPI_CompleteFlow`（171行目）
- `user` → `dmUser`（変数名、91行目など）
- `post` → `dmPost`（変数名、203行目など）
- コメント内の`User` → `DmUser`、`Post` → `DmPost`

#### 2.2.13 server/internal/repository/interfaces.go
コメント内の修正漏れ：
- コメント内の`User` → `DmUser`（9行目など）

#### 2.2.14 server/internal/repository/dm_user_repository_gorm_test.go
以下のテスト関数名が修正漏れ：
- `TestUserRepositoryGORM_Create` → `TestDmUserRepositoryGORM_Create`（15行目）
- `TestUserRepositoryGORM_GetByID` → `TestDmUserRepositoryGORM_GetByID`（37行目）
- `TestUserRepositoryGORM_GetByID_NotFound` → `TestDmUserRepositoryGORM_GetByID_NotFound`（61行目）
- `TestUserRepositoryGORM_Update` → `TestDmUserRepositoryGORM_Update`（74行目）
- `TestUserRepositoryGORM_Delete` → `TestDmUserRepositoryGORM_Delete`（107行目）
- `TestUserRepositoryGORM_List` → `TestDmUserRepositoryGORM_List`（132行目）

#### 2.2.15 server/internal/repository/dm_post_repository.go
コメント内の修正漏れ：
- コメント内の`User` → `DmUser`（109行目など）

#### 2.2.16 server/internal/repository/dm_user_repository_test.go
以下のテスト関数名・変数名が修正漏れ：
- `TestUserRepository_Create` → `TestDmUserRepository_Create`（15行目）
- `TestUserRepository_GetByID` → `TestDmUserRepository_GetByID`（37行目）
- `TestUserRepository_GetByID_NotFound` → `TestDmUserRepository_GetByID_NotFound`（61行目）
- `TestUserRepository_Update` → `TestDmUserRepository_Update`（74行目）
- `TestUserRepository_Delete` → `TestDmUserRepository_Delete`（107行目）
- `TestUserRepository_List` → `TestDmUserRepository_List`（132行目）
- `users` → `dmUsers`（155行目）

#### 2.2.17 server/internal/admin/sharding.go
以下の関数名が修正漏れ：
- `FindUserAcrossShards` → `FindDmUserAcrossShards`（54行目）
- `CountUsersAcrossShards` → `CountDmUsersAcrossShards`（101行目）
- `FindPostAcrossShards` → `FindDmPostAcrossShards`（59行目）
- `CountPostsAcrossShards` → `CountDmPostsAcrossShards`（106行目）
- コメント内の`User` → `DmUser`、`Post` → `DmPost`

#### 2.2.18 server/test/integration/dm_post_flow_gorm_test.go
以下のテスト関数名が修正漏れ：
- `TestPostCRUDFlowGORM` → `TestDmPostCRUDFlowGORM`（15行目）
- `TestCrossShardJoinGORM` → `TestDmCrossShardJoinGORM`（78行目）

コメント内の修正漏れ：
- コメント内の`Post` → `DmPost`

#### 2.2.19 server/test/fixtures/dm_posts.go
以下の関数名が修正漏れ：
- `CreateTestPost` → `CreateTestDmPost`（14行目）
- `CreateTestPostWithContent` → `CreateTestDmPostWithContent`（26行目）
- `CreateMultipleTestPosts` → `CreateMultipleTestDmPosts`（38行目）

以下の変数名が修正漏れ：
- `post` → `dmPost`（21行目、33行目）
- `posts` → `dmPosts`（40行目、43行目）

コメント内の修正漏れ：
- コメント内の`Post` → `DmPost`

### 2.3 課題点
1. **命名の一貫性不足**: 一部の変数名・関数名が旧命名規則のまま残っている
2. **コードの可読性低下**: 命名が統一されていないため、コードの理解が困難
3. **保守性の低下**: 命名の不統一により、将来の修正時に混乱を招く可能性

### 2.4 本実装による改善点
1. **命名の一貫性確保**: 全ての変数名・関数名を`dm_`プレフィックス付きの命名規則に統一
2. **コードの可読性向上**: 統一された命名により、コードの理解が容易になる
3. **保守性の向上**: 命名の統一により、将来の修正時の混乱を防ぐ

## 3. 機能要件

### 3.1 server/cmd/generate-sample-data/main.go の修正

#### 3.1.1 関数名の変更
- **関数定義**: `generateUsers` → `generateDmUsers`（67行目）
- **関数呼び出し**: `generateUsers` → `generateDmUsers`（45行目）
- **関数定義**: `insertUsersBatch` → `insertDmUsersBatch`（201行目）
- **関数呼び出し**: `insertUsersBatch` → `insertDmUsersBatch`（97行目）
- **関数定義**: `fetchUserIDs` → `fetchDmUserIDs`（292行目）
- **関数呼び出し**: `fetchUserIDs` → `fetchDmUserIDs`（102行目）

#### 3.1.2 変数名の変更
- `userIDs` → `dmUserIDs`（45行目、70行目、102行目、106行目、117行目、139行目）
- `users` → `dmUsers`（84行目、92行目、97行目、219行目、221行目）
- `user` → `dmUser`（86行目、219行目）

#### 3.1.3 関数名の変更（post, news）
- **関数定義**: `generatePosts` → `generateDmPosts`（117行目）
- **関数呼び出し**: `generatePosts` → `generateDmPosts`（51行目）
- **関数定義**: `insertPostsBatch` → `insertDmPostsBatch`（235行目）
- **関数呼び出し**: `insertPostsBatch` → `insertDmPostsBatch`（153行目）
- **関数定義**: `generateNews` → `generateDmNews`（165行目）
- **関数呼び出し**: `generateNews` → `generateDmNews`（56行目）
- **関数定義**: `insertNewsBatch` → `insertDmNewsBatch`（269行目）
- **関数呼び出し**: `insertNewsBatch` → `insertDmNewsBatch`（191行目）

#### 3.1.4 変数名の変更（post, news）
- `posts` → `dmPosts`（136行目、148行目、152行目、153行目、241行目、243行目、244行目、246行目、253行目）
- `post` → `dmPost`（141行目、253行目）
- `news` → `dmNews`（173行目、186行目、190行目、191行目、275行目、277行目、278行目、280行目）

#### 3.1.5 コメントの更新
- 関数コメント内の`generateUsers` → `generateDmUsers`
- 関数コメント内の`insertUsersBatch` → `insertDmUsersBatch`
- 関数コメント内の`fetchUserIDs` → `fetchDmUserIDs`
- 関数コメント内の`generatePosts` → `generateDmPosts`
- 関数コメント内の`insertPostsBatch` → `insertDmPostsBatch`
- 関数コメント内の`generateNews` → `generateDmNews`
- 関数コメント内の`insertNewsBatch` → `insertDmNewsBatch`
- コメント内の`user_id` → `dm_user_id`（66行目）
- コメント内の`users` → `dm_users`、`posts` → `dm_posts`、`news` → `dm_news`

### 3.2 server/cmd/admin/main.go の修正

#### 3.2.1 関数呼び出しの変更
- `pages.UserRegisterPage` → `pages.DmUserRegisterPage`（91行目）
- `pages.UserRegisterCompletePage` → `pages.DmUserRegisterCompletePage`（94行目）

#### 3.2.2 関数定義の変更（必要に応じて）
- `server/internal/admin/pages/dm_user_register.go`: `UserRegisterPage` → `DmUserRegisterPage`（17行目）
- `server/internal/admin/pages/dm_user_register_complete.go`: `UserRegisterCompletePage` → `DmUserRegisterCompletePage`（14行目）
- `server/internal/admin/pages/pages.go`: `UserRegisterCompletePage` → `DmUserRegisterCompletePage`（53行目）

**注意**: 関数定義のリネームが必要な場合は、他の参照箇所も全て更新する必要がある

### 3.3 server/cmd/server/main.go の修正

#### 3.3.1 変数名の変更
- `userRepo` → `dmUserRepo`（42行目）
- `userService` → `dmUserService`（46行目）
- `userHandler` → `dmUserHandler`（50行目）

#### 3.3.2 変数名の変更（post）
- `postRepo` → `dmPostRepo`（43行目）
- `postService` → `dmPostService`（47行目）
- `postHandler` → `dmPostHandler`（51行目）

#### 3.3.3 変数参照の更新
- `userRepo`の参照を`dmUserRepo`に変更（47行目）
- `userService`の参照を`dmUserService`に変更（50行目）
- `userHandler`の参照を`dmUserHandler`に変更（55行目）
- `postRepo`の参照を`dmPostRepo`に変更（47行目）
- `postService`の参照を`dmPostService`に変更（51行目）
- `postHandler`の参照を`dmPostHandler`に変更（55行目）

### 3.4 server/cmd/list-dm-users/main.go の修正

#### 3.4.1 関数名の変更
- **関数定義**: `printUsersTSV` → `printDmUsersTSV`（93行目）
- **関数呼び出し**: `printUsersTSV` → `printDmUsersTSV`（75行目）

### 3.5 server/cmd/list-dm-users/main_test.go の修正

#### 3.5.1 関数呼び出しの変更
- `printUsersTSV` → `printDmUsersTSV`（72行目、130行目）

#### 3.5.2 変数名の変更
- `users` → `dmUsers`（18行目、24行目、30行目、44行目、92行目、115行目）
- `user` → `dmUser`（92行目、98行目、101行目、104行目）

#### 3.5.3 テスト関数名の変更（必須）
- `TestPrintUsersTSV` → `TestPrintDmUsersTSV`（15行目）
- `TestPrintUsersTSV_RFC3339Format` → `TestPrintDmUsersTSV_RFC3339Format`（114行目）

**注意**: テスト関数名の変更は必須。旧命名規則のままでは開発の邪魔になるため。

### 3.6 server/internal/repository/dm_post_repository_gorm.go の修正

#### 3.6.1 変数名の変更
- `tableUserPosts` → `tableDmUserPosts`（149行目、165行目、170行目）

### 3.7 server/test/integration/dm_user_flow_gorm_test.go の修正

#### 3.7.1 テスト関数名の変更
- `TestUserCRUDFlowGORM` → `TestDmUserCRUDFlowGORM`（15行目）
- `TestUserCrossShardOperationsGORM` → `TestDmUserCrossShardOperationsGORM`（77行目）

#### 3.7.2 コメントの更新
- コメント内の`User` → `DmUser`

### 3.8 server/test/integration/dm_user_flow_test.go の修正

#### 3.8.1 テスト関数名の変更
- `TestUserCRUDFlow` → `TestDmUserCRUDFlow`（16行目）
- `TestUserCrossShardOperations` → `TestDmUserCrossShardOperations`（78行目）

#### 3.8.2 コメントの更新
- コメント内の`User` → `DmUser`

### 3.9 server/test/integration/dm_post_flow_test.go の修正

#### 3.9.1 テスト関数名の変更
- `TestPostCRUDFlow` → `TestDmPostCRUDFlow`（17行目）
- `TestCrossShardJoin` → `TestDmCrossShardJoin`（78行目）

#### 3.9.2 コメントの更新
- コメント内の`User` → `DmUser`、`Post` → `DmPost`

### 3.10 server/test/integration/sharding_test.go の修正

#### 3.10.1 テスト関数名の変更
- `TestCrossTableQueryUsers` → `TestCrossTableQueryDmUsers`（114行目）

#### 3.10.2 変数名の変更
- `testUsers` → `testDmUsers`（121行目）
- `u` → `dmU`（ループ変数、135行目など）

#### 3.10.3 コメントの更新
- コメント内の`User` → `DmUser`

### 3.11 server/test/fixtures/dm_users.go の修正

#### 3.11.1 関数名の変更
- `CreateTestUser` → `CreateTestDmUser`（14行目）
- `CreateTestUserWithEmail` → `CreateTestDmUserWithEmail`（25行目）
- `CreateMultipleTestUsers` → `CreateMultipleTestDmUsers`（36行目）

#### 3.11.2 コメントの更新
- コメント内の`User` → `DmUser`

### 3.12 server/test/e2e/api_test.go の修正

#### 3.12.1 テスト関数名の変更
- `TestUserAPI_CreateAndRetrieve` → `TestDmUserAPI_CreateAndRetrieve`（75行目）
- `TestUserAPI_UpdateAndDelete` → `TestDmUserAPI_UpdateAndDelete`（123行目）
- `TestPostAPI_CompleteFlow` → `TestDmPostAPI_CompleteFlow`（171行目）

#### 3.12.2 変数名の変更
- `user` → `dmUser`（変数名、91行目など）
- `post` → `dmPost`（変数名、203行目など）

#### 3.12.3 コメントの更新
- コメント内の`User` → `DmUser`、`Post` → `DmPost`

### 3.13 server/internal/repository/interfaces.go の修正

#### 3.13.1 コメントの更新
- コメント内の`User` → `DmUser`

### 3.14 server/internal/repository/dm_user_repository_gorm_test.go の修正

#### 3.14.1 テスト関数名の変更
- `TestUserRepositoryGORM_Create` → `TestDmUserRepositoryGORM_Create`（15行目）
- `TestUserRepositoryGORM_GetByID` → `TestDmUserRepositoryGORM_GetByID`（37行目）
- `TestUserRepositoryGORM_GetByID_NotFound` → `TestDmUserRepositoryGORM_GetByID_NotFound`（61行目）
- `TestUserRepositoryGORM_Update` → `TestDmUserRepositoryGORM_Update`（74行目）
- `TestUserRepositoryGORM_Delete` → `TestDmUserRepositoryGORM_Delete`（107行目）
- `TestUserRepositoryGORM_List` → `TestDmUserRepositoryGORM_List`（132行目）

### 3.15 server/internal/repository/dm_post_repository.go の修正

#### 3.15.1 コメントの更新
- コメント内の`User` → `DmUser`

### 3.16 server/internal/repository/dm_user_repository_test.go の修正

#### 3.16.1 テスト関数名の変更
- `TestUserRepository_Create` → `TestDmUserRepository_Create`（15行目）
- `TestUserRepository_GetByID` → `TestDmUserRepository_GetByID`（37行目）
- `TestUserRepository_GetByID_NotFound` → `TestDmUserRepository_GetByID_NotFound`（61行目）
- `TestUserRepository_Update` → `TestDmUserRepository_Update`（74行目）
- `TestUserRepository_Delete` → `TestDmUserRepository_Delete`（107行目）
- `TestUserRepository_List` → `TestDmUserRepository_List`（132行目）

#### 3.16.2 変数名の変更
- `users` → `dmUsers`（155行目）

### 3.17 server/internal/admin/sharding.go の修正

#### 3.17.1 関数名の変更
- `FindUserAcrossShards` → `FindDmUserAcrossShards`（54行目）
- `CountUsersAcrossShards` → `CountDmUsersAcrossShards`（101行目）
- `FindPostAcrossShards` → `FindDmPostAcrossShards`（59行目）
- `CountPostsAcrossShards` → `CountDmPostsAcrossShards`（106行目）

#### 3.17.2 コメントの更新
- コメント内の`User` → `DmUser`、`Post` → `DmPost`

### 3.18 server/test/integration/dm_post_flow_gorm_test.go の修正

#### 3.18.1 テスト関数名の変更
- `TestPostCRUDFlowGORM` → `TestDmPostCRUDFlowGORM`（15行目）
- `TestCrossShardJoinGORM` → `TestDmCrossShardJoinGORM`（78行目）

#### 3.18.2 コメントの更新
- コメント内の`Post` → `DmPost`

### 3.19 server/test/fixtures/dm_posts.go の修正

#### 3.19.1 関数名の変更
- `CreateTestPost` → `CreateTestDmPost`（14行目）
- `CreateTestPostWithContent` → `CreateTestDmPostWithContent`（26行目）
- `CreateMultipleTestPosts` → `CreateMultipleTestDmPosts`（38行目）

#### 3.19.2 変数名の変更
- `post` → `dmPost`（21行目、33行目）
- `posts` → `dmPosts`（40行目、43行目）

#### 3.19.3 コメントの更新
- コメント内の`Post` → `DmPost`

## 4. 非機能要件

### 4.1 既存機能への影響
- 既存の機能は全て正常に動作すること
- 変数名・関数名の変更のみで、ロジックの変更は行わない
- 既存のテストが全て通過すること

### 4.2 テスト
- 既存のテストが全て通過すること
- 修正した関数名・変数名に関する新しいテストは不要（既存テストの更新のみ）

### 4.3 パフォーマンス
- 変数名・関数名の変更によるパフォーマンスへの影響はないこと

### 4.4 コード品質
- 命名規則の一貫性が確保されること
- 全ての修正漏れが修正されること
- コードレビューで指摘されないこと

## 5. 制約事項

### 5.1 技術的制約
- 既存のロジックを変更しないこと（変数名・関数名の変更のみ）
- 既存のAPIインターフェースを変更しないこと
- 既存のテストパターンを維持すること

### 5.2 プロジェクト制約
- 既存のコーディング規約に従うこと
- 既存の命名規則に従うこと（`dm_`プレフィックス付き）

### 5.3 命名規則
- 変数名: 
  - `user*` → `dmUser*`、`*User*` → `*DmUser*`
  - `post*` → `dmPost*`、`*Post*` → `*DmPost*`
  - `news*` → `dmNews*`、`*News*` → `*DmNews*`
- 関数名: 
  - `*User*` → `*DmUser*`、`*Users*` → `*DmUsers*`
  - `*Post*` → `*DmPost*`、`*Posts*` → `*DmPosts*`
  - `*News*` → `*DmNews*`
- 一貫性: 全ての変数名・関数名を統一された命名規則に従うこと

## 6. 受け入れ基準

### 6.1 server/cmd/generate-sample-data/main.go
- [ ] `generateUsers`が`generateDmUsers`に変更されている
- [ ] `insertUsersBatch`が`insertDmUsersBatch`に変更されている
- [ ] `fetchUserIDs`が`fetchDmUserIDs`に変更されている
- [ ] `generatePosts`が`generateDmPosts`に変更されている
- [ ] `insertPostsBatch`が`insertDmPostsBatch`に変更されている
- [ ] `generateNews`が`generateDmNews`に変更されている
- [ ] `insertNewsBatch`が`insertDmNewsBatch`に変更されている
- [ ] `userIDs`が`dmUserIDs`に変更されている
- [ ] `users`が`dmUsers`に変更されている
- [ ] `user`が`dmUser`に変更されている
- [ ] `posts`が`dmPosts`に変更されている
- [ ] `post`が`dmPost`に変更されている
- [ ] `news`が`dmNews`に変更されている
- [ ] 全ての関数呼び出しが更新されている
- [ ] 全ての変数参照が更新されている
- [ ] コメントが更新されている

### 6.2 server/cmd/admin/main.go
- [ ] `pages.UserRegisterPage`が`pages.DmUserRegisterPage`に変更されている
- [ ] `pages.UserRegisterCompletePage`が`pages.DmUserRegisterCompletePage`に変更されている
- [ ] 関数定義もリネームされている（必要に応じて）

### 6.3 server/cmd/server/main.go
- [ ] `userRepo`が`dmUserRepo`に変更されている
- [ ] `userService`が`dmUserService`に変更されている
- [ ] `userHandler`が`dmUserHandler`に変更されている
- [ ] `postRepo`が`dmPostRepo`に変更されている
- [ ] `postService`が`dmPostService`に変更されている
- [ ] `postHandler`が`dmPostHandler`に変更されている
- [ ] 全ての変数参照が更新されている

### 6.4 server/cmd/list-dm-users/main.go
- [ ] `printUsersTSV`が`printDmUsersTSV`に変更されている
- [ ] 全ての関数呼び出しが更新されている

### 6.5 server/cmd/list-dm-users/main_test.go
- [ ] `printUsersTSV`が`printDmUsersTSV`に変更されている
- [ ] `users`が`dmUsers`に変更されている
- [ ] `user`が`dmUser`に変更されている
- [ ] テスト関数名も更新されている（必須）

### 6.6 server/internal/repository/dm_post_repository_gorm.go
- [ ] `tableUserPosts`が`tableDmUserPosts`に変更されている
- [ ] 全ての変数参照が更新されている

### 6.7 server/test/integration/dm_user_flow_gorm_test.go
- [ ] テスト関数名が更新されている
- [ ] コメントが更新されている

### 6.8 server/test/integration/dm_user_flow_test.go
- [ ] テスト関数名が更新されている
- [ ] コメントが更新されている

### 6.9 server/test/integration/dm_post_flow_test.go
- [ ] テスト関数名が更新されている
- [ ] コメントが更新されている

### 6.10 server/test/integration/sharding_test.go
- [ ] テスト関数名が更新されている
- [ ] 変数名が更新されている
- [ ] コメントが更新されている

### 6.11 server/test/fixtures/dm_users.go
- [ ] 関数名が更新されている
- [ ] コメントが更新されている

### 6.12 server/test/e2e/api_test.go
- [ ] テスト関数名が更新されている（User, Post）
- [ ] 変数名が更新されている（user, post）
- [ ] コメントが更新されている

### 6.13 server/internal/repository/interfaces.go
- [ ] コメントが更新されている

### 6.14 server/internal/repository/dm_user_repository_gorm_test.go
- [ ] テスト関数名が更新されている

### 6.15 server/internal/repository/dm_post_repository.go
- [ ] コメントが更新されている

### 6.16 server/internal/repository/dm_user_repository_test.go
- [ ] テスト関数名が更新されている
- [ ] 変数名が更新されている

### 6.17 server/internal/admin/sharding.go
- [ ] 関数名が更新されている（User, Post）
- [ ] コメントが更新されている

### 6.18 server/test/integration/dm_post_flow_gorm_test.go
- [ ] テスト関数名が更新されている
- [ ] コメントが更新されている

### 6.19 server/test/fixtures/dm_posts.go
- [ ] 関数名が更新されている
- [ ] 変数名が更新されている
- [ ] コメントが更新されている

### 6.18 テスト
- [ ] 全テストが通過する
- [ ] 修正した関数名・変数名に関するテストが正常に動作する

### 6.8 コード品質
- [ ] 命名規則の一貫性が確保されている
- [ ] 全ての修正漏れが修正されている
- [ ] コードレビューで指摘されない

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### CLIツール
- `server/cmd/generate-sample-data/main.go`: 関数名の変更（7箇所: user 3, post 2, news 2）、変数名の変更（複数箇所: user, post, news）
- `server/cmd/admin/main.go`: 関数呼び出しの変更（2箇所）
- `server/cmd/server/main.go`: 変数名の変更（6箇所: user 3, post 3）
- `server/cmd/list-dm-users/main.go`: 関数名の変更（1箇所）
- `server/cmd/list-dm-users/main_test.go`: 関数呼び出しの変更（2箇所）

#### Repository層
- `server/internal/repository/dm_post_repository_gorm.go`: 変数名の変更（1箇所）
- `server/internal/repository/interfaces.go`: コメントの更新
- `server/internal/repository/dm_post_repository.go`: コメントの更新
- `server/internal/repository/dm_user_repository_gorm_test.go`: テスト関数名の変更（6箇所）
- `server/internal/repository/dm_user_repository_test.go`: テスト関数名・変数名の変更（7箇所）

#### テストコード
- `server/test/integration/dm_user_flow_gorm_test.go`: テスト関数名の変更（2箇所）
- `server/test/integration/dm_user_flow_test.go`: テスト関数名の変更（2箇所）
- `server/test/integration/dm_post_flow_test.go`: テスト関数名の変更（2箇所）、コメントの更新
- `server/test/integration/dm_post_flow_gorm_test.go`: テスト関数名の変更（2箇所）、コメントの更新
- `server/test/integration/sharding_test.go`: テスト関数名・変数名の変更（複数箇所）
- `server/test/fixtures/dm_users.go`: 関数名の変更（3箇所）
- `server/test/fixtures/dm_posts.go`: 関数名の変更（3箇所）、変数名の変更（複数箇所）
- `server/test/e2e/api_test.go`: テスト関数名・変数名の変更（複数箇所: user, post）

#### 管理画面
- `server/internal/admin/pages/dm_user_register.go`: 関数定義の変更（必要に応じて）
- `server/internal/admin/pages/dm_user_register_complete.go`: 関数定義の変更（必要に応じて）
- `server/internal/admin/pages/pages.go`: 関数参照の変更（必要に応じて）
- `server/internal/admin/sharding.go`: 関数名の変更（4箇所: user 2, post 2）

### 7.2 新規追加が必要なファイル
なし（既存ファイルの変更のみ）

### 7.3 削除されるファイル
なし（既存ファイルは変更のみ）

### 7.4 再利用する既存機能
- 既存のロジック（変更なし）
- 既存のテスト（更新のみ）

## 8. 実装上の注意事項

### 8.1 関数定義のリネーム
- `UserRegisterPage`、`UserRegisterCompletePage`の関数定義をリネームする場合、他の参照箇所も全て更新する必要がある
- `server/internal/admin/pages/pages.go`の`RegisterCustomPages`関数内の参照も更新が必要

### 8.2 テスト関数名の変更
- テスト関数名の変更は必須。旧命名規則のままでは開発の邪魔になるため。
- テスト関数名を変更する場合は、テスト実行時に正しく認識されることを確認する

### 8.3 変数参照の更新
- 変数名を変更する場合、その変数を参照している全ての箇所を更新する必要がある
- 特に`server/cmd/server/main.go`では、変数参照が複数箇所にあるため注意

### 8.4 コメントの更新
- 関数コメント内の関数名参照も更新する
- コードの可読性を維持するため、コメントの更新も重要

### 8.5 テスト関数名の変更（必須）
- テスト関数名の変更は必須。旧命名規則のままでは開発の邪魔になるため。
- テスト関数名を変更する場合、テスト実行時に正しく認識されることを確認する
- 特に統合テストやE2Eテストでは、テスト関数名がテスト結果に表示されるため、変更は必須

### 8.6 フィクスチャ関数名の変更
- フィクスチャ関数名を変更する場合、その関数を使用している全てのテストコードを更新する必要がある
- `server/test/fixtures/dm_users.go`の関数名を変更する場合、以下のファイルで使用されている可能性がある：
  - `server/test/integration/dm_post_flow_test.go`
  - その他の統合テストファイル

### 8.7 管理画面関数名の変更
- `FindUserAcrossShards`、`CountUsersAcrossShards`の関数名を変更する場合、これらの関数を使用している全ての箇所を更新する必要がある
- 管理画面のコードで使用されている可能性があるため、慎重に確認する

### 8.8 変数名の変更
- ループ変数（`u`など）の変更は、可読性を考慮して実施する
- 旧命名規則のままでは開発の邪魔になるため、変更は必須
- 本実装では、一貫性のため変更は必須

### 8.9 コメントの更新
- コメント内の`User`、`users`などの旧命名規則も全て更新する
- コードの可読性を維持するため、コメントの更新も重要
- 特に関数の説明コメントやテストの説明コメントは、正確に更新する

### 8.10 一貫性の確保
- 全ての修正漏れを漏れなく修正する
- 命名規則の一貫性を確保する
- コードレビューで指摘されないよう、慎重に実装する
- 修正後は、grep等で`\buser\b`、`\bUser\b`などの旧命名規則が残っていないか確認する

## 9. 参考情報

### 9.1 関連Issue
- GitHub Issue #54: テーブル名を変更する (users, posts, news)

### 9.2 関連Feature
- Feature 0026-chtablename: テーブル名変更機能（前回の実装）

### 9.3 既存実装
- `server/cmd/generate-sample-data/main.go`: サンプルデータ生成ツール
- `server/cmd/admin/main.go`: 管理画面サーバー
- `server/cmd/server/main.go`: APIサーバー
- `server/cmd/list-dm-users/main.go`: ユーザー一覧表示ツール
- `server/internal/repository/dm_post_repository_gorm.go`: 投稿Repository実装
- `server/internal/admin/pages/dm_user_register.go`: ユーザー登録ページ
- `server/internal/admin/pages/dm_user_register_complete.go`: ユーザー登録完了ページ

### 9.4 技術スタック
- **Go**: 1.21+
- **GORM**: v1.25.12
- **データベース**: SQLite3（開発環境）

### 9.5 変更パターン
- 変数名: 
  - `user*` → `dmUser*`、`*User*` → `*DmUser*`
  - `post*` → `dmPost*`、`*Post*` → `*DmPost*`
  - `news*` → `dmNews*`、`*News*` → `*DmNews*`
- 関数名: 
  - `*User*` → `*DmUser*`、`*Users*` → `*DmUsers*`
  - `*Post*` → `*DmPost*`、`*Posts*` → `*DmPosts*`
  - `*News*` → `*DmNews*`
- 一貫性: 全ての変数名・関数名を統一された命名規則に従う
