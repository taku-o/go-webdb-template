# server/Dockerfile* の更新要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #90
- **親Issue番号**: #85
- **Issueタイトル**: server/Dockerfile* の更新
- **Feature名**: 0047-pgmain-docker
- **作成日**: 2025-01-27
- **対象ブランチ**: switch-to-postgresql

### 1.2 目的
SQLiteを利用しないことになったため、`server/Dockerfile.develop`と`server/Dockerfile.admin.develop`を削除し、`server/Dockerfile`と`server/Dockerfile.admin`に統合する。また、CGO_ENABLED=0を採用できるため、DebianベースからAlpineベースのイメージに書き換えて、イメージサイズを削減し、ビルド時間を短縮する。

### 1.3 スコープ
- `server/Dockerfile.develop`の削除
- `server/Dockerfile.admin.develop`の削除
- `server/Dockerfile`のAlpineベースへの書き換え（CGO_ENABLED=0維持）
- `server/Dockerfile.admin`のAlpineベースへの書き換え（CGO_ENABLED=0維持）
- `docker-compose.api.develop.yml`の更新（`Dockerfile.develop`参照を`Dockerfile`に変更）
- `docker-compose.admin.develop.yml`の更新（`Dockerfile.admin.develop`参照を`Dockerfile.admin`に変更）
- ドキュメントの更新（`docs/Docker.md`）

**本実装の範囲外**:
- クライアント側のDockerfile変更（`client/Dockerfile`は変更しない）
- その他のdocker-composeファイルの変更（staging/production環境用は変更しない）
- PostgreSQL/Redisコンテナの設定変更

## 2. 背景・現状分析

### 2.1 現在の実装
- **Dockerfile構成**:
  - `server/Dockerfile`: staging/production用（CGO_ENABLED=0、Debianベース）
  - `server/Dockerfile.develop`: develop用（CGO_ENABLED=1、Debianベース、SQLite対応）
  - `server/Dockerfile.admin`: staging/production用（CGO_ENABLED=0、Debianベース）
  - `server/Dockerfile.admin.develop`: develop用（CGO_ENABLED=1、Debianベース、SQLite対応）
- **ベースイメージ**: 
  - ビルドステージ: `golang:1.24-bookworm`
  - 実行ステージ: `debian:bookworm-slim`
- **CGO設定**:
  - `Dockerfile`/`Dockerfile.admin`: CGO_ENABLED=0
  - `Dockerfile.develop`/`Dockerfile.admin.develop`: CGO_ENABLED=1（SQLite用）
- **docker-compose設定**:
  - `docker-compose.api.develop.yml`: `Dockerfile.develop`を参照
  - `docker-compose.admin.develop.yml`: `Dockerfile.admin.develop`を参照
  - `docker-compose.api.staging.yml`: `Dockerfile`を参照
  - `docker-compose.api.production.yml`: `Dockerfile`を参照
  - `docker-compose.admin.staging.yml`: `Dockerfile.admin`を参照
  - `docker-compose.admin.production.yml`: `Dockerfile.admin`を参照

### 2.2 課題点
1. **SQLite依存の不要なDockerfile**: SQLiteを利用しないことになったため、`Dockerfile.develop`と`Dockerfile.admin.develop`が不要になった
2. **CGO依存の削除**: SQLiteを使用しないため、CGO_ENABLED=1は不要で、CGO_ENABLED=0を採用できる
3. **イメージサイズの最適化**: DebianベースからAlpineベースに変更することで、イメージサイズを削減できる
4. **ビルド時間の短縮**: Alpineベースのイメージは軽量で、ビルド時間を短縮できる
5. **設定の重複**: develop環境とstaging/production環境で異なるDockerfileを管理する必要があり、保守性が低い

### 2.3 本実装による改善点
1. **Dockerfileの統合**: develop環境とstaging/production環境で同じDockerfileを使用し、保守性を向上
2. **CGO_ENABLED=0の採用**: SQLite依存を削除し、CGO_ENABLED=0でビルドすることで、クロスコンパイルが容易になり、セキュリティも向上
3. **Alpineベースへの移行**: イメージサイズを削減し、ビルド時間を短縮
4. **設定の簡素化**: Dockerfileの数を削減し、管理を簡素化

## 3. 機能要件

### 3.1 Dockerfileの統合と削除

