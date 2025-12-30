# CSVダウンロード機能実装タスク一覧

## 概要
CSVダウンロード機能の実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: APIサーバー全体のタイムアウト設定

#### - [ ] タスク 1.1: ServerConfigにIdleTimeoutフィールドを追加
**目的**: 設定構造体にIdleTimeoutフィールドを追加し、設定ファイルから読み込めるようにする

**作業内容**:
- `server/internal/config/config.go`の`ServerConfig`構造体に`IdleTimeout time.Duration`フィールドを追加
- `mapstructure:"idle_timeout"`タグを追加
- 既存の`ReadTimeout`と`WriteTimeout`と同じ形式で実装

**受け入れ基準**:
- `ServerConfig`構造体に`IdleTimeout time.Duration`フィールドが追加されている
- `mapstructure:"idle_timeout"`タグが設定されている
- 既存のコードスタイルに従っている

---

#### - [ ] タスク 1.2: 設定ファイルにidle_timeoutを追加
**目的**: 設定ファイルにIdleTimeoutの設定値を追加する

**作業内容**:
- `config/develop/config.yaml`の`server`セクションに`idle_timeout: 120s`を追加
- `config/production/config.yaml.example`の`server`セクションに`idle_timeout: 120s`を追加
- 既存の`read_timeout`と`write_timeout`と同じ形式で記述

**受け入れ基準**:
- `config/develop/config.yaml`に`idle_timeout: 120s`が追加されている
- `config/production/config.yaml.example`に`idle_timeout: 120s`が追加されている
- YAML形式が正しい

---

#### - [ ] タスク 1.3: サーバー起動時にIdleTimeoutを設定
**目的**: Echoサーバー起動時にIdleTimeoutを設定する

**作業内容**:
- `server/cmd/server/main.go`の`e.Server.ReadTimeout`と`e.Server.WriteTimeout`設定の後に、`e.Server.IdleTimeout = cfg.Server.IdleTimeout`を追加
- 既存のタイムアウト設定と同じ場所に配置

**受け入れ基準**:
- `server/cmd/server/main.go`に`e.Server.IdleTimeout = cfg.Server.IdleTimeout`が追加されている
- 設定が正しく読み込まれ、サーバーに適用されている
- 既存のコードスタイルに従っている

---

### Phase 2: サーバー側のCSVダウンロードエンドポイント実装

#### - [ ] タスク 2.1: CSVダウンロードエンドポイントの実装
**目的**: Huma APIにCSVダウンロードエンドポイントを登録する

**作業内容**:
- `server/internal/api/handler/dm_user_handler.go`の`RegisterDmUserEndpoints`関数内に、CSVダウンロードエンドポイントを追加
- `huma.Register`を使用して`GET /api/dm-users/csv`エンドポイントを登録
- `OperationID`: `"download-users-csv"`
- `Summary`: `"ユーザー情報をCSV形式でダウンロード"`
- `Description`: `"**Access Level:** `public` (Public API Key JWT または Auth0 JWT でアクセス可能)"`
- `Tags`: `[]string{"users"}`
- `Security`: `[]map[string][]string{{"bearerAuth": {}}}`
- 認証チェック: `auth.CheckAccessLevel(ctx, auth.AccessLevelPublic)`
- エラーハンドリング: 認証エラー時は`huma.Error403Forbidden`を返す

**受け入れ基準**:
- `GET /api/dm-users/csv`エンドポイントが登録されている
- 認証チェックが正しく実装されている
- エラーハンドリングが適切に実装されている
- 既存のエンドポイントと同じ形式で実装されている

---

#### - [ ] タスク 2.2: ユーザー情報取得処理の実装
**目的**: CSVダウンロード用にユーザー情報20件を取得する

**作業内容**:
- CSVダウンロードエンドポイントのハンドラー関数内で、`DmUserService.ListDmUsers(ctx, 20, 0)`を呼び出し
- エラーハンドリング: データ取得エラー時は`huma.Error500InternalServerError`を返す
- 取得件数: 20件（固定）
- 取得順序: `DmUserService.ListDmUsers()`の実装を確認し、必要に応じて修正（作成日時の降順）

**受け入れ基準**:
- `DmUserService.ListDmUsers(ctx, 20, 0)`が正しく呼び出されている
- エラーハンドリングが適切に実装されている
- 20件のユーザー情報が取得できること

---

#### - [ ] タスク 2.3: ストリーミングCSV生成の実装
**目的**: HumaのStreamResponseを使用してCSVをストリーミング生成する

