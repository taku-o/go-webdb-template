# server/Dockerfile* の更新実装タスク一覧

## 概要
SQLite依存を削除してDockerfileを統合し、Alpineベースのイメージに書き換えるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: DockerfileのAlpineベースへの書き換え

#### タスク 1.1: server/DockerfileのビルドステージをAlpineベースに書き換え
**目的**: `server/Dockerfile`のビルドステージをDebianベースからAlpineベースに変更。

**作業内容**:
- `server/Dockerfile`を開く
- ビルドステージのベースイメージを変更:
  - `FROM golang:1.24-bookworm AS builder` → `FROM golang:1.24-alpine AS builder`
- CGO_ENABLED=0を維持（既存の設定を維持）
- ビルドコマンドを確認:
  - `RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server/main.go`
  - 変更不要（CGO_ENABLED=0を維持）

**受け入れ基準**:
- ビルドステージのベースイメージが`golang:1.24-alpine`に変更されている
- CGO_ENABLED=0が維持されている
- ビルドコマンドが正しく設定されている

- _Requirements: 3.2.1_
- _Design: 3.1.1_

---

#### タスク 1.2: server/Dockerfileの実行ステージをAlpineベースに書き換え
**目的**: `server/Dockerfile`の実行ステージをDebianベースからAlpineベースに変更。

**作業内容**:
- `server/Dockerfile`を開く
- 実行ステージのベースイメージを変更:
  - `FROM debian:bookworm-slim` → `FROM alpine:latest`
- パッケージインストールコマンドを変更:
  - `RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates tzdata && rm -rf /var/lib/apt/lists/*`
  - → `RUN apk add --no-cache ca-certificates tzdata`
- 非rootユーザー作成コマンドを変更:
  - `RUN groupadd -g 1000 appuser && useradd -u 1000 -g appuser -m appuser`
  - → `RUN addgroup -g 1000 appuser && adduser -u 1000 -G appuser -D appuser`
- ディレクトリ構造の作成コマンドを確認（変更不要）:
  - `RUN mkdir -p /app/server/data /app/config /app/logs && chown -R appuser:appuser /app`

**修正後のコード例**:
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

**受け入れ基準**:
- 実行ステージのベースイメージが`alpine:latest`に変更されている
- パッケージインストールコマンドが`apk add --no-cache`に変更されている
- 非rootユーザー作成コマンドが`addgroup`/`adduser`に変更されている
- ディレクトリ構造の作成コマンドが正しく設定されている

- _Requirements: 3.2.1_
- _Design: 3.1.2_

---

#### タスク 1.3: server/Dockerfile.adminのビルドステージをAlpineベースに書き換え
**目的**: `server/Dockerfile.admin`のビルドステージをDebianベースからAlpineベースに変更。

**作業内容**:
- `server/Dockerfile.admin`を開く
- ビルドステージのベースイメージを変更:
  - `FROM golang:1.24-bookworm AS builder` → `FROM golang:1.24-alpine AS builder`
- CGO_ENABLED=0を維持（既存の設定を維持）
- ビルドコマンドを確認:
  - `RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o admin ./cmd/admin/main.go`
  - 変更不要（CGO_ENABLED=0を維持）

**受け入れ基準**:
- ビルドステージのベースイメージが`golang:1.24-alpine`に変更されている
- CGO_ENABLED=0が維持されている
- ビルドコマンドが正しく設定されている

- _Requirements: 3.2.2_
- _Design: 3.2.1_

---

#### タスク 1.4: server/Dockerfile.adminの実行ステージをAlpineベースに書き換え
**目的**: `server/Dockerfile.admin`の実行ステージをDebianベースからAlpineベースに変更。

**作業内容**:
- `server/Dockerfile.admin`を開く
- 実行ステージのベースイメージを変更:
  - `FROM debian:bookworm-slim` → `FROM alpine:latest`
- パッケージインストールコマンドを変更:
  - `RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates tzdata && rm -rf /var/lib/apt/lists/*`
  - → `RUN apk add --no-cache ca-certificates tzdata`
