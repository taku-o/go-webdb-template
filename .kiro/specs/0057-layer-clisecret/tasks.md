# server/cmd/generate-secretの構造修正の実装タスク一覧

## 概要
`server/cmd/generate-secret`の実装を、APIサーバーと同じレイヤー構造（usecase -> service -> auth）に変更するためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: 秘密鍵生成処理の共通化（auth層）

#### タスク 1.1: secret.goの実装
**目的**: 秘密鍵生成処理を共通ライブラリとして実装する。

**作業内容**:
- `server/internal/auth/secret.go`を作成
- `GenerateSecretKey()`関数を実装
- 32バイト（256ビット）のランダムな秘密鍵を生成してBase64エンコードして返す

**実装内容**:
- ファイルパス: `server/internal/auth/secret.go`
- パッケージ名: `auth`
- 関数定義:
  ```go
  func GenerateSecretKey() (string, error)
  ```
- 実装内容:
  - `crypto/rand`パッケージを使用して32バイトのランダムな秘密鍵を生成
  - `encoding/base64`パッケージを使用してBase64エンコード
  - エラーハンドリング: 乱数生成に失敗した場合は適切なエラーを返す
- 既存の`main.go`の実装と同じロジックを使用

**受け入れ基準**:
- [ ] `server/internal/auth/secret.go`が作成されている
- [ ] `GenerateSecretKey() (string, error)`関数が実装されている
- [ ] 32バイト（256ビット）のランダムな秘密鍵を生成している
- [ ] Base64エンコードして返している
- [ ] `crypto/rand`と`encoding/base64`を使用している
- [ ] エラーハンドリングが適切に実装されている

- _Requirements: 3.1.1, 6.1_
- _Design: 3.1.1_

---

#### タスク 1.2: secret_test.goの実装
**目的**: 秘密鍵生成処理の単体テストを実装する。

**作業内容**:
- `server/internal/auth/secret_test.go`を作成
- 正常系のテストケースを実装
- 一意性のテストケースを実装

**実装内容**:
- ファイルパス: `server/internal/auth/secret_test.go`
- パッケージ名: `auth`
- テストケース:
  1. 正常系: 秘密鍵が正常に生成される場合
  2. 正常系: Base64デコードして32バイトであることを確認
  3. 正常系: 複数回生成して、それぞれ異なる秘密鍵が生成されることを確認
- テーブル駆動テストを使用
- `github.com/stretchr/testify/assert`を使用

**受け入れ基準**:
- [ ] `server/internal/auth/secret_test.go`が作成されている
- [ ] 正常系のテストケースが実装されている
- [ ] Base64デコードして32バイトであることを確認するテストが実装されている
- [ ] 一意性のテストケースが実装されている
- [ ] 全てのテストケースが通過する

- _Requirements: 6.1, 6.7_
- _Design: 3.1.2_

---

### Phase 2: Service層の作成

#### タスク 2.1: secret_service.goの実装
**目的**: 秘密鍵生成用のservice層を実装する。

**作業内容**:
- `server/internal/service/secret_service.go`を作成
- `SecretService`構造体を定義
- `SecretServiceInterface`を定義
- `GenerateSecretKey()`メソッドを実装

**実装内容**:
- ファイルパス: `server/internal/service/secret_service.go`
- パッケージ名: `service`
- インターフェース定義:
  ```go
  type SecretServiceInterface interface {
      GenerateSecretKey(ctx context.Context) (string, error)
  }
  ```
- 構造体定義:
  ```go
  type SecretService struct{}
  ```
- コンストラクタ:
  ```go
  func NewSecretService() *SecretService
  ```
- メソッド:
  ```go
  func (s *SecretService) GenerateSecretKey(ctx context.Context) (string, error)
  ```
- 共通ライブラリ（`auth.GenerateSecretKey()`）を呼び出して結果を返す
- エラーハンドリング: 共通ライブラリから返されたエラーをそのまま返す

