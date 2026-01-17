# 不足しているテストコードの作成の要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Feature名**: 0073-missing-test
- **作成日**: 2026-01-27
- **関連Issue**: https://github.com/taku-o/go-webdb-template/issues/151

### 1.2 目的
テストのない箇所にテストコードを用意する。特に、clientアプリのテストが不足している箇所に対して、画面が表示されているか確認する程度の簡易なテストでも良いので、テストコードを作成する。簡単にテストが作れるようならしっかりテストを作成する。

### 1.3 スコープ
- クライアント側（`client/`ディレクトリ）のNext.jsコード内のテストが不足している箇所の特定
- テストが不足しているページコンポーネントに対するテストコードの作成
- テストが不足しているコンポーネントに対するテストコードの作成
- テストが不足しているカスタムフックやユーティリティ関数に対するテストコードの作成
- サーバー側（`server/`ディレクトリ）のGoコード内のテストが不足している箇所の特定
- テストが不足しているハンドラー、サービス、リポジトリ、ユースケースなどに対するテストコードの作成

**本実装の範囲外**:
- 外部ライブラリのコード
- 既にテストが存在する箇所の再実装

## 2. 背景・現状分析

### 2.1 現状

#### 2.1.1 クライアント側（Next.js）
- Next.js 14+ (App Router)を使用したクライアントアプリケーションが存在
- 既存のテストファイル:
  - Jest統合テスト:
    - `client/src/__tests__/integration/users-page.test.tsx` (dm-users/page.tsx)
    - `client/src/__tests__/integration/dm-jobqueue-page.test.tsx` (dm-jobqueue/page.tsx)
  - Jestコンポーネントテスト:
    - `client/src/__tests__/components/TodayApiButton.test.tsx`
  - Jestライブラリテスト:
    - `client/src/__tests__/lib/api.test.ts`
  - Playwright E2Eテスト:
    - `client/e2e/auth-flow.spec.ts`
    - `client/e2e/cross-shard.spec.ts`
    - `client/e2e/csv-download.spec.ts`
    - `client/e2e/email-send.spec.ts`
    - `client/e2e/post-flow.spec.ts`
    - `client/e2e/user-flow.spec.ts`

#### 2.1.2 サーバー側（Go）
- Go 1.21+を使用したサーバーアプリケーションが存在
- 既存のテストファイル（主要なもの）:
  - Handlerテスト:
    - `server/internal/api/handler/today_handler_test.go`
    - `server/internal/api/handler/email_handler_test.go`
    - `server/internal/api/handler/dm_jobqueue_handler_test.go`
    - `server/internal/api/handler/upload_handler_test.go`
    - `server/internal/api/handler/dm_user_handler_huma_test.go`
    - `server/internal/api/handler/dm_post_handler_huma_test.go`
  - Usecaseテスト:
    - `server/internal/usecase/api/dm_user_usecase_test.go`
    - `server/internal/usecase/api/dm_post_usecase_test.go`
    - `server/internal/usecase/api/dm_jobqueue_usecase_test.go`
    - `server/internal/usecase/api/email_usecase_test.go`
    - `server/internal/usecase/api/today_usecase_test.go`
    - `server/internal/usecase/admin/dm_user_register_usecase_test.go`
    - `server/internal/usecase/admin/api_key_usecase_test.go`
    - `server/internal/usecase/cli/list_dm_users_usecase_test.go`
    - `server/internal/usecase/cli/generate_sample_usecase_test.go`
    - `server/internal/usecase/cli/generate_secret_usecase_test.go`
  - Serviceテスト:
    - `server/internal/service/api_key_service_test.go`
    - `server/internal/service/secret_service_test.go`
    - `server/internal/service/generate_sample_service_test.go`
    - `server/internal/service/date_service_test.go`
    - `server/internal/service/email/email_service_test.go`
    - `server/internal/service/jobqueue/client_test.go`
    - `server/internal/service/jobqueue/server_test.go`
    - `server/internal/service/jobqueue/processor_test.go`
  - Repositoryテスト:
    - `server/internal/repository/dm_user_repository_test.go`
    - `server/internal/repository/dm_post_repository_test.go`
    - `server/internal/repository/dm_news_repository_test.go`
  - その他のテスト:
    - `server/internal/auth/jwt_test.go`
    - `server/internal/auth/auth0_validator_test.go`
    - `server/internal/auth/middleware_test.go`
    - `server/internal/db/sharding_test.go`
    - `server/internal/db/group_manager_test.go`
    - `server/internal/config/config_test.go`
    - `server/test/integration/`配下の統合テスト
    - `server/test/e2e/`配下のE2Eテスト

