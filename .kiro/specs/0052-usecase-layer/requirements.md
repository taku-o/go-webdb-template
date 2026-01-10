# APIサーバーにusecase層を導入する要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #107
- **Issueタイトル**: APIサーバーにusecase層を導入する
- **Feature名**: 0052-usecase-layer
- **作成日**: 2026-01-10

### 1.2 目的
APIサーバーのレイヤー構造を改善し、ビジネスロジックとドメインロジックを明確に分離する。これにより、コードの可読性、保守性、テスタビリティを向上させる。

### 1.3 スコープ
- usecase層の導入とディレクトリ構造の作成
- 各レイヤーの役割の明確化
- 既存のすべてのAPIエンドポイントのhandlerからserviceへの直接呼び出しを、handler -> usecase -> service の構造に変更
- 以下のすべてのAPIエンドポイントにusecase層を導入:
  - today API (`today_handler.go`)
  - dm_user API (`dm_user_handler.go`)
  - dm_post API (`dm_post_handler.go`)
  - email API (`email_handler.go`)
  - dm_jobqueue API (`dm_jobqueue_handler.go`)
  - upload API (`upload_handler.go`)

**本実装の範囲外**:
- usecase層の詳細な設計パターン（薄いusecase層として実装）
- 複雑なビジネスロジックの実装（現時点では複雑な処理を行っていないため）

## 2. 背景・現状分析

### 2.1 現在の状況
- **レイヤー構造**: `api(controller) -> service(ビジネスロジック) -> repository -> model`
- **handler層**: `server/internal/api/handler/` に各種ハンドラーが存在
  - `today_handler.go`: 今日の日付を返すAPI（認証チェックのみ）
  - `dm_user_handler.go`: ユーザーCRUD API（service層を直接呼び出し）
  - `dm_post_handler.go`: 投稿CRUD API（service層を直接呼び出し）
  - その他のハンドラー
- **service層**: `server/internal/service/` に各種サービスが存在
  - `dm_user_service.go`: ユーザーのビジネスロジックとドメインロジックが混在
  - `dm_post_service.go`: 投稿のビジネスロジックとドメインロジックが混在
  - その他のサービス

### 2.2 課題点
1. **レイヤー責任の混在**: service層にビジネスロジックとドメインロジックが混在している
2. **拡張性の制約**: 複数のドメインサービスを組み合わせたビジネスロジックを実装しにくい
3. **テストの困難さ**: ビジネスロジックとドメインロジックが分離されていないため、テストが複雑になる
4. **コードの可読性**: ビジネスロジックとドメインロジックの境界が不明確

### 2.3 本実装による改善点
1. **レイヤー責任の明確化**: usecase層でビジネスロジック、service層でドメインロジックを担当
2. **拡張性の向上**: 複数のserviceを組み合わせたビジネスロジックをusecase層で実装可能
3. **テストの容易さ**: 各レイヤーの責任が明確になることで、テストが容易になる
4. **コードの可読性**: レイヤー構造が明確になり、コードの理解が容易になる

## 3. 機能要件

### 3.1 usecase層の導入

#### 3.1.1 ディレクトリ構造
- **ディレクトリパス**: `server/internal/usecase/`
- **命名規則**: `{機能名}_usecase.go`（例: `today_usecase.go`）

#### 3.1.2 レイヤー構造の変更
- **変更前**: `api(controller) -> service(ビジネスロジック) -> repository -> model`
- **変更後**: `api(controller) -> usecase(ビジネスロジック) -> service(ドメインロジック) -> repository -> model`

### 3.2 各レイヤーの役割定義

#### 3.2.1 api層（handler）
- **役割**: 入出力の制御、バリデーション
- **責務**:
  - HTTPリクエスト/レスポンスの処理
  - 認証・認可チェック
  - 入力値のバリデーション（形式チェックなど）
  - エラーハンドリングとHTTPステータスコードの設定
