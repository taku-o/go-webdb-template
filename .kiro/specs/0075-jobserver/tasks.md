# ジョブキューサーバー分離実装タスク一覧

## 概要
ジョブキューサーバー分離の実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: JobQueueサーバーのエントリーポイント作成

#### - [ ] タスク 1.1: `server/cmd/jobqueue/main.go`の作成
**目的**: JobQueueサーバーのエントリーポイントを作成し、設定読み込み、サーバー起動、Graceful shutdownを実装

**作業内容**:
- `server/cmd/jobqueue/`ディレクトリを作成
- `server/cmd/jobqueue/main.go`を作成
- 既存の`server/cmd/server/main.go`を参考に実装
- 設定ファイルの読み込み（`config.Load()`）
- Asynqサーバーの初期化（`jobqueue.NewServer()`）
- Redis接続エラーハンドリング（エラーでも起動を継続、警告ログを出力）
- ジョブ処理サーバーの起動（バックグラウンドgoroutine）
- Graceful shutdownの実装（SIGINT、SIGTERMシグナル処理）
- シグナル待機とサーバー停止処理（30秒のタイムアウト）

**受け入れ基準**:
- `server/cmd/jobqueue/main.go`が作成されている
- `cd server && APP_ENV=develop go run ./cmd/jobqueue/main.go`でJobQueueサーバーが起動する
- 設定ファイル（`config/develop.yaml`等）が正しく読み込まれる
- Redis接続エラー時でもサーバーが停止せず、起動を継続する
- Graceful shutdownが正常に動作する（SIGINT、SIGTERMで停止）

---

### Phase 2: ジョブ処理フローの実装

#### - [ ] タスク 2.1: `server/internal/service/delay_print_service.go`の作成
**目的**: ビジネスユーティリティロジック層を実装し、標準出力への文字列出力を提供

**作業内容**:
- `server/internal/service/delay_print_service.go`を作成
- `DelayPrintService`構造体を定義
- `NewDelayPrintService`コンストラクタを実装
- `PrintMessage`メソッドを実装:
  - タイムスタンプの付与（`time.Now().Format("2006-01-02 15:04:05")`）
  - 標準出力への文字列出力（`fmt.Printf`）
  - バッファのフラッシュ（`os.Stdout.Sync()`）
- 既存の`ProcessDelayPrintJob`の処理ロジックを移行

**受け入れ基準**:
- `server/internal/service/delay_print_service.go`が作成されている
- `DelayPrintService`構造体が定義されている
- `PrintMessage`メソッドが実装されている
- 標準出力にタイムスタンプ付きで文字列が出力される
- バッファが正しくフラッシュされる

---

#### - [ ] タスク 2.2: `server/internal/usecase/jobqueue/delay_print.go`の作成
**目的**: ビジネスロジック層を実装し、サービス層を呼び出して処理を実現

**作業内容**:
- `server/internal/usecase/jobqueue/`ディレクトリを作成
- `server/internal/usecase/jobqueue/delay_print.go`を作成
- `DelayPrintUsecase`構造体を定義
- `NewDelayPrintUsecase`コンストラクタを実装（サービス層への依存を注入）
- `Execute`メソッドを実装:
  - ペイロードからメッセージを取得
  - デフォルトメッセージの設定（空文字列の場合: "Job executed successfully"）
  - サービス層（`DelayPrintService`）の呼び出し
  - エラーハンドリング

**受け入れ基準**:
- `server/internal/usecase/jobqueue/delay_print.go`が作成されている
- `DelayPrintUsecase`構造体が定義されている
- `Execute`メソッドが実装されている
- デフォルトメッセージが正しく設定される
- サービス層が正しく呼び出される

---

#### - [ ] タスク 2.3: `server/internal/service/jobqueue/processor.go`の修正
**目的**: 入出力制御とusecase層の呼び出しを実装

**作業内容**:
- 既存の`server/internal/service/jobqueue/processor.go`を修正
- `ProcessDelayPrintJob`関数を修正:
  - ペイロードの解析とバリデーション（既存の処理を維持）
  - usecase層（`DelayPrintUsecase`）の呼び出しに変更
  - usecase層への依存を注入（コンストラクタまたは関数パラメータ）
  - エラーハンドリングの維持
- 既存の標準出力への直接出力処理を削除

**受け入れ基準**:
- `ProcessDelayPrintJob`関数がusecase層を呼び出すように修正されている
- ペイロードの解析とバリデーションが維持されている
- usecase層への依存が正しく注入されている
- エラーハンドリングが維持されている

---

#### - [ ] タスク 2.4: `server/internal/service/jobqueue/server.go`の確認
**目的**: 既存のAsynqサーバー実装を確認し、ジョブハンドラーの登録を確認

