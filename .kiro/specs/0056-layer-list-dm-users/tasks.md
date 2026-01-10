# server/cmd/list-dm-usersのレイヤー構造修正の実装タスク一覧

## 概要
`server/cmd/list-dm-users`の実装を、APIサーバーと同じレイヤー構造（usecase -> service -> repository -> model）に変更するためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: CLI用usecase層の作成

#### タスク 1.1: usecase/cliディレクトリの作成
**目的**: CLI用のusecase層を配置するためのディレクトリを作成する。

**作業内容**:
- `server/internal/usecase/cli`ディレクトリを作成

**実装内容**:
- ディレクトリパス: `server/internal/usecase/cli/`
- 既存の`server/internal/usecase/`ディレクトリの下に`cli`サブディレクトリを作成

**受け入れ基準**:
- [ ] `server/internal/usecase/cli`ディレクトリが存在する

- _Requirements: 3.1.1, 6.1_
- _Design: 3.1.1_

---

#### タスク 1.2: ListDmUsersUsecaseの実装
**目的**: CLI用のdm_user一覧取得usecaseを実装する。

**作業内容**:
- `server/internal/usecase/cli/list_dm_users_usecase.go`を作成
- `ListDmUsersUsecase`構造体を定義
- `DmUserServiceInterface`を依存として注入
- `ListDmUsers()`メソッドを実装

**実装内容**:
- ファイルパス: `server/internal/usecase/cli/list_dm_users_usecase.go`
- パッケージ名: `cli`
- 構造体定義:
  ```go
  type ListDmUsersUsecase struct {
      dmUserService usecase.DmUserServiceInterface
  }
  ```
- コンストラクタ:
  ```go
  func NewListDmUsersUsecase(dmUserService usecase.DmUserServiceInterface) *ListDmUsersUsecase
  ```
- メソッド:
  ```go
  func (u *ListDmUsersUsecase) ListDmUsers(ctx context.Context, limit, offset int) ([]*model.DmUser, error)
  ```
- 既存の`DmUserServiceInterface`を使用（`internal/usecase/dm_user_usecase.go`で定義済み）
- service層の`ListDmUsers()`を呼び出して結果を返す
- エラーハンドリング: service層から返されたエラーをそのまま返す

**受け入れ基準**:
- [ ] `server/internal/usecase/cli/list_dm_users_usecase.go`が作成されている
- [ ] `ListDmUsersUsecase`構造体が定義されている
- [ ] `DmUserServiceInterface`を依存として注入している
- [ ] `ListDmUsers(ctx context.Context, limit, offset int) ([]*model.DmUser, error)`メソッドが実装されている
- [ ] service層の`ListDmUsers()`を呼び出して結果を返している
- [ ] エラーハンドリングが適切に実装されている

- _Requirements: 3.1.2, 6.1_
- _Design: 3.1.2_

---

### Phase 2: main.goの簡素化

#### タスク 2.1: usecase層のインポートと初期化の追加
**目的**: main.goでusecase層を使用できるようにインポートと初期化を追加する。

**作業内容**:
- `server/cmd/list-dm-users/main.go`を修正
- `internal/usecase/cli`パッケージをインポート
- usecase層の初期化を追加
- service層の直接呼び出しをusecase層の呼び出しに変更

**実装内容**:
- 修正対象: `server/cmd/list-dm-users/main.go`
- インポート追加:
  ```go
  "github.com/taku-o/go-webdb-template/internal/usecase/cli"
  ```
- 初期化の追加（Service層の初期化の後）:
  ```go
  // Usecase層の初期化
  listDmUsersUsecase := cli.NewListDmUsersUsecase(dmUserService)
  ```
- service層の直接呼び出しを変更:
  - 修正前:
    ```go
    dmUsers, err := dmUserService.ListDmUsers(ctx, validatedLimit, 0)
    ```
  - 修正後:
    ```go
    dmUsers, err := listDmUsersUsecase.ListDmUsers(ctx, validatedLimit, 0)
    ```

