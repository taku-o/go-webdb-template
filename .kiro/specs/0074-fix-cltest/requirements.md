# Requirements Document

## Introduction

`npm run test:client` コマンドでクライアント側のテストを実行した際に、5件のテストファイルで `SyntaxError: Unexpected token 'export'` エラーが発生している。このエラーは、MSW（Mock Service Worker）ライブラリの依存パッケージ `until-async` がESM形式でエクスポートしているため、Jestがそれを正しくトランスパイルできないことが原因である。

### 失敗しているテストファイル
1. `src/__tests__/integration/dm-email-send-page.test.tsx`
2. `src/__tests__/integration/dm-posts-page.test.tsx`
3. `src/__tests__/integration/dm-jobqueue-page.test.tsx`
4. `src/__tests__/integration/users-page.test.tsx`
5. `src/__tests__/integration/dm-user-posts-page.test.tsx`

### 成功しているテストファイルとの違い
- 失敗しているテスト: `msw/node` の `setupServer` を使用してAPIをモック
- 成功しているテスト: `jest.mock('@/lib/api', ...)` でAPIクライアントを直接モック

### エラー内容
```
/Users/taku-o/Documents/workspaces/go-webdb-template/client/node_modules/.pnpm/until-async@3.0.2/node_modules/until-async/lib/index.js:23
export { until };
^^^^^^

SyntaxError: Unexpected token 'export'
```

### 制約事項
- テストを無視したり、抑制したりする対応は不可
- テストを削除して無効化する対応は不可
- 既存のテストとの実装の統一性を配慮する必要がある

## Requirements

### Requirement 1: MSW依存テストのモック方式統一

**Objective:** As a 開発者, I want MSWを使用しているテストファイルのモック方式を統一したい, so that テストが正常に実行され、保守性が向上する

#### Acceptance Criteria
1. When `npm run test:client` を実行した時, the テストシステム shall 全てのテストファイルがエラーなく実行される
2. When 失敗しているテストファイルを修正した時, the 修正されたテストファイル shall 成功しているテストファイル（例: `dm-feed-post-page.test.tsx`）と同じモック方式（`jest.mock`）を使用する
3. The 修正後のテスト shall 元のテストと同等のテストケースをカバーする（テスト項目を削減しない）
4. The 修正後のテスト shall 他の成功しているテストファイルと一貫したコーディングスタイルを維持する

### Requirement 2: テストの機能的等価性の維持

**Objective:** As a 開発者, I want モック方式の変更後も同等のテスト品質を維持したい, so that テストの網羅性が損なわれない

#### Acceptance Criteria
1. When テストを修正した時, the 修正後のテスト shall 以下のテストケースを引き続きカバーする:
   - ページタイトルの表示確認
   - フォームの表示確認
   - API成功時の動作確認
   - APIエラー時のエラーハンドリング確認
   - ローディング状態の確認（該当する場合）
2. If 元のテストにナビゲーションリンクのテストが含まれている場合, then the 修正後のテスト shall 同等のナビゲーションリンクテストを含む
3. The 修正後のテスト shall モックデータの構造を元のテストと同等に維持する

### Requirement 3: 既存成功テストへの影響回避

**Objective:** As a 開発者, I want 修正が既存の成功しているテストに影響を与えないことを確認したい, so that リグレッションを防止できる

#### Acceptance Criteria
1. When テスト修正を行った後, the テストシステム shall 既存の成功しているテスト（17件）が引き続き成功する
2. The 修正 shall `jest.config.js` や `jest.setup.js` などの共通設定ファイルを変更しない（または変更が他のテストに悪影響を与えない）
3. When `npm run test:client` を実行した時, the テストシステム shall 全22件のテストが成功する
