# APIサーバーのusecaseのソースコードの位置を変更するの設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、APIサーバー用のusecaseファイルを`server/internal/usecase/`から`server/internal/usecase/api/`に移動し、関連する全てのimport文を修正するための詳細設計を定義する。これにより、APIサーバー用、Adminアプリ用、CLI用のusecaseを明確に分離し、ディレクトリ構造の一貫性を向上させる。

### 1.2 設計の範囲
- API用usecaseディレクトリ（`server/internal/usecase/api`）の作成設計
- usecaseファイルの移動設計（10ファイル）
- パッケージ名の変更設計（`package usecase` → `package api`）
- import文の修正設計（14ファイル）
- テスト設計
- ドキュメント更新の設計

### 1.3 設計方針
- **一貫性**: Adminアプリ用、CLI用のusecaseと同じディレクトリ構造を採用
- **パッケージ名の変更**: パッケージ名を`package usecase`から`package api`に変更（admin、cliと実装を統一するため）
- **ファイル内容の変更**: パッケージ名を変更するため、ファイル内容も修正が必要
- **全ての参照箇所の修正**: 移動したusecaseを参照している全てのファイルのimport文と型参照を修正
- **テストの維持**: 既存のテストが正常に動作することを確認

## 2. アーキテクチャ設計

### 2.1 全体構成

#### 2.1.1 ディレクトリ構造の変更

**変更前**:
```
server/internal/usecase/
├── dm_user_usecase.go
├── dm_user_usecase_test.go
├── dm_post_usecase.go
├── dm_post_usecase_test.go
├── dm_jobqueue_usecase.go
├── dm_jobqueue_usecase_test.go
├── email_usecase.go
├── email_usecase_test.go
├── today_usecase.go
├── today_usecase_test.go
├── admin/
│   └── ...
└── cli/
    └── ...
```

**変更後**:
```
server/internal/usecase/
├── api/                          # 新規作成
│   ├── dm_user_usecase.go        # 移動
│   ├── dm_user_usecase_test.go    # 移動
│   ├── dm_post_usecase.go        # 移動
│   ├── dm_post_usecase_test.go   # 移動
│   ├── dm_jobqueue_usecase.go    # 移動
│   ├── dm_jobqueue_usecase_test.go # 移動
│   ├── email_usecase.go          # 移動
│   ├── email_usecase_test.go    # 移動
│   ├── today_usecase.go          # 移動
│   └── today_usecase_test.go     # 移動
├── admin/
│   └── ...
└── cli/
    └── ...
```

#### 2.1.2 パッケージ名の変更

- **パッケージ名**: 移動後のファイルは`package api`に変更
- **理由**: admin、cliと実装を統一し、パッケージ名で用途を明確に識別できるようにするため
- **影響**: import文は`internal/usecase/api`に変更し、パッケージ名は`api`を使用。型の参照は`api.DmUserUsecase`のように変更

### 2.2 ファイル移動設計

#### 2.2.1 移動対象ファイル一覧

| 現在のパス | 移動後のパス | パッケージ名 |
|-----------|-------------|------------|
| `server/internal/usecase/dm_user_usecase.go` | `server/internal/usecase/api/dm_user_usecase.go` | `package api` |
| `server/internal/usecase/dm_user_usecase_test.go` | `server/internal/usecase/api/dm_user_usecase_test.go` | `package api` |
| `server/internal/usecase/dm_post_usecase.go` | `server/internal/usecase/api/dm_post_usecase.go` | `package api` |
| `server/internal/usecase/dm_post_usecase_test.go` | `server/internal/usecase/api/dm_post_usecase_test.go` | `package api` |
| `server/internal/usecase/dm_jobqueue_usecase.go` | `server/internal/usecase/api/dm_jobqueue_usecase.go` | `package api` |
| `server/internal/usecase/dm_jobqueue_usecase_test.go` | `server/internal/usecase/api/dm_jobqueue_usecase_test.go` | `package api` |
| `server/internal/usecase/email_usecase.go` | `server/internal/usecase/api/email_usecase.go` | `package api` |
| `server/internal/usecase/email_usecase_test.go` | `server/internal/usecase/api/email_usecase_test.go` | `package api` |
| `server/internal/usecase/today_usecase.go` | `server/internal/usecase/api/today_usecase.go` | `package api` |
| `server/internal/usecase/today_usecase_test.go` | `server/internal/usecase/api/today_usecase_test.go` | `package api` |

#### 2.2.2 ファイル移動の手順

