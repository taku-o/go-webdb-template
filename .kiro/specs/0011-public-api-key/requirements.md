# Public API Key認証機能要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #21
- **Issueタイトル**: APIサーバーに、PublicなAPIキーによる認証を追加する
- **Feature名**: 0011-public-api-key
- **作成日**: 2025-01-27

### 1.2 目的
APIサーバーにJWT形式のPublic APIキーによる認証機能を実装する。
これにより、外部クライアントがAPIキーを使用してAPIサーバーにアクセスできるようになり、セキュアなAPIアクセスを実現する。

### 1.3 スコープ
- JWT形式のAPIキーによる認証機能の実装
- Public APIキーとPrivate APIキーの2種類のサポート（Privateは将来実装）
- AuthorizationヘッダーでのBearerトークン受信と検証
- APIキーの無効化機能（versionベース）
- GoAdmin管理画面でのPublic JWTキー発行機能
- クライアント側（TypeScript/Next.js）でのAPIキー送信実装
- テスト用ダミーAPIキーの提供

## 2. 背景・現状分析

### 2.1 現在の実装
- **APIサーバー**: `server/internal/api/router/router.go`でルーティングを管理
- **認証機能**: 現在、認証機能は実装されていない
- **APIエンドポイント**: `/api/users`, `/api/posts`等のエンドポイントが存在
- **CORS設定**: `server/internal/config/config.go`でCORS設定が管理されている
- **GoAdmin管理画面**: `server/internal/admin/`に実装済み（0004-goadmin）
- **クライアント**: `client/src/lib/api.ts`でTypeScript/Next.jsクライアントが実装済み
- **設定管理**: `server/internal/config/config.go`で設定を管理
- **環境別設定**: `config/{env}/config.yaml`で環境別設定を管理

### 2.2 課題点
1. **認証機能の不足**: APIサーバーに認証機能がなく、誰でもアクセス可能な状態
2. **セキュリティリスク**: 認証なしでAPIにアクセスできるため、不正アクセスのリスクがある
3. **APIキー管理の不足**: APIキーを発行・管理する機能がない
4. **クライアント側の対応不足**: クライアント側でAPIキーを送信する実装がない

### 2.3 本実装による改善点
1. **セキュリティ向上**: JWT形式のAPIキーによる認証により、セキュアなAPIアクセスを実現
2. **APIキー管理**: GoAdmin管理画面からPublic APIキーを発行・管理できる
3. **無効化機能**: versionベースの無効化により、脆弱性発見時に迅速に対応可能
4. **クライアント対応**: クライアント側でAPIキーを送信する実装により、エンドツーエンドで認証が機能

## 3. 機能要件

### 3.1 JWT APIキー認証機能

#### 3.1.1 基本機能
- AuthorizationヘッダーでBearerトークンとしてAPIキーを受信
- JWT形式のAPIキーを検証
- Public APIキーとPrivate APIキーの2種類をサポート（Privateは将来実装）
- すべてのAPIエンドポイント（`/api/*`）に認証ミドルウェアを適用
- 認証に失敗した場合は401 Unauthorizedを返す

#### 3.1.2 JWTデータ構造
JWTのペイロードには以下の情報を含む：

```json
{
  "iss": "go-webdb-template",
  "sub": "1234566790",                        // PrivateキーならユーザーID、Publicなら "public_client"
  "type": "public" | "private",               // キーの種別
  "scope": ["read", "write"],                 // 許可する操作（固定）
  "iat": 1715654400,                          // 発行日
  "version": "v1",                           // 無効化ロジックに利用
  "env": "develop" | "staging" | "production"  // 環境
}
```

**注意事項**:
- Public APIキーは有効期限（`exp`）を設けない（無期限）
- Private APIキーは将来実装時に有効期限を設ける可能性があるが、今回は対象外

#### 3.1.3 認証処理フロー
1. クライアントがリクエストに`Authorization: Bearer <JWT_API_KEY>`ヘッダーを付与
2. サーバー側の認証ミドルウェアがJWTトークンを検証
3. 以下の検証を実施:
   - JWT署名の検証（秘密鍵による検証）
   - `iss`が`"go-webdb-template"`であること
   - `type`が`"public"`または`"private"`であること
   - `version`が無効バージョンリストに含まれていないこと
   - `env`が現在の環境と一致すること