**受け入れ基準**:
- [ ] `server/internal/service/secret_service.go`が作成されている
- [ ] `SecretService`構造体が定義されている
- [ ] `SecretServiceInterface`が定義されている
- [ ] `GenerateSecretKey(ctx context.Context) (string, error)`メソッドが実装されている
- [ ] 共通ライブラリ（`auth.GenerateSecretKey()`）を呼び出して結果を返している
- [ ] エラーハンドリングが適切に実装されている

- _Requirements: 3.2.1, 6.2_
- _Design: 3.2.1_

---

#### タスク 2.2: secret_service_test.goの実装
**目的**: Service層の単体テストを実装する。

**作業内容**:
- `server/internal/service/secret_service_test.go`を作成
- 正常系のテストケースを実装

**実装内容**:
- ファイルパス: `server/internal/service/secret_service_test.go`
- パッケージ名: `service`
- テストケース:
  1. 正常系: 秘密鍵が正常に生成される場合
- テーブル駆動テストを使用
- `github.com/stretchr/testify/assert`を使用

**受け入れ基準**:
- [ ] `server/internal/service/secret_service_test.go`が作成されている
- [ ] 正常系のテストケースが実装されている
- [ ] 全てのテストケースが通過する

- _Requirements: 6.2, 6.7_
- _Design: 3.2.2_

---

### Phase 3: CLI用usecase層の作成

#### タスク 3.1: generate_secret_usecase.goの実装
**目的**: CLI用の秘密鍵生成usecaseを実装する。

**作業内容**:
- `server/internal/usecase/cli/generate_secret_usecase.go`を作成
- `GenerateSecretUsecase`構造体を定義
- `SecretServiceInterface`を依存として注入
- `GenerateSecret()`メソッドを実装

**実装内容**:
- ファイルパス: `server/internal/usecase/cli/generate_secret_usecase.go`
- パッケージ名: `cli`
- 構造体定義:
  ```go
  type GenerateSecretUsecase struct {
      secretService service.SecretServiceInterface
  }
  ```
- コンストラクタ:
  ```go
  func NewGenerateSecretUsecase(secretService service.SecretServiceInterface) *GenerateSecretUsecase
  ```
- メソッド:
  ```go
  func (u *GenerateSecretUsecase) GenerateSecret(ctx context.Context) (string, error)
  ```
- service層の`GenerateSecretKey()`を呼び出して結果を返す
- エラーハンドリング: service層から返されたエラーをそのまま返す

**受け入れ基準**:
- [ ] `server/internal/usecase/cli/generate_secret_usecase.go`が作成されている
- [ ] `GenerateSecretUsecase`構造体が定義されている
- [ ] `SecretServiceInterface`を依存として注入している
- [ ] `GenerateSecret(ctx context.Context) (string, error)`メソッドが実装されている
- [ ] service層の`GenerateSecretKey()`を呼び出して結果を返している
- [ ] エラーハンドリングが適切に実装されている

- _Requirements: 3.3.2, 6.3_
- _Design: 3.3.2_

---

#### タスク 3.2: generate_secret_usecase_test.goの実装
**目的**: CLI用usecase層の単体テストを実装する。

**作業内容**:
- `server/internal/usecase/cli/generate_secret_usecase_test.go`を作成
- モックの実装（`MockSecretServiceInterface`）
- テストケースの実装（正常系、異常系）

**実装内容**:
- ファイルパス: `server/internal/usecase/cli/generate_secret_usecase_test.go`
- パッケージ名: `cli`
- モック実装:
  ```go
  type MockSecretServiceInterface struct {
      GenerateSecretKeyFunc func(ctx context.Context) (string, error)
  }
  ```
- テストケース:
  1. 正常系: 秘密鍵が正常に生成される場合
  2. 異常系: service層でエラーが発生した場合
- テーブル駆動テストを使用
- `github.com/stretchr/testify/assert`を使用

