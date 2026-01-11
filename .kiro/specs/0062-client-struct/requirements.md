# クライアントアプリケーション設計改善の要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Feature名**: 0062-client-struct
- **作成日**: 2026-01-27

### 1.2 目的
クライアントアプリケーション（Next.js）の設計を改善し、認証処理とAPI呼び出しの責務を適切に分離する。これにより、コードの保守性、テスタビリティ、再利用性を向上させ、コンポーネントとビジネスロジックの関心を分離する。

### 1.3 スコープ
- 認証トークン取得ロジックを共通化（`client/src/lib/auth.ts`の作成）
- API Clientの改善（`client/src/lib/api.ts`の修正）
- コンポーネントの簡素化（`client/src/components/TodayApiButton.tsx`の修正）
- Auth0 SDKの標準的な方法を積極的に利用
- テストコードの更新

**基本方針**:
- 使えそうなライブラリ、フレームワークがあるなら積極的に利用する
- 後方互換性は不要（既存コードの変更も可）

**本実装の範囲外**:
- 既存のAPIエンドポイントの変更（サーバー側の変更は不要）
- UI/UXの変更（機能的な変更のみ）

## 2. 背景・現状分析

### 2.1 現在の状況

#### 2.1.1 TodayApiButton.tsxの現状
- **実装場所**: `client/src/components/TodayApiButton.tsx`に直接実装されている
- **認証処理**: コンポーネント内で直接`/auth/token`を呼び出している
- **API呼び出し**: コンポーネント内で直接`fetch`を呼び出している
- **問題点**:
  - 認証ロジックがコンポーネント内に散在している
  - API呼び出しの詳細がコンポーネントに露出している
  - 認証トークンの取得方法が統一されていない

#### 2.1.2 api.tsの現状
- **実装場所**: `client/src/lib/api.ts`に実装されている
- **認証方式**: API Keyを前提としている
- **問題点**:
  - `getToday`メソッドはJWTを引数で受け取る形になっている
  - 認証トークンの取得ロジックが共通化されていない
  - 認証方式（API Key vs JWT）の切り替えが不自然

#### 2.1.3 他のコンポーネントの現状
- **dm-users/page.tsx**: `apiClient`を使用しており、適切な設計
- **dm-posts/page.tsx**: `apiClient`を使用しており、適切な設計
- **dm-user-posts/page.tsx**: `apiClient`を使用しており、適切な設計
- **dm-jobqueue/page.tsx**: `apiClient`を使用しており、適切な設計
- **dm_email/send/page.tsx**: `apiClient`を使用しており、適切な設計

### 2.2 現在の処理内容

#### 2.2.1 TodayApiButton.tsxの処理内容
1. **認証トークンの取得**: 
   - ログイン中: `/auth/token`を呼び出してAuth0 JWTを取得
   - 未ログイン: `process.env.NEXT_PUBLIC_API_KEY`を使用
2. **API呼び出し**: 
   - `fetch`を直接呼び出して`/api/today`エンドポイントにリクエスト
   - Authorizationヘッダーにトークンを付与
3. **エラーハンドリング**: コンポーネント内でエラーを処理

#### 2.2.2 api.tsの処理内容
1. **API Keyの取得**: コンストラクタで`process.env.NEXT_PUBLIC_API_KEY`を取得
2. **リクエスト処理**: `request`メソッドでAPI呼び出しを実行
3. **認証**: API KeyをAuthorizationヘッダーに付与
4. **getTodayメソッド**: JWTを引数で受け取り、`request`メソッドに渡す

### 2.3 課題点
1. **認証処理の分散**: 認証トークン取得ロジックがコンポーネント内に直接実装されている
2. **API呼び出しの分散**: コンポーネント内で直接`fetch`を呼び出している
3. **責務の混在**: コンポーネントが認証処理とAPI呼び出しの詳細を知る必要がある
4. **再利用性の低さ**: 認証処理がコンポーネント内に直接実装されており、他のコンポーネントで再利用できない
5. **テストの困難さ**: 認証処理とAPI呼び出しがコンポーネント内に直接実装されているため、単体テストが困難
6. **コードの重複**: 認証トークン取得ロジックが複数箇所に散在する可能性がある

