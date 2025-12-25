# Public API Key認証機能実装タスク一覧

## 概要
Public API Key認証機能の実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: 設定の拡張

#### - [ ] タスク 1.1: APIConfig構造体の追加
**目的**: APIキー設定項目を追加

**作業内容**:
- `server/internal/config/config.go`の`Config`構造体に`APIConfig`フィールドを追加
- `APIConfig`構造体を定義:
  ```go
  type APIConfig struct {
      CurrentVersion  string   `mapstructure:"current_version"`
      PublicKey       string   `mapstructure:"public_key"`       // オプション
      SecretKey       string   `mapstructure:"secret_key"`      // 必須
      InvalidVersions []string `mapstructure:"invalid_versions"`
  }
  ```
- `Config`構造体に`API APIConfig`フィールドを追加

**受け入れ基準**:
- `APIConfig`構造体が定義されている
- `Config`構造体に`API`フィールドが追加されている
- コンパイルエラーがない

---

#### - [ ] タスク 1.2: 設定ファイルの更新（develop環境）
**目的**: develop環境の設定ファイルにAPIキー設定を追加

**作業内容**:
- `config/develop/config.yaml`の`api`セクションに以下の項目を追加:
  ```yaml
  api:
    current_version: "v2"
    secret_key: "<SECRET_KEY>"  # 秘密鍵生成ツールで生成した値を設定
    invalid_versions:
      - "v1"
  ```
- 注意: `secret_key`は後で秘密鍵生成ツールで生成した値を設定する

**受け入れ基準**:
- `config/develop/config.yaml`に`api`セクションが追加されている
- YAML形式が正しい
- `current_version`が`"v2"`に設定されている
- `invalid_versions`に`"v1"`が含まれている

---

#### - [ ] タスク 1.3: 設定ファイルの更新（staging環境）
**目的**: staging環境の設定ファイルにAPIキー設定を追加

**作業内容**:
- `config/staging/config.yaml`の`api`セクションに以下の項目を追加:
  ```yaml
  api:
    current_version: "v2"
    secret_key: "<SECRET_KEY>"  # 秘密鍵生成ツールで生成した値を設定
    invalid_versions:
      - "v1"
  ```
- 注意: `secret_key`は後で秘密鍵生成ツールで生成した値を設定する

**受け入れ基準**:
- `config/staging/config.yaml`に`api`セクションが追加されている
- YAML形式が正しい

---

#### - [ ] タスク 1.4: 設定ファイルの更新（production環境）
**目的**: production環境の設定ファイル例にAPIキー設定を追加

**作業内容**:
- `config/production/config.yaml.example`の`api`セクションに以下の項目を追加（コメント付きで説明）:
  ```yaml
  api:
    current_version: "v2"
    # public_key: "<JWT_TOKEN>"  # オプション: 発行済みのPublic APIキー
    secret_key: "<SECRET_KEY>"  # 必須: 秘密鍵生成ツールで生成した値を設定
    invalid_versions:
      - "v1"
  ```

**受け入れ基準**:
- `config/production/config.yaml.example`に`api`セクションが追加されている
- コメントで説明が記載されている

---

#### - [ ] タスク 1.5: .gitignoreの更新
**目的**: staging環境の設定ファイルをcommit不可にする

**作業内容**:
- `.gitignore`に`config/staging/config.yaml`を追加
- 注意: `config/production/config.yaml`は既に`.gitignore`に追加済み

**受け入れ基準**:
- `.gitignore`に`config/staging/config.yaml`が追加されている

---

### Phase 2: JWT検証機能の実装

#### - [ ] タスク 2.1: authディレクトリの作成
**目的**: 認証機能用のディレクトリを作成

**作業内容**:
- `server/internal/auth/`ディレクトリを作成

**受け入れ基準**:
- `server/internal/auth/`ディレクトリが作成されている

---

#### - [ ] タスク 2.2: JWTClaims構造体の定義
**目的**: JWTクレーム構造体を定義

**作業内容**:
- `server/internal/auth/jwt.go`を作成
- `JWTClaims`構造体を定義:
  ```go
  type JWTClaims struct {
      Issuer   string   `json:"iss"`
      Subject  string   `json:"sub"`
      Type     string   `json:"type"`
      Scope    []string `json:"scope"`
      IssuedAt int64    `json:"iat"`
      Version  string   `json:"version"`
      Env      string   `json:"env"`
      jwt.RegisteredClaims
  }
  ```
