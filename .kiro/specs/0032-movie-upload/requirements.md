# 動画ファイルアップロード機能要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #64
- **Issueタイトル**: ファイルアップロード機能の実装
- **Feature名**: 0032-movie-upload
- **作成日**: 2025-01-27

### 1.2 目的
大容量動画ファイルのアップロード機能を実装する。TUSプロトコルを使用してアップロードの中断と再開をサポートし、ユーザーが大容量ファイルを安全にアップロードできるようにする。

### 1.3 スコープ
- クライアント側: `/dm_movie/upload` ページの実装
- APIサーバー側: `/api/upload/dm_movie` エンドポイントの実装（TUSプロトコル対応）
- フロントエンド: uppy または tus-js-client を使用
- バックエンド: tusd ライブラリを使用
- 開発環境: ローカルファイルシステム保存
- 本番環境: AWS S3保存

**本実装の範囲外**:
- 動画ファイルの再生機能
- 動画ファイルの編集機能
- 動画ファイルの削除機能
- 動画ファイルの一覧表示機能
- 他のファイルタイプ（画像、ドキュメントなど）のアップロード機能

## 2. 背景・現状分析

### 2.1 現在の実装
- **アップロード機能**: 現在、ファイルアップロード機能は実装されていない
- **TUSプロトコル**: 現在、TUSプロトコルサポートは実装されていない
- **ファイルストレージ**: 現在、ファイルストレージ機能は実装されていない

### 2.2 課題点
1. **大容量ファイルアップロードの不在**: 動画ファイルなどの大容量ファイルをアップロードする機能が存在しない
2. **中断・再開機能の不在**: ネットワークエラーなどでアップロードが中断された場合、最初からやり直す必要がある
3. **専用アップロード画面の不在**: ファイルアップロード専用のユーザーインターフェースが存在しない

### 2.3 本実装による改善点
1. **大容量ファイルアップロード機能の提供**: TUSプロトコルを使用して大容量ファイルを安全にアップロードできるようになる
2. **中断・再開機能の実装**: ネットワークエラーなどで中断されても、最後に成功したチャンクから再開できる
3. **専用アップロード画面の提供**: ユーザーフレンドリーなアップロードインターフェースを提供
4. **環境別ストレージ対応**: 開発環境ではローカルファイルシステム、本番環境ではAWS S3を使用

## 3. 機能要件

### 3.1 クライアント側の実装

#### 3.1.1 アップロード画面の実装
- **ファイル**: `client/src/app/dm_movie/upload/page.tsx`
- **URL**: `/dm_movie/upload`
- **実装内容**:
  - uppy または tus-js-client を使用したアップロードUI
  - ファイル選択機能
  - アップロード進捗表示
  - エラー表示機能

#### 3.1.2 TUSクライアントの設定
- **実装内容**:
  - uppy を使用する場合: `@uppy/core` と `@uppy/tus` を使用
  - tus-js-client を使用する場合: `tus-js-client` ライブラリを使用
  - エンドポイント: `http://localhost:8080/api/upload/dm_movie` (開発環境)
  - チャンクサイズ: 5MB
  - リトライ遅延: [0, 1000, 3000, 5000] ms

### 3.2 APIサーバー側の実装

#### 3.2.1 TUSエンドポイントの実装
- **ファイル**: `server/internal/api/handler/upload_handler.go` (新規作成)
- **エンドポイント**: `/api/upload/dm_movie` (TUSプロトコル対応)
- **アクセスレベル**: `public` (Public API Key JWT または Auth0 JWT でアクセス可能)
- **実装内容**:
  - tusd ライブラリ (`github.com/tus/tusd/v2/pkg/handler`) を使用
  - TUSプロトコル (OPTIONS, POST, PATCH, HEAD, DELETE) をサポート
  - アップロード完了時のHook処理（PostFinishフックを使用）
  - ファイルサイズ制限の検証（設定可能な最大ファイルサイズ）
  - ファイル拡張子の制限（許可された拡張子のみアップロード可能）

#### 3.2.2 ファイルストレージの実装
- **開発環境**:
  - ローカルファイルシステム (`./uploads`)
  - `.gitignore` に `uploads/` を追加
