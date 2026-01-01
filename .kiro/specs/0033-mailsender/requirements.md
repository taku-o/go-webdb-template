# メール送信機能要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #66
- **Issueタイトル**: メール送信機能の実装
- **Feature名**: 0033-mailsender
- **作成日**: 2025-01-27

### 1.2 目的
メール送信機能を実装する。開発環境では標準出力への出力、Mailpitによる確認、本番環境ではAWS SESによる実際のメール送信をサポートし、テンプレート機能により動的なメール送信を可能にする。

### 1.3 スコープ
- サーバー側: メール送信インターフェースと3つの実装（標準出力、Mailpit、AWS SES）
- サーバー側: メールテンプレート機能の実装
- サーバー側: メール送信APIエンドポイントの実装
- クライアント側: メール送信画面の実装
- 開発環境: 標準出力への出力（デフォルト）
- 開発環境: Mailpitによるメール確認（必要な時に切り替えて確認に使用する）
- 本番環境: AWS SESによるメール送信（デフォルト）

**本実装の範囲外**:
- メール受信機能
- メール一覧表示機能
- メール送信履歴の保存・管理機能
- 添付ファイルの送信機能
- HTMLメールの送信機能（テキストメールのみ）
- メール送信のスケジューリング機能

## 2. 背景・現状分析

### 2.1 現在の実装
- **メール送信機能**: 現在、メール送信機能は実装されていない
- **メールテンプレート**: 現在、メールテンプレート機能は実装されていない
- **Mailpit**: 現在、Mailpitの設定は存在しない

### 2.2 課題点
1. **メール送信機能の不在**: アプリケーションからメールを送信する機能が存在しない
2. **環境別の送信方式の不在**: 開発環境と本番環境で異なる送信方式を切り替える仕組みが存在しない
3. **メールテンプレート機能の不在**: 動的なメール本文を生成する機能が存在しない
4. **メール送信確認手段の不在**: 開発環境でメール内容を確認する手段が存在しない

### 2.3 本実装による改善点
1. **メール送信機能の提供**: アプリケーションからメールを送信できるようになる
2. **環境別送信方式の実装**: 開発環境では標準出力、本番環境ではAWS SESを使用できる
3. **メールテンプレート機能の実装**: テンプレートに動的な値を埋め込んでメールを送信できる
4. **開発環境での確認手段の提供**: Mailpitを使用して開発環境でメール内容を確認できる

## 3. 機能要件

### 3.1 サーバー側の実装

#### 3.1.1 メール送信インターフェースの実装
- **ファイル**: `server/internal/service/email/email_sender.go` (新規作成)
- **実装内容**:
  - `EmailSender` インターフェースの定義
  - `Send(to []string, subject, body string) error` メソッドの定義

#### 3.1.2 標準出力送信実装（MockSender）
- **ファイル**: `server/internal/service/email/mock_sender.go` (新規作成)
- **実装内容**:
  - `MockSender` 構造体の実装
  - `EmailSender` インターフェースの実装
  - 標準出力にメール内容を出力する機能
  - 出力フォーマット: `[Mock Email] To: %v | Subject: %s\nBody: %s\n`

#### 3.1.3 Mailpit送信実装（MailpitSender）
- **ファイル**: `server/internal/service/email/mailpit_sender.go` (新規作成)
- **実装内容**:
  - `MailpitSender` 構造体の実装
  - `EmailSender` インターフェースの実装
  - gomailライブラリ (`gopkg.in/mail.v2`) を使用したSMTP経由のメール送信機能
  - MailpitのSMTP設定（デフォルト: `localhost:1025`）

#### 3.1.4 AWS SES送信実装（SESSender）
- **ファイル**: `server/internal/service/email/ses_sender.go` (新規作成)
- **実装内容**:
  - `SESSender` 構造体の実装
  - `EmailSender` インターフェースの実装
  - AWS SES SDK (`github.com/aws/aws-sdk-go-v2/service/ses`) を使用したメール送信機能
  - 送信元メールアドレスの設定（設定ファイルから読み込み）

#### 3.1.5 メール送信方式の選択
- **ファイル**: `server/internal/service/email/email_service.go` (新規作成)
- **実装内容**:
  - 環境変数 `APP_ENV` に基づいて送信方式を選択
  - `develop` 環境: `MockSender` をデフォルトで使用
  - `staging`, `production` 環境: `SESSender` をデフォルトで使用
  - 設定ファイルで送信方式を明示的に指定可能

