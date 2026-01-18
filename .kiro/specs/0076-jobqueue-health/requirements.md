# JobQueueサーバー死活監視エンドポイントの要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Feature名**: 0076-jobqueue-health
- **作成日**: 2026-01-17

### 1.2 目的
APIサーバー、クライアントサーバー、Adminサーバーには`/health`という死活監視用のエンドポイントが用意されている。新しく作成したJobQueueサーバーにも同様の`/health`エンドポイントを実装し、他のサーバーと一貫性を保つ。これにより、JobQueueサーバーの死活監視が可能になり、docker-compose等のヘルスチェック機能を有効化できる。

### 1.3 スコープ
- JobQueueサーバーにHTTPサーバーを追加し、`/health`エンドポイントを実装する
- エンドポイントは認証不要でアクセス可能とする
- 設定ファイルにJobQueueサーバー用のポート設定を追加する
- 他のサーバー（API、Admin）と同様の実装パターンを使用する

**本実装の範囲外**:
- 他のサーバーの`/health`エンドポイントの変更
- その他のエンドポイントの追加や変更
- ヘルスチェックの詳細な診断機能（現時点では単純なOK応答のみ）
- ジョブ処理機能の変更

## 2. 背景・現状分析

### 2.1 現在の状況
- **APIサーバー**: `server/internal/api/router/router.go`に`/health`エンドポイントが実装されている
  - 実装内容: `GET /health`で`200 OK`と`"OK"`という文字列を返す
  - 認証: 不要
  - ポート: 8080
  - フレームワーク: Echo
- **Adminサーバー**: `server/cmd/admin/main.go`に`/health`エンドポイントが実装されている
  - 実装内容: `GET /health`で`200 OK`と`"OK"`という文字列を返す
  - 認証: 不要
  - ポート: 8081
  - ルーター: Gorilla Mux Router
- **JobQueueサーバー**: `/health`エンドポイントが未実装
  - 現在の実装: HTTPサーバーを起動していない（Asynqサーバーのみ）
  - ポート: 未設定（HTTPサーバーが存在しないため）
  - エントリーポイント: `server/cmd/jobqueue/main.go`

### 2.2 課題点
1. **JobQueueサーバーのヘルスチェック不可**: docker-compose等のヘルスチェックが動作していない
2. **一貫性の欠如**: 他のサーバー（API、Admin）には`/health`エンドポイントがあるが、JobQueueサーバーにはない
3. **監視ツールとの連携不可**: 外部監視ツールがJobQueueサーバーの死活監視を行うことができない
4. **運用上の問題**: JobQueueサーバーが正常に動作しているかどうかをHTTPリクエストで確認できない

### 2.3 本実装による改善点
1. **ヘルスチェック機能の有効化**: docker-compose等のヘルスチェックが正常に動作する
2. **一貫性の確保**: すべてのサーバー（API、クライアント、Admin、JobQueue）に`/health`エンドポイントが存在する
3. **監視ツールとの連携**: 外部監視ツールがJobQueueサーバーの死活監視を行えるようになる
4. **運用の改善**: HTTPリクエストでJobQueueサーバーの状態を確認できるようになる

## 3. 機能要件

### 3.1 HTTPサーバーの追加

#### 3.1.1 HTTPサーバーの起動
- **実装場所**: `server/cmd/jobqueue/main.go`
- **実装方法**: 既存のAdminサーバー（`server/cmd/admin/main.go`）を参考に実装
- **機能**:
  - HTTPサーバーを起動する
  - `/health`エンドポイントを登録する
  - Asynqサーバーと並行して動作する（goroutineで起動）
  - Graceful shutdownを実装する

#### 3.1.2 ポート設定の追加
- **設定ファイル**: `config/develop.yaml`等の環境別設定ファイル
- **設定項目**: `jobqueue.port`（例: 8082）
- **設定構造**: 既存の`server.port`、`admin.port`と同様の構造
- **デフォルト値**: 8082（他のサーバーと重複しないポート番号）

### 3.2 `/health`エンドポイントの実装

#### 3.1.1 エンドポイント仕様
- **パス**: `/health`
- **メソッド**: `GET`
- **認証**: 不要（認証ミドルウェアを通過しない）
- **レスポンス**: 
  - ステータスコード: `200 OK`
  - レスポンスボディ: `"OK"`（文字列）
  - Content-Type: `text/plain`

#### 3.1.2 実装場所
- **ファイル**: `server/cmd/jobqueue/main.go`
- **ルーター**: 標準ライブラリの`net/http`を使用（シンプルな実装）
- **実装方法**: Adminサーバーと同様のシンプルな実装
  - 認証ミドルウェアを通過しない
  - 最小限の実装（`"OK"`を返すのみ）