### 2.2 問題点
- クライアント側:
  - 多くのページコンポーネントにテストが存在しない
  - 多くのコンポーネントにテストが存在しない
  - カスタムフックやユーティリティ関数の一部にテストが存在しない
- サーバー側:
  - 一部のハンドラーにテストが存在しない可能性がある（例: `dm_user_handler.go`, `dm_post_handler.go`の非Huma版）
  - 一部のサービスにテストが存在しない可能性がある（例: `dm_user_service.go`, `dm_post_service.go`）
  - その他のコンポーネントにテストが不足している可能性がある
- 全体的に:
  - テストカバレッジが不十分な可能性がある

### 2.3 必要性
- コードの品質保証とリグレッション防止
- リファクタリング時の安全性確保
- テストカバレッジの向上
- 開発効率の向上（テストによる動作確認の自動化）

### 2.4 テスト方針
- クライアント側:
  - clientアプリのテストの場合、テストの内容は画面が表示されているか確認する程度の簡易なものでも良い
  - 簡単にテストが作れるようならしっかりテストを作成する
  - 既存のテストパターンに従う（Jest + React Testing Library、Playwright）
- サーバー側:
  - Goの標準テストフレームワーク（`testing`パッケージ）と`testify`を使用
  - テーブル駆動テストパターンに従う
  - 既存のテストパターンに従う

## 3. 機能要件

### 3.1 テストが不足しているページコンポーネントの特定とテスト作成

#### 3.1.1 対象ファイル（テストが不足している可能性があるページ）
- `client/app/page.tsx`: トップページ
- `client/app/dm-posts/page.tsx`: 投稿一覧ページ
- `client/app/dm-user-posts/page.tsx`: ユーザー投稿一覧ページ
- `client/app/dm_email/send/page.tsx`: メール送信ページ
- `client/app/dm_feed/page.tsx`: フィードページ
- `client/app/dm_feed/[userId]/page.tsx`: ユーザー別フィードページ
- `client/app/dm_feed/[userId]/[postId]/page.tsx`: 投稿詳細ページ
- `client/app/dm_movie/upload/page.tsx`: 動画アップロードページ
- `client/app/dm_videoplayer/page.tsx`: 動画プレーヤーページ

#### 3.1.2 テスト要件
- 各ページが正常に表示されることを確認するテスト
- 必要に応じて、主要な機能（フォーム送信、データ表示など）のテスト
- 簡易なテスト（画面表示確認）でも良いが、可能であればしっかりしたテストを作成

### 3.2 テストが不足しているコンポーネントの特定とテスト作成

#### 3.2.1 対象ファイル（テストが不足している可能性があるコンポーネント）
- `client/components/feed/feed-form.tsx`: フィード投稿フォーム
- `client/components/feed/feed-post-card.tsx`: フィード投稿カード
- `client/components/feed/feed-reply-card.tsx`: フィード返信カード
- `client/components/feed/reply-form.tsx`: 返信フォーム
- `client/components/home/card.tsx`: ホームカード
- `client/components/home/component-grid.tsx`: コンポーネントグリッド
- `client/components/home/demo-modal.tsx`: デモモーダル
- `client/components/home/web-vitals.tsx`: Web Vitals
- `client/components/layout/footer.tsx`: フッター
- `client/components/layout/navbar-client.tsx`: ナビゲーションバー（クライアント）
- `client/components/layout/navbar.tsx`: ナビゲーションバー
- `client/components/shared/counting-numbers.tsx`: カウント数値表示
- `client/components/shared/error-alert.tsx`: エラーアラート
- `client/components/shared/loading-overlay.tsx`: ローディングオーバーレイ
- `client/components/shared/loading-spinner.tsx`: ローディングスピナー
- `client/components/shared/modal.tsx`: モーダル
- `client/components/shared/popover.tsx`: ポップオーバー
- `client/components/shared/tooltip.tsx`: ツールチップ
- `client/components/shared/icons/*.tsx`: アイコンコンポーネント群
- `client/components/video-player/video-player.tsx`: 動画プレーヤー

