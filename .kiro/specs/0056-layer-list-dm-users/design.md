# server/cmd/list-dm-usersのレイヤー構造修正の設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、`server/cmd/list-dm-users`の実装を、APIサーバーと同じレイヤー構造（usecase -> service -> repository -> model）に変更するための詳細設計を定義する。これにより、CLIコマンドとAPIサーバーで一貫したアーキテクチャを実現し、コードの保守性と再利用性を向上させる。

### 1.2 設計の範囲
- CLI用usecase層（`server/internal/usecase/cli`）の設計
- `server/cmd/list-dm-users/main.go`の簡素化設計
- 依存関係の注入設計
- テスト設計
- ドキュメント更新の設計

### 1.3 設計方針
- **一貫性**: APIサーバーと同じレイヤー構造を採用
- **既存コードの活用**: 既存のservice、repository、model層をそのまま使用
- **責務の明確化**: 各レイヤーの責務を明確に分離
- **テスト容易性**: usecase層を独立してテストできる設計
- **後方互換性**: 既存のCLIコマンドの動作（引数、出力形式、エラーメッセージ）を維持

## 2. アーキテクチャ設計

### 2.1 全体構成

```
┌─────────────────────────────────────────────────────────────┐
│              CLI Layer (cmd/list-dm-users/main.go)          │
│  • コマンドライン引数の解析                                  │
│  • 引数のバリデーション                                      │
│  • 設定ファイルの読み込み                                    │
│  • GroupManagerの初期化                                     │
│  • レイヤーの初期化（Repository → Service → Usecase）      │
│  • usecase層の呼び出し                                      │
│  • 結果の出力（TSV形式）                                    │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│         Usecase Layer (internal/usecase/cli)                  │
│  • ListDmUsersUsecase                                        │
│  • ビジネスロジックの調整（CLI用）                           │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│         Service Layer (internal/service)                      │
│  • DmUserService                                            │
│  • ドメインロジック                                          │
│  • クロスシャード操作                                        │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│      Repository Layer (internal/repository)                  │
│  • DmUserRepository                                         │
│  • データアクセスの抽象化                                    │
│  • CRUD操作                                                 │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│           DB Layer (internal/db)                              │
│  • GroupManager                                             │
│  • シャーディング戦略                                        │
│  • 接続プール管理                                            │
└────────────────────────┬────────────────────────────────────┘
                         │
          ┌──────────────┴──────────────┐
          ▼                              ▼
    ┌─────────┐                    ┌─────────┐
    │ Shard 1 │                    │ Shard 2 │
    └─────────┘                    └─────────┘
```

### 2.2 データフロー

```
main.go
  ↓
コマンドライン引数の解析（flag.Parse()）
  ↓
引数のバリデーション（validateLimit()）
  ↓
設定ファイルの読み込み（config.Load()）
  ↓
GroupManagerの初期化（db.NewGroupManager()）
  ↓
Repository層の初期化（repository.NewDmUserRepository()）
  ↓
Service層の初期化（service.NewDmUserService()）
  ↓
Usecase層の初期化（cli.NewListDmUsersUsecase()）
  ↓
usecase.ListDmUsers(ctx, limit, offset)
  ↓
service.ListDmUsers(ctx, limit, offset)
  ↓
repository.List(ctx, limit, offset) [クロスシャードクエリ]
  ↓
[]*model.DmUser を返却
  ↓
結果の出力（printDmUsersTSV()）
```

### 2.3 レイヤー構造の比較

#### 修正前
```
main.go
  ↓ (直接呼び出し)
service.DmUserService.ListDmUsers()
  ↓
repository.DmUserRepository.List()
  ↓
model.DmUser
```

#### 修正後
```
main.go
  ↓
usecase/cli.ListDmUsersUsecase.ListDmUsers()
  ↓
service.DmUserService.ListDmUsers()
  ↓
repository.DmUserRepository.List()
  ↓
model.DmUser
```

