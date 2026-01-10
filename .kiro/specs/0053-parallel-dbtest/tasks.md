# 並列データベーステスト失敗対策の実装タスク一覧

## 概要
並列実行されるデータベーステストの失敗を解決するため、`gofrs/flock`ライブラリを使用したロックファイル機構を実装するためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: 依存関係の追加

#### タスク 1.1: gofrs/flockライブラリの追加
**目的**: ロックファイル管理に必要な`gofrs/flock`ライブラリを`go.mod`に追加する。

**作業内容**:
- `server/go.mod`に`gofrs/flock`ライブラリを追加
- `go mod tidy`を実行して依存関係を解決

**実装内容**:
- `cd server`でディレクトリに移動
- `go get github.com/gofrs/flock`を実行
- `go mod tidy`を実行

**受け入れ基準**:
- `server/go.mod`に`github.com/gofrs/flock`が追加されている
- `server/go.sum`が更新されている（自動生成）
- `go mod tidy`が正常に実行される

- _Requirements: 6.1, 7.1_
- _Design: 3.3_

---

### Phase 2: ロックファイル管理ユーティリティの実装

#### タスク 2.1: lock.goファイルの作成とgetProjectRoot()関数の実装
**目的**: プロジェクトルート（`go.mod`が存在するディレクトリ）を判定する関数を実装する。

**作業内容**:
- `server/test/testutil/lock.go`を作成
- `getProjectRoot()`関数を実装

**実装内容**:
- パッケージ名: `testutil`
- 関数名: `getProjectRoot() (string, error)`
- 実装方法:
  - `runtime.Caller(0)`で現在のファイルの位置を取得
  - そのディレクトリから`go.mod`を探す（最大5階層まで上に遡る）
  - `go.mod`が見つかったディレクトリを返す

**ファイル構成**:
```go
package testutil

import (
    "fmt"
    "os"
    "path/filepath"
    "runtime"
)

func getProjectRoot() (string, error) {
    _, filename, _, ok := runtime.Caller(0)
    if !ok {
        return "", fmt.Errorf("failed to get current file path")
    }
    
    dir := filepath.Dir(filename)
    
    // go.modを探す（最大5階層まで上に遡る）
    for i := 0; i < 5; i++ {
        goModPath := filepath.Join(dir, "go.mod")
        if _, err := os.Stat(goModPath); err == nil {
            return dir, nil
        }
        parent := filepath.Dir(dir)
        if parent == dir {
            break
        }
        dir = parent
    }
    
    return "", fmt.Errorf("go.mod not found")
}
```

**受け入れ基準**:
- `server/test/testutil/lock.go`が作成されている
- `getProjectRoot()`関数が実装されている
- `getProjectRoot()`関数が`go.mod`が存在するディレクトリを返す
- `go.mod`が見つからない場合、適切なエラーを返す

- _Requirements: 6.1_
- _Design: 3.1.2, 6.1_

---

#### タスク 2.2: AcquireTestLock()関数の実装
**目的**: ロックファイルを取得する関数を実装する。タイムアウト付きでロック取得を試行し、適切なエラーメッセージを返す。

**作業内容**:
- `AcquireTestLock()`関数を実装
- タイムアウト処理（30秒）を実装
- エラーハンドリングを実装

**実装内容**:
- 関数名: `AcquireTestLock(t *testing.T) (*flock.Flock, error)`
- 処理フロー:
  1. `getProjectRoot()`でプロジェクトルートを取得
  2. ロックファイルのパスを構築（`{projectRoot}/test-db.lock`）
  3. `flock.New(lockPath)`でロックオブジェクトを作成
  4. `context.WithTimeout(context.Background(), 30*time.Second)`でタイムアウト付きコンテキストを作成
  5. `TryLockContext(ctx)`でロック取得を試行
  6. エラーハンドリング:
     - タイムアウト: `"{ロックファイルPATH}のロックが取れなかったのでタイムアウトしました"`
     - その他のエラー: `"ロックファイルの取得に失敗しました ({ロックファイルPATH}): {エラー詳細}"`
  7. 成功したら`flock.Flock`オブジェクトを返す

**ファイル構成**:
```go
package testutil

import (
    "context"
    "fmt"
    "path/filepath"
    "testing"
    "time"

    "github.com/gofrs/flock"
)

// AcquireTestLock acquires a file lock for database tests
// Returns a flock.Flock object that should be unlocked with defer fileLock.Unlock()
func AcquireTestLock(t *testing.T) (*flock.Flock, error) {
    // プロジェクトルートを取得
    projectRoot, err := getProjectRoot()
    if err != nil {
        return nil, fmt.Errorf("failed to get project root: %w", err)
    }
    
    // ロックファイルのパスを構築
    lockPath := filepath.Join(projectRoot, "test-db.lock")
    
    // ロックオブジェクトを作成
    fileLock := flock.New(lockPath)
    
    // タイムアウト付きコンテキストを作成
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // ロック取得を試行
    err = fileLock.TryLockContext(ctx)
    if err != nil {
        // エラーハンドリング
        if err == context.DeadlineExceeded {
            return nil, fmt.Errorf("%sのロックが取れなかったのでタイムアウトしました", lockPath)
        }
        return nil, fmt.Errorf("ロックファイルの取得に失敗しました (%s): %w", lockPath, err)
    }
    
    return fileLock, nil
}
```

