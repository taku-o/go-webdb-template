# メール送信機能設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、メール送信機能の詳細設計を定義する。3つの送信方式（標準出力、Mailpit、AWS SES）をサポートし、テンプレート機能により動的なメール送信を可能にするシステムを構築する。

### 1.2 設計の範囲
- サーバー側: メール送信インターフェースと3つの実装、テンプレート機能、APIエンドポイント
- クライアント側: メール送信画面の実装
- 設定: MailpitのDocker設定と起動スクリプト
- 認証・認可: Public API Key JWT または Auth0 JWT による認証
- エラーハンドリング: サーバー側とクライアント側のエラーハンドリング
- テスト: 単体テストとE2Eテストの実装

### 1.3 設計方針
- **インターフェースベースの設計**: EmailSenderインターフェースにより送信方式を切り替え可能
- **環境別送信方式**: 開発環境では標準出力、本番環境ではAWS SESをデフォルトで使用
- **テンプレート機能**: ソースコードに直書きでテンプレートを管理
- **既存システムとの統合**: 既存のHuma API、Echo、Next.jsのパターンに従う
- **拡張性の確保**: 将来的なHTMLメール、添付ファイルなどに対応できる設計

## 2. アーキテクチャ設計

### 2.1 ディレクトリ構造

#### 2.1.1 変更前の構造
```
client/
├── src/
│   ├── app/
│   │   └── (既存のページ)
│   └── lib/
│       └── api.ts

server/
├── internal/
│   ├── api/
│   │   ├── handler/
│   │   │   └── (既存のハンドラー)
│   │   └── router/
│   │       └── router.go
│   ├── service/
│   │   └── (既存のサービス)
│   └── config/
│       └── config.go
```

#### 2.1.2 変更後の構造
```
client/
├── src/
│   ├── app/
│   │   └── dm_email/
│   │       └── send/
│   │           └── page.tsx        # メール送信画面（新規作成）
│   └── lib/
│       └── api.ts                  # メール送信API呼び出し追加

server/
├── internal/
│   ├── api/
│   │   ├── handler/
│   │   │   └── email_handler.go   # メール送信ハンドラー（新規作成）
│   │   └── router/
│   │       └── router.go          # メール送信エンドポイント登録追加
│   ├── service/
│   │   └── email/                 # メール送信サービス（新規作成）
│   │       ├── email_sender.go     # EmailSenderインターフェース
│   │       ├── mock_sender.go      # MockSender実装
│   │       ├── mailpit_sender.go   # MailpitSender実装
│   │       ├── ses_sender.go       # SESSender実装
│   │       ├── email_service.go    # EmailService（送信方式選択）
│   │       └── template.go         # テンプレート機能
│   └── config/
│       └── config.go               # EmailConfig追加
├── cmd/
│   └── server/
│       └── main.go                 # EmailService初期化追加
├── scripts/
│   └── start-mailpit.sh            # Mailpit起動スクリプト（新規作成）
├── docker-compose.mailpit.yml      # Mailpit設定（新規作成）
└── config/
    ├── develop/
    │   └── config.yaml              # email設定追加
    ├── staging/
    │   └── config.yaml              # email設定追加
    └── production/
        └── config.yaml.example      # email設定追加
```

### 2.2 ファイル構成

#### 2.2.1 変更ファイル
- **サーバー側**:
  - `server/internal/api/router/router.go`: メール送信エンドポイントの登録を追加
  - `server/internal/config/config.go`: `EmailConfig`構造体を追加
  - `server/cmd/server/main.go`: EmailServiceの初期化を追加
  - `config/develop/config.yaml`: メール送信設定を追加
  - `config/staging/config.yaml`: メール送信設定を追加
  - `config/production/config.yaml.example`: メール送信設定を追加
- **クライアント側**:
  - `client/src/lib/api.ts`: メール送信API呼び出しを追加