### 2.4 本実装による改善点
1. **関心の分離**: コンポーネント、API Client、Auth Serviceの責務を明確化
2. **コードの再利用性**: 認証処理を他のコンポーネントでも再利用可能
3. **テスタビリティ**: 各モジュールを独立してテスト可能
4. **保守性**: 認証ロジックの変更が一箇所で済む
5. **一貫性**: API呼び出しパターンが統一される

## 3. 機能要件

### 3.1 認証サービスの共通化

#### 3.1.1 Auth0 SDKの標準エンドポイントの活用
- **目的**: Auth0 SDKの標準的な方法を積極的に利用してコードを簡素化
- **実装内容**:
  - `@auth0/nextjs-auth0`は既にインストール済みであることを確認
  - `handleAuth`を使用した標準的な認証エンドポイント（`/api/auth/[auth0]/route.ts`）を作成
  - これにより、`/api/auth/login`、`/api/auth/logout`などの標準エンドポイントが自動的に作成される
  - 既存の`/auth/token`エンドポイントは必要に応じて維持または改善（アクセストークン取得用）
  - Auth0 SDKの標準的な方法を優先的に利用

#### 3.1.2 auth.tsの作成
- **目的**: 認証トークン取得ロジックを一箇所に集約
- **実装内容**:
  - `client/src/lib/auth.ts`を新規作成
  - `getAuthToken(auth0user: Auth0User | undefined): Promise<string>`関数を実装
  - ログイン中: `/auth/token`を呼び出してAuth0 JWTを取得（既存のエンドポイントを使用）
  - 未ログイン: `process.env.NEXT_PUBLIC_API_KEY`を使用
  - エラーハンドリングを実装
  - **注意**: Auth0 SDKの標準的な方法を優先的に利用し、コードを簡素化する

#### 3.1.3 getAuthToken関数の仕様
- **引数**: `auth0user: Auth0User | undefined`（Auth0のユーザー情報）
- **戻り値**: `Promise<string>`（認証トークン）
- **処理内容**:
  - `auth0user`が存在する場合: `/auth/token`を呼び出してAuth0 JWTを取得（既存のエンドポイントを使用）
  - `auth0user`が存在しない場合: `process.env.NEXT_PUBLIC_API_KEY`を返す
  - エラー発生時: 適切なエラーメッセージを投げる

### 3.2 API Clientの改善

#### 3.2.1 api.tsの修正
- **目的**: 認証トークンの取得を`ApiClient`内部で処理するように変更
- **実装内容**:
  - `request`メソッドに`auth0user`パラメータを追加
  - `getAuthToken`関数をインポート
  - `request`メソッド内で`getAuthToken`を呼び出してトークンを取得
  - 既存のメソッドも必要に応じて修正（後方互換性は不要）

#### 3.2.2 requestメソッドの修正
- **目的**: 認証トークンの取得を内部で処理
- **実装内容**:
  - `request`メソッドのシグネチャを変更: `request<T>(endpoint: string, options?: RequestInit, auth0user?: Auth0User | undefined): Promise<T>`
  - `getAuthToken(auth0user)`を呼び出してトークンを取得
  - 取得したトークンをAuthorizationヘッダーに付与
  - 既存のエラーハンドリングを維持

#### 3.2.3 getTodayメソッドの修正
- **目的**: `auth0user`パラメータを受け取るように変更
- **実装内容**:
  - `getToday`メソッドのシグネチャを変更: `getToday(auth0user?: Auth0User | undefined): Promise<{ date: string }>`
  - `request`メソッドに`auth0user`パラメータを渡す
  - 既存の動作を維持

#### 3.2.4 既存メソッドの改善
- **目的**: 既存のAPI呼び出しメソッドも統一されたパターンに変更
- **実装内容**:
  - 既存のメソッド（`getDmUsers`、`getDmPosts`など）も必要に応じて修正（メソッド名を本来あるべき名前に変更）
  - `auth0user`パラメータを追加し、統一された認証処理を使用
  - `auth0user`パラメータが未指定の場合、`getAuthToken(undefined)`が呼び出され、API Keyが使用される
  - 後方互換性は不要（既存コードの変更も可）

### 3.3 コンポーネントの簡素化