- **本番環境**:
  - AWS S3 を使用
  - S3バケットとクレデンシャルの設定

#### 3.2.3 アップロード完了時のHook処理
- **実装内容**:
  - tusdのHookシステムを使用してアップロード完了時の処理を実装
  - `hooks.NewHandler()`でHookハンドラーを作成
  - `PostFinish`フックを設定してアップロード完了時に処理を実行
  - 現時点では処理は空実装（将来的にアップロード完了状態の表示などに使用予定）
  - HookイベントからファイルID、パス、サイズなどの情報を取得可能

## 4. 非機能要件

### 4.1 パフォーマンス
- チャンクサイズ: 5MB単位でアップロード
- リトライ機能: 指数バックオフによる自動リトライ

### 4.2 セキュリティ
- 認証: Public API Key JWT または Auth0 JWT による認証が必要
- ファイルサイズ制限: 設定可能な最大ファイルサイズを超えるファイルのアップロードを拒否
- ファイル拡張子制限: 許可された拡張子のみアップロード可能（動画ファイル形式を想定）
- ファイル検証: アップロードされたファイルの検証は必要だが、処理が重いため本番環境ではAWSの機能を使用してS3上でアップロード完了後に実施する。コード上での実装は不要。

### 4.3 可用性
- 中断・再開機能: ネットワークエラーなどで中断されても再開可能
- エラーハンドリング: 適切なエラーメッセージとログ記録

## 5. 技術仕様

### 5.1 フロントエンド技術スタック
- **ライブラリ**: uppy または tus-js-client
- **フレームワーク**: Next.js 14+ (App Router)
- **言語**: TypeScript 5+

### 5.2 バックエンド技術スタック
- **ライブラリ**: `github.com/tus/tusd/v2/pkg/handler`
- **言語**: Go 1.21+
- **プロトコル**: TUS (Resumable Upload Protocol)

### 5.3 ストレージ
- **開発環境**: ローカルファイルシステム
- **本番環境**: AWS S3

## 6. 受け入れ基準

### 6.1 機能要件
1. **アップロード機能**: ユーザーが動画ファイルを選択し、アップロードを開始できること
2. **中断・再開機能**: アップロードが中断されても、最後に成功したチャンクから再開できること
3. **進捗表示**: アップロードの進捗がリアルタイムで表示されること
4. **エラーハンドリング**: エラーが発生した場合、適切なエラーメッセージが表示されること
5. **認証**: Public API Key JWT または Auth0 JWT による認証が正しく機能すること
6. **ファイルサイズ制限**: 設定された最大ファイルサイズを超えるファイルのアップロードが拒否されること
7. **ファイル拡張子制限**: 許可された拡張子以外のファイルのアップロードが拒否されること

### 6.2 非機能要件
1. **パフォーマンス**: 5MBチャンクで効率的にアップロードできること
2. **セキュリティ**: 認証なしでアップロードできないこと
3. **可用性**: ネットワークエラーなどで中断されても再開できること

## 7. 制約事項

1. **ファイルタイプ**: 動画ファイルを想定（将来的に他のファイルタイプにも対応可能）
2. **同時アップロード数**: 現時点では1ファイルずつ（将来的な拡張項目）

## 8. 将来の拡張項目（現時点では未実装）

以下の機能は将来の拡張として検討されていますが、現時点では実装対象外です：

- 動画ファイルの再生機能
- 動画ファイルの編集機能
- 動画ファイルの削除機能
- 動画ファイルの一覧表示機能
- アップロード速度の制限
- アップロード履歴の管理

## Project Description (Input)

ファイルアップロード機能を実装する。
ただし、想定する対象は動画ファイルであり、
大容量であるため、アップロードの中断と再開をサポートしたい。

### フロントエンド
* uppy または tus-js-client を使用する。
```
import Uppy from '@uppy/core';
import Tus from '@uppy/tus';
import { Dashboard } from '@uppy/react';

const uppy = new Uppy().use(Tus, {
  endpoint: 'http://localhost:8080/api/upload/dm_movie', // Goバックエンド
  retryDelays: [0, 1000, 3000, 5000],
  chunkSize: 5 * 1024 * 1024, // 5MB単位
});

// コンポーネント内で <Dashboard uppy={uppy} /> を表示
```

