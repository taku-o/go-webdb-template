# Clientサーバー死活監視エンドポイントの実装タスク一覧

## 概要
Clientサーバー（Next.js）に`/health`エンドポイントを実装するためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: Route Handlerの実装

#### - [ ] タスク 1.1: `client/app/health/route.ts`の作成
**目的**: Next.js 14のApp RouterのRoute Handlerを使用して`/health`エンドポイントを実装する。

**作業内容**:
- `client/app/health/`ディレクトリを作成
- `client/app/health/route.ts`ファイルを作成
- `NextResponse`をインポート
- `GET`関数をエクスポートして、`200 OK`と`"OK"`を返す実装を追加

**実装コード**:
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

**受け入れ基準**:
- `client/app/health/route.ts`ファイルが作成されている
- `GET`関数がエクスポートされている
- ステータスコードが`200 OK`である
- レスポンスボディが`"OK"`である
- Content-Typeが`text/plain`である
- `/health`でアクセスできる（`/api/health`ではない）

- _Requirements: 3.1.1, 3.1.2, 3.1.3, 6.1_
- _Design: 3.1.1, 3.1.2, 3.1.3, 3.1.4_

---

### Phase 2: 動作確認

#### - [ ] タスク 2.1: 開発サーバーでの動作確認
**目的**: 開発サーバーで`/health`エンドポイントが正常に動作することを確認する。

**作業内容**:
- `client`ディレクトリに移動
- `npm run dev`で開発サーバーを起動
- 別のターミナルで`curl http://localhost:3000/health`を実行
- レスポンスが`OK`であることを確認
- ステータスコードが`200 OK`であることを確認
- Content-Typeが`text/plain`であることを確認

**確認コマンド**:
```bash
cd client
npm run dev

# 別のターミナルで
curl http://localhost:3000/health
curl -v http://localhost:3000/health
```

**期待される結果**:
```
OK
```

**レスポンスヘッダー**:
```
HTTP/1.1 200 OK
Content-Type: text/plain
...
```

**受け入れ基準**:
- 開発サーバーで`/health`エンドポイントが正常に動作する
- レスポンスが`OK`である
- ステータスコードが`200 OK`である
- Content-Typeが`text/plain`である

- _Requirements: 6.2_
- _Design: 5.1.1_

---

#### - [ ] タスク 2.2: 本番サーバーでの動作確認
**目的**: 本番サーバーで`/health`エンドポイントが正常に動作することを確認する。

**作業内容**:
- `client`ディレクトリに移動
- `npm run build`でビルドを実行
- `npm start`で本番サーバーを起動
- 別のターミナルで`curl http://localhost:3000/health`を実行
- レスポンスが`OK`であることを確認
- ステータスコードが`200 OK`であることを確認
- Content-Typeが`text/plain`であることを確認

**確認コマンド**:
```bash
cd client
npm run build
npm start

# 別のターミナルで
curl http://localhost:3000/health
curl -v http://localhost:3000/health
```

**期待される結果**:
```
OK
```

**レスポンスヘッダー**:
```
HTTP/1.1 200 OK
Content-Type: text/plain
...
```

**受け入れ基準**:
- 本番サーバーで`/health`エンドポイントが正常に動作する
- レスポンスが`OK`である
- ステータスコードが`200 OK`である
- Content-Typeが`text/plain`である

- _Requirements: 6.2_
- _Design: 5.1.2_

---

#### - [ ] タスク 2.3: 既存機能への影響確認
**目的**: 既存のClientサーバーの機能が正常に動作することを確認する。

**作業内容**:
- 既存のページが正常に動作することを確認
- 既存のAPI Routeが正常に動作することを確認
- 既存の認証機能が正常に動作することを確認
- `/health`エンドポイントが既存の機能と競合しないことを確認

**確認項目**:
- 既存のページ（例: `/`, `/dm-users`, `/dm-posts`など）が正常に表示される
- 既存のAPI Route（例: `/api/auth/[...nextauth]`）が正常に動作する
- 認証機能が正常に動作する
- `/health`エンドポイントが既存のページやAPI Routeと競合しない

**受け入れ基準**:
- 既存のページが正常に動作する
- 既存のAPI Routeが正常に動作する
- 既存の認証機能が正常に動作する
- `/health`エンドポイントが既存の機能と競合しない

- _Requirements: 6.2, 7.2_
- _Design: 5.3_

---

### Phase 3: 一貫性の確認

#### - [ ] タスク 3.1: 他のサーバーとの一貫性確認
**目的**: 他のサーバー（API、Admin、JobQueue）と同様のレスポンス形式であることを確認する。

**作業内容**:
- APIサーバーの`/health`エンドポイントを確認（`curl http://localhost:8080/health`）
- Adminサーバーの`/health`エンドポイントを確認（`curl http://localhost:8081/health`）
- JobQueueサーバーの`/health`エンドポイントを確認（`curl http://localhost:8082/health`）
- Clientサーバーの`/health`エンドポイントを確認（`curl http://localhost:3000/health`）
- すべてのサーバーが同じレスポンス形式（`200 OK`と`"OK"`）を返すことを確認

**確認コマンド**:
```bash
# APIサーバー
curl http://localhost:8080/health

# Adminサーバー
curl http://localhost:8081/health

# JobQueueサーバー
curl http://localhost:8082/health

# Clientサーバー
curl http://localhost:3000/health
```

**期待される結果**:
すべてのサーバーが`OK`を返す。

**受け入れ基準**:
- すべてのサーバーが`200 OK`を返す
- すべてのサーバーが`"OK"`を返す
- すべてのサーバーが`text/plain`のContent-Typeを返す
- すべてのサーバーが認証不要でアクセス可能である