**受け入れ基準**:
- [ ] `server/internal/usecase/cli/generate_secret_usecase_test.go`が作成されている
- [ ] モック（`MockSecretServiceInterface`）が実装されている
- [ ] 正常系のテストケースが実装されている
- [ ] 異常系のテストケースが実装されている
- [ ] 全てのテストケースが通過する

- _Requirements: 6.3, 6.7_
- _Design: 3.3.3_

---

### Phase 4: main.goの簡素化

#### タスク 4.1: main.goの修正
**目的**: main.goをエントリーポイントと入出力制御のみに限定する。

**作業内容**:
- `server/cmd/generate-secret/main.go`を修正
- `internal/service`パッケージをインポート
- `internal/usecase/cli`パッケージをインポート
- Service層の初期化を追加
- Usecase層の初期化を追加
- 秘密鍵生成処理をusecase層の呼び出しに変更
- 既存の出力形式を維持

**実装内容**:
- 修正対象: `server/cmd/generate-secret/main.go`
- インポート追加:
  ```go
  "github.com/taku-o/go-webdb-template/internal/service"
  "github.com/taku-o/go-webdb-template/internal/usecase/cli"
  ```
- 初期化の追加:
  ```go
  // Service層の初期化
  secretService := service.NewSecretService()

  // Usecase層の初期化
  generateSecretUsecase := cli.NewGenerateSecretUsecase(secretService)
  ```
- 秘密鍵生成処理の変更:
  - 修正前:
    ```go
    secretKey := make([]byte, 32)
    if _, err := rand.Read(secretKey); err != nil {
        log.Fatalf("Failed to generate secret key: %v", err)
    }
    encoded := base64.StdEncoding.EncodeToString(secretKey)
    fmt.Println(encoded)
    ```
  - 修正後:
    ```go
    ctx := context.Background()
    secretKey, err := generateSecretUsecase.GenerateSecret(ctx)
    if err != nil {
        log.Fatalf("Failed to generate secret key: %v", err)
    }
    fmt.Println(secretKey)
    ```
- 削除する内容:
  - `crypto/rand`と`encoding/base64`の直接使用
  - 秘密鍵生成処理の直接実装

**受け入れ基準**:
- [ ] `internal/service`パッケージがインポートされている
- [ ] `internal/usecase/cli`パッケージがインポートされている
- [ ] Service層が適切に初期化されている
- [ ] Usecase層が適切に初期化されている
- [ ] 秘密鍵生成処理がusecase層の呼び出しに変更されている
- [ ] 既存の出力形式（標準出力にBase64エンコードされた秘密鍵を出力）が維持されている
- [ ] 既存のエラーハンドリング（`log.Fatalf`）が維持されている
- [ ] `crypto/rand`と`encoding/base64`の直接使用が削除されている

- _Requirements: 3.4.1, 3.5.1, 6.4, 6.5_
- _Design: 3.4.1_

---

### Phase 5: ビルドと動作確認

#### タスク 5.1: ビルドの確認
**目的**: ビルドが正常に完了し、バイナリが正しい場所に出力されることを確認する。

**作業内容**:
- `server/bin/`ディレクトリが存在することを確認（存在しない場合は作成）
- ビルドコマンドを実行
- バイナリが`server/bin/generate-secret`に出力されることを確認

**実装内容**:
- ビルドコマンド:
  ```bash
  cd server
  go build -o bin/generate-secret ./cmd/generate-secret
  ```
- 確認項目:
  - `server/bin/generate-secret`が存在する
  - バイナリが実行可能である

**受け入れ基準**:
- [ ] `server/bin/`ディレクトリが存在する（または作成された）
- [ ] ビルドが正常に完了する
- [ ] バイナリが`server/bin/generate-secret`に出力される
- [ ] バイナリが実行可能である

- _Requirements: 6.4_
- _Design: 3.4.2_

---

#### タスク 5.2: 動作確認
**目的**: CLIコマンドが正常に動作することを確認する。

**作業内容**:
- ローカル環境でCLIコマンドを実行
- 既存の出力形式（Base64エンコードされた秘密鍵）が維持されていることを確認
- 既存のエラーメッセージが維持されていることを確認

