# メール送信機能実装タスク一覧

## 概要
メール送信機能の実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: 設定とインフラの準備

#### - [ ] タスク 1.1: EmailConfig構造体の追加
**目的**: メール送信機能の設定を管理する構造体を追加する

**作業内容**:
- `server/internal/config/config.go`の`Config`構造体に`Email EmailConfig`フィールドを追加
- `mapstructure:"email"`タグを追加
- `EmailConfig`構造体を定義:
  - `SenderType string` - 送信方式（"mock", "mailpit", "ses"）
  - `Mock MockConfig` - MockSender設定
  - `Mailpit MailpitConfig` - MailpitSender設定
  - `SES SESConfig` - SESSender設定
- `MockConfig`構造体を定義（現時点では空）
- `MailpitConfig`構造体を定義:
  - `SMTPHost string` - SMTPホスト（デフォルト: "localhost"）
  - `SMTPPort int` - SMTPポート（デフォルト: 1025）
- `SESConfig`構造体を定義:
  - `From string` - 送信元メールアドレス
  - `Region string` - AWSリージョン

**受け入れ基準**:
- `Config`構造体に`Email EmailConfig`フィールドが追加されている
- `EmailConfig`、`MockConfig`、`MailpitConfig`、`SESConfig`構造体が定義されている
- 既存のコードスタイルに従っている

---

#### - [ ] タスク 1.2: 設定ファイルにメール送信設定を追加
**目的**: 各環境の設定ファイルにメール送信設定を追加する

**作業内容**:
- `config/develop/config.yaml`の`email`セクションを追加:
  - `sender_type: "mock"`
  - `mailpit.smtp_host: "localhost"`
  - `mailpit.smtp_port: 1025`
  - `ses.from: "sender@example.com"`
  - `ses.region: "us-east-1"`
- `config/staging/config.yaml`の`email`セクションを追加:
  - `sender_type: "ses"`
  - 上記と同様の設定
- `config/production/config.yaml.example`の`email`セクションを追加:
  - `sender_type: "ses"`
  - 上記と同様の設定

**受け入れ基準**:
- 3つの設定ファイルに`email`セクションが追加されている
- 設定値が適切に設定されている
- YAML形式が正しい

---

#### - [ ] タスク 1.3: MailpitのDocker Compose設定を作成
**目的**: MailpitをDockerで起動できるように設定を作成する

**作業内容**:
- `docker-compose.mailpit.yml`を新規作成
- Mailpitサービスの定義:
  - イメージ: `axllent/mailpit`
  - コンテナ名: `mailpit`
  - ポート設定:
    - `1025:1025` (SMTP)
    - `8025:8025` (Web UI)
  - 環境変数:
    - `MP_SMTP_BIND_ADDR=0.0.0.0:1025`
    - `MP_WEB_BIND_ADDR=0.0.0.0:8025`

**受け入れ基準**:
- `docker-compose.mailpit.yml`が作成されている
- Mailpitサービスが正しく定義されている
- ポート設定が正しい

---

#### - [ ] タスク 1.4: Mailpit起動スクリプトの作成
**目的**: Mailpitを簡単に起動・停止できるスクリプトを作成する

**作業内容**:
- `scripts/start-mailpit.sh`を新規作成
- `start`コマンド`: `docker-compose -f docker-compose.mailpit.yml up -d`
- `stop`コマンド`: `docker-compose -f docker-compose.mailpit.yml down`
- 実行権限を設定: `chmod +x scripts/start-mailpit.sh`
- 使用方法のコメントを追加

**受け入れ基準**:
- `scripts/start-mailpit.sh`が作成されている
- `start`と`stop`コマンドが実装されている
- 実行権限が設定されている

---

### Phase 2: EmailSenderインターフェースと基本実装

#### - [ ] タスク 2.1: EmailSenderインターフェースの定義
**目的**: メール送信の統一インターフェースを定義する

**作業内容**:
- `server/internal/service/email/email_sender.go`を新規作成
- `EmailSender`インターフェースを定義:
  - `Send(to []string, subject, body string) error`メソッド
- パッケージコメントを追加

**受け入れ基準**:
- `EmailSender`インターフェースが定義されている
- `Send`メソッドのシグネチャが正しい
- 既存のコードスタイルに従っている

---

#### - [ ] タスク 2.2: MockSender実装
**目的**: 標準出力にメール内容を出力するMockSenderを実装する

