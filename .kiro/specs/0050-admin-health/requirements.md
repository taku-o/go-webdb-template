# GoAdmin死活監視エンドポイントの要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #103
- **Issueタイトル**: GoAdmin死活監視エンドポイントの作成
- **Feature名**: 0050-admin-health
- **作成日**: 2026-01-10

### 1.2 目的
APIサーバーとクライアントサーバーには`/health`という死活監視用のエンドポイントが用意されている。この`/health`エンドポイントをAdminサーバーにも実装し、docker-compose.admin.ymlのヘルスチェック機能を有効化する。

### 1.3 スコープ
- Adminサーバーに`/health`エンドポイントを実装する
- エンドポイントは認証不要でアクセス可能とする
- docker-compose.admin.ymlのヘルスチェック機能が正常に動作することを確認する

**本実装の範囲外**:
- APIサーバーやクライアントサーバーの`/health`エンドポイントの変更
- その他のエンドポイントの追加や変更
- ヘルスチェックの詳細な診断機能（現時点では単純なOK応答のみ）

## 2. 背景・現状分析

### 2.1 現在の状況
- **APIサーバー**: `server/internal/api/router/router.go`に`/health`エンドポイントが実装されている
  - 実装内容: `GET /health`で`200 OK`と`"OK"`という文字列を返す
  - 認証: 不要
  - ポート: 8080
- **クライアントサーバー**: `/health`エンドポイントが実装されている（詳細は未確認）
- **Adminサーバー**: `/health`エンドポイントが未実装
  - ポート: 8081
  - ルーター: Gorilla Mux Routerを使用
  - フレームワーク: GoAdmin Engineを使用
- **docker-compose.admin.yml**: 既にヘルスチェック設定が存在するが、エンドポイントが未実装のため動作していない
  ```yaml
  healthcheck:
    test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8081/health"]
    interval: 30s
    timeout: 10s
    retries: 3
    start_period: 40s
  ```

### 2.2 課題点
1. **Adminサーバーのヘルスチェック不可**: docker-compose.admin.ymlのヘルスチェックが動作していない
2. **一貫性の欠如**: APIサーバーとクライアントサーバーには`/health`エンドポイントがあるが、Adminサーバーにはない
3. **監視ツールとの連携不可**: 外部監視ツールがAdminサーバーの死活監視を行うことができない

### 2.3 本実装による改善点
1. **ヘルスチェック機能の有効化**: docker-compose.admin.ymlのヘルスチェックが正常に動作する
2. **一貫性の確保**: 3つのサーバー（API、クライアント、Admin）すべてに`/health`エンドポイントが存在する
3. **監視ツールとの連携**: 外部監視ツールがAdminサーバーの死活監視を行えるようになる

## 3. 機能要件

### 3.1 `/health`エンドポイントの実装

#### 3.1.1 エンドポイント仕様
- **パス**: `/health`
- **メソッド**: `GET`
- **認証**: 不要（認証ミドルウェアを通過しない）
- **レスポンス**: 
  - ステータスコード: `200 OK`
  - レスポンスボディ: `"OK"`（文字列）
  - Content-Type: `text/plain`

#### 3.1.2 実装場所
- **ファイル**: `server/cmd/admin/main.go`
- **ルーター**: Gorilla Mux Router（`app := mux.NewRouter()`）
- **実装方法**: APIサーバーと同様のシンプルな実装
  - GoAdmin Engineのミドルウェアチェーンを通過しない
  - 認証ミドルウェアを通過しない
  - アクセスログミドルウェアは通過する可能性がある（実装による）

#### 3.1.3 実装の詳細
- Gorilla Mux Routerに直接エンドポイントを登録する
- ハンドラー関数はシンプルに`200 OK`と`"OK"`を返す
- エラーハンドリングは不要（常に成功を返す）

### 3.2 docker-compose.admin.ymlとの連携

#### 3.2.1 ヘルスチェック設定の確認
- 既存のヘルスチェック設定が正常に動作することを確認する
- 設定内容:
  ```yaml
  healthcheck:
    test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8081/health"]
    interval: 30s
    timeout: 10s
    retries: 3
    start_period: 40s
  ```

#### 3.2.2 動作確認
- docker-compose.admin.ymlを使用してコンテナを起動
- ヘルスチェックが正常に動作することを確認
- コンテナのステータスが`healthy`になることを確認

## 4. 非機能要件

### 4.1 パフォーマンス
- **レスポンス時間**: 1ms以下（シンプルな実装のため）
- **リソース使用量**: 最小限（追加のリソース消費は不要）

### 4.2 可用性
- **可用性**: サーバーが起動している限り、常に`200 OK`を返す
- **エラー処理**: エラーが発生しない実装（常に成功を返す）

### 4.3 セキュリティ
- **認証**: 不要（ヘルスチェック用のため）
- **情報漏洩**: サーバーの内部情報を返さない（`"OK"`のみ）