- 非rootユーザー作成コマンドを変更:
  - `RUN groupadd -g 1000 appuser && useradd -u 1000 -g appuser -m appuser`
  - → `RUN addgroup -g 1000 appuser && adduser -u 1000 -G appuser -D appuser`
- ディレクトリ構造の作成コマンドを確認（変更不要）:
  - `RUN mkdir -p /app/server/data /app/config /app/logs && chown -R appuser:appuser /app`

**修正後のコード例**:
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

**受け入れ基準**:
- 実行ステージのベースイメージが`alpine:latest`に変更されている
- パッケージインストールコマンドが`apk add --no-cache`に変更されている
- 非rootユーザー作成コマンドが`addgroup`/`adduser`に変更されている
- ディレクトリ構造の作成コマンドが正しく設定されている

- _Requirements: 3.2.2_
- _Design: 3.2.2_

---

### Phase 2: docker-composeファイルの更新

#### タスク 2.1: docker-compose.api.develop.ymlの更新
**目的**: `docker-compose.api.develop.yml`で`Dockerfile.develop`参照を`Dockerfile`に変更。

**作業内容**:
- `docker-compose.api.develop.yml`を開く
- `build.dockerfile`を変更:
  - `dockerfile: Dockerfile.develop` → `dockerfile: Dockerfile`
- 他の設定（ports、environment、volumes、networks等）は変更しないことを確認

**修正後のコード例**:
```yaml
services:
  api:
    build:
      context: ./server
      dockerfile: Dockerfile  # Dockerfile.develop → Dockerfile に変更
    # 他の設定は変更なし
```

**受け入れ基準**:
- `dockerfile: Dockerfile.develop`が`dockerfile: Dockerfile`に変更されている
- 他の設定（ports、environment、volumes、networks等）が変更されていない

- _Requirements: 3.3.1_
- _Design: 3.3.1_

---

#### タスク 2.2: docker-compose.admin.develop.ymlの更新
**目的**: `docker-compose.admin.develop.yml`で`Dockerfile.admin.develop`参照を`Dockerfile.admin`に変更。

**作業内容**:
- `docker-compose.admin.develop.yml`を開く
- `build.dockerfile`を変更:
  - `dockerfile: Dockerfile.admin.develop` → `dockerfile: Dockerfile.admin`
- 他の設定（ports、environment、volumes、networks等）は変更しないことを確認

**修正後のコード例**:
```yaml
services:
  admin:
    build:
      context: ./server
      dockerfile: Dockerfile.admin  # Dockerfile.admin.develop → Dockerfile.admin に変更
    # 他の設定は変更なし
```

**受け入れ基準**:
- `dockerfile: Dockerfile.admin.develop`が`dockerfile: Dockerfile.admin`に変更されている
- 他の設定（ports、environment、volumes、networks等）が変更されていない

- _Requirements: 3.3.2_
- _Design: 3.3.2_

---

### Phase 3: ビルドテスト

#### タスク 3.1: server/Dockerfileのビルドテスト
**目的**: `server/Dockerfile`がAlpineベースで正常にビルドできることを確認。

**作業内容**:
- `server`ディレクトリに移動
- Dockerfileをビルド:
  ```bash
  cd server
  docker build -f Dockerfile -t test-api .
  ```
- ビルドログを確認:
  - エラーが発生していないことを確認
  - ビルドが正常に完了していることを確認
- イメージが作成されていることを確認:
  ```bash
  docker images | grep test-api
  ```

**受け入れ基準**:
- ビルドが正常に完了し、エラーが発生していない
- イメージが作成されている
- ビルドログにエラーが表示されていない

- _Requirements: 6.3_
- _Design: 6.1.1_

---

#### タスク 3.2: server/Dockerfile.adminのビルドテスト
**目的**: `server/Dockerfile.admin`がAlpineベースで正常にビルドできることを確認。

**作業内容**:
- `server`ディレクトリに移動
- Dockerfile.adminをビルド:
  ```bash
  cd server
  docker build -f Dockerfile.admin -t test-admin .
  ```
- ビルドログを確認:
  - エラーが発生していないことを確認
  - ビルドが正常に完了していることを確認
- イメージが作成されていることを確認:
  ```bash
  docker images | grep test-admin
  ```

