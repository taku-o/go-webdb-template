# Echo・Huma導入実装タスク一覧

## 概要
Echo・Huma導入の実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: 依存関係の追加

#### - [ ] タスク 1.1: EchoとHumaの依存関係を追加
**目的**: EchoとHumaフレームワークの依存関係をgo.modに追加

**作業内容**:
- `server/go.mod`に以下の依存関係を追加:
  - `github.com/labstack/echo/v4`
  - `github.com/danielgtaylor/huma/v2`
  - `github.com/danielgtaylor/huma/v2/adapters/humaecho`
- `go mod tidy`を実行して依存関係を解決

**受け入れ基準**:
- すべての依存関係が追加されている
- `go mod tidy`が正常に実行される
- コンパイルエラーがない

---

### Phase 2: Echoフレームワークの導入

#### - [ ] タスク 2.1: Echoインスタンスの作成と基本設定
**目的**: Echoインスタンスを作成し、基本的な設定を行う

**作業内容**:
- `server/cmd/server/main.go`を更新
- `echo.New()`でEchoインスタンスを作成
- デバッグモードの設定（開発環境のみ）
- Recoverミドルウェアの追加

**受け入れ基準**:
- Echoインスタンスが作成されている
- デバッグモードが適切に設定されている
- Recoverミドルウェアが追加されている

---

#### - [ ] タスク 2.2: 認証ミドルウェアのEcho形式への変換
**目的**: 既存の認証ミドルウェアをEcho形式に変換

**作業内容**:
- `server/internal/auth/middleware.go`を更新
- `echo.MiddlewareFunc`を返す関数に変更
- 既存のJWT検証ロジックを維持
- エラーレスポンスをEcho形式に変更

**受け入れ基準**:
- 認証ミドルウェアがEcho形式に変換されている
- 既存のJWT検証ロジックが維持されている
- エラーレスポンスが適切に返される

---

#### - [ ] タスク 2.3: CORS設定のEcho形式への移行
**目的**: 既存のCORS設定をEchoのCORSミドルウェアに移行

**作業内容**:
- `server/internal/api/router/router.go`を更新
- EchoのCORSミドルウェアを使用
- 既存のCORS設定をEcho形式に変換

**受け入れ基準**:
- CORS設定がEcho形式に移行されている
- 既存のCORS設定が維持されている

---

#### - [ ] タスク 2.4: アクセスログのEcho形式への統合
**目的**: 既存のアクセスログ機能をEcho形式に統合

**作業内容**:
- `server/internal/logging/middleware.go`を更新
- Echoミドルウェア形式に変換
- 既存のログ機能を維持

**受け入れ基準**:
- アクセスログがEcho形式に統合されている
- 既存のログ機能が維持されている

---

### Phase 3: Humaフレームワークの導入

#### - [ ] タスク 3.1: Huma APIインスタンスの作成
**目的**: Huma APIインスタンスを作成し、基本設定を行う

**作業内容**:
- `server/internal/api/router/router.go`を更新
- `humaecho.New(e, config)`でHuma APIインスタンスを作成
- Huma設定（API名、バージョン、OpenAPIパスなど）を設定

**受け入れ基準**:
- Huma APIインスタンスが作成されている
- Huma設定が適切に設定されている

---

### Phase 4: リクエスト/レスポンス構造体の定義

#### - [ ] タスク 4.1: Huma用ディレクトリとパッケージの作成
**目的**: Huma用の構造体定義用ディレクトリとパッケージを作成

**作業内容**:
- `server/internal/api/huma/`ディレクトリを作成
- `inputs.go`と`outputs.go`ファイルを作成
- パッケージ名を`humaapi`に設定

**受け入れ基準**:
- ディレクトリとファイルが作成されている
- パッケージ名が適切に設定されている

---

#### - [ ] タスク 4.2: ユーザーエンドポイントのリクエスト構造体の定義
**目的**: ユーザーエンドポイント用のリクエスト構造体を定義

**作業内容**:
- `server/internal/api/huma/inputs.go`に以下を定義:
  - `CreateUserInput`
  - `GetUserInput`
  - `UpdateUserInput`
  - `DeleteUserInput`
  - `ListUsersInput`