#### 3.1.1 Dockerfile.developの削除
- **対象ファイル**: `server/Dockerfile.develop`
- **削除理由**: SQLiteを利用しないため、CGO_ENABLED=1のDockerfileが不要になった
- **削除後の対応**: `docker-compose.api.develop.yml`を`Dockerfile`を参照するように変更

#### 3.1.2 Dockerfile.admin.developの削除
- **対象ファイル**: `server/Dockerfile.admin.develop`
- **削除理由**: SQLiteを利用しないため、CGO_ENABLED=1のDockerfileが不要になった
- **削除後の対応**: `docker-compose.admin.develop.yml`を`Dockerfile.admin`を参照するように変更

### 3.2 DockerfileのAlpineベースへの書き換え

#### 3.2.1 server/Dockerfileの書き換え
- **ベースイメージ変更**:
  - ビルドステージ: `golang:1.24-bookworm` → `golang:1.24-alpine`
  - 実行ステージ: `debian:bookworm-slim` → `alpine:latest`
- **CGO設定**: CGO_ENABLED=0を維持
- **パッケージインストール**:
  - Debian: `apt-get` → Alpine: `apk add`
  - 必要なパッケージ: `ca-certificates`, `tzdata`
  - 非rootユーザー作成: Alpine用のコマンドに変更
- **ディレクトリ構造**: 既存のディレクトリ構造を維持
- **ビルドコマンド**: CGO_ENABLED=0を維持し、Alpine環境で正常にビルドできることを確認

#### 3.2.2 server/Dockerfile.adminの書き換え
- **ベースイメージ変更**:
  - ビルドステージ: `golang:1.24-bookworm` → `golang:1.24-alpine`
  - 実行ステージ: `debian:bookworm-slim` → `alpine:latest`
- **CGO設定**: CGO_ENABLED=0を維持
- **パッケージインストール**:
  - Debian: `apt-get` → Alpine: `apk add`
  - 必要なパッケージ: `ca-certificates`, `tzdata`
  - 非rootユーザー作成: Alpine用のコマンドに変更
- **ディレクトリ構造**: 既存のディレクトリ構造を維持
- **ビルドコマンド**: CGO_ENABLED=0を維持し、Alpine環境で正常にビルドできることを確認

### 3.3 docker-composeファイルの更新

#### 3.3.1 docker-compose.api.develop.ymlの更新
- **変更内容**: `dockerfile: Dockerfile.develop` → `dockerfile: Dockerfile`
- **確認事項**: 他の設定（ports、environment、volumes、networks等）は変更しない

#### 3.3.2 docker-compose.admin.develop.ymlの更新
- **変更内容**: `dockerfile: Dockerfile.admin.develop` → `dockerfile: Dockerfile.admin`
- **確認事項**: 他の設定（ports、environment、volumes、networks等）は変更しない

### 3.4 ドキュメントの更新

#### 3.4.1 docs/Docker.mdの更新
- **Dockerfile構成表の更新**: 
  - `Dockerfile.develop`と`Dockerfile.admin.develop`の行を削除
  - `Dockerfile`と`Dockerfile.admin`のベースイメージをAlpineに更新
  - CGO設定をCGO_ENABLED=0に統一
- **ビルド・起動・停止コマンド**: 変更なし（既存のコマンドが正常に動作することを確認）

## 4. 非機能要件

### 4.1 イメージサイズ
- **目標**: Alpineベースに変更することで、イメージサイズを削減
- **現状**: Debianベース（`debian:bookworm-slim`）のイメージサイズ
- **期待値**: Alpineベース（`alpine:latest`）のイメージサイズは、Debianベースより大幅に小さい

### 4.2 ビルド時間
- **目標**: Alpineベースに変更することで、ビルド時間を短縮
- **現状**: Debianベースでのビルド時間
- **期待値**: Alpineベースでのビルド時間は、Debianベースより短い

### 4.3 互換性
- **既存機能の維持**: Alpineベースに変更しても、既存の機能が正常に動作することを確認
- **バイナリの互換性**: CGO_ENABLED=0でビルドしたバイナリが正常に動作することを確認
- **ネットワーク設定**: 既存のdocker-compose設定との互換性を維持

