# server/Dockerfile* の更新設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、SQLite依存を削除してDockerfileを統合し、Alpineベースのイメージに書き換えるための詳細設計を定義する。CGO_ENABLED=0を維持し、イメージサイズを削減し、ビルド時間を短縮する。

### 1.2 設計の範囲
- `server/Dockerfile.develop`の削除
- `server/Dockerfile.admin.develop`の削除
- `server/Dockerfile`のAlpineベースへの書き換え
- `server/Dockerfile.admin`のAlpineベースへの書き換え
- `docker-compose.api.develop.yml`の更新
- `docker-compose.admin.develop.yml`の更新
- `docs/Docker.md`の更新

### 1.3 設計方針
- **Alpineベースへの移行**: DebianベースからAlpineベースに変更し、イメージサイズを削減
- **CGO_ENABLED=0の維持**: SQLite依存を削除し、CGO_ENABLED=0でビルド
- **既存機能の維持**: Alpineベースに変更しても、既存の機能が正常に動作することを確認
- **設定の簡素化**: develop環境とstaging/production環境で同じDockerfileを使用
- **互換性の確保**: 既存のdocker-compose設定との互換性を維持

## 2. アーキテクチャ設計

### 2.1 既存アーキテクチャの分析

#### 2.1.1 現在の構成
- **Dockerfile構成**:
  - `server/Dockerfile`: staging/production用（CGO_ENABLED=0、Debianベース）
  - `server/Dockerfile.develop`: develop用（CGO_ENABLED=1、Debianベース、SQLite対応）→ 削除対象
  - `server/Dockerfile.admin`: staging/production用（CGO_ENABLED=0、Debianベース）
  - `server/Dockerfile.admin.develop`: develop用（CGO_ENABLED=1、Debianベース、SQLite対応）→ 削除対象
- **ベースイメージ**: 
  - ビルドステージ: `golang:1.24-bookworm`（Debianベース）
  - 実行ステージ: `debian:bookworm-slim`（Debianベース）
- **パッケージマネージャー**: `apt-get`（Debian）
- **非rootユーザー作成**: `groupadd`と`useradd`（Debian）
- **必要なパッケージ**: `ca-certificates`, `tzdata`
- **docker-compose設定**:
  - `docker-compose.api.develop.yml`: `Dockerfile.develop`を参照 → `Dockerfile`に変更
  - `docker-compose.admin.develop.yml`: `Dockerfile.admin.develop`を参照 → `Dockerfile.admin`に変更

#### 2.1.2 変更後の構成
- **Dockerfile構成**:
  - `server/Dockerfile`: 全環境共通（CGO_ENABLED=0、Alpineベース）
  - `server/Dockerfile.admin`: 全環境共通（CGO_ENABLED=0、Alpineベース）
- **ベースイメージ**: 
  - ビルドステージ: `golang:1.24-alpine`（Alpineベース）
  - 実行ステージ: `alpine:latest`（Alpineベース）
- **パッケージマネージャー**: `apk add`（Alpine）
- **非rootユーザー作成**: `addgroup`と`adduser`（Alpine）
- **必要なパッケージ**: `ca-certificates`, `tzdata`（Alpineパッケージ名）
- **docker-compose設定**:
  - `docker-compose.api.develop.yml`: `Dockerfile`を参照
  - `docker-compose.admin.develop.yml`: `Dockerfile.admin`を参照

### 2.2 システム構成図

```
┌─────────────────────────────────────────────────────────────┐
│                    Dockerfile構成（変更後）                   │
│                                                              │
│  ┌────────────────────────────────────────────────────┐   │
│  │ server/Dockerfile                                    │   │
│  │   - ビルドステージ: golang:1.24-alpine              │   │
│  │   - 実行ステージ: alpine:latest                     │   │
│  │   - CGO_ENABLED=0                                   │   │
│  │   - 用途: 全環境（develop/staging/production）      │   │
│  └────────────────────────────────────────────────────┘   │
│                          │                                  │
│                          ▼                                  │
│  ┌────────────────────────────────────────────────────┐   │
│  │ server/Dockerfile.admin                            │   │
│  │   - ビルドステージ: golang:1.24-alpine              │   │
│  │   - 実行ステージ: alpine:latest                     │   │
│  │   - CGO_ENABLED=0                                   │   │
│  │   - 用途: 全環境（develop/staging/production）      │   │
│  └────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│              docker-compose設定（変更後）                      │
│                                                              │
│  ┌────────────────────────────────────────────────────┐   │
│  │ docker-compose.api.develop.yml                      │   │
│  │   - dockerfile: Dockerfile                          │   │
│  └────────────────────────────────────────────────────┘   │
│                          │                                  │
│  ┌────────────────────────────────────────────────────┐   │
│  │ docker-compose.admin.develop.yml                    │   │
│  │   - dockerfile: Dockerfile.admin                    │   │
│  └────────────────────────────────────────────────────┘   │
│                          │                                  │
│  ┌────────────────────────────────────────────────────┐   │
│  │ docker-compose.api.staging.yml                      │   │
│  │   - dockerfile: Dockerfile（変更なし）              │   │
│  └────────────────────────────────────────────────────┘   │
│                          │                                  │
│  ┌────────────────────────────────────────────────────┐   │
│  │ docker-compose.api.production.yml                   │   │
│  │   - dockerfile: Dockerfile（変更なし）              │   │
│  └────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### 2.3 削除対象ファイル

```
server/
├── Dockerfile.develop          → 削除
└── Dockerfile.admin.develop    → 削除
```

## 3. コンポーネント設計

### 3.1 server/Dockerfile（Alpineベース）

#### 3.1.1 ビルドステージ
```dockerfile
# ビルドステージ
FROM golang:1.24-alpine AS builder