4. 検証成功時はリクエストを処理、失敗時は401 Unauthorizedを返す

#### 3.1.4 スコープ検証
- `scope`フィールドに`"read"`が含まれている場合、GETリクエストを許可
- `scope`フィールドに`"write"`が含まれている場合、POST/PUT/DELETEリクエストを許可
- スコープが不足している場合は403 Forbiddenを返す

### 3.2 APIキー無効化機能

#### 3.2.1 versionベースの無効化
- JWTの`version`フィールドを確認
- 設定ファイルに無効バージョンリストを定義
- 無効バージョンのAPIキーは401 Unauthorizedを返す

#### 3.2.2 無効バージョン管理
- 設定ファイル（`config/{env}/config.yaml`）に`api.invalid_versions`フィールドを追加
- 無効バージョンリストは文字列配列で定義
- 例: `["v1", "v2"]`のように複数のバージョンを無効化可能

#### 3.2.3 無効化の例
```yaml
api:
  invalid_versions:
    - "v1"  # v1のキーに脆弱性が見つかった場合
```

### 3.3 GoAdmin管理画面統合

#### 3.3.1 Public JWTキー発行ページ
- GoAdmin管理画面にPublic JWTキーを発行するカスタムページを追加
- ページパス: `/admin/api-key`（既存のカスタムページパターンに合わせる。`RegisterCustomPages`で`/api-key`と登録すると、GoAdminが自動的に`/admin`プレフィックスを追加する）
- キー生成・表示機能を実装
- 管理画面のメニューに「APIキー発行」メニュー項目を追加

#### 3.3.1.1 メニュー項目の追加
- GoAdmin管理画面のメニューに「APIキー発行」項目を追加
- メニューの配置: 既存の「カスタムページ」カテゴリに追加
- メニュー項目の設定:
  - 親カテゴリ: 「カスタムページ」（`parent_id`は「カスタムページ」カテゴリのID）
  - タイトル: 「APIキー発行」
  - アイコン: `fa-key`（Font Awesome）
  - URI: `/api-key`（カスタムページのパス。GoAdminが自動的に`/admin`プレフィックスを追加し、実際のURLは`/admin/api-key`になる）
  - 順序: 既存メニュー（「ユーザー登録」）の後に配置（例: order = 2）
- データベースマイグレーションファイルまたは初期データとして追加

#### 3.3.2 キー発行機能
- ボタンクリックで新しいPublic JWTキーを生成
- 生成されたJWTトークンを画面に表示
- JWTのペイロード（デコードされたJSON）を画面に表示
- JWTトークンのダウンロードボタンを用意（テキストファイルとしてダウンロード）
- キーのコピー機能（オプション）
- 発行日時、有効期限（オプション）の表示

#### 3.3.2.1 JWTペイロードの表示
- JWTトークンをデコードしてペイロードを取得
- ペイロードのJSONを整形して表示（読みやすい形式）
- 表示項目:
  - `iss`: 発行者
  - `sub`: サブジェクト
  - `type`: キーの種別
  - `scope`: スコープ（配列）
  - `iat`: 発行日時（Unix timestamp → 人間が読める形式に変換）
  - `version`: バージョン
  - `env`: 環境

#### 3.3.2.2 JWTダウンロード機能
- ダウンロードボタンをクリックするとJWTトークンをテキストファイルとしてダウンロード
- ファイル名: `api-key-{timestamp}.txt`（例: `api-key-20250127-143045.txt`）
- ファイル内容: JWTトークンのみ（1行）
- Content-Type: `text/plain`

#### 3.3.3 キー生成パラメータ
- `iss`: `"go-webdb-template"`（固定）
- `sub`: `"public_client"`（固定）
- `type`: `"public"`（固定）
- `scope`: `["read", "write"]`（固定）
- `iat`: 現在時刻（Unix timestamp）
- `version`: 設定ファイルの`api.current_version`から取得（設定ファイルに定義）
- `env`: 現在の環境（`APP_ENV`から取得）
- `exp`: 未定義（Public APIキーは無期限）

### 3.4 クライアント側対応