### 4.4 セキュリティ
- **非rootユーザー**: 既存の非rootユーザー（appuser）設定を維持
- **パッケージの更新**: Alpineのパッケージを最新の状態に保つ
- **CGO無効化**: CGO_ENABLED=0により、セキュリティリスクを低減

### 4.5 環境別対応
- **開発環境**: `docker-compose.api.develop.yml`と`docker-compose.admin.develop.yml`で`Dockerfile`と`Dockerfile.admin`を使用
- **staging環境**: 既存の`docker-compose.api.staging.yml`と`docker-compose.admin.staging.yml`で`Dockerfile`と`Dockerfile.admin`を使用（変更なし）
- **production環境**: 既存の`docker-compose.api.production.yml`と`docker-compose.admin.production.yml`で`Dockerfile`と`Dockerfile.admin`を使用（変更なし）

## 5. 制約事項

### 5.1 既存システムとの関係
- **PostgreSQL接続**: Alpineベースに変更しても、PostgreSQL接続が正常に動作することを確認
- **Redis接続**: Alpineベースに変更しても、Redis接続が正常に動作することを確認
- **設定ファイル**: 既存の設定ファイル（`config/{env}/database.yaml`等）との互換性を維持
- **ボリュームマウント**: 既存のボリュームマウント設定との互換性を維持

### 5.2 技術スタック
- **Goバージョン**: Go 1.24を維持
- **ベースイメージ**: `golang:1.24-alpine`（ビルドステージ）、`alpine:latest`（実行ステージ）
- **CGO設定**: CGO_ENABLED=0を維持
- **ビルドツール**: Alpine環境で必要なビルドツールをインストール

### 5.3 依存関係
- **PostgreSQLドライバー**: `gorm.io/driver/postgres`はCGO不要で動作する
- **その他の依存関係**: CGO_ENABLED=0でビルドできることを確認

### 5.4 運用上の制約
- **ビルド環境**: Alpineベースのイメージをビルドできる環境が必要
- **実行環境**: Alpineベースのイメージが実行できる環境が必要
- **互換性**: 既存のdocker-compose設定との互換性を維持

## 6. 受け入れ基準

### 6.1 Dockerfileの統合と削除
- [ ] `server/Dockerfile.develop`が削除されている
- [ ] `server/Dockerfile.admin.develop`が削除されている
- [ ] `docker-compose.api.develop.yml`が`Dockerfile`を参照している
- [ ] `docker-compose.admin.develop.yml`が`Dockerfile.admin`を参照している

### 6.2 DockerfileのAlpineベースへの書き換え
- [ ] `server/Dockerfile`がAlpineベースに書き換えられている
- [ ] `server/Dockerfile.admin`がAlpineベースに書き換えられている
- [ ] ビルドステージが`golang:1.24-alpine`を使用している
- [ ] 実行ステージが`alpine:latest`を使用している
- [ ] CGO_ENABLED=0でビルドされている
- [ ] 必要なパッケージ（`ca-certificates`, `tzdata`）がインストールされている
- [ ] 非rootユーザー（appuser）が作成されている
- [ ] ディレクトリ構造が既存と同じである

### 6.3 ビルドと動作確認
- [ ] `docker-compose.api.develop.yml`でイメージが正常にビルドできる
- [ ] `docker-compose.admin.develop.yml`でイメージが正常にビルドできる
- [ ] ビルドしたイメージからコンテナが正常に起動できる
- [ ] APIサーバーが正常に動作する（PostgreSQL接続確認）
- [ ] Adminサーバーが正常に動作する（PostgreSQL接続確認）
- [ ] ヘルスチェックが正常に動作する

### 6.4 イメージサイズとビルド時間
- [ ] AlpineベースのイメージサイズがDebianベースより小さいことを確認
- [ ] Alpineベースのビルド時間がDebianベースより短いことを確認（オプション）

### 6.5 ドキュメント
- [ ] `docs/Docker.md`のDockerfile構成表が更新されている
- [ ] `Dockerfile.develop`と`Dockerfile.admin.develop`の記述が削除されている
- [ ] `Dockerfile`と`Dockerfile.admin`のベースイメージがAlpineに更新されている
- [ ] CGO設定がCGO_ENABLED=0に統一されている

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### 削除するファイル
- `server/Dockerfile.develop`
- `server/Dockerfile.admin.develop`

