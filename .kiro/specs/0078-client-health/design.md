# Clientサーバー死活監視エンドポイントの設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、Clientサーバー（Next.js）に`/health`エンドポイントを実装するための詳細設計を定義する。他のサーバー（API、Admin、JobQueue）と同様のシンプルな実装を行い、一貫性を保つ。

### 1.2 設計の範囲
- Clientサーバーに`/health`エンドポイントを実装する
- Next.js 14のApp RouterのRoute Handlerを使用する
- 認証不要でアクセス可能とする
- 単体テストを実装する

### 1.3 設計方針
- **シンプルな実装**: 他のサーバーと同様のシンプルな実装を維持する
- **認証不要**: ヘルスチェック用のため、認証ミドルウェアを通過しない
- **一貫性**: 他のサーバー（API、Admin、JobQueue）の実装パターンに合わせる
- **Next.js標準機能**: Next.js 14のApp RouterのRoute Handlerを使用する
- **テスト**: 単体テストを実装する

## 2. アーキテクチャ設計

### 2.1 サーバー構成

```
┌─────────────────────────────────────────────────────────────┐
│              Clientサーバー (Port 3000)                    │
│              Next.js 14 (App Router)                        │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  Next.js App Router              │
        │  - ページルーティング              │
        │  - API Route Handler              │
        │  - Route Handler                  │
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  /health Route Handler           │
        │  - app/health/route.ts           │
        │  - GET関数をエクスポート           │
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  HTTPレスポンス                    │
        │  - Status: 200 OK                 │
        │  - Body: "OK"                    │
        │  - Content-Type: text/plain       │
        └─────────────────────────────────┘
```

### 2.2 リクエストフロー

```
┌─────────────────────────────────────────────────────────────┐
│              /health エンドポイントのリクエストフロー           │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  HTTPリクエスト: GET /health     │
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  Next.js App Router               │
        │  - ルーティング処理                │
        │  - app/health/route.tsを解決      │
        │  - 認証ミドルウェアなし             │
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  Route Handler (GET関数)         │
        │  - NextResponseを返す           │
        │  - Status: 200 OK               │
        │  - Body: "OK"                   │
        │  - Content-Type: text/plain      │
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  HTTPレスポンス                  │
        │  - Status: 200 OK               │
        │  - Body: "OK"                   │
        │  - Content-Type: text/plain      │
        └─────────────────────────────────┘
```

### 2.3 ファイル構造

```
client/
└── app/
    ├── api/
    │   └── auth/
    │       └── [...nextauth]/
    │           └── route.ts  (既存)
    └── health/
        └── route.ts  (新規作成)
```

## 3. 実装設計

### 3.1 Route Handlerの実装

#### 3.1.1 実装場所
- **ファイル**: `client/app/health/route.ts`（新規作成）
- **ディレクトリ**: `client/app/health/`（新規作成）

#### 3.1.2 実装コード

```typescript
// client/app/health/route.ts
import { NextResponse } from 'next/server'

export async function GET() {
  return new NextResponse('OK', {
    status: 200,
    headers: {
      'Content-Type': 'text/plain',
    },
  })
}
```

#### 3.1.3 実装の詳細
- **Route Handler**: Next.js 14のApp RouterのRoute Handlerを使用
- **エクスポート**: `GET`関数をエクスポートすることで、`GET /health`エンドポイントを実装
- **レスポンス**: `NextResponse`を使用して、適切なステータスコードとContent-Typeを設定
- **認証**: 不要（認証ミドルウェアを通過しない）
- **エラーハンドリング**: 不要（常に成功を返す）

#### 3.1.4 パスの実現方法
- `app/health/route.ts`を作成することで、`/health`エンドポイントを実装
- `app/api/health/route.ts`ではない（`/api/health`ではなく`/health`でアクセスできるようにするため）

### 3.2 コード変更箇所

#### 3.2.1 新規作成するファイル

**ファイル**: `client/app/health/route.ts`

