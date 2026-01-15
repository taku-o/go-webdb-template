**[日本語](../ja/Send-Mail.md) | [English]**

# Email Sending Feature Guide

## Overview

This document explains the usage of the email sending feature in go-webdb-template.

The email sending feature supports the following three sending methods:

1. **Standard Output (MockSender)**: Outputs email content to standard output in development environment
2. **Mailpit (MailpitSender)**: View email content in Mailpit in development environment
3. **AWS SES (SESSender)**: Send actual emails in production environment

## Feature Description

### Sender Type Selection

The sender type is automatically selected based on the `APP_ENV` environment variable and the `email.sender_type` configuration:

- **Development environment (`APP_ENV=develop`)**: Uses `MockSender` by default (outputs to standard output)
- **Staging/Production environment (`APP_ENV=staging` or `production`)**: Uses `SESSender` by default (sends via AWS SES)

You can change the sender type by explicitly specifying `sender_type` in the configuration file.

### Template Feature

Email body is generated using templates. Templates can include dynamic values that are replaced with actual values at send time.

## Usage

### Client-Side (Web Interface)

#### 1. Access Email Sending Page

Access the following URL in your browser:

```
http://localhost:3000/dm_email/send
```

#### 2. Fill in Email Form

1. **Email Address**: Enter the recipient's email address
2. **Name**: Enter the recipient's name

#### 3. Click Send Button

Click the "Send" button to send the email.

#### 4. Check Send Result

- **Success**: Success message is displayed
- **Failure**: Error message is displayed

### API Usage

#### Endpoint

**POST** `/api/email/send`

#### Authentication

This endpoint requires authentication. Use one of the following authentication methods:

- **Public API Key JWT**: `Authorization: Bearer <PUBLIC_API_KEY_JWT>`
- **Auth0 JWT**: `Authorization: Bearer <AUTH0_JWT>`

#### Request Example

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

#### Request Body

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

**Field Description**:
- `to` (required): Array of recipient email addresses
- `template` (required): Template name to use
- `data` (required): Data to embed in the template (object)

#### Response Examples

**Success (200 OK)**:
```json
{
  "success": true,
  "message": "Email sent successfully"
}
```

**Error (400 Bad Request)**:
```json
{
  "error": "Invalid email address format"
}
```

**Error (500 Internal Server Error)**:
```json
{
  "error": "Failed to send email"
}
```

## Environment-Specific Configuration and Usage

### Development Environment (Standard Output)

In development environment, email content is output to standard output by default.

#### Configuration

`config/develop/config.yaml`:

```yaml
email:
  sender_type: "mock"  # default
  mock: {}
```

#### Usage

1. Start server:
   ```bash
   APP_ENV=develop go run ./cmd/server/main.go
   ```

2. Send email (via web interface or API)

3. Check server's standard output:
   ```
   [Mock Email] To: [recipient@example.com] | Subject: Welcome
   Body: Hello, John Doe.

   Email: john@example.com

   Thank you for registering.
   ```

### Development Environment (Mailpit)

Use Mailpit when you want to view email content in a web UI in development environment.

#### 1. Start Mailpit

```bash
./scripts/start-mailpit.sh start
```

Or use Docker Compose directly:

```bash
docker-compose -f docker-compose.mailpit.yml up -d
```

#### 2. Access Mailpit Web UI

Access the following URL in your browser:

```
http://localhost:8025
```

#### 3. Update Configuration

`config/develop/config.yaml`:

```yaml
email:
  sender_type: "mailpit"  # Use Mailpit
  mailpit:
    smtp_host: "localhost"
    smtp_port: 1025
```

#### 4. Restart Server

Restart the server after configuration changes.

#### 5. Send Email

When you send an email, you can view its content in the Mailpit Web UI.

#### 6. Stop Mailpit

```bash
./scripts/start-mailpit.sh stop
```

Or:

```bash
docker-compose -f docker-compose.mailpit.yml down
```

### Production Environment (AWS SES)

In production environment, actual emails are sent using AWS SES.

#### 1. Configure AWS SES

1. Verify sender email address in AWS SES
2. Configure AWS credentials (environment variables or AWS config file)

#### 2. Configure Settings File

`config/production/config.yaml`:

```yaml
email:
  sender_type: "ses"  # default
  ses:
    from: "sender@example.com"  # Verified sender email address
    region: "us-east-1"  # AWS region
```

#### 3. Start Server

```bash
APP_ENV=production go run ./cmd/server/main.go
```

#### 4. Send Email

When you send an email, it will be sent via AWS SES.

## Template Usage

### Available Templates

Currently, the following templates are available:

- `welcome`: Welcome email

### Template Usage Examples

#### welcome Template

**Request**:
```json
{
  "to": ["user@example.com"],
  "template": "welcome",
  "data": {
    "name": "Taro Yamada",
    "email": "user@example.com"
  }
}
```

**Generated Email Body**:
```
Hello, Taro Yamada.

Email: user@example.com

Thank you for registering.
```

**Subject**: `Welcome`

### Adding New Templates

To add a new template, edit `server/internal/service/email/template.go`.

## Troubleshooting

### Email Not Being Sent

#### Development Environment (MockSender)

- Check server's standard output
- Check logs for error messages

#### Development Environment (Mailpit)

1. **Verify Mailpit is running**:
   ```bash
   docker ps | grep mailpit
   ```

2. **Verify Mailpit Web UI is accessible**:
   ```
   http://localhost:8025
   ```

3. **Check configuration file**:
   - Is `sender_type` set to `"mailpit"`?
   - Is `smtp_host` set to `"localhost"`?
   - Is `smtp_port` set to `1025`?

4. **Check server logs**:
   - Check for Mailpit connection errors

#### Production Environment (AWS SES)

1. **Verify AWS credentials**:
   - Are credentials set in environment variables or AWS config file?
   - Are credentials valid?

2. **Verify sender email address**:
   - Is sender email address verified in AWS SES?
   - Is `from` correct in configuration file?

3. **Check AWS SES limits**:
   - In sandbox environment, can only send to verified email addresses
   - Check if sending limit has been reached

4. **Check server logs**:
   - Check for AWS SES errors

### Authentication Error (401 Unauthorized)

- Verify authentication token is correctly set
- Verify token is not expired
- Verify Public API Key JWT or Auth0 JWT is correctly configured

### Validation Error (400 Bad Request)

- Verify request body format is correct
- Verify email address format is correct
- Verify all required fields (`to`, `template`, `data`) are included

### Template Error (400 Bad Request)

- Verify template name is correct
- Verify all required data is included in template data

### Mailpit Connection Error (503 Service Unavailable)

- Verify Mailpit is running
- Verify SMTP port (1025) is correctly configured
- Check firewall and network settings

## Related Documentation

- [API Documentation](./API.md): API endpoint details
- [Architecture](./Architecture.md): System architecture description

## Reference Information

### Mailpit

- **Web UI**: http://localhost:8025
- **SMTP**: localhost:1025
- **Official Documentation**: https://github.com/axllent/mailpit

### AWS SES

- **AWS SES Documentation**: https://docs.aws.amazon.com/ses/
- **Sender Email Verification**: Perform in AWS SES Console
