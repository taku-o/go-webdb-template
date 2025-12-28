# Auth0 API呼び出し機能実装タスク一覧

## 概要
Auth0 API呼び出し機能の実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: サーバー側 - ライブラリの追加と設定

#### - [ ] タスク 1.1: Auth0 JWT検証ライブラリの追加
**目的**: `github.com/MicahParks/keyfunc`ライブラリを追加

**作業内容**:
- `server/go.mod`に`github.com/MicahParks/keyfunc`を追加
- `go get github.com/MicahParks/keyfunc`を実行
- ライブラリのバージョンを確認

**受け入れ基準**:
- `server/go.mod`に`github.com/MicahParks/keyfunc`が追加されている
- `go get`が正常に完了している
- `server/go.sum`が更新されている

---

#### - [ ] タスク 1.2: 設定構造体の拡張
**目的**: `AUTH0_ISSUER_BASE_URL`設定を追加

**作業内容**:
- `server/internal/config/config.go`の`APIConfig`構造体に`Auth0IssuerBaseURL`フィールドを追加
- `mapstructure:"auth0_issuer_base_url"`タグを追加

**受け入れ基準**:
- `APIConfig`構造体に`Auth0IssuerBaseURL string`フィールドが追加されている
- `mapstructure:"auth0_issuer_base_url"`タグが設定されている

---

#### - [ ] タスク 1.3: 設定ファイルへの追加
**目的**: 環境別設定ファイルに`AUTH0_ISSUER_BASE_URL`を追加

**作業内容**:
- `config/develop/config.yaml`に`auth0_issuer_base_url: "https://dev-oaa5vtzmld4dsxtd.jp.auth0.com"`を追加
- `config/staging/config.yaml`に`auth0_issuer_base_url`を追加（設定値は環境に応じて決定）
- `config/production/config.yaml`に`auth0_issuer_base_url`を追加（設定値は環境に応じて決定）

**受け入れ基準**:
- すべての環境別設定ファイルに`auth0_issuer_base_url`が追加されている
- develop環境の設定値が正しく設定されている

---

### Phase 2: サーバー側 - JWT種類の判別機能

#### - [ ] タスク 2.1: JWT種類の判別機能の実装
**目的**: JWTの種類（Auth0 JWT / Public API Key JWT）を判別する機能を実装

**作業内容**:
- `server/internal/auth/jwt.go`に`JWTType`型を追加
- `JWTTypeAuth0`、`JWTTypePublicAPIKey`、`JWTTypeUnknown`定数を定義
- `DetectJWTType()`関数を実装
  - 署名検証なしでJWTをパース（`ParseUnverified`を使用）
  - `iss`（Issuer）による判別
  - Auth0ドメインパターンのチェック

**受け入れ基準**:
- `JWTType`型が定義されている
- `DetectJWTType()`関数が実装されている
- Auth0 JWTとPublic API Key JWTを正しく判別できる
- 未知のJWTタイプの場合、`JWTTypeUnknown`を返す

---

### Phase 3: サーバー側 - Auth0 JWT検証機能

#### - [ ] タスク 3.1: Auth0Validator構造体の実装
**目的**: Auth0 JWT検証機能を提供する構造体を実装

**作業内容**:
- `server/internal/auth/auth0_validator.go`を作成
- `Auth0Validator`構造体を定義
- `NewAuth0Validator()`関数を実装
  - JWKS URLの構築（`{AUTH0_ISSUER_BASE_URL}/.well-known/jwks.json`）
  - `keyfunc.Options`の設定（12時間の定期更新、5分間隔の再取得制限、10秒のタイムアウト、未知のKIDが来たら再取得）
  - `keyfunc.Get()`を使用してJWKSを取得・キャッシュ
- `ValidateAuth0JWT()`メソッドを実装
  - `jwt.Parse()`と`jwks.Keyfunc`を使用してJWTを検証
- `Close()`メソッドを実装（リソースの解放）

