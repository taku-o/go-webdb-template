# 並列データベーステスト失敗対策の設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、並列実行されるデータベーステストの失敗を解決するため、`gofrs/flock`ライブラリを使用したロックファイル機構を実装する詳細設計を定義する。

### 1.2 設計の範囲
- ロックファイル管理ユーティリティ関数の設計
- `SetupTestGroupManager()`関数へのロック取得・解放処理の追加設計
- プロジェクトルートの判定方法の設計
- エラーハンドリングの詳細設計
- タイムアウト処理の設計
- テスト設計

### 1.3 設計方針
- **最小限の変更**: 既存のテストコードを変更せず、`SetupTestGroupManager()`内で処理
- **確実なロック解放**: `defer`を使用して確実にロックを解放
- **明確なエラーメッセージ**: ロックファイルのパスを含むエラーメッセージで、消し忘れに気づけるようにする
- **シンプルな実装**: `gofrs/flock`ライブラリを活用し、シンプルに実装

## 2. アーキテクチャ設計

### 2.1 全体構成

```
テスト実行
  ↓
SetupTestGroupManager()
  ↓
AcquireTestLock()  ← 新規作成
  ↓
flock.New() + TryLockContext()  ← gofrs/flock使用
  ↓
ロック取得成功
  ↓
データベースセットアップ処理
  ↓
defer fileLock.Unlock()  ← ロック解放
```

### 2.2 ディレクトリ構造

```
server/
├── go.mod                    # gofrs/flockを追加
├── go.sum                    # 依存関係（自動生成）
├── test-db.lock             # ロックファイル（git管理対象外）
└── test/
    └── testutil/
        ├── db.go            # SetupTestGroupManager()を修正
        └── lock.go          # 新規作成: AcquireTestLock()
```

## 3. 詳細設計

### 3.1 ロックファイル管理ユーティリティ (`server/test/testutil/lock.go`)

#### 3.1.1 ファイル構成

```go
package testutil

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
    "runtime"
    "testing"
    "time"

    "github.com/gofrs/flock"
)
```

#### 3.1.2 プロジェクトルートの判定

プロジェクトルート（`server/go.mod`が存在するディレクトリ）を判定する関数を実装する。

**実装方法**:
- 現在のファイルの位置（`runtime.Caller()`）から`go.mod`を探す
- または、`go.mod`が存在するディレクトリを再帰的に探す
- または、環境変数や固定パスを使用（シンプルな方法）

**推奨実装**:
```go
// getProjectRoot returns the project root directory (where go.mod exists)
func getProjectRoot() (string, error) {
    // 現在のファイルのディレクトリから開始
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

#### 3.1.3 `AcquireTestLock()`関数の設計

**関数シグネチャ**:
```go
func AcquireTestLock(t *testing.T) (*flock.Flock, error)
```

**処理フロー**:
1. プロジェクトルートを取得
2. ロックファイルのパスを構築（`{projectRoot}/test-db.lock`）
3. `flock.New(lockPath)`でロックオブジェクトを作成
4. `context.WithTimeout(context.Background(), 30*time.Second)`でタイムアウト付きコンテキストを作成
5. `TryLockContext(ctx)`でロック取得を試行
6. エラーハンドリング:
   - タイムアウト: `"{ロックファイルPATH}のロックが取れなかったのでタイムアウトしました"`
   - その他のエラー: `"ロックファイルの取得に失敗しました ({ロックファイルPATH}): {エラー詳細}"`
7. 成功したら`flock.Flock`オブジェクトを返す

**実装コード**:
```go
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

### 3.2 `SetupTestGroupManager()`関数の修正設計

#### 3.2.1 修正箇所

`server/test/testutil/db.go`の`SetupTestGroupManager()`関数の先頭にロック取得処理を追加する。

#### 3.2.2 修正内容

**修正前**:
```go
func SetupTestGroupManager(t *testing.T, dbCount int, tablesPerDB int) *db.GroupManager {
    // 設定ファイルから読み込む
    cfg, err := LoadTestConfig()
    require.NoError(t, err)
    
    // ... 既存の処理
}
```

**修正後**:
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
    
    // ... 既存の処理
}
```

#### 3.2.3 エラーハンドリング

- **ロック取得失敗**: `t.Fatalf()`でテストを失敗させる
- **ロック解放失敗**: `t.Logf()`で警告をログに出力（テストは継続）

### 3.3 依存関係の追加

#### 3.3.1 `go.mod`への追加

```bash
cd server
go get github.com/gofrs/flock
go mod tidy
```

#### 3.3.2 バージョン指定

最新の安定版を使用（バージョン指定は不要、`go get`で最新版を取得）。

## 4. エラーハンドリング設計

### 4.1 エラーケースと対応

| エラーケース | エラーメッセージ | 対応 |
|------------|----------------|------|
| プロジェクトルート取得失敗 | `"failed to get project root: {エラー詳細}"` | テストを失敗させる |
| ロック取得タイムアウト | `"{ロックファイルPATH}のロックが取れなかったのでタイムアウトしました"` | テストを失敗させる |
| ロック取得失敗（権限不足など） | `"ロックファイルの取得に失敗しました ({ロックファイルPATH}): {エラー詳細}"` | テストを失敗させる |
| ロック解放失敗 | `"Warning: failed to unlock test lock: {エラー詳細}"` | 警告をログに出力（テストは継続） |

### 4.2 エラーメッセージの形式

- **タイムアウト**: `"{ロックファイルPATH}のロックが取れなかったのでタイムアウトしました"`
  - 例: `/path/to/server/test-db.lockのロックが取れなかったのでタイムアウトしました`
- **その他のエラー**: `"ロックファイルの取得に失敗しました ({ロックファイルPATH}): {エラー詳細}"`
  - 例: `ロックファイルの取得に失敗しました (/path/to/server/test-db.lock): permission denied`

## 5. タイムアウト設計

### 5.1 タイムアウト時間

- **デフォルト**: 30秒
- **理由**: テストの実行時間を考慮し、適切な待機時間を設定

### 5.2 タイムアウト処理

- `context.WithTimeout(context.Background(), 30*time.Second)`を使用
- `TryLockContext(ctx)`でタイムアウトを検知
- `context.DeadlineExceeded`エラーを検知して、適切なエラーメッセージを返す

## 6. プロジェクトルート判定の詳細設計

### 6.1 判定方法

1. 現在のファイル（`lock.go`）の位置を取得（`runtime.Caller(0)`）
2. そのディレクトリから`go.mod`を探す
3. 見つからない場合は親ディレクトリに遡る（最大5階層まで）
4. `go.mod`が見つかったディレクトリをプロジェクトルートとする

### 6.2 実装の詳細

```go
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

