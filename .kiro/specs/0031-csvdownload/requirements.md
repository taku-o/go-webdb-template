# CSVダウンロード機能要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #62
- **Issueタイトル**: CSVダウンロード機能を用意する
- **Feature名**: 0031-csvdownload
- **作成日**: 2025-01-27

### 1.2 目的
サンプル実装としてCSVダウンロード機能を用意する。ファイルを作成してからダウンロードするのではなく、ストリームで結果を返す方式を採用し、メモリ効率の良い実装を実現する。

### 1.3 スコープ
- クライアント側: `/dm-users/csv` エンドポイントの実装
- APIサーバー側: `/api/dm-users/csv` エンドポイントの実装（public API）
- ユーザー情報20件をCSV形式でダウンロード
- ストリーミング方式によるCSV生成とダウンロード

**本実装の範囲外**:
- 他のエンティティ（dm_postsなど）のCSVダウンロード機能
- CSVファイルのアップロード機能
- 大量データのCSVダウンロード（ページネーション、フィルタリングなど）
- CSV形式のカスタマイズ（区切り文字、エンコーディングなど）

## 2. 背景・現状分析

### 2.1 現在の実装
- **dm-usersページ**: `client/src/app/dm-users/page.tsx`にユーザー一覧表示機能が実装されている
- **Get Todayボタン**: `client/src/app/page.tsx`にTodayApiButtonコンポーネントが配置されている
- **APIクライアント**: `client/src/lib/api.ts`にApiClientクラスが実装されている
- **dm-users API**: `server/internal/api/handler/dm_user_handler.go`にユーザーCRUD APIが実装されている
- **アクセスレベル**: public APIとして実装されており、Public API Key JWTまたはAuth0 JWTでアクセス可能

### 2.2 課題点
1. **CSVダウンロード機能の不在**: ユーザー情報をCSV形式でダウンロードする機能が存在しない
2. **データエクスポート手段の不足**: ユーザーがデータを外部で利用するためのエクスポート機能がない

### 2.3 本実装による改善点
1. **データエクスポート機能の提供**: ユーザー情報をCSV形式でダウンロードできるようになる
2. **ストリーミング方式の採用**: メモリ効率の良い実装により、大量データにも対応可能な基盤を構築
3. **サンプル実装の提供**: 他のエンティティのCSVダウンロード機能実装時の参考となる

## 3. 機能要件

### 3.1 クライアント側の実装

#### 3.1.1 CSVダウンロードボタンの配置
- **ファイル**: `client/src/app/dm-users/page.tsx`
- **配置位置**: Get Todayボタンの下辺り（`client/src/app/page.tsx`のTodayApiButtonコンポーネントの下）
- **実装内容**:
  - ダウンロードボタンを追加
  - ボタンをクリックするとCSVダウンロードを開始
  - ダウンロード中はローディング状態を表示
  - エラー発生時はエラーメッセージを表示

#### 3.1.2 APIクライアントの拡張
- **ファイル**: `client/src/lib/api.ts`
- **実装内容**:
  - `downloadUsersCSV()`メソッドを追加
  - `/api/dm-users/csv`エンドポイントにGETリクエストを送信
  - レスポンスをBlobとして取得
  - Content-Dispositionヘッダーからファイル名を取得（またはデフォルト名を使用）
  - Blobをダウンロードリンクとして生成し、自動ダウンロードを実行

#### 3.1.3 ダウンロード処理の実装
- **実装内容**:
  - `fetch()`を使用してAPIを呼び出し
  - レスポンスを`blob()`で取得
  - `URL.createObjectURL()`でBlob URLを生成
  - `<a>`要素を作成し、`download`属性を設定
  - プログラム的にクリックイベントを発火してダウンロードを開始
  - ダウンロード後、Blob URLを`URL.revokeObjectURL()`で解放

### 3.2 APIサーバー側の実装