**作業内容**:
- `huma.StreamResponse`を使用してストリーミングレスポンスを返す
- `ContentType`: `"text/csv; charset=utf-8"`
- `Headers`: `map[string]string{"Content-Disposition": `attachment; filename="dm-users.csv"`}`
- `Body`: `func(w io.Writer)`形式の関数を実装
- `http.ResponseWriter`を取得し、`http.NewResponseController`を使用してタイムアウトを設定
- タイムアウト時間: 3分（`time.Now().Add(3 * time.Minute)`）
- `encoding/csv`パッケージの`csv.NewWriter`を使用してCSVエンコーダーを作成
- ヘッダー行を書き込み: `["ID", "Name", "Email", "Created At", "Updated At"]`
- ユーザー情報を1件ずつCSV行として書き込み:
  - `user.ID`
  - `user.Name`
  - `user.Email`
  - `user.CreatedAt.Format(time.RFC3339)`
  - `user.UpdatedAt.Format(time.RFC3339)`
- `csvWriter.Flush()`を`defer`で実行

**受け入れ基準**:
- `huma.StreamResponse`が正しく実装されている
- Content-TypeとContent-Dispositionが正しく設定されている
- タイムアウト設定が3分に設定されている
- CSVヘッダー行が正しく書き込まれている
- ユーザーデータが正しくCSV行として書き込まれている
- 日時フォーマットがISO 8601形式（RFC3339）であること

---

#### - [ ] タスク 2.4: ストリーミングCSV生成のエラーハンドリング
**目的**: CSV生成時のエラーを適切に処理する

**作業内容**:
- CSV書き込み時のエラーをログに記録（ストリーミング開始後はHTTPステータスコードを変更できないため）
- `csvWriter.Write()`のエラーをチェック
- エラー発生時はログに記録し、処理を中断

**受け入れ基準**:
- CSV書き込み時のエラーがログに記録される
- エラー発生時に処理が適切に中断される

---

#### - [ ] タスク 2.5: サーバー側の単体テスト実装
**目的**: CSVダウンロードエンドポイントの単体テストを作成する

**作業内容**:
- `server/internal/api/handler/dm_user_handler_test.go`にCSVダウンロードエンドポイントのテストを追加
- テストケース:
  1. 正常系: CSVダウンロードが正常に動作する
     - ユーザー情報20件を取得
     - CSV形式で正しく出力される
     - Content-TypeとContent-Dispositionが正しく設定される
  2. 認証エラー: 認証が失敗した場合
     - 403 Forbiddenが返される
  3. データベースエラー: データ取得が失敗した場合
     - 500 Internal Server Errorが返される
  4. 空データ: ユーザーが0件の場合
     - ヘッダー行のみのCSVが返される

**受け入れ基準**:
- `dm_user_handler_test.go`にCSVダウンロードエンドポイントのテストが追加されている
- 上記のテストケースが全て実装されている
- テストが全て通過する

---

### Phase 3: クライアント側の実装

#### - [ ] タスク 3.1: APIクライアントにdownloadUsersCSVメソッドを追加
**目的**: APIクライアントにCSVダウンロード用のメソッドを追加する

**作業内容**:
- `client/src/lib/api.ts`の`ApiClient`クラスに`downloadUsersCSV(): Promise<void>`メソッドを追加
- 実装内容:
  - `fetch()`を使用して`/api/dm-users/csv`エンドポイントにGETリクエストを送信
  - `Authorization`ヘッダーに`Bearer ${token}`を設定（`this.apiKey`を使用）
  - レスポンスのステータスをチェック
  - 401/403エラー: JSONレスポンスからエラーメッセージを取得してスロー
  - その他のエラー: テキストレスポンスからエラーメッセージを取得してスロー
  - レスポンスを`blob()`で取得
  - Content-Dispositionヘッダーからファイル名を取得（デフォルト: `dm-users.csv`）
  - `URL.createObjectURL()`でBlob URLを生成
  - `<a>`要素を作成し、`download`属性を設定
  - プログラム的にクリックイベントを発火してダウンロードを開始
  - ダウンロード後、`document.body.removeChild(link)`で要素を削除
  - `URL.revokeObjectURL()`でBlob URLを解放

**受け入れ基準**:
- `downloadUsersCSV()`メソッドが実装されている
- エラーハンドリングが適切に実装されている
- Blob URLの生成と解放が正しく実装されている
- ダウンロードが正常に実行される

---

#### - [ ] タスク 3.2: CSVダウンロードボタンの追加
**目的**: dm-usersページにCSVダウンロードボタンを追加する

**作業内容**:
- `client/src/app/dm-users/page.tsx`にCSVダウンロードボタンを追加
- 状態管理:
  - `downloading`: ダウンロード中の状態を管理（`useState<boolean>`）
  - `downloadError`: エラーメッセージを管理（`useState<string | null>`）
