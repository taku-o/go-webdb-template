# テーブル名変更機能修正実装タスク一覧（リトライ）

## 概要
Issue #54の対応で発生した変数名・関数名の修正漏れを全て修正する機能の実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: CLIツールの修正

#### タスク 1.1: server/cmd/generate-sample-data/main.go の修正
**目的**: 関数名・変数名・コメントの修正漏れを修正

**作業内容**:
- 関数名の変更:
  - `generateUsers` → `generateDmUsers`（関数定義67行目、呼び出し45行目）
  - `insertUsersBatch` → `insertDmUsersBatch`（関数定義201行目、呼び出し97行目）
  - `fetchUserIDs` → `fetchDmUserIDs`（関数定義292行目、呼び出し102行目）
  - `generatePosts` → `generateDmPosts`（関数定義117行目、呼び出し51行目）
  - `insertPostsBatch` → `insertDmPostsBatch`（関数定義235行目、呼び出し153行目）
  - `generateNews` → `generateDmNews`（関数定義165行目、呼び出し56行目）
  - `insertNewsBatch` → `insertDmNewsBatch`（関数定義269行目、呼び出し191行目）
- 変数名の変更:
  - `userIDs` → `dmUserIDs`（45, 70, 102, 106, 117, 139行目）
  - `users` → `dmUsers`（84, 92, 97, 219, 221行目）
  - `user` → `dmUser`（86, 219行目）
  - `posts` → `dmPosts`（136, 148, 152, 153, 241, 243, 244, 246, 253行目）
  - `post` → `dmPost`（141, 253行目）
  - `news` → `dmNews`（173, 186, 190, 191, 275, 277, 278, 280行目）
- コメントの更新:
  - 関数コメント内の旧命名規則を新命名規則に変更
  - コメント内の`users` → `dm_users`、`posts` → `dm_posts`、`news` → `dm_news`
- コンパイルエラーの確認
- テストの実行（該当する場合）

**受け入れ基準**:
- 全ての関数名が新命名規則に変更されている
- 全ての変数名が新命名規則に変更されている
- 全てのコメントが更新されている
- コンパイルエラーがない
- テストが通過する

---

#### タスク 1.2: server/cmd/server/main.go の修正
**目的**: 変数名の修正漏れを修正

**作業内容**:
- 変数名の変更:
  - `userRepo` → `dmUserRepo`（42行目）
  - `userService` → `dmUserService`（46行目）
  - `userHandler` → `dmUserHandler`（50行目）
  - `postRepo` → `dmPostRepo`（43行目）
  - `postService` → `dmPostService`（47行目）
  - `postHandler` → `dmPostHandler`（51行目）
- 変数参照の更新:
  - `userRepo`の参照を`dmUserRepo`に変更（47行目）
  - `userService`の参照を`dmUserService`に変更（50行目）
  - `userHandler`の参照を`dmUserHandler`に変更（55行目）
  - `postRepo`の参照を`dmPostRepo`に変更（47行目）
  - `postService`の参照を`dmPostService`に変更（51行目）
  - `postHandler`の参照を`dmPostHandler`に変更（55行目）
- コンパイルエラーの確認
- テストの実行（該当する場合）

**受け入れ基準**:
- 全ての変数名が新命名規則に変更されている
- 全ての変数参照が更新されている
- コンパイルエラーがない
- テストが通過する

---

#### タスク 1.3: server/cmd/admin/main.go の修正
**目的**: 関数呼び出しの修正漏れを修正

**作業内容**:
- 関数呼び出しの変更:
  - `pages.UserRegisterPage` → `pages.DmUserRegisterPage`（91行目）
  - `pages.UserRegisterCompletePage` → `pages.DmUserRegisterCompletePage`（94行目）
- コンパイルエラーの確認
- 関数定義のリネームが必要か確認（必要に応じてPhase 4で対応）

**受け入れ基準**:
- 関数呼び出しが新命名規則に変更されている
- コンパイルエラーがない

---

#### タスク 1.4: server/cmd/list-dm-users/main.go の修正
**目的**: 関数名の修正漏れを修正

**作業内容**:
- 関数名の変更:
  - `printUsersTSV` → `printDmUsersTSV`（関数定義93行目、呼び出し75行目）
