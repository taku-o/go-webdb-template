**[日本語](../ja/Logging.md) | [English]**

# Logging Feature Guide

## Overview

This document explains the usage of the logging feature in go-webdb-template.

This project supports the following three types of logs:

1. **Access Log**: HTTP request/response logs
2. **Mail Sending Log**: Email sending process logs
3. **SQL Log**: Database query logs

## Configuration

### Configuration File

Logging configuration is done in the `logging` section of `config/{env}/config.yaml`.

```yaml
logging:
  level: debug                    # Log level (debug, info, warn, error)
  format: json                    # Log format (json, text)
  output: stdout                  # Output destination (stdout, file)
  output_dir: ../logs             # Log file output directory
  sql_log_enabled: true           # SQL log enabled/disabled
  mail_log_enabled: true          # Mail sending log enabled/disabled
  mail_log_output_dir: ../logs    # Mail log output directory (optional, uses output_dir if not set)
```

### Log Levels

- `debug`: All logs including debug information
- `info`: Information level logs
- `warn`: Warning level logs
- `error`: Error level logs

### Log Formats

- `json`: JSON format output (structured logs)
- `text`: Text format output (human-readable format)

### Output Destinations

- `stdout`: Output to standard output
- `file`: Output to file (directory specified in `output_dir`)

## Access Log

### Overview

Access log records information about all HTTP requests and responses.

### Output Format (JSON)

```json
{
  "timestamp": "2025-01-27 10:30:45",
  "method": "POST",
  "path": "/api/users",
  "protocol": "HTTP/1.1",
  "status_code": 201,
  "response_time_ms": 45.23,
  "remote_ip": "192.168.1.100",
  "user_agent": "Mozilla/5.0...",
  "headers": "Content-Type: application/json; Authorization: Bearer ...",
  "request_body": "{\"name\":\"John\",\"email\":\"john@example.com\"}"
}
```

**Field Descriptions**:
- `timestamp`: Request received time
- `method`: HTTP method (GET, POST, PUT, DELETE, etc.)
- `path`: Request path
- `protocol`: HTTP protocol version
- `status_code`: HTTP status code
- `response_time_ms`: Response time (milliseconds)
- `remote_ip`: Client IP address
- `user_agent`: User agent
- `headers`: Request headers (semicolon separated)
- `request_body`: Request body (POST/PUT/PATCH only, max 1MB)

### Request Body Logging Conditions

Request body is logged only when the following conditions are met:

1. HTTP method is POST, PUT, or PATCH
2. Content-Length is 1MB or less
3. Content-Type is not image/video/audio/multipart

### Log File Location

- Output: `{output_dir}/access.log`
- Rotation: Daily file split (`access-2025-01-27.log` format)

## Mail Sending Log

### Overview

Mail sending log records success/failure of email sending processes.

### Configuration

To enable mail sending log, set `mail_log_enabled: true`.

```yaml
logging:
  mail_log_enabled: true
  mail_log_output_dir: ../logs  # Optional, uses output_dir if not set
```

### Output Format (JSON)

```json
{
  "timestamp": "2025-01-27 10:30:45",
  "to": ["recipient@example.com"],
  "subject": "Welcome",
  "body": "Hello, John!",
  "body_truncated": false,
  "sender_type": "mock",
  "success": true,
  "error": ""
}
```

**On Error**:

```json
{
  "timestamp": "2025-01-27 10:30:45",
  "to": ["recipient@example.com"],
  "subject": "Welcome",
  "body": "Hello, John!",
  "body_truncated": false,
  "sender_type": "ses",
  "success": false,
  "error": "failed to send email: connection timeout"
}
```

**Field Descriptions**:
- `timestamp`: Email sending time
- `to`: Array of recipient email addresses
- `subject`: Email subject
- `body`: Email body
- `body_truncated`: Whether email body was truncated (if too long)
- `sender_type`: Sending method (`mock`, `mailpit`, `ses`)
- `success`: Success/failure
- `error`: Error message (only on failure)

### Log File Location

- Output: `{mail_log_output_dir}/mail.log` (or `{output_dir}/mail.log` if not set)
- Rotation: Daily file split (`mail-2025-01-27.log` format)

## SQL Log

### Overview

SQL log records database query execution. Queries using GORM are logged.

### Configuration

To enable SQL log, set `sql_log_enabled: true`.

```yaml
logging:
  sql_log_enabled: true
  output_dir: ../logs
```

### Output Format (Text)

```
[2025-01-27 10:30:45] [INFO] [master] [postgres] SELECT * FROM "dm_news" WHERE "dm_news"."id" = 'abc123' LIMIT 1
[2025-01-27 10:30:45] [INFO] [sharding] [postgres] INSERT INTO "dm_users_000" ("id","name","email","created_at","updated_at") VALUES ('def456','John','john@example.com','2025-01-27 10:30:45','2025-01-27 10:30:45')
```

**Format**: `[Time] [Level] [Group Name] [Driver] SQL Query`

### Log File Location

- Output: `{output_dir}/sql-{date}.log` (e.g., `sql-2025-01-27.log`)
- Rotation: Daily file split (1 file per day)

### Log Level

SQL logs are output at `Info` level. Follows GORM's log level settings.

## Log Rotation

### Access Log & Mail Sending Log

Log rotation is implemented using the `lumberjack` library.

- **Daily Split**: 1 file per day (`access-2025-01-27.log` format)
- **Size Limit**: None (daily split only)
- **Retention Period**: Unlimited (manual deletion required)

### SQL Log

- **Daily Split**: 1 file per day (`sql-2025-01-27.log` format)
- **Retention Period**: Unlimited (manual deletion required)

## Checking Log Files

### Development Environment

```bash
# Check access log
tail -f logs/access-$(date +%Y-%m-%d).log

# Check mail sending log
tail -f logs/mail-$(date +%Y-%m-%d).log

# Check SQL log
tail -f logs/sql-$(date +%Y-%m-%d).log
```

### Format JSON Logs for Viewing

```bash
# Format JSON logs using jq
tail -f logs/access-$(date +%Y-%m-%d).log | jq .

# Filter by specific conditions
tail -f logs/access-$(date +%Y-%m-%d).log | jq 'select(.status_code >= 400)'
```

## Performance Impact

### Access Log

- Request body reading is done in memory
- Large request bodies (over 1MB) are not logged
- File I/O is done asynchronously, but may impact performance under high load

### Mail Sending Log

- Executed as part of email sending process, minimal overhead
- Email body is truncated if too long

### SQL Log

- All SQL queries are logged, so log files may grow large under high load
- Setting `sql_log_enabled: false` in production is advised

## Production Environment Suggested Settings

```yaml
logging:
  level: info                    # Don't output debug info
  format: json                   # Structured logs for easy integration with log analysis tools
  output: file                   # Output to file
  output_dir: /var/log/app       # Specify appropriate log directory
  sql_log_enabled: false         # Disable SQL log (for performance and security)
  mail_log_enabled: true         # Enable mail sending log (for auditing)
```

## Notes

1. **Log File Disk Space**: Log files are not automatically deleted, periodic cleanup is required
2. **Security**: Logs may contain auth tokens and request bodies, set appropriate access controls
3. **Performance**: Log output may impact performance in high-load environments
4. **SQL Log**: Disabling SQL log in production is advised (for security and performance)
