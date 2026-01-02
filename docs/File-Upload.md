# ファイルアップロード機能利用手順

## 概要

このドキュメントでは、go-webdb-templateのファイルアップロード機能の利用手順を説明します。

ファイルアップロード機能は、**TUSプロトコル**（Resumable Upload Protocol）を使用した大容量ファイルのアップロードをサポートしています。TUSプロトコルにより、ネットワーク切断時でもアップロードを再開できるため、大容量ファイルのアップロードに適しています。

## 機能説明

### TUSプロトコルについて

TUS（Tus Resumable Upload Protocol）は、HTTPベースのファイルアップロードプロトコルです。以下の特徴があります：

- **再開可能**: ネットワーク切断後もアップロードを再開できる
- **大容量ファイル対応**: 大きなファイルを効率的にアップロードできる
- **進捗追跡**: アップロードの進捗状況を追跡できる

### ストレージタイプ

以下の2つのストレージタイプをサポートしています：

1. **ローカルストレージ（local）**: サーバーのローカルファイルシステムに保存
2. **S3ストレージ（s3）**: AWS S3に保存

## 設定方法

### 設定ファイル

設定は`config/{env}/config.yaml`の`upload`セクションで行います。

#### ローカルストレージ設定例

```yaml
upload:
  base_path: "/api/upload/dm_movie"  # TUSエンドポイントのベースパス
  max_file_size: 2147483648          # 最大ファイルサイズ（バイト、例: 2GB）
  allowed_extensions:                 # 許可された拡張子リスト
    - "mp4"
  storage:
    type: "local"
    local:
      path: "./uploads"               # ローカル保存パス
```

#### S3ストレージ設定例

```yaml
upload:
  base_path: "/api/upload/dm_movie"
  max_file_size: 2147483648
  allowed_extensions:
    - "mp4"
  storage:
    type: "s3"
    s3:
      bucket: "my-upload-bucket"     # S3バケット名
      region: "us-east-1"            # AWSリージョン
```

**注意**: S3ストレージを使用する場合、AWS認証情報は環境変数またはAWS設定ファイル（`~/.aws/credentials`）から自動的に取得されます。

### CORS設定

TUSプロトコルを使用するため、CORS設定に以下のヘッダーを追加する必要があります：

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

## 利用方法

### 認証

ファイルアップロードエンドポイントは認証が必要です。以下のいずれかの認証方式を使用してください：

- **Public API Key JWT**: `Authorization: Bearer <PUBLIC_API_KEY_JWT>`
- **Auth0 JWT**: `Authorization: Bearer <AUTH0_JWT>`

### TUSプロトコルの基本フロー

#### 1. OPTIONSリクエスト（サーバー機能確認）

```bash
curl -X OPTIONS http://localhost:8080/api/upload/dm_movie \
  -H "Authorization: Bearer <YOUR_TOKEN>"
```

**レスポンスヘッダー**:
- `Tus-Version`: サポートされているTUSプロトコルバージョン
- `Tus-Extension`: サポートされている拡張機能
- `Tus-Max-Size`: 最大ファイルサイズ

#### 2. POSTリクエスト（アップロード開始）

```bash
curl -X POST http://localhost:8080/api/upload/dm_movie \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -H "Upload-Length: 1048576" \
  -H "Upload-Metadata: filename dGVzdC5tcDQ="
```

**リクエストヘッダー**:
- `Upload-Length`: ファイルサイズ（バイト）
- `Upload-Metadata`: Base64エンコードされたメタデータ（例: `filename dGVzdC5tcDQ=` は `filename test.mp4`）

**レスポンス**:
- `201 Created`: アップロードが作成された
- `Location`: アップロードリソースのURL（例: `/api/upload/dm_movie/abc123`）

#### 3. PATCHリクエスト（データアップロード）

```bash
curl -X PATCH http://localhost:8080/api/upload/dm_movie/abc123 \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -H "Upload-Offset: 0" \
  -H "Content-Type: application/offset+octet-stream" \
  --data-binary @file.mp4
```

**リクエストヘッダー**:
- `Upload-Offset`: 現在のアップロード位置（バイト）
- `Content-Type`: `application/offset+octet-stream`

**レスポンスヘッダー**:
- `204 No Content`: アップロードが成功した
- `Upload-Offset`: 新しいアップロード位置

#### 4. HEADリクエスト（アップロード状態確認）

```bash
curl -X HEAD http://localhost:8080/api/upload/dm_movie/abc123 \
  -H "Authorization: Bearer <YOUR_TOKEN>"
```

**レスポンスヘッダー**:
- `Upload-Offset`: 現在のアップロード位置
- `Upload-Length`: ファイル全体のサイズ

### クライアント側の実装例（JavaScript）

```javascript
// TUSクライアントライブラリを使用（例: tus-js-client）
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

// アップロード開始
upload.start();
```

## エラーレスポンス

### 400 Bad Request

```json
{
  "error": "Invalid Upload-Length header"
}
```

**原因**: `Upload-Length`ヘッダーが無効

### 400 Bad Request（拡張子エラー）

```json
{
  "error": "File extension 'avi' is not allowed. Allowed extensions: [mp4]"
}
```

**原因**: 許可されていない拡張子のファイルをアップロードしようとした

### 413 Request Entity Too Large

```json
{
  "error": "File size exceeds maximum allowed size of 2147483648 bytes"
}
```

**原因**: ファイルサイズが`max_file_size`を超えている

### 401 Unauthorized

**原因**: 認証トークンが無効または未設定

## ファイル検証

アップロード時に以下の検証が行われます：

1. **ファイルサイズ**: `Upload-Length`ヘッダーで指定されたサイズが`max_file_size`以下であることを確認
2. **拡張子**: `Upload-Metadata`ヘッダーから取得したファイル名の拡張子が`allowed_extensions`に含まれていることを確認

## アップロード完了時の処理

アップロードが完了すると、TUSハンドラーが自動的にファイルを保存します。

- **ローカルストレージ**: 指定されたパス（`storage.local.path`）に保存
- **S3ストレージ**: 指定されたS3バケットに保存

## 注意事項

1. **認証必須**: すべてのTUSエンドポイントは認証が必要です
2. **ファイルサイズ制限**: `max_file_size`で設定されたサイズを超えるファイルはアップロードできません
3. **拡張子制限**: `allowed_extensions`に含まれていない拡張子のファイルはアップロードできません
4. **ストレージパス**: ローカルストレージを使用する場合、指定されたパスのディレクトリが存在しない場合は自動的に作成されます
5. **S3認証**: S3ストレージを使用する場合、AWS認証情報（アクセスキー、シークレットキー）が環境変数またはAWS設定ファイルに設定されている必要があります

## 関連ドキュメント

- [TUSプロトコル仕様](https://tus.io/protocols/resumable-upload.html)
- [tus-js-client](https://github.com/tus/tus-js-client) - JavaScriptクライアントライブラリ
