---
layout: default
title: セットアップ手順
lang: ja
---

# セットアップ手順

クライアントサーバーを動作させるまでの詳細なセットアップ手順を説明します。

---

## 1. 初期セットアップ

### パッケージアプリケーションのインストール

- **Docker**: [https://www.docker.com/ja-jp/](https://www.docker.com/ja-jp/)
- **Cursor**: [https://cursor.com/ja](https://cursor.com/ja)
- **Go**: [https://go.dev/dl/](https://go.dev/dl/)

### Homebrewのインストール

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
eval "$(/opt/homebrew/bin/brew shellenv)"
```

#### GitHub CLIのインストール

```bash
brew install gh
gh auth login
gh auth status
```

#### Atlasのインストール

```bash
brew install ariga/tap/atlas
```

### Node.js（nvm）のインストール

```bash
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.4/install.sh | bash
nvm ls-remote
nvm install v22.14.0
nvm use v22.14.0
nvm alias default v22.14.0
```

`.bashrc`に以下を追加：

```bash
if [ -f ~/.nvm/nvm.sh ]
then
  source ~/.nvm/nvm.sh
fi
```

#### Claude Codeのインストール

```bash
npm install -g @anthropic-ai/claude-code
```

### uvのインストール

```bash
brew install uv
```

#### Serenaの設定

プロジェクトディレクトリで以下を実行：

```bash
claude mcp add serena -- uvx --from git+https://github.com/oraios/serena serena-mcp-server --context ide-assistant --enable-web-dashboard false --project $(pwd)
```

Serenaインデックスの更新：

```bash
uvx --from git+https://github.com/oraios/serena index-project
```

---

## 2. 依存関係のインストール

### サーバー側

```bash
cd server
go mod download
```

### クライアント側

```bash
cd client
npm install --legacy-peer-deps
```

---

## 3. データベースのセットアップ

### PostgreSQLの起動

```bash
./scripts/start-postgres.sh start
```

**接続情報**

| データベース | ホスト | ポート | ユーザー | パスワード | データベース名 |
|------------|--------|--------|---------|-----------|--------------|
| Master | localhost | 5432 | webdb | webdb | webdb_master |
| Sharding 1 | localhost | 5433 | webdb | webdb | webdb_sharding_1 |
| Sharding 2 | localhost | 5434 | webdb | webdb | webdb_sharding_2 |
| Sharding 3 | localhost | 5435 | webdb | webdb | webdb_sharding_3 |
| Sharding 4 | localhost | 5436 | webdb | webdb | webdb_sharding_4 |

### マイグレーションの適用

```bash
./scripts/migrate.sh all
```

---

## 4. Redisの起動

```bash
# Redisを起動
./scripts/start-redis.sh start
```

- Redis: http://localhost:6379

---

## 5. Auth0アカウントの設定

次の資料を参考にAuth0アカウントをセットアップします。
- [Auth0 外部ID連携 導入・開発ガイド](https://github.com/taku-o/go-webdb-template/blob/master/docs/Partner-Idp-Auth0-Login.md)

---

## 6. クライアント環境変数の設定

### .env.localの作成

`client/.env.develop`を`client/.env.local`にリネームして以下の環境変数を設定：

```
# Auth0設定
AUTH0_ISSUER=https://your-tenant.auth0.com
AUTH0_CLIENT_ID=your-client-id
AUTH0_CLIENT_SECRET=your-client-secret
```

---

## 7. サーバーの起動

### APIサーバーの起動

```bash
cd server
APP_ENV=develop go run cmd/server/main.go
```

### Adminサーバーの起動

```bash
cd server
APP_ENV=develop go run cmd/admin/main.go
```

### クライアントサーバーの起動

```bash
cd client
npm run dev
```

---

## 8. 各種URL情報

| サービス | URL | 備考 |
|---------|-----|------|
| クライアント | http://localhost:3000 | Next.jsアプリケーション |
| APIサーバー doc | http://localhost:8080/docs | API Documentation UI |
| Adminサーバー | http://localhost:8081/admin | 管理画面 |

### Adminサーバー認証情報

| 項目 | 値 |
|------|-----|
| ID | admin |
| Password | admin123 |

---

## ナビゲーション

- [ホーム]({{ site.baseurl }}/pages/ja/)
- [プロジェクト概要]({{ site.baseurl }}/pages/ja/about)
- [English]({{ site.baseurl }}/pages/en/setup)