#### 2.2.2 新規作成ファイル
- **クライアント側**:
  - `client/src/app/dm_email/send/page.tsx`: メール送信画面
- **サーバー側**:
  - `server/internal/api/handler/email_handler.go`: メール送信ハンドラー
  - `server/internal/service/email/email_sender.go`: EmailSenderインターフェース
  - `server/internal/service/email/mock_sender.go`: MockSender実装
  - `server/internal/service/email/mailpit_sender.go`: MailpitSender実装
  - `server/internal/service/email/ses_sender.go`: SESSender実装
  - `server/internal/service/email/email_service.go`: EmailService
  - `server/internal/service/email/template.go`: テンプレート機能
- **設定・スクリプト**:
  - `docker-compose.mailpit.yml`: Mailpit Docker設定
  - `scripts/start-mailpit.sh`: Mailpit起動スクリプト

### 2.3 システム構成図

```
┌─────────────────────────────────────────────────────────┐
│              クライアント（Next.js）                        │
│  ┌──────────────────────────────────────────────────┐  │
│  │  dm_email/send/page.tsx                          │  │
│  │  - メールアドレス入力                              │  │
│  │  - 名前入力                                       │  │
│  │  - 送信ボタン                                     │  │
│  └──────────────────┬───────────────────────────────┘  │
│                     │                                    │
│                     │ HTTP POST                          │
│                     │ /api/email/send                   │
│                     │ Authorization: Bearer {JWT}        │
│                     ▼                                    │
└─────────────────────┼──────────────────────────────────┘
                    │ HTTP
                    │ /api/email/send
                    │ Authorization: Bearer {JWT}
                    ▼
┌─────────────────────────────────────────────────────────┐
│              APIサーバー（Go + Echo + Huma）               │
│  ┌──────────────────────────────────────────────────┐  │
│  │  router.go                                      │  │
│  │  - メール送信エンドポイント登録                    │  │
│  │  - 認証ミドルウェア適用                           │  │
│  └──────────────────┬───────────────────────────────┘  │
│                     │                                    │
│                     ▼                                    │
│  ┌──────────────────────────────────────────────────┐  │
│  │  email_handler.go                               │  │
│  │  - リクエストバリデーション                        │  │
│  │  - テンプレート処理                               │  │
│  │  - EmailService呼び出し                          │  │
│  └──────────────────┬───────────────────────────────┘  │
│                     │                                    │
│                     ▼                                    │
│  ┌──────────────────────────────────────────────────┐  │
│  │  email_service.go                               │  │
│  │  - 送信方式の選択                                 │  │
│  │  - EmailSender呼び出し                           │  │
│  └──────────────────┬───────────────────────────────┘  │
│                     │                                    │
│                     ▼                                    │
│  ┌──────────────────────────────────────────────────┐  │
│  │  EmailSender実装                                 │  │
│  │  - MockSender (標準出力)                         │  │
│  │  - MailpitSender (gomail + SMTP)                 │  │
│  │  - SESSender (AWS SES SDK)                       │  │
│  └──────────────────┬───────────────────────────────┘  │
│                     │                                    │
│                     ▼                                    │
│  ┌──────────────────────────────────────────────────┐  │
│  │  送信先                                            │  │
│  │  - 開発環境: 標準出力                              │  │
│  │  - 開発環境（Mailpit）: localhost:1025 (SMTP)    │  │
│  │  - 本番環境: AWS SES                              │  │
│  └──────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

### 2.4 データフロー

#### 2.4.1 メール送信フロー
```
1. ユーザーがメール送信画面でメールアドレスと名前を入力
    ↓
2. クライアントがPOSTリクエストを送信（/api/email/send）
    ↓
3. APIサーバー側で認証チェック（public API）
    ↓
4. email_handler.goでリクエストバリデーション
    ↓
5. template.goでテンプレートを選択し、データを置換
    ↓
6. email_service.goで送信方式を選択（環境に基づく）
    ↓