#### 3.3.1 TodayApiButton.tsxの修正
- **目的**: コンポーネントから認証処理とAPI呼び出しの詳細を削除
- **実装内容**:
  - 認証処理（`/auth/token`の呼び出し）を削除
  - API呼び出し（`fetch`の直接呼び出し）を削除
  - `apiClient.getToday(auth0user || undefined)`を呼び出すように変更（`useUser()`から取得した変数名を`auth0user`に変更）
  - UI表示と状態管理のみに集中
  - 既存のエラーハンドリングを維持

#### 3.3.2 コンポーネントの責務
- **UI表示**: コンポーネントのレンダリング
- **状態管理**: `useState`を使用した状態管理
- **イベントハンドリング**: ボタンクリックなどのイベント処理
- **API呼び出し**: `apiClient`を使用したAPI呼び出し（詳細は隠蔽）

### 3.4 テストの更新

#### 3.4.1 api.test.tsの更新
- **目的**: 新しい認証処理に対応したテストを追加
- **実装内容**:
  - `getAuthToken`関数のモックを追加
  - `getToday`メソッドのテストを更新
  - `auth0user`パラメータが渡された場合のテストを追加

#### 3.4.2 TodayApiButtonのテスト更新
- **目的**: 新しい実装に対応したテストを更新
- **実装内容**:
  - `apiClient.getToday`のモックを追加
  - 認証処理のテストを削除（認証処理は`auth.ts`に移動）
  - API呼び出しのテストを更新

## 4. 非機能要件

### 4.1 パフォーマンス
- **既存機能の維持**: 既存のAPI呼び出しのパフォーマンスを維持
- **オーバーヘッド**: 認証処理の共通化によるパフォーマンスオーバーヘッドは無視できるレベル
- **トークン取得**: `/auth/token`の呼び出し回数は既存と同等（各API呼び出しごとに1回）

### 4.2 信頼性
- **エラーハンドリング**: 適切なエラーハンドリングを実装
- **認証エラー**: 認証トークン取得失敗時のエラーハンドリングを適切に実装
- **ライブラリの利用**: 標準的なライブラリ・フレームワークを積極的に利用し、信頼性を向上

### 4.3 保守性
- **コードの可読性**: 各モジュールの責務が明確で、コードが読みやすい
- **一貫性**: 他のコンポーネントと同じAPI呼び出しパターンを採用
- **テスト容易性**: 各モジュールを独立してテスト可能
- **再利用性**: 認証処理を他のコンポーネントでも再利用可能

### 4.4 互換性
- **TypeScript型**: TypeScript型定義を適切に実装
- **ライブラリの利用**: 標準的なライブラリ・フレームワークを積極的に利用

## 5. 制約事項

### 5.1 技術的制約
- **既存のAPIエンドポイント**: 既存のAPIエンドポイント（`/api/today`など）を変更しない
- **認証方式**: Auth0 JWTとAPI Keyの両方に対応する必要がある
- **Next.js App Router**: Next.js 14+のApp Routerを使用
- **TypeScript**: TypeScript 5+を使用

### 5.2 実装上の制約
- **ディレクトリ構造**: 既存のディレクトリ構造に従う（`client/src/lib/`、`client/src/components/`）
- **命名規則**: 既存の命名規則に従う（`getAuthToken`、`apiClient`など）
- **ライブラリの利用**: 使えそうなライブラリ、フレームワークがあるなら積極的に利用する
- **後方互換性**: 不要（既存コードの変更も可）

### 5.3 動作環境
- **ローカル環境**: ローカル環境でクライアントアプリが正常に動作することを確認
- **CI環境**: CI環境でもクライアントアプリが正常に動作することを確認（該当する場合）
- **ブラウザ**: モダンブラウザ（Chrome、Firefox、Safari、Edge）で動作することを前提

## 6. 受け入れ基準

### 6.1 認証サービスの共通化
- [ ] `@auth0/nextjs-auth0`がインストールされていることを確認
- [ ] `handleAuth`を使用した標準的な認証エンドポイント（`/api/auth/[auth0]/route.ts`）が作成されている
- [ ] `client/src/lib/auth.ts`が作成されている
- [ ] `getAuthToken(auth0user: Auth0User | undefined): Promise<string>`関数が実装されている
- [ ] ログイン中に`/auth/token`を呼び出してAuth0 JWTを取得する処理が実装されている（既存のエンドポイントを使用）
- [ ] 未ログイン時に`process.env.NEXT_PUBLIC_API_KEY`を使用する処理が実装されている
- [ ] エラーハンドリングが適切に実装されている

