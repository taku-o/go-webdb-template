# CSVダウンロード機能設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、ユーザー情報をCSV形式でダウンロードする機能の詳細設計を定義する。ストリーミング方式によるメモリ効率の良い実装を実現し、クライアント側とAPIサーバー側の両方で適切に動作するシステムを構築する。

### 1.2 設計の範囲
- クライアント側: CSVダウンロードボタンの実装、APIクライアントの拡張、ダウンロード処理の実装
- APIサーバー側: CSVダウンロードエンドポイントの実装、ストリーミングCSV生成の実装、タイムアウト設定
- エラーハンドリング: クライアント側とサーバー側のエラーハンドリング
- テスト: 単体テストとE2Eテストの実装

### 1.3 設計方針
- **ストリーミング方式の採用**: ファイルを作成せず、HTTPレスポンスに直接ストリーミングすることでメモリ効率を向上
- **既存システムとの統合**: 既存のHuma API、Echo、Next.jsのパターンに従う
- **タイムアウト設定**: ストリーミングレスポンスの書き込みタイムアウトを3分に設定
- **エラーハンドリング**: 適切なエラーメッセージとHTTPステータスコードを返す
- **既存機能の維持**: 既存のdm-users API機能に影響を与えない

## 2. アーキテクチャ設計

### 2.1 ディレクトリ構造

#### 2.1.1 変更前の構造
```
client/
├── src/
│   ├── app/
│   │   ├── page.tsx                    # Get Todayボタンあり
│   │   └── dm-users/
│   │       └── page.tsx                # ユーザー一覧表示のみ
│   ├── lib/
│   │   └── api.ts                      # APIクライアント（CSVダウンロードなし）
│   └── components/
│       └── TodayApiButton.tsx          # Get Todayボタン

server/
├── internal/
│   ├── api/
│   │   ├── handler/
│   │   │   └── dm_user_handler.go      # CRUD APIのみ
│   │   └── huma/
│   │       └── inputs.go                # 入力定義
│   ├── service/
│   │   └── dm_user_service.go          # ビジネスロジック
│   └── repository/
│       └── dm_user_repository.go       # データアクセス
```

#### 2.1.2 変更後の構造
```
client/
├── src/
│   ├── app/
│   │   ├── page.tsx                    # Get Todayボタンあり（変更なし）
│   │   └── dm-users/
│   │       └── page.tsx                # CSVダウンロードボタン追加
│   ├── lib/
│   │   └── api.ts                      # downloadUsersCSV()メソッド追加
│   └── components/
│       └── TodayApiButton.tsx          # 変更なし

server/
├── internal/
│   ├── api/
│   │   ├── handler/
│   │   │   └── dm_user_handler.go      # CSVダウンロードエンドポイント追加
│   │   └── huma/
│   │       └── inputs.go                # 変更なし
│   ├── config/
│   │   └── config.go                   # ServerConfigにIdleTimeout追加
│   ├── service/
│   │   └── dm_user_service.go          # 変更なし（既存メソッドを使用）
│   └── repository/
│       └── dm_user_repository.go       # 変更なし（既存メソッドを使用）
├── cmd/
│   └── server/
│       └── main.go                     # IdleTimeout設定追加
└── config/
    ├── develop/
    │   └── config.yaml                 # idle_timeout追加
    └── production/
        └── config.yaml.example         # idle_timeout追加
```

### 2.2 ファイル構成

#### 2.2.1 変更ファイル
- **クライアント側**:
  - `client/src/app/dm-users/page.tsx`: CSVダウンロードボタンを追加
  - `client/src/lib/api.ts`: `downloadUsersCSV()`メソッドを追加
- **サーバー側**:
  - `server/internal/api/handler/dm_user_handler.go`: CSVダウンロードエンドポイントを追加
  - `server/internal/config/config.go`: `ServerConfig`に`IdleTimeout`フィールドを追加
  - `server/cmd/server/main.go`: `IdleTimeout`設定を追加
  - `config/develop/config.yaml`: `idle_timeout`設定を追加
  - `config/production/config.yaml.example`: `idle_timeout`設定を追加

#### 2.2.2 新規作成ファイル
- なし（既存ファイルに追加のみ）

### 2.3 システム構成図