7. 選択されたEmailSender実装（MockSender/MailpitSender/SESSender）でメール送信
    ↓
8. 送信結果をクライアントに返却
    ↓
9. クライアントで送信結果を表示
```

## 3. コンポーネント設計

### 3.1 クライアント側コンポーネント

#### 3.1.1 メール送信画面（dm_email/send/page.tsx）

**責任と境界**
- **主要責任**: メール送信UIの提供、フォーム入力、送信結果の表示
- **ドメイン境界**: フロントエンドUI層
- **データ所有**: フォーム状態（メールアドレス、名前、送信状態、エラー）

**依存関係**
- **インバウンド**: なし（ページコンポーネント）
- **アウトバウンド**: api.ts（メール送信API呼び出し）
- **外部**: `/api/email/send` エンドポイント

**インターフェース設計**

```typescript
// コンポーネント構造
export default function EmailSendPage() {
  const [email, setEmail] = useState('')
  const [name, setName] = useState('')
  const [loading, setLoading] = useState(false)
  const [result, setResult] = useState<{ success: boolean; message: string } | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    // API呼び出し
  }

  return (
    <form onSubmit={handleSubmit}>
      {/* メールアドレス入力 */}
      {/* 名前入力 */}
      {/* 送信ボタン */}
      {/* 結果表示 */}
    </form>
  )
}
```

**主要機能**
- メールアドレス入力フィールド（バリデーション付き）
- 名前入力フィールド
- 送信ボタン（ローディング状態表示）
- 送信結果の表示（成功/失敗）
- エラーメッセージの表示

### 3.2 サーバー側コンポーネント

#### 3.2.1 メール送信ハンドラー（email_handler.go）

**責任と境界**
- **主要責任**: HTTPリクエスト/レスポンスの処理、バリデーション、EmailService呼び出し
- **ドメイン境界**: API層
- **データ所有**: リクエストデータ、レスポンスデータ

**依存関係**
- **インバウンド**: router.go（エンドポイント登録）
- **アウトバウンド**: EmailService、TemplateService、authパッケージ
- **外部**: なし

**インターフェース設計**

```go
// EmailHandler はメール送信APIのハンドラー
type EmailHandler struct {
    emailService *email.EmailService
    templateService *email.TemplateService
}

// NewEmailHandler は新しいEmailHandlerを作成
func NewEmailHandler(emailService *email.EmailService, templateService *email.TemplateService) *EmailHandler

// SendEmail はメール送信リクエストを処理
func (h *EmailHandler) SendEmail(c echo.Context) error
```

**主要機能**
- リクエストボディのデコードとバリデーション
- メールアドレスの形式検証
- テンプレートの選択とデータ置換
- EmailServiceによるメール送信
- レスポンスのフォーマットと返却
- エラーハンドリングとHTTPステータスコードマッピング

#### 3.2.2 EmailSenderインターフェース（email_sender.go）

**責任と境界**
- **主要責任**: メール送信の抽象化インターフェース定義
- **ドメイン境界**: サービス層
- **データ所有**: なし（インターフェース）

**依存関係**
- **インバウンド**: EmailService
- **アウトバウンド**: なし（インターフェース）

**インターフェース設計**

```go
// EmailSender はメール送信のインターフェース
type EmailSender interface {
    Send(to []string, subject, body string) error
}
```

**主要機能**
- メール送信の統一インターフェース提供
- 複数の送信方式実装の抽象化

#### 3.2.3 MockSender実装（mock_sender.go）

**責任と境界**
- **主要責任**: 標準出力へのメール内容出力
- **ドメイン境界**: サービス層
- **データ所有**: なし（ステートレス）

**依存関係**
- **インバウンド**: EmailService
- **アウトバウンド**: 標準出力（fmt.Printf）

**インターフェース設計**

```go
// MockSender は標準出力にメールを出力する送信実装
type MockSender struct{}

// NewMockSender は新しいMockSenderを作成
func NewMockSender() *MockSender