- Humaタグを適切に設定

**受け入れ基準**:
- すべてのリクエスト構造体が定義されている
- Humaタグが適切に設定されている
- バリデーションタグが設定されている

---

#### - [ ] タスク 4.3: ユーザーエンドポイントのレスポンス構造体の定義
**目的**: ユーザーエンドポイント用のレスポンス構造体を定義

**作業内容**:
- `server/internal/api/huma/outputs.go`に以下を定義:
  - `UserOutput`
  - `UsersOutput`
  - `DeleteUserOutput`（204 No Content用）

**受け入れ基準**:
- すべてのレスポンス構造体が定義されている
- `Body`フィールドが適切に設定されている

---

#### - [ ] タスク 4.4: 投稿エンドポイントのリクエスト構造体の定義
**目的**: 投稿エンドポイント用のリクエスト構造体を定義

**作業内容**:
- `server/internal/api/huma/inputs.go`に以下を定義:
  - `CreatePostInput`
  - `GetPostInput`
  - `UpdatePostInput`
  - `DeletePostInput`
  - `ListPostsInput`
  - `GetUserPostsInput`
- Humaタグを適切に設定

**受け入れ基準**:
- すべてのリクエスト構造体が定義されている
- Humaタグが適切に設定されている
- バリデーションタグが設定されている

---

#### - [ ] タスク 4.5: 投稿エンドポイントのレスポンス構造体の定義
**目的**: 投稿エンドポイント用のレスポンス構造体を定義

**作業内容**:
- `server/internal/api/huma/outputs.go`に以下を定義:
  - `PostOutput`
  - `PostsOutput`
  - `UserPostsOutput`
  - `DeletePostOutput`（204 No Content用）

**受け入れ基準**:
- すべてのレスポンス構造体が定義されている
- `Body`フィールドが適切に設定されている

---

### Phase 5: ユーザーエンドポイントのHuma化

#### - [ ] タスク 5.1: CreateUserエンドポイントのHuma化
**目的**: POST /api/usersエンドポイントをHumaで実装

**作業内容**:
- `server/internal/api/handler/user_handler.go`を更新
- `huma.Register()`でエンドポイントを登録
- 既存のService層呼び出しを維持
- エラーハンドリングをHuma形式に変更

**受け入れ基準**:
- エンドポイントがHumaで実装されている
- 既存のService層呼び出しが維持されている
- エラーハンドリングが適切に実装されている

---

#### - [ ] タスク 5.2: GetUserエンドポイントのHuma化
**目的**: GET /api/users/{id}エンドポイントをHumaで実装

**作業内容**:
- `server/internal/api/handler/user_handler.go`を更新
- `huma.Register()`でエンドポイントを登録
- パスパラメータの取得方法を変更
- 既存のService層呼び出しを維持

**受け入れ基準**:
- エンドポイントがHumaで実装されている
- パスパラメータが正しく取得される
- 既存のService層呼び出しが維持されている

---

#### - [ ] タスク 5.3: ListUsersエンドポイントのHuma化
**目的**: GET /api/usersエンドポイントをHumaで実装

**作業内容**:
- `server/internal/api/handler/user_handler.go`を更新
- `huma.Register()`でエンドポイントを登録
- クエリパラメータの取得方法を変更
- 既存のService層呼び出しを維持

**受け入れ基準**:
- エンドポイントがHumaで実装されている
- クエリパラメータが正しく取得される
- 既存のService層呼び出しが維持されている

---

#### - [ ] タスク 5.4: UpdateUserエンドポイントのHuma化
**目的**: PUT /api/users/{id}エンドポイントをHumaで実装

**作業内容**:
- `server/internal/api/handler/user_handler.go`を更新
- `huma.Register()`でエンドポイントを登録
- パスパラメータとリクエストボディの取得方法を変更
- 既存のService層呼び出しを維持

**受け入れ基準**:
- エンドポイントがHumaで実装されている
- パスパラメータとリクエストボディが正しく取得される
- 既存のService層呼び出しが維持されている

