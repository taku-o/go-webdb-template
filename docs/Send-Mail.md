# メール送信機能利用手順

## 概要

このドキュメントでは、go-webdb-templateのメール送信機能の利用手順を説明します。

メール送信機能は、以下の3つの送信方式をサポートしています：

1. **標準出力送信（MockSender）**: 開発環境でメール内容を標準出力に出力
2. **Mailpit送信（MailpitSender）**: 開発環境でメール内容をMailpitで確認
3. **AWS SES送信（SESSender）**: 本番環境で実際にメールを送信

## 機能説明

### 送信方式の選択

送信方式は環境変数 `APP_ENV` と設定ファイルの `email.sender_type` に基づいて自動的に選択されます：

- **開発環境（`APP_ENV=develop`）**: デフォルトで `MockSender` を使用（標準出力に出力）
- **ステージング/本番環境（`APP_ENV=staging` または `production`）**: デフォルトで `SESSender` を使用（AWS SESで送信）

設定ファイルで `sender_type` を明示的に指定することで、送信方式を変更できます。

### テンプレート機能

メール本文はテンプレートを使用して生成されます。テンプレートには動的な値を埋め込むことができ、送信時に実際の値に置換されます。

## 利用方法

### クライアント側（Web画面）

#### 1. メール送信画面にアクセス

ブラウザで以下のURLにアクセスします：

```
http://localhost:3000/dm_email/send
```

#### 2. メール送信フォームに入力

1. **メールアドレス**: 送信先のメールアドレスを入力
2. **名前**: 受信者の名前を入力

#### 3. 送信ボタンをクリック

「送信」ボタンをクリックすると、メールが送信されます。

#### 4. 送信結果の確認

- **成功時**: 成功メッセージが表示されます
- **失敗時**: エラーメッセージが表示されます

### API経由での利用

#### エンドポイント

**POST** `/api/email/send`

#### 認証

このエンドポイントは認証が必要です。以下のいずれかの認証方式を使用してください：

- **Public API Key JWT**: `Authorization: Bearer <PUBLIC_API_KEY_JWT>`
- **Auth0 JWT**: `Authorization: Bearer <AUTH0_JWT>`

#### リクエスト例

```bash
curl -X POST http://localhost:8080/api/email/send \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -d '{
    "to": ["recipient@example.com"],
    "template": "welcome",
    "data": {
      "name": "John Doe",
      "email": "john@example.com"
    }
  }'
```

#### リクエストボディ

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

**フィールド説明**:
- `to` (必須): 送信先メールアドレスの配列
- `template` (必須): 使用するテンプレート名
- `data` (必須): テンプレートに埋め込むデータ（オブジェクト）

#### レスポンス例

**成功時（200 OK）**:
```json
{
  "success": true,
  "message": "Email sent successfully"
}
```

**エラー時（400 Bad Request）**:
```json
{
  "error": "Invalid email address format"
}
```

**エラー時（500 Internal Server Error）**:
```json
{
  "error": "Failed to send email"
}
```

## 環境別の設定と使い方

### 開発環境（標準出力送信）

開発環境では、デフォルトでメール内容が標準出力に出力されます。

#### 設定

`config/develop/config.yaml`:

```yaml
email:
  sender_type: "mock"  # デフォルト
  mock: {}
```

#### 使い方

1. サーバーを起動:
   ```bash
   APP_ENV=develop go run ./cmd/server/main.go
   ```

2. メールを送信（Web画面またはAPI経由）

3. サーバーの標準出力を確認:
   ```
   [Mock Email] To: [recipient@example.com] | Subject: ようこそ
   Body: こんにちは、John Doeさん。

   メールアドレス: john@example.com

   ご登録ありがとうございます。
   ```

### 開発環境（Mailpit送信）

開発環境でメール内容をWeb UIで確認したい場合は、Mailpitを使用します。

#### 1. Mailpitの起動

```bash
./scripts/start-mailpit.sh start
```

または、直接Docker Composeを使用:

```bash
docker-compose -f docker-compose.mailpit.yml up -d
```

#### 2. Mailpit Web UIにアクセス

ブラウザで以下のURLにアクセス:

```
http://localhost:8025
```

#### 3. 設定を変更

`config/develop/config.yaml`:

```yaml
email:
  sender_type: "mailpit"  # Mailpitを使用
  mailpit:
    smtp_host: "localhost"
    smtp_port: 1025
```