// Send はメール内容を標準出力に出力
func (s *MockSender) Send(to []string, subject, body string) error
```

**主要機能**
- メール内容を標準出力に出力
- 出力フォーマット: `[Mock Email] To: %v | Subject: %s\nBody: %s\n`
- エラーハンドリング（基本的にエラーは発生しない）

#### 3.2.4 MailpitSender実装（mailpit_sender.go）

**責任と境界**
- **主要責任**: gomailライブラリを使用したSMTP経由のMailpitへのメール送信
- **ドメイン境界**: サービス層
- **データ所有**: SMTP設定（ホスト、ポート）

**依存関係**
- **インバウンド**: EmailService
- **アウトバウンド**: gomailライブラリ (`gopkg.in/mail.v2`)、Mailpit SMTPサーバー

**インターフェース設計**

```go
// MailpitSender はMailpitにメールを送信する実装
type MailpitSender struct {
    smtpHost string
    smtpPort int
}

// NewMailpitSender は新しいMailpitSenderを作成
func NewMailpitSender(smtpHost string, smtpPort int) *MailpitSender

// Send はgomailを使用してMailpitにメールを送信
func (s *MailpitSender) Send(to []string, subject, body string) error
```

**主要機能**
- gomailライブラリを使用したSMTP経由のメール送信
- Mailpit SMTPサーバー（デフォルト: `localhost:1025`）への接続
- 送信元メールアドレスの設定
- エラーハンドリング（Mailpitが起動していない場合など）

#### 3.2.5 SESSender実装（ses_sender.go）

**責任と境界**
- **主要責任**: AWS SES SDKを使用したメール送信
- **ドメイン境界**: サービス層
- **データ所有**: AWS SESクライアント、送信元メールアドレス

**依存関係**
- **インバウンド**: EmailService
- **アウトバウンド**: AWS SES SDK (`github.com/aws/aws-sdk-go-v2/service/ses`)、AWS認証情報

**インターフェース設計**

```go
// SESSender はAWS SESにメールを送信する実装
type SESSender struct {
    client *ses.Client
    from   string
}

// NewSESSender は新しいSESSenderを作成
func NewSESSender(client *ses.Client, from string) *SESSender

// Send はAWS SES SDKを使用してメールを送信
func (s *SESSender) Send(to []string, subject, body string) error
```

**主要機能**
- AWS SES SDKを使用したメール送信
- 送信元メールアドレスの設定（設定ファイルから読み込み）
- テキストメールの送信（HTMLメールは将来的な拡張項目）
- エラーハンドリング（AWS認証情報エラーなど）

#### 3.2.6 EmailService（email_service.go）

**責任と境界**
- **主要責任**: 環境に基づく送信方式の選択、EmailSenderの管理
- **ドメイン境界**: サービス層
- **データ所有**: 現在のEmailSender実装、設定情報

**依存関係**
- **インバウンド**: email_handler.go
- **アウトバウンド**: EmailSender実装（MockSender/MailpitSender/SESSender）、configパッケージ

**インターフェース設計**

```go
// EmailService はメール送信サービス
type EmailService struct {
    sender EmailSender
}

// NewEmailService は新しいEmailServiceを作成
func NewEmailService(cfg *config.EmailConfig) (*EmailService, error)

// SendEmail はメールを送信
func (s *EmailService) SendEmail(to []string, subject, body string) error
```

**主要機能**
- 環境変数 `APP_ENV` に基づく送信方式の選択
- 設定ファイルでの送信方式の明示的指定
- EmailSender実装の初期化と管理
- メール送信の統一インターフェース提供

#### 3.2.7 テンプレート機能（template.go）

**責任と境界**
- **主要責任**: メールテンプレートの定義とデータ置換
- **ドメイン境界**: サービス層
- **データ所有**: テンプレート定義（ソースコードに直書き）

**依存関係**
- **インバウンド**: email_handler.go
- **アウトバウンド**: Go標準ライブラリ `text/template`

**インターフェース設計**

```go
// TemplateService はメールテンプレートサービス
type TemplateService struct {
    templates map[string]*template.Template
}

