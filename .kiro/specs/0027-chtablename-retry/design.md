# テーブル名変更機能修正設計書（リトライ）

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、Issue #54の対応で発生した変数名・関数名の修正漏れを全て修正する機能の詳細設計を定義する。命名規則の一貫性を確保し、旧命名規則（`user`, `post`, `news`）が残らないようにする。

### 1.2 設計の範囲
- 変数名の修正漏れの修正設計
- 関数名の修正漏れの修正設計
- テストコード内の関数名・変数名参照の修正設計
- コメント内の旧命名規則の修正設計
- 実装順序と検証方法

### 1.3 設計方針
- **一貫性の確保**: 全ての変数名・関数名を`dm_`プレフィックス付きの命名規則に統一する
- **漏れの防止**: 旧命名規則（`user`, `post`, `news`）が残らないよう、体系的に修正する
- **段階的実装**: ファイル単位で段階的に実装し、各段階でテストを実行する
- **既存ロジックの維持**: 変数名・関数名の変更のみで、ロジックの変更は行わない

## 2. アーキテクチャ設計

### 2.1 命名規則の統一

#### 2.1.1 変数名の変更パターン

```
変更前 → 変更後:
- user* → dmUser*
- *User* → *DmUser*
- users → dmUsers
- user → dmUser
- post* → dmPost*
- *Post* → *DmPost*
- posts → dmPosts
- post → dmPost
- news* → dmNews*
- *News* → *DmNews*
- news → dmNews
```

#### 2.1.2 関数名の変更パターン

```
変更前 → 変更後:
- *User* → *DmUser*
- *Users* → *DmUsers*
- *Post* → *DmPost*
- *Posts* → *DmPosts*
- *News* → *DmNews*
```

#### 2.1.3 テスト関数名の変更パターン

```
変更前 → 変更後:
- Test*User* → Test*DmUser*
- Test*Post* → Test*DmPost*
- Test*News* → Test*DmNews*
```

### 2.2 変更対象ファイルの分類

#### 2.2.1 CLIツール
- `server/cmd/generate-sample-data/main.go`: 関数名・変数名の変更
- `server/cmd/admin/main.go`: 関数呼び出しの変更
- `server/cmd/server/main.go`: 変数名の変更
- `server/cmd/list-dm-users/main.go`: 関数名の変更
- `server/cmd/list-dm-users/main_test.go`: 関数呼び出し・テスト関数名・変数名の変更

#### 2.2.2 Repository層
- `server/internal/repository/dm_post_repository_gorm.go`: 変数名の変更
- `server/internal/repository/interfaces.go`: コメントの更新
- `server/internal/repository/dm_post_repository.go`: コメントの更新
- `server/internal/repository/dm_user_repository_gorm_test.go`: テスト関数名の変更
- `server/internal/repository/dm_user_repository_test.go`: テスト関数名・変数名の変更

#### 2.2.3 テストコード
- `server/test/integration/dm_user_flow_gorm_test.go`: テスト関数名の変更
- `server/test/integration/dm_user_flow_test.go`: テスト関数名の変更
- `server/test/integration/dm_post_flow_test.go`: テスト関数名の変更、コメントの更新
- `server/test/integration/dm_post_flow_gorm_test.go`: テスト関数名の変更、コメントの更新
- `server/test/integration/sharding_test.go`: テスト関数名・変数名の変更
- `server/test/fixtures/dm_users.go`: 関数名の変更、コメントの更新
- `server/test/fixtures/dm_posts.go`: 関数名の変更、変数名の変更、コメントの更新
- `server/test/e2e/api_test.go`: テスト関数名・変数名の変更、コメントの更新

#### 2.2.4 管理画面
- `server/internal/admin/pages/dm_user_register.go`: 関数定義の変更（必要に応じて）
- `server/internal/admin/pages/dm_user_register_complete.go`: 関数定義の変更（必要に応じて）
- `server/internal/admin/pages/pages.go`: 関数参照の変更（必要に応じて）
- `server/internal/admin/sharding.go`: 関数名の変更、コメントの更新

## 3. 実装設計

### 3.1 実装順序

#### Phase 1: CLIツールの修正
1. `server/cmd/generate-sample-data/main.go`
   - 関数名の変更（`generateUsers`, `insertUsersBatch`, `fetchUserIDs`, `generatePosts`, `insertPostsBatch`, `generateNews`, `insertNewsBatch`）
   - 変数名の変更（`userIDs`, `users`, `user`, `posts`, `post`, `news`）
   - コメントの更新
2. `server/cmd/server/main.go`
   - 変数名の変更（`userRepo`, `userService`, `userHandler`, `postRepo`, `postService`, `postHandler`）
   - 変数参照の更新
