# APIドキュメント改善実装タスク一覧

## 概要
APIドキュメント改善機能の実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: SecurityScheme設定の追加

#### - [ ] タスク 1.1: router.goにSecurityScheme設定を追加
**目的**: Huma ConfigにBearer認証スキームを定義

**作業内容**:
- `server/internal/api/router/router.go`の`humaConfig`設定部分を修正
- `humaConfig.Components`に`SecuritySchemes`を追加
- SecurityScheme名は"bearerAuth"を使用
- SecuritySchemeの設定内容:
  - Type: "http"
  - Scheme: "bearer"
  - BearerFormat: "JWT"

**受け入れ基準**:
- `humaConfig.Components.SecuritySchemes`に"bearerAuth"が定義されている
- SecuritySchemeの設定内容が正しい（Type: "http", Scheme: "bearer", BearerFormat: "JWT"）
- 既存の設定が上書きされていない

---

### Phase 2: UserエンドポイントへのSecurityとTags追加

#### - [ ] タスク 2.1: POST /api/usersにSecurityとTags追加
**目的**: ユーザー作成エンドポイントにSecurityプロパティとアクセスレベルTagを追加

**作業内容**:
- `server/internal/api/handler/user_handler.go`の`create-user`エンドポイントを修正
- `Security`プロパティを追加: `[]map[string][]string{{"bearerAuth": []}}`
- `Tags`に"Public API"を追加: `[]string{"users", "Public API"}`

**受け入れ基準**:
- `Security`プロパティが追加されている
- `Tags`に"Public API"が追加されている
- 既存の"users"タグが維持されている
- 既存の処理ロジックが変更されていない

---

#### - [ ] タスク 2.2: GET /api/users/{id}にSecurityとTags追加
**目的**: ユーザー取得エンドポイントにSecurityプロパティとアクセスレベルTagを追加

**作業内容**:
- `server/internal/api/handler/user_handler.go`の`get-user`エンドポイントを修正
- `Security`プロパティを追加: `[]map[string][]string{{"bearerAuth": []}}`
- `Tags`に"Public API"を追加: `[]string{"users", "Public API"}`

**受け入れ基準**:
- `Security`プロパティが追加されている
- `Tags`に"Public API"が追加されている
- 既存の"users"タグが維持されている
- 既存の処理ロジックが変更されていない

---

#### - [ ] タスク 2.3: GET /api/usersにSecurityとTags追加
**目的**: ユーザー一覧取得エンドポイントにSecurityプロパティとアクセスレベルTagを追加

**作業内容**:
- `server/internal/api/handler/user_handler.go`の`list-users`エンドポイントを修正
- `Security`プロパティを追加: `[]map[string][]string{{"bearerAuth": []}}`
- `Tags`に"Public API"を追加: `[]string{"users", "Public API"}`

**受け入れ基準**:
- `Security`プロパティが追加されている
- `Tags`に"Public API"が追加されている
- 既存の"users"タグが維持されている
- 既存の処理ロジックが変更されていない

---

#### - [ ] タスク 2.4: PUT /api/users/{id}にSecurityとTags追加
**目的**: ユーザー更新エンドポイントにSecurityプロパティとアクセスレベルTagを追加

**作業内容**:
- `server/internal/api/handler/user_handler.go`の`update-user`エンドポイントを修正
- `Security`プロパティを追加: `[]map[string][]string{{"bearerAuth": []}}`
- `Tags`に"Public API"を追加: `[]string{"users", "Public API"}`

**受け入れ基準**:
- `Security`プロパティが追加されている
- `Tags`に"Public API"が追加されている
- 既存の"users"タグが維持されている
- 既存の処理ロジックが変更されていない

---

#### - [ ] タスク 2.5: DELETE /api/users/{id}にSecurityとTags追加
**目的**: ユーザー削除エンドポイントにSecurityプロパティとアクセスレベルTagを追加

