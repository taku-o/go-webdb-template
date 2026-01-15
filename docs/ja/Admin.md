**[日本語]** | [English](../en/Admin.md)

# 管理画面ドキュメント

## 概要

GoAdminフレームワークを使用した管理画面です。メインサービス（ポート8080）とは独立したサービスとしてポート8081で動作します。

### 主な機能

- **テーブル管理**: Users/Postsテーブルの一覧表示・CRUD操作
- **シャーディング対応**: 全シャードのデータを統合表示
- **認証・認可**: GoAdmin組み込み認証機能
- **カスタムページ**: ダッシュボード、ユーザー登録フォーム

## 起動方法

### 前提条件

- データベースがセットアップ済みであること
- GoAdmin用のマイグレーションが適用済みであること

```bash
# PostgreSQLコンテナを起動
./scripts/start-postgres.sh start

# マイグレーションを適用（初回のみ）
./scripts/migrate.sh all
```

### 管理画面の起動

```bash
cd server
APP_ENV=develop go run cmd/admin/main.go
```

管理画面は http://localhost:8081/admin でアクセスできます。

## ログイン方法

### 開発環境のデフォルト認証情報

| 項目 | 値 |
|------|-----|
| URL | http://localhost:8081/admin/login |
| ユーザー名 | `admin` |
| パスワード | `admin123` |

### 認証情報の変更

認証情報は設定ファイル（`config/develop.yaml`）で変更できます。

```yaml
admin:
  port: 8081
  auth:
    username: admin
    password: your_password_here
  session:
    lifetime: 7200  # セッション有効期限（秒）
```

### 本番環境での設定

本番環境では環境変数または別の設定ファイルを使用してください。

```yaml
# config/production.yaml
admin:
  port: 8081
  auth:
    username: ${ADMIN_USERNAME}
    password: ${ADMIN_PASSWORD}
  session:
    lifetime: 3600
```

## テーブル管理

### ユーザー管理（Users）

**URL**: http://localhost:8081/admin/info/users

#### 一覧表示

| カラム | 説明 | 操作 |
|--------|------|------|
| ID | ユーザーID | ソート可能 |
| 名前 | ユーザー名 | ソート・フィルタ可能 |
| メールアドレス | メールアドレス | ソート・フィルタ可能 |
| 作成日時 | 登録日時 | ソート可能 |
| 更新日時 | 最終更新日時 | ソート可能 |

#### 操作

- **新規作成**: 「新規」ボタンをクリック
- **編集**: 一覧の「編集」アイコンをクリック
- **削除**: 一覧の「削除」アイコンをクリック
- **エクスポート**: 「エクスポート」ボタンでCSV出力

### 投稿管理（Posts）

**URL**: http://localhost:8081/admin/info/posts

#### 一覧表示

| カラム | 説明 | 操作 |
|--------|------|------|
| ID | 投稿ID | ソート可能 |
| ユーザーID | 投稿者のユーザーID | ソート・フィルタ可能 |
| タイトル | 投稿タイトル | ソート・フィルタ可能 |
| 内容 | 投稿本文 | - |
| 作成日時 | 投稿日時 | ソート可能 |
| 更新日時 | 最終更新日時 | ソート可能 |

#### 操作

- **新規作成**: 「新規」ボタンをクリック
- **編集**: 一覧の「編集」アイコンをクリック
- **削除**: 一覧の「削除」アイコンをクリック
- **エクスポート**: 「エクスポート」ボタンでCSV出力

## カスタムページ

### ダッシュボード

**URL**: http://localhost:8081/admin/

ログイン後に表示されるトップページです。

#### 表示内容

- **統計情報**: ユーザー数、投稿数
- **クイックアクション**: ユーザー登録、投稿作成へのリンク
- **システム情報**: プロジェクト名、GoAdminバージョン

### ユーザー登録ページ

**URL**: http://localhost:8081/admin/user/register

カスタムフォームによるユーザー登録ページです。

#### 入力項目