#### 3.1.6 メールテンプレート機能の実装
- **ファイル**: `server/internal/service/email/template.go` (新規作成)
- **実装内容**:
  - メールテンプレートの定義（ソースコードに直書き）
  - 複数のテンプレートを定義可能
  - テンプレート内のプレースホルダーを動的に置換する機能
  - テンプレート名によるテンプレートの選択機能

#### 3.1.7 メール送信APIエンドポイントの実装
- **ファイル**: `server/internal/api/handler/email_handler.go` (新規作成)
- **エンドポイント**: `/api/email/send`
- **HTTPメソッド**: POST
- **アクセスレベル**: `public` (Public API Key JWT または Auth0 JWT でアクセス可能)
- **リクエストボディ**:
  ```json
  {
    "to": ["recipient@example.com"],
    "template": "welcome",
    "data": {
      "name": "John Doe",
      "email": "john@example.com"
    }
  }
  ```
- **レスポンス**:
  ```json
  {
    "success": true,
    "message": "Email sent successfully"
  }
  ```

### 3.2 クライアント側の実装

#### 3.2.1 メール送信画面の実装
- **ファイル**: `client/src/app/dm_email/send/page.tsx` (新規作成)
- **URL**: `/dm_email/send`
- **実装内容**:
  - メールアドレス入力フィールド
  - 名前入力フィールド
  - 送信ボタン
  - 送信結果の表示（成功/失敗）
  - エラーメッセージの表示

### 3.3 Mailpitの設定

#### 3.3.1 Docker Compose設定
- **ファイル**: `docker-compose.mailpit.yml` (新規作成)
- **実装内容**:
  - Mailpitサービスの定義
  - ポート設定: 1025 (SMTP), 8025 (Web UI)
  - 環境変数の設定

#### 3.3.2 Mailpit起動スクリプト
- **ファイル**: `scripts/start-mailpit.sh` (新規作成)
- **実装内容**:
  - Mailpitの起動コマンド
  - 停止コマンド
  - 実行権限の設定

## 4. 非機能要件

### 4.1 パフォーマンス
- メール送信は非同期処理を想定しない（同期的に送信）
- メール送信のタイムアウト設定（デフォルト: 30秒）

### 4.2 セキュリティ
- 認証: Public API Key JWT または Auth0 JWT による認証が必要
- メールアドレスのバリデーション: 送信先メールアドレスの形式検証
- AWS SES認証情報: 環境変数または設定ファイルから安全に読み込み

### 4.3 可用性
- エラーハンドリング: 適切なエラーメッセージとログ記録
- 送信失敗時のリトライ機能は不要（将来的な拡張項目）

### 4.4 保守性
- テンプレートはソースコードに直書きで管理
- 複数のテンプレートを一箇所にまとめて定義

## 5. 技術仕様

### 5.1 サーバー側技術スタック
- **言語**: Go 1.21+
- **AWS SDK**: `github.com/aws/aws-sdk-go-v2/service/ses`
- **gomail**: `gopkg.in/mail.v2` (Mailpit送信用)
- **テンプレート処理**: Go標準ライブラリの `text/template` または `html/template`

### 5.2 クライアント側技術スタック
- **フレームワーク**: Next.js 14+ (App Router)
- **言語**: TypeScript 5+

### 5.3 外部サービス
- **Mailpit**: Dockerコンテナで動作するメールテストツール
- **AWS SES**: 本番環境でのメール送信サービス

## 6. 受け入れ基準

### 6.1 機能要件
1. **メール送信機能**: 3つの送信方式（標準出力、Mailpit、AWS SES）が正しく動作すること
2. **環境別送信方式**: 環境変数に基づいて適切な送信方式が選択されること
3. **メールテンプレート機能**: テンプレートに動的な値を埋め込んでメール本文を生成できること
4. **APIエンドポイント**: `/api/email/send` エンドポイントが正しく動作すること
5. **クライアント画面**: メールアドレスと名前を入力してメールを送信できること
6. **Mailpit起動**: Mailpit起動スクリプトでMailpitを起動できること

### 6.2 非機能要件
1. **認証**: 認証なしでメール送信できないこと
2. **エラーハンドリング**: 適切なエラーメッセージが返されること
3. **ログ記録**: メール送信のログが記録されること

## 7. 制約事項

1. **メール形式**: テキストメールのみ（HTMLメールは将来的な拡張項目）
2. **添付ファイル**: 添付ファイルの送信は将来的な拡張項目
3. **送信履歴**: メール送信履歴の保存は将来的な拡張項目
4. **リトライ機能**: 送信失敗時の自動リトライ機能は将来的な拡張項目

## 8. 将来の拡張項目（現時点では未実装）