**作業内容**:
- `server/internal/api/handler/user_handler.go`の`delete-user`エンドポイントを修正
- `Security`プロパティを追加: `[]map[string][]string{{"bearerAuth": []}}`
- `Tags`に"Public API"を追加: `[]string{"users", "Public API"}`

**受け入れ基準**:
- `Security`プロパティが追加されている
- `Tags`に"Public API"が追加されている
- 既存の"users"タグが維持されている
- 既存の処理ロジックが変更されていない

---

### Phase 3: PostエンドポイントへのSecurityとTags追加

#### - [ ] タスク 3.1: POST /api/postsにSecurityとTags追加
**目的**: 投稿作成エンドポイントにSecurityプロパティとアクセスレベルTagを追加

**作業内容**:
- `server/internal/api/handler/post_handler.go`の`create-post`エンドポイントを修正
- `Security`プロパティを追加: `[]map[string][]string{{"bearerAuth": []}}`
- `Tags`に"Public API"を追加: `[]string{"posts", "Public API"}`

**受け入れ基準**:
- `Security`プロパティが追加されている
- `Tags`に"Public API"が追加されている
- 既存の"posts"タグが維持されている
- 既存の処理ロジックが変更されていない

---

#### - [ ] タスク 3.2: GET /api/posts/{id}にSecurityとTags追加
**目的**: 投稿取得エンドポイントにSecurityプロパティとアクセスレベルTagを追加

**作業内容**:
- `server/internal/api/handler/post_handler.go`の`get-post`エンドポイントを修正
- `Security`プロパティを追加: `[]map[string][]string{{"bearerAuth": []}}`
- `Tags`に"Public API"を追加: `[]string{"posts", "Public API"}`

**受け入れ基準**:
- `Security`プロパティが追加されている
- `Tags`に"Public API"が追加されている
- 既存の"posts"タグが維持されている
- 既存の処理ロジックが変更されていない

---

#### - [ ] タスク 3.3: GET /api/postsにSecurityとTags追加
**目的**: 投稿一覧取得エンドポイントにSecurityプロパティとアクセスレベルTagを追加

**作業内容**:
- `server/internal/api/handler/post_handler.go`の`list-posts`エンドポイントを修正
- `Security`プロパティを追加: `[]map[string][]string{{"bearerAuth": []}}`
- `Tags`に"Public API"を追加: `[]string{"posts", "Public API"}`

**受け入れ基準**:
- `Security`プロパティが追加されている
- `Tags`に"Public API"が追加されている
- 既存の"posts"タグが維持されている
- 既存の処理ロジックが変更されていない

---

#### - [ ] タスク 3.4: PUT /api/posts/{id}にSecurityとTags追加
**目的**: 投稿更新エンドポイントにSecurityプロパティとアクセスレベルTagを追加

**作業内容**:
- `server/internal/api/handler/post_handler.go`の`update-post`エンドポイントを修正
- `Security`プロパティを追加: `[]map[string][]string{{"bearerAuth": []}}`
- `Tags`に"Public API"を追加: `[]string{"posts", "Public API"}`

**受け入れ基準**:
- `Security`プロパティが追加されている
- `Tags`に"Public API"が追加されている
- 既存の"posts"タグが維持されている
- 既存の処理ロジックが変更されていない

---

#### - [ ] タスク 3.5: DELETE /api/posts/{id}にSecurityとTags追加
**目的**: 投稿削除エンドポイントにSecurityプロパティとアクセスレベルTagを追加

**作業内容**:
- `server/internal/api/handler/post_handler.go`の`delete-post`エンドポイントを修正
- `Security`プロパティを追加: `[]map[string][]string{{"bearerAuth": []}}`
- `Tags`に"Public API"を追加: `[]string{"posts", "Public API"}`

**受け入れ基準**:
- `Security`プロパティが追加されている
- `Tags`に"Public API"が追加されている
- 既存の"posts"タグが維持されている
- 既存の処理ロジックが変更されていない