```
┌─────────────────────────────────────────────────────────┐
│              クライアント（Next.js）                      │
│  ┌──────────────────────────────────────────────────┐  │
│  │  dm-users/page.tsx                               │  │
│  │  - CSVダウンロードボタン                           │  │
│  └──────────────────┬───────────────────────────────┘  │
│                     │                                    │
│                     ▼                                    │
│  ┌──────────────────────────────────────────────────┐  │
│  │  api.ts (ApiClient)                               │  │
│  │  - downloadUsersCSV()                             │  │
│  │  - fetch() → blob() → ダウンロード                  │  │
│  └──────────────────┬───────────────────────────────┘  │
└─────────────────────┼──────────────────────────────────┘
                    │ HTTP GET /api/dm-users/csv
                    │ Authorization: Bearer {JWT}
                    ▼
┌─────────────────────────────────────────────────────────┐
│              APIサーバー（Go + Huma）                     │
│  ┌──────────────────────────────────────────────────┐  │
│  │  dm_user_handler.go                               │  │
│  │  - RegisterDmUserEndpoints()                    │  │
│  │  - GET /api/dm-users/csv                         │  │
│  │  - huma.StreamResponse                            │  │
│  └──────────────────┬───────────────────────────────┘  │
│                     │                                    │
│                     ▼                                    │
│  ┌──────────────────────────────────────────────────┐  │
│  │  dm_user_service.go                              │  │
│  │  - ListDmUsers(ctx, 20, 0)                      │  │
│  └──────────────────┬───────────────────────────────┘  │
│                     │                                    │
│                     ▼                                    │
│  ┌──────────────────────────────────────────────────┐  │
│  │  dm_user_repository.go                          │  │
│  │  - List(ctx, 20, 0)                             │  │
│  │  - 全シャードからデータ取得                        │  │
│  └──────────────────┬───────────────────────────────┘  │
│                     │                                    │
│                     ▼                                    │
│  ┌──────────────────────────────────────────────────┐  │
│  │  データベース（シャーディング）                     │  │
│  │  - dm_users_000 ～ dm_users_031                  │  │
│  └──────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────────────────┐
│            CSVストリーミングレスポンス                     │
│  - Content-Type: text/csv; charset=utf-8                │
│  - Content-Disposition: attachment; filename="..."     │
│  - タイムアウト: 3分                                      │
│  - ヘッダー行 + データ行（ストリーミング）                  │
└─────────────────────────────────────────────────────────┘
```

### 2.4 データフロー

#### 2.4.1 CSVダウンロードフロー
```
1. ユーザーがCSVダウンロードボタンをクリック
    ↓
2. apiClient.downloadUsersCSV()を呼び出し
    ↓
3. fetch('/api/dm-users/csv')でGETリクエスト送信
    ↓
4. APIサーバー側で認証チェック（public API）
    ↓
5. DmUserService.ListDmUsers(ctx, 20, 0)を呼び出し
    ↓
6. DmUserRepository.List(ctx, 20, 0)で全シャードからデータ取得
    ↓
7. huma.StreamResponseでCSVをストリーミング生成
    ↓
8. タイムアウト設定（3分）
    ↓
9. CSVヘッダー行を書き込み
    ↓
10. ユーザーデータを1件ずつCSV行として書き込み
    ↓
11. クライアント側でBlobとして受信
    ↓
12. URL.createObjectURL()でBlob URLを生成
    ↓
13. <a>要素でダウンロードを実行
    ↓
14. URL.revokeObjectURL()でBlob URLを解放
```

## 3. コンポーネント設計

### 3.1 クライアント側コンポーネント

#### 3.1.1 CSVダウンロードボタン（dm-users/page.tsx）

**実装内容**:
```typescript
const [downloading, setDownloading] = useState(false)
const [downloadError, setDownloadError] = useState<string | null>(null)

const handleDownloadCSV = async () => {
  try {
    setDownloading(true)
    setDownloadError(null)
    await apiClient.downloadUsersCSV()
  } catch (err) {
    setDownloadError(err instanceof Error ? err.message : 'Failed to download CSV')
  } finally {
    setDownloading(false)
  }
}

// ボタンの配置（Get Todayボタンの下辺り）
<button
  onClick={handleDownloadCSV}
  disabled={downloading}
  className="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600 disabled:opacity-50"
>
  {downloading ? 'ダウンロード中...' : 'CSVダウンロード'}
</button>
```

**配置位置**: Get Todayボタン（`client/src/app/page.tsx`のTodayApiButtonコンポーネント）の下辺り