**受け入れ基準**:
- `server/internal/auth/auth0_validator.go`が作成されている
- `Auth0Validator`構造体が定義されている
- `NewAuth0Validator()`関数が実装されている
- `ValidateAuth0JWT()`メソッドが実装されている
- `Close()`メソッドが実装されている
- JWKSキャッシュが正しく動作する

---

### Phase 4: サーバー側 - 認証ミドルウェアの拡張

#### - [ ] タスク 4.1: 認証ミドルウェアの拡張
**目的**: JWT種類の判別とAuth0 JWT検証をミドルウェアに統合

**作業内容**:
- `server/internal/auth/middleware.go`の`NewHumaAuthMiddleware()`関数を修正
- `auth0IssuerBaseURL`パラメータを追加
- `Auth0Validator`の初期化処理を追加
- JWT種類の判別処理を追加（`DetectJWTType()`を使用）
- JWT種類に応じた検証処理を追加
  - Auth0 JWTの場合: `auth0Validator.ValidateAuth0JWT()`を使用
  - Public API Key JWTの場合: 既存の`validator.ValidateJWT()`を使用
- JWTの許容する公開レベルを判定
  - Auth0 JWT → `"private"`（publicとprivateの両方にアクセス可能）
  - Public API Key JWT → `"public"`（publicなAPIのみアクセス可能）
- 許容する公開レベルをコンテキストに設定（`huma.ContextWithValue()`を使用）

**受け入れ基準**:
- `NewHumaAuthMiddleware()`関数に`auth0IssuerBaseURL`パラメータが追加されている
- `Auth0Validator`の初期化処理が実装されている
- JWT種類の判別処理が実装されている
- JWT種類に応じた検証処理が実装されている
- JWTの許容する公開レベルが正しく判定される
- 許容する公開レベルがコンテキストに設定される

---

#### - [ ] タスク 4.2: ルーターでのミドルウェア設定の更新
**目的**: ミドルウェアの呼び出し時に`AUTH0_ISSUER_BASE_URL`を渡す

**作業内容**:
- `server/internal/api/router/router.go`を修正
- `NewHumaAuthMiddleware()`の呼び出し時に`cfg.API.Auth0IssuerBaseURL`を渡す
- 設定が空の場合のエラーハンドリングを追加
  - `AUTH0_ISSUER_BASE_URL`が空の場合、エラーを投げてサーバーを起動しない
  - サービスの方針として、必要な設定が不十分な時はエラーを投げてしまって良い

**受け入れ基準**:
- `NewHumaAuthMiddleware()`の呼び出し時に`cfg.API.Auth0IssuerBaseURL`が渡されている
- 設定が正しく読み込まれている
- `AUTH0_ISSUER_BASE_URL`が空の場合、サーバー起動時にエラーが発生する

---

### Phase 5: サーバー側 - 新規private APIエンドポイント

#### - [ ] タスク 5.1: Today APIハンドラーの実装
**目的**: `/api/today`エンドポイントを実装

**作業内容**:
- `server/internal/api/handler/today_handler.go`を作成
- `TodayHandler`構造体を定義
- `NewTodayHandler()`関数を実装
- `RegisterTodayEndpoints()`関数を実装
  - `GET /api/today`エンドポイントを登録
  - エンドポイントの公開レベルを`"private"`として定義
  - ハンドラー関数内で公開レベルのチェックを実装
    - コンテキストからJWTの許容する公開レベルを取得
    - エンドポイントの公開レベル（`"private"`）と比較
    - Public API Key JWTでアクセスした場合、403 Forbiddenを返す
  - 今日の日付をYYYY-MM-DD形式で返す

**受け入れ基準**:
- `server/internal/api/handler/today_handler.go`が作成されている
- `TodayHandler`構造体が定義されている
- `GET /api/today`エンドポイントが登録されている
- エンドポイントの公開レベルが`"private"`として定義されている
- 公開レベルのチェックが実装されている
- 今日の日付がYYYY-MM-DD形式で返される