---

#### - [ ] タスク 3.6: GET /api/user-postsにSecurityとTags追加
**目的**: ユーザーと投稿のJOIN結果取得エンドポイントにSecurityプロパティとアクセスレベルTagを追加

**作業内容**:
- `server/internal/api/handler/post_handler.go`の`get-user-posts`エンドポイントを修正
- `Security`プロパティを追加: `[]map[string][]string{{"bearerAuth": []}}`
- `Tags`に"Public API"を追加: `[]string{"posts", "Public API"}`

**受け入れ基準**:
- `Security`プロパティが追加されている
- `Tags`に"Public API"が追加されている
- 既存の"posts"タグが維持されている
- 既存の処理ロジックが変更されていない

---

### Phase 4: TodayエンドポイントへのSecurityとTags追加

#### - [ ] タスク 4.1: GET /api/todayにSecurityとTags追加
**目的**: 今日の日付取得エンドポイントにSecurityプロパティとアクセスレベルTagを追加

**作業内容**:
- `server/internal/api/handler/today_handler.go`の`get-today`エンドポイントを修正
- `Security`プロパティを追加: `[]map[string][]string{{"bearerAuth": []}}`
- `Tags`に"Private API"を追加: `[]string{"today", "Private API"}`

**受け入れ基準**:
- `Security`プロパティが追加されている
- `Tags`に"Private API"が追加されている
- 既存の"today"タグが維持されている
- 既存の処理ロジックが変更されていない

---

### Phase 5: 動作確認とテスト

#### - [ ] タスク 5.1: サーバー起動確認
**目的**: サーバーが正常に起動することを確認

**作業内容**:
- サーバーを起動
- エラーが発生しないことを確認
- ログにエラーが出力されていないことを確認

**受け入れ基準**:
- サーバーが正常に起動する
- エラーが発生しない
- ログにエラーが出力されていない

---

#### - [ ] タスク 5.2: APIドキュメントの表示確認
**目的**: Swagger UIでAPIドキュメントが正しく表示されることを確認

**作業内容**:
- `http://localhost:8080/docs`にアクセス
- Swagger UIが正常に表示されることを確認
- すべてのエンドポイントが表示されることを確認

**受け入れ基準**:
- Swagger UIが正常に表示される
- すべてのエンドポイント（12個）が表示される
- エラーが発生しない

---

#### - [ ] タスク 5.3: SecurityScheme設定の確認
**目的**: SecuritySchemeが正しく設定されていることを確認

**作業内容**:
- Swagger UIでSecuritySchemeが定義されていることを確認
- "Authorize"ボタンが表示されることを確認
- "Authorize"ボタンをクリックしてトークン入力欄が表示されることを確認

**受け入れ基準**:
- SecurityScheme（bearerAuth）が定義されている
- "Authorize"ボタンが表示される
- トークン入力欄が表示される

---

#### - [ ] タスク 5.4: Request Sampleの確認
**目的**: Request Sample（curl）にAuthorization Headerが表示されることを確認

**作業内容**:
- 各エンドポイントのRequest Sampleを確認
- curlコマンドに`-H 'Authorization: Bearer <token>'`が含まれていることを確認
- すべてのエンドポイント（12個）で確認

**受け入れ基準**:
- すべてのエンドポイントのRequest Sampleに`-H 'Authorization: Bearer <token>'`が含まれている
- curlコマンドの形式が正しい

---

#### - [ ] タスク 5.5: Tags表示の確認
**目的**: アクセスレベルTagが正しく表示されることを確認

**作業内容**:
- 各エンドポイントのTagsを確認
- Userエンドポイント（5つ）に"Public API"タグが表示されることを確認
- Postエンドポイント（6つ）に"Public API"タグが表示されることを確認
- Todayエンドポイント（1つ）に"Private API"タグが表示されることを確認
- 既存の機能タグ（"users", "posts", "today"）も表示されることを確認

