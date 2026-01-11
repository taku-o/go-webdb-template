# CLIツールドキュメント

## 概要

バッチ処理やcronジョブとして実行するためのCLIツールです。既存のサービス層・リポジトリ層を再利用し、非対話的な実行をサポートします。

### 主な機能

- **ユーザー一覧出力**: 全シャードからユーザー一覧を取得しTSV形式で出力
- **件数制限**: `--limit`フラグで出力件数を制御
- **cron対応**: 非対話的実行、適切な終了コード

## ディレクトリ構造

```
server/
├── cmd/
│   ├── list-users/
│   │   ├── main.go          # CLIツール本体
│   │   └── main_test.go     # ユニットテスト
│   └── generate-sample-data/
│       └── main.go          # サンプルデータ生成ツール
└── bin/                      # ビルド後の実行ファイル（.gitignore対象）
    ├── list-users
    └── generate-sample-data
```

## ビルド方法

### 開発環境

```bash
cd server
go build -o bin/list-users ./cmd/list-users
```

### 本番環境（クロスコンパイル）

```bash
cd server
GOOS=linux GOARCH=amd64 go build -o bin/list-users ./cmd/list-users
```

### リリースビルド（最適化）

```bash
cd server
go build -ldflags="-s -w" -o bin/list-users ./cmd/list-users
```

## list-users コマンド

### 概要

全シャードからユーザー一覧を取得し、TSV形式で標準出力に出力します。

### 使用方法

```bash
APP_ENV=<環境名> ./bin/list-users [オプション]
```

### オプション

| オプション | 説明 | デフォルト | 有効範囲 |
|-----------|------|----------|---------|
| `--limit` | 出力件数 | 20 | 1〜100 |

### 実行例

```bash
# デフォルト（20件）
APP_ENV=develop ./bin/list-users

# 件数を指定
APP_ENV=develop ./bin/list-users --limit 50

# 最大件数（100件）
APP_ENV=develop ./bin/list-users --limit 100

# ファイルに出力
APP_ENV=develop ./bin/list-users --limit 100 > users.tsv
```

### 出力形式

TSV（タブ区切り）形式で出力されます。

```
ID	Name	Email	CreatedAt	UpdatedAt
1234567890123456789	John Doe	john@example.com	2025-01-27T10:30:00Z	2025-01-27T10:30:00Z
1234567890123456790	Jane Smith	jane@example.com	2025-01-27T11:00:00Z	2025-01-27T11:00:00Z
```

| フィールド | 型 | 説明 |
|----------|-----|------|
| ID | int64 | ユーザーID（タイムスタンプベース） |
| Name | string | ユーザー名 |
| Email | string | メールアドレス |
| CreatedAt | RFC3339 | 作成日時 |
| UpdatedAt | RFC3339 | 更新日時 |

### 終了コード

| コード | 説明 |
|-------|------|
| 0 | 正常終了 |
| 1 | エラー終了（設定エラー、DB接続エラー、引数エラーなど） |

### エラーメッセージ

エラーメッセージは標準エラー出力に出力されます。

```bash
# limit値が不正な場合
$ APP_ENV=develop ./bin/list-users --limit 0
2025/01/27 10:30:00 Error: limit must be at least 1

# limit値が最大値を超えた場合（警告）
$ APP_ENV=develop ./bin/list-users --limit 200
2025/01/27 10:30:00 Warning: limit exceeds maximum (100), using 100
ID	Name	Email	CreatedAt	UpdatedAt
...
```

## cron設定例

### 毎日午前3時にユーザー一覧をバックアップ

```cron
0 3 * * * cd /path/to/server && APP_ENV=production ./bin/list-users --limit 100 > /var/log/users_$(date +\%Y\%m\%d).tsv 2>> /var/log/list-users.log
```

### 環境変数の設定

cronで実行する場合、環境変数を明示的に設定する必要があります。

```cron
APP_ENV=production
PATH=/usr/local/go/bin:/usr/bin:/bin

0 3 * * * cd /path/to/server && ./bin/list-users --limit 100 > /var/log/users.tsv 2>&1
```

## テスト

### ユニットテストの実行

```bash
cd server
go test -v ./cmd/list-users/...
```

### テストカバレッジの確認

```bash
cd server
go test -cover ./cmd/list-users/...
```

## アーキテクチャ

CLIツールは既存のレイヤードアーキテクチャを再利用しています。APIサーバーと同様に、usecase層を介してservice層を呼び出す構成になっています。

### list-dm-users コマンド

```
┌─────────────────────────────────────────────────────────────┐
│                    list-dm-users コマンド                     │
│                    (cmd/list-dm-users/main.go)              │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│         Usecase層 (internal/usecase/cli)                     │
│         - ListDmUsersUsecase.ListDmUsers()                  │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Service層 (internal/service)                    │
│              - DmUserService.ListDmUsers()                  │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Repository層 (internal/repository)              │
│              - DmUserRepository.List()                      │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              DB層 (internal/db)                              │
│              - GroupManager                                  │
└────────────────────────┬──────────────────────────────────┘
                         │
          ┌──────────────┴──────────────┐
          ▼                              ▼
    ┌─────────┐                    ┌─────────┐
    │ Shard 1 │                    │ Shard 2 │
    └─────────┘                    └─────────┘
```

### generate-sample-data コマンド

```
┌─────────────────────────────────────────────────────────────┐
│               generate-sample-data コマンド                   │
│               (cmd/generate-sample-data/main.go)            │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│         Usecase層 (internal/usecase/cli)                     │
│         - GenerateSampleUsecase.GenerateSampleData()        │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Service層 (internal/service)                    │
│              - GenerateSampleService.GenerateDmUsers()      │
│              - GenerateSampleService.GenerateDmPosts()      │
│              - GenerateSampleService.GenerateDmNews()       │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Repository層 (internal/repository)              │
│              - DmUserRepository.InsertDmUsersBatch()        │
│              - DmPostRepository.InsertDmPostsBatch()        │
│              - DmNewsRepository.InsertDmNewsBatch()         │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              DB層 (internal/db)                              │
│              - GroupManager                                  │
│              - TableSelector                                 │
└────────────────────────┬──────────────────────────────────┘
                         │
          ┌──────────────┼──────────────┐
          ▼              ▼              ▼
    ┌─────────┐    ┌─────────┐    ┌─────────┐
    │ Master  │    │ Shard 1 │    │ Shard 2 │  ...
    │(dm_news)│    │(dm_users│    │(dm_users│
    └─────────┘    │ dm_posts)│   │ dm_posts)│
                   └─────────┘    └─────────┘
```

### レイヤー構造

| レイヤー | ディレクトリ | 役割 |
|---------|-------------|------|
| CLI層 | cmd/list-dm-users/main.go | エントリーポイント、バリデーション、入出力制御 |
| Usecase層 | internal/usecase/cli/ | CLI用ビジネスロジック調整 |
| Service層 | internal/service/ | ドメインロジック、クロスシャード操作 |
| Repository層 | internal/repository/ | データアクセス抽象化 |
| DB層 | internal/db/ | シャーディング戦略、接続管理 |

## 関連ドキュメント

- [Architecture.md](Architecture.md) - アーキテクチャ詳細
- [Sharding.md](Sharding.md) - シャーディング戦略
- [Testing.md](Testing.md) - テスト戦略