**実装内容**:
- 実行コマンド:
  ```bash
  cd server
  ./bin/generate-secret
  ```
- 確認項目:
  - 標準出力にBase64エンコードされた秘密鍵が出力される
  - 出力される秘密鍵が32バイト（Base64デコード後）である
  - エラーが発生しない

**受け入れ基準**:
- [ ] ローカル環境でCLIコマンドが正常に動作する
- [ ] 既存の出力形式（Base64エンコードされた秘密鍵）が維持されている
- [ ] 出力される秘密鍵が32バイト（Base64デコード後）である
- [ ] エラーが発生しない

- _Requirements: 6.6_

---

#### タスク 5.3: 既存テストの確認
**目的**: 既存のテストが全て通過することを確認する。

**作業内容**:
- 全てのテストを実行
- 既存のテストが全て通過することを確認

**実装内容**:
- テスト実行:
  ```bash
  cd server
  go test ./...
  ```
- 確認項目:
  - 全てのテストが通過する
  - 新規作成したテストも含めて全て通過する

**受け入れ基準**:
- [ ] 全てのテストが通過する
- [ ] 既存のテストが全て通過する
- [ ] 新規作成したテストも含めて全て通過する

- _Requirements: 6.6, 6.7_

---

### Phase 6: ドキュメントの更新

#### タスク 6.1: Architecture.mdの更新
**目的**: アーキテクチャドキュメントにCLIコマンドのレイヤー構造を追加する。

**作業内容**:
- `docs/Architecture.md`を修正
- CLIコマンドのレイヤー構造を追加（usecase層を含む）
- CLIコマンドのアーキテクチャ図を更新
- CLI用usecase層の説明を追加
- 秘密鍵生成処理の共通ライブラリの説明を追加

**実装内容**:
- 修正対象: `docs/Architecture.md`
- 追加内容:
  - CLIコマンドのレイヤー構造の説明
  - CLIコマンドのアーキテクチャ図（usecase層を含む）
  - CLI用usecase層の説明
  - 秘密鍵生成処理の共通ライブラリ（`auth.GenerateSecretKey()`）の説明

**受け入れ基準**:
- [ ] `docs/Architecture.md`にCLIコマンドのレイヤー構造が追加されている
- [ ] CLIコマンドのアーキテクチャ図が更新されている（usecase層を含む）
- [ ] CLI用usecase層の説明が追加されている
- [ ] 秘密鍵生成処理の共通ライブラリの説明が追加されている

- _Requirements: 3.6.1, 6.8_

---

#### タスク 6.2: Project-Structure.mdの更新
**目的**: プロジェクト構造ドキュメントに新規作成するファイルを追加する。

**作業内容**:
- `docs/Project-Structure.md`を修正
- `server/internal/usecase/cli/generate_secret_usecase.go`を追加
- `server/internal/service/secret_service.go`を追加
- `server/internal/auth/secret.go`を追加

**実装内容**:
- 修正対象: `docs/Project-Structure.md`
- 追加内容:
  - `server/internal/usecase/cli/generate_secret_usecase.go`の説明
  - `server/internal/service/secret_service.go`の説明
  - `server/internal/auth/secret.go`の説明

**受け入れ基準**:
- [ ] `docs/Project-Structure.md`に新規作成するファイルが追加されている
- [ ] `server/internal/usecase/cli/generate_secret_usecase.go`が追加されている
- [ ] `server/internal/service/secret_service.go`が追加されている
- [ ] `server/internal/auth/secret.go`が追加されている

- _Requirements: 3.6.2, 6.8_

---

#### タスク 6.3: Command-Line-Tool.mdの更新（存在する場合）
**目的**: CLIツールドキュメントのアーキテクチャ図を更新してusecase層を含める。

**作業内容**:
- `docs/Command-Line-Tool.md`が存在する場合は修正
- CLIコマンドのアーキテクチャ図を更新（usecase層を追加）
- レイヤー構造の説明を更新（main.go → usecase → service → auth → 出力）