---

#### - [ ] タスク 5.5: DeleteUserエンドポイントのHuma化
**目的**: DELETE /api/users/{id}エンドポイントをHumaで実装

**作業内容**:
- `server/internal/api/handler/user_handler.go`を更新
- `huma.Register()`でエンドポイントを登録
- パスパラメータの取得方法を変更
- 既存のService層呼び出しを維持
- 204 No Contentレスポンスを返す

**受け入れ基準**:
- エンドポイントがHumaで実装されている
- パスパラメータが正しく取得される
- 204 No Contentレスポンスが返される

---

### Phase 6: 投稿エンドポイントのHuma化

#### - [ ] タスク 6.1: CreatePostエンドポイントのHuma化
**目的**: POST /api/postsエンドポイントをHumaで実装

**作業内容**:
- `server/internal/api/handler/post_handler.go`を更新
- `huma.Register()`でエンドポイントを登録
- 既存のService層呼び出しを維持
- エラーハンドリングをHuma形式に変更

**受け入れ基準**:
- エンドポイントがHumaで実装されている
- 既存のService層呼び出しが維持されている
- エラーハンドリングが適切に実装されている

---

#### - [ ] タスク 6.2: GetPostエンドポイントのHuma化
**目的**: GET /api/posts/{id}エンドポイントをHumaで実装

**作業内容**:
- `server/internal/api/handler/post_handler.go`を更新
- `huma.Register()`でエンドポイントを登録
- パスパラメータとクエリパラメータの取得方法を変更
- 既存のService層呼び出しを維持

**受け入れ基準**:
- エンドポイントがHumaで実装されている
- パスパラメータとクエリパラメータが正しく取得される
- 既存のService層呼び出しが維持されている

---

#### - [ ] タスク 6.3: ListPostsエンドポイントのHuma化
**目的**: GET /api/postsエンドポイントをHumaで実装

**作業内容**:
- `server/internal/api/handler/post_handler.go`を更新
- `huma.Register()`でエンドポイントを登録
- クエリパラメータの取得方法を変更
- `user_id`クエリパラメータの処理を実装
- 既存のService層呼び出しを維持

**受け入れ基準**:
- エンドポイントがHumaで実装されている
- クエリパラメータが正しく取得される
- `user_id`が指定された場合の処理が正しく動作する

---

#### - [ ] タスク 6.4: UpdatePostエンドポイントのHuma化
**目的**: PUT /api/posts/{id}エンドポイントをHumaで実装

**作業内容**:
- `server/internal/api/handler/post_handler.go`を更新
- `huma.Register()`でエンドポイントを登録
- パスパラメータ、クエリパラメータ、リクエストボディの取得方法を変更
- 既存のService層呼び出しを維持

**受け入れ基準**:
- エンドポイントがHumaで実装されている
- パスパラメータ、クエリパラメータ、リクエストボディが正しく取得される
- 既存のService層呼び出しが維持されている

---

#### - [ ] タスク 6.5: DeletePostエンドポイントのHuma化
**目的**: DELETE /api/posts/{id}エンドポイントをHumaで実装

**作業内容**:
- `server/internal/api/handler/post_handler.go`を更新
- `huma.Register()`でエンドポイントを登録
- パスパラメータとクエリパラメータの取得方法を変更
- 既存のService層呼び出しを維持
- 204 No Contentレスポンスを返す

**受け入れ基準**:
- エンドポイントがHumaで実装されている
- パスパラメータとクエリパラメータが正しく取得される
- 204 No Contentレスポンスが返される

---

#### - [ ] タスク 6.6: GetUserPostsエンドポイントのHuma化
**目的**: GET /api/user-postsエンドポイントをHumaで実装

**作業内容**:
- `server/internal/api/handler/post_handler.go`を更新
- `huma.Register()`でエンドポイントを登録
- クエリパラメータの取得方法を変更
- 既存のService層呼び出しを維持

**受け入れ基準**:
- エンドポイントがHumaで実装されている
- クエリパラメータが正しく取得される
- 既存のService層呼び出しが維持されている