**作業内容**:
- `server/internal/service/email/mock_sender.go`を新規作成
- `MockSender`構造体を定義（フィールドなし）
- `NewMockSender() *MockSender`関数を実装
- `Send(to []string, subject, body string) error`メソッドを実装:
  - `fmt.Printf`を使用して標準出力に出力
  - 出力フォーマット: `[Mock Email] To: %v | Subject: %s\nBody: %s\n`
  - エラーは発生しないため、常に`nil`を返す

**受け入れ基準**:
- `MockSender`構造体が定義されている
- `NewMockSender`関数が実装されている
- `Send`メソッドが`EmailSender`インターフェースを実装している
- 出力フォーマットが正しい
- エラーハンドリングが適切に実装されている

---

### Phase 3: Mailpit送信実装

#### - [ ] タスク 3.1: gomailライブラリの依存関係追加
**目的**: gomailライブラリをプロジェクトに追加する

**作業内容**:
- `server/go.mod`に`gopkg.in/mail.v2`を追加
- `go get gopkg.in/mail.v2`を実行
- `go mod tidy`を実行

**受け入れ基準**:
- `go.mod`に`gopkg.in/mail.v2`が追加されている
- `go.sum`が更新されている

---

#### - [ ] タスク 3.2: MailpitSender実装
**目的**: gomailライブラリを使用してMailpitにメールを送信する実装を作成する

**作業内容**:
- `server/internal/service/email/mailpit_sender.go`を新規作成
- `MailpitSender`構造体を定義:
  - `smtpHost string` - SMTPホスト
  - `smtpPort int` - SMTPポート
- `NewMailpitSender(smtpHost string, smtpPort int) *MailpitSender`関数を実装
- `Send(to []string, subject, body string) error`メソッドを実装:
  - `mail.NewMessage()`でメッセージを作成
  - `SetHeader("From", ...)`で送信元を設定
  - `SetHeader("To", ...)`で送信先を設定
  - `SetHeader("Subject", subject)`で件名を設定
  - `SetBody("text/plain", body)`で本文を設定
  - `mail.NewDialer(smtpHost, smtpPort, "", "")`でダイアラーを作成
  - `DialAndSend`でメールを送信
  - エラーハンドリング: Mailpitが起動していない場合などのエラーを適切に処理

**受け入れ基準**:
- `MailpitSender`構造体が定義されている
- `NewMailpitSender`関数が実装されている
- `Send`メソッドが`EmailSender`インターフェースを実装している
- gomailライブラリが正しく使用されている
- エラーハンドリングが適切に実装されている

---

### Phase 4: AWS SES送信実装

#### - [ ] タスク 4.1: AWS SES SDKの依存関係追加
**目的**: AWS SES SDKをプロジェクトに追加する

**作業内容**:
- `server/go.mod`に`github.com/aws/aws-sdk-go-v2/service/ses`を追加
- `go get github.com/aws/aws-sdk-go-v2/service/ses`を実行
- `go get github.com/aws/aws-sdk-go-v2/config`を実行（設定用）
- `go mod tidy`を実行

**受け入れ基準**:
- `go.mod`に必要なAWS SDKパッケージが追加されている
- `go.sum`が更新されている

---

#### - [ ] タスク 4.2: SESSender実装
**目的**: AWS SES SDKを使用してメールを送信する実装を作成する

**作業内容**:
- `server/internal/service/email/ses_sender.go`を新規作成
- `SESSender`構造体を定義:
  - `client *ses.Client` - AWS SESクライアント
  - `from string` - 送信元メールアドレス
- `NewSESSender(region, from string) (*SESSender, error)`関数を実装:
  - AWS設定を読み込み（`config.LoadDefaultConfig`）
  - リージョンを設定
  - `ses.NewFromConfig`でSESクライアントを作成
  - エラーハンドリング: AWS認証情報エラーを適切に処理
- `Send(to []string, subject, body string) error`メソッドを実装:
  - `ses.SendEmailInput`を作成:
    - `Destination.ToAddresses`に送信先を設定
    - `Message.Subject.Data`に件名を設定
    - `Message.Body.Text.Data`に本文を設定（テキストメール）
    - `Source`に送信元メールアドレスを設定
  - `client.SendEmail(context.TODO(), input)`でメールを送信
  - エラーハンドリング: AWS SESエラーを適切に処理

