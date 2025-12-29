# サンプルデータ生成機能

## 概要

開発用のサンプルデータを大量に生成するCLIツールです。Gofakeitライブラリを使用して、リアルなランダムデータを生成します。

## ビルド方法

```bash
cd server
go build -o bin/generate-sample-data ./cmd/generate-sample-data
```

## 実行方法

```bash
APP_ENV=develop ./bin/generate-sample-data
```

## 生成されるデータ

### usersテーブル

- **対象**: 32分割テーブル（users_000～users_031）
- **生成件数**: 合計約96件（各テーブルに3件ずつ）
- **生成フィールド**:
  - `name`: ランダムな名前
  - `email`: ランダムなメールアドレス
  - `created_at`, `updated_at`: 現在時刻

### postsテーブル

- **対象**: 32分割テーブル（posts_000～posts_031）
- **生成件数**: 合計約96件（各テーブルに3件ずつ）
- **生成フィールド**:
  - `user_id`: 既存のusersテーブルからランダムに選択
  - `title`: 5単語程度のランダムな文
  - `content`: 3～5文、各文10単語程度のランダムな段落
  - `created_at`, `updated_at`: 現在時刻

### newsテーブル

- **対象**: master DBのnewsテーブル
- **生成件数**: 100件
- **生成フィールド**:
  - `title`: 5単語程度のランダムな文
  - `content`: 3～5文、各文10単語程度のランダムな段落
  - `author_id`: ランダムな整数
  - `published_at`: ランダムな日時
  - `created_at`, `updated_at`: 現在時刻

## 実行例

```
$ APP_ENV=develop ./bin/generate-sample-data
2025/12/29 17:48:57 Starting sample data generation...
2025/12/29 17:48:57 Generated 3 users in users_000
2025/12/29 17:48:57 Generated 3 users in users_001
...
2025/12/29 17:48:57 Generated 3 users in users_031
2025/12/29 17:48:57 Generated 3 posts in posts_000
2025/12/29 17:48:57 Generated 3 posts in posts_001
...
2025/12/29 17:48:57 Generated 3 posts in posts_031
2025/12/29 17:48:57 Generated 100 news articles
2025/12/29 17:48:57 Sample data generation completed successfully
```

## 注意事項

- develop環境での使用を想定しています
- 既存データの削除は行いません（追加のみ）
- データ生成量は固定です（変更不可）
- 複数回実行するとデータが追加されます

## 技術仕様

- **バッチサイズ**: 500件ずつ
- **シャーディング**: 32分割テーブルに均等に分散
- **ライブラリ**: `github.com/brianvoe/gofakeit/v6`

## 関連ドキュメント

- [Command-Line-Tool.md](./Command-Line-Tool.md) - 既存のCLIツールドキュメント
- [Sharding.md](./Sharding.md) - シャーディング戦略の詳細
