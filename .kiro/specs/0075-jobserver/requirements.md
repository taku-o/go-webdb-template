# ジョブキューサーバー分離要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #153
- **Issueタイトル**: ジョブキューサーバーを別に起動できるようにする
- **Feature名**: 0075-jobserver
- **作成日**: 2026-01-17

### 1.2 目的
現在、このserverアプリではAsynqとRedisでジョブキューを処理している。ジョブの登録はserverアプリで行うが、ジョブの消化処理はAPIサーバー内で実行されている。

本要件では、ジョブの消化処理のみを担当する独立したサーバー（`server/cmd/jobqueue`）を作成し、APIサーバーからジョブの消化処理を分離することを目的とする。これにより、ジョブ処理サーバーを独立してスケールアウトできるようになり、APIサーバーとジョブ処理サーバーの責務が明確に分離される。

### 1.3 スコープ
- JobQueueサーバー（`server/cmd/jobqueue`）の作成
- ジョブ処理フローの実装（processor → usecase → service）
- 既存のジョブタイプ定数（`JobTypeDelayPrint`）の使用
- 現時点のジョブ処理実装（標準出力への文字列出力）
- APIサーバーからのジョブ消化処理の削除
- Redis接続の堅牢性の実装

**本実装の範囲外**:
- 新しいジョブタイプの追加（現時点では標準出力への文字列出力のみ）
- ジョブ処理の最適化
- ジョブ処理の監視機能
- ジョブ処理のメトリクス収集
- ジョブ処理の優先度制御
- ジョブ処理の並列度制御（既存の設定を維持）

## 2. 背景・現状分析

### 2.1 現在の実装
- **ジョブキュー構成**:
  - Redis + Asynqライブラリを使用したジョブキューシステム
  - ジョブの登録: APIサーバー（`server/cmd/server/main.go`）で実装
  - ジョブの消化処理: APIサーバー内でAsynqサーバーを起動して処理
- **ジョブ処理の実装**:
  - `server/internal/service/jobqueue/server.go`: Asynqサーバーの実装
  - `server/internal/service/jobqueue/processor.go`: ジョブ処理の実装（`ProcessDelayPrintJob`）
  - `server/internal/service/jobqueue/constants.go`: ジョブタイプ定数の定義
- **APIサーバーの実装**:
  - `server/cmd/server/main.go`でAsynqサーバーを起動（102-120行目）
  - ジョブ登録用のクライアント（`jobqueue.Client`）を使用
  - ジョブ登録API: `POST /api/dm-jobqueue/register`
- **設定管理**:
  - 環境別設定ファイル: `config/develop.yaml`等
  - Redis接続設定: `config/{env}/cacheserver.yaml`

### 2.2 課題点
1. **責務の混在**: APIサーバーがジョブの登録と消化処理の両方を担当している
2. **スケーラビリティの制約**: ジョブ処理のスケールアウトがAPIサーバーと連動している
3. **リソース競合**: APIサーバーのリソースとジョブ処理のリソースが競合する可能性がある
4. **運用の柔軟性の不足**: ジョブ処理サーバーを独立してスケールアウトできない

### 2.3 本実装による改善点
1. **責務の分離**: APIサーバーとJobQueueサーバーの責務が明確に分離される
2. **スケーラビリティの向上**: ジョブ処理サーバーを独立してスケールアウト可能
3. **運用の柔軟性**: ジョブ処理サーバーを独立して起動・停止・再起動可能
4. **保守性の向上**: レイヤードアーキテクチャに従った実装により、保守性が向上

## 3. 機能要件

### 3.1 JobQueueサーバーの作成

#### 3.1.1 エントリーポイントの作成
- **ファイル**: `server/cmd/jobqueue/main.go`（新規作成）
- **起動コマンド**: `cd server && APP_ENV=develop go run ./cmd/jobqueue/main.go`
- **機能**:
  - 設定ファイルの読み込み（既存の`config.Load()`を使用）
  - Asynqサーバーの初期化と起動
  - ジョブ処理の登録
  - Graceful shutdownの実装
  - Redis接続エラーハンドリング

#### 3.1.2 設定ファイルの使用
- **設定ファイル**: 既存の`config/develop.yaml`等を使用
- **環境変数**: `APP_ENV`環境変数で環境を切り替え（既存システムと同様）
- **Redis設定**: `config/{env}/cacheserver.yaml`からジョブキュー用Redis接続設定を取得

#### 3.1.3 Redis接続の堅牢性
- **起動時の動作**: Redisが起動していない場合でも、JobQueueサーバーは停止せず、起動を継続する
- **接続エラーハンドリング**: Redis接続エラーをログに記録し、処理を一時停止する
- **自動復旧**: Redisが再開した場合、処理を自動的に再開する
- **実装**: 既存のAPIサーバーと同様のRedis接続エラーハンドリングを実装

