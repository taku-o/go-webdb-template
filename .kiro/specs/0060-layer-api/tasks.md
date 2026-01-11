# APIサーバーのusecaseのソースコードの位置を変更するの実装タスク一覧

## 概要
APIサーバー用のusecaseファイルを`server/internal/usecase/`から`server/internal/usecase/api/`に移動し、パッケージ名を`package usecase`から`package api`に変更するためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: ディレクトリ作成とファイル移動

#### タスク 1.1: API用usecaseディレクトリの作成
**目的**: APIサーバー用のusecaseディレクトリを新規作成する。

**作業内容**:
- `server/internal/usecase/api`ディレクトリを新規作成

**実装内容**:
- ディレクトリ作成: `server/internal/usecase/api`

**受け入れ基準**:
- [ ] `server/internal/usecase/api`ディレクトリが作成されている

- _Requirements: 3.1.1_
- _Design: 3.1.1_

---

#### タスク 1.2: dm_user_usecaseファイルの移動とパッケージ名変更
**目的**: `dm_user_usecase.go`と`dm_user_usecase_test.go`を移動し、パッケージ名を変更する。

**作業内容**:
- `server/internal/usecase/dm_user_usecase.go`を`server/internal/usecase/api/dm_user_usecase.go`に移動
- `server/internal/usecase/dm_user_usecase_test.go`を`server/internal/usecase/api/dm_user_usecase_test.go`に移動
- パッケージ名を`package usecase`から`package api`に変更

**実装内容**:
- ファイル移動:
  - `server/internal/usecase/dm_user_usecase.go` → `server/internal/usecase/api/dm_user_usecase.go`
  - `server/internal/usecase/dm_user_usecase_test.go` → `server/internal/usecase/api/dm_user_usecase_test.go`
- パッケージ名変更:
  - `package usecase` → `package api`

**受け入れ基準**:
- [ ] `server/internal/usecase/dm_user_usecase.go`が`server/internal/usecase/api/dm_user_usecase.go`に移動されている
- [ ] `server/internal/usecase/dm_user_usecase_test.go`が`server/internal/usecase/api/dm_user_usecase_test.go`に移動されている
- [ ] 移動後のファイルのパッケージ名が`package api`に変更されている

- _Requirements: 3.2.1, 3.2.2_
- _Design: 2.2.1, 2.2.2_

---

#### タスク 1.3: dm_post_usecaseファイルの移動とパッケージ名変更
**目的**: `dm_post_usecase.go`と`dm_post_usecase_test.go`を移動し、パッケージ名を変更する。

**作業内容**:
- `server/internal/usecase/dm_post_usecase.go`を`server/internal/usecase/api/dm_post_usecase.go`に移動
- `server/internal/usecase/dm_post_usecase_test.go`を`server/internal/usecase/api/dm_post_usecase_test.go`に移動
- パッケージ名を`package usecase`から`package api`に変更

**実装内容**:
- ファイル移動:
  - `server/internal/usecase/dm_post_usecase.go` → `server/internal/usecase/api/dm_post_usecase.go`
  - `server/internal/usecase/dm_post_usecase_test.go` → `server/internal/usecase/api/dm_post_usecase_test.go`
- パッケージ名変更:
  - `package usecase` → `package api`

**受け入れ基準**:
- [ ] `server/internal/usecase/dm_post_usecase.go`が`server/internal/usecase/api/dm_post_usecase.go`に移動されている
- [ ] `server/internal/usecase/dm_post_usecase_test.go`が`server/internal/usecase/api/dm_post_usecase_test.go`に移動されている
- [ ] 移動後のファイルのパッケージ名が`package api`に変更されている

- _Requirements: 3.2.1, 3.2.2_
- _Design: 2.2.1, 2.2.2_

---

#### タスク 1.4: dm_jobqueue_usecaseファイルの移動とパッケージ名変更
**目的**: `dm_jobqueue_usecase.go`と`dm_jobqueue_usecase_test.go`を移動し、パッケージ名を変更する。