**受け入れ基準**:
- `AcquireTestLock()`関数が実装されている
- タイムアウト処理（30秒）が実装されている
- エラーメッセージにロックファイルのパスが含まれている
- タイムアウト時のエラーメッセージが正しい形式である
- その他のエラー時のエラーメッセージが正しい形式である

- _Requirements: 3.2.1, 6.1, 6.4_
- _Design: 3.1.3, 4.1, 4.2_

---

### Phase 3: SetupTestGroupManager()の修正

#### タスク 3.1: SetupTestGroupManager()にロック取得・解放処理を追加
**目的**: `SetupTestGroupManager()`関数の先頭にロック取得処理を追加し、`defer`でロック解放処理を追加する。

**作業内容**:
- `server/test/testutil/db.go`の`SetupTestGroupManager()`関数を修正
- 関数の先頭に`AcquireTestLock()`の呼び出しを追加
- `defer`でロック解放処理を追加
- エラーハンドリングを実装

**実装内容**:
- 修正対象: `server/test/testutil/db.go`の`SetupTestGroupManager()`関数
- 修正箇所: 関数の先頭（設定ファイル読み込みの前）
- 追加内容:
  - `AcquireTestLock(t)`の呼び出し
  - `defer fileLock.Unlock()`でロック解放
  - ロック取得失敗時のエラーハンドリング（`t.Fatalf()`）
  - ロック解放失敗時の警告ログ（`t.Logf()`）

**修正後のコードイメージ**:
```go
func SetupTestGroupManager(t *testing.T, dbCount int, tablesPerDB int) *db.GroupManager {
    // ロックを取得
    fileLock, err := AcquireTestLock(t)
    if err != nil {
        t.Fatalf("Failed to acquire test lock: %v", err)
    }
    defer func() {
        if err := fileLock.Unlock(); err != nil {
            t.Logf("Warning: failed to unlock test lock: %v", err)
        }
    }()
    
    // 設定ファイルから読み込む
    cfg, err := LoadTestConfig()
    require.NoError(t, err)
    
    // ... 既存の処理を続行
}
```

**受け入れ基準**:
- `SetupTestGroupManager()`関数にロック取得処理が追加されている
- `SetupTestGroupManager()`関数にロック解放処理が追加されている（`defer`を使用）
- ロック取得失敗時に`t.Fatalf()`でテストを失敗させる
- ロック解放失敗時に`t.Logf()`で警告をログに出力する
- 既存の処理が正常に動作する

- _Requirements: 3.3.1, 6.2_
- _Design: 3.2_

---

### Phase 4: 動作確認

#### タスク 4.1: 単体テストの実行
**目的**: ロックファイル管理ユーティリティが正常に動作することを確認する。

**作業内容**:
- `server/test/testutil/lock_test.go`を作成（オプション）
- 単体テストを実行して動作確認

**受け入れ基準**:
- ロックファイル管理ユーティリティが正常に動作する
- `getProjectRoot()`関数が正しくプロジェクトルートを返す
- `AcquireTestLock()`関数が正常にロックを取得できる

- _Requirements: 6.1_
- _Design: 7.1_

---

#### タスク 4.2: 並列実行テストの実行
**目的**: 並列実行時にテストが正常に動作することを確認する。

**作業内容**:
- `go test -parallel 4 ./test/integration/... ./test/e2e/..`を実行
- テストが正常に完了することを確認
- ロックファイルが適切に作成・削除されることを確認

**受け入れ基準**:
- `go test -parallel 4 ./test/integration/... ./test/e2e/..`が正常に実行できる
- 並列実行時にテストが失敗しない
- 既存のテストが全て失敗しないことを確認
- ロックファイルが適切に作成・削除されることを確認（`git status`で確認）

- _Requirements: 6.3, 6.5_
- _Design: 7.2_

---

#### タスク 4.3: エラーハンドリングの確認
**目的**: エラーハンドリングが適切に実装されていることを確認する。

**作業内容**:
- ロック取得失敗時のエラーメッセージを確認
- タイムアウト時のエラーメッセージを確認
- エラーメッセージにロックファイルのパスが含まれていることを確認