- **制約**: 極力処理を行わない（ビジネスロジックはusecase層に委譲）

#### 3.2.2 usecase層
- **役割**: ビジネスロジック
- **責務**:
  - 複数のserviceを組み合わせたビジネスロジックの実装
  - トランザクション管理（必要に応じて）
  - ビジネスルールの適用
- **制約**:
  - usecaseから別のusecaseは呼び出さない
  - usecaseは複数、もしくは一つのserviceを呼び出して処理を行う
  - repository層を直接呼び出さない（service層を経由）

#### 3.2.3 service層
- **役割**: ドメインロジック
- **責務**:
  - 特定の分類（ドメイン）の処理を行う
  - ドメイン固有のバリデーション
  - ドメイン固有のビジネスルール
- **制約**: 単一のドメインに特化した処理を行う

### 3.3 全APIエンドポイントのusecase層実装

#### 3.3.1 実装対象
以下のすべてのAPIエンドポイントにusecase層を導入:
- **today API**: `server/internal/api/handler/today_handler.go`
  - 現在の実装: handler内で直接処理（認証チェックと日付取得）
  - 変更内容: usecase層とservice層を追加
- **dm_user API**: `server/internal/api/handler/dm_user_handler.go`
  - 現在の実装: handlerからservice層を直接呼び出し
  - 変更内容: usecase層を追加し、handler -> usecase -> service の構造に変更
- **dm_post API**: `server/internal/api/handler/dm_post_handler.go`
  - 現在の実装: handlerからservice層を直接呼び出し
  - 変更内容: usecase層を追加し、handler -> usecase -> service の構造に変更
- **email API**: `server/internal/api/handler/email_handler.go`
  - 現在の実装: handlerからservice層を直接呼び出し
  - 変更内容: usecase層を追加し、handler -> usecase -> service の構造に変更
- **dm_jobqueue API**: `server/internal/api/handler/dm_jobqueue_handler.go`
  - 現在の実装: handlerからservice層を直接呼び出し
  - 変更内容: usecase層を追加し、handler -> usecase -> service の構造に変更
- **upload API**: `server/internal/api/handler/upload_handler.go`
  - 現在の実装: handler内で直接処理
  - 変更内容: usecase層を追加し、handler -> usecase -> service の構造に変更

#### 3.3.2 usecase層の実装
各APIエンドポイントに対応するusecase層を作成:
- **ファイル**: `server/internal/usecase/{機能名}_usecase.go`
  - `today_usecase.go`: 今日の日付を取得するビジネスロジック
  - `dm_user_usecase.go`: ユーザー関連のビジネスロジック
  - `dm_post_usecase.go`: 投稿関連のビジネスロジック
  - `email_usecase.go`: メール送信関連のビジネスロジック
  - `dm_jobqueue_usecase.go`: ジョブキュー関連のビジネスロジック
  - `upload_usecase.go`: アップロード関連のビジネスロジック
- **実装内容**:
  - service層を呼び出して処理を実行
  - 必要に応じて複数のserviceを組み合わせたビジネスロジックを実装
  - 現時点では複雑な処理を行っていないため、薄いusecase層として実装

#### 3.3.3 service層の実装
既存のservice層をドメインロジックのみを担当するように整理:
- **既存のservice**: ドメインロジックのみを担当するように整理
- **新規作成が必要なservice**: 
  - `date_service.go`: 日付関連のドメインロジック（today API用）
- **実装内容**:
  - ドメイン固有の処理を実装
  - ビジネスロジックはusecase層に移譲

#### 3.3.4 handler層の修正
すべてのhandler層をusecase層を呼び出すように修正:
- **変更内容**:
  - usecase層を呼び出すように変更
  - 認証チェックはhandler層で継続して行う
  - 入力値の形式チェックはhandler層で行う
  - ビジネスロジックはusecase層に委譲

### 3.4 依存関係の注入

#### 3.4.1 依存関係の構造
- handlerはusecaseに依存
- usecaseはserviceに依存
- serviceはrepositoryに依存