### 3.2 ジョブ処理フローの実装

#### 3.2.1 レイヤードアーキテクチャの適用
ジョブ処理は以下の順序で実行される：

1. **ジョブハンドラーの登録**
   - `server/internal/service/jobqueue/server.go`: ジョブハンドラーの登録
   - AsynqサーバーがRedisからジョブを取得
   - ジョブタイプを特定（タスクの分類用のキーを使用）

2. **入出力制御とusecase層の呼び出し**
   - `server/internal/service/jobqueue/processor.go`: 入出力制御とusecase層の呼び出し
   - ジョブのペイロードの解析とバリデーション

3. **ビジネスロジックの実行**
   - `server/internal/usecase/jobqueue/delay_print.go`: サービス層を呼び出して処理を実現する。ビジネスロジック。

4. **ビジネスユーティリティロジックの実行**
   - `server/internal/service/delay_print_service.go`: ビジネスユーティリティロジック
   - 実際の処理（標準出力への文字列出力など）を実行

#### 3.2.2 既存アーキテクチャパターンの遵守
- 既存のレイヤードアーキテクチャパターンに従う
- 既存のエラーハンドリングパターンに従う
- 既存の命名規則に従う

### 3.3 既存ジョブタイプ定数の使用

#### 3.3.1 既存定数の使用
- **既存定数**: `server/internal/service/jobqueue/constants.go`に定義されている`JobTypeDelayPrint = "demo:delay_print"`を使用
- **変更不要**: 既存の定数定義を変更する必要はない
- **使用箇所**: JobQueueサーバーでジョブ処理を登録する際に、既存の`JobTypeDelayPrint`定数を使用

#### 3.3.2 定数の使用
- ジョブ登録時に既存のジョブタイプ定数（`JobTypeDelayPrint`）を使用
- ジョブ処理時に既存のジョブタイプ定数を使用して処理を特定
- タイプミスを防ぐため、定数を使用

### 3.4 現時点のジョブ処理実装

#### 3.4.1 標準出力への文字列出力処理
- **機能**: ジョブが処理された時、標準出力に文字列を出力する
- **実装ファイル**:
  - `server/internal/usecase/jobqueue/delay_print.go`: ジョブ処理のusecase層（新規作成）
  - `server/internal/service/delay_print_service.go`: ジョブ処理のservice層（新規作成）
- **既存実装との関係**: 既存の`ProcessDelayPrintJob`と同等の機能を提供
- **レイヤードアーキテクチャ**: usecase → service の順で処理を実行

#### 3.4.2 ジョブ処理の登録
- **登録場所**: `server/cmd/jobqueue/main.go`でジョブ処理を登録
- **登録方法**: Asynqの`ServeMux`を使用してジョブ処理を登録
- **ジョブタイプ**: `JobTypeDelayPrint`を使用

### 3.5 APIサーバーからのジョブ消化処理の削除

#### 3.5.1 Asynqサーバーの起動処理の削除
- **ファイル**: `server/cmd/server/main.go`
- **削除対象**: Asynqサーバーの初期化と起動処理（102-120行目付近）
- **削除内容**:
  - `jobqueue.NewServer()`の呼び出し
  - `jobQueueServer.Start()`の呼び出し
  - `jobQueueServer.Shutdown()`の呼び出し

#### 3.5.2 ジョブ登録機能の維持
- **維持対象**: ジョブ登録用のクライアント（`jobqueue.Client`）
- **維持対象**: ジョブ登録API（`POST /api/dm-jobqueue/register`）
- **維持対象**: ジョブ登録用のUsecase（`DmJobqueueUsecase`）

## 4. 非機能要件

### 4.1 起動コマンド
- **起動コマンド**: `cd server && APP_ENV=develop go run ./cmd/jobqueue/main.go`
- **環境変数**: `APP_ENV`環境変数で環境を切り替え（develop/staging/production）
- **デフォルト環境**: 環境変数が未設定の場合は`develop`環境

### 4.2 Redis接続の堅牢性
- **起動時の動作**: Redisが起動していない場合でも、JobQueueサーバーは停止せず、起動を継続する
- **接続エラーハンドリング**: Redis接続エラーをログに記録し、処理を一時停止する
- **自動復旧**: Redisが再開した場合、処理を自動的に再開する
- **実装**: 既存のAPIサーバーと同様のRedis接続エラーハンドリングを実装

### 4.3 設定ファイルの共有
- **設定ファイル**: 既存の`config/develop.yaml`等を使用
- **環境別設定**: `APP_ENV`環境変数で環境を切り替え
- **Redis設定**: `config/{env}/cacheserver.yaml`からジョブキュー用Redis接続設定を取得