---

#### - [ ] タスク 5.2: Today APIのHuma API定義の追加
**目的**: `/api/today`エンドポイントのHuma API定義を追加

**作業内容**:
- `server/internal/api/huma/`ディレクトリに`GetTodayInput`と`TodayOutput`の定義を追加
- または既存のファイルに追加

**受け入れ基準**:
- `GetTodayInput`と`TodayOutput`が定義されている
- Huma API定義が正しく生成される

---

#### - [ ] タスク 5.3: ルーターでのToday API登録
**目的**: Today APIハンドラーをルーターに登録

**作業内容**:
- `server/internal/api/router/router.go`を修正
- `RegisterTodayEndpoints()`を呼び出してToday APIを登録

**受け入れ基準**:
- `RegisterTodayEndpoints()`が呼び出されている
- `/api/today`エンドポイントが正しく登録されている

---

### Phase 6: サーバー側 - 既存APIハンドラーの公開レベル定義

#### - [ ] タスク 6.1: User APIハンドラーの公開レベル定義
**目的**: User APIの各エンドポイントに公開レベル（`"public"`）を定義

**作業内容**:
- `server/internal/api/handler/user_handler.go`を修正
- 各エンドポイントのハンドラー関数内で公開レベルのチェックを実装
  - コンテキストからJWTの許容する公開レベルを取得
  - エンドポイントの公開レベル（`"public"`）と比較
  - Public API Key JWTでprivateなAPIにアクセスした場合、403 Forbiddenを返す（この場合は該当しないが、実装パターンとして追加）

**受け入れ基準**:
- すべてのUser APIエンドポイントに公開レベル（`"public"`）が定義されている
- 公開レベルのチェックが実装されている

---

#### - [ ] タスク 6.2: Post APIハンドラーの公開レベル定義
**目的**: Post APIの各エンドポイントに公開レベル（`"public"`）を定義

**作業内容**:
- `server/internal/api/handler/post_handler.go`を修正
- 各エンドポイントのハンドラー関数内で公開レベルのチェックを実装
  - コンテキストからJWTの許容する公開レベルを取得
  - エンドポイントの公開レベル（`"public"`）と比較

**受け入れ基準**:
- すべてのPost APIエンドポイントに公開レベル（`"public"`）が定義されている
- 公開レベルのチェックが実装されている

---

### Phase 7: クライアント側 - ApiClientクラスの修正

#### - [ ] タスク 7.1: ApiClientクラスのJWT取得ロジックの追加
**目的**: Auth0 JWTとPublic API Keyの切り替えを実装

**作業内容**:
- `client/src/lib/api.ts`を修正
- `getJWT()`メソッドの実装方法を検討
  - 注意: `useUser()`はReactフックのため、クラスメソッド内では直接使用できない
  - 実装方法の選択肢:
    1. ApiClientの`request()`メソッドにJWTを引数として渡す
    2. ApiClientのインスタンス作成時にJWT取得関数を渡す
    3. ApiClientを関数として実装し、`useUser()`を使用できるようにする
  - 要件定義書の意図を考慮して最適な方法を選択
- 選択した方法でJWT取得ロジックを実装
  - ログイン状態を確認（`useUser()`フックを使用可能な方法で）
  - ログイン中は`getAccessToken()`でJWTを取得
  - 未ログイン時は`NEXT_PUBLIC_API_KEY`を使用
- `request()`メソッドでJWTを`Authorization: Bearer <JWT>`ヘッダーで送信

**受け入れ基準**:
- JWT取得ロジックが実装されている
- ログイン中はAuth0 JWTを使用できる
- 未ログイン時は`NEXT_PUBLIC_API_KEY`を使用できる
- すべてのAPIリクエストに適切なJWTが送信される

---

### Phase 8: クライアント側 - TodayApiButtonコンポーネント

