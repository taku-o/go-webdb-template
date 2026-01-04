# ビルドステージ
FROM golang:1.24-bookworm AS builder

WORKDIR /build

# 依存関係をコピー
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# ビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server/main.go

# 実行ステージ
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates tzdata && rm -rf /var/lib/apt/lists/*

# 非rootユーザーを作成
RUN groupadd -g 1000 appuser && \
    useradd -u 1000 -g appuser -m appuser

# ディレクトリ構造を作成（server/から起動される想定）
# /app/server/ - バイナリ配置・作業ディレクトリ
# /app/config/ - 設定ファイル（../config/でアクセス）
# /app/logs/   - ログファイル（../logs/でアクセス）
RUN mkdir -p /app/server/data /app/config /app/logs && \
    chown -R appuser:appuser /app

WORKDIR /app/server

# ビルド成果物をコピー
COPY --from=builder /build/server .

USER appuser

EXPOSE 8080

CMD ["./server"]