**受け入れ基準**:
- すべてのエンドポイントに適切なアクセスレベルTagが表示される
- 既存の機能タグも表示される
- Tag名が正確である（大文字小文字を含む）

---

#### - [ ] タスク 5.6: 既存API動作の確認
**目的**: 既存のAPI動作に影響がないことを確認

**作業内容**:
- 既存のAPIクライアントでAPIを呼び出す
- 認証が正常に動作することを確認
- アクセス制御が正常に動作することを確認
- レスポンスが正常に返ることを確認

**受け入れ基準**:
- 既存のAPI動作に影響がない
- 認証が正常に動作する
- アクセス制御が正常に動作する
- レスポンスが正常に返る

---

#### - [ ] タスク 5.7: Swagger UIでのAPIテスト
**目的**: Swagger UIから直接APIをテストできることを確認

**作業内容**:
- "Authorize"ボタンからトークンを入力
- 各エンドポイントをSwagger UIから直接テスト
- レスポンスが正常に返ることを確認

**受け入れ基準**:
- "Authorize"ボタンからトークンを入力できる
- Swagger UIからAPIをテストできる
- レスポンスが正常に返る

---

## 実装上の注意事項

### SecurityScheme設定
- `huma.DefaultConfig()`で作成したConfigに`Components`を追加する際は、既存の設定を上書きしないよう注意
- SecurityScheme名は"bearerAuth"を使用（Issue #37の記載に従う）

### Securityプロパティの適用
- すべてのエンドポイントに一貫して適用する
- Securityプロパティの形式は`[]map[string][]string{{"bearerAuth": []}}`とする
- 空の配列`[]`を指定することで、すべてのスコープを許可する

### Tagsによるアクセスレベル表示
- 既存の機能タグ（"users", "posts", "today"）は維持し、アクセスレベルTagを追加する
- Tag名は大文字小文字を区別するため、"Public API"と"Private API"を正確に記述する
- 各エンドポイントのアクセスレベルは、既存の`auth.CheckAccessLevel()`の呼び出しから判断する

### 既存機能との互換性
- 認証ロジックやAPI動作は一切変更しない
- 既存の`auth.CheckAccessLevel()`によるアクセス制御は維持
- ドキュメント表示のみを改善する

## 受け入れ基準（全体）

### SecurityScheme設定
- [ ] `huma.Config`の`Components.SecuritySchemes`にBearer認証スキームが定義されている
- [ ] SecuritySchemeの設定内容が正しい（Type: "http", Scheme: "bearer", BearerFormat: "JWT"）

### 各エンドポイントへのSecurity適用
- [ ] すべてのエンドポイント（12個）に`Security`プロパティが追加されている
- [ ] Securityプロパティの形式が正しい（`[]map[string][]string{{"bearerAuth": []}}`）

### APIアクセスレベルの表示
- [ ] すべてのエンドポイントに適切なアクセスレベルTagが追加されている
  - Userエンドポイント（5つ）："Public API"
  - Postエンドポイント（6つ）："Public API"
  - Todayエンドポイント（1つ）："Private API"
- [ ] 既存の機能タグ（"users", "posts", "today"）が維持されている

### APIドキュメントの表示確認
- [ ] `http://localhost:8080/docs`にアクセスして、Request Sample（curl）に`-H 'Authorization: Bearer <token>'`が表示される
- [ ] Swagger UIに「Authorize」ボタンが表示される
- [ ] 「Authorize」ボタンをクリックしてトークンを入力できる
- [ ] 各APIにアクセスレベルのTag（"Public API" / "Private API"）が表示される
- [ ] 既存の機能タグ（"users", "posts", "today"）も表示される

### 既存機能の動作確認
- [ ] 既存のAPI動作に影響がないことを確認（認証、アクセス制御が正常に動作）
- [ ] 既存のAPIクライアントへの影響がないことを確認
