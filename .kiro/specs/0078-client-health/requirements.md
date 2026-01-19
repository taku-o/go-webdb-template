# Clientサーバー死活監視エンドポイントの要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Feature名**: 0078-client-health
- **作成日**: 2026-01-17

### 1.2 目的
APIサーバー、Adminサーバー、JobQueueサーバーには`/health`という死活監視用のエンドポイントが用意されている。Clientサーバー（Next.js）にも同様の`/health`エンドポイントを実装し、他のサーバーと一貫性を保つ。これにより、Clientサーバーの死活監視が可能になり、docker-compose等のヘルスチェック機能を有効化できる。

### 1.3 スコープ
- Clientサーバー（Next.js）に`/health`エンドポイントを実装する
- エンドポイントは認証不要でアクセス可能とする
- 他のサーバー（API、Admin、JobQueue）と同様の実装パターンを使用する

**本実装の範囲外**:
- 他のサーバーの`/health`エンドポイントの変更
- その他のエンドポイントの追加や変更
- ヘルスチェックの詳細な診断機能（現時点では単純なOK応答のみ）
- Clientサーバーの既存機能の変更

## 2. 背景・現状分析

### 2.1 現在の状況
- **APIサーバー**: `server/internal/api/router/router.go`に`/health`エンドポイントが実装されている
  - 実装内容: `GET /health`で`200 OK`と`"OK"`という文字列を返す
  - 認証: 不要
  - ポート: 8080
  - フレームワーク: Echo
- **Adminサーバー**: `server/cmd/admin/main.go`に`/health`エンドポイントが実装されている
  - 実装内容: `GET /health`で`200 OK`と`"OK"`という文字列を返す
  - 認証: 不要
  - ポート: 8081
  - ルーター: Gorilla Mux Router
- **JobQueueサーバー**: `server/cmd/jobqueue/main.go`に`/health`エンドポイントが実装されている
  - 実装内容: `GET /health`で`200 OK`と`"OK"`という文字列を返す
  - 認証: 不要
  - ポート: 8082
  - ルーター: 標準ライブラリの`net/http`
- **Clientサーバー**: `/health`エンドポイントが未実装
  - 現在の実装: Next.js 14のApp Routerを使用
  - ポート: 3000
  - エントリーポイント: `client/app/`
  - 既存のAPI Route: `client/app/api/auth/[...nextauth]/route.ts`など

### 2.2 課題点
1. **Clientサーバーのヘルスチェック不可**: docker-compose等のヘルスチェックが動作していない
2. **一貫性の欠如**: 他のサーバー（API、Admin、JobQueue）には`/health`エンドポイントがあるが、Clientサーバーにはない
3. **監視ツールとの連携不可**: 外部監視ツールがClientサーバーの死活監視を行うことができない
4. **運用上の問題**: Clientサーバーが正常に動作しているかどうかをHTTPリクエストで確認できない
5. **起動サーバー一覧表示機能との整合性**: 起動サーバー一覧表示機能（0077-listapp）では、Clientサーバーはポート接続のみで判定しているが、他のサーバーと同様に`/health`エンドポイントで判定できるようにする

### 2.3 本実装による改善点
1. **ヘルスチェック機能の有効化**: docker-compose等のヘルスチェックが正常に動作する
2. **一貫性の確保**: すべてのサーバー（API、Client、Admin、JobQueue）に`/health`エンドポイントが存在する
3. **監視ツールとの連携**: 外部監視ツールがClientサーバーの死活監視を行えるようになる
4. **運用の改善**: HTTPリクエストでClientサーバーの状態を確認できるようになる
5. **起動サーバー一覧表示機能の改善**: ポート接続だけでなく、`/health`エンドポイントでも判定できるようになる

## 3. 機能要件

### 3.1 `/health`エンドポイントの実装

#### 3.1.1 エンドポイント仕様
- **パス**: `/health`
- **メソッド**: `GET`
- **認証**: 不要（認証ミドルウェアを通過しない）
- **レスポンス**: 
  - ステータスコード: `200 OK`
  - レスポンスボディ: `"OK"`（文字列）
  - Content-Type: `text/plain`

#### 3.1.2 実装場所
- **ファイル**: `client/app/api/health/route.ts`（新規作成）
- **フレームワーク**: Next.js 14のApp Router
- **実装方法**: Next.jsのRoute Handlerを使用
  - 認証ミドルウェアを通過しない
  - 最小限の実装（`"OK"`を返すのみ）