#### - [ ] タスク 8.1: TodayApiButtonコンポーネントの作成
**目的**: private API（`/api/today`）を呼び出すボタンコンポーネントを作成

**作業内容**:
- `client/src/components/TodayApiButton.tsx`を作成
- `useUser()`フックを使用してログイン状態を確認
- `getAccessToken()`でJWTを取得（ログイン中の場合）
- `NEXT_PUBLIC_API_KEY`を使用（未ログイン時）
- `/api/today`エンドポイントを呼び出す
- エラーハンドリングとエラーメッセージの表示を実装
- ローディング状態の表示を実装
- 日付の表示を実装

**受け入れ基準**:
- `client/src/components/TodayApiButton.tsx`が作成されている
- ログイン状態に応じて適切なJWTを使用できる
- エラーハンドリングが実装されている
- ローディング状態が表示される
- 日付が正しく表示される

---

#### - [ ] タスク 8.2: トップページへのTodayApiButtonの追加
**目的**: トップページにTodayApiButtonコンポーネントを追加

**作業内容**:
- `client/src/app/page.tsx`を修正
- `TodayApiButton`コンポーネントをインポート
- `TodayApiButton`コンポーネントを表示

**受け入れ基準**:
- `TodayApiButton`コンポーネントがトップページに表示される
- 既存のコンテンツに影響がない

---

### Phase 9: テスト

#### - [ ] タスク 9.1: JWT種類の判別テスト
**目的**: JWT種類の判別機能をテスト

**作業内容**:
- `server/internal/auth/jwt_test.go`を作成または修正
- `TestDetectJWTType()`関数を実装
- Auth0 JWTとPublic API Key JWTの判別テストを追加
- 未知のJWTタイプのテストを追加

**受け入れ基準**:
- `TestDetectJWTType()`関数が実装されている
- Auth0 JWTとPublic API Key JWTを正しく判別できる
- 未知のJWTタイプの場合、`JWTTypeUnknown`を返す

---

#### - [ ] タスク 9.2: Auth0 JWT検証テスト
**目的**: Auth0 JWT検証機能をテスト

**作業内容**:
- `server/internal/auth/auth0_validator_test.go`を作成
- `TestNewAuth0Validator()`関数を実装
- `TestValidateAuth0JWT()`関数を実装
- 実際のAuth0 JWTを使用したテストを追加（テスト用のAuth0 JWTを生成）

**受け入れ基準**:
- `TestNewAuth0Validator()`関数が実装されている
- `TestValidateAuth0JWT()`関数が実装されている
- Auth0 JWTを正しく検証できる

---

#### - [ ] タスク 9.3: 認証ミドルウェアの統合テスト
**目的**: 認証ミドルウェアの統合テストを実装

**作業内容**:
- `server/internal/auth/middleware_test.go`を作成または修正
- Auth0 JWTを使用したAPI呼び出しのテストを追加
- Public API Key JWTを使用したAPI呼び出しのテストを追加
- publicなAPIへのアクセステストを追加
- privateなAPIへのアクセステストを追加
- Public API Key JWTでprivateなAPIにアクセスした場合の403 Forbiddenテストを追加

**受け入れ基準**:
- Auth0 JWTを使用したAPI呼び出しが正常に動作する
- Public API Key JWTを使用したAPI呼び出しが正常に動作する
- publicなAPIへのアクセスが正常に動作する
- privateなAPIへのアクセスが正常に動作する（Auth0 JWTの場合）
- Public API Key JWTでprivateなAPIにアクセスした場合、403 Forbiddenを返す

---

#### - [ ] タスク 9.4: Today APIのテスト
**目的**: `/api/today`エンドポイントをテスト

**作業内容**:
- `server/internal/api/handler/today_handler_test.go`を作成
- `TestRegisterTodayEndpoints()`関数を実装
- Auth0 JWTを使用したアクセステストを追加
- Public API Key JWTを使用したアクセステストを追加（403 Forbiddenを期待）
- レスポンス形式のテストを追加

