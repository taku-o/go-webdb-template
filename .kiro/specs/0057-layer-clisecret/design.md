# server/cmd/generate-secretの構造修正の設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、`server/cmd/generate-secret`の実装を、APIサーバーと同じレイヤー構造（usecase -> service -> auth）に変更するための詳細設計を定義する。これにより、CLIコマンドとAPIサーバーで一貫したアーキテクチャを実現し、秘密鍵生成処理を共通化してコードの保守性と再利用性を向上させる。

### 1.2 設計の範囲
- 秘密鍵生成処理の共通化（`server/internal/auth/secret.go`）の設計
- Service層（`server/internal/service/secret_service.go`）の設計
- CLI用usecase層（`server/internal/usecase/cli/generate_secret_usecase.go`）の設計
- `server/cmd/generate-secret/main.go`の簡素化設計
- 依存関係の注入設計
- テスト設計
- ドキュメント更新の設計

### 1.3 設計方針
- **一貫性**: APIサーバーと同じレイヤー構造を採用
- **共通化**: 秘密鍵生成処理を共通ライブラリとして実装
- **責務の明確化**: 各レイヤーの責務を明確に分離
- **テスト容易性**: 各層を独立してテストできる設計
- **後方互換性**: 既存のCLIコマンドの動作（出力形式、エラーメッセージ）を維持
- **セキュリティ**: `crypto/rand`パッケージを使用して安全な乱数生成を行う

## 2. アーキテクチャ設計

### 2.1 全体構成

```
┌─────────────────────────────────────────────────────────────┐
│          CLI Layer (cmd/generate-secret/main.go)              │
│  • エントリーポイント                                         │
│  • レイヤーの初期化（Service → Usecase）                     │
│  • usecase層の呼び出し                                       │
│  • 結果の出力（標準出力）                                     │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│      Usecase Layer (internal/usecase/cli)                     │
│  • GenerateSecretUsecase                                     │
│  • ビジネスロジックの調整（CLI用）                           │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│         Service Layer (internal/service)                        │
│  • SecretService                                             │
│  • ドメインロジック                                           │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│         Auth Layer (internal/auth)                             │
│  • GenerateSecretKey()                                        │
│  • 秘密鍵生成処理（共通ライブラリ）                           │
│  • crypto/rand + encoding/base64                             │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 データフロー

```
main.go
  ↓
Service層の初期化（service.NewSecretService()）
  ↓
Usecase層の初期化（cli.NewGenerateSecretUsecase(secretService)）
  ↓
usecase.GenerateSecret(ctx)
  ↓
service.GenerateSecretKey(ctx)
  ↓
auth.GenerateSecretKey()
  ↓
32バイトのランダムな秘密鍵を生成（crypto/rand）
  ↓
Base64エンコード（encoding/base64）
  ↓
string を返却
  ↓
結果の出力（標準出力にBase64エンコードされた秘密鍵を出力）
```

### 2.3 レイヤー構造の比較

#### 修正前
```
main.go
  ↓ (直接実装)
crypto/rand.Read()
  ↓
encoding/base64.EncodeToString()
  ↓
標準出力に出力
```

#### 修正後
```
main.go
  ↓
usecase/cli.GenerateSecretUsecase.GenerateSecret()
  ↓
service.SecretService.GenerateSecretKey()
  ↓
auth.GenerateSecretKey()
  ↓