- ボタンの実装:
  - `handleDownloadCSV`関数を作成
  - `apiClient.downloadUsersCSV()`を呼び出し
  - エラーハンドリング: エラー発生時は`setDownloadError`でエラーメッセージを設定
  - ローディング状態: `setDownloading(true/false)`で管理
- ボタンの配置: ユーザー一覧セクションの上部（「ユーザー一覧」見出しの下）
- ボタンのスタイル: 既存のボタンスタイルに合わせる（例: `bg-green-500 text-white rounded hover:bg-green-600`）
- エラーメッセージの表示: エラー発生時はエラーメッセージを表示（既存のエラー表示と同じ形式）

**受け入れ基準**:
- CSVダウンロードボタンが表示される
- ボタンをクリックするとCSVダウンロードが開始される
- ダウンロード中はローディング状態が表示される
- エラー発生時はエラーメッセージが表示される
- 既存のUIデザインと一貫性がある

---

#### - [ ] タスク 3.3: クライアント側の単体テスト実装
**目的**: APIクライアントのCSVダウンロード機能の単体テストを作成する

**作業内容**:
- `client/src/lib/__tests__/api.test.ts`に`downloadUsersCSV()`メソッドのテストを追加
- テストケース:
  1. 正常系: CSVダウンロードが正常に動作する
     - `fetch()`が正常に呼び出される
     - Blobが正しく生成される
     - ダウンロードが実行される
  2. エラーハンドリング: エラーが発生した場合
     - エラーが正しくスローされる
  3. ファイル名の取得: Content-Dispositionヘッダーからファイル名を取得できること

**受け入れ基準**:
- `api.test.ts`に`downloadUsersCSV()`メソッドのテストが追加されている
- 上記のテストケースが全て実装されている
- テストが全て通過する

---

### Phase 4: 統合テストと動作確認

#### - [ ] タスク 4.1: E2Eテストの実装
**目的**: クライアント側からAPIサーバー側までのE2Eテストを実装する

**作業内容**:
- `client/e2e/csv-download.spec.ts`を作成
- テストケース:
  1. 正常系: CSVダウンロードが正常に動作する
     - dm-usersページにアクセス
     - CSVダウンロードボタンをクリック
     - CSVファイルがダウンロードされる
     - CSVファイルの内容を検証（ヘッダー行、データ行）
  2. エラーハンドリング: エラーが発生した場合
     - エラーメッセージが表示される

**受け入れ基準**:
- `csv-download.spec.ts`が作成されている
- 上記のテストケースが全て実装されている
- テストが全て通過する

---

#### - [ ] タスク 4.2: 動作確認
**目的**: 実装した機能が正しく動作することを確認する

**作業内容**:
- 開発サーバーを起動
- ブラウザでdm-usersページにアクセス
- CSVダウンロードボタンをクリック
- CSVファイルがダウンロードされることを確認
- CSVファイルの内容を確認（ExcelやGoogleスプレッドシートで開く）
- データが正しく表示されることを確認
- エラーハンドリングが正しく動作することを確認

**受け入れ基準**:
- CSVダウンロードボタンが表示される
- ボタンをクリックするとCSVファイルがダウンロードされる
- CSVファイルの内容が正しい（ヘッダー行、データ行、エンコーディング）
- CSVファイルが一般的なCSVリーダーで正しく読み込める

---

#### - [ ] タスク 4.3: 既存テストの確認
**目的**: 既存のテストが全て成功することを確認する

**作業内容**:
- サーバー側の既存テストを実行（`go test ./...`）
- クライアント側の既存テストを実行（`npm test`）
- E2Eテストを実行（`npm run test:e2e`）
- 全てのテストが成功することを確認

**受け入れ基準**:
- サーバー側の既存テストが全て成功する
- クライアント側の既存テストが全て成功する
- E2Eテストが全て成功する

---

## 実装時の注意事項

### タイムアウト設定
- `http.NewResponseController`はGo 1.20以降で利用可能。プロジェクトのGoバージョンを確認する必要がある
- CSVダウンロード用の個別タイムアウト設定（3分）により、既存のWriteTimeout（30秒）を上書きできる

### ストリーミング中のエラー処理
- ストリーミング開始後はHTTPステータスコードを変更できない
- エラーはログに記録するのみ（ユーザーへの通知は困難）

### CSV形式
- BOM（Byte Order Mark）は追加しない
- 標準的なCSV形式（UTF-8 without BOM）を使用する
- 日時フォーマットはISO 8601形式（RFC3339）を使用する

### データ取得の順序
- `DmUserService.ListDmUsers()`の実装を確認し、必要に応じて修正（作成日時の降順）

### 既存機能の維持
- 既存のdm-users API機能に影響を与えない
- 既存のコードスタイルに従う