// NewTemplateService は新しいTemplateServiceを作成
func NewTemplateService() *TemplateService

// Render はテンプレート名とデータからメール本文を生成
func (s *TemplateService) Render(templateName string, data interface{}) (string, error)

// GetSubject はテンプレート名に基づいて件名を取得
func (s *TemplateService) GetSubject(templateName string) string
```

**主要機能**
- メールテンプレートの定義（ソースコードに直書き）
- 複数のテンプレートを一箇所にまとめて定義
- テンプレート内のプレースホルダーを動的に置換
- テンプレート名によるテンプレートの選択
- 件名の取得

#### 3.2.8 設定管理（config/config.go）

**責任と境界**
- **主要責任**: メール送信機能の設定管理
- **ドメイン境界**: 設定層
- **データ所有**: メール送信設定値

**依存関係**
- **インバウンド**: EmailService
- **アウトバウンド**: viper（設定ファイル読み込み）

**インターフェース設計**

```go
// EmailConfig はメール送信機能の設定
type EmailConfig struct {
    SenderType string      `mapstructure:"sender_type"` // "mock", "mailpit", "ses"
    Mock       MockConfig  `mapstructure:"mock"`
    Mailpit    MailpitConfig `mapstructure:"mailpit"`
    SES        SESConfig   `mapstructure:"ses"`
}

// MockConfig はMockSenderの設定
type MockConfig struct {
    // 設定項目なし（標準出力に出力するだけ）
}

// MailpitConfig はMailpitSenderの設定
type MailpitConfig struct {
    SMTPHost string `mapstructure:"smtp_host"` // デフォルト: "localhost"
    SMTPPort int    `mapstructure:"smtp_port"` // デフォルト: 1025
}