---

### Phase 7: ヘルスチェックエンドポイントとルーター設定

#### - [ ] タスク 7.1: ヘルスチェックエンドポイントの実装
**目的**: GET /healthエンドポイントをEchoで実装

**作業内容**:
- `server/internal/api/router/router.go`を更新
- Echoの通常ハンドラーとして実装
- 200 OKレスポンスを返す

**受け入れ基準**:
- ヘルスチェックエンドポイントがEchoで実装されている
- 200 OKレスポンスが返される

---

#### - [ ] タスク 7.2: ルーターの統合とEchoサーバー起動
**目的**: ルーターを統合し、Echoサーバーを起動するように変更

**作業内容**:
- `server/internal/api/router/router.go`を更新
- 既存のGorilla MuxルーターをEchoルーターに置き換え
- すべてのエンドポイントを登録
- ミドルウェアの適用順序を設定
- `server/cmd/server/main.go`を更新
- Echoサーバーの起動方法に変更

**受け入れ基準**:
- ルーターがEcho形式に置き換えられている
- すべてのエンドポイントが登録されている
- ミドルウェアが適切な順序で適用されている
- Echoサーバーが正常に起動する

---

### Phase 8: テストと検証

#### - [ ] タスク 8.1: 既存テストの更新
**目的**: 既存のテストコードをEcho/Huma形式に対応

**作業内容**:
- `server/test/e2e/api_test.go`を更新
- Echo/Humaのテストユーティリティを使用
- テストが正常に実行されることを確認

**受け入れ基準**:
- 既存のテストが正常に実行される
- すべてのテストが通過する

---

#### - [ ] タスク 8.2: 新しいテストの追加
**目的**: Humaの動作を検証する新しいテストを追加

**作業内容**:
- ハンドラーの単体テストを追加
- リクエスト/レスポンス構造体のテストを追加
- エラーハンドリングのテストを追加

**受け入れ基準**:
- 新しいテストが追加されている
- すべてのテストが通過する

---

#### - [ ] タスク 8.3: OpenAPI仕様の検証
**目的**: OpenAPI仕様が正しく生成されることを確認

**作業内容**:
- `/openapi.json`エンドポイントにアクセス
- `/openapi.yaml`エンドポイントにアクセス
- OpenAPI仕様の妥当性を確認
- 全エンドポイントが含まれていることを確認

**受け入れ基準**:
- OpenAPI仕様が正しく生成されている
- 全エンドポイントが含まれている
- リクエスト/レスポンススキーマが正しく定義されている

---

### Phase 9: ドキュメント更新

#### - [ ] タスク 9.1: APIドキュメントの更新
**目的**: APIドキュメントにOpenAPI仕様への参照を追加

**作業内容**:
- `docs/API.md`を更新
- OpenAPI仕様への参照を追加
- 変更内容を文書化

**受け入れ基準**:
- APIドキュメントが更新されている
- OpenAPI仕様への参照が追加されている

---

#### - [ ] タスク 9.2: READMEの更新
**目的**: READMEに依存関係の更新を反映

**作業内容**:
- `README.md`を更新（必要に応じて）
- 依存関係の変更を反映

**受け入れ基準**:
- READMEが更新されている（必要に応じて）

---

## 実装順序の推奨

1. **Phase 1**: 依存関係の追加
2. **Phase 2**: Echoフレームワークの導入（ミドルウェアの移行）
3. **Phase 3**: Humaフレームワークの導入
4. **Phase 4**: リクエスト/レスポンス構造体の定義
5. **Phase 5-6**: エンドポイントのHuma化（ユーザー→投稿の順）
6. **Phase 7**: ヘルスチェックとルーター統合
7. **Phase 8**: テストと検証
8. **Phase 9**: ドキュメント更新

## 注意事項

- 各フェーズの完了後に動作確認を行うこと
- 既存のService層とRepository層は変更しないこと
- 既存のAPI仕様（リクエスト/レスポンス形式）を維持すること
- エラーハンドリングはHumaの標準形式を使用すること
- テストは段階的に追加・更新すること

