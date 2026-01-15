**[日本語](../ja/File-Upload.md) | [English]**

# File Upload Feature Guide

## Overview

This document explains the usage of the file upload feature in go-webdb-template.

The file upload feature supports large file uploads using the **TUS protocol** (Resumable Upload Protocol). With the TUS protocol, uploads can be resumed after network disconnection, making it suitable for large file uploads.

## Feature Description

### About TUS Protocol

TUS (Tus Resumable Upload Protocol) is an HTTP-based file upload protocol with the following features:

- **Resumable**: Can resume uploads after network disconnection
- **Large File Support**: Efficiently upload large files
- **Progress Tracking**: Track upload progress

### Storage Types

The following two storage types are supported:

1. **Local Storage (local)**: Saves to the server's local file system
2. **S3 Storage (s3)**: Saves to AWS S3

## Configuration

### Configuration File

Configuration is done in the `upload` section of `config/{env}/config.yaml`.

#### Local Storage Configuration Example

```yaml
upload:
  base_path: "/api/upload/dm_movie"  # TUS endpoint base path
  max_file_size: 2147483648          # Max file size (bytes, e.g., 2GB)
  allowed_extensions:                 # Allowed extension list
    - "mp4"
  storage:
    type: "local"
    local:
      path: "./uploads"               # Local save path
```

#### S3 Storage Configuration Example

```yaml
upload:
  base_path: "/api/upload/dm_movie"
  max_file_size: 2147483648
  allowed_extensions:
    - "mp4"
  storage:
    type: "s3"
    s3:
      bucket: "my-upload-bucket"     # S3 bucket name
      region: "us-east-1"            # AWS region
```

**Note**: When using S3 storage, AWS credentials are automatically retrieved from environment variables or AWS configuration files (`~/.aws/credentials`).

### CORS Configuration

For using the TUS protocol, the following headers must be added to CORS configuration:

```yaml
cors:
  allowed_headers:
    - Content-Type
    - Authorization
    - Tus-Resumable
    - Upload-Length
    - Upload-Offset
    - Upload-Metadata
  expose_headers:
    - Location
    - Upload-Offset
    - Upload-Length
    - Tus-Resumable
    - Tus-Version
    - Tus-Extension
    - Tus-Max-Size
```

## Usage

### Authentication

File upload endpoints require authentication. Use one of the following authentication methods:

- **Public API Key JWT**: `Authorization: Bearer <PUBLIC_API_KEY_JWT>`
- **Auth0 JWT**: `Authorization: Bearer <AUTH0_JWT>`

### TUS Protocol Basic Flow

#### 1. OPTIONS Request (Server Capability Check)

```bash
curl -X OPTIONS http://localhost:8080/api/upload/dm_movie \
  -H "Authorization: Bearer <YOUR_TOKEN>"
```

**Response Headers**:
- `Tus-Version`: Supported TUS protocol versions
- `Tus-Extension`: Supported extensions
- `Tus-Max-Size`: Maximum file size

#### 2. POST Request (Start Upload)

```bash
curl -X POST http://localhost:8080/api/upload/dm_movie \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -H "Upload-Length: 1048576" \
  -H "Upload-Metadata: filename dGVzdC5tcDQ="
```

**Request Headers**:
- `Upload-Length`: File size (bytes)
- `Upload-Metadata`: Base64-encoded metadata (e.g., `filename dGVzdC5tcDQ=` is `filename test.mp4`)

**Response**:
- `201 Created`: Upload was created
- `Location`: Upload resource URL (e.g., `/api/upload/dm_movie/abc123`)

#### 3. PATCH Request (Upload Data)

```bash
curl -X PATCH http://localhost:8080/api/upload/dm_movie/abc123 \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -H "Upload-Offset: 0" \
  -H "Content-Type: application/offset+octet-stream" \
  --data-binary @file.mp4
```

**Request Headers**:
- `Upload-Offset`: Current upload position (bytes)
- `Content-Type`: `application/offset+octet-stream`

**Response Headers**:
- `204 No Content`: Upload succeeded
- `Upload-Offset`: New upload position

#### 4. HEAD Request (Check Upload Status)

```bash
curl -X HEAD http://localhost:8080/api/upload/dm_movie/abc123 \
  -H "Authorization: Bearer <YOUR_TOKEN>"
```

**Response Headers**:
- `Upload-Offset`: Current upload position
- `Upload-Length`: Total file size

### Client-side Implementation Example (JavaScript)

```javascript
// Using TUS client library (e.g., tus-js-client)
import * as tus from 'tus-js-client';

const file = document.querySelector('input[type="file"]').files[0];

const upload = new tus.Upload(file, {
  endpoint: 'http://localhost:8080/api/upload/dm_movie',
  retryDelays: [0, 3000, 5000, 10000, 20000],
  headers: {
    'Authorization': 'Bearer <YOUR_TOKEN>'
  },
  metadata: {
    filename: file.name,
    filetype: file.type
  },
  onError: (error) => {
    console.error('Upload failed:', error);
  },
  onProgress: (bytesUploaded, bytesTotal) => {
    const percentage = (bytesUploaded / bytesTotal * 100).toFixed(2);
    console.log(`Upload progress: ${percentage}%`);
  },
  onSuccess: () => {
    console.log('Upload finished:', upload.url);
  }
});

// Start upload
upload.start();
```

## Error Responses

### 400 Bad Request

```json
{
  "error": "Invalid Upload-Length header"
}
```

**Cause**: Invalid `Upload-Length` header

### 400 Bad Request (Extension Error)

```json
{
  "error": "File extension 'avi' is not allowed. Allowed extensions: [mp4]"
}
```

**Cause**: Attempted to upload a file with an unauthorized extension

### 413 Request Entity Too Large

```json
{
  "error": "File size exceeds maximum allowed size of 2147483648 bytes"
}
```

**Cause**: File size exceeds `max_file_size`

### 401 Unauthorized

**Cause**: Authentication token is invalid or not set

## File Validation

The following validations are performed during upload:

1. **File Size**: Verifies that the size specified in `Upload-Length` header is within `max_file_size`
2. **Extension**: Verifies that the file extension from `Upload-Metadata` header is in `allowed_extensions`

## Upload Completion Handling

When upload is finished, the TUS handler automatically saves the file.

- **Local Storage**: Saved to the specified path (`storage.local.path`)
- **S3 Storage**: Saved to the specified S3 bucket

## Notes

1. **Authentication Required**: All TUS endpoints require authentication
2. **File Size Limit**: Files exceeding the size set in `max_file_size` cannot be uploaded
3. **Extension Restriction**: Files with extensions not in `allowed_extensions` cannot be uploaded
4. **Storage Path**: When using local storage, the directory at the specified path is automatically created if it doesn't exist
5. **S3 Authentication**: When using S3 storage, AWS credentials (access key, secret key) must be set in environment variables or AWS configuration files

## Related Documentation

- [TUS Protocol Specification](https://tus.io/protocols/resumable-upload.html)
- [tus-js-client](https://github.com/tus/tus-js-client) - JavaScript client library