**実装内容**:
- 修正対象: `docs/Command-Line-Tool.md`（存在する場合）
- 更新内容:
  - CLIコマンドのアーキテクチャ図を更新（usecase層を追加）
  - レイヤー構造の説明を更新

**受け入れ基準**:
- [ ] `docs/Command-Line-Tool.md`が存在する場合は、アーキテクチャ図が更新されている
- [ ] レイヤー構造の説明が更新されている（usecase層を含む）

- _Requirements: 3.6.3, 6.8_

---

#### タスク 6.4: structure.mdの更新
**目的**: ファイル組織ドキュメントに新規作成するファイルを追加する。

**作業内容**:
- `.kiro/steering/structure.md`を修正
- `server/internal/usecase/cli/generate_secret_usecase.go`を追加
- `server/internal/service/secret_service.go`を追加
- `server/internal/auth/secret.go`を追加

**実装内容**:
- 修正対象: `.kiro/steering/structure.md`
- 追加内容:
  - `server/internal/usecase/cli/generate_secret_usecase.go`の説明
  - `server/internal/service/secret_service.go`の説明
  - `server/internal/auth/secret.go`の説明

**受け入れ基準**:
- [ ] `.kiro/steering/structure.md`に新規作成するファイルが追加されている
- [ ] `server/internal/usecase/cli/generate_secret_usecase.go`が追加されている
- [ ] `server/internal/service/secret_service.go`が追加されている
- [ ] `server/internal/auth/secret.go`が追加されている

- _Requirements: 3.6.4, 6.8_

---

## タスクの依存関係

```
Phase 1: 秘密鍵生成処理の共通化（auth層）
  ├─ タスク 1.1: secret.goの実装
  └─ タスク 1.2: secret_test.goの実装

Phase 2: Service層の作成
  ├─ タスク 2.1: secret_service.goの実装（Phase 1に依存）
  └─ タスク 2.2: secret_service_test.goの実装（タスク 2.1に依存）

Phase 3: CLI用usecase層の作成
  ├─ タスク 3.1: generate_secret_usecase.goの実装（Phase 2に依存）
  └─ タスク 3.2: generate_secret_usecase_test.goの実装（タスク 3.1に依存）

Phase 4: main.goの簡素化
  └─ タスク 4.1: main.goの修正（Phase 3に依存）

Phase 5: ビルドと動作確認
  ├─ タスク 5.1: ビルドの確認（Phase 4に依存）
  ├─ タスク 5.2: 動作確認（タスク 5.1に依存）
  └─ タスク 5.3: 既存テストの確認（Phase 1-3に依存）

Phase 6: ドキュメントの更新
  ├─ タスク 6.1: Architecture.mdの更新（Phase 4に依存）
  ├─ タスク 6.2: Project-Structure.mdの更新（Phase 1-3に依存）
  ├─ タスク 6.3: Command-Line-Tool.mdの更新（Phase 4に依存）
  └─ タスク 6.4: structure.mdの更新（Phase 1-3に依存）
```

## 実装順序の推奨

1. **Phase 1**: 秘密鍵生成処理の共通化（auth層）
   - タスク 1.1 → タスク 1.2

2. **Phase 2**: Service層の作成
   - タスク 2.1 → タスク 2.2

3. **Phase 3**: CLI用usecase層の作成
   - タスク 3.1 → タスク 3.2

4. **Phase 4**: main.goの簡素化
   - タスク 4.1

5. **Phase 5**: ビルドと動作確認
   - タスク 5.1 → タスク 5.2 → タスク 5.3

6. **Phase 6**: ドキュメントの更新
   - タスク 6.1, 6.2, 6.3, 6.4（並行実行可能）

## 注意事項

- 各タスクの実装前に、対応する要件定義書と設計書の該当セクションを確認すること
- テストは各フェーズの実装と同時に実装すること
- ドキュメントの更新は実装完了後に実施すること
- ビルド出力先は`server/bin/generate-secret`に統一すること