**受け入れ基準**:
- ビルドが正常に完了し、エラーが発生していない
- イメージが作成されている
- ビルドログにエラーが表示されていない

- _Requirements: 6.3_
- _Design: 6.1.1_

---

#### タスク 3.3: docker-compose.api.develop.ymlのビルドテスト
**目的**: `docker-compose.api.develop.yml`が正常にビルドできることを確認。

**作業内容**:
- プロジェクトルートディレクトリに移動
- docker-composeでビルド:
  ```bash
  docker-compose -f docker-compose.api.develop.yml build
  ```
- ビルドログを確認:
  - エラーが発生していないことを確認
  - ビルドが正常に完了していることを確認
  - `Dockerfile`が使用されていることを確認

**受け入れ基準**:
- ビルドが正常に完了し、エラーが発生していない
- `Dockerfile`が使用されている
- ビルドログにエラーが表示されていない

- _Requirements: 6.3_
- _Design: 6.1.2_

---

#### タスク 3.4: docker-compose.admin.develop.ymlのビルドテスト
**目的**: `docker-compose.admin.develop.yml`が正常にビルドできることを確認。

**作業内容**:
- プロジェクトルートディレクトリに移動
- docker-composeでビルド:
  ```bash
  docker-compose -f docker-compose.admin.develop.yml build
  ```
- ビルドログを確認:
  - エラーが発生していないことを確認
  - ビルドが正常に完了していることを確認
  - `Dockerfile.admin`が使用されていることを確認

**受け入れ基準**:
- ビルドが正常に完了し、エラーが発生していない
- `Dockerfile.admin`が使用されている
- ビルドログにエラーが表示されていない

- _Requirements: 6.3_
- _Design: 6.1.2_

---

### Phase 4: 動作テスト

#### タスク 4.1: PostgreSQLコンテナの起動確認
**目的**: PostgreSQLコンテナが正常に起動していることを確認。

**作業内容**:
- PostgreSQLコンテナを起動:
  ```bash
  docker-compose -f docker-compose.postgres.yml up -d
  ```
- コンテナの状態を確認:
  ```bash
  docker-compose -f docker-compose.postgres.yml ps
  ```
- コンテナのログを確認:
  ```bash
  docker-compose -f docker-compose.postgres.yml logs
  ```
- エラーが発生していないことを確認

**受け入れ基準**:
- PostgreSQLコンテナが正常に起動している
- コンテナのログにエラーが表示されていない
- すべてのPostgreSQLコンテナ（master 1台、sharding 4台）が起動している

- _Requirements: 6.3_
- _Design: 6.2.2_

---

#### タスク 4.2: Redisコンテナの起動確認
**目的**: Redisコンテナが正常に起動していることを確認。

**作業内容**:
- Redisコンテナを起動:
  ```bash
  docker-compose -f docker-compose.redis.yml up -d
  ```
- コンテナの状態を確認:
  ```bash
  docker-compose -f docker-compose.redis.yml ps
  ```
- コンテナのログを確認:
  ```bash
  docker-compose -f docker-compose.redis.yml logs
  ```
- エラーが発生していないことを確認

**受け入れ基準**:
- Redisコンテナが正常に起動している
- コンテナのログにエラーが表示されていない

- _Requirements: 6.3_
- _Design: 6.2.3_

---

#### タスク 4.3: APIサーバーコンテナの起動と動作確認
**目的**: APIサーバーコンテナが正常に起動し、PostgreSQLとRedisに接続できることを確認。

**作業内容**:
- APIサーバーコンテナを起動:
  ```bash
  docker-compose -f docker-compose.api.develop.yml up -d
  ```
- コンテナの状態を確認:
  ```bash
  docker-compose -f docker-compose.api.develop.yml ps
  ```
- コンテナのログを確認:
  ```bash
  docker-compose -f docker-compose.api.develop.yml logs
  ```
- エラーが発生していないことを確認:
  - PostgreSQL接続エラーが発生していないことを確認
  - Redis接続エラーが発生していないことを確認
- ヘルスチェックを実行:
  ```bash
  curl http://localhost:8080/health
  ```
- ヘルスチェックが正常に応答することを確認