**状態管理**:
- `downloading`: ダウンロード中の状態を管理
- `downloadError`: エラーメッセージを管理

**エラーハンドリング**:
- ネットワークエラー: エラーメッセージを表示
- HTTPエラー: エラーレスポンスのメッセージを表示

#### 3.1.2 APIクライアント拡張（api.ts）

**実装内容**:
```typescript
async downloadUsersCSV(): Promise<void> {
  const url = `${this.baseURL}/api/dm-users/csv`
  const token = this.apiKey

  const response = await fetch(url, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  })

  if (!response.ok) {
    if (response.status === 401 || response.status === 403) {
      const errorData = await response.json().catch(() => ({}))
      throw new Error(errorData.message || response.statusText)
    }
    const error = await response.text()
    throw new Error(error || response.statusText)
  }

  // Blobとして取得
  const blob = await response.blob()

  // Content-Dispositionヘッダーからファイル名を取得
  const contentDisposition = response.headers.get('Content-Disposition')
  let filename = 'dm-users.csv'
  if (contentDisposition) {
    const filenameMatch = contentDisposition.match(/filename="?(.+?)"?$/i)
    if (filenameMatch) {
      filename = filenameMatch[1]
    }
  }

  // Blob URLを生成
  const blobUrl = URL.createObjectURL(blob)

  // <a>要素を作成してダウンロード
  const link = document.createElement('a')
  link.href = blobUrl
  link.download = filename
  document.body.appendChild(link)
  link.click()

  // クリーンアップ
  document.body.removeChild(link)
  URL.revokeObjectURL(blobUrl)
}
```

**実装の詳細**:
- `fetch()`を使用してAPIを呼び出し
- レスポンスを`blob()`で取得
- Content-Dispositionヘッダーからファイル名を取得（デフォルト: `dm-users.csv`）
- `URL.createObjectURL()`でBlob URLを生成
- `<a>`要素を作成し、`download`属性を設定
- プログラム的にクリックイベントを発火してダウンロードを開始
- ダウンロード後、Blob URLを`URL.revokeObjectURL()`で解放

### 3.2 サーバー側コンポーネント

#### 3.2.1 CSVダウンロードエンドポイント（dm_user_handler.go）

**実装内容**:
```go
// GET /api/dm-users/csv - CSVダウンロード
huma.Register(api, huma.Operation{
    OperationID: "download-users-csv",
    Method:      http.MethodGet,
    Path:        "/api/dm-users/csv",
    Summary:     "ユーザー情報をCSV形式でダウンロード",
    Description: "**Access Level:** `public` (Public API Key JWT または Auth0 JWT でアクセス可能)",
    Tags:        []string{"users"},
    Security: []map[string][]string{
        {"bearerAuth": {}},
    },
}, func(ctx context.Context, input *struct{}) (*huma.StreamResponse, error) {
    // 公開レベルのチェック（publicエンドポイント）
    if err := auth.CheckAccessLevel(ctx, auth.AccessLevelPublic); err != nil {
        return nil, huma.Error403Forbidden(err.Error())
    }

    // ユーザー情報20件を取得
    users, err := h.dmUserService.ListDmUsers(ctx, 20, 0)
    if err != nil {
        return nil, huma.Error500InternalServerError(err.Error())
    }

    // ストリーミングレスポンスを返す
    return &huma.StreamResponse{
        ContentType: "text/csv; charset=utf-8",
        Headers: map[string]string{
            "Content-Disposition": `attachment; filename="dm-users.csv"`,
        },
        Body: func(w io.Writer) {
            // http.ResponseWriterを取り出す
            rw, ok := w.(http.ResponseWriter)
            if ok {
                rc := http.NewResponseController(rw)
                rc.SetWriteDeadline(time.Now().Add(3 * time.Minute))
            }

            // CSVエンコーダーを作成
            csvWriter := csv.NewWriter(w)
            defer csvWriter.Flush()

            // ヘッダー行を書き込み
            if err := csvWriter.Write([]string{
                "ID",
                "Name",
                "Email",
                "Created At",
                "Updated At",
            }); err != nil {
                // エラーはログに記録（ストリーミング中はエラーレスポンスを返せない）
                return
            }

            // ユーザーデータを1件ずつCSV行として書き込み
            for _, user := range users {
                if err := csvWriter.Write([]string{
                    user.ID,
                    user.Name,
                    user.Email,
                    user.CreatedAt.Format(time.RFC3339),
                    user.UpdatedAt.Format(time.RFC3339),
                }); err != nil {
                    // エラーはログに記録
                    return
                }
            }
        },
    }, nil
})
```

