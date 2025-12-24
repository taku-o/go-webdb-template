# Go DB Project Sample

Go + Next.js + Database Sharding対応のサンプルプロジェクトです。大規模プロジェクト向けの構成を採用しています。

## プロジェクト概要

- **サーバー**: Go言語、レイヤードアーキテクチャ、Database Sharding対応
- **クライアント**: Next.js 14 (App Router)、TypeScript
- **データベース**: SQLite (開発環境)、PostgreSQL/MySQL (本番想定)
- **テスト**: Go testing、Jest、Playwright

## 特徴

- ✅ **Sharding対応**: Hash-based shardingで複数DBへデータ分散
- ✅ **レイヤー分離**: API層、Service層、Repository層、DB層で責務を明確化
- ✅ **環境別設定**: develop/staging/production環境で設定切り替え
- ✅ **型安全**: TypeScriptによる型定義
- ✅ **テスト**: ユニット/統合/E2Eテスト対応

## セットアップ

### 前提条件

- Go 1.21+
- Node.js 18+
- SQLite3

### 1. 依存関係のインストール

#### サーバー側
```bash
cd server
go mod download
```

#### クライアント側
```bash
cd client
npm install
```

### 2. データベースのセットアップ

```bash
mkdir -p server/data
sqlite3 server/data/shard1.db < db/migrations/shard1/001_init.sql
sqlite3 server/data/shard2.db < db/migrations/shard2/001_init.sql
```

### 3. サーバー起動

```bash
cd server
APP_ENV=develop go run cmd/server/main.go
```

サーバーは http://localhost:8080 で起動します。

### 4. クライアント起動

```bash
cd client
npm run dev
```

クライアントは http://localhost:3000 で起動します。

## API エンドポイント

- `GET /api/users` - ユーザー一覧
- `POST /api/users` - ユーザー作成
- `GET /api/posts` - 投稿一覧
- `POST /api/posts` - 投稿作成
- `GET /api/user-posts` - ユーザーと投稿をJOIN（クロスシャードクエリ）

詳細は [プロジェクト構造計画](docs/plans/project-structure.md) を参照してください。

## Sharding戦略

Hash-based shardingを採用しています。

```
shard_id = hash(user_id) % shard_count + 1
```

- 同一ユーザーのデータは常に同じShardに配置
- クロスシャードクエリは各Shardから並列取得してマージ

## ライセンス

MIT License