標準出力に出力
```

## 3. 詳細設計

### 3.1 秘密鍵生成処理の共通化設計

#### 3.1.1 `secret.go`の設計

**ファイルパス**: `server/internal/auth/secret.go`

**実装内容**:

```go
package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateSecretKey は32バイト（256ビット）のランダムな秘密鍵を生成してBase64エンコードして返す
func GenerateSecretKey() (string, error) {
	// 32バイト（256ビット）のランダムな秘密鍵を生成
	secretKey := make([]byte, 32)
	if _, err := rand.Read(secretKey); err != nil {
		return "", fmt.Errorf("failed to generate secret key: %w", err)
	}

	// Base64エンコード
	encoded := base64.StdEncoding.EncodeToString(secretKey)

	return encoded, nil
}
```

**設計のポイント**:
- `crypto/rand`パッケージを使用して安全な乱数生成を行う
- 32バイト（256ビット）のランダムな秘密鍵を生成
- Base64エンコードして文字列として返す
- エラーハンドリング: 乱数生成に失敗した場合は適切なエラーを返す
- 既存の`main.go`の実装と同じロジックを使用

#### 3.1.2 `secret_test.go`の設計

**ファイルパス**: `server/internal/auth/secret_test.go`

**実装内容**:

```go
package auth

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSecretKey(t *testing.T) {
	tests := []struct {
		name      string
		wantError bool
	}{
		{
			name:      "success",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateSecretKey()

			if tt.wantError {
				assert.Error(t, err)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, got)

				// Base64デコードして32バイトであることを確認
				decoded, err := base64.StdEncoding.DecodeString(got)
				assert.NoError(t, err)
				assert.Equal(t, 32, len(decoded))
			}
		})
	}
}

func TestGenerateSecretKey_Uniqueness(t *testing.T) {
	// 複数回生成して、それぞれ異なる秘密鍵が生成されることを確認
	secret1, err1 := GenerateSecretKey()
	assert.NoError(t, err1)

	secret2, err2 := GenerateSecretKey()
	assert.NoError(t, err2)

	assert.NotEqual(t, secret1, secret2, "generated secrets should be unique")
}
```

**設計のポイント**:
- 正常系のテスト: 秘密鍵が正常に生成されることを確認
- Base64デコードして32バイトであることを確認
- 一意性のテスト: 複数回生成して、それぞれ異なる秘密鍵が生成されることを確認
- `github.com/stretchr/testify/assert`を使用してアサーション

### 3.2 Service層の設計

#### 3.2.1 `secret_service.go`の設計

**ファイルパス**: `server/internal/service/secret_service.go`

**実装内容**:

```go
package service

import (
	"context"

	"github.com/taku-o/go-webdb-template/internal/auth"
)

// SecretServiceInterface は秘密鍵生成サービスのインターフェース
type SecretServiceInterface interface {
	GenerateSecretKey(ctx context.Context) (string, error)
}

// SecretService は秘密鍵生成のビジネスロジックを担当
type SecretService struct{}

// NewSecretService は新しいSecretServiceを作成
func NewSecretService() *SecretService {
	return &SecretService{}
}

// GenerateSecretKey は秘密鍵を生成
func (s *SecretService) GenerateSecretKey(ctx context.Context) (string, error) {
	return auth.GenerateSecretKey()
}
```

**設計のポイント**:
- `SecretServiceInterface`を新規作成（usecase層で使用するため）
- 現時点では依存関係なし（将来的に拡張可能な設計）
- 共通ライブラリ（`auth.GenerateSecretKey()`）を呼び出して結果を返す
- エラーハンドリング: 共通ライブラリから返されたエラーをそのまま返す（エラーのラップは不要）
- `context.Context`を受け取る（将来の拡張性のため）

#### 3.2.2 `secret_service_test.go`の設計

**ファイルパス**: `server/internal/service/secret_service_test.go`

**実装内容**:

```go
package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecretService_GenerateSecretKey(t *testing.T) {
	tests := []struct {
		name      string
		wantError bool
	}{
		{
			name:      "success",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewSecretService()
			ctx := context.Background()

			got, err := service.GenerateSecretKey(ctx)

			if tt.wantError {
				assert.Error(t, err)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, got)
			}
		})
	}
}
```

**設計のポイント**:
- 正常系のテスト: 秘密鍵が正常に生成されることを確認
- `github.com/stretchr/testify/assert`を使用してアサーション

### 3.3 CLI用usecase層の設計

#### 3.3.1 ディレクトリ構造

```
server/internal/usecase/cli/
├── list_dm_users_usecase.go
├── list_dm_users_usecase_test.go
├── generate_secret_usecase.go
└── generate_secret_usecase_test.go
```

#### 3.3.2 `generate_secret_usecase.go`の設計

**ファイルパス**: `server/internal/usecase/cli/generate_secret_usecase.go`

**実装内容**:

```go
package cli

import (
	"context"

	"github.com/taku-o/go-webdb-template/internal/service"
)