- 必要なインポートを追加（`github.com/golang-jwt/jwt/v5`等）

**受け入れ基準**:
- `server/internal/auth/jwt.go`ファイルが作成されている
- `JWTClaims`構造体が定義されている
- コンパイルエラーがない

---

#### - [ ] タスク 2.3: JWTValidator構造体の定義
**目的**: JWT検証機能を提供する構造体を定義

**作業内容**:
- `JWTValidator`構造体を定義:
  ```go
  type JWTValidator struct {
      secretKey       string
      invalidVersions []string
      currentEnv      string
  }
  ```
- `NewJWTValidator`関数を実装

**受け入れ基準**:
- `JWTValidator`構造体が定義されている
- `NewJWTValidator`関数が実装されている
- コンパイルエラーがない

---

#### - [ ] タスク 2.4: ValidateJWT関数の実装
**目的**: JWTトークンの検証機能を実装

**作業内容**:
- `ValidateJWT`関数を実装:
  - JWTトークンのパース
  - 署名アルゴリズムの検証（HS256）
  - 署名の検証（秘密鍵による検証）
  - クレームの取得と検証
- `validateClaims`関数を実装:
  - `iss`の検証（`"go-webdb-template"`であること）
  - `type`の検証（`"public"`または`"private"`であること）
  - `version`の検証（無効バージョンリストとの照合）
  - `env`の検証（現在の環境と一致すること）

**受け入れ基準**:
- `ValidateJWT`関数が実装されている
- `validateClaims`関数が実装されている
- すべての検証項目が実装されている
- コンパイルエラーがない

---

#### - [ ] タスク 2.5: IsVersionInvalid関数の実装
**目的**: バージョン無効化チェック機能を実装

**作業内容**:
- `IsVersionInvalid`関数を実装:
  - 無効バージョンリストと照合
  - 一致する場合は`true`を返す

**受け入れ基準**:
- `IsVersionInvalid`関数が実装されている
- 無効バージョンリストとの照合が正しく動作する
- コンパイルエラーがない

---

#### - [ ] タスク 2.6: ParseJWTClaims関数の実装
**目的**: JWTトークンからクレームをパース（表示用）

**作業内容**:
- `ParseJWTClaims`関数を実装:
  - 署名検証なしでJWTトークンをパース
  - クレームを取得して返す
  - エラーハンドリングを実装

**受け入れ基準**:
- `ParseJWTClaims`関数が実装されている
- 署名検証なしでパースできる
- エラーハンドリングが実装されている
- コンパイルエラーがない

---

### Phase 3: 認証ミドルウェアの実装

#### - [ ] タスク 3.1: AuthMiddleware構造体の定義
**目的**: 認証ミドルウェア構造体を定義

**作業内容**:
- `server/internal/auth/middleware.go`を作成
- `AuthMiddleware`構造体を定義:
  ```go
  type AuthMiddleware struct {
      validator *JWTValidator
  }
  ```
- `NewAuthMiddleware`関数を実装

**受け入れ基準**:
- `server/internal/auth/middleware.go`ファイルが作成されている
- `AuthMiddleware`構造体が定義されている
- `NewAuthMiddleware`関数が実装されている
- コンパイルエラーがない

---

#### - [ ] タスク 3.2: Middleware関数の実装
**目的**: HTTPミドルウェア関数を実装

**作業内容**:
- `Middleware`関数を実装:
  - AuthorizationヘッダーからJWTトークンを取得
  - Bearerトークンの形式を検証
  - JWT検証を実行
  - スコープ検証を実行
  - 次のハンドラーを実行
- エラーレスポンスの実装:
  - 401 Unauthorized: `{"code": 401, "message": "..."}`
  - 403 Forbidden: `{"code": 403, "message": "..."}`

**受け入れ基準**:
- `Middleware`関数が実装されている
- Authorizationヘッダーの取得と検証が実装されている
- JWT検証が実装されている
- エラーレスポンスが正しい形式で返される
- コンパイルエラーがない

---