- コメントの更新:
  - 関数コメント内の`users` → `dmUsers`
- コンパイルエラーの確認

**受け入れ基準**:
- 関数名が新命名規則に変更されている
- コメントが更新されている
- コンパイルエラーがない

---

#### タスク 1.5: server/cmd/list-dm-users/main_test.go の修正
**目的**: 関数呼び出し・テスト関数名・変数名の修正漏れを修正

**作業内容**:
- 関数呼び出しの変更:
  - `printUsersTSV` → `printDmUsersTSV`（72行目、130行目）
- テスト関数名の変更:
  - `TestPrintUsersTSV` → `TestPrintDmUsersTSV`（15行目）
  - `TestPrintUsersTSV_RFC3339Format` → `TestPrintDmUsersTSV_RFC3339Format`（114行目）
- 変数名の変更:
  - `users` → `dmUsers`（18, 24, 30, 44, 92, 115行目）
  - `user` → `dmUser`（92, 98, 101, 104行目）
- テストの実行

**受け入れ基準**:
- 関数呼び出しが新命名規則に変更されている
- テスト関数名が新命名規則に変更されている
- 変数名が新命名規則に変更されている
- テストが通過する

---

### Phase 2: Repository層の修正

#### タスク 2.1: server/internal/repository/dm_post_repository_gorm.go の修正
**目的**: 変数名の修正漏れを修正

**作業内容**:
- 変数名の変更:
  - `tableUserPosts` → `tableDmUserPosts`（149, 165, 170行目）
- コンパイルエラーの確認
- テストの実行（該当する場合）

**受け入れ基準**:
- 変数名が新命名規則に変更されている
- 全ての変数参照が更新されている
- コンパイルエラーがない
- テストが通過する

---

#### タスク 2.2: server/internal/repository/interfaces.go の修正
**目的**: コメントの修正漏れを修正

**作業内容**:
- コメントの更新:
  - コメント内の`User` → `DmUser`、`Post` → `DmPost`

**受け入れ基準**:
- コメントが更新されている

---

#### タスク 2.3: server/internal/repository/dm_post_repository.go の修正
**目的**: コメントの修正漏れを修正

**作業内容**:
- コメントの更新:
  - コメント内の`User` → `DmUser`、`Post` → `DmPost`

**受け入れ基準**:
- コメントが更新されている

---

#### タスク 2.4: server/internal/repository/dm_user_repository_gorm_test.go の修正
**目的**: テスト関数名の修正漏れを修正

**作業内容**:
- テスト関数名の変更:
  - `TestUserRepositoryGORM_Create` → `TestDmUserRepositoryGORM_Create`（15行目）
  - `TestUserRepositoryGORM_GetByID` → `TestDmUserRepositoryGORM_GetByID`（37行目）
  - `TestUserRepositoryGORM_GetByID_NotFound` → `TestDmUserRepositoryGORM_GetByID_NotFound`（61行目）
  - `TestUserRepositoryGORM_Update` → `TestDmUserRepositoryGORM_Update`（74行目）
  - `TestUserRepositoryGORM_Delete` → `TestDmUserRepositoryGORM_Delete`（107行目）
  - `TestUserRepositoryGORM_List` → `TestDmUserRepositoryGORM_List`（132行目）
- テストの実行

**受け入れ基準**:
- テスト関数名が新命名規則に変更されている
- テストが通過する

---

#### タスク 2.5: server/internal/repository/dm_user_repository_test.go の修正
**目的**: テスト関数名・変数名の修正漏れを修正

**作業内容**:
- テスト関数名の変更:
  - `TestUserRepository_Create` → `TestDmUserRepository_Create`（15行目）
  - `TestUserRepository_GetByID` → `TestDmUserRepository_GetByID`（37行目）
  - `TestUserRepository_GetByID_NotFound` → `TestDmUserRepository_GetByID_NotFound`（61行目）
  - `TestUserRepository_Update` → `TestDmUserRepository_Update`（74行目）
  - `TestUserRepository_Delete` → `TestDmUserRepository_Delete`（107行目）
  - `TestUserRepository_List` → `TestDmUserRepository_List`（132行目）
- 変数名の変更:
  - `users` → `dmUsers`（155行目）