## 3. 詳細設計

### 3.1 CLI用usecase層の設計

#### 3.1.1 ディレクトリ構造

```
server/internal/usecase/cli/
├── list_dm_users_usecase.go
└── list_dm_users_usecase_test.go
```

#### 3.1.2 `list_dm_users_usecase.go`の設計

**ファイルパス**: `server/internal/usecase/cli/list_dm_users_usecase.go`

**実装内容**:

```go
package cli

import (
	"context"

	"github.com/taku-o/go-webdb-template/internal/model"
	"github.com/taku-o/go-webdb-template/internal/usecase"
)

// ListDmUsersUsecase はCLI用のdm_user一覧取得usecase
type ListDmUsersUsecase struct {
	dmUserService usecase.DmUserServiceInterface
}

// NewListDmUsersUsecase は新しいListDmUsersUsecaseを作成
func NewListDmUsersUsecase(dmUserService usecase.DmUserServiceInterface) *ListDmUsersUsecase {
	return &ListDmUsersUsecase{
		dmUserService: dmUserService,
	}
}

// ListDmUsers はユーザー一覧を取得
func (u *ListDmUsersUsecase) ListDmUsers(ctx context.Context, limit, offset int) ([]*model.DmUser, error) {
	return u.dmUserService.ListDmUsers(ctx, limit, offset)
}
```

**設計のポイント**:
- 既存の`DmUserServiceInterface`を使用（`internal/usecase/dm_user_usecase.go`で定義済み）
- コンストラクタで依存関係を注入
- service層のメソッドをそのまま呼び出す（CLI用の特別な処理は不要）
- エラーハンドリングはservice層から返されたエラーをそのまま返す

#### 3.1.3 `list_dm_users_usecase_test.go`の設計

**ファイルパス**: `server/internal/usecase/cli/list_dm_users_usecase_test.go`

**実装内容**:

```go
package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/taku-o/go-webdb-template/internal/model"
	"github.com/taku-o/go-webdb-template/internal/usecase"
)

// MockDmUserServiceInterface はDmUserServiceInterfaceのモック
type MockDmUserServiceInterface struct {
	ListDmUsersFunc func(ctx context.Context, limit, offset int) ([]*model.DmUser, error)
}

func (m *MockDmUserServiceInterface) CreateDmUser(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error) {
	return nil, nil
}

func (m *MockDmUserServiceInterface) GetDmUser(ctx context.Context, id string) (*model.DmUser, error) {
	return nil, nil
}

func (m *MockDmUserServiceInterface) ListDmUsers(ctx context.Context, limit, offset int) ([]*model.DmUser, error) {
	if m.ListDmUsersFunc != nil {
		return m.ListDmUsersFunc(ctx, limit, offset)
	}
	return nil, nil
}

func (m *MockDmUserServiceInterface) UpdateDmUser(ctx context.Context, id string, req *model.UpdateDmUserRequest) (*model.DmUser, error) {
	return nil, nil
}

func (m *MockDmUserServiceInterface) DeleteDmUser(ctx context.Context, id string) error {
	return nil
}

func TestListDmUsersUsecase_ListDmUsers(t *testing.T) {
	tests := []struct {
		name        string
		limit       int
		offset      int
		mockFunc    func(ctx context.Context, limit, offset int) ([]*model.DmUser, error)
		wantUsers   []*model.DmUser
		wantError   bool
		expectedErr string
	}{
		{
			name:   "success with users",
			limit:  20,
			offset: 0,
			mockFunc: func(ctx context.Context, limit, offset int) ([]*model.DmUser, error) {
				return []*model.DmUser{
					{ID: "1", Name: "User 1", Email: "user1@example.com"},
					{ID: "2", Name: "User 2", Email: "user2@example.com"},
				}, nil
			},
			wantUsers: []*model.DmUser{
				{ID: "1", Name: "User 1", Email: "user1@example.com"},
				{ID: "2", Name: "User 2", Email: "user2@example.com"},
			},
			wantError: false,
		},
		{
			name:   "success with empty list",
			limit:  20,
			offset: 0,
			mockFunc: func(ctx context.Context, limit, offset int) ([]*model.DmUser, error) {
				return []*model.DmUser{}, nil
			},
			wantUsers: []*model.DmUser{},
			wantError: false,
		},
		{
			name:   "service error",
			limit:  20,
			offset: 0,
			mockFunc: func(ctx context.Context, limit, offset int) ([]*model.DmUser, error) {
				return nil, errors.New("database error")
			},
			wantUsers:   nil,
			wantError:   true,
			expectedErr: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockDmUserServiceInterface{
				ListDmUsersFunc: tt.mockFunc,
			}

			usecase := NewListDmUsersUsecase(mockService)

			ctx := context.Background()
			gotUsers, err := usecase.ListDmUsers(ctx, tt.limit, tt.offset)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, gotUsers)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUsers, gotUsers)
			}
		})
	}
}
```