**作業内容**:
- `server/internal/usecase/dm_jobqueue_usecase.go`を`server/internal/usecase/api/dm_jobqueue_usecase.go`に移動
- `server/internal/usecase/dm_jobqueue_usecase_test.go`を`server/internal/usecase/api/dm_jobqueue_usecase_test.go`に移動
- パッケージ名を`package usecase`から`package api`に変更

**実装内容**:
- ファイル移動:
  - `server/internal/usecase/dm_jobqueue_usecase.go` → `server/internal/usecase/api/dm_jobqueue_usecase.go`
  - `server/internal/usecase/dm_jobqueue_usecase_test.go` → `server/internal/usecase/api/dm_jobqueue_usecase_test.go`
- パッケージ名変更:
  - `package usecase` → `package api`

**受け入れ基準**:
- [ ] `server/internal/usecase/dm_jobqueue_usecase.go`が`server/internal/usecase/api/dm_jobqueue_usecase.go`に移動されている
- [ ] `server/internal/usecase/dm_jobqueue_usecase_test.go`が`server/internal/usecase/api/dm_jobqueue_usecase_test.go`に移動されている
- [ ] 移動後のファイルのパッケージ名が`package api`に変更されている

- _Requirements: 3.2.1, 3.2.2_
- _Design: 2.2.1, 2.2.2_

---

#### タスク 1.5: email_usecaseファイルの移動とパッケージ名変更
**目的**: `email_usecase.go`と`email_usecase_test.go`を移動し、パッケージ名を変更する。

**作業内容**:
- `server/internal/usecase/email_usecase.go`を`server/internal/usecase/api/email_usecase.go`に移動
- `server/internal/usecase/email_usecase_test.go`を`server/internal/usecase/api/email_usecase_test.go`に移動
- パッケージ名を`package usecase`から`package api`に変更

**実装内容**:
- ファイル移動:
  - `server/internal/usecase/email_usecase.go` → `server/internal/usecase/api/email_usecase.go`
  - `server/internal/usecase/email_usecase_test.go` → `server/internal/usecase/api/email_usecase_test.go`
- パッケージ名変更:
  - `package usecase` → `package api`

**受け入れ基準**:
- [ ] `server/internal/usecase/email_usecase.go`が`server/internal/usecase/api/email_usecase.go`に移動されている
- [ ] `server/internal/usecase/email_usecase_test.go`が`server/internal/usecase/api/email_usecase_test.go`に移動されている
- [ ] 移動後のファイルのパッケージ名が`package api`に変更されている

- _Requirements: 3.2.1, 3.2.2_
- _Design: 2.2.1, 2.2.2_

---

#### タスク 1.6: today_usecaseファイルの移動とパッケージ名変更
**目的**: `today_usecase.go`と`today_usecase_test.go`を移動し、パッケージ名を変更する。

**作業内容**:
- `server/internal/usecase/today_usecase.go`を`server/internal/usecase/api/today_usecase.go`に移動
- `server/internal/usecase/today_usecase_test.go`を`server/internal/usecase/api/today_usecase_test.go`に移動
- パッケージ名を`package usecase`から`package api`に変更

**実装内容**:
- ファイル移動:
  - `server/internal/usecase/today_usecase.go` → `server/internal/usecase/api/today_usecase.go`
  - `server/internal/usecase/today_usecase_test.go` → `server/internal/usecase/api/today_usecase_test.go`
- パッケージ名変更:
  - `package usecase` → `package api`

**受け入れ基準**:
- [ ] `server/internal/usecase/today_usecase.go`が`server/internal/usecase/api/today_usecase.go`に移動されている
- [ ] `server/internal/usecase/today_usecase_test.go`が`server/internal/usecase/api/today_usecase_test.go`に移動されている
- [ ] 移動後のファイルのパッケージ名が`package api`に変更されている

- _Requirements: 3.2.1, 3.2.2_
- _Design: 2.2.1, 2.2.2_

---

### Phase 2: import文の修正