以下の機能は将来の拡張として検討されていますが、現時点では実装対象外です：

- HTMLメールの送信機能
- 添付ファイルの送信機能
- メール送信履歴の保存・管理機能
- メール送信のスケジューリング機能
- 送信失敗時の自動リトライ機能
- メール送信の非同期処理
- メール送信のバッチ処理

## Project Description (Input)

メール送信機能を実装する。

### 3つの送信方式
* 標準出力にメッセージを出力するだけのバージョン
* Mailpitを利用して、実際のメールを送信しないが、送信内容を確認する方法
* AWS SESを利用して、メールを送信する方法
    * 開発環境は標準出力の形式。staging、productionはAWS SESの方式をデフォルトの設定としたい。

### メールテンプレート
* メール本文のテンプレート設定を持っておいて、テンプレートの文面の一部を動的に置換して、送信メッセージを作り出したい。
    * テンプレートは設定ファイルでなくて、ソースコードに直書きでも良い。複数のテンプレートの定義があるとして、それらがまとまっていると嬉しい。

### 実装の方針 (構想)
* interfaceを用意
```
type EmailSender interface {
    Send(to []string, subject, body string) error
}
```

* 出力先にinterfaceを実装
```
type MockSender struct{}

func (s *MockSender) Send(to []string, subject, body string) error {
    fmt.Printf("[Mock Email] To: %v | Subject: %s\nBody: %s\n", to, subject, body)
    return nil
}
```

```
import (
    "context"
    "github.com/aws/aws-sdk-go-v2/service/ses"
    "github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type SESSender struct {
    client *ses.Client
}

func (s *SESSender) Send(to []string, subject, body string) error {
    input := &ses.SendEmailInput{
        Destination: &types.Destination{ToAddresses: to},
        Message: &types.Message{
            Subject: &types.Content{Data: &subject},
            Body:    &types.Body{Html: &types.Content{Data: &body}},
        },
        Source: aws.String("sender@example.com"),
    }
    _, err := s.client.SendEmail(context.TODO(), input)
    return err
}
```

### Mailpit
* MailpitはDockerで動作させる。
* 普段は停止しておいて、Mailpitを使いたいときに起動する。
* Mailpitの起動スクリプトが欲しい。

### このアプリでの処理の流れ
* クライアントアプリで、メールアドレス、名前を入力する。
* APIサーバーで受け取って、メールアドレス、名前の部分を置換して、メールを送信。
    * エンドポイントや画面の作りは提案を出して。
        * 参考コードとしての実装なのでエンドポイントは邪魔にならないURLがいい。

## Requirements

### Requirement 1: メール送信インターフェースの実装
**Objective:** As a system, I want to define a common interface for email sending, so that different email sending implementations can be used interchangeably.

#### Acceptance Criteria
1. WHEN an EmailSender interface is defined THEN it SHALL have a `Send(to []string, subject, body string) error` method
2. IF different email sending implementations are created THEN they SHALL all implement the EmailSender interface
3. WHERE the interface is used THEN the system SHALL allow switching between different implementations without changing calling code

### Requirement 2: 標準出力送信実装（MockSender）の実装
**Objective:** As a developer, I want to send emails to standard output in development, so that I can test email functionality without actually sending emails.

#### Acceptance Criteria
1. WHEN MockSender is implemented THEN it SHALL implement the EmailSender interface
2. IF Send method is called THEN the system SHALL output email content to standard output
3. WHERE email content is output THEN it SHALL use the format: `[Mock Email] To: %v | Subject: %s\nBody: %s\n`
4. WHEN Send method completes successfully THEN it SHALL return nil error
5. IF an error occurs THEN the system SHALL return an appropriate error

### Requirement 3: Mailpit送信実装（MailpitSender）の実装
**Objective:** As a developer, I want to send emails through Mailpit, so that I can view email content in a web interface without actually sending emails.

#### Acceptance Criteria
1. WHEN MailpitSender is implemented THEN it SHALL implement the EmailSender interface
2. IF Send method is called THEN the system SHALL send email to Mailpit using gomail library
3. WHERE gomail library is used THEN the system SHALL use `gopkg.in/mail.v2` package
4. WHERE Mailpit SMTP is configured THEN the system SHALL use `localhost:1025` as default SMTP address
5. WHEN email is sent to Mailpit THEN the system SHALL use SMTP protocol through gomail
6. IF Mailpit is not running THEN the system SHALL return an appropriate error
7. WHEN Send method completes successfully THEN it SHALL return nil error

### Requirement 4: AWS SES送信実装（SESSender）の実装
**Objective:** As a system, I want to send emails through AWS SES, so that emails can be delivered to recipients in production.