**受け入れ基準**:
- [ ] `internal/usecase/cli`パッケージがインポートされている
- [ ] usecase層が適切に初期化されている
- [ ] service層の直接呼び出しがusecase層の呼び出しに変更されている
- [ ] 既存のコードの動作が維持されている

- _Requirements: 3.2.1, 3.3.1, 6.2, 6.3_
- _Design: 3.2.1_

---

#### タスク 2.2: 既存関数の維持確認
**目的**: 既存のバリデーション関数と出力関数が維持されていることを確認する。

**作業内容**:
- `validateLimit()`関数が維持されていることを確認
- `printDmUsersTSV()`関数が維持されていることを確認
- 既存のバリデーションロジックが維持されていることを確認
- 既存の出力形式（TSV形式）が維持されていることを確認

**実装内容**:
- 確認対象: `server/cmd/list-dm-users/main.go`
- `validateLimit()`関数: 変更なし（維持）
- `printDmUsersTSV()`関数: 変更なし（維持）

**受け入れ基準**:
- [ ] `validateLimit()`関数が維持されている
- [ ] `printDmUsersTSV()`関数が維持されている
- [ ] 既存のバリデーションロジックが維持されている
- [ ] 既存の出力形式（TSV形式）が維持されている

- _Requirements: 3.2.2, 3.2.3, 6.2_
- _Design: 3.2.2, 3.2.3_

---

### Phase 3: テストの実装

#### タスク 3.1: usecase層のテストファイルの作成
**目的**: CLI用usecase層の単体テストを実装する。

**作業内容**:
- `server/internal/usecase/cli/list_dm_users_usecase_test.go`を作成
- モックの実装（`MockDmUserServiceInterface`）
- テストケースの実装（正常系、異常系）

**実装内容**:
- ファイルパス: `server/internal/usecase/cli/list_dm_users_usecase_test.go`
- パッケージ名: `cli`
- モック実装:
  ```go
  type MockDmUserServiceInterface struct {
      ListDmUsersFunc func(ctx context.Context, limit, offset int) ([]*model.DmUser, error)
  }
  ```
- テストケース:
  1. 正常系: ユーザー一覧が取得できる場合
  2. 正常系: 空のリストが返される場合
  3. 異常系: service層でエラーが発生した場合
- テーブル駆動テストを使用
- `github.com/stretchr/testify/assert`を使用

**受け入れ基準**:
- [ ] `server/internal/usecase/cli/list_dm_users_usecase_test.go`が作成されている
- [ ] モック（`MockDmUserServiceInterface`）が実装されている
- [ ] 正常系のテストケースが実装されている
- [ ] 異常系のテストケースが実装されている
- [ ] 全てのテストケースが通過する

- _Requirements: 6.5_
- _Design: 3.4.1_

---

#### タスク 3.2: 既存テストの確認
**目的**: 既存のテスト（`main_test.go`）が全て通過することを確認する。

**作業内容**:
- `server/cmd/list-dm-users/main_test.go`の既存テストを実行
- 全てのテストケースが通過することを確認

**実装内容**:
- テスト実行:
  ```bash
  go test -v ./cmd/list-dm-users/...
  ```
- 確認対象:
  - `TestPrintDmUsersTSV`: TSV形式での出力のテスト
  - `TestValidateLimit`: バリデーション関数のテスト

**受け入れ基準**:
- [ ] 既存のテスト（`server/cmd/list-dm-users/main_test.go`）が全て通過する
- [ ] `TestPrintDmUsersTSV`が通過する
- [ ] `TestValidateLimit`が通過する

- _Requirements: 6.4, 6.5_
- _Design: 3.4.2_

---

### Phase 4: ドキュメントの更新

#### タスク 4.1: Architecture.mdの更新
**目的**: CLIコマンドのレイヤー構造をアーキテクチャドキュメントに反映する。

**作業内容**:
- `docs/Architecture.md`を修正
- CLIコマンドのレイヤー構造を追加
- CLIコマンドのアーキテクチャ図を更新
- CLI用usecase層の説明を追加

**実装内容**:
- 修正対象: `docs/Architecture.md`
- 追加するセクション:
  - CLI Layer (`cmd/list-dm-users`)の説明
  - CLI Usecase Layer (`internal/usecase/cli`)の説明