3. `server/cmd/admin/main.go`
   - 関数呼び出しの変更（`UserRegisterPage`, `UserRegisterCompletePage`）
4. `server/cmd/list-dm-users/main.go`
   - 関数名の変更（`printUsersTSV`）
5. `server/cmd/list-dm-users/main_test.go`
   - 関数呼び出しの変更、テスト関数名の変更、変数名の変更

#### Phase 2: Repository層の修正
1. `server/internal/repository/dm_post_repository_gorm.go`
   - 変数名の変更（`tableUserPosts`）
2. `server/internal/repository/interfaces.go`
   - コメントの更新
3. `server/internal/repository/dm_post_repository.go`
   - コメントの更新
4. `server/internal/repository/dm_user_repository_gorm_test.go`
   - テスト関数名の変更（6箇所）
5. `server/internal/repository/dm_user_repository_test.go`
   - テスト関数名の変更（6箇所）、変数名の変更（`users`）

#### Phase 3: テストコードの修正
1. `server/test/integration/dm_user_flow_gorm_test.go`
   - テスト関数名の変更（2箇所）、コメントの更新
2. `server/test/integration/dm_user_flow_test.go`
   - テスト関数名の変更（2箇所）、コメントの更新
3. `server/test/integration/dm_post_flow_test.go`
   - テスト関数名の変更（2箇所）、コメントの更新
4. `server/test/integration/dm_post_flow_gorm_test.go`
   - テスト関数名の変更（2箇所）、コメントの更新
5. `server/test/integration/sharding_test.go`
   - テスト関数名の変更、変数名の変更、コメントの更新
6. `server/test/fixtures/dm_users.go`
   - 関数名の変更（3箇所）、コメントの更新
7. `server/test/fixtures/dm_posts.go`
   - 関数名の変更（3箇所）、変数名の変更、コメントの更新
8. `server/test/e2e/api_test.go`
   - テスト関数名の変更（3箇所）、変数名の変更、コメントの更新

#### Phase 4: 管理画面の修正
1. `server/internal/admin/pages/dm_user_register.go`
   - 関数定義の変更（`UserRegisterPage` → `DmUserRegisterPage`、必要に応じて）
2. `server/internal/admin/pages/dm_user_register_complete.go`
   - 関数定義の変更（`UserRegisterCompletePage` → `DmUserRegisterCompletePage`、必要に応じて）
3. `server/internal/admin/pages/pages.go`
   - 関数参照の変更（必要に応じて）
4. `server/internal/admin/sharding.go`
   - 関数名の変更（`FindUserAcrossShards`, `CountUsersAcrossShards`, `FindPostAcrossShards`, `CountPostsAcrossShards`）
   - コメントの更新

### 3.2 実装方法

#### 3.2.1 変数名の変更
1. 変数定義の変更
2. その変数を参照している全ての箇所を更新
3. コメント内の参照も更新

#### 3.2.2 関数名の変更
1. 関数定義の変更
2. その関数を呼び出している全ての箇所を更新
3. コメント内の参照も更新
4. テスト関数名も更新（テスト関数の場合）

#### 3.2.3 コメントの更新
1. 関数コメント内の関数名参照を更新
2. 変数コメント内の変数名参照を更新
3. 一般的な説明コメント内の旧命名規則を更新

### 3.3 検証方法

#### 3.3.1 各Phase完了時の検証
1. コンパイルエラーの確認
   ```bash
   go build ./...
   ```
2. リンターエラーの確認
   ```bash
   go vet ./...
   ```
3. テストの実行
   ```bash
   go test ./...
   ```

#### 3.3.2 全体完了時の検証
1. 全テストの実行
   ```bash
   go test ./...
   ```
2. 旧命名規則の残存確認
   ```bash
   grep -r "\buser\b" server/ --exclude-dir=vendor
   grep -r "\bpost\b" server/ --exclude-dir=vendor
   grep -r "\bnews\b" server/ --exclude-dir=vendor
   ```
   （ただし、`dmUser`, `dmPost`, `dmNews`などの新命名規則は除外）

## 4. エラーハンドリング

### 4.1 コンパイルエラー
- 変数名・関数名の変更時に、参照箇所の更新漏れが発生する可能性がある
- コンパイルエラーが発生した場合は、エラーメッセージを確認し、参照箇所を特定して修正する

### 4.2 テストエラー
- テスト関数名の変更時に、テストが認識されない可能性がある
- テストが実行されない場合は、テスト関数名が正しく変更されているか確認する

### 4.3 参照漏れ
- 変数名・関数名の変更時に、参照箇所の更新漏れが発生する可能性がある
- grep等で旧命名規則が残っていないか確認する