**設計のポイント**:
- 関数ポインタを使用してモックを実装（既存のテストパターンに合わせる）
- テーブル駆動テストを使用
- 正常系と異常系の両方をテスト
- `github.com/stretchr/testify/assert`を使用してアサーション

### 3.2 main.goの簡素化設計

#### 3.2.1 修正後のmain.goの構造

**ファイルパス**: `server/cmd/list-dm-users/main.go`

**実装内容**:

```go
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/taku-o/go-webdb-template/internal/config"
	"github.com/taku-o/go-webdb-template/internal/db"
	"github.com/taku-o/go-webdb-template/internal/model"
	"github.com/taku-o/go-webdb-template/internal/repository"
	"github.com/taku-o/go-webdb-template/internal/service"
	"github.com/taku-o/go-webdb-template/internal/usecase/cli"
)

func main() {
	// コマンドライン引数の解析
	limit := flag.Int("limit", 20, "Number of users to output (default: 20, max: 100)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	// 引数のバリデーション
	validatedLimit, err, warning := validateLimit(*limit)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	if warning {
		log.Printf("Warning: limit exceeds maximum (100), using 100")
	}

	// 設定ファイルの読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// GroupManagerの初期化
	groupManager, err := db.NewGroupManager(cfg)
	if err != nil {
		log.Fatalf("Failed to create group manager: %v", err)
	}
	defer groupManager.CloseAll()

	// すべてのデータベースへの接続確認
	if err := groupManager.PingAll(); err != nil {
		log.Fatalf("Failed to ping databases: %v", err)
	}

	// Repository層の初期化
	dmUserRepo := repository.NewDmUserRepository(groupManager)

	// Service層の初期化
	dmUserService := service.NewDmUserService(dmUserRepo)

	// Usecase層の初期化
	listDmUsersUsecase := cli.NewListDmUsersUsecase(dmUserService)

	// ユーザー一覧の取得
	ctx := context.Background()
	dmUsers, err := listDmUsersUsecase.ListDmUsers(ctx, validatedLimit, 0)
	if err != nil {
		log.Fatalf("Failed to list users: %v", err)
	}

	// limit件に制限（クロスシャードクエリのため）
	if len(dmUsers) > validatedLimit {
		dmUsers = dmUsers[:validatedLimit]
	}

	// TSV形式での出力
	printDmUsersTSV(dmUsers)

	os.Exit(0)
}

// validateLimit validates the limit parameter and returns the validated limit,
// an error if invalid, and a boolean indicating if a warning was issued.
func validateLimit(limit int) (int, error, bool) {
	if limit < 1 {
		return 0, errors.New("limit must be at least 1"), false
	}
	if limit > 100 {
		return 100, nil, true
	}
	return limit, nil, false
}

// printDmUsersTSV prints dm_users in TSV format to stdout.
func printDmUsersTSV(dmUsers []*model.DmUser) {
	// ヘッダー行の出力
	fmt.Println("ID\tName\tEmail\tCreatedAt\tUpdatedAt")

	// 各ユーザー情報の出力
	for _, dmUser := range dmUsers {
		fmt.Printf("%s\t%s\t%s\t%s\t%s\n",
			dmUser.ID,
			dmUser.Name,
			dmUser.Email,
			dmUser.CreatedAt.Format(time.RFC3339),
			dmUser.UpdatedAt.Format(time.RFC3339),
		)
	}
}
```