#### 3.4.1 TypeScript/Next.jsクライアント修正
- `client/src/lib/api.ts`を修正
- すべてのAPIリクエストに`Authorization: Bearer <API_KEY>`ヘッダーを付与
- 環境変数（`NEXT_PUBLIC_API_KEY`）からAPIキーを取得

#### 3.4.2 APIキー設定方法
- 開発環境: `.env.local`に`NEXT_PUBLIC_API_KEY`を設定
- 本番環境: 環境変数で`NEXT_PUBLIC_API_KEY`を設定
- APIキーが設定されていない場合のエラーハンドリング

#### 3.4.3 エラーハンドリング
- APIキーが設定されていない場合: エラーを投げてリクエストを送信しない
- 401 Unauthorized: APIキーが無効または認証失敗
- 403 Forbidden: スコープが不足している場合
- 適切なエラーメッセージを表示

### 3.5 Public APIキーの保存場所

#### 3.5.1 設定ファイルへの保存
- Public APIキーは`config/{env}/config.yaml`に保存
- 設定構造:
  ```yaml
  api:
    current_version: "v1"              # 現在のバージョン（新規発行時に使用）
    public_key: "<JWT_TOKEN>"          # 発行済みのPublic APIキー（オプション）
    secret_key: "<SECRET_KEY_FOR_SIGNING>"  # JWT署名用の秘密鍵（手動設定）
    invalid_versions:                  # 無効化されたバージョンリスト
      - "v1"
  ```

#### 3.5.1.2 秘密鍵の生成方法
- コマンドラインツールで秘密鍵を生成
- 生成された秘密鍵を設定ファイル（`config/{env}/config.yaml`）に手動で記述
- 複数のWebサーバーで動作する場合、各サーバーで同じ秘密鍵を使用する必要がある
- 秘密鍵は環境別に分離（`config/{env}/`に保存）

#### 3.5.1.1 バージョン管理の仕組み
- `current_version`: 新規にAPIキーを発行する際に使用するバージョン
- `invalid_versions`: 無効化されたバージョンのリスト
- 新規キー発行時は`current_version`の値をJWTの`version`フィールドに設定
- バージョンアップ時は`current_version`を更新（例: "v1" → "v2"）
- 脆弱性発見時は該当バージョンを`invalid_versions`に追加

#### 3.5.2 環境別の扱い
- **develop環境**: Public APIキーはcommit可能（実装作業効率化のため）
- **staging環境**: Public APIキーは`.gitignore`に追加してcommit不可
- **production環境**: Public APIキーは`.gitignore`に追加してcommit不可（既に追加済み）

#### 3.5.3 .gitignore設定
- `config/staging/config.yaml`を`.gitignore`に追加
- `config/production/config.yaml`は既に`.gitignore`に追加済み

## 4. 非機能要件

### 4.1 セキュリティ要件
- JWT署名の検証を必須とする（秘密鍵による検証）
- 秘密鍵は環境別に分離（`config/{env}/`に保存）
- APIキーは設定ファイルに保存（DB管理は行わない）
- 無効バージョンリストによる迅速な無効化対応

### 4.2 パフォーマンス要件
- 認証ミドルウェアのオーバーヘッドを最小化
- JWT検証処理は高速に実行（キャッシュは不要、毎回検証）
- 認証処理がAPIレスポンス時間に大きな影響を与えない

### 4.3 エラーハンドリング
- 認証失敗時: 401 Unauthorizedを返す
- スコープ不足時: 403 Forbiddenを返す
- JWT形式が不正な場合: 401 Unauthorizedを返す
- エラーレスポンス形式: JSON形式で`{"code": 401, "message": "..."}`の形式
  - 例: `{"code": 401, "message": "Invalid API key"}`
  - 例: `{"code": 403, "message": "Insufficient scope"}`

### 4.4 既存機能への影響
- 既存のAPIエンドポイントへの影響を最小化
- 既存のCORS設定との共存
- 既存のGoAdmin管理画面との統合
- 既存の設定ファイル構造への影響を最小化

### 4.5 テスト要件
- テスト用のダミーAPIキーを`server/internal/config/testdata/develop/`に配置
- テスト実行時は`testdata/`の設定ファイルを使用
- 既存のテストが正常に動作することを確認
- 認証機能のテストを実装

## 5. 制約事項