#### 4. サーバーを再起動

設定変更後、サーバーを再起動します。

#### 5. メール送信

メールを送信すると、Mailpit Web UIでメール内容を確認できます。

#### 6. Mailpitの停止

```bash
./scripts/start-mailpit.sh stop
```

または:

```bash
docker-compose -f docker-compose.mailpit.yml down
```

### 本番環境（AWS SES送信）

本番環境では、AWS SESを使用して実際にメールを送信します。

#### 1. AWS SESの設定

1. AWS SESで送信元メールアドレスを検証
2. AWS認証情報を設定（環境変数またはAWS設定ファイル）

#### 2. 設定ファイルの設定

`config/production/config.yaml`:

```yaml
email:
  sender_type: "ses"  # デフォルト
  ses:
    from: "sender@example.com"  # 検証済みの送信元メールアドレス
    region: "us-east-1"  # AWSリージョン
```

#### 3. サーバーを起動

```bash
APP_ENV=production go run ./cmd/server/main.go
```

#### 4. メール送信

メールを送信すると、AWS SES経由で実際にメールが送信されます。

## テンプレートの使い方

### 利用可能なテンプレート

現在、以下のテンプレートが利用可能です：

- `welcome`: ウェルカムメール

### テンプレートの使用例

#### welcomeテンプレート

**リクエスト**:
```json
{
  "to": ["user@example.com"],
  "template": "welcome",
  "data": {
    "name": "山田太郎",
    "email": "user@example.com"
  }
}
```

**生成されるメール本文**:
```
こんにちは、山田太郎さん。

メールアドレス: user@example.com

ご登録ありがとうございます。
```

**件名**: `ようこそ`

### 新しいテンプレートの追加

新しいテンプレートを追加する場合は、`server/internal/service/email/template.go` を編集してください。

## トラブルシューティング

### メールが送信されない

#### 開発環境（MockSender）

- サーバーの標準出力を確認してください
- ログにエラーメッセージがないか確認してください

#### 開発環境（Mailpit）

1. **Mailpitが起動しているか確認**:
   ```bash
   docker ps | grep mailpit
   ```

2. **Mailpit Web UIにアクセスできるか確認**:
   ```
   http://localhost:8025
   ```

3. **設定ファイルの確認**:
   - `sender_type` が `"mailpit"` になっているか
   - `smtp_host` が `"localhost"` になっているか
   - `smtp_port` が `1025` になっているか

4. **サーバーログを確認**:
   - Mailpit接続エラーがないか確認

#### 本番環境（AWS SES）

1. **AWS認証情報の確認**:
   - 環境変数またはAWS設定ファイルに認証情報が設定されているか
   - 認証情報が有効か

2. **送信元メールアドレスの確認**:
   - AWS SESで送信元メールアドレスが検証されているか
   - 設定ファイルの `from` が正しいか

3. **AWS SESの制限確認**:
   - サンドボックス環境の場合、検証済みメールアドレスにのみ送信可能
   - 送信制限に達していないか

4. **サーバーログを確認**:
   - AWS SESエラーがないか確認

### 認証エラー（401 Unauthorized）

- 認証トークンが正しく設定されているか確認
- トークンが有効期限内か確認
- Public API Key JWT または Auth0 JWT が正しく設定されているか確認

### バリデーションエラー（400 Bad Request）

- リクエストボディの形式が正しいか確認
- メールアドレスの形式が正しいか確認
- 必須フィールド（`to`, `template`, `data`）がすべて含まれているか確認

### テンプレートエラー（400 Bad Request）

- テンプレート名が正しいか確認
- テンプレートデータに必要なデータがすべて含まれているか確認

### Mailpit接続エラー（503 Service Unavailable）

- Mailpitが起動しているか確認
- SMTPポート（1025）が正しく設定されているか確認
- ファイアウォールやネットワーク設定を確認

## 関連ドキュメント

- [API Documentation](./API.md): APIエンドポイントの詳細
- [Architecture](./Architecture.md): システムアーキテクチャの説明

## 参考情報

### Mailpit

- **Web UI**: http://localhost:8025
- **SMTP**: localhost:1025
- **公式ドキュメント**: https://github.com/axllent/mailpit

### AWS SES

- **AWS SES ドキュメント**: https://docs.aws.amazon.com/ses/
- **送信元メールアドレスの検証**: AWS SESコンソールで実施