**変更点**:
1. `internal/usecase/cli`パッケージをインポート
2. usecase層の初期化を追加（`cli.NewListDmUsersUsecase(dmUserService)`）
3. service層の直接呼び出しをusecase層の呼び出しに変更
4. `validateLimit()`関数と`printDmUsersTSV()`関数は維持（変更なし）

#### 3.2.2 バリデーション関数（変更なし）

`validateLimit()`関数は既存の実装をそのまま維持します。

#### 3.2.3 出力関数（変更なし）

`printDmUsersTSV()`関数は既存の実装をそのまま維持します。

### 3.3 依存関係の注入設計

#### 3.3.1 初期化の順序

```
1. config.Load()
   ↓
2. db.NewGroupManager(cfg)
   ↓
3. repository.NewDmUserRepository(groupManager)
   ↓
4. service.NewDmUserService(dmUserRepo)
   ↓
5. cli.NewListDmUsersUsecase(dmUserService)
```

#### 3.3.2 依存関係の図

```
ListDmUsersUsecase
  └── DmUserServiceInterface
        └── DmUserRepositoryInterface
              └── GroupManager
                    └── Config
```

### 3.4 テスト設計

#### 3.4.1 usecase層のテスト

**テストファイル**: `server/internal/usecase/cli/list_dm_users_usecase_test.go`

**テストケース**:
1. 正常系: ユーザー一覧が取得できる場合
2. 正常系: 空のリストが返される場合
3. 異常系: service層でエラーが発生した場合

**テスト手法**:
- モックを使用してservice層をモック化
- テーブル駆動テストを使用
- `github.com/stretchr/testify`を使用

#### 3.4.2 main.goのテスト（既存テストの維持）

**テストファイル**: `server/cmd/list-dm-users/main_test.go`

既存のテストはそのまま維持します：
- `TestPrintDmUsersTSV`: TSV形式での出力のテスト
- `TestValidateLimit`: バリデーション関数のテスト

**注意点**:
- usecase層を使用するように変更したが、テスト対象の関数（`validateLimit()`、`printDmUsersTSV()`）は変更されないため、既存のテストはそのまま動作する

### 3.5 ドキュメント更新の設計

#### 3.5.1 `docs/Architecture.md`の更新

**更新内容**:
1. CLIコマンドのレイヤー構造を追加
2. CLIコマンドのアーキテクチャ図を更新（usecase層を含む）
3. CLI用usecase層の説明を追加

**追加するセクション**:

```markdown
### CLI Layer (`cmd/list-dm-users`)

**Location**: `cmd/list-dm-users/main.go`

**Responsibilities**:
- Command-line argument parsing
- Input validation
- Configuration loading
- Layer initialization (Repository → Service → Usecase)
- Usecase layer invocation
- Output formatting (TSV)

**Key Components**:
- `main()`: Entry point
- `validateLimit()`: Input validation
- `printDmUsersTSV()`: Output formatting

### CLI Usecase Layer (`internal/usecase/cli`)

**Location**: `internal/usecase/cli/`

**Responsibilities**:
- CLI-specific business logic coordination
- Service layer invocation for CLI commands

**Key Components**:
- `ListDmUsersUsecase`: User list retrieval for CLI

**Constraints**:
- Uses existing service layer interfaces
- Does not contain domain logic (delegates to service layer)
```