**受け入れ基準**:
- `SESSender`構造体が定義されている
- `NewSESSender`関数が実装されている
- `Send`メソッドが`EmailSender`インターフェースを実装している
- AWS SES SDKが正しく使用されている
- エラーハンドリングが適切に実装されている

---

### Phase 5: EmailServiceと送信方式選択

#### - [ ] タスク 5.1: EmailService実装
**目的**: 環境に基づいて適切な送信方式を選択するEmailServiceを実装する

**作業内容**:
- `server/internal/service/email/email_service.go`を新規作成
- `EmailService`構造体を定義:
  - `sender EmailSender` - 現在の送信実装
- `NewEmailService(cfg *config.EmailConfig) (*EmailService, error)`関数を実装:
  - 設定ファイルの`sender_type`を確認
  - 環境変数`APP_ENV`を確認（設定ファイルが優先）
  - `sender_type`が空の場合は環境に基づいてデフォルトを選択:
    - `APP_ENV=develop`: `MockSender`を使用
    - `APP_ENV=staging`または`production`: `SESSender`を使用
  - `sender_type`に基づいて適切な送信実装を初期化:
    - `"mock"`: `NewMockSender()`を使用
    - `"mailpit"`: `NewMailpitSender(cfg.Mailpit.SMTPHost, cfg.Mailpit.SMTPPort)`を使用
    - `"ses"`: `NewSESSender(cfg.SES.Region, cfg.SES.From)`を使用
  - エラーハンドリング: 無効な`sender_type`の場合はエラーを返す
- `SendEmail(to []string, subject, body string) error`メソッドを実装:
  - 内部の`sender.Send`を呼び出す

**受け入れ基準**:
- `EmailService`構造体が定義されている
- `NewEmailService`関数が正しく実装されている
- 環境に基づいて適切な送信方式が選択される
- 設定ファイルの設定が優先される
- エラーハンドリングが適切に実装されている

---

#### - [ ] タスク 5.2: EmailServiceの初期化をmain.goに追加
**目的**: アプリケーション起動時にEmailServiceを初期化する

**作業内容**:
- `server/cmd/server/main.go`を修正
- 設定読み込み後にEmailServiceを初期化:
  - `emailService, err := email.NewEmailService(&cfg.Email)`
  - エラーハンドリング: 初期化エラーを適切に処理
- EmailServiceを後続の処理で使用できるように保持（ハンドラーに渡すため）

**受け入れ基準**:
- `main.go`でEmailServiceが初期化されている
- エラーハンドリングが適切に実装されている

---

### Phase 6: テンプレート機能

#### - [ ] タスク 6.1: TemplateService実装
**目的**: メールテンプレートの定義とデータ置換機能を実装する

**作業内容**:
- `server/internal/service/email/template.go`を新規作成
- `TemplateService`構造体を定義:
  - `templates map[string]*template.Template` - テンプレートマップ
  - `subjects map[string]string` - 件名マップ
- `NewTemplateService() *TemplateService`関数を実装:
  - テンプレートを定義（ソースコードに直書き）:
    - `welcome`テンプレート: 名前とメールアドレスを置換
    - 必要に応じて他のテンプレートも追加
  - `text/template`を使用してテンプレートをパース
  - テンプレートマップに登録
  - 件名マップに登録
- `Render(templateName string, data interface{}) (string, error)`メソッドを実装:
  - テンプレート名からテンプレートを取得
  - テンプレートが存在しない場合はエラーを返す
  - `Execute`でデータを置換してメール本文を生成
  - エラーハンドリング: テンプレート処理エラーを適切に処理
- `GetSubject(templateName string) (string, error)`メソッドを実装:
  - テンプレート名から件名を取得
  - テンプレートが存在しない場合はエラーを返す

**受け入れ基準**:
- `TemplateService`構造体が定義されている
- `NewTemplateService`関数が実装されている
- テンプレートがソースコードに直書きで定義されている
- `Render`メソッドが正しく実装されている
- `GetSubject`メソッドが正しく実装されている
- エラーハンドリングが適切に実装されている

---

### Phase 7: APIエンドポイント実装

#### - [ ] タスク 7.1: EmailHandler実装
**目的**: メール送信APIエンドポイントのハンドラーを実装する