#### 3.1.3 実装の詳細
- Next.js 14のApp RouterのRoute Handlerを使用
- `app/api/health/route.ts`に`GET`関数をエクスポート
- ハンドラー関数はシンプルに`200 OK`と`"OK"`を返す
- エラーハンドリングは不要（常に成功を返す）

### 3.2 実装方式

#### 3.2.1 Next.js App RouterのRoute Handler
Next.js 14のApp Routerでは、`app/api/health/route.ts`ファイルを作成し、`GET`関数をエクスポートすることで、`GET /api/health`エンドポイントを実装できる。

ただし、他のサーバーと一貫性を保つため、`/health`（`/api/health`ではなく）でアクセスできるようにする必要がある。

#### 3.2.2 パスの実現方法
Next.jsのApp Routerでは、`app/health/route.ts`を作成することで、`/health`エンドポイントを実装できる。これは`app/api/health/route.ts`とは異なり、`/api`プレフィックスが付かない。

#### 3.2.3 実装例
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

## 4. 非機能要件

### 4.1 パフォーマンス
- **レスポンス時間**: 1ms以下（シンプルな実装のため）
- **リソース使用量**: 最小限（追加のリソース消費は不要）
- **Next.jsのオーバーヘッド**: 最小限（Route Handlerを使用）

### 4.2 可用性
- **可用性**: サーバーが起動している限り、常に`200 OK`を返す
- **エラー処理**: エラーが発生しない実装（常に成功を返す）
- **Next.jsの起動状態**: Next.jsの開発サーバーまたは本番サーバーが起動している限り、エンドポイントが利用可能

### 4.3 セキュリティ
- **認証**: 不要（ヘルスチェック用のため）
- **情報漏洩**: サーバーの内部情報を返さない（`"OK"`のみ）
- **ポート**: 必要に応じてファイアウォールで制限可能

### 4.4 保守性
- **コードの簡潔性**: 他のサーバーと同様のシンプルな実装
- **一貫性**: 他のサーバー（API、Admin、JobQueue）の実装と同様のパターンを使用
- **Next.jsの標準機能**: Next.jsの標準機能（Route Handler）を使用

### 4.5 動作環境
- **開発環境**: Next.jsの開発サーバー（`npm run dev`）で動作する
- **本番環境**: Next.jsの本番サーバー（`npm run build && npm start`）で動作する
- **Docker環境**: docker-compose等で動作することを確認（将来の拡張）

## 5. 制約事項

### 5.1 技術的制約
- **Next.jsのApp Router**: Next.js 14のApp Routerを使用（既存の実装と一致）
- **ポート**: 3000を使用（既存の設定と一致）
- **既存機能への影響**: Clientサーバーの既存機能に影響を与えない

### 5.2 実装上の制約
- **認証ミドルウェア**: `/health`エンドポイントは認証を通過しない
- **既存コードへの影響**: 既存のClientサーバーの実装に影響を与えない
- **パス**: `/health`でアクセスできるようにする（`/api/health`ではない）

### 5.3 動作環境
- **ローカル環境**: ローカル環境でも動作することを確認
- **Docker環境**: docker-compose等で動作することを確認（将来の拡張）
- **開発サーバーと本番サーバー**: 両方の環境で動作することを確認

## 6. 受け入れ基準

### 6.1 `/health`エンドポイントの実装
- [ ] `client/app/health/route.ts`に`GET`関数が実装されている
- [ ] エンドポイントが認証なしでアクセス可能である
- [ ] エンドポイントが`200 OK`と`"OK"`を返す
- [ ] エンドポイントが`text/plain`のContent-Typeを返す
- [ ] エンドポイントが`/health`でアクセスできる（`/api/health`ではない）

### 6.2 動作確認
- [ ] ローカル環境で`curl http://localhost:3000/health`が正常に動作する
- [ ] 開発サーバー（`npm run dev`）でエンドポイントが動作する
- [ ] 本番サーバー（`npm run build && npm start`）でエンドポイントが動作する
- [ ] 既存のClientサーバーの機能が正常に動作することを確認

### 6.3 テスト
- [ ] 単体テストが実装されている（該当する場合）
- [ ] 統合テストが実装されている（該当する場合）
- [ ] 既存のテストが全て失敗しないことを確認

