# 開発コマンド一覧

## サーバー起動

### メインサーバー（ポート8080）
```bash
cd server
APP_ENV=develop go run cmd/server/main.go
```

### 管理画面サーバー（ポート8081）
```bash
cd server
APP_ENV=develop go run cmd/admin/main.go
```

## テスト

### サーバー側テスト（全体）
```bash
cd server
go test ./... -v
```

### サーバー側テスト（特定パッケージ）
```bash
cd server
go test ./internal/... -v
```

### テストカバレッジ
```bash
cd server
go test ./internal/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### クライアント側テスト
```bash
cd client
npm test
```

### E2Eテスト
```bash
cd client
npx playwright test
```

## ビルド

### サーバービルド
```bash
cd server
go build ./cmd/server/...
go build ./cmd/admin/...
```

## 依存関係

### サーバー依存インストール
```bash
cd server
go mod download
go mod tidy
```

### クライアント依存インストール
```bash
cd client
npm install
```

## データベース

### マイグレーション適用（SQLite）
```bash
mkdir -p server/data
sqlite3 server/data/shard1.db < db/migrations/shard1/001_init.sql
sqlite3 server/data/shard1.db < db/migrations/shard1/002_goadmin.sql
sqlite3 server/data/shard2.db < db/migrations/shard2/001_init.sql
```

## Git関連
```bash
git status
git diff
git log --oneline -10
git branch -a
```