#### 3.5.2 `docs/Project-Structure.md`の更新

**更新内容**:
1. `server/internal/usecase/cli`ディレクトリを追加
2. `server/internal/usecase/cli/list_dm_users_usecase.go`を追加

**追加する行**:

```markdown
│   │   ├── usecase/            # ビジネスロジック層
│   │   │   ├── dm_user_usecase.go
│   │   │   ├── dm_post_usecase.go
│   │   │   ├── email_usecase.go
│   │   │   ├── cli/            # CLI用usecase層
│   │   │   │   └── list_dm_users_usecase.go
│   │   │   └── ...
```

#### 3.5.3 `docs/Command-Line-Tool.md`の更新

**更新内容**:
1. アーキテクチャ図を更新（usecase層を追加）
2. レイヤー構造の説明を更新

**更新するアーキテクチャ図**:

```markdown
## アーキテクチャ

CLIツールは既存のレイヤードアーキテクチャを再利用しています。

```
┌─────────────────────────────────────────────────────────────┐
│                    list-dm-users コマンド                     │
│                    (cmd/list-dm-users/main.go)                │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│         Usecase層 (internal/usecase/cli)                     │
│         - ListDmUsersUsecase.ListDmUsers()                  │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Service層 (internal/service)                    │
│              - DmUserService.ListDmUsers()                  │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Repository層 (internal/repository)              │
│              - DmUserRepository.List()                      │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              DB層 (internal/db)                              │
│              - GroupManager                                  │
└────────────────────────┬──────────────────────────────────┘
                         │
          ┌──────────────┴──────────────┐
          ▼                              ▼
    ┌─────────┐                    ┌─────────┐
    │ Shard 1 │                    │ Shard 2 │
    └─────────┘                    └─────────┘
```
```

#### 3.5.4 `.kiro/steering/structure.md`の更新

**更新内容**:
1. `server/internal/usecase/cli`ディレクトリを追加
2. CLI用usecase層の説明を追加

**追加する行**:

```markdown
│   │   ├── usecase/            # ビジネスロジック層
│   │   │   ├── dm_user_usecase.go
│   │   │   ├── dm_post_usecase.go
│   │   │   ├── email_usecase.go
│   │   │   ├── cli/            # CLI用usecase層
│   │   │   │   └── list_dm_users_usecase.go
│   │   │   └── ...
```

## 4. インターフェース設計

### 4.1 既存インターフェースの使用

**使用するインターフェース**: `DmUserServiceInterface`

**定義場所**: `server/internal/usecase/dm_user_usecase.go`

**インターフェース定義**:

```go
type DmUserServiceInterface interface {
	CreateDmUser(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error)
	GetDmUser(ctx context.Context, id string) (*model.DmUser, error)
	ListDmUsers(ctx context.Context, limit, offset int) ([]*model.DmUser, error)
	UpdateDmUser(ctx context.Context, id string, req *model.UpdateDmUserRequest) (*model.DmUser, error)
	DeleteDmUser(ctx context.Context, id string) error
}
```

**設計方針**:
- 新規インターフェースは作成しない
- 既存の`DmUserServiceInterface`をそのまま使用
- これにより、既存のservice層の実装をそのまま使用できる

## 5. エラーハンドリング設計

### 5.1 エラーハンドリングの流れ

```
usecase層
  ↓ (エラーをそのまま返す)
service層
  ↓ (エラーをラップして返す)
repository層
  ↓ (データベースエラーを返す)
main.go
  ↓ (log.Fatalf()でエラーを出力して終了)
```

### 5.2 エラーメッセージ

既存のエラーメッセージを維持します：
- 設定ファイルの読み込みエラー: `"Failed to load config: %v"`
- GroupManagerの初期化エラー: `"Failed to create group manager: %v"`
- データベース接続エラー: `"Failed to ping databases: %v"`
- ユーザー一覧取得エラー: `"Failed to list users: %v"`
- 引数エラー: `"Error: %v"`