#### - [ ] タスク 3.3: validateScope関数の実装
**目的**: スコープ検証機能を実装

**作業内容**:
- `validateScope`関数を実装:
  - GETリクエストには`read`スコープが必要
  - POST/PUT/DELETEリクエストには`write`スコープが必要
  - スコープが不足している場合はエラーを返す

**受け入れ基準**:
- `validateScope`関数が実装されている
- すべてのHTTPメソッドに対するスコープ検証が実装されている
- コンパイルエラーがない

---

#### - [ ] タスク 3.4: writeErrorResponse関数の実装
**目的**: エラーレスポンス書き込み機能を実装

**作業内容**:
- `writeErrorResponse`関数を実装:
  - JSON形式でエラーレスポンスを返す
  - 形式: `{"code": <HTTPステータスコード>, "message": "<エラーメッセージ>"}`

**受け入れ基準**:
- `writeErrorResponse`関数が実装されている
- エラーレスポンスが正しい形式で返される
- コンパイルエラーがない

---

### Phase 4: 認証ミドルウェアの適用

#### - [ ] タスク 4.1: router.goの修正
**目的**: 認証ミドルウェアをAPIルーターに適用

**作業内容**:
- `server/internal/api/router/router.go`を修正:
  - `auth`パッケージをインポート
  - `NewRouter`関数で認証ミドルウェアを作成
  - `/api/*`パスに認証ミドルウェアを適用
  - 環境変数（`APP_ENV`）から現在の環境を取得
  - ミドルウェアの適用順序: CORS → 認証 → ハンドラー

**受け入れ基準**:
- 認証ミドルウェアが`/api/*`パスに適用されている
- `/health`エンドポイントは認証不要（`/api/*`の外側にあるため）
- コンパイルエラーがない

---

### Phase 5: GoAdminキー発行ページの実装

#### - [ ] タスク 5.1: api_key.goファイルの作成
**目的**: APIキー発行ページのファイルを作成

**作業内容**:
- `server/internal/admin/pages/api_key.go`を作成
- パッケージ宣言とインポート文を追加

**受け入れ基準**:
- `server/internal/admin/pages/api_key.go`ファイルが作成されている
- パッケージ宣言とインポート文が正しい
- コンパイルエラーがない

---

#### - [ ] タスク 5.2: APIKeyPage関数の実装
**目的**: APIキー発行ページのハンドラーを実装

**作業内容**:
- `APIKeyPage`関数を実装:
  - GETリクエスト: フォーム表示
  - POSTリクエスト: キー生成
  - ダウンロードリクエスト（`?download=true`）: ダウンロード処理
- 設定の取得とエラーハンドリング

**受け入れ基準**:
- `APIKeyPage`関数が実装されている
- GET/POST/ダウンロードリクエストの処理が実装されている
- エラーハンドリングが実装されている
- コンパイルエラーがない

---

#### - [ ] タスク 5.3: generatePublicAPIKey関数の実装
**目的**: Public JWTキー生成機能を実装

**作業内容**:
- `generatePublicAPIKey`関数を実装:
  - JWTクレームの作成:
    - `iss`: `"go-webdb-template"`（固定）
    - `sub`: `"public_client"`（固定）
    - `type`: `"public"`（固定）
    - `scope`: `["read", "write"]`（固定）
    - `iat`: 現在時刻（Unix timestamp）
    - `version`: 設定ファイルの`current_version`から取得
    - `env`: 現在の環境（`APP_ENV`から取得）
  - JWTトークンの署名（HS256、秘密鍵を使用）
  - エラーハンドリング

**受け入れ基準**:
- `generatePublicAPIKey`関数が実装されている
- すべてのJWTクレームが正しく設定されている
- JWTトークンが正しく署名されている
- コンパイルエラーがない

---

#### - [ ] タスク 5.4: handleGenerateKey関数の実装
**目的**: キー生成処理を実装

**作業内容**:
- `handleGenerateKey`関数を実装:
  - 現在の環境を取得
  - JWTトークンを生成
  - ペイロードをデコード
  - 生成結果を表示

**受け入れ基準**:
- `handleGenerateKey`関数が実装されている
- JWTトークンの生成とデコードが正しく動作する
- コンパイルエラーがない

---

#### - [ ] タスク 5.5: handleDownload関数の実装
**目的**: JWTトークンダウンロード機能を実装

