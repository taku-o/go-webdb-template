# Gofakeitによる開発用サンプルデータ生成機能実装タスク一覧

## 概要
Gofakeitによる開発用サンプルデータ生成機能の実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: 依存関係の追加

#### - [ ] タスク 1.1: Gofakeitライブラリの追加
**目的**: サンプルデータ生成に必要なGofakeitライブラリを追加

**作業内容**:
- `go.mod`に`github.com/brianvoe/gofakeit/v6`を追加
- `go get github.com/brianvoe/gofakeit/v6`を実行
- `go mod tidy`を実行して依存関係を整理

**受け入れ基準**:
- `go.mod`に`github.com/brianvoe/gofakeit/v6`が追加されている
- `go mod tidy`が正常に実行される
- 依存関係のエラーがない

---

### Phase 2: CLIツールのディレクトリ構造作成

#### - [ ] タスク 2.1: CLIツールディレクトリの作成
**目的**: サンプルデータ生成CLIツールを配置するディレクトリを作成

**作業内容**:
- `server/cmd/generate-sample-data/`ディレクトリを作成

**受け入れ基準**:
- `server/cmd/generate-sample-data/`ディレクトリが存在する

---

### Phase 3: CLIツールの基本実装

#### - [ ] タスク 3.1: main.goの基本構造の作成
**目的**: CLIツールの基本構造を作成

**作業内容**:
- `server/cmd/generate-sample-data/main.go`を作成
- パッケージ宣言とインポート文を追加:
  - `package main`
  - 必要なパッケージのインポート（`context`, `fmt`, `log`, `os`, `time`, `strings`）
  - プロジェクト内部パッケージのインポート（`config`, `db`, `model`）
  - `github.com/brianvoe/gofakeit/v6`のインポート
- `main()`関数の空の実装を作成

**受け入れ基準**:
- `server/cmd/generate-sample-data/main.go`が存在する
- コンパイルエラーがない

---

#### - [ ] タスク 3.2: 設定読み込みとDB接続の実装
**目的**: 設定ファイルの読み込みとGroupManagerの初期化を実装

**作業内容**:
- `main()`関数に設定読み込み処理を追加:
  - `config.Load()`で設定を読み込み
  - エラーハンドリングを実装
- GroupManagerの初期化処理を追加:
  - `db.NewGroupManager(cfg)`でGroupManagerを作成
  - `groupManager.PingAll()`でデータベース接続確認
  - `defer groupManager.CloseAll()`でリソースのクリーンアップ

**受け入れ基準**:
- 設定ファイルが正常に読み込まれる
- GroupManagerが正常に初期化される
- データベース接続確認が正常に実行される
- エラー発生時に適切なエラーメッセージが表示される

---

### Phase 4: バッチ挿入関数の実装

#### - [ ] タスク 4.1: 汎用バッチ挿入関数の実装
**目的**: 動的テーブル名に対応したバッチ挿入関数を実装

**作業内容**:
- `insertUsersBatch`関数を実装:
  - 引数: `db *gorm.DB`, `tableName string`, `users []*model.User`
  - 生SQLでバッチ挿入を実装（動的テーブル名対応）
  - バッチサイズ500件ずつに分割
  - トランザクション管理を実装
- `insertPostsBatch`関数を実装:
  - 引数: `db *gorm.DB`, `tableName string`, `posts []*model.Post`
  - 生SQLでバッチ挿入を実装（動的テーブル名対応）
  - バッチサイズ500件ずつに分割
  - トランザクション管理を実装
- `insertNewsBatch`関数を実装:
  - 引数: `db *gorm.DB`, `news []*model.News`
  - GORMの`CreateInBatches`を使用（固定テーブル名）
  - バッチサイズ500件ずつに分割
  - トランザクション管理を実装

**受け入れ基準**:
- `insertUsersBatch`関数が実装されている
- `insertPostsBatch`関数が実装されている
- `insertNewsBatch`関数が実装されている
- バッチサイズ500件ずつに分割されている
- トランザクション管理が実装されている
- エラーハンドリングが適切に実装されている

---

