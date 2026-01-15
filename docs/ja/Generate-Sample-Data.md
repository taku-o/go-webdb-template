**[日本語]** | [English](../en/Generate-Sample-Data.md)

# サンプルデータ生成機能

## 概要

開発用のサンプルデータを大量に生成するCLIツールです。Gofakeitライブラリを使用して、リアルなランダムデータを生成します。

PostgreSQLを使用してmaster/shardingデータベースにデータを投入します。

## 前提条件

コマンド実行前に以下の準備が必要です。

### 1. PostgreSQLコンテナの起動

```bash
./scripts/start-postgres.sh start
```

**接続情報**（開発環境）:

| データベース | ホスト | ポート | ユーザー | パスワード | データベース名 |
|------------|--------|--------|---------|-----------|--------------|
| Master | localhost | 5432 | webdb | webdb | webdb_master |
| Sharding 1 | localhost | 5433 | webdb | webdb | webdb_sharding_1 |
| Sharding 2 | localhost | 5434 | webdb | webdb | webdb_sharding_2 |
| Sharding 3 | localhost | 5435 | webdb | webdb | webdb_sharding_3 |
| Sharding 4 | localhost | 5436 | webdb | webdb | webdb_sharding_4 |

### 2. マイグレーションの適用

```bash
./scripts/migrate.sh
```

## ビルド方法

```bash
cd server
go build -o bin/generate-sample-data ./cmd/generate-sample-data
```

## 実行方法

```bash
# go runで実行
cd server
APP_ENV=develop go run cmd/generate-sample-data/main.go

# ビルド済みバイナリで実行
APP_ENV=develop ./bin/generate-sample-data
```

## 生成されるデータ

### dm_usersテーブル

- **対象**: shardingデータベース（4台に分散）
- **テーブル**: 32分割テーブル（dm_users_000〜dm_users_031）
- **生成件数**: 合計100件（UUIDに基づいて各テーブルに分散）
- **生成フィールド**:
  - `id`: UUIDv7（ハイフン抜き小文字、32文字）
  - `name`: ランダムな名前
  - `email`: ランダムなメールアドレス
  - `created_at`, `updated_at`: 現在時刻

### dm_postsテーブル

- **対象**: shardingデータベース（4台に分散）
- **テーブル**: 32分割テーブル（dm_posts_000〜dm_posts_031）
- **生成件数**: 合計100件（user_idに基づいて各テーブルに分散）
- **生成フィールド**:
  - `id`: UUIDv7（ハイフン抜き小文字、32文字）
  - `user_id`: 既存のdm_usersテーブルからランダムに選択
  - `title`: 5単語程度のランダムな文
  - `content`: 3〜5文、各文10単語程度のランダムな段落
  - `created_at`, `updated_at`: 現在時刻

### dm_newsテーブル

- **対象**: masterデータベース（webdb_master）
- **テーブル**: dm_news（固定テーブル名）
- **生成件数**: 100件
- **生成フィールド**:
  - `title`: 5単語程度のランダムな文
  - `content`: 3〜5文、各文10単語程度のランダムな段落
  - `author_id`: ランダムな32ビット整数
  - `published_at`: ランダムな日時
  - `created_at`, `updated_at`: 現在時刻

## 実行例

```
$ APP_ENV=develop go run cmd/generate-sample-data/main.go
2026/01/09 23:39:14 Starting sample data generation...
2026/01/09 23:39:14 Generated 2 dm_users in dm_users_001
2026/01/09 23:39:14 Generated 3 dm_users in dm_users_023
...
2026/01/09 23:39:14 Generated 5 dm_posts in dm_posts_010
2026/01/09 23:39:14 Generated 1 dm_posts in dm_posts_022
...
2026/01/09 23:39:14 Generated 100 dm_news articles
2026/01/09 23:39:14 Sample data generation completed successfully
```

## データの確認

### CloudBeaverで確認

```bash
./scripts/cloudbeaver-start.sh
```

ブラウザで http://localhost:8978 にアクセスして、生成されたデータを確認できます。

### コマンドラインで確認

```bash
# dm_newsの件数確認
docker exec postgres-master psql -U webdb -d webdb_master -c "SELECT COUNT(*) FROM dm_news;"

# dm_usersの件数確認（sharding-1のテーブル）
docker exec postgres-sharding-1 psql -U webdb -d webdb_sharding_1 -c "SELECT COUNT(*) FROM dm_users_000;"
```

## トラブルシューティング

### 接続エラー

```
Failed to create group manager: ...
```

**原因**: PostgreSQLコンテナが起動していない

**解決方法**:
```bash
./scripts/start-postgres.sh start
```

### テーブルが存在しないエラー

```
relation "dm_users_000" does not exist
```

**原因**: マイグレーションが適用されていない

**解決方法**:
```bash
./scripts/migrate.sh
```

## 注意事項

- develop環境での使用を想定しています
- 既存データの削除は行いません（追加のみ）
- 複数回実行するとデータが追加されます

## アーキテクチャ

generate-sample-dataコマンドは、APIサーバーと同じレイヤードアーキテクチャを使用しています。usecase層を介してservice層を呼び出す構成になっています。

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
| CLI層 | cmd/generate-sample-data/main.go | エントリーポイント、入出力制御 |
| Usecase層 | internal/usecase/cli/ | CLI用ビジネスロジック調整 |
| Service層 | internal/service/ | データ生成ロジック、gofakeitの使用 |
| Repository層 | internal/repository/ | バッチ挿入、データアクセス抽象化 |
| DB層 | internal/db/ | シャーディング戦略、接続管理 |

## 技術仕様

- **バッチサイズ**: 500件ずつ
- **シャーディング**: UUIDに基づいて32分割テーブルに分散
- **ライブラリ**: `github.com/brianvoe/gofakeit/v6`
- **データベース**: PostgreSQL（master 1台 + sharding 4台）

## 関連ドキュメント

- [Command-Line-Tool.md](./Command-Line-Tool.md) - 既存のCLIツールドキュメント
- [Sharding.md](./Sharding.md) - シャーディング戦略の詳細
- [Database-Viewer.md](./Database-Viewer.md) - CloudBeaverの使い方