- テストの実行

**受け入れ基準**:
- テスト関数名が新命名規則に変更されている
- 変数名が新命名規則に変更されている
- テストが通過する

---

### Phase 3: テストコードの修正

#### タスク 3.1: server/test/integration/dm_user_flow_gorm_test.go の修正
**目的**: テスト関数名・コメントの修正漏れを修正

**作業内容**:
- テスト関数名の変更:
  - `TestUserCRUDFlowGORM` → `TestDmUserCRUDFlowGORM`（15行目）
  - `TestUserCrossShardOperationsGORM` → `TestDmUserCrossShardOperationsGORM`（77行目）
- コメントの更新:
  - コメント内の`User` → `DmUser`
- テストの実行

**受け入れ基準**:
- テスト関数名が新命名規則に変更されている
- コメントが更新されている
- テストが通過する

---

#### タスク 3.2: server/test/integration/dm_user_flow_test.go の修正
**目的**: テスト関数名・コメントの修正漏れを修正

**作業内容**:
- テスト関数名の変更:
  - `TestUserCRUDFlow` → `TestDmUserCRUDFlow`（16行目）
  - `TestUserCrossShardOperations` → `TestDmUserCrossShardOperations`（78行目）
- コメントの更新:
  - コメント内の`User` → `DmUser`
- テストの実行

**受け入れ基準**:
- テスト関数名が新命名規則に変更されている
- コメントが更新されている
- テストが通過する

---

#### タスク 3.3: server/test/integration/dm_post_flow_test.go の修正
**目的**: テスト関数名・コメントの修正漏れを修正

**作業内容**:
- テスト関数名の変更:
  - `TestPostCRUDFlow` → `TestDmPostCRUDFlow`（17行目）
  - `TestCrossShardJoin` → `TestDmCrossShardJoin`（78行目）
- コメントの更新:
  - コメント内の`User` → `DmUser`、`Post` → `DmPost`
- テストの実行

**受け入れ基準**:
- テスト関数名が新命名規則に変更されている
- コメントが更新されている
- テストが通過する

---

#### タスク 3.4: server/test/integration/dm_post_flow_gorm_test.go の修正
**目的**: テスト関数名・コメントの修正漏れを修正

**作業内容**:
- テスト関数名の変更:
  - `TestPostCRUDFlowGORM` → `TestDmPostCRUDFlowGORM`（15行目）
  - `TestCrossShardJoinGORM` → `TestDmCrossShardJoinGORM`（78行目）
- コメントの更新:
  - コメント内の`Post` → `DmPost`
- テストの実行

**受け入れ基準**:
- テスト関数名が新命名規則に変更されている
- コメントが更新されている
- テストが通過する

---

#### タスク 3.5: server/test/integration/sharding_test.go の修正
**目的**: テスト関数名・変数名・コメントの修正漏れを修正

**作業内容**:
- テスト関数名の変更:
  - `TestCrossTableQueryUsers` → `TestCrossTableQueryDmUsers`（114行目）
- 変数名の変更:
  - `testUsers` → `testDmUsers`（121行目）
  - `u` → `dmU`（ループ変数、135行目など）
- コメントの更新:
  - コメント内の`User` → `DmUser`
- テストの実行

**受け入れ基準**:
- テスト関数名が新命名規則に変更されている
- 変数名が新命名規則に変更されている
- コメントが更新されている
- テストが通過する

---

#### タスク 3.6: server/test/fixtures/dm_users.go の修正
**目的**: 関数名・コメントの修正漏れを修正

**作業内容**:
- 関数名の変更:
  - `CreateTestUser` → `CreateTestDmUser`（14行目）
  - `CreateTestUserWithEmail` → `CreateTestDmUserWithEmail`（25行目）
  - `CreateMultipleTestUsers` → `CreateMultipleTestDmUsers`（36行目）
- コメントの更新:
  - コメント内の`User` → `DmUser`
- この関数を使用している全てのテストコードの更新確認
- コンパイルエラーの確認
- テストの実行

**受け入れ基準**:
- 関数名が新命名規則に変更されている
- コメントが更新されている
- この関数を使用している全てのテストコードが更新されている
- コンパイルエラーがない
- テストが通過する