#### 3.2.1 CSVダウンロードエンドポイントの実装
- **ファイル**: `server/internal/api/handler/dm_user_handler.go`
- **エンドポイント**: `GET /api/dm-users/csv`
- **アクセスレベル**: `public` (Public API Key JWT または Auth0 JWT でアクセス可能)
- **実装内容**:
  - Huma APIにエンドポイントを登録
  - ユーザー情報20件を取得
  - CSV形式でストリーミングレスポンスを返す
  - Content-Type: `text/csv; charset=utf-8`
  - Content-Disposition: `attachment; filename="dm-users.csv"`

#### 3.2.2 ストリーミングCSV生成の実装
- **実装方式**: Humaの`StreamResponse`を使用して、ファイルを作成せずにHTTPレスポンスに直接ストリーミング
- **実装内容**:
  - Humaの`huma.StreamResponse`を使用してストリーミングレスポンスを返す
  - `StreamResponse.Body`に`func(w io.Writer)`形式の関数を設定
  - `http.ResponseWriter`を取得し、`http.NewResponseController`を使用してタイムアウトを設定
  - タイムアウト時間: 3分（`time.Now().Add(3 * time.Minute)`）
  - `encoding/csv`パッケージを使用してCSVエンコーダーを作成
  - HTTPレスポンスライターに直接書き込み
  - ヘッダー行を書き込み（ID, Name, Email, Created At, Updated At）
  - ユーザー情報を1件ずつ取得し、CSV行として書き込み
  - エラー発生時は適切なエラーレスポンスを返す
- **実装コード例**:
  ```go
  func(ctx context.Context, input *struct{}) (*huma.StreamResponse, error) {
      return &huma.StreamResponse{
          Body: func(w io.Writer) {
              // http.ResponseWriterを取り出す
              rw, ok := w.(http.ResponseWriter)
              if ok {
                  rc := http.NewResponseController(rw)
                  rc.SetWriteDeadline(time.Now().Add(3 * time.Minute))
              }

              // CSV書き込み処理...
          },
      }, nil
  }
  ```

#### 3.2.3 ユーザー情報の取得
- **実装内容**:
  - `DmUserService`を使用してユーザー情報を取得
  - 取得件数: 20件（固定）
  - 取得順序: 作成日時の降順（最新の20件）
  - シャーディング対応: 全シャードからデータを取得（`GetAll`メソッドを使用）

#### 3.2.4 CSV形式の仕様
- **文字エンコーディング**: UTF-8
- **区切り文字**: カンマ（`,`）
- **改行文字**: LF（`\n`）
- **ヘッダー行**: 必須（ID, Name, Email, Created At, Updated At）
- **日時フォーマット**: ISO 8601形式（`2006-01-02T15:04:05Z07:00`）
- **特殊文字のエスケープ**: `encoding/csv`パッケージの標準的なエスケープ処理に従う

### 3.3 エラーハンドリング

#### 3.3.1 クライアント側のエラーハンドリング
- **実装内容**:
  - ネットワークエラー: エラーメッセージを表示
  - HTTPエラー（4xx, 5xx）: エラーレスポンスのメッセージを表示
  - ダウンロード失敗: ユーザーに分かりやすいエラーメッセージを表示

#### 3.3.2 サーバー側のエラーハンドリング
- **実装内容**:
  - データベースエラー: 500 Internal Server Errorを返す
  - 認証エラー: 403 Forbiddenを返す
  - CSV生成エラー: 500 Internal Server Errorを返す
  - ストリーミング中のエラー: 適切なエラーレスポンスを返す

## 4. 非機能要件

### 4.1 パフォーマンス要件
- **レスポンス時間**: 20件のデータ取得とCSV生成は1秒以内に完了すること
- **メモリ使用量**: ストリーミング方式により、メモリ使用量を最小限に抑えること
- **同時リクエスト**: 複数の同時リクエストに対応できること
- **タイムアウト設定**: ストリーミングレスポンスの書き込みタイムアウトを3分に設定すること