#### Acceptance Criteria
1. WHEN SESSender is implemented THEN it SHALL implement the EmailSender interface
2. IF Send method is called THEN the system SHALL send email using AWS SES SDK
3. WHERE AWS SES client is configured THEN the system SHALL use `github.com/aws/aws-sdk-go-v2/service/ses`
4. WHEN email is sent THEN the system SHALL include destination addresses, subject, and body
5. WHERE sender email address is configured THEN the system SHALL read it from configuration file
6. IF AWS credentials are invalid THEN the system SHALL return an appropriate error
7. WHEN Send method completes successfully THEN it SHALL return nil error

### Requirement 5: 環境別メール送信方式の選択
**Objective:** As a system, I want to select email sending method based on environment, so that appropriate method is used for each environment.

#### Acceptance Criteria
1. WHEN APP_ENV is "develop" THEN the system SHALL use MockSender as default
2. WHEN APP_ENV is "staging" or "production" THEN the system SHALL use SESSender as default
3. WHERE email sending method is configured in config file THEN the system SHALL use the configured method
4. IF environment variable and config file both specify method THEN the system SHALL prioritize config file setting
5. WHEN email service is initialized THEN the system SHALL select appropriate sender based on environment

### Requirement 6: メールテンプレート機能の実装
**Objective:** As a system, I want to use email templates with dynamic content, so that personalized emails can be sent efficiently.

#### Acceptance Criteria
1. WHEN email templates are defined THEN they SHALL be written directly in source code
2. IF multiple templates are needed THEN they SHALL be defined in a single location
3. WHERE template contains placeholders THEN the system SHALL replace them with actual values
4. WHEN template is used THEN the system SHALL allow selecting template by name
5. IF template name is invalid THEN the system SHALL return an appropriate error
6. WHERE template data is provided THEN the system SHALL replace all placeholders in template

### Requirement 7: メール送信APIエンドポイントの実装
**Objective:** As a client application, I want to send emails through API endpoint, so that emails can be sent from frontend.

#### Acceptance Criteria
1. WHEN POST request is sent to `/api/email/send` THEN the system SHALL process email sending request
2. IF request body contains recipient, template, and data THEN the system SHALL use them for email sending
3. WHERE authentication is required THEN the system SHALL verify Public API Key JWT or Auth0 JWT
4. WHEN authentication fails THEN the system SHALL return HTTP 401 Unauthorized
5. IF request body is invalid THEN the system SHALL return HTTP 400 Bad Request
6. WHEN email is sent successfully THEN the system SHALL return HTTP 200 OK with success message
7. IF email sending fails THEN the system SHALL return HTTP 500 Internal Server Error with error message

### Requirement 8: クライアント側メール送信画面の実装
**Objective:** As a user, I want to input email address and name to send email, so that I can send personalized emails easily.

#### Acceptance Criteria
1. WHEN user navigates to `/email/send` THEN the system SHALL display email sending form
2. IF form is displayed THEN it SHALL contain email address input field and name input field
3. WHERE user submits form THEN the system SHALL send POST request to `/api/email/send`
4. WHEN email is sent successfully THEN the system SHALL display success message
5. IF email sending fails THEN the system SHALL display error message
6. WHERE form validation is needed THEN the system SHALL validate email address format

### Requirement 9: MailpitのDocker設定と起動スクリプト
**Objective:** As a developer, I want to start Mailpit easily, so that I can test email functionality.

#### Acceptance Criteria
1. WHEN docker-compose.mailpit.yml is created THEN it SHALL define Mailpit service
2. IF Mailpit service is defined THEN it SHALL expose port 8025 for SMTP and Web UI
3. WHERE Mailpit is started THEN the system SHALL be accessible at `http://localhost:8025`
4. WHEN start-mailpit.sh script is created THEN it SHALL contain command to start Mailpit
5. IF script is executed THEN it SHALL start Mailpit using docker-compose
6. WHERE script has execute permission THEN it SHALL be executable

### Requirement 10: エラーハンドリングとログ記録
**Objective:** As a system, I want to handle errors gracefully and log important events, so that issues can be diagnosed and resolved.

#### Acceptance Criteria
1. WHEN email sending fails THEN the system SHALL log the error with appropriate context
2. IF Mailpit is not available THEN the system SHALL return appropriate error message
3. WHEN AWS SES credentials are invalid THEN the system SHALL return appropriate error message
4. IF template processing fails THEN the system SHALL return appropriate error message
5. WHERE errors occur THEN the system SHALL log them with sufficient context for debugging