**作業内容**:
- `handleDownload`関数を実装:
  - クエリパラメータからトークンを取得
  - ファイル名を生成（`api-key-{timestamp}.txt`）
  - ダウンロードレスポンスを設定:
    - `Content-Disposition`: `attachment; filename=...`
    - `Content-Type`: `text/plain`
  - トークンをレスポンスボディに書き込む

**受け入れ基準**:
- `handleDownload`関数が実装されている
- ダウンロードレスポンスが正しく設定されている
- コンパイルエラーがない

---

#### - [ ] タスク 5.6: renderAPIKeyPage関数の実装
**目的**: APIキー発行フォームのレンダリングを実装

**作業内容**:
- `renderAPIKeyPage`関数を実装:
  - HTMLフォームを生成
  - フォームのactionは`/admin/api-key`（POST）
  - 「APIキーを発行」ボタンを配置
  - GoAdminのテンプレートスタイルに合わせる

**受け入れ基準**:
- `renderAPIKeyPage`関数が実装されている
- HTMLフォームが正しく生成される
- コンパイルエラーがない

---

#### - [ ] タスク 5.7: renderAPIKeyResult関数の実装
**目的**: 生成結果表示のレンダリングを実装

**作業内容**:
- `renderAPIKeyResult`関数を実装:
  - JWTトークンを表示（textarea、readonly）
  - JWTペイロードを表示（preタグ、JSON整形）
  - 発行日時を表示（Unix timestamp → 人間が読める形式）
  - バージョンを表示
  - 環境を表示
  - ダウンロードボタンを配置（`/admin/api-key?download=true&token=...`）

**受け入れ基準**:
- `renderAPIKeyResult`関数が実装されている
- すべての情報が正しく表示される
- ダウンロードボタンが正しく動作する
- コンパイルエラーがない

---

#### - [ ] タスク 5.8: RegisterCustomPagesへの追加
**目的**: キー発行ページをGoAdminに登録

**作業内容**:
- `server/internal/admin/pages/pages.go`の`RegisterCustomPages`関数を修正:
  - `"/api-key"`パスを追加
  - `APIKeyPage`関数をハンドラーとして登録

**受け入れ基準**:
- `RegisterCustomPages`関数に`"/api-key"`が追加されている
- コンパイルエラーがない

---

### Phase 6: メニュー項目の追加

#### - [ ] タスク 6.1: マイグレーションファイルの作成
**目的**: メニュー項目追加用のマイグレーションファイルを作成

**作業内容**:
- `db/migrations/shard1/004_api_key_menu.sql`を作成
- SQL文を記述:
  ```sql
  -- APIキー発行（カスタムページの子メニュー）
  INSERT INTO goadmin_menu (parent_id, type, "order", title, icon, uri, plugin_name, created_at, updated_at)
  VALUES (
      (SELECT id FROM goadmin_menu WHERE title = 'カスタムページ'),
      1, 2, 'APIキー発行', 'fa-key', '/api-key', '', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
  );
  ```

**受け入れ基準**:
- `db/migrations/shard1/004_api_key_menu.sql`ファイルが作成されている
- SQL文が正しい
- 既存の「カスタムページ」カテゴリに追加されている

---

#### - [ ] タスク 6.2: 他のシャードへのマイグレーションコピー
**目的**: すべてのシャードに同じマイグレーションを適用

**作業内容**:
- `db/migrations/shard2/004_api_key_menu.sql`を作成（shard1と同じ内容）
- `db/migrations/shard3/004_api_key_menu.sql`を作成（shard1と同じ内容）
- `db/migrations/shard4/004_api_key_menu.sql`を作成（shard1と同じ内容）

**受け入れ基準**:
- すべてのシャード（shard1, shard2, shard3, shard4）にマイグレーションファイルが作成されている
- すべてのマイグレーションファイルの内容が同じである

---

### Phase 7: 秘密鍵生成ツールの実装

#### - [ ] タスク 7.1: generate-secretディレクトリの作成
**目的**: 秘密鍵生成ツール用のディレクトリを作成

**作業内容**:
- `server/cmd/generate-secret/`ディレクトリを作成

**受け入れ基準**:
- `server/cmd/generate-secret/`ディレクトリが作成されている