### 4.2 セキュリティ要件
- **認証**: Public API Key JWTまたはAuth0 JWTによる認証が必要
- **認可**: publicアクセスレベルでアクセス可能
- **入力検証**: エンドポイントはパラメータを受け取らないため、入力検証は不要

### 4.3 互換性要件
- **ブラウザ対応**: モダンブラウザ（Chrome, Firefox, Safari, Edge）で動作すること
- **文字エンコーディング**: UTF-8でエンコードされたCSVファイルを生成すること
- **CSV形式**: 一般的なCSVリーダー（Excel、Googleスプレッドシートなど）で読み込めること

### 4.4 保守性要件
- **コード品質**: 既存のコードスタイルに従うこと
- **テスト**: 単体テストとE2Eテストを実装すること
- **ドキュメント**: APIドキュメント（OpenAPI仕様）にエンドポイントを追加すること

## 5. 受け入れ基準

### 5.1 機能的な受け入れ基準
1. **CSVダウンロードボタンの表示**: dm-usersページにCSVダウンロードボタンが表示されること
2. **CSVダウンロードの実行**: ボタンをクリックすると、ユーザー情報20件がCSV形式でダウンロードされること
3. **CSV形式の正確性**: ダウンロードされたCSVファイルが正しい形式（ヘッダー行、データ行、エンコーディング）であること
4. **データの正確性**: CSVファイルに含まれるデータが、データベースのユーザー情報と一致すること
5. **ストリーミング動作**: ファイルを作成せずに、ストリーミング方式でCSVを生成すること
6. **タイムアウト設定**: ストリーミングレスポンスの書き込みタイムアウトが3分に設定されていること
7. **エラーハンドリング**: エラー発生時に適切なエラーメッセージが表示されること

### 5.2 非機能的な受け入れ基準
1. **パフォーマンス**: 20件のデータ取得とCSV生成が1秒以内に完了すること
2. **メモリ効率**: ストリーミング方式により、メモリ使用量が適切に管理されること
3. **セキュリティ**: Public API Key JWTまたはAuth0 JWTによる認証が正しく機能すること
4. **互換性**: ダウンロードされたCSVファイルが、一般的なCSVリーダーで正しく読み込めること

### 5.3 テスト要件
1. **単体テスト**: APIハンドラーの単体テストを実装すること
2. **E2Eテスト**: クライアント側からAPIサーバー側までのE2Eテストを実装すること
3. **既存テストの維持**: 既存のテストが全て成功すること

## 6. 制約事項

### 6.1 技術的制約
- **フレームワーク**: 既存のHuma API、Echo、Next.jsを使用すること
- **データベース**: 既存のシャーディング構成を維持すること
- **認証**: 既存の認証・認可システムを使用すること

### 6.2 実装上の制約
- **ファイル作成禁止**: 一時ファイルを作成せず、ストリーミング方式で実装すること
- **データ件数**: 20件固定（ページネーションやフィルタリングは実装しない）
- **CSV形式**: 標準的なCSV形式（カンマ区切り、UTF-8エンコーディング）を使用すること

## 7. 関連ドキュメント

### 7.1 既存の実装参考
- **クライアント側APIクライアント**: `client/src/lib/api.ts`
- **サーバー側APIハンドラー**: `server/internal/api/handler/dm_user_handler.go`
- **Today API実装**: `server/internal/api/handler/today_handler.go`
- **dm-usersページ**: `client/src/app/dm-users/page.tsx`
- **TodayApiButtonコンポーネント**: `client/src/components/TodayApiButton.tsx`

### 7.2 技術ドキュメント
- **Huma API**: https://huma.dev/
- **Huma StreamResponse**: https://pkg.go.dev/github.com/danielgtaylor/huma/v2#StreamResponse
- **Go encoding/csv**: https://pkg.go.dev/encoding/csv
- **Go http.ResponseController**: https://pkg.go.dev/net/http#ResponseController
- **Next.js**: https://nextjs.org/docs
- **Fetch API**: https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API