1. **ディレクトリの作成**: `server/internal/usecase/api`ディレクトリを新規作成
2. **ファイルの移動**: 10ファイルを`server/internal/usecase/`から`server/internal/usecase/api/`に移動
3. **パッケージ名の変更**: 移動後のファイルのパッケージ名を`package usecase`から`package api`に変更
4. **ファイル内容の確認**: パッケージ名の変更が正しく行われていることを確認

### 2.3 import文修正設計

#### 2.3.1 修正対象ファイル一覧

| ファイルパス | 現在のimport | 修正後のimport | 使用箇所（修正後） |
|------------|------------|--------------|-----------------|
| `server/internal/api/handler/dm_user_handler.go` | `internal/usecase` | `internal/usecase/api` | `*api.DmUserUsecase` |
| `server/internal/api/handler/dm_user_handler_test.go` | `internal/usecase` | `internal/usecase/api` | 該当する場合 |
| `server/internal/api/handler/dm_post_handler.go` | `internal/usecase` | `internal/usecase/api` | `*api.DmPostUsecase` |
| `server/internal/api/handler/dm_post_handler_test.go` | `internal/usecase` | `internal/usecase/api` | 該当する場合 |
| `server/internal/api/handler/dm_jobqueue_handler.go` | `internal/usecase` | `internal/usecase/api` | `*api.DmJobqueueUsecase` |
| `server/internal/api/handler/dm_jobqueue_handler_test.go` | `internal/usecase` | `internal/usecase/api` | `*api.DmJobqueueUsecase` |
| `server/internal/api/handler/email_handler.go` | `internal/usecase` | `internal/usecase/api` | `*api.EmailUsecase` |
| `server/internal/api/handler/email_handler_test.go` | `internal/usecase` | `internal/usecase/api` | `*api.EmailUsecase` |
| `server/internal/api/handler/today_handler.go` | `internal/usecase` | `internal/usecase/api` | `*api.TodayUsecase` |
| `server/internal/api/handler/today_handler_test.go` | `internal/usecase` | `internal/usecase/api` | `*api.TodayUsecase` |
| `server/cmd/server/main.go` | `internal/usecase` | `internal/usecase/api` | `*api.DmJobqueueUsecase` |
| `server/test/testutil/db.go` | `internal/usecase` | `internal/usecase/api` | `api.EmailServiceInterface`, `api.TemplateServiceInterface` |
| `server/internal/usecase/admin/dm_user_register_usecase.go` | `internal/usecase` | `internal/usecase/api` | `api.DmUserServiceInterface` |
| `server/internal/usecase/cli/list_dm_users_usecase.go` | `internal/usecase` | `internal/usecase/api` | `api.DmUserServiceInterface` |

#### 2.3.2 import文修正の詳細設計

##### 2.3.2.1 API Handler層のimport文修正

**修正例: `server/internal/api/handler/dm_user_handler.go`**

```go
// 修正前
import (
    "github.com/taku-o/go-webdb-template/internal/usecase"
)

var dmUserUsecase *usecase.DmUserUsecase

// 修正後
import (
    "github.com/taku-o/go-webdb-template/internal/usecase/api"
)

var dmUserUsecase *api.DmUserUsecase
```

**注意点**:
- importパスは`internal/usecase/api`に変更
- パッケージ名は`api`に変更されるため、型の参照は`api.DmUserUsecase`に変更
- コンストラクタの呼び出しも`api.NewDmUserUsecase()`に変更

##### 2.3.2.2 main.goのimport文修正

**修正例: `server/cmd/server/main.go`**

```go
// 修正前
import (
    "github.com/taku-o/go-webdb-template/internal/usecase"
)

var dmJobqueueUsecase *usecase.DmJobqueueUsecase
dmJobqueueUsecase = usecase.NewDmJobqueueUsecase(jobQueueClientAdapter)

// 修正後
import (
    "github.com/taku-o/go-webdb-template/internal/usecase/api"
)

var dmJobqueueUsecase *api.DmJobqueueUsecase
dmJobqueueUsecase = api.NewDmJobqueueUsecase(jobQueueClientAdapter)
```

**注意点**:
- importパスは`internal/usecase/api`に変更
- パッケージ名は`api`に変更されるため、型の参照は`api.DmJobqueueUsecase`に変更
- コンストラクタの呼び出しも`api.NewDmJobqueueUsecase()`に変更

##### 2.3.2.3 テストユーティリティのimport文修正

**修正例: `server/test/testutil/db.go`**