**受け入れ基準**:
- `/api/today`エンドポイントが正しく登録されている
- Auth0 JWTでアクセスした場合、正常に動作する
- Public API Key JWTでアクセスした場合、403 Forbiddenを返す
- レスポンス形式が正しい（`{"date": "YYYY-MM-DD"}`）

---

#### - [ ] タスク 9.5: クライアント側のE2Eテスト
**目的**: クライアント側のE2Eテストを実装

**作業内容**:
- `client/e2e/auth0-api-call.spec.ts`を作成
- 未ログイン状態でToday APIボタンをクリックした場合のエラーテストを追加
- Auth0ログイン後のToday APIボタンの動作テストを追加（テスト用のAuth0アカウントを使用）
  - テスト用の認証情報を環境変数で管理（`TEST_AUTH0_USERNAME`、`TEST_AUTH0_PASSWORD`）
  - 実際のAuth0ログインフローをテスト

**受け入れ基準**:
- 未ログイン状態でToday APIボタンをクリックした場合、エラーメッセージが表示される
- Auth0ログイン後のToday APIボタンが正常に動作する
- 日付が正しく表示される

---

### Phase 10: 動作確認とドキュメント更新

#### - [ ] タスク 10.1: 動作確認
**目的**: 実装した機能が正常に動作することを確認

**作業内容**:
- 開発環境でサーバーを起動
- クライアントを起動
- Auth0でログイン
- Auth0 JWTを使用してAPIを呼び出すことを確認
- Public API Key JWTを使用してAPIを呼び出すことを確認
- publicなAPIへのアクセスを確認
- privateなAPIへのアクセスを確認（Auth0 JWT）
- Public API Key JWTでprivateなAPIにアクセスした場合の403 Forbiddenを確認

**受け入れ基準**:
- すべての機能が正常に動作する
- エラーが発生しない

---

#### - [ ] タスク 10.2: ドキュメント更新
**目的**: 実装内容をドキュメントに反映

**作業内容**:
- README.mdにAuth0 API呼び出し機能の説明を追加（必要に応じて）
- APIドキュメントに`/api/today`エンドポイントの説明を追加（必要に応じて）

**受け入れ基準**:
- ドキュメントが更新されている
- 実装内容が正しく記載されている

---

## タスクの依存関係

### 必須の依存関係
- Phase 1 → Phase 2 → Phase 3 → Phase 4 → Phase 5 → Phase 6
- Phase 1 → Phase 7 → Phase 8
- Phase 4 → Phase 9（サーバー側テスト）
- Phase 5 → Phase 9（Today APIテスト）
- Phase 8 → Phase 9（クライアント側E2Eテスト）

### 並行実行可能
- Phase 5とPhase 6は並行実行可能
- Phase 7とPhase 8は並行実行可能
- Phase 9の各テストは並行実行可能

## 注意事項

### 実装時の注意点
1. **JWT種類の判別**: 署名検証前にパースするため、`ParseUnverified`を使用すること
2. **Auth0 JWT検証**: JWKSキャッシュの設定（12時間の定期更新、5分間隔の再取得制限、10秒のタイムアウト、未知のKIDが来たら再取得）を正しく実装すること
3. **API公開レベルの定義**: 各ハンドラーファイルでエンドポイントの公開レベルを定義すること
4. **クライアント側のJWT取得**: `useUser()`はReactフックのため、クラスメソッド内では直接使用できない。実装方法を慎重に検討すること
5. **エラーハンドリング**: 適切なエラーメッセージを返すこと

### テスト時の注意点
1. **Auth0 JWTのテスト**: 実際のAuth0 JWTを使用するか、テスト用のJWTを生成する必要がある
2. **テスト用のAuth0アカウント**: E2Eテストで使用するテスト用のAuth0アカウント（`go-webdb-template-test@nanasi.jp`）の認証情報を環境変数で管理すること
3. **テスト環境**: テスト環境と本番環境が完全に分離されていることを確認すること