---

#### - [ ] タスク 7.2: main.goの実装
**目的**: 秘密鍵生成ツールを実装

**作業内容**:
- `server/cmd/generate-secret/main.go`を作成
- ランダムな秘密鍵を生成（32バイト = 256ビット）
- Base64エンコード
- 標準出力に表示
- エラーハンドリング

**受け入れ基準**:
- `server/cmd/generate-secret/main.go`ファイルが作成されている
- ランダムな秘密鍵が生成される
- Base64エンコードが正しく動作する
- コンパイルエラーがない

---

### Phase 8: クライアント側の実装

#### - [ ] タスク 8.1: api.tsの修正
**目的**: APIクライアントに認証ヘッダーを追加

**作業内容**:
- `client/src/lib/api.ts`を修正:
  - 環境変数（`NEXT_PUBLIC_API_KEY`）からAPIキーを取得
  - コンストラクタでAPIキーが設定されていない場合、エラーを投げる
  - `request`メソッドに`Authorization: Bearer <API_KEY>`ヘッダーを追加
  - 401/403エラー時のエラーハンドリングを改善

**受け入れ基準**:
- 環境変数からAPIキーを取得できる
- APIキーが設定されていない場合、エラーを投げる
- すべてのAPIリクエストに`Authorization`ヘッダーが付与される
- 401/403エラー時の処理が実装されている
- コンパイルエラーがない

---

### Phase 9: テスト用ダミーAPIキーの作成

#### - [ ] タスク 9.1: testdataディレクトリの作成
**目的**: テスト用設定ファイル用のディレクトリを作成

**作業内容**:
- `server/internal/config/testdata/develop/`ディレクトリが存在することを確認（既存の可能性あり）
- 存在しない場合は作成

**受け入れ基準**:
- `server/internal/config/testdata/develop/`ディレクトリが存在する

---

#### - [ ] タスク 9.2: テスト用設定ファイルの作成
**目的**: テスト用のダミーAPIキー設定を作成

**作業内容**:
- `server/internal/config/testdata/develop/api_key.yaml`を作成
- テスト用の秘密鍵を設定
- テスト用のダミーAPIキーを生成（JWT形式）
- 設定内容:
  ```yaml
  api:
    current_version: "v2"
    secret_key: "<TEST_SECRET_KEY>"  # テスト用の秘密鍵
    invalid_versions:
      - "v1"
  ```
- 注意: テスト用の秘密鍵は固定値で良い

**受け入れ基準**:
- `server/internal/config/testdata/develop/api_key.yaml`ファイルが作成されている
- テスト用の秘密鍵が設定されている
- YAML形式が正しい

---

### Phase 10: テストの実装

#### - [ ] タスク 10.1: JWT検証機能のテスト
**目的**: JWT検証機能のユニットテストを実装

**作業内容**:
- `server/internal/auth/jwt_test.go`を作成
- テストケース:
  - 正常系: 有効なJWTトークンの検証
  - 異常系: 無効な署名のJWTトークン
  - 異常系: 不正なiss
  - 異常系: 不正なtype
  - 異常系: 無効バージョンのキー
  - 異常系: 環境不一致
- テスト用のダミーAPIキーを使用

**受け入れ基準**:
- `server/internal/auth/jwt_test.go`ファイルが作成されている
- すべてのテストケースが実装されている
- テストが正常に実行できる

---

#### - [ ] タスク 10.2: 認証ミドルウェアのテスト
**目的**: 認証ミドルウェアのユニットテストを実装

**作業内容**:
- `server/internal/auth/middleware_test.go`を作成
- テストケース:
  - 正常系: 有効なAPIキーでのリクエスト
  - 異常系: Authorizationヘッダーなし
  - 異常系: 無効なAPIキー
  - 異常系: スコープ不足（GETリクエストにreadスコープなし）
  - 異常系: スコープ不足（POSTリクエストにwriteスコープなし）
- テスト用のダミーAPIキーを使用

**受け入れ基準**:
- `server/internal/auth/middleware_test.go`ファイルが作成されている
- すべてのテストケースが実装されている
- テストが正常に実行できる

---

#### - [ ] タスク 10.3: API認証の統合テスト
**目的**: API認証の統合テストを実装

