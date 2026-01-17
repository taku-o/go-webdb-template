# 不足しているテストコードの作成の実装タスク一覧

## 概要
テストのない箇所にテストコードを用意する実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: テストが不足している箇所の特定

#### - [ ] タスク 1.1: クライアント側のページコンポーネントのテスト不足箇所を特定
**目的**: クライアント側のページコンポーネントでテストが不足している箇所を特定する

**作業内容**:
- `client/app/`配下の各ページファイルを確認
- 対応するテストファイルが`client/src/__tests__/integration/`配下に存在するか確認
- テストが不足しているページのリストを作成
- 対象ページ:
  - `client/app/page.tsx`
  - `client/app/dm-posts/page.tsx`
  - `client/app/dm-user-posts/page.tsx`
  - `client/app/dm_email/send/page.tsx`
  - `client/app/dm_feed/page.tsx`
  - `client/app/dm_feed/[userId]/page.tsx`
  - `client/app/dm_feed/[userId]/[postId]/page.tsx`
  - `client/app/dm_movie/upload/page.tsx`
  - `client/app/dm_videoplayer/page.tsx`

**受け入れ基準**:
- テストが不足しているページコンポーネントのリストが作成されている
- 各ページに対して、テストが必要か不要かの判断が記録されている

_Requirements: 6.1, Design: Phase 1_

---

#### - [ ] タスク 1.2: クライアント側のコンポーネントのテスト不足箇所を特定
**目的**: クライアント側のコンポーネントでテストが不足している箇所を特定する

**作業内容**:
- `client/components/`配下の各コンポーネントファイルを確認
- 対応するテストファイルが`client/src/__tests__/components/`配下に存在するか確認
- テストが不足しているコンポーネントのリストを作成
- 対象ディレクトリ:
  - `client/components/feed/`
  - `client/components/home/`
  - `client/components/layout/`
  - `client/components/shared/`
  - `client/components/video-player/`

**受け入れ基準**:
- テストが不足しているコンポーネントのリストが作成されている
- 各コンポーネントに対して、テストが必要か不要かの判断が記録されている

_Requirements: 6.1, Design: Phase 1_

---

#### - [ ] タスク 1.3: クライアント側のカスタムフックとユーティリティ関数のテスト不足箇所を特定
**目的**: クライアント側のカスタムフックとユーティリティ関数でテストが不足している箇所を特定する

**作業内容**:
- `client/lib/hooks/`配下の各カスタムフックファイルを確認
- `client/lib/utils.ts`を確認
- 対応するテストファイルが存在するか確認
- テストが不足しているフックとユーティリティ関数のリストを作成
- 対象ファイル:
  - `client/lib/hooks/use-intersection-observer.ts`
  - `client/lib/hooks/use-local-storage.ts`
  - `client/lib/hooks/use-media-query.ts`
  - `client/lib/hooks/use-scroll.ts`
  - `client/lib/utils.ts`

**受け入れ基準**:
- テストが不足しているカスタムフックとユーティリティ関数のリストが作成されている
- 各フック・関数に対して、テストが必要か不要かの判断が記録されている

_Requirements: 6.1, Design: Phase 1_

---

#### - [ ] タスク 1.4: サーバー側のハンドラー、サービス、リポジトリ、ユースケースのテスト不足箇所を特定
**目的**: サーバー側のハンドラー、サービス、リポジトリ、ユースケースでテストが不足している箇所を特定する

**作業内容**:
- `server/internal/api/handler/`配下の各ハンドラーファイルを確認
- `server/internal/service/`配下の各サービスファイルを確認
- `server/internal/repository/`配下の各リポジトリファイルを確認
- `server/internal/usecase/`配下の各ユースケースファイルを確認
- 対応するテストファイルが存在するか確認
- テストが不足しているファイルのリストを作成
- 特に確認すべきファイル:
  - `server/internal/api/handler/dm_user_handler.go`
  - `server/internal/api/handler/dm_post_handler.go`
  - `server/internal/service/dm_user_service.go`
  - `server/internal/service/dm_post_service.go`