### 4.4 保守性
- **コードの簡潔性**: APIサーバーと同様のシンプルな実装
- **一貫性**: APIサーバーの実装と同様のパターンを使用

## 5. 制約事項

### 5.1 技術的制約
- **ルーター**: Gorilla Mux Routerを使用（既存の実装に合わせる）
- **フレームワーク**: GoAdmin Engineを使用（既存の実装に合わせる）
- **ポート**: 8081（既存の設定に合わせる）

### 5.2 実装上の制約
- **認証ミドルウェア**: `/health`エンドポイントは認証を通過しない
- **GoAdmin Engine**: GoAdmin Engineのミドルウェアチェーンを通過しない
- **既存コードへの影響**: 既存のエンドポイントや機能に影響を与えない

### 5.3 動作環境
- **Docker環境**: docker-compose.admin.ymlで動作することを確認
- **ローカル環境**: ローカル環境でも動作することを確認

## 6. 受け入れ基準

### 6.1 `/health`エンドポイントの実装
- [ ] `GET /health`エンドポイントが実装されている
- [ ] エンドポイントが認証なしでアクセス可能である
- [ ] エンドポイントが`200 OK`と`"OK"`を返す
- [ ] エンドポイントが`text/plain`のContent-Typeを返す

### 6.2 docker-compose.admin.ymlとの連携
- [ ] docker-compose.admin.ymlを使用してコンテナを起動できる
- [ ] ヘルスチェックが正常に動作する
- [ ] コンテナのステータスが`healthy`になる
- [ ] ヘルスチェックのログにエラーが表示されない

### 6.3 動作確認
- [ ] ローカル環境で`curl http://localhost:8081/health`が正常に動作する
- [ ] Docker環境で`wget --quiet --tries=1 --spider http://localhost:8081/health`が正常に動作する
- [ ] 既存のエンドポイント（`/admin`など）が正常に動作することを確認

### 6.4 テスト
- [ ] 単体テストが実装されている（該当する場合）
- [ ] 統合テストが実装されている（該当する場合）
- [ ] 既存のテストが全て失敗しないことを確認

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### 新規作成・修正するファイル
- `server/cmd/admin/main.go`: `/health`エンドポイントの実装を追加

#### 確認が必要なファイル
- `docker-compose.admin.yml`: 既存のヘルスチェック設定が正常に動作することを確認（変更不要）

### 7.2 既存機能への影響
- **既存のエンドポイント**: 影響なし（新規エンドポイントの追加のみ）
- **GoAdmin Engine**: 影響なし（GoAdmin Engineのミドルウェアチェーンを通過しない）
- **認証機能**: 影響なし（認証ミドルウェアを通過しない）

### 7.3 テストへの影響
- **既存のテスト**: 影響なし（新規エンドポイントの追加のみ）
- **新規テスト**: `/health`エンドポイントのテストを追加する可能性がある

## 8. 実装上の注意事項

### 8.1 エンドポイント実装の注意事項
- **ルーターへの登録**: Gorilla Mux Routerに直接エンドポイントを登録する
- **認証ミドルウェア**: `/health`エンドポイントは認証ミドルウェアを通過しないようにする
- **GoAdmin Engine**: GoAdmin Engineのミドルウェアチェーンを通過しないようにする
- **実装の簡潔性**: APIサーバーと同様のシンプルな実装を維持する

### 8.2 テストの注意事項
- **単体テスト**: エンドポイントが正常に動作することを確認するテストを追加する
- **統合テスト**: docker-compose.admin.ymlを使用した統合テストを実施する
- **既存テスト**: 既存のテストが全て失敗しないことを確認する

### 8.3 動作確認の注意事項
- **ローカル環境**: ローカル環境で`curl http://localhost:8081/health`が正常に動作することを確認
- **Docker環境**: docker-compose.admin.ymlを使用してコンテナを起動し、ヘルスチェックが正常に動作することを確認
- **既存機能**: 既存のエンドポイント（`/admin`など）が正常に動作することを確認

## 9. 参考情報

### 9.1 関連Issue
- GitHub Issue #103: GoAdmin死活監視エンドポイントの作成

### 9.2 既存実装の参考
- **APIサーバー**: `server/internal/api/router/router.go`の`/health`エンドポイント実装
  ```go
  e.GET("/health", func(c echo.Context) error {
      return c.String(http.StatusOK, "OK")
  })
  ```
- **docker-compose.api.yml**: APIサーバーのヘルスチェック設定
  ```yaml
  healthcheck:
    test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
  ```

### 9.3 技術スタック
- **言語**: Go
- **ルーター**: Gorilla Mux Router
- **フレームワーク**: GoAdmin Engine
- **コンテナ管理**: Docker Compose
- **ヘルスチェックツール**: wget

### 9.4 関連ドキュメント
- `docker-compose.admin.yml`: AdminサーバーのDocker Compose設定
- `server/cmd/admin/main.go`: Adminサーバーのメインエントリーポイント
- `server/internal/api/router/router.go`: APIサーバーのルーター実装（参考）