**受け入れ基準**:
- APIサーバーコンテナが正常に起動している
- PostgreSQL接続エラーが発生していない
- Redis接続エラーが発生していない
- ヘルスチェックが正常に応答する（HTTP 200 OK）

- _Requirements: 6.3_
- _Design: 6.2.1, 6.2.2, 6.2.3_

---

#### タスク 4.4: Adminサーバーコンテナの起動と動作確認
**目的**: Adminサーバーコンテナが正常に起動し、PostgreSQLに接続できることを確認。

**作業内容**:
- Adminサーバーコンテナを起動:
  ```bash
  docker-compose -f docker-compose.admin.develop.yml up -d
  ```
- コンテナの状態を確認:
  ```bash
  docker-compose -f docker-compose.admin.develop.yml ps
  ```
- コンテナのログを確認:
  ```bash
  docker-compose -f docker-compose.admin.develop.yml logs
  ```
- エラーが発生していないことを確認:
  - PostgreSQL接続エラーが発生していないことを確認
  - サーバー起動エラーが発生していないことを確認
- ポートが開いていることを確認:
  ```bash
  nc -zv localhost 8081
  ```
  または
  ```bash
  curl -I http://localhost:8081
  ```
- ポートが開いていることを確認（接続可能であることを確認）

**受け入れ基準**:
- Adminサーバーコンテナが正常に起動している
- PostgreSQL接続エラーが発生していない
- サーバー起動エラーが発生していない
- ポート8081が開いている（接続可能である）

- _Requirements: 6.3_
- _Design: 6.2.1, 6.2.2_

---

#### タスク 4.5: イメージサイズの確認
**目的**: AlpineベースのイメージサイズがDebianベースより小さいことを確認（オプション）。

**作業内容**:
- イメージサイズを確認:
  ```bash
  docker images | grep -E "(test-api|test-admin|go-webdb-template-api|go-webdb-template-admin)"
  ```
- イメージサイズを比較:
  - Alpineベースのイメージサイズを記録
  - Debianベースのイメージサイズと比較（既存のイメージがある場合）
- イメージサイズが削減されていることを確認

**受け入れ基準**:
- イメージサイズが確認できる
- AlpineベースのイメージサイズがDebianベースより小さい（比較可能な場合）

- _Requirements: 6.4_
- _Design: 6.3.1_

---

### Phase 5: 不要ファイルの削除

#### タスク 5.1: server/Dockerfile.developの削除
**目的**: 不要になった`server/Dockerfile.develop`を削除。

**作業内容**:
- `server/Dockerfile.develop`が存在することを確認
- ファイルを削除:
  ```bash
  rm server/Dockerfile.develop
  ```
- ファイルが削除されたことを確認:
  ```bash
  ls server/Dockerfile.develop
  ```
- エラー（ファイルが存在しない）が表示されることを確認

**受け入れ基準**:
- `server/Dockerfile.develop`が削除されている
- ファイルが存在しないことが確認できる

- _Requirements: 3.1.1, 6.1_
- _Design: 2.3_

---

#### タスク 5.2: server/Dockerfile.admin.developの削除
**目的**: 不要になった`server/Dockerfile.admin.develop`を削除。

**作業内容**:
- `server/Dockerfile.admin.develop`が存在することを確認
- ファイルを削除:
  ```bash
  rm server/Dockerfile.admin.develop
  ```
- ファイルが削除されたことを確認:
  ```bash
  ls server/Dockerfile.admin.develop
  ```
- エラー（ファイルが存在しない）が表示されることを確認

**受け入れ基準**:
- `server/Dockerfile.admin.develop`が削除されている
- ファイルが存在しないことが確認できる

- _Requirements: 3.1.2, 6.1_
- _Design: 2.3_

---

### Phase 6: ドキュメントの更新

#### タスク 6.1: docs/Docker.mdのDockerfile構成表の更新
**目的**: `docs/Docker.md`のDockerfile構成表を更新し、Alpineベースへの変更と削除ファイルを反映。