- アーキテクチャ図の更新（usecase層を含む）

**受け入れ基準**:
- [ ] `docs/Architecture.md`にCLIコマンドのレイヤー構造が追加されている
- [ ] `docs/Architecture.md`のCLIコマンドのアーキテクチャ図が更新されている
- [ ] CLI用usecase層の説明が追加されている

- _Requirements: 3.4.1, 6.6_
- _Design: 3.5.1_

---

#### タスク 4.2: Project-Structure.mdの更新
**目的**: 新規作成する`server/internal/usecase/cli`ディレクトリをプロジェクト構造に反映する。

**作業内容**:
- `docs/Project-Structure.md`を修正
- `server/internal/usecase/cli`ディレクトリを追加
- `server/internal/usecase/cli/list_dm_users_usecase.go`を追加

**実装内容**:
- 修正対象: `docs/Project-Structure.md`
- 追加する行:
  ```markdown
  │   │   ├── usecase/            # ビジネスロジック層
  │   │   │   ├── dm_user_usecase.go
  │   │   │   ├── dm_post_usecase.go
  │   │   │   ├── email_usecase.go
  │   │   │   ├── cli/            # CLI用usecase層
  │   │   │   │   └── list_dm_users_usecase.go
  │   │   │   └── ...
  ```

**受け入れ基準**:
- [ ] `docs/Project-Structure.md`に`server/internal/usecase/cli`ディレクトリが追加されている
- [ ] `server/internal/usecase/cli/list_dm_users_usecase.go`が追加されている

- _Requirements: 3.4.2, 6.6_
- _Design: 3.5.2_

---

#### タスク 4.3: Command-Line-Tool.mdの更新
**目的**: CLIコマンドのアーキテクチャ図を更新してusecase層を含める。

**作業内容**:
- `docs/Command-Line-Tool.md`を修正
- CLIコマンドのアーキテクチャ図を更新（usecase層を追加）
- レイヤー構造の説明を更新

**実装内容**:
- 修正対象: `docs/Command-Line-Tool.md`
- 更新するアーキテクチャ図:
  - usecase層を追加
  - レイヤー構造の説明を更新（main.go → usecase → service → repository → model）

**受け入れ基準**:
- [ ] `docs/Command-Line-Tool.md`のアーキテクチャ図が更新されている（usecase層を含む）
- [ ] `docs/Command-Line-Tool.md`のレイヤー構造の説明が更新されている

- _Requirements: 3.4.3, 6.6_
- _Design: 3.5.3_

---

#### タスク 4.4: structure.mdの更新
**目的**: 新規作成する`server/internal/usecase/cli`ディレクトリをファイル組織に反映する。

**作業内容**:
- `.kiro/steering/structure.md`を修正
- `server/internal/usecase/cli`ディレクトリを追加
- CLI用usecase層の説明を追加

**実装内容**:
- 修正対象: `.kiro/steering/structure.md`
- 追加する行:
  ```markdown
  │   │   ├── usecase/            # ビジネスロジック層
  │   │   │   ├── dm_user_usecase.go
  │   │   │   ├── dm_post_usecase.go
  │   │   │   ├── email_usecase.go
  │   │   │   ├── cli/            # CLI用usecase層
  │   │   │   │   └── list_dm_users_usecase.go
  │   │   │   └── ...
  ```

**受け入れ基準**:
- [ ] `.kiro/steering/structure.md`に`server/internal/usecase/cli`ディレクトリが追加されている
- [ ] CLI用usecase層の説明が追加されている

- _Requirements: 3.4.4, 6.6_
- _Design: 3.5.4_

---

### Phase 5: 動作確認

#### タスク 5.1: ローカル環境での動作確認
**目的**: ローカル環境でCLIコマンドが正常に動作することを確認する。

**作業内容**:
- ローカル環境でCLIコマンドを実行
- 既存の引数（`-limit`）が正常に動作することを確認
- 既存の出力形式（TSV形式）が維持されていることを確認
- 既存のエラーメッセージが維持されていることを確認

