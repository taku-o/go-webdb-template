# CLIツール対応要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #9
- **Issueタイトル**: バッチ処理対応
- **Feature名**: 0007-cli
- **作成日**: 2025-01-27

### 1.2 目的
cronなどから処理を実行できるようにするため、CLIツールの基盤を構築する。
サンプルとして、ユーザー一覧を出力するCLIツールを実装し、今後のバッチ処理実装の参考とする。

### 1.3 スコープ
- CLIツールの配置場所と実行ファイル生成場所の決定
- ユーザー一覧を出力するCLIツールの実装
- limitパラメータによる出力件数制限機能の実装
- 既存の設定読み込み機能とDB接続機能の再利用

## 2. 背景・現状分析

### 2.1 現在の実装
- **コマンド構造**: `server/cmd/`配下に`server/main.go`と`admin/main.go`が存在
- **設定読み込み**: `internal/config/config.go`の`Load()`関数で環境変数`APP_ENV`に基づいて設定を読み込み
- **DB接続**: `internal/db/manager.go`の`NewGORMManager()`でGORM接続を初期化
- **ユーザー一覧取得**: `internal/service/user_service.go`の`ListUsers()`メソッドでユーザー一覧を取得可能
- **Repository層**: `internal/repository/user_repository_gorm.go`の`List()`メソッドでクロスシャードクエリを実行

### 2.2 課題点
1. **バッチ処理実行手段の欠如**: cronなどから処理を実行するためのCLIツールが存在しない
2. **CLIツール配置場所の未定義**: バッチ処理用のコマンドを配置する場所が明確でない
3. **実行ファイル生成場所の未定義**: ビルド後の実行ファイルの配置場所が明確でない

### 2.3 本実装による改善点
1. **バッチ処理実行基盤の構築**: CLIツールの基盤を構築し、今後のバッチ処理実装を容易にする
2. **明確なディレクトリ構造**: CLIツールの配置場所と実行ファイル生成場所を明確化
3. **再利用可能な実装**: 既存の設定読み込みとDB接続機能を再利用し、一貫性を保つ

## 3. 機能要件

### 3.1 CLIツールの配置場所

#### 3.1.1 ディレクトリ構造
既存の`server/cmd/`配下にCLIツールを配置する：
- `server/cmd/list-users/main.go`: ユーザー一覧を出力するCLIツール

#### 3.1.2 命名規則
- 各CLIツールは`server/cmd/{command-name}/main.go`の形式で配置
- コマンド名はkebab-case（例: `list-users`, `export-data`）

### 3.2 実行ファイルの生成場所

#### 3.2.1 ビルド出力先
- 実行ファイルは`server/bin/`ディレクトリに生成する
- ファイル名は`{command-name}`（例: `list-users`）

#### 3.2.2 ビルドコマンド
```bash
# 開発環境でのビルド
cd server
go build -o bin/list-users ./cmd/list-users

# 本番環境でのビルド（クロスコンパイル対応）
GOOS=linux GOARCH=amd64 go build -o bin/list-users ./cmd/list-users
```

### 3.3 ユーザー一覧出力機能

#### 3.3.1 基本機能
- すべてのシャードからユーザー一覧を取得し、標準出力に出力する
- 既存の`UserService.ListUsers()`メソッドを利用する

#### 3.3.2 出力形式
- デフォルトはTSV（タブ区切り）形式で出力
- 出力項目: ID, Name, Email, CreatedAt, UpdatedAt
- ヘッダー行を含める

#### 3.3.3 limitパラメータ
- `--limit`フラグで出力件数を制限できる
- デフォルト値: 20件
- 最大値: 100件（既存の`UserService.ListUsers()`の制限に準拠）
- 最小値: 1件

#### 3.3.4 コマンドライン引数
```bash
# デフォルト（20件）
./bin/list-users

# limit指定
./bin/list-users --limit 50

# 環境指定
APP_ENV=production ./bin/list-users --limit 100
```

### 3.4 設定読み込みとDB接続

#### 3.4.1 設定読み込み
- 既存の`config.Load()`関数を使用
- 環境変数`APP_ENV`で環境を切り替え（develop/staging/production）

#### 3.4.2 DB接続
- 既存の`db.NewGORMManager()`でGORM接続を初期化
- すべてのシャードへの接続確認（`PingAll()`）を実行
- 処理終了時に`CloseAll()`で接続をクローズ

#### 3.4.3 エラーハンドリング
- 設定読み込みエラー、DB接続エラー、クエリエラーを適切に処理
- エラー時は標準エラー出力にエラーメッセージを出力し、非ゼロの終了コードで終了

## 4. 非機能要件

### 4.1 cron実行対応
- 標準出力と標準エラー出力を適切に分離
- 非対話的な実行を前提とする（プロンプト表示なし）
- 終了コードを適切に返す（成功: 0、エラー: 1以上）

### 4.2 パフォーマンス
- 大量データ取得時のメモリ使用量を考慮
- limitパラメータによる制限で、メモリ使用量を制御