**作業内容**:
- `docs/Docker.md`を開く
- Dockerfile構成表を更新:
  - `Dockerfile.develop`の行を削除
  - `Dockerfile.admin.develop`の行を削除
  - `Dockerfile`のベースイメージを`golang:1.24-bookworm → debian:bookworm-slim`から`golang:1.24-alpine → alpine:latest`に更新
  - `Dockerfile.admin`のベースイメージを`golang:1.24-bookworm → debian:bookworm-slim`から`golang:1.24-alpine → alpine:latest`に更新
  - CGO設定をCGO_ENABLED=0に統一
  - 用途を「全環境（develop/staging/production）」に変更

**更新後の表例**:
```markdown
| ファイル | 用途 | CGO | ベースイメージ |
|---------|------|-----|---------------|
| `server/Dockerfile` | 全環境 | 0 | golang:1.24-alpine → alpine:latest |
| `server/Dockerfile.admin` | 全環境 | 0 | golang:1.24-alpine → alpine:latest |
```

**受け入れ基準**:
- `Dockerfile.develop`と`Dockerfile.admin.develop`の行が削除されている
- `Dockerfile`と`Dockerfile.admin`のベースイメージが`golang:1.24-alpine → alpine:latest`に更新されている
- CGO設定がCGO_ENABLED=0に統一されている
- 用途が「全環境（develop/staging/production）」に変更されている

- _Requirements: 3.4.1, 6.5_
- _Design: 3.4.1_

---

#### タスク 6.2: docs/Docker.mdの説明文の更新
**目的**: `docs/Docker.md`の説明文を更新し、Alpineベースへの変更と削除ファイルを明記。

**作業内容**:
- `docs/Docker.md`を開く
- 説明文を更新:
  - Alpineベースに変更したことを明記
  - CGO_ENABLED=0に統一したことを明記
  - `Dockerfile.develop`と`Dockerfile.admin.develop`を削除したことを明記
  - develop環境とstaging/production環境で同じDockerfileを使用することを明記

**受け入れ基準**:
- Alpineベースに変更したことが明記されている
- CGO_ENABLED=0に統一したことが明記されている
- `Dockerfile.develop`と`Dockerfile.admin.develop`を削除したことが明記されている
- develop環境とstaging/production環境で同じDockerfileを使用することが明記されている

- _Requirements: 3.4.1, 6.5_
- _Design: 3.4.1, 7.4.2_

---

### Phase 7: 最終確認

#### タスク 7.1: 全タスクの受け入れ基準確認
**目的**: すべてのタスクの受け入れ基準が満たされていることを確認。

**作業内容**:
- 要件定義書の受け入れ基準（6.1〜6.5）を確認:
  - 6.1: Dockerfileの統合と削除
  - 6.2: DockerfileのAlpineベースへの書き換え
  - 6.3: ビルドと動作確認
  - 6.4: イメージサイズとビルド時間
  - 6.5: ドキュメント
- 各受け入れ基準が満たされていることを確認
- 未達成の受け入れ基準がある場合は、対応タスクを再実行

**受け入れ基準**:
- すべての受け入れ基準が満たされている
- 未達成の受け入れ基準がない

- _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5_

---

## タスク実行順序

1. **Phase 1**: DockerfileのAlpineベースへの書き換え（タスク 1.1〜1.4）
2. **Phase 2**: docker-composeファイルの更新（タスク 2.1〜2.2）
3. **Phase 3**: ビルドテスト（タスク 3.1〜3.4）
4. **Phase 4**: 動作テスト（タスク 4.1〜4.5）
5. **Phase 5**: 不要ファイルの削除（タスク 5.1〜5.2）
6. **Phase 6**: ドキュメントの更新（タスク 6.1〜6.2）
7. **Phase 7**: 最終確認（タスク 7.1）

## 注意事項

- タスク 3.1〜3.4のビルドテストは、タスク 1.1〜1.4とタスク 2.1〜2.2の完了後に実行する
- タスク 4.1〜4.4の動作テストは、タスク 3.1〜3.4のビルドテストの完了後に実行する
- タスク 5.1〜5.2の不要ファイルの削除は、タスク 4.1〜4.4の動作テストの完了後に実行する
- タスク 6.1〜6.2のドキュメントの更新は、タスク 5.1〜5.2の不要ファイルの削除の完了後に実行する
- タスク 7.1の最終確認は、すべてのタスクの完了後に実行する