#### 3.2.2 テスト要件
- 各コンポーネントが正常にレンダリングされることを確認するテスト
- 必要に応じて、主要な機能（クリック、入力、表示切り替えなど）のテスト
- 簡易なテスト（画面表示確認）でも良いが、可能であればしっかりしたテストを作成

### 3.3 テストが不足しているカスタムフックやユーティリティ関数の特定とテスト作成

#### 3.3.1 対象ファイル（テストが不足している可能性があるカスタムフック・ユーティリティ）
- `client/lib/hooks/use-intersection-observer.ts`: Intersection Observerフック
- `client/lib/hooks/use-local-storage.ts`: localStorageフック
- `client/lib/hooks/use-media-query.ts`: メディアクエリフック
- `client/lib/hooks/use-scroll.ts`: スクロールフック
- `client/lib/utils.ts`: ユーティリティ関数

#### 3.3.2 テスト要件
- 各フックやユーティリティ関数が正常に動作することを確認するテスト
- 可能であれば、しっかりしたテストを作成

### 3.4 サーバー側（Go）のテストが不足している箇所の特定とテスト作成

#### 3.4.1 対象ファイル（テストが不足している可能性があるファイル）
- `server/internal/api/handler/dm_user_handler.go`: ユーザーハンドラー（Humaを使用、テストが不足している可能性）
- `server/internal/api/handler/dm_post_handler.go`: 投稿ハンドラー（Humaを使用、テストが不足している可能性）
- `server/internal/service/dm_user_service.go`: ユーザーサービス
- `server/internal/service/dm_post_service.go`: 投稿サービス
- その他、テストが不足している可能性があるファイル

#### 3.4.2 テスト要件
- 各ハンドラー、サービス、リポジトリ、ユースケースが正常に動作することを確認するテスト
- テーブル駆動テストパターンに従う
- 既存のテストパターンに従う

### 3.5 テストの種類と優先順位

#### 3.5.1 テストの種類
1. **クライアント側**:
   - Jest + React Testing Library: コンポーネントのユニットテスト、統合テスト
   - Playwright: E2Eテスト（既存のE2Eテストでカバーされている場合は不要）
2. **サーバー側**:
   - Go標準テスト（`testing`パッケージ）: ユニットテスト
   - `testify`: アサーション、モック
   - 統合テスト、E2Eテスト

#### 3.5.2 優先順位
1. **高優先度**: 
   - クライアント側: ページコンポーネント（ユーザーが直接アクセスするページ）
   - サーバー側: ハンドラー、サービス（主要なビジネスロジック）
2. **中優先度**: 
   - クライアント側: 主要なコンポーネント（フォーム、カード、モーダルなど）
   - サーバー側: リポジトリ、ユースケース
3. **低優先度**: 
   - クライアント側: ユーティリティコンポーネント（アイコン、ローディングスピナーなど）
   - サーバー側: ユーティリティ関数、ヘルパー関数

## 4. 非機能要件

### 4.1 テスト品質
- 既存のテストパターンに従う
- テストが失敗した場合、原因を特定しやすいエラーメッセージ
- テストの実行時間を考慮した実装

### 4.2 保守性
- テストコードの可読性を向上
- テストの構造を明確にする
- 既存のテスト構造に沿った実装

### 4.3 互換性
- 既存のテストが正常に動作することを維持
- 既存のテストフレームワークとの互換性を保つ

## 5. 制約事項

### 5.1 技術的制約
- クライアント側:
  - Next.js App Routerの特性を考慮
  - Server ComponentsとClient Componentsの使い分け
  - SSR（Server-Side Rendering）を考慮したテスト実装
- サーバー側:
  - Go 1.21+の制約
  - データベース接続の管理（テスト環境でのDB接続）
  - シャーディング対応のテスト実装

### 5.2 実装上の制約
- 既存のプロジェクト構造に沿って実装する
- 既存のテストパターンに従う
- 既存のテストが正常に動作することを維持