// GenerateSecretUsecase はCLI用の秘密鍵生成usecase
type GenerateSecretUsecase struct {
	secretService service.SecretServiceInterface
}

// NewGenerateSecretUsecase は新しいGenerateSecretUsecaseを作成
func NewGenerateSecretUsecase(secretService service.SecretServiceInterface) *GenerateSecretUsecase {
	return &GenerateSecretUsecase{
		secretService: secretService,
	}
}

// GenerateSecret は秘密鍵を生成
func (u *GenerateSecretUsecase) GenerateSecret(ctx context.Context) (string, error) {
	return u.secretService.GenerateSecretKey(ctx)
}
```

**設計のポイント**:
- `SecretServiceInterface`を使用（依存関係の注入）
- コンストラクタで依存関係を注入
- service層のメソッドをそのまま呼び出す（CLI用の特別な処理は不要）
- エラーハンドリング: service層から返されたエラーをそのまま返す（エラーのラップは不要）

#### 3.3.3 `generate_secret_usecase_test.go`の設計

**ファイルパス**: `server/internal/usecase/cli/generate_secret_usecase_test.go`

**実装内容**:

```go
package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/taku-o/go-webdb-template/internal/service"
)

// MockSecretServiceInterface はSecretServiceInterfaceのモック
type MockSecretServiceInterface struct {
	GenerateSecretKeyFunc func(ctx context.Context) (string, error)
}

func (m *MockSecretServiceInterface) GenerateSecretKey(ctx context.Context) (string, error) {
	if m.GenerateSecretKeyFunc != nil {
		return m.GenerateSecretKeyFunc(ctx)
	}
	return "", nil
}

func TestGenerateSecretUsecase_GenerateSecret(t *testing.T) {
	tests := []struct {
		name        string
		mockFunc    func(ctx context.Context) (string, error)
		wantSecret  string
		wantError   bool
		expectedErr string
	}{
		{
			name: "success",
			mockFunc: func(ctx context.Context) (string, error) {
				return "test-secret-key", nil
			},
			wantSecret: "test-secret-key",
			wantError:  false,
		},
		{
			name: "service error",
			mockFunc: func(ctx context.Context) (string, error) {
				return "", errors.New("failed to generate secret key")
			},
			wantSecret:   "",
			wantError:   true,
			expectedErr: "failed to generate secret key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockSecretServiceInterface{
				GenerateSecretKeyFunc: tt.mockFunc,
			}

			usecase := NewGenerateSecretUsecase(mockService)

			ctx := context.Background()
			gotSecret, err := usecase.GenerateSecret(ctx)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Empty(t, gotSecret)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantSecret, gotSecret)
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

### 3.4 main.goの簡素化設計

#### 3.4.1 修正後のmain.goの構造

**ファイルパス**: `server/cmd/generate-secret/main.go`

**実装内容**:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/taku-o/go-webdb-template/internal/service"
	"github.com/taku-o/go-webdb-template/internal/usecase/cli"
)

func main() {
	// Service層の初期化
	secretService := service.NewSecretService()

	// Usecase層の初期化
	generateSecretUsecase := cli.NewGenerateSecretUsecase(secretService)

	// 秘密鍵の生成
	ctx := context.Background()
	secretKey, err := generateSecretUsecase.GenerateSecret(ctx)
	if err != nil {
		log.Fatalf("Failed to generate secret key: %v", err)
	}

	// 標準出力に表示
	fmt.Println(secretKey)

	os.Exit(0)
}
```

**変更点**:
1. `internal/service`パッケージをインポート
2. `internal/usecase/cli`パッケージをインポート
3. Service層の初期化を追加（`service.NewSecretService()`）
4. Usecase層の初期化を追加（`cli.NewGenerateSecretUsecase(secretService)`）
5. 秘密鍵生成処理をusecase層の呼び出しに変更
6. 既存の出力形式（標準出力にBase64エンコードされた秘密鍵を出力）を維持
7. 既存のエラーハンドリング（`log.Fatalf`）を維持

**削除した内容**:
- `crypto/rand`と`encoding/base64`の直接使用（共通ライブラリに移動）
- 秘密鍵生成処理の直接実装（usecase層に移動）

#### 3.4.2 ビルド方法

**ビルドコマンド**:
```bash
cd server
go build -o bin/generate-secret ./cmd/generate-secret
```

**ビルド出力先**: `server/bin/generate-secret`

**実行方法**:
```bash
cd server
./bin/generate-secret
```

**注意事項**:
- ビルド出力先は`server/bin/`ディレクトリに統一
- `server/bin/`ディレクトリは`.gitignore`に含まれているため、ビルド成果物はGitにコミットされない
- ビルド前に`server/bin/`ディレクトリが存在しない場合は作成する必要がある

### 3.5 依存関係の注入設計

#### 3.5.1 初期化の順序

```
1. service.NewSecretService()
   ↓