### 4.4 Graceful Shutdown
- **シグナル処理**: SIGINT、SIGTERMシグナルを受信した場合、Graceful shutdownを実行
- **タイムアウト**: 30秒のタイムアウトを設定
- **処理中のジョブ**: 処理中のジョブは完了まで待機

### 4.5 ログ出力
- **標準出力**: ジョブ処理の結果を標準出力に出力
- **エラーログ**: エラーをログに記録
- **起動ログ**: サーバーの起動・停止をログに記録

## 5. 制約事項

### 5.1 既存システムとの関係
- **既存機能の維持**: ジョブ登録API（`POST /api/dm-jobqueue/register`）は既存の動作を維持する
- **設定ファイルの共有**: JobQueueサーバーは既存の設定ファイル（`config/develop.yaml`等）を使用する
- **レイヤードアーキテクチャの遵守**: 既存のアーキテクチャパターンに従う
- **後方互換性**: 既存のジョブタイプ（`JobTypeDelayPrint`）は引き続き動作する

### 5.2 環境別の対応
- **環境切り替え**: `APP_ENV`環境変数で環境を切り替え（既存システムと同様）
- **開発環境**: 本実装は開発環境を優先
- **ステージング・本番環境**: 環境別設定ファイルを使用

### 5.3 技術スタック
- **Go**: 既存のGoバージョン（1.23.4）を維持
- **Asynq**: 既存のAsynqライブラリ（v0.25.1）を維持
- **Redis**: 既存のRedis接続設定を維持
- **データベース**: 既存のデータベース設定を維持

### 5.4 データベース接続
- **接続管理**: 既存の`GroupManager`を使用
- **シャーディング**: 既存のシャーディング戦略を維持
- **接続プール**: 既存の接続プール設定を維持

## 6. 受け入れ基準

### 6.1 JobQueueサーバーの作成
- [ ] `server/cmd/jobqueue/main.go`が作成されている
- [ ] `cd server && APP_ENV=develop go run ./cmd/jobqueue/main.go`でJobQueueサーバーが起動する
- [ ] JobQueueサーバーがRedisからジョブを取得して処理を開始する
- [ ] Redisが起動していない場合でも、JobQueueサーバーは停止せず、起動を継続する
- [ ] Redisが途中で停止して再開した場合、処理を自動的に再開する
- [ ] JobQueueサーバーが既存のAPIサーバーと同じ設定ファイル（`config/develop.yaml`等）を使用する

### 6.2 ジョブ処理フローの実装
- [ ] ジョブ処理が以下の順序で実行される:
  - `server/internal/service/jobqueue/server.go`: ジョブハンドラーの登録
  - `server/internal/service/jobqueue/processor.go`: 入出力制御とusecase層の呼び出し
  - `server/internal/usecase/jobqueue/delay_print.go`: サービス層を呼び出して処理を実現する。ビジネスロジック。
  - `server/internal/service/delay_print_service.go`: ビジネスユーティリティロジック
- [ ] ジョブ処理が既存のレイヤードアーキテクチャパターンに従う
- [ ] ジョブ処理がエラーハンドリングを適切に実装する

### 6.3 既存ジョブタイプ定数の使用
- [ ] 既存のジョブタイプ定数（`JobTypeDelayPrint`）が`server/internal/service/jobqueue/constants.go`に定義されていることを確認
- [ ] JobQueueサーバーで既存の`JobTypeDelayPrint`定数を使用してジョブ処理を登録する
- [ ] ジョブ登録時に既存の`JobTypeDelayPrint`定数を使用する

### 6.4 現時点のジョブ処理実装
- [ ] ジョブが処理された時、標準出力に文字列を出力する
- [ ] `server/internal/usecase/jobqueue/delay_print.go`が作成されている
- [ ] `server/internal/service/delay_print_service.go`が作成されている
- [ ] ジョブ処理が既存の`ProcessDelayPrintJob`と同等の機能を提供する
- [ ] ジョブ処理がレイヤードアーキテクチャに従って実装される

### 6.5 APIサーバーからのジョブ消化処理の削除
- [ ] APIサーバーが起動した時、ジョブの消化処理を開始しない
- [ ] APIサーバーがジョブの登録機能のみを提供する
- [ ] `server/cmd/server/main.go`からAsynqサーバーの起動処理が削除されている
- [ ] APIサーバーがジョブ登録用のクライアント（`jobqueue.Client`）を引き続き使用する
- [ ] ジョブ登録API（`POST /api/dm-jobqueue/register`）が既存の動作を維持する

### 6.6 Redis接続の堅牢性
- [ ] Redisが起動していない状態でJobQueueサーバーを起動した時、エラーで停止せず、起動を継続する
- [ ] Redisが途中で停止した時、エラーをログに記録し、処理を一時停止する
- [ ] Redisが再開した時、自動的に処理を再開する
- [ ] JobQueueサーバーが既存のAPIサーバーと同様のRedis接続エラーハンドリングを実装する