**実装の詳細**:
- **エンドポイント**: `GET /api/dm-users/csv`
- **アクセスレベル**: `public` (Public API Key JWT または Auth0 JWT でアクセス可能)
- **認証チェック**: `auth.CheckAccessLevel(ctx, auth.AccessLevelPublic)`でチェック
- **データ取得**: `DmUserService.ListDmUsers(ctx, 20, 0)`で20件取得
- **ストリーミング**: `huma.StreamResponse`を使用
- **タイムアウト設定**: `http.NewResponseController`で3分のタイムアウトを設定
- **CSV生成**: `encoding/csv`パッケージの`csv.NewWriter`を使用
- **ヘッダー設定**: Content-TypeとContent-Dispositionを設定

#### 3.2.2 CSV形式の仕様

**実装内容**:
- **文字エンコーディング**: UTF-8
- **区切り文字**: カンマ（`,`）
- **改行文字**: LF（`\n`）- `encoding/csv`パッケージが自動的に処理
- **ヘッダー行**: `ID, Name, Email, Created At, Updated At`
- **日時フォーマット**: ISO 8601形式（`time.RFC3339` = `2006-01-02T15:04:05Z07:00`）
- **特殊文字のエスケープ**: `encoding/csv`パッケージが自動的に処理（ダブルクォートで囲む）

**CSV出力例**:
```csv
ID,Name,Email,Created At,Updated At
550e8400e29b41d4a716446655440000,John Doe,john@example.com,2025-01-27T10:00:00Z,2025-01-27T10:00:00Z
550e8400e29b41d4a716446655440001,Jane Smith,jane@example.com,2025-01-27T10:01:00Z,2025-01-27T10:01:00Z
```

## 4. エラーハンドリング設計

### 4.1 クライアント側のエラーハンドリング

#### 4.1.1 ネットワークエラー
- **エラー種別**: `fetch()`のネットワークエラー
- **処理**: `catch`ブロックでエラーをキャッチし、エラーメッセージを表示
- **ユーザーへの表示**: "Failed to download CSV: {エラーメッセージ}"

#### 4.1.2 HTTPエラー（4xx, 5xx）
- **エラー種別**: HTTPステータスコードが200以外
- **処理**: 
  - 401/403エラー: JSONレスポンスからエラーメッセージを取得
  - その他のエラー: テキストレスポンスからエラーメッセージを取得
- **ユーザーへの表示**: エラーメッセージを表示

#### 4.1.3 ダウンロード失敗
- **エラー種別**: Blob生成、URL生成、ダウンロード実行時のエラー
- **処理**: `catch`ブロックでエラーをキャッチし、エラーメッセージを表示
- **ユーザーへの表示**: "Failed to download CSV: {エラーメッセージ}"

### 4.2 サーバー側のエラーハンドリング

#### 4.2.1 認証エラー
- **エラー種別**: `auth.CheckAccessLevel()`がエラーを返す
- **処理**: `huma.Error403Forbidden(err.Error())`を返す
- **HTTPステータス**: 403 Forbidden

#### 4.2.2 データベースエラー
- **エラー種別**: `DmUserService.ListDmUsers()`がエラーを返す
- **処理**: `huma.Error500InternalServerError(err.Error())`を返す
- **HTTPステータス**: 500 Internal Server Error

#### 4.2.3 CSV生成エラー
- **エラー種別**: CSV書き込み時のエラー（`csvWriter.Write()`がエラーを返す）
- **処理**: エラーをログに記録（ストリーミング中はエラーレスポンスを返せないため）
- **注意**: ストリーミング開始後はHTTPステータスコードを変更できないため、エラーはログに記録するのみ

#### 4.2.4 タイムアウトエラー
- **エラー種別**: 3分のタイムアウトに達した場合
- **処理**: `http.ResponseController.SetWriteDeadline()`により自動的に接続が切断される
- **HTTPステータス**: クライアント側でタイムアウトエラーとして検出される

## 5. テスト設計

### 5.1 単体テスト

#### 5.1.1 サーバー側の単体テスト

**テストファイル**: `server/internal/api/handler/dm_user_handler_test.go`