### 5.1 技術的制約
- **JWTライブラリ**: `github.com/golang-jwt/jwt/v5`を使用
- **DB管理は行わない**: APIキーは設定ファイルベースで管理（Issue #21の要件）
- **既存のアーキテクチャ**: レイヤードアーキテクチャを維持
- **既存の設定構造**: `config.Load()`を使用した設定読み込みを維持

### 5.2 プロジェクト制約
- 既存のAPIルーター構造（`server/internal/api/router/router.go`）を維持
- 既存のGoAdmin管理画面実装（`server/internal/admin/`）との統合
- 既存のクライアント実装（`client/src/lib/api.ts`）への影響を最小化
- 既存の設定ファイル構造への影響を最小化

### 5.3 セキュリティ制約
- 秘密鍵は環境別に分離（`config/{env}/`に保存）
- staging/productionのAPIキーは`.gitignore`でcommit不可
- 秘密鍵の漏洩を防ぐため、適切な権限管理を実施

### 5.4 テスト制約
- テスト用のダミーAPIキーは`testdata/`に配置
- テスト実行時は`testdata/`の設定ファイルを使用
- 既存のテストが`t.Skipf()`でスキップできることを維持

## 6. 受け入れ基準

### 6.1 機能要件
- [ ] Public APIキーでAPIアクセス可能
- [ ] Private APIキーでAPIアクセス可能（将来実装、今回は対象外）
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

### 6.2 非機能要件
- [ ] JWT署名の検証が正常に動作する
- [ ] 認証処理がAPIレスポンス時間に大きな影響を与えない
- [ ] 認証失敗時に適切なエラーレスポンス（401/403）を返す
- [ ] エラーレスポンスが`{"code": <HTTPステータスコード>, "message": "..."}`形式である
- [ ] 既存のAPIエンドポイントが正常に動作する
- [ ] 既存のCORS設定が正常に動作する
- [ ] 既存のGoAdmin管理画面が正常に動作する

### 6.3 設定
- [ ] `config/{env}/config.yaml`にAPIキー設定が追加されている
- [ ] `api.current_version`が設定ファイルに定義されている
- [ ] `api.secret_key`が設定ファイルに定義されている
- [ ] 秘密鍵生成ツールが正常に動作する
- [ ] `api.invalid_versions`が設定ファイルで管理できる
- [ ] `config/staging/config.yaml`が`.gitignore`に追加されている
- [ ] テスト用のダミーAPIキーが`testdata/`に配置されている

### 6.4 テスト
- [ ] 既存のテストが正常に動作する
- [ ] 認証機能のテストが実装されている
- [ ] テスト用のダミーAPIキーでテストが正常に実行できる
- [ ] CI/CDでのテスト実行が正常に動作する

### 6.5 クライアント側
- [ ] TypeScriptクライアントがAPIキーを正しく送信する
- [ ] 環境変数からAPIキーを取得できる
- [ ] APIキーなしのリクエストは適切にエラーを返す
- [ ] 401/403エラー時の適切なエラーハンドリングが実装されている

## 7. 影響範囲

### 7.1 新規追加が必要なディレクトリ・ファイル

#### サーバー側
- `server/internal/auth/jwt.go`: JWT検証機能
  - `ValidateJWT`関数: JWTトークンの検証
  - `ParseJWTClaims`関数: JWTクレームのパース
  - `IsVersionInvalid`関数: バージョン無効化チェック

#### コマンドラインツール
- `server/cmd/generate-secret/main.go`: 秘密鍵生成ツール（新規作成）
  - ランダムな秘密鍵を生成
  - 生成された秘密鍵を標準出力に表示
  - 環境別の秘密鍵生成に対応
  
- `server/internal/auth/middleware.go`: 認証ミドルウェア
  - `AuthMiddleware`関数: HTTPミドルウェア実装
  - スコープ検証機能
  
- `server/internal/admin/pages/api_key.go`: GoAdminキー発行ページ
  - `APIKeyPage`関数: キー発行ページの実装（GET/POSTリクエスト処理）
  - `GeneratePublicAPIKey`関数: Public JWTキーの生成
  - `DownloadAPIKey`関数: JWTトークンのダウンロード処理
  - `DecodeJWTPayload`関数: JWTペイロードのデコード（表示用）
  - HTMLテンプレートの実装