## 7. 影響範囲

### 7.1 新規追加が必要なディレクトリ・ファイル

#### ディレクトリ
- `server/cmd/jobqueue/`: JobQueueサーバーのエントリーポイント
- `server/internal/usecase/jobqueue/`: ジョブ処理のusecase層

#### ファイル
- `server/cmd/jobqueue/main.go`: JobQueueサーバーのエントリーポイント（新規作成）
- `server/internal/usecase/jobqueue/delay_print.go`: ジョブ処理のusecase層（新規作成）
- `server/internal/service/delay_print_service.go`: ジョブ処理のservice層（新規作成）

### 7.2 変更が必要なファイル

#### サーバー実装
- `server/cmd/server/main.go`: Asynqサーバーの起動処理を削除

### 7.3 既存ファイルの扱い
- `server/internal/service/jobqueue/server.go`: 既存の実装を維持（JobQueueサーバーで使用）
- `server/internal/service/jobqueue/processor.go`: 既存の実装を維持（参考として使用）
- `server/internal/service/jobqueue/client.go`: 既存の実装を維持（APIサーバーで使用）
- `server/internal/usecase/api/dm_jobqueue_usecase.go`: 既存の実装を維持（APIサーバーで使用）

## 8. 実装上の注意事項

### 8.1 JobQueueサーバーの実装
- 既存の`server/cmd/server/main.go`を参考に実装
- Asynqサーバーの初期化と起動処理を実装
- Graceful shutdownを実装
- Redis接続エラーハンドリングを実装
- 既存の設定ファイル読み込み処理を使用

### 8.2 ジョブ処理フローの実装
- 既存のレイヤードアーキテクチャパターンに従う
- processor → usecase → service の順で処理を実行
  - `server/internal/service/jobqueue/processor.go`: 入出力制御とusecase層の呼び出し
  - `server/internal/usecase/jobqueue/delay_print.go`: サービス層を呼び出して処理を実現する。ビジネスロジック。
  - `server/internal/service/delay_print_service.go`: ビジネスユーティリティロジック
- 既存のエラーハンドリングパターンに従う
- 既存の命名規則に従う

### 8.3 既存ジョブタイプ定数の使用
- 既存の`JobTypeDelayPrint`定数（`server/internal/service/jobqueue/constants.go`に定義）を使用
- 既存の定数定義を変更する必要はない
- JobQueueサーバーでジョブ処理を登録する際に、既存の`JobTypeDelayPrint`定数を使用

### 8.4 現時点のジョブ処理実装
- `server/internal/usecase/jobqueue/delay_print.go`を新規作成
- `server/internal/service/delay_print_service.go`を新規作成
- 既存の`ProcessDelayPrintJob`と同等の機能を提供
- レイヤードアーキテクチャに従って実装

### 8.5 APIサーバーからのジョブ消化処理の削除
- `server/cmd/server/main.go`からAsynqサーバーの起動処理を削除
- ジョブ登録用のクライアント（`jobqueue.Client`）は維持
- ジョブ登録API（`POST /api/dm-jobqueue/register`）は維持
- ジョブ登録用のUsecase（`DmJobqueueUsecase`）は維持

### 8.6 Redis接続の堅牢性
- 既存のAPIサーバーと同様のRedis接続エラーハンドリングを実装
- Redis接続エラーをログに記録
- Redis接続エラー時もサーバーを停止しない
- Redis接続が復旧した場合、処理を自動的に再開

### 8.7 テストの実装
- JobQueueサーバーの起動テストを実装
- ジョブ処理のユニットテストを実装
- ジョブ処理の統合テストを実装
- Redis接続エラーのテストを実装

## 9. 参考情報

### 9.1 関連Issue
- GitHub Issue #153: ジョブキューサーバーを別に起動できるようにする

### 9.2 既存ドキュメント
- `docs/ja/Queue-Job.md`: ジョブキュー機能の利用手順
- `docs/ja/Project-Structure.md`: プロジェクト構造の説明
- `docs/ja/Architecture.md`: システムアーキテクチャの説明

### 9.3 既存実装
- `server/cmd/server/main.go`: APIサーバーのエントリーポイント
- `server/internal/service/jobqueue/server.go`: Asynqサーバーの実装
- `server/internal/service/jobqueue/processor.go`: ジョブ処理の実装
- `server/internal/service/jobqueue/constants.go`: ジョブタイプ定数の定義
- `server/internal/service/jobqueue/client.go`: ジョブ登録用のクライアント

### 9.4 技術スタック
- **Asynq**: https://github.com/hibiken/asynq
- **Redis**: https://redis.io/
- **Go**: https://golang.org/