// SESConfig はSESSenderの設定
type SESConfig struct {
    From   string `mapstructure:"from"`   // 送信元メールアドレス
    Region string `mapstructure:"region"` // AWSリージョン
    // AWS認証情報は環境変数から取得
}
```

## 4. データモデル

### 4.1 リクエストデータモデル

```go
// SendEmailRequest はメール送信リクエスト
type SendEmailRequest struct {
    To       []string               `json:"to" validate:"required,min=1,dive,email"`
    Template string                 `json:"template" validate:"required"`
    Data     map[string]interface{} `json:"data" validate:"required"`
}
```

### 4.2 レスポンスデータモデル

```go
// SendEmailResponse はメール送信レスポンス
type SendEmailResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}
```

### 4.3 テンプレートデータモデル

テンプレートはソースコードに直書きで定義します。例：

```go
// テンプレート定義例
const (
    templateWelcome = `こんにちは、{{.Name}}さん。

メールアドレス: {{.Email}}

ご登録ありがとうございます。`
    
    templateWelcomeSubject = "ようこそ"
)
```

## 5. エラーハンドリング

### 5.1 エラー戦略

#### 5.1.1 クライアント側エラー
- **ネットワークエラー**: エラーメッセージを表示
- **認証エラー**: エラーメッセージを表示し、ログイン画面へリダイレクト
- **バリデーションエラー**: フィールドごとのエラーメッセージを表示
- **送信失敗**: エラーメッセージを表示

#### 5.1.2 サーバー側エラー
- **認証エラー**: HTTP 401 Unauthorized
- **バリデーションエラー**: HTTP 400 Bad Request
- **テンプレートエラー**: HTTP 400 Bad Request（無効なテンプレート名など）
- **送信エラー**: HTTP 500 Internal Server Error
- **Mailpit接続エラー**: HTTP 503 Service Unavailable（Mailpitが起動していない場合）

### 5.2 エラーカテゴリとレスポンス

**ユーザーエラー（4xx）**
- **400 Bad Request**: 無効なリクエスト、バリデーションエラー、無効なテンプレート名
- **401 Unauthorized**: 認証エラー

**システムエラー（5xx）**
- **500 Internal Server Error**: メール送信エラー、予期しないエラー
- **503 Service Unavailable**: Mailpitが起動していない、AWS SES接続エラー

### 5.3 モニタリング

- **ログ記録**: メール送信開始、完了、エラーをログに記録
- **メトリクス**: メール送信成功率、平均送信時間、エラー率
- **アラート**: 送信エラーの増加時にアラート

## 6. テスト戦略

### 6.1 ユニットテスト

#### 6.1.1 サーバー側
- **email_sender.go**: インターフェース定義のテスト
- **mock_sender.go**: 標準出力への出力テスト
- **mailpit_sender.go**: gomail設定とSMTP接続テスト（モック使用）
- **ses_sender.go**: AWS SES SDK呼び出しテスト（モック使用）
- **email_service.go**: 送信方式選択ロジックのテスト
- **template.go**: テンプレート処理とデータ置換のテスト
- **email_handler.go**: リクエストバリデーション、エラーハンドリングのテスト

#### 6.1.2 クライアント側
- **dm_email/send/page.tsx**: コンポーネントのレンダリングテスト、フォームバリデーションテスト

### 6.2 統合テスト

#### 6.2.1 サーバー側
- **メール送信統合**: MockSender、MailpitSender、SESSenderの統合テスト
- **認証統合**: Public API Key JWT、Auth0 JWTの認証テスト
- **テンプレート統合**: テンプレート選択とデータ置換の統合テスト

#### 6.2.2 クライアント側
- **メール送信フロー**: フォーム入力から送信完了までのテスト
- **エラーハンドリング**: ネットワークエラー、認証エラー、バリデーションエラーのテスト

### 6.3 E2Eテスト

- **メール送信成功シナリオ**: フォーム入力 → 送信 → 成功通知
- **エラーシナリオ**: バリデーションエラー、認証エラー、送信エラー
- **環境別送信方式**: 開発環境（MockSender）、Mailpit、本番環境（SES）のテスト

## 7. セキュリティ考慮事項

### 7.1 認証・認可
- **認証方式**: Public API Key JWT または Auth0 JWT
- **認証チェック**: メール送信エンドポイントへの全リクエストで認証を要求
- **認証ミドルウェア**: 既存のauth.NewHumaAuthMiddlewareを使用

### 7.2 入力検証
- **メールアドレス検証**: 送信先メールアドレスの形式検証
- **テンプレート名検証**: 無効なテンプレート名の拒否
- **データ検証**: テンプレートデータの型検証

### 7.3 認証情報管理
- **AWS認証情報**: 環境変数またはAWS認証情報ファイルから取得
- **設定ファイル**: 機密情報は設定ファイルに保存（gitから除外）

### 7.4 メール送信制限
- **送信先制限**: 現時点では制限なし（将来的な拡張項目）
- **送信頻度制限**: 現時点では制限なし（将来的な拡張項目）

## 8. パフォーマンスとスケーラビリティ

### 8.1 パフォーマンス
- **同期送信**: メール送信は同期的に実行（非同期処理は将来的な拡張項目）
- **タイムアウト**: メール送信のタイムアウト設定（デフォルト: 30秒）
- **接続プール**: Mailpit、AWS SESの接続管理

### 8.2 スケーラビリティ
- **水平スケーリング**: 複数のAPIサーバーインスタンスで動作可能
- **AWS SESスケーリング**: AWS SESの自動スケーリング機能を活用
- **負荷分散**: ロードバランサー経由での複数インスタンスへの分散

## 9. 技術スタックと設計決定

### 9.1 技術スタック

#### 9.1.1 フロントエンド
- **フレームワーク**: Next.js 14+ (App Router)
- **言語**: TypeScript 5+

#### 9.1.2 バックエンド
- **言語**: Go 1.21+
- **AWS SDK**: `github.com/aws/aws-sdk-go-v2/service/ses`
- **gomail**: `gopkg.in/mail.v2` (Mailpit送信用)
- **テンプレート処理**: Go標準ライブラリの `text/template`

### 9.2 主要な設計決定

#### 決定1: インターフェースベースの設計
- **コンテキスト**: 複数の送信方式（標準出力、Mailpit、AWS SES）をサポートする必要がある
- **決定**: EmailSenderインターフェースを定義し、各送信方式を実装
- **利点**: 送信方式の切り替えが容易、テストが容易、拡張性が高い
- **代替案**: 条件分岐による送信方式の選択（拡張性が低い）

#### 決定2: 環境別送信方式の自動選択
- **コンテキスト**: 開発環境と本番環境で異なる送信方式を使用する必要がある
- **決定**: 環境変数 `APP_ENV` に基づいて自動的に送信方式を選択
- **利点**: 設定が簡潔、環境ごとの適切な送信方式が自動選択される
- **代替案**: 常に設定ファイルで明示的に指定（設定が煩雑）

#### 決定3: テンプレートのソースコード直書き
- **コンテキスト**: テンプレートを設定ファイルではなくソースコードに直書きする
- **決定**: template.goファイルにテンプレートを定義
- **利点**: テンプレートの変更がコードレビューで確認できる、バージョン管理が容易
- **代替案**: 設定ファイルやデータベースに保存（柔軟性が高いが複雑）

#### 決定4: gomailライブラリの使用
- **コンテキスト**: Mailpitへのメール送信にSMTPプロトコルを使用する
- **決定**: gomailライブラリ (`gopkg.in/mail.v2`) を使用
- **利点**: 標準的なSMTPライブラリ、シンプルなAPI
- **代替案**: net/smtp標準ライブラリ（より低レベルで複雑）

## 10. 設定ファイル例

### 10.1 開発環境（develop/config.yaml）

```yaml
email:
  sender_type: "mock"  # デフォルト: "mock"
  mock: {}
  mailpit:
    smtp_host: "localhost"
    smtp_port: 1025
  ses:
    from: "sender@example.com"
    region: "us-east-1"