#### タスク 2.1: dm_user_handlerのimport文修正
**目的**: `dm_user_handler.go`と`dm_user_handler_test.go`のimport文と型参照を修正する。

**作業内容**:
- `server/internal/api/handler/dm_user_handler.go`のimport文と型参照を修正
- `server/internal/api/handler/dm_user_handler_test.go`のimport文と型参照を修正（該当する場合）

**実装内容**:
- 修正対象: `server/internal/api/handler/dm_user_handler.go`、`server/internal/api/handler/dm_user_handler_test.go`
- import文の修正:
  ```go
  // 修正前
  import (
      "github.com/taku-o/go-webdb-template/internal/usecase"
  )
  
  // 修正後
  import (
      "github.com/taku-o/go-webdb-template/internal/usecase/api"
  )
  ```
- 型参照の修正:
  ```go
  // 修正前
  var dmUserUsecase *usecase.DmUserUsecase
  dmUserUsecase = usecase.NewDmUserUsecase(dmUserService)
  
  // 修正後
  var dmUserUsecase *api.DmUserUsecase
  dmUserUsecase = api.NewDmUserUsecase(dmUserService)
  ```

**受け入れ基準**:
- [ ] `server/internal/api/handler/dm_user_handler.go`のimport文が`internal/usecase/api`に変更されている
- [ ] `server/internal/api/handler/dm_user_handler.go`の型参照が`api.DmUserUsecase`に変更されている
- [ ] `server/internal/api/handler/dm_user_handler.go`のコンストラクタ呼び出しが`api.NewDmUserUsecase()`に変更されている
- [ ] `server/internal/api/handler/dm_user_handler_test.go`のimport文と型参照が修正されている（該当する場合）

- _Requirements: 3.3.1_
- _Design: 2.3.2.1_

---

#### タスク 2.2: dm_post_handlerのimport文修正
**目的**: `dm_post_handler.go`と`dm_post_handler_test.go`のimport文と型参照を修正する。

**作業内容**:
- `server/internal/api/handler/dm_post_handler.go`のimport文と型参照を修正
- `server/internal/api/handler/dm_post_handler_test.go`のimport文と型参照を修正（該当する場合）

**実装内容**:
- 修正対象: `server/internal/api/handler/dm_post_handler.go`、`server/internal/api/handler/dm_post_handler_test.go`
- import文の修正:
  ```go
  // 修正前
  import (
      "github.com/taku-o/go-webdb-template/internal/usecase"
  )
  
  // 修正後
  import (
      "github.com/taku-o/go-webdb-template/internal/usecase/api"
  )
  ```
- 型参照の修正:
  ```go
  // 修正前
  var dmPostUsecase *usecase.DmPostUsecase
  dmPostUsecase = usecase.NewDmPostUsecase(dmPostService)
  
  // 修正後
  var dmPostUsecase *api.DmPostUsecase
  dmPostUsecase = api.NewDmPostUsecase(dmPostService)
  ```

**受け入れ基準**:
- [ ] `server/internal/api/handler/dm_post_handler.go`のimport文が`internal/usecase/api`に変更されている
- [ ] `server/internal/api/handler/dm_post_handler.go`の型参照が`api.DmPostUsecase`に変更されている
- [ ] `server/internal/api/handler/dm_post_handler.go`のコンストラクタ呼び出しが`api.NewDmPostUsecase()`に変更されている
- [ ] `server/internal/api/handler/dm_post_handler_test.go`のimport文と型参照が修正されている（該当する場合）

- _Requirements: 3.3.1_
- _Design: 2.3.2.1_

---

#### タスク 2.3: dm_jobqueue_handlerのimport文修正
**目的**: `dm_jobqueue_handler.go`と`dm_jobqueue_handler_test.go`のimport文と型参照を修正する。

**作業内容**:
- `server/internal/api/handler/dm_jobqueue_handler.go`のimport文と型参照を修正
- `server/internal/api/handler/dm_jobqueue_handler_test.go`のimport文と型参照を修正