**作業内容**:
- 既存の`server/internal/service/jobqueue/server.go`を確認
- `JobTypeDelayPrint`定数の使用を確認
- ジョブハンドラーの登録（`mux.HandleFunc(JobTypeDelayPrint, ProcessDelayPrintJob)`）を確認
- 変更が必要な場合は修正（基本的には変更不要）

**受け入れ基準**:
- `server.go`の実装が確認されている
- `JobTypeDelayPrint`定数が正しく使用されている
- ジョブハンドラーが正しく登録されている

---

### Phase 3: APIサーバーからのジョブ消化処理の削除

#### - [ ] タスク 3.1: `server/cmd/server/main.go`の修正
**目的**: APIサーバーからAsynqサーバーの起動処理を削除し、ジョブ登録機能のみを維持

**作業内容**:
- `server/cmd/server/main.go`を修正
- Asynqサーバーの初期化処理を削除（102-120行目付近）:
  - `jobqueue.NewServer()`の呼び出しを削除
  - `jobQueueServer.Start()`の呼び出しを削除
  - `jobQueueServer.Shutdown()`の呼び出しを削除
  - 関連する変数宣言を削除（`jobQueueServer`）
- ジョブ登録用のクライアント（`jobqueue.Client`）は維持
- ジョブ登録用のUsecase（`DmJobqueueUsecase`）は維持
- ジョブ登録API（`POST /api/dm-jobqueue/register`）は維持

**受け入れ基準**:
- `server/cmd/server/main.go`からAsynqサーバーの起動処理が削除されている
- APIサーバーが起動した時、ジョブの消化処理を開始しない
- ジョブ登録用のクライアントが維持されている
- ジョブ登録APIが正常に動作する

---

### Phase 4: テストの実装

#### - [ ] タスク 4.1: `server/internal/service/delay_print_service_test.go`の作成
**目的**: `DelayPrintService`のユニットテストを実装

**作業内容**:
- `server/internal/service/delay_print_service_test.go`を作成
- `DelayPrintService`のユニットテストを実装:
  - `PrintMessage`メソッドのテスト
  - 標準出力への出力確認（`os.Stdout`をキャプチャ）
  - タイムスタンプの確認
  - バッファフラッシュの確認
  - エラーハンドリングのテスト（該当する場合）

**受け入れ基準**:
- `delay_print_service_test.go`が作成されている
- `PrintMessage`メソッドのテストが実装されている
- 標準出力への出力が確認できる
- タイムスタンプが正しく付与される

---

#### - [ ] タスク 4.2: `server/internal/usecase/jobqueue/delay_print_test.go`の作成
**目的**: `DelayPrintUsecase`のユニットテストを実装

**作業内容**:
- `server/internal/usecase/jobqueue/delay_print_test.go`を作成
- `DelayPrintUsecase`のユニットテストを実装:
  - `Execute`メソッドのテスト
  - デフォルトメッセージの設定テスト（空文字列の場合）
  - サービス層の呼び出しテスト（モック使用）
  - エラーハンドリングのテスト

**受け入れ基準**:
- `delay_print_test.go`が作成されている
- `Execute`メソッドのテストが実装されている
- デフォルトメッセージが正しく設定されることを確認
- サービス層が正しく呼び出されることを確認（モック使用）

---

#### - [ ] タスク 4.3: `server/internal/service/jobqueue/processor_test.go`の作成
**目的**: `ProcessDelayPrintJob`のユニットテストを実装

**作業内容**:
- `server/internal/service/jobqueue/processor_test.go`を作成
- `ProcessDelayPrintJob`のユニットテストを実装:
  - ペイロード解析のテスト（正常系、異常系）
  - usecase層の呼び出しテスト（モック使用）
  - エラーハンドリングのテスト

**受け入れ基準**:
- `processor_test.go`が作成されている
- `ProcessDelayPrintJob`のテストが実装されている
- ペイロード解析が正しく動作することを確認
- usecase層が正しく呼び出されることを確認（モック使用）

---

#### - [ ] タスク 4.4: `server/cmd/jobqueue/main_test.go`の作成
**目的**: JobQueueサーバーの統合テストを実装

**作業内容**:
- `server/cmd/jobqueue/main_test.go`を作成
- JobQueueサーバーの統合テストを実装:
  - 設定ファイル読み込みのテスト
  - Redis接続エラーのテスト（エラーでも起動を継続）
  - Graceful shutdownのテスト（シグナル処理）

**受け入れ基準**:
- `main_test.go`が作成されている
- 設定ファイル読み込みのテストが実装されている
- Redis接続エラーのテストが実装されている
- Graceful shutdownのテストが実装されている