**作業内容**:
- `server/internal/api/handler/email_handler.go`を新規作成
- `EmailHandler`構造体を定義:
  - `emailService *email.EmailService` - メール送信サービス
  - `templateService *email.TemplateService` - テンプレートサービス
- `NewEmailHandler(emailService *email.EmailService, templateService *email.TemplateService) *EmailHandler`関数を実装
- `SendEmailRequest`構造体を定義:
  - `To []string` - 送信先メールアドレスリスト
  - `Template string` - テンプレート名
  - `Data map[string]interface{}` - テンプレートデータ
- `SendEmailResponse`構造体を定義:
  - `Success bool` - 送信成功フラグ
  - `Message string` - メッセージ
- `SendEmail(c echo.Context) error`メソッドを実装:
  - リクエストボディをデコード
  - バリデーション:
    - `To`が空でないことを確認
    - メールアドレスの形式を検証（`net/mail`パッケージを使用）
    - `Template`が空でないことを確認
    - `Data`が空でないことを確認
  - テンプレートサービスでメール本文を生成:
    - `templateService.Render(req.Template, req.Data)`
    - 件名を取得: `templateService.GetSubject(req.Template)`
  - EmailServiceでメールを送信:
    - `emailService.SendEmail(req.To, subject, body)`
  - レスポンスを返却:
    - 成功時: HTTP 200 OK、`SendEmailResponse`を返す
    - エラー時: 適切なHTTPステータスコードとエラーメッセージを返す
  - エラーハンドリング:
    - バリデーションエラー: HTTP 400 Bad Request
    - テンプレートエラー: HTTP 400 Bad Request
    - 送信エラー: HTTP 500 Internal Server Error
    - 認証エラー: HTTP 401 Unauthorized（認証ミドルウェアで処理）

**受け入れ基準**:
- `EmailHandler`構造体が定義されている
- `NewEmailHandler`関数が実装されている
- `SendEmail`メソッドが正しく実装されている
- リクエストバリデーションが適切に実装されている
- エラーハンドリングが適切に実装されている
- HTTPステータスコードが正しく設定されている

---

#### - [ ] タスク 7.2: メール送信エンドポイントの登録
**目的**: Echoルーターにメール送信エンドポイントを登録する

**作業内容**:
- `server/internal/api/router/router.go`の`NewRouter`関数を修正
- `RegisterEmailEndpoints(e *echo.Echo, h *handler.EmailHandler) error`関数を実装:
  - `POST /api/email/send`エンドポイントを登録
  - 認証ミドルウェアを適用（既存の`auth.NewHumaAuthMiddleware`を使用）
  - ハンドラーを`h.SendEmail`に設定
- `NewRouter`関数内で`RegisterEmailEndpoints`を呼び出す
- EmailServiceとTemplateServiceを初期化してハンドラーに渡す
- エラーハンドリング: エンドポイント登録エラーを適切に処理

**受け入れ基準**:
- `RegisterEmailEndpoints`関数が正しく実装されている
- メール送信エンドポイントがEchoルーターに登録されている
- 認証ミドルウェアが適用されている
- エラーハンドリングが適切に実装されている

---

### Phase 8: クライアント側実装

#### - [ ] タスク 8.1: メール送信API呼び出し関数の追加
**目的**: クライアント側でメール送信APIを呼び出す関数を追加する

**作業内容**:
- `client/src/lib/api.ts`を修正
- `sendEmail`関数を追加:
  - パラメータ: `to: string[], template: string, data: Record<string, any>`
  - `POST /api/email/send`にリクエストを送信
  - 認証トークンをヘッダーに含める
  - レスポンスを返す: `{ success: boolean; message: string }`
  - エラーハンドリング: ネットワークエラー、認証エラー、バリデーションエラーを適切に処理

**受け入れ基準**:
- `sendEmail`関数が実装されている
- APIエンドポイントが正しく設定されている
- 認証トークンが正しく送信されている
- エラーハンドリングが適切に実装されている

---

#### - [ ] タスク 8.2: メール送信画面の基本実装
**目的**: メール送信画面の基本構造を実装する

**作業内容**:
- `client/src/app/dm_email/send/page.tsx`を新規作成
- Next.js App Routerのページコンポーネントとして実装
- 基本的なレイアウトとスタイリングを実装
- ページタイトルと説明を表示
- フォームの基本構造を準備

**受け入れ基準**:
- `/dm_email/send`にアクセスできる
- 基本的なレイアウトが表示される
- Next.js App Routerのパターンに従っている