**実装内容**:
- 修正対象: `server/internal/api/handler/dm_jobqueue_handler.go`、`server/internal/api/handler/dm_jobqueue_handler_test.go`
- import文の修正:
  ```go
  // 修正前
  import (
      "github.com/taku-o/go-webdb-template/internal/usecase"
  )
  
  // 修正後
  import (
      "github.com/taku-o/go-webdb-template/internal/usecase/api"
  )
  ```
- 型参照の修正:
  ```go
  // 修正前
  var dmJobqueueUsecase *usecase.DmJobqueueUsecase
  dmJobqueueUsecase = usecase.NewDmJobqueueUsecase(jobQueueClientAdapter)
  
  // 修正後
  var dmJobqueueUsecase *api.DmJobqueueUsecase
  dmJobqueueUsecase = api.NewDmJobqueueUsecase(jobQueueClientAdapter)
  ```

**受け入れ基準**:
- [ ] `server/internal/api/handler/dm_jobqueue_handler.go`のimport文が`internal/usecase/api`に変更されている
- [ ] `server/internal/api/handler/dm_jobqueue_handler.go`の型参照が`api.DmJobqueueUsecase`に変更されている
- [ ] `server/internal/api/handler/dm_jobqueue_handler.go`のコンストラクタ呼び出しが`api.NewDmJobqueueUsecase()`に変更されている
- [ ] `server/internal/api/handler/dm_jobqueue_handler_test.go`のimport文と型参照が修正されている

- _Requirements: 3.3.1_
- _Design: 2.3.2.1_

---

#### タスク 2.4: email_handlerのimport文修正
**目的**: `email_handler.go`と`email_handler_test.go`のimport文と型参照を修正する。

**作業内容**:
- `server/internal/api/handler/email_handler.go`のimport文と型参照を修正
- `server/internal/api/handler/email_handler_test.go`のimport文と型参照を修正

**実装内容**:
- 修正対象: `server/internal/api/handler/email_handler.go`、`server/internal/api/handler/email_handler_test.go`
- import文の修正:
  ```go
  // 修正前
  import (
      "github.com/taku-o/go-webdb-template/internal/usecase"
  )
  
  // 修正後
  import (
      "github.com/taku-o/go-webdb-template/internal/usecase/api"
  )
  ```
- 型参照の修正:
  ```go
  // 修正前
  var emailUsecase *usecase.EmailUsecase
  emailUsecase = usecase.NewEmailUsecase(emailService, templateService)
  
  // 修正後
  var emailUsecase *api.EmailUsecase
  emailUsecase = api.NewEmailUsecase(emailService, templateService)
  ```

**受け入れ基準**:
- [ ] `server/internal/api/handler/email_handler.go`のimport文が`internal/usecase/api`に変更されている
- [ ] `server/internal/api/handler/email_handler.go`の型参照が`api.EmailUsecase`に変更されている
- [ ] `server/internal/api/handler/email_handler.go`のコンストラクタ呼び出しが`api.NewEmailUsecase()`に変更されている
- [ ] `server/internal/api/handler/email_handler_test.go`のimport文と型参照が修正されている

- _Requirements: 3.3.1_
- _Design: 2.3.2.1_

---

#### タスク 2.5: today_handlerのimport文修正
**目的**: `today_handler.go`と`today_handler_test.go`のimport文と型参照を修正する。

**作業内容**:
- `server/internal/api/handler/today_handler.go`のimport文と型参照を修正
- `server/internal/api/handler/today_handler_test.go`のimport文と型参照を修正

**実装内容**:
- 修正対象: `server/internal/api/handler/today_handler.go`、`server/internal/api/handler/today_handler_test.go`
- import文の修正:
  ```go
  // 修正前
  import (
      "github.com/taku-o/go-webdb-template/internal/usecase"
  )
  
  // 修正後
  import (
      "github.com/taku-o/go-webdb-template/internal/usecase/api"
  )
  ```