```typescript
import { NextResponse } from 'next/server'

export async function GET() {
  return new NextResponse('OK', {
    status: 200,
    headers: {
      'Content-Type': 'text/plain',
    },
  })
}
```

#### 3.2.2 変更不要なファイル
- `client/next.config.js`: 変更不要（Next.jsの標準機能を使用）
- `client/app/api/auth/[...nextauth]/route.ts`: 変更不要（既存のAPI Routeと競合しない）
- その他の既存ファイル: 変更不要

### 3.3 実装の詳細設計

#### 3.3.1 NextResponseの使用理由
- Next.js 14のApp Routerでは、`NextResponse`を使用してHTTPレスポンスを返す
- `Response`オブジェクトを直接返すことも可能だが、`NextResponse`を使用することで、Next.jsの機能を最大限に活用できる

#### 3.3.2 Content-Typeの設定
- 他のサーバー（API、Admin、JobQueue）と同様に`text/plain`を設定
- レスポンスボディは`"OK"`という文字列のみ

#### 3.3.3 エラーハンドリング
- ヘルスチェックエンドポイントは常に成功を返す
- エラーハンドリングは不要（Next.jsの起動状態を確認するため）

## 4. テスト設計

### 4.1 単体テスト

#### 4.1.1 テストファイル
- **ファイル**: `client/src/__tests__/api/health-route.test.ts`（新規作成）

#### 4.1.2 テスト内容

```typescript
// client/src/__tests__/api/health-route.test.ts
import { GET } from '@/app/health/route'
import { NextRequest } from 'next/server'

describe('GET /health', () => {
  it('should return 200 OK with "OK" body', async () => {
    const request = new NextRequest('http://localhost:3000/health')
    const response = await GET()
    
    expect(response.status).toBe(200)
    expect(await response.text()).toBe('OK')
    expect(response.headers.get('Content-Type')).toBe('text/plain')
  })

  it('should not require authentication', async () => {
    const request = new NextRequest('http://localhost:3000/health')
    const response = await GET()
    
    // 認証エラーが発生しないことを確認
    expect(response.status).toBe(200)
  })
})
```

#### 4.1.3 テストの詳細
- **ステータスコード**: `200 OK`を返すことを確認
- **レスポンスボディ**: `"OK"`を返すことを確認
- **Content-Type**: `text/plain`を返すことを確認
- **認証**: 認証不要であることを確認

### 4.2 統合テスト（オプション）

#### 4.2.1 Playwrightテスト（オプション）
起動サーバー一覧表示機能（0077-listapp）で使用されるため、統合テストは不要と判断。必要に応じて、Playwrightテストを追加することも可能。

## 5. 動作確認設計

### 5.1 ローカル環境での動作確認

#### 5.1.1 開発サーバーでの確認
```bash
# 開発サーバーを起動
cd client
npm run dev

# 別のターミナルで確認
curl http://localhost:3000/health
```

期待される結果:
```
OK
```

#### 5.1.2 本番サーバーでの確認
```bash
# 本番サーバーをビルドして起動
cd client
npm run build
npm start

# 別のターミナルで確認
curl http://localhost:3000/health
```

期待される結果:
```
OK
```

### 5.2 レスポンスヘッダーの確認

```bash
curl -v http://localhost:3000/health
```

期待される結果:
```
HTTP/1.1 200 OK
Content-Type: text/plain
...
OK
```

### 5.3 既存機能への影響確認

- 既存のページが正常に動作することを確認
- 既存のAPI Routeが正常に動作することを確認
- 既存の認証機能が正常に動作することを確認

## 6. 一貫性の確認

### 6.1 他のサーバーとの比較

| サーバー | パス | メソッド | 認証 | ステータス | ボディ | Content-Type |
|---------|------|---------|------|-----------|-------|--------------|
| API | `/health` | GET | 不要 | 200 OK | "OK" | text/plain |
| Admin | `/health` | GET | 不要 | 200 OK | "OK" | text/plain |
| JobQueue | `/health` | GET | 不要 | 200 OK | "OK" | text/plain |
| Client | `/health` | GET | 不要 | 200 OK | "OK" | text/plain |