### Phase 5: usersテーブルへのデータ生成実装

#### - [ ] タスク 5.1: generateUsers関数の実装
**目的**: usersテーブルにデータを生成する関数を実装

**作業内容**:
- `generateUsers`関数を実装:
  - 引数: `groupManager *db.GroupManager`, `totalCount int`
  - 戻り値: `[]int64`（生成されたuser_idのリスト）
  - 32分割テーブル（users_000～users_031）に均等に分散
  - 各テーブルに約3～4件ずつ生成（100件 ÷ 32テーブル）
  - Gofakeitを使用してデータ生成:
    - `name`: `gofakeit.Name()`
    - `email`: `gofakeit.Email()`
    - `created_at`, `updated_at`: `time.Now()`
  - 各テーブル番号（0～31）をループしてデータ生成
  - `groupManager.GetShardingConnection(tableNumber)`で接続を取得
  - `insertUsersBatch`関数を使用してバッチ挿入
  - 生成されたuser_idを収集してリスト化
  - 進捗状況をログ出力

**受け入れ基準**:
- `generateUsers`関数が実装されている
- 32分割テーブルに均等に分散されている
- 各テーブルに約3～4件ずつ生成されている
- Gofakeitを使用してデータが生成されている
- 生成されたuser_idがリスト化されている
- 進捗状況がログ出力されている
- エラーハンドリングが適切に実装されている

---

### Phase 6: postsテーブルへのデータ生成実装

#### - [ ] タスク 6.1: generatePosts関数の実装
**目的**: postsテーブルにデータを生成する関数を実装

**作業内容**:
- `generatePosts`関数を実装:
  - 引数: `groupManager *db.GroupManager`, `userIDs []int64`, `totalCount int`
  - 32分割テーブル（posts_000～posts_031）に均等に分散
  - 各テーブルに約3～4件ずつ生成（100件 ÷ 32テーブル）
  - 既存のusersテーブルからuser_idを参照:
    - `userIDs`リストからランダムに選択
    - `gofakeit.IntRange(0, len(userIDs)-1)`でインデックスを取得
  - Gofakeitを使用してデータ生成:
    - `user_id`: 既存のusersテーブルからランダムに選択
    - `title`: `gofakeit.Sentence(5)`
    - `content`: `gofakeit.Paragraph(3, 5, 10, "\n")`
    - `created_at`, `updated_at`: `time.Now()`
  - 各テーブル番号（0～31）をループしてデータ生成
  - `groupManager.GetShardingConnection(tableNumber)`で接続を取得
  - `insertPostsBatch`関数を使用してバッチ挿入
  - 進捗状況をログ出力

**受け入れ基準**:
- `generatePosts`関数が実装されている
- 32分割テーブルに均等に分散されている
- 各テーブルに約3～4件ずつ生成されている
- user_idが既存のusersテーブルから参照されている
- Gofakeitを使用してデータが生成されている
- 進捗状況がログ出力されている
- エラーハンドリングが適切に実装されている

---

### Phase 7: newsテーブルへのデータ生成実装

#### - [ ] タスク 7.1: generateNews関数の実装
**目的**: newsテーブルにデータを生成する関数を実装

**作業内容**:
- `generateNews`関数を実装:
  - 引数: `groupManager *db.GroupManager`, `totalCount int`
  - master DBのnewsテーブルに直接100件を生成
  - Gofakeitを使用してデータ生成:
    - `title`: `gofakeit.Sentence(5)`
    - `content`: `gofakeit.Paragraph(3, 5, 10, "\n")`
    - `author_id`: `gofakeit.Int64()`（ランダムな整数、NULL許容）
    - `published_at`: `gofakeit.Date()`（ランダムな日時、NULL許容）
    - `created_at`, `updated_at`: `time.Now()`
  - `groupManager.GetMasterConnection()`でmaster接続を取得
  - `insertNewsBatch`関数を使用してバッチ挿入
  - 進捗状況をログ出力

**受け入れ基準**:
- `generateNews`関数が実装されている
- master DBのnewsテーブルに100件が生成されている
- Gofakeitを使用してデータが生成されている
- 進捗状況がログ出力されている
- エラーハンドリングが適切に実装されている