---

### Phase 5: 動作確認と検証

#### - [ ] タスク 5.1: JobQueueサーバーの起動確認
**目的**: JobQueueサーバーが正常に起動することを確認

**作業内容**:
- `cd server && APP_ENV=develop go run ./cmd/jobqueue/main.go`で起動
- 起動ログを確認
- Redis接続状態を確認
- エラーがないことを確認

**受け入れ基準**:
- JobQueueサーバーが正常に起動する
- 起動ログにエラーが表示されない
- Redis接続が正常に確立される（Redisが起動している場合）

---

#### - [ ] タスク 5.2: ジョブ処理の動作確認
**目的**: ジョブが正常に処理されることを確認

**作業内容**:
- APIサーバーからジョブを登録（`POST /api/dm-jobqueue/register`）
- JobQueueサーバーでジョブが処理されることを確認
- 標準出力への文字列出力を確認
- タイムスタンプが正しく付与されることを確認

**受け入れ基準**:
- ジョブが正常に登録される
- JobQueueサーバーでジョブが処理される
- 標準出力にタイムスタンプ付きで文字列が出力される

---

#### - [ ] タスク 5.3: APIサーバーの動作確認
**目的**: APIサーバーが正常に動作し、ジョブの消化処理が開始されないことを確認

**作業内容**:
- APIサーバーを起動（`cd server && APP_ENV=develop go run ./cmd/server/main.go`）
- 起動ログを確認（Asynqサーバーの起動ログが表示されないことを確認）
- ジョブ登録API（`POST /api/dm-jobqueue/register`）が正常に動作することを確認
- ジョブが正常に登録されることを確認

**受け入れ基準**:
- APIサーバーが正常に起動する
- 起動ログにAsynqサーバーの起動メッセージが表示されない
- ジョブ登録APIが正常に動作する
- ジョブが正常に登録される

---

#### - [ ] タスク 5.4: Redis接続エラーの動作確認
**目的**: Redis接続エラー時でもJobQueueサーバーが正常に動作することを確認

**作業内容**:
- Redisを停止
- Redisが起動していない状態でJobQueueサーバーを起動
- エラーで停止せず、起動を継続することを確認（警告ログが出力されることを確認）
- Redisを起動
- Redis起動後、処理が自動的に再開することを確認

**受け入れ基準**:
- Redisが起動していない状態でJobQueueサーバーが起動を継続する
- 警告ログが出力される
- Redis起動後、処理が自動的に再開する

---

#### - [ ] タスク 5.5: 全体テストの実行
**目的**: 既存のテストと新規作成したテストが全て成功することを確認

**作業内容**:
- `APP_ENV=test go test ./...`を実行
- 既存のテストが全て成功することを確認
- 新規作成したテストが成功することを確認
- テストカバレッジを確認（オプション）

**受け入れ基準**:
- 既存のテストが全て成功する
- 新規作成したテストが全て成功する
- テストエラーが0件である

---

### Phase 6: ドキュメントの更新

#### - [ ] タスク 6.1: `docs/ja/Queue-Job.md`の更新
**目的**: ジョブキュー機能の利用手順を更新し、JobQueueサーバーの起動手順を追加

**作業内容**:
- 「環境構築」セクションにJobQueueサーバーの起動手順を追加:
  - `cd server && APP_ENV=develop go run ./cmd/jobqueue/main.go`
- 「ジョブの実行確認」セクションを更新:
  - APIサーバーの標準出力ではなく、JobQueueサーバーの標準出力に出力されることを明記
- 「停止手順」セクションにJobQueueサーバーの停止を追加:
  - Ctrl+C または `kill <PID>` で停止
- 「注意事項」セクションの「迷子のAsynqサーバー設定問題」を更新:
  - APIサーバーではなくJobQueueサーバーについての説明に変更
  - `ps aux | grep "go run ./cmd/jobqueue/main.go"` のコマンドに変更
- 「Redis接続エラー時の動作」セクションを更新:
  - JobQueueサーバーについての説明を追加

**受け入れ基準**:
- JobQueueサーバーの起動手順が記載されている
- ジョブの実行確認の説明が正しく更新されている
- JobQueueサーバーの停止手順が記載されている
- 注意事項が正しく更新されている

---

#### - [ ] タスク 6.2: `docs/ja/Project-Structure.md`の更新
**目的**: プロジェクト構造にJobQueueサーバーと新規作成したファイルを追加

**作業内容**:
- `server/cmd/`セクションに`jobqueue/main.go`を追加
- `server/internal/usecase/`セクションに`jobqueue/`ディレクトリを追加:
  - `delay_print.go`
  - `delay_print_test.go`