### 4.3 保守性
- 既存のアーキテクチャ（Repository層、Service層）を再利用
- コードの重複を最小化
- 明確なエラーメッセージを出力

### 4.4 拡張性
- 将来的に他のCLIツールを追加しやすい構造
- 共通処理（設定読み込み、DB接続）を再利用可能にする

## 5. 制約事項

### 5.1 技術的制約
- Go言語の標準ライブラリ（`flag`パッケージ）を使用してCLI引数を処理
- 既存の設定読み込み機能（`internal/config`）を変更しない
- 既存のDB接続機能（`internal/db`）を変更しない

### 5.2 プロジェクト制約
- 既存のアーキテクチャ（レイヤードアーキテクチャ）を維持
- 既存のRepository層、Service層を再利用
- 既存のテストコードへの影響を最小化

### 5.3 ディレクトリ構造
- `server/cmd/`配下にCLIツールを配置
- `server/bin/`ディレクトリに実行ファイルを生成（`.gitignore`に追加）

## 6. 受け入れ基準

### 6.1 機能要件
- [ ] `server/cmd/list-users/main.go`が作成されている
- [ ] `server/bin/`ディレクトリが作成されている（`.gitignore`に追加）
- [ ] `go build`コマンドで実行ファイルが生成できる
- [ ] デフォルトで20件のユーザー一覧がTSV形式で出力される
- [ ] `--limit`フラグで出力件数を制限できる
- [ ] limit値が1未満の場合はエラーを出力する
- [ ] limit値が100を超える場合は100に制限される
- [ ] 環境変数`APP_ENV`で環境を切り替えられる
- [ ] すべてのシャードからユーザーを取得できる（クロスシャードクエリ）

### 6.2 非機能要件
- [ ] cronから実行可能（非対話的実行）
- [ ] エラー時は適切な終了コードを返す
- [ ] 標準出力と標準エラー出力が適切に分離されている
- [ ] 既存のアーキテクチャを維持している
- [ ] 既存のテストコードが正常に動作する

### 6.3 ドキュメント
- [ ] README.mdにCLIツールの使用方法が記載されている（必要に応じて）
- [ ] `.gitignore`に`server/bin/`が追加されている

## 7. 影響範囲

### 7.1 新規追加が必要なディレクトリ・ファイル

#### ディレクトリ
- `server/cmd/list-users/`: ユーザー一覧出力CLIツール
- `server/bin/`: 実行ファイル生成先（`.gitignore`に追加）

#### ファイル
- `server/cmd/list-users/main.go`: ユーザー一覧出力CLIツールのメインファイル

### 7.2 変更が必要なファイル

#### 設定ファイル
- `.gitignore`: `server/bin/`を追加

#### ドキュメント
- `README.md`: CLIツールの使用方法を追加（必要に応じて）

### 7.3 削除されるファイル
なし

## 8. 実装上の注意事項

### 8.1 CLI引数の処理
- Go言語の標準ライブラリ`flag`パッケージを使用
- `flag.Int()`で`--limit`フラグを定義
- `flag.Parse()`で引数を解析

### 8.2 出力形式
- TSV形式で出力（タブ区切り）
- ヘッダー行を最初に出力
- 各ユーザーの情報を1行ずつ出力
- 日時はRFC3339形式で出力

### 8.3 エラーハンドリング
- 設定読み込みエラー: `log.Fatalf()`でエラーメッセージを出力して終了
- DB接続エラー: `log.Fatalf()`でエラーメッセージを出力して終了
- クエリエラー: `log.Printf()`でエラーメッセージを出力し、`os.Exit(1)`で終了
- 引数エラー: `flag.Usage()`で使用方法を表示し、`os.Exit(1)`で終了

### 8.4 リソース管理
- DB接続は`defer gormManager.CloseAll()`で確実にクローズ
- コンテキストは`context.Background()`を使用

### 8.5 既存コードの再利用
- `config.Load()`で設定を読み込み
- `db.NewGORMManager(cfg)`でDB接続を初期化
- `repository.NewUserRepositoryGORM(gormManager)`でRepositoryを初期化
- `service.NewUserService(userRepo)`でServiceを初期化
- `userService.ListUsers(ctx, limit, 0)`でユーザー一覧を取得

## 9. 参考情報

### 9.1 関連Issue
- GitHub Issue #9: バッチ処理対応

### 9.2 既存ドキュメント
- `server/cmd/server/main.go`: サーバー起動コマンドの実装例
- `server/cmd/admin/main.go`: 管理画面起動コマンドの実装例
- `server/internal/config/config.go`: 設定読み込み実装
- `server/internal/db/manager.go`: DB接続管理実装
- `server/internal/service/user_service.go`: ユーザーService実装

### 9.3 既存実装
- `server/internal/repository/user_repository_gorm.go`: ユーザーRepository実装（GORM版）
- `server/internal/model/user.go`: ユーザーモデル定義

### 9.4 Go言語標準ライブラリ
- `flag`パッケージ: https://pkg.go.dev/flag
- `os`パッケージ: https://pkg.go.dev/os
- `context`パッケージ: https://pkg.go.dev/context