**受け入れ基準**:
- ロック取得失敗時に適切なエラーメッセージが表示される
- タイムアウト時に適切なエラーメッセージが表示される（形式: `"{ロックファイルPATH}のロックが取れなかったのでタイムアウトしました"`）
- エラーメッセージにロックファイルのパスが含まれている

- _Requirements: 6.4_
- _Design: 4.1, 4.2_

---

### Phase 5: .gitignoreの確認

#### タスク 5.1: .gitignoreの確認
**目的**: ロックファイルが`.gitignore`に追加されていないことを確認する。

**作業内容**:
- `.gitignore`ファイルを確認
- ロックファイル（`test-db.lock`）が`.gitignore`に追加されていないことを確認
- `git status`でロックファイルが表示されることを確認（ロックファイルが残っている場合）

**受け入れ基準**:
- `.gitignore`にロックファイル（`test-db.lock`）が追加されていない
- `git status`でロックファイルが表示される（ロックファイルが残っている場合）

- _Requirements: 6.4, 7.1_
- _Design: 8.1_

---

### Phase 6: ドキュメントの更新

#### タスク 6.1: Testing.mdにロックファイル機構の記載を追加
**目的**: テストドキュメントにロックファイル機構の説明を追加し、開発者が理解できるようにする。

**作業内容**:
- `docs/Testing.md`の「Test Utilities」セクションにロックファイル機構の説明を追加
- 並列実行時の動作について説明
- ロックファイルの場所と管理方法について説明
- エラーメッセージの確認方法について説明

**実装内容**:
- 追加場所: `docs/Testing.md`の「Test Utilities」セクション（`SetupTestGroupManager()`の説明の後）
- 追加内容:
  - ロックファイル機構の概要
  - 並列実行時の排他制御の仕組み
  - ロックファイルの場所（`server/test-db.lock`）
  - タイムアウト（30秒）
  - エラーメッセージの形式
  - ロックファイルの消し忘れに気づく方法（`git status`）

**追加する内容のイメージ**:
```markdown
### Lock File Mechanism for Parallel Tests

When running tests in parallel with `go test -parallel 4`, database tests use a file lock mechanism to prevent conflicts. The `SetupTestGroupManager()` function automatically acquires a lock before setting up the test database.

**Lock File Location**: `server/test-db.lock`

**How It Works**:
1. Before setting up the test database, `SetupTestGroupManager()` acquires a file lock
2. If another test is already running, it waits up to 30 seconds for the lock to be released
3. After the test completes, the lock is automatically released

**Error Messages**:
- Timeout: `"{lock_file_path}のロックが取れなかったのでタイムアウトしました"`
- Other errors: `"ロックファイルの取得に失敗しました ({lock_file_path}): {error_details}"`

**Checking for Leftover Lock Files**:
If a test is interrupted, the lock file might remain. Check with:
```bash
git status
```

If `test-db.lock` appears, manually delete it:
```bash
rm server/test-db.lock
```

**Note**: The lock file is intentionally not added to `.gitignore` so that leftover files are visible in `git status`.
```

**受け入れ基準**:
- `docs/Testing.md`にロックファイル機構の説明が追加されている
- 並列実行時の動作が説明されている
- ロックファイルの場所と管理方法が説明されている
- エラーメッセージの確認方法が説明されている
- ロックファイルの消し忘れに気づく方法が説明されている

- _Requirements: 7.1_
- _Design: 8.1_

---

## 実装順序の推奨

1. **Phase 1**: 依存関係の追加
2. **Phase 2**: ロックファイル管理ユーティリティの実装
3. **Phase 3**: `SetupTestGroupManager()`の修正
4. **Phase 4**: 動作確認
5. **Phase 5**: `.gitignore`の確認
6. **Phase 6**: ドキュメントの更新

## 注意事項

### 実装時の注意点

1. **ロックファイルのパス**: プロジェクトルート直下の`test-db.lock`を使用
2. **タイムアウト**: 30秒のタイムアウトを設定
3. **エラーメッセージ**: ロックファイルのパスを含める（消し忘れに気づけるように）
4. **ロック解放**: `defer`を使用して確実にロックを解放
5. **既存テストへの影響**: 既存のテストコードを変更する必要はない（`SetupTestGroupManager()`内で処理）

### テスト時の注意点

1. **並列実行**: `go test -parallel 4`で並列実行を確認
2. **ロックファイルの確認**: テスト実行後に`git status`でロックファイルが残っていないことを確認
3. **エラーメッセージの確認**: タイムアウト時やエラー時に適切なメッセージが表示されることを確認

## 参考情報

- 要件定義書: `requirements.md`
- 設計書: `design.md`
- 既存のテストユーティリティ: `server/test/testutil/db.go`
- 使用ライブラリ: https://github.com/gofrs/flock