```go
// 修正前
import (
    "github.com/taku-o/go-webdb-template/internal/usecase"
)

func CreateEmailHandler(emailService usecase.EmailServiceInterface, templateService usecase.TemplateServiceInterface) *handler.EmailHandler

// 修正後
import (
    "github.com/taku-o/go-webdb-template/internal/usecase/api"
)

func CreateEmailHandler(emailService api.EmailServiceInterface, templateService api.TemplateServiceInterface) *handler.EmailHandler
```

**注意点**:
- importパスは`internal/usecase/api`に変更
- パッケージ名は`api`に変更されるため、インターフェースの参照は`api.EmailServiceInterface`、`api.TemplateServiceInterface`に変更

##### 2.3.2.4 Admin/CLI用usecaseのimport文修正

**修正例: `server/internal/usecase/admin/dm_user_register_usecase.go`**

```go
// 修正前
import (
    "github.com/taku-o/go-webdb-template/internal/usecase"
)

type DmUserRegisterUsecase struct {
    dmUserService usecase.DmUserServiceInterface
}

// 修正後
import (
    "github.com/taku-o/go-webdb-template/internal/usecase/api"
)

type DmUserRegisterUsecase struct {
    dmUserService api.DmUserServiceInterface
}
```

**注意点**:
- importパスは`internal/usecase/api`に変更
- パッケージ名は`api`に変更されるため、インターフェースの参照は`api.DmUserServiceInterface`に変更

### 2.4 パッケージ名の変更設計

#### 2.4.1 パッケージ名の変更理由

- **一貫性**: admin、cliと実装を統一し、パッケージ名で用途を明確に識別できるようにするため
- **明確性**: パッケージ名が`api`になることで、APIサーバー用のusecaseであることが明確になる
- **後方互換性の排除**: 後方互換性を維持すると、admin、cliと実装が異なり害の方が大きいため、パッケージ名を変更する

#### 2.4.2 パッケージ名の使用例

```go
// import文
import (
    "github.com/taku-o/go-webdb-template/internal/usecase/api"
)

// 型の参照（パッケージ名はapiに変更）
var dmUserUsecase *api.DmUserUsecase

// コンストラクタの呼び出し（パッケージ名はapiに変更）
dmUserUsecase = api.NewDmUserUsecase(dmUserService)

// インターフェースの参照（パッケージ名はapiに変更）
type DmUserRegisterUsecase struct {
    dmUserService api.DmUserServiceInterface
}
```

## 3. 実装設計

### 3.1 ディレクトリ作成

#### 3.1.1 ディレクトリ構造の作成

```bash
# ディレクトリの作成
mkdir -p server/internal/usecase/api
```

### 3.2 ファイル移動

#### 3.2.1 ファイル移動の手順

1. **ディレクトリの作成**: `server/internal/usecase/api`ディレクトリを新規作成
2. **ファイルの移動**: 10ファイルを`server/internal/usecase/`から`server/internal/usecase/api/`に移動
   - `mv server/internal/usecase/dm_user_usecase.go server/internal/usecase/api/dm_user_usecase.go`
   - `mv server/internal/usecase/dm_user_usecase_test.go server/internal/usecase/api/dm_user_usecase_test.go`
   - `mv server/internal/usecase/dm_post_usecase.go server/internal/usecase/api/dm_post_usecase.go`
   - `mv server/internal/usecase/dm_post_usecase_test.go server/internal/usecase/api/dm_post_usecase_test.go`
   - `mv server/internal/usecase/dm_jobqueue_usecase.go server/internal/usecase/api/dm_jobqueue_usecase.go`
   - `mv server/internal/usecase/dm_jobqueue_usecase_test.go server/internal/usecase/api/dm_jobqueue_usecase_test.go`
   - `mv server/internal/usecase/email_usecase.go server/internal/usecase/api/email_usecase.go`
   - `mv server/internal/usecase/email_usecase_test.go server/internal/usecase/api/email_usecase_test.go`
   - `mv server/internal/usecase/today_usecase.go server/internal/usecase/api/today_usecase.go`
   - `mv server/internal/usecase/today_usecase_test.go server/internal/usecase/api/today_usecase_test.go`
3. **パッケージ名の変更**: 移動後のファイルのパッケージ名を`package usecase`から`package api`に変更
4. **ファイル内容の確認**: パッケージ名の変更が正しく行われていることを確認

### 3.3 import文の修正

#### 3.3.1 import文修正の手順

1. **API Handler層のimport文修正**: 5つのhandlerファイルとそのテストファイルのimport文を修正
2. **main.goのimport文修正**: `server/cmd/server/main.go`のimport文を修正
3. **テストユーティリティのimport文修正**: `server/test/testutil/db.go`のimport文を修正
4. **Admin/CLI用usecaseのimport文修正**: `server/internal/usecase/admin/dm_user_register_usecase.go`と`server/internal/usecase/cli/list_dm_users_usecase.go`のimport文を修正