**受け入れ基準**:
- テストが不足しているハンドラー、サービス、リポジトリ、ユースケースのリストが作成されている
- 各ファイルに対して、テストが必要か不要かの判断が記録されている

_Requirements: 6.1, Design: Phase 1_

---

### Phase 2: 高優先度のテスト作成（クライアント側: ページコンポーネント）

#### - [ ] タスク 2.1: `client/app/page.tsx`のテスト作成
**目的**: トップページのテストを作成する（簡易なテストでも良い）

**作業内容**:
- `client/src/__tests__/integration/page-page.test.tsx`を作成
- MSWを使用してAPIモックを設定
- ページが正常に表示されることを確認するテストを実装
- 必要に応じて、主要な機能のテストを追加
- **useEffectを極力使わない**（`waitFor`や`findBy*`クエリを使用）
- テストを実行して動作確認

**受け入れ基準**:
- テストファイルが作成されている
- ページが正常に表示されることを確認するテストが実装されている
- テストが正常に実行される
- 既存のテストが正常に動作する

_Requirements: 6.2, Design: 1.1, Phase 2_

---

#### - [ ] タスク 2.2: `client/app/dm-posts/page.tsx`のテスト作成
**目的**: 投稿一覧ページのテストを作成する（簡易なテストでも良い）

**作業内容**:
- `client/src/__tests__/integration/dm-posts-page.test.tsx`を作成
- MSWを使用してAPIモックを設定
- ページが正常に表示されることを確認するテストを実装
- 必要に応じて、主要な機能（投稿表示、フォーム送信など）のテストを追加
- **useEffectを極力使わない**（`waitFor`や`findBy*`クエリを使用）
- テストを実行して動作確認

**受け入れ基準**:
- テストファイルが作成されている
- ページが正常に表示されることを確認するテストが実装されている
- テストが正常に実行される
- 既存のテストが正常に動作する

_Requirements: 6.2, Design: 1.1, Phase 2_

---

#### - [ ] タスク 2.3: `client/app/dm-user-posts/page.tsx`のテスト作成
**目的**: ユーザー投稿一覧ページのテストを作成する（簡易なテストでも良い）

**作業内容**:
- `client/src/__tests__/integration/dm-user-posts-page.test.tsx`を作成
- MSWを使用してAPIモックを設定
- ページが正常に表示されることを確認するテストを実装
- 必要に応じて、主要な機能のテストを追加
- **useEffectを極力使わない**（`waitFor`や`findBy*`クエリを使用）
- テストを実行して動作確認

**受け入れ基準**:
- テストファイルが作成されている
- ページが正常に表示されることを確認するテストが実装されている
- テストが正常に実行される
- 既存のテストが正常に動作する

_Requirements: 6.2, Design: 1.1, Phase 2_

---

#### - [ ] タスク 2.4: `client/app/dm_email/send/page.tsx`のテスト作成
**目的**: メール送信ページのテストを作成する（簡易なテストでも良い）

**作業内容**:
- `client/src/__tests__/integration/dm-email-send-page.test.tsx`を作成
- MSWを使用してAPIモックを設定
- ページが正常に表示されることを確認するテストを実装
- 必要に応じて、主要な機能（メール送信フォームなど）のテストを追加
- **useEffectを極力使わない**（`waitFor`や`findBy*`クエリを使用）
- テストを実行して動作確認

**受け入れ基準**:
- テストファイルが作成されている
- ページが正常に表示されることを確認するテストが実装されている
- テストが正常に実行される
- 既存のテストが正常に動作する

_Requirements: 6.2, Design: 1.1, Phase 2_

---

#### - [ ] タスク 2.5: `client/app/dm_feed/page.tsx`のテスト作成
**目的**: フィードページのテストを作成する（簡易なテストでも良い）