---

#### - [ ] タスク 8.3: メール送信フォームの実装
**目的**: メール送信フォームの入力フィールドと送信機能を実装する

**作業内容**:
- メールアドレス入力フィールドを追加:
  - `input type="email"`を使用
  - バリデーション: メールアドレスの形式を検証
  - エラーメッセージを表示
- 名前入力フィールドを追加:
  - `input type="text"`を使用
  - バリデーション: 必須項目
  - エラーメッセージを表示
- 送信ボタンを追加:
  - ローディング状態を表示
  - 送信中は無効化
- フォーム送信処理を実装:
  - `handleSubmit`関数を実装
  - `api.sendEmail`を呼び出し
  - 送信結果を表示（成功/失敗）
  - エラーメッセージを表示
- 状態管理:
  - `useState`でフォーム状態を管理
  - ローディング状態を管理
  - エラー状態を管理

**受け入れ基準**:
- メールアドレス入力フィールドが実装されている
- 名前入力フィールドが実装されている
- 送信ボタンが実装されている
- フォームバリデーションが実装されている
- 送信処理が実装されている
- 送信結果の表示が実装されている
- エラーメッセージの表示が実装されている

---

### Phase 9: テスト実装

#### - [ ] タスク 9.1: MockSenderのユニットテスト
**目的**: MockSenderの動作を検証するユニットテストを実装する

**作業内容**:
- `server/internal/service/email/mock_sender_test.go`を新規作成
- `TestMockSender_Send`テストを実装:
  - MockSenderを作成
  - `Send`メソッドを呼び出し
  - 標準出力に正しいフォーマットで出力されることを確認
  - エラーが返されないことを確認

**受け入れ基準**:
- ユニットテストが実装されている
- テストが成功する

---

#### - [ ] タスク 9.2: MailpitSenderのユニットテスト
**目的**: MailpitSenderの動作を検証するユニットテストを実装する

**作業内容**:
- `server/internal/service/email/mailpit_sender_test.go`を新規作成
- `TestMailpitSender_Send`テストを実装:
  - MailpitSenderを作成（テスト用のSMTP設定）
  - モックSMTPサーバーを使用（または統合テストとして実装）
  - `Send`メソッドを呼び出し
  - メールが正しく送信されることを確認
  - エラーハンドリングをテスト（Mailpitが起動していない場合など）

**受け入れ基準**:
- ユニットテストが実装されている
- テストが成功する

---

#### - [ ] タスク 9.3: SESSenderのユニットテスト
**目的**: SESSenderの動作を検証するユニットテストを実装する

**作業内容**:
- `server/internal/service/email/ses_sender_test.go`を新規作成
- `TestSESSender_Send`テストを実装:
  - AWS SESクライアントをモック化
  - SESSenderを作成
  - `Send`メソッドを呼び出し
  - AWS SES SDKが正しく呼び出されることを確認
  - エラーハンドリングをテスト（AWS認証情報エラーなど）

**受け入れ基準**:
- ユニットテストが実装されている
- テストが成功する

---

#### - [ ] タスク 9.4: EmailServiceのユニットテスト
**目的**: EmailServiceの送信方式選択ロジックを検証するユニットテストを実装する

**作業内容**:
- `server/internal/service/email/email_service_test.go`を新規作成
- `TestEmailService_NewEmailService`テストを実装:
  - 環境変数`APP_ENV`を設定してテスト
  - 設定ファイルの`sender_type`を設定してテスト
  - 適切な送信実装が選択されることを確認
- `TestEmailService_SendEmail`テストを実装:
  - EmailServiceを作成
  - `SendEmail`メソッドを呼び出し
  - 内部の送信実装が正しく呼び出されることを確認

**受け入れ基準**:
- ユニットテストが実装されている
- テストが成功する

---

#### - [ ] タスク 9.5: TemplateServiceのユニットテスト
**目的**: TemplateServiceのテンプレート処理を検証するユニットテストを実装する

**作業内容**:
- `server/internal/service/email/template_test.go`を新規作成
- `TestTemplateService_Render`テストを実装:
  - TemplateServiceを作成
  - `Render`メソッドを呼び出し
  - テンプレートが正しく置換されることを確認
  - 無効なテンプレート名でエラーが返されることを確認