**テストケース**:
1. **正常系**: CSVダウンロードが正常に動作する
   - ユーザー情報20件を取得
   - CSV形式で正しく出力される
   - Content-TypeとContent-Dispositionが正しく設定される

2. **認証エラー**: 認証が失敗した場合
   - 403 Forbiddenが返される

3. **データベースエラー**: データ取得が失敗した場合
   - 500 Internal Server Errorが返される

4. **空データ**: ユーザーが0件の場合
   - ヘッダー行のみのCSVが返される

**実装例**:
```go
func TestDmUserHandler_DownloadUsersCSV(t *testing.T) {
    // テストの実装
    // - モックサービスを作成
    // - ハンドラーを作成
    // - リクエストを送信
    // - レスポンスを検証
}
```

#### 5.1.2 クライアント側の単体テスト

**テストファイル**: `client/src/lib/__tests__/api.test.ts`

**テストケース**:
1. **正常系**: CSVダウンロードが正常に動作する
   - `fetch()`が正常に呼び出される
   - Blobが正しく生成される
   - ダウンロードが実行される

2. **エラーハンドリング**: エラーが発生した場合
   - エラーが正しくスローされる

**実装例**:
```typescript
describe('downloadUsersCSV', () => {
  it('should download CSV file', async () => {
    // テストの実装
    // - fetch()をモック
    // - Blobをモック
    // - ダウンロードを検証
  })
})
```

### 5.2 E2Eテスト

#### 5.2.1 Playwrightテスト

**テストファイル**: `client/e2e/csv-download.spec.ts`

**テストケース**:
1. **正常系**: CSVダウンロードが正常に動作する
   - dm-usersページにアクセス
   - CSVダウンロードボタンをクリック
   - CSVファイルがダウンロードされる
   - CSVファイルの内容を検証

2. **エラーハンドリング**: エラーが発生した場合
   - エラーメッセージが表示される

**実装例**:
```typescript
import { test, expect } from '@playwright/test'

test('should download CSV file', async ({ page }) => {
  // テストの実装
  // - ページにアクセス
  // - ボタンをクリック
  // - ダウンロードを待機
  // - ファイルの内容を検証
})
```

## 6. パフォーマンス設計

### 6.1 メモリ効率

#### 6.1.1 ストリーミング方式の採用
- **理由**: ファイルを作成せず、HTTPレスポンスに直接ストリーミングすることでメモリ使用量を最小限に抑える
- **実装**: `huma.StreamResponse`を使用し、データを1件ずつCSV行として書き込み

#### 6.1.2 データ取得の最適化
- **取得件数**: 20件固定（要件定義に従う）
- **取得方法**: `DmUserService.ListDmUsers(ctx, 20, 0)`を使用
- **シャーディング対応**: 全シャードからデータを取得（既存の`List`メソッドを使用）

### 6.2 レスポンス時間

#### 6.2.1 目標値
- **目標**: 20件のデータ取得とCSV生成は1秒以内に完了すること
- **実現方法**: 
  - ストリーミング方式により、データ取得とCSV生成を並行して実行
  - データベースクエリの最適化（既存のインデックスを活用）

### 6.3 タイムアウト設定

#### 6.3.1 タイムアウト時間
- **設定値**: 3分（`time.Now().Add(3 * time.Minute)`）
- **理由**: 大量データのダウンロードにも対応できるよう、十分な時間を確保
- **実装**: `http.NewResponseController.SetWriteDeadline()`を使用

## 7. セキュリティ設計

### 7.1 認証・認可

#### 7.1.1 認証方式
- **方式**: Public API Key JWTまたはAuth0 JWTによる認証
- **実装**: 既存の`auth.CheckAccessLevel(ctx, auth.AccessLevelPublic)`を使用

#### 7.1.2 アクセスレベル
- **レベル**: `public` (Public API Key JWT または Auth0 JWT でアクセス可能)
- **理由**: 既存のdm-users APIと同じアクセスレベルを維持

### 7.2 入力検証

#### 7.2.1 パラメータ検証
- **パラメータ**: なし（エンドポイントはパラメータを受け取らない）
- **検証**: 不要

### 7.3 出力検証

#### 7.3.1 CSV形式の検証
- **検証内容**: 
  - 特殊文字のエスケープ（`encoding/csv`パッケージが自動的に処理）
  - 文字エンコーディング（UTF-8）
  - 改行文字（LF）

## 8. 保守性設計