### 6.4 一貫性の確認
- [ ] 他のサーバー（API、Admin、JobQueue）と同様のレスポンス形式である
- [ ] 他のサーバーと同様に認証不要でアクセス可能である
- [ ] 起動サーバー一覧表示機能（0077-listapp）で`/health`エンドポイントを使用できる

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### 新規作成するファイル
- `client/app/health/route.ts`: `/health`エンドポイントの実装（新規作成）

#### 確認が必要なファイル
- `client/app/api/auth/[...nextauth]/route.ts`: 既存のAPI Routeの実装を確認（参考）
- `client/next.config.js`: Next.jsの設定を確認（変更不要）

### 7.2 既存機能への影響
- **既存のClientサーバー**: 影響なし（新規エンドポイントの追加のみ）
- **既存のAPI Route**: 影響なし（`/health`は既存のAPI Routeと競合しない）
- **既存のページ**: 影響なし（`/health`は既存のページと競合しない）

### 7.3 テストへの影響
- **既存のテスト**: 影響なし（新規エンドポイントの追加のみ）
- **新規テスト**: `/health`エンドポイントのテストを追加する可能性がある

### 7.4 起動サーバー一覧表示機能への影響
- **0077-listapp**: Clientサーバーの確認方法をポート接続から`/health`エンドポイントに変更できる（将来的な改善）

## 8. 実装上の注意事項

### 8.1 Next.js App Router実装の注意事項
- **Route Handler**: Next.js 14のRoute Handlerを使用する
- **パス**: `app/health/route.ts`を作成することで、`/health`エンドポイントを実装できる
- **レスポンス形式**: `NextResponse`を使用して、適切なContent-Typeを設定する
- **エラーハンドリング**: エラーが発生しない実装（常に成功を返す）

### 8.2 一貫性の注意事項
- **レスポンス形式**: 他のサーバーと同様に`"OK"`を返す
- **Content-Type**: 他のサーバーと同様に`text/plain`を返す
- **認証**: 他のサーバーと同様に認証不要とする

### 8.3 テストの注意事項
- **単体テスト**: エンドポイントが正常に動作することを確認するテストを追加する
- **統合テスト**: Next.jsの開発サーバーと本番サーバーで動作することを確認する
- **既存テスト**: 既存のテストが全て失敗しないことを確認する

### 8.4 動作確認の注意事項
- **ローカル環境**: ローカル環境で`curl http://localhost:3000/health`が正常に動作することを確認
- **開発サーバー**: 開発サーバー（`npm run dev`）でエンドポイントが動作することを確認
- **本番サーバー**: 本番サーバー（`npm run build && npm start`）でエンドポイントが動作することを確認
- **既存機能**: 既存のClientサーバーの機能が正常に動作することを確認

## 9. 参考情報

### 9.1 既存実装の参考
- **APIサーバー**: `server/internal/api/router/router.go`の`/health`エンドポイント実装
  ```go
  e.GET("/health", func(c echo.Context) error {
      return c.String(http.StatusOK, "OK")
  })
  ```
- **Adminサーバー**: `server/cmd/admin/main.go`の`/health`エンドポイント実装
  ```go
  app.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
      w.Header().Set("Content-Type", "text/plain")
      w.WriteHeader(http.StatusOK)
      w.Write([]byte("OK"))
  }).Methods("GET")
  ```
- **JobQueueサーバー**: `server/cmd/jobqueue/main.go`の`/health`エンドポイント実装
  ```go
  mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
      w.Header().Set("Content-Type", "text/plain")
      w.WriteHeader(http.StatusOK)
      w.Write([]byte("OK"))
  })
  ```

### 9.2 Next.js App Routerの参考
- **Route Handler**: Next.js 14のRoute Handlerを使用してAPIエンドポイントを実装
- **パス**: `app/health/route.ts`を作成することで、`/health`エンドポイントを実装できる
- **既存のAPI Route**: `client/app/api/auth/[...nextauth]/route.ts`を参考にする

### 9.3 技術スタック
- **フレームワーク**: Next.js 14+ (App Router)
- **言語**: TypeScript 5+
- **Route Handler**: Next.jsの標準機能

### 9.4 関連ドキュメント
- `client/app/api/auth/[...nextauth]/route.ts`: 既存のAPI Routeの実装（参考）
- `client/next.config.js`: Next.jsの設定ファイル
- `.kiro/specs/0076-jobqueue-health/requirements.md`: JobQueueサーバーの`/health`エンドポイント実装の要件定義書（参考）
- `.kiro/specs/0077-listapp/requirements.md`: 起動サーバー一覧表示機能の要件定義書（参考）
