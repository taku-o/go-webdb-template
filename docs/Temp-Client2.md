# Client2 App

新clientアプリケーション（precedentテンプレートベース）

既存の`client`アプリケーションの機能を、NextAuth (Auth.js) v5とshadcn/uiを活用したモダンな実装に移植したアプリケーションです。

## セットアップ

1. 依存関係のインストール
   ```bash
   cd client2
   npm install --legacy-peer-deps
   ```

   **注意**: peer dependencyの競合がある場合は`--legacy-peer-deps`フラグを使用してください。

2. 環境変数の設定
   
   **AUTH_SECRETの生成**:
   ```bash
   # プロジェクトルートで実行
   npm run cli:generate-secret
   ```
   このコマンドで生成された秘密鍵をコピーします。
   
   `.env.local`を作成して以下の環境変数を設定：
   ```
   # NextAuth (Auth.js)
   AUTH_SECRET=<npm run cli:generate-secretで生成した秘密鍵>
   AUTH_URL=http://localhost:3000

   # Auth0設定
   AUTH0_ISSUER=https://your-tenant.auth0.com
   AUTH0_CLIENT_ID=your-client-id
   AUTH0_CLIENT_SECRET=your-client-secret
   AUTH0_AUDIENCE=https://your-api-audience

   # API設定
   NEXT_PUBLIC_API_KEY=your-api-key
   NEXT_PUBLIC_API_BASE_URL=http://localhost:8080

   # テスト環境用（テスト実行時に必要）
   APP_ENV=test
   ```

   **注意**: 
   - `AUTH_SECRET`は`npm run cli:generate-secret`コマンドで生成します（`server/cmd/generate-secret`を使用）。
   - `APP_ENV=test`はテスト実行時に必要です（`npm test`、`npm run e2e`実行時）。

3. Auth0アプリケーション設定
   Auth0ダッシュボード（`Applications > [対象アプリ] > Settings`）で以下のURLを設定：

   **Allowed Callback URLs:**
   ```
   http://localhost:3000/api/auth/callback/auth0
   ```

   **Allowed Logout URLs:**
   ```
   http://localhost:3000
   ```

   **Allowed Web Origins:**
   ```
   http://localhost:3000
   ```

4. 開発サーバーの起動
   ```bash
   npm run dev
   ```
   
   開発サーバーは`http://localhost:3000`で起動します。

## 技術スタック

- **フレームワーク**: Next.js 14+ (App Router)
- **言語**: TypeScript 5+
- **UIコンポーネント**: shadcn/ui
- **認証**: NextAuth (Auth.js) v5
- **スタイリング**: Tailwind CSS
- **フォーム管理**: react-hook-form
- **バリデーション**: zod
- **ファイルアップロード**: Uppy (TUSプロトコル)
- **テスト**: Playwright (E2E), Jest (単体・統合), MSW (APIモック)

## 実装されている機能

### 認証機能
- NextAuth (Auth.js) v5によるログイン/ログアウト
- Auth0プロバイダーとの統合
- 認証トークンの取得（`/api/auth/token`）
- プロフィール取得（`/api/auth/profile`）
- セッション管理（サーバーサイド・クライアントサイド対応）

### ページ機能

#### トップページ (`/`)
- 機能一覧の表示（shadcn/uiの`card`コンポーネント）
- 認証状態の表示
- TodayApiButtonコンポーネント（プライベートAPIテスト用）

#### ユーザー管理ページ (`/dm-users`)
- ユーザー一覧の表示
- ユーザーの作成・編集・削除（CRUD）
- CSVダウンロード機能

#### 投稿管理ページ (`/dm-posts`)
- 投稿一覧の表示
- 投稿の作成・編集・削除（CRUD）
- ユーザー選択機能（ドロップダウン）

#### ユーザーと投稿のJOINページ (`/dm-user-posts`)
- ユーザーと投稿をJOINして表示
- クロスシャードクエリ対応
- 空状態の表示

#### メール送信ページ (`/dm_email/send`)
- ウェルカムメールの送信
- 送信結果の表示

#### 動画アップロードページ (`/dm_movie/upload`)
- 動画ファイルのアップロード（TUSプロトコル）
- Uppyダッシュボードを使用したアップロードUI
- アップロード進捗の表示

#### ジョブキューページ (`/dm-jobqueue`)
- 遅延ジョブの登録（参考実装）
- ジョブ登録結果の表示

### コンポーネント

#### TodayApiButton
- プライベートAPIテスト用のコンポーネント
- 認証状態に応じた表示制御

#### レイアウトコンポーネント
- `Navbar`: ナビゲーションバー（認証状態表示、ログイン/ログアウトボタン）
- `Footer`: フッター（現在は使用していない）

### APIクライアント

`lib/api.ts`に実装されているAPIクライアント機能：

- ユーザー管理API（`getDmUsers`, `createDmUser`, `updateDmUser`, `deleteDmUser`, `downloadDmUsersCSV`）
- 投稿管理API（`getDmPosts`, `createDmPost`, `updateDmPost`, `deleteDmPost`）
- クロスシャードクエリAPI（`getDmUserPosts`）
- メール送信API（`sendWelcomeEmail`）
- ジョブキューAPI（`registerJob`）
- 動画アップロードAPI（`createMovieUploader` - Uppyインスタンス作成）
- プライベートAPI（`getToday`）

すべてのAPI呼び出しは、NextAuthの認証トークンを自動的に含めます。

### 型定義

- `types/dm_post.ts`: 投稿関連の型定義
- `types/dm_user.ts`: ユーザー関連の型定義
- `types/jobqueue.ts`: ジョブキュー関連の型定義
- `types/next-auth.d.ts`: NextAuthのセッション型拡張