| 項目 | 必須 | 説明 |
|------|------|------|
| 名前 | ○ | 100文字以内 |
| メールアドレス | ○ | 255文字以内、有効な形式 |

#### バリデーション

- 名前は必須、100文字以内
- メールアドレスは必須、有効な形式、255文字以内
- メールアドレスの重複チェック

#### 登録完了ページ

登録成功後、登録したユーザー情報を表示する完了ページにリダイレクトされます。

## メニュー構造

```
管理画面
├── ダッシュボード（/admin/）
├── ユーザー
│   ├── 一覧（/admin/info/users）
│   └── 登録（/admin/user/register）
└── 投稿
    └── 一覧（/admin/info/posts）
```

## トラブルシューティング

### ログインできない

1. **認証情報の確認**: 設定ファイル（`config/develop.yaml`）のusername/passwordを確認
2. **データベースの確認**: GoAdmin用マイグレーション（002_goadmin.sql）が適用されているか確認
3. **ブラウザキャッシュ**: Cookieをクリアしてリトライ

```bash
# GoAdmin用テーブルの確認（PostgreSQL）
psql -h localhost -p 5432 -U webdb -d webdb_master -c "SELECT * FROM goadmin_users;"
```

### 管理画面が起動しない

1. **ポート競合**: ポート8081が使用中でないか確認

```bash
lsof -i :8081
```

2. **設定ファイルの確認**: `config/develop.yaml`のフォーマットエラーがないか確認

3. **依存関係の確認**: 必要なパッケージがインストールされているか確認

```bash
cd server
go mod tidy
```

### データが表示されない

1. **データベース接続の確認**: サーバーログでDB接続エラーがないか確認
2. **テーブルの確認**: dm_usersテーブル、dm_postsテーブルにデータがあるか確認

```bash
# PostgreSQLで確認
psql -h localhost -p 5433 -U webdb -d webdb_sharding_1 -c "SELECT COUNT(*) FROM dm_users_000;"
psql -h localhost -p 5433 -U webdb -d webdb_sharding_1 -c "SELECT COUNT(*) FROM dm_posts_000;"
```

### セッションがすぐ切れる

設定ファイルの`session.lifetime`を確認してください。値は秒単位です。

```yaml
admin:
  session:
    lifetime: 7200  # 2時間
```

## 技術仕様

### 使用ライブラリ

- GoAdmin v1.2.26
- AdminLTE テーマ
- Gorilla Mux（HTTPルーター）

### アーキテクチャ

```
┌─────────────────────────────────────────────────────────────┐
│              Admin Service (Port 8081)                       │
│                                                               │
│  ┌────────────────────────────────────────────────────────┐ │
│  │                  GoAdmin Engine                         │ │
│  │  • Admin Plugin (CRUD自動生成)                          │ │
│  │  • Custom Pages (カスタムページ)                         │ │
│  │  • Authentication (認証・認可)                          │ │
│  └──────────────────────┬─────────────────────────────────┘ │
│                         │                                     │
│  ┌──────────────────────▼─────────────────────────────────┐ │
│  │                 GORM Manager                            │ │
│  │  (既存の接続管理を再利用)                                 │ │
│  └──────────────────────┬─────────────────────────────────┘ │
└─────────────────────────┼───────────────────────────────────┘
                          │
         ┌────────────────┴────────────────┐
         ▼                                  ▼
    ┌─────────┐                        ┌─────────┐
    │ Shard 1 │                        │ Shard 2 │
    └─────────┘                        └─────────┘
```

### ファイル構成

```
server/
├── cmd/
│   └── admin/
│       └── main.go           # エントリーポイント
└── internal/
    └── admin/
        ├── config.go         # GoAdmin設定
        ├── tables.go         # テーブルジェネレータ
        ├── sharding.go       # クロスシャードクエリ
        ├── auth/
        │   ├── auth.go       # 認証ロジック
        │   └── session.go    # セッション管理
        └── pages/
            ├── pages.go               # カスタムページ基盤
            ├── home.go                # ダッシュボード
            ├── user_register.go       # ユーザー登録
            └── user_register_complete.go  # 登録完了
```