**作業内容**:
- `client/src/__tests__/integration/dm-feed-page.test.tsx`を作成
- MSWを使用してAPIモックを設定
- ページが正常に表示されることを確認するテストを実装
- 必要に応じて、主要な機能のテストを追加
- **useEffectを極力使わない**（`waitFor`や`findBy*`クエリを使用）
- テストを実行して動作確認

**受け入れ基準**:
- テストファイルが作成されている
- ページが正常に表示されることを確認するテストが実装されている
- テストが正常に実行される
- 既存のテストが正常に動作する

_Requirements: 6.2, Design: 1.1, Phase 2_

---

#### - [ ] タスク 2.6: `client/app/dm_feed/[userId]/page.tsx`のテスト作成
**目的**: ユーザー別フィードページのテストを作成する（簡易なテストでも良い）

**作業内容**:
- `client/src/__tests__/integration/dm-feed-user-page.test.tsx`を作成
- MSWを使用してAPIモックを設定
- 動的ルートパラメータ（userId）を考慮したテストを実装
- ページが正常に表示されることを確認するテストを実装
- 必要に応じて、主要な機能のテストを追加
- **useEffectを極力使わない**（`waitFor`や`findBy*`クエリを使用）
- テストを実行して動作確認

**受け入れ基準**:
- テストファイルが作成されている
- ページが正常に表示されることを確認するテストが実装されている
- テストが正常に実行される
- 既存のテストが正常に動作する

_Requirements: 6.2, Design: 1.1, Phase 2_

---

#### - [ ] タスク 2.7: `client/app/dm_feed/[userId]/[postId]/page.tsx`のテスト作成
**目的**: 投稿詳細ページのテストを作成する（簡易なテストでも良い）

**作業内容**:
- `client/src/__tests__/integration/dm-feed-post-page.test.tsx`を作成
- MSWを使用してAPIモックを設定
- 動的ルートパラメータ（userId、postId）を考慮したテストを実装
- ページが正常に表示されることを確認するテストを実装
- 必要に応じて、主要な機能のテストを追加
- **useEffectを極力使わない**（`waitFor`や`findBy*`クエリを使用）
- テストを実行して動作確認

**受け入れ基準**:
- テストファイルが作成されている
- ページが正常に表示されることを確認するテストが実装されている
- テストが正常に実行される
- 既存のテストが正常に動作する

_Requirements: 6.2, Design: 1.1, Phase 2_

---

#### - [ ] タスク 2.8: `client/app/dm_movie/upload/page.tsx`のテスト作成
**目的**: 動画アップロードページのテストを作成する（簡易なテストでも良い）

**作業内容**:
- `client/src/__tests__/integration/dm-movie-upload-page.test.tsx`を作成
- MSWを使用してAPIモックを設定
- ページが正常に表示されることを確認するテストを実装
- 必要に応じて、主要な機能（アップロードフォームなど）のテストを追加
- **useEffectを極力使わない**（`waitFor`や`findBy*`クエリを使用）
- テストを実行して動作確認

**受け入れ基準**:
- テストファイルが作成されている
- ページが正常に表示されることを確認するテストが実装されている
- テストが正常に実行される
- 既存のテストが正常に動作する

_Requirements: 6.2, Design: 1.1, Phase 2_

---

#### - [ ] タスク 2.9: `client/app/dm_videoplayer/page.tsx`のテスト作成
**目的**: 動画プレーヤーページのテストを作成する（簡易なテストでも良い）

**作業内容**:
- `client/src/__tests__/integration/dm-videoplayer-page.test.tsx`を作成
- MSWを使用してAPIモックを設定
- ページが正常に表示されることを確認するテストを実装
- 必要に応じて、主要な機能のテストを追加
- **useEffectを極力使わない**（`waitFor`や`findBy*`クエリを使用）
- テストを実行して動作確認