---

#### タスク 3.7: server/test/fixtures/dm_posts.go の修正
**目的**: 関数名・変数名・コメントの修正漏れを修正

**作業内容**:
- 関数名の変更:
  - `CreateTestPost` → `CreateTestDmPost`（14行目）
  - `CreateTestPostWithContent` → `CreateTestDmPostWithContent`（26行目）
  - `CreateMultipleTestPosts` → `CreateMultipleTestDmPosts`（38行目）
- 変数名の変更:
  - `post` → `dmPost`（21, 33行目）
  - `posts` → `dmPosts`（40, 43行目）
- コメントの更新:
  - コメント内の`Post` → `DmPost`
- この関数を使用している全てのテストコードの更新確認
- コンパイルエラーの確認
- テストの実行

**受け入れ基準**:
- 関数名が新命名規則に変更されている
- 変数名が新命名規則に変更されている
- コメントが更新されている
- この関数を使用している全てのテストコードが更新されている
- コンパイルエラーがない
- テストが通過する

---

#### タスク 3.8: server/test/e2e/api_test.go の修正
**目的**: テスト関数名・変数名・コメントの修正漏れを修正

**作業内容**:
- テスト関数名の変更:
  - `TestUserAPI_CreateAndRetrieve` → `TestDmUserAPI_CreateAndRetrieve`（75行目）
  - `TestUserAPI_UpdateAndDelete` → `TestDmUserAPI_UpdateAndDelete`（123行目）
  - `TestPostAPI_CompleteFlow` → `TestDmPostAPI_CompleteFlow`（171行目）
- 変数名の変更:
  - `user` → `dmUser`（91行目など）
  - `post` → `dmPost`（203行目など）
- コメントの更新:
  - コメント内の`User` → `DmUser`、`Post` → `DmPost`
- テストの実行

**受け入れ基準**:
- テスト関数名が新命名規則に変更されている
- 変数名が新命名規則に変更されている
- コメントが更新されている
- テストが通過する

---

### Phase 4: 管理画面の修正

#### タスク 4.1: server/internal/admin/pages/dm_user_register.go の修正
**目的**: 関数定義の修正漏れを修正（必要に応じて）

**作業内容**:
- 関数定義の変更:
  - `UserRegisterPage` → `DmUserRegisterPage`（17行目）
- この関数を呼び出している全ての箇所の更新確認
- コンパイルエラーの確認

**受け入れ基準**:
- 関数定義が新命名規則に変更されている（必要に応じて）
- この関数を呼び出している全ての箇所が更新されている
- コンパイルエラーがない

---

#### タスク 4.2: server/internal/admin/pages/dm_user_register_complete.go の修正
**目的**: 関数定義の修正漏れを修正（必要に応じて）

**作業内容**:
- 関数定義の変更:
  - `UserRegisterCompletePage` → `DmUserRegisterCompletePage`（14行目）
- この関数を呼び出している全ての箇所の更新確認
- コンパイルエラーの確認

**受け入れ基準**:
- 関数定義が新命名規則に変更されている（必要に応じて）
- この関数を呼び出している全ての箇所が更新されている
- コンパイルエラーがない

---

#### タスク 4.3: server/internal/admin/pages/pages.go の修正
**目的**: 関数参照の修正漏れを修正（必要に応じて）

**作業内容**:
- 関数参照の変更:
  - `UserRegisterCompletePage` → `DmUserRegisterCompletePage`（53行目）
- コンパイルエラーの確認

**受け入れ基準**:
- 関数参照が新命名規則に変更されている（必要に応じて）
- コンパイルエラーがない

---

#### タスク 4.4: server/internal/admin/sharding.go の修正
**目的**: 関数名・コメントの修正漏れを修正

**作業内容**:
- 関数名の変更:
  - `FindUserAcrossShards` → `FindDmUserAcrossShards`（54行目）
  - `CountUsersAcrossShards` → `CountDmUsersAcrossShards`（101行目）
  - `FindPostAcrossShards` → `FindDmPostAcrossShards`（59行目）
  - `CountPostsAcrossShards` → `CountDmPostsAcrossShards`（106行目）
- コメントの更新:
  - コメント内の`User` → `DmUser`、`Post` → `DmPost`