### 6.2 API Clientの改善
- [ ] `client/src/lib/api.ts`の`request`メソッドが`auth0user`パラメータを受け取るように修正されている
- [ ] `request`メソッド内で`getAuthToken`を呼び出してトークンを取得する処理が実装されている
- [ ] `getToday`メソッドが`auth0user`パラメータを受け取るように修正されている
- [ ] 既存のメソッド（`getDmUsers`、`getDmPosts`など）も必要に応じて修正されている（メソッド名を本来あるべき名前に変更）
- [ ] `auth0user`パラメータが未指定の場合、API Keyが使用される

### 6.3 コンポーネントの簡素化
- [ ] `client/src/components/TodayApiButton.tsx`から認証処理（`/auth/token`の呼び出し）が削除されている
- [ ] `client/src/components/TodayApiButton.tsx`からAPI呼び出し（`fetch`の直接呼び出し）が削除されている
- [ ] `apiClient.getToday(auth0user || undefined)`を呼び出すように変更されている（`useUser()`から取得した変数名を`auth0user`に変更）
- [ ] UI表示と状態管理のみに集中している
- [ ] 既存のエラーハンドリングが維持されている

### 6.4 動作確認
- [ ] ローカル環境でクライアントアプリが正常に動作する
- [ ] `TodayApiButton`コンポーネントが正常に動作する
- [ ] ログイン中にAuth0 JWTが使用される
- [ ] 未ログイン時にAPI Keyが使用される
- [ ] 全てのAPI呼び出し（`dm-users`、`dm-posts`など）が正常に動作する
- [ ] 全てのテストが通過する

### 6.5 テスト
- [ ] `client/src/lib/__tests__/api.test.ts`が更新されている
- [ ] `getAuthToken`関数のテストが実装されている（該当する場合）
- [ ] `getToday`メソッドのテストが更新されている
- [ ] `TodayApiButton`のテストが更新されている（該当する場合）
- [ ] 既存のテストが全て通過する

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### 新規作成が必要なファイル
- `client/src/lib/auth.ts`: 認証トークン取得ロジックを集約
- `client/src/app/api/auth/[auth0]/route.ts`: Auth0 SDKの標準的な認証エンドポイント（`handleAuth`を使用）

#### 修正が必要なファイル
- `client/src/lib/api.ts`: `request`メソッドと`getToday`メソッドを修正
- `client/src/components/TodayApiButton.tsx`: 認証処理とAPI呼び出しを削除し、`apiClient`を使用するように変更
- `client/src/lib/__tests__/api.test.ts`: 新しい認証処理に対応したテストを追加

### 7.2 既存機能への影響
- **既存のAPI呼び出し**: 必要に応じて修正（後方互換性は不要）
- **既存のコンポーネント**: 必要に応じて修正（後方互換性は不要）
- **認証方式**: Auth0 JWTとAPI Keyの両方に対応
- **ライブラリの利用**: 標準的なライブラリ・フレームワークを積極的に利用

## 8. 実装上の注意事項

### 8.1 認証サービスの実装
- **Auth0 SDKの積極的利用**: `@auth0/nextjs-auth0`の標準的な方法を積極的に利用し、コードを簡素化する
- **標準エンドポイント**: `handleAuth`を使用した標準的な認証エンドポイント（`/api/auth/[auth0]/route.ts`）を作成
- **既存エンドポイント**: 既存の`/auth/token`エンドポイントは必要に応じて維持または改善（アクセストークン取得用）
- **ライブラリの利用**: 使えそうなライブラリ、フレームワークがあるなら積極的に利用する
- **エラーハンドリング**: `/auth/token`の呼び出し失敗時は適切なエラーメッセージを投げる
- **環境変数**: `process.env.NEXT_PUBLIC_API_KEY`が設定されていない場合のエラーハンドリング
- **型定義**: `Auth0User`型は`@auth0/nextjs-auth0`から`User`としてインポートし、`Auth0User`としてエイリアス定義する