**受け入れ基準**:
- テストファイルが作成されている
- ページが正常に表示されることを確認するテストが実装されている
- テストが正常に実行される
- 既存のテストが正常に動作する

_Requirements: 6.2, Design: 1.1, Phase 2_

---

### Phase 3: 高優先度のテスト作成（サーバー側: ハンドラー、サービス）

#### - [ ] タスク 3.1: `server/internal/api/handler/dm_user_handler.go`のテスト作成
**目的**: ユーザーハンドラーの包括的なテストを作成する

**作業内容**:
- `server/internal/api/handler/dm_user_handler_test.go`を作成（または既存のテストを拡張）
- Huma APIを使用したテストを実装
- テーブル駆動テストパターンを使用
- 各エンドポイント（作成、取得、一覧、更新、削除、CSVダウンロード）のテストを実装
- モックを使用してUsecaseをモック化
- テストを実行して動作確認（`APP_ENV=test go test ./...`）

**受け入れ基準**:
- テストファイルが作成されている（または既存のテストが拡張されている）
- 各エンドポイントのテストが実装されている
- テーブル駆動テストパターンに従っている
- テストが正常に実行される（`APP_ENV=test`を指定）
- 既存のテストが正常に動作する

_Requirements: 6.5, Design: 2.1, Phase 2_

---

#### - [ ] タスク 3.2: `server/internal/api/handler/dm_post_handler.go`のテスト作成
**目的**: 投稿ハンドラーの包括的なテストを作成する

**作業内容**:
- `server/internal/api/handler/dm_post_handler_test.go`を作成（または既存のテストを拡張）
- Huma APIを使用したテストを実装
- テーブル駆動テストパターンを使用
- 各エンドポイント（作成、取得、一覧、更新、削除、ユーザー投稿JOIN）のテストを実装
- モックを使用してUsecaseをモック化
- テストを実行して動作確認（`APP_ENV=test go test ./...`）

**受け入れ基準**:
- テストファイルが作成されている（または既存のテストが拡張されている）
- 各エンドポイントのテストが実装されている
- テーブル駆動テストパターンに従っている
- テストが正常に実行される（`APP_ENV=test`を指定）
- 既存のテストが正常に動作する

_Requirements: 6.5, Design: 2.1, Phase 2_

---

#### - [ ] タスク 3.3: `server/internal/service/dm_user_service.go`のテスト作成
**目的**: ユーザーサービスのテストを作成する

**作業内容**:
- `server/internal/service/dm_user_service_test.go`を作成
- テーブル駆動テストパターンを使用
- 各メソッド（CreateUser、GetUser、ListUsers、UpdateUser、DeleteUserなど）のテストを実装
- モックを使用してRepositoryをモック化
- テストを実行して動作確認（`APP_ENV=test go test ./...`）

**受け入れ基準**:
- テストファイルが作成されている
- 各メソッドのテストが実装されている
- テーブル駆動テストパターンに従っている
- テストが正常に実行される（`APP_ENV=test`を指定）
- 既存のテストが正常に動作する

_Requirements: 6.5, Design: 2.2, Phase 2_

---

#### - [ ] タスク 3.4: `server/internal/service/dm_post_service.go`のテスト作成
**目的**: 投稿サービスのテストを作成する

**作業内容**:
- `server/internal/service/dm_post_service_test.go`を作成
- テーブル駆動テストパターンを使用
- 各メソッド（CreatePost、GetPost、ListPosts、UpdatePost、DeletePostなど）のテストを実装
- モックを使用してRepositoryをモック化
- テストを実行して動作確認（`APP_ENV=test go test ./...`）

**受け入れ基準**:
- テストファイルが作成されている
- 各メソッドのテストが実装されている
- テーブル駆動テストパターンに従っている
- テストが正常に実行される（`APP_ENV=test`を指定）
- 既存のテストが正常に動作する