- 型参照の修正:
  ```go
  // 修正前
  var todayUsecase *usecase.TodayUsecase
  todayUsecase = usecase.NewTodayUsecase(dateService)
  
  // 修正後
  var todayUsecase *api.TodayUsecase
  todayUsecase = api.NewTodayUsecase(dateService)
  ```

**受け入れ基準**:
- [ ] `server/internal/api/handler/today_handler.go`のimport文が`internal/usecase/api`に変更されている
- [ ] `server/internal/api/handler/today_handler.go`の型参照が`api.TodayUsecase`に変更されている
- [ ] `server/internal/api/handler/today_handler.go`のコンストラクタ呼び出しが`api.NewTodayUsecase()`に変更されている
- [ ] `server/internal/api/handler/today_handler_test.go`のimport文と型参照が修正されている

- _Requirements: 3.3.1_
- _Design: 2.3.2.1_

---

#### タスク 2.6: main.goのimport文修正
**目的**: `server/cmd/server/main.go`のimport文と型参照を修正する。

**作業内容**:
- `server/cmd/server/main.go`のimport文と型参照を修正

**実装内容**:
- 修正対象: `server/cmd/server/main.go`
- import文の修正:
  ```go
  // 修正前
  import (
      "github.com/taku-o/go-webdb-template/internal/usecase"
  )
  
  // 修正後
  import (
      "github.com/taku-o/go-webdb-template/internal/usecase/api"
  )
  ```
- 型参照の修正:
  ```go
  // 修正前
  var dmJobqueueUsecase *usecase.DmJobqueueUsecase
  dmJobqueueUsecase = usecase.NewDmJobqueueUsecase(jobQueueClientAdapter)
  
  // 修正後
  var dmJobqueueUsecase *api.DmJobqueueUsecase
  dmJobqueueUsecase = api.NewDmJobqueueUsecase(jobQueueClientAdapter)
  ```

**受け入れ基準**:
- [ ] `server/cmd/server/main.go`のimport文が`internal/usecase/api`に変更されている
- [ ] `server/cmd/server/main.go`の型参照が`api.DmJobqueueUsecase`に変更されている
- [ ] `server/cmd/server/main.go`のコンストラクタ呼び出しが`api.NewDmJobqueueUsecase()`に変更されている

- _Requirements: 3.3.2_
- _Design: 2.3.2.2_

---

#### タスク 2.7: testutil/db.goのimport文修正
**目的**: `server/test/testutil/db.go`のimport文とインターフェース参照を修正する。

**作業内容**:
- `server/test/testutil/db.go`のimport文とインターフェース参照を修正

**実装内容**:
- 修正対象: `server/test/testutil/db.go`
- import文の修正:
  ```go
  // 修正前
  import (
      "github.com/taku-o/go-webdb-template/internal/usecase"
  )
  
  // 修正後
  import (
      "github.com/taku-o/go-webdb-template/internal/usecase/api"
  )
  ```
- インターフェース参照の修正:
  ```go
  // 修正前
  func CreateEmailHandler(emailService usecase.EmailServiceInterface, templateService usecase.TemplateServiceInterface) *handler.EmailHandler
  
  // 修正後
  func CreateEmailHandler(emailService api.EmailServiceInterface, templateService api.TemplateServiceInterface) *handler.EmailHandler
  ```

**受け入れ基準**:
- [ ] `server/test/testutil/db.go`のimport文が`internal/usecase/api`に変更されている
- [ ] `server/test/testutil/db.go`のインターフェース参照が`api.EmailServiceInterface`、`api.TemplateServiceInterface`に変更されている

- _Requirements: 3.3.3_
- _Design: 2.3.2.3_

---

#### タスク 2.8: admin/dm_user_register_usecase.goのimport文修正
**目的**: `server/internal/usecase/admin/dm_user_register_usecase.go`のimport文とインターフェース参照を修正する。

**作業内容**:
- `server/internal/usecase/admin/dm_user_register_usecase.go`のimport文とインターフェース参照を修正