### 8.2 API Clientの実装
- **既存メソッドの改善**: 既存のメソッド（`getDmUsers`、`getDmPosts`など）も必要に応じて修正（メソッド名を本来あるべき名前に変更、後方互換性は不要）
- **`auth0user`パラメータ**: オプショナルとし、未指定時は`undefined`を`getAuthToken`に渡す
- **エラーハンドリング**: 適切なエラーハンドリングを実装
- **ライブラリの利用**: 標準的なライブラリ・フレームワークを積極的に利用

### 8.3 コンポーネントの実装
- **認証処理の削除**: `/auth/token`の呼び出しを完全に削除
- **API呼び出しの削除**: `fetch`の直接呼び出しを完全に削除
- **`apiClient`の使用**: `apiClient.getToday(auth0user || undefined)`を使用
- **エラーハンドリング**: 既存のエラーハンドリングを維持

### 8.4 テストの実装
- **モック**: `getAuthToken`関数と`apiClient`のモックを使用
- **テストケース**: ログイン中と未ログイン時の両方のケースをテスト
- **テストの更新**: 必要に応じて既存のテストも更新（後方互換性は不要）
- **ライブラリの利用**: テスト用のライブラリも積極的に利用

## 9. 参考情報

### 9.1 関連ドキュメント
- `.kiro/specs/0062-client-struct/client-architecture-proposal.md`: 設計提案ドキュメント
- `.kiro/steering/structure.md`: ファイル組織とコードパターン
- `.kiro/steering/tech.md`: 技術スタックとアーキテクチャ

### 9.2 既存実装の参考
- `client/src/lib/api.ts`: 既存のAPI Clientの実装パターン
- `client/src/app/dm-users/page.tsx`: 既存のコンポーネントの実装パターン（適切な設計）
- `client/src/app/dm-posts/page.tsx`: 既存のコンポーネントの実装パターン（適切な設計）

### 9.3 技術スタック
- **フレームワーク**: Next.js 14+ (App Router)
- **言語**: TypeScript 5+
- **認証**: Auth0 Next.js SDK (`@auth0/nextjs-auth0`)
- **テスト**: Jest、React Testing Library

### 9.4 Auth0 SDKの標準的な方法
- **`handleAuth`**: `app/api/auth/[auth0]/route.ts`を作成し、`handleAuth()`をエクスポートすることで、標準的な認証エンドポイント（`/api/auth/login`、`/api/auth/logout`など）が自動的に作成される。積極的に利用する。
- **`useUser`**: クライアントコンポーネントでユーザー情報を取得するためのフック（既に使用中）
- **`getAccessToken`**: サーバーコンポーネントでアクセストークンを取得するための関数（既存の`/auth/token`エンドポイントで使用中）
- **ライブラリの積極的利用**: 使えそうなライブラリ、フレームワークがあるなら積極的に利用する

### 9.5 アーキテクチャの比較

| 項目 | 現在（修正前） | 修正後 |
|------|---------------|--------|
| 認証処理 | コンポーネント内に直接実装 | `auth.ts`に集約 |
| API呼び出し | コンポーネント内で直接`fetch` | `apiClient`を使用 |
| 責務 | コンポーネントが認証とAPI呼び出しの詳細を知る | コンポーネントはUI表示と状態管理のみ |
| 再利用性 | 認証処理がコンポーネント内に直接実装 | 認証処理を他のコンポーネントでも再利用可能 |
| テスタビリティ | 認証処理とAPI呼び出しがコンポーネント内に直接実装 | 各モジュールを独立してテスト可能 |

### 9.5 設計原則

1. **関心の分離 (Separation of Concerns)**
   - コンポーネント: UI表示と状態管理のみ
   - API Client: API呼び出しの抽象化
   - Auth Service: 認証トークンの取得ロジック

2. **単一責任の原則 (Single Responsibility Principle)**
   - 各モジュールは一つの責任のみを持つ
   - 認証処理は`auth.ts`に集約
   - API呼び出しは`api.ts`に集約

3. **DRY原則 (Don't Repeat Yourself)**
   - 認証トークン取得ロジックを一箇所に集約
   - API呼び出しパターンを統一

4. **依存性の逆転 (Dependency Inversion)**
   - コンポーネントは`apiClient`に依存
   - `apiClient`は`auth`に依存
   - 実装の詳細は隠蔽される