### 5.3 テストの制約
- 簡易なテスト（画面表示確認）でも良い
- 簡単にテストが作れるようならしっかりテストを作成する
- テストのない箇所を特定し、テストを作成する

## 6. 受け入れ基準

### 6.1 テストが不足している箇所の特定
- [ ] クライアント側: テストが不足しているページコンポーネントを特定した
- [ ] クライアント側: テストが不足しているコンポーネントを特定した
- [ ] クライアント側: テストが不足しているカスタムフックやユーティリティ関数を特定した
- [ ] サーバー側: テストが不足しているハンドラー、サービス、リポジトリ、ユースケースを特定した

### 6.2 ページコンポーネントのテスト作成
- [ ] `client/app/page.tsx`のテストが作成されている（または、テストが不要であることを確認）
- [ ] `client/app/dm-posts/page.tsx`のテストが作成されている（または、テストが不要であることを確認）
- [ ] `client/app/dm-user-posts/page.tsx`のテストが作成されている（または、テストが不要であることを確認）
- [ ] `client/app/dm_email/send/page.tsx`のテストが作成されている（または、テストが不要であることを確認）
- [ ] `client/app/dm_feed/page.tsx`のテストが作成されている（または、テストが不要であることを確認）
- [ ] `client/app/dm_feed/[userId]/page.tsx`のテストが作成されている（または、テストが不要であることを確認）
- [ ] `client/app/dm_feed/[userId]/[postId]/page.tsx`のテストが作成されている（または、テストが不要であることを確認）
- [ ] `client/app/dm_movie/upload/page.tsx`のテストが作成されている（または、テストが不要であることを確認）
- [ ] `client/app/dm_videoplayer/page.tsx`のテストが作成されている（または、テストが不要であることを確認）

### 6.3 コンポーネントのテスト作成
- [ ] `client/components/feed/`配下のコンポーネントのテストが作成されている（または、テストが不要であることを確認）
- [ ] `client/components/home/`配下のコンポーネントのテストが作成されている（または、テストが不要であることを確認）
- [ ] `client/components/layout/`配下のコンポーネントのテストが作成されている（または、テストが不要であることを確認）
- [ ] `client/components/shared/`配下の主要なコンポーネントのテストが作成されている（または、テストが不要であることを確認）
- [ ] `client/components/video-player/`配下のコンポーネントのテストが作成されている（または、テストが不要であることを確認）

### 6.4 カスタムフックやユーティリティ関数のテスト作成
- [ ] `client/lib/hooks/`配下のカスタムフックのテストが作成されている（または、テストが不要であることを確認）
- [ ] `client/lib/utils.ts`のテストが作成されている（または、テストが不要であることを確認）

### 6.5 サーバー側（Go）のテスト作成
- [ ] `server/internal/api/handler/dm_user_handler.go`のテストが作成されている（または、テストが不要であることを確認）
  - 注: このハンドラーはHumaを使用しており、`dm_user_handler_huma_test.go`が存在するが、より包括的なテストが必要な可能性がある
- [ ] `server/internal/api/handler/dm_post_handler.go`のテストが作成されている（または、テストが不要であることを確認）
  - 注: このハンドラーはHumaを使用しており、`dm_post_handler_huma_test.go`が存在するが、より包括的なテストが必要な可能性がある
- [ ] `server/internal/service/dm_user_service.go`のテストが作成されている（または、テストが不要であることを確認）
- [ ] `server/internal/service/dm_post_service.go`のテストが作成されている（または、テストが不要であることを確認）
- [ ] その他、特定したテストが不足している箇所のテストが作成されている（または、テストが不要であることを確認）

### 6.6 テストの実行
- [ ] 作成したテストが正常に実行される（クライアント側）
- [ ] 作成したテストが正常に実行される（サーバー側、`APP_ENV=test`を指定）
- [ ] 既存のテストが正常に動作する
- [ ] テストエラーが発生しない

### 6.7 テストの品質
- [ ] テストコードが既存のテストパターンに従っている
- [ ] テストコードが可読性が高い
- [ ] テストが適切な粒度で作成されている

## 7. 影響範囲

### 7.1 新規作成されるファイル
- クライアント側:
  - `client/src/__tests__/integration/`配下の新しいテストファイル
  - `client/src/__tests__/components/`配下の新しいテストファイル
  - `client/src/__tests__/lib/`配下の新しいテストファイル