**作業内容**:
- `server/test/integration/api_auth_test.go`を作成（または既存の統合テストファイルに追加）
- テストケース:
  - 正常系: 有効なAPIキーでAPIアクセス
  - 異常系: 無効なAPIキーでAPIアクセス
  - 異常系: APIキーなしでAPIアクセス
  - 異常系: スコープ不足でのAPIアクセス
- テスト用のダミーAPIキーを使用

**受け入れ基準**:
- 統合テストファイルが作成されている
- すべてのテストケースが実装されている
- テストが正常に実行できる

---

#### - [ ] タスク 10.4: クライアント側のテスト
**目的**: クライアント側のテストを実装

**作業内容**:
- `client/src/lib/__tests__/api.test.ts`を更新（または作成）
- テストケース:
  - 正常系: APIキーを設定してAPIアクセス
  - 異常系: APIキー未設定時のエラー
  - 異常系: 401/403エラー時の処理
- モックを使用してAPIリクエストをテスト

**受け入れ基準**:
- クライアント側のテストファイルが作成されている
- すべてのテストケースが実装されている
- テストが正常に実行できる

---

### Phase 11: ドキュメント更新

#### - [ ] タスク 11.1: READMEの更新
**目的**: READMEにAPIキー認証機能の説明を追加

**作業内容**:
- `README.md`を更新（または該当するドキュメントファイル）
- APIキー認証機能の説明を追加
- 秘密鍵生成ツールの使用方法を追加
- クライアント側の環境変数設定方法を追加

**受け入れ基準**:
- READMEにAPIキー認証機能の説明が追加されている
- 秘密鍵生成ツールの使用方法が記載されている
- クライアント側の環境変数設定方法が記載されている

---

## 受け入れ基準（全体）

### 機能要件
- [ ] Public APIキーでAPIアクセス可能
- [ ] 無効なversionのキーは拒否される（401 Unauthorized）
- [ ] スコープが不足している場合は403 Forbiddenを返す
- [ ] GoAdmin管理画面でPublicキーを発行できる
- [ ] 管理画面のメニューに「APIキー発行」項目が表示される
- [ ] メニューからキー発行ページにアクセスできる
- [ ] 生成されたJWTのペイロードが画面に表示される
- [ ] JWTトークンをダウンロードできる
- [ ] 認証ミドルウェアがすべてのAPIエンドポイントに適用される
- [ ] クライアント側（TypeScript）がAPIキーを正しく送信する
- [ ] 環境変数（`NEXT_PUBLIC_API_KEY`）からAPIキーを取得できる
- [ ] APIキーが設定されていない場合、エラーを投げてリクエストを送信しない

### 非機能要件
- [ ] JWT署名の検証が正常に動作する
- [ ] 認証処理がAPIレスポンス時間に大きな影響を与えない
- [ ] 認証失敗時に適切なエラーレスポンス（401/403）を返す
- [ ] エラーレスポンスが`{"code": <HTTPステータスコード>, "message": "..."}`形式である
- [ ] 既存のAPIエンドポイントが正常に動作する
- [ ] 既存のCORS設定が正常に動作する
- [ ] 既存のGoAdmin管理画面が正常に動作する

### 設定
- [ ] `config/{env}/config.yaml`にAPIキー設定が追加されている
- [ ] `api.current_version`が設定ファイルに定義されている
- [ ] `api.secret_key`が設定ファイルに定義されている
- [ ] 秘密鍵生成ツールが正常に動作する
- [ ] `api.invalid_versions`が設定ファイルで管理できる
- [ ] `config/staging/config.yaml`が`.gitignore`に追加されている
- [ ] テスト用のダミーAPIキーが`testdata/`に配置されている

### テスト
- [ ] 既存のテストが正常に動作する
- [ ] 認証機能のテストが実装されている
- [ ] テスト用のダミーAPIキーでテストが正常に実行できる
- [ ] CI/CDでのテスト実行が正常に動作する

### クライアント側
- [ ] TypeScriptクライアントがAPIキーを正しく送信する
- [ ] 環境変数からAPIキーを取得できる
- [ ] APIキーなしのリクエストは適切にエラーを返す
- [ ] 401/403エラー時の適切なエラーハンドリングが実装されている