```

### 10.2 ステージング環境（staging/config.yaml）

```yaml
email:
  sender_type: "ses"  # デフォルト: "ses"
  mock: {}
  mailpit:
    smtp_host: "localhost"
    smtp_port: 1025
  ses:
    from: "sender@example.com"
    region: "us-east-1"
```

### 10.3 本番環境（production/config.yaml.example）

```yaml
email:
  sender_type: "ses"  # デフォルト: "ses"
  mock: {}
  mailpit:
    smtp_host: "localhost"
    smtp_port: 1025
  ses:
    from: "sender@example.com"
    region: "us-east-1"
```

## 11. Mailpit設定

### 11.1 Docker Compose設定（docker-compose.mailpit.yml）

```yaml
version: '3.8'

services:
  mailpit:
    image: axllent/mailpit
    container_name: mailpit
    ports:
      - "1025:1025"  # SMTP
      - "8025:8025"  # Web UI
    environment:
      - MP_SMTP_BIND_ADDR=0.0.0.0:1025
      - MP_WEB_BIND_ADDR=0.0.0.0:8025
```

### 11.2 起動スクリプト（scripts/start-mailpit.sh）

```bash
#!/bin/bash

# Mailpit起動スクリプト

case "$1" in
  start)
    echo "Starting Mailpit..."
    docker-compose -f docker-compose.mailpit.yml up -d
    echo "Mailpit started. Web UI: http://localhost:8025"
    ;;
  stop)
    echo "Stopping Mailpit..."
    docker-compose -f docker-compose.mailpit.yml down
    echo "Mailpit stopped."
    ;;
  *)
    echo "Usage: $0 {start|stop}"
    exit 1
    ;;
esac
```
