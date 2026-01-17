# Implementation Plan

## Acceptance Criteria (requirements.mdより)

### Requirement 1: MSW依存テストのモック方式統一
1. `npm run test:client` 実行時に全テストがエラーなく実行される
2. 失敗テストが成功テストと同じ `jest.mock` 方式を使用する
3. 修正後のテストが元のテストケースをカバーする
4. 修正後のテストが一貫したコーディングスタイルを維持する

### Requirement 2: テストの機能的等価性の維持
1. テストケースカバレッジ維持（ページタイトル、フォーム、API成功/エラー、ローディング状態）
2. ナビゲーションリンクテスト維持（該当する場合）
3. モックデータ構造維持

### Requirement 3: 既存成功テストへの影響回避
1. 既存の成功テスト（17件）が引き続き成功する
2. 共通設定ファイルを変更しない
3. 全22件のテストが成功する

## Tasks

- [ ] 1. テストファイルのモック方式変更（5件）

- [ ] 1.1 (P) メール送信ページテストのモック方式変更
  - MSWの `setupServer` と関連インポートを削除
  - `jest.mock('@/lib/api', ...)` でapiClientをモック
  - `mockSendEmail` 関数を定義し、各テストケースで `mockResolvedValue` / `mockRejectedValue` を設定
  - `beforeEach` で `jest.clearAllMocks()` を呼び出し
  - 6件のテストケース（ページタイトル、フォーム表示、戻るリンク、送信成功、APIエラー、ローディング状態）を維持
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 2.1_

- [ ] 1.2 (P) 投稿管理ページテストのモック方式変更
  - MSWの `setupServer` と関連インポートを削除
  - `jest.mock('@/lib/api', ...)` で `getDmPosts`, `getDmUsers`, `createDmPost`, `deleteDmPost` をモック
  - モック関数を定義し、`beforeEach` で初期値を設定
  - `server.use()` を使用していた箇所をモック関数の再設定に変更
  - 8件のテストケース（ページタイトル、投稿一覧、ローディング、戻るリンク、作成フォーム、投稿件数、APIエラー、空状態）を維持
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 2.1, 2.2_

- [ ] 1.3 (P) ジョブキューページテストのモック方式変更
  - MSWの `setupServer` と関連インポートを削除
  - `jest.mock('@/lib/api', ...)` で `registerJob` をモック
  - `mockRegisterJob` 関数を定義し、各テストケースで戻り値を設定
  - 4件のテストケース（フォーム表示、ジョブ登録成功、APIエラー、カスタムメッセージ入力）を維持
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 2.1_

- [ ] 1.4 (P) ユーザー管理ページテストのモック方式変更
  - MSWの `setupServer` と関連インポートを削除
  - `jest.mock('@/lib/api', ...)` で `getDmUsers`, `createDmUser`, `deleteDmUser` をモック
  - モック関数を定義し、`beforeEach` で初期値を設定
  - 動的なモックデータ更新をテストケース内での `mockResolvedValue` 再設定に変更
  - 4件のテストケース（ユーザー一覧表示、新規作成、APIエラー、ローディング状態）を維持
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 2.1_

- [ ] 1.5 (P) ユーザー投稿JOINページテストのモック方式変更
  - MSWの `setupServer` と関連インポートを削除
  - `jest.mock('@/lib/api', ...)` で `getDmUserPosts` をモック
  - `mockGetDmUserPosts` 関数を定義し、各テストケースで戻り値を設定
  - 8件のテストケース（ページタイトル、クロスシャード説明、一覧表示、ローディング、戻るリンク、シャーディング情報、APIエラー、空状態）を維持
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 2.1, 2.2_

- [ ] 2. 全体テスト実行と検証
  - `npm run test:client` を実行し、全22件のテストが成功することを確認
  - 修正した5件のテストが新たに成功することを確認
  - 既存の成功テスト17件が引き続き成功することを確認
  - テストケース数が修正前後で変わっていないことを確認（合計30テストケース）
  - _Requirements: 3.1, 3.2, 3.3_