### 6.3 エッジケースの考慮

- **`go.mod`が見つからない場合**: エラーを返す
- **シンボリックリンク**: `filepath.Dir()`と`os.Stat()`で適切に処理される
- **権限不足**: `os.Stat()`でエラーが返される場合は、エラーを返す

## 7. テスト設計

### 7.1 単体テスト

#### 7.1.1 `AcquireTestLock()`のテスト

**テストケース**:
1. 正常系: ロック取得に成功する
2. 異常系: プロジェクトルートが見つからない場合
3. 異常系: ロック取得がタイムアウトする場合（別プロセスでロックを保持）

**テストファイル**: `server/test/testutil/lock_test.go`（新規作成）

**実装イメージ**:
```go
func TestAcquireTestLock_Success(t *testing.T) {
    fileLock, err := AcquireTestLock(t)
    require.NoError(t, err)
    defer fileLock.Unlock()
    
    // ロックが取得できていることを確認
    assert.NotNil(t, fileLock)
}

func TestAcquireTestLock_Timeout(t *testing.T) {
    // 別のロックを先に取得
    fileLock1, err := AcquireTestLock(t)
    require.NoError(t, err)
    defer fileLock1.Unlock()
    
    // 別のテストでロック取得を試みる（タイムアウトする）
    // 注意: 並列実行では同じプロセス内なので、実際にはタイムアウトしない可能性がある
    // このテストは統合テストで確認する
}
```

### 7.2 統合テスト

#### 7.2.1 並列実行テスト

**テストケース**:
- `go test -parallel 4 ./test/integration/... ./test/e2e/..`が正常に実行できる
- 並列実行時にテストが失敗しない
- ロックファイルが適切に作成・削除される

**確認方法**:
- 実際に並列実行してテストが成功することを確認
- ロックファイルが残っていないことを確認（`git status`で確認）

## 8. 実装上の注意事項

### 8.1 ロックファイルの管理

- **作成場所**: プロジェクトルート直下（`server/test-db.lock`）
- **Git管理**: `.gitignore`に追加しない（消し忘れに気づけるようにする）
- **削除**: テスト終了時に`defer fileLock.Unlock()`で自動削除

### 8.2 並列実行時の動作

- **同一プロセス内**: `gofrs/flock`が適切に排他制御を行う
- **異なるプロセス**: ファイルロックにより排他制御される
- **タイムアウト**: 30秒以内にロック取得できない場合はエラー

### 8.3 エラーメッセージの重要性

- **ロックファイルのパスを含める**: 消し忘れに気づけるようにする
- **明確なメッセージ**: タイムアウトとその他のエラーを区別する

### 8.4 既存テストへの影響

- **自動適用**: `SetupTestGroupManager()`を使用しているすべてのテストに自動的に適用される
- **コード変更不要**: 既存のテストコードを変更する必要はない

## 9. 実装順序

### 9.1 実装ステップ

1. **依存関係の追加**
   - `server/go.mod`に`gofrs/flock`を追加
   - `go mod tidy`を実行

2. **ロックファイル管理ユーティリティの実装**
   - `server/test/testutil/lock.go`を作成
   - `getProjectRoot()`関数を実装
   - `AcquireTestLock()`関数を実装

3. **`SetupTestGroupManager()`の修正**
   - `server/test/testutil/db.go`を修正
   - ロック取得・解放処理を追加

4. **動作確認**
   - 単体テストを実行
   - 並列実行テストを実行
   - ロックファイルが適切に作成・削除されることを確認

5. **`.gitignore`の確認**
   - ロックファイル（`test-db.lock`）が`.gitignore`に追加されていないことを確認

## 10. 参考情報

### 10.1 使用ライブラリ

- **gofrs/flock**: https://github.com/gofrs/flock
- ファイルロックライブラリ
- クロスプラットフォーム対応（Linux, macOS, Windows）

### 10.2 関連ドキュメント

- 要件定義書: `requirements.md`
- 既存のテストユーティリティ: `server/test/testutil/db.go`

### 10.3 技術スタック

- **言語**: Go
- **テストフレームワーク**: Go標準の`testing`パッケージ
- **ロックライブラリ**: `gofrs/flock`
- **ファイルシステム**: `os`, `path/filepath`, `runtime`パッケージ