#### 3.3.2 import文修正の詳細

各ファイルで以下の変更を行う：

```go
// 修正前
import (
    "github.com/taku-o/go-webdb-template/internal/usecase"
)

var dmUserUsecase *usecase.DmUserUsecase
dmUserUsecase = usecase.NewDmUserUsecase(dmUserService)

// 修正後
import (
    "github.com/taku-o/go-webdb-template/internal/usecase/api"
)

var dmUserUsecase *api.DmUserUsecase
dmUserUsecase = api.NewDmUserUsecase(dmUserService)
```

**注意点**:
- importパスは`internal/usecase/api`に変更
- パッケージ名は`api`に変更されるため、型の参照は`api.DmUserUsecase`に変更
- コンストラクタの呼び出しも`api.NewDmUserUsecase()`に変更
- インターフェースの参照も`api.DmUserServiceInterface`のように変更

### 3.4 動作確認

#### 3.4.1 ビルド確認

```bash
# ビルド確認
cd server
go build ./...
```

#### 3.4.2 テスト実行

```bash
# テスト実行
cd server
go test ./...
```

#### 3.4.3 動作確認

```bash
# サーバーの起動確認
cd server/cmd/server
go run main.go
```

## 4. テスト設計

### 4.1 テスト方針

- **既存テストの維持**: 移動後のファイルのテストが正常に動作することを確認
- **import文の確認**: 全てのテストファイルのimport文が正しく修正されていることを確認
- **動作確認**: 既存のテストが全て通過することを確認

### 4.2 テスト項目

#### 4.2.1 ビルドテスト

- [ ] 全てのパッケージが正常にビルドできる
- [ ] importエラーが発生しない

#### 4.2.2 単体テスト

- [ ] `server/internal/usecase/api/dm_user_usecase_test.go`が正常に動作する
- [ ] `server/internal/usecase/api/dm_post_usecase_test.go`が正常に動作する
- [ ] `server/internal/usecase/api/dm_jobqueue_usecase_test.go`が正常に動作する
- [ ] `server/internal/usecase/api/email_usecase_test.go`が正常に動作する
- [ ] `server/internal/usecase/api/today_usecase_test.go`が正常に動作する

#### 4.2.3 統合テスト

- [ ] API Handler層のテストが正常に動作する
- [ ] 既存の統合テストが全て通過する

#### 4.2.4 動作確認

- [ ] ローカル環境でAPIサーバーが正常に動作する
- [ ] 既存のAPIエンドポイントが正常に動作する

## 5. ドキュメント更新設計

### 5.1 更新対象ドキュメント

#### 5.1.1 プロジェクト構造ドキュメント

**ファイル**: `docs/Project-Structure.md`

**更新内容**:
- `server/internal/usecase/api`ディレクトリを追加
- 移動したファイルのパスを更新

**更新例**:
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

#### 5.1.2 ファイル組織ドキュメント

**ファイル**: `.kiro/steering/structure.md`

**更新内容**:
- `server/internal/usecase/api`ディレクトリを追加
- 移動したファイルのパスを更新

**更新例**:
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

## 6. 実装上の注意事項

### 6.1 ファイル移動時の注意事項

- **パッケージ名の変更**: 移動後のファイルのパッケージ名を`package usecase`から`package api`に変更
- **ファイル内容の変更**: パッケージ名を変更するため、ファイル内容も修正が必要
- **テストファイルの同時移動**: テストファイルも同時に移動し、パッケージ名も変更する必要がある

### 6.2 import文修正時の注意事項

- **全ての参照箇所の修正**: 移動したusecaseを参照している全てのファイルのimport文と型参照を修正する必要がある
- **パッケージ名の変更**: importパスは`internal/usecase/api`に変更し、パッケージ名は`api`を使用
- **型参照の変更**: 型の参照は`usecase.DmUserUsecase`から`api.DmUserUsecase`に変更
- **インターフェース参照の変更**: インターフェースの参照も`usecase.DmUserServiceInterface`から`api.DmUserServiceInterface`に変更

### 6.3 テスト時の注意事項

- **既存テストの動作確認**: 移動後、既存のテストが全て正常に動作することを確認
- **import文の確認**: 全てのテストファイルのimport文が正しく修正されていることを確認

### 6.4 ドキュメント更新時の注意事項