**実装内容**:
- テスト実行:
  ```bash
  APP_ENV=develop go run ./cmd/list-dm-users/main.go --limit 20
  ```
- 確認項目:
  - コマンドが正常に実行される
  - 引数（`-limit`）が正常に動作する
  - 出力形式（TSV形式）が維持されている
  - エラーメッセージが維持されている

**受け入れ基準**:
- [ ] ローカル環境でCLIコマンドが正常に動作する
- [ ] 既存の引数（`-limit`）が正常に動作する
- [ ] 既存の出力形式（TSV形式）が維持されている
- [ ] 既存のエラーメッセージが維持されている

- _Requirements: 6.4_
- _Design: 9.1_

---

#### タスク 5.2: テストの実行と確認
**目的**: 全てのテストが通過することを確認する。

**作業内容**:
- usecase層のテストを実行
- main.goのテストを実行
- 全てのテストが通過することを確認

**実装内容**:
- テスト実行:
  ```bash
  # usecase層のテスト
  go test -v ./internal/usecase/cli/...
  
  # main.goのテスト
  go test -v ./cmd/list-dm-users/...
  
  # カバレッジ確認
  go test -cover ./internal/usecase/cli/...
  ```

**受け入れ基準**:
- [ ] usecase層のテストが全て通過する
- [ ] main.goのテストが全て通過する
- [ ] テストカバレッジが適切である（usecase層: 80%以上を目標）

- _Requirements: 6.4, 6.5_
- _Design: 8.2_

---

#### タスク 5.3: CI環境での動作確認（該当する場合）
**目的**: CI環境でもCLIコマンドが正常に動作することを確認する。

**作業内容**:
- CI環境でCLIコマンドを実行
- 全てのテストが通過することを確認

**実装内容**:
- CI環境でのテスト実行
- テスト結果の確認

**受け入れ基準**:
- [ ] CI環境でCLIコマンドが正常に動作する（該当する場合）
- [ ] CI環境で全てのテストが通過する（該当する場合）

- _Requirements: 6.4_
- _Design: 9.1_

---

## 実装順序

1. **Phase 1: CLI用usecase層の作成**（最優先）
   - タスク 1.1: usecase/cliディレクトリの作成
   - タスク 1.2: ListDmUsersUsecaseの実装

2. **Phase 2: main.goの簡素化**
   - タスク 2.1: usecase層のインポートと初期化の追加
   - タスク 2.2: 既存関数の維持確認

3. **Phase 3: テストの実装**
   - タスク 3.1: usecase層のテストファイルの作成
   - タスク 3.2: 既存テストの確認

4. **Phase 4: ドキュメントの更新**
   - タスク 4.1: Architecture.mdの更新
   - タスク 4.2: Project-Structure.mdの更新
   - タスク 4.3: Command-Line-Tool.mdの更新
   - タスク 4.4: structure.mdの更新

5. **Phase 5: 動作確認**
   - タスク 5.1: ローカル環境での動作確認
   - タスク 5.2: テストの実行と確認
   - タスク 5.3: CI環境での動作確認（該当する場合）

## 注意事項

### 実装時の注意点

1. **既存コードの維持**: 既存のservice、repository、model層は変更しない
2. **テストの維持**: 既存のテスト（main_test.go）は全て通過することを確認
3. **後方互換性**: 既存のCLIコマンドの動作（引数、出力形式、エラーメッセージ）を維持
4. **一貫性**: 既存のusecase層の実装パターンに従う

### リスクと対策

1. **既存テストの失敗**: main.goを修正したことで既存のテストが失敗する可能性
   - **対策**: 既存のテスト関数（`validateLimit()`、`printDmUsersTSV()`）は変更しないため、既存のテストはそのまま動作する

2. **パフォーマンスの劣化**: usecase層の追加によるパフォーマンスの劣化
   - **対策**: usecase層はservice層のメソッドをそのまま呼び出すだけなので、パフォーマンスへの影響は無視できるレベル

## 関連ドキュメント

- 要件定義書: `requirements.md`
- 設計書: `design.md`
- 関連Issue: https://github.com/taku-o/go-webdb-template/issues/116