#### 3.1.3 実装の詳細
- 標準ライブラリの`net/http`を使用してHTTPサーバーを起動
- `http.HandleFunc`を使用して`/health`エンドポイントを登録
- ハンドラー関数はシンプルに`200 OK`と`"OK"`を返す
- エラーハンドリングは不要（常に成功を返す）

### 3.3 設定ファイルの拡張

#### 3.3.1 設定構造の追加
- **設定ファイル**: `config/develop.yaml`等
- **設定項目**: `jobqueue`セクションを追加
  ```yaml
  jobqueue:
    port: 8082
    read_timeout: 30s
    write_timeout: 30s
  ```
- **設定構造体**: `server/internal/config/config.go`に`JobQueueConfig`構造体を追加
- **設定読み込み**: 既存の設定読み込み処理を使用

#### 3.3.2 環境別設定
- **開発環境**: `config/develop/config.yaml`に追加
- **ステージング環境**: `config/staging/config.yaml`に追加
- **本番環境**: `config/production/config.yaml.example`に追加

## 4. 非機能要件

### 4.1 パフォーマンス
- **レスポンス時間**: 1ms以下（シンプルな実装のため）
- **リソース使用量**: 最小限（追加のリソース消費は不要）
- **HTTPサーバーのオーバーヘッド**: 最小限（標準ライブラリの`net/http`を使用）

### 4.2 可用性
- **可用性**: サーバーが起動している限り、常に`200 OK`を返す
- **エラー処理**: エラーが発生しない実装（常に成功を返す）
- **Asynqサーバーとの関係**: HTTPサーバーはAsynqサーバーと独立して動作する

### 4.3 セキュリティ
- **認証**: 不要（ヘルスチェック用のため）
- **情報漏洩**: サーバーの内部情報を返さない（`"OK"`のみ）
- **ポート**: 必要に応じてファイアウォールで制限可能

### 4.4 保守性
- **コードの簡潔性**: Adminサーバーと同様のシンプルな実装
- **一貫性**: 他のサーバー（API、Admin）の実装と同様のパターンを使用
- **設定の一貫性**: 既存の設定構造（`server`、`admin`）と同様の構造を使用

### 4.5 起動と停止
- **起動**: HTTPサーバーとAsynqサーバーを並行して起動
- **停止**: Graceful shutdownを実装（既存の実装を維持）
- **シグナル処理**: SIGINT、SIGTERMシグナルを受信した場合、両方のサーバーを停止

## 5. 制約事項

### 5.1 技術的制約
- **HTTPサーバー**: 標準ライブラリの`net/http`を使用（シンプルな実装のため）
- **ポート**: 8082を使用（他のサーバーと重複しない）
- **既存機能への影響**: Asynqサーバーの動作に影響を与えない

### 5.2 実装上の制約
- **認証ミドルウェア**: `/health`エンドポイントは認証を通過しない
- **既存コードへの影響**: 既存のAsynqサーバーの実装に影響を与えない
- **設定ファイル**: 既存の設定ファイル構造を維持する

### 5.3 動作環境
- **ローカル環境**: ローカル環境でも動作することを確認
- **Docker環境**: docker-compose等で動作することを確認（将来の拡張）

## 6. 受け入れ基準

### 6.1 HTTPサーバーの追加
- [ ] `server/cmd/jobqueue/main.go`にHTTPサーバーの起動処理が追加されている
- [ ] HTTPサーバーがAsynqサーバーと並行して動作する
- [ ] Graceful shutdownが正常に動作する（SIGINT、SIGTERMで停止）

### 6.2 `/health`エンドポイントの実装
- [ ] `GET /health`エンドポイントが実装されている
- [ ] エンドポイントが認証なしでアクセス可能である
- [ ] エンドポイントが`200 OK`と`"OK"`を返す
- [ ] エンドポイントが`text/plain`のContent-Typeを返す

### 6.3 設定ファイルの拡張
- [ ] `config/develop/config.yaml`に`jobqueue`セクションが追加されている
- [ ] `config/staging/config.yaml`に`jobqueue`セクションが追加されている
- [ ] `config/production/config.yaml.example`に`jobqueue`セクションが追加されている
- [ ] `server/internal/config/config.go`に`JobQueueConfig`構造体が追加されている
- [ ] 設定ファイルからポート番号が正しく読み込まれる