WORKDIR /build

# 依存関係をコピー
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# ビルド（CGO_ENABLED=0）
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server/main.go
```

**変更点**:
- `golang:1.24-bookworm` → `golang:1.24-alpine`
- CGO_ENABLED=0を維持（既存と同じ）

#### 3.1.2 実行ステージ
```dockerfile
# 実行ステージ
FROM alpine:latest

# 必要なパッケージをインストール
RUN apk add --no-cache ca-certificates tzdata

# 非rootユーザーを作成（Alpine用コマンド）
RUN addgroup -g 1000 appuser && \
    adduser -u 1000 -G appuser -D appuser

# ディレクトリ構造を作成
RUN mkdir -p /app/server/data /app/config /app/logs && \
    chown -R appuser:appuser /app

WORKDIR /app/server

# ビルド成果物をコピー
COPY --from=builder /build/server .

USER appuser

EXPOSE 8080

CMD ["./server"]
```

**変更点**:
- `debian:bookworm-slim` → `alpine:latest`
- `apt-get` → `apk add`
- `groupadd`/`useradd` → `addgroup`/`adduser`
- `--no-install-recommends` → `--no-cache`（Alpineのオプション）
- パッケージ名は同じ（`ca-certificates`, `tzdata`）

### 3.2 server/Dockerfile.admin（Alpineベース）

#### 3.2.1 ビルドステージ
```dockerfile
# ビルドステージ
FROM golang:1.24-alpine AS builder

WORKDIR /build

# 依存関係をコピー
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# ビルド（CGO_ENABLED=0）
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o admin ./cmd/admin/main.go
```

**変更点**:
- `golang:1.24-bookworm` → `golang:1.24-alpine`
- CGO_ENABLED=0を維持（既存と同じ）

#### 3.2.2 実行ステージ
```dockerfile
# 実行ステージ
FROM alpine:latest

# 必要なパッケージをインストール
RUN apk add --no-cache ca-certificates tzdata

# 非rootユーザーを作成（Alpine用コマンド）
RUN addgroup -g 1000 appuser && \
    adduser -u 1000 -G appuser -D appuser

# ディレクトリ構造を作成
RUN mkdir -p /app/server/data /app/config /app/logs && \
    chown -R appuser:appuser /app

WORKDIR /app/server

# ビルド成果物をコピー
COPY --from=builder /build/admin .

USER appuser

EXPOSE 8081

CMD ["./admin"]
```

**変更点**:
- `debian:bookworm-slim` → `alpine:latest`
- `apt-get` → `apk add`
- `groupadd`/`useradd` → `addgroup`/`adduser`
- `--no-install-recommends` → `--no-cache`（Alpineのオプション）
- パッケージ名は同じ（`ca-certificates`, `tzdata`）

### 3.3 docker-composeファイルの更新

#### 3.3.1 docker-compose.api.develop.yml
```yaml
services:
  api:
    build:
      context: ./server
      dockerfile: Dockerfile  # Dockerfile.develop → Dockerfile に変更
    # 他の設定は変更なし
```

#### 3.3.2 docker-compose.admin.develop.yml
```yaml
services:
  admin:
    build:
      context: ./server
      dockerfile: Dockerfile.admin  # Dockerfile.admin.develop → Dockerfile.admin に変更
    # 他の設定は変更なし