---

### Phase 8: main関数の統合

#### - [ ] タスク 8.1: main関数の統合実装
**目的**: すべてのデータ生成関数をmain関数に統合

**作業内容**:
- `main()`関数にデータ生成処理を統合:
  1. 設定読み込みとDB接続（タスク3.2で実装済み）
  2. `generateUsers`関数を呼び出し（userIDsを取得）
  3. `generatePosts`関数を呼び出し（userIDsを渡す）
  4. `generateNews`関数を呼び出し
  5. 生成完了メッセージを表示
  6. 正常終了（終了コード: 0）
- エラーハンドリングを実装:
  - 各データ生成関数のエラーを適切に処理
  - エラー発生時は非ゼロの終了コードを返す
  - 詳細なエラーメッセージを表示

**受け入れ基準**:
- `main()`関数にすべてのデータ生成処理が統合されている
- データ生成の順序が正しい（users → posts → news）
- エラーハンドリングが適切に実装されている
- 生成完了メッセージが表示される
- 正常終了時に終了コード0を返す
- エラー発生時に非ゼロの終了コードを返す

---

### Phase 9: ビルドと動作確認

#### - [ ] タスク 9.1: ビルドの確認
**目的**: CLIツールが正常にビルドできることを確認

**作業内容**:
- `cd server && go build -o bin/generate-sample-data ./cmd/generate-sample-data`を実行
- ビルドエラーがないことを確認
- 実行ファイルが生成されることを確認

**受け入れ基準**:
- ビルドが正常に完了する
- ビルドエラーがない
- `server/bin/generate-sample-data`が生成される

---

#### - [ ] タスク 9.2: 動作確認
**目的**: CLIツールが正常に動作することを確認

**作業内容**:
- `APP_ENV=develop ./bin/generate-sample-data`を実行
- データ生成が正常に実行されることを確認
- 各テーブルにデータが生成されることを確認
- 進捗状況が表示されることを確認
- 生成完了メッセージが表示されることを確認

**受け入れ基準**:
- コマンドが正常に実行される
- usersテーブルに100件程度のデータが生成される
- postsテーブルに100件程度のデータが生成される
- newsテーブルに100件程度のデータが生成される
- シャーディングテーブル（users, posts）が適切に分散される
- postsテーブルのuser_idが既存のusersテーブルから参照されている
- 進捗状況が表示される
- 生成完了メッセージが表示される
- エラーが発生しない

---

### Phase 10: ドキュメント作成

#### - [ ] タスク 10.1: Generate-Sample-Data.mdの作成
**目的**: サンプルデータ生成機能のドキュメントを作成

**作業内容**:
- `docs/Generate-Sample-Data.md`を作成
- ドキュメントに以下の内容を記載:
  - 機能概要
  - ビルド方法
  - 実行方法
  - 生成されるデータの説明
  - オプション（現時点ではなし）
  - トラブルシューティング

**受け入れ基準**:
- `docs/Generate-Sample-Data.md`が作成されている
- コマンドの使用方法が記載されている
- 生成されるデータの説明が記載されている
- ビルド方法が記載されている
- 実行方法が記載されている

---

## 実装順序の推奨

1. Phase 1: 依存関係の追加
2. Phase 2: CLIツールのディレクトリ構造作成
3. Phase 3: CLIツールの基本実装
4. Phase 4: バッチ挿入関数の実装
5. Phase 5: usersテーブルへのデータ生成実装
6. Phase 6: postsテーブルへのデータ生成実装
7. Phase 7: newsテーブルへのデータ生成実装
8. Phase 8: main関数の統合
9. Phase 9: ビルドと動作確認
10. Phase 10: ドキュメント作成

## 注意事項

- 各フェーズは順番に実装すること
- 各タスクの受け入れ基準を満たしてから次のタスクに進むこと
- エラーハンドリングを適切に実装すること
- 進捗状況を適切に表示すること
- 既存のCLI実装パターン（`server/cmd/list-users/`）に従うこと