### バックエンド
* 公式の tusd ライブラリをハンドラーとして組み込む
    * github.com/tus/tusd/v2/pkg/handler
```
// 簡易的な実装例
store := filestore.FileStore{Path: "./uploads"}
composer := handler.NewStoreComposer()
store.UseIn(composer)

// HTTPフックサーバーを設定
hookHandler := hooks.NewHandler()
hookHandler.PostFinish = func(event hooks.HookEvent) {
    // アップロード完了時の処理
    // ファイルID、パス、サイズなどの情報が取得可能
    // 現時点では空実装（将来的にアップロード完了状態の表示などに使用予定）
}

handler, err := handler.NewHandler(handler.Config{
    BasePath:              "/api/upload/dm_movie",
    StoreComposer:         composer,
})
```

### ファイルの保存先
* 開発環境ではローカルファイル。保存先は.gitignoreに追加する。
* 本番環境、staging環境では、AWS S3か、Tencent Cloudの類似機能になる。
    * ひとまずAWS S3を想定して実装。

### クライアントの画面
* 専用のアップロード画面を用意したい。

### URL
* クライアント側
    * /dm_movie/upload
* サーバー側API
    * /api/upload/dm_movie
    * public API権限
* あたりで。

## Requirements

### Requirement 1: 動画ファイルアップロード機能の提供
**Objective:** As a user, I want to upload large video files to the system, so that I can store and manage video content.

#### Acceptance Criteria
1. WHEN a user selects a video file THEN the system SHALL initiate the upload process using the TUS protocol
2. IF the upload is interrupted THEN the system SHALL support resuming the upload from the last successful chunk
3. WHILE uploading a file THE system SHALL display upload progress to the user
4. WHEN a chunk upload fails THEN the system SHALL automatically retry with exponential backoff delays (0ms, 1000ms, 3000ms, 5000ms)
5. WHERE the chunk size is configured THEN the system SHALL use 5MB chunks for file transfer
6. WHEN the upload completes successfully THEN the system SHALL notify the user and provide the uploaded file information

### Requirement 2: フロントエンドアップロードUIの実装
**Objective:** As a user, I want to use a dedicated upload page with a user-friendly interface, so that I can easily upload video files.

#### Acceptance Criteria
1. WHEN a user navigates to `/dm_movie/upload` THEN the system SHALL display the upload interface
2. IF uppy or tus-js-client is available THEN the system SHALL use it for file upload functionality
3. WHEN the upload interface is displayed THEN the system SHALL show a file selection area and upload progress
4. WHILE a file is being uploaded THE system SHALL display real-time progress percentage and transfer speed
5. WHEN an upload error occurs THEN the system SHALL display an error message to the user
6. WHERE the TUS endpoint is configured THEN the system SHALL connect to `http://localhost:8080/api/upload/dm_movie` (or configured endpoint)

### Requirement 3: バックエンドTUSプロトコルサポート
**Objective:** As a system, I want to support the TUS protocol for resumable uploads, so that large files can be uploaded reliably.

#### Acceptance Criteria
1. WHEN a TUS upload request is received THEN the system SHALL handle it using the tusd library (`github.com/tus/tusd/v2/pkg/handler`)
2. IF the request is a TUS OPTIONS request THEN the system SHALL return TUS protocol capabilities
3. WHEN a TUS POST request is received THEN the system SHALL create an upload resource and return the upload URL
4. WHEN a TUS PATCH request is received THEN the system SHALL append the chunk to the existing upload
5. IF the upload is complete THEN the system SHALL trigger the completion notification hook
6. WHERE the BasePath is configured THEN the system SHALL use `/api/upload/dm_movie` as the TUS endpoint base path
7. WHEN a TUS HEAD request is received THEN the system SHALL return the current upload offset and metadata

### Requirement 4: ファイル保存先の環境別対応
**Objective:** As a system, I want to store uploaded files in different locations based on the environment, so that development and production requirements are met.