**実装内容**:
- 修正対象: `server/internal/usecase/admin/dm_user_register_usecase.go`
- import文の修正:
  ```go
  // 修正前
  import (
      "github.com/taku-o/go-webdb-template/internal/usecase"
  )
  
  // 修正後
  import (
      "github.com/taku-o/go-webdb-template/internal/usecase/api"
  )
  ```
- インターフェース参照の修正:
  ```go
  // 修正前
  type DmUserRegisterUsecase struct {
      dmUserService usecase.DmUserServiceInterface
  }
  
  // 修正後
  type DmUserRegisterUsecase struct {
      dmUserService api.DmUserServiceInterface
  }
  ```

**受け入れ基準**:
- [ ] `server/internal/usecase/admin/dm_user_register_usecase.go`のimport文が`internal/usecase/api`に変更されている
- [ ] `server/internal/usecase/admin/dm_user_register_usecase.go`のインターフェース参照が`api.DmUserServiceInterface`に変更されている

- _Requirements: 3.3.4_
- _Design: 2.3.2.4_

---

#### タスク 2.9: cli/list_dm_users_usecase.goのimport文修正
**目的**: `server/internal/usecase/cli/list_dm_users_usecase.go`のimport文とインターフェース参照を修正する。

**作業内容**:
- `server/internal/usecase/cli/list_dm_users_usecase.go`のimport文とインターフェース参照を修正

**実装内容**:
- 修正対象: `server/internal/usecase/cli/list_dm_users_usecase.go`
- import文の修正:
  ```go
  // 修正前
  import (
      "github.com/taku-o/go-webdb-template/internal/usecase"
  )
  
  // 修正後
  import (
      "github.com/taku-o/go-webdb-template/internal/usecase/api"
  )
  ```
- インターフェース参照の修正:
  ```go
  // 修正前
  type ListDmUsersUsecase struct {
      dmUserService usecase.DmUserServiceInterface
  }
  
  // 修正後
  type ListDmUsersUsecase struct {
      dmUserService api.DmUserServiceInterface
  }
  ```

**受け入れ基準**:
- [ ] `server/internal/usecase/cli/list_dm_users_usecase.go`のimport文が`internal/usecase/api`に変更されている
- [ ] `server/internal/usecase/cli/list_dm_users_usecase.go`のインターフェース参照が`api.DmUserServiceInterface`に変更されている

- _Requirements: 3.3.4_
- _Design: 2.3.2.4_

---

### Phase 3: ビルドとテスト確認

#### タスク 3.1: ビルド確認
**目的**: 全てのパッケージが正常にビルドできることを確認する。

**作業内容**:
- ビルドコマンドを実行してエラーがないことを確認

**実装内容**:
- ビルドコマンド:
  ```bash
  cd server
  go build ./...
  ```

**受け入れ基準**:
- [ ] 全てのパッケージが正常にビルドできる
- [ ] importエラーが発生しない

- _Requirements: 6.4_
- _Design: 3.4.1, 4.2.1_

---

#### タスク 3.2: テスト実行
**目的**: 既存のテストが全て通過することを確認する。

**作業内容**:
- テストコマンドを実行して全てのテストが通過することを確認

**実装内容**:
- テストコマンド:
  ```bash
  cd server
  go test ./...
  ```

**受け入れ基準**:
- [ ] 既存のテストが全て通過する
- [ ] 移動後のファイルのテストが正常に動作する

- _Requirements: 6.4, 6.5_
- _Design: 3.4.2, 4.2.2, 4.2.3_

---

### Phase 4: ドキュメント更新

#### タスク 4.1: Project-Structure.mdの更新
**目的**: `docs/Project-Structure.md`に新規作成するディレクトリと移動したファイルのパスを反映する。

**作業内容**:
- `docs/Project-Structure.md`を修正
- `server/internal/usecase/api`ディレクトリを追加
- 移動したファイルのパスを更新