```

### 3.4 ドキュメント更新

#### 3.4.1 docs/Docker.mdの更新
- **Dockerfile構成表の更新**:
  - `Dockerfile.develop`と`Dockerfile.admin.develop`の行を削除
  - `Dockerfile`と`Dockerfile.admin`のベースイメージを`golang:1.24-alpine → alpine:latest`に更新
  - CGO設定をCGO_ENABLED=0に統一
  - 用途を「全環境（develop/staging/production）」に変更

## 4. データモデル設計

本実装ではデータモデルの変更はありません。

## 5. エラーハンドリング

### 5.1 ビルド時のエラー

#### 5.1.1 Alpineパッケージのインストールエラー
- **原因**: パッケージ名の誤りやネットワークエラー
- **対応**: パッケージ名を確認し、`apk update`を実行してから`apk add`を実行

#### 5.1.2 ビルドエラー
- **原因**: CGO_ENABLED=0でビルドできない依存関係がある
- **対応**: すべての依存関係がCGO不要であることを確認（既に確認済み）

### 5.2 実行時のエラー

#### 5.2.1 パッケージ不足エラー
- **原因**: 必要なパッケージ（`ca-certificates`, `tzdata`）がインストールされていない
- **対応**: Dockerfileでパッケージをインストールしていることを確認

#### 5.2.2 非rootユーザー作成エラー
- **原因**: Alpine用のコマンド（`addgroup`, `adduser`）の構文エラー
- **対応**: Alpine用のコマンド構文を確認

### 5.3 互換性エラー

#### 5.3.1 PostgreSQL接続エラー
- **原因**: Alpineベースに変更したことで、PostgreSQL接続に問題が発生
- **対応**: PostgreSQLドライバー（`gorm.io/driver/postgres`）はCGO不要で動作するため、問題ないことを確認

#### 5.3.2 Redis接続エラー
- **原因**: Alpineベースに変更したことで、Redis接続に問題が発生
- **対応**: Redis接続はネットワーク接続のみで、CGO不要であるため、問題ないことを確認

## 6. テスト戦略

### 6.1 ビルドテスト

#### 6.1.1 Dockerfileのビルドテスト
- **テスト内容**: `server/Dockerfile`と`server/Dockerfile.admin`が正常にビルドできることを確認
- **テスト方法**: 
  ```bash
  cd server
  docker build -f Dockerfile -t test-api .
  docker build -f Dockerfile.admin -t test-admin .
  ```
- **期待結果**: ビルドが正常に完了し、イメージが作成される

#### 6.1.2 docker-composeのビルドテスト
- **テスト内容**: `docker-compose.api.develop.yml`と`docker-compose.admin.develop.yml`が正常にビルドできることを確認
- **テスト方法**: 
  ```bash
  docker-compose -f docker-compose.api.develop.yml build
  docker-compose -f docker-compose.admin.develop.yml build
  ```
- **期待結果**: ビルドが正常に完了し、イメージが作成される

### 6.2 動作テスト

#### 6.2.1 コンテナ起動テスト
- **テスト内容**: ビルドしたイメージからコンテナが正常に起動できることを確認
- **テスト方法**: 
  ```bash
  docker-compose -f docker-compose.api.develop.yml up -d
  docker-compose -f docker-compose.admin.develop.yml up -d
  ```
- **期待結果**: コンテナが正常に起動し、エラーが発生しない

#### 6.2.2 PostgreSQL接続テスト
- **テスト内容**: APIサーバーとAdminサーバーがPostgreSQLに正常に接続できることを確認
- **テスト方法**: 
  - コンテナのログを確認
  - APIサーバー: ヘルスチェックエンドポイント（`/health`）にアクセス
  - Adminサーバー: ポートが開いているか確認（`nc -zv localhost 8081`または`curl -I http://localhost:8081`）
- **期待結果**: PostgreSQL接続が正常に確立され、エラーが発生しない

#### 6.2.3 Redis接続テスト
- **テスト内容**: APIサーバーがRedisに正常に接続できることを確認
- **テスト方法**: 
  - コンテナのログを確認
  - ヘルスチェックエンドポイントにアクセス
- **期待結果**: Redis接続が正常に確立され、エラーが発生しない

### 6.3 イメージサイズテスト

#### 6.3.1 イメージサイズ比較
- **テスト内容**: AlpineベースのイメージサイズがDebianベースより小さいことを確認
- **テスト方法**: 
  ```bash
  docker images | grep test-api
  docker images | grep test-admin
  ```
- **期待結果**: AlpineベースのイメージサイズがDebianベースより小さい

### 6.4 互換性テスト

#### 6.4.1 既存機能の動作確認
- **テスト内容**: Alpineベースに変更しても、既存の機能が正常に動作することを確認
- **テスト方法**: 
  - APIサーバーのエンドポイントにアクセス
  - Adminサーバーのエンドポイントにアクセス