### 8.1 コード品質

#### 8.1.1 コードスタイル
- **方針**: 既存のコードスタイルに従う
- **実装**: 
  - Go: `gofmt`、`golint`に従う
  - TypeScript: ESLint、Prettierに従う

#### 8.1.2 コメント
- **方針**: 既存のコメントスタイルに従う
- **実装**: 
  - Go: パブリック関数にはコメントを追加
  - TypeScript: 複雑な処理にはコメントを追加

### 8.2 ドキュメント

#### 8.2.1 APIドキュメント
- **実装**: Huma APIのOpenAPI仕様に自動的に追加される
- **内容**: 
  - エンドポイント: `GET /api/dm-users/csv`
  - 説明: ユーザー情報をCSV形式でダウンロード
  - アクセスレベル: `public`

#### 8.2.2 コードドキュメント
- **実装**: 既存のコードドキュメントスタイルに従う
- **内容**: 
  - 関数の説明
  - パラメータの説明
  - 戻り値の説明

## 9. 実装順序

### 9.1 フェーズ1: サーバー側の実装
1. APIサーバー全体のタイムアウト設定の追加（IdleTimeout）
2. CSVダウンロードエンドポイントの実装
3. ストリーミングCSV生成の実装
4. CSVダウンロード用の個別タイムアウト設定の実装
5. 単体テストの実装

### 9.2 フェーズ2: クライアント側の実装
1. APIクライアントの拡張（`downloadUsersCSV()`メソッド）
2. CSVダウンロードボタンの追加
3. エラーハンドリングの実装
4. 単体テストの実装

### 9.3 フェーズ3: 統合テスト
1. E2Eテストの実装
2. 動作確認
3. 既存テストの確認

## 10. 注意事項

### 10.1 ストリーミング中のエラー処理
- **制約**: ストリーミング開始後はHTTPステータスコードを変更できない
- **対応**: エラーはログに記録するのみ（ユーザーへの通知は困難）

### 10.2 タイムアウト設定

#### 10.2.1 既存のタイムアウト設定
- **ReadTimeout**: 30秒（`config/develop/config.yaml`で設定）
- **WriteTimeout**: 30秒（`config/develop/config.yaml`で設定）
- **IdleTimeout**: 未設定（設定が必要）

#### 10.2.2 APIサーバー全体のタイムアウト設定
- **実装場所**: `server/cmd/server/main.go`
- **既存設定**: `e.Server.ReadTimeout`と`e.Server.WriteTimeout`が設定されている
- **追加設定**: `e.Server.IdleTimeout`を設定する必要がある
- **デフォルト値**: 
  - `IdleTimeout`: 120秒（2分）- Keep-Alive接続のアイドルタイムアウト

#### 10.2.3 CSVダウンロード用の個別タイムアウト設定
- **実装方式**: `http.NewResponseController.SetWriteDeadline()`を使用
- **タイムアウト時間**: 3分（`time.Now().Add(3 * time.Minute)`）
- **理由**: 既存のWriteTimeout（30秒）では不十分なため、個別にタイムアウトを設定
- **注意**: 
  - `http.NewResponseController`はGo 1.20以降で利用可能
  - プロジェクトのGoバージョンを確認する必要がある
  - 個別タイムアウト設定により、既存のWriteTimeoutを上書きできる

#### 10.2.4 設定ファイルへの追加
- **ファイル**: `config/develop/config.yaml`、`config/production/config.yaml.example`
- **追加項目**: 
  ```yaml
  server:
    port: 8080
    read_timeout: 30s
    write_timeout: 30s
    idle_timeout: 120s  # 新規追加
  ```
- **設定構造体**: `server/internal/config/config.go`の`ServerConfig`に`IdleTimeout`フィールドを追加
- **実装**: `server/cmd/server/main.go`で`e.Server.IdleTimeout = cfg.Server.IdleTimeout`を設定

### 10.3 データ取得の順序
- **要件**: 作成日時の降順（最新の20件）
- **実装**: `DmUserService.ListDmUsers()`の実装を確認し、必要に応じて修正

### 10.4 CSV形式の互換性
- **BOM（Byte Order Mark）**: BOMは追加しない
- **理由**: 標準的なCSV形式（UTF-8 without BOM）を使用する
- **注意**: Excelなどの一部のアプリケーションはBOMを期待する場合があるが、本実装ではBOMは追加しない