## 利用可能なスクリプト

### 開発
- `npm run dev` - 開発サーバーを起動（ポート3000）
- `npm run build` - プロダクションビルドを実行
- `npm run start` - プロダクションビルドを起動（ポート3000）
- `npm run lint` - ESLintを実行
- `npm run type-check` - TypeScript型チェックを実行
- `npm run format` - Prettierでフォーマットを確認
- `npm run format:write` - Prettierでフォーマットを適用

### テスト
- `npm test` - Jestテストを実行（単体・統合テスト）
- `npm run test:watch` - Jestテストをウォッチモードで実行
- `npm run test:coverage` - Jestテストのカバレッジを取得
- `npm run e2e` - Playwright E2Eテストを実行
- `npm run e2e:ui` - Playwright E2EテストをUIモードで実行
- `npm run e2e:headed` - Playwright E2Eテストをヘッドモードで実行

**注意**: テスト実行時は`APP_ENV=test`が自動的に設定されます（`package.json`のスクリプトに含まれています）。

## ディレクトリ構造

```
client2/
├── app/
│   ├── api/
│   │   └── auth/
│   │       ├── [...nextauth]/route.ts    # NextAuth認証ルート
│   │       ├── profile/route.ts           # プロフィール取得API
│   │       └── token/route.ts             # トークン取得API
│   ├── dm_email/send/page.tsx             # メール送信ページ
│   ├── dm_movie/upload/page.tsx           # 動画アップロードページ
│   ├── dm-jobqueue/page.tsx               # ジョブキューページ
│   ├── dm-posts/page.tsx                  # 投稿管理ページ
│   ├── dm-user-posts/page.tsx            # ユーザーと投稿のJOINページ
│   ├── dm-users/page.tsx                 # ユーザー管理ページ
│   ├── layout.tsx                         # ルートレイアウト
│   ├── page.tsx                           # トップページ
│   └── globals.css                        # グローバルスタイル
├── components/
│   ├── ui/                                # shadcn/uiコンポーネント
│   ├── layout/                            # レイアウトコンポーネント
│   └── TodayApiButton.tsx                # TodayApiButtonコンポーネント
├── lib/
│   ├── api.ts                             # APIクライアント
│   └── auth.ts                            # 認証ヘルパー
├── types/                                 # 型定義
│   ├── dm_post.ts
│   ├── dm_user.ts
│   ├── jobqueue.ts
│   └── next-auth.d.ts
├── e2e/                                   # E2Eテスト
│   ├── auth-flow.spec.ts
│   ├── user-flow.spec.ts
│   ├── post-flow.spec.ts
│   ├── cross-shard.spec.ts
│   ├── email-send.spec.ts
│   └── csv-download.spec.ts
├── src/__tests__/                        # 単体・統合テスト
│   ├── integration/
│   ├── components/
│   └── lib/
├── auth.ts                                # NextAuth設定
├── jest.config.js                         # Jest設定
├── jest.setup.js                          # Jestセットアップ
├── jest.polyfills.js                      # Jestポリフィル
├── playwright.config.ts                   # Playwright設定
└── package.json
```

## 認証フロー

1. ユーザーがログインボタンをクリック
2. NextAuthがAuth0にリダイレクト
3. Auth0で認証後、`/api/auth/callback/auth0`にリダイレクト
4. NextAuthがセッションを作成し、アクセストークンをJWTに保存
5. セッション情報がクライアントに返される
6. API呼び出し時は`lib/auth.ts`の`getAuthToken()`でトークンを取得し、Authorizationヘッダーに含める

## テスト

### E2Eテスト（Playwright）
- 認証フロー（`auth-flow.spec.ts`）
- ユーザー管理フロー（`user-flow.spec.ts`）
- 投稿管理フロー（`post-flow.spec.ts`）
- クロスシャードクエリフロー（`cross-shard.spec.ts`）
- メール送信フロー（`email-send.spec.ts`）
- CSVダウンロードフロー（`csv-download.spec.ts`）

### 統合テスト（Jest + MSW）
- ユーザー管理ページ（`users-page.test.tsx`）
- ジョブキューページ（`dm-jobqueue-page.test.tsx`）

### 単体テスト（Jest）
- TodayApiButtonコンポーネント（`TodayApiButton.test.tsx`）
- APIクライアント（`api.test.ts`）

## 環境変数の詳細

### NextAuth (Auth.js)
- `AUTH_SECRET`: セッション暗号化用の秘密鍵（必須）
- `AUTH_URL`: アプリケーションのベースURL（開発環境では`http://localhost:3000`）

### Auth0
- `AUTH0_ISSUER`: Auth0テナントのURL（例: `https://your-tenant.auth0.com`）
- `AUTH0_CLIENT_ID`: Auth0アプリケーションのクライアントID
- `AUTH0_CLIENT_SECRET`: Auth0アプリケーションのクライアントシークレット
- `AUTH0_AUDIENCE`: APIのオーディエンス（例: `https://your-api-audience`）

### API
- `NEXT_PUBLIC_API_KEY`: バックエンドAPIのAPIキー
- `NEXT_PUBLIC_API_BASE_URL`: バックエンドAPIのベースURL（例: `http://localhost:8080`）

### テスト環境
- `APP_ENV`: テスト実行時に`test`を設定（`package.json`のスクリプトで自動設定）

## 注意事項

- このドキュメントは一時的なもので、`client`から`client2`への移行が完了したら、この内容をREADMEに移植する想定です。
- テスト実行時は`APP_ENV=test`が設定されていることを確認してください（認証エラーが発生する場合があります）。
- 開発サーバーとAPIサーバーが両方起動している必要があります（E2Eテスト実行時）。