- **期待結果**: 既存の機能が正常に動作する

## 7. 実装上の注意事項

### 7.1 Alpineベースへの移行

#### 7.1.1 パッケージマネージャーの違い
- **Debian**: `apt-get update && apt-get install -y --no-install-recommends <package> && rm -rf /var/lib/apt/lists/*`
- **Alpine**: `apk add --no-cache <package>`
- **注意点**: Alpineでは`--no-cache`オプションを使用することで、パッケージキャッシュを削除し、イメージサイズを削減

#### 7.1.2 非rootユーザー作成の違い
- **Debian**: `groupadd -g 1000 appuser && useradd -u 1000 -g appuser -m appuser`
- **Alpine**: `addgroup -g 1000 appuser && adduser -u 1000 -G appuser -D appuser`
- **注意点**: 
  - Alpineでは`addgroup`と`adduser`を使用
  - `-D`オプションでパスワードなしのユーザーを作成
  - `-G`オプションでグループを指定

#### 7.1.3 パッケージ名の確認
- **ca-certificates**: DebianとAlpineで同じパッケージ名
- **tzdata**: DebianとAlpineで同じパッケージ名
- **注意点**: パッケージ名が異なる場合は、適切なパッケージ名を確認

### 7.2 CGO_ENABLED=0の確認

#### 7.2.1 ビルド確認
- **確認内容**: CGO_ENABLED=0でビルドできることを確認
- **確認方法**: ビルドログを確認し、エラーが発生しないことを確認

#### 7.2.2 依存関係確認
- **確認内容**: すべての依存関係がCGO不要であることを確認
- **確認方法**: `go.mod`を確認し、CGO依存のライブラリが存在しないことを確認

### 7.3 ビルドと動作確認

#### 7.3.1 ビルド順序
1. `server/Dockerfile`をAlpineベースに書き換え
2. `server/Dockerfile.admin`をAlpineベースに書き換え
3. ビルドテストを実行
4. `docker-compose.api.develop.yml`を更新
5. `docker-compose.admin.develop.yml`を更新
6. docker-composeビルドテストを実行
7. 動作テストを実行
8. `server/Dockerfile.develop`を削除
9. `server/Dockerfile.admin.develop`を削除
10. `docs/Docker.md`を更新

#### 7.3.2 動作確認の順序
1. PostgreSQLコンテナを起動
2. Redisコンテナを起動
3. APIサーバーコンテナを起動
4. Adminサーバーコンテナを起動
5. ヘルスチェックを実行
6. ログを確認

### 7.4 ドキュメント整備

#### 7.4.1 Dockerfile構成表の更新
- **更新内容**: 
  - `Dockerfile.develop`と`Dockerfile.admin.develop`の行を削除
  - `Dockerfile`と`Dockerfile.admin`のベースイメージを`golang:1.24-alpine → alpine:latest`に更新
  - CGO設定をCGO_ENABLED=0に統一
  - 用途を「全環境（develop/staging/production）」に変更

#### 7.4.2 説明文の更新
- **更新内容**: 
  - Alpineベースに変更したことを明記
  - CGO_ENABLED=0に統一したことを明記
  - `Dockerfile.develop`と`Dockerfile.admin.develop`を削除したことを明記

## 8. 参考情報

### 8.1 関連Issue
- GitHub Issue #85: 開発環境はPostgreSQLを利用する前提とする
- GitHub Issue #90: server/Dockerfile* の更新

### 8.2 既存ドキュメント
- `docs/Docker.md`: Docker環境の説明
- `server/Dockerfile`: 既存のAPIサーバー用Dockerfile
- `server/Dockerfile.admin`: 既存のAdminサーバー用Dockerfile
- `docker-compose.api.develop.yml`: 開発環境用APIサーバーのdocker-compose設定
- `docker-compose.admin.develop.yml`: 開発環境用Adminサーバーのdocker-compose設定

### 8.3 技術スタック
- **Go**: 1.24
- **ベースイメージ**: `golang:1.24-alpine`（ビルドステージ）、`alpine:latest`（実行ステージ）
- **CGO**: CGO_ENABLED=0
- **PostgreSQLドライバー**: `gorm.io/driver/postgres`（CGO不要）

### 8.4 参考リンク
- Alpine Linux公式ドキュメント: https://alpinelinux.org/documentation/
- Docker公式ドキュメント: https://docs.docker.com/
- Go公式ドキュメント: https://go.dev/doc/
- Alpine Linuxパッケージ検索: https://pkgs.alpinelinux.org/packages