### 6.4 動作確認
- [ ] ローカル環境で`curl http://localhost:8082/health`が正常に動作する
- [ ] JobQueueサーバーが起動した時、HTTPサーバーとAsynqサーバーの両方が動作する
- [ ] 既存のAsynqサーバーの機能が正常に動作することを確認
- [ ] Graceful shutdownが正常に動作する

### 6.5 テスト
- [ ] 単体テストが実装されている（該当する場合）
- [ ] 統合テストが実装されている（該当する場合）
- [ ] 既存のテストが全て失敗しないことを確認

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### 新規作成・修正するファイル
- `server/cmd/jobqueue/main.go`: HTTPサーバーの起動処理と`/health`エンドポイントの実装を追加
- `server/internal/config/config.go`: `JobQueueConfig`構造体を追加
- `config/develop/config.yaml`: `jobqueue`セクションを追加
- `config/staging/config.yaml`: `jobqueue`セクションを追加
- `config/production/config.yaml.example`: `jobqueue`セクションを追加

#### 確認が必要なファイル
- `server/cmd/jobqueue/main.go`: 既存のAsynqサーバーの実装を確認（変更不要）

### 7.2 既存機能への影響
- **既存のAsynqサーバー**: 影響なし（HTTPサーバーは独立して動作）
- **ジョブ処理機能**: 影響なし（既存の実装を維持）
- **設定ファイル**: 既存の設定項目に影響なし（新規追加のみ）

### 7.3 テストへの影響
- **既存のテスト**: 影響なし（新規エンドポイントの追加のみ）
- **新規テスト**: `/health`エンドポイントのテストを追加する可能性がある

## 8. 実装上の注意事項

### 8.1 HTTPサーバー実装の注意事項
- **並行実行**: HTTPサーバーとAsynqサーバーをgoroutineで並行して起動
- **Graceful shutdown**: 両方のサーバーを適切に停止する
- **エラーハンドリング**: HTTPサーバーの起動エラーを適切に処理する
- **実装の簡潔性**: Adminサーバーと同様のシンプルな実装を維持する

### 8.2 設定ファイル実装の注意事項
- **設定構造**: 既存の`ServerConfig`、`AdminConfig`と同様の構造を使用
- **デフォルト値**: ポート番号のデフォルト値を適切に設定
- **環境別設定**: 各環境の設定ファイルに適切な値を設定

### 8.3 テストの注意事項
- **単体テスト**: エンドポイントが正常に動作することを確認するテストを追加する
- **統合テスト**: HTTPサーバーとAsynqサーバーが並行して動作することを確認する
- **既存テスト**: 既存のテストが全て失敗しないことを確認する

### 8.4 動作確認の注意事項
- **ローカル環境**: ローカル環境で`curl http://localhost:8082/health`が正常に動作することを確認
- **既存機能**: 既存のAsynqサーバーの機能が正常に動作することを確認
- **Graceful shutdown**: 両方のサーバーが適切に停止することを確認

## 9. 参考情報

### 9.1 既存実装の参考
- **Adminサーバー**: `server/cmd/admin/main.go`の`/health`エンドポイント実装
  ```go
  app.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
      w.Header().Set("Content-Type", "text/plain")
      w.WriteHeader(http.StatusOK)
      w.Write([]byte("OK"))
  }).Methods("GET")
  ```
- **APIサーバー**: `server/internal/api/router/router.go`の`/health`エンドポイント実装
  ```go
  e.GET("/health", func(c echo.Context) error {
      return c.String(http.StatusOK, "OK")
  })
  ```

### 9.2 設定ファイルの参考
- **APIサーバー設定**: `config/develop/config.yaml`の`server`セクション
  ```yaml
  server:
    port: 8080
    read_timeout: 30s
    write_timeout: 30s
  ```
- **Adminサーバー設定**: `config/develop/config.yaml`の`admin`セクション
  ```yaml
  admin:
    port: 8081
    read_timeout: 30s
    write_timeout: 30s
  ```

### 9.3 技術スタック
- **言語**: Go
- **HTTPサーバー**: 標準ライブラリの`net/http`
- **設定管理**: `github.com/spf13/viper`（既存システムと同様）
- **コンテナ管理**: Docker Compose（将来の拡張）

### 9.4 関連ドキュメント
- `server/cmd/jobqueue/main.go`: JobQueueサーバーのメインエントリーポイント
- `server/cmd/admin/main.go`: Adminサーバーのメインエントリーポイント（参考）
- `server/internal/config/config.go`: 設定構造体の定義
- `config/develop/config.yaml`: 開発環境の設定ファイル