- **一貫性**: 全てのドキュメントで同じディレクトリ構造を記載
- **パスの正確性**: 移動したファイルのパスを正確に更新

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### 7.1.1 移動が必要なファイル（10ファイル）

- `server/internal/usecase/dm_user_usecase.go` → `server/internal/usecase/api/dm_user_usecase.go`
- `server/internal/usecase/dm_user_usecase_test.go` → `server/internal/usecase/api/dm_user_usecase_test.go`
- `server/internal/usecase/dm_post_usecase.go` → `server/internal/usecase/api/dm_post_usecase.go`
- `server/internal/usecase/dm_post_usecase_test.go` → `server/internal/usecase/api/dm_post_usecase_test.go`
- `server/internal/usecase/dm_jobqueue_usecase.go` → `server/internal/usecase/api/dm_jobqueue_usecase.go`
- `server/internal/usecase/dm_jobqueue_usecase_test.go` → `server/internal/usecase/api/dm_jobqueue_usecase_test.go`
- `server/internal/usecase/email_usecase.go` → `server/internal/usecase/api/email_usecase.go`
- `server/internal/usecase/email_usecase_test.go` → `server/internal/usecase/api/email_usecase_test.go`
- `server/internal/usecase/today_usecase.go` → `server/internal/usecase/api/today_usecase.go`
- `server/internal/usecase/today_usecase_test.go` → `server/internal/usecase/api/today_usecase_test.go`

#### 7.1.2 修正が必要なファイル（14ファイル）

- `server/internal/api/handler/dm_user_handler.go`: import文の修正
- `server/internal/api/handler/dm_user_handler_test.go`: import文の修正（該当する場合）
- `server/internal/api/handler/dm_post_handler.go`: import文の修正
- `server/internal/api/handler/dm_post_handler_test.go`: import文の修正（該当する場合）
- `server/internal/api/handler/dm_jobqueue_handler.go`: import文の修正
- `server/internal/api/handler/dm_jobqueue_handler_test.go`: import文の修正
- `server/internal/api/handler/email_handler.go`: import文の修正
- `server/internal/api/handler/email_handler_test.go`: import文の修正
- `server/internal/api/handler/today_handler.go`: import文の修正
- `server/internal/api/handler/today_handler_test.go`: import文の修正
- `server/cmd/server/main.go`: import文の修正
- `server/test/testutil/db.go`: import文の修正
- `server/internal/usecase/admin/dm_user_register_usecase.go`: import文の修正
- `server/internal/usecase/cli/list_dm_users_usecase.go`: import文の修正

#### 7.1.3 更新が必要なドキュメント（2ファイル）

- `docs/Project-Structure.md`: ディレクトリ構造の更新
- `.kiro/steering/structure.md`: ディレクトリ構造の更新

### 7.2 既存機能への影響

- **既存のAPIサーバー**: 動作は維持されるが、import文の修正が必要
- **既存のAdmin/CLI用usecase**: インターフェースの参照は維持されるが、import文の修正が必要
- **既存のビジネスロジック**: 影響なし（ロジックは維持される）

## 8. 参考情報

### 8.1 関連ドキュメント

- `docs/Project-Structure.md`: プロジェクト構造ドキュメント
- `.kiro/steering/structure.md`: ファイル組織とコードパターン
- `.kiro/steering/tech.md`: 技術スタックとアーキテクチャ

### 8.2 既存実装の参考

- `server/internal/usecase/admin/`: Adminアプリ用のusecase層の実装パターン
- `server/internal/usecase/cli/`: CLI用のusecase層の実装パターン
- `server/internal/api/handler/`: API Handler層の実装パターン

### 8.3 技術スタック

- **言語**: Go
- **アーキテクチャ**: レイヤードアーキテクチャ（handler -> usecase -> service -> repository -> db）
- **テスト**: `testing`（標準ライブラリ）、`github.com/stretchr/testify`（アサーション、モック）

### 8.4 ディレクトリ構造の比較

| 項目 | 現在（修正前） | 修正後 |
|------|---------------|--------|
| APIサーバー用usecase | `usecase/dm_user_usecase.go` | `usecase/api/dm_user_usecase.go` |
| Adminアプリ用usecase | `usecase/admin/dm_user_register_usecase.go` | `usecase/admin/dm_user_register_usecase.go`（変更なし） |
| CLI用usecase | `usecase/cli/list_dm_users_usecase.go` | `usecase/cli/list_dm_users_usecase.go`（変更なし） |
| パッケージ名 | `package usecase` | `package api` |
| importパス | `internal/usecase` | `internal/usecase/api` |
| 型参照 | `usecase.DmUserUsecase` | `api.DmUserUsecase` |
