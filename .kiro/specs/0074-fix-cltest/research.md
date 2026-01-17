# Research & Design Decisions

## Summary
- **Feature**: `0074-fix-cltest`
- **Discovery Scope**: Simple Addition（既存テストのモック方式変更）
- **Key Findings**:
  - MSW依存テストが `until-async` パッケージのESMエクスポート問題でJest実行時にエラー
  - 成功しているテストは `jest.mock('@/lib/api', ...)` を使用しており、MSWを使用していない
  - 5件の失敗テストを成功パターンに合わせて修正する必要がある

## Research Log

### MSW + Jest + ESM互換性問題
- **Context**: `npm run test:client` 実行時に5件のテストが `SyntaxError: Unexpected token 'export'` で失敗
- **Sources Consulted**: テスト実行ログ、jest.config.js、失敗テストファイル
- **Findings**:
  - エラーは `node_modules/.pnpm/until-async@3.0.2/node_modules/until-async/lib/index.js:23` で発生
  - `until-async` はMSWの依存パッケージ
  - `jest.config.js` の `transformIgnorePatterns` に `until-async` が含まれているが、Next.jsのJest設定と競合している可能性がある
- **Implications**: MSWの使用を避け、`jest.mock` 方式に統一することで問題を解決可能

### 成功テストパターンの分析
- **Context**: 成功しているテストファイルのモック方式を分析
- **Sources Consulted**: `dm-feed-post-page.test.tsx`, `dm-feed-user-page.test.tsx`
- **Findings**:
  - `jest.mock('@/lib/api', ...)` でAPIクライアントをモック
  - モック関数は `jest.fn()` で作成し、`mockResolvedValue` / `mockRejectedValue` で動作を制御
  - `beforeEach` で `jest.clearAllMocks()` を呼び出してモック状態をリセット
  - テストケースごとにモックの戻り値を変更可能
- **Implications**: この方式を失敗テストに適用することで、一貫性と保守性を確保

### 失敗テストで使用されているAPIメソッドの特定
- **Context**: 各失敗テストファイルで必要なAPIモックを特定
- **Sources Consulted**: 5件の失敗テストファイル、`lib/api.ts`
- **Findings**:
  | テストファイル | 使用APIメソッド |
  |---------------|----------------|
  | dm-email-send-page.test.tsx | `sendEmail` |
  | dm-posts-page.test.tsx | `getDmPosts`, `getDmUsers`, `createDmPost`, `deleteDmPost` |
  | dm-jobqueue-page.test.tsx | `registerJob` |
  | users-page.test.tsx | `getDmUsers`, `createDmUser`, `deleteDmUser` |
  | dm-user-posts-page.test.tsx | `getDmUserPosts` |
- **Implications**: 各テストファイルで必要なAPIメソッドのみをモックする

## Architecture Pattern Evaluation

| Option | Description | Strengths | Risks / Limitations | Notes |
|--------|-------------|-----------|---------------------|-------|
| MSW継続 + transformIgnorePatterns修正 | Jest設定の修正でMSWを使用し続ける | MSWの機能をフル活用可能 | 設定が複雑、他のテストに影響の可能性 | 既存の成功テストとの統一性が低い |
| jest.mock方式への統一 | 全テストを `jest.mock` 方式に変更 | シンプル、既存成功テストと一貫性あり | MSWの機能は使用不可 | **選択** |

## Design Decisions

### Decision: `jest.mock` 方式への統一
- **Context**: MSWを使用したテストがESM互換性問題で失敗している
- **Alternatives Considered**:
  1. Jest設定の修正（transformIgnorePatterns、moduleNameMapper等）
  2. MSWのバージョンダウングレード
  3. `jest.mock` 方式への変更
- **Selected Approach**: `jest.mock('@/lib/api', ...)` を使用したAPIモックに変更
- **Rationale**:
  - 成功しているテストと同じ方式で、一貫性を確保
  - Jest設定の変更は他のテストに影響を与える可能性がある
  - シンプルで保守しやすい
- **Trade-offs**:
  - MSWのネットワークレベルモック機能は使用不可
  - ただし、このプロジェクトではAPIクライアントレベルのモックで十分
- **Follow-up**: 修正後、全22件のテストが成功することを確認

## Risks & Mitigations
- **リスク1**: テストケースの漏れ — 元のテストと修正後のテストでテストケース数を比較して検証
- **リスク2**: モックデータの不一致 — 元のテストのモックデータ構造を維持
- **リスク3**: 既存テストへの影響 — 共通設定ファイルは変更しない方針

## References
- [Jest Mock Functions](https://jestjs.io/docs/mock-functions) — jest.mockの使用方法
- 既存成功テスト: `client/src/__tests__/integration/dm-feed-post-page.test.tsx`