#### Acceptance Criteria
1. WHEN the environment is development THEN the system SHALL store files in the local filesystem at `./uploads`
2. WHEN the environment is production or staging THEN the system SHALL store files in AWS S3
3. IF the upload directory does not exist THEN the system SHALL create it automatically
4. WHEN files are stored locally THEN the system SHALL add the upload directory to `.gitignore`
5. WHERE AWS S3 is configured THEN the system SHALL use appropriate S3 bucket and credentials
6. WHEN a file is saved THEN the system SHALL preserve the original filename and metadata

### Requirement 5: アップロード完了時のHook処理
**Objective:** As a system, I want to execute custom logic when an upload completes, so that post-upload processing can be implemented in the future.

#### Acceptance Criteria
1. WHEN an upload completes successfully THEN the system SHALL trigger the PostFinish hook
2. IF PostFinish hook is configured THEN the system SHALL execute the hook handler function
3. WHEN the hook is executed THEN the system SHALL provide upload metadata (file ID, path, size, etc.) in the hook event
4. WHERE PostFinish hook is implemented THEN the system SHALL allow future extensions for upload completion status display
5. IF the hook handler is empty THEN the system SHALL still complete the upload successfully (hook execution does not block upload completion)

### Requirement 6: 認証・認可の実装
**Objective:** As a system, I want to secure the upload endpoint with authentication, so that only authorized users can upload files.

#### Acceptance Criteria
1. WHEN an upload request is received THEN the system SHALL verify the request has valid authentication
2. IF the request has a Public API Key JWT THEN the system SHALL allow the upload
3. IF the request has an Auth0 JWT THEN the system SHALL allow the upload
4. WHEN authentication fails THEN the system SHALL return HTTP 401 Unauthorized
5. WHERE the endpoint is configured as public API THEN the system SHALL accept Public API Key JWT or Auth0 JWT

### Requirement 7: ファイルサイズ制限と拡張子制限の実装
**Objective:** As a system, I want to enforce file size limits and file extension restrictions, so that only appropriate files within size constraints can be uploaded.

#### Acceptance Criteria
1. WHEN a file upload is initiated THEN the system SHALL check the file size against the configured maximum file size limit
2. IF the file size exceeds the maximum limit THEN the system SHALL reject the upload and return an appropriate error response
3. WHEN a file upload is initiated THEN the system SHALL check the file extension against the allowed extensions list
4. IF the file extension is not in the allowed list THEN the system SHALL reject the upload and return an appropriate error response
5. WHERE file size limit is configured THEN the system SHALL enforce the limit before accepting the upload
6. WHERE allowed file extensions are configured THEN the system SHALL only accept files with those extensions

### Requirement 8: エラーハンドリングとログ記録
**Objective:** As a system, I want to handle errors gracefully and log important events, so that issues can be diagnosed and resolved.

#### Acceptance Criteria
1. WHEN an upload error occurs THEN the system SHALL log the error with appropriate context
2. IF a chunk upload fails THEN the system SHALL return an appropriate HTTP error status
3. WHEN storage operations fail THEN the system SHALL return HTTP 500 Internal Server Error
4. IF invalid file metadata is provided THEN the system SHALL return HTTP 400 Bad Request
5. WHEN TUS protocol violations occur THEN the system SHALL return appropriate TUS error responses

### Requirement 9: ファイル検証の実施
**Objective:** As a system, I want to validate uploaded files for security and integrity, so that malicious or corrupted files are not stored.

#### Acceptance Criteria
1. WHEN a file upload completes in production environment THEN the system SHALL validate the file using AWS functionality on S3
2. IF file validation is required THEN the system SHALL use AWS services (e.g., AWS Lambda, AWS Rekognition, or other AWS validation services) to perform validation after upload completion
3. WHEN file validation fails THEN the system SHALL handle the failure appropriately (e.g., delete the file, notify administrators)
4. WHERE the environment is development THEN the system SHALL skip file validation (code implementation for validation is not required)
5. IF file validation is needed THEN the system SHALL NOT implement validation logic in application code (AWS functionality shall be used instead)