- この関数を使用している全ての箇所の更新確認
- コンパイルエラーの確認
- テストの実行（該当する場合）

**受け入れ基準**:
- 関数名が新命名規則に変更されている
- コメントが更新されている
- この関数を使用している全ての箇所が更新されている
- コンパイルエラーがない
- テストが通過する

---

### Phase 5: 最終検証

#### タスク 5.1: 全テストの実行
**目的**: 全ての修正が正常に動作することを確認

**作業内容**:
- 全テストの実行:
  ```bash
  go test ./...
  ```
- テスト結果の確認
- 失敗したテストの修正（該当する場合）

**受け入れ基準**:
- 全テストが通過する
- テストエラーがない

---

#### タスク 5.2: 旧命名規則の残存確認
**目的**: 旧命名規則が残っていないことを確認

**作業内容**:
- 旧命名規則の検索:
  ```bash
  grep -r "\buser\b" server/ --exclude-dir=vendor | grep -v "dmUser" | grep -v "DmUser"
  grep -r "\bpost\b" server/ --exclude-dir=vendor | grep -v "dmPost" | grep -v "DmPost"
  grep -r "\bnews\b" server/ --exclude-dir=vendor | grep -v "dmNews" | grep -v "DmNews"
  ```
- 検索結果の確認
- 残存している旧命名規則の修正（該当する場合）

**受け入れ基準**:
- 旧命名規則が残っていない（新命名規則は除外）
- 検索結果が空である

---

#### タスク 5.3: コンパイルエラーの確認
**目的**: コンパイルエラーがないことを確認

**作業内容**:
- コンパイルの実行:
  ```bash
  go build ./...
  ```
- コンパイルエラーの確認
- エラーの修正（該当する場合）

**受け入れ基準**:
- コンパイルエラーがない
- ビルドが成功する

---

#### タスク 5.4: リンターエラーの確認
**目的**: リンターエラーがないことを確認

**作業内容**:
- リンターの実行:
  ```bash
  go vet ./...
  ```
- リンターエラーの確認
- エラーの修正（該当する場合）

**受け入れ基準**:
- リンターエラーがない

---

## 実装上の注意事項

### 参照箇所の更新
- 変数名・関数名を変更する場合、その変数・関数を参照している全ての箇所を更新する必要がある
- 特に複数のファイルにまたがる参照がある場合は、慎重に確認する

### フィクスチャ関数の変更
- `CreateTestUser`、`CreateTestPost`などのフィクスチャ関数名を変更する場合、その関数を使用している全てのテストコードを更新する必要がある
- 特に`server/test/integration/dm_post_flow_test.go`などで使用されている可能性がある

### 管理画面関数の変更
- `FindUserAcrossShards`、`CountUsersAcrossShards`、`FindPostAcrossShards`、`CountPostsAcrossShards`の関数名を変更する場合、これらの関数を使用している全ての箇所を更新する必要がある
- 管理画面のコードで使用されている可能性があるため、慎重に確認する

### テスト関数名の変更
- テスト関数名を変更する場合、テスト実行時に正しく認識されることを確認する
- テスト関数名の変更は必須。旧命名規則のままでは開発の邪魔になるため。

### コメントの更新
- コメント内の`User`、`users`、`Post`、`posts`、`News`、`news`などの旧命名規則も全て更新する
- コードの可読性を維持するため、コメントの更新も重要

## 参考情報

### 関連ドキュメント
- 要件定義書: `requirements.md`
- 設計書: `design.md`
- 前回実装: Feature 0026-chtablename

### 技術スタック
- **Go**: 1.21+
- **GORM**: v1.25.12
- **データベース**: SQLite3（開発環境）

### 変更パターン
- 変数名: 
  - `user*` → `dmUser*`、`*User*` → `*DmUser*`
  - `post*` → `dmPost*`、`*Post*` → `*DmPost*`
  - `news*` → `dmNews*`、`*News*` → `*DmNews*`
- 関数名: 
  - `*User*` → `*DmUser*`、`*Users*` → `*DmUsers*`
  - `*Post*` → `*DmPost*`、`*Posts*` → `*DmPosts*`
  - `*News*` → `*DmNews*`