#### 3.4.2 初期化処理
- 各レイヤーの初期化処理を修正
- 依存関係を適切に注入

## 4. 非機能要件

### 4.1 パフォーマンス
- **レイヤー追加によるオーバーヘッド**: 最小限に抑える（薄いusecase層として実装）
- **既存のパフォーマンス**: 既存のパフォーマンスに影響を与えない

### 4.2 保守性
- **コードの可読性**: レイヤー構造が明確になり、コードの理解が容易になる
- **拡張性**: 新しいビジネスロジックを追加しやすい構造
- **一貫性**: 既存のコードスタイルと一貫性を保つ

### 4.3 テスタビリティ
- **単体テスト**: 各レイヤーを独立してテスト可能
- **モック**: 各レイヤーをモック化してテスト可能
- **既存テスト**: 既存のテストが正常に動作する

## 5. 制約事項

### 5.1 技術的制約
- **既存のコードスタイル**: 既存のコードスタイルを維持
- **既存のテスト**: 既存のテストが正常に動作することを確認
- **後方互換性**: 既存のAPIの動作に影響を与えない

### 5.2 実装上の制約
- **薄いusecase層**: 現在の実装は複雑なことをしていないため、単純に薄いusecase層を差し込む
- **全APIエンドポイントへの導入**: すべてのAPIエンドポイントにusecase層を導入する
- **既存機能への影響**: 既存の機能に影響を与えない（APIの動作は変わらない）

### 5.3 動作環境
- **ローカル環境**: ローカル環境で正常に動作することを確認
- **既存のAPI**: 既存のAPIが正常に動作することを確認

## 6. 受け入れ基準

### 6.1 usecase層の導入
- [ ] `server/internal/usecase/` ディレクトリが作成されている
- [ ] すべてのAPIエンドポイントに対応するusecase層が作成されている:
  - [ ] `today_usecase.go`
  - [ ] `dm_user_usecase.go`
  - [ ] `dm_post_usecase.go`
  - [ ] `email_usecase.go`
  - [ ] `dm_jobqueue_usecase.go`
  - [ ] `upload_usecase.go`
- [ ] usecase層の構造が適切に定義されている

### 6.2 service層の拡張
- [ ] 既存のservice層がドメインロジックのみを担当するように整理されている
- [ ] 新規作成が必要なservice（`date_service.go`等）が作成されている
- [ ] service層がドメインロジックのみを担当している

### 6.3 全APIエンドポイントの実装
- [ ] すべてのhandler層がusecase層を呼び出すように修正されている:
  - [ ] `today_handler.go`
  - [ ] `dm_user_handler.go`
  - [ ] `dm_post_handler.go`
  - [ ] `email_handler.go`
  - [ ] `dm_jobqueue_handler.go`
  - [ ] `upload_handler.go`
- [ ] すべてのusecase層がservice層を呼び出している
- [ ] 既存のすべてのAPIの動作が変わらない（同じレスポンスを返す）

### 6.4 依存関係の注入
- [ ] 各レイヤーの依存関係が適切に注入されている
- [ ] 初期化処理が適切に実装されている

### 6.5 テスト
- [ ] usecase層のテストが作成されている
- [ ] service層のテストが作成されている（または既存のテストが正常に動作する）
- [ ] handler層のテストが正常に動作する
- [ ] 既存のテストが全て失敗しないことを確認

### 6.6 動作確認
- [ ] すべてのAPIエンドポイントが正常に動作することを確認:
  - [ ] today API
  - [ ] dm_user API
  - [ ] dm_post API
  - [ ] email API
  - [ ] dm_jobqueue API
  - [ ] upload API