- _Requirements: 6.4_
- _Design: 6.1, 6.2_

---

### Phase 4: テスト

#### - [ ] タスク 4.1: 単体テストの実装
**目的**: `/health`エンドポイントの単体テストを実装する。

**作業内容**:
- `client/src/__tests__/routes/health-route.test.ts`ファイルを作成
- `GET`関数をインポート
- ステータスコード、レスポンスボディ、Content-Typeを確認するテストを実装
- 認証不要であることを確認するテストを実装

**実装コード**:
```typescript
// client/src/__tests__/routes/health-route.test.ts
import { GET } from '@/app/health/route'

describe('GET /health', () => {
  it('should return 200 OK with "OK" body', async () => {
    const response = await GET()
    
    expect(response.status).toBe(200)
    expect(await response.text()).toBe('OK')
    expect(response.headers.get('Content-Type')).toBe('text/plain')
  })

  it('should not require authentication', async () => {
    const response = await GET()
    
    // 認証エラーが発生しないことを確認
    expect(response.status).toBe(200)
  })
})
```

**受け入れ基準**:
- 単体テストが実装されている
- ステータスコードが`200 OK`であることを確認するテストが含まれている
- レスポンスボディが`"OK"`であることを確認するテストが含まれている
- Content-Typeが`text/plain`であることを確認するテストが含まれている
- 認証不要であることを確認するテストが含まれている
- テストが正常に実行される

- _Requirements: 6.3_
- _Design: 4.1_

---

#### - [ ] タスク 4.2: 既存テストの確認
**目的**: 既存のテストが全て失敗しないことを確認する。

**作業内容**:
- `client`ディレクトリに移動
- `npm test`を実行
- 既存のテストが全て正常に実行されることを確認
- 新規追加した`/health`エンドポイントが既存のテストに影響を与えないことを確認

**確認コマンド**:
```bash
cd client
npm test
```

**受け入れ基準**:
- 既存のテストが全て正常に実行される
- 新規追加した`/health`エンドポイントが既存のテストに影響を与えない

- _Requirements: 6.3_
- _Design: 4.1.3_

---

## 受け入れ基準の確認

### 要件定義書の受け入れ基準

#### 6.1 `/health`エンドポイントの実装
- [ ] `client/app/health/route.ts`に`GET`関数が実装されている
- [ ] エンドポイントが認証なしでアクセス可能である
- [ ] エンドポイントが`200 OK`と`"OK"`を返す
- [ ] エンドポイントが`text/plain`のContent-Typeを返す
- [ ] エンドポイントが`/health`でアクセスできる（`/api/health`ではない）

#### 6.2 動作確認
- [ ] ローカル環境で`curl http://localhost:3000/health`が正常に動作する
- [ ] 開発サーバー（`npm run dev`）でエンドポイントが動作する
- [ ] 本番サーバー（`npm run build && npm start`）でエンドポイントが動作する
- [ ] 既存のClientサーバーの機能が正常に動作することを確認

#### 6.3 テスト
- [ ] 単体テストが実装されている（該当する場合）
- [ ] 統合テストが実装されている（該当する場合）
- [ ] 既存のテストが全て失敗しないことを確認

#### 6.4 一貫性の確認
- [ ] 他のサーバー（API、Admin、JobQueue）と同様のレスポンス形式である
- [ ] 他のサーバーと同様に認証不要でアクセス可能である
- [ ] 起動サーバー一覧表示機能（0077-listapp）で`/health`エンドポイントを使用できる

---

## 実装順序

1. **Phase 1: Route Handlerの実装**
   - タスク 1.1: `client/app/health/route.ts`の作成

2. **Phase 2: 動作確認**
   - タスク 2.1: 開発サーバーでの動作確認
   - タスク 2.2: 本番サーバーでの動作確認
   - タスク 2.3: 既存機能への影響確認

3. **Phase 3: 一貫性の確認**
   - タスク 3.1: 他のサーバーとの一貫性確認

4. **Phase 4: テスト**
   - タスク 4.1: 単体テストの実装
   - タスク 4.2: 既存テストの確認

---

## 注意事項

### 実装時の注意事項
- Next.js 14のApp RouterのRoute Handlerを使用する
- `app/health/route.ts`を作成することで、`/health`エンドポイントを実装できる（`/api/health`ではない）
- 他のサーバーと同様に`"OK"`を返す
- 他のサーバーと同様に`text/plain`のContent-Typeを返す
- 認証不要でアクセス可能にする

### テスト時の注意事項
- 既存のテストが全て失敗しないことを確認する
- 必要に応じて、Playwrightテストを追加することも可能

### 動作確認時の注意事項
- 開発サーバーと本番サーバーの両方で動作確認を行う
- 既存の機能が正常に動作することを確認する
- 他のサーバーとの一貫性を確認する

---

## 参考情報

### 関連ドキュメント
- `.kiro/specs/0078-client-health/requirements.md`: 要件定義書
- `.kiro/specs/0078-client-health/design.md`: 設計書
- `.kiro/specs/0076-jobqueue-health/tasks.md`: JobQueueサーバーのタスク定義書（参考）

### 既存実装の参考
- `client/app/api/auth/[...nextauth]/route.ts`: 既存のAPI Routeの実装（参考）
- `server/internal/api/router/router.go`: APIサーバーの`/health`エンドポイント実装（参考）
- `server/cmd/admin/main.go`: Adminサーバーの`/health`エンドポイント実装（参考）
- `server/cmd/jobqueue/main.go`: JobQueueサーバーの`/health`エンドポイント実装（参考）