#### 変更するファイル
- `server/Dockerfile`: Alpineベースに書き換え
- `server/Dockerfile.admin`: Alpineベースに書き換え
- `docker-compose.api.develop.yml`: `Dockerfile.develop`参照を`Dockerfile`に変更
- `docker-compose.admin.develop.yml`: `Dockerfile.admin.develop`参照を`Dockerfile.admin`に変更
- `docs/Docker.md`: Dockerfile構成表と説明を更新

### 7.2 既存ファイルの扱い
- `server/Dockerfile`: 既存の設定をAlpineベースに書き換え（CGO_ENABLED=0維持）
- `server/Dockerfile.admin`: 既存の設定をAlpineベースに書き換え（CGO_ENABLED=0維持）
- `docker-compose.api.develop.yml`: `dockerfile`指定のみ変更
- `docker-compose.admin.develop.yml`: `dockerfile`指定のみ変更
- `docker-compose.api.staging.yml`: 変更なし（既に`Dockerfile`を参照）
- `docker-compose.api.production.yml`: 変更なし（既に`Dockerfile`を参照）
- `docker-compose.admin.staging.yml`: 変更なし（既に`Dockerfile.admin`を参照）
- `docker-compose.admin.production.yml`: 変更なし（既に`Dockerfile.admin`を参照）

## 8. 実装上の注意事項

### 8.1 Alpineベースへの移行
- **パッケージマネージャー**: `apt-get`から`apk add`に変更
- **パッケージ名**: DebianとAlpineでパッケージ名が異なる場合があるため、適切なパッケージ名を確認
- **非rootユーザー作成**: Alpine用のコマンド（`addgroup`, `adduser`）を使用
- **タイムゾーンデータ**: `tzdata`パッケージをインストール

### 8.2 CGO_ENABLED=0の確認
- **ビルド確認**: CGO_ENABLED=0でビルドできることを確認
- **依存関係確認**: すべての依存関係がCGO不要であることを確認
- **動作確認**: ビルドしたバイナリが正常に動作することを確認

### 8.3 ビルドと動作確認
- **ビルド確認**: `docker-compose.api.develop.yml`と`docker-compose.admin.develop.yml`でイメージが正常にビルドできることを確認
- **起動確認**: ビルドしたイメージからコンテナが正常に起動できることを確認
- **接続確認**: PostgreSQLとRedisへの接続が正常に動作することを確認
- **ヘルスチェック**: ヘルスチェックが正常に動作することを確認

### 8.4 ドキュメント整備
- **Dockerfile構成表**: `docs/Docker.md`のDockerfile構成表を更新
- **ベースイメージ**: Alpineベースに変更したことを明記
- **CGO設定**: CGO_ENABLED=0に統一したことを明記
- **削除ファイル**: `Dockerfile.develop`と`Dockerfile.admin.develop`を削除したことを明記

### 8.5 互換性確認
- **既存機能**: Alpineベースに変更しても、既存の機能が正常に動作することを確認
- **設定ファイル**: 既存の設定ファイルとの互換性を維持
- **ネットワーク設定**: 既存のdocker-compose設定との互換性を維持

## 9. 参考情報

### 9.1 関連Issue
- GitHub Issue #85: 開発環境はPostgreSQLを利用する前提とする
- GitHub Issue #90: server/Dockerfile* の更新

### 9.2 既存ドキュメント
- `docs/Docker.md`: Docker環境の説明
- `server/Dockerfile`: 既存のAPIサーバー用Dockerfile
- `server/Dockerfile.admin`: 既存のAdminサーバー用Dockerfile
- `docker-compose.api.develop.yml`: 開発環境用APIサーバーのdocker-compose設定
- `docker-compose.admin.develop.yml`: 開発環境用Adminサーバーのdocker-compose設定

### 9.3 技術スタック
- **Go**: 1.24
- **ベースイメージ**: `golang:1.24-alpine`（ビルドステージ）、`alpine:latest`（実行ステージ）
- **CGO**: CGO_ENABLED=0
- **PostgreSQLドライバー**: `gorm.io/driver/postgres`（CGO不要）

### 9.4 参考リンク
- Alpine Linux公式ドキュメント: https://alpinelinux.org/documentation/
- Docker公式ドキュメント: https://docs.docker.com/
- Go公式ドキュメント: https://go.dev/doc/