## 5. 注意事項

### 5.1 関数定義のリネーム
- `UserRegisterPage`、`UserRegisterCompletePage`の関数定義をリネームする場合、他の参照箇所も全て更新する必要がある
- `server/internal/admin/pages/pages.go`の`RegisterCustomPages`関数内の参照も更新が必要

### 5.2 フィクスチャ関数名の変更
- `CreateTestUser`、`CreateTestPost`などのフィクスチャ関数名を変更する場合、その関数を使用している全てのテストコードを更新する必要がある
- 特に`server/test/integration/dm_post_flow_test.go`などで使用されている可能性がある

### 5.3 管理画面関数名の変更
- `FindUserAcrossShards`、`CountUsersAcrossShards`、`FindPostAcrossShards`、`CountPostsAcrossShards`の関数名を変更する場合、これらの関数を使用している全ての箇所を更新する必要がある
- 管理画面のコードで使用されている可能性があるため、慎重に確認する

### 5.4 変数名の変更
- ループ変数（`u`など）の変更は、可読性を考慮して実施する
- 旧命名規則のままでは開発の邪魔になるため、変更は必須

### 5.5 コメントの更新
- コメント内の`User`、`users`、`Post`、`posts`、`News`、`news`などの旧命名規則も全て更新する
- コードの可読性を維持するため、コメントの更新も重要

### 5.6 一貫性の確保
- 全ての修正漏れを漏れなく修正する
- 命名規則の一貫性を確保する
- 修正後は、grep等で`\buser\b`、`\bUser\b`、`\bpost\b`、`\bPost\b`、`\bnews\b`、`\bNews\b`などの旧命名規則が残っていないか確認する

## 6. 実装チェックリスト

### 6.1 Phase 1: CLIツールの修正
- [ ] `server/cmd/generate-sample-data/main.go`の修正
- [ ] `server/cmd/server/main.go`の修正
- [ ] `server/cmd/admin/main.go`の修正
- [ ] `server/cmd/list-dm-users/main.go`の修正
- [ ] `server/cmd/list-dm-users/main_test.go`の修正
- [ ] Phase 1完了時の検証

### 6.2 Phase 2: Repository層の修正
- [ ] `server/internal/repository/dm_post_repository_gorm.go`の修正
- [ ] `server/internal/repository/interfaces.go`の修正
- [ ] `server/internal/repository/dm_post_repository.go`の修正
- [ ] `server/internal/repository/dm_user_repository_gorm_test.go`の修正
- [ ] `server/internal/repository/dm_user_repository_test.go`の修正
- [ ] Phase 2完了時の検証

### 6.3 Phase 3: テストコードの修正
- [ ] `server/test/integration/dm_user_flow_gorm_test.go`の修正
- [ ] `server/test/integration/dm_user_flow_test.go`の修正
- [ ] `server/test/integration/dm_post_flow_test.go`の修正
- [ ] `server/test/integration/dm_post_flow_gorm_test.go`の修正
- [ ] `server/test/integration/sharding_test.go`の修正
- [ ] `server/test/fixtures/dm_users.go`の修正
- [ ] `server/test/fixtures/dm_posts.go`の修正
- [ ] `server/test/e2e/api_test.go`の修正
- [ ] Phase 3完了時の検証

### 6.4 Phase 4: 管理画面の修正
- [ ] `server/internal/admin/pages/dm_user_register.go`の修正
- [ ] `server/internal/admin/pages/dm_user_register_complete.go`の修正
- [ ] `server/internal/admin/pages/pages.go`の修正
- [ ] `server/internal/admin/sharding.go`の修正
- [ ] Phase 4完了時の検証

### 6.5 最終検証
- [ ] 全テストの実行
- [ ] 旧命名規則の残存確認
- [ ] コンパイルエラーの確認
- [ ] リンターエラーの確認

## 7. 参考情報

### 7.1 関連ドキュメント
- 要件定義書: `requirements.md`
- 前回実装: Feature 0026-chtablename

### 7.2 技術スタック
- **Go**: 1.21+
- **GORM**: v1.25.12
- **データベース**: SQLite3（開発環境）

### 7.3 変更パターン
- 変数名: 
  - `user*` → `dmUser*`、`*User*` → `*DmUser*`
  - `post*` → `dmPost*`、`*Post*` → `*DmPost*`
  - `news*` → `dmNews*`、`*News*` → `*DmNews*`
- 関数名: 
  - `*User*` → `*DmUser*`、`*Users*` → `*DmUsers*`
  - `*Post*` → `*DmPost*`、`*Posts*` → `*DmPosts*`
  - `*News*` → `*DmNews*`