- `TestTemplateService_GetSubject`テストを実装:
  - `GetSubject`メソッドを呼び出し
  - 正しい件名が返されることを確認
  - 無効なテンプレート名でエラーが返されることを確認

**受け入れ基準**:
- ユニットテストが実装されている
- テストが成功する

---

#### - [ ] タスク 9.6: EmailHandlerの統合テスト
**目的**: EmailHandlerのAPIエンドポイントを検証する統合テストを実装する

**作業内容**:
- `server/internal/api/handler/email_handler_test.go`を新規作成
- `TestEmailHandler_SendEmail`テストを実装:
  - EmailHandlerを作成（モックサービスを使用）
  - HTTPリクエストを作成
  - レスポンスを検証:
    - 成功時: HTTP 200 OK
    - バリデーションエラー時: HTTP 400 Bad Request
    - テンプレートエラー時: HTTP 400 Bad Request
    - 送信エラー時: HTTP 500 Internal Server Error
  - リクエストボディのバリデーションをテスト
  - メールアドレスの形式検証をテスト

**受け入れ基準**:
- 統合テストが実装されている
- テストが成功する

---

#### - [ ] タスク 9.7: クライアント側のユニットテスト
**目的**: メール送信画面のコンポーネントを検証するユニットテストを実装する

**作業内容**:
- `client/src/app/dm_email/send/__tests__/page.test.tsx`を新規作成
- コンポーネントのレンダリングテストを実装
- フォームバリデーションテストを実装:
  - メールアドレスの形式検証
  - 必須項目の検証
- 送信処理のテストを実装（モックAPIを使用）

**受け入れ基準**:
- ユニットテストが実装されている
- テストが成功する

---

#### - [ ] タスク 9.8: E2Eテストの実装
**目的**: メール送信機能のE2Eテストを実装する

**作業内容**:
- `client/e2e/email-send.spec.ts`を新規作成
- メール送信成功シナリオのテスト:
  - メール送信画面にアクセス
  - メールアドレスと名前を入力
  - 送信ボタンをクリック
  - 成功メッセージが表示されることを確認
- エラーシナリオのテスト:
  - バリデーションエラー（無効なメールアドレス）
  - 認証エラー（認証トークンなし）
  - 送信エラー（Mailpitが起動していない場合など）

**受け入れ基準**:
- E2Eテストが実装されている
- テストが成功する

---

### Phase 10: ドキュメントと最終確認

#### - [ ] タスク 10.1: 実装の最終確認
**目的**: 実装が要件定義書と設計書に準拠していることを確認する

**作業内容**:
- 要件定義書の各要件を確認
- 設計書の各コンポーネントを確認
- コードレビュー:
  - コードスタイルの確認
  - エラーハンドリングの確認
  - テストカバレッジの確認
- 動作確認:
  - 開発環境でMockSenderが動作することを確認
  - Mailpitでメール送信が動作することを確認
  - 本番環境（SES）での動作確認（可能な場合）

**受け入れ基準**:
- 要件定義書の要件がすべて実装されている
- 設計書の設計に準拠している
- コードレビューで問題がない
- 動作確認が成功している

---

#### - [ ] タスク 10.2: 既存テストの確認
**目的**: 既存のテストがすべて成功することを確認する

**作業内容**:
- サーバー側の既存テストを実行: `cd server && go test ./...`
- クライアント側の既存テストを実行: `cd client && npm test`
- テストが失敗する場合は修正

**受け入れ基準**:
- 既存のテストがすべて成功する

---

## 実装順序の推奨

1. **Phase 1**: 設定とインフラの準備（必須）
2. **Phase 2**: EmailSenderインターフェースと基本実装（必須）
3. **Phase 3**: Mailpit送信実装（開発環境での確認用）
4. **Phase 4**: AWS SES送信実装（本番環境用）
5. **Phase 5**: EmailServiceと送信方式選択（必須）
6. **Phase 6**: テンプレート機能（必須）
7. **Phase 7**: APIエンドポイント実装（必須）
8. **Phase 8**: クライアント側実装（必須）
9. **Phase 9**: テスト実装（品質保証）
10. **Phase 10**: ドキュメントと最終確認（品質保証）

## 注意事項

- 各タスクは独立して実装可能なように設計されていますが、依存関係がある場合は順序を守ってください
- テストは実装と並行して進めることを推奨します
- エラーハンドリングは各タスクで適切に実装してください
- コードスタイルは既存のコードに従ってください