### 6.2 実装パターンの比較

- **APIサーバー**: Echoフレームワークを使用
- **Adminサーバー**: Gorilla Mux Routerを使用
- **JobQueueサーバー**: 標準ライブラリの`net/http`を使用
- **Clientサーバー**: Next.js 14のApp RouterのRoute Handlerを使用

すべてのサーバーで同じレスポンス形式（`200 OK`と`"OK"`）を返すことで、一貫性を保つ。

## 7. 影響範囲

### 7.1 新規作成するファイル

- `client/app/health/route.ts`: `/health`エンドポイントの実装
- `client/src/__tests__/api/health-route.test.ts`: 単体テスト（オプション）

### 7.2 変更不要なファイル

- `client/next.config.js`: 変更不要
- `client/app/api/auth/[...nextauth]/route.ts`: 変更不要
- その他の既存ファイル: 変更不要

### 7.3 既存機能への影響

- **既存のページ**: 影響なし（`/health`は既存のページと競合しない）
- **既存のAPI Route**: 影響なし（`/health`は既存のAPI Routeと競合しない）
- **既存の認証機能**: 影響なし（`/health`は認証を通過しない）

## 8. 実装上の注意事項

### 8.1 Next.js App Routerの注意事項

- **Route Handler**: Next.js 14のApp RouterのRoute Handlerを使用する
- **パス**: `app/health/route.ts`を作成することで、`/health`エンドポイントを実装できる
- **エクスポート**: `GET`関数をエクスポートすることで、`GET /health`エンドポイントを実装
- **レスポンス**: `NextResponse`を使用して、適切なContent-Typeを設定する

### 8.2 一貫性の注意事項

- **レスポンス形式**: 他のサーバーと同様に`"OK"`を返す
- **Content-Type**: 他のサーバーと同様に`text/plain`を返す
- **認証**: 他のサーバーと同様に認証不要とする

### 8.3 テストの注意事項

- **単体テスト**: エンドポイントが正常に動作することを確認するテストを追加する
- **既存テスト**: 既存のテストが全て失敗しないことを確認する

### 8.4 動作確認の注意事項

- **開発サーバー**: 開発サーバー（`npm run dev`）でエンドポイントが動作することを確認
- **本番サーバー**: 本番サーバー（`npm run build && npm start`）でエンドポイントが動作することを確認
- **既存機能**: 既存のClientサーバーの機能が正常に動作することを確認

## 9. 参考情報

### 9.1 Next.js App Routerの参考

- **Route Handler**: Next.js 14のApp RouterのRoute Handlerを使用してAPIエンドポイントを実装
- **パス**: `app/health/route.ts`を作成することで、`/health`エンドポイントを実装できる
- **既存のAPI Route**: `client/app/api/auth/[...nextauth]/route.ts`を参考にする

### 9.2 既存実装の参考

- **APIサーバー**: `server/internal/api/router/router.go`の`/health`エンドポイント実装
- **Adminサーバー**: `server/cmd/admin/main.go`の`/health`エンドポイント実装
- **JobQueueサーバー**: `server/cmd/jobqueue/main.go`の`/health`エンドポイント実装

### 9.3 技術スタック

- **フレームワーク**: Next.js 14+ (App Router)
- **言語**: TypeScript 5+
- **Route Handler**: Next.jsの標準機能

### 9.4 関連ドキュメント

- `client/app/api/auth/[...nextauth]/route.ts`: 既存のAPI Routeの実装（参考）
- `client/next.config.js`: Next.jsの設定ファイル
- `.kiro/specs/0076-jobqueue-health/design.md`: JobQueueサーバーの`/health`エンドポイント実装の設計書（参考）
- `.kiro/specs/0078-client-health/requirements.md`: 本機能の要件定義書