_Requirements: 6.5, Design: 2.2, Phase 2_

---

### Phase 4: 中優先度のテスト作成（クライアント側: 主要なコンポーネント）

#### - [ ] タスク 4.1: `client/components/feed/`配下のコンポーネントのテスト作成
**目的**: フィード関連コンポーネントのテストを作成する（簡易なテストでも良い）

**作業内容**:
- `client/src/__tests__/components/feed-form.test.tsx`を作成
- `client/src/__tests__/components/feed-post-card.test.tsx`を作成
- `client/src/__tests__/components/feed-reply-card.test.tsx`を作成
- `client/src/__tests__/components/reply-form.test.tsx`を作成
- 各コンポーネントが正常にレンダリングされることを確認するテストを実装
- 必要に応じて、主要な機能（フォーム送信、クリックなど）のテストを追加
- **useEffectを極力使わない**（`waitFor`や`findBy*`クエリを使用）
- テストを実行して動作確認

**受け入れ基準**:
- 各コンポーネントのテストファイルが作成されている
- 各コンポーネントが正常にレンダリングされることを確認するテストが実装されている
- テストが正常に実行される
- 既存のテストが正常に動作する

_Requirements: 6.3, Design: 1.2, Phase 3_

---

#### - [ ] タスク 4.2: `client/components/home/`配下のコンポーネントのテスト作成
**目的**: ホーム関連コンポーネントのテストを作成する（簡易なテストでも良い）

**作業内容**:
- `client/src/__tests__/components/home-card.test.tsx`を作成
- `client/src/__tests__/components/component-grid.test.tsx`を作成
- `client/src/__tests__/components/demo-modal.test.tsx`を作成
- `client/src/__tests__/components/web-vitals.test.tsx`を作成（必要に応じて）
- 各コンポーネントが正常にレンダリングされることを確認するテストを実装
- 必要に応じて、主要な機能のテストを追加
- **useEffectを極力使わない**（`waitFor`や`findBy*`クエリを使用）
- テストを実行して動作確認

**受け入れ基準**:
- 各コンポーネントのテストファイルが作成されている
- 各コンポーネントが正常にレンダリングされることを確認するテストが実装されている
- テストが正常に実行される
- 既存のテストが正常に動作する

_Requirements: 6.3, Design: 1.2, Phase 3_

---

#### - [ ] タスク 4.3: `client/components/layout/`配下のコンポーネントのテスト作成
**目的**: レイアウト関連コンポーネントのテストを作成する（簡易なテストでも良い）

**作業内容**:
- `client/src/__tests__/components/footer.test.tsx`を作成
- `client/src/__tests__/components/navbar-client.test.tsx`を作成
- `client/src/__tests__/components/navbar.test.tsx`を作成
- 各コンポーネントが正常にレンダリングされることを確認するテストを実装
- 必要に応じて、主要な機能のテストを追加
- **useEffectを極力使わない**（`waitFor`や`findBy*`クエリを使用）
- テストを実行して動作確認

**受け入れ基準**:
- 各コンポーネントのテストファイルが作成されている
- 各コンポーネントが正常にレンダリングされることを確認するテストが実装されている
- テストが正常に実行される
- 既存のテストが正常に動作する

_Requirements: 6.3, Design: 1.2, Phase 3_

---

#### - [ ] タスク 4.4: `client/components/shared/`配下の主要なコンポーネントのテスト作成
**目的**: 共有コンポーネントのテストを作成する（簡易なテストでも良い）

**作業内容**:
- `client/src/__tests__/components/counting-numbers.test.tsx`を作成
- `client/src/__tests__/components/error-alert.test.tsx`を作成
- `client/src/__tests__/components/loading-overlay.test.tsx`を作成
- `client/src/__tests__/components/loading-spinner.test.tsx`を作成
- `client/src/__tests__/components/modal.test.tsx`を作成
- `client/src/__tests__/components/popover.test.tsx`を作成
- `client/src/__tests__/components/tooltip.test.tsx`を作成
- 各コンポーネントが正常にレンダリングされることを確認するテストを実装
- 必要に応じて、主要な機能のテストを追加
- **useEffectを極力使わない**（`waitFor`や`findBy*`クエリを使用）
- テストを実行して動作確認