#### データベースマイグレーション
- `db/migrations/shard1/004_api_key_menu.sql`: メニュー項目の追加
  - 「APIキー発行」メニュー項目を`goadmin_menu`テーブルに追加
  - 既存の「カスタムページ」カテゴリに追加、または新しいカテゴリを作成

#### クライアント側
- 変更のみ（新規ファイルなし）

#### テスト用
- `server/internal/config/testdata/develop/api_key.yaml`: テスト用ダミーAPIキー設定

### 7.2 変更が必要なファイル

#### サーバー側
- `server/internal/config/config.go`: APIキー設定の追加
  - `APIConfig`構造体の追加
  - `Config`構造体に`API`フィールドを追加
  
- `server/internal/api/router/router.go`: 認証ミドルウェアの適用
  - `NewRouter`関数に認証ミドルウェアを追加
  - `/api/*`パスに認証ミドルウェアを適用
  
- `server/cmd/admin/main.go`: GoAdminキー発行ページの登録
  - `RegisterCustomPages`関数にキー発行ページを追加

#### コマンドラインツール
- `server/cmd/generate-secret/main.go`: 秘密鍵生成ツール
  - ランダムな秘密鍵を生成
  - 標準出力に表示

#### 設定ファイル
- `config/develop/config.yaml`: Public APIキー設定を追加
- `config/staging/config.yaml`: Public APIキー設定を追加（.gitignoreに追加）
- `config/production/config.yaml.example`: Public APIキー設定例を追加

#### クライアント側
- `client/src/lib/api.ts`: Authorizationヘッダーの追加
  - `request`メソッドに`Authorization`ヘッダーを追加
  - 環境変数（`NEXT_PUBLIC_API_KEY`）からAPIキーを取得
  - APIキーが設定されていない場合、エラーを投げてリクエストを送信しない

#### その他
- `.gitignore`: `config/staging/config.yaml`を追加

### 7.3 変更不要なファイル
- `server/internal/api/handler/*.go`: ハンドラーは変更不要（ミドルウェアで認証）
- `server/internal/service/*.go`: サービス層は変更不要
- `server/internal/repository/*.go`: リポジトリ層は変更不要
- `server/internal/db/*.go`: データベース接続は変更不要

### 7.4 削除されるファイル
なし

## 8. 実装上の注意事項

### 8.1 JWTライブラリの使用
- **ライブラリ**: `github.com/golang-jwt/jwt/v5`を使用
- **署名アルゴリズム**: HS256（HMAC-SHA256）を使用
- **秘密鍵**: 環境別に分離（`config/{env}/config.yaml`に保存）

### 8.2 JWT検証処理
- JWT署名の検証を必須とする
- `iss`、`type`、`version`、`env`の各フィールドを検証
- 無効バージョンリストとの照合
- スコープの検証（read/write）

### 8.3 認証ミドルウェアの実装
- HTTPミドルウェアとして実装
- `/api/*`パスに適用（すべてのAPIエンドポイント）
- `/health`エンドポイントは認証不要（既存の動作を維持、`/api/*`の外側にあるため）
- 認証失敗時は401 Unauthorizedを返す（エラーレスポンス: `{"code": 401, "message": "..."}`）
- スコープ不足時は403 Forbiddenを返す（エラーレスポンス: `{"code": 403, "message": "..."}`）

### 8.4 GoAdminキー発行ページの実装
- 既存のカスタムページ実装例（`.kiro/specs/0004-goadmin/`）に従う
- `server/internal/admin/pages/api_key.go`に実装
- JWT生成処理を実装
- 生成されたJWTトークンを画面に表示
- JWTペイロードをデコードして表示（JSON整形）
- JWTダウンロード機能を実装（テキストファイルとしてダウンロード）

### 8.4.1 メニュー項目の追加
- データベースマイグレーションファイル（`db/migrations/shard1/004_api_key_menu.sql`）を作成
- `goadmin_menu`テーブルに「APIキー発行」メニュー項目を追加
- メニューの設定:
  - `parent_id`: 既存の「カスタムページ」カテゴリのID（`SELECT id FROM goadmin_menu WHERE title = 'カスタムページ'`で取得）
  - `type`: 1（メニュー項目）
  - `order`: 既存メニュー（「ユーザー登録」）の後に配置（例: 2）
  - `title`: 「APIキー発行」
  - `icon`: `fa-key`
  - `uri`: `/api-key`（GoAdminが自動的に`/admin`プレフィックスを追加し、実際のURLは`/admin/api-key`になる）