- [ ] 既存のAPIの動作が変わらないことを確認（影響がないことを確認）
- [ ] レイヤー構造が適切に実装されていることを確認

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### 新規作成するファイル
- `server/internal/usecase/today_usecase.go`: today API用のusecase層
- `server/internal/usecase/dm_user_usecase.go`: dm_user API用のusecase層
- `server/internal/usecase/dm_post_usecase.go`: dm_post API用のusecase層
- `server/internal/usecase/email_usecase.go`: email API用のusecase層
- `server/internal/usecase/dm_jobqueue_usecase.go`: dm_jobqueue API用のusecase層
- `server/internal/usecase/upload_usecase.go`: upload API用のusecase層
- `server/internal/service/date_service.go`: 日付関連のservice層（新規作成、today API用）

#### 修正が必要なファイル
- `server/internal/api/handler/today_handler.go`: usecase層を呼び出すように修正
- `server/internal/api/handler/dm_user_handler.go`: usecase層を呼び出すように修正
- `server/internal/api/handler/dm_post_handler.go`: usecase層を呼び出すように修正
- `server/internal/api/handler/email_handler.go`: usecase層を呼び出すように修正
- `server/internal/api/handler/dm_jobqueue_handler.go`: usecase層を呼び出すように修正
- `server/internal/api/handler/upload_handler.go`: usecase層を呼び出すように修正
- 初期化処理（router.go等）: usecase層の初期化を追加

#### 確認が必要なファイル
- 既存のserviceファイル: ドメインロジックのみを担当するように整理されていることを確認
- 既存のテストファイル: 正常に動作することを確認

### 7.2 既存機能への影響
- **すべてのAPIエンドポイント**: usecase層を経由するように変更するが、APIの動作は変わらない
- **既存のテスト**: 既存のテストが正常に動作することを確認

### 7.3 将来の拡張への影響
- **パターンの確立**: すべてのAPIエンドポイントでusecase層の実装パターンが確立される
- **レイヤー構造の明確化**: 今後の開発でレイヤー構造が明確になる
- **一貫性**: すべてのAPIエンドポイントで一貫したレイヤー構造が維持される

## 8. 実装上の注意事項

### 8.1 usecase層の実装
- **薄い実装**: 現在の実装は複雑なことをしていないため、単純に薄いusecase層を差し込む
- **ビジネスロジック**: 将来的にビジネスロジックが複雑になった場合に備えて、usecase層で実装する
- **service層の呼び出し**: usecase層から複数のserviceを呼び出すことが可能

### 8.2 service層の実装
- **ドメインロジック**: service層はドメインロジックのみを担当
- **既存のservice**: 既存のservice層をドメインロジックのみを担当するように整理（必要に応じて）
- **新規作成が必要なservice**: today API用の`date_service.go`等、必要に応じて新規作成

### 8.3 handler層の実装
- **認証チェック**: 認証チェックはhandler層で継続して行う
- **バリデーション**: 入力値の形式チェックはhandler層で行う
- **エラーハンドリング**: エラーハンドリングとHTTPステータスコードの設定はhandler層で行う

### 8.4 テストの実装
- **各レイヤーのテスト**: 各レイヤーを独立してテスト可能にする
- **モックの使用**: 各レイヤーをモック化してテストする
- **既存テストの確認**: 既存のテストが正常に動作することを確認

## 9. 参考情報

### 9.1 関連Issue
- GitHub Issue #107: APIサーバーにusecase層を導入する

### 9.2 既存実装の参考
- **handler層**: `server/internal/api/handler/today_handler.go`
- **service層**: `server/internal/service/dm_user_service.go`, `server/internal/service/dm_post_service.go`
- **repository層**: `server/internal/repository/`

### 9.3 技術スタック
- **言語**: Go
- **フレームワーク**: Huma API, Echo
- **アーキテクチャ**: レイヤードアーキテクチャ（handler -> usecase -> service -> repository -> model）

### 9.4 関連ドキュメント
- `server/internal/api/handler/today_handler.go`: 現在のtoday APIの実装
- `server/internal/service/`: 既存のservice層の実装
- `server/internal/repository/`: 既存のrepository層の実装