**受け入れ基準**:
- 各コンポーネントのテストファイルが作成されている
- 各コンポーネントが正常にレンダリングされることを確認するテストが実装されている
- テストが正常に実行される
- 既存のテストが正常に動作する

_Requirements: 6.3, Design: 1.2, Phase 3_

---

#### - [ ] タスク 4.5: `client/components/video-player/`配下のコンポーネントのテスト作成
**目的**: 動画プレーヤーコンポーネントのテストを作成する（簡易なテストでも良い）

**作業内容**:
- `client/src/__tests__/components/video-player.test.tsx`を作成
- コンポーネントが正常にレンダリングされることを確認するテストを実装
- 必要に応じて、主要な機能のテストを追加
- **useEffectを極力使わない**（`waitFor`や`findBy*`クエリを使用）
- テストを実行して動作確認

**受け入れ基準**:
- テストファイルが作成されている
- コンポーネントが正常にレンダリングされることを確認するテストが実装されている
- テストが正常に実行される
- 既存のテストが正常に動作する

_Requirements: 6.3, Design: 1.2, Phase 3_

---

### Phase 5: 中優先度のテスト作成（サーバー側: リポジトリ、ユースケース）

#### - [ ] タスク 5.1: サーバー側のリポジトリのテスト不足箇所を確認し、テストを作成
**目的**: リポジトリでテストが不足している箇所のテストを作成する

**作業内容**:
- Phase 1で特定したリポジトリのテスト不足箇所を確認
- 各リポジトリに対して、テーブル駆動テストパターンを使用したテストを作成
- テスト用データベースを使用（`test/testutil`のヘルパー関数を使用）
- 各メソッド（Create、Get、List、Update、Deleteなど）のテストを実装
- テストを実行して動作確認（`APP_ENV=test go test ./...`）

**受け入れ基準**:
- テストが不足しているリポジトリのテストが作成されている
- テーブル駆動テストパターンに従っている
- テストが正常に実行される（`APP_ENV=test`を指定）
- 既存のテストが正常に動作する

_Requirements: 6.5, Design: 2.3, Phase 3_

---

#### - [ ] タスク 5.2: サーバー側のユースケースのテスト不足箇所を確認し、テストを作成
**目的**: ユースケースでテストが不足している箇所のテストを作成する

**作業内容**:
- Phase 1で特定したユースケースのテスト不足箇所を確認
- 各ユースケースに対して、テーブル駆動テストパターンを使用したテストを作成
- モックを使用してServiceをモック化
- 各メソッドのテストを実装
- テストを実行して動作確認（`APP_ENV=test go test ./...`）

**受け入れ基準**:
- テストが不足しているユースケースのテストが作成されている
- テーブル駆動テストパターンに従っている
- テストが正常に実行される（`APP_ENV=test`を指定）
- 既存のテストが正常に動作する

_Requirements: 6.5, Design: 2.4, Phase 3_

---

### Phase 6: 低優先度のテスト作成（クライアント側: カスタムフック、ユーティリティ関数）

#### - [ ] タスク 6.1: `client/lib/hooks/`配下のカスタムフックのテスト作成
**目的**: カスタムフックのテストを作成する

**作業内容**:
- `client/src/__tests__/lib/hooks/use-intersection-observer.test.ts`を作成
- `client/src/__tests__/lib/hooks/use-local-storage.test.ts`を作成
- `client/src/__tests__/lib/hooks/use-media-query.test.ts`を作成
- `client/src/__tests__/lib/hooks/use-scroll.test.ts`を作成
- `renderHook`を使用して各フックのテストを実装
- 各フックが正常に動作することを確認するテストを実装
- **useEffectを極力使わない**
- テストを実行して動作確認