- すべてのシャード（shard1, shard2, shard3, shard4）に同じマイグレーションを適用
- 実装例:
  ```sql
  -- APIキー発行（カスタムページの子メニュー）
  INSERT INTO goadmin_menu (parent_id, type, "order", title, icon, uri, plugin_name, created_at, updated_at)
  VALUES (
      (SELECT id FROM goadmin_menu WHERE title = 'カスタムページ'),
      1, 2, 'APIキー発行', 'fa-key', '/api-key', '', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
  );
  ```

### 8.5 設定ファイルの構造
- `api`セクションを追加
- `current_version`: 現在のバージョン（新規発行時に使用、文字列、例: "v1"）
- `public_key`: Public APIキー（JWTトークン、オプション）
- `secret_key`: JWT署名用の秘密鍵（必須）
- `invalid_versions`: 無効バージョンリスト（文字列配列）

### 8.6 クライアント側の実装
- `client/src/lib/api.ts`の`request`メソッドを修正
- 環境変数（`NEXT_PUBLIC_API_KEY`）からAPIキーを取得
- すべてのAPIリクエストに`Authorization: Bearer <API_KEY>`ヘッダーを付与
- APIキーが設定されていない場合のエラーハンドリング

### 8.7 テスト用ダミーAPIキー
- `server/internal/config/testdata/develop/api_key.yaml`に配置
- テスト実行時は`testdata/`の設定ファイルを使用
- テスト用の秘密鍵を使用してJWTを生成
- 既存のテストが正常に動作することを確認

### 8.8 エラーハンドリング
- 認証失敗時: 401 Unauthorized（`{"code": 401, "message": "Invalid API key"}`）
- スコープ不足時: 403 Forbidden（`{"code": 403, "message": "Insufficient scope"}`）
- JWT形式が不正な場合: 401 Unauthorized（`{"code": 401, "message": "Invalid token format"}`）
- エラーレスポンスはJSON形式で返す（`{"code": <HTTPステータスコード>, "message": "<エラーメッセージ>"}`）

### 8.9 秘密鍵生成ツールの実装
- `server/cmd/generate-secret/main.go`を作成
- ランダムな秘密鍵（32文字以上の文字列）を生成
- 生成された秘密鍵を標準出力に表示
- 環境変数（`APP_ENV`）に対応（オプション）
- 使用例: `go run server/cmd/generate-secret/main.go`

## 9. 参考情報

### 9.1 関連Issue
- GitHub Issue #21: APIサーバーに、PublicなAPIキーによる認証を追加する

### 9.2 既存ドキュメント
- `server/internal/api/router/router.go`: APIルーター実装
- `server/internal/config/config.go`: 設定読み込み実装
- `server/internal/admin/pages/`: GoAdminカスタムページ実装例
- `client/src/lib/api.ts`: TypeScriptクライアント実装
- `config/develop/config.yaml`: 開発環境設定ファイル
- `config/staging/config.yaml`: ステージング環境設定ファイル
- `config/production/config.yaml.example`: 本番環境設定ファイル例

### 9.3 既存実装
- `server/internal/api/router/router.go`: `NewRouter`関数
- `server/internal/config/config.go`: `Load`関数、`Config`構造体
- `server/internal/admin/pages/`: カスタムページ実装例（`.kiro/specs/0004-goadmin/`）
- `client/src/lib/api.ts`: `ApiClient`クラス

### 9.4 JWTライブラリ
- **`github.com/golang-jwt/jwt/v5`**: Go言語用JWTライブラリ
- **公式ドキュメント**: https://github.com/golang-jwt/jwt
- **署名アルゴリズム**: HS256（HMAC-SHA256）を使用

### 9.5 Go言語標準ライブラリ
- `net/http`パッケージ: HTTPミドルウェア実装
- `context`パッケージ: コンテキスト管理
- `time`パッケージ: タイムスタンプ処理