- `server/internal/service/`セクションに以下を追加:
  - `delay_print_service.go`
  - `delay_print_service_test.go`

**受け入れ基準**:
- `server/cmd/jobqueue/main.go`が記載されている
- `server/internal/usecase/jobqueue/`ディレクトリが記載されている
- `server/internal/service/delay_print_service.go`が記載されている

---

#### - [ ] タスク 6.3: `docs/ja/Architecture.md`の更新
**目的**: システムアーキテクチャにJobQueueサーバーの説明を追加

**作業内容**:
- 「System Architecture」セクションにJobQueueサーバーの説明を追加
- 「Layer Responsibilities」セクションの「Usecase Layer」に以下を追加:
  - `DelayPrintUsecase`: ジョブ処理のビジネスロジック
- 「Layer Responsibilities」セクションの「Service Layer」に以下を追加:
  - `DelayPrintService`: ジョブ処理のビジネスユーティリティロジック
- ジョブ処理フローの説明を追加（processor → usecase → service）

**受け入れ基準**:
- JobQueueサーバーの説明が追加されている
- `DelayPrintUsecase`が記載されている
- `DelayPrintService`が記載されている
- ジョブ処理フローの説明が追加されている

---

#### - [ ] タスク 6.4: `docs/en/Queue-Job.md`の更新
**目的**: 英語版のジョブキュー機能の利用手順を更新

**作業内容**:
- `docs/ja/Queue-Job.md`と同様の更新を英語版に適用
- JobQueueサーバーの起動手順を追加
- ジョブの実行確認の説明を更新
- JobQueueサーバーの停止手順を追加
- 注意事項を更新

**受け入れ基準**:
- 英語版のドキュメントが日本語版と同様に更新されている
- JobQueueサーバーに関する情報が正しく記載されている

---

#### - [ ] タスク 6.5: `docs/en/Project-Structure.md`の更新
**目的**: 英語版のプロジェクト構造を更新

**作業内容**:
- `docs/ja/Project-Structure.md`と同様の更新を英語版に適用
- `server/cmd/jobqueue/main.go`を追加
- `server/internal/usecase/jobqueue/`ディレクトリを追加
- `server/internal/service/delay_print_service.go`を追加

**受け入れ基準**:
- 英語版のドキュメントが日本語版と同様に更新されている
- 新規作成したファイルが正しく記載されている

---

#### - [ ] タスク 6.6: `docs/en/Architecture.md`の更新
**目的**: 英語版のシステムアーキテクチャを更新

**作業内容**:
- `docs/ja/Architecture.md`と同様の更新を英語版に適用
- JobQueueサーバーの説明を追加
- `DelayPrintUsecase`と`DelayPrintService`を追加
- ジョブ処理フローの説明を追加

**受け入れ基準**:
- 英語版のドキュメントが日本語版と同様に更新されている
- JobQueueサーバーに関する情報が正しく記載されている

---

#### - [ ] タスク 6.7: `README.md`の更新
**目的**: プロジェクトの概要にJobQueueサーバーの起動手順を追加

**作業内容**:
- 「Start Server」セクションにJobQueueサーバーの起動手順を追加:
  - `cd server && APP_ENV=develop go run ./cmd/jobqueue/main.go`
- 開発環境サーバー構成の説明を更新:
  - APIサーバー、クライアント、Adminサーバーに加えて、JobQueueサーバーを追加
- 必要に応じて、ジョブキュー機能の説明を更新

**受け入れ基準**:
- JobQueueサーバーの起動手順が記載されている
- 開発環境サーバー構成が正しく更新されている

---

#### - [ ] タスク 6.8: `README.ja.md`の更新
**目的**: 日本語版のプロジェクト概要にJobQueueサーバーの起動手順を追加

**作業内容**:
- 「3. サーバー起動」セクションにJobQueueサーバーの起動手順を追加:
  - `cd server && APP_ENV=develop go run ./cmd/jobqueue/main.go`
- 「7. Redisの起動（ジョブキュー機能用）」セクションの後に、JobQueueサーバーの起動セクションを追加:
  - セクション番号を調整（「8. Mailpitの起動」以降を1つずつ繰り下げ）
  - 「8. JobQueueサーバーの起動（ジョブキュー機能用）」セクションを追加
  - JobQueueサーバーの起動手順を記載
  - ジョブ処理の説明を追加
- 必要に応じて、ジョブキュー機能の説明を更新

**受け入れ基準**:
- JobQueueサーバーの起動手順が記載されている
- セクション番号が正しく調整されている
- ジョブキュー機能の説明が正しく更新されている