## 6. パフォーマンス設計

### 6.1 オーバーヘッド

usecase層の追加によるオーバーヘッドは無視できるレベルです：
- usecase層はservice層のメソッドをそのまま呼び出すだけ
- 追加の処理は行わない
- メモリ使用量の増加も最小限

### 6.2 既存機能の維持

既存のservice層とrepository層のパフォーマンスは維持されます：
- 既存の実装をそのまま使用
- 追加の処理は行わない

## 7. セキュリティ設計

### 7.1 入力検証

既存の入力検証を維持します：
- `validateLimit()`関数でlimit値の検証
- 最小値チェック（1以上）
- 最大値チェック（100以下）

### 7.2 エラーメッセージ

既存のエラーメッセージを維持します：
- 機密情報を含まないエラーメッセージ
- ユーザーフレンドリーなエラーメッセージ

## 8. テスト戦略

### 8.1 テストレベル

1. **ユニットテスト**: usecase層の単体テスト
2. **統合テスト**: 既存のテスト（main_test.go）を維持

### 8.2 テストカバレッジ

- usecase層: 80%以上のカバレッジを目標
- 既存のテスト: 全て通過することを確認

### 8.3 テスト実行

```bash
# usecase層のテスト
go test -v ./internal/usecase/cli/...

# main.goのテスト
go test -v ./cmd/list-dm-users/...

# カバレッジ確認
go test -cover ./internal/usecase/cli/...
```

## 9. 実装順序

### 9.1 実装の優先順位

1. **usecase層の実装**（最優先）
   - `server/internal/usecase/cli/list_dm_users_usecase.go`の作成
   - `server/internal/usecase/cli/list_dm_users_usecase_test.go`の作成

2. **main.goの修正**
   - usecase層を使用するように修正
   - 既存のテストが通過することを確認

3. **ドキュメントの更新**
   - `docs/Architecture.md`の更新
   - `docs/Project-Structure.md`の更新
   - `docs/Command-Line-Tool.md`の更新
   - `.kiro/steering/structure.md`の更新

### 9.2 実装の注意点

1. **既存コードの維持**: 既存のservice、repository、model層は変更しない
2. **テストの維持**: 既存のテスト（main_test.go）は全て通過することを確認
3. **後方互換性**: 既存のCLIコマンドの動作（引数、出力形式、エラーメッセージ）を維持

## 10. リスクと対策

### 10.1 リスク

1. **既存テストの失敗**: main.goを修正したことで既存のテストが失敗する可能性
2. **パフォーマンスの劣化**: usecase層の追加によるパフォーマンスの劣化

### 10.2 対策

1. **既存テストの維持**: 既存のテスト関数（`validateLimit()`、`printDmUsersTSV()`）は変更しないため、既存のテストはそのまま動作する
2. **パフォーマンステスト**: 既存のパフォーマンステストを実行して、パフォーマンスの劣化がないことを確認

## 11. 参考情報

### 11.1 既存実装の参考

- `server/internal/usecase/dm_user_usecase.go`: 既存のusecase層の実装パターン
- `server/cmd/list-dm-users/main.go`: 既存のCLI実装
- `server/internal/service/dm_user_service.go`: 既存のservice層の実装
- `server/internal/repository/dm_user_repository.go`: 既存のrepository層の実装

### 11.2 関連ドキュメント

- `docs/Architecture.md`: アーキテクチャドキュメント
- `docs/Project-Structure.md`: プロジェクト構造ドキュメント
- `docs/Command-Line-Tool.md`: CLIツールドキュメント
- `.kiro/steering/structure.md`: ファイル組織とコードパターン
- `.kiro/steering/tech.md`: 技術スタックとアーキテクチャ