- サーバー側:
  - `server/internal/api/handler/`配下の新しいテストファイル
  - `server/internal/service/`配下の新しいテストファイル
  - `server/internal/repository/`配下の新しいテストファイル
  - `server/internal/usecase/`配下の新しいテストファイル

### 7.2 既存ファイルへの影響
- 既存のテストファイルへの影響は最小限
- 既存のテストが正常に動作することを維持

### 7.3 既存機能への影響
- 既存の機能への影響はない（テストコードの追加のみ）

## 8. 実装上の注意事項

### 8.1 テストの作成方針
- 簡易なテスト（画面表示確認）でも良い
- 簡単にテストが作れるようならしっかりテストを作成する
- 既存のテストパターンに従う

### 8.2 テストフレームワークの使用
- クライアント側:
  - Jest + React Testing Library: コンポーネントのユニットテスト、統合テスト
  - MSW: APIモック（既存のテストで使用されているパターンに従う）
  - Playwright: E2Eテスト（既存のE2Eテストでカバーされている場合は不要）
- サーバー側:
  - Go標準テスト（`testing`パッケージ）: ユニットテスト
  - `testify/assert`: アサーション
  - `testify/mock`: モック（必要に応じて）
  - `net/http/httptest`: HTTPテスト
  - テーブル駆動テストパターンを使用

### 8.3 テストの構造
- 既存のテストファイルの構造に従う
- テストのグループ化（describeブロック）を適切に行う
- テストの名前を明確にする

### 8.4 テストデータの管理
- モックデータの管理方法を既存のテストパターンに従う
- MSWを使用したAPIモックの実装

### 8.5 テストの実行
- クライアント側: テストが正常に実行されることを確認
- サーバー側: `APP_ENV=test go test ./...`でテストが正常に実行されることを確認
- 既存のテストが正常に動作することを確認

## 9. 参考情報

### 9.1 関連ドキュメント
- Next.js App Routerドキュメント
- React Testing Libraryドキュメント
- Jestドキュメント
- Playwrightドキュメント
- 既存のプロジェクトドキュメント

### 9.2 関連Issue
- https://github.com/taku-o/go-webdb-template/issues/151: 本要件定義書の元となったIssue

### 9.3 技術スタック
- **クライアント側**:
  - フレームワーク: Next.js 14+ (App Router)
  - 言語: TypeScript 5+
  - テストフレームワーク: Jest、React Testing Library、Playwright
  - APIモック: MSW (Mock Service Worker)
- **サーバー側**:
  - 言語: Go 1.21+
  - テストフレームワーク: Go標準テスト（`testing`パッケージ）、`testify`
  - HTTPテスト: `net/http/httptest`

### 9.4 既存のテストパターン
- クライアント側:
  - Jest統合テスト: `client/src/__tests__/integration/users-page.test.tsx`を参考
  - Jestコンポーネントテスト: `client/src/__tests__/components/TodayApiButton.test.tsx`を参考
  - Jestライブラリテスト: `client/src/__tests__/lib/api.test.ts`を参考
  - Playwright E2Eテスト: `client/e2e/`配下のファイルを参考
- サーバー側:
  - Handlerテスト: `server/internal/api/handler/today_handler_test.go`を参考
  - Usecaseテスト: `server/internal/usecase/api/dm_user_usecase_test.go`を参考
  - Serviceテスト: `server/internal/service/api_key_service_test.go`を参考
  - Repositoryテスト: `server/internal/repository/dm_user_repository_test.go`を参考
  - 統合テスト: `server/test/integration/`配下のファイルを参考

### 9.5 実装の流れ
1. テストが不足している箇所を特定（クライアント側・サーバー側）
2. 各箇所に対してテストの必要性を判断
3. テストを作成（簡易なテストでも良いが、可能であればしっかりしたテストを作成）
   - クライアント側: Jest + React Testing Libraryを使用
   - サーバー側: Go標準テスト + testifyを使用（テーブル駆動テストパターン）
4. テストの実行と動作確認
   - クライアント側: `npm test`などで実行
   - サーバー側: `APP_ENV=test go test ./...`で実行
5. 既存のテストが正常に動作することを確認