**実装内容**:
- 修正対象: `docs/Project-Structure.md`
- 更新内容:
  ```markdown
  │   │   ├── usecase/            # ビジネスロジック層
  │   │   │   ├── api/            # APIサーバー用usecase層（新規）
  │   │   │   │   ├── dm_user_usecase.go
  │   │   │   │   ├── dm_user_usecase_test.go
  │   │   │   │   ├── dm_post_usecase.go
  │   │   │   │   ├── dm_post_usecase_test.go
  │   │   │   │   ├── dm_jobqueue_usecase.go
  │   │   │   │   ├── dm_jobqueue_usecase_test.go
  │   │   │   │   ├── email_usecase.go
  │   │   │   │   ├── email_usecase_test.go
  │   │   │   │   ├── today_usecase.go
  │   │   │   │   └── today_usecase_test.go
  │   │   │   ├── admin/          # Admin用usecase層
  │   │   │   └── cli/            # CLI用usecase層
  ```

**受け入れ基準**:
- [ ] `docs/Project-Structure.md`に`server/internal/usecase/api`ディレクトリが追加されている
- [ ] `docs/Project-Structure.md`の移動したファイルのパスが更新されている

- _Requirements: 3.4.1, 6.6_
- _Design: 5.1.1_

---

#### タスク 4.2: structure.mdの更新
**目的**: `.kiro/steering/structure.md`に新規作成するディレクトリと移動したファイルのパスを反映する。

**作業内容**:
- `.kiro/steering/structure.md`を修正
- `server/internal/usecase/api`ディレクトリを追加
- 移動したファイルのパスを更新

**実装内容**:
- 修正対象: `.kiro/steering/structure.md`
- 更新内容:
  ```markdown
  │   │   ├── usecase/            # ビジネスロジック層
  │   │   │   ├── api/            # APIサーバー用usecase層（新規）
  │   │   │   │   ├── dm_user_usecase.go
  │   │   │   │   ├── dm_user_usecase_test.go
  │   │   │   │   ├── dm_post_usecase.go
  │   │   │   │   ├── dm_post_usecase_test.go
  │   │   │   │   ├── dm_jobqueue_usecase.go
  │   │   │   │   ├── dm_jobqueue_usecase_test.go
  │   │   │   │   ├── email_usecase.go
  │   │   │   │   ├── email_usecase_test.go
  │   │   │   │   ├── today_usecase.go
  │   │   │   │   └── today_usecase_test.go
  │   │   │   ├── admin/          # Admin用usecase層
  │   │   │   └── cli/            # CLI用usecase層
  ```

**受け入れ基準**:
- [ ] `.kiro/steering/structure.md`に`server/internal/usecase/api`ディレクトリが追加されている
- [ ] `.kiro/steering/structure.md`の移動したファイルのパスが更新されている

- _Requirements: 3.4.2, 6.6_
- _Design: 5.1.2_

---

### Phase 5: 動作確認

#### タスク 5.1: ローカル環境での動作確認
**目的**: ローカル環境でAPIサーバーが正常に動作することを確認する。

**作業内容**:
- サーバーを起動してエラーがないことを確認
- 既存のAPIエンドポイントが正常に動作することを確認

**実装内容**:
- サーバー起動:
  ```bash
  cd server/cmd/server
  go run main.go
  ```
- APIエンドポイントの動作確認

**受け入れ基準**:
- [ ] ローカル環境でAPIサーバーが正常に動作する
- [ ] 既存のAPIエンドポイントが正常に動作する

- _Requirements: 6.4_
- _Design: 3.4.3, 4.2.4_

---

## タスクの依存関係

- Phase 1のタスクは順次実行可能（並列実行も可能）
- Phase 2のタスクはPhase 1完了後に実行
- Phase 3のタスクはPhase 2完了後に実行
- Phase 4のタスクはPhase 1完了後に実行可能（Phase 2と並列実行可能）
- Phase 5のタスクはPhase 3完了後に実行

## 注意事項

- ファイル移動時は、パッケージ名の変更も同時に行う
- import文の修正時は、型参照とコンストラクタ呼び出しも同時に修正する
- テストファイルも同様に移動とパッケージ名変更が必要
- 全ての参照箇所を漏れなく修正する必要がある