2. cli.NewGenerateSecretUsecase(secretService)
   ↓
3. usecase.GenerateSecret(ctx)
```

#### 3.5.2 依存関係の図

```
GenerateSecretUsecase
  └── SecretServiceInterface
        └── (依存なし)
              └── auth.GenerateSecretKey()
```

## 4. テスト設計

### 4.1 秘密鍵生成処理のテスト

**テストファイル**: `server/internal/auth/secret_test.go`

**テストケース**:
1. 正常系: 秘密鍵が正常に生成される場合
2. 正常系: Base64デコードして32バイトであることを確認
3. 正常系: 複数回生成して、それぞれ異なる秘密鍵が生成されることを確認

**テスト手法**:
- テーブル駆動テストを使用
- `github.com/stretchr/testify`を使用

### 4.2 Service層のテスト

**テストファイル**: `server/internal/service/secret_service_test.go`

**テストケース**:
1. 正常系: 秘密鍵が正常に生成される場合

**テスト手法**:
- テーブル駆動テストを使用
- `github.com/stretchr/testify`を使用

### 4.3 usecase層のテスト

**テストファイル**: `server/internal/usecase/cli/generate_secret_usecase_test.go`

**テストケース**:
1. 正常系: 秘密鍵が正常に生成される場合
2. 異常系: service層でエラーが発生した場合

**テスト手法**:
- モックを使用してservice層をモック化
- テーブル駆動テストを使用
- `github.com/stretchr/testify`を使用

### 4.4 main.goのテスト

**テストファイル**: なし（main.goのテストは不要）

**理由**:
- main.goはエントリーポイントと入出力制御のみを担当
- ビジネスロジックは全てusecase層、service層、auth層に移動
- 各層のテストで十分にカバーされる

## 5. ドキュメント更新の設計

### 5.1 アーキテクチャドキュメントの更新

**修正対象**: `docs/Architecture.md`

**更新内容**:
- CLIコマンドのレイヤー構造を追加（usecase層を含む）
- CLIコマンドのアーキテクチャ図を更新
- CLI用usecase層の説明を追加
- 秘密鍵生成処理の共通ライブラリの説明を追加

### 5.2 プロジェクト構造ドキュメントの更新

**修正対象**: `docs/Project-Structure.md`

**更新内容**:
- `server/internal/usecase/cli/generate_secret_usecase.go`を追加
- `server/internal/service/secret_service.go`を追加
- `server/internal/auth/secret.go`を追加

### 5.3 CLIツールドキュメントの更新

**修正対象**: `docs/Command-Line-Tool.md`（存在する場合）

**更新内容**:
- CLIコマンドのアーキテクチャ図を更新（usecase層を追加）
- レイヤー構造の説明を更新（main.go → usecase → service → auth → 出力）

### 5.4 ファイル組織ドキュメントの更新

**修正対象**: `.kiro/steering/structure.md`

**更新内容**:
- `server/internal/usecase/cli/generate_secret_usecase.go`を追加
- `server/internal/service/secret_service.go`を追加
- `server/internal/auth/secret.go`を追加

## 6. 実装上の注意事項

### 6.1 秘密鍵生成処理の共通化
- **セキュリティ**: `crypto/rand`パッケージを使用して安全な乱数生成を行う
- **出力形式**: Base64エンコードされた文字列を返す（既存の動作を維持）
- **エラーハンドリング**: 乱数生成に失敗した場合は適切なエラーを返す
- **既存実装との互換性**: 既存の`main.go`の実装と同じロジックを使用

### 6.2 Service層の実装
- **インターフェースの定義**: `SecretServiceInterface`を新規作成
- **依存関係の注入**: コンストラクタで依存関係を注入（現時点では依存関係なし）
- **エラーハンドリング**: 共通ライブラリから返されたエラーをそのまま返す（エラーのラップは不要）
- **将来の拡張性**: `context.Context`を受け取る（将来の拡張性のため）

### 6.3 usecase層の実装
- **インターフェースの使用**: 新規作成する`SecretServiceInterface`を使用
- **依存関係の注入**: コンストラクタで`SecretServiceInterface`を注入
- **エラーハンドリング**: service層から返されたエラーをそのまま返す（エラーのラップは不要）

### 6.4 main.goの修正
- **usecase層の初期化**: Service → Usecaseの順で初期化
- **既存の出力**: 標準出力にBase64エンコードされた秘密鍵を出力（既存の動作を維持）
- **エラーハンドリング**: 既存のエラーハンドリング（`log.Fatalf`）を維持

### 6.5 テストの実装
- **共通ライブラリのテスト**: 秘密鍵生成処理のテストを実装
- **service層のテスト**: 正常系のテストを実装
- **usecase層のテスト**: `SecretServiceInterface`のモックを使用してテスト

### 6.6 ドキュメントの更新
- **アーキテクチャドキュメント**: CLIコマンドのレイヤー構造を明確に記載
- **プロジェクト構造ドキュメント**: 新規作成するファイルを反映
- **CLIツールドキュメント**: アーキテクチャ図を更新してusecase層を含める
- **ファイル組織ドキュメント**: 新規作成するファイルを反映
- **一貫性**: 全てのドキュメントで同じレイヤー構造を記載

## 7. 参考情報

### 7.1 関連ドキュメント
- `docs/Architecture.md`: アーキテクチャドキュメント
- `docs/Project-Structure.md`: プロジェクト構造ドキュメント
- `.kiro/steering/structure.md`: ファイル組織とコードパターン
- `.kiro/steering/tech.md`: 技術スタックとアーキテクチャ

### 7.2 既存実装の参考
- `server/internal/usecase/cli/list_dm_users_usecase.go`: 既存のCLI用usecase層の実装パターン
- `server/cmd/generate-secret/main.go`: 既存のCLI実装
- `server/internal/auth/jwt.go`: 既存のauth層の実装パターン
- `server/internal/service/dm_user_service.go`: 既存のservice層の実装パターン

### 7.3 技術スタック
- **言語**: Go
- **アーキテクチャ**: レイヤードアーキテクチャ（usecase -> service -> auth -> 出力）
- **テスト**: `testing`（標準ライブラリ）、`github.com/stretchr/testify`（アサーション、モック）
- **セキュリティ**: `crypto/rand`（安全な乱数生成）、`encoding/base64`（Base64エンコード）

### 7.4 レイヤー構造の比較

| 項目 | 現在（修正前） | 修正後 |
|------|---------------|--------|
| CLI層 | main.go（秘密鍵生成処理が直接実装） | main.go（エントリーポイント、入出力） |
| Usecase層 | なし | `usecase/cli/GenerateSecretUsecase` |
| Service層 | なし | `service.SecretService` |
| Auth層 | なし | `auth.GenerateSecretKey()` |
| 出力 | main.go内で直接出力 | main.go内で出力（usecase層から取得） |

### 7.5 APIサーバーとの比較

| 項目 | APIサーバー | CLI（修正後） |
|------|------------|--------------|
| エントリーポイント | `server/cmd/server/main.go` | `server/cmd/generate-secret/main.go` |
| バリデーション | API Layer（Handler） | CLI層（main.go、現時点では不要） |
| Usecase層 | `usecase.DmUserUsecase` | `usecase/cli.GenerateSecretUsecase` |
| Service層 | `service.DmUserService` | `service.SecretService` |
| Auth層 | `auth.GeneratePublicAPIKey` | `auth.GenerateSecretKey` |
| 出力 | HTTPレスポンス | 標準出力 |