**受け入れ基準**:
- 各カスタムフックのテストファイルが作成されている
- 各フックが正常に動作することを確認するテストが実装されている
- テストが正常に実行される
- 既存のテストが正常に動作する

_Requirements: 6.4, Design: 1.3, Phase 4_

---

#### - [ ] タスク 6.2: `client/lib/utils.ts`のテスト作成
**目的**: ユーティリティ関数のテストを作成する

**作業内容**:
- `client/src/__tests__/lib/utils.test.ts`を作成
- 各ユーティリティ関数のテストを実装
- 各関数が正常に動作することを確認するテストを実装
- テストを実行して動作確認

**受け入れ基準**:
- テストファイルが作成されている
- 各ユーティリティ関数のテストが実装されている
- テストが正常に実行される
- 既存のテストが正常に動作する

_Requirements: 6.4, Design: 1.4, Phase 4_

---

### Phase 7: テストの実行と確認

#### - [ ] タスク 7.1: クライアント側のテストの実行と確認
**目的**: 作成したクライアント側のテストが正常に実行されることを確認する

**作業内容**:
- `npm test`でクライアント側のテストを実行
- 作成したテストが正常に実行されることを確認
- 既存のテストが正常に動作することを確認
- テストエラーが発生しないことを確認
- 必要に応じて、テストを修正

**受け入れ基準**:
- 作成したテストが正常に実行される
- 既存のテストが正常に動作する
- テストエラーが発生しない

_Requirements: 6.6, Design: Phase 5_

---

#### - [ ] タスク 7.2: サーバー側のテストの実行と確認
**目的**: 作成したサーバー側のテストが正常に実行されることを確認する

**作業内容**:
- `APP_ENV=test go test ./...`でサーバー側のテストを実行
- 作成したテストが正常に実行されることを確認
- 既存のテストが正常に動作することを確認
- テストエラーが発生しないことを確認
- 必要に応じて、テストを修正

**受け入れ基準**:
- 作成したテストが正常に実行される（`APP_ENV=test`を指定）
- 既存のテストが正常に動作する
- テストエラーが発生しない

_Requirements: 6.6, Design: Phase 5_

---

#### - [ ] タスク 7.3: テストの品質確認
**目的**: 作成したテストの品質を確認する

**作業内容**:
- テストコードが既存のテストパターンに従っていることを確認
- テストコードの可読性を確認
- テストが適切な粒度で作成されていることを確認
- 必要に応じて、テストを改善

**受け入れ基準**:
- テストコードが既存のテストパターンに従っている
- テストコードが可読性が高い
- テストが適切な粒度で作成されている

_Requirements: 6.7, Design: Phase 5_

---

## 注意事項

### クライアント側のテスト
- **useEffectを極力使わない**: テストコード内でもuseEffectの使用を避け、`waitFor`や`findBy*`クエリを使用して非同期処理を待つ
- MSWを使用してAPIモックを実装
- 既存のテストパターンに従う

### サーバー側のテスト
- **`APP_ENV=test`を必ず指定**: テスト実行時に`APP_ENV=test`を指定しないと認証エラーが発生する
- テーブル駆動テストパターンを使用
- 既存のテストパターンに従う

### テストの優先順位
1. **高優先度**: ページコンポーネント（クライアント側）、ハンドラー・サービス（サーバー側）
2. **中優先度**: 主要なコンポーネント（クライアント側）、リポジトリ・ユースケース（サーバー側）
3. **低優先度**: ユーティリティコンポーネント・フック（クライアント側）、ユーティリティ関数（サーバー側）

### テストの作成方針
- 簡易なテスト（画面表示確認）でも良い
- 簡単にテストが作れるようならしっかりテストを作成する
- テストが不要な箇所は、その旨を記録する
